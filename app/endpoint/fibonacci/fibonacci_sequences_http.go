package fibonacci

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

type request struct {
	X uint64 `json:"x"`
	Y uint64 `json:"y"`
}

func (f *FibonacciEndpoint) FibonacciSequencesHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		var req *request

		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			f.responseWriter(http.StatusInternalServerError, map[string]interface{}{
				"error": err.Error(),
			}, w)

			return
		}

		err = json.Unmarshal(body, &req)
		if err != nil {
			f.responseWriter(http.StatusBadRequest, map[string]interface{}{
				"error": err.Error(),
			}, w)

			return
		}

		err = f.service.Validate(req.X, req.Y)
		if err != nil {
			f.responseWriter(http.StatusBadRequest, map[string]interface{}{
				"error": err.Error(),
			}, w)

			return
		}

		fib, err := f.service.FibonacciSequences(r.Context(), req.X, req.Y)
		if err != nil {
			f.responseWriter(http.StatusInternalServerError, map[string]interface{}{
				"error": err.Error(),
			}, w)

			return
		}

		f.responseWriter(http.StatusOK, map[string]interface{}{
			"fibonacci_sequences": fib,
		}, w)

		return
	} else {
		f.responseWriter(http.StatusMethodNotAllowed, map[string]interface{}{
			"error": "method not allowed",
		}, w)

		return
	}
}

func (f *FibonacciEndpoint) responseWriter(statusCode int, data interface{}, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")

	w.WriteHeader(statusCode)

	json, err := json.Marshal(data)
	if err != nil {
		f.responseWriter(http.StatusInternalServerError, map[string]interface{}{
			"error": err.Error(),
		}, w)

		return
	}

	_, err = w.Write(json)
	if err != nil {
		f.responseWriter(http.StatusInternalServerError, map[string]interface{}{
			"error": err.Error(),
		}, w)

		return
	}
}
