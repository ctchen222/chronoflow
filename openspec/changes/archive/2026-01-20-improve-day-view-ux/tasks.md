# Tasks: Improve Day View UX

## Implementation Order

### Phase 1: Day View List Mode

1. [x] **Add DayViewMode type to calendar**
   - File: `pkg/calendar/calendar.go`
   - Add `DayViewMode` type (List, Timeline)
   - Add `dayViewMode` field to Model
   - Default to `DayViewModeList`

2. [x] **Implement List Mode rendering**
   - File: `pkg/calendar/calendar.go`
   - Add `renderDayViewList()` method
   - Render card-style task items with title + description
   - Sort by priority, completed tasks at bottom

3. [x] **Add T key toggle**
   - File: `cmd/chronoflow/main.go`
   - Handle `T` key in Day View to toggle mode
   - Call `calendar.ToggleDayViewMode()`

4. [x] **Update renderDayView() dispatcher**
   - File: `pkg/calendar/calendar.go`
   - Check `dayViewMode` and call appropriate renderer
   - List Mode → `renderDayViewList()`
   - Timeline Mode → existing `renderDayViewTimeline()`

### Phase 2: Dynamic Help Bar

5. [x] **Extend help bar for Day View modes**
   - File: `internal/ui/views.go`
   - Update `RenderHelpBar()` to accept Day View context
   - Add cases for List Mode and Timeline Mode (with focus)

6. [x] **Pass Day View state to help bar**
   - File: `cmd/chronoflow/main.go`
   - Pass `dayViewMode` and `dayViewFocus` to ViewRenderer
   - Update View() to include Day View state

### Phase 3: Empty State Hints

7. [x] **Add empty state rendering for List Mode**
   - File: `pkg/calendar/calendar.go`
   - Show hint when no tasks: "Press 'a' to add a task"

8. [x] **Add empty state hints for Timeline Mode**
   - File: `pkg/calendar/calendar.go`
   - Empty timeline: "Select a task and press Enter"
   - Empty unscheduled: "All tasks scheduled!"

### Phase 4: Task Movement

9. [x] **Add MoveMinutes to TimelineConfig**
   - File: `internal/domain/config.go`
   - Add `MoveMinutes int` field
   - Default to `SlotMinutes` value

10. [x] **Add RescheduleTodo service method**
    - File: `internal/service/todo_service.go`
    - Add `RescheduleTodo(date, index, deltaMinutes, dayStart, dayEnd)`
    - Validate time boundaries

11. [x] **Implement Shift+J/K handlers**
    - File: `cmd/chronoflow/main.go`
    - Handle `J` (Shift+j) and `K` (Shift+k) in Timeline focus
    - Calculate new times based on MoveMinutes
    - Show boundary error messages

12. [x] **Add timeline cursor for scheduled tasks**
    - File: `pkg/calendar/calendar.go`
    - Track selected scheduled task index
    - Highlight selected task in timeline

### Phase 5: Testing

13. [x] **Update golden files**
    - Updated golden files for List Mode (now default)
    - All tests passing

14. [x] **Tests verified**
    - Build succeeds: `go build ./cmd/chronoflow`
    - All tests pass: `go test ./...`

## Validation

```bash
go test ./...
go build ./cmd/chronoflow
```

## Dependencies

- Phase 1 can start immediately
- Phase 2 depends on Phase 1 (needs mode state)
- Phase 3 depends on Phase 1 (needs mode rendering)
- Phase 4 can start after Phase 1
- Phase 5 after all implementation phases
