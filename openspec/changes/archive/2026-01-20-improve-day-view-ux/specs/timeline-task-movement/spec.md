# Timeline Task Movement

Enable moving scheduled tasks on the timeline using keyboard shortcuts.

## ADDED Requirements

### Requirement: Timeline Configuration Extension

The TimelineConfig SHALL support a configurable movement granularity.

- TimelineConfig SHALL include a `MoveMinutes` field
- `MoveMinutes` SHALL default to the value of `SlotMinutes` if not specified
- `MoveMinutes` SHALL be persisted in the config file

#### Scenario: Default movement granularity

- Given the user has not configured MoveMinutes
- And SlotMinutes is set to 30
- Then MoveMinutes SHALL default to 30

#### Scenario: Custom movement granularity

- Given the user has configured MoveMinutes to 15
- Then task movement SHALL use 15-minute increments

### Requirement: Move Task Earlier

The system SHALL allow users to move a scheduled task to an earlier time slot.

- Pressing Shift+K SHALL move the selected scheduled task earlier by MoveMinutes
- The task's StartTime and EndTime SHALL both be adjusted by the same amount
- The system SHALL NOT allow moving a task before DayStart
- If movement would exceed boundaries, the system SHALL display an error message

#### Scenario: Move task earlier successfully

- Given the user is in Day View Timeline Mode with Timeline focus
- And a task is scheduled at 10:00-11:00
- And MoveMinutes is 30
- When the user presses Shift+K
- Then the task SHALL be rescheduled to 09:30-10:30

#### Scenario: Cannot move task before DayStart

- Given the user is in Day View Timeline Mode with Timeline focus
- And a task is scheduled at 08:00-09:00
- And DayStart is 08:00
- When the user presses Shift+K
- Then the task SHALL remain at 08:00-09:00
- And the system SHALL display "Cannot move earlier"

### Requirement: Move Task Later

The system SHALL allow users to move a scheduled task to a later time slot.

- Pressing Shift+J SHALL move the selected scheduled task later by MoveMinutes
- The task's StartTime and EndTime SHALL both be adjusted by the same amount
- The system SHALL NOT allow moving a task past DayEnd
- If movement would exceed boundaries, the system SHALL display an error message

#### Scenario: Move task later successfully

- Given the user is in Day View Timeline Mode with Timeline focus
- And a task is scheduled at 10:00-11:00
- And MoveMinutes is 30
- When the user presses Shift+J
- Then the task SHALL be rescheduled to 10:30-11:30

#### Scenario: Cannot move task past DayEnd

- Given the user is in Day View Timeline Mode with Timeline focus
- And a task is scheduled at 17:00-18:00
- And DayEnd is 18:00
- When the user presses Shift+J
- Then the task SHALL remain at 17:00-18:00
- And the system SHALL display "Cannot move later"

### Requirement: Timeline Task Selection

The system SHALL allow users to select scheduled tasks on the timeline for movement.

- The timeline SHALL track a selected scheduled task index
- The selected task SHALL be visually highlighted
- Navigation keys (j/k) in Timeline focus SHALL move between scheduled tasks
- If no scheduled tasks exist, the timeline cursor SHALL control the time slot position

#### Scenario: Navigate between scheduled tasks

- Given the user is in Day View Timeline Mode with Timeline focus
- And there are multiple scheduled tasks
- When the user presses j or k
- Then the selection SHALL move to the next or previous scheduled task

#### Scenario: Selected task highlighting

- Given the user is in Day View Timeline Mode with Timeline focus
- And a task is selected
- Then the selected task SHALL be displayed with a highlighted border or background
