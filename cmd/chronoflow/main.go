package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"ctchen222/chronoflow/pkg/calendar"
	"ctchen222/chronoflow/pkg/todo"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const appName = "chronoflow"

// getDataDir returns the application data directory path
func getDataDir() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "."
	}
	return filepath.Join(homeDir, "."+appName)
}

// getDataFilePath returns the full path to the todos.json file
func getDataFilePath() string {
	return filepath.Join(getDataDir(), "todos.json")
}

// Priority levels
const (
	PriorityNone   = 0
	PriorityLow    = 1
	PriorityMedium = 2
	PriorityHigh   = 3
)

// item represents a to-do item in our list.
type item struct {
	ItemTitle    string `json:"title"`
	ItemDesc     string `json:"desc"`
	ItemComplete bool   `json:"completed"`
	ItemPriority int    `json:"priority"` // 0=none, 1=low, 2=medium, 3=high
	isOverdue    bool   // not saved, calculated at runtime
}

func (i item) priorityIcon() string {
	switch i.ItemPriority {
	case PriorityHigh:
		return "!!!"
	case PriorityMedium:
		return "!!"
	case PriorityLow:
		return "!"
	default:
		return ""
	}
}

func (i item) priorityStyle() lipgloss.Style {
	switch i.ItemPriority {
	case PriorityHigh:
		return lipgloss.NewStyle().Foreground(lipgloss.Color("#FF6B6B")).Bold(true)
	case PriorityMedium:
		return lipgloss.NewStyle().Foreground(lipgloss.Color("#FFB86C"))
	case PriorityLow:
		return lipgloss.NewStyle().Foreground(lipgloss.Color("#8BE9FD"))
	default:
		return lipgloss.NewStyle()
	}
}

func (i item) Title() string {
	checkbox := "☐ "
	style := lipgloss.NewStyle()

	if i.ItemComplete {
		checkbox = "☑ "
		style = style.Foreground(lipgloss.Color("#666")).Strikethrough(true)
	} else if i.isOverdue {
		checkbox = "⚠ "
		style = style.Foreground(lipgloss.Color("#FF6B6B")).Bold(true)
	} else if i.ItemPriority > 0 {
		style = i.priorityStyle()
	}

	title := style.Render(i.ItemTitle)
	if i.ItemPriority > 0 && !i.ItemComplete {
		title = i.priorityStyle().Render(i.priorityIcon()) + " " + title
	}

	return checkbox + title
}

func (i item) Description() string {
	if i.isOverdue && !i.ItemComplete {
		return lipgloss.NewStyle().Foreground(lipgloss.Color("#FF6B6B")).Render("Overdue! " + i.ItemDesc)
	}
	return i.ItemDesc
}

func (i item) FilterValue() string { return i.ItemTitle }

// todos is our in-memory data store.
var todos = make(map[string][]item)

// AppState defines the current state of the application.
type appState int

const (
	viewing appState = iota
	editing
	confirmingDelete
	searching
)

// appFocus determines which panel is currently focused.
type appFocus int

const (
	calendarFocus appFocus = iota
	todoFocus
)

// editFocus determines which input is focused in the editing view.
type editFocus int

const (
	titleFocus editFocus = iota
	descFocus
)

// searchResult represents a todo item found in search
type searchResult struct {
	dateKey string
	index   int
	item    item
}

type model struct {
	calendar        *calendar.Model
	todo            todo.Model
	titleInput      textinput.Model
	descInput       textarea.Model
	searchInput     textinput.Model
	state           appState
	focus           appFocus
	editFocus       editFocus // Focus for editing view
	editingIndex    int
	editingPriority int    // Priority for the item being edited
	deletingIndex   int    // Index of item being deleted
	deletingTitle   string // Title of item being deleted (for display)
	searchResults   []searchResult
	searchIndex     int // Currently selected search result
	width           int
	height          int
}

