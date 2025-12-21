package repository

import (
	"encoding/json"
	"os"
	"path/filepath"
	"time"

	"ctchen222/chronoflow/internal/domain"
)

// TodoRepository defines the interface for todo data access
type TodoRepository interface {
	// GetAll returns all todos grouped by date
	GetAll() map[string][]domain.Todo

	// GetByDate returns todos for a specific date
	GetByDate(date time.Time) []domain.Todo

	// Save saves a todo for a specific date (creates or updates)
	Save(date time.Time, index int, todo domain.Todo) error

	// Add adds a new todo for a specific date
	Add(date time.Time, todo domain.Todo) error

	// Delete removes a todo at the given index for a specific date
	Delete(date time.Time, index int) error

	// Reorder swaps two todos at the given indices
	Reorder(date time.Time, fromIndex, toIndex int) error

	// Load loads todos from persistent storage
	Load() error

	// Persist saves todos to persistent storage
	Persist() error
}

// JSONTodoRepository implements TodoRepository using JSON file storage
type JSONTodoRepository struct {
	todos    map[string][]domain.Todo
	filePath string
}

// NewJSONTodoRepository creates a new JSON-based repository
func NewJSONTodoRepository(filePath string) *JSONTodoRepository {
	return &JSONTodoRepository{
		todos:    make(map[string][]domain.Todo),
		filePath: filePath,
	}
}

// dateKey formats a time.Time to the standard date key format
func dateKey(date time.Time) string {
	return date.Format("2006-01-02")
}

func (r *JSONTodoRepository) GetAll() map[string][]domain.Todo {
	return r.todos
}

func (r *JSONTodoRepository) GetByDate(date time.Time) []domain.Todo {
	key := dateKey(date)
	if todos, ok := r.todos[key]; ok {
		// Return a copy to prevent external mutation
		result := make([]domain.Todo, len(todos))
		copy(result, todos)
		return result
	}
	return []domain.Todo{}
}

func (r *JSONTodoRepository) Save(date time.Time, index int, todo domain.Todo) error {
	key := dateKey(date)
	if todos, ok := r.todos[key]; ok && index >= 0 && index < len(todos) {
		r.todos[key][index] = todo
		return nil
	}
	return nil
}

func (r *JSONTodoRepository) Add(date time.Time, todo domain.Todo) error {
	key := dateKey(date)
	r.todos[key] = append(r.todos[key], todo)
	return nil
}

func (r *JSONTodoRepository) Delete(date time.Time, index int) error {
	key := dateKey(date)
	if todos, ok := r.todos[key]; ok && index >= 0 && index < len(todos) {
		r.todos[key] = append(todos[:index], todos[index+1:]...)
	}
	return nil
}

func (r *JSONTodoRepository) Reorder(date time.Time, fromIndex, toIndex int) error {
	key := dateKey(date)
	if todos, ok := r.todos[key]; ok {
		if fromIndex >= 0 && fromIndex < len(todos) && toIndex >= 0 && toIndex < len(todos) {
			todos[fromIndex], todos[toIndex] = todos[toIndex], todos[fromIndex]
		}
	}
	return nil
}

func (r *JSONTodoRepository) Load() error {
	data, err := os.ReadFile(r.filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil // No file, start fresh
		}
		return err
	}
	if len(data) == 0 {
		return nil
	}
	return json.Unmarshal(data, &r.todos)
}

func (r *JSONTodoRepository) Persist() error {
	// Ensure directory exists
	dir := filepath.Dir(r.filePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	data, err := json.MarshalIndent(r.todos, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(r.filePath, data, 0644)
}
