package ui

import (
	"ctchen222/chronoflow/internal/domain"
	"ctchen222/chronoflow/internal/service"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/lipgloss"
)

// TodoItem wraps a domain.Todo for display in the UI list
// Implements list.Item interface for bubbles/list
type TodoItem struct {
	domain.Todo
	IsOverdue bool
}

// Title returns the formatted title for display
func (i TodoItem) Title() string {
	checkbox := "☐ "
	style := lipgloss.NewStyle()

	if i.Complete {
		checkbox = "☑ "
		style = style.Foreground(lipgloss.Color("#666")).Strikethrough(true)
	} else if i.IsOverdue {
		checkbox = "⚠ "
		style = style.Foreground(lipgloss.Color("#FF6B6B")).Bold(true)
	} else if i.Priority > 0 {
		style = priorityStyle(i.Priority)
	}

	title := style.Render(i.Todo.Title)
	if i.Priority > 0 && !i.Complete {
		title = priorityStyle(i.Priority).Render(i.Priority.Icon()) + " " + title
	}

	return checkbox + title
}

// Description returns the formatted description for display
func (i TodoItem) Description() string {
	if i.IsOverdue && !i.Complete {
		return lipgloss.NewStyle().Foreground(lipgloss.Color("#FF6B6B")).Render("Overdue! " + i.Desc)
	}
	return i.Desc
}

// FilterValue returns the value used for filtering
func (i TodoItem) FilterValue() string {
	return i.Todo.Title
}

// priorityStyle returns the lipgloss style for a priority level
func priorityStyle(p domain.Priority) lipgloss.Style {
	switch p {
	case domain.PriorityHigh:
		return lipgloss.NewStyle().Foreground(lipgloss.Color("#FF6B6B")).Bold(true)
	case domain.PriorityMedium:
		return lipgloss.NewStyle().Foreground(lipgloss.Color("#FFB86C"))
	case domain.PriorityLow:
		return lipgloss.NewStyle().Foreground(lipgloss.Color("#8BE9FD"))
	default:
		return lipgloss.NewStyle()
	}
}

// TodoPresenter converts service data to UI-ready format
type TodoPresenter struct{}

// NewTodoPresenter creates a new TodoPresenter
func NewTodoPresenter() *TodoPresenter {
	return &TodoPresenter{}
}

// ToListItems converts TodoWithStatus slice to list.Item slice
func (p *TodoPresenter) ToListItems(todos []service.TodoWithStatus) []list.Item {
	items := make([]list.Item, len(todos))
	for i, td := range todos {
		items[i] = TodoItem{
			Todo:      td.Todo,
			IsOverdue: td.IsOverdue,
		}
	}
	return items
}

// PriorityOption represents a priority option for the UI
type PriorityOption struct {
	Level domain.Priority
	Label string
	Color string
}

// GetPriorityOptions returns the available priority options
func (p *TodoPresenter) GetPriorityOptions() []PriorityOption {
	return []PriorityOption{
		{domain.PriorityNone, "None", "#666"},
		{domain.PriorityLow, "Low", "#8BE9FD"},
		{domain.PriorityMedium, "Medium", "#FFB86C"},
		{domain.PriorityHigh, "High", "#FF6B6B"},
	}
}