// updateTodos sets the items for the todo list based on the selected date.
func (m *model) updateTodos() {
	dateKey := m.calendar.Cursor().Format("2006-01-02")
	cursorDate := m.calendar.Cursor()

	// Check if this date is in the past
	now := time.Now()
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.Local)
	selectedDate := time.Date(cursorDate.Year(), cursorDate.Month(), cursorDate.Day(), 0, 0, 0, 0, time.Local)
	isPastDate := selectedDate.Before(today)

	items := []list.Item{}
	if foundTodos, ok := todos[dateKey]; ok {
		for _, td := range foundTodos {
			// Mark as overdue if it's a past date and incomplete
			td.isOverdue = isPastDate && !td.ItemComplete
			items = append(items, td)
		}
	}
	m.todo.SetItems(items)
	m.todo.SetTitle(fmt.Sprintf("To-Do on %s", m.calendar.Cursor().Format("2006-01-02")))

	// Calculate and set statistics
	stats := m.calculateStats()
	m.todo.SetStats(stats)

	// Update calendar with dates that have todos
	m.syncCalendarTodos()
}

// calculateStats calculates todo statistics based on the current view mode
func (m *model) calculateStats() todo.Stats {
	now := time.Now()
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.Local)
	cursorDate := m.calendar.Cursor()

	// Determine the period range based on view mode
	var periodStart, periodEnd time.Time
	var periodLabel string

	if m.calendar.GetViewMode() == calendar.WeekView {
		// Week view: calculate week boundaries
		weekday := int(cursorDate.Weekday())
		periodStart = time.Date(cursorDate.Year(), cursorDate.Month(), cursorDate.Day()-weekday, 0, 0, 0, 0, time.Local)
		periodEnd = periodStart.AddDate(0, 0, 7)
		periodLabel = "This Week"
	} else {
		// Month view: calculate month boundaries
		periodStart = time.Date(cursorDate.Year(), cursorDate.Month(), 1, 0, 0, 0, 0, time.Local)
		periodEnd = periodStart.AddDate(0, 1, 0)
		periodLabel = "This Month"
	}

	stats := todo.Stats{
		PeriodLabel: periodLabel,
	}

	for dateKey, items := range todos {
		date, err := time.Parse("2006-01-02", dateKey)
		if err != nil {
			continue
		}
		// Convert to local time for comparison
		date = time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, time.Local)
		isPast := date.Before(today)
		inPeriod := !date.Before(periodStart) && date.Before(periodEnd)

		for _, item := range items {
			stats.TotalAll++
			if item.ItemComplete {
				stats.CompletedAll++
			} else if isPast {
				stats.OverdueAll++
			}

			// Period stats
			if inPeriod {
				stats.TotalPeriod++
				if item.ItemComplete {
					stats.CompletedPeriod++
				} else if isPast {
					stats.OverduePeriod++
				}
			}
		}
	}
	return stats
}

// syncCalendarTodos updates the calendar with todo status for each date
func (m *model) syncCalendarTodos() {
	todoStatus := make(map[string]calendar.TodoStatus)
	now := time.Now()
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.Local)

	for dateKey, items := range todos {
		if len(items) == 0 {
			continue
		}

		// Parse the date
		date, err := time.Parse("2006-01-02", dateKey)
		if err != nil {
			continue
		}

		// Check if all todos are complete and if there are incomplete ones
		allComplete := true
		hasIncomplete := false
		for _, item := range items {
			if !item.ItemComplete {
				allComplete = false
				hasIncomplete = true
				break
			}
		}

		// Determine if overdue (past date with incomplete todos)
		isOverdue := date.Before(today) && hasIncomplete

		// Convert items to calendar.TodoItem for week view display
		calendarItems := make([]calendar.TodoItem, len(items))
		for i, it := range items {
			calendarItems[i] = calendar.TodoItem{
				Title:    it.ItemTitle,
				Complete: it.ItemComplete,
				Priority: it.ItemPriority,
			}
		}

		todoStatus[dateKey] = calendar.TodoStatus{
			HasTodos:    true,
			HasOverdue:  isOverdue,
			AllComplete: allComplete,
			Count:       len(items),
			Items:       calendarItems,
		}
	}
	m.calendar.SetTodoStatus(todoStatus)
}

