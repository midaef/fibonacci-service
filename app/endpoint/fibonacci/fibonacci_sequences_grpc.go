package fibonacci

import (
	"context"
	"github.com/midaef/fibonacci-service/extra/fibonacci"
)

func (f *FibonacciEndpoint) FibonacciSequences(ctx context.Context, req *fibonacci.FibonacciSequencesRequest) (*fibonacci.FibonacciSequencesResponse, error) {
	fib, err := f.service.FibonacciSequences(ctx, req.GetX(), req.GetY())
	if err != nil {
		return nil, err
	}

	return &fibonacci.FibonacciSequencesResponse{
		FibonacciSequences: fib,
	}, nil
}
