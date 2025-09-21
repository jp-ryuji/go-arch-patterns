package dto

// ListCarsInput represents the input data for listing cars
type ListCarsInput struct {
	TenantID  string `validate:"required"`
	PageSize  int32
	PageToken string
}
