package service

import (
	"time"

	"ctchen222/chronoflow/internal/domain"
)

// ViewMode represents the calendar view mode
type ViewMode int

const (
	MonthView ViewMode = iota
	WeekView
)

// Stats holds statistics about todos
type Stats struct {
	TotalAll        int
	CompletedAll    int
	OverdueAll      int
	TotalPeriod     int
	CompletedPeriod int
	OverduePeriod   int
	PeriodLabel     string
}

// StatsCalculator calculates todo statistics
type StatsCalculator struct {
	timeProvider TimeProvider
}

// NewStatsCalculator creates a new StatsCalculator
func NewStatsCalculator(tp TimeProvider) *StatsCalculator {
	return &StatsCalculator{timeProvider: tp}
}

// CalculateStats calculates statistics for all todos based on view mode and cursor date
func (sc *StatsCalculator) CalculateStats(
	todos map[string][]domain.Todo,
	viewMode ViewMode,
	cursorDate time.Time,
) Stats {
	today := sc.timeProvider.Today()
	periodStart, periodEnd, periodLabel := sc.getPeriodBounds(viewMode, cursorDate)

	stats := Stats{
		PeriodLabel: periodLabel,
	}

	for dateKey, items := range todos {
		date, err := time.Parse("2006-01-02", dateKey)
		if err != nil {
			continue
		}
		date = time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, time.Local)

		isPast := date.Before(today)
		inPeriod := !date.Before(periodStart) && date.Before(periodEnd)

		for _, item := range items {
			stats.TotalAll++
			if item.Complete {
				stats.CompletedAll++
			} else if isPast {
				stats.OverdueAll++
			}

			if inPeriod {
				stats.TotalPeriod++
				if item.Complete {
					stats.CompletedPeriod++
				} else if isPast {
					stats.OverduePeriod++
				}
			}
		}
	}

	return stats
}

// getPeriodBounds returns the start and end dates for the current period
func (sc *StatsCalculator) getPeriodBounds(viewMode ViewMode, cursorDate time.Time) (start, end time.Time, label string) {
	if viewMode == WeekView {
		weekday := int(cursorDate.Weekday())
		start = time.Date(cursorDate.Year(), cursorDate.Month(), cursorDate.Day()-weekday, 0, 0, 0, 0, time.Local)
		end = start.AddDate(0, 0, 7)
		label = "This Week"
	} else {
		start = time.Date(cursorDate.Year(), cursorDate.Month(), 1, 0, 0, 0, 0, time.Local)
		end = start.AddDate(0, 1, 0)
		label = "This Month"
	}
	return
}

// CalculateDateStats calculates stats for a single date's todos
func (sc *StatsCalculator) CalculateDateStats(todos []domain.Todo, todoDate time.Time) (total, completed, overdue int) {
	today := sc.timeProvider.Today()
	total = len(todos)

	for _, t := range todos {
		if t.Complete {
			completed++
		} else if t.IsOverdue(todoDate, today) {
			overdue++
		}
	}
	return
}

// IsDateOverdue checks if a date has any overdue incomplete todos
func (sc *StatsCalculator) IsDateOverdue(todos []domain.Todo, todoDate time.Time) bool {
	today := sc.timeProvider.Today()
	if !todoDate.Before(today) {
		return false
	}

	for _, t := range todos {
		if !t.Complete {
			return true
		}
	}
	return false
}

// AreAllComplete checks if all todos for a date are completed
func (sc *StatsCalculator) AreAllComplete(todos []domain.Todo) bool {
	if len(todos) == 0 {
		return false
	}
	for _, t := range todos {
		if !t.Complete {
			return false
		}
	}
	return true
}
