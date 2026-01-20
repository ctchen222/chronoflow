# Tasks: Add Help Modal

## Implementation Order

### Phase 1: State Definition

1. [x] **Add StateHelp to state.go**
   - File: `internal/ui/state.go`
   - Add `StateHelp` constant to AppState

### Phase 2: Help Modal Rendering

2. [x] **Implement RenderHelp function**
   - File: `internal/ui/views.go`
   - Create categorized keyboard reference layout
   - Use accent color border and dimmed background

3. [x] **Update RenderMain to handle StateHelp**
   - File: `cmd/chronoflow/main.go` (in View() method)
   - Render help modal overlay when in StateHelp

### Phase 3: Keyboard Handling

4. [x] **Add ? key handler**
   - File: `cmd/chronoflow/main.go`
   - Handle `?` key in StateViewing to switch to StateHelp

5. [x] **Add any-key-to-close handler**
   - File: `cmd/chronoflow/main.go`
   - Handle any key in StateHelp to return to StateViewing
   - Ctrl+C still quits from help modal

### Phase 4: Help Bar Update

6. [x] **Add ? help to help bar**
   - File: `internal/ui/views.go`
   - Update RenderHelpBar to include `? help`

### Phase 5: Testing

7. [x] **Add visual regression test for help modal**
   - File: `cmd/chronoflow/ui_test.go`
   - Test help modal display

8. [x] **Update existing golden files**
   - Run tests with -update flag

9. [x] **Manual testing**
   - Verify ? opens help
   - Verify any key closes help
   - Verify all shortcuts are listed

## Validation

```bash
go test ./...
go test ./cmd/chronoflow -update
```
