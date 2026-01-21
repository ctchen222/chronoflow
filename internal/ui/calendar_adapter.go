package ui

import (
	"time"

	"ctchen222/chronoflow/internal/domain"
	"ctchen222/chronoflow/internal/service"
	"ctchen222/chronoflow/pkg/calendar"
)

// CalendarAdapter adapts todo data for calendar display
type CalendarAdapter struct {
	statsCalc *service.StatsCalculator
}

// NewCalendarAdapter creates a new CalendarAdapter
func NewCalendarAdapter(statsCalc *service.StatsCalculator) *CalendarAdapter {
	return &CalendarAdapter{
		statsCalc: statsCalc,
	}
}

// BuildTodoStatus builds the todo status map for calendar display
func (a *CalendarAdapter) BuildTodoStatus(allTodos map[string][]domain.Todo) map[string]calendar.TodoStatus {
	todoStatus := make(map[string]calendar.TodoStatus)

	for dateKey, items := range allTodos {
		if len(items) == 0 {
			continue
		}

		// Parse the date
		date, err := time.Parse("2006-01-02", dateKey)
		if err != nil {
			continue
		}

		// Use stats calculator helper methods
		isOverdue := a.statsCalc.IsDateOverdue(items, date)
		allComplete := a.statsCalc.AreAllComplete(items)

		// Convert items to calendar.TodoItem for week/day view display
		calendarItems := make([]calendar.TodoItem, len(items))
		for i, it := range items {
			startTime := ""
			endTime := ""
			if it.StartTime != nil {
				startTime = *it.StartTime
			}
			if it.EndTime != nil {
				endTime = *it.EndTime
			}
			calendarItems[i] = calendar.TodoItem{
				Title:     it.Title,
				Desc:      it.Desc,
				Complete:  it.Complete,
				Priority:  int(it.Priority),
				StartTime: startTime,
				EndTime:   endTime,
			}
		}

		todoStatus[dateKey] = calendar.TodoStatus{
			HasTodos:    true,
			HasOverdue:  isOverdue,
			AllComplete: allComplete,
			Count:       len(items),
			Items:       calendarItems,
		}
	}

	return todoStatus
}

// ConvertViewMode converts calendar.ViewMode to service.ViewMode
func ConvertViewMode(vm calendar.ViewMode) service.ViewMode {
	if vm == calendar.WeekView {
		return service.WeekView
	}
	return service.MonthView
}
