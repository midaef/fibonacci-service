package fibonacci

import (
	"context"
	"github.com/midaef/fibonacci-service/config"
)

type IFibonacciService interface {
	Fibonacci(ctx context.Context, n uint64) (uint64, error)

	GetFibonacciSequences(ctx context.Context, key string) (string, error)
	SetFibonacciSequences(ctx context.Context, key string, value string) error
	FibonacciSequences(ctx context.Context, x uint64, y uint64) ([]uint64, error)
	Validate(x uint64, y uint64) error
}

type FibonacciEndpoint struct {
	service IFibonacciService
}

func NewFibonacciEndpoint(service IFibonacciService, config *config.Config) *FibonacciEndpoint {
	return &FibonacciEndpoint{
		service: service,
	}
}
