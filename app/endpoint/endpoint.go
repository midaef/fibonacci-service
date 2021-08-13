package endpoint

import (
	"context"
	"github.com/midaef/fibonacci-service/extra/fibonacci"
	"net/http"
)

type EndpointContainer struct {
	FibonacciService FibonacciServiceInter
}

func NewEndpointContainer(fibonacci FibonacciServiceInter) *EndpointContainer {
	return &EndpointContainer{
		FibonacciService: fibonacci,
	}
}

type FibonacciServiceInter interface {
	FibonacciSequences(ctx context.Context, req *fibonacci.FibonacciSequencesRequest) (*fibonacci.FibonacciSequencesResponse, error)
	FibonacciSequencesHTTP(w http.ResponseWriter, r *http.Request)
}
