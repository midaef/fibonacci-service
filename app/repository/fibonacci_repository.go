package repository

import (
	"context"
	"github.com/go-redis/redis/v8"
)

type Fibonacci struct {
	rdb *redis.Client
}

func NewFibonacciRepository(rdb *redis.Client) *Fibonacci {
	return &Fibonacci{
		rdb: rdb,
	}
}

func (f *Fibonacci) GetFibonacciSequences(ctx context.Context, key string) (string, error) {
	value, err := f.rdb.Get(ctx, key).Result()
	switch {
	case err == redis.Nil:
		return "", nil
	case err != nil:
		return "", err
	case value == "":
		return "", nil
	}

	return value, nil
}

func (f *Fibonacci) SetFibonacciSequences(ctx context.Context, key string, value string) error {
	err := f.rdb.Set(ctx, key, value, 0).Err()
	if err != nil {
		return err
	}

	return nil
}
