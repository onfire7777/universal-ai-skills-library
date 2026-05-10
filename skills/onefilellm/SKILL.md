---
name: onefilellm
description: Use when aggregating local files, directories, GitHub repositories/issues/PRs, documentation pages, PDFs, arXiv/DOI/PMID sources, YouTube transcripts, stdin, or clipboard text into a single LLM-ready XML context file using the locally installed jimmc414/onefilellm CLI. Use for context packaging, source bundling, repo-to-text conversion, and multi-source research ingestion.
license: MIT
---

# OneFileLLM

Use this skill to collect multiple sources into one XML-style context payload for LLM analysis.

## Source

- Upstream repository: `https://github.com/jimmc414/onefilellm`
- Integrated source commit: `99c51a2cbe8cc01c0db037a9f800ca31fae9c2cd`
- Local checkout: `%USERPROFILE%\.onefilellm\onefilellm`
- Isolated venv CLI: `%USERPROFILE%\.onefilellm\venv\Scripts\onefilellm.exe`
- Windows shim: `%USERPROFILE%\.onefilellm\bin\onefilellm.cmd`

## README-Grounded Usage

The upstream README defines OneFileLLM as a command-line tool and Python API for aggregating local files, GitHub repos, web pages, PDFs, YouTube transcripts, arXiv/DOI/PMID sources, stdin, and clipboard content into a single structured XML output.

Run:

```bash
python scripts/run_onefilellm.py ./docs README.md
python scripts/run_onefilellm.py https://github.com/user/project
python scripts/run_onefilellm.py https://docs.python.org/3/tutorial/
python scripts/run_onefilellm.py --help-topic examples
```

Equivalent direct commands:

```bash
%USERPROFILE%\.onefilellm\venv\Scripts\onefilellm.exe --help
%USERPROFILE%\.onefilellm\bin\onefilellm.cmd --help-topic crawling
```

## References

- `references/source-readme.md` contains the full upstream README.
- `references/architecture.md` contains upstream architecture notes.
- `references/requirements.txt` and `references/requirements-lock.txt` document installed dependencies.
- `references/source-metadata.json` records the integrated commit and local paths.

## Safety And Quality Rules

- Prefer local files/directories when possible.
- Treat network sources as active fetches. Only run against trusted URLs or explicit user-provided URLs.
- Do not pass secrets, private tokens, or sensitive local directories unless the user explicitly asks and the scope is clear.
- Use `GITHUB_TOKEN` only when private GitHub access or higher rate limits are required.
- Use `OFFLINE_MODE=1` for local-only dry runs or when network fetches should be blocked.
- For YouTube-only transcript work, prefer the `youtube-transcript-english` skill because it enforces English output more strictly than OneFileLLM's built-in transcript fallback.
