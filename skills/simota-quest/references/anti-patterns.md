# Anti-Patterns

Common game design pitfalls, anti-fun patterns, dark patterns, and balance traps. Every Quest deliverable must be checked against this reference.

---

## Zileas / Riot Games Anti-Fun Patterns

Tom "Zileas" Cadwell (former Riot VP of Design) identified patterns that feel unfun even when numerically balanced:

| # | Pattern | Description | Prevention |
|---|---------|-------------|------------|
| 1 | **Burden of Knowledge** | Mechanic requires external research to understand | Teach through gameplay; clear tooltips; progressive complexity |
| 2 | **Invisible Power** | Stat bonuses player can't perceive | Make power gains visible and felt; use clear feedback |
| 3 | **Unintuitive Counters** | Counter-play that doesn't make logical sense | Counters should be discoverable through logic, not memorization |
| 4 | **Fun Denial** | Mechanic that removes another player's fun (long CC, invulnerability) | Cap CC duration; provide counterplay; reduce frustration time |
| 5 | **No Counterplay** | Mechanic with no possible response | Every ability must have at least one counterplay option |
| 6 | **False Choice** | Options where one is always optimal | Ensure each option has a context where it's best |
| 7 | **Anti-Comeback** | Winning makes winning easier exponentially | Add rubber-banding; scale rewards inversely with lead |
| 8 | **Grind Without Mastery** | Time investment without skill growth | Tie progression to skill milestones, not just hours |
| 9 | **Toxic Incentive** | Mechanic that encourages anti-social behavior | Align individual incentives with team/community goals |
| 10 | **Unclear Optimization** | Player can't tell what's making them succeed/fail | Clear feedback on what works; readable cause-effect |
| 11 | **Power Without Decisions** | Automatic power gains without meaningful choice | Require decisions at growth points; tradeoffs, not just upgrades |

---

## MMORPG Anti-Patterns

| Pattern | Description | Prevention |
|---------|-------------|------------|
| **Gold Inflation** | Currency generation outpaces sinks | See `economy-design.md` inflation control strategies |
| **Power Creep** | Each content update invalidates previous gear | Horizontal progression; sidegrades; stat squishes |
| **Content Drought** | Gaps between updates cause player exodus | Content cadence planning; evergreen systems; UGC |
| **Dailyitis** | Players feel obligated, not motivated | Remove FOMO; allow catch-up; respect player time |
| **Alt-Unfriendly** | Progress doesn't transfer between characters | Account-wide unlocks; roster systems |
| **Gear Treadmill** | Constant item level increases reset progress | Meaningful endgame; horizontal progression; transmog |
| **Dead Content** | Old zones/raids become irrelevant | Level sync; transmog incentives; achievement rewards |
| **Toxicity Loop** | Competitive systems breed harassment | Positive reinforcement; honor systems; smart moderation |
| **Whale Dependency** | Revenue relies on <1% of players spending excessively | Diversified monetization; spending caps; ethical limits |
| **Vertical-Only Endgame** | Only path is "get bigger numbers" | Collections, cosmetics, social features, housing, PvP modes |

---

## Mobile Game Dark Patterns

| Pattern | Description | Ethical Alternative |
|---------|-------------|---------------------|
| **Artificial Scarcity** | Limited-time offers creating urgency | Rotating availability with advance notice |
| **Social Pressure** | Notifications shaming non-play or exposing to friends | Opt-in social features; no shame mechanics |
| **Skinner Box** | Pure variable-ratio reward with no mastery | Add skill component; reduce randomness dependency |
| **Energy Systems** | Time-gates that sell removal | Generous energy; no purchase required for daily play |
| **Notification Spam** | Constant push notifications to re-engage | Respectful notification frequency; easy opt-out |
| **Sunk Cost Exploitation** | Making players feel they'll "waste" prior investment | Allow graceful exit; respect player autonomy |
| **Pay-to-Skip** | Selling solutions to designed frustration | Don't design frustration; sell cosmetics instead |
| **Loot Box Obscurity** | Hiding odds or manipulating perceived value | Transparent odds; pity systems; clear pricing |
| **Anchoring** | Showing inflated "original price" to make deals seem valuable | Honest pricing; no fake discounts |
| **Incomplete Collections** | Giving 99/100 pieces to trigger completion drive | Allow collections to complete through gameplay |

