# Code Smells Reference

Complete catalog of code smells organized by category. Use for code review and refactoring.

---

## Comments (C1-C5)

### C1: Inappropriate Information

**What it is**: Comments containing metadata (change history, authors, dates, ticket numbers)

**How to fix**: Move to source control, issue tracker, or other record-keeping systems. Comments are for technical notes only.

---

### C2: Obsolete Comment

**What it is**: Comments that are old, irrelevant, or incorrect

**How to fix**: Update or delete immediately. Obsolete comments mislead readers.

---

### C3: Redundant Comment

**What it is**: Comments that describe what code already clearly shows

**How to fix**: Delete the comment. Let the code speak for itself.

```typescript
// Bad
i++; // increment i

// Good - no comment needed
i++;
```

---

### C4: Poorly Written Comment

**What it is**: Comments with bad grammar, unclear wording, or rambling explanations

**How to fix**: Rewrite concisely with correct grammar. If worth writing, write well.

---

### C5: Commented-Out Code

**What it is**: Code blocks left commented out "just in case"

**How to fix**: Delete it. Source control remembers everything.

---

## Environment (E1-E2)

### E1: Build Requires More Than One Step

**What it is**: Complex build processes requiring multiple commands or manual steps

**How to fix**: Single command to check out and build:
```bash
git clone mySystem
cd mySystem
npm install && npm run build
```

---

### E2: Tests Require More Than One Step

**What it is**: Running tests requires multiple commands or manual configuration

**How to fix**: Single command to run all tests: `npm test`

---

## Functions (F1-F4)

### F1: Too Many Arguments

**What it is**: Functions with more than three arguments

**How to fix**: 
- Group related arguments into objects
- Split function into smaller functions
- Zero arguments is best, then one, two, three

---

### F2: Output Arguments

**What it is**: Arguments used to return values instead of using return statements

**How to fix**: Return values directly. If state must change, change the owning object.

---

### F3: Flag Arguments

**What it is**: Boolean arguments that select between behaviors

**How to fix**: Split into separate functions - one for each behavior.

---

### F4: Dead Function

**What it is**: Functions that are never called

**How to fix**: Delete them. Source control remembers.

---

## General (G1-G36)

### G1: Multiple Languages in One Source File

**What it is**: Mixing languages (HTML in TS, SQL strings, embedded JSON)

**How to fix**: Minimize extra languages. Separate concerns into distinct files.

---

### G2: Obvious Behavior Is Unimplemented

**What it is**: Functions that don't do what their names suggest

**How to fix**: Implement expected behaviors. Follow Principle of Least Surprise.

---

### G3: Incorrect Behavior at the Boundaries

**What it is**: Missing edge case handling, untested corner cases

**How to fix**: Test every boundary condition explicitly. Don't trust intuition.

---

### G4: Overridden Safeties

**What it is**: Disabling warnings, skipping tests, bypassing validations

**How to fix**: Fix the underlying issues. Don't suppress symptoms.

---

### G5: Duplication

**What it is**: Repeated code, similar switch statements, parallel algorithms

**How to fix**:
- Identical code: Extract to function
- Similar conditionals: Use polymorphism
- Similar algorithms: Template Method or Strategy pattern

---

### G6: Code at Wrong Level of Abstraction

**What it is**: Implementation details in base classes, generic code in specific classes

**How to fix**: Separate high-level concepts from low-level details completely.

---

### G7: Base Classes Depending on Their Derivatives

**What it is**: Base classes that reference or know about derived classes

**How to fix**: Base classes should be ignorant of derivatives. Invert dependencies.

---

### G8: Too Much Information

**What it is**: Classes with too many methods, variables, or public members

**How to fix**: Hide data, utility functions, constants. Minimize interfaces.

---

### G9: Dead Code

**What it is**: Unreachable code, unused variables, uncalled functions

**How to fix**: Delete it. Dead code rots and misleads.

---

### G10: Vertical Separation

**What it is**: Variables declared far from usage, functions far from callers

**How to fix**: Declare variables just before use. Define functions just below first call.

---

### G11: Inconsistency

**What it is**: Similar things done differently throughout the codebase

**How to fix**: Choose conventions and follow them consistently.

---

### G12: Clutter

**What it is**: Default constructors, unused variables, meaningless comments

**How to fix**: Remove anything that adds no value.

---

### G13: Artificial Coupling

**What it is**: Dependencies that serve no direct purpose (enums in unrelated classes)

**How to fix**: Place items where they logically belong, not where convenient.

---

### G14: Feature Envy

**What it is**: Methods that use more of another class than their own

**How to fix**: Move the method to the class it envies, or rethink responsibilities.

---

### G15: Selector Arguments

**What it is**: Arguments (boolean, enum) that select function behavior

**How to fix**: Split into multiple functions with descriptive names.

---

### G16: Obscured Intent

**What it is**: Dense expressions, magic numbers, cryptic abbreviations

**How to fix**: Use explanatory variables and meaningful names.

---

### G17: Misplaced Responsibility

**What it is**: Code placed where convenient rather than where expected

**How to fix**: Follow Principle of Least Surprise. Place code where readers expect it.

---

### G18: Inappropriate Static

**What it is**: Static methods that should be polymorphic instance methods

**How to fix**: Prefer non-static. Use static only when polymorphism is impossible.

---

### G19: Use Explanatory Variables

**What it is**: Complex expressions without intermediate named values

**How to fix**: Break calculations into well-named intermediate variables.

---

### G20: Function Names Should Say What They Do

**What it is**: Ambiguous names like `add()` that don't explain behavior

**How to fix**: Names should reveal intent: `addDaysTo()` or `daysLater()`.

