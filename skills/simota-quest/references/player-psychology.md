# Player Psychology

Motivation frameworks, engagement models, behavioral psychology, and ethical engagement design for games.

---

## Bartle's Player Types

Richard Bartle's taxonomy (1996) classifies players by two axes: acting vs interacting, players vs world.

| Type | Motivation | % of Players | Design Implications |
|------|-----------|-------------|---------------------|
| **Achiever** | Goals, points, status | ~10% | Progression systems, leaderboards, completion tracking |
| **Explorer** | Discovery, knowledge | ~10% | Hidden content, lore, secrets, open worlds |
| **Socializer** | Relationships, community | ~75% | Chat, guilds, co-op, social spaces |
| **Killer** | Competition, dominance | ~5% | PvP, rankings, competitive modes |

### Extended Taxonomy (Bartle 2003)

Each type has implicit (experiencing) and explicit (acting) variants:
- **Achiever**: Planner (explicit) vs Opportunist (implicit)
- **Explorer**: Scientist (explicit) vs Hacker (implicit)
- **Socializer**: Networker (explicit) vs Friend (implicit)
- **Killer**: Politician (explicit) vs Griefer (implicit)

### Design Rules

- Don't design for one type exclusively; most players are multi-type
- Socializers are the majority — social features have the broadest impact
- Killers are few but vocal; their experience affects everyone
- Achievers generate content (guides, tier lists) that retains other types

---

## Flow Theory (Csikszentmihalyi)

Flow is the optimal state where challenge matches skill:

```
                    Anxiety
                   /
                  /
     Challenge   / ← FLOW CHANNEL
                /
               /
              /  Boredom
             /
            Skill →
```

### Flow Conditions

1. **Clear goals**: Player knows what to do next
2. **Immediate feedback**: Actions have visible results
3. **Challenge-skill balance**: Not too easy, not too hard
4. **Sense of control**: Player feels in command
5. **Loss of self-consciousness**: Immersion
6. **Time distortion**: "One more turn" effect

### Flow Channel Design

| Scenario | Player State | Design Response |
|----------|-------------|-----------------|
| Challenge >> Skill | Anxiety/Frustration | Lower difficulty, add hints, provide shortcuts |
| Challenge ≈ Skill | Flow | Maintain pace, gradual escalation |
| Challenge << Skill | Boredom | Increase challenge, add optional hard modes |

### Practical Application

- **Tutorial**: Start below flow channel (safe learning), quickly enter it
- **Mid-game**: Oscillate within flow channel (tension → release cycles)
- **Boss fights**: Brief anxiety (above channel) → mastery → flow
- **Post-game**: Optional challenges for players who've surpassed normal flow channel

---

## Self-Determination Theory (SDT)

Deci & Ryan's framework identifies three innate psychological needs:

| Need | Definition | Game Design |
|------|-----------|-------------|
| **Autonomy** | Feeling of choice and control | Multiple viable paths, customization, player expression |
| **Competence** | Feeling of mastery and growth | Clear progression, skill development, difficulty curve |
| **Relatedness** | Feeling of connection | Co-op, guilds, shared experiences, narrative bonds |

### SDT Application Guidelines

- **Autonomy**: Offer 2-3 meaningful choices per decision point. Forced linearity undermines autonomy.
- **Competence**: Every session should contain at least one "I got better" moment. Track and surface improvement.
- **Relatedness**: Enable players to share achievements, help each other, and build shared memories.

### Intrinsic vs Extrinsic Motivation

| Type | Driver | Game Example | Sustainability |
|------|--------|-------------|---------------|
| **Intrinsic** | Fun, curiosity, mastery | Exploring a beautiful world | High — self-sustaining |
| **Extrinsic** | Rewards, points, pressure | Daily login rewards | Low — requires escalation |

**Overjustification effect**: Adding extrinsic rewards to intrinsically fun activities can reduce intrinsic motivation. Don't reward exploration with XP if the exploration itself is the reward.

---

## Octalysis Framework (Yu-kai Chou)

8 core drives that motivate human behavior:

| # | Core Drive | Type | Game Example |
|---|-----------|------|-------------|
| 1 | **Epic Meaning & Calling** | White Hat | "You're the Chosen One" narrative |
| 2 | **Development & Accomplishment** | White Hat | XP bars, leveling, achievements |
| 3 | **Empowerment of Creativity** | White Hat | Minecraft building, character customization |
| 4 | **Ownership & Possession** | White Hat | Inventory, collections, housing |
| 5 | **Social Influence & Relatedness** | Neutral | Guilds, leaderboards, co-op |
| 6 | **Scarcity & Impatience** | Black Hat | Limited-time events, energy systems |
| 7 | **Unpredictability & Curiosity** | Black Hat | Loot boxes, random events, mystery |
| 8 | **Loss & Avoidance** | Black Hat | Streak timers, decay systems, FOMO |

### White Hat vs Black Hat

- **White Hat** (drives 1-4): Make players feel powerful, creative, meaningful. Sustainable engagement.
- **Black Hat** (drives 6-8): Create urgency, anxiety, obsession. Effective but potentially harmful.
- **Best practice**: Lead with White Hat; use Black Hat sparingly and ethically.

