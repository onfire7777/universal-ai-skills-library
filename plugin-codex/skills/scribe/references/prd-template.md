# PRD (Product Requirements Document) Template

Purpose: Use this file for business-facing requirements, user value, goals, non-goals, and testable acceptance criteria.

Contents:

- full PRD template
- quick PRD template
- PRD quality checklist

## Template

```markdown
# PRD: [Feature Name]

## Document Info
| Field | Value |
|-------|-------|
| Version | v0.1 |
| Status | Draft |
| Author | [name] |
| Reviewers | [PM, Tech Lead, QA] |
| Audience | PM, Dev, QA |
| Related Docs | [links] |

## Change History
| Date | Version | Author | Change |
|------|---------|--------|--------|
| YYYY-MM-DD | v0.1 | [name] | Initial draft |

## 1. Overview
### 1.1 Purpose
[Why this document exists]

### 1.2 Background
[Business context, problem, timing]

### 1.3 Goals
- [Goal 1]
- [Goal 2]

### 1.4 Non-Goals (Out of Scope)
- [Non-goal 1]
- [Non-goal 2]

### 1.5 Success Metrics
| Metric | Target | Measurement |
|--------|--------|-------------|
| Activation rate | +15% | product analytics |
| Completion rate | >= 70% | funnel report |

## 2. User Stories
### 2.1 Target Users
- Persona: [who]
- Context: [when and why]

### 2.2 User Stories
As a [user]
I want [goal]
So that [benefit]

## 3. Functional Requirements
### 3.1 [Feature Area]
**REQ-001**: [Requirement title]
- Description: [what must happen]
- Priority: Must / Should / Could
- Trigger: [action or event]
- Output: [observable result]
- Dependencies: [REQ-XXX or external dependency]
- Acceptance Criteria: `AC-001`

**REQ-002**: [Requirement title]
- Description: ...
- Acceptance Criteria: `AC-002`

## 4. Non-Functional Requirements
### 4.1 Performance
**NFR-001**: Response time
- p95 latency: [target]

### 4.2 Security
**NFR-002**: Authentication
- JWT token expiry: `24 hours`

### 4.3 Scalability
**NFR-003**: Load capacity
- [concurrency / throughput target]

### 4.4 Availability
**NFR-004**: Uptime
- [availability target]

### 4.5 Accessibility
**NFR-005**: WCAG compliance
- [target level]

## 5. Acceptance Criteria
**AC-001**: [Criterion for `REQ-001`]
Given [precondition]
When [action]
Then [result]
And [additional check]

**AC-002**: [Criterion for `REQ-002`]
Given ...
When ...
Then ...

## 6. Edge Cases & Error Handling
- [edge case]
- [failure mode]
- [fallback]

## 7. UI/UX Requirements
### 7.1 Wireframes
[link or note]

### 7.2 UI Specifications
- states
- accessibility needs
- copy constraints

### 7.3 Interaction Flow
[step sequence]

## 8. Technical Constraints
### 8.1 Technology Stack
- [relevant constraints only]

### 8.2 Integration Points
- [system or API]

### 8.3 Data Requirements
- [storage, retention, privacy, migration]

## 9. Dependencies
### 9.1 Internal Dependencies
- [team or component]

### 9.2 External Dependencies
- [vendor, SDK, contract]

## 10. Risks & Mitigations
| Risk | Impact | Mitigation |
|------|--------|------------|
| [risk] | [level] | [action] |

## 11. Timeline
| Milestone | Date | Owner |
|-----------|------|-------|
| Draft approved | YYYY-MM-DD | [owner] |

## 12. Open Questions
- [question]

## 13. Appendix
### 13.1 Glossary
- [term]

### 13.2 References
- [doc link]

### 13.3 Traceability Matrix
| Requirement | Design | Test | Code / Doc Target |
|-------------|--------|------|-------------------|
| `REQ-001` | `HLD-3.1` | `TC-001` | `auth.service.ts` |
```

## Quick Template (Minimal)

```markdown
# PRD: [Feature Name]

## Overview
[problem, user, value]

## User Story
As a [user], I want [goal], so that [benefit].

## Requirements
| ID | Description | Priority | Acceptance Criteria |
|----|-------------|----------|---------------------|
| `REQ-001` | [Description] | Must | [Criteria] |
| `REQ-002` | [Description] | Should | [Criteria] |

## Edge Cases
- [edge case]

## Out of Scope
- [non-goal]

## Open Questions
- [question]
```

## PRD Quality Checklist

- [ ] All requirements have IDs (`REQ-XXX`)
- [ ] Every requirement maps to acceptance criteria
- [ ] `Non-Goals` is present
- [ ] Success metrics are measurable
- [ ] PRD explains `What` and `Why`, not implementation `How`
- [ ] Edge cases and dependencies are explicit
- [ ] Traceability matrix exists when downstream design or tests are expected
