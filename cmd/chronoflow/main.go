package main

import (
	"fmt"
	"os"

	"ctchen222/chronoflow/pkg/calendar"
	tea "github.com/charmbracelet/bubbletea"
)

type model struct {
	calendar *calendar.Model
}

func (m *model) Init() tea.Cmd {
	return m.calendar.Init()
}

func (m *model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.calendar.SetSize(msg.Width, msg.Height)
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
	return m.calendar.View()
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