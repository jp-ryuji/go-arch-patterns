package output

import "time"

// CreateCar represents the response data for car creation
type CreateCar struct {
	ID        string    `json:"id"`
	TenantID  string    `json:"tenant_id"`
	Model     string    `json:"model"`
	CreatedAt time.Time `json:"created_at"`
}

// GetCar represents the response data for getting a single car
type GetCar struct {
	ID        string    `json:"id"`
	TenantID  string    `json:"tenant_id"`
	Model     string    `json:"model"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// ListCars represents the response data for listing cars
type ListCars struct {
	Cars          []CarSummary `json:"cars"`
	NextPageToken string       `json:"next_page_token,omitempty"`
	TotalCount    int32        `json:"total_count,omitempty"`
}

// CarSummary represents a summary view of a car for listing
type CarSummary struct {
	ID    string `json:"id"`
	Model string `json:"model"`
}
