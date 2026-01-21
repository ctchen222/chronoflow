package service

import (
	"strings"
	"time"

	"ctchen222/chronoflow/internal/domain"
	"ctchen222/chronoflow/internal/repository"
)

// TodoService provides business logic operations for todos
type TodoService struct {
	repo     repository.TodoRepository
	timeProv TimeProvider
}

// NewTodoService creates a new TodoService
func NewTodoService(repo repository.TodoRepository, timeProv TimeProvider) *TodoService {
	return &TodoService{
		repo:     repo,
		timeProv: timeProv,
	}
}

// GetTodosForDate returns todos for a specific date with overdue status calculated
func (s *TodoService) GetTodosForDate(date time.Time) []TodoWithStatus {
	todos := s.repo.GetByDate(date)
	today := s.timeProv.Today()

	result := make([]TodoWithStatus, len(todos))
	for i, td := range todos {
		result[i] = TodoWithStatus{
			Todo:      td,
			IsOverdue: td.IsOverdue(date, today),
		}
	}
	return result
}

// TodoWithStatus wraps a Todo with calculated display status
type TodoWithStatus struct {
	domain.Todo
	IsOverdue bool
}

// ToggleComplete toggles the completion status of a todo
func (s *TodoService) ToggleComplete(date time.Time, index int) error {
	todos := s.repo.GetByDate(date)
	if index < 0 || index >= len(todos) {
		return nil
	}
	todos[index].Complete = !todos[index].Complete
	return s.repo.Save(date, index, todos[index])
}

// SetPriority sets the priority of a todo
func (s *TodoService) SetPriority(date time.Time, index int, priority domain.Priority) error {
	todos := s.repo.GetByDate(date)
	if index < 0 || index >= len(todos) {
		return nil
	}
	todos[index].Priority = priority
	return s.repo.Save(date, index, todos[index])
}

// Add creates a new todo for a date
func (s *TodoService) Add(date time.Time, title, desc string, priority domain.Priority) error {
	if title == "" {
		return nil
	}
	todo := domain.Todo{
		Title:    title,
		Desc:     desc,
		Priority: priority,
	}
	return s.repo.Add(date, todo)
}

// Update updates an existing todo
func (s *TodoService) Update(date time.Time, index int, title, desc string, priority domain.Priority) error {
	todos := s.repo.GetByDate(date)
	if index < 0 || index >= len(todos) {
		return nil
	}
	todos[index].Title = title
	todos[index].Desc = desc
	todos[index].Priority = priority
	return s.repo.Save(date, index, todos[index])
}

// Delete removes a todo
func (s *TodoService) Delete(date time.Time, index int) error {
	return s.repo.Delete(date, index)
}

// MoveUp moves a todo up in the list (swap with previous)
func (s *TodoService) MoveUp(date time.Time, index int) error {
	if index <= 0 {
		return nil
	}
	return s.repo.Reorder(date, index, index-1)
}

// MoveDown moves a todo down in the list (swap with next)
func (s *TodoService) MoveDown(date time.Time, index int) error {
	todos := s.repo.GetByDate(date)
	if index < 0 || index >= len(todos)-1 {
		return nil
	}
	return s.repo.Reorder(date, index, index+1)
}

// SearchResult represents a search result
type SearchResult struct {
	DateKey string
	Index   int
	Todo    domain.Todo
}

// Search searches all todos for the given query
func (s *TodoService) Search(query string) []SearchResult {
	if query == "" {
		return nil
	}

	query = strings.ToLower(query)
	allTodos := s.repo.GetAll()

	var results []SearchResult
	for dateKey, items := range allTodos {
		for idx, td := range items {
			if strings.Contains(strings.ToLower(td.Title), query) ||
				strings.Contains(strings.ToLower(td.Desc), query) {
				results = append(results, SearchResult{
					DateKey: dateKey,
					Index:   idx,
					Todo:    td,
				})
			}
		}
	}
	return results
}

