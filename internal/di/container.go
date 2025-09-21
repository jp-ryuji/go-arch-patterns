package di

import (
	"github.com/jp-ryuji/go-arch-patterns/internal/application/service"
	"github.com/jp-ryuji/go-arch-patterns/internal/infrastructure/postgres/entgen"
	"github.com/jp-ryuji/go-arch-patterns/internal/infrastructure/postgres/repository"
	"github.com/jp-ryuji/go-arch-patterns/internal/interface/http"
)

// Container holds all the dependencies
type Container struct {
	Client     *entgen.Client
	CarService service.CarService
	HTTPServer *http.Server
	grpcPort   int
	httpPort   int
}

// NewContainer creates a new dependency injection container with an existing client
func NewContainer(client *entgen.Client, grpcPort, httpPort int) (*Container, error) {
	// Create repositories
	carRepo := repository.NewCarRepository(client)
	outboxRepo := repository.NewOutboxRepository(client)

	// Create transaction manager
	txManager := repository.NewTransactionManager(client)

	// Create application services
	carService := service.NewCarService(carRepo, outboxRepo, txManager)

	// Create HTTP server with gRPC-Gateway
	server := http.NewServer(grpcPort, httpPort, carService)

	return &Container{
		Client:     client,
		CarService: carService,
		HTTPServer: server,
		grpcPort:   grpcPort,
		httpPort:   httpPort,
	}, nil
}

// Close closes all resources in the container
func (c *Container) Close() {
	if c.Client != nil {
		c.Client.Close()
	}
}
