# Collaboration Examples

Scenarios demonstrating good and bad collaboration patterns in software teams.

## Bad Scenarios

### The Printer Company Fiasco

**Situation**: A company built high-end printers with many components (feeders, printers, stackers, staplers, cutters).

**What happened**:
- Each programmer worked on "their" device
- One person owned the feeder code, another owned the stapler code
- Each kept their technology to themselves
- No one could touch anyone else's code
- Political clout was tied to how much the business valued each device
- The printer programmer was "unassailable"
- Salary reviews were tied to device importance

**Problems observed**:
- Massive duplication across the codebase
- Completely skewed interfaces between modules
- No amount of consulting could convince them to change
- Business incentives reinforced bad technical practices

**Root cause**: Owned code culture backed by misaligned incentives.

---

### The Cubicle Corner Anti-Pattern

**Situation**: A team of programmers working in a modern office.

**What it looks like**:
- Programmers sitting in cubicle corners
- Backs turned to each other
- Staring at screens
- Wearing headphones all day
- No spontaneous conversation

**Problems**:
- No serendipitous communication
- Can't sense when teammates are struggling
- No knowledge sharing happening
- Not actually functioning as a team
- Might as well be remote contractors

**The metaphor**: "Rubbing cerebellums" - the cerebellum is at the back of the brain, so to rub them you'd face away from each other. That's not collaboration.

---

### The "I Work Better Alone" Fallacy

**Situation**: A developer believes they are more productive working solo.

**The problem**:
- Even if individually true, the *team* doesn't work better
- Creates knowledge silos
- No code review happening
- When that person is sick or leaves, their code is a mystery
- Other team members can't help or contribute

**Reality check**: Programming is about working with people, not avoiding them.

## Good Scenarios

### Collective Ownership in Practice

**Situation**: A team that practices true collective ownership.

**What it looks like**:
- Any team member can check out any module
- Changes are made where needed, not based on "ownership"
- Team members learn by working on unfamiliar parts
- No single points of failure in knowledge
- Code reviews happen naturally through pairing

**Benefits**:
- Less duplication (multiple eyes see patterns)
- Better interfaces (designed for team, not individual)
- Resilient team (anyone can cover for anyone)
- Continuous learning built into daily work

---

### Effective Pairing

**Situation**: Two developers working together on a challenging problem.

**What it looks like**:
- Sharing a workstation (or screen in remote)
- Discussing approaches before coding
- One types while other reviews in real-time
- Frequent role swaps
- Knowledge transfers naturally

**Why it works**:
- "Two heads are better than one"
- Real-time code review (most efficient form)
- Knowledge sharing happens automatically
- No code goes unreviewed
- Both can now work on this part of the system

**When to pair**:
- Complex problems (most efficient approach)
- Unfamiliar code areas (learning opportunity)
- Critical features (need the review)
- Emergency fixes (why do we pair then but not normally?)

---

### The Collaborative Workspace

**Situation**: A team arranged for maximum collaboration.

**What it looks like**:
- Team members sitting around tables
- Facing each other, not walls
- Can see body language and facial expressions
- Overhear when someone is frustrated
- Spontaneous conversations happen naturally
- Headphones off most of the time

**Benefits**:
- Serendipitous communication
- Problems surface quickly
- Help offered before it's requested
- Team feels like a unit
- Shared context and awareness

## Comparison Table

| Aspect | Bad Pattern | Good Pattern |
|--------|-------------|--------------|
| Code Ownership | Individual ("my code") | Collective (team code) |
| Physical Setup | Corners, backs turned | Tables, facing each other |
| Knowledge | Silos, single experts | Shared, everyone learns |
| Code Review | Rare, formal, after-the-fact | Constant through pairing |
| Communication | Scheduled meetings only | Continuous, serendipitous |
| Incentives | Tied to individual ownership | Tied to team outcomes |

## Key Takeaways

1. **Owned code is a disaster** - It creates duplication, bad interfaces, and political dysfunction
2. **Collective ownership works** - When everyone can work anywhere, the code improves
3. **Pairing is efficient** - If it works in emergencies, it works for regular work
4. **Physical collaboration matters** - Face each other, communicate constantly
5. **Programming is about people** - Accept it and learn to work well with others
