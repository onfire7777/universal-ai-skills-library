# /skill

Load a skill on demand through the universal router.

## Usage

```text
/skill <name-or-query>
```

## Behavior

1. If `<name-or-query>` matches an exact skill name, run `skill-router skill <name>` and follow the output.
2. If no exact match is obvious, run `skill-router skill search <name-or-query>`.
3. Load only the best matching skill needed for the current task.

Historical command aliases are not installed by default. Use `skill-router skill <name>`.
