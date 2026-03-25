# Naming Examples

Code examples demonstrating meaningful naming principles in TypeScript.

## Bad Examples

### Non-Revealing Names

```typescript
const d = 0; // elapsed time in days

function getThem(): number[][] {
  const list1: number[][] = [];
  for (const x of theList) {
    if (x[0] === 4) list1.push(x);
  }
  return list1;
}
```

**Problems**: `d`, `getThem`, `list1`, `x` reveal nothing; magic number `4` unexplained

### Disinformative Names

```typescript
const accountList = new Set<Account>(); // not actually a list!
let a = l;  // l looks like 1
if (O === l) a = O1;  // O looks like 0
```

**Problems**: `accountList` lies about type; `l` and `O` are visually confusing

### Number-Series and Noise Words

```typescript
function copyChars(a1: string[], a2: string[]): void {
  for (let i = 0; i < a1.length; i++) a2[i] = a1[i];
}

class Product {}
class ProductInfo {}  // What's the difference?
class ProductData {}  // Info and Data add no meaning
```

**Problems**: `a1`/`a2` are meaningless; noise words create false distinctions

### Unpronounceable and Unsearchable

```typescript
class DtaRcrd102 {
  private genymdhms: Date;  // "gen why emm dee aich emm ess"
  private readonly pszqint = "102";
}

for (let j = 0; j < 34; j++) {
  s += (t[j] * 4) / 5;  // Can't search for 4 or 5
}
```

**Problems**: Can't discuss verbally; can't find magic numbers in codebase

### Encoded Names

```typescript
class Part {
  private m_dsc: string;  // m_ prefix is clutter
  setName(name: string): void { this.m_dsc = name; }
}

let strName: string;  // Hungarian notation duplicates type system
let phoneString: PhoneNumber;  // Encoding lies when type changes
```

**Problems**: Prefixes add noise; encodings become lies when types change

## Good Examples

### Intention-Revealing Names

```typescript
const elapsedTimeInDays = 0;
const daysSinceCreation = 0;

function getFlaggedCells(): Cell[] {
  const flaggedCells: Cell[] = [];
  for (const cell of gameBoard) {
    if (cell.isFlagged()) flaggedCells.push(cell);
  }
  return flaggedCells;
}
```

**Why it works**: Names reveal units, purpose, and domain context

### Pronounceable and Searchable

```typescript
class Customer {
  private generationTimestamp: Date;
  private modificationTimestamp: Date;
}

const WORK_DAYS_PER_WEEK = 5;
for (let j = 0; j < NUMBER_OF_TASKS; j++) {
  const realTaskWeeks = taskEstimate[j] / WORK_DAYS_PER_WEEK;
  sum += realTaskWeeks;
}
```

**Why it works**: Can discuss verbally; constants are searchable

### Clean Class Without Encodings

```typescript
class Part {
  private description: string;
  setDescription(description: string): void {
    this.description = description;
  }
}
```

**Why it works**: No prefixes; name matches concept; method matches field

## Refactoring Walkthrough

### Before: Variables Without Context

```typescript
function printGuessStatistics(candidate: string, count: number): void {
  let number: string, verb: string, pluralModifier: string;
  if (count === 0) { number = "no"; verb = "are"; pluralModifier = "s"; }
  else if (count === 1) { number = "1"; verb = "is"; pluralModifier = ""; }
  else { number = count.toString(); verb = "are"; pluralModifier = "s"; }
  console.log(`There ${verb} ${number} ${candidate}${pluralModifier}`);
}
```

### After: Context Through Class

```typescript
class GuessStatisticsMessage {
  private number: string;
  private verb: string;
  private pluralModifier: string;

  make(candidate: string, count: number): string {
    this.createPluralDependentMessageParts(count);
    return `There ${this.verb} ${this.number} ${candidate}${this.pluralModifier}`;
  }

  private createPluralDependentMessageParts(count: number): void {
    if (count === 0) this.thereAreNoLetters();
    else if (count === 1) this.thereIsOneLetter();
    else this.thereAreManyLetters(count);
  }

  private thereAreManyLetters(count: number): void {
    this.number = count.toString(); this.verb = "are"; this.pluralModifier = "s";
  }
  private thereIsOneLetter(): void {
    this.number = "1"; this.verb = "is"; this.pluralModifier = "";
  }
  private thereAreNoLetters(): void {
    this.number = "no"; this.verb = "are"; this.pluralModifier = "s";
  }
}
```

**Changes**: Class provides context; method names reveal intent; each method does one thing

## Method Naming Patterns

```typescript
// Accessors: get prefix
const name = employee.getName();

// Mutators: set prefix
customer.setName("Mike");

// Predicates: is/has prefix
if (paycheck.isPosted()) { }
if (user.hasPermission("admin")) { }

// Factory methods: describe what they create from
const point = Complex.fromRealNumber(23.0);
const user = User.createGuest();
```
