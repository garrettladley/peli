package tui

import (
	"context"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/garrettladley/peli/internal/db"
	"github.com/garrettladley/peli/internal/ranking"
)

func (m model) handleKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch m.screen {
	case screenHome:
		return m.handleHomeKey(msg)
	case screenAdd:
		return m.handleAddKey(msg)
	case screenCompare:
		return m.handleCompareKey(msg)
	case screenResult:
		return m.handleResultKey(msg)
	}
	return m, nil
}

func (m model) handleHomeKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "q", keyQuit:
		return m, tea.Quit
	case "a":
		m.screen = screenAdd
		m.nameInput.Reset()
		m.cuisineInput.Reset()
		m.nameInput.Focus()
		m.cuisineInput.Blur()
		m.inputFocus = 0
		return m, nil
	case "up", "k":
		m.viewport.MoveUp()
	case "down", "j":
		m.viewport.MoveDown()
	case "pgup", "ctrl+u":
		m.viewport.PageUp()
	case "pgdown", "ctrl+d":
		m.viewport.PageDown()
	case "home", "g":
		m.viewport.GoToStart()
	case "end", "G":
		m.viewport.GoToEnd()
	case "d", "backspace":
		if len(m.restaurants) > 0 {
			return m.deleteRestaurant()
		}
	case "c":
		return m.clearNonSeedData()
	}
	return m, nil
}

func (m model) deleteRestaurant() (tea.Model, tea.Cmd) {
	if m.viewport.Cursor >= len(m.restaurants) {
		return m, nil
	}
	ctx := context.Background()
	r := m.restaurants[m.viewport.Cursor]
	if err := m.svc.Delete(ctx, r.ID); err != nil {
		m.err = err
		return m, nil
	}
	if m.viewport.Cursor > 0 {
		m.viewport.Cursor--
	}
	m.viewport.ClampOffset()
	m.viewport.EnsureVisible()
	return m, m.loadRestaurants
}

func (m model) clearNonSeedData() (tea.Model, tea.Cmd) {
	ctx := context.Background()
	if err := m.svc.Reset(ctx); err != nil {
		m.err = err
		return m, nil
	}
	m.viewport.Cursor = 0
	m.viewport.Offset = 0
	return m, m.loadRestaurants
}

func (m model) handleAddKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case keyEsc:
		m.screen = screenHome
		return m, nil
	case "tab", "down":
		m.inputFocus = (m.inputFocus + 1) % 2
		m.updateInputFocus()
		return m, nil
	case "shift+tab", "up":
		m.inputFocus = 1 - m.inputFocus // fixed: was (m.inputFocus + 1) % 2
		m.updateInputFocus()
		return m, nil
	case keyEnter:
		if m.nameInput.Value() != "" && m.cuisineInput.Value() != "" {
			return m.startRanking()
		}
		return m, nil
	default:
		return m.updateInputs(msg)
	}
}

func (m *model) updateInputFocus() {
	if m.inputFocus == 0 {
		m.nameInput.Focus()
		m.cuisineInput.Blur()
	} else {
		m.nameInput.Blur()
		m.cuisineInput.Focus()
	}
}

func (m model) updateInputs(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	if m.inputFocus == 0 {
		m.nameInput, cmd = m.nameInput.Update(msg)
	} else {
		m.cuisineInput, cmd = m.cuisineInput.Update(msg)
	}
	return m, cmd
}

func (m model) startRanking() (tea.Model, tea.Cmd) {
	ctx := context.Background()

	newR, err := m.svc.Create(ctx, db.CreateRestaurantParams{
		Name:     m.nameInput.Value(),
		Cuisine:  m.cuisineInput.Value(),
		Elo:      m.cfg.InitialElo,
		Score:    m.cfg.InitialScore,
		Position: int64(len(m.restaurants)),
	})
	if err != nil {
		m.err = err
		return m, nil
	}

	m.newRestaurant = newR
	m.binarySearch = ranking.NewBinarySearch(newR, m.restaurants)

	if m.binarySearch.Done {
		return m.finishRanking()
	}

	m.screen = screenCompare
	m.selected = 0
	return m, nil
}

func (m model) handleCompareKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "q", keyQuit:
		return m, tea.Quit
	case keyEsc:
		ctx := context.Background()
		_ = m.svc.Delete(ctx, m.newRestaurant.ID)
		m.screen = screenHome
		return m, m.loadRestaurants
	case "left", "h", "1":
		m.selected = 0
		return m, nil
	case "right", "l", "2":
		m.selected = 1
		return m, nil
	case keyEnter, " ":
		newWins := m.selected == 0
		m.binarySearch.Step(newWins)

		if m.binarySearch.Done {
			return m.finishRanking()
		}
		return m, nil
	}
	return m, nil
}

func (m model) finishRanking() (tea.Model, tea.Cmd) {
	ctx := context.Background()
	insertPos := m.binarySearch.InsertAt

	if err := m.svc.ShiftPositions(ctx, int64(insertPos)); err != nil {
		m.err = err
		return m, nil
	}

	if err := m.svc.UpdatePosition(ctx, m.newRestaurant.ID, int64(insertPos)); err != nil {
		m.err = err
		return m, nil
	}

	for _, comp := range m.binarySearch.Comparisons {
		if err := m.svc.RecordComparison(ctx, db.RecordComparisonParams{
			WinnerID:   comp.WinnerID,
			WinnerName: comp.WinnerName,
			LoserID:    comp.LoserID,
			LoserName:  comp.LoserName,
		}); err != nil {
			m.err = err
			return m, nil
		}
	}

	m.insertedAt = insertPos
	m.screen = screenResult

	return m, m.loadRestaurants
}

func (m model) handleResultKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "q", keyQuit:
		return m, tea.Quit
	case keyEnter, keyEsc, " ":
		m.screen = screenHome
		return m, nil
	}
	return m, nil
}
