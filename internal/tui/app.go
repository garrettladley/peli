package tui

import (
	"context"
	"fmt"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/garrettladley/peli/internal/db"
	"github.com/garrettladley/peli/internal/ranking"
	"github.com/garrettladley/peli/internal/service/restaurant"
	"github.com/garrettladley/peli/internal/tui/components"
	"github.com/garrettladley/peli/internal/tui/views"
)

type screen int

const (
	screenHome screen = iota
	screenAdd
	screenCompare
	screenResult
)

const (
	keyQuit  = "ctrl+c"
	keyEsc   = "esc"
	keyEnter = "enter"
)

type model struct {
	screen      screen
	restaurants []db.Restaurant
	svc         restaurant.Service
	cfg         Config
	viewport    *components.Viewport
	err         error
	width       int
	height      int

	// add screen
	nameInput    textinput.Model
	cuisineInput textinput.Model
	inputFocus   int // 0 = name, 1 = cuisine

	// compare screen
	binarySearch *ranking.BinarySearchState
	selected     int // 0 = new (left), 1 = opponent (right)

	// result screen
	insertedAt    int
	newRestaurant db.Restaurant
}

func newModel(svc restaurant.Service, opts ...Option) model {
	cfg := DefaultConfig()
	cfg.Apply(opts...)

	nameInput := textinput.New()
	nameInput.Placeholder = "Restaurant name"
	nameInput.CharLimit = cfg.NameCharLimit
	nameInput.Width = cfg.InputWidth

	cuisineInput := textinput.New()
	cuisineInput.Placeholder = "Cuisine type"
	cuisineInput.CharLimit = cfg.CuisineCharLimit
	cuisineInput.Width = cfg.InputWidth

	return model{
		screen:       screenHome,
		svc:          svc,
		cfg:          cfg,
		viewport:     components.NewViewport(1, 0),
		nameInput:    nameInput,
		cuisineInput: cuisineInput,
	}
}

// Run starts the TUI application.
func Run(svc restaurant.Service, opts ...Option) error {
	p := tea.NewProgram(newModel(svc, opts...), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		return fmt.Errorf("run program: %w", err)
	}
	return nil
}

func (m model) Init() tea.Cmd {
	return m.loadRestaurants
}

type restaurantsLoadedMsg struct {
	restaurants []db.Restaurant
}

type errMsg struct {
	err error
}

func (m model) loadRestaurants() tea.Msg {
	ctx := context.Background()
	restaurants, err := m.svc.GetAll(ctx)
	if err != nil {
		return errMsg{err: err}
	}
	return restaurantsLoadedMsg{restaurants: restaurants}
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		// reserve: title(1) + divider(1) + blank(1) + colheader(1) + separator(1) + footer(3)
		viewportHeight := max(1, m.height-m.cfg.ViewportPadding)
		m.viewport.SetHeight(viewportHeight)
		return m, nil
	case tea.KeyMsg:
		return m.handleKey(msg)
	case restaurantsLoadedMsg:
		m.restaurants = msg.restaurants
		m.viewport.SetTotal(len(m.restaurants))
		return m, nil
	case errMsg:
		m.err = msg.err
		return m, nil
	}
	return m, nil
}

func (m model) View() string {
	var content string

	if m.err != nil {
		content = fmt.Sprintf("Error: %v\n\nPress 'q' to quit.", m.err)
	} else {
		switch m.screen {
		case screenHome:
			content = m.viewHome()
		case screenAdd:
			content = m.viewAdd()
		case screenCompare:
			content = m.viewCompare()
		case screenResult:
			content = m.viewResult()
		}
	}

	if m.width == 0 || m.height == 0 {
		return content
	}

	return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, content)
}

func (m model) viewHome() string {
	return views.RenderHome(views.HomeParams{
		Restaurants:     m.restaurants,
		Viewport:        m.viewport,
		NameTruncate:    m.cfg.NameTruncateWidth,
		CuisineTruncate: m.cfg.CuisineTruncateWidth,
		ScoreBarWidth:   15,
	})
}

func (m model) viewAdd() string {
	return views.RenderAdd(views.AddParams{
		NameInput:    m.nameInput,
		CuisineInput: m.cuisineInput,
		InputFocus:   m.inputFocus,
	})
}

func (m model) viewCompare() string {
	return views.RenderCompare(views.CompareParams{
		NewRestaurant:    m.newRestaurant,
		BinarySearch:     m.binarySearch,
		Selected:         m.selected,
		ContentWidth:     contentWidth,
		CardWidth:        m.cfg.CardWidth,
		ProgressBarWidth: m.cfg.ProgressBarWidth,
	})
}

func (m model) viewResult() string {
	return views.RenderResult(views.ResultParams{
		Restaurants:     m.restaurants,
		InsertedAt:      m.insertedAt,
		NewRestaurant:   m.newRestaurant,
		ContentWidth:    contentWidth,
		WindowContext:   m.cfg.ResultWindowContext,
		ComparisonCount: len(m.binarySearch.Comparisons),
	})
}
