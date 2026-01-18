# Project Context

## Purpose
ChronoFlow is a terminal-based calendar and todo manager that provides an efficient, keyboard-driven interface for task management. It combines calendar navigation with todo lists, featuring vim-style controls, priority management, and persistent local storage. Designed for power users who prefer terminal workflows over GUI applications.

**Current Version**: 0.1.0

## Tech Stack
- **Language**: Go 1.24.2+
- **TUI Framework**: Charmbracelet ecosystem
  - Bubble Tea - Elm-inspired TUI framework
  - Lip Gloss - Terminal styling
  - Bubbles - Common TUI components (list, textarea, textinput)
- **Data Storage**: JSON file (`~/.chronoflow/todos.json`)
- **Build System**: Make

## Project Conventions

### Code Style
- Follow standard Go formatting (`gofmt`)
- Use meaningful, descriptive variable names
- Keep functions focused and small (single responsibility)
- Error handling: return errors, don't panic
- Use interfaces for testability (`TimeProvider`, `TodoRepository`)

### Architecture Patterns

**Layered Architecture**:
```
cmd/chronoflow/     → Entry point + UI state machine
internal/domain/    → Pure data models (no dependencies)
internal/repository/→ Data persistence (JSON file)
internal/service/   → Business logic
internal/ui/        → Presentation layer
pkg/                → Reusable TUI components
```

**Design Patterns Used**:
- **Dependency Injection**: Services receive dependencies through constructors
- **Interface Abstraction**: Enables testing with mock implementations
- **Presenter Pattern**: UI formatting separated from Bubble Tea components
- **Adapter Pattern**: `CalendarAdapter` bridges domain models and widgets
- **State Machine**: Application states managed in main.go Update()

**Bubble Tea Patterns**:
- `Update()` handles all input, returns new model + command
- `View()` is pure rendering with no side effects
- Use `tea.Cmd` for async operations
- Keep state in main model, components receive data via methods

### Testing Strategy
- **Domain tests**: Pure unit tests for models (`Todo.IsOverdue()`, `Priority` methods)
- **Service tests**: Unit tests with `MockTimeProvider` for deterministic time
- **Repository tests**: Integration tests with temporary files
- **UI tests**: Golden file tests using teatest for visual regression
- Run tests: `go test ./...`

### Git Workflow
- **Main branch**: `main` (production-ready)
- **Development branch**: `develop` (integration)
- **Feature branches**: `feature/<name>` or `chore/<name>`
- **Commit messages**: Conventional commits style
  - `feat:` for new features
  - `fix:` for bug fixes
  - `chore:` for maintenance
  - `docs:` for documentation

## Domain Context

**Core Concepts**:
- **Todo**: Task with title, description, completion status, priority, and date
- **Priority**: Four levels (None, Low, Medium, High) with visual indicators (`!`, `!!`, `!!!`)
- **Overdue**: Incomplete todo on a past date
- **Period**: Time range for statistics (current week or month based on view mode)

**Application States**:
| State | Description |
|-------|-------------|
| `StateViewing` | Main dual-panel view (calendar + todos) |
| `StateEditing` | Modal for adding/editing todos |
| `StateConfirmingDelete` | Delete confirmation modal |
| `StateSearching` | Search modal with results |

**Data Flow**:
```
User Input → main.go (Update) → TodoService → TodoRepository → JSON file
                                     ↓
                              StatsCalculator
                                     ↓
                         TodoPresenter / CalendarAdapter
                                     ↓
                            pkg/calendar + pkg/todo
                                     ↓
                              ViewRenderer → Terminal
```

## Important Constraints

1. **Offline-only**: No network access, fully local storage
2. **Single-user**: No authentication or multi-user support
3. **Terminal requirements**:
   - Minimum size: 80x24 columns/rows
   - True color support recommended for proper styling
4. **Cross-platform**: macOS and Linux (ARM64 + AMD64)
5. **No external services**: No database, no API dependencies

## External Dependencies

**Go Modules** (key dependencies):
| Package | Purpose |
|---------|---------|
| `github.com/charmbracelet/bubbletea` | TUI framework |
| `github.com/charmbracelet/lipgloss` | Terminal styling |
| `github.com/charmbracelet/bubbles` | TUI components |

**File System**:
- Data directory: `~/.chronoflow/`
- Data file: `~/.chronoflow/todos.json`
- Auto-created on first run

**Build Dependencies**:
- Go 1.24.2+
- Make (for build automation)
