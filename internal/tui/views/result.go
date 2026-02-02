package views

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/garrettladley/peli/internal/db"
	"github.com/garrettladley/peli/internal/xstrings"
)

var (
	resultTitleStyle = lipgloss.NewStyle().
				Bold(true).
				Foreground(lipgloss.Color("#4a90d9"))

	resultHelpStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#6b7280"))

	resultBoxStyle = lipgloss.NewStyle().
			Border(lipgloss.DoubleBorder()).
			BorderForeground(lipgloss.Color("#22c55e")).
			Padding(1, 3).
			Align(lipgloss.Center).
			Foreground(lipgloss.Color("#ffffff"))

	resultSectionHeader = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#ffffff")).
				Bold(true)

	resultPositionStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#ffffff")).
				Bold(true)

	resultInactiveStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#4b5563"))
)

// ResultParams holds parameters for rendering the result view.
type ResultParams struct {
	Restaurants     []db.Restaurant
	InsertedAt      int
	NewRestaurant   db.Restaurant
	ContentWidth    int
	WindowContext   int
	ComparisonCount int
}

// RenderResult renders the result screen after ranking is complete.
func RenderResult(p ResultParams) string {
	var b strings.Builder

	b.WriteString(renderResultHeader("RANKED!", "enter done", p.ContentWidth))
	b.WriteString("\n")
	b.WriteString(strings.Repeat("─", p.ContentWidth))
	b.WriteString("\n\n")

	resultContent := fmt.Sprintf("%s\n\nRanked #%d of %d",
		p.NewRestaurant.Name,
		p.InsertedAt+1,
		len(p.Restaurants))
	b.WriteString(lipgloss.PlaceHorizontal(p.ContentWidth, lipgloss.Center, resultBoxStyle.Render(resultContent)))
	b.WriteString("\n\n")

	b.WriteString(resultSectionHeader.Render("◆ FINAL RANKING"))
	b.WriteString("\n")

	// show window around inserted position
	windowStart := max(0, p.InsertedAt-p.WindowContext)
	windowEnd := min(len(p.Restaurants)-1, p.InsertedAt+p.WindowContext)

	if windowStart > 0 {
		b.WriteString("... ")
	}
	for i := windowStart; i <= windowEnd; i++ {
		r := p.Restaurants[i]
		style := resultInactiveStyle
		cell := fmt.Sprintf("[%d]", i+1)
		if r.ID == p.NewRestaurant.ID {
			style = resultPositionStyle
			cell = fmt.Sprintf("*%d*", i+1)
		}
		b.WriteString(style.Render(cell))
		b.WriteString(" ")
	}
	if windowEnd < len(p.Restaurants)-1 {
		b.WriteString("...")
	}
	b.WriteString("\n")

	// restaurant names under positions
	if windowStart > 0 {
		b.WriteString("    ")
	}
	for i := windowStart; i <= windowEnd; i++ {
		r := p.Restaurants[i]
		style := resultInactiveStyle
		name := xstrings.Truncate(r.Name, 4)
		if r.ID == p.NewRestaurant.ID {
			style = resultPositionStyle
		}
		b.WriteString(style.Render(name))
		b.WriteString(" ")
	}
	b.WriteString("\n\n")

	b.WriteString(resultHelpStyle.Render(fmt.Sprintf("Completed in %d comparison(s)", p.ComparisonCount)))

	return b.String()
}

func renderResultHeader(title, help string, width int) string {
	titleRendered := resultTitleStyle.Render(title)
	helpRendered := resultHelpStyle.Render(help)

	titleWidth := lipgloss.Width(titleRendered)
	helpWidth := lipgloss.Width(helpRendered)
	gap := max(width-titleWidth-helpWidth, 1)

	return titleRendered + strings.Repeat(" ", gap) + helpRendered
}
