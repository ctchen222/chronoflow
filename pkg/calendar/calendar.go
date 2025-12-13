package calendar

import (
	"fmt"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	// Style for the highlighted day in the calendar.
	dayHighlight = lipgloss.NewStyle().
		Background(lipgloss.Color("#7D56F4")).
		Foreground(lipgloss.Color("#FFFFFF")).
		Bold(true)
)

type Model struct {
	selectedDate time.Time
	cursor       time.Time
	width        int
	height       int
}

func New() *Model {
	now := time.Now()
	return &Model{
		selectedDate: now,
		cursor:       now,
	}
}

func (m *Model) SetSize(width, height int) {
	m.width = width
	m.height = height
}

func (m *Model) Init() tea.Cmd {
	return nil
}

func (m *Model) Cursor() time.Time {
	return m.cursor
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "h", "left":
			m.cursor = m.cursor.AddDate(0, 0, -1)
		case "l", "right":
			m.cursor = m.cursor.AddDate(0, 0, 1)
		case "k", "up":
			m.cursor = m.cursor.AddDate(0, 0, -7)
		case "j", "down":
			m.cursor = m.cursor.AddDate(0, 0, 7)
		case "b":
			m.selectedDate = m.selectedDate.AddDate(0, -1, 0)
			m.cursor = m.selectedDate
		case "n":
			m.selectedDate = m.selectedDate.AddDate(0, 1, 0)
			m.cursor = m.selectedDate
		}
	}
	if m.cursor.Month() != m.selectedDate.Month() || m.cursor.Year() != m.selectedDate.Year() {
		m.selectedDate = m.cursor
	}
	return m, nil
}

func (m *Model) View() string {

	if m.width == 0 || m.height == 0 {

		return "loading..."

	}

	var s strings.Builder

	// --- RENDER AND APPEND, WITH METICULOUS NEWLINE MANAGEMENT ---

	// Main Header (single line) -> requires a newline after

	header := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#FAFAFA")).
		Background(lipgloss.Color("#7D56F4")).
		Padding(0, 1).
		Width(m.width).
		Align(lipgloss.Center).
		Render(m.cursor.Format("January 2006"))

	s.WriteString(header)

	s.WriteString("\n")

	// Sub-header (single line) -> requires a newline after

	subHeader := lipgloss.NewStyle().
		Foreground(lipgloss.Color("240")).
		Padding(0, 1).
		Width(m.width).
		Align(lipgloss.Center).
		Render(m.cursor.Format("Monday, Jan 2, 2006"))

	s.WriteString(subHeader)

	s.WriteString("\n")

	// Weekday headers (single line) -> requires a newline after

	weekdays := []string{"Sun", "Mon", "Tue", "Wed", "Thu", "Fri", "Sat"}

	baseCellWidth := m.width / 7

	remainder := m.width % 7

	cellWidths := make([]int, 7)

	for i := 0; i < 7; i++ {

		cellWidths[i] = baseCellWidth

		if remainder > 0 {

			cellWidths[i]++

			remainder--

		}

	}

			weekdayViews := make([]string, 7)

			for i, day := range weekdays {

				style := lipgloss.NewStyle().

					Width(cellWidths[i]).

					Align(lipgloss.Center).

					Foreground(lipgloss.Color("240")).

					BorderBottom(true).

					BorderForeground(lipgloss.Color("238"))

				weekdayViews[i] = style.Render(day)

			}

	weekdayHeader := lipgloss.JoinHorizontal(lipgloss.Top, weekdayViews...)

	s.WriteString(weekdayHeader)

	s.WriteString("\n")

	// Calendar grid (multi-line blocks) -> DO NOT add newlines after

	firstDay := time.Date(m.selectedDate.Year(), m.selectedDate.Month(), 1, 0, 0, 0, 0, time.UTC)

	daysInMonth := time.Date(m.selectedDate.Year(), m.selectedDate.Month()+1, 0, 0, 0, 0, 0, time.UTC).Day()

	day := 1

	for i := 0; i < 6; i++ {

		var rowViews []string

		for j := 0; j < 7; j++ {

			style := lipgloss.NewStyle().
				Width(cellWidths[j]).
				Height(3). // Hardcoded height for stability

				Border(lipgloss.NormalBorder(), false).
				Align(lipgloss.Center, lipgloss.Center)

			style = style.BorderBottom(true)

			if j != 6 {

				style = style.BorderRight(true)

			}

			if (i == 0 && j < int(firstDay.Weekday())) || day > daysInMonth {

				rowViews = append(rowViews, style.Render(""))

			} else {

								if day == m.cursor.Day() && m.selectedDate.Month() == m.cursor.Month() && m.selectedDate.Year() == m.cursor.Year() {

									style = style.Copy().Inherit(dayHighlight)

								}

				rowViews = append(rowViews, style.Render(fmt.Sprintf("%d", day)))

				day++

			}

		}

				s.WriteString(lipgloss.JoinHorizontal(lipgloss.Top, rowViews...))

				// Add a newline after each row to stack them vertically.

				if i < 5 {

					s.WriteString("\n")

				}

			}

	return s.String()

}
