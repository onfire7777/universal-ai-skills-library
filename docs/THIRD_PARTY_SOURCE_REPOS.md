# Third-Party Source Repos

The universal stack can expose large third-party skill/tool projects without
copying them into every AI client. Keep one upstream checkout, index its skills
read-only, and install only the compact `universal-ai-skills` wrapper into local
agent roots.

## GStack

- Upstream: <https://github.com/garrytan/gstack>
- Local source: `C:\Users\burni\.gstack\gstack`
- Runtime artifacts:
  - `C:\Users\burni\.gstack\gstack\browse\dist\browse.exe`
  - `C:\Users\burni\.gstack\gstack\design\dist\design.exe`
  - `C:\Users\burni\.gstack\gstack\make-pdf\dist\pdf.exe`
- Router sources:
  - `gstack-gbrain`: `~\.gstack\gstack\.gbrain\skills`
  - `gstack-codex`: `~\.gstack\gstack\.agents\skills`
  - `gstack-openclaw`: `~\.gstack\gstack\openclaw\skills`

Use namespaced skill names such as `gstack-review`, `gstack-qa`,
`gstack-ship`, and `gstack-cso` to avoid collisions with generic skills.

## GBrain

- Upstream: <https://github.com/garrytan/gbrain>
- Local source: `C:\Users\burni\gbrain`
- Runtime state: `C:\Users\burni\.gbrain`
- CLI: `gbrain` from `C:\Users\burni\.bun\bin`
- Router sources:
  - `gbrain-source`: `~\gbrain\skills`
  - `gbrain-user`: `~\.gbrain\skills`

GBrain is a persistent local knowledge system. It should be called through its
CLI for brain state operations and through `skill-router skill <name>` for
skillpack instructions.

## Verification

```powershell
skill-router skills sources --refresh
skill-router skill search gstack
skill-router skill search gbrain
skill-router gstack status
skill-router gbrain status
gbrain --version
```

Warnings for missing optional API keys are acceptable when vector embeddings or
provider-backed query expansion are not configured.
