package ui

import (
	"fmt"
	"strings"

	"ctchen222/chronoflow/internal/domain"

	"github.com/charmbracelet/lipgloss"
)

// Constants for responsive modal
const (
	ModalWidthPercent = 0.70 // 70% of terminal width
	ModalMinWidth     = 50   // Minimum modal width
	ModalMaxWidth     = 140  // Maximum modal width
	MinPreviewWidth   = 30   // Minimum width to show preview
)

// Minimum terminal size constants
const (
	MinTerminalWidth  = 80
	MinTerminalHeight = 24
)

// ModalDimensions holds calculated modal dimensions
type ModalDimensions struct {
	TotalWidth   int
	EditorWidth  int
	PreviewWidth int
	InputWidth   int
	ShowPreview  bool
}

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

// CalculateModalDimensions calculates responsive modal dimensions
func (v *ViewRenderer) CalculateModalDimensions(previewEnabled bool) ModalDimensions {
	// Calculate base modal width (70% of terminal)
	modalWidth := int(float64(v.width) * ModalWidthPercent)

	// Apply min/max bounds
	if modalWidth < ModalMinWidth {
		modalWidth = ModalMinWidth
	}
	if modalWidth > ModalMaxWidth {
		modalWidth = ModalMaxWidth
	}

	// Account for modal padding and borders (2 padding + 2 border each side = 8)
	innerWidth := modalWidth - 8

	dims := ModalDimensions{
		TotalWidth:  modalWidth,
		ShowPreview: false,
	}

	// Determine if we can show split view
	// Need at least MinPreviewWidth for each pane plus 3 for separator
	minSplitWidth := (MinPreviewWidth * 2) + 3

	if previewEnabled && innerWidth >= minSplitWidth {
		dims.ShowPreview = true
		// Calculate split (50/50)
		halfWidth := (innerWidth - 3) / 2 // -3 for separator
		dims.EditorWidth = halfWidth
		dims.PreviewWidth = innerWidth - halfWidth - 3
		dims.InputWidth = halfWidth - 4 // Account for input padding/border
	} else {
		// Full-width editor (no preview)
		dims.EditorWidth = innerWidth
		dims.PreviewWidth = 0
		dims.InputWidth = innerWidth - 4
	}

	return dims
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

// RenderEditing renders the editing modal with optional markdown preview
func (v *ViewRenderer) RenderEditing(state EditingState) string {
	var accentColor, headerIcon string
	if state.IsNew {
		accentColor = "#50FA7B" // Green for new
		headerIcon = "+"
	} else {
		accentColor = "#8BE9FD" // Cyan for edit
		headerIcon = "~"
	}

	// Calculate responsive dimensions
	dims := v.CalculateModalDimensions(state.PreviewEnabled)

	// Header and date on same line
	headerText := "New Todo"
	if !state.IsNew {
		headerText = "Edit Todo"
	}
	headerStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color(accentColor))
	headerLeft := headerStyle.Render(headerIcon + "  " + headerText)

	dateText := state.Date.Format("Mon, Jan 2, 2006")
	dateStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#888"))
	dateRight := dateStyle.Render(dateText)

	// Calculate spacing between header and date
	innerWidth := dims.TotalWidth - 8 // Account for modal padding and border
	headerWidth := lipgloss.Width(headerLeft)
	dateWidth := lipgloss.Width(dateRight)
	spacerWidth := innerWidth - headerWidth - dateWidth
	if spacerWidth < 1 {
		spacerWidth = 1
	}
	spacer := lipgloss.NewStyle().Width(spacerWidth).Render("")
	headerLine := lipgloss.JoinHorizontal(lipgloss.Center, headerLeft, spacer, dateRight)

	// Build editor column
	editorColumn := v.renderEditorColumn(state, accentColor, dims.InputWidth)

	var contentArea string
	if dims.ShowPreview {
		// Fixed height for both columns to ensure alignment
		boxHeight := 15
		editorInnerWidth := dims.EditorWidth
		previewInnerWidth := dims.PreviewWidth

		// Preview column with dashed border
		previewColumn := v.renderPreviewColumn(state.PreviewContent, previewInnerWidth-4, accentColor)

		// Place content in fixed-size containers
		editorContent := lipgloss.Place(editorInnerWidth, boxHeight, lipgloss.Left, lipgloss.Top, editorColumn)
		previewContent := lipgloss.Place(previewInnerWidth, boxHeight, lipgloss.Left, lipgloss.Top, previewColumn)

		// Join columns horizontally (no extra box borders)
		contentArea = lipgloss.JoinHorizontal(
			lipgloss.Top,
			editorContent,
			"  ", // Gap between columns
			previewContent,
		)
	} else {
		// Single column layout
		contentArea = editorColumn
	}

	// Integrated help bar
	helpBar := v.renderModalHelpBar(innerWidth)

	// Combine all parts
	modalContent := lipgloss.JoinVertical(lipgloss.Left,
		headerLine,
		"", // Empty line for spacing
		contentArea,
		"", // Empty line before help
		helpBar,
	)

	// Modal box with background
	modalBox := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color(accentColor)).
		Padding(1, 2).
		Width(dims.TotalWidth).
		Render(modalContent)

	// Center the modal in the available space
	bgHeight := v.height - 1 // exclude help bar
	return lipgloss.Place(v.width, bgHeight,
		lipgloss.Center, lipgloss.Center,
		modalBox,
		lipgloss.WithWhitespaceChars(" "),
		lipgloss.WithWhitespaceForeground(lipgloss.Color("#333")))
}

