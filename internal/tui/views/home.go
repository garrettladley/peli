package views

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/garrettladley/peli/internal/db"
	"github.com/garrettladley/peli/internal/tui/components"
	"github.com/garrettladley/peli/internal/xstrings"
)

const contentWidth = 70

var (
	green   = lipgloss.Color("#22c55e")
	white   = lipgloss.Color("#ffffff")
	gray    = lipgloss.Color("#6b7280")
	dimGray = lipgloss.Color("#4b5563")

	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#4a90d9"))

	subtitleStyle = lipgloss.NewStyle().
			Foreground(gray).
			Italic(true)

	listItemStyle = lipgloss.NewStyle().
			PaddingLeft(2).
			Foreground(white)

	selectedItemStyle = lipgloss.NewStyle().
				PaddingLeft(2).
				Foreground(white).
				Bold(true)

	helpStyle = lipgloss.NewStyle().
			Foreground(gray)

	scoreBarFilled = lipgloss.NewStyle().
			Foreground(green)

	scoreBarEmpty = lipgloss.NewStyle().
			Foreground(dimGray)
)

// HomeParams holds parameters for rendering the home view.
type HomeParams struct {
	Restaurants     []db.Restaurant
	Viewport        *components.Viewport
	NameTruncate    int
	CuisineTruncate int
	ScoreBarWidth   int
}

// RenderHome renders the home screen view.
func RenderHome(p HomeParams) string {
	var b strings.Builder

	b.WriteString(centerInContent(titleStyle.Render("PELI RANKINGS")))
	b.WriteString("\n")
	b.WriteString(renderDivider())
	b.WriteString("\n\n")

	if len(p.Restaurants) == 0 {
		b.WriteString(centerInContent(subtitleStyle.Render("No restaurants yet. Press 'a' to add one!")))
		b.WriteString("\n")
	} else {
		b.WriteString(fmt.Sprintf("   %-4s %-24s %-14s %-20s %s\n", "#", "Restaurant", "Cuisine", "Score", ""))
		b.WriteString("   " + strings.Repeat("─", contentWidth-5) + "\n")

		start, end := p.Viewport.VisibleRange()

		for i := start; i < end; i++ {
			r := p.Restaurants[i]
			cursor := "  "
			style := listItemStyle
			if i == p.Viewport.Cursor {
				cursor = "> "
				style = selectedItemStyle
			}

			scoreBar := renderScoreBar(r.Score, p.ScoreBarWidth)
			line := fmt.Sprintf("%s%-4d %-24s %-14s %s %.1f",
				cursor,
				r.Position+1,
				xstrings.Truncate(r.Name, p.NameTruncate),
				xstrings.Truncate(r.Cuisine, p.CuisineTruncate),
				scoreBar,
				r.Score,
			)
			b.WriteString(style.Render(line))
			b.WriteString("\n")
		}

		if len(p.Restaurants) > p.Viewport.Height {
			scrollInfo := fmt.Sprintf("  [%d-%d of %d]", start+1, end, len(p.Restaurants))
			b.WriteString(helpStyle.Render(scrollInfo))
			b.WriteString("\n")
		}
	}

	b.WriteString("\n")
	b.WriteString(renderDivider())
	b.WriteString("\n")
	b.WriteString(centerInContent(helpStyle.Render("↑↓ navigate • a add • d delete • c clear added • q quit")))

	return b.String()
}

func renderScoreBar(score float64, width int) string {
	filled := int((score / 10.0) * float64(width))
	empty := width - filled

	return scoreBarFilled.Render(strings.Repeat("█", filled)) +
		scoreBarEmpty.Render(strings.Repeat("░", empty))
}

func renderDivider() string {
	return strings.Repeat("─", contentWidth)
}

func centerInContent(s string) string {
	return lipgloss.PlaceHorizontal(contentWidth, lipgloss.Center, s)
}
