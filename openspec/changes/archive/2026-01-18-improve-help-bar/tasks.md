# Tasks: Improve Help Bar

## Implementation Order

1. [x] **Update Todo panel help bar**
   - File: `internal/ui/views.go`
   - Change `1/2/3 priority` to `0-3 priority`
   - Add `Esc back` before `? help`

2. [x] **Add ? help to all help bars**
   - File: `internal/ui/views.go`
   - Already completed in `add-help-modal` implementation

3. [x] **Update golden files**
   - Run `go test ./cmd/chronoflow -update`

4. [x] **Manual testing**
   - Verify help bar changes in both panels

## Validation

```bash
go test ./...
```
