package main

import (
	"fmt"
	"os"

	"ctchen222/chronoflow/pkg/calendar"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type model struct {
	calendar *calendar.Model
	width    int
	height   int
}

func (m *model) Init() tea.Cmd {
	return m.calendar.Init()
}

func (m *model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		calendarWidth := int(float64(m.width) * 0.7)
		m.calendar.SetSize(calendarWidth, m.height)
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		}
	}
	newCalendar, cmd := m.calendar.Update(msg)
	m.calendar = newCalendar.(*calendar.Model)
	return m, cmd
}

func (m *model) View() string {
	// Re-calculate widths here to ensure they are always in sync.
	calendarWidth := int(float64(m.width) * 0.7)
	rightPanelWidth := m.width - calendarWidth

	// The calendar view string, which should have a width of `calendarWidth`.
	calendarView := m.calendar.View()

	// --- Right Panel ---
	selectedDateStr := m.calendar.Cursor().Format("2006-01-02 Monday")
	rightPanelContent := lipgloss.NewStyle().
		Bold(true).
		Padding(1, 2).
		Render("Selected Date:\n" + selectedDateStr)

	// Use a styled panel with a rounded border and color to make it distinct.
	rightPanel := lipgloss.NewStyle().
		Width(rightPanelWidth).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("238")).
		Padding(1, 2).
		Render(rightPanelContent)

	// --- Join horizontally ---
	return lipgloss.JoinHorizontal(lipgloss.Top,
		calendarView,
		rightPanel,
	)
}

func main() {
	m := &model{
		calendar: calendar.New(),
	}
	p := tea.NewProgram(m, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
