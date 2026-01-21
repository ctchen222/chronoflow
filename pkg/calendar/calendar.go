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

// TodoItem represents a single todo for display in week/day view
type TodoItem struct {
	Title     string
	Desc      string
	Complete  bool
	Priority  int // 0=none, 1=low, 2=medium, 3=high
	StartTime string // "HH:MM" format, empty = unscheduled
	EndTime   string // "HH:MM" format
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
	DayView
)

// DayViewFocus represents which panel has focus in Day View
type DayViewFocus int

const (
	DayViewFocusUnscheduled DayViewFocus = iota
	DayViewFocusTimeline
)

// DayViewMode represents the display mode for Day View
type DayViewMode int

const (
	DayViewModeList     DayViewMode = iota // Default: clean card-style list
	DayViewModeTimeline                    // Split-panel with timeline + unscheduled
)

// TimelineConfig holds timeline display settings
type TimelineConfig struct {
	DayStart    string // "HH:MM" format
	DayEnd      string // "HH:MM" format
	SlotMinutes int    // display granularity
	MoveMinutes int    // movement granularity for Shift+J/K
}

// DefaultTimelineConfig returns default timeline settings
func DefaultTimelineConfig() TimelineConfig {
	return TimelineConfig{
		DayStart:    "08:00",
		DayEnd:      "18:00",
		SlotMinutes: 30,
		MoveMinutes: 30, // Default to same as SlotMinutes
	}
}

type Model struct {
	selectedDate time.Time
	cursor       time.Time
	width        int
	height       int
	todoStatus   map[string]TodoStatus // todo status by date (format: "2006-01-02")
	viewMode     ViewMode

	// Day View state
	dayViewMode         DayViewMode  // List or Timeline mode
	dayViewFocus        DayViewFocus // Which panel has focus (Timeline mode only)
	timelineCursor      int          // Current time slot index (0 = first slot)
	selectedUnscheduled int          // Selected index in unscheduled list
	selectedListItem    int          // Selected index in list mode
	selectedScheduled   int          // Selected index in scheduled list (for task movement)
	timelineConfig      TimelineConfig
}

