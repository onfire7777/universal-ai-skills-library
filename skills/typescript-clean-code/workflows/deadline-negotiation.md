# Deadline Negotiation Workflow

Handling unrealistic deadline requests professionally.

## When to Use

- Asked to commit to an impossible deadline
- Deadline is set without your input
- Scope exceeds available time
- Pressure to "just make it work"

## Prerequisites

- Understanding of the work required
- Estimate of realistic timeline (see `estimation.md`)
- Knowledge of your capacity

**Reference**: `saying-no/rules.md`, `commitment/rules.md`, `estimation/rules.md`, `pressure/rules.md`

---

## The Golden Rules

1. **Never say "I'll try"** - It's dishonest and sets you up to fail
2. **Provide data, not emotions** - Use estimates and facts
3. **Offer alternatives** - Don't just say no, propose solutions
4. **Escalate professionally** - Not as a threat, but as information

---

## Workflow Steps

### Step 1: Gather Information

**Goal**: Understand what's being asked before responding.

- [ ] What exactly is the deadline?
- [ ] What is the full scope of work?
- [ ] Why is this deadline set?
- [ ] What are the consequences of missing it?
- [ ] Who set the deadline?

**Ask clarifying questions**:
```
"Help me understand - what's driving the March 1st date?"
"What's included in 'complete'? Does that include testing and documentation?"
"What happens if we deliver on March 15th instead?"
```

---

### Step 2: Do Your Estimation

**Goal**: Know the realistic timeline before negotiating.

**Reference**: `estimation.md` workflow

- [ ] Break down the work into tasks
- [ ] Apply PERT estimation
- [ ] Calculate expected time and uncertainty
- [ ] Document assumptions

**Example**:
```
Feature: User Export

Estimated: 25-32 hours (95% confidence)
Calendar time: 4-5 days (accounting for meetings, etc.)
Realistic delivery: 1.5 weeks from start

Requested deadline: End of this week (3 days)
Gap: 2-4 days short
```

---

### Step 3: Identify Options

**Goal**: Prepare alternatives before the conversation.

**Option types**:

1. **Reduce scope**: What can be cut or deferred?
2. **Extend deadline**: What's the realistic date?
3. **Add resources**: Can anyone help? (Be realistic about ramp-up)
4. **Change approach**: Is there a simpler solution?

**Example options**:
```
Option A: Reduce Scope
- Deliver JSON export only (cut CSV)
- Skip documentation
- Delivery: End of week (meets deadline)
- Trade-off: CSV comes in week 2

Option B: Extend Deadline
- Full feature with JSON + CSV
- Complete documentation
- Delivery: End of next week
- Trade-off: 1 week later

Option C: Partial Delivery
- JSON export by end of week
- CSV export by mid next week
- Trade-off: Two releases needed
```

---

### Step 4: Have the Conversation

**Goal**: Communicate clearly and professionally.

**Structure**:
1. Acknowledge the request
2. Share your analysis
3. Present the gap
4. Offer alternatives
5. Ask for decision

**Example conversation**:

**Bad approach** (don't do this):
```
"There's no way I can do this by Friday. It's impossible."
```

**Good approach**:
```
"I understand the deadline is Friday. I've broken down the work,
and with testing and documentation, I estimate 4-5 days of work.

Given we're starting today (Tuesday), that puts realistic 
completion at next Tuesday, not Friday.

I see three options:
1. We can deliver JSON export only by Friday, and add CSV next week
2. We can extend the deadline to next Tuesday for the full feature
3. We can discuss if there's scope I'm not understanding correctly

Which option works best for the business?"
```

**Key phrases**:
- "Based on my analysis..."
- "Here are the options I see..."
- "What would you like to prioritize?"
- "I want to give you a commitment I can keep"

---

### Step 5: Handle Pushback

**Goal**: Stay professional when pressured.

**Common pushback and responses**:

**"Can't you just try?"**
```
Reference: commitment/rules.md - The "Try" problem

Response: "I want to be honest with you. If I say I'll try, I'm really 
saying I don't think I can do it but I don't want to say no. 
That's not fair to either of us. Here's what I CAN commit to..."
```

**"We promised the client"**
```
Reference: saying-no/rules.md - High stakes

Response: "I understand there's a commitment to the client. That makes 
it even more important that we're realistic. Would you rather tell 
them now about a 3-day delay, or on Friday that we missed the deadline?
Let's figure out the best path forward together."
```

**"Everyone else manages to deliver on time"**
```
Response: "I want to deliver on time too. Can we look at the scope 
together? Maybe I'm including work that isn't needed, or maybe 
there's something I can learn from how others approach this."
```

**"This is critical for the business"**
```
Response: "I understand this is important. That's exactly why I want 
to give you accurate information. Committing to Friday and delivering 
broken code or missing on Tuesday would hurt the business more than 
being clear now about realistic timing."
```

---

### Step 6: Document the Agreement

**Goal**: Have a clear record of what was agreed.

After the conversation:
- [ ] Send a follow-up email/message summarizing:
  - What was agreed
  - What's in scope / out of scope
  - The delivery date
  - Any assumptions or dependencies

**Example follow-up**:
```
Hi [Manager],

Thanks for the discussion. To confirm our agreement:

Scope: User export feature - JSON format only
Deadline: Friday, March 1st
Out of scope: CSV format (scheduled for March 8th)

Assumptions:
- No other high-priority work assigned this week
- API contract is finalized
- Test environment is available

Let me know if I've missed anything.
```

---

### Step 7: Handle Escalation (if needed)

**Goal**: Escalate professionally if pushback continues.

**When to escalate**:
- You're being asked to commit to something impossible
- Your concerns are being dismissed
- The risk to the project/company is significant

**How to escalate**:
```
"I want to make sure we're making the right decision here. I'm 
not comfortable committing to Friday because [specific reasons].

I think we should involve [manager/stakeholder] to make sure we're 
aligned on the trade-offs. Would you like me to set up that meeting?"
```

**Reference**: `saying-no/rules.md` - "Threaten professionally"

---

## Quick Reference

### What NOT to Say

| Don't Say | Why | Say Instead |
|-----------|-----|-------------|
| "I'll try" | Dishonest, sets up failure | "I can commit to X by Y" |
| "That's impossible" | Sounds defensive | "Based on my analysis, that would require..." |
| "You don't understand" | Confrontational | "Help me understand the constraints..." |
| "Fine, I'll do it" | Resentful commitment | "Here are the options..." |

### What TO Say

| Situation | Response |
|-----------|----------|
| Unrealistic deadline | "Based on my estimate, realistic delivery is [date]. Here are options..." |
| Pressure to commit | "I want to give you a commitment I can keep. Can we discuss scope?" |
| Scope creep | "If we add X, the deadline moves to Y. Which do you prefer?" |
| Being dismissed | "I understand the pressure. Let's involve [stakeholder] to decide." |

---

## Exit Criteria

Negotiation is complete when:
- [ ] Clear agreement on scope and deadline
- [ ] Agreement documented in writing
- [ ] You can honestly commit to the deadline
- [ ] Trade-offs are understood by all parties
- [ ] Escalation path is clear if things change