// renderEditorColumn renders the editor inputs column
func (v *ViewRenderer) renderEditorColumn(state EditingState, accentColor string, inputWidth int) string {
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
		Width(inputWidth)
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
		Width(inputWidth)
	descSection := lipgloss.JoinVertical(lipgloss.Left,
		descLabelStyle.Render("Description (Markdown supported)"),
		descInputStyle.Render(state.DescView),
	)

	// Priority selector
	prioritySection := v.renderPrioritySelector(state.Priority, accentColor)

	return lipgloss.JoinVertical(lipgloss.Left,
		titleSection,
		descSection,
		prioritySection,
	)
}

// renderPreviewColumn renders the markdown preview column with dashed border
func (v *ViewRenderer) renderPreviewColumn(previewContent string, width int, accentColor string) string {
	// Preview header
	previewHeaderStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#7D56F4")).
		Bold(true)
	previewHeader := previewHeaderStyle.Render("Preview")

	// Dashed border style for read-only indicator
	dashedBorder := lipgloss.Border{
		Top:         "─",
		Bottom:      "─",
		Left:        "┊",
		Right:       "┊",
		TopLeft:     "┌",
		TopRight:    "┐",
		BottomLeft:  "└",
		BottomRight: "┘",
	}

	// Preview content area with dashed border
	previewStyle := lipgloss.NewStyle().
		Border(dashedBorder).
		BorderForeground(lipgloss.Color("#555")).
		Padding(0, 1).
		Width(width).
		Height(10)

	previewArea := previewStyle.Render(previewContent)

	return lipgloss.JoinVertical(lipgloss.Left,
		previewHeader,
		previewArea,
	)
}

// renderPrioritySelector renders the priority selection row with radio-style buttons
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
		var indicator string
		if opt.level == selected {
			indicator = "●"
			style = style.Bold(true)
		} else {
			indicator = "○"
		}
		item := style.Render(indicator + " " + opt.label)
		priorityItems = append(priorityItems, item+"   ") // Spacing between items
	}

	priorityRow := lipgloss.JoinHorizontal(lipgloss.Center, priorityItems...)
	return lipgloss.JoinVertical(lipgloss.Left,
		priorityLabelStyle.Render("Priority (←/→ or Ctrl+0/1/2/3)"),
		priorityRow,
	)
}

