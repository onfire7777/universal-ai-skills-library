# Economy Design

Currency systems, monetization models, inflation control, and ethical guidelines for game economies.

---

## Taps & Sinks Model

Every game economy is a flow system:

- **Taps** (sources): Where currency/resources enter the economy
- **Sinks** (drains): Where currency/resources leave the economy
- **Equilibrium**: Taps ≈ Sinks over time → stable economy

### Common Taps

| Tap | Type | Notes |
|-----|------|-------|
| Quest rewards | Finite | One-time; depletes over content lifecycle |
| Monster drops | Infinite | Must be rate-limited |
| Daily login | Infinite | Predictable income floor |
| Crafting/selling | Infinite | Player-driven; harder to control |
| Real-money purchase | Infinite | External injection; inflationary risk |
| Achievement rewards | Finite | Milestone-based |

### Common Sinks

| Sink | Type | Effectiveness |
|------|------|---------------|
| Consumables | Recurring | High — constant drain |
| Repair costs | Recurring | Medium — tied to activity |
| Fast travel fees | Recurring | Low — convenience tax |
| Auction house tax | Recurring | High — removes currency per transaction |
| Cosmetic purchases | One-time | Medium — voluntary but desirable |
| Upgrade costs (exponential) | Recurring | High — exponential scaling absorbs inflation |
| Crafting materials | Recurring | Medium — resource transformation |
| Guild/clan fees | Recurring | Low-medium — social pressure to pay |
| Prestige reset | One-time | High — voluntarily resetting progress for status |

---

## Castronova's Three Pillars

Edward Castronova's framework for virtual economies:

1. **Production**: How goods/currency are created (taps, crafting, drops)
2. **Distribution**: How goods move between players (trade, auction, gifts)
3. **Consumption**: How goods are destroyed/used (sinks, durability, consumables)

A healthy economy needs all three pillars balanced. Most economy failures come from weak consumption (insufficient sinks).

---

## Inflation Control Strategies

| # | Strategy | Mechanism | Example |
|---|----------|-----------|---------|
| 1 | **Transaction tax** | Remove % of currency per trade | 5% auction house fee (WoW) |
| 2 | **Repair costs** | Currency drain tied to activity | Equipment degradation (Dark Souls) |
| 3 | **Exponential upgrade costs** | Absorb accumulated wealth | Upgrade costs doubling per tier |
| 4 | **Limited inventory** | Cap hoarding | Bank slot limits, weight systems |
| 5 | **Seasonal resets** | Periodic economy wipe | Ladder seasons (Path of Exile) |
| 6 | **Bind-on-pickup** | Remove items from trade | Soulbound gear (WoW) |
| 7 | **Consumable meta** | Ensure constant resource drain | Potion/buff requirements for endgame |
| 8 | **Currency caps** | Hard limit on holdings | Maximum gold per character |
| 9 | **NPC price scaling** | Prices rise with server wealth | Dynamic vendor pricing |
| 10 | **Prestige systems** | Convert wealth to status | Spending currency for titles/cosmetics |

### Inflation Detection Metrics

- **CPI equivalent**: Track price of a basket of common items over time
- **Velocity of money**: How quickly currency changes hands
- **Gini coefficient**: Wealth distribution inequality (>0.6 = problematic)
- **New player purchasing power**: Can a new player afford basic gear?

---

## Dual Currency Design

Most F2P games use two currencies:

| Currency | Earned By | Spent On | Design Purpose |
|----------|-----------|----------|----------------|
| **Soft** (gold, coins) | Gameplay | Progression, consumables | Core loop reward |
| **Hard** (gems, crystals) | Real money (+ small gameplay) | Cosmetics, convenience, gacha | Monetization |

### Design Rules

- Soft currency should always be sufficient for core progression
- Hard currency should never be required for completing content
- Exchange rate (hard → soft) should be one-directional to prevent arbitrage
- Small hard currency trickle from gameplay (5-10% of what purchasers get) maintains F2P player engagement

---

## Gacha & Pity Systems

### Gacha Models

