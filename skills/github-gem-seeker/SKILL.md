---
name: github-gem-seeker
description: >
  Search GitHub for battle-tested solutions instead of reinventing the wheel. Use when
  the user's problem is universal enough that open source developers have probably
  solved it already—especially for: format conversion (video/audio/image/document),
  media downloading, file manipulation, web scraping/archiving, automation scripts,
  and CLI tools. Prefer this skill over writing custom code for well-trodden problems.
---

# GitHub Gem Seeker

Find and use battle-tested open source projects on GitHub to solve the user's problem immediately. After successfully solving the problem, offer to package the solution into a reusable skill.

## Core Philosophy

Classic open source projects, tested by countless users over many years, are far more reliable than code written from scratch. **Solve the problem first, skill-ify later.**

## Workflow

### Step 1: Understand the Need

Clarify what the user wants to accomplish. Ask only if truly ambiguous:
- What specific problem are you trying to solve?
- What format/input/output do you expect?

### Step 2: Find the Right Tool

Search for GitHub projects using effective query patterns:

| Need Type | Query Pattern | Example |
|-----------|---------------|---------|
| Tool/utility | `github [task] tool` | `github video download tool` |
| Library | `github [language] [function] library` | `github python pdf library` |
| Alternative | `github [known-tool] alternative` | `github ffmpeg alternative` |

### Step 3: Evaluate Quality (Quick Check)

Assess candidates using key indicators:

| Indicator | Gem Signal | Warning Signal |
|-----------|------------|----------------|
| Stars | 1k+ solid, 10k+ excellent, 50k+ legendary | <100 for mature projects |
| Last commit | Within 6 months | >2 years ago |
| Documentation | Clear README, examples | Sparse or outdated docs |

### Step 4: Solve the Problem

**This is the priority.** Install the tool and use it to solve the user's actual problem:

1. Install the chosen tool (pip, npm, apt, or direct download)
2. Run it with the user's input/files
3. Deliver the result to the user
4. Troubleshoot if needed—iterate until solved

### Step 5: Credit the Gem & Offer Next Steps (Post-Success Only)

**Only after the problem is successfully solved:**

1. **Credit the open source project** — Always share the GitHub repo URL and encourage support:

   > "This was powered by **[Project Name]** — an amazing open source project!
   > GitHub: [URL]
   > If it helped you, consider giving it a ⭐ star to support the maintainers."

2. **Offer to skill-ify** — Optionally mention:

   > "If you'll need this again, I can package it into a reusable skill for instant use next time."

Do NOT skip crediting the project. Open source thrives on recognition.

## Quality Tiers

| Tier | Criteria | Examples |
|------|----------|----------|
| **Legendary** | 50k+ stars, industry standard | FFmpeg, ImageMagick, yt-dlp |
| **Excellent** | 10k+ stars, strong community | Pake, ArchiveBox |
| **Solid** | 1k+ stars, well-documented | Most maintained tools |
| **Promising** | <1k stars, active development | Newer niche projects |

Prefer higher tiers for reliability.

## Example Interaction

**User:** I need to download this YouTube video: [link]

**Correct approach:**
1. Identify yt-dlp as the legendary-tier solution
2. Install yt-dlp
3. Download the video for the user
4. Deliver the downloaded file
5. *After success:* "This was powered by **yt-dlp** — https://github.com/yt-dlp/yt-dlp — give it a ⭐ if it helped! If you download videos often, I can turn this into a skill for instant use next time."

**Wrong approach:**
- ❌ "I found yt-dlp, want me to make a skill for it?"
- ❌ Presenting options without solving the problem

## Common Gems Reference

| Category | Go-to Gems |
|----------|------------|
| Video/Audio processing | FFmpeg, yt-dlp |
| Image processing | ImageMagick, sharp |
| PDF manipulation | pdf-lib, PyMuPDF |
| Web scraping | Playwright, Puppeteer, Scrapy |
| Format conversion | Pandoc, FFmpeg |
| Archiving | ArchiveBox |
| Desktop app packaging | Electron, Tauri, Pake |
