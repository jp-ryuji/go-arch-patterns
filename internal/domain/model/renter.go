package model

// Renter is an interface that represents a renter entity.
// Both Company and Individual implement this interface.
type Renter interface {
	// GetID returns the unique identifier of the renter.
	GetID() string
	// GetTenantID returns the tenant identifier associated with the renter.
	GetTenantID() string
	// GetType returns the type of the renter (e.g., "company" or "individual").
	GetType() RenterModel
}
