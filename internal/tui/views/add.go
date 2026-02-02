package views

import (
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/lipgloss"
)

const formWidth = 50

var (
	inputLabelStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#ffffff")).
			Bold(true)

	inputBoxStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#4b5563")).
			Padding(0, 1).
			Width(44)

	inputBoxFocusedStyle = lipgloss.NewStyle().
				Border(lipgloss.RoundedBorder()).
				BorderForeground(lipgloss.Color("#ffffff")).
				Padding(0, 1).
				Width(44)

	addTitleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#4a90d9"))

	addHelpStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#6b7280"))
)

// AddParams holds parameters for rendering the add view.
type AddParams struct {
	NameInput    textinput.Model
	CuisineInput textinput.Model
	InputFocus   int // 0 = name, 1 = cuisine
}

// RenderAdd renders the add restaurant screen.
func RenderAdd(p AddParams) string {
	var b strings.Builder

	b.WriteString(lipgloss.PlaceHorizontal(formWidth, lipgloss.Center, addTitleStyle.Render("ADD RESTAURANT")))
	b.WriteString("\n")
	b.WriteString(strings.Repeat("─", formWidth))
	b.WriteString("\n\n")

	// name input
	b.WriteString(inputLabelStyle.Render("Name"))
	b.WriteString("\n")
	nameBox := inputBoxStyle
	if p.InputFocus == 0 {
		nameBox = inputBoxFocusedStyle
	}
	b.WriteString(nameBox.Render(p.NameInput.View()))
	b.WriteString("\n\n")

	// cuisine input
	b.WriteString(inputLabelStyle.Render("Cuisine"))
	b.WriteString("\n")
	cuisineBox := inputBoxStyle
	if p.InputFocus == 1 {
		cuisineBox = inputBoxFocusedStyle
	}
	b.WriteString(cuisineBox.Render(p.CuisineInput.View()))
	b.WriteString("\n\n")

	// help text
	var help string
	if p.NameInput.Value() != "" && p.CuisineInput.Value() != "" {
		help = "enter submit • esc cancel"
	} else {
		help = "tab next • esc cancel"
	}
	b.WriteString(lipgloss.PlaceHorizontal(formWidth, lipgloss.Center, addHelpStyle.Render(help)))

	return b.String()
}
