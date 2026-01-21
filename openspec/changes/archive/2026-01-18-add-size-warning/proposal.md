# Proposal: Add Minimum Terminal Size Warning

## Change ID
`add-size-warning`

## Summary
當終端尺寸小於最小建議尺寸時，顯示警告訊息。

## Context
在 80x24 的最小終端尺寸下，日曆格子會很擁擠，Week view 的 todos 會被大量截斷。目前沒有提示使用者終端太小。

## Proposed Solution

### 1. 定義最小尺寸
- 最小寬度：80 columns
- 最小高度：24 rows
- 建議尺寸：100x30

### 2. 警告顯示
當終端尺寸低於最小值時，在畫面中央顯示警告：

```
╭─────────────────────────────────────╮
│  Terminal size too small            │
│                                     │
│  Current: 60x20                     │
│  Minimum: 80x24                     │
│                                     │
│  Please resize your terminal        │
╰─────────────────────────────────────╯
```

## Impact Analysis

### Files to Modify
| File | Change |
|------|--------|
| `cmd/chronoflow/main.go` | 在 View() 中檢查尺寸 |
| `internal/ui/views.go` | 新增 RenderSizeWarning() |

## Acceptance Criteria
1. 終端小於 80x24 時顯示警告
2. 警告顯示當前尺寸和最小需求
3. 終端放大後警告自動消失

## References
- TUI Inspection Report: Low Priority #9
