# /skill

Load a skill on demand through the universal router.

## Usage

```text
/skill <name-or-query>
```

## Behavior

1. If `<name-or-query>` matches an exact skill name, run `skill-router skill <name>` and follow the output.
2. If no exact match is obvious, run `skill-router skill search <name-or-query>`.
3. For a single requirement, load its matching skill. For multi-part work, use `skill-router skills workflow plan` to select and sequence contributions.

Historical command aliases are not installed by default. Use `skill-router skill <name>`.
