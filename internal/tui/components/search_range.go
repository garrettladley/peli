package components

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/garrettladley/peli/internal/db"
	"github.com/garrettladley/peli/internal/xstrings"
)

var (
	searchRangeActive = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#ffffff")).
				Bold(true)

	searchRangeCurrent = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#22c55e")).
				Bold(true)

	searchRangeBracket = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#ffffff"))

	searchRangeEliminated = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#ef4444")).
				Faint(true)

	progressFilled = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#22c55e"))

	progressEmpty = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#4b5563"))

	helpStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#6b7280"))
)

// SearchRange visualizes a binary search progress.
type SearchRange struct {
	Ranked           []db.Restaurant
	Low              int
	High             int
	CurrentMid       int
	ContentWidth     int
	ProgressBarWidth int
	WindowContext    int
}

// NewSearchRange creates a new search range visualization.
func NewSearchRange(ranked []db.Restaurant, low, high, mid int) SearchRange {
	return SearchRange{
		Ranked:           ranked,
		Low:              low,
		High:             high,
		CurrentMid:       mid,
		ContentWidth:     70,
		ProgressBarWidth: 30,
		WindowContext:    2,
	}
}

// CalculateProgress returns the percentage of the search completed.
func (s SearchRange) CalculateProgress() int {
	total := len(s.Ranked)
	if total <= 1 {
		return 0
	}
	currentRange := s.High - s.Low + 1
	return int(float64(total-currentRange) / float64(total-1) * 100)
}

// Render returns the rendered search range, automatically choosing rich or compact.
func (s SearchRange) Render() string {
	total := len(s.Ranked)
	estimatedWidth := total * 5 // each position takes ~5 chars

	if estimatedWidth <= s.ContentWidth-4 {
		return s.RenderRich()
	}
	return s.RenderCompact()
}

// RenderRich renders a detailed visualization showing all positions.
func (s SearchRange) RenderRich() string {
	var b strings.Builder
	progress := s.CalculateProgress()

	// progress bar
	b.WriteString("  ")
	bar := NewProgressBar(s.ProgressBarWidth, float64(progress),
		WithFillStyle(progressFilled),
		WithEmptyStyle(progressEmpty),
	)
	b.WriteString(bar.Render())
	b.WriteString(helpStyle.Render(fmt.Sprintf(" %d%% narrowed", progress)))
	b.WriteString("\n\n")

	// position numbers
	b.WriteString("  ")
	for i := range s.Ranked {
		inRange := i >= s.Low && i <= s.High
		isCurrent := i == s.CurrentMid

		var cell string
		var style lipgloss.Style

		if isCurrent {
			cell = fmt.Sprintf("⟨%d⟩", i+1)
			style = searchRangeCurrent
		} else if inRange {
			cell = fmt.Sprintf("┃%d┃", i+1)
			style = searchRangeActive
		} else {
			cell = fmt.Sprintf(" %d ", i+1)
			style = searchRangeEliminated
		}

		b.WriteString(style.Render(cell))
		b.WriteString(" ")
	}
	b.WriteString("\n")

	// range indicator line
	b.WriteString("  ")
	for i := range s.Ranked {
		inRange := i >= s.Low && i <= s.High
		isCurrent := i == s.CurrentMid
		indicator, style := s.getRangeIndicator(i, isCurrent, inRange)
		b.WriteString(style.Render(indicator))
	}
	b.WriteString("\n")

	// restaurant names
	b.WriteString("  ")
	for i, r := range s.Ranked {
		inRange := i >= s.Low && i <= s.High
		isCurrent := i == s.CurrentMid

		name := xstrings.Truncate(r.Name, 4)
		padded := fmt.Sprintf("%-4s", name)

		var style lipgloss.Style
		if isCurrent {
			style = searchRangeCurrent
		} else if inRange {
			style = searchRangeActive
		} else {
			style = searchRangeEliminated
		}

		b.WriteString(style.Render(padded))
		b.WriteString(" ")
	}

	b.WriteString("\n\n")

	// legend
	b.WriteString("  ")
	b.WriteString(searchRangeCurrent.Render("⟨N⟩"))
	b.WriteString(helpStyle.Render(" comparing  "))
	b.WriteString(searchRangeActive.Render("┃N┃"))
	b.WriteString(helpStyle.Render(" in range  "))
	b.WriteString(searchRangeEliminated.Render(" N "))
	b.WriteString(helpStyle.Render(" eliminated"))

	return b.String()
}

// RenderCompact renders a condensed visualization for large lists.
func (s SearchRange) RenderCompact() string {
	var b strings.Builder
	total := len(s.Ranked)
	progress := s.CalculateProgress()

	// progress bar
	b.WriteString("  ")
	bar := NewProgressBar(s.ProgressBarWidth, float64(progress),
		WithFillStyle(progressFilled),
		WithEmptyStyle(progressEmpty),
	)
	b.WriteString(bar.Render())
	b.WriteString(helpStyle.Render(fmt.Sprintf(" %d%% narrowed", progress)))
	b.WriteString("\n\n")

	b.WriteString(fmt.Sprintf("  Range: #%d - #%d of %d\n\n", s.Low+1, s.High+1, total))

	// show window around active range
	windowStart := max(0, s.Low-s.WindowContext)
	windowEnd := min(total-1, s.High+s.WindowContext)

	if windowStart > 0 {
		b.WriteString(helpStyle.Render(fmt.Sprintf("  ... %d more above\n", windowStart)))
	}

	for i := windowStart; i <= windowEnd; i++ {
		r := s.Ranked[i]
		inRange := i >= s.Low && i <= s.High
		isCurrent := i == s.CurrentMid

		var status string
		var style lipgloss.Style
		if isCurrent {
			status = "► COMPARING"
			style = searchRangeCurrent
		} else if inRange {
			status = "  in range "
			style = searchRangeActive
		} else {
			status = "  eliminated"
			style = searchRangeEliminated
		}

		line := fmt.Sprintf("  #%-3d %s  %s", i+1, status, xstrings.TruncateWithEllipsis(r.Name, 20))
		b.WriteString(style.Render(line))
		b.WriteString("\n")
	}

	remaining := total - 1 - windowEnd
	if remaining > 0 {
		b.WriteString(helpStyle.Render(fmt.Sprintf("  ... %d more below\n", remaining)))
	}

	return b.String()
}

func (s SearchRange) getRangeIndicator(i int, isCurrent, inRange bool) (string, lipgloss.Style) {
	if isCurrent {
		return " ▲  ", searchRangeCurrent
	}
	if !inRange {
		return " ·  ", searchRangeEliminated
	}
	indicator := s.getRangeBracket(i)
	return indicator, searchRangeBracket
}

func (s SearchRange) getRangeBracket(i int) string {
	switch {
	case i == s.Low && i == s.High:
		return "═══ "
	case i == s.Low:
		return "╔══ "
	case i == s.High:
		return "══╗ "
	default:
		return "═══ "
	}
}
