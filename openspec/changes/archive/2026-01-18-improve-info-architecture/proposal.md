# Change: Improve Information Architecture

## Why

TUI 設計檢查報告 (2026-01-18) 指出資訊架構評分為 4/5，有兩個主要的改進建議：
1. 缺少導航階層的麵包屑顯示（Month > Week）
2. 日曆標題未顯示待辦事項數量

這些改進將幫助使用者更快速地理解目前所在的視圖模式，並提供更好的上下文資訊。

## What Changes

- **Add view mode indicator (breadcrumb)**: 在日曆標題區域顯示目前的視圖模式（MONTH VIEW 或 WEEK VIEW）
- **Add todo count in calendar header**: 在子標題顯示當日的待辦事項數量，讓使用者快速掌握任務量

## Impact

- Affected specs: `specs/ui-navigation` (new capability)
- Affected code:
  - `pkg/calendar/calendar.go:172-212` - Month view header rendering
  - `pkg/calendar/calendar.go:394-420` - Week view header rendering
  - `pkg/calendar/calendar.go:50-57` - Model struct (need to expose todo count)
