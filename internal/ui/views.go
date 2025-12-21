package ui

import (
	"fmt"

	"ctchen222/chronoflow/internal/domain"

	"github.com/charmbracelet/lipgloss"
)

// ViewRenderer handles all view rendering logic
type ViewRenderer struct {
	width  int
	height int
}

// NewViewRenderer creates a new ViewRenderer
func NewViewRenderer() *ViewRenderer {
	return &ViewRenderer{}
}

// SetSize updates the viewport dimensions
func (v *ViewRenderer) SetSize(width, height int) {
	v.width = width
	v.height = height
}

// Width returns the current width
func (v *ViewRenderer) Width() int {
	return v.width
}

// Height returns the current height
func (v *ViewRenderer) Height() int {
	return v.height
}

// RenderMain renders the main viewing state with calendar and todo panels
func (v *ViewRenderer) RenderMain(state MainViewState) string {
	panelHeight := v.height - 1 // reserve 1 line for help bar
	calendarWidth := int(float64(v.width) * 0.7)
	todoWidth := v.width - calendarWidth

	// Inner content dimensions (subtract 2 for border on each side)
	calInnerW := calendarWidth - 2
	calInnerH := panelHeight - 2
	todoInnerW := todoWidth - 2
	todoInnerH := panelHeight - 2

	// Get content and place it in fixed-size container
	calContent := lipgloss.Place(calInnerW, calInnerH, lipgloss.Left, lipgloss.Top, state.CalendarView)
	todoContent := lipgloss.Place(todoInnerW, todoInnerH, lipgloss.Left, lipgloss.Top, state.TodoView)

	// Border colors based on focus
	calBorderColor := lipgloss.Color("#444")
	todoBorderColor := lipgloss.Color("#444")
	if state.Focus == FocusCalendar {
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

	return lipgloss.JoinHorizontal(lipgloss.Top, calView, todoView)
}

// RenderEditing renders the editing modal
func (v *ViewRenderer) RenderEditing(state EditingState) string {
	var accentColor, headerIcon string
	if state.IsNew {
		accentColor = "#50FA7B" // Green for new
		headerIcon = "+"
	} else {
		accentColor = "#8BE9FD" // Cyan for edit
		headerIcon = "~"
	}

	// Header
	headerText := "New Todo"
	if !state.IsNew {
		headerText = "Edit Todo"
	}
	headerStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color(accentColor)).
		MarginBottom(1)
	header := headerStyle.Render(headerIcon + "  " + headerText)

	// Date
	dateText := state.Date.Format("Mon, Jan 2, 2006")
	dateStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#888")).
		MarginBottom(1)
	date := dateStyle.Render(dateText)

	// Title input with label
	titleLabelStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#888"))
	titleBorderColor := lipgloss.Color("#444")
	if state.Focus == FocusTitle {
		titleBorderColor = lipgloss.Color(accentColor)
	}
	titleInputStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(titleBorderColor).
		Padding(0, 1).
		Width(60)
	titleSection := lipgloss.JoinVertical(lipgloss.Left,
		titleLabelStyle.Render("Title"),
		titleInputStyle.Render(state.TitleView),
	)

	// Description input with label
	descLabelStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#888")).
		MarginTop(1)
	descBorderColor := lipgloss.Color("#444")
	if state.Focus == FocusDesc {
		descBorderColor = lipgloss.Color(accentColor)
	}
	descInputStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(descBorderColor).
		Padding(0, 1).
		Width(60)
	descSection := lipgloss.JoinVertical(lipgloss.Left,
		descLabelStyle.Render("Description (optional)"),
		descInputStyle.Render(state.DescView),
	)

	// Priority selector
	prioritySection := v.renderPrioritySelector(state.Priority, accentColor)

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
	bgHeight := v.height - 1 // exclude help bar
	return lipgloss.Place(v.width, bgHeight,
		lipgloss.Center, lipgloss.Center,
		modalBox,
		lipgloss.WithWhitespaceChars(" "),
		lipgloss.WithWhitespaceForeground(lipgloss.Color("#333")))
}

// renderPrioritySelector renders the priority selection row
func (v *ViewRenderer) renderPrioritySelector(selected domain.Priority, accentColor string) string {
	priorityLabelStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#888")).
		MarginTop(1)

	priorityOptions := []struct {
		level domain.Priority
		label string
		color string
	}{
		{domain.PriorityNone, "None", "#666"},
		{domain.PriorityLow, "Low", "#8BE9FD"},
		{domain.PriorityMedium, "Medium", "#FFB86C"},
		{domain.PriorityHigh, "High", "#FF6B6B"},
	}

	var priorityItems []string
	for _, opt := range priorityOptions {
		style := lipgloss.NewStyle().Foreground(lipgloss.Color(opt.color))
		label := opt.label
		if opt.level == selected {
			label = "[" + label + "]"
			style = style.Bold(true)
		} else {
			label = " " + label + " "
		}
		priorityItems = append(priorityItems, style.Render(label))
	}

	priorityRow := lipgloss.JoinHorizontal(lipgloss.Center, priorityItems...)
	return lipgloss.JoinVertical(lipgloss.Left,
		priorityLabelStyle.Render("Priority (Ctrl+0/1/2/3)"),
		priorityRow,
	)
}

