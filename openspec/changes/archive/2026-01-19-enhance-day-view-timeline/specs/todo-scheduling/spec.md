# Todo Scheduling

This capability defines task scheduling with time slots.

## ADDED Requirements

### Requirement: Todo Time Fields

Todos MUST support optional start and end times for scheduling purposes.

#### Scenario: Todo with scheduled time

- **Given** a todo with StartTime "09:00" and EndTime "10:30"
- **When** IsScheduled() is called
- **Then** it returns true

#### Scenario: Todo without scheduled time

- **Given** a todo with StartTime nil
- **When** IsScheduled() is called
- **Then** it returns false

#### Scenario: Todo duration calculation

- **Given** a todo with StartTime "09:00" and EndTime "10:30"
- **When** Duration() is called
- **Then** it returns 90 minutes

### Requirement: Cursor-Based Scheduling

Users SHALL be able to schedule tasks by positioning cursor on timeline and pressing Enter.

#### Scenario: Schedule task via cursor positioning

- **Given** a task is selected in the unscheduled panel
- **And** the timeline cursor is at 09:00
- **When** the user presses Enter
- **Then** the task is scheduled with StartTime "09:00"
- **And** the task EndTime is set to "10:00" (default 1-hour duration)
- **And** the task moves from unscheduled panel to timeline

### Requirement: Quick Schedule Input

Users SHALL be able to schedule tasks by entering time directly via the 's' key.

#### Scenario: Schedule task with start time only

- **Given** a task is selected in the unscheduled panel
- **When** the user presses 's'
- **And** enters "09:00"
- **And** presses Enter
- **Then** the task is scheduled with StartTime "09:00" and EndTime "10:00"

#### Scenario: Schedule task with time range

- **Given** a task is selected in the unscheduled panel
- **When** the user presses 's'
- **And** enters "09:00-10:30"
- **And** presses Enter
- **Then** the task is scheduled with StartTime "09:00" and EndTime "10:30"

#### Scenario: Invalid time format shows error

- **Given** the schedule input modal is open
- **When** the user enters "9am"
- **And** presses Enter
- **Then** an error message is displayed
- **And** the modal remains open

### Requirement: Unschedule Task

Users MUST be able to remove scheduling from a task using the 'u' key.

#### Scenario: Unschedule via keyboard shortcut

- **Given** a scheduled task is selected on the timeline
- **When** the user presses 'u'
- **Then** the task StartTime and EndTime are set to nil
- **And** the task moves from timeline to unscheduled panel

### Requirement: Adjust Task Duration

Users SHALL be able to adjust scheduled task duration with '+' and '-' keys.

#### Scenario: Extend task duration

- **Given** a task scheduled from 09:00-10:00 is selected
- **When** the user presses '+'
- **Then** the task EndTime changes to "10:30" (extended by one slot)

#### Scenario: Shrink task duration

- **Given** a task scheduled from 09:00-10:00 is selected
- **When** the user presses '-'
- **Then** the task EndTime changes to "09:30" (shrunk by one slot)

#### Scenario: Minimum duration enforced

- **Given** a task scheduled from 09:00-09:30 is selected (minimum duration)
- **When** the user presses '-'
- **Then** the task EndTime remains "09:30"
- **And** a feedback message indicates minimum duration reached
