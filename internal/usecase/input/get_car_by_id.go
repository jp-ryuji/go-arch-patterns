package input

// GetCarByIDInput represents the input data for retrieving a car by ID
type GetCarByIDInput struct {
	ID string `validate:"required"`
}
