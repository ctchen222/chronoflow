# Day View Modes

Day View supports two viewing modes to serve different user needs.

## ADDED Requirements

### Requirement: List Mode Display

The system SHALL provide a List Mode for Day View that displays tasks in a clean card-style format.

- List Mode SHALL be the default mode when entering Day View
- Each task card SHALL display the full title on the first line
- Each task card SHALL display the description on subsequent lines (up to 3 lines, truncated with ellipsis)
- Tasks SHALL be sorted by priority (high to low)
- Completed tasks SHALL appear at the bottom of the list with muted styling
- List Mode SHALL NOT display any time information

#### Scenario: User enters Day View

- Given the user is in Month View or Week View
- When the user presses 'd' to enter Day View
- Then the Day View SHALL display in List Mode by default

#### Scenario: List Mode task display

- Given the user is in Day View List Mode
- And there are tasks for the selected date
- Then each task SHALL be displayed as a card with checkbox, priority indicator, title, and description

### Requirement: Timeline Mode Display

The system SHALL provide a Timeline Mode for Day View that displays a time-based scheduling interface.

- Timeline Mode SHALL display a split-panel layout
- The left panel SHALL show time slots with scheduled tasks
- The right panel SHALL show unscheduled tasks
- Timeline Mode SHALL preserve all existing timeline functionality (scheduling, duration adjustment, unscheduling)

#### Scenario: User views Timeline Mode

- Given the user is in Day View Timeline Mode
- Then the left panel SHALL display time slots from DayStart to DayEnd
- And the right panel SHALL display unscheduled tasks

### Requirement: Mode Toggle

The system SHALL allow users to toggle between List Mode and Timeline Mode.

- Pressing 'T' SHALL toggle between List Mode and Timeline Mode
- The current mode SHALL persist while the user remains in Day View
- The mode SHALL reset to List Mode when the user leaves Day View

#### Scenario: Toggle from List to Timeline

- Given the user is in Day View List Mode
- When the user presses 'T'
- Then the view SHALL switch to Timeline Mode

#### Scenario: Toggle from Timeline to List

- Given the user is in Day View Timeline Mode
- When the user presses 'T'
- Then the view SHALL switch to List Mode

#### Scenario: Mode resets on exit

- Given the user is in Day View Timeline Mode
- When the user presses 'm' to switch to Month View
- And then presses 'd' to return to Day View
- Then the Day View SHALL display in List Mode
