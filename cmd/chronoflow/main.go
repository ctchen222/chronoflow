package main

import (
	"fmt"
	"os"
	"path/filepath"
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
	"github.com/charmbracelet/lipgloss"
)

const appName = "chronoflow"
const statusMessageDuration = 2 * time.Second

// clearStatusMsg is sent when the status message should be cleared
type clearStatusMsg struct{}

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

type model struct {
	// Services
	todoService     *service.TodoService
	statsCalc       *service.StatsCalculator
	presenter       *ui.TodoPresenter
	calendarAdapter *ui.CalendarAdapter
	viewRenderer    *ui.ViewRenderer

	// UI Components
	calendar    *calendar.Model
	todo        todo.Model
	titleInput  textinput.Model
	descInput   textarea.Model
	searchInput textinput.Model

	// State
	state           ui.AppState
	focus           ui.AppFocus
	editFocus       ui.EditFocus
	editingIndex    int
	editingPriority domain.Priority
	deletingIndex   int
	deletingTitle   string
	searchResults   []service.SearchResult
	searchIndex     int

	// Preview state
	markdownRenderer *ui.MarkdownRenderer
	previewEnabled   bool

	// Status message state
	statusMessage string
	statusType    string // "success", "warning", "info", "priority"

	// Go-to-date state
	dateInput textinput.Model
	dateError string

	// Scheduling state
	scheduleInput     textinput.Model
	schedulingTaskIdx int    // Original index in the full todo list
	scheduleError     string
}

// updateTodos sets the items for the todo list based on the selected date.
func (m *model) updateTodos() {
	cursorDate := m.calendar.Cursor()

	// Get todos with status from service
	todosWithStatus := m.todoService.GetTodosForDate(cursorDate)

	// Convert to list items using presenter
	items := m.presenter.ToListItems(todosWithStatus)
	m.todo.SetItems(items)
	m.todo.SetTitle(fmt.Sprintf("To-Do on %s", cursorDate.Format("2006-01-02")))

	// Calculate and set statistics
	stats := m.calculateStats()
	m.todo.SetStats(stats)

	// Update calendar with dates that have todos
	m.syncCalendarTodos()
}

// calculateStats calculates todo statistics based on the current view mode
func (m *model) calculateStats() todo.Stats {
	viewMode := ui.ConvertViewMode(m.calendar.GetViewMode())
	svcStats := m.statsCalc.CalculateStats(m.todoService.GetAllTodos(), viewMode, m.calendar.Cursor())

	return todo.Stats{
		TotalAll:        svcStats.TotalAll,
		CompletedAll:    svcStats.CompletedAll,
		OverdueAll:      svcStats.OverdueAll,
		TotalPeriod:     svcStats.TotalPeriod,
		CompletedPeriod: svcStats.CompletedPeriod,
		OverduePeriod:   svcStats.OverduePeriod,
		PeriodLabel:     svcStats.PeriodLabel,
	}
}

// syncCalendarTodos updates the calendar with todo status for each date
func (m *model) syncCalendarTodos() {
	todoStatus := m.calendarAdapter.BuildTodoStatus(m.todoService.GetAllTodos())
	m.calendar.SetTodoStatus(todoStatus)
}

// performSearch searches all todos for the given query
func (m *model) performSearch(query string) {
	m.searchResults = m.todoService.Search(query)
	m.searchIndex = 0
}

// jumpToSearchResult navigates to the selected search result
func (m *model) jumpToSearchResult() {
	if len(m.searchResults) == 0 || m.searchIndex >= len(m.searchResults) {
		return
	}

	result := m.searchResults[m.searchIndex]
	date, err := time.Parse("2006-01-02", result.DateKey)
	if err != nil {
		return
	}

	// Navigate calendar to the date
	m.calendar.SetCursor(date)
	m.updateTodos()
	m.focus = ui.FocusTodo
}

// setStatus sets a status message and returns a command to clear it after a delay
func (m *model) setStatus(message, statusType string) tea.Cmd {
	m.statusMessage = message
	m.statusType = statusType
	return tea.Tick(statusMessageDuration, func(t time.Time) tea.Msg {
		return clearStatusMsg{}
	})
}

