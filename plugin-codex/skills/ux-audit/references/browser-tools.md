# Browser Tool Reference

## Tool Selection

| Tool | Best For | Auth Support | Setup |
|------|----------|-------------|-------|
| **Chrome MCP** | Authenticated apps | Full — uses your logged-in Chrome session | Claude Code extension in Chrome |
| **Playwright MCP** | Public apps, parallel sessions | Manual login required | MCP plugin |
| **playwright-cli** | Scripted flows, sub-agent tasks | Persistent session state | npm package |

**Default**: Use Chrome MCP when available. It uses your real browser session — OAuth, cookies, and all auth state just work. For unauthenticated or public apps, Playwright is fine.

## Chrome MCP Commands

| Action | Tool |
|--------|------|
| See current tabs | `mcp__claude-in-chrome__tabs_context_mcp` |
| Open new tab | `mcp__claude-in-chrome__tabs_create_mcp` with URL |
| Read page content | `mcp__claude-in-chrome__read_page` |
| Get page text | `mcp__claude-in-chrome__get_page_text` |
| Click element | `mcp__claude-in-chrome__computer` with click action |
| Fill form field | `mcp__claude-in-chrome__form_input` |
| Navigate | `mcp__claude-in-chrome__navigate` |
| Take screenshot | `mcp__claude-in-chrome__computer` with screenshot action |
| Run JavaScript | `mcp__claude-in-chrome__javascript_tool` |
| Resize window | `mcp__claude-in-chrome__resize_window` |
| Record GIF | `mcp__claude-in-chrome__gif_creator` |

**Important**: Call `tabs_context_mcp` first to see what tabs exist. Avoid triggering JS alerts/confirms — they block the extension.

## Playwright MCP Commands

| Action | Tool |
|--------|------|
| Navigate | `mcp__plugin_playwright_playwright__browser_navigate` |
| Take screenshot | `mcp__plugin_playwright_playwright__browser_take_screenshot` |
| Click | `mcp__plugin_playwright_playwright__browser_click` |
| Fill form | `mcp__plugin_playwright_playwright__browser_fill_form` |
| Get page snapshot | `mcp__plugin_playwright_playwright__browser_snapshot` |
| Resize | `mcp__plugin_playwright_playwright__browser_resize` |
| Run code | `mcp__plugin_playwright_playwright__browser_run_code` |

## playwright-cli Commands

```bash
playwright-cli -s=audit open https://app.example.com
playwright-cli -s=audit snapshot          # Get element refs
playwright-cli -s=audit click e5          # Click by ref
playwright-cli -s=audit screenshot --filename=issue-1.png
playwright-cli -s=audit close
```

Use `-s=NAME` for named sessions. Use `--persistent` to keep login state.

## Mobile Viewport Testing

```
# Chrome MCP
mcp__claude-in-chrome__resize_window — width: 375, height: 812

# Playwright MCP
mcp__plugin_playwright_playwright__browser_resize — width: 375, height: 812

# playwright-cli
playwright-cli -s=audit resize 375 812
```
