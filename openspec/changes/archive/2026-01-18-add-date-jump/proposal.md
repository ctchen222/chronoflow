# Proposal: Add Direct Date Jump

## Change ID
`add-date-jump`

## Summary
新增 `g` 鍵開啟日期輸入 modal，讓使用者可以直接跳轉到指定日期。

## Context
目前使用者只能透過 h/j/k/l 和 b/n 導航，無法快速跳轉到特定日期。

## Proposed Solution

### 1. 新增 StateGoToDate 狀態
處理日期輸入 modal。

### 2. 日期輸入 Modal
```
╭─────── Go to Date ───────╮
│                          │
│  Enter date (YYYY-MM-DD) │
│  ┌────────────────────┐  │
│  │ 2026-02-15         │  │
│  └────────────────────┘  │
│                          │
│  Enter confirm │ Esc cancel │
╰──────────────────────────╯
```

### 3. 日期解析
支援格式：`YYYY-MM-DD`、`MM-DD`（當年）、`DD`（當月）

## Impact Analysis

### Files to Modify
| File | Change |
|------|--------|
| `internal/ui/state.go` | 新增 `StateGoToDate` |
| `internal/ui/views.go` | 新增 `RenderGoToDate()` |
| `cmd/chronoflow/main.go` | 新增 `g` 鍵處理和日期解析 |

## Acceptance Criteria
1. 按 `g` 開啟日期輸入 modal
2. 輸入有效日期後跳轉到該日期
3. 支援多種日期格式
4. 無效日期顯示錯誤訊息

## References
- TUI Inspection Report: Medium Priority #5
