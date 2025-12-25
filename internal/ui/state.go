package ui

import (
	"time"

	"ctchen222/chronoflow/internal/domain"
	"ctchen222/chronoflow/internal/service"
)

// AppState defines the current state of the application
type AppState int

const (
	StateViewing AppState = iota
	StateEditing
	StateConfirmingDelete
	StateSearching
)

// AppFocus determines which panel is currently focused
type AppFocus int

const (
	FocusCalendar AppFocus = iota
	FocusTodo
)

// EditFocus determines which input is focused in the editing view
type EditFocus int

const (
	FocusTitle EditFocus = iota
	FocusDesc
)

// EditingState holds the state for the editing modal
type EditingState struct {
	IsNew       bool
	Date        time.Time
	TitleValue  string
	DescValue   string
	Priority    domain.Priority
	Focus       EditFocus
	TitleView   string // rendered title input
	DescView    string // rendered desc input
	// Preview fields
	PreviewEnabled bool   // whether preview pane is visible
	PreviewContent string // rendered markdown preview
}

// DeleteState holds the state for the delete confirmation modal
type DeleteState struct {
	Title string
}

// SearchState holds the state for the search modal
type SearchState struct {
	InputView   string // rendered search input
	InputValue  string
	Results     []service.SearchResult
	SelectedIdx int
}

// MainViewState holds the state for the main view
type MainViewState struct {
	CalendarView string
	TodoView     string
	Focus        AppFocus
}