// performSearch searches all todos for the given query
func (m *model) performSearch(query string) {
	m.searchResults = nil
	m.searchIndex = 0

	if query == "" {
		return
	}

	query = strings.ToLower(query)

	// Collect all dates and sort them
	var dates []string
	for dateKey := range todos {
		dates = append(dates, dateKey)
	}
	sort.Strings(dates)

	// Search through all todos
	for _, dateKey := range dates {
		items := todos[dateKey]
		for idx, it := range items {
			// Search in title and description
			if strings.Contains(strings.ToLower(it.ItemTitle), query) ||
				strings.Contains(strings.ToLower(it.ItemDesc), query) {
				m.searchResults = append(m.searchResults, searchResult{
					dateKey: dateKey,
					index:   idx,
					item:    it,
				})
			}
		}
	}
}

// jumpToSearchResult navigates to the selected search result
func (m *model) jumpToSearchResult() {
	if len(m.searchResults) == 0 || m.searchIndex >= len(m.searchResults) {
		return
	}

	result := m.searchResults[m.searchIndex]
	date, err := time.Parse("2006-01-02", result.dateKey)
	if err != nil {
		return
	}

	// Navigate calendar to the date
	m.calendar.SetCursor(date)
	m.updateTodos()
	m.focus = todoFocus
}

func (m *model) Init() tea.Cmd {
	m.syncCalendarTodos()
	m.updateTodos()
	return m.calendar.Init()
}