// renderModalHelpBar renders the integrated help bar for the editing modal
func (v *ViewRenderer) renderModalHelpBar(width int) string {
	// Separator line
	separator := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#444")).
		Width(width).
		Render(strings.Repeat("─", width))

	// Help text
	helpStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#666"))
	keyStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#888"))

	helpItems := []string{
		keyStyle.Render("Tab") + helpStyle.Render(" switch"),
		keyStyle.Render("Ctrl+P") + helpStyle.Render(" preview"),
		keyStyle.Render("Enter") + helpStyle.Render(" save"),
		keyStyle.Render("Esc") + helpStyle.Render(" cancel"),
	}

	helpText := strings.Join(helpItems, helpStyle.Render("  │  "))

	return lipgloss.JoinVertical(lipgloss.Left,
		separator,
		helpText,
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
// highlightMatch highlights the matching substring in text with the accent color
func highlightMatch(text, query string, baseStyle lipgloss.Style, highlightColor string) string {
	if query == "" {
		return baseStyle.Render(text)
	}

	lowerText := strings.ToLower(text)
	lowerQuery := strings.ToLower(query)
	idx := strings.Index(lowerText, lowerQuery)

	if idx == -1 {
		return baseStyle.Render(text)
	}

	highlightStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color(highlightColor)).
		Bold(true)

	before := text[:idx]
	match := text[idx : idx+len(query)]
	after := text[idx+len(query):]

	return baseStyle.Render(before) + highlightStyle.Render(match) + baseStyle.Render(after)
}

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

	// Highlight color for matches
	highlightColor := "#7D56F4" // Accent color

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

		// Highlight matching text in title
		highlightedTitle := highlightMatch(r.Todo.Title, state.InputValue, titleStyle, highlightColor)

		line := prefix + dateStyle.Render(r.DateKey) + " " + status + " " + highlightedTitle
		resultLines = append(resultLines, line)
	}

	resultsHeader := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#888")).
		MarginTop(1).
		Render(fmt.Sprintf("Results (%d found)", len(state.Results)))

	return lipgloss.JoinVertical(lipgloss.Left,
		append([]string{resultsHeader}, resultLines...)...)
}

