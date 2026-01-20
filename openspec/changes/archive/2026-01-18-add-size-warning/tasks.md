# Tasks: Add Minimum Terminal Size Warning

## Implementation Order

1. [x] **Define minimum size constants**
   - File: `internal/ui/views.go`
   - MinTerminalWidth = 80, MinTerminalHeight = 24

2. [x] **Implement RenderSizeWarning**
   - File: `internal/ui/views.go`
   - Display current size and minimum requirements
   - Added `IsTooSmall()` helper method

3. [x] **Add size check in View()**
   - File: `cmd/chronoflow/main.go`
   - Check size before rendering normal UI

4. [x] **Testing**
   - Test with various terminal sizes

## Validation

```bash
go test ./...
```
