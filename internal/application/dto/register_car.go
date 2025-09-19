package dto

// RegisterCarInput represents the input data for registering a car
type RegisterCarInput struct {
	TenantID string `validate:"required"`
	Model    string `validate:"required"`
}