// RenderHelp renders the keyboard reference modal
func (v *ViewRenderer) RenderHelp() string {
	accentColor := "#7D56F4"

	// Header
	headerStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color(accentColor))
	header := headerStyle.Render("Keyboard Reference")

	// Category header style
	categoryStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#FFFFFF")).
		MarginTop(1)

	// Underline style
	underlineStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#666"))

	// Key style
	keyStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color(accentColor)).
		Bold(true).
		Width(12)

	// Description style
	descStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#AAAAAA"))

	// Build keyboard shortcuts content
	// Navigation column
	navHeader := categoryStyle.Render("Navigation")
	navUnderline := underlineStyle.Render("──────────")
	navItems := []string{
		keyStyle.Render("h/j/k/l") + descStyle.Render("Move cursor"),
		keyStyle.Render("b/n") + descStyle.Render("Prev/next month"),
		keyStyle.Render("w") + descStyle.Render("Toggle week view"),
		keyStyle.Render("d") + descStyle.Render("Day view"),
		keyStyle.Render("m") + descStyle.Render("Month view"),
		keyStyle.Render("t") + descStyle.Render("Jump to today"),
		keyStyle.Render("g") + descStyle.Render("Go to date"),
	}
	navColumn := lipgloss.JoinVertical(lipgloss.Left,
		append([]string{navHeader, navUnderline}, navItems...)...)

	// Todo Actions column
	todoHeader := categoryStyle.Render("Todo Actions")
	todoUnderline := underlineStyle.Render("────────────")
	todoItems := []string{
		keyStyle.Render("Space/x") + descStyle.Render("Toggle done"),
		keyStyle.Render("a") + descStyle.Render("Add todo"),
		keyStyle.Render("e/Enter") + descStyle.Render("Edit todo"),
		keyStyle.Render("d") + descStyle.Render("Delete todo"),
		keyStyle.Render("1/2/3/0") + descStyle.Render("Set priority"),
		keyStyle.Render("J/K") + descStyle.Render("Reorder"),
		"",
	}
	todoColumn := lipgloss.JoinVertical(lipgloss.Left,
		append([]string{todoHeader, todoUnderline}, todoItems...)...)

	// Edit Mode column
	editHeader := categoryStyle.Render("Edit Mode")
	editUnderline := underlineStyle.Render("─────────")
	editItems := []string{
		keyStyle.Render("Tab") + descStyle.Render("Switch field"),
		keyStyle.Render("Ctrl+P") + descStyle.Render("Preview markdown"),
		keyStyle.Render("Enter") + descStyle.Render("Save"),
		keyStyle.Render("Esc") + descStyle.Render("Cancel"),
	}
	editColumn := lipgloss.JoinVertical(lipgloss.Left,
		append([]string{editHeader, editUnderline}, editItems...)...)

	// General column
	generalHeader := categoryStyle.Render("General")
	generalUnderline := underlineStyle.Render("───────")
	generalItems := []string{
		keyStyle.Render("Tab") + descStyle.Render("Switch panel"),
		keyStyle.Render("Esc") + descStyle.Render("Back/Cancel"),
		keyStyle.Render("/") + descStyle.Render("Search"),
		keyStyle.Render("?") + descStyle.Render("This help"),
		keyStyle.Render("q") + descStyle.Render("Quit"),
	}
	generalColumn := lipgloss.JoinVertical(lipgloss.Left,
		append([]string{generalHeader, generalUnderline}, generalItems...)...)

	// Create two rows of columns
	columnGap := "    " // 4 spaces between columns
	topRow := lipgloss.JoinHorizontal(lipgloss.Top, navColumn, columnGap, todoColumn)
	bottomRow := lipgloss.JoinHorizontal(lipgloss.Top, editColumn, columnGap, generalColumn)

	// Footer
	footerStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#666")).
		Italic(true).
		MarginTop(2).
		Align(lipgloss.Center)
	footer := footerStyle.Render("Press any key to close")

	// Combine all content
	modalContent := lipgloss.JoinVertical(lipgloss.Center,
		header,
		"",
		topRow,
		"",
		bottomRow,
		footer,
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

// IsTooSmall checks if the terminal is below minimum size
func (v *ViewRenderer) IsTooSmall() bool {
	return v.width < MinTerminalWidth || v.height < MinTerminalHeight
}

// RenderSizeWarning renders a warning when terminal is too small
func (v *ViewRenderer) RenderSizeWarning() string {
	accentColor := "#FFB86C" // Orange for warning

	// Header
	headerStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color(accentColor))
	header := headerStyle.Render("Terminal size too small")

	// Current size
	currentStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FFFFFF"))
	current := currentStyle.Render(fmt.Sprintf("Current: %dx%d", v.width, v.height))

	// Minimum size
	minStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#888"))
	minimum := minStyle.Render(fmt.Sprintf("Minimum: %dx%d", MinTerminalWidth, MinTerminalHeight))

	// Instructions
	instructionStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#888")).
		Italic(true)
	instruction := instructionStyle.Render("Please resize your terminal")

	// Combine modal content
	modalContent := lipgloss.JoinVertical(lipgloss.Center,
		header,
		"",
		current,
		minimum,
		"",
		instruction,
	)

	// Modal box
	modalBox := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color(accentColor)).
		Padding(1, 3).
		Render(modalContent)

	// Center the modal
	return lipgloss.Place(v.width, v.height,
		lipgloss.Center, lipgloss.Center,
		modalBox)
}

// RenderStatusMessage renders a status message with appropriate color
func (v *ViewRenderer) RenderStatusMessage(message, statusType string) string {
	if message == "" {
		return ""
	}

	var color string
	switch statusType {
	case "success":
		color = "#50FA7B" // Green
	case "warning":
		color = "#FFB86C" // Orange
	case "info":
		color = "#888888" // Gray
	case "priority":
		color = "#7D56F4" // Accent purple
	default:
		color = "#FFFFFF" // White
	}

	style := lipgloss.NewStyle().
		Foreground(lipgloss.Color(color)).
		Bold(true).
		Align(lipgloss.Center).
		Width(v.width)

	return style.Render(message)
}