### Octalysis Scoring

Rate each drive 0-10 for your game design. A balanced game has:
- Total score > 40 (out of 80)
- No single drive < 2 (gaps feel incomplete)
- No more than 2 Black Hat drives > 7 (risk of exploitation)

---

## Hook Model (Nir Eyal)

Four-phase habit-forming loop:

```
TRIGGER → ACTION → VARIABLE REWARD → INVESTMENT
   ↑                                       |
   └───────────────────────────────────────┘
```

### Phase Details

| Phase | Definition | Game Example |
|-------|-----------|-------------|
| **Trigger** | External (notification) or internal (boredom, curiosity) | Push notification, "one more turn" urge |
| **Action** | Simplest behavior in anticipation of reward | Tap to open, start a match |
| **Variable Reward** | Unpredictable positive outcome | Random loot drop, match result uncertainty |
| **Investment** | User puts something in (time, data, effort, money) | Character progress, friend list, customization |

### Ethical Hook Design

- Triggers should invite, not interrupt (no notification spam)
- Actions should be genuinely fun, not just habitual
- Rewards should feel fair and transparent
- Investment should create genuine value for the player

---

## Prospect Theory & Loss Aversion

Daniel Kahneman & Amos Tversky's finding: Losses feel approximately 2× as painful as equivalent gains feel good.

```
perceived_value(loss) ≈ -2 × perceived_value(equivalent_gain)
```

### Game Design Implications

| Concept | Application | Example |
|---------|------------|---------|
| **Loss aversion** | Players hate losing progress more than they enjoy gaining it | Souls-like death mechanics risk player inventory |
| **Endowment effect** | Players overvalue what they already own | Trading systems, item attachment |
| **Sunk cost** | Players continue due to prior investment | Long-running campaigns, gacha collections |
| **Framing effect** | "Save 20%" feels different than "Lose 80%" | Positive framing for difficulty; negative framing sparingly |

### Ethical Guidelines

- Don't exploit loss aversion for monetization (e.g., "Buy now or lose your streak!")
- Offer "save points" and "undo" to reduce unfair loss
- If using roguelike permadeath, make losses feel fair and educational
- Frame failure as learning, not punishment

---

## Engagement Loops

### Four Loop Levels

| Level | Duration | Example | Design Focus |
|-------|----------|---------|--------------|
| **Micro loop** | Seconds | Shoot → hit → feedback | Game feel, juice, responsiveness |
| **Session loop** | Minutes-hours | Quest → reward → upgrade → harder quest | Core loop, progression feel |
| **Progression loop** | Days-weeks | Level up → new area → new abilities | Content pacing, unlock cadence |
| **Meta loop** | Weeks-months | Season → rank → reset → new season | Long-term retention, social status |

### Loop Design Rules

1. **Every loop must complete**: No dead ends or infinite grinds without payoff
2. **Loops nest**: Micro loops compose into session loops, etc.
3. **Variable reward at each level**: Something unpredictable at every timescale
4. **Rest points**: Allow natural stopping points (save between sessions, not mid-boss)
5. **Escalation**: Each loop iteration should feel slightly different (new enemies, mechanics, stakes)

### Session Length Targets

| Platform | Target Session | Design Strategy |
|----------|---------------|-----------------|
| Mobile casual | 3-5 minutes | Bite-sized levels, auto-save |
| Mobile mid-core | 10-20 minutes | Match-based, energy systems |
| PC/Console casual | 30-60 minutes | Chapter structure, frequent saves |
| PC/Console core | 1-3 hours | Quest chains, dungeon runs |
| MMO/Live service | Variable | Multiple activity types for any length |

---

## The Ethical Line

### Engagement vs Compulsion

| Engagement (Ethical) | Compulsion (Unethical) |
|---------------------|----------------------|
| Player plays because it's fun | Player plays to avoid punishment |
| Natural stopping points exist | Stopping feels like losing |
| Rewards are transparent | Rewards exploit cognitive biases |
| Player feels in control | Player feels manipulated |
| Time spent is valued | Time spent is extracted |
| Absence has no penalty | Absence causes loss (streak, decay) |

### Questions to Ask

For every engagement mechanic, answer:
1. Would I be comfortable if the player knew exactly how this mechanic works?
2. Does this mechanic respect the player's time and autonomy?
3. Would I be comfortable if my child played this game for 4 hours?
4. Does the player benefit from engagement, or just the developer?

If any answer is "no" → flag as ethical concern in Quest deliverable.

---

## Quick Reference: Framework Selection Guide

| Design Question | Recommended Framework |
|----------------|----------------------|
| Who are my players? | Bartle types + Cast personas |
| Is my game fun moment-to-moment? | Flow theory |
| Why do players keep playing? | SDT (intrinsic) + Octalysis |
| How do I build habits? | Hook Model |
| How do I frame rewards/losses? | Prospect theory |
| How do I structure play sessions? | Engagement loops |
| Is my design ethical? | Ethical line checklist |
