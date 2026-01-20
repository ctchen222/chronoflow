# Proposal: Add Shift+Tab Reverse Navigation

## Change ID
`add-shift-tab`

## Summary
新增 Shift+Tab 支援反向面板導航（Todo → Calendar）。

## Context
目前 Tab 只能從 Calendar 切換到 Todo，要返回需要按 Esc。Shift+Tab 是更直覺的反向導航方式。

## Proposed Solution
- Tab: Calendar → Todo
- Shift+Tab: Todo → Calendar

## Impact Analysis

### Files to Modify
| File | Change |
|------|--------|
| `cmd/chronoflow/main.go` | 新增 Shift+Tab 處理 |

## Acceptance Criteria
1. Shift+Tab 從 Todo 切換到 Calendar
2. 在 Calendar 時 Shift+Tab 無動作

## References
- TUI Inspection Report: Low Priority #7
