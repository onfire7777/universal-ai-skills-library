# Diagram Tools Comparison

Purpose: Read this when Mermaid is not sufficient and you need to choose another diagram format or tool.

## Contents

- Tool comparison
- Selection rules

## Comparison

| Tool | Best For | Tradeoff |
|------|----------|----------|
| Mermaid | Default docs-native diagrams | Limited pixel-level control |
| draw.io | Editable diagrams, workshops, handoff artifacts | XML overhead |
| D2 | Layout-heavy text diagrams | Extra toolchain |
| PlantUML | UML-heavy diagrams | Less native in modern docs flows |
| Structurizr | C4 at scale | Stronger setup dependency |
| ASCII | Terminal/comments/accessibility fallback | Low visual density |

## Selection Rules

- Use Mermaid by default.
- Use draw.io when the user needs editable output.
- Use D2 or PlantUML only when the user explicitly requests them or their feature set solves a real limitation.
- Use Structurizr when the task is large-scale C4 architecture and the environment already supports it.
- Use ASCII for plain text, comments, or fallback.

## D2 Advanced Features

| Feature | Description | When useful |
|---------|-------------|-------------|
| Layout engines | dagre (default), ELK, TALA (commercial) | TALA for best auto-layout on complex graphs |
| Animations | Animated connections and transitions | Interactive presentations, live demos |
| Tooltips | Hover-triggered detail panels | Dense diagrams where detail-on-demand helps |
| Code snippets | Embedded syntax-highlighted code blocks | Architecture diagrams that reference implementation |
| LaTeX | Mathematical notation in labels | Algorithm or data-science diagrams |
| Vars & imports | Variables, reusable components, multi-file composition | Large diagram sets with shared elements |

### When to consider D2 over Mermaid

- `100+` nodes where layout quality matters — D2 with ELK/TALA handles dense graphs better.
- Layout control is the primary concern (precise positioning, edge routing).
- Multi-file composition is needed (imports, variables, reusable components).
- **Stay with Mermaid** when: GitHub/GitLab native rendering matters, team familiarity is higher, or the diagram is `<50` nodes.

## Architecture-as-Code Pattern

### Definition

Store diagram source (DSL) in Git alongside code. CI/CD renders diagrams automatically on change.

### Workflow

```
Code change → Update diagram DSL → PR review (diff-friendly) → CI renders → Docs updated
```

### Canvas role

Canvas generates and updates the DSL files. CI/CD rendering and deployment are project-side concerns.

### Recommended structure

```
docs/diagrams/
├── system-context.mmd      # Mermaid source
├── data-flow.mmd
├── deployment.d2           # D2 source (when needed)
└── README.md               # Index and rendering instructions
```

### Tool references

| Tool | Purpose |
|------|---------|
| Swark | Automated code-to-Mermaid generation |
| D2 CLI | `d2 input.d2 output.svg` rendering |
| Mermaid CLI | `mmdc -i input.mmd -o output.svg` rendering |
| GitHub Actions | `mermaid-js/mermaid-cli` action for CI rendering |
