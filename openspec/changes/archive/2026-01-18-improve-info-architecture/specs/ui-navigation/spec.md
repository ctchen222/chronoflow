## ADDED Requirements

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
