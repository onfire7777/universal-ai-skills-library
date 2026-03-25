# Coding Practices Examples

Scenarios demonstrating professional vs unprofessional coding behavior.

## Bad Scenarios

### The 3 AM Hero

**Situation**: Developer works 18-hour days, 60-70 hour weeks, feeling dedicated and professional. At 3 AM, solves a timing problem by having code send messages to itself.

**What happened**:
- The "solution" instituted a faulty design structure
- Everyone had to work around it constantly
- Caused strange timing errors and infinite loops
- Became team cruft surrounded by workarounds
- Years of accumulated patches and side effects

**Problems**:
- Tired brain chose wrong solution that "looked good enough"
- Confusion of long hours with professionalism
- Technical debt that lasted years
- Team morale impact (became a running joke)

**Lesson**: Dedication is about discipline, not hours. Eight good hours beats eighteen exhausted ones.

---

### The Worried Coder

**Situation**: Developer has argument with spouse, then tries to code. Sits with eyes on screen, fingers on keyboard, doing nothing.

**What happened**:
- Background process runs in mind reviewing the argument
- Physical stress felt in chest and stomach
- Forces self to write a line or two, but can't sustain
- Descends into "stupefied insensibility"
- Any code produced is trash

**Problems**:
- Wasted time producing nothing
- Code produced needs rework
- Prolonged personal issue by not addressing it

**Better approach**: Spend a dedicated hour addressing the worry, then return to coding with clearer mind.

---

### The Zone Addict

**Situation**: Developer measures self-worth by time spent in the Zone, enjoying the euphoria and sense of conquest.

**What happened**:
- Wrote more code, went through TDD loops faster
- Lost big-picture perspective
- Made decisions that needed reversal later
- Snapped at colleagues who interrupted

**Problems**:
- Net productivity lower due to rework
- Damaged team relationships
- Missed better design solutions

**Better approach**: Walk away when feeling the Zone coming on; find a pair partner.

---

### The Hope-Based Estimate

**Situation**: Trade show in 10 days. Developer's three-number estimate is 8/12/20 days. Developer hopes to make it in 10.

**What happened**:
- Didn't communicate reality to stakeholders
- Took shortcuts, worked extra hours
- Team had false hope
- Delayed necessary tough decisions
- Missed deadline anyway

**Problems**:
- Hope destroyed schedule and reputation
- Everyone avoided facing the real issue
- No fall-back plan was created
- Rushed code created technical debt

**Better approach**: Immediately communicate the 12-day nominal estimate, request scope reduction or fall-back plan.

---

### The False Delivery

**Situation**: Pressure to show progress. Developer convinces self that feature is "done enough" and moves to next task.

**What happened**:
- Practice became contagious on team
- Definition of "done" stretched further and further
- One team defined "done" as "checked-in" (didn't have to compile)
- Managers heard everything was fine
- "Freight train of unfinished work" hit the team

**Problems**:
- Status reports became fiction
- No visibility into real project state
- Catastrophic surprise when truth emerged

**Better approach**: Create independent definition of "done" with automated acceptance tests.

---

## Good Scenarios

### The Partitioned Worrier

**Situation**: Developer has financial worries nagging during work.

**What they did**:
- Recognized the background process consuming focus
- Dedicated one hour specifically to financial planning
- Didn't solve the problem but reduced anxiety
- Returned to coding with quieted background process
- Produced quality code for rest of day

**Why it works**:
- Acknowledges human reality
- One hour of worry-time beats full day of half-focused coding
- Professional allocation of mental resources

---

### The Gracious Interruptee

**Situation**: Developer deep in complex problem when colleague asks for help.

**What they did**:
- Stopped what they were doing politely
- Helped colleague work through their issue
- Spent 30 minutes pairing on the problem
- Returned to own work refreshed with new perspective
- Used failing test to reconstruct context

**Why it works**:
- Professional ethics require mutual help
- Fresh perspective often solves problems quickly
- TDD maintained context for return
- Built team trust and collaboration

---

### The Honest Estimator

**Situation**: Manager asks developer to "try" to make impossible deadline.

**What they did**:
- Held to original estimates
- Presented three-number estimate (best/nominal/worst)
- Suggested scope reduction as only viable option

**Why it works**: Original estimates more accurate than pressure-modified ones; no false hope.

---

### The Block Breaker

**Situation**: Developer sits at workstation but code won't come.

**What they did**:
- Recognized writer's block symptoms
- Found a pair partner
- Felt the physiological change that breaks blockage

**Why it works**: Pairing has chemical/physiological effect on brain that breaks mental logjams.

---

### The Strategic Walker

**Situation**: Developer stuck on problem late in the day, tempted to stay until solved.

**What they did**:
- Went home at normal time instead of pushing through
- Let subconscious work on problem overnight
- Solved problem in shower next morning

**Why it works**: Disengagement allows creative hunting; fresh perspective sees options intensity missed.

---

## Key Patterns

| Scenario | Unprofessional | Professional |
|----------|---------------|--------------|
| Tired | Push through, "heroic" hours | Go home, get sleep |
| Worried | Force yourself to code | Partition time for the worry |
| In Zone | Stay there, feel productive | Walk away, find partner |
| Interrupted | Snap, glare, protect turf | Help graciously, use TDD for context |
| Behind schedule | Hope, rush, false delivery | Honest estimates, reduce scope |
| Stuck | Stay stuck, avoid asking | Ask for help, pair up |
