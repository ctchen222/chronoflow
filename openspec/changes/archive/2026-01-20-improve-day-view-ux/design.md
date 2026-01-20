# Design: Improve Day View UX

## Overview

This design covers three related improvements to the Day View:
1. Dual viewing modes (List/Timeline)
2. Discoverability improvements (help bar + hints)
3. Task movement on timeline

## 1. Day View Dual Modes

### Mode Definitions

**List Mode (Default)**
- Clean, content-focused view
- All tasks displayed as cards with full title and multi-line description
- No time information shown
- Sorted by: Priority (high → low), then completed tasks at bottom
- Single scrollable list

```
┌─ DAY VIEW ─────────────────────────────────────────┐
│              Sunday, January 18, 2026              │
│          📅 January 2026 › Week 3 › Sun 18         │
│                                                    │
│  ╭────────────────────────────────────────────╮    │
│  │ ☐ !!! Complete TUI visual testing          │    │
│  │     Add teatest support for golden file    │    │
│  │     testing and visual regression...       │    │
│  ╰────────────────────────────────────────────╯    │
│                                                    │
│  ╭────────────────────────────────────────────╮    │
│  │ ☐ !! Review pull request                   │    │
│  │     Check the new API changes and...       │    │
│  ╰────────────────────────────────────────────╯    │
│                                                    │
│  ╭────────────────────────────────────────────╮    │
│  │ ☑ Setup project structure                  │    │
│  │     Initialize go modules and deps         │    │
│  ╰────────────────────────────────────────────╯    │
└────────────────────────────────────────────────────┘
```

**Timeline Mode**
- Existing split-panel layout
- Left: Timeline with time slots and scheduled tasks
- Right: Unscheduled tasks list
- Focus on time-based planning

### State Changes

Add to `calendar.Model`:
```go
type DayViewMode int

const (
    DayViewModeList DayViewMode = iota     // Default
    DayViewModeTimeline
)

// In Model struct
dayViewMode DayViewMode
```

### Toggle Behavior

- `T` key toggles `dayViewMode` between List and Timeline
- Mode persists during Day View session
- Resets to List Mode when leaving Day View

## 2. Discoverability Improvements

### Dynamic Help Bar

Update `RenderHelpBar()` to check Day View mode and focus:

| Context | Help Bar Content |
|---------|------------------|
| Day View - List Mode | `j/k nav │ Space done │ e edit │ a add │ T timeline │ ? help` |
| Day View - Timeline (Unscheduled) | `Tab switch │ j/k nav │ s schedule │ Enter assign │ T list │ ? help` |
| Day View - Timeline (Timeline) | `Tab switch │ j/k nav │ J/K move │ +/- duration │ u unschedule │ T list │ ? help` |

### Empty State Hints

Show contextual hints in empty panels:

**List Mode - No Tasks:**
```
        No tasks for today

      Press 'a' to add a task
```

**Timeline Mode - Empty Timeline:**
```
    No scheduled tasks

  Select a task and press Enter
     or press 's' to schedule
```

**Timeline Mode - Empty Unscheduled:**
```
    All tasks scheduled!

   Press 'a' to add more tasks
```

## 3. Timeline Task Movement

### Key Bindings

| Key | Action |
|-----|--------|
| `Shift+J` | Move selected task later (down) |
| `Shift+K` | Move selected task earlier (up) |

### Configuration

Extend `TimelineConfig`:
```go
type TimelineConfig struct {
    DayStart       string // "HH:MM" format
    DayEnd         string // "HH:MM" format
    SlotMinutes    int    // Display granularity (default: 30)
    MoveMinutes    int    // Movement granularity (default: SlotMinutes)
}
```

### Movement Logic

1. Get current task's StartTime and EndTime
2. Calculate new times: `newStart = oldStart ± MoveMinutes`
3. Validate boundaries:
   - `newStart >= DayStart`
   - `newEnd <= DayEnd`
4. Update task via `TodoService.RescheduleTodo()`
5. Refresh view

### Boundary Handling

- Cannot move task before `DayStart` - show status message "Cannot move earlier"
- Cannot move task past `DayEnd` - show status message "Cannot move later"

## Implementation Order

1. **Phase 1: List Mode** - New default view, T toggle
2. **Phase 2: Help Bar** - Dynamic context-aware help
3. **Phase 3: Empty Hints** - Contextual empty states
4. **Phase 4: Task Movement** - Shift+J/K with config

## Testing Strategy

- Update golden files for List Mode rendering
- Add tests for mode toggle behavior
- Add tests for task movement boundaries
- Verify help bar changes per context