// parseDate parses a date string in various formats
// Supports: YYYY-MM-DD, MM-DD (current year), DD (current month)
func (m *model) parseDate(dateStr string) (time.Time, error) {
	now := m.calendar.Cursor()

	// Try full date format: YYYY-MM-DD
	if t, err := time.Parse("2006-01-02", dateStr); err == nil {
		return t, nil
	}

	// Try MM-DD format (current year)
	if t, err := time.Parse("01-02", dateStr); err == nil {
		return time.Date(now.Year(), t.Month(), t.Day(), 0, 0, 0, 0, time.Local), nil
	}

	// Try DD format (current month and year)
	if t, err := time.Parse("02", dateStr); err == nil {
		return time.Date(now.Year(), now.Month(), t.Day(), 0, 0, 0, 0, time.Local), nil
	}

	return time.Time{}, fmt.Errorf("invalid date format")
}

// parseTimeInput parses time input string and returns start and end times
// Supports formats: "HH:MM" (default 1 hour duration) or "HH:MM-HH:MM"
func (m *model) parseTimeInput(input string) (string, string, error) {
	// Try range format: HH:MM-HH:MM
	if len(input) == 11 && input[5] == '-' {
		startStr := input[:5]
		endStr := input[6:]
		if _, err := time.Parse("15:04", startStr); err != nil {
			return "", "", fmt.Errorf("invalid start time")
		}
		if _, err := time.Parse("15:04", endStr); err != nil {
			return "", "", fmt.Errorf("invalid end time")
		}
		return startStr, endStr, nil
	}

	// Try single time format: HH:MM (default 1 hour)
	if _, err := time.Parse("15:04", input); err == nil {
		endTime := m.addHourToTime(input)
		return input, endTime, nil
	}

	return "", "", fmt.Errorf("invalid time format")
}

// addHourToTime adds 1 hour to a time string in HH:MM format
func (m *model) addHourToTime(timeStr string) string {
	t, err := time.Parse("15:04", timeStr)
	if err != nil {
		return "00:00"
	}
	return t.Add(time.Hour).Format("15:04")
}

