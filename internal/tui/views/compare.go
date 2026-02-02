package views

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/garrettladley/peli/internal/db"
	"github.com/garrettladley/peli/internal/ranking"
	"github.com/garrettladley/peli/internal/tui/components"
	"github.com/garrettladley/peli/internal/xstrings"
)

const (
	decisionNameMaxLen = 15
)

var (
	compareTitleStyle = lipgloss.NewStyle().
				Bold(true).
				Foreground(lipgloss.Color("#4a90d9"))

	compareHelpStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#6b7280"))

	sectionHeader = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#ffffff")).
			Bold(true)

	vsStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#6b7280")).
		Bold(true).
		Padding(0, 2)

	positionStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#ffffff")).
			Bold(true)

	decisionWinStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#22c55e"))

	decisionLoseStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#ef4444"))
)

// CompareParams holds parameters for rendering the compare view.
type CompareParams struct {
	NewRestaurant    db.Restaurant
	BinarySearch     *ranking.BinarySearchState
	Selected         int // 0 = new (left), 1 = opponent (right)
	ContentWidth     int
	CardWidth        int
	ProgressBarWidth int
}

// RenderCompare renders the comparison screen.
func RenderCompare(p CompareParams) string {
	var b strings.Builder

	opponent := p.BinarySearch.CurrentOpponent()

	comparisonInfo := fmt.Sprintf("Comparison %d of ~%d",
		p.BinarySearch.ComparisonNumber(),
		p.BinarySearch.ExpectedComparisons())
	b.WriteString(renderCompareHeader(fmt.Sprintf("RANKING: %s", p.NewRestaurant.Name), comparisonInfo, p.ContentWidth))
	b.WriteString("\n")
	b.WriteString(strings.Repeat("─", p.ContentWidth))
	b.WriteString("\n\n")

	b.WriteString(lipgloss.PlaceHorizontal(p.ContentWidth, lipgloss.Center, "Which do you prefer?"))
	b.WriteString("\n\n")

	// render cards
	leftCard := components.NewCard(p.NewRestaurant, true, p.Selected == 0).WithWidth(p.CardWidth)
	rightCard := components.NewCard(opponent, false, p.Selected == 1).WithWidth(p.CardWidth)
	cards := components.RenderCards(leftCard, rightCard, vsStyle)
	b.WriteString(lipgloss.PlaceHorizontal(p.ContentWidth, lipgloss.Center, cards))
	b.WriteString("\n\n")

	// arrow hints
	leftHint := "←"
	rightHint := "→"
	if p.Selected == 0 {
		leftHint = positionStyle.Render(leftHint)
		rightHint = compareHelpStyle.Render(rightHint)
	} else {
		leftHint = compareHelpStyle.Render(leftHint)
		rightHint = positionStyle.Render(rightHint)
	}
	cardSpacing := p.CardWidth + 4 + p.CardWidth
	arrowLine := lipgloss.JoinHorizontal(lipgloss.Center,
		lipgloss.PlaceHorizontal(p.CardWidth, lipgloss.Center, leftHint),
		strings.Repeat(" ", 4),
		lipgloss.PlaceHorizontal(p.CardWidth, lipgloss.Center, rightHint),
	)
	b.WriteString(lipgloss.PlaceHorizontal(p.ContentWidth, lipgloss.Center,
		lipgloss.PlaceHorizontal(cardSpacing, lipgloss.Center, arrowLine)))
	b.WriteString("\n")

	b.WriteString("\n")
	b.WriteString(strings.Repeat("─", p.ContentWidth))
	b.WriteString("\n")
	b.WriteString(sectionHeader.Render("◆ SEARCH PROGRESS"))
	b.WriteString("\n\n")

	// search range visualization
	sr := components.NewSearchRange(
		p.BinarySearch.Ranked,
		p.BinarySearch.Low,
		p.BinarySearch.High,
		p.BinarySearch.CurrentMid,
	)
	sr.ContentWidth = p.ContentWidth
	sr.ProgressBarWidth = p.ProgressBarWidth
	b.WriteString(sr.Render())
	b.WriteString("\n")

	// decisions section - center header and tree as a unit
	b.WriteString("\n")
	b.WriteString(lipgloss.PlaceHorizontal(p.ContentWidth, lipgloss.Center, sectionHeader.Render("◆ DECISIONS")))
	b.WriteString("\n")
	if len(p.BinarySearch.Comparisons) > 0 {
		tree := renderDecisions(p.BinarySearch.Comparisons, p.NewRestaurant.ID)
		b.WriteString(lipgloss.PlaceHorizontal(p.ContentWidth, lipgloss.Center, strings.TrimSuffix(tree, "\n")))
		b.WriteString("\n")
	} else {
		b.WriteString(lipgloss.PlaceHorizontal(p.ContentWidth, lipgloss.Center, compareHelpStyle.Render("(none yet)")))
		b.WriteString("\n")
	}

	b.WriteString("\n")
	b.WriteString(lipgloss.PlaceHorizontal(p.ContentWidth, lipgloss.Center, compareHelpStyle.Render("← → select • enter confirm • esc cancel")))

	return b.String()
}

func renderCompareHeader(title, help string, width int) string {
	titleRendered := compareTitleStyle.Render(title)
	helpRendered := compareHelpStyle.Render(help)

	titleWidth := lipgloss.Width(titleRendered)
	helpWidth := lipgloss.Width(helpRendered)
	gap := max(width-titleWidth-helpWidth, 1)

	return titleRendered + strings.Repeat(" ", gap) + helpRendered
}

func renderDecisions(comparisons []ranking.CompareResult, newRestID int64) string {
	var result strings.Builder
	for _, comp := range comparisons {
		if comp.WinnerID == newRestID {
			result.WriteString(renderWinDecision(comp.LoserName))
		} else {
			result.WriteString(renderLoseDecision(comp.WinnerName))
		}
		result.WriteString("\n")
	}
	return result.String()
}

// fixed width: "▲ beat    " (10) + name (15) = 25
const decisionLineWidth = 10 + decisionNameMaxLen

func renderWinDecision(opponent string) string {
	line := decisionWinStyle.Render("▲ beat    ") +
		compareHelpStyle.Render(xstrings.Truncate(opponent, decisionNameMaxLen))
	if w := lipgloss.Width(line); w < decisionLineWidth {
		line += strings.Repeat(" ", decisionLineWidth-w)
	}
	return line
}

func renderLoseDecision(opponent string) string {
	line := decisionLoseStyle.Render("▼ lost to ") +
		compareHelpStyle.Render(xstrings.Truncate(opponent, decisionNameMaxLen))
	if w := lipgloss.Width(line); w < decisionLineWidth {
		line += strings.Repeat(" ", decisionLineWidth-w)
	}
	return line
}
