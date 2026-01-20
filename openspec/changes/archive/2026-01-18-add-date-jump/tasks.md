# Tasks: Add Direct Date Jump

## Implementation Order

1. [x] **Add StateGoToDate to state.go**
2. [x] **Add date input field to model**
3. [x] **Implement RenderGoToDate modal**
4. [x] **Add g key handler**
5. [x] **Implement date parsing logic**
   - Supports: YYYY-MM-DD, MM-DD (current year), DD (current month)
6. [x] **Add error handling for invalid dates**
7. [x] **Update help bar with g shortcut**
8. [x] **Add visual regression test**
   - Golden files updated
9. [x] **Manual testing**

## Validation

```bash
go test ./...
```
