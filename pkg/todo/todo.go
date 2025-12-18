package todo

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	docStyle = lipgloss.NewStyle().Padding(0, 1)
)

// Stats holds statistics about todos
type Stats struct {
	TotalAll       int    // total todos across all dates
	CompletedAll   int    // completed todos across all dates
	OverdueAll     int    // overdue todos across all dates
	TotalPeriod    int    // todos for current period (week/month)
	CompletedPeriod int   // completed for current period
	PeriodLabel    string // "This Week" or "This Month"
}

type Model struct {
	list  list.Model
	stats Stats
	title string // Store title separately to avoid mutation issues
	width int    // Store width for progress bar rendering
}

func New() Model {
	l := list.New([]list.Item{}, list.NewDefaultDelegate(), 0, 0)
	l.Title = "" // Disable built-in title, we render our own
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false)
	l.SetShowHelp(false) // We have our own help bar

	// Clear title bar styles to prevent blue square artifact
	l.Styles.TitleBar = lipgloss.NewStyle()
	l.Styles.Title = lipgloss.NewStyle()

	return Model{list: l, title: "To-Do List"}
}

// SetStats updates the statistics
func (m *Model) SetStats(stats Stats) {
	m.stats = stats
}

// SetShowHelp controls whether the built-in help is shown.
func (m *Model) SetShowHelp(show bool) {
	m.list.SetShowHelp(show)
}

func (m *Model) SetSize(w, h int) {
	m.width = w
	// Account for docStyle padding (1 char left + 1 char right = 2)
	// Reserve space for title (1 line) + stats (2 lines) + progress bar (1 line) + margin
	listHeight := h - 5
	if listHeight < 1 {
		listHeight = 1
	}
	m.list.SetSize(w-2, listHeight)
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m Model) View() string {
	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#FAFAFA"))

	// Render our custom title
	content := titleStyle.Render(m.title)

	// Build stats and progress bar based on current period
	if m.stats.TotalPeriod > 0 {
		// Calculate completion percentage for the period
		pct := m.stats.CompletedPeriod * 100 / m.stats.TotalPeriod

		// Progress bar
		barWidth := m.width - 6 // Account for padding and brackets
		if barWidth < 10 {
			barWidth = 10
		}
		filledWidth := barWidth * m.stats.CompletedPeriod / m.stats.TotalPeriod
		emptyWidth := barWidth - filledWidth

		filledStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#50FA7B"))
		emptyStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#444"))
		pctStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#888"))

		progressBar := filledStyle.Render(strings.Repeat("█", filledWidth)) +
			emptyStyle.Render(strings.Repeat("░", emptyWidth))
		progressLine := fmt.Sprintf("[%s] %s", progressBar, pctStyle.Render(fmt.Sprintf("%d%%", pct)))

		// Stats text with period label
		statsStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#888"))
		completedStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#50FA7B"))
		overdueStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#FF6B6B"))
		labelStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#7D56F4"))

		statsText := labelStyle.Render(m.stats.PeriodLabel) + ": " +
			completedStyle.Render(fmt.Sprintf("%d/%d", m.stats.CompletedPeriod, m.stats.TotalPeriod))
		if m.stats.OverdueAll > 0 {
			statsText += "  " + overdueStyle.Render(fmt.Sprintf("Overdue: %d", m.stats.OverdueAll))
		}

		content = lipgloss.JoinVertical(lipgloss.Left,
			content,
			progressLine,
			statsStyle.Render(statsText),
		)
	}

	if len(m.list.Items()) == 0 {
		// Better empty state with icon and helpful tips
		iconStyle := lipgloss.NewStyle().
			Foreground(lipgloss.Color("#7D56F4")).
			Bold(true)
		msgStyle := lipgloss.NewStyle().
			Foreground(lipgloss.Color("#888"))
		tipStyle := lipgloss.NewStyle().
			Foreground(lipgloss.Color("#666")).
			Italic(true)

		emptyContent := lipgloss.JoinVertical(lipgloss.Center,
			iconStyle.Render("[ ]"),
			"",
			msgStyle.Render("No tasks for this day"),
			"",
			tipStyle.Render("Press 'a' to add a new task"),
			tipStyle.Render("Use h/l to browse other days"),
		)

		content = lipgloss.JoinVertical(lipgloss.Left, content, "", emptyContent)
		return docStyle.Render(content)
	}

	// For non-empty list, append the list view (title is already disabled)
	listView := m.list.View()
	content = lipgloss.JoinVertical(lipgloss.Left, content, listView)
	return docStyle.Render(content)
}

func (m *Model) SetTitle(title string) {
	m.title = title
}

// SetItems replaces the current list of items with a new one.
func (m *Model) SetItems(items []list.Item) {
	m.list.SetItems(items)
}

// SelectedItem returns the currently selected item in the list.
func (m Model) SelectedItem() list.Item {
	return m.list.SelectedItem()
}

// ListIndex returns the index of the currently selected item.
func (m Model) ListIndex() int {
	return m.list.Index()
}