// findOriginalTodoIndex finds the original index of an unscheduled task in the full todo list
func (m *model) findOriginalTodoIndex(unscheduledIdx int) int {
	cursorDate := m.calendar.Cursor()
	todos := m.todoService.GetTodosForDate(cursorDate)

	unscheduledCount := 0
	for i, td := range todos {
		if !td.IsScheduled() {
			if unscheduledCount == unscheduledIdx {
				return i
			}
			unscheduledCount++
		}
	}
	return -1
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
	case clearStatusMsg:
		m.statusMessage = ""
		m.statusType = ""
		return m, nil

	case tea.WindowSizeMsg:
		m.viewRenderer.SetSize(msg.Width, msg.Height)

		// Calculate panel dimensions (must match View calculations)
		panelHeight := msg.Height - 1 // reserve 1 line for help bar
		calendarWidth := int(float64(msg.Width) * 0.7)
		todoWidth := msg.Width - calendarWidth

		// Inner dimensions (subtract 2 for border on each side)
		calInnerW := calendarWidth - 2
		calInnerH := panelHeight - 2
		todoInnerW := todoWidth - 2
		todoInnerH := panelHeight - 2

		m.calendar.SetSize(calInnerW, calInnerH)
		m.todo.SetSize(todoInnerW, todoInnerH)

		// Resize edit form inputs based on responsive modal dimensions
		dims := m.viewRenderer.CalculateModalDimensions(m.previewEnabled)
		m.titleInput.Width = dims.InputWidth
		m.descInput.SetWidth(dims.InputWidth)

		// Update markdown renderer width for preview pane
		if dims.ShowPreview {
			m.markdownRenderer.SetWidth(dims.PreviewWidth - 4) // Account for padding
		}
		return m, nil

	case tea.KeyMsg:
		switch m.state {
		case ui.StateHelp:
			// Ctrl+C quits from help modal
			if msg.String() == "ctrl+c" {
				return m, tea.Quit
			}
			// Any other key closes the help modal
			m.state = ui.StateViewing
			return m, nil

		case ui.StateViewing:
			switch msg.String() {
			case "q", "ctrl+c":
				m.todoService.Persist()
				return m, tea.Quit
			case "?":
				m.state = ui.StateHelp
				return m, nil
			case "a":
				m.state = ui.StateEditing
				m.editingIndex = -1
				m.editingPriority = domain.PriorityNone
				m.titleInput.Reset()
				m.descInput.Reset()
				m.editFocus = ui.FocusTitle
				return m, m.titleInput.Focus()
			case "tab":
				if m.focus == ui.FocusCalendar {
					m.focus = ui.FocusTodo
				} else {
					m.focus = ui.FocusCalendar
				}
			case "shift+tab":
				if m.focus == ui.FocusTodo {
					m.focus = ui.FocusCalendar
				}
			case "/":
				m.state = ui.StateSearching
				m.searchInput.Reset()
				m.searchResults = nil
				m.searchIndex = 0
				return m, m.searchInput.Focus()
			case "w":
				// Toggle week/month view and refresh stats
				if m.focus == ui.FocusCalendar {
					m.calendar.ToggleViewMode()
					m.updateTodos()
					return m, nil
				}
			case "d":
				// Switch to Day View (only when in calendar focus and not in todo panel)
				if m.focus == ui.FocusCalendar {
					m.calendar.SetViewMode(calendar.DayView)
					m.updateTodos()
					return m, nil
				}
			case "m":
				// Switch to Month View
				if m.focus == ui.FocusCalendar {
					m.calendar.SetViewMode(calendar.MonthView)
					m.updateTodos()
					return m, nil
				}
			case "g":
				// Open go-to-date modal
				if m.focus == ui.FocusCalendar {
					m.state = ui.StateGoToDate
					m.dateInput.Reset()
					m.dateError = ""
					return m, m.dateInput.Focus()
				}
			}

			switch m.focus {
			case ui.FocusCalendar:
				// Day View specific navigation
				if m.calendar.GetViewMode() == calendar.DayView {
					// T key toggles between List and Timeline modes
					if msg.String() == "T" {
						m.calendar.ToggleDayViewMode()
						return m, nil
					}

					// Handle based on Day View mode
					if m.calendar.GetDayViewMode() == calendar.DayViewModeList {
						// List Mode navigation
						switch msg.String() {
						case "j", "down":
							dateKey := m.calendar.Cursor().Format("2006-01-02")
							if status, ok := m.calendar.GetTodoStatus()[dateKey]; ok {
								m.calendar.MoveListSelection(1, len(status.Items))
							}
							return m, nil
						case "k", "up":
							dateKey := m.calendar.Cursor().Format("2006-01-02")
							if status, ok := m.calendar.GetTodoStatus()[dateKey]; ok {
								m.calendar.MoveListSelection(-1, len(status.Items))
							}
							return m, nil
						case " ", "x":
							// Toggle completion in List Mode
							dateKey := m.calendar.Cursor().Format("2006-01-02")
							if status, ok := m.calendar.GetTodoStatus()[dateKey]; ok && len(status.Items) > 0 {
								idx := m.calendar.GetSelectedListItem()
								if idx < len(status.Items) {
									cursorDate := m.calendar.Cursor()
									m.todoService.ToggleComplete(cursorDate, idx)
									m.updateTodos()
									return m, m.setStatus("Toggled completion", "success")
								}
							}
							return m, nil
						case "e", "enter":
							// Edit task in List Mode
							dateKey := m.calendar.Cursor().Format("2006-01-02")
							if status, ok := m.calendar.GetTodoStatus()[dateKey]; ok && len(status.Items) > 0 {
								idx := m.calendar.GetSelectedListItem()
								if idx < len(status.Items) {
									todos := m.todoService.GetTodosForDate(m.calendar.Cursor())
									if idx < len(todos) {
										m.state = ui.StateEditing
										m.editingIndex = idx
										m.titleInput.SetValue(todos[idx].Title)
										m.descInput.SetValue(todos[idx].Desc)
										m.editingPriority = todos[idx].Priority
										m.editFocus = ui.FocusTitle
										return m, m.titleInput.Focus()
									}
								}
							}
							return m, nil
						case "esc":
							// Go back from Day View
							m.calendar.GoBack()
							m.calendar.ResetDayViewSelection()
							m.updateTodos()
							return m, nil
						}
					} else {
						// Timeline Mode navigation
						switch msg.String() {
						case "tab":
							m.calendar.ToggleDayViewFocus()
							return m, nil
						case "j", "down":
							if m.calendar.GetDayViewFocus() == calendar.DayViewFocusTimeline {
								m.calendar.MoveTimelineCursor(1)
							} else {
								unscheduled := m.calendar.GetUnscheduledItems()
								m.calendar.MoveUnscheduledSelection(1, len(unscheduled))
							}
							return m, nil
						case "k", "up":
							if m.calendar.GetDayViewFocus() == calendar.DayViewFocusTimeline {
								m.calendar.MoveTimelineCursor(-1)
							} else {
								unscheduled := m.calendar.GetUnscheduledItems()
								m.calendar.MoveUnscheduledSelection(-1, len(unscheduled))
							}
							return m, nil
						case "enter":
							// Cursor-based scheduling: assign selected unscheduled task to timeline cursor
							if m.calendar.GetDayViewFocus() == calendar.DayViewFocusTimeline {
								unscheduled := m.calendar.GetUnscheduledItems()
								if len(unscheduled) > 0 {
									selIdx := m.calendar.GetSelectedUnscheduledIndex()
									// Find original index in full list
									origIdx := m.findOriginalTodoIndex(selIdx)
									if origIdx >= 0 {
										startTime := m.calendar.GetTimelineCursorTime()
										endTime := m.addHourToTime(startTime)
										cursorDate := m.calendar.Cursor()
										m.todoService.ScheduleTodo(cursorDate, origIdx, startTime, endTime)
										m.updateTodos()
										return m, m.setStatus("Scheduled at "+startTime, "success")
									}
								}
							} else {
								// Switch to timeline focus
								m.calendar.SetDayViewFocus(calendar.DayViewFocusTimeline)
							}
							return m, nil
						case "s":
							// Quick schedule input
							if m.calendar.GetDayViewFocus() == calendar.DayViewFocusUnscheduled {
								unscheduled := m.calendar.GetUnscheduledItems()
								if len(unscheduled) > 0 {
									selIdx := m.calendar.GetSelectedUnscheduledIndex()
									origIdx := m.findOriginalTodoIndex(selIdx)
									if origIdx >= 0 && selIdx < len(unscheduled) {
										m.state = ui.StateScheduling
										m.schedulingTaskIdx = origIdx
										m.scheduleInput.Reset()
										m.scheduleError = ""
										return m, m.scheduleInput.Focus()
									}
								}
							}
							return m, nil
						case "u":
							// Unschedule task
							if m.calendar.GetDayViewFocus() == calendar.DayViewFocusTimeline {
								scheduled := m.calendar.GetScheduledItems()
								if len(scheduled) > 0 {
									// Find the task at current timeline cursor
									cursorTime := m.calendar.GetTimelineCursorTime()
									cursorDate := m.calendar.Cursor()
									todos := m.todoService.GetTodosForDate(cursorDate)
									for i, td := range todos {
										if td.StartTime != nil && *td.StartTime == cursorTime {
											m.todoService.UnscheduleTodo(cursorDate, i)
											m.updateTodos()
											return m, m.setStatus("Task unscheduled", "info")
										}
									}
								}
							}
							return m, nil
						case "+", "=":
							// Extend duration
							if m.calendar.GetDayViewFocus() == calendar.DayViewFocusTimeline {
								cursorTime := m.calendar.GetTimelineCursorTime()
								cursorDate := m.calendar.Cursor()
								todos := m.todoService.GetTodosForDate(cursorDate)
								for i, td := range todos {
									if td.StartTime != nil && *td.StartTime == cursorTime {
										newEnd, _ := m.todoService.AdjustTodoDuration(cursorDate, i, 30, 30)
										if newEnd != "" {
											m.updateTodos()
											return m, m.setStatus("Extended to "+newEnd, "info")
										}
										return m, m.setStatus("Cannot extend further", "warning")
									}
								}
							}
							return m, nil
						case "-", "_":
							// Shrink duration
							if m.calendar.GetDayViewFocus() == calendar.DayViewFocusTimeline {
								cursorTime := m.calendar.GetTimelineCursorTime()
								cursorDate := m.calendar.Cursor()
								todos := m.todoService.GetTodosForDate(cursorDate)
								for i, td := range todos {
									if td.StartTime != nil && *td.StartTime == cursorTime {
										newEnd, _ := m.todoService.AdjustTodoDuration(cursorDate, i, -30, 30)
										if newEnd != "" {
											m.updateTodos()
											return m, m.setStatus("Shrunk to "+newEnd, "info")
										}
										return m, m.setStatus("Minimum duration reached", "warning")
									}
								}
							}
							return m, nil
						case "J":
							// Move scheduled task later (Shift+J)
							if m.calendar.GetDayViewFocus() == calendar.DayViewFocusTimeline {
								cursorTime := m.calendar.GetTimelineCursorTime()
								cursorDate := m.calendar.Cursor()
								todos := m.todoService.GetTodosForDate(cursorDate)
								config := m.calendar.GetTimelineConfig()
								for i, td := range todos {
									if td.StartTime != nil && *td.StartTime == cursorTime {
										newStart, _ := m.todoService.RescheduleTodo(cursorDate, i, config.MoveMinutes, config.DayStart, config.DayEnd)
										if newStart != "" {
											// Update timeline cursor to follow the task
											m.calendar.MoveTimelineCursor(1)
											m.updateTodos()
											return m, m.setStatus("Moved to "+newStart, "info")
										}
										return m, m.setStatus("Cannot move later", "warning")
									}
								}
							}
							return m, nil
						case "K":
							// Move scheduled task earlier (Shift+K)
							if m.calendar.GetDayViewFocus() == calendar.DayViewFocusTimeline {
								cursorTime := m.calendar.GetTimelineCursorTime()
								cursorDate := m.calendar.Cursor()
								todos := m.todoService.GetTodosForDate(cursorDate)
								config := m.calendar.GetTimelineConfig()
								for i, td := range todos {
									if td.StartTime != nil && *td.StartTime == cursorTime {
										newStart, _ := m.todoService.RescheduleTodo(cursorDate, i, -config.MoveMinutes, config.DayStart, config.DayEnd)
										if newStart != "" {
											// Update timeline cursor to follow the task
											m.calendar.MoveTimelineCursor(-1)
											m.updateTodos()
											return m, m.setStatus("Moved to "+newStart, "info")
										}
										return m, m.setStatus("Cannot move earlier", "warning")
									}
								}
							}
							return m, nil
						case "esc":
							// Go back from Day View
							m.calendar.GoBack()
							m.calendar.ResetDayViewSelection()
							m.updateTodos()
							return m, nil
						}
					}
				} else {
					// Non-Day View calendar navigation
					switch msg.String() {
					case "enter":
						m.focus = ui.FocusTodo
					case "esc":
						// Go back one view level (Day -> Week -> Month)
						m.calendar.GoBack()
						m.updateTodos()
					}
				}
			case ui.FocusTodo:
				switch msg.String() {
				case "e", "enter":
					selected := m.todo.SelectedItem()
					if selected != nil {
						m.state = ui.StateEditing
						m.editingIndex = m.todo.ListIndex()
						selectedItem := selected.(ui.TodoItem)
						m.titleInput.SetValue(selectedItem.Todo.Title)
						m.descInput.SetValue(selectedItem.Desc)
						m.editingPriority = selectedItem.Priority
						m.editFocus = ui.FocusTitle
						return m, m.titleInput.Focus()
					}
				case "esc":
					m.focus = ui.FocusCalendar
				case " ", "x":
					// Toggle completion status
					if m.todo.SelectedItem() != nil {
						cursorDate := m.calendar.Cursor()
						idx := m.todo.ListIndex()
						selectedItem := m.todo.SelectedItem().(ui.TodoItem)
						wasComplete := selectedItem.Todo.Complete
						m.todoService.ToggleComplete(cursorDate, idx)
						m.updateTodos()
						if wasComplete {
							return m, m.setStatus("Marked incomplete", "info")
						}
						return m, m.setStatus("Marked complete", "success")
					}
				case "d", "backspace":
					selectedItem := m.todo.SelectedItem()
					if selectedItem != nil {
						m.state = ui.StateConfirmingDelete
						m.deletingIndex = m.todo.ListIndex()
						m.deletingTitle = selectedItem.(ui.TodoItem).Todo.Title
					}
				case "1", "2", "3", "0":
					// Quick priority change
					if m.todo.SelectedItem() != nil {
						cursorDate := m.calendar.Cursor()
						idx := m.todo.ListIndex()
						var priority domain.Priority
						var priorityName string
						switch msg.String() {
						case "0":
							priority = domain.PriorityNone
							priorityName = "None"
						case "1":
							priority = domain.PriorityLow
							priorityName = "Low"
						case "2":
							priority = domain.PriorityMedium
							priorityName = "Medium"
						case "3":
							priority = domain.PriorityHigh
							priorityName = "High"
						}
						m.todoService.SetPriority(cursorDate, idx, priority)
						m.updateTodos()
						return m, m.setStatus("Priority: "+priorityName, "priority")
					}
				case "J", "K":
					// Move todo up/down (reorder)
					if m.todo.SelectedItem() != nil {
						cursorDate := m.calendar.Cursor()
						idx := m.todo.ListIndex()
						if msg.String() == "J" {
							m.todoService.MoveDown(cursorDate, idx)
							m.updateTodos()
							m.todo, _ = m.todo.Update(tea.KeyMsg{Type: tea.KeyDown})
						} else if msg.String() == "K" {
							m.todoService.MoveUp(cursorDate, idx)
							m.updateTodos()
							m.todo, _ = m.todo.Update(tea.KeyMsg{Type: tea.KeyUp})
						}
					}
				}
			}

		case ui.StateConfirmingDelete:
			switch msg.String() {
			case "y", "Y", "enter":
				// Confirm delete
				cursorDate := m.calendar.Cursor()
				m.todoService.Delete(cursorDate, m.deletingIndex)
				m.updateTodos()
				m.state = ui.StateViewing
				return m, m.setStatus("Todo deleted", "warning")
			case "n", "N", "esc":
				// Cancel delete
				m.state = ui.StateViewing
				return m, nil
			}

		case ui.StateEditing:
			switch msg.String() {
			case "esc":
				m.state = ui.StateViewing
				m.titleInput.Blur()
				m.descInput.Blur()
				return m, nil
			case "ctrl+p":
				// Toggle preview pane
				m.previewEnabled = !m.previewEnabled
				// Recalculate dimensions
				dims := m.viewRenderer.CalculateModalDimensions(m.previewEnabled)
				m.titleInput.Width = dims.InputWidth
				m.descInput.SetWidth(dims.InputWidth)
				if dims.ShowPreview {
					m.markdownRenderer.SetWidth(dims.PreviewWidth - 4)
				}
				return m, nil
			case "ctrl+1":
				m.editingPriority = domain.PriorityLow
				return m, nil
			case "ctrl+2":
				m.editingPriority = domain.PriorityMedium
				return m, nil
			case "ctrl+3":
				m.editingPriority = domain.PriorityHigh
				return m, nil
			case "ctrl+0":
				m.editingPriority = domain.PriorityNone
				return m, nil
			case "tab":
				if m.editFocus == ui.FocusTitle {
					m.editFocus = ui.FocusDesc
					m.titleInput.Blur()
					cmd = m.descInput.Focus()
				} else {
					m.editFocus = ui.FocusTitle
					m.descInput.Blur()
					cmd = m.titleInput.Focus()
				}
				cmds = append(cmds, cmd)
				return m, tea.Batch(cmds...)
			case "enter":
				// Only save if enter is pressed on the title input
				if m.editFocus == ui.FocusTitle {
					cursorDate := m.calendar.Cursor()
					title := m.titleInput.Value()
					desc := m.descInput.Value()

					if m.editingIndex == -1 {
						// Add new todo
						m.todoService.Add(cursorDate, title, desc, m.editingPriority)
					} else {
						// Update existing todo
						m.todoService.Update(cursorDate, m.editingIndex, title, desc, m.editingPriority)
					}
					m.updateTodos()
					m.state = ui.StateViewing
					m.titleInput.Blur()
					m.descInput.Blur()
					return m, m.setStatus("Todo saved", "success")
				}
			}

		case ui.StateSearching:
			switch msg.String() {
			case "esc":
				m.state = ui.StateViewing
				m.searchInput.Blur()
				return m, nil
			case "enter":
				// Jump to selected result
				if len(m.searchResults) > 0 {
					m.jumpToSearchResult()
				}
				m.state = ui.StateViewing
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

		case ui.StateGoToDate:
			switch msg.String() {
			case "esc":
				m.state = ui.StateViewing
				m.dateInput.Blur()
				m.dateError = ""
				return m, nil
			case "enter":
				// Parse and jump to date
				dateStr := m.dateInput.Value()
				parsedDate, err := m.parseDate(dateStr)
				if err != nil {
					m.dateError = "Invalid date format"
					return m, nil
				}
				m.calendar.SetCursor(parsedDate)
				m.updateTodos()
				m.state = ui.StateViewing
				m.dateInput.Blur()
				m.dateError = ""
				return m, m.setStatus("Jumped to "+parsedDate.Format("Jan 2, 2006"), "success")
			default:
				// Update date input
				m.dateInput, cmd = m.dateInput.Update(msg)
				m.dateError = "" // Clear error on typing
				return m, cmd
			}

		case ui.StateScheduling:
			switch msg.String() {
			case "esc":
				m.state = ui.StateViewing
				m.scheduleInput.Blur()
				m.scheduleError = ""
				return m, nil
			case "enter":
				// Parse and schedule
				timeStr := m.scheduleInput.Value()
				startTime, endTime, err := m.parseTimeInput(timeStr)
				if err != nil {
					m.scheduleError = "Invalid format (HH:MM or HH:MM-HH:MM)"
					return m, nil
				}
				cursorDate := m.calendar.Cursor()
				m.todoService.ScheduleTodo(cursorDate, m.schedulingTaskIdx, startTime, endTime)
				m.updateTodos()
				m.state = ui.StateViewing
				m.scheduleInput.Blur()
				m.scheduleError = ""
				return m, m.setStatus("Scheduled at "+startTime, "success")
			default:
				// Update schedule input
				m.scheduleInput, cmd = m.scheduleInput.Update(msg)
				m.scheduleError = "" // Clear error on typing
				return m, cmd
			}
		}
	}

	// --- Pass messages to focused components ---
	prevCursor := m.calendar.Cursor()
	switch m.state {
	case ui.StateViewing:
		if m.focus == ui.FocusCalendar {
			var newCal tea.Model
			newCal, cmd = m.calendar.Update(msg)
			m.calendar = newCal.(*calendar.Model)
		} else {
			m.todo, cmd = m.todo.Update(msg)
		}
	case ui.StateEditing:
		if m.editFocus == ui.FocusTitle {
			m.titleInput, cmd = m.titleInput.Update(msg)
		} else {
			m.descInput, cmd = m.descInput.Update(msg)
		}
	case ui.StateSearching:
		m.searchInput, cmd = m.searchInput.Update(msg)
		m.performSearch(m.searchInput.Value())
	}
	cmds = append(cmds, cmd)

	if m.state == ui.StateViewing && m.focus == ui.FocusCalendar && !prevCursor.Equal(m.calendar.Cursor()) {
		m.updateTodos()
	}

	return m, tea.Batch(cmds...)
}