func New() *Model {
	now := time.Now()
	return &Model{
		selectedDate:        now,
		cursor:              now,
		todoStatus:          make(map[string]TodoStatus),
		viewMode:            MonthView,
		dayViewMode:         DayViewModeList, // Default to List Mode
		dayViewFocus:        DayViewFocusUnscheduled,
		timelineCursor:      0,
		selectedUnscheduled: 0,
		selectedListItem:    0,
		selectedScheduled:   0,
		timelineConfig:      DefaultTimelineConfig(),
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

// SetViewMode sets the view mode directly
func (m *Model) SetViewMode(mode ViewMode) {
	m.viewMode = mode
}

// GetViewMode returns the current view mode
func (m *Model) GetViewMode() ViewMode {
	return m.viewMode
}

// GoBack returns to the previous view level (Day -> Week -> Month)
func (m *Model) GoBack() {
	switch m.viewMode {
	case DayView:
		m.viewMode = WeekView
	case WeekView:
		m.viewMode = MonthView
	}
}

// SetTimelineConfig updates the timeline configuration
func (m *Model) SetTimelineConfig(config TimelineConfig) {
	m.timelineConfig = config
}

// GetDayViewFocus returns the current day view focus
func (m *Model) GetDayViewFocus() DayViewFocus {
	return m.dayViewFocus
}

// SetDayViewFocus sets the day view focus
func (m *Model) SetDayViewFocus(focus DayViewFocus) {
	m.dayViewFocus = focus
}

// ToggleDayViewFocus switches focus between timeline and unscheduled panels
func (m *Model) ToggleDayViewFocus() {
	if m.dayViewFocus == DayViewFocusTimeline {
		m.dayViewFocus = DayViewFocusUnscheduled
	} else {
		m.dayViewFocus = DayViewFocusTimeline
	}
}

// MoveTimelineCursor moves the timeline cursor by delta slots
func (m *Model) MoveTimelineCursor(delta int) {
	// Calculate total slots
	startHour, startMin := m.parseTime(m.timelineConfig.DayStart)
	endHour, endMin := m.parseTime(m.timelineConfig.DayEnd)
	startMinutes := startHour*60 + startMin
	endMinutes := endHour*60 + endMin
	totalSlots := (endMinutes - startMinutes) / m.timelineConfig.SlotMinutes

	newPos := m.timelineCursor + delta
	if newPos < 0 {
		newPos = 0
	}
	if newPos >= totalSlots {
		newPos = totalSlots - 1
	}
	m.timelineCursor = newPos
}

// GetTimelineCursorTime returns the time string for the current cursor position
func (m *Model) GetTimelineCursorTime() string {
	startHour, startMin := m.parseTime(m.timelineConfig.DayStart)
	startMinutes := startHour*60 + startMin
	cursorMinutes := startMinutes + m.timelineCursor*m.timelineConfig.SlotMinutes
	hour := cursorMinutes / 60
	minute := cursorMinutes % 60
	return fmt.Sprintf("%02d:%02d", hour, minute)
}

// GetSelectedUnscheduledIndex returns the selected index in unscheduled list
func (m *Model) GetSelectedUnscheduledIndex() int {
	return m.selectedUnscheduled
}

// MoveUnscheduledSelection moves the selection in unscheduled panel
func (m *Model) MoveUnscheduledSelection(delta int, maxItems int) {
	if maxItems == 0 {
		return
	}
	newPos := m.selectedUnscheduled + delta
	if newPos < 0 {
		newPos = 0
	}
	if newPos >= maxItems {
		newPos = maxItems - 1
	}
	m.selectedUnscheduled = newPos
}

// ResetDayViewSelection resets the day view selection state
func (m *Model) ResetDayViewSelection() {
	m.dayViewMode = DayViewModeList // Reset to List Mode when leaving Day View
	m.dayViewFocus = DayViewFocusUnscheduled
	m.timelineCursor = 0
	m.selectedUnscheduled = 0
	m.selectedListItem = 0
	m.selectedScheduled = 0
}

// GetDayViewMode returns the current day view mode
func (m *Model) GetDayViewMode() DayViewMode {
	return m.dayViewMode
}

// SetDayViewMode sets the day view mode
func (m *Model) SetDayViewMode(mode DayViewMode) {
	m.dayViewMode = mode
}

// ToggleDayViewMode switches between List and Timeline modes
func (m *Model) ToggleDayViewMode() {
	if m.dayViewMode == DayViewModeList {
		m.dayViewMode = DayViewModeTimeline
	} else {
		m.dayViewMode = DayViewModeList
	}
}

// GetSelectedListItem returns the selected index in list mode
func (m *Model) GetSelectedListItem() int {
	return m.selectedListItem
}

// MoveListSelection moves the selection in list mode
func (m *Model) MoveListSelection(delta int, maxItems int) {
	if maxItems == 0 {
		return
	}
	newPos := m.selectedListItem + delta
	if newPos < 0 {
		newPos = 0
	}
	if newPos >= maxItems {
		newPos = maxItems - 1
	}
	m.selectedListItem = newPos
}

// GetSelectedScheduledIndex returns the selected index in scheduled list
func (m *Model) GetSelectedScheduledIndex() int {
	return m.selectedScheduled
}

// MoveScheduledSelection moves the selection in scheduled list
func (m *Model) MoveScheduledSelection(delta int, maxItems int) {
	if maxItems == 0 {
		return
	}
	newPos := m.selectedScheduled + delta
	if newPos < 0 {
		newPos = 0
	}
	if newPos >= maxItems {
		newPos = maxItems - 1
	}
	m.selectedScheduled = newPos
}

// GetTimelineConfig returns the timeline configuration
func (m *Model) GetTimelineConfig() TimelineConfig {
	return m.timelineConfig
}

// GetUnscheduledItems returns unscheduled items for the current date
func (m *Model) GetUnscheduledItems() []TodoItem {
	dateKey := m.cursor.Format("2006-01-02")
	status := m.todoStatus[dateKey]
	_, unscheduled := m.splitItemsBySchedule(status.Items)
	return unscheduled
}

// GetScheduledItems returns scheduled items for the current date
func (m *Model) GetScheduledItems() []TodoItem {
	dateKey := m.cursor.Format("2006-01-02")
	status := m.todoStatus[dateKey]
	scheduled, _ := m.splitItemsBySchedule(status.Items)
	return scheduled
}

// GetTodoCountForDate returns the number of todos for a specific date
func (m *Model) GetTodoCountForDate(date time.Time) int {
	dateKey := date.Format("2006-01-02")
	if status, ok := m.todoStatus[dateKey]; ok {
		return status.Count
	}
	return 0
}

// renderViewModeIndicator renders the view mode badge (MONTH VIEW, WEEK VIEW, or DAY VIEW)
func (m *Model) renderViewModeIndicator() string {
	var modeText string
	switch m.viewMode {
	case MonthView:
		modeText = "MONTH VIEW"
	case WeekView:
		modeText = "WEEK VIEW"
	case DayView:
		modeText = "DAY VIEW"
	}

	return lipgloss.NewStyle().
		Foreground(lipgloss.Color("#7D56F4")).
		Bold(true).
		Render(modeText)
}

// renderBreadcrumb renders the navigation breadcrumb based on current view mode
func (m *Model) renderBreadcrumb() string {
	separatorStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("245"))
	separator := separatorStyle.Render(" › ")

	monthYear := m.cursor.Format("January 2006")

	switch m.viewMode {
	case MonthView:
		return "📅 " + monthYear
	case WeekView:
		// Calculate week number within the month
		_, weekNum := m.cursor.ISOWeek()
		return "📅 " + monthYear + separator + fmt.Sprintf("Week %d", weekNum)
	case DayView:
		_, weekNum := m.cursor.ISOWeek()
		dayStr := m.cursor.Format("Mon 2")
		return "📅 " + monthYear + separator + fmt.Sprintf("Week %d", weekNum) + separator + dayStr
	}
	return ""
}

// formatSubHeaderWithCount formats the sub-header text with optional todo count
func (m *Model) formatSubHeaderWithCount() string {
	dateStr := m.cursor.Format("Monday, Jan 2, 2006")
	count := m.GetTodoCountForDate(m.cursor)

	if count == 0 {
		return dateStr
	}

	taskWord := "tasks"
	if count == 1 {
		taskWord = "task"
	}
	return fmt.Sprintf("%s (%d %s)", dateStr, count, taskWord)
}

// SetTodoStatus updates the todo status for each date
func (m *Model) SetTodoStatus(status map[string]TodoStatus) {
	m.todoStatus = status
}

// GetTodoStatus returns the todo status map
func (m *Model) GetTodoStatus() map[string]TodoStatus {
	return m.todoStatus
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
		return ""
	}

	// Render day view if in day mode
	if m.viewMode == DayView {
		return m.renderDayView()
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



	// View mode indicator line
	viewModeIndicator := lipgloss.NewStyle().
		Width(m.width).
		Align(lipgloss.Center).
		Render(m.renderViewModeIndicator())
	s.WriteString(viewModeIndicator)
	s.WriteString("\n")

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

	// Sub-header with todo count (single line) -> requires a newline after
	subHeader := lipgloss.NewStyle().
		Foreground(lipgloss.Color("240")).
		Padding(0, 1).
		Width(m.width).
		Align(lipgloss.Center).
		Render(m.formatSubHeaderWithCount())
	s.WriteString(subHeader)
	s.WriteString("\n")

	// Breadcrumb line
	breadcrumb := lipgloss.NewStyle().
		Foreground(lipgloss.Color("245")).
		Width(m.width).
		Align(lipgloss.Center).
		Render(m.renderBreadcrumb())
	s.WriteString(breadcrumb)
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
	// Header overhead: view mode indicator (1) + main header (1) + subheader (1) + breadcrumb (1) + weekday header with border (2) + 5 newlines = 11
	// Plus 1 for top border of first row + 1 for margins = 12 total
	availableHeight := m.height - 12

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

	// View mode indicator line
	viewModeIndicator := lipgloss.NewStyle().
		Width(m.width).
		Align(lipgloss.Center).
		Render(m.renderViewModeIndicator())
	s.WriteString(viewModeIndicator)
	s.WriteString("\n")

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

	// Sub-header with current selected date and todo count
	subHeader := lipgloss.NewStyle().
		Foreground(lipgloss.Color("240")).
		Padding(0, 1).
		Width(m.width).
		Align(lipgloss.Center).
		Render(m.formatSubHeaderWithCount())

	s.WriteString(subHeader)
	s.WriteString("\n")

	// Breadcrumb line
	breadcrumb := lipgloss.NewStyle().
		Foreground(lipgloss.Color("245")).
		Width(m.width).
		Align(lipgloss.Center).
		Render(m.renderBreadcrumb())
	s.WriteString(breadcrumb)
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
	// Header overhead: view mode indicator (1) + main header (1) + subheader (1) + breadcrumb (1) + weekday header with border (2) + 5 newlines = 10
	availableHeight := m.height - 10

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

// renderDayView dispatches to the appropriate day view renderer based on mode
func (m *Model) renderDayView() string {
	if m.dayViewMode == DayViewModeList {
		return m.renderDayViewList()
	}
	return m.renderDayViewTimeline()
}

// renderDayViewTimeline renders the day view with split-panel layout (timeline + unscheduled)
func (m *Model) renderDayViewTimeline() string {
	var s strings.Builder

	dateKey := m.cursor.Format("2006-01-02")
	status := m.todoStatus[dateKey]

	// View mode indicator line
	viewModeIndicator := lipgloss.NewStyle().
		Width(m.width).
		Align(lipgloss.Center).
		Render(m.renderViewModeIndicator())
	s.WriteString(viewModeIndicator)
	s.WriteString("\n")

	// Main Header with full date
	header := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#FAFAFA")).
		Background(lipgloss.Color("#7D56F4")).
		Padding(0, 1).
		Width(m.width).
		Align(lipgloss.Center).
		Render(m.cursor.Format("Monday, January 2, 2006"))
	s.WriteString(header)
	s.WriteString("\n")

	// Breadcrumb
	breadcrumb := m.renderBreadcrumb()
	breadcrumbLine := lipgloss.NewStyle().
		Foreground(lipgloss.Color("240")).
		Width(m.width).
		Align(lipgloss.Center).
		Render(breadcrumb)
	s.WriteString(breadcrumbLine)
	s.WriteString("\n")

	// Separator items for scheduled/unscheduled
	scheduled, unscheduled := m.splitItemsBySchedule(status.Items)

	// Panel widths (60% timeline, 40% unscheduled)
	timelineWidth := m.width * 60 / 100
	unscheduledWidth := m.width - timelineWidth - 1 // -1 for border

	// Available height for content
	// Header overhead: view mode (1) + header (1) + breadcrumb (1) + panel headers (1) + 4 newlines = 8
	availableHeight := m.height - 8

	// Render panel headers
	timelineHeaderStyle := lipgloss.NewStyle().
		Width(timelineWidth).
		Foreground(lipgloss.Color("#7D56F4")).
		Bold(true)
	if m.dayViewFocus == DayViewFocusTimeline {
		timelineHeaderStyle = timelineHeaderStyle.Background(lipgloss.Color("#3D3D3D"))
	}

	unscheduledHeaderStyle := lipgloss.NewStyle().
		Width(unscheduledWidth).
		Foreground(lipgloss.Color("#7D56F4")).
		Bold(true)
	if m.dayViewFocus == DayViewFocusUnscheduled {
		unscheduledHeaderStyle = unscheduledHeaderStyle.Background(lipgloss.Color("#3D3D3D"))
	}

	timelineHeader := timelineHeaderStyle.Render(" TIMELINE")
	unscheduledHeader := unscheduledHeaderStyle.Render(fmt.Sprintf(" UNSCHEDULED (%d)", len(unscheduled)))

	panelHeaders := lipgloss.JoinHorizontal(lipgloss.Top,
		timelineHeader,
		lipgloss.NewStyle().Foreground(lipgloss.Color("238")).Render("│"),
		unscheduledHeader,
	)
	s.WriteString(panelHeaders)
	s.WriteString("\n")

	// Separator
	separator := lipgloss.NewStyle().
		Foreground(lipgloss.Color("238")).
		Render(strings.Repeat("─", timelineWidth) + "┼" + strings.Repeat("─", unscheduledWidth))
	s.WriteString(separator)
	s.WriteString("\n")

	// Content height after panel headers and separator
	contentHeight := availableHeight - 2

	// Render timeline panel
	timelinePanel := m.renderTimelinePanel(scheduled, timelineWidth, contentHeight)

	// Render unscheduled panel
	unscheduledPanel := m.renderUnscheduledPanel(unscheduled, unscheduledWidth, contentHeight)

	// Join panels horizontally
	borderStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("238")).
		Height(contentHeight)
	panels := lipgloss.JoinHorizontal(lipgloss.Top,
		timelinePanel,
		borderStyle.Render(strings.Repeat("│\n", contentHeight)),
		unscheduledPanel,
	)
	s.WriteString(panels)

	return s.String()
}

// renderDayViewList renders the day view in List Mode (card-style)
func (m *Model) renderDayViewList() string {
	var s strings.Builder

	dateKey := m.cursor.Format("2006-01-02")
	status := m.todoStatus[dateKey]

	// View mode indicator line
	viewModeIndicator := lipgloss.NewStyle().
		Width(m.width).
		Align(lipgloss.Center).
		Render(m.renderViewModeIndicator())
	s.WriteString(viewModeIndicator)
	s.WriteString("\n")

	// Main Header with full date
	header := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#FAFAFA")).
		Background(lipgloss.Color("#7D56F4")).
		Padding(0, 1).
		Width(m.width).
		Align(lipgloss.Center).
		Render(m.cursor.Format("Monday, January 2, 2006"))
	s.WriteString(header)
	s.WriteString("\n")

	// Breadcrumb
	breadcrumb := m.renderBreadcrumb()
	breadcrumbLine := lipgloss.NewStyle().
		Foreground(lipgloss.Color("240")).
		Width(m.width).
		Align(lipgloss.Center).
		Render(breadcrumb)
	s.WriteString(breadcrumbLine)
	s.WriteString("\n")

	// Available height for content
	// Header overhead: view mode (1) + header (1) + breadcrumb (1) + 3 newlines = 6
	availableHeight := m.height - 6

	// Sort items: by priority (high to low), then completed at bottom
	items := m.sortItemsForListMode(status.Items)

	// Render task cards
	content := m.renderListModeContent(items, m.width, availableHeight)
	s.WriteString(content)

	return s.String()
}

// sortItemsForListMode sorts items by priority (high to low), completed at bottom
func (m *Model) sortItemsForListMode(items []TodoItem) []TodoItem {
	if len(items) == 0 {
		return items
	}

	// Separate completed and incomplete
	var incomplete, completed []TodoItem
	for _, item := range items {
		if item.Complete {
			completed = append(completed, item)
		} else {
			incomplete = append(incomplete, item)
		}
	}

	// Sort incomplete by priority (high to low)
	for i := 0; i < len(incomplete)-1; i++ {
		for j := i + 1; j < len(incomplete); j++ {
			if incomplete[j].Priority > incomplete[i].Priority {
				incomplete[i], incomplete[j] = incomplete[j], incomplete[i]
			}
		}
	}

	// Sort completed by priority too
	for i := 0; i < len(completed)-1; i++ {
		for j := i + 1; j < len(completed); j++ {
			if completed[j].Priority > completed[i].Priority {
				completed[i], completed[j] = completed[j], completed[i]
			}
		}
	}

	// Combine: incomplete first, then completed
	return append(incomplete, completed...)
}

// renderListModeContent renders the task cards for list mode
func (m *Model) renderListModeContent(items []TodoItem, width, height int) string {
	var lines []string

	// Empty state
	if len(items) == 0 {
		emptyStyle := lipgloss.NewStyle().
			Foreground(lipgloss.Color("240")).
			Italic(true)
		hintStyle := lipgloss.NewStyle().
			Foreground(lipgloss.Color("#7D56F4"))

		// Center the empty message
		emptyLines := []string{
			"",
			"",
			emptyStyle.Render("No tasks for today"),
			"",
			hintStyle.Render("Press 'a' to add a task"),
		}

		centeredContent := lipgloss.NewStyle().
			Width(width).
			Height(height).
			Align(lipgloss.Center, lipgloss.Center).
			Render(strings.Join(emptyLines, "\n"))

		return centeredContent
	}

	// Card styling
	cardWidth := width - 4 // Padding on sides
	if cardWidth > 80 {
		cardWidth = 80 // Max width for readability
	}

	// Calculate how many lines each card takes (title + 2-3 desc lines + border)
	linesPerCard := 5

	for i, item := range items {
		if len(lines) >= height-linesPerCard {
			remaining := len(items) - i
			lines = append(lines, "")
			lines = append(lines, lipgloss.NewStyle().
				Foreground(lipgloss.Color("245")).
				Italic(true).
				Render(fmt.Sprintf("  +%d more tasks...", remaining)))
			break
		}

		// Selection indicator
		isSelected := i == m.selectedListItem

		// Build card content
		cardContent := m.renderTaskCard(item, cardWidth-4, isSelected)
		lines = append(lines, cardContent)
		lines = append(lines, "") // Space between cards
	}

	// Pad to fill height
	for len(lines) < height {
		lines = append(lines, "")
	}

	return strings.Join(lines, "\n")
}

// renderTaskCard renders a single task as a card
func (m *Model) renderTaskCard(item TodoItem, contentWidth int, isSelected bool) string {
	// Checkbox
	checkbox := "☐"
	if item.Complete {
		checkbox = "☑"
	}

	// Priority indicator
	priorityStr := ""
	priorityColor := lipgloss.Color("#FFFFFF")
	switch item.Priority {
	case 3:
		priorityStr = " !!!"
		priorityColor = lipgloss.Color("#FF6B6B")
	case 2:
		priorityStr = " !!"
		priorityColor = lipgloss.Color("#FFB86C")
	case 1:
		priorityStr = " !"
		priorityColor = lipgloss.Color("#8BE9FD")
	}

	// Title styling
	titleStyle := lipgloss.NewStyle()
	if item.Complete {
		titleStyle = titleStyle.Foreground(lipgloss.Color("#666")).Strikethrough(true)
	} else {
		titleStyle = titleStyle.Foreground(priorityColor).Bold(true)
	}

	// Title line
	title := item.Title
	if len(title) > contentWidth-6 {
		title = title[:contentWidth-9] + "…"
	}
	titleLine := checkbox + " " + titleStyle.Render(title) + priorityStr

	// Description lines (up to 2 lines)
	var descLines []string
	if item.Desc != "" {
		descStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("245"))
		if item.Complete {
			descStyle = descStyle.Foreground(lipgloss.Color("#555"))
		}

		desc := item.Desc
		// Split into multiple lines if needed
		words := strings.Fields(desc)
		var currentLine string
		for _, word := range words {
			testLine := currentLine
			if testLine != "" {
				testLine += " "
			}
			testLine += word

			if len(testLine) > contentWidth-4 {
				if currentLine != "" {
					descLines = append(descLines, descStyle.Render("    "+currentLine))
					if len(descLines) >= 2 {
						// Add ellipsis to last line
						lastIdx := len(descLines) - 1
						descLines[lastIdx] = descStyle.Render("    " + currentLine + "...")
						break
					}
				}
				currentLine = word
			} else {
				currentLine = testLine
			}
		}
		if currentLine != "" && len(descLines) < 2 {
			descLines = append(descLines, descStyle.Render("    "+currentLine))
		}
	}

	// Build card content
	var cardLines []string
	cardLines = append(cardLines, titleLine)
	cardLines = append(cardLines, descLines...)

	cardContent := strings.Join(cardLines, "\n")

	// Card border styling
	borderColor := lipgloss.Color("#444")
	if isSelected {
		borderColor = lipgloss.Color("#7D56F4")
	} else if !item.Complete {
		switch item.Priority {
		case 3:
			borderColor = lipgloss.Color("#FF6B6B")
		case 2:
			borderColor = lipgloss.Color("#FFB86C")
		case 1:
			borderColor = lipgloss.Color("#8BE9FD")
		}
	}

	cardStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(borderColor).
		Padding(0, 1).
		Width(contentWidth + 4).
		MarginLeft(2)

	if isSelected {
		cardStyle = cardStyle.BorderForeground(lipgloss.Color("#7D56F4")).Bold(true)
	}

	return cardStyle.Render(cardContent)
}

