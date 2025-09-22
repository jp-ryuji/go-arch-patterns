package output

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
