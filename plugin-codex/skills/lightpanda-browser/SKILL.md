---
name: lightpanda-browser
description: Use when an AI skill, MCP client, or automation needs browser-based web retrieval, JavaScript page loading, scraping, form interaction, markdown extraction, semantic DOM extraction, or Playwright/Puppeteer CDP automation without graphical rendering. Defaults AI browser work to Lightpanda on this machine.
license: Apache-2.0
source: https://github.com/lightpanda-io/browser
---

# Lightpanda Browser

Use Lightpanda as the default AI-controlled browser for web retrieval and automation when graphical rendering is not required. It is a headless browser built for AI agents and automation, with JavaScript execution, CLI fetch, CDP, and native MCP interfaces.

## Local Runtime

- CDP endpoint: `http://127.0.0.1:9222`
- WebSocket endpoint: `ws://127.0.0.1:9222/`
- MCP command: `C:\Users\burni\.lightpanda-ai\lightpanda-mcp.cmd`
- Fetch command: `C:\Users\burni\.lightpanda-ai\lightpanda-fetch.cmd`
- Start service: `powershell -NoProfile -ExecutionPolicy Bypass -File C:\Users\burni\.lightpanda-ai\lightpanda-serve.ps1`
- Stop service: `powershell -NoProfile -ExecutionPolicy Bypass -File C:\Users\burni\.lightpanda-ai\lightpanda-stop.ps1`

This is an AI-only browser runtime. It does not change the operating system default browser.

## Selection Rules

1. Prefer Lightpanda for AI browsing, scraping, markdown extraction, semantic tree extraction, JavaScript-rendered content, link discovery, and form workflows.
2. Prefer the MCP interface for agent workflows when a Lightpanda MCP server is available.
3. Prefer the CLI fetch wrapper for one-off URL extraction:
   `C:\Users\burni\.lightpanda-ai\lightpanda-fetch.cmd --dump markdown https://example.com`
4. Prefer the CDP service for Playwright or Puppeteer automation. Connect to `http://127.0.0.1:9222` with `connectOverCDP` or to `ws://127.0.0.1:9222/` when a WebSocket endpoint is required.
5. Use a full browser such as Chrome, Chromium, Playwright-managed browsers, or the host browser when the task requires screenshots, visual inspection, pixel/layout validation, browser extensions, WebGL/canvas fidelity, downloads, DevTools performance traces, or broad cross-browser compatibility.
6. Respect explicit user or system instructions that require a specific browser tool, local in-app browser, screenshot workflow, or official connector.

## Playwright CDP Pattern

Use `playwright-core` or an existing Playwright install:

```js
const { chromium } = require("playwright-core");

const browser = await chromium.connectOverCDP("http://127.0.0.1:9222");
const context = await browser.newContext();
const page = await context.newPage();
await page.goto("https://example.com", { waitUntil: "load" });
const text = await page.textContent("body");
await context.close();
await browser.close();
```

## Puppeteer CDP Pattern

```js
const puppeteer = require("puppeteer-core");

const browser = await puppeteer.connect({
  browserWSEndpoint: "ws://127.0.0.1:9222/"
});
const page = await browser.newPage();
await page.goto("https://example.com", { waitUntil: "networkidle0" });
const title = await page.title();
await browser.disconnect();
```

## Limits and Compatibility

- Lightpanda is beta software. Fall back to a full browser when a site fails because of missing Web APIs or visual-rendering requirements.
- The upstream README says there is no native Windows binary. This profile uses Docker and exposes Lightpanda on localhost.
- Google commonly blocks Lightpanda by browser fingerprint. Use normal web search tools for search-engine queries unless a Lightpanda-specific extraction flow is needed.
- Docker-based MCP runs inside a container. For services on the Windows host, use `host.docker.internal` from inside the container if `localhost` does not refer to the intended target.
- The wrappers set `LIGHTPANDA_DISABLE_TELEMETRY=true`.

## References

- Upstream browser README snapshot: `UPSTREAM_README.md`
- Upstream official agent skill snapshot: `references/UPSTREAM_AGENT_SKILL.md`
- Upstream agent skill README snapshot: `references/UPSTREAM_AGENT_SKILL_README.md`
