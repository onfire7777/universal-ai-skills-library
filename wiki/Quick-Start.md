# Quick Start

Get the `skill-router` binary installed and routing in a couple of minutes.

Source: [`docs/QUICKSTART.md`](https://github.com/onfire7777/universal-ai-skills-library/blob/main/docs/QUICKSTART.md)

---

## Prerequisites

- **Go** `1.25+` (the installer builds `skill-router` from source).
- **Git**.
- Optional: **Node.js 22** if you'll work on the [registry generator](Manifest-and-Registry.md) or run parity checks.

---

## 1. Clone

```bash
git clone https://github.com/onfire7777/universal-ai-skills-library.git
cd universal-ai-skills-library
```

## 2. Install

**macOS / Linux / WSL:**
```bash
bash install.sh
```

**Windows (PowerShell):**
```powershell
.\install.ps1
```

Want only the router, not the full local AI stack? On Windows use `.\install.ps1 -SkipStackInstall`. See [Installation & Setup](Installation-and-Setup.md) for all flags and install modes.

## 3. Verify

```bash
skill-router --version
skill-router doctor          # diagnose the install
skill-router doctor --json   # machine-readable
```

## 4. Try routing

```bash
# Search the corpus
skill-router skill search debugging

# Load one skill (prints its SKILL.md)
skill-router skill universal-ai-skills

# Ask the router to decide
skill-router route --explain "help me write pytest fixtures"

# The deterministic preflight decision, as JSON
skill-router preflight --json "help me write pytest fixtures"
```

## 5. Validate the catalog

```bash
skill-router skills validate-manifest
```

---

## Wire it into your AI client

The point of UASL is that your agent calls the router automatically. The universal contract every client follows:

1. On a real user prompt, run
   `skill-router preflight --hook-event UserPromptSubmit --json "<prompt>"`.
2. If the decision is **`route`**, load one skill with `skill-router skill <name>`.
3. If **`ambiguous`**, let the host AI pick from the listed candidates.
4. If **`no_route`**, continue normally.

To install the compact wrapper into every detected client:

```bash
skill-router sync installed     # wrapper-only, safe
skill-router sync matrix        # preview what's detected (read-only)
```

Per-client details (Codex, Claude, Cursor, Hermes, Paperclip, …) are on the [Agent Support Matrix](Agent-Support-Matrix.md) and [Installation & Setup](Installation-and-Setup.md) pages.

---

## Next steps

- [Installation & Setup](Installation-and-Setup.md) — install modes, flags, where things land
- [Skill Router CLI](Skill-Router-CLI.md) — every command
- [Architecture](Architecture.md) — how routing actually works
- [Skills Corpus](Skills-Corpus.md) — what a skill is
