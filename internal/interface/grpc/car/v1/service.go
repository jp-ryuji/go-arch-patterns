package car

import (
	"context"

	carv1 "github.com/jp-ryuji/go-arch-patterns/api/generated/car/v1"
	"github.com/jp-ryuji/go-arch-patterns/internal/application/dto"
	"github.com/jp-ryuji/go-arch-patterns/internal/application/service"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// CarServiceServer implements the gRPC service for car operations
type CarServiceServer struct {
	carv1.UnimplementedCarServiceServer
	carService service.CarService
}

// NewCarServiceServer creates a new CarServiceServer
func NewCarServiceServer(carService service.CarService) *CarServiceServer {
	return &CarServiceServer{
		carService: carService,
	}
}

// CreateCar creates a new car
func (s *CarServiceServer) CreateCar(ctx context.Context, req *carv1.CreateCarRequest) (*carv1.CreateCarResponse, error) {
	// Convert gRPC request to application DTO
	input := dto.RegisterCarInput{
		TenantID: req.GetTenantId(),
		Model:    req.GetModel(),
	}

	// Call application service
	car, err := s.carService.Register(ctx, input)
	if err != nil {
		return nil, err
	}

	// Convert entity to gRPC response
	response := &carv1.CreateCarResponse{
		Car: &carv1.Car{
			Id:        car.ID,
			TenantId:  car.TenantID,
			Model:     car.Model,
			CreatedAt: timestamppb.New(car.CreatedAt),
			UpdatedAt: timestamppb.New(car.UpdatedAt),
		},
	}

	return response, nil
}

// GetCar retrieves a car by ID
func (s *CarServiceServer) GetCar(ctx context.Context, req *carv1.GetCarRequest) (*carv1.GetCarResponse, error) {
	// Convert gRPC request to application DTO
	input := dto.GetCarByIDInput{
		ID: req.GetId(),
	}

	// Call application service
	car, err := s.carService.GetByID(ctx, input)
	if err != nil {
		return nil, err
	}

	// Convert entity to gRPC response
	response := &carv1.GetCarResponse{
		Car: &carv1.Car{
			Id:        car.ID,
			TenantId:  car.TenantID,
			Model:     car.Model,
			CreatedAt: timestamppb.New(car.CreatedAt),
			UpdatedAt: timestamppb.New(car.UpdatedAt),
		},
	}

	return response, nil
}

// ListCars retrieves a list of cars
func (s *CarServiceServer) ListCars(ctx context.Context, req *carv1.ListCarsRequest) (*carv1.ListCarsResponse, error) {
	// Convert gRPC request to application DTO
	input := dto.ListCarsInput{
		TenantID:  req.GetTenantId(),
		PageSize:  req.GetPageSize(),
		PageToken: req.GetPageToken(),
	}

	// Call application service
	cars, err := s.carService.List(ctx, input)
	if err != nil {
		return nil, err
	}

	// Convert entities to gRPC response
	grpcCars := make([]*carv1.Car, len(*cars))
	for i, car := range *cars {
		grpcCars[i] = &carv1.Car{
			Id:        car.ID,
			TenantId:  car.TenantID,
			Model:     car.Model,
			CreatedAt: timestamppb.New(car.CreatedAt),
			UpdatedAt: timestamppb.New(car.UpdatedAt),
		}
	}

	// TODO: Implement proper pagination token generation
	response := &carv1.ListCarsResponse{
		Cars:          grpcCars,
		NextPageToken: "",
	}

	return response, nil
}
