package tests

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/midaef/fibonacci-service/app"
	"github.com/midaef/fibonacci-service/app/endpoint"
	"github.com/midaef/fibonacci-service/app/repository"
	"github.com/midaef/fibonacci-service/app/service"
	"github.com/midaef/fibonacci-service/config"
	"github.com/midaef/fibonacci-service/dependers/redis"
	app_fibonacci "github.com/midaef/fibonacci-service/extra/fibonacci"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log"
	"net"
	"net/http"
	"net/http/httptest"
	"testing"
)

func init() {
	go initTestGRPCServer()
}

const httpPort = "8080"

const grpcPort = "7001"

const ip = "0.0.0.0"

func TestExpectedErrorsGRPCFibonacci(t *testing.T) {
	client, err := initTestGRPCClient()
	if err != nil {
		t.Fatal(err)
	}

	type request struct {
		X uint64
		Y uint64
	}

	testCasesInvalid := []struct {
		name          string
		payload       *request
		expectedError error
	}{
		{
			name: "invalid",
			payload: &request{
				X: 5,
				Y: 5,
			},
			expectedError: status.Error(codes.InvalidArgument, "x cannot be equal to y"),
		},
		{
			name: "invalid",
			payload: &request{
				X: 10,
				Y: 1,
			},
			expectedError: status.Error(codes.InvalidArgument, "y cannot be less than x"),
		},
		{
			name: "invalid",
			payload: &request{
				X: 1,
				Y: 94,
			},
			expectedError: status.Error(codes.InvalidArgument, "y cannot be greater than 93"),
		},
	}

	for _, test := range testCasesInvalid {
		t.Run(test.name, func(t *testing.T) {
			_, err := client.FibonacciSequences(context.Background(), &app_fibonacci.FibonacciSequencesRequest{
				X: test.payload.X,
				Y: test.payload.Y,
			})
			if err != nil {
				assert.EqualError(t, err, test.expectedError.Error())
			}
		})
	}
}

func TestExpectedValidDataGRPCFibonacci(t *testing.T) {
	client, err := initTestGRPCClient()
	if err != nil {
		t.Fatal(err)
	}

	type request struct {
		X uint64
		Y uint64
	}

	testCasesValid := []struct {
		name         string
		payload      *request
		expectedData []uint64
	}{
		{
			name: "valid",
			payload: &request{
				X: 5,
				Y: 15,
			},
			expectedData: []uint64{
				5,
				8,
				13,
				21,
				34,
				55,
				89,
				144,
				233,
				377,
				610,
			},
		},
		{
			name: "valid",
			payload: &request{
				X: 1,
				Y: 5,
			},
			expectedData: []uint64{
				1,
				1,
				2,
				3,
				5,
			},
		},
		{
			name: "valid",
			payload: &request{
				X: 0,
				Y: 3,
			},
			expectedData: []uint64{
				0,
				1,
				1,
				2,
			},
		},
		{
			name: "valid",
			payload: &request{
				X: 90,
				Y: 93,
			},
			expectedData: []uint64{
				2880067194370816120,
				4660046610375530309,
				7540113804746346429,
				12200160415121876738,
			},
		},
	}

	for _, test := range testCasesValid {
		t.Run(test.name, func(t *testing.T) {
			resp, err := client.FibonacciSequences(context.Background(), &app_fibonacci.FibonacciSequencesRequest{
				X: test.payload.X,
				Y: test.payload.Y,
			})
			if err != nil {
				t.Fatal(err)
			}

			assert.Equal(t, resp.GetFibonacciSequences(), test.expectedData)
		})
	}
}

