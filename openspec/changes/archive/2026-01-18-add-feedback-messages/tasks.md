# Tasks: Add Feedback Messages

## Implementation Order

1. [x] **Add status message fields to model**
   - File: `cmd/chronoflow/main.go`
   - Add `statusMessage string` and `statusType string` fields

2. [x] **Create status message tick command**
   - File: `cmd/chronoflow/main.go`
   - Implement `clearStatusMsg` message type
   - Add `tea.Tick` for 2-second auto-dismiss

3. [x] **Add setStatus helper function**
   - File: `cmd/chronoflow/main.go`
   - Helper to set message and start timer

4. [x] **Add feedback to save action**
   - File: `cmd/chronoflow/main.go`
   - Call setStatus("Todo saved", "success") after save

5. [x] **Add feedback to delete action**
   - File: `cmd/chronoflow/main.go`
   - Call setStatus("Todo deleted", "warning") after delete

6. [x] **Add feedback to toggle action**
   - File: `cmd/chronoflow/main.go`
   - Show "Marked complete" or "Marked incomplete"

7. [x] **Add feedback to priority action**
   - File: `cmd/chronoflow/main.go`
   - Show "Priority: [level]" with appropriate color

8. [x] **Implement RenderStatusMessage**
   - File: `internal/ui/views.go`
   - Render message with appropriate color based on type

9. [x] **Integrate status message into main view**
   - File: `cmd/chronoflow/main.go` (View function)
   - Display above help bar

10. [x] **Testing**
    - Update golden files
    - Manual testing of all feedback scenarios

## Validation

```bash
go test ./...
```