// RenderGoToDate renders the go-to-date modal
func (v *ViewRenderer) RenderGoToDate(state GoToDateState) string {
	accentColor := "#7D56F4"

	// Header
	headerStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color(accentColor)).
		MarginBottom(1)
	header := headerStyle.Render("Go to Date")

	// Instructions
	instructionStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#888"))
	instruction := instructionStyle.Render("Enter date (YYYY-MM-DD, MM-DD, or DD)")

	// Input field
	inputBorderColor := lipgloss.Color(accentColor)
	inputStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(inputBorderColor).
		Padding(0, 1).
		Width(24)
	inputField := inputStyle.Render(state.InputView)

	// Error message (if any)
	var errorLine string
	if state.ErrorMsg != "" {
		errorStyle := lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FF6B6B")).
			Italic(true)
		errorLine = errorStyle.Render(state.ErrorMsg)
	}

	// Combine modal content
	var modalContent string
	if errorLine != "" {
		modalContent = lipgloss.JoinVertical(lipgloss.Center,
			header,
			instruction,
			inputField,
			errorLine,
		)
	} else {
		modalContent = lipgloss.JoinVertical(lipgloss.Center,
			header,
			instruction,
			inputField,
		)
	}

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

// RenderScheduling renders the schedule input modal
func (v *ViewRenderer) RenderScheduling(state SchedulingState) string {
	accentColor := "#7D56F4"

	// Header
	headerStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color(accentColor)).
		MarginBottom(1)
	header := headerStyle.Render("Schedule Task")

	// Task name
	taskStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FFFFFF")).
		MarginBottom(1)
	taskLine := taskStyle.Render("Task: " + state.TaskTitle)

	// Instructions
	instructionStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#888"))
	instruction := instructionStyle.Render("Enter time (HH:MM or HH:MM-HH:MM)")

	// Default duration hint
	hintStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#666")).
		Italic(true)
	hint := hintStyle.Render("Default duration: 1 hour")

	// Input field
	inputBorderColor := lipgloss.Color(accentColor)
	inputStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(inputBorderColor).
		Padding(0, 1).
		Width(24)
	inputField := inputStyle.Render(state.InputView)

	// Error message (if any)
	var errorLine string
	if state.ErrorMsg != "" {
		errorStyle := lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FF6B6B")).
			Italic(true)
		errorLine = errorStyle.Render(state.ErrorMsg)
	}

	// Combine modal content
	var modalContent string
	if errorLine != "" {
		modalContent = lipgloss.JoinVertical(lipgloss.Center,
			header,
			taskLine,
			instruction,
			inputField,
			errorLine,
			hint,
		)
	} else {
		modalContent = lipgloss.JoinVertical(lipgloss.Center,
			header,
			taskLine,
			instruction,
			inputField,
			hint,
		)
	}

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

// RenderHelpBar renders the help bar at the bottom
func (v *ViewRenderer) RenderHelpBar(state AppState, focus AppFocus) string {
	return v.RenderHelpBarWithDayView(state, focus, DayViewContext{})
}

