package restaurant

import (
	"context"
	"fmt"

	"github.com/garrettladley/peli/internal/db"
	"github.com/garrettladley/peli/internal/ranking"
)

type service struct {
	queries db.Querier
}

// New creates a new restaurant service.
func New(queries db.Querier) Service {
	return &service{queries: queries}
}

func (s *service) GetAll(ctx context.Context) ([]db.Restaurant, error) {
	restaurants, err := s.queries.GetAllRestaurants(ctx)
	if err != nil {
		return nil, fmt.Errorf("get all restaurants: %w", err)
	}
	return restaurants, nil
}

func (s *service) Create(ctx context.Context, params db.CreateRestaurantParams) (db.Restaurant, error) {
	restaurant, err := s.queries.CreateRestaurant(ctx, params)
	if err != nil {
		return db.Restaurant{}, fmt.Errorf("create restaurant: %w", err)
	}
	return restaurant, nil
}

func (s *service) UpdatePosition(ctx context.Context, id int64, position int64) error {
	restaurants, err := s.queries.GetAllRestaurants(ctx)
	if err != nil {
		return fmt.Errorf("get restaurants: %w", err)
	}

	var target *db.Restaurant
	for i := range restaurants {
		if restaurants[i].ID == id {
			target = &restaurants[i]
			break
		}
	}
	if target == nil {
		return fmt.Errorf("restaurant %d not found", id)
	}

	if err := s.queries.UpdateRestaurant(ctx, db.UpdateRestaurantParams{
		ID:       id,
		Elo:      target.Elo,
		Score:    target.Score,
		Position: position,
	}); err != nil {
		return fmt.Errorf("update position: %w", err)
	}

	return s.recalculateScores(ctx)
}

func (s *service) ShiftPositions(ctx context.Context, fromPosition int64) error {
	if err := s.queries.UpdatePositions(ctx, fromPosition); err != nil {
		return fmt.Errorf("shift positions: %w", err)
	}
	return nil
}

func (s *service) Delete(ctx context.Context, id int64) error {
	if err := s.queries.DeleteRestaurant(ctx, id); err != nil {
		return fmt.Errorf("delete restaurant: %w", err)
	}
	return s.renumberAndRecalculate(ctx)
}

func (s *service) Seed(ctx context.Context, restaurants []SeedRestaurant) error {
	if err := s.queries.DeleteAllComparisons(ctx); err != nil {
		return fmt.Errorf("delete comparisons: %w", err)
	}
	if err := s.queries.DeleteAllRestaurants(ctx); err != nil {
		return fmt.Errorf("delete restaurants: %w", err)
	}

	n := len(restaurants)
	for _, r := range restaurants {
		// calculate score based on position: top = 10, bottom = 10/n
		posFromBottom := float64(n - int(r.Position))
		score := 10.0 * (posFromBottom / float64(n))

		if _, err := s.queries.CreateSeedRestaurant(ctx, db.CreateSeedRestaurantParams{
			Name:     r.Name,
			Cuisine:  r.Cuisine,
			Elo:      1500.0,
			Score:    score,
			Position: r.Position,
		}); err != nil {
			return fmt.Errorf("create seed restaurant %q: %w", r.Name, err)
		}
	}

	return nil
}

func (s *service) Reset(ctx context.Context) error {
	if err := s.queries.DeleteNonSeedRestaurants(ctx); err != nil {
		return fmt.Errorf("delete non-seed: %w", err)
	}
	return s.renumberAndRecalculate(ctx)
}

func (s *service) RecordComparison(ctx context.Context, params db.RecordComparisonParams) error {
	if err := s.queries.RecordComparison(ctx, params); err != nil {
		return fmt.Errorf("record comparison: %w", err)
	}
	return nil
}

// renumberAndRecalculate renumbers positions (closing gaps) and recalculates scores.
func (s *service) renumberAndRecalculate(ctx context.Context) error {
	restaurants, err := s.queries.GetAllRestaurants(ctx)
	if err != nil {
		return fmt.Errorf("get restaurants: %w", err)
	}

	if len(restaurants) == 0 {
		return nil
	}

	// renumber positions to close gaps (0, 1, 2, ...)
	for i := range restaurants {
		restaurants[i].Position = int64(i)
	}

	// recalculate scores based on new positions
	ranking.RecalculateAllScores(restaurants)

	// persist both position and score changes
	for _, r := range restaurants {
		if err := s.queries.UpdateRestaurant(ctx, db.UpdateRestaurantParams{
			ID:       r.ID,
			Elo:      r.Elo,
			Score:    r.Score,
			Position: r.Position,
		}); err != nil {
			return fmt.Errorf("update restaurant %d: %w", r.ID, err)
		}
	}

	return nil
}

// recalculateScores fetches all restaurants and updates their scores based on position.
func (s *service) recalculateScores(ctx context.Context) error {
	restaurants, err := s.queries.GetAllRestaurants(ctx)
	if err != nil {
		return fmt.Errorf("get restaurants: %w", err)
	}

	if len(restaurants) == 0 {
		return nil
	}

	ranking.RecalculateAllScores(restaurants)

	for _, r := range restaurants {
		if err := s.queries.UpdateRestaurant(ctx, db.UpdateRestaurantParams{
			ID:       r.ID,
			Elo:      r.Elo,
			Score:    r.Score,
			Position: r.Position,
		}); err != nil {
			return fmt.Errorf("update score for %d: %w", r.ID, err)
		}
	}

	return nil
}
