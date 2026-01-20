# Tasks: Add Fuzzy Search with Highlighting

## Implementation Order

1. [x] **Update search to be case-insensitive**
   - File: `internal/service/todo_service.go`
   - Already implemented with `strings.ToLower()`

2. [x] **Implement match highlighting**
   - File: `internal/ui/views.go`
   - Added `highlightMatch()` helper function
   - Highlights matched substring with accent color (#7D56F4)

3. [x] **Testing**
   - Test various search patterns
   - Verify highlighting works correctly

## Validation

```bash
go test ./...
```
