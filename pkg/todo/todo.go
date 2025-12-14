package todo

import (
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	docStyle = lipgloss.NewStyle().Margin(1, 2)
)

type Model struct {
	list list.Model
}

func New() Model {
	m := Model{
		list: list.New([]list.Item{}, list.NewDefaultDelegate(), 0, 0),
	}
	m.list.Title = "To-Do List"
	return m
}

func (m *Model) SetSize(w, h int) {
	m.list.SetSize(w, h)
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m Model) View() string {
	return docStyle.Render(m.list.View())
}

func (m *Model) SetTitle(title string) {
	m.list.Title = title
}

// SetItems replaces the current list of items with a new one.
func (m *Model) SetItems(items []list.Item) {
	m.list.SetItems(items)
}

// SelectedItem returns the currently selected item in the list.
func (m Model) SelectedItem() list.Item {
	return m.list.SelectedItem()
}

// ListIndex returns the index of the currently selected item.
func (m Model) ListIndex() int {
	return m.list.Index()
}