| Model | Description | Player Feel |
|-------|-------------|-------------|
| **Standard** | Fixed probability per pull | Unpredictable, can feel unfair |
| **Soft pity** | Rate increases after threshold | Building excitement |
| **Hard pity** | Guaranteed at N pulls | Safety net, budget planning |
| **Spark** | Choose any item after N pulls | Full control after investment |
| **Step-up** | Increasing value per sequential pull | Encourages multi-pull |

### Pity Math

```
P(getting item by pull N) = 1 - (1 - base_rate)^N        # without pity
P(getting item by pull N) = 1 - Π(1 - rate_i) for i=1..N  # with escalating pity
```

**Hard pity guarantee**: At pull N, rate = 100%.
**Soft pity example**: Base 0.6%, +6% per pull after pull 73, hard pity at 90. Expected pulls ≈ 62.

### Ethical Gacha Guidelines

- [ ] Display exact probabilities for all outcomes
- [ ] Implement pity/spark system (hard pity ≤ 200 pulls recommended)
- [ ] Show pull history and pity counter
- [ ] No "bait and switch" — featured items must have stated rates
- [ ] Allow earning gacha currency through gameplay (not just purchase)
- [ ] Age-gate and spending limits for minors
- [ ] No limited-time-only items that create FOMO pressure

---

## Monetization Models

### Comparison

| Model | Revenue Source | Player Trust | Design Constraint |
|-------|---------------|-------------|-------------------|
| **Premium (B2P)** | Upfront purchase | High | Content must justify price |
| **Subscription** | Monthly fee | High | Constant content updates needed |
| **F2P + Cosmetics** | Optional cosmetics | Medium-high | Core game must be fully free |
| **F2P + Convenience** | Time-savers, boosts | Medium | Must not feel required |
| **F2P + Gacha** | Random reward pulls | Low-medium | Requires ethical safeguards |
| **F2P + P2W** | Power advantages | Low | Destroys competitive integrity |
| **Hybrid** | B2P + cosmetic MTX | Medium-high | No pay-for-power post-purchase |

### Ethical Monetization Principles

1. **No pay-to-win**: Purchased items must not provide competitive advantage
2. **Transparency**: All odds, rates, and contents clearly disclosed
3. **Respect time**: Paying accelerates, never gates core content
4. **No predatory FOMO**: Limited-time offers must not exploit urgency
5. **Child protection**: Age verification, spending caps, parental controls
6. **Earned, not rented**: Purchased items persist; no expiring purchases
7. **Value clarity**: Players understand exactly what they're buying

---

## Real-World Economy Failures

### Diablo II (2000s)
- **Problem**: Rampant duplication exploits flooded economy with gold and items
- **Result**: Gold became worthless; Stone of Jordan (SoJ) became de facto currency
- **Lesson**: Economy security is as important as economy design

### Diablo III Auction House (2012–2014)
- **Problem**: Real Money Auction House (RMAH) made playing for loot less satisfying than buying it
- **Result**: Core gameplay loop undermined; RMAH removed in 2014
- **Lesson**: Monetization must not shortcut the core game loop

### Venezuela's RuneScape Gold Farming (2017+)
- **Problem**: Real-world economic crisis drove mass gold farming
- **Result**: Hyperinflation in-game; legitimate players priced out
- **Lesson**: Virtual economies are not isolated from real-world pressures

### Star Wars Battlefront II (2017)
- **Problem**: Loot boxes provided gameplay advantages (P2W)
- **Result**: Massive backlash, government investigations, legislation changes
- **Lesson**: P2W + loot boxes = regulatory and reputational risk

---

## Economy Design Checklist

- [ ] All taps and sinks identified and documented
- [ ] Tap:sink ratio modeled for 30/90/365 day horizons
- [ ] New player purchasing power calculated
- [ ] Inflation detection metrics defined
- [ ] Currency cap or overflow handling designed
- [ ] Dual currency exchange rules specified (if applicable)
- [ ] Monetization model selected with ethical review
- [ ] Gacha probabilities disclosed and pity system designed (if applicable)
- [ ] Economy simulation spec prepared for Monte Carlo validation
