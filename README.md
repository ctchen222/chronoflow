# Chronoflow

A terminal-based calendar and todo manager built with Go and [Bubble Tea](https://github.com/charmbracelet/bubbletea).

![Demo](https://img.shields.io/badge/version-0.1.0-blue)
![Go](https://img.shields.io/badge/Go-1.21+-00ADD8?logo=go)
![License](https://img.shields.io/badge/license-MIT-green)

## Features

- **Calendar Views** - Switch between month and week views
- **Todo Management** - Create, edit, delete, and complete tasks
- **Priority Levels** - Set High/Medium/Low priority for tasks
- **Progress Tracking** - Visual progress bar and completion statistics
- **Search** - Search across all todos instantly
- **Overdue Detection** - Automatically highlights overdue tasks
- **Vim-style Navigation** - Fast keyboard-driven interface
- **Persistent Storage** - Todos saved locally in JSON format

## Installation

### Homebrew (macOS/Linux)

```bash
brew tap ctchen222/tap
brew install chronoflow
```

### Manual Download

Download the latest release for your platform from [Releases](https://github.com/ctchen222/chronoflow/releases).

```bash
# macOS Apple Silicon
curl -LO https://github.com/ctchen222/chronoflow/releases/download/v0.1.0/chronoflow-0.1.0-darwin-arm64.tar.gz
tar -xzf chronoflow-0.1.0-darwin-arm64.tar.gz
sudo mv chronoflow-darwin-arm64 /usr/local/bin/chronoflow

# macOS Intel
curl -LO https://github.com/ctchen222/chronoflow/releases/download/v0.1.0/chronoflow-0.1.0-darwin-amd64.tar.gz
tar -xzf chronoflow-0.1.0-darwin-amd64.tar.gz
sudo mv chronoflow-darwin-amd64 /usr/local/bin/chronoflow
```

### Build from Source

```bash
git clone https://github.com/ctchen222/chronoflow.git
cd chronoflow
make install
```

## Usage

```bash
chronoflow
```

### Keyboard Shortcuts

#### Calendar (Left Panel)

| Key | Action |
|-----|--------|
| `h` / `l` | Previous / Next day |
| `j` / `k` | Next / Previous week |
| `b` / `n` | Previous / Next month |
| `t` | Jump to today |
| `w` | Toggle week/month view |
| `Tab` | Switch to todo panel |
| `/` | Search todos |

#### Todo List (Right Panel)

| Key | Action |
|-----|--------|
| `j` / `k` | Navigate up/down |
| `J` / `K` | Move todo up/down (reorder) |
| `Space` / `x` | Toggle completion |
| `a` | Add new todo |
| `e` / `Enter` | Edit todo |
| `d` / `Backspace` | Delete todo |
| `1` / `2` / `3` | Set priority (Low/Medium/High) |
| `0` | Remove priority |
| `Esc` | Back to calendar |

#### Edit Mode

| Key | Action |
|-----|--------|
| `Tab` | Switch between title/description |
| `Ctrl+0/1/2/3` | Set priority |
| `Enter` | Save |
| `Esc` | Cancel |

#### Search Mode

| Key | Action |
|-----|--------|
| `Up` / `Down` | Navigate results |
| `Enter` | Jump to selected todo |
| `Esc` | Cancel search |

### General

| Key | Action |
|-----|--------|
| `q` / `Ctrl+C` | Quit |

## Data Storage

Todos are stored in `todos.json` in the current directory. The file is automatically created on first use.

## Development

```bash
# Build
make build

# Run
make run

# Build release binaries
make release VERSION=0.1.0

# Clean
make clean
```

## License

MIT License - see [LICENSE](LICENSE) for details.

## Acknowledgments

Built with:
- [Bubble Tea](https://github.com/charmbracelet/bubbletea) - TUI framework
- [Lip Gloss](https://github.com/charmbracelet/lipgloss) - Style definitions
- [Bubbles](https://github.com/charmbracelet/bubbles) - TUI components
