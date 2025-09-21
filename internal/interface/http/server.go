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
	carServiceServer := car.NewCarServiceServer(s.carService)
	carv1.RegisterCarServiceServer(s.grpcServer, carServiceServer)
	fmt.Printf("Registered CarServiceServer with gRPC server\n")

	// Create context for Listen calls
	ctx := context.Background()

	// Start gRPC server in a goroutine
	listenConfig := &net.ListenConfig{}
	grpcAddr := fmt.Sprintf(":%d", s.grpcPort)
	grpcListener, err := listenConfig.Listen(ctx, "tcp", grpcAddr)
	if err != nil {
		return fmt.Errorf("failed to listen on port %d for gRPC: %w", s.grpcPort, err)
	}
	fmt.Printf("gRPC server listening on %s\n", grpcAddr)

	go func() {
		fmt.Printf("Starting gRPC server\n")
		if err := s.grpcServer.Serve(grpcListener); err != nil {
			fmt.Printf("gRPC server error: %v\n", err)
		} else {
			fmt.Printf("gRPC server stopped\n")
		}
	}()
	fmt.Printf("gRPC server goroutine started\n")

	// Give the gRPC server a moment to start
	time.Sleep(100 * time.Millisecond)
	fmt.Printf("Waited for gRPC server to start\n")

	// Create HTTP server with gRPC-Gateway
	// Create a new context for the gateway (don't cancel the previous one)
	gwCtx := context.Background()

	// Create gRPC-Gateway mux
	gwmux := runtime.NewServeMux()
	fmt.Printf("Created gRPC-Gateway mux\n")

	// Set up connection to gRPC server
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
	grpcEndpoint := fmt.Sprintf("localhost:%d", s.grpcPort)
	fmt.Printf("Connecting to gRPC server at %s\n", grpcEndpoint)
	err = carv1.RegisterCarServiceHandlerFromEndpoint(gwCtx, gwmux, grpcEndpoint, opts)
	if err != nil {
		fmt.Printf("Failed to register car service handler: %v\n", err)
		return fmt.Errorf("failed to register car service handler: %w", err)
	}
	fmt.Printf("Registered car service handler with gRPC-Gateway\n")

	// Create HTTP server with timeout configuration
	s.httpServer = &http.Server{
		Addr:              fmt.Sprintf(":%d", s.httpPort),
		Handler:           gwmux,
		ReadHeaderTimeout: 5 * time.Second, // Add timeout to prevent Slowloris attacks
	}
	fmt.Printf("Created HTTP server on port %d\n", s.httpPort)

	// Start HTTP server
	httpAddr := fmt.Sprintf(":%d", s.httpPort)
	httpListener, err := listenConfig.Listen(ctx, "tcp", httpAddr)
	if err != nil {
		return fmt.Errorf("failed to listen on port %d for HTTP: %w", s.httpPort, err)
	}
	fmt.Printf("HTTP server listening on %s\n", httpAddr)

	go func() {
		fmt.Printf("Starting HTTP server\n")
		if err := s.httpServer.Serve(httpListener); err != nil && err != http.ErrServerClosed {
			fmt.Printf("HTTP server error: %v\n", err)
		} else {
			fmt.Printf("HTTP server stopped\n")
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
