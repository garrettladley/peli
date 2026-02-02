package components

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
	"github.com/garrettladley/peli/internal/db"
)

var (
	dimGray = lipgloss.Color("#4b5563")
	white   = lipgloss.Color("#ffffff")

	cardBaseStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			Padding(1, 3).
			Align(lipgloss.Center).
			Foreground(white)

	cardNormalBorder   = dimGray
	cardSelectedBorder = white
)

// Card renders a comparison card for a restaurant.
type Card struct {
	Restaurant db.Restaurant
	IsNew      bool
	Selected   bool
	Width      int
}

// NewCard creates a new comparison card.
func NewCard(r db.Restaurant, isNew, selected bool) Card {
	return Card{
		Restaurant: r,
		IsNew:      isNew,
		Selected:   selected,
		Width:      24,
	}
}

// WithWidth sets the card width.
func (c Card) WithWidth(w int) Card {
	c.Width = w
	return c
}

// Render returns the rendered card string.
func (c Card) Render() string {
	style := cardBaseStyle.Width(c.Width)

	if c.Selected {
		style = style.BorderForeground(cardSelectedBorder).Bold(true)
	} else {
		style = style.BorderForeground(cardNormalBorder)
	}

	var content string
	if c.IsNew {
		content = fmt.Sprintf("%s\n(new)", c.Restaurant.Name)
	} else {
		content = fmt.Sprintf("%s\n%.1f", c.Restaurant.Name, c.Restaurant.Score)
	}

	return style.Render(content)
}

// RenderCards renders two cards side by side with a "vs" separator.
func RenderCards(left, right Card, vsStyle lipgloss.Style) string {
	leftRendered := left.Render()
	rightRendered := right.Render()
	vs := vsStyle.Render("vs")

	return lipgloss.JoinHorizontal(lipgloss.Center, leftRendered, vs, rightRendered)
}
