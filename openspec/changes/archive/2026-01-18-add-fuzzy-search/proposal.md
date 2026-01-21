# Proposal: Add Fuzzy Search with Highlighting

## Change ID
`add-fuzzy-search`

## Summary
改進搜尋功能，支援模糊搜尋並高亮匹配的文字。

## Context
目前搜尋是精確匹配，沒有高亮顯示匹配部分，使用者難以識別為何結果匹配。

## Proposed Solution

### 1. 模糊搜尋
使用簡單的子字串匹配或更進階的 fuzzy matching 算法。

### 2. 匹配高亮
在搜尋結果中，用不同顏色高亮匹配的部分。

```
Search: "doc"
─────────────
│ Update documentation   ← "doc" highlighted
│ Add Docker support     ← "Doc" highlighted
```

## Impact Analysis

### Files to Modify
| File | Change |
|------|--------|
| `internal/service/todo_service.go` | 改進搜尋邏輯 |
| `internal/ui/views.go` | 新增高亮渲染 |

## Acceptance Criteria
1. 搜尋支援不區分大小寫
2. 匹配文字以不同顏色高亮
3. 搜尋結果按相關性排序（可選）

## References
- TUI Inspection Report: Low Priority #8
