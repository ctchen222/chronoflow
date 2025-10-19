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
	calendarWidth := int(float64(m.width) * 0.7)
	rightPanelWidth := m.width - calendarWidth

	rightPanel := lipgloss.NewStyle().
		Width(rightPanelWidth).
		Height(m.height).
		Border(lipgloss.NormalBorder(), true).
		Render("Future Component")

	return lipgloss.JoinHorizontal(lipgloss.Top, m.calendar.View(), rightPanel)
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
