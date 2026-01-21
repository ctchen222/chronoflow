# Implementation Tasks

## 1. Add View Mode Indicator

- [x] 1.1 Add `RenderViewModeIndicator()` method in `pkg/calendar/calendar.go` to render "MONTH VIEW" or "WEEK VIEW" badge
- [x] 1.2 Update `View()` method in month view to include the view mode indicator in header area
- [x] 1.3 Update `renderWeekView()` method to include the view mode indicator in header area
- [x] 1.4 Style the indicator with accent color and consistent formatting

## 2. Add Todo Count in Calendar Header

- [x] 2.1 Add helper method `GetTodoCountForDate(date time.Time) int` in `pkg/calendar/calendar.go`
- [x] 2.2 Update month view sub-header to show todo count: "Monday, Jan 2, 2006 (3 tasks)"
- [x] 2.3 Update week view sub-header to show todo count for cursor date
- [x] 2.4 Handle zero count case gracefully (no count displayed when 0 tasks)

## 3. Testing & Validation

- [x] 3.1 Verify view mode indicator displays correctly in both month and week views
- [x] 3.2 Verify todo count updates when navigating between dates
- [x] 3.3 Verify styling is consistent with existing UI patterns
- [x] 3.4 Test edge cases: no todos, many todos, terminal resize
