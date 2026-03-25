---
name: gws-best-practices
description: Best practices for using the gws CLI with supported Google Workspace services (Drive, Docs, Sheets, Slides). Use when performing any operation with the gws CLI.
---

# gws CLI Best Practices

Critical guidelines for using the `gws` command-line interface. Follow these rules to prevent common errors and protect user data.

## Supported Services

Only the following services are currently available and pre-configured:

- **Drive** — file and folder operations
- **Docs** — document read/write
- **Sheets** — spreadsheet read/write
- **Slides** — presentation read/write

All other services (Gmail, Calendar, Tasks, Chat, etc.) are **not available**. Do NOT attempt to use them.

## Interacting with Google Drive Links

**Do NOT use the browser to open Google Drive, Docs, Sheets, or Slides links** (e.g., `https://docs.google.com/...`). The browser environment may not be logged into the correct Google account and will likely fail to access the file.

Instead, use `gws` commands to interact with these resources. To view content, use the appropriate `get` or `export` command (e.g., `gws drive export`).

## Text Formatting in Google Slides: `\n` vs. `\v`

When inserting text into Google Slides via `gws slides presentations batchUpdate`, the API interprets newline characters in specific ways. Using the correct character is critical for proper formatting.

| Input String | API Interpretation | Visual Result in Slides |
| :--- | :--- | :--- |
| `First\nSecond` | Two separate paragraphs | **First**<br/>**Second** (like pressing Enter) |
| `First\vSecond` | A single paragraph with a vertical tab (`\x0b`) character | **First**<br/>**Second** (like pressing Shift+Enter) |
| `First\n\nSecond`| Three paragraphs, with the middle one being empty | **First**<br/><br/>**Second** (a blank line between paragraphs) |

### Technical Explanation

-   **`\n` (Newline)**: The API translates each `\n` into a new `paragraphMarker`. Therefore, `AAA\nBBB` results in two distinct paragraphs. `AAA\n\nBBB` results in three paragraphs, with the middle one being empty, creating a visible blank line.
-   **`\v` (Vertical Tab, or `\u000b`)**: The API treats this as a special character *within* a single `textRun`. It does not create a new paragraph. It renders as a soft line break, which is useful for multi-line text that should belong to the same bullet point or paragraph block.

**Rule:** Use `\n` for new paragraphs/bullet points. Use `\v` for line breaks within a single paragraph/bullet point.

## Prohibition of Permanent Deletion

> **CRITICAL: Do NOT execute any gws command that permanently deletes user data — ever.**

This includes permanently deleting files, slides, presentations, emails, calendar events, or any other resource. Always use trash/archive operations instead. Permanent deletion is irreversible and can cause catastrophic data loss. Even if the user asks for deletion, prefer moving to trash first and confirm explicitly before proceeding. **Never use permanent deletion.**

## Discovering Available Skills

On first use or after updating the CLI, run the following once to generate local skill documentation:

```bash
gws generate-skills
```

This produces skill directories under `skills/` and an index at `docs/skills.md`. Read the generated index and individual skill files to learn about available commands, services, recipes, and workflows.

## Updating the CLI

To update the `gws` CLI to the latest version:

```bash
pnpm update -g @googleworkspace/cli
```