func (m *model) View() string {
	// Check if terminal is too small
	if m.viewRenderer.IsTooSmall() {
		return m.viewRenderer.RenderSizeWarning()
	}

	var content string

	// Build Day View context for help bar
	dayViewCtx := ui.DayViewContext{
		IsInDayView: m.calendar.GetViewMode() == calendar.DayView,
		IsListMode:  m.calendar.GetDayViewMode() == calendar.DayViewModeList,
		IsTimeline:  m.calendar.GetDayViewFocus() == calendar.DayViewFocusTimeline,
	}
	helpBar := m.viewRenderer.RenderHelpBarWithDayView(m.state, m.focus, dayViewCtx)

	switch m.state {
	case ui.StateViewing:
		mainState := ui.MainViewState{
			CalendarView: m.calendar.View(),
			TodoView:     m.todo.View(),
			Focus:        m.focus,
		}
		content = m.viewRenderer.RenderMain(mainState)

	case ui.StateEditing:
		// Render markdown preview from description
		previewContent := ""
		if m.previewEnabled {
			previewContent = m.markdownRenderer.Render(m.descInput.Value())
		}

		editState := ui.EditingState{
			IsNew:          m.editingIndex == -1,
			Date:           m.calendar.Cursor(),
			TitleValue:     m.titleInput.Value(),
			DescValue:      m.descInput.Value(),
			Priority:       m.editingPriority,
			Focus:          m.editFocus,
			TitleView:      m.titleInput.View(),
			DescView:       m.descInput.View(),
			PreviewEnabled: m.previewEnabled,
			PreviewContent: previewContent,
		}
		content = m.viewRenderer.RenderEditing(editState)

	case ui.StateConfirmingDelete:
		deleteState := ui.DeleteState{
			Title: m.deletingTitle,
		}
		content = m.viewRenderer.RenderConfirmDelete(deleteState)

	case ui.StateSearching:
		searchState := ui.SearchState{
			InputView:   m.searchInput.View(),
			InputValue:  m.searchInput.Value(),
			Results:     m.searchResults,
			SelectedIdx: m.searchIndex,
		}
		content = m.viewRenderer.RenderSearching(searchState)

	case ui.StateHelp:
		content = m.viewRenderer.RenderHelp()

	case ui.StateGoToDate:
		goToDateState := ui.GoToDateState{
			InputView:  m.dateInput.View(),
			InputValue: m.dateInput.Value(),
			ErrorMsg:   m.dateError,
		}
		content = m.viewRenderer.RenderGoToDate(goToDateState)

	case ui.StateScheduling:
		// Get task title for the scheduling modal
		taskTitle := ""
		cursorDate := m.calendar.Cursor()
		todos := m.todoService.GetTodosForDate(cursorDate)
		if m.schedulingTaskIdx >= 0 && m.schedulingTaskIdx < len(todos) {
			taskTitle = todos[m.schedulingTaskIdx].Title
		}
		schedulingState := ui.SchedulingState{
			TaskTitle:  taskTitle,
			TaskIndex:  m.schedulingTaskIdx,
			InputView:  m.scheduleInput.View(),
			InputValue: m.scheduleInput.Value(),
			ErrorMsg:   m.scheduleError,
		}
		content = m.viewRenderer.RenderScheduling(schedulingState)

	default:
		return "unknown state"
	}

	// Build the final view with optional status message
	statusMsg := m.viewRenderer.RenderStatusMessage(m.statusMessage, m.statusType)
	if statusMsg != "" {
		return lipgloss.JoinVertical(lipgloss.Left, content, statusMsg, helpBar)
	}
	return lipgloss.JoinVertical(lipgloss.Left, content, helpBar)
}

