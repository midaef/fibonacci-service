package service

import (
	"context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"strconv"
)

func (s *Service) GetFibonacciSequences(ctx context.Context, key string) (string, error) {
	return s.repository.FibonacciRepository.GetFibonacciSequences(ctx, key)
}

func (s *Service) SetFibonacciSequences(ctx context.Context, key string, value string) error {
	return s.repository.FibonacciRepository.SetFibonacciSequences(ctx, key, value)
}

func (s *Service) FibonacciSequences(ctx context.Context, x uint64, y uint64) ([]uint64, error) {
	fib := make([]uint64, 0, (y-x)+1)

	for i := x; i <= y; i++ {
		value, err := s.GetFibonacciSequences(ctx, strconv.Itoa(int(i)))
		if err != nil {
			return nil, err
		}

		if value == "" {
			n, err := s.Fibonacci(ctx, i)
			if err != nil {
				return nil, err
			}

			err = s.SetFibonacciSequences(ctx, strconv.Itoa(int(i)), strconv.Itoa(int(n)))
			if err != nil {
				return nil, err
			}

			value = strconv.Itoa(int(n))
		}

		val, err := strconv.Atoi(value)
		if err != nil {
			return nil, err
		}

		fib = append(fib, uint64(val))
	}

	return fib, nil
}

func (s *Service) Fibonacci(ctx context.Context, n uint64) (uint64, error) {
	if n == 0 {
		return 0, nil
	}

	if n == 1 {
		return 1, nil
	}

	n1, err := s.repository.FibonacciRepository.GetFibonacciSequences(ctx, strconv.Itoa(int(n)-1))
	if err != nil {
		return 0, err
	}

	if n1 == "" {
		nLocal, err := s.Fibonacci(ctx, n-1)
		if err != nil {
			return 0, err
		}

		n1 = strconv.Itoa(int(nLocal))
	}

	firstFibonacci, err := strconv.Atoi(n1)
	if err != nil {
		return 0, err
	}

	n2, err := s.repository.FibonacciRepository.GetFibonacciSequences(ctx, strconv.Itoa(int(n)-2))
	if err != nil {
		return 0, err
	}

	if n2 == "" {
		nLocal, err := s.Fibonacci(ctx, n-1)
		if err != nil {
			return 0, err
		}

		n2 = strconv.Itoa(int(nLocal))
	}

	secondFibonacci, err := strconv.Atoi(n2)
	if err != nil {
		return 0, err
	}

	return uint64(firstFibonacci) + uint64(secondFibonacci), nil
}

func (s *Service) Validate(x uint64, y uint64) error {
	if x < 0 || y < 0 {
		return status.Error(codes.InvalidArgument, "x and y cannot be less than 0")
	} else if x > y {
		return status.Error(codes.InvalidArgument, "y cannot be less than x")
	} else if x == y {
		return status.Error(codes.InvalidArgument, "x cannot be equal to y")
	} else if y > 93 {
		return status.Error(codes.InvalidArgument, "y cannot be greater than 93")
	}

	return nil
}
