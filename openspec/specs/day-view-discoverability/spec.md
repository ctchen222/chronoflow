# day-view-discoverability Specification

## Purpose
TBD - created by archiving change improve-day-view-ux. Update Purpose after archive.
## Requirements
### Requirement: Dynamic Help Bar for Day View

The system SHALL display context-aware keyboard shortcuts in the help bar based on the current Day View mode and focus.

- In List Mode, the help bar SHALL show: navigation, task actions, and mode toggle shortcuts
- In Timeline Mode with Unscheduled focus, the help bar SHALL show: panel switch, navigation, scheduling shortcuts
- In Timeline Mode with Timeline focus, the help bar SHALL show: panel switch, navigation, movement, duration, and unschedule shortcuts

#### Scenario: Help bar in List Mode

- Given the user is in Day View List Mode
- Then the help bar SHALL display shortcuts including: `j/k nav`, `Space done`, `e edit`, `a add`, `T timeline`

#### Scenario: Help bar in Timeline Unscheduled focus

- Given the user is in Day View Timeline Mode
- And the focus is on the Unscheduled panel
- Then the help bar SHALL display shortcuts including: `Tab switch`, `j/k nav`, `s schedule`, `Enter assign`, `T list`

#### Scenario: Help bar in Timeline focus

- Given the user is in Day View Timeline Mode
- And the focus is on the Timeline panel
- Then the help bar SHALL display shortcuts including: `Tab switch`, `j/k nav`, `J/K move`, `+/- duration`, `u unschedule`, `T list`

### Requirement: Empty State Hints

The system SHALL display helpful hints when Day View panels are empty.

- Empty List Mode SHALL show a hint to add tasks
- Empty Timeline panel SHALL show a hint about scheduling
- Empty Unscheduled panel SHALL show a message indicating all tasks are scheduled

#### Scenario: Empty List Mode

- Given the user is in Day View List Mode
- And there are no tasks for the selected date
- Then the view SHALL display "No tasks for today" with hint "Press 'a' to add a task"

#### Scenario: Empty Timeline panel

- Given the user is in Day View Timeline Mode
- And there are no scheduled tasks
- Then the Timeline panel SHALL display a hint: "Select a task and press Enter" or "Press 's' to schedule"

#### Scenario: Empty Unscheduled panel

- Given the user is in Day View Timeline Mode
- And all tasks are scheduled
- Then the Unscheduled panel SHALL display "All tasks scheduled!" with hint "Press 'a' to add more tasks"

