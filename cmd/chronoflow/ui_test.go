package main

import (
	"io"
	"testing"
	"time"

	"ctchen222/chronoflow/internal/domain"
	"ctchen222/chronoflow/internal/repository"
	"ctchen222/chronoflow/internal/service"
	"ctchen222/chronoflow/internal/ui"
	"ctchen222/chronoflow/pkg/calendar"
	"ctchen222/chronoflow/pkg/todo"

	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/x/exp/teatest"
)

// newTestModel creates a model configured for testing with sample data
func newTestModel(t *testing.T) *model {
	t.Helper()

	// Use temp file for repository
	tmpFile := t.TempDir() + "/test_todos.json"
	repo := repository.NewJSONTodoRepository(tmpFile)

	// Fixed time for consistent output
	fixedTime := time.Date(2026, 1, 18, 10, 0, 0, 0, time.UTC)
	timeProv := service.NewMockTimeProvider(fixedTime)

	// Add sample todos for today
	repo.Add(fixedTime, domain.Todo{
		Title:    "Complete TUI visual testing",
		Desc:     "Add teatest support",
		Complete: false,
		Priority: domain.PriorityHigh,
	})
	repo.Add(fixedTime, domain.Todo{
		Title:    "Review pull request",
		Desc:     "",
		Complete: true,
		Priority: domain.PriorityMedium,
	})
	repo.Add(fixedTime, domain.Todo{
		Title:    "Update documentation",
		Desc:     "Add examples for new features",
		Complete: false,
		Priority: domain.PriorityLow,
	})

	statsCalc := service.NewStatsCalculator(timeProv)
	todoService := service.NewTodoService(repo, timeProv)
	presenter := ui.NewTodoPresenter()
	calendarAdapter := ui.NewCalendarAdapter(statsCalc)
	viewRenderer := ui.NewViewRenderer()

	// Initialize inputs
	ti := textinput.New()
	ti.Placeholder = "Buy milk..."
	ti.CharLimit = 256
	ti.Width = 56

	ta := textarea.New()
	ta.Placeholder = "Add detailed notes..."
	ta.SetWidth(56)
	ta.SetHeight(8)

	si := textinput.New()
	si.Placeholder = "Search todos..."
	si.CharLimit = 100
	si.Width = 38

	mdRenderer := ui.NewMarkdownRenderer(40)

	// Create calendar and set to fixed date
	cal := calendar.New()
	cal.SetCursor(fixedTime)

	m := &model{
		todoService:      todoService,
		statsCalc:        statsCalc,
		presenter:        presenter,
		calendarAdapter:  calendarAdapter,
		viewRenderer:     viewRenderer,
		calendar:         cal,
		todo:             todo.New(),
		titleInput:       ti,
		descInput:        ta,
		searchInput:      si,
		state:            ui.StateViewing,
		focus:            ui.FocusCalendar,
		editingIndex:     -1,
		editFocus:        ui.FocusTitle,
		markdownRenderer: mdRenderer,
		previewEnabled:   true,
	}

	return m
}

// readOutput reads all bytes from an io.Reader
func readOutput(t *testing.T, r io.Reader) []byte {
	t.Helper()
	data, err := io.ReadAll(r)
	if err != nil {
		t.Fatalf("failed to read output: %v", err)
	}
	return data
}

func TestMainView_MonthView(t *testing.T) {
	m := newTestModel(t)

	tm := teatest.NewTestModel(t, m,
		teatest.WithInitialTermSize(100, 30),
	)

	// Wait for render
	time.Sleep(200 * time.Millisecond)

	// Send quit to end the program
	tm.Send(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'q'}})

	// Get final output
	out := readOutput(t, tm.FinalOutput(t, teatest.WithFinalTimeout(5*time.Second)))

	// Compare against golden file
	teatest.RequireEqualOutput(t, out)
}

func TestMainView_TodoFocus(t *testing.T) {
	m := newTestModel(t)

	tm := teatest.NewTestModel(t, m,
		teatest.WithInitialTermSize(100, 30),
	)

	// Switch focus to todo panel
	tm.Send(tea.KeyMsg{Type: tea.KeyTab})
	time.Sleep(100 * time.Millisecond)

	// Send quit
	tm.Send(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'q'}})

	out := readOutput(t, tm.FinalOutput(t, teatest.WithFinalTimeout(5*time.Second)))
	teatest.RequireEqualOutput(t, out)
}

func TestMainView_WeekView(t *testing.T) {
	m := newTestModel(t)

	tm := teatest.NewTestModel(t, m,
		teatest.WithInitialTermSize(100, 30),
	)

	// Toggle to week view
	tm.Send(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'w'}})
	time.Sleep(100 * time.Millisecond)

	// Send quit
	tm.Send(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'q'}})

	out := readOutput(t, tm.FinalOutput(t, teatest.WithFinalTimeout(5*time.Second)))
	teatest.RequireEqualOutput(t, out)
}

func TestMainView_Navigation(t *testing.T) {
	m := newTestModel(t)

	tm := teatest.NewTestModel(t, m,
		teatest.WithInitialTermSize(100, 30),
	)

	// Navigate down in calendar
	tm.Send(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'j'}})
	tm.Send(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'j'}})
	time.Sleep(100 * time.Millisecond)

	// Send quit
	tm.Send(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'q'}})

	out := readOutput(t, tm.FinalOutput(t, teatest.WithFinalTimeout(5*time.Second)))
	teatest.RequireEqualOutput(t, out)
}

func TestMainView_DayView(t *testing.T) {
	m := newTestModel(t)

	tm := teatest.NewTestModel(t, m,
		teatest.WithInitialTermSize(100, 30),
	)

	// Switch to day view
	tm.Send(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'d'}})
	time.Sleep(100 * time.Millisecond)

	// Send quit
	tm.Send(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'q'}})

	out := readOutput(t, tm.FinalOutput(t, teatest.WithFinalTimeout(5*time.Second)))
	teatest.RequireEqualOutput(t, out)
}

func TestMainView_DayViewEmpty(t *testing.T) {
	m := newTestModel(t)

	tm := teatest.NewTestModel(t, m,
		teatest.WithInitialTermSize(100, 30),
	)

	// Navigate to a date without todos then switch to day view
	tm.Send(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'l'}}) // Move to Jan 19
	tm.Send(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'d'}}) // Day view
	time.Sleep(100 * time.Millisecond)

	// Send quit
	tm.Send(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'q'}})

	out := readOutput(t, tm.FinalOutput(t, teatest.WithFinalTimeout(5*time.Second)))
	teatest.RequireEqualOutput(t, out)
}

func TestMainView_HelpModal(t *testing.T) {
	m := newTestModel(t)

	tm := teatest.NewTestModel(t, m,
		teatest.WithInitialTermSize(100, 30),
	)

	// Open help modal with ?
	tm.Send(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'?'}})
	time.Sleep(100 * time.Millisecond)

	// Press Ctrl+C to force quit while help modal is displayed
	// This captures the help modal output
	tm.Send(tea.KeyMsg{Type: tea.KeyCtrlC})

	out := readOutput(t, tm.FinalOutput(t, teatest.WithFinalTimeout(5*time.Second)))
	teatest.RequireEqualOutput(t, out)
}
