# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

ChronoFlow is a terminal-based calendar and todo manager built with Go and the Charmbracelet TUI ecosystem (Bubble Tea, Lip Gloss, Bubbles). It features vim-style navigation, multiple calendar views (month/week), priority management, and persistent JSON storage.

**Current Version**: 0.1.0

## Essential Commands

### Development

```bash
# Build binary for current platform
make build

# Build and run locally
make run

# Install to GOPATH/bin
make install
```

### Build & Release

```bash
# Build release binaries for all platforms (creates dist/)
make release VERSION=0.1.0

# Clean build artifacts
make clean
```

## Architecture Overview

### Project Structure

```
chronoflow/
├── cmd/chronoflow/main.go     # Entry point + UI state machine
├── internal/
│   ├── domain/                # Core models (Todo, Priority)
│   ├── repository/            # Data persistence (JSON file)
│   ├── service/               # Business logic (TodoService, StatsCalculator)
│   └── ui/                    # Presentation layer (ViewRenderer, Presenter)
├── pkg/
│   ├── calendar/              # Calendar widget component
│   └── todo/                  # Todo list widget component
├── docs/                      # Documentation
└── Makefile                   # Build automation
```

### Layer Responsibilities

**Domain Layer** (`internal/domain/`):
- Pure data structures with no external dependencies
- `Todo`: Core todo item with title, description, completion, priority
- `Priority`: Enum type (None, Low, Medium, High) with display helpers

**Repository Layer** (`internal/repository/`):
- `TodoRepository` interface for CRUD operations
- `JSONTodoRepository` persists to `~/.chronoflow/todos.json`

**Service Layer** (`internal/service/`):
- `TodoService`: Business logic (toggle, priority, reorder, search)
- `StatsCalculator`: Calculates statistics by period
- `TimeProvider`: Interface for testability (Real + Mock implementations)

**UI Layer** (`internal/ui/`):
- `ViewRenderer`: Renders all UI states (Main, Editing, Delete, Search)
- `TodoPresenter`: Converts domain models to UI list items
- `CalendarAdapter`: Builds todo status map for calendar display
- `state.go`: Defines application states and focus modes

**Package Layer** (`pkg/`):
- `calendar.Model`: Bubble Tea calendar widget (month/week views)
- `todo.Model`: Bubble Tea todo list widget with stats display

### Application States

Defined in `internal/ui/state.go`:

| State | Description |
|-------|-------------|
| `StateViewing` | Main view with calendar + todo panel |
| `StateEditing` | Add/edit todo modal |
| `StateConfirmingDelete` | Delete confirmation modal |
| `StateSearching` | Search modal with results |

### Data Flow

```
User Input → main.go (Update)
       ↓
  TodoService
       ↓
  TodoRepository → ~/.chronoflow/todos.json
       ↓
  StatsCalculator
       ↓
  TodoPresenter / CalendarAdapter
       ↓
  pkg/calendar + pkg/todo
       ↓
  ViewRenderer → Terminal Output
```

## UI Development

### Key UI Files

| File | Purpose |
|------|---------|
| `cmd/chronoflow/main.go` | State machine, keyboard handling, main Update/View |
| `internal/ui/views.go` | All view rendering (RenderMain, RenderEditing, etc.) |
| `internal/ui/state.go` | State and focus type definitions |
| `pkg/calendar/calendar.go` | Calendar component with month/week views |
| `pkg/todo/todo.go` | Todo list component with progress bar |

### Styling

Uses Lip Gloss for terminal styling:
- Primary accent color: `#7D56F4` (purple)
- Focused borders use accent color
- Different colors for priorities (cyan/orange/red)
- Gray for secondary/muted text

### Adding New UI Elements

1. Define new state in `internal/ui/state.go` if needed
2. Add render function in `internal/ui/views.go`
3. Handle keyboard input in `cmd/chronoflow/main.go` Update method
4. Update help bar in `RenderHelpBar()` for new shortcuts

## Testing

### Testing Strategy

- **Domain tests**: Pure unit tests for `Todo.IsOverdue()`, `Priority` methods
- **Service tests**: Unit tests with `MockTimeProvider` for deterministic time
- **Repository tests**: Integration tests with temp files
- **UI tests**: Golden file tests using teatest for visual regression

### Running Tests

```bash
go test ./...
```

## Environment Configuration

### Data Storage

- Location: `~/.chronoflow/todos.json`
- Auto-created on first run
- No external dependencies (no database, no network)

### Requirements

- Go 1.24.2+
- No environment variables required

## Code Style & Conventions

### Go Conventions

- Follow standard Go formatting (`gofmt`)
- Use meaningful variable names
- Keep functions focused and small
- Error handling: return errors, don't panic

### Project-Specific Patterns

- **Dependency Injection**: Services receive dependencies through constructors
- **Interface Abstraction**: Use interfaces for testability (`TimeProvider`, `TodoRepository`)
- **Separation of Concerns**: Domain models have no UI or persistence knowledge
- **Presenter Pattern**: UI formatting separated from Bubble Tea components
- **Adapter Pattern**: `CalendarAdapter` bridges `TodoService` and Calendar widget

### Bubble Tea Patterns

- `Update()` method handles all input and returns new model + command
- `View()` method is pure rendering, no side effects
- Use `tea.Cmd` for async operations
- Keep state in the main model, components receive data via methods

## Important Notes

1. **Go Version**: Requires Go 1.24.2+
2. **No External Services**: Fully offline, file-based storage
3. **Cross-Platform**: Builds for macOS and Linux (ARM64 + AMD64)
4. **Terminal Size**: UI adapts to terminal dimensions automatically
5. **Vim-style Navigation**: h/j/k/l and other vim bindings throughout

## Debugging Tips

1. Use `log.Printf()` to debug - output goes to stderr
2. Run with `bubbletea.WithAltScreen()` disabled for easier debugging
3. Check `~/.chronoflow/todos.json` for data issues
4. Use `go run ./cmd/chronoflow` for quick iteration
5. Terminal must support true color for proper styling

## Security Considerations

1. Data stored locally only - no network access
2. No sensitive information in todo storage
3. File permissions follow user defaults
4. No authentication required (single-user application)

## General Interaction Rules

1. If you are under 95% confidence of what to do or which direction to go, ask clarifying questions.
2. When showing options, display them in hierarchical format:
   ```
   1.0
   1.1
   1.1.1
   2.0
   ```
3. For UI changes, always read the relevant files first before suggesting modifications.
4. Follow existing patterns in the codebase - check neighboring files for conventions.
