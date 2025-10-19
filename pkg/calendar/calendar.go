package calendar

import (
	"fmt"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	dayHighlight = lipgloss.Color("#888888")
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
		case "n":
			m.selectedDate = m.selectedDate.AddDate(0, 1, 0)
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

	// Header
	header := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#FAFAFA")).
		Background(lipgloss.Color("#7D56F4")).
		Padding(0, 1).
		Width(m.width).
		Render(m.cursor.Format("January 2006"))

	// Weekday headers
	weekdays := []string{"Sun", "Mon", "Tue", "Wed", "Thu", "Fri", "Sat"}

	// Distribute width
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
			Foreground(lipgloss.Color("240")).
			Width(cellWidths[i]).
			Align(lipgloss.Center)
		weekdayViews[i] = style.Render(day)
	}
	weekdayHeader := lipgloss.JoinHorizontal(lipgloss.Top, weekdayViews...)

	// Calendar grid
	firstDay := time.Date(m.selectedDate.Year(), m.selectedDate.Month(), 1, 0, 0, 0, 0, time.UTC)
	daysInMonth := time.Date(m.selectedDate.Year(), m.selectedDate.Month()+1, 0, 0, 0, 0, 0, time.UTC).Day()

	// Distribute height
	baseCellHeight := (m.height - 2) / 6 // -2 for header and weekday header
	heightRemainder := (m.height - 2) % 6
	cellHeights := make([]int, 6)
	for i := 0; i < 6; i++ {
		cellHeights[i] = baseCellHeight
		if heightRemainder > 0 {
			cellHeights[i]++
			heightRemainder--
		}
	}

	var allRows []string
	day := 1
	for i := 0; i < 6; i++ {
		var rowViews []string
		for j := 0; j < 7; j++ {
			style := lipgloss.NewStyle().
				Width(cellWidths[j]).
				Height(cellHeights[i]).
				Border(lipgloss.NormalBorder(), false).
				Align(lipgloss.Center)
			style = style.BorderBottom(true)
			if j != 6 {
				style = style.BorderRight(true)
			}

			if (i == 0 && j < int(firstDay.Weekday())) || day > daysInMonth {
				rowViews = append(rowViews, style.Render(""))
			} else {
				if day == m.cursor.Day() && m.selectedDate.Month() == m.cursor.Month() && m.selectedDate.Year() == m.cursor.Year() {
					style = style.Background(dayHighlight)
				}
				rowViews = append(rowViews, style.Render(fmt.Sprintf("%d", day)))
				day++
			}
		}
		allRows = append(allRows, lipgloss.JoinHorizontal(lipgloss.Top, rowViews...))
		if day > daysInMonth {
			break
		}
	}

	return lipgloss.JoinVertical(lipgloss.Left, header, weekdayHeader, lipgloss.JoinVertical(lipgloss.Left, allRows...))
}