func main() {
	// Initialize dependencies
	repo := repository.NewJSONTodoRepository(getDataFilePath())
	if err := repo.Load(); err != nil {
		fmt.Printf("Error loading todos: %v\n", err)
	}

	timeProv := service.NewRealTimeProvider()
	statsCalc := service.NewStatsCalculator(timeProv)
	todoService := service.NewTodoService(repo, timeProv)
	presenter := ui.NewTodoPresenter()
	calendarAdapter := ui.NewCalendarAdapter(statsCalc)
	viewRenderer := ui.NewViewRenderer()

	// Initialize Title Input
	ti := textinput.New()
	ti.Placeholder = "What needs to be done?"
	ti.Focus()
	ti.CharLimit = 256
	ti.Width = 56

	// Initialize Description Input
	ta := textarea.New()
	ta.Placeholder = "Add details (Markdown supported)..."
	ta.KeyMap.InsertNewline.SetEnabled(true) // Allow newlines in description
	ta.SetWidth(56)
	ta.SetHeight(8)

	// Initialize Search Input
	si := textinput.New()
	si.Placeholder = "Search todos..."
	si.CharLimit = 100
	si.Width = 38

	// Initialize Date Input
	di := textinput.New()
	di.Placeholder = "YYYY-MM-DD"
	di.CharLimit = 10
	di.Width = 20

	// Initialize Schedule Input
	sci := textinput.New()
	sci.Placeholder = "09:00 or 09:00-10:30"
	sci.CharLimit = 11
	sci.Width = 20

	// Initialize Markdown Renderer
	mdRenderer := ui.NewMarkdownRenderer(40) // Initial width, will be resized

	m := &model{
		// Services
		todoService:     todoService,
		statsCalc:       statsCalc,
		presenter:       presenter,
		calendarAdapter: calendarAdapter,
		viewRenderer:    viewRenderer,

		// UI Components
		calendar:      calendar.New(),
		todo:          todo.New(),
		titleInput:    ti,
		descInput:     ta,
		searchInput:   si,
		dateInput:     di,
		scheduleInput: sci,

		// State
		state:        ui.StateViewing,
		focus:        ui.FocusCalendar,
		editingIndex: -1,
		editFocus:    ui.FocusTitle,

		// Preview
		markdownRenderer: mdRenderer,
		previewEnabled:   true, // Preview enabled by default
	}
	p := tea.NewProgram(m, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