// RenderConfirmDelete renders the delete confirmation modal
func (v *ViewRenderer) RenderConfirmDelete(state DeleteState) string {
	// Truncate title if too long
	title := state.Title
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
	bgHeight := v.height - 1
	return lipgloss.Place(v.width, bgHeight,
		lipgloss.Center, lipgloss.Center,
		modalBox,
		lipgloss.WithWhitespaceChars(" "),
		lipgloss.WithWhitespaceForeground(lipgloss.Color("#333")))
}

// RenderSearching renders the search modal
func (v *ViewRenderer) RenderSearching(state SearchState) string {
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
		inputStyle.Render(state.InputView),
	)

	// Results
	resultsContent := v.renderSearchResults(state, accentColor)

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
	bgHeight := v.height - 1
	return lipgloss.Place(v.width, bgHeight,
		lipgloss.Center, lipgloss.Center,
		modalBox,
		lipgloss.WithWhitespaceChars(" "),
		lipgloss.WithWhitespaceForeground(lipgloss.Color("#333")))
}

// renderSearchResults renders the search results list
func (v *ViewRenderer) renderSearchResults(state SearchState, accentColor string) string {
	if len(state.Results) == 0 {
		if state.InputValue == "" {
			return lipgloss.NewStyle().
				Foreground(lipgloss.Color("#666")).
				Italic(true).
				Render("Type to search...")
		}
		return lipgloss.NewStyle().
			Foreground(lipgloss.Color("#666")).
			Italic(true).
			Render("No results found")
	}

	var resultLines []string
	maxResults := 8 // Show max 8 results
	start := 0
	if state.SelectedIdx >= maxResults {
		start = state.SelectedIdx - maxResults + 1
	}
	end := start + maxResults
	if end > len(state.Results) {
		end = len(state.Results)
	}

	for i := start; i < end; i++ {
		r := state.Results[i]
		dateStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#7D56F4"))
		titleStyle := lipgloss.NewStyle()

		prefix := "  "
		if i == state.SelectedIdx {
			prefix = "> "
			titleStyle = titleStyle.Bold(true).Foreground(lipgloss.Color(accentColor))
		}

		// Show completion status
		status := "☐"
		if r.Todo.Complete {
			status = "☑"
			titleStyle = titleStyle.Foreground(lipgloss.Color("#666"))
		}

		line := prefix + dateStyle.Render(r.DateKey) + " " + status + " " + titleStyle.Render(r.Todo.Title)
		resultLines = append(resultLines, line)
	}

	resultsHeader := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#888")).
		MarginTop(1).
		Render(fmt.Sprintf("Results (%d found)", len(state.Results)))

	return lipgloss.JoinVertical(lipgloss.Left,
		append([]string{resultsHeader}, resultLines...)...)
}

// RenderHelpBar renders the help bar at the bottom
func (v *ViewRenderer) RenderHelpBar(state AppState, focus AppFocus) string {
	keyStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#7D56F4")).
		Bold(true)
	descStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#666"))
	sepStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#444"))

	sep := sepStyle.Render(" │ ")

	var keys string
	switch state {
	case StateViewing:
		if focus == FocusCalendar {
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
	case StateEditing:
		keys = keyStyle.Render("Tab") + descStyle.Render(" switch field") + sep +
			keyStyle.Render("Enter") + descStyle.Render(" save") + sep +
			keyStyle.Render("Esc") + descStyle.Render(" cancel")
	case StateConfirmingDelete:
		keys = keyStyle.Render("y/Enter") + descStyle.Render(" confirm") + sep +
			keyStyle.Render("n/Esc") + descStyle.Render(" cancel")
	case StateSearching:
		keys = keyStyle.Render("Up/Down") + descStyle.Render(" navigate") + sep +
			keyStyle.Render("Enter") + descStyle.Render(" go to") + sep +
			keyStyle.Render("Esc") + descStyle.Render(" cancel")
	}

	return lipgloss.NewStyle().
		Width(v.width).
		Align(lipgloss.Center).
		Background(lipgloss.Color("#1a1a1a")).
		Padding(0, 1).
		Render(keys)
}
