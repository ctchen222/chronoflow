# ui-help Spec Delta

## ADDED Requirements

### Requirement: Help Modal

The application SHALL provide a comprehensive keyboard reference modal accessible via the `?` key.

#### Scenario: Help modal opens on ? key

- **WHEN** the user presses the `?` key in any viewing state
- **THEN** a help modal SHALL be displayed
- **AND** the modal SHALL overlay the current view

#### Scenario: Help modal displays categorized shortcuts

- **WHEN** the help modal is displayed
- **THEN** shortcuts SHALL be organized into categories:
  - Navigation (h/j/k/l, b/n, w, d, m, t, g)
  - Todo Actions (Space, a, e, d, 1/2/3/0, J/K)
  - Edit Mode (Tab, Ctrl+P, Enter, Esc)
  - General (Tab, Esc, /, ?, q)

#### Scenario: Help modal closes on any key

- **WHEN** the help modal is displayed
- **AND** the user presses any key
- **THEN** the modal SHALL close
- **AND** the application SHALL return to the previous state

#### Scenario: Help modal styling

- **WHEN** the help modal is displayed
- **THEN** it SHALL use the accent color (#7D56F4) for the border
- **AND** the background SHALL be dimmed (#333)
- **AND** category headers SHALL be visually distinct

---

### Requirement: Help Indicator in Help Bar

The help bar SHALL indicate the availability of the help modal.

#### Scenario: Help bar shows ? shortcut

- **WHEN** the user is in the main viewing state
- **THEN** the help bar SHALL include `? help` indicator
