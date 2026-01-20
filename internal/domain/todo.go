package domain

import "time"

// Todo represents a todo item (pure data, no UI concerns)
type Todo struct {
	Title     string   `json:"title"`
	Desc      string   `json:"desc"`
	Complete  bool     `json:"completed"`
	Priority  Priority `json:"priority"`
	StartTime *string  `json:"start_time,omitempty"` // "HH:MM" format, nil = unscheduled
	EndTime   *string  `json:"end_time,omitempty"`   // "HH:MM" format
}

// IsScheduled returns true if the todo has a scheduled start time
func (t Todo) IsScheduled() bool {
	return t.StartTime != nil
}

// Duration returns the duration of a scheduled todo
// Returns 0 if the todo is not scheduled or has invalid times
func (t Todo) Duration() time.Duration {
	if t.StartTime == nil || t.EndTime == nil {
		return 0
	}
	start, err := time.Parse("15:04", *t.StartTime)
	if err != nil {
		return 0
	}
	end, err := time.Parse("15:04", *t.EndTime)
	if err != nil {
		return 0
	}
	return end.Sub(start)
}

// NewTodo creates a new Todo with the given title
func NewTodo(title string) Todo {
	return Todo{
		Title:    title,
		Priority: PriorityNone,
	}
}

// IsOverdue checks if a todo is overdue based on its date and current time
func (t Todo) IsOverdue(todoDate, now time.Time) bool {
	if t.Complete {
		return false
	}
	// Normalize both dates to midnight for comparison
	todoMidnight := time.Date(todoDate.Year(), todoDate.Month(), todoDate.Day(), 0, 0, 0, 0, time.Local)
	nowMidnight := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.Local)
	return todoMidnight.Before(nowMidnight)
}

// TodoList represents a collection of todos for a specific date
type TodoList struct {
	Date  time.Time
	Todos []Todo
}

// Stats returns statistics for this todo list
func (tl TodoList) Stats(now time.Time) (total, completed, overdue int) {
	total = len(tl.Todos)
	for _, t := range tl.Todos {
		if t.Complete {
			completed++
		} else if t.IsOverdue(tl.Date, now) {
			overdue++
		}
	}
	return
}
