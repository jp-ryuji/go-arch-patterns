package output

import (
	"github.com/jp-ryuji/go-arch-patterns/internal/domain/entity"
)

// CarEntityToSummary converts a domain Car entity to CarSummary DTO
func CarEntityToSummary(car *entity.Car) CarSummary {
	return CarSummary{
		ID:    car.ID,
		Model: car.Model,
	}
}

// CarEntitiesToList converts multiple Car entities to ListCars output DTO
func CarEntitiesToList(cars []*entity.Car, nextPageToken string, totalCount int32) *ListCars {
	summaries := make([]CarSummary, len(cars))
	for i, car := range cars {
		summaries[i] = CarEntityToSummary(car)
	}

	return &ListCars{
		Cars:          summaries,
		NextPageToken: nextPageToken,
		TotalCount:    totalCount,
	}
}
