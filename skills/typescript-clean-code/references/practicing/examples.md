# Practicing Examples

Practice exercises and kata examples for deliberate skill development.

## Classic Kata Examples

### The Bowling Game Kata

A 30-minute TDD exercise that calculates bowling scores.

**What you practice**:
- TDD red-green-refactor cycle
- Incremental design
- Handling edge cases (strikes, spares)

```typescript
// Start: Write a failing test
describe('BowlingGame', () => {
  let game: BowlingGame;

  beforeEach(() => {
    game = new BowlingGame();
  });

  it('should score zero for gutter game', () => {
    rollMany(20, 0);
    expect(game.score()).toBe(0);
  });

  it('should score 20 for all ones', () => {
    rollMany(20, 1);
    expect(game.score()).toBe(20);
  });

  // Continue with spare, strike, perfect game...
});
```

### Prime Factors Kata

Decompose a number into its prime factors.

**What you practice**:
- Algorithm development through TDD
- Simple incremental steps
- Refactoring as patterns emerge

```typescript
// Tests drive the implementation
describe('PrimeFactors', () => {
  it('returns empty for 1', () => {
    expect(primeFactors(1)).toEqual([]);
  });

  it('returns [2] for 2', () => {
    expect(primeFactors(2)).toEqual([2]);
  });

  it('returns [3] for 3', () => {
    expect(primeFactors(3)).toEqual([3]);
  });

  it('returns [2, 2] for 4', () => {
    expect(primeFactors(4)).toEqual([2, 2]);
  });

  // Pattern emerges, refactor...
});
```

### Word Wrap Kata

Wrap text at a given column width.

**What you practice**:
- String manipulation
- Edge case handling
- Clean function design

```typescript
describe('WordWrap', () => {
  it('returns empty for null', () => {
    expect(wrap(null, 10)).toBe('');
  });

  it('returns text unchanged if shorter than width', () => {
    expect(wrap('hello', 10)).toBe('hello');
  });

  it('wraps at word boundary', () => {
    expect(wrap('hello world', 7)).toBe('hello\nworld');
  });
});
```

## Wasa (Ping-Pong) Session Example

### Setup

Two developers, one problem, alternating roles.

**Round 1**: Developer A writes test
```typescript
it('should return fizz for multiples of 3', () => {
  expect(fizzBuzz(3)).toBe('fizz');
});
```

**Round 2**: Developer B makes it pass, writes next test
```typescript
function fizzBuzz(n: number): string {
  if (n % 3 === 0) return 'fizz';
  return String(n);
}

it('should return buzz for multiples of 5', () => {
  expect(fizzBuzz(5)).toBe('buzz');
});
```

**Round 3**: Developer A makes it pass, writes next test
```typescript
function fizzBuzz(n: number): string {
  if (n % 3 === 0) return 'fizz';
  if (n % 5 === 0) return 'buzz';
  return String(n);
}

it('should return fizzbuzz for multiples of 15', () => {
  expect(fizzBuzz(15)).toBe('fizzbuzz');
});
```

## Practice Routines

### Daily Warm-up (15 minutes)

```
1. Pick a familiar kata
2. Set a timer for 15 minutes
3. Focus on smooth keystrokes and minimal mouse use
4. Track how far you get each day
```

### Weekly Deep Practice (1-2 hours)

```
1. Choose a kata you don't know well
2. Work through it slowly, understanding each step
3. Repeat 3-5 times
4. By end of session, aim for fluent execution
```

### Monthly Challenge

```
1. Learn a new kata in a language you don't use at work
2. Contribute to an open source project
3. Attend or host a coding dojo session
```

## Kata Resources

| Resource | URL |
|----------|-----|
| Kata Catalog | katas.softwarecraftsmanship.org |
| Code Kata | codekata.pragprog.com |
| Coding Dojo | codingdojo.org |

## Good Practice Habits

### Before Practice

- Clear workspace
- Close distractions (email, chat)
- Set a time limit
- Choose a specific focus

### During Practice

- Focus on form, not just completion
- Notice inefficiencies in keystrokes
- Experiment with IDE shortcuts
- Stay in flow state

### After Practice

- Reflect on what felt awkward
- Note improvements to try next time
- Track progress over weeks

## TypeScript-Friendly Kata Ideas

| Kata | Focus Area |
|------|------------|
| String Calculator | TDD basics, parsing |
| Roman Numerals | Algorithm, edge cases |
| Tennis Scoring | State management |
| Bank Account | OOP, immutability |
| Gilded Rose | Refactoring legacy code |
| Mars Rover | Command pattern, state |