---

## Balance Anti-Patterns

| Pattern | Description | Prevention |
|---------|-------------|------------|
| **Number Inflation** | Stats grow so large they lose meaning | Use percentage scaling; periodic stat squish; diminishing returns |
| **Meta Lock** | One dominant strategy with no viable alternatives | Regular balance patches; multiple viable builds; counter systems |
| **Pay-to-Win** | Purchased advantages in competitive play | Cosmetics-only monetization; separate competitive queues |
| **RNG Dependency** | Outcomes determined entirely by randomness | Add skill modifiers; bad luck protection; player agency in outcomes |
| **Nerf Cycle** | Repeatedly nerfing the strongest option | Buff the weak instead; smaller adjustments; consider the ecosystem |
| **Complexity Creep** | Each update adds mechanics without removing | Complexity budget; remove or consolidate before adding |
| **Gear > Skill** | Equipment matters more than player ability | Skill-based multipliers; gear as enabler, not determinant |
| **Binary Balance** | Character/item is either overpowered or useless | Granular tuning; role differentiation; niche strengths |
| **One-Dimensional Scaling** | Single stat determines everything | Multi-stat interactions; role-specific scaling |
| **Snowball Effect** | Early advantage guarantees victory | Comeback mechanics; diminishing returns on leads |

---

## UX Anti-Patterns

| Pattern | Description | Prevention |
|---------|-------------|------------|
| **Forced Tutorial** | Unskippable, lengthy onboarding | Progressive disclosure; learn-by-doing; skip option for veterans |
| **Information Overload** | All mechanics explained at once | Introduce one system at a time; in-context teaching |
| **Menu Labyrinth** | Important features buried in nested menus | Max 3 levels deep; most-used features on main screen |
| **Notification Overload** | Constant alerts competing for attention | Prioritize notifications; batch low-priority items |
| **False Urgency** | Countdown timers on non-time-sensitive content | Reserve timers for genuinely time-limited events |
| **Confirmation Fatigue** | "Are you sure?" on every action | Confirm only destructive/irreversible actions; undo instead |
| **Hidden Depth** | Core features discoverable only by accident | Feature discovery through gameplay; clear UI signaling |

---

## Anti-Pattern Audit Checklist

For every Quest deliverable, verify:

### Balance
- [ ] No single dominant strategy without counterplay
- [ ] Numbers justified with formulas, not intuition
- [ ] Edge cases tested (min/max level, best/worst gear)
- [ ] Progression has both vertical and horizontal elements
- [ ] Comeback mechanics exist for competitive modes

### Economy
- [ ] Taps and sinks identified and balanced
- [ ] Inflation projections modeled
- [ ] New player purchasing power preserved
- [ ] No P2W without explicit request + ethical flag

### Narrative
- [ ] No ludonarrative dissonance
- [ ] Player choices have meaningful consequences
- [ ] No false choices (options leading to same outcome)
- [ ] Exposition distributed, not dumped

### Psychology
- [ ] Engagement loops based on mastery, not exploitation
- [ ] No predatory FOMO mechanics
- [ ] Respect for player time (no artificial time-gates as monetization)
- [ ] Transparent odds for all random elements

### UX
- [ ] Progressive disclosure of complexity
- [ ] No forced unskippable tutorials
- [ ] Important features not hidden
- [ ] Notification frequency respectful

### Ethics
- [ ] Dark patterns flagged and disclosed
- [ ] Child protection measures if applicable
- [ ] Spending limits and transparency for monetization
- [ ] Variable-ratio rewards used responsibly with alternatives
