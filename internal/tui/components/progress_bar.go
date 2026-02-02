package components

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

var (
	defaultFilledStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#22c55e"))
	defaultEmptyStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("#4b5563"))
)

// ProgressBar renders a horizontal progress bar.
type ProgressBar struct {
	Width      int
	Percentage float64
	FillStyle  lipgloss.Style
	EmptyStyle lipgloss.Style
	FillChar   string
	EmptyChar  string
}

// ProgressBarOption is a functional option for configuring a ProgressBar.
type ProgressBarOption func(*ProgressBar)

// NewProgressBar creates a new progress bar with the given width and percentage.
func NewProgressBar(width int, percentage float64, opts ...ProgressBarOption) ProgressBar {
	p := ProgressBar{
		Width:      width,
		Percentage: percentage,
		FillStyle:  defaultFilledStyle,
		EmptyStyle: defaultEmptyStyle,
		FillChar:   "━",
		EmptyChar:  "━",
	}
	for _, opt := range opts {
		opt(&p)
	}
	return p
}

// WithFillStyle sets the style for the filled portion.
func WithFillStyle(s lipgloss.Style) ProgressBarOption {
	return func(p *ProgressBar) {
		p.FillStyle = s
	}
}

// WithEmptyStyle sets the style for the empty portion.
func WithEmptyStyle(s lipgloss.Style) ProgressBarOption {
	return func(p *ProgressBar) {
		p.EmptyStyle = s
	}
}

// WithFillChar sets the character used for the filled portion.
func WithFillChar(c string) ProgressBarOption {
	return func(p *ProgressBar) {
		p.FillChar = c
	}
}

// WithEmptyChar sets the character used for the empty portion.
func WithEmptyChar(c string) ProgressBarOption {
	return func(p *ProgressBar) {
		p.EmptyChar = c
	}
}

// Render returns the rendered progress bar string.
func (p ProgressBar) Render() string {
	filled := min(int(p.Percentage/100.0*float64(p.Width)), p.Width)
	if filled < 0 {
		filled = 0
	}
	empty := p.Width - filled

	return p.FillStyle.Render(strings.Repeat(p.FillChar, filled)) +
		p.EmptyStyle.Render(strings.Repeat(p.EmptyChar, empty))
}
