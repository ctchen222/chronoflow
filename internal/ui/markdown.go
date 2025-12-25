package ui

import (
	"strings"
	"sync"

	"github.com/charmbracelet/glamour"
	"github.com/charmbracelet/lipgloss"
)

// MarkdownRenderer handles markdown rendering with caching for performance
type MarkdownRenderer struct {
	renderer   *glamour.TermRenderer
	width      int
	cache      string
	cacheInput string
	cacheMutex sync.RWMutex
}

// NewMarkdownRenderer creates a new markdown renderer with the given width
func NewMarkdownRenderer(width int) *MarkdownRenderer {
	r, _ := glamour.NewTermRenderer(
		glamour.WithStylePath("dark"), // Use fixed dark style (faster than WithAutoStyle)
		glamour.WithWordWrap(width),
	)
	return &MarkdownRenderer{
		renderer: r,
		width:    width,
	}
}

// SetWidth updates the renderer width (recreates renderer if width changed)
func (m *MarkdownRenderer) SetWidth(width int) {
	if m.width == width {
		return
	}
	m.width = width
	m.renderer, _ = glamour.NewTermRenderer(
		glamour.WithStylePath("dark"),
		glamour.WithWordWrap(width),
	)
	// Invalidate cache on width change
	m.cacheMutex.Lock()
	m.cache = ""
	m.cacheInput = ""
	m.cacheMutex.Unlock()
}

// Render renders markdown content with caching
func (m *MarkdownRenderer) Render(content string) string {
	if content == "" {
		return lipgloss.NewStyle().
			Foreground(lipgloss.Color("#666")).
			Italic(true).
			Render("Preview will appear here...")
	}

	// Check cache first
	m.cacheMutex.RLock()
	if m.cacheInput == content {
		cached := m.cache
		m.cacheMutex.RUnlock()
		return cached
	}
	m.cacheMutex.RUnlock()

	// Render new content
	rendered, err := m.renderer.Render(content)
	if err != nil {
		return content // Fallback to raw content on error
	}

	// Trim trailing whitespace/newlines
	rendered = strings.TrimSpace(rendered)

	// Update cache
	m.cacheMutex.Lock()
	m.cache = rendered
	m.cacheInput = content
	m.cacheMutex.Unlock()

	return rendered
}
