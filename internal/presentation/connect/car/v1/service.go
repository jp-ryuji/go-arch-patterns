package car

import (
	"context"

	"connectrpc.com/connect"
	carv1 "github.com/jp-ryuji/go-arch-patterns/api/generated/car/v1"
	"github.com/jp-ryuji/go-arch-patterns/internal/application/input"
	"github.com/jp-ryuji/go-arch-patterns/internal/application/service"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// CarServiceHandler implements the Connect service for car operations
type CarServiceHandler struct {
	carService service.CarService
}

// NewCarServiceHandler creates a new CarServiceHandler
func NewCarServiceHandler(carService service.CarService) *CarServiceHandler {
	return &CarServiceHandler{
		carService: carService,
	}
}

// CreateCar creates a new car
func (h *CarServiceHandler) CreateCar(ctx context.Context, req *connect.Request[carv1.CreateCarRequest]) (*connect.Response[carv1.CreateCarResponse], error) {
	// Convert Connect request to application DTO
	input := input.CreateCar{
		TenantID: req.Msg.GetTenantId(),
		Model:    req.Msg.GetModel(),
	}

	// Call application service
	carOutput, err := h.carService.Create(ctx, input)
	if err != nil {
		return nil, err
	}

	// Convert DTO to Connect response
	response := &carv1.CreateCarResponse{
		Car: &carv1.Car{
			Id:        carOutput.ID,
			TenantId:  carOutput.TenantID,
			Model:     carOutput.Model,
			CreatedAt: timestamppb.New(carOutput.CreatedAt),
		},
	}

	return connect.NewResponse(response), nil
}

// GetCar retrieves a car by ID
func (h *CarServiceHandler) GetCar(ctx context.Context, req *connect.Request[carv1.GetCarRequest]) (*connect.Response[carv1.GetCarResponse], error) {
	// Convert Connect request to application DTO
	input := input.GetCarByID{
		ID: req.Msg.GetId(),
	}

	// Call application service
	carOutput, err := h.carService.GetByID(ctx, input)
	if err != nil {
		return nil, err
	}

	// Convert DTO to Connect response
	response := &carv1.GetCarResponse{
		Car: &carv1.Car{
			Id:        carOutput.ID,
			TenantId:  carOutput.TenantID,
			Model:     carOutput.Model,
			CreatedAt: timestamppb.New(carOutput.CreatedAt),
			UpdatedAt: timestamppb.New(carOutput.UpdatedAt),
		},
	}

	return connect.NewResponse(response), nil
}

// ListCars retrieves a list of cars
func (h *CarServiceHandler) ListCars(ctx context.Context, req *connect.Request[carv1.ListCarsRequest]) (*connect.Response[carv1.ListCarsResponse], error) {
	// Convert Connect request to application DTO
	input := input.ListCars{
		TenantID:  req.Msg.GetTenantId(),
		PageSize:  req.Msg.GetPageSize(),
		PageToken: req.Msg.GetPageToken(),
	}

	// Call application service
	listOutput, err := h.carService.List(ctx, input)
	if err != nil {
		return nil, err
	}

	// Convert DTOs to Connect response
	grpcCars := make([]*carv1.Car, len(listOutput.Cars))
	for i, carSummary := range listOutput.Cars {
		grpcCars[i] = &carv1.Car{
			Id:    carSummary.ID,
			Model: carSummary.Model,
		}
	}

	response := &carv1.ListCarsResponse{
		Cars:          grpcCars,
		NextPageToken: listOutput.NextPageToken,
	}

	return connect.NewResponse(response), nil
}
