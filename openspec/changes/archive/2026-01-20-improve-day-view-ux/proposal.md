# Proposal: Improve Day View UX

## Summary

Enhance the Day View with better usability through dual viewing modes, improved discoverability, and task movement capabilities.

## Why

The current Day View timeline implementation has several UX issues:
1. **Lack of discoverability** - Users don't know how to use the timeline features (Tab, s, Enter, +/-, u)
2. **No task movement** - Scheduled tasks cannot be easily repositioned on the timeline
3. **Information overload** - Timeline mode is always shown even when users just want to see task details

## What Changes

### 1. Day View Dual Modes (day-view-modes)

Add two distinct viewing modes for Day View:

| Mode | Purpose | Default |
|------|---------|---------|
| **List Mode** | Clean card-style view showing task details (title + description) | Yes |
| **Timeline Mode** | Current split-panel with timeline + unscheduled tasks | No |

- Press `T` to toggle between modes
- List Mode: Focus on content, no time display, sorted by priority
- Timeline Mode: Focus on scheduling, shows time slots

### 2. Day View Discoverability (day-view-discoverability)

Improve feature discovery through:

**A. Dynamic Help Bar**
- Show context-specific shortcuts based on current mode and focus
- List Mode: `j/k nav │ Space done │ e edit │ T timeline │ ...`
- Timeline Mode (Unscheduled): `Tab switch │ s schedule │ Enter assign │ ...`
- Timeline Mode (Timeline): `Tab switch │ J/K move │ +/- duration │ u unschedule │ ...`

**B. Empty State Hints**
- Show helpful hints when panels are empty
- "Press 's' to schedule a task" in empty timeline
- "Press 'a' to add a task" in empty unscheduled list

### 3. Timeline Task Movement (timeline-task-movement)

Enable moving scheduled tasks on the timeline:

- `Shift+J`: Move task down (later time)
- `Shift+K`: Move task up (earlier time)
- Movement granularity: Configurable in TimelineConfig (default: same as SlotMinutes)
- Boundary handling: Cannot move before DayStart or after DayEnd

## Acceptance Criteria

- [ ] Day View defaults to List Mode showing card-style task details
- [ ] Pressing T toggles between List Mode and Timeline Mode
- [ ] Help bar updates based on current mode and focus
- [ ] Empty states show contextual hints
- [ ] Shift+J/K moves scheduled tasks on timeline
- [ ] Movement granularity is configurable
- [ ] All existing Day View functionality remains working
