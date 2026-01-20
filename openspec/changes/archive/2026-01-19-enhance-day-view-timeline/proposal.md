# Proposal: Enhance Day View with Timeline

## Why

Currently, the Day View only displays task titles without descriptions, and tasks have no time-based scheduling capability. Users cannot:
1. See full task details (descriptions) in the Day View
2. Plan their day by assigning tasks to specific time slots
3. Visualize how their time is allocated throughout the day

This limits ChronoFlow's usefulness as a daily planning tool, forcing users to mentally track when they intend to work on tasks.

## What Changes

### 1. Split-Panel Day View Layout
Transform Day View into a split-panel layout:
- **Left panel (60%)**: Vertical timeline showing scheduled tasks as proportional blocks
- **Right panel (40%)**: List of unscheduled tasks with full details (title + description)

### 2. Todo Scheduling Capability
Extend the Todo domain model to support optional time scheduling:
- Add `StartTime` and `EndTime` fields (nullable, `HH:MM` format)
- Tasks without times remain "unscheduled" and appear in the right panel
- Scheduled tasks display as blocks on the timeline with height proportional to duration

### 3. Timeline Configuration
Allow users to customize the timeline display:
- Configurable day start/end times (default: 08:00-18:00)
- Configurable time slot granularity (default: 30 minutes)
- Settings stored in `~/.chronoflow/config.json`

### 4. Scheduling Interactions
Two methods to schedule tasks:
- **Cursor positioning**: Select task in unscheduled list → move cursor on timeline → Enter to assign
- **Quick input**: Select task → press `s` → type time (e.g., `09:00` or `09:00-10:30`)

### Out of Scope
- Week View changes (deferred for future iteration)
- Recurring task scheduling
- Calendar sync/export
