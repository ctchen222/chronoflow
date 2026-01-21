# day-view-timeline Specification

## Purpose
TBD - created by archiving change enhance-day-view-timeline. Update Purpose after archive.
## Requirements
### Requirement: Split-Panel Day View Layout

The Day View SHALL display a split-panel layout with timeline on the left and unscheduled tasks on the right.

#### Scenario: Day View renders split-panel layout

- **Given** the user is in Day View
- **When** the view renders
- **Then** the left panel (60% width) shows a vertical timeline
- **And** the right panel (40% width) shows unscheduled tasks

#### Scenario: Timeline displays time markers

- **Given** the user is in Day View
- **When** the timeline renders
- **Then** time markers are shown at configured intervals (default 30 min)
- **And** markers display in HH:MM format (e.g., "08:00", "08:30")

#### Scenario: Scheduled tasks appear as blocks on timeline

- **Given** a task with StartTime "09:00" and EndTime "10:30"
- **When** the Day View renders
- **Then** the task appears as a block on the timeline starting at 09:00
- **And** the block height is proportional to the 1.5-hour duration

### Requirement: Unscheduled Task List

The right panel MUST display tasks without scheduled times, showing full details including title and description.

#### Scenario: Unscheduled tasks show title and description

- **Given** a task with title "回覆郵件" and description "確認報價"
- **And** the task has no StartTime set
- **When** the unscheduled panel renders
- **Then** the task title "回覆郵件" is displayed
- **And** the description "確認報價" is displayed below the title

#### Scenario: Unscheduled panel shows count in header

- **Given** 3 tasks without scheduled times
- **When** the Day View renders
- **Then** the unscheduled panel header shows "UNSCHEDULED (3)"

### Requirement: Day View Panel Focus

Users MUST be able to switch focus between timeline and unscheduled panels using Tab key.

#### Scenario: Tab switches focus between panels

- **Given** the user is in Day View with focus on unscheduled panel
- **When** the user presses Tab
- **Then** focus moves to the timeline panel
- **And** the focused panel is visually highlighted

#### Scenario: Timeline cursor navigation

- **Given** the user is in Day View with focus on timeline
- **When** the user presses j/k
- **Then** the timeline cursor moves to the next/previous time slot
- **And** the current cursor position is visually indicated

