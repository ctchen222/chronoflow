# Proposal: Add Feedback Messages

## Change ID
`add-feedback-messages`

## Summary
新增操作回饋訊息，讓使用者在儲存、刪除、切換完成狀態等操作後收到視覺確認。

## Context
TUI inspection report 指出：
- 儲存 todo 後沒有確認訊息
- 切換完成狀態時沒有視覺回饋
- 刪除確認後是靜默成功

## Proposed Solution

### 1. 新增 Status Message 組件
在 help bar 上方或 panel 底部顯示短暫的狀態訊息。

### 2. 訊息類型
| 操作 | 訊息 | 顏色 |
|------|------|------|
| 儲存 todo | "Todo saved" | Green |
| 刪除 todo | "Todo deleted" | Orange |
| 切換完成 | "Marked complete" / "Marked incomplete" | Green/Gray |
| 優先級變更 | "Priority set to High/Medium/Low/None" | Priority color |

### 3. 自動消失
訊息顯示 2 秒後自動消失。

## Impact Analysis

### Files to Modify
| File | Change |
|------|--------|
| `cmd/chronoflow/main.go` | 新增 statusMessage 和 statusTimer 欄位 |
| `cmd/chronoflow/main.go` | 在操作完成後設定訊息 |
| `internal/ui/views.go` | 新增 RenderStatusMessage() |
| `internal/ui/views.go` | 更新 RenderMain() 包含狀態訊息 |

### Dependencies
- 需要使用 `tea.Tick` 實作自動消失計時器

## Acceptance Criteria
1. 儲存 todo 後顯示 "Todo saved"
2. 刪除 todo 後顯示 "Todo deleted"
3. 切換完成狀態顯示對應訊息
4. 訊息 2 秒後自動消失
5. 訊息使用適當的顏色編碼

## References
- TUI Inspection Report: Section 4 (Feedback & Affordances)
- Medium Priority #4