// splitItemsBySchedule separates scheduled and unscheduled items
func (m *Model) splitItemsBySchedule(items []TodoItem) (scheduled, unscheduled []TodoItem) {
	for _, item := range items {
		if item.StartTime != "" {
			scheduled = append(scheduled, item)
		} else {
			unscheduled = append(unscheduled, item)
		}
	}
	return
}

// renderTimelinePanel renders the left panel with timeline
func (m *Model) renderTimelinePanel(scheduled []TodoItem, width, height int) string {
	var lines []string

	// Calculate time slots
	startHour, startMin := m.parseTime(m.timelineConfig.DayStart)
	endHour, endMin := m.parseTime(m.timelineConfig.DayEnd)

	startMinutes := startHour*60 + startMin
	endMinutes := endHour*60 + endMin
	totalSlots := (endMinutes - startMinutes) / m.timelineConfig.SlotMinutes

	// Calculate lines per slot based on available height
	linesPerSlot := height / totalSlots
	if linesPerSlot < 1 {
		linesPerSlot = 1
	}

	// Build a map of scheduled items by start time
	scheduledMap := make(map[string][]TodoItem)
	for _, item := range scheduled {
		scheduledMap[item.StartTime] = append(scheduledMap[item.StartTime], item)
	}

	timeStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("245"))
	cursorStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#7D56F4")).Bold(true)
	emptyLineStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("238"))

	for slot := 0; slot < totalSlots && len(lines) < height; slot++ {
		slotMinutes := startMinutes + slot*m.timelineConfig.SlotMinutes
		hour := slotMinutes / 60
		minute := slotMinutes % 60
		timeStr := fmt.Sprintf("%02d:%02d", hour, minute)

		// Check if cursor is on this slot
		isCursor := m.dayViewFocus == DayViewFocusTimeline && slot == m.timelineCursor

		// Check for scheduled items at this time
		items, hasItem := scheduledMap[timeStr]

		if hasItem && len(items) > 0 {
			// Render scheduled task block
			item := items[0] // For simplicity, show first item
			blockLines := m.renderTaskBlock(item, width-7, linesPerSlot)

			for i, blockLine := range blockLines {
				var line string
				if i == 0 {
					if isCursor {
						line = cursorStyle.Render("▶"+timeStr) + " " + blockLine
					} else {
						line = timeStyle.Render(" "+timeStr) + " " + blockLine
					}
				} else {
					line = "       " + blockLine
				}
				lines = append(lines, line)
				if len(lines) >= height {
					break
				}
			}
		} else {
			// Empty slot
			for i := 0; i < linesPerSlot && len(lines) < height; i++ {
				var line string
				if i == 0 {
					if isCursor {
						line = cursorStyle.Render("▶"+timeStr) + " " + emptyLineStyle.Render(strings.Repeat("┄", width-8))
					} else {
						line = timeStyle.Render(" "+timeStr) + " " + emptyLineStyle.Render(strings.Repeat("┄", width-8))
					}
				} else {
					line = strings.Repeat(" ", width)
				}
				lines = append(lines, line)
			}
		}
	}

	// Pad to fill height
	for len(lines) < height {
		lines = append(lines, strings.Repeat(" ", width))
	}

	return lipgloss.NewStyle().Width(width).Height(height).Render(strings.Join(lines, "\n"))
}

