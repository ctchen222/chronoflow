# ui-layout Spec Delta

## ADDED Requirements

### Requirement: Terminal Size Warning

The application SHALL warn users when the terminal size is below minimum requirements.

#### Scenario: Display warning for small terminal

- **WHEN** the terminal width is less than 80 columns
- **OR** the terminal height is less than 24 rows
- **THEN** a warning message SHALL be displayed
- **AND** the warning SHALL show current and minimum sizes

#### Scenario: Warning disappears when resized

- **WHEN** the terminal is resized to meet minimum requirements
- **THEN** the warning SHALL automatically disappear
- **AND** the normal UI SHALL be rendered
