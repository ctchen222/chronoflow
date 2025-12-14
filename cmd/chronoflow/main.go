package main

import (
	"encoding/json"
	"fmt"
	"os"

	"ctchen222/chronoflow/pkg/calendar"
	"ctchen222/chronoflow/pkg/todo"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const dbFile = "todos.json"

// item represents a to-do item in our list.
type item struct {
	ItemTitle string `json:"title"`
	ItemDesc  string `json:"desc"`
}

func (i item) Title() string       { return i.ItemTitle }
func (i item) Description() string { return i.ItemDesc }
func (i item) FilterValue() string { return i.ItemTitle }

// todos is our in-memory data store.
var todos = make(map[string][]item)

// AppState defines the current state of the application.
type appState int

const (
	viewing appState = iota
	editing
)

// appFocus determines which panel is currently focused.
type appFocus int

const (
	calendarFocus appFocus = iota
	todoFocus
)

type model struct {
	calendar     *calendar.Model
	todo         todo.Model
	textarea     textarea.Model
	state        appState
	focus        appFocus
	editingIndex int // Index of the todo item being edited. -1 means not editing.
	width        int
	height       int
}

// updateTodos sets the items for the todo list based on the selected date.
func (m *model) updateTodos() {
	dateKey := m.calendar.Cursor().Format("2006-01-02")
	items := []list.Item{}
	if foundTodos, ok := todos[dateKey]; ok {
		for _, td := range foundTodos {
			items = append(items, td)
		}
	}
	m.todo.SetItems(items)
	m.todo.SetTitle(fmt.Sprintf("To-Do on %s", m.calendar.Cursor().Format("2006-01-02")))
}

func (m *model) Init() tea.Cmd {
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
		calendarWidth := int(float64(m.width) * 0.7)
		rightPanelWidth := m.width - calendarWidth

		// For viewing state
		m.calendar.SetSize(calendarWidth, m.height)
		m.todo.SetSize(rightPanelWidth, m.height)

		// For editing state
		m.textarea.SetWidth(m.width - 4) // Some padding
		m.textarea.SetHeight(m.height/3 - 2)
		return m, nil

	case tea.KeyMsg:
		// Switch between main application states (viewing, editing)
		switch m.state {
		case viewing:
			// Keys that are always active in viewing mode
			switch msg.String() {
			case "q", "ctrl+c":
				saveTodos() // Save on quit
				return m, tea.Quit
			case "a": // 'a' to add a new todo
				m.editingIndex = -1 // Explicitly set to -1 for adding a new item
				m.state = editing
				m.textarea.Reset()
				return m, m.textarea.Focus()
			case "tab":
				if m.focus == calendarFocus {
					m.focus = todoFocus
				} else {
					m.focus = calendarFocus
				}
			}

			// Handle keys based on which panel is focused
			switch m.focus {
			case calendarFocus:
				switch msg.String() {
				case "enter":
					m.focus = todoFocus
				}
			case todoFocus:
				switch msg.String() {
				case "e", "enter":
					selectedItem := m.todo.SelectedItem()
					if selectedItem != nil {
						m.state = editing
						m.editingIndex = m.todo.ListIndex()
						m.textarea.SetValue(selectedItem.FilterValue())
						return m, m.textarea.Focus()
					}
				case "esc":
					m.focus = calendarFocus
				case "d", "backspace":
					selectedItem := m.todo.SelectedItem()
					if selectedItem != nil {
						dateKey := m.calendar.Cursor().Format("2006-01-02")
						// Filter out the item to be deleted
						newItems := []item{}
						for _, i := range todos[dateKey] {
							if i.Title() != selectedItem.FilterValue() {
								newItems = append(newItems, i)
							}
						}
						todos[dateKey] = newItems
						m.updateTodos() // Refresh the list
					}
				}
			}

		case editing:
			switch msg.String() {
			case "esc":
				m.state = viewing
				m.textarea.Blur()
				m.editingIndex = -1 // Reset editing state
				return m, nil
			case "enter":
				newTodoText := m.textarea.Value()
				dateKey := m.calendar.Cursor().Format("2006-01-02")

				if m.editingIndex == -1 {
					// Add new item
					if newTodoText != "" {
						todos[dateKey] = append(todos[dateKey], item{ItemTitle: newTodoText, ItemDesc: ""})
					}
				} else {
					// Update existing item
					if items, ok := todos[dateKey]; ok {
						if len(items) > m.editingIndex {
							items[m.editingIndex].ItemTitle = newTodoText
							todos[dateKey] = items // Re-assign the modified slice
						}
					}
				}

				m.updateTodos() // Refresh the todo list
				m.state = viewing
				m.textarea.Blur()
				m.textarea.Reset()
				m.editingIndex = -1 // Reset editing state
				return m, nil
			}
		}
	}

	// --- Pass messages to focused components ---
	prevCursor := m.calendar.Cursor() // Store cursor before updates

	switch m.state {
	case viewing:
		switch m.focus {
		        case calendarFocus:
		            // Pass updates to the calendar model.
		            newCalendar, calendarCmd := m.calendar.Update(msg)
		            m.calendar = newCalendar.(*calendar.Model)
		            cmd = calendarCmd
		        case todoFocus:
		            // Pass updates to the todo model.
		            m.todo, cmd = m.todo.Update(msg)
		        }
		        cmds = append(cmds, cmd)
		
		    case editing:
		        m.textarea, cmd = m.textarea.Update(msg)
		        cmds = append(cmds, cmd)
		    }
	// If calendar navigation changed the date, update the todos list
	if m.state == viewing && m.focus == calendarFocus && !prevCursor.Equal(m.calendar.Cursor()) {
		m.updateTodos()
	}

	return m, tea.Batch(cmds...)
}

func (m *model) View() string {
	switch m.state {
	case viewing:
		// Recalculate widths here to ensure they are always in sync.
		calendarWidth := int(float64(m.width) * 0.7)
		rightPanelWidth := m.width - calendarWidth
		
		// Ensure child components have the correct size.
		m.calendar.SetSize(calendarWidth, m.height)
		m.todo.SetSize(rightPanelWidth, m.height) // Pass full height to todo now, let its internal lipgloss handle padding/borders

		return lipgloss.JoinHorizontal(lipgloss.Top,
			m.calendar.View(),
			m.todo.View(),
		)
	case editing:
		return lipgloss.NewStyle().
			Width(m.width).
			Height(m.height).
			Align(lipgloss.Center, lipgloss.Center).
			Render(m.textarea.View())
	default:
		return "unknown state"
	}
}

// saveTodos saves the current to-do items to the JSON database file.
func saveTodos() {
	data, err := json.MarshalIndent(todos, "", "  ")
	if err != nil {
		fmt.Printf("Error marshaling todos to JSON: %v\n", err)
		return
	}

	err = os.WriteFile(dbFile, data, 0644)
	if err != nil {
		fmt.Printf("Error writing to db file: %v\n", err)
	}
}

// loadTodos loads the to-do items from the JSON database file.
func loadTodos() {
	data, err := os.ReadFile(dbFile)
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
	loadTodos() // Load on startup

	ta := textarea.New()
	ta.Placeholder = "Enter new to-do item (Enter to save, Esc to cancel)..."
	ta.Focus()
	ta.KeyMap.InsertNewline.SetEnabled(false)

	m := &model{
		calendar:     calendar.New(),
		todo:         todo.New(),
		textarea:     ta,
		state:        viewing,
		focus:        calendarFocus,
		editingIndex: -1,
	}
	p := tea.NewProgram(m, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
