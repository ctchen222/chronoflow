# Design: Enhanced Day View with Timeline

## Architecture Overview

This feature spans multiple layers of the application:

```
┌─────────────────────────────────────────────────────────────────┐
│                        UI Layer Changes                         │
├─────────────────────────────────────────────────────────────────┤
│  cmd/chronoflow/main.go                                         │
│  - Add DayViewFocus enum (Timeline, Unscheduled)                │
│  - Handle scheduling key bindings (s, u, Enter)                 │
│  - Manage cursor position on timeline                           │
├─────────────────────────────────────────────────────────────────┤
│  pkg/calendar/calendar.go                                       │
│  - Extend renderDayView() for split-panel layout                │
│  - Add timeline rendering with time markers                     │
│  - Add unscheduled task list rendering                          │
│  - Track timeline cursor position                               │
├─────────────────────────────────────────────────────────────────┤
│  internal/ui/views.go                                           │
│  - Add RenderScheduleInput() for time input modal               │
│  - Update help bar for Day View context                         │
├─────────────────────────────────────────────────────────────────┤
│                      Domain Layer Changes                        │
├─────────────────────────────────────────────────────────────────┤
│  internal/domain/todo.go                                        │
│  - Add StartTime, EndTime fields to Todo struct                 │
│  - Add IsScheduled(), Duration() helper methods                 │
├─────────────────────────────────────────────────────────────────┤
│  internal/domain/config.go (NEW)                                │
│  - TimelineConfig struct for user preferences                   │
├─────────────────────────────────────────────────────────────────┤
│                    Repository Layer Changes                      │
├─────────────────────────────────────────────────────────────────┤
│  internal/repository/config_repository.go (NEW)                 │
│  - Load/save config from ~/.chronoflow/config.json              │
├─────────────────────────────────────────────────────────────────┤
│                     Service Layer Changes                        │
├─────────────────────────────────────────────────────────────────┤
│  internal/service/todo_service.go                               │
│  - Add ScheduleTodo(index, startTime, endTime)                  │
│  - Add UnscheduleTodo(index)                                    │
└─────────────────────────────────────────────────────────────────┘
```

## Data Model

### Todo Extension

```go
type Todo struct {
    Title     string   `json:"title"`
    Desc      string   `json:"desc"`
    Complete  bool     `json:"completed"`
    Priority  Priority `json:"priority"`
    StartTime *string  `json:"start_time,omitempty"` // "HH:MM" format, nil = unscheduled
    EndTime   *string  `json:"end_time,omitempty"`   // "HH:MM" format
}
```

**Design Decisions**:
- Use `*string` (pointer) to distinguish between "unscheduled" (nil) and "00:00" (midnight)
- `omitempty` prevents unnecessary JSON fields for unscheduled tasks
- Time stored as string for simplicity and JSON readability
- No date component in time fields (date is already the map key in storage)

### Timeline Configuration

```go
type TimelineConfig struct {
    DayStart    string `json:"day_start"`    // "08:00" default
    DayEnd      string `json:"day_end"`      // "18:00" default
    SlotMinutes int    `json:"slot_minutes"` // 30 default (minimum granularity)
}
```

**Storage**: `~/.chronoflow/config.json`
```json
{
  "timeline": {
    "day_start": "08:00",
    "day_end": "18:00",
    "slot_minutes": 30
  }
}
```

## UI Layout

### Split-Panel Day View

```
┌─────────────────────────────────────────────────────────────────────┐
│                           DAY VIEW                                   │
│                   Monday, January 20, 2026                          │
│                  📅 January 2026 › Week 4 › Mon 20                  │
├────────────────────────────────────┬────────────────────────────────┤
│         TIMELINE                   │       UNSCHEDULED (3)          │
├────────────────────────────────────┼────────────────────────────────┤
│ 08:00 ┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄ │ > ☐ 回覆客戶郵件 !!!           │
│       ┌──────────────────────┐     │     需要確認報價細節            │
│ 09:00 │ ☐ 寫週報 !!          │     │                                │
│       │   整理本週進度        │     │   ☐ 買咖啡豆                   │
│ 10:00 │                      │     │                                │
│       └──────────────────────┘     │   ☑ 繳電話費                   │
│ 11:00 ┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄ │                                │
│       ┌──────────────────────┐     │                                │
│ 12:00 │ ☐ 午餐會議           │     │                                │
│       └──────────────────────┘     │                                │
│ 13:00 ┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄ │                                │
│  ...                               │                                │
├────────────────────────────────────┴────────────────────────────────┤
│ Tab panel │ j/k move │ Enter assign │ s schedule │ u unschedule    │
└─────────────────────────────────────────────────────────────────────┘
```