// GetAllTodos returns all todos (for stats calculation)
func (s *TodoService) GetAllTodos() map[string][]domain.Todo {
	return s.repo.GetAll()
}

// Persist saves all todos to persistent storage
func (s *TodoService) Persist() error {
	return s.repo.Persist()
}

// ScheduleTodo assigns start and end times to a todo
func (s *TodoService) ScheduleTodo(date time.Time, index int, startTime, endTime string) error {
	todos := s.repo.GetByDate(date)
	if index < 0 || index >= len(todos) {
		return nil
	}
	todos[index].StartTime = &startTime
	todos[index].EndTime = &endTime
	return s.repo.Save(date, index, todos[index])
}

// UnscheduleTodo removes the scheduled time from a todo
func (s *TodoService) UnscheduleTodo(date time.Time, index int) error {
	todos := s.repo.GetByDate(date)
	if index < 0 || index >= len(todos) {
		return nil
	}
	todos[index].StartTime = nil
	todos[index].EndTime = nil
	return s.repo.Save(date, index, todos[index])
}

// AdjustTodoDuration extends or shrinks a scheduled todo's duration
// deltaMinutes can be positive (extend) or negative (shrink)
// Returns the new end time or empty string if adjustment failed
func (s *TodoService) AdjustTodoDuration(date time.Time, index int, deltaMinutes int, minSlotMinutes int) (string, error) {
	todos := s.repo.GetByDate(date)
	if index < 0 || index >= len(todos) {
		return "", nil
	}
	todo := todos[index]
	if todo.StartTime == nil || todo.EndTime == nil {
		return "", nil
	}

	start, err := time.Parse("15:04", *todo.StartTime)
	if err != nil {
		return "", nil
	}
	end, err := time.Parse("15:04", *todo.EndTime)
	if err != nil {
		return "", nil
	}

	newEnd := end.Add(time.Duration(deltaMinutes) * time.Minute)
	minEnd := start.Add(time.Duration(minSlotMinutes) * time.Minute)

	// Enforce minimum duration
	if newEnd.Before(minEnd) {
		return "", nil
	}

	newEndStr := newEnd.Format("15:04")
	todos[index].EndTime = &newEndStr
	if err := s.repo.Save(date, index, todos[index]); err != nil {
		return "", err
	}
	return newEndStr, nil
}

// RescheduleTodo moves a scheduled todo to a new time by delta minutes
// deltaMinutes can be positive (later) or negative (earlier)
// dayStart and dayEnd are boundaries in "HH:MM" format
// Returns the new start time, or empty string if movement would exceed boundaries
func (s *TodoService) RescheduleTodo(date time.Time, index int, deltaMinutes int, dayStart, dayEnd string) (string, error) {
	todos := s.repo.GetByDate(date)
	if index < 0 || index >= len(todos) {
		return "", nil
	}
	todo := todos[index]
	if todo.StartTime == nil || todo.EndTime == nil {
		return "", nil
	}

	start, err := time.Parse("15:04", *todo.StartTime)
	if err != nil {
		return "", nil
	}
	end, err := time.Parse("15:04", *todo.EndTime)
	if err != nil {
		return "", nil
	}

	// Parse boundaries
	dayStartTime, err := time.Parse("15:04", dayStart)
	if err != nil {
		return "", nil
	}
	dayEndTime, err := time.Parse("15:04", dayEnd)
	if err != nil {
		return "", nil
	}

	// Calculate new times
	delta := time.Duration(deltaMinutes) * time.Minute
	newStart := start.Add(delta)
	newEnd := end.Add(delta)

	// Check boundaries
	if newStart.Before(dayStartTime) || newEnd.After(dayEndTime) {
		return "", nil // Cannot move past boundaries
	}

	newStartStr := newStart.Format("15:04")
	newEndStr := newEnd.Format("15:04")
	todos[index].StartTime = &newStartStr
	todos[index].EndTime = &newEndStr
	if err := s.repo.Save(date, index, todos[index]); err != nil {
		return "", err
	}
	return newStartStr, nil
}
