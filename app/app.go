package app

import (
	"github.com/midaef/fibonacci-service/app/endpoint"
	"github.com/midaef/fibonacci-service/app/endpoint/fibonacci"
	"github.com/midaef/fibonacci-service/app/repository"
	"github.com/midaef/fibonacci-service/app/service"
	"github.com/midaef/fibonacci-service/config"
	"github.com/midaef/fibonacci-service/dependers/redis"
	app_fibonacci "github.com/midaef/fibonacci-service/extra/fibonacci"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"net"
	"net/http"
	"time"
)

// App ...
type App struct {
	logger *zap.Logger
	config *config.Config
}

// NewApp ...
func NewApp(logger *zap.Logger, config *config.Config) *App {
	return &App{
		logger: logger,
		config: config,
	}
}

// StartApp ...
func (app *App) StartApp() error {
	startTime := time.Now().UnixNano()

	rdb := redis.NewRedisClient(app.config)

	app.logger.Info("redis successfully connected",
		zap.String("host", app.config.Redis.Host+":"+app.config.Redis.Port),
	)

	store := repository.NewRepository(rdb)

	service := service.NewService(store, app.config)

	serviceContainer := app.InitEndpointContainer(service)

	go app.StartHTTPServer(serviceContainer)

	listener, err := net.Listen("tcp", ":"+app.config.AppConfig.GRPCPort)
	if err != nil {
		return err
	}

	grpcServer := grpc.NewServer()

	app_fibonacci.RegisterFibonacciServer(grpcServer, serviceContainer.FibonacciService)

	app.logger.Info("fibonacci-service successfully started",
		zap.String("grpc_host", app.config.AppConfig.IP+":"+app.config.AppConfig.GRPCPort),
		zap.String("http_host", app.config.AppConfig.IP+":"+app.config.AppConfig.HTTPPort),
		zap.Int64("duration", time.Now().UnixNano()-startTime),
	)

	if err := grpcServer.Serve(listener); err != nil {
		return err
	}

	return nil
}

func (app *App) StartHTTPServer(serviceContainer *endpoint.EndpointContainer) {
	mux := http.NewServeMux()

	mux.HandleFunc("/api/v1/fibonacci_sequences", serviceContainer.FibonacciService.FibonacciSequencesHTTP)

	err := http.ListenAndServe(":"+app.config.AppConfig.HTTPPort, mux)
	if err != nil {
		app.logger.Fatal("error started http server")
	}
}

func (app *App) InitEndpointContainer(service *service.Service) *endpoint.EndpointContainer {
	fibonacciService := fibonacci.NewFibonacciEndpoint(service, app.config)

	serviceContainer := endpoint.NewEndpointContainer(
		fibonacciService,
	)

	return serviceContainer
}
