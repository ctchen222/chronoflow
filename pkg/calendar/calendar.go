package calendar

import (
	"fmt"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	// Style for the highlighted day (cursor) in the calendar.
	dayHighlight = lipgloss.NewStyle().
		Background(lipgloss.Color("#7D56F4")).
		Foreground(lipgloss.Color("#FFFFFF")).
		Bold(true)

	// Style for today's date (when not selected).
	todayStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#7D56F4")).
		Bold(true).
		Underline(true)
)

// TodoItem represents a single todo for display in week view
type TodoItem struct {
	Title    string
	Complete bool
	Priority int // 0=none, 1=low, 2=medium, 3=high
}

// TodoStatus represents the status of todos for a date
type TodoStatus struct {
	HasTodos    bool
	HasOverdue  bool       // has incomplete todos on a past date
	AllComplete bool       // all todos are completed
	Count       int        // total number of todos for this date
	Items       []TodoItem // actual todo items for week view display
}

// ViewMode represents the calendar view mode
type ViewMode int

const (
	MonthView ViewMode = iota
	WeekView
)

type Model struct {
	selectedDate time.Time
	cursor       time.Time
	width        int
	height       int
	todoStatus   map[string]TodoStatus // todo status by date (format: "2006-01-02")
	viewMode     ViewMode
}

func New() *Model {
	now := time.Now()
	return &Model{
		selectedDate: now,
		cursor:       now,
		todoStatus:   make(map[string]TodoStatus),
		viewMode:     MonthView,
	}
}

// ToggleViewMode switches between month and week view
func (m *Model) ToggleViewMode() {
	if m.viewMode == MonthView {
		m.viewMode = WeekView
	} else {
		m.viewMode = MonthView
	}
}

// GetViewMode returns the current view mode
func (m *Model) GetViewMode() ViewMode {
	return m.viewMode
}

