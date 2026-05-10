# Estimation Examples

Examples demonstrating estimation principles and calculations.

## Bad Examples

### Vague Single-Number Estimate

```
Mike: "What is your estimate for completing the Frazzle task?"
Peter: "Three days."
```

**Problems**:
- No information about likelihood
- Mike doesn't know if 3 days is certain or optimistic
- No way to plan for delays
- Creates implied expectation of completion in 3 days

### The "Try" Trap

```
Mike: "Peter, can you give me a hard date when you'll be done?"
Peter: "No, Mike. Like I said, it'll probably be done in three, 
        maybe four, days."
Mike: "Can we say four then?"
Peter: "No, it could be five or six."
Mike: "OK, Peter, but can you try to make it no more than six days?"
Peter: "Sure, I'll try."
```

**Problems**:
- Peter just made a commitment disguised as "trying"
- If not done in 6 days, Mike can accuse him of "not trying hard enough"
- Peter is now implicitly agreeing to work overtime, weekends
- The word "try" is loaded - agreeing to try is agreeing to succeed

## Good Examples

### Communicating Uncertainty

```
Mike: "What is your estimate for completing the Frazzle task?"
Peter: "I'd estimate three days as most likely."
Mike: "How likely is it that you'll be done in three days?"
Peter: "Fifty or sixty percent."
Mike: "So there's a good chance that it'll take you four days."
Peter: "Yes, in fact it might even take me five or six, though I doubt it."
Mike: "How much do you doubt it?"
Peter: "I'm ninety-five percent certain I'll be done before six days. 
        If everything goes wrong, it could take ten or eleven days. 
        But it's not very likely that so much will go wrong."
```

**Why it works**:
- Peter shares the full probability distribution
- Mike understands the range of possibilities
- No false precision or hidden commitments
- Mike can plan appropriately for the uncertainty

### Declining a Commitment Request

```
Mike: "Peter, can you give me a hard date when you'll be done?"
Peter: "No, Mike. Like I said, it'll probably be done in three, 
        maybe four, days."
Mike: "Can we say four then?"
Peter: "No, it could be five or six."
Mike: "Can you try to make it no more than six days?"
Peter: "I can't commit to that, Mike. I've given you my best estimate - 
        probably 3-4 days, could be 5-6, outside chance of longer."
```

**Why it works**:
- Peter avoids the "try" trap
- Clear distinction between estimate and commitment
- Repeats the probability information
- Professional and honest

## PERT Calculation Examples

### Single Task Estimate

**Scenario**: Peter estimates the Frazzle task

| Estimate Type | Value | Meaning |
|--------------|-------|---------|
| O (Optimistic) | 1 day | If everything goes perfectly |
| N (Nominal) | 3 days | Most likely duration |
| P (Pessimistic) | 12 days | If everything goes wrong |

**Calculations**:
- Expected duration (mu) = (1 + 4(3) + 12) / 6 = (1 + 12 + 12) / 6 = **4.2 days**
- Standard deviation (sigma) = (12 - 1) / 6 = **1.8 days**

**Interpretation**: Task will likely take ~4 days, could reasonably take 6 days (mu + sigma), possibly even 8 days (mu + 2*sigma).

### Multiple Tasks Combined

**Scenario**: Peter has three sequential tasks

| Task | O | N | P | mu | sigma |
|------|---|---|---|-----|-------|
| Alpha | 1 | 3 | 12 | 4.2 | 1.8 |
| Beta | 1 | 1 | 14 | 3.5 | 2.2 |
| Gamma | 3 | 6 | 8 | 6.5 | 1.3 |

**Combined Calculations**:
- Total expected = 4.2 + 3.5 + 6.5 = **14 days**
- Combined sigma = sqrt(1.8^2 + 2.2^2 + 1.3^2)
  - = sqrt(3.24 + 4.84 + 1.69)
  - = sqrt(9.77)
  - = **~3 days**

**Interpretation**:
- Likely completion: 14 days
- Could take 17 days (1 sigma)
- Possibly 20 days (2 sigma)

**The Surprise**: 
- Optimistic totals: 1 + 1 + 3 = 5 days
- Nominal totals: 3 + 1 + 6 = 10 days
- Yet realistic expectation: 14+ days

This is why projects estimated optimistically take 3-5x longer than hoped!

## Estimation Technique Examples

### Flying Fingers

```
Moderator: "Let's estimate the user authentication task. 
            Scale is days. Discussion first..."

[Team discusses requirements, complications, implementation approaches]

Moderator: "Hands under the table... 1, 2, 3, show!"

Alice: 3 fingers
Bob: 4 fingers  
Carol: 3 fingers
Dave: 1 finger  <-- outlier

Moderator: "Dave, why 1 day?"
Dave: "I thought we could reuse the OAuth library from project X."
Alice: "Oh, I didn't know about that. That changes things."

[Further discussion, then re-vote]

All: 2 fingers
Moderator: "Consensus at 2 days."
```

### Planning Poker with Trivariate

```
Moderator: "Task: Implement payment processing. 
            First, show cards for OPTIMISTIC estimate."

[Team shows cards]
Results: 2, 1, 2, 1, 2
Take lowest: O = 1 day

Moderator: "Now NOMINAL estimate."
Results: 5, 5, 3, 5, 5
Take consensus: N = 5 days

Moderator: "Now PESSIMISTIC estimate."
Results: 10, 14, 12, 10, 8
Take highest: P = 14 days

Expected = (1 + 20 + 14) / 6 = 5.8 days
Sigma = (14 - 1) / 6 = 2.2 days
```

## Key Takeaways

| Situation | Bad Response | Good Response |
|-----------|--------------|---------------|
| "How long?" | "3 days" | "Probably 3, could be 5-6, possibly longer" |
| "Can you commit?" | "I'll try" | "I can't commit, but here's my estimate..." |
| "Just give me a number" | Single number | O/N/P with calculation |
| Estimating alone | Trust your gut | Use team consensus techniques |
