package input

// GetCarByID represents the input data for retrieving a car by ID
type GetCarByID struct {
	ID string `validate:"required"`
}

// ListCars represents the input data for listing cars
type ListCars struct {
	TenantID  string `validate:"required"`
	PageSize  int32
	PageToken string
}

// CreateCar represents the input data for creating a car
type CreateCar struct {
	TenantID string `validate:"required"`
	Model    string `validate:"required"`
}
