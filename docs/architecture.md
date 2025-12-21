# ChronoFlow Architecture

## Project Structure

```
chronoflow/
├── cmd/
│   └── chronoflow/
│       └── main.go              # Application entry point and UI state management
├── internal/
│   ├── domain/                  # Core domain models (no external dependencies)
│   │   ├── priority.go          # Priority type with Icon() and String() methods
│   │   └── todo.go              # Todo struct with IsOverdue() method
│   ├── repository/              # Data persistence layer
│   │   └── todo_repository.go   # TodoRepository interface + JSON implementation
│   ├── service/                 # Business logic layer
│   │   ├── stats_service.go     # StatsCalculator for todo statistics
│   │   ├── time_provider.go     # TimeProvider interface (Real + Mock)
│   │   └── todo_service.go      # TodoService for todo operations
│   └── ui/                      # UI presentation layer
│       ├── calendar_adapter.go  # Adapts todo data for calendar display
│       └── presenter.go         # TodoPresenter and TodoItem for list display
├── pkg/                         # Reusable UI components
│   ├── calendar/
│   │   └── calendar.go          # Calendar widget (Bubble Tea component)
│   └── todo/
│       └── todo.go              # Todo list widget (Bubble Tea component)
├── tests/                       # Test files
├── docs/                        # Documentation
├── go.mod
├── go.sum
├── Makefile
└── README.md
```

## Layer Responsibilities

### Domain Layer (`internal/domain/`)
Pure data structures and business rules with no external dependencies.

- **Todo**: Core todo item with title, description, completion status, and priority
- **Priority**: Enum-like type (None, Low, Medium, High) with display helpers

### Repository Layer (`internal/repository/`)
Data persistence abstraction.

- **TodoRepository interface**: Defines CRUD operations for todos
- **JSONTodoRepository**: Implements persistence to `~/.chronoflow/todos.json`

### Service Layer (`internal/service/`)
Business logic and operations.

- **TimeProvider**: Abstracts time for testability (RealTimeProvider, MockTimeProvider)
- **TodoService**: Wraps repository with business logic (toggle, priority, reorder, search)
- **StatsCalculator**: Calculates statistics (total, completed, overdue) by period

### UI Layer (`internal/ui/`)
Presentation logic separate from Bubble Tea components.

- **TodoPresenter**: Converts service data to UI list items
- **TodoItem**: Wraps domain.Todo with display formatting (implements list.Item)
- **CalendarAdapter**: Builds todo status map for calendar display

### Package Layer (`pkg/`)
Reusable Bubble Tea components.

- **calendar.Model**: Interactive calendar widget with month/week views
- **todo.Model**: Todo list widget with statistics display

## Data Flow

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
         main.go (View) → Terminal Output
```

## Key Design Decisions

1. **Dependency Injection**: Services receive their dependencies through constructors
2. **Interface Abstraction**: TimeProvider and TodoRepository are interfaces for testability
3. **Separation of Concerns**: Domain models have no UI or persistence knowledge
4. **Presenter Pattern**: UI formatting logic separated from Bubble Tea components
5. **Adapter Pattern**: CalendarAdapter bridges TodoService and Calendar widget

## Testing Strategy

- **Domain tests**: Pure unit tests for Todo.IsOverdue(), Priority methods
- **Service tests**: Unit tests with MockTimeProvider for deterministic time
- **Repository tests**: Integration tests with temp files
- **UI tests**: Golden file tests using teatest for visual regression