func (m *model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

		// Calculate panel dimensions (must match View calculations)
		panelHeight := m.height - 1 // reserve 1 line for help bar
		calendarWidth := int(float64(m.width) * 0.7)
		todoWidth := m.width - calendarWidth

		// Inner dimensions (subtract 2 for border on each side)
		calInnerW := calendarWidth - 2
		calInnerH := panelHeight - 2
		todoInnerW := todoWidth - 2
		todoInnerH := panelHeight - 2

		m.calendar.SetSize(calInnerW, calInnerH)
		m.todo.SetSize(todoInnerW, todoInnerH)

		// Also resize edit form inputs
		inputWidth := m.width / 3
		m.titleInput.Width = inputWidth
		m.descInput.SetWidth(inputWidth)
		return m, nil

	case tea.KeyMsg:
		switch m.state {
		case viewing:
			switch msg.String() {
			case "q", "ctrl+c":
				saveTodos()
				return m, tea.Quit
			case "a":
				m.state = editing
				m.editingIndex = -1
				m.editingPriority = PriorityNone
				m.titleInput.Reset()
				m.descInput.Reset()
				m.editFocus = titleFocus
				return m, m.titleInput.Focus()
			case "tab":
				if m.focus == calendarFocus {
					m.focus = todoFocus
				} else {
					m.focus = calendarFocus
				}
			case "/":
				m.state = searching
				m.searchInput.Reset()
				m.searchResults = nil
				m.searchIndex = 0
				return m, m.searchInput.Focus()
			case "w":
				// Toggle week/month view and refresh stats
				if m.focus == calendarFocus {
					m.calendar.ToggleViewMode()
					m.updateTodos()
					return m, nil
				}
			}

			switch m.focus {
			case calendarFocus:
				if msg.String() == "enter" {
					m.focus = todoFocus
				}
			case todoFocus:
				switch msg.String() {
				case "e", "enter":
					selected := m.todo.SelectedItem()
					if selected != nil {
						m.state = editing
						m.editingIndex = m.todo.ListIndex()
						selectedItem := selected.(item)
						m.titleInput.SetValue(selectedItem.ItemTitle)
						m.descInput.SetValue(selectedItem.ItemDesc)
						m.editingPriority = selectedItem.ItemPriority
						m.editFocus = titleFocus
						return m, m.titleInput.Focus()
					}
				case "esc":
					m.focus = calendarFocus
				case " ", "x":
					// Toggle completion status
					selected := m.todo.SelectedItem()
					if selected != nil {
						dateKey := m.calendar.Cursor().Format("2006-01-02")
						idx := m.todo.ListIndex()
						if items, ok := todos[dateKey]; ok && idx < len(items) {
							items[idx].ItemComplete = !items[idx].ItemComplete
							todos[dateKey] = items
							m.updateTodos()
						}
					}
				case "d", "backspace":
					selectedItem := m.todo.SelectedItem()
					if selectedItem != nil {
						m.state = confirmingDelete
						m.deletingIndex = m.todo.ListIndex()
						m.deletingTitle = selectedItem.(item).ItemTitle
					}
				case "1", "2", "3", "0":
					// Quick priority change
					selected := m.todo.SelectedItem()
					if selected != nil {
						dateKey := m.calendar.Cursor().Format("2006-01-02")
						idx := m.todo.ListIndex()
						if items, ok := todos[dateKey]; ok && idx < len(items) {
							switch msg.String() {
							case "0":
								items[idx].ItemPriority = PriorityNone
							case "1":
								items[idx].ItemPriority = PriorityLow
							case "2":
								items[idx].ItemPriority = PriorityMedium
							case "3":
								items[idx].ItemPriority = PriorityHigh
							}
							todos[dateKey] = items
							m.updateTodos()
						}
					}
				case "J", "K":
					// Move todo up/down (reorder)
					selected := m.todo.SelectedItem()
					if selected != nil {
						dateKey := m.calendar.Cursor().Format("2006-01-02")
						idx := m.todo.ListIndex()
						if items, ok := todos[dateKey]; ok && len(items) > 1 {
							var newIdx int
							if msg.String() == "J" && idx < len(items)-1 {
								// Move down
								newIdx = idx + 1
								items[idx], items[newIdx] = items[newIdx], items[idx]
								todos[dateKey] = items
								m.updateTodos()
								// Move cursor down to follow the item
								m.todo, _ = m.todo.Update(tea.KeyMsg{Type: tea.KeyDown})
							} else if msg.String() == "K" && idx > 0 {
								// Move up
								newIdx = idx - 1
								items[idx], items[newIdx] = items[newIdx], items[idx]
								todos[dateKey] = items
								m.updateTodos()
								// Move cursor up to follow the item
								m.todo, _ = m.todo.Update(tea.KeyMsg{Type: tea.KeyUp})
							}
						}
					}
				}
			}

		case confirmingDelete:
			switch msg.String() {
			case "y", "Y", "enter":
				// Confirm delete
				dateKey := m.calendar.Cursor().Format("2006-01-02")
				if items, ok := todos[dateKey]; ok && len(items) > m.deletingIndex {
					todos[dateKey] = append(items[:m.deletingIndex], items[m.deletingIndex+1:]...)
				}
				m.updateTodos()
				m.state = viewing
				return m, nil
			case "n", "N", "esc":
				// Cancel delete
				m.state = viewing
				return m, nil
			}

		case editing:
			switch msg.String() {
			case "esc":
				m.state = viewing
				m.titleInput.Blur()
				m.descInput.Blur()
				return m, nil
			case "ctrl+1":
				m.editingPriority = PriorityLow
				return m, nil
			case "ctrl+2":
				m.editingPriority = PriorityMedium
				return m, nil
			case "ctrl+3":
				m.editingPriority = PriorityHigh
				return m, nil
			case "ctrl+0":
				m.editingPriority = PriorityNone
				return m, nil
			case "tab":
				if m.editFocus == titleFocus {
					m.editFocus = descFocus
					m.titleInput.Blur()
					cmd = m.descInput.Focus()
				} else {
					m.editFocus = titleFocus
					m.descInput.Blur()
					cmd = m.titleInput.Focus()
				}
				cmds = append(cmds, cmd)
				return m, tea.Batch(cmds...)
			case "enter":
				// Only save if enter is pressed on the title input
				if m.editFocus == titleFocus {
					dateKey := m.calendar.Cursor().Format("2006-01-02")
					title := m.titleInput.Value()
					desc := m.descInput.Value()

					if m.editingIndex == -1 {
						if title != "" {
							todos[dateKey] = append(todos[dateKey], item{
								ItemTitle:    title,
								ItemDesc:     desc,
								ItemPriority: m.editingPriority,
							})
						}
					} else {
						if items, ok := todos[dateKey]; ok && len(items) > m.editingIndex {
							items[m.editingIndex].ItemTitle = title
							items[m.editingIndex].ItemDesc = desc
							items[m.editingIndex].ItemPriority = m.editingPriority
							todos[dateKey] = items
						}
					}
					m.updateTodos()
					m.state = viewing
					m.titleInput.Blur()
					m.descInput.Blur()
					return m, nil
				}
			}

		case searching:
			switch msg.String() {
			case "esc":
				m.state = viewing
				m.searchInput.Blur()
				return m, nil
			case "enter":
				// Jump to selected result
				if len(m.searchResults) > 0 {
					m.jumpToSearchResult()
				}
				m.state = viewing
				m.searchInput.Blur()
				return m, nil
			case "up", "ctrl+p":
				// Previous result
				if m.searchIndex > 0 {
					m.searchIndex--
				}
				return m, nil
			case "down", "ctrl+n":
				// Next result
				if m.searchIndex < len(m.searchResults)-1 {
					m.searchIndex++
				}
				return m, nil
			default:
				// Update search input and perform search
				m.searchInput, cmd = m.searchInput.Update(msg)
				m.performSearch(m.searchInput.Value())
				return m, cmd
			}
		}
	}

	// --- Pass messages to focused components ---
	prevCursor := m.calendar.Cursor()
	switch m.state {
	case viewing:
		if m.focus == calendarFocus {
			var newCal tea.Model
			newCal, cmd = m.calendar.Update(msg)
			m.calendar = newCal.(*calendar.Model)
		} else {
			m.todo, cmd = m.todo.Update(msg)
		}
	case editing:
		if m.editFocus == titleFocus {
			m.titleInput, cmd = m.titleInput.Update(msg)
		} else {
			m.descInput, cmd = m.descInput.Update(msg)
		}
	case searching:
		m.searchInput, cmd = m.searchInput.Update(msg)
		m.performSearch(m.searchInput.Value())
	}
	cmds = append(cmds, cmd)

	if m.state == viewing && m.focus == calendarFocus && !prevCursor.Equal(m.calendar.Cursor()) {
		m.updateTodos()
	}

	return m, tea.Batch(cmds...)
}