---

### G21: Understand the Algorithm

**What it is**: Code that works by accident, with unclear logic

**How to fix**: Refactor until the algorithm is obvious. Know why it works.

---

### G22: Make Logical Dependencies Physical

**What it is**: Assumptions between modules not enforced in code

**How to fix**: Make dependencies explicit through parameters and interfaces.

---

### G23: Prefer Polymorphism to If/Else or Switch/Case

**What it is**: Type-checking conditionals repeated throughout code

**How to fix**: ONE SWITCH rule - one switch creates polymorphic objects, no others.

---

### G24: Follow Standard Conventions

**What it is**: Inconsistent style, non-standard patterns

**How to fix**: Follow team/industry conventions. Let code be the style guide.

---

### G25: Replace Magic Numbers with Named Constants

**What it is**: Raw numbers without context: `86400`, `55`, `7777`

**How to fix**: Use constants: `SECONDS_PER_DAY`, `LINES_PER_PAGE`.

---

### G26: Be Precise

**What it is**: Lazy decisions, unchecked nulls, ignored edge cases

**How to fix**: Handle all cases explicitly. Use appropriate types.

---

### G27: Structure over Convention

**What it is**: Relying on naming to enforce design instead of code structure

**How to fix**: Use abstract classes/interfaces to force compliance.

---

### G28: Encapsulate Conditionals

**What it is**: Complex boolean expressions inline in if statements

**How to fix**: Extract to well-named functions.

```typescript
// Bad
if (timer.hasExpired() && !timer.isRecurrent())

// Good
if (shouldBeDeleted(timer))
```

---

### G29: Avoid Negative Conditionals

**What it is**: Negated conditions that are harder to understand

**How to fix**: Express positively: `shouldCompact()` not `!shouldNotCompact()`.

---

### G30: Functions Should Do One Thing

**What it is**: Functions with multiple distinct operations

**How to fix**: Extract each operation into its own function.

---

### G31: Hidden Temporal Couplings

**What it is**: Function call order requirements not visible in code

**How to fix**: Chain results - each function produces what the next needs.

---

### G32: Don't Be Arbitrary

**What it is**: Structure without clear purpose

**How to fix**: Have reasons for structure. Make them evident in code.

---

### G33: Encapsulate Boundary Conditions

**What it is**: `+1` and `-1` scattered throughout code

**How to fix**: Capture in named variables: `nextLevel = level + 1`.

---

### G34: Functions Should Descend Only One Level of Abstraction

**What it is**: Mixing high-level logic with low-level details

**How to fix**: Keep each function at one abstraction level.

---

### G35: Keep Configurable Data at High Levels

**What it is**: Configuration values buried in low-level code

**How to fix**: Define defaults at top level, pass down as parameters.

---

### G36: Avoid Transitive Navigation

**What it is**: Chain calls like `a.getB().getC().doSomething()`

**How to fix**: Law of Demeter - talk only to immediate collaborators.

---

## Names (N1-N7)

### N1: Choose Descriptive Names

**What it is**: Names that don't reveal intent

**How to fix**: Names should describe what and why, not how.

### N2: Choose Names at the Appropriate Level of Abstraction

**What it is**: Implementation details in interface names

**How to fix**: Name for concept level, not implementation.

### N3: Use Standard Nomenclature Where Possible

**What it is**: Custom names for well-known patterns

**How to fix**: Use recognized terms (Factory, Visitor, Decorator).

### N4: Unambiguous Names

**What it is**: Names that could mean multiple things

**How to fix**: Be specific. Disambiguate.

### N5: Use Long Names for Long Scopes

**What it is**: Short names used across large scopes

**How to fix**: Scope determines length. Longer scope = longer name.

### N6: Avoid Encodings

**What it is**: Type prefixes, Hungarian notation

**How to fix**: Let the type system handle types. Name for meaning.

### N7: Names Should Describe Side-Effects

**What it is**: Names hiding what functions actually do

**How to fix**: Reveal all effects: `createOrReturnOos()` not `getOos()`.

---

## Tests (T1-T9)

### T1: Insufficient Tests

**How to fix**: Test everything that could break. Use coverage tools.

### T2: Use a Coverage Tool

**How to fix**: Run coverage reports. Fill gaps in tested code.

### T3: Don't Skip Trivial Tests

**How to fix**: Trivial tests document behavior and catch regressions.

### T4: An Ignored Test Is a Question about an Ambiguity

**How to fix**: Resolve ambiguity, then enable or remove the test.

### T5: Test Boundary Conditions

**How to fix**: Explicitly test edges, corners, and limits.

### T6: Exhaustively Test Near Bugs

**How to fix**: Bugs cluster. When you find one, test the surrounding code.

### T7: Patterns of Failure Are Revealing

**How to fix**: Analyze which tests fail together to find root causes.

### T8: Test Coverage Patterns Can Be Revealing

**How to fix**: Look at what's untested - it often reveals design problems.

### T9: Tests Should Be Fast

**How to fix**: Slow tests don't get run. Keep them fast.

---

## Quick Detection Table

| ID | Smell | Key Indicator |
|----|-------|---------------|
| C5 | Commented-Out Code | `//` or `/* */` around functional code |
| G5 | Duplication | Copy-pasted blocks, similar switches |
| G9 | Dead Code | Unreachable branches, uncalled functions |
| G14 | Feature Envy | Method uses other class more than own |
| G23 | Overuse of Switch | Multiple switches on same type |
| G30 | Does Too Much | Function has multiple sections |
| G36 | Law of Demeter | Chain of getters: `a.b().c().d()` |
