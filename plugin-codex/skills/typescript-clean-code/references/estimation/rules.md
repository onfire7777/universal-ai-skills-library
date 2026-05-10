# Estimation Rules

Rules for providing professional estimates and avoiding commitment traps.

## Core Rules

### 1. Never Confuse Estimates with Commitments

Estimates are guesses; commitments are promises. Keep them separate.

- Only commit when you are certain you can deliver
- Communicate estimates as probability distributions, not single numbers
- Make the distinction explicit in conversations with stakeholders

### 2. Never Make Implied Commitments

Watch for language that creates hidden commitments.

- Never agree to "try" - it's a commitment in disguise
- If you say yes to "try," you're committing to work extra hours, weekends, etc.
- If pressed for commitment, decline clearly rather than hedge

**Example**:
```
// Bad - creates implied commitment
Mike: "Can you try to make it no more than six days?"
Peter: "I'll try."

// Good - maintains honesty
Mike: "Can you try to make it no more than six days?"
Peter: "I can't commit to that. It could be five or six days, 
        possibly longer if things go wrong."
```

### 3. Provide Trivariate Estimates (PERT)

Give three numbers, not one.

- **O (Optimistic)**: Everything goes perfectly (< 1% chance)
- **N (Nominal)**: Most likely duration (highest probability)
- **P (Pessimistic)**: Everything goes wrong (< 1% chance)

Calculate expected duration: **(O + 4N + P) / 6**
Calculate uncertainty: **(P - O) / 6**

### 4. Use Team Estimation Techniques

Don't estimate alone - use the people around you.

- **Wideband Delphi**: Team discusses and estimates until consensus
- **Flying Fingers**: Show 0-5 fingers simultaneously after discussion
- **Planning Poker**: Use cards (0, 1, 3, 5, 10) to estimate simultaneously
- **Affinity Estimation**: Silently sort task cards by size, then discuss

### 5. Break Large Tasks into Smaller Ones

The Law of Large Numbers improves accuracy.

- Estimate small tasks independently
- Errors tend to cancel out (though underestimation bias persists)
- Breaking tasks up reveals hidden complexity

## Guidelines

### Communicating Estimates

- Always include the probability distribution, not just the peak
- "Probably 3 days, could be 5-6, unlikely but possibly 10-11"
- Let managers understand the uncertainty to make appropriate plans

### When Asked for Commitment

- Ask yourself: Am I *certain* I can achieve this?
- If not certain, decline to commit
- Offer an estimate with probability distribution instead

### Combining Multiple Task Estimates

For a sequence of tasks:
- Expected total = sum of individual expected durations
- Combined sigma = sqrt(sum of squared individual sigmas)

## Exceptions

- **Crisis situations**: May need to provide quick rough estimates, but clarify they are rough
- **Trivial tasks**: Very small, well-understood tasks may not need trivariate analysis

## Quick Reference

| Rule | Summary |
|------|---------|
| Estimate != Commitment | Keep guesses and promises separate |
| Never "try" | Agreeing to try is agreeing to commit |
| Use PERT | Three numbers: O, N, P -> calculate expected |
| Team estimation | Use consensus techniques, don't estimate alone |
| Break it down | Smaller tasks = better accuracy |
| Communicate uncertainty | Share the distribution, not just the peak |

## Estimation Techniques Summary

| Technique | How It Works | Best For |
|-----------|--------------|----------|
| Flying Fingers | Hold up 0-5 fingers simultaneously | Quick team estimates |
| Planning Poker | Cards (0,1,3,5,10) shown simultaneously | Distributed teams |
| Affinity Estimation | Silent sorting, then discussion | Many tasks at once |
| Trivariate (PERT) | O + 4N + P / 6 | Individual task probability |
