package repository

import (
	"context"
	"github.com/go-redis/redis/v8"
)

type FibonacciRepository interface {
	GetFibonacciSequences(ctx context.Context, key string) (string, error)
	SetFibonacciSequences(ctx context.Context, key string, value string) error
}

type Repository struct {
	FibonacciRepository FibonacciRepository
}

func NewRepository(rdb *redis.Client) *Repository {
	return &Repository{
		FibonacciRepository: NewFibonacciRepository(rdb),
	}
}
