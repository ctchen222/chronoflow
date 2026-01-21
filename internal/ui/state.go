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
	StateHelp
	StateGoToDate
	StateScheduling
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

// GoToDateState holds the state for the go-to-date modal
type GoToDateState struct {
	InputView  string // rendered date input
	InputValue string
	ErrorMsg   string // error message for invalid date
}

// SchedulingState holds the state for the schedule input modal
type SchedulingState struct {
	TaskTitle  string // title of the task being scheduled
	TaskIndex  int    // original index in the todo list
	InputView  string // rendered time input
	InputValue string
	ErrorMsg   string // error message for invalid time
}

// DayViewContext holds Day View state for help bar rendering
type DayViewContext struct {
	IsInDayView bool // whether we're in Day View
	IsListMode  bool // true = List Mode, false = Timeline Mode
	IsTimeline  bool // Timeline focus (vs Unscheduled focus) in Timeline Mode
}
