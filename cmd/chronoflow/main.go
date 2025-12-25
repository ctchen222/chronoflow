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
		case ui.StateViewing:
			switch msg.String() {
			case "q", "ctrl+c":
				m.todoService.Persist()
				return m, tea.Quit
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
			}

			switch m.focus {
			case ui.FocusCalendar:
				if msg.String() == "enter" {
					m.focus = ui.FocusTodo
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
						m.todoService.ToggleComplete(cursorDate, idx)
						m.updateTodos()
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
						switch msg.String() {
						case "0":
							priority = domain.PriorityNone
						case "1":
							priority = domain.PriorityLow
						case "2":
							priority = domain.PriorityMedium
						case "3":
							priority = domain.PriorityHigh
						}
						m.todoService.SetPriority(cursorDate, idx, priority)
						m.updateTodos()
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
				return m, nil
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
					return m, nil
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
	var content string
	helpBar := m.viewRenderer.RenderHelpBar(m.state, m.focus)

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

	default:
		return "unknown state"
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
		calendar:    calendar.New(),
		todo:        todo.New(),
		titleInput:  ti,
		descInput:   ta,
		searchInput: si,

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
