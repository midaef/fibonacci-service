package service

import (
	"github.com/midaef/fibonacci-service/app/repository"
	"github.com/midaef/fibonacci-service/config"
)

type Service struct {
	repository *repository.Repository
	config     *config.Config
}

func NewService(repository *repository.Repository, config *config.Config) *Service {
	return &Service{
		repository: repository,
		config:     config,
	}
}
