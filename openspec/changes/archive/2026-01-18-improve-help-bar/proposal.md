# Proposal: Improve Help Bar

## Change ID
`improve-help-bar`

## Summary
改進 help bar 以顯示更多快捷鍵，包括目前隱藏的 `0`（清除優先級）和 `Esc`（返回日曆）。

## Context
TUI inspection report 指出部分快捷鍵不易被發現：
- Todo panel 中 `0` 可以清除優先級，但 help bar 只顯示 `1/2/3`
- `Esc` 可從 todo panel 返回日曆，但未顯示在 help bar

## Proposed Solution

### 1. 更新 Todo Focus Help Bar
從：
```
j/k nav │ J/K move │ Space done │ 1/2/3 priority │ / search │ a add │ e edit │ q quit
```

改為：
```
j/k nav │ J/K move │ Space done │ 0-3 priority │ / search │ a add │ e edit │ Esc back │ q quit
```

### 2. 新增 `? help` 提示
在所有 help bar 末端新增 `? help`（配合 help modal 功能）

## Impact Analysis

### Files to Modify
| File | Change |
|------|--------|
| `internal/ui/views.go` | 更新 `RenderHelpBar()` |

### Dependencies
- 無

## Acceptance Criteria
1. Todo panel help bar 顯示 `0-3 priority`
2. Todo panel help bar 顯示 `Esc back`
3. 所有 help bar 顯示 `? help`
4. 視覺回歸測試通過

## References
- TUI Inspection Report: High Priority #2
