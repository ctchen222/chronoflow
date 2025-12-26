package domain

import "time"

// Todo represents a todo item (pure data, no UI concerns)
type Todo struct {
	Title    string   `json:"title"`
	Desc     string   `json:"desc"`
	Complete bool     `json:"completed"`
	Priority Priority `json:"priority"`
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