**Panel Proportions**:
- Timeline panel: 60% width
- Unscheduled panel: 40% width
- Minimum terminal width: 100 columns recommended for this view

### Timeline Block Rendering

Block height is proportional to task duration:
- 1 hour = base height (e.g., 4 lines)
- 30 min = half base height (2 lines)
- 2 hours = double base height (8 lines)

Formula: `blockHeight = (durationMinutes / slotMinutes) * linesPerSlot`

### Schedule Input Modal

```
┌────────────────────────────────────┐
│         Schedule Task              │
│                                    │
│  Task: 回覆客戶郵件                 │
│                                    │
│  Time: [09:00-10:30    ]           │
│                                    │
│  Format: HH:MM or HH:MM-HH:MM      │
│  Default duration: 1 hour          │
│                                    │
│      Enter confirm │ Esc cancel    │
└────────────────────────────────────┘
```

## Interaction Flows

### Flow 1: Cursor-Based Scheduling

```
1. User in Day View, focus on Unscheduled panel
2. j/k to select task (highlighted with ">")
3. Tab to switch focus to Timeline panel
4. j/k to move cursor to desired time slot (cursor shown as "▶")
5. Enter to assign selected task to cursor position
6. Task moves from Unscheduled to Timeline with default 1-hour duration
```

### Flow 2: Quick Schedule Input

```
1. User in Day View, focus on Unscheduled panel
2. j/k to select task
3. Press 's' to open schedule input modal
4. Type time: "09:00" (1-hour default) or "09:00-10:30" (specific range)
5. Enter to confirm
6. Task moves from Unscheduled to Timeline
```

### Flow 3: Unschedule Task

```
1. User in Day View, focus on Timeline panel
2. j/k to select scheduled task
3. Press 'u' to unschedule
4. Task moves back to Unscheduled panel (StartTime/EndTime set to nil)
```

### Flow 4: Adjust Duration

```
1. User in Day View, focus on Timeline panel
2. j/k to select scheduled task
3. Press '+' to extend by one slot (30 min)
4. Press '-' to shrink by one slot (30 min, minimum 1 slot)
```

## State Management

### New State Types

```go
// Day View focus area
type DayViewFocus int

const (
    DayViewFocusTimeline DayViewFocus = iota
    DayViewFocusUnscheduled
)

// App state for schedule input
const (
    StateScheduling AppState = iota + existing_states
)
```

### Model Extensions

```go
type model struct {
    // ... existing fields

    // Day View state
    dayViewFocus       DayViewFocus
    timelineCursor     int       // Current time slot index (0 = first slot)
    selectedUnscheduled int      // Selected index in unscheduled list

    // Schedule input
    scheduleInput      textinput.Model
    schedulingTaskIdx  int       // Index of task being scheduled
}
```

## Edge Cases

1. **Overlapping tasks**: Allow overlaps, render side-by-side or stacked
2. **Task extends beyond day end**: Truncate visual display, show indicator
3. **Empty unscheduled list**: Show "All tasks scheduled" message
4. **No scheduled tasks**: Show empty timeline with time markers
5. **Terminal too narrow**: Fall back to single-panel view or show warning

## Testing Strategy

1. **Domain tests**: `Todo.IsScheduled()`, `Todo.Duration()` methods
2. **Service tests**: `ScheduleTodo()`, `UnscheduleTodo()` with mock repository
3. **UI tests**: Golden file tests for timeline rendering at various states
4. **Integration tests**: Full scheduling flow from UI to persistence