// renderTaskBlock renders a scheduled task as a visual block
func (m *Model) renderTaskBlock(item TodoItem, width, minHeight int) []string {
	// Calculate block height based on duration
	blockHeight := minHeight
	if item.StartTime != "" && item.EndTime != "" {
		startH, startM := m.parseTime(item.StartTime)
		endH, endM := m.parseTime(item.EndTime)
		durationMin := (endH*60 + endM) - (startH*60 + startM)
		slots := durationMin / m.timelineConfig.SlotMinutes
		if slots > 1 {
			blockHeight = slots * minHeight
		}
	}

	var lines []string

	// Checkbox
	checkbox := "☐"
	if item.Complete {
		checkbox = "☑"
	}

	// Priority indicator
	priorityStr := ""
	switch item.Priority {
	case 3:
		priorityStr = " !!!"
	case 2:
		priorityStr = " !!"
	case 1:
		priorityStr = " !"
	}

	// Style based on completion/priority
	blockStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#7D56F4")).
		Width(width - 2)

	if item.Complete {
		blockStyle = blockStyle.Foreground(lipgloss.Color("#666"))
	} else {
		switch item.Priority {
		case 3:
			blockStyle = blockStyle.BorderForeground(lipgloss.Color("#FF6B6B"))
		case 2:
			blockStyle = blockStyle.BorderForeground(lipgloss.Color("#FFB86C"))
		case 1:
			blockStyle = blockStyle.BorderForeground(lipgloss.Color("#8BE9FD"))
		}
	}

	// Title line
	title := item.Title
	if len(title) > width-6 {
		title = title[:width-9] + "…"
	}
	titleLine := checkbox + " " + title + priorityStr

	// Description (if fits)
	content := titleLine
	if item.Desc != "" && blockHeight > 2 {
		desc := item.Desc
		if len(desc) > width-4 {
			desc = desc[:width-7] + "…"
		}
		content = titleLine + "\n" + desc
	}

	block := blockStyle.Render(content)
	blockLines := strings.Split(block, "\n")

	for _, line := range blockLines {
		lines = append(lines, line)
	}

	return lines
}