// RenderHelpBarWithDayView renders the help bar with Day View context
func (v *ViewRenderer) RenderHelpBarWithDayView(state AppState, focus AppFocus, dayView DayViewContext) string {
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
			// Check if in Day View
			if dayView.IsInDayView {
				if dayView.IsListMode {
					// Day View - List Mode
					keys = keyStyle.Render("j/k") + descStyle.Render(" nav") + sep +
						keyStyle.Render("Space") + descStyle.Render(" done") + sep +
						keyStyle.Render("e") + descStyle.Render(" edit") + sep +
						keyStyle.Render("a") + descStyle.Render(" add") + sep +
						keyStyle.Render("T") + descStyle.Render(" timeline") + sep +
						keyStyle.Render("Esc") + descStyle.Render(" back") + sep +
						keyStyle.Render("?") + descStyle.Render(" help") + sep +
						keyStyle.Render("q") + descStyle.Render(" quit")
				} else {
					// Day View - Timeline Mode
					if dayView.IsTimeline {
						// Timeline focus
						keys = keyStyle.Render("Tab") + descStyle.Render(" switch") + sep +
							keyStyle.Render("j/k") + descStyle.Render(" nav") + sep +
							keyStyle.Render("J/K") + descStyle.Render(" move") + sep +
							keyStyle.Render("+/-") + descStyle.Render(" duration") + sep +
							keyStyle.Render("u") + descStyle.Render(" unschedule") + sep +
							keyStyle.Render("T") + descStyle.Render(" list") + sep +
							keyStyle.Render("Esc") + descStyle.Render(" back") + sep +
							keyStyle.Render("?") + descStyle.Render(" help")
					} else {
						// Unscheduled focus
						keys = keyStyle.Render("Tab") + descStyle.Render(" switch") + sep +
							keyStyle.Render("j/k") + descStyle.Render(" nav") + sep +
							keyStyle.Render("s") + descStyle.Render(" schedule") + sep +
							keyStyle.Render("Enter") + descStyle.Render(" assign") + sep +
							keyStyle.Render("T") + descStyle.Render(" list") + sep +
							keyStyle.Render("Esc") + descStyle.Render(" back") + sep +
							keyStyle.Render("?") + descStyle.Render(" help")
					}
				}
			} else {
				// Non-Day View calendar
				keys = keyStyle.Render("h/j/k/l") + descStyle.Render(" nav") + sep +
					keyStyle.Render("b/n") + descStyle.Render(" month") + sep +
					keyStyle.Render("w") + descStyle.Render(" week") + sep +
					keyStyle.Render("d") + descStyle.Render(" day") + sep +
					keyStyle.Render("m") + descStyle.Render(" month") + sep +
					keyStyle.Render("t") + descStyle.Render(" today") + sep +
					keyStyle.Render("g") + descStyle.Render(" jump") + sep +
					keyStyle.Render("/") + descStyle.Render(" search") + sep +
					keyStyle.Render("Tab") + descStyle.Render(" todos") + sep +
					keyStyle.Render("?") + descStyle.Render(" help") + sep +
					keyStyle.Render("q") + descStyle.Render(" quit")
			}
		} else {
			keys = keyStyle.Render("j/k") + descStyle.Render(" nav") + sep +
				keyStyle.Render("J/K") + descStyle.Render(" move") + sep +
				keyStyle.Render("Space") + descStyle.Render(" done") + sep +
				keyStyle.Render("0-3") + descStyle.Render(" priority") + sep +
				keyStyle.Render("/") + descStyle.Render(" search") + sep +
				keyStyle.Render("a") + descStyle.Render(" add") + sep +
				keyStyle.Render("e") + descStyle.Render(" edit") + sep +
				keyStyle.Render("Esc") + descStyle.Render(" back") + sep +
				keyStyle.Render("?") + descStyle.Render(" help") + sep +
				keyStyle.Render("q") + descStyle.Render(" quit")
		}
	case StateEditing:
		keys = keyStyle.Render("Tab") + descStyle.Render(" switch field") + sep +
			keyStyle.Render("Ctrl+P") + descStyle.Render(" preview") + sep +
			keyStyle.Render("Enter") + descStyle.Render(" save") + sep +
			keyStyle.Render("Esc") + descStyle.Render(" cancel")
	case StateConfirmingDelete:
		keys = keyStyle.Render("y/Enter") + descStyle.Render(" confirm") + sep +
			keyStyle.Render("n/Esc") + descStyle.Render(" cancel")
	case StateSearching:
		keys = keyStyle.Render("Up/Down") + descStyle.Render(" navigate") + sep +
			keyStyle.Render("Enter") + descStyle.Render(" go to") + sep +
			keyStyle.Render("Esc") + descStyle.Render(" cancel")
	case StateHelp:
		keys = keyStyle.Render("any key") + descStyle.Render(" close")
	case StateGoToDate:
		keys = keyStyle.Render("Enter") + descStyle.Render(" go to date") + sep +
			keyStyle.Render("Esc") + descStyle.Render(" cancel")
	}

	return lipgloss.NewStyle().
		Width(v.width).
		Align(lipgloss.Center).
		Background(lipgloss.Color("#1a1a1a")).
		Padding(0, 1).
		Render(keys)
}