func (m *model) helpBar() string {
	keyStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#7D56F4")).
		Bold(true)
	descStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#666"))
	sepStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#444"))

	sep := sepStyle.Render(" │ ")

	var keys string
	switch m.state {
	case viewing:
		if m.focus == calendarFocus {
			keys = keyStyle.Render("h/j/k/l") + descStyle.Render(" nav") + sep +
				keyStyle.Render("b/n") + descStyle.Render(" month") + sep +
				keyStyle.Render("w") + descStyle.Render(" week") + sep +
				keyStyle.Render("t") + descStyle.Render(" today") + sep +
				keyStyle.Render("/") + descStyle.Render(" search") + sep +
				keyStyle.Render("Tab") + descStyle.Render(" todos") + sep +
				keyStyle.Render("q") + descStyle.Render(" quit")
		} else {
			keys = keyStyle.Render("j/k") + descStyle.Render(" nav") + sep +
				keyStyle.Render("J/K") + descStyle.Render(" move") + sep +
				keyStyle.Render("Space") + descStyle.Render(" done") + sep +
				keyStyle.Render("1/2/3") + descStyle.Render(" priority") + sep +
				keyStyle.Render("/") + descStyle.Render(" search") + sep +
				keyStyle.Render("a") + descStyle.Render(" add") + sep +
				keyStyle.Render("e") + descStyle.Render(" edit") + sep +
				keyStyle.Render("q") + descStyle.Render(" quit")
		}
	case editing:
		keys = keyStyle.Render("Tab") + descStyle.Render(" switch field") + sep +
			keyStyle.Render("Enter") + descStyle.Render(" save") + sep +
			keyStyle.Render("Esc") + descStyle.Render(" cancel")
	case confirmingDelete:
		keys = keyStyle.Render("y/Enter") + descStyle.Render(" confirm") + sep +
			keyStyle.Render("n/Esc") + descStyle.Render(" cancel")
	case searching:
		keys = keyStyle.Render("Up/Down") + descStyle.Render(" navigate") + sep +
			keyStyle.Render("Enter") + descStyle.Render(" go to") + sep +
			keyStyle.Render("Esc") + descStyle.Render(" cancel")
	}

	return lipgloss.NewStyle().
		Width(m.width).
		Align(lipgloss.Center).
		Background(lipgloss.Color("#1a1a1a")).
		Padding(0, 1).
		Render(keys)
}

