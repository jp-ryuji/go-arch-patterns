package http

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"time"

	"connectrpc.com/grpchealth"
	"connectrpc.com/grpcreflect"
	"github.com/jp-ryuji/go-arch-patterns/api/generated/car/v1/carv1connect"
	"github.com/jp-ryuji/go-arch-patterns/internal/application/service"
	connectcar "github.com/jp-ryuji/go-arch-patterns/internal/presentation/connect/car/v1"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
)

// Server represents the HTTP server with gRPC Connect
type Server struct {
	httpServer *http.Server
	grpcPort   int
	httpPort   int
	carService service.CarService
}

// NewServer creates a new HTTP server with gRPC Connect
func NewServer(grpcPort, httpPort int, carService service.CarService) *Server {
	return &Server{
		grpcPort:   grpcPort,
		httpPort:   httpPort,
		carService: carService,
	}
}

// Start starts the HTTP server with gRPC Connect
func (s *Server) Start() error {
	// Create context for Listen calls
	ctx := context.Background()

	// Create HTTP server with gRPC Connect
	// Create a mux for Connect handlers
	mux := http.NewServeMux()

	// Register Connect handlers
	connectCarServiceHandler := connectcar.NewCarServiceHandler(s.carService)
	path, handler := carv1connect.NewCarServiceHandler(connectCarServiceHandler)
	mux.Handle(path, handler)

	// Register health and reflection handlers
	mux.Handle(grpchealth.NewHandler(grpchealth.NewStaticChecker(carv1connect.CarServiceName)))
	mux.Handle(grpcreflect.NewHandlerV1(grpcreflect.NewStaticReflector(carv1connect.CarServiceName)))
	mux.Handle(grpcreflect.NewHandlerV1Alpha(grpcreflect.NewStaticReflector(carv1connect.CarServiceName)))

	fmt.Printf("Registered car service handler with gRPC Connect\n")

	// Create HTTP server with timeout configuration
	s.httpServer = &http.Server{
		Addr:              fmt.Sprintf(":%d", s.httpPort),
		Handler:           h2c.NewHandler(mux, &http2.Server{}),
		ReadHeaderTimeout: 5 * time.Second, // Add timeout to prevent Slowloris attacks
	}
	fmt.Printf("Created HTTP server on port %d\n", s.httpPort)

	// Start HTTP server
	httpAddr := fmt.Sprintf(":%d", s.httpPort)
	httpListener, err := (&net.ListenConfig{}).Listen(ctx, "tcp", httpAddr)
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

	fmt.Printf("HTTP server with gRPC Connect started on port %d\n", s.httpPort)

	return nil
}

// Stop stops the HTTP server
func (s *Server) Stop() error {
	// Stop HTTP server
	if err := s.httpServer.Close(); err != nil {
		return fmt.Errorf("failed to stop HTTP server: %w", err)
	}

	return nil
}