// renderUnscheduledPanel renders the right panel with unscheduled tasks
func (m *Model) renderUnscheduledPanel(unscheduled []TodoItem, width, height int) string {
	var lines []string

	if len(unscheduled) == 0 {
		emptyStyle := lipgloss.NewStyle().
			Foreground(lipgloss.Color("240")).
			Italic(true)
		hintStyle := lipgloss.NewStyle().
			Foreground(lipgloss.Color("#7D56F4"))
		lines = append(lines, "")
		lines = append(lines, emptyStyle.Render("  All tasks scheduled!"))
		lines = append(lines, "")
		lines = append(lines, hintStyle.Render("  Press 'a' to add"))
		lines = append(lines, hintStyle.Render("  more tasks"))
	} else {
		for i, item := range unscheduled {
			if len(lines) >= height-1 {
				remaining := len(unscheduled) - i
				lines = append(lines, fmt.Sprintf("  +%d more...", remaining))
				break
			}

			// Selection indicator
			selector := "  "
			if m.dayViewFocus == DayViewFocusUnscheduled && i == m.selectedUnscheduled {
				selector = "> "
			}

			// Checkbox
			checkbox := "☐"
			if item.Complete {
				checkbox = "☑"
			}

			// Priority
			priorityStr := ""
			switch item.Priority {
			case 3:
				priorityStr = " !!!"
			case 2:
				priorityStr = " !!"
			case 1:
				priorityStr = " !"
			}

			// Style
			titleStyle := lipgloss.NewStyle()
			if item.Complete {
				titleStyle = titleStyle.Foreground(lipgloss.Color("#666")).Strikethrough(true)
			} else {
				switch item.Priority {
				case 3:
					titleStyle = titleStyle.Foreground(lipgloss.Color("#FF6B6B")).Bold(true)
				case 2:
					titleStyle = titleStyle.Foreground(lipgloss.Color("#FFB86C"))
				case 1:
					titleStyle = titleStyle.Foreground(lipgloss.Color("#8BE9FD"))
				}
			}

			// Title line
			title := item.Title
			maxTitleWidth := width - 10
			if len(title) > maxTitleWidth {
				title = title[:maxTitleWidth-1] + "…"
			}

			line := selector + checkbox + " " + titleStyle.Render(title) + priorityStr
			lines = append(lines, line)

			// Description on next line
			if item.Desc != "" && len(lines) < height-1 {
				desc := item.Desc
				maxDescWidth := width - 6
				if len(desc) > maxDescWidth {
					desc = desc[:maxDescWidth-1] + "…"
				}
				descStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("245"))
				lines = append(lines, "    "+descStyle.Render(desc))
			}
		}
	}

	// Pad to fill height
	for len(lines) < height {
		lines = append(lines, "")
	}

	return lipgloss.NewStyle().Width(width).Height(height).Render(strings.Join(lines, "\n"))
}

// parseTime parses "HH:MM" format and returns hour, minute
func (m *Model) parseTime(timeStr string) (int, int) {
	if len(timeStr) < 5 {
		return 0, 0
	}
	var hour, minute int
	fmt.Sscanf(timeStr, "%d:%d", &hour, &minute)
	return hour, minute
}
