# Testing Strategies Knowledge

Core concepts and foundational understanding for professional testing strategies.

## Overview

Professional developers need more than unit tests and acceptance tests - they need a comprehensive testing strategy. The goal is that QA should find nothing wrong, achieved through a hierarchy of automated tests at different levels of the system.

## Key Concepts

### QA Should Find Nothing

**Definition**: The development team's goal should be that QA discovers zero defects.

Every time QA finds something, the development team should react with concern and take steps to prevent similar issues in the future. While this goal may not always be achieved, it should drive the team's testing practices.

**Key points**:
- QA finding bugs indicates gaps in the testing strategy
- Each discovered bug is an opportunity to improve prevention
- Development owns quality, not QA

### QA as Part of the Team

**Definition**: QA and Development work together toward quality, not as adversaries.

QA serves two critical roles that complement development:

**Key points**:
- QA acts as specifiers (creating acceptance tests from requirements)
- QA acts as characterizers (exploring actual system behavior)
- Business writes happy-path tests; QA writes edge cases

### The Test Automation Pyramid

**Definition**: A hierarchical structure of automated tests, with unit tests forming the large base and exploratory tests at the small top.

The pyramid represents both the quantity and granularity of tests needed at each level. More tests at the bottom (fast, focused), fewer at the top (slow, broad).

**Key points**:
- Unit tests form the foundation (highest quantity)
- Each level builds on the confidence of levels below
- Higher levels test integration, not business rules

## Terminology

| Term | Definition |
|------|------------|
| Unit Test | Tests written by programmers to specify system at lowest level |
| Component Test | Acceptance tests for individual components and business rules |
| Integration Test | Tests how component assemblies communicate together |
| System Test | Tests against entire integrated system |
| Exploratory Test | Manual, unscripted testing to find unexpected behaviors |
| Test Doubles | Mocks/stubs used to isolate components under test |
| Choreography Test | Tests that verify components work together (integration) |

## How It Relates To

- **TDD**: Unit tests from the pyramid's base come from Test-Driven Development
- **Acceptance Testing**: Component tests are the acceptance tests for business rules
- **Continuous Integration**: Lower-level tests run on every commit; higher levels run periodically

## Common Misconceptions

- **Myth**: More tests at every level is always better
  **Reality**: The pyramid shape is intentional - more unit tests, fewer system tests

- **Myth**: QA is responsible for finding all bugs
  **Reality**: Development should prevent bugs; QA validates the prevention worked

- **Myth**: Manual testing can be replaced entirely by automation
  **Reality**: Exploratory testing requires human creativity and cannot be scripted

## Quick Reference

| Concept | One-Line Summary |
|---------|-----------------|
| QA Finds Nothing | Development's goal is zero defects reaching QA |
| QA as Specifiers | QA translates requirements into acceptance tests |
| QA as Characterizers | QA identifies actual system behavior through exploration |
| Test Pyramid | Unit (90%+) > Component (50%) > Integration > System (10%) > Exploratory |
