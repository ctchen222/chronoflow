# ui-help Spec Delta

## ADDED Requirements

### Requirement: Todo Panel Help Bar Content

The help bar SHALL display all relevant shortcuts when focused on the todo panel.

#### Scenario: Todo panel help bar shows all priority shortcuts

- **WHEN** the user is focused on the todo panel
- **THEN** the help bar SHALL display `0-3 priority`
- **AND** users can understand that `0` clears priority

#### Scenario: Todo panel help bar shows Esc shortcut

- **WHEN** the user is focused on the todo panel
- **THEN** the help bar SHALL display `Esc back`
- **AND** users can discover how to return to calendar
