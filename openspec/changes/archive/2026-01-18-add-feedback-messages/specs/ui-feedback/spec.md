# ui-feedback Spec Delta

## ADDED Requirements

### Requirement: Status Messages

The application SHALL display brief status messages after user actions to provide feedback.

#### Scenario: Save feedback

- **WHEN** the user saves a todo
- **THEN** a "Todo saved" message SHALL be displayed
- **AND** the message SHALL use green color (#50FA7B)

#### Scenario: Delete feedback

- **WHEN** the user deletes a todo
- **THEN** a "Todo deleted" message SHALL be displayed
- **AND** the message SHALL use orange color (#FFB86C)

#### Scenario: Toggle complete feedback

- **WHEN** the user toggles a todo's completion status
- **THEN** "Marked complete" or "Marked incomplete" SHALL be displayed
- **AND** complete SHALL use green, incomplete SHALL use gray

#### Scenario: Priority change feedback

- **WHEN** the user changes a todo's priority
- **THEN** "Priority: High/Medium/Low/None" SHALL be displayed
- **AND** the message SHALL use the priority's color

#### Scenario: Message auto-dismiss

- **WHEN** a status message is displayed
- **THEN** it SHALL automatically disappear after 2 seconds
