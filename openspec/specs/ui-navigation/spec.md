# ui-navigation Specification

## Purpose
TBD - created by archiving change improve-info-architecture. Update Purpose after archive.
## Requirements
### Requirement: View Mode Indicator

The calendar component SHALL display a view mode indicator that shows the current navigation context (Month View or Week View).

#### Scenario: Month view indicator displayed

- **WHEN** the user is viewing the calendar in month view
- **THEN** a "MONTH VIEW" indicator SHALL be visible in the calendar header area

#### Scenario: Week view indicator displayed

- **WHEN** the user is viewing the calendar in week view
- **THEN** a "WEEK VIEW" indicator SHALL be visible in the calendar header area

#### Scenario: View mode indicator styling

- **WHEN** the view mode indicator is displayed
- **THEN** it SHALL use the accent color (#7D56F4) for visual consistency
- **AND** it SHALL be positioned near the main header for easy visibility

---

### Requirement: Todo Count in Calendar Header

The calendar component SHALL display the number of todos for the currently selected date in the sub-header.

#### Scenario: Todo count displayed for date with tasks

- **WHEN** the user navigates to a date that has todos
- **THEN** the sub-header SHALL show the todo count in format "(N tasks)" or "(1 task)"
- **AND** the count SHALL appear after the date string

#### Scenario: No count displayed for date without tasks

- **WHEN** the user navigates to a date that has no todos
- **THEN** the sub-header SHALL NOT display any count indicator
- **AND** only the date string SHALL be shown

#### Scenario: Todo count updates on navigation

- **WHEN** the user navigates to a different date
- **THEN** the todo count SHALL update to reflect the new date's tasks

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

### Requirement: Reverse Panel Navigation

The application SHALL support Shift+Tab for reverse panel navigation between Calendar and Todo panels.

#### Scenario: Shift+Tab switches from Todo to Calendar

- **WHEN** the user is focused on the Todo panel
- **AND** the user presses Shift+Tab
- **THEN** focus SHALL move to the Calendar panel

