package tui

import (
	"context"

	"github.com/garrettladley/peli/internal/service/restaurant"
)

// Deps holds external dependencies for the TUI.
type Deps struct {
	Ctx     context.Context
	Service restaurant.Service
}

// NewDeps creates a new Deps with the given service.
func NewDeps(svc restaurant.Service) Deps {
	return Deps{
		Ctx:     context.Background(),
		Service: svc,
	}
}
