# Proposal: Add Help Modal

## Change ID
`add-help-modal`

## Summary
新增完整的鍵盤參考 modal，讓使用者按 `?` 鍵即可查看所有可用的快捷鍵。這是 TUI inspection report 中評分最低的項目（Help & Discoverability: 3/5）的主要改進。

## Context
目前的 help bar 空間有限，無法顯示所有快捷鍵，導致一些功能不易被發現：
- `Ctrl+P` markdown 預覽
- `J/K` 重新排序 todos
- `0` 清除優先級
- `Esc` 從 todo panel 返回日曆

## Proposed Solution

### 1. 新增 StateHelp 狀態
在 `internal/ui/state.go` 新增 `StateHelp` 狀態。

### 2. 實作 Help Modal UI
顯示分類的鍵盤快捷鍵參考：

```
╭─────────────────── Keyboard Reference ───────────────────╮
│                                                          │
│  Navigation                    Todo Actions              │
│  ──────────                    ────────────              │
│  h/j/k/l    Move cursor        Space/x   Toggle done     │
│  b/n        Prev/next month    a         Add todo        │
│  w          Toggle week view   e/Enter   Edit todo       │
│  d          Day view           d         Delete todo     │
│  m          Month view         1/2/3/0   Set priority    │
│  t          Jump to today      J/K       Reorder         │
│  g          Go to date                                   │
│                                                          │
│  Edit Mode                     General                   │
│  ─────────                     ───────                   │
│  Tab        Switch field       Tab       Switch panel    │
│  Ctrl+P     Preview markdown   Esc       Back/Cancel     │
│  Enter      Save               /         Search          │
│  Esc        Cancel             ?         This help       │
│                                q         Quit            │
│                                                          │
│                    Press any key to close                │
╰──────────────────────────────────────────────────────────╯
```

### 3. 鍵盤處理
- `?` 鍵開啟 help modal
- 任意鍵關閉 modal

## Impact Analysis

### Files to Modify
| File | Change |
|------|--------|
| `internal/ui/state.go` | 新增 `StateHelp` 狀態 |
| `internal/ui/views.go` | 新增 `RenderHelp()` 函數 |
| `cmd/chronoflow/main.go` | 新增 `?` 鍵處理和 StateHelp 邏輯 |

### Dependencies
- 無新依賴

## Acceptance Criteria
1. 按 `?` 鍵顯示 help modal
2. Modal 顯示所有分類的快捷鍵
3. 任意鍵可關閉 modal
4. Help bar 顯示 `?` 快捷鍵提示
5. 視覺回歸測試通過

## References
- TUI Inspection Report: Section 7 (Help & Discoverability)
- High Priority Recommendation #1