func (m *model) View() string {
	switch m.state {
	case viewing:
		// Calculate dimensions
		panelHeight := m.height - 1 // reserve 1 line for help bar
		calendarWidth := int(float64(m.width) * 0.7)
		todoWidth := m.width - calendarWidth

		// Inner content dimensions (subtract 2 for border on each side)
		calInnerW := calendarWidth - 2
		calInnerH := panelHeight - 2
		todoInnerW := todoWidth - 2
		todoInnerH := panelHeight - 2

		// Get content and place it in fixed-size container
		calContent := lipgloss.Place(calInnerW, calInnerH, lipgloss.Left, lipgloss.Top, m.calendar.View())
		todoContent := lipgloss.Place(todoInnerW, todoInnerH, lipgloss.Left, lipgloss.Top, m.todo.View())

		// Border colors based on focus
		calBorderColor := lipgloss.Color("#444")
		todoBorderColor := lipgloss.Color("#444")
		if m.focus == calendarFocus {
			calBorderColor = lipgloss.Color("#7D56F4")
		} else {
			todoBorderColor = lipgloss.Color("#7D56F4")
		}

		// Apply borders
		calView := lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(calBorderColor).
			Render(calContent)
		todoView := lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(todoBorderColor).
			Render(todoContent)

		panels := lipgloss.JoinHorizontal(lipgloss.Top, calView, todoView)
		return lipgloss.JoinVertical(lipgloss.Left, panels, m.helpBar())
	case editing:
		// Determine colors based on mode (new vs edit)
		isNew := m.editingIndex == -1
		var accentColor, headerIcon string
		if isNew {
			accentColor = "#50FA7B" // Green for new
			headerIcon = "+"
		} else {
			accentColor = "#8BE9FD" // Cyan for edit
			headerIcon = "~"
		}

		// Header
		headerText := "New Todo"
		if !isNew {
			headerText = "Edit Todo"
		}
		headerStyle := lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color(accentColor)).
			MarginBottom(1)
		header := headerStyle.Render(headerIcon + "  " + headerText)

		// Date
		dateText := m.calendar.Cursor().Format("Mon, Jan 2, 2006")
		dateStyle := lipgloss.NewStyle().
			Foreground(lipgloss.Color("#888")).
			MarginBottom(1)
		date := dateStyle.Render(dateText)

		// Title input with label
		titleLabelStyle := lipgloss.NewStyle().
			Foreground(lipgloss.Color("#888"))
		titleBorderColor := lipgloss.Color("#444")
		if m.editFocus == titleFocus {
			titleBorderColor = lipgloss.Color(accentColor)
		}
		titleInputStyle := lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(titleBorderColor).
			Padding(0, 1).
			Width(60)
		titleSection := lipgloss.JoinVertical(lipgloss.Left,
			titleLabelStyle.Render("Title"),
			titleInputStyle.Render(m.titleInput.View()),
		)

		// Description input with label
		descLabelStyle := lipgloss.NewStyle().
			Foreground(lipgloss.Color("#888")).
			MarginTop(1)
		descBorderColor := lipgloss.Color("#444")
		if m.editFocus == descFocus {
			descBorderColor = lipgloss.Color(accentColor)
		}
		descInputStyle := lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(descBorderColor).
			Padding(0, 1).
			Width(60)
		descSection := lipgloss.JoinVertical(lipgloss.Left,
			descLabelStyle.Render("Description (optional)"),
			descInputStyle.Render(m.descInput.View()),
		)

		// Priority selector
		priorityLabelStyle := lipgloss.NewStyle().
			Foreground(lipgloss.Color("#888")).
			MarginTop(1)

		priorityOptions := []struct {
			level int
			label string
			color string
		}{
			{PriorityNone, "None", "#666"},
			{PriorityLow, "Low", "#8BE9FD"},
			{PriorityMedium, "Medium", "#FFB86C"},
			{PriorityHigh, "High", "#FF6B6B"},
		}

		var priorityItems []string
		for _, opt := range priorityOptions {
			style := lipgloss.NewStyle().Foreground(lipgloss.Color(opt.color))
			label := opt.label
			if opt.level == m.editingPriority {
				label = "[" + label + "]"
				style = style.Bold(true)
			} else {
				label = " " + label + " "
			}
			priorityItems = append(priorityItems, style.Render(label))
		}

		priorityRow := lipgloss.JoinHorizontal(lipgloss.Center, priorityItems...)
		prioritySection := lipgloss.JoinVertical(lipgloss.Left,
			priorityLabelStyle.Render("Priority (Ctrl+0/1/2/3)"),
			priorityRow,
		)

		// Combine all modal content
		modalContent := lipgloss.JoinVertical(lipgloss.Left,
			header,
			date,
			titleSection,
			descSection,
			prioritySection,
		)

		// Modal box with background
		modalBox := lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color(accentColor)).
			Padding(1, 2).
			Render(modalContent)

		// Center the modal in the available space
		bgHeight := m.height - 1 // exclude help bar
		centered := lipgloss.Place(m.width, bgHeight,
			lipgloss.Center, lipgloss.Center,
			modalBox,
			lipgloss.WithWhitespaceChars(" "),
			lipgloss.WithWhitespaceForeground(lipgloss.Color("#333")))

		return lipgloss.JoinVertical(lipgloss.Left, centered, m.helpBar())

	case confirmingDelete:
		// Truncate title if too long
		title := m.deletingTitle
		if len(title) > 35 {
			title = title[:32] + "..."
		}

		accentColor := "#FF6B6B" // Red for delete

		// Header
		headerStyle := lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color(accentColor)).
			MarginBottom(1)
		header := headerStyle.Render("x  Delete Todo?")

		// Todo title being deleted
		titleStyle := lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFFFFF")).
			Bold(true).
			Padding(1, 0)
		todoTitle := titleStyle.Render("\"" + title + "\"")

		// Warning message
		warningStyle := lipgloss.NewStyle().
			Foreground(lipgloss.Color("#888")).
			Italic(true)
		warning := warningStyle.Render("This action cannot be undone.")

		// Modal content
		modalContent := lipgloss.JoinVertical(lipgloss.Center,
			header,
			todoTitle,
			warning,
		)

		// Modal box
		modalBox := lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color(accentColor)).
			Padding(1, 3).
			Render(modalContent)

		// Center the modal
		bgHeight := m.height - 1
		centered := lipgloss.Place(m.width, bgHeight,
			lipgloss.Center, lipgloss.Center,
			modalBox,
			lipgloss.WithWhitespaceChars(" "),
			lipgloss.WithWhitespaceForeground(lipgloss.Color("#333")))

		return lipgloss.JoinVertical(lipgloss.Left, centered, m.helpBar())

	case searching:
		accentColor := "#FFB86C" // Orange for search

		// Header
		headerStyle := lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color(accentColor)).
			MarginBottom(1)
		header := headerStyle.Render("/  Search Todos")

		// Search input
		inputLabelStyle := lipgloss.NewStyle().
			Foreground(lipgloss.Color("#888"))
		inputBorderColor := lipgloss.Color(accentColor)
		inputStyle := lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(inputBorderColor).
			Padding(0, 1).
			Width(40)
		inputSection := lipgloss.JoinVertical(lipgloss.Left,
			inputLabelStyle.Render("Search query"),
			inputStyle.Render(m.searchInput.View()),
		)

		// Results
		var resultsContent string
		if len(m.searchResults) == 0 {
			if m.searchInput.Value() == "" {
				resultsContent = lipgloss.NewStyle().
					Foreground(lipgloss.Color("#666")).
					Italic(true).
					Render("Type to search...")
			} else {
				resultsContent = lipgloss.NewStyle().
					Foreground(lipgloss.Color("#666")).
					Italic(true).
					Render("No results found")
			}
		} else {
			var resultLines []string
			maxResults := 8 // Show max 8 results
			start := 0
			if m.searchIndex >= maxResults {
				start = m.searchIndex - maxResults + 1
			}
			end := start + maxResults
			if end > len(m.searchResults) {
				end = len(m.searchResults)
			}

			for i := start; i < end; i++ {
				r := m.searchResults[i]
				dateStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#7D56F4"))
				titleStyle := lipgloss.NewStyle()

				prefix := "  "
				if i == m.searchIndex {
					prefix = "> "
					titleStyle = titleStyle.Bold(true).Foreground(lipgloss.Color(accentColor))
				}

				// Show completion status
				status := "☐"
				if r.item.ItemComplete {
					status = "☑"
					titleStyle = titleStyle.Foreground(lipgloss.Color("#666"))
				}

				line := prefix + dateStyle.Render(r.dateKey) + " " + status + " " + titleStyle.Render(r.item.ItemTitle)
				resultLines = append(resultLines, line)
			}

			resultsHeader := lipgloss.NewStyle().
				Foreground(lipgloss.Color("#888")).
				MarginTop(1).
				Render(fmt.Sprintf("Results (%d found)", len(m.searchResults)))

			resultsContent = lipgloss.JoinVertical(lipgloss.Left,
				append([]string{resultsHeader}, resultLines...)...)
		}

		// Combine all modal content
		modalContent := lipgloss.JoinVertical(lipgloss.Left,
			header,
			inputSection,
			resultsContent,
		)

		// Modal box
		modalBox := lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color(accentColor)).
			Padding(1, 2).
			Render(modalContent)

		// Center the modal
		bgHeight := m.height - 1
		centered := lipgloss.Place(m.width, bgHeight,
			lipgloss.Center, lipgloss.Center,
			modalBox,
			lipgloss.WithWhitespaceChars(" "),
			lipgloss.WithWhitespaceForeground(lipgloss.Color("#333")))

		return lipgloss.JoinVertical(lipgloss.Left, centered, m.helpBar())

	default:
		return "unknown state"
	}
}

