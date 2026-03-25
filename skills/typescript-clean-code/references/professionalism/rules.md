# Professionalism Rules

Guidelines for professional behavior in software development.

## Core Rules

### 1. Do No Harm to Function

Create software that works correctly. You are responsible for bugs even though they are inevitable.

- Never release code you aren't certain about
- Apologize for bugs, then prevent recurrence
- Your error rate should decrease over your career toward zero

### 2. Do No Harm to Structure

Never sacrifice code structure for short-term function delivery.

- Flexible structure enables future changes
- Compromising structure compromises the future
- You must be able to make changes without exorbitant costs

### 3. QA Should Find Nothing

Release only code you expect to pass QA completely.

- Never use QA as bug catchers
- Sending known-faulty code is unprofessional
- Be surprised and chagrined when QA finds issues
- Every escaped bug requires root cause analysis

### 4. Know It Works

Test thoroughly before release.

- Test it up, down, and seven ways to Sunday
- Automate your tests for quick, repeatable execution
- 100% test coverage is demanded, not suggested
- Design code to be easy to test (write tests first)

### 5. Maintain Automated QA

Your test suite should be your release gate.

- If tests pass, you should be confident to ship
- Tests should run quickly (minutes, not hours)
- At minimum, tests should indicate high probability of passing QA

### 6. Practice Merciless Refactoring

Continuously improve code structure.

- Make small improvements every time you touch code
- The Boy Scout Rule: leave code cleaner than you found it
- Tests enable fearless refactoring
- Treat code like clay - continuously shape it

### 7. Own Your Career

Your professional development is your responsibility.

- Don't rely on employer for training, books, or conferences
- Invest ~20 hours/week in your own growth
- 40 hours for employer's problems, 20 for your career
- If employer provides learning opportunities, be grateful but don't expect them

### 8. Know Your Field

Maintain broad knowledge of software development history and techniques.

- Know design patterns (all 24 GOF patterns, POSA patterns)
- Know design principles (SOLID, component principles)
- Know methods (XP, Scrum, Lean, Kanban, Waterfall, Structured Analysis)
- Know disciplines (TDD, OO design, CI, Pair Programming)
- Know artifacts (UML, DFDs, State Diagrams, decision tables)

### 9. Never Stop Learning

Continuous learning is mandatory for professionals.

- Read books, articles, blogs
- Attend conferences and user groups
- Learn outside your comfort zone
- Architects who stop coding become irrelevant

### 10. Practice Deliberately

Performance is not practice. Practice skills outside daily work.

- Use kata exercises for skill sharpening
- 10-minute warm-up in morning, cool-down in evening
- Practice in multiple languages
- Train fingers and brain through repetition

### 11. Collaborate and Teach

Learning accelerates through interaction with others.

- Program together, practice together, design together
- Mentoring juniors is a professional responsibility
- Teaching drives knowledge deeper into your own understanding
- Balance collaboration with necessary alone time

### 12. Know Your Domain

Understand the business context of your software.

- Read books on the domain
- Interview customers and users
- Challenge specification errors with domain knowledge
- Never code blindly from specs

### 13. Identify with Employer/Customer

Their problems are your problems.

- Work toward solutions that address real needs
- Avoid "us versus them" mentality
- Put yourself in employer's shoes when developing features

### 14. Practice Humility

Balance confidence with awareness of your fallibility.

- Know your job and take pride in your work
- Take bold, calculated risks based on confidence
- Accept that you will sometimes fail
- Be first to laugh when you're the butt of a joke
- Never demean others for mistakes

## Guidelines

Less strict recommendations for professional behavior:

- Use lunch hours for reading
- Listen to podcasts during commute
- Spend 90 minutes daily learning something new
- Do kata to maintain skills in multiple languages
- Pair program regularly but preserve alone time
- Sit with new team members to show them the ropes

## Exceptions

When these rules may be relaxed:

- **Time investment**: Life circumstances may require temporary adjustment, but this should be exception not norm
- **Test coverage**: Some mission-critical systems may need additional QA beyond automated tests
- **Domain expertise**: You needn't be a domain expert, but due diligence is required

## Quick Reference

| Rule | Summary |
|------|---------|
| Do No Harm - Function | Don't create bugs; own the ones that escape |
| Do No Harm - Structure | Never sacrifice design for features |
| QA Finds Nothing | Release only code you're certain about |
| Know It Works | Test everything, automate everything |
| Merciless Refactoring | Improve code every time you touch it |
| Own Your Career | Invest 20 hours/week in growth |
| Know Your Field | Master patterns, principles, methods, disciplines |
| Never Stop Learning | Read, attend, participate continuously |
| Practice | Do kata daily to sharpen skills |
| Collaborate | Program and learn together |
| Know Domain | Understand the business you're coding for |
| Stay Humble | Accept you will fail; laugh at yourself |