func TestExpectedCodesHTTPFibonacci(t *testing.T) {
	serviceContainer, err := initTestHTTPServer()
	if err != nil {
		t.Fatal(err)
	}

	testCases := []struct {
		name         string
		payload      interface{}
		expectedCode int
	}{
		{
			name: "valid",
			payload: map[string]uint64{
				"x": 5,
				"y": 15,
			},
			expectedCode: http.StatusOK,
		},
		{
			name: "valid",
			payload: map[string]uint64{
				"x": 1,
				"y": 93,
			},
			expectedCode: http.StatusOK,
		},
		{
			name: "valid",
			payload: map[string]uint64{
				"x": 0,
				"y": 93,
			},
			expectedCode: http.StatusOK,
		},
		{
			name: "valid",
			payload: map[string]uint64{
				"x": 50,
				"y": 60,
			},
			expectedCode: http.StatusOK,
		},
		{
			name: "invalid",
			payload: map[string]uint64{
				"x": 50,
			},
			expectedCode: http.StatusBadRequest,
		},
		{
			name: "invalid",
			payload: map[string]uint64{
				"x": 5,
				"y": 5,
			},
			expectedCode: http.StatusBadRequest,
		},
		{
			name: "invalid",
			payload: map[string]uint64{
				"x": 6,
				"y": 2,
			},
			expectedCode: http.StatusBadRequest,
		},
		{
			name: "invalid",
			payload: map[string]int{
				"x": -1,
				"y": 1,
			},
			expectedCode: http.StatusBadRequest,
		},
		{
			name: "invalid",
			payload: map[string]int{
				"x": -1,
				"y": -1,
			},
			expectedCode: http.StatusBadRequest,
		},
		{
			name: "invalid",
			payload: map[string]int{
				"x": 1,
				"y": -1,
			},
			expectedCode: http.StatusBadRequest,
		},
		{
			name: "invalid",
			payload: map[string]int{
				"x": 1,
				"y": 94,
			},
			expectedCode: http.StatusBadRequest,
		},
		{
			name: "invalid",
			payload: map[string]string{
				"x": "1",
				"y": "15",
			},
			expectedCode: http.StatusBadRequest,
		},
	}

	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			rec := httptest.NewRecorder()

			jsonTest, err := json.Marshal(test.payload)
			if err != nil {
				t.Fatal(err)
			}

			req, err := http.NewRequest("POST", "/api/v1/fibonacci_sequences", bytes.NewBuffer(jsonTest))
			if err != nil {
				t.Fatal(err)
			}

			req.Header.Set("Content-Type", "application/json")

			handler := http.HandlerFunc(serviceContainer.FibonacciService.FibonacciSequencesHTTP)
			handler.ServeHTTP(rec, req)

			assert.Equal(t, test.expectedCode, rec.Code)
		})
	}
}

func initTestHTTPServer() (*endpoint.EndpointContainer, error) {
	_, serviceContainer, app := initTestServer()

	go app.StartHTTPServer(serviceContainer)

	return serviceContainer, nil
}

func initTestServer() (*config.Config, *endpoint.EndpointContainer, *app.App) {
	cfg := initTestConfig()

	logger, _ := zap.NewProduction()
	defer logger.Sync()

	app := app.NewApp(logger, cfg)

	rdb := redis.NewRedisClient(cfg)

	store := repository.NewRepository(rdb)

	service := service.NewService(store, cfg)

	serviceContainer := app.InitEndpointContainer(service)

	return cfg, serviceContainer, app
}

func initTestGRPCServer() {
	config, serviceContainer, _ := initTestServer()

	listener, _ := net.Listen("tcp", ":"+config.AppConfig.GRPCPort)

	grpcServer := grpc.NewServer()

	app_fibonacci.RegisterFibonacciServer(grpcServer, serviceContainer.FibonacciService)

	err := grpcServer.Serve(listener)
	if err != nil {
		log.Fatal(err)
	}
}

func initTestConfig() *config.Config {
	return &config.Config{
		AppConfig: config.AppConfig{
			IP:       ip,
			GRPCPort: grpcPort,
			HTTPPort: httpPort,
		},
		Redis: config.Redis{
			DB:       0,
			Host:     ip,
			Port:     "6379",
			Password: "",
		},
	}
}

func initTestGRPCClient() (app_fibonacci.FibonacciClient, error) {
	conn, err := grpc.Dial(ip+":"+grpcPort, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}

	client := app_fibonacci.NewFibonacciClient(conn)

	return client, nil
}
