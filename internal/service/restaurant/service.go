package restaurant

import (
	"context"

	"github.com/garrettladley/peli/internal/db"
)

// Service defines restaurant operations with automatic score management.
type Service interface {
	// GetAll returns all restaurants ordered by position.
	GetAll(ctx context.Context) ([]db.Restaurant, error)

	// Create creates a new restaurant with initial values.
	// Position should be set to len(existing) as a placeholder.
	Create(ctx context.Context, params db.CreateRestaurantParams) (db.Restaurant, error)

	// UpdatePosition updates a restaurant's position and recalculates all scores.
	UpdatePosition(ctx context.Context, id int64, position int64) error

	// ShiftPositions increments positions >= the given position.
	ShiftPositions(ctx context.Context, fromPosition int64) error

	// Delete removes a restaurant and recalculates scores.
	Delete(ctx context.Context, id int64) error

	// Seed clears all data and inserts seed restaurants with recalculated scores.
	Seed(ctx context.Context, restaurants []SeedRestaurant) error

	// Reset deletes non-seed restaurants and recalculates scores.
	Reset(ctx context.Context) error

	// RecordComparison records a comparison result.
	RecordComparison(ctx context.Context, params db.RecordComparisonParams) error
}

// SeedRestaurant represents a restaurant for seeding.
type SeedRestaurant struct {
	Name     string
	Cuisine  string
	Position int64 // 0-indexed
}
