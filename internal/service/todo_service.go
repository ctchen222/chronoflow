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
