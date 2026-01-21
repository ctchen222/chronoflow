# Tasks: Enhance Day View with Timeline

## Implementation Order

### Phase 1: Domain Layer Foundation

1. [x] **Extend Todo struct with time fields**
   - File: `internal/domain/todo.go`
   - Add `StartTime *string` and `EndTime *string` fields
   - Add `IsScheduled() bool` method
   - Add `Duration() time.Duration` method

2. [x] **Create TimelineConfig struct**
   - File: `internal/domain/config.go` (new)
   - Define `TimelineConfig` with `DayStart`, `DayEnd`, `SlotMinutes`
   - Define `Config` wrapper struct

3. [x] **Create ConfigRepository**
   - File: `internal/repository/config_repository.go` (new)
   - Implement `Load()` and `Save()` for `~/.chronoflow/config.json`
   - Handle auto-creation with defaults

### Phase 2: Service Layer

4. [x] **Add scheduling methods to TodoService**
   - File: `internal/service/todo_service.go`
   - Add `ScheduleTodo(dateKey string, index int, startTime, endTime string)`
   - Add `UnscheduleTodo(dateKey string, index int)`
   - Add `AdjustTodoDuration(dateKey string, index int, deltaMinutes int)`

### Phase 3: Calendar Component Updates

5. [x] **Extend TodoItem struct for calendar**
   - File: `pkg/calendar/calendar.go`
   - Add `Desc`, `StartTime`, `EndTime` fields to `TodoItem`

6. [x] **Update CalendarAdapter**
   - File: `internal/ui/calendar_adapter.go`
   - Pass `Desc`, `StartTime`, `EndTime` to calendar `TodoItem`

7. [x] **Implement split-panel Day View layout**
   - File: `pkg/calendar/calendar.go`
   - Modify `renderDayView()` to render two panels
   - Add timeline panel rendering with time markers
   - Add unscheduled panel rendering with descriptions

8. [x] **Add timeline cursor and navigation**
   - File: `pkg/calendar/calendar.go`
   - Add `timelineCursor int` field
   - Add methods: `MoveTimelineCursor(delta int)`, `GetTimelineCursorTime() string`

### Phase 4: Main Application Integration

9. [x] **Add Day View focus management**
   - File: `cmd/chronoflow/main.go`
   - Add `DayViewFocus` type (Timeline, Unscheduled)
   - Add `dayViewFocus`, `selectedUnscheduled` fields to model
   - Handle Tab to switch focus

10. [x] **Implement cursor-based scheduling**
    - File: `cmd/chronoflow/main.go`
    - Handle Enter in timeline focus to assign selected task
    - Update todo via service and refresh view

11. [x] **Implement quick schedule input (s key)**
    - File: `cmd/chronoflow/main.go`
    - Add `StateScheduling` app state
    - Add `scheduleInput` textinput model
    - Parse time input: `HH:MM` or `HH:MM-HH:MM`

12. [x] **Add schedule input modal rendering**
    - File: `internal/ui/views.go`
    - Implement `RenderScheduling()` modal

13. [x] **Implement unschedule (u key)**
    - File: `cmd/chronoflow/main.go`
    - Handle 'u' key on timeline to unschedule task

14. [x] **Implement duration adjustment (+/- keys)**
    - File: `cmd/chronoflow/main.go`
    - Handle '+' and '-' to adjust EndTime by slot increment

15. [x] **Update state.go with new states**
    - File: `internal/ui/state.go`
    - Add `StateScheduling` app state
    - Add `SchedulingState` struct

### Phase 5: Testing

16. [x] **UI golden file tests**
    - Updated golden files for new Day View layout
    - All tests passing

## Validation

```bash
go test ./...
go build ./cmd/chronoflow
```

## Dependencies

- Tasks 1-3 (domain/repository) can be done in parallel
- Task 4 depends on tasks 1-3
- Tasks 5-6 depend on task 1
- Tasks 7-8 depend on tasks 5-6
- Tasks 9-15 depend on tasks 4, 7-8
- Tasks 16-18 can begin after their respective implementation tasks