// SetTodoStatus updates the todo status for each date
func (m *Model) SetTodoStatus(status map[string]TodoStatus) {
	m.todoStatus = status
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

// SetCursor sets the cursor to a specific date
func (m *Model) SetCursor(date time.Time) {
	m.cursor = date
	m.selectedDate = date
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
		case "t":
			// Jump to today
			now := time.Now()
			m.cursor = now
			m.selectedDate = now
		case "w":
			// Toggle week/month view
			m.ToggleViewMode()
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

	// Render week view if in week mode
	if m.viewMode == WeekView {
		return m.renderWeekView()
	}

	// Month view
	const minCalendarHeight = 20
	if m.height < minCalendarHeight {
		return lipgloss.NewStyle().
			Width(m.width).
			Height(m.height).
			Align(lipgloss.Center, lipgloss.Center).
			Render("Terminal too small")
	}



	var s strings.Builder



	// --- RENDER AND APPEND, WITH CORRECT NEWLINE MANAGEMENT ---



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



	// --- Grid ---
	// Available height calculation:
	// Header overhead: main header (1) + subheader (1) + weekday header with border (2) + 3 newlines = 7
	// Plus 1 for top border of first row + 1 for margins = 9 total
	availableHeight := m.height - 9

	baseCellHeight := availableHeight / 6

	heightRemainder := availableHeight % 6

	cellHeights := make([]int, 6)

	for i := 0; i < 6; i++ {

		cellHeights[i] = baseCellHeight

		if heightRemainder > 0 {

			cellHeights[i]++

			heightRemainder--

		}

	}



	firstDay := time.Date(m.selectedDate.Year(), m.selectedDate.Month(), 1, 0, 0, 0, 0, time.UTC)

	daysInMonth := time.Date(m.selectedDate.Year(), m.selectedDate.Month()+1, 0, 0, 0, 0, 0, time.UTC).Day()

	day := 1

	for i := 0; i < 6; i++ {

		var rowViews []string

		for j := 0; j < 7; j++ {

			style := lipgloss.NewStyle().
				Width(cellWidths[j]).
				Height(cellHeights[i]).
				Border(lipgloss.NormalBorder(), false).
				Align(lipgloss.Center, lipgloss.Center).
				BorderBottom(true).
				BorderForeground(lipgloss.Color("238"))

			// Add top border for first row
			if i == 0 {
				style = style.BorderTop(true)
			}

			if j != 6 {
				style = style.BorderRight(true)
			}



			if (i == 0 && j < int(firstDay.Weekday())) || day > daysInMonth {

				rowViews = append(rowViews, style.Render(""))

			} else {

				now := time.Now()
				isToday := day == now.Day() && m.selectedDate.Month() == now.Month() && m.selectedDate.Year() == now.Year()
				isCursor := day == m.cursor.Day() && m.selectedDate.Month() == m.cursor.Month() && m.selectedDate.Year() == m.cursor.Year()

				// Check todo status for this date
				dateKey := time.Date(m.selectedDate.Year(), m.selectedDate.Month(), day, 0, 0, 0, 0, time.UTC).Format("2006-01-02")
				status := m.todoStatus[dateKey]

				if isCursor {
					style = style.Copy().Inherit(dayHighlight)
				} else if isToday {
					style = style.Copy().Inherit(todayStyle)
				}

				// Apply color based on todo status (only if not cursor/today which have their own style)
				dayNum := fmt.Sprintf("%d", day)

				if status.HasTodos && !isCursor {
					if status.HasOverdue {
						// Red for overdue
						style = style.Copy().Foreground(lipgloss.Color("#FF6B6B")).Bold(true)
					} else if status.AllComplete {
						// Gray for all complete
						style = style.Copy().Foreground(lipgloss.Color("#666"))
					} else {
						// Green for has pending todos
						style = style.Copy().Foreground(lipgloss.Color("#50FA7B"))
					}
				}

				// Show todo count below the day number if there are todos
				cellContent := dayNum
				if status.Count > 0 {
					cellContent = fmt.Sprintf("%s\n(%d)", dayNum, status.Count)
				}

				rowViews = append(rowViews, style.Render(cellContent))

				day++

			}

		}

		s.WriteString(lipgloss.JoinHorizontal(lipgloss.Top, rowViews...))

		if i < 5 {

			s.WriteString("\n") // Add necessary newline between rows

		}

	}



	return s.String()

}

// renderWeekView renders the week view of the calendar
func (m *Model) renderWeekView() string {
	var s strings.Builder

	// Find the start of the week (Sunday) containing the cursor
	weekStart := m.cursor.AddDate(0, 0, -int(m.cursor.Weekday()))

	// Main Header
	weekEnd := weekStart.AddDate(0, 0, 6)
	headerText := fmt.Sprintf("Week of %s - %s", weekStart.Format("Jan 2"), weekEnd.Format("Jan 2, 2006"))

	header := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#FAFAFA")).
		Background(lipgloss.Color("#7D56F4")).
		Padding(0, 1).
		Width(m.width).
		Align(lipgloss.Center).
		Render(headerText)

	s.WriteString(header)
	s.WriteString("\n")

	// Sub-header with current selected date
	subHeader := lipgloss.NewStyle().
		Foreground(lipgloss.Color("240")).
		Padding(0, 1).
		Width(m.width).
		Align(lipgloss.Center).
		Render(m.cursor.Format("Monday, Jan 2, 2006"))

	s.WriteString(subHeader)
	s.WriteString("\n")

	// Calculate cell widths for 7 columns
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

	// Weekday headers
	weekdays := []string{"Sun", "Mon", "Tue", "Wed", "Thu", "Fri", "Sat"}
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

	// Calculate available height for the week row
	// Header overhead: main header (1) + subheader (1) + weekday header with border (2) + 3 newlines = 7
	availableHeight := m.height - 7

	// Calculate how many todo items can fit per cell
	// Reserve 2 lines for day number and margin, rest for todos
	maxTodoLines := availableHeight - 2
	if maxTodoLines < 1 {
		maxTodoLines = 1
	}

	// Render the 7 days in a single row
	now := time.Now()
	var dayViews []string

	for i := 0; i < 7; i++ {
		currentDay := weekStart.AddDate(0, 0, i)
		dateKey := currentDay.Format("2006-01-02")
		status := m.todoStatus[dateKey]

		isToday := currentDay.Day() == now.Day() &&
			currentDay.Month() == now.Month() &&
			currentDay.Year() == now.Year()
		isCursor := currentDay.Day() == m.cursor.Day() &&
			currentDay.Month() == m.cursor.Month() &&
			currentDay.Year() == m.cursor.Year()

		style := lipgloss.NewStyle().
			Width(cellWidths[i]).
			Height(availableHeight).
			Align(lipgloss.Left, lipgloss.Top).
			PaddingTop(1).
			PaddingLeft(1)

		if i != 6 {
			style = style.BorderRight(true).
				Border(lipgloss.NormalBorder(), false, true, false, false).
				BorderForeground(lipgloss.Color("238"))
		}

		// Apply cursor/today styling
		if isCursor {
			style = style.Copy().Inherit(dayHighlight)
		} else if isToday {
			style = style.Copy().Inherit(todayStyle)
		}

		// Display day number and month if it's the 1st
		dayText := fmt.Sprintf("%d", currentDay.Day())
		if currentDay.Day() == 1 {
			dayText = currentDay.Format("Jan 2")
		}

		// Build cell content with todo items
		var lines []string
		lines = append(lines, dayText)

		if len(status.Items) > 0 {
			lines = append(lines, "") // Empty line after day number

			// Calculate max width for todo text (cell width minus padding)
			maxTodoWidth := cellWidths[i] - 4
			if maxTodoWidth < 5 {
				maxTodoWidth = 5
			}

			// Show todo items
			todosToShow := maxTodoLines - 1 // Reserve 1 line for "+N more" if needed
			if len(status.Items) <= maxTodoLines {
				todosToShow = len(status.Items)
			}

			for j := 0; j < todosToShow && j < len(status.Items); j++ {
				todo := status.Items[j]

				// Checkbox
				checkbox := "☐"
				if todo.Complete {
					checkbox = "☑"
				}

				// Truncate title to fit
				title := todo.Title
				if len(title) > maxTodoWidth {
					title = title[:maxTodoWidth-1] + "…"
				}

				lines = append(lines, checkbox+" "+title)
			}

			// Show "+N more" if there are more items
			remaining := len(status.Items) - todosToShow
			if remaining > 0 {
				lines = append(lines, fmt.Sprintf("+%d more", remaining))
			}
		}

		cellContent := strings.Join(lines, "\n")

		// Apply todo status colors to the entire cell (only if not cursor)
		if status.HasTodos && !isCursor && !isToday {
			if status.HasOverdue {
				style = style.Copy().Foreground(lipgloss.Color("#FF6B6B"))
			} else if status.AllComplete {
				style = style.Copy().Foreground(lipgloss.Color("#666"))
			}
		}

		dayViews = append(dayViews, style.Render(cellContent))
	}

	s.WriteString(lipgloss.JoinHorizontal(lipgloss.Top, dayViews...))

	return s.String()
}
