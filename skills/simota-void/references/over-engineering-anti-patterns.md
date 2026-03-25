# Over-Engineering Anti-Patterns

Purpose: Use this file when a target feels more elaborate than the problem it solves.

Contents:
- Ten common over-engineering patterns
- YAGNI / KISS / DRY tension and Rule of Three
- Root causes, detection signals, and prevention rules

## 10 Common Anti-Patterns

| ID | Anti-pattern | Symptom | Void question |
|----|--------------|---------|---------------|
| `OE-01` | Premature abstraction | Interface or abstract class used in one place | Is this abstraction used in 2+ real places? |
| `OE-02` | Future-proofing by speculation | Extension points that never get used | Is there a concrete near-term plan? |
| `OE-03` | Pattern worship | Factory / Strategy / Observer where an `if` would do | Is the pattern smaller than the problem? |
| `OE-04` | Homegrown framework | Custom infrastructure where a mature library exists | Why is the standard option insufficient? |
| `OE-05` | Over-configurability | Endless options and flags | Does this option actually change in practice? |
| `OE-06` | Premature optimization | Complexity added before measuring a bottleneck | Do we have performance evidence? |
| `OE-07` | Type-system maze | Deeply nested generics or conditional types | Are types improving comprehension? |
| `OE-08` | Microservice overuse | Tiny services with huge operational cost | Would a modular monolith be enough? |
| `OE-09` | DRY obsession | Coincidental similarity becomes forced coupling | Do these things change for the same reason? |
| `OE-10` | Excessive defensive programming | Internal paths full of redundant checks | Is this actually a system boundary? |

## YAGNI / KISS / DRY Tension

Rules:
- Prefer `YAGNI` over speculative generality.
- Prefer `KISS` when DRY adds indirection without durable payoff.
- Use `DRY` only when duplication changes for the same reason.

### Rule of Three

```text
1st time: write the straightforward implementation
2nd time: tolerate duplication while watching the pattern
3rd time: extract or abstract if the change reason is truly shared
```

## Root Causes

1. Fear of future rework
2. Status from architectural sophistication
3. Cargo-culted best practices
4. Lack of usage or performance evidence
5. Confusing flexibility with value

## Detection Signals

| Signal | Threshold | Meaning |
|--------|-----------|---------|
| Single-use abstraction | `1` implementation | `OE-01` likely |
| Unchanged config options | `>50%` never changed | `OE-05` likely |
| Design discussion vs implementation time | `>50%` of total effort is design debate | over-design likely |
| Generics depth | `3+` nested levels | `OE-07` likely |

## Prevention Rules

- Ask for evidence before adding flexibility.
- Prefer concrete code until the third real repetition.
- Keep configuration only when the default is not enough for a meaningful share of use cases.
- Reject optimization work without measured bottlenecks.
- Review "future use" comments as subtraction candidates.

## Void Use

Use this reference to:
- flag `OE-01` to `OE-10` during `QUESTION`
- map overhead into the `Cognitive Load` dimension during `WEIGH`
- prefer `Pattern Simplification` or `Abstraction Collapse` during `SUBTRACT`

Quality gates:
- single-use abstraction -> warn on `OE-01`
- `"TODO: future use"` -> flag `OE-02`
- `3+` generic nesting levels -> consider simplification
- `50%+` unchanged config options -> consider hardcoding defaults

Sources: [Martin Fowler: YAGNI](https://martinfowler.com/bliki/Yagni.html) · [Sandi Metz: The Wrong Abstraction](https://sandimetz.com/blog/2016/1/20/the-wrong-abstraction) · [Joel Spolsky: Things You Should Never Do](https://www.joelonsoftware.com/2000/04/06/things-you-should-never-do-part-i/)
