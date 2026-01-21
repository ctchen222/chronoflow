# ui-navigation Spec Delta

## ADDED Requirements

### Requirement: Direct Date Jump

The calendar SHALL support direct navigation to a specific date via the `g` key.

#### Scenario: Open date jump modal

- **WHEN** the user presses `g` key in calendar view
- **THEN** a date input modal SHALL be displayed

#### Scenario: Jump to valid date

- **WHEN** the user enters a valid date and presses Enter
- **THEN** the calendar SHALL navigate to that date
- **AND** the modal SHALL close

#### Scenario: Handle invalid date

- **WHEN** the user enters an invalid date
- **THEN** an error message SHALL be displayed
- **AND** the modal SHALL remain open for correction

#### Scenario: Supported date formats

- **WHEN** parsing user input
- **THEN** the following formats SHALL be supported:
  - `YYYY-MM-DD` (full date)
  - `MM-DD` (current year assumed)
  - `DD` (current month and year assumed)
