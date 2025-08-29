package input

// UpdateCarInput represents the input data for updating a car
type UpdateCarInput struct {
	ID       string `validate:"required"`
	TenantID string `validate:"required"`
	Model    string `validate:"required"`
}