// saveTodos saves the current to-do items to the JSON database file.
func saveTodos() {
	// Ensure data directory exists
	if err := os.MkdirAll(getDataDir(), 0755); err != nil {
		fmt.Printf("Error creating data directory: %v\n", err)
		return
	}

	data, err := json.MarshalIndent(todos, "", "  ")
	if err != nil {
		fmt.Printf("Error marshaling todos to JSON: %v\n", err)
		return
	}

	err = os.WriteFile(getDataFilePath(), data, 0644)
	if err != nil {
		fmt.Printf("Error writing to db file: %v\n", err)
	}
}

// loadTodos loads the to-do items from the JSON database file.
func loadTodos() {
	data, err := os.ReadFile(getDataFilePath())
	if err != nil {
		if os.IsNotExist(err) {
			return // No file, nothing to load
		}
		fmt.Printf("Error reading db file: %v\n", err)
		return
	}
	if len(data) == 0 {
		return // Empty file
	}
	err = json.Unmarshal(data, &todos)
	if err != nil {
		fmt.Printf("Error parsing db file: %v\n", err)
		todos = make(map[string][]item) // Start fresh if file is corrupt
	}
}

func main() {
	loadTodos() // Load todos on startup

	// Initialize Title Input
	ti := textinput.New()
	ti.Placeholder = "Buy milk..."
	ti.Focus()
	ti.CharLimit = 256
	ti.Width = 56

	// Initialize Description Input
	ta := textarea.New()
	ta.Placeholder = "Add detailed notes, links, or any information..."
	ta.KeyMap.InsertNewline.SetEnabled(true) // Allow newlines in description
	ta.SetWidth(56)
	ta.SetHeight(8)

	// Initialize Search Input
	si := textinput.New()
	si.Placeholder = "Search todos..."
	si.CharLimit = 100
	si.Width = 38

	m := &model{
		calendar:     calendar.New(),
		todo:         todo.New(),
		titleInput:   ti,
		descInput:    ta,
		searchInput:  si,
		state:        viewing,
		focus:        calendarFocus,
		editingIndex: -1,
		editFocus:    titleFocus,
	}
	p := tea.NewProgram(m, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
