package http

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	carv1 "github.com/jp-ryuji/go-arch-patterns/api/generated/car/v1"
	"github.com/jp-ryuji/go-arch-patterns/internal/application/service"
	"github.com/jp-ryuji/go-arch-patterns/internal/interface/grpc/car/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// Server represents the HTTP server with gRPC-Gateway
type Server struct {
	grpcServer *grpc.Server
	httpServer *http.Server
	grpcPort   int
	httpPort   int
	carService service.CarService
}

// NewServer creates a new HTTP server with gRPC-Gateway
func NewServer(grpcPort, httpPort int, carService service.CarService) *Server {
	return &Server{
		grpcPort:   grpcPort,
		httpPort:   httpPort,
		carService: carService,
	}
}

// Start starts the gRPC and HTTP servers
func (s *Server) Start() error {
	// Create gRPC server
	s.grpcServer = grpc.NewServer()

	// Register gRPC services
	carv1.RegisterCarServiceServer(s.grpcServer, car.NewCarServiceServer(s.carService))

	// Create context for Listen calls
	ctx := context.Background()

	// Start gRPC server in a goroutine
	listenConfig := &net.ListenConfig{}
	grpcListener, err := listenConfig.Listen(ctx, "tcp", fmt.Sprintf(":%d", s.grpcPort))
	if err != nil {
		return fmt.Errorf("failed to listen on port %d for gRPC: %w", s.grpcPort, err)
	}

	go func() {
		if err := s.grpcServer.Serve(grpcListener); err != nil {
			fmt.Printf("gRPC server error: %v\n", err)
		}
	}()

	// Create HTTP server with gRPC-Gateway
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	// Create gRPC-Gateway mux
	gwmux := runtime.NewServeMux()

	// Set up connection to gRPC server
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
	err = carv1.RegisterCarServiceHandlerFromEndpoint(ctx, gwmux, fmt.Sprintf("localhost:%d", s.grpcPort), opts)
	if err != nil {
		return fmt.Errorf("failed to register car service handler: %w", err)
	}

	// Create HTTP server with timeout configuration
	s.httpServer = &http.Server{
		Addr:              fmt.Sprintf(":%d", s.httpPort),
		Handler:           gwmux,
		ReadHeaderTimeout: 5 * time.Second, // Add timeout to prevent Slowloris attacks
	}

	// Start HTTP server
	httpListener, err := listenConfig.Listen(ctx, "tcp", fmt.Sprintf(":%d", s.httpPort))
	if err != nil {
		return fmt.Errorf("failed to listen on port %d for HTTP: %w", s.httpPort, err)
	}

	go func() {
		if err := s.httpServer.Serve(httpListener); err != nil && err != http.ErrServerClosed {
			fmt.Printf("HTTP server error: %v\n", err)
		}
	}()

	fmt.Printf("gRPC server started on port %d\n", s.grpcPort)
	fmt.Printf("HTTP server with gRPC-Gateway started on port %d\n", s.httpPort)

	return nil
}

// Stop stops the gRPC and HTTP servers
func (s *Server) Stop() error {
	// Stop gRPC server
	s.grpcServer.GracefulStop()

	// Stop HTTP server
	if err := s.httpServer.Close(); err != nil {
		return fmt.Errorf("failed to stop HTTP server: %w", err)
	}

	return nil
}
