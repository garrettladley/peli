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

	decisionTreeBranch = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#ffffff"))
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

	// decisions section
	b.WriteString("\n")
	b.WriteString(sectionHeader.Render("◆ DECISIONS"))
	b.WriteString("\n")
	if len(p.BinarySearch.Comparisons) > 0 {
		b.WriteString(renderDecisionTree(p.BinarySearch.Comparisons, p.NewRestaurant))
	} else {
		b.WriteString(compareHelpStyle.Render("  (none yet)"))
		b.WriteString("\n")
	}

	b.WriteString("\n")
	b.WriteString(compareHelpStyle.Render("← → select • enter confirm • esc cancel"))

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

func renderDecisionTree(comparisons []ranking.CompareResult, newRest db.Restaurant) string {
	var b strings.Builder
	total := len(comparisons)

	for i, comp := range comparisons {
		var branch string
		if i == total-1 {
			branch = "  └──"
		} else {
			branch = "  ├──"
		}
		b.WriteString(decisionTreeBranch.Render(branch))
		b.WriteString(" ")

		if comp.WinnerID == newRest.ID {
			b.WriteString(decisionWinStyle.Render("▲ "))
			b.WriteString(decisionWinStyle.Render(xstrings.Truncate(newRest.Name, 15)))
			b.WriteString(compareHelpStyle.Render(" beats "))
			b.WriteString(decisionLoseStyle.Render(xstrings.Truncate(comp.LoserName, 15)))
		} else {
			b.WriteString(decisionLoseStyle.Render("▼ "))
			b.WriteString(decisionLoseStyle.Render(xstrings.Truncate(newRest.Name, 15)))
			b.WriteString(compareHelpStyle.Render(" under "))
			b.WriteString(decisionWinStyle.Render(xstrings.Truncate(comp.WinnerName, 15)))
		}
		b.WriteString("\n")
	}

	return b.String()
}
