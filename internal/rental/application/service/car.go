package rentalservice

import (
	"context"

	rental "github.com/jp-ryuji/go-arch-patterns/internal/rental"
	rentalinput "github.com/jp-ryuji/go-arch-patterns/internal/rental/application/input"
	rentaloutput "github.com/jp-ryuji/go-arch-patterns/internal/rental/application/output"
)

// CarService defines the interface for car-related business logic
//
//go:generate go run go.uber.org/mock/mockgen -source=$GOFILE -destination=mock/$GOFILE -package=mockrentalservice
type CarService interface {
	Create(ctx context.Context, input rentalinput.CreateCar) (*rental.Car, error)
	GetByID(ctx context.Context, input rentalinput.GetCarByID) (*rental.Car, error)
	GetByIDWithTenant(ctx context.Context, input rentalinput.GetCarByID) (*rental.Car, error)
	List(ctx context.Context, input rentalinput.ListCars) (*rentaloutput.ListCars, error)
}
