# Normalization Checklist (Audit Reference)

**Purpose:** Defines the 16-item checklist with PASS/PARTIAL/FAIL criteria and P0-P3 priority classification.
**Read when:** Starting any SCAN or CLASSIFY phase.
**Source:** `.agents/skill-normalization-checklist.md` (ecosystem master checklist).

---

## 16-Item Checklist

### F1: YAML Frontmatter

| Status | Criteria |
|--------|----------|
| PASS | `---` delimiters present, contains `name:` and `description:` fields |
| PARTIAL | Delimiters present but missing `name:` or `description:` |
| FAIL | No YAML frontmatter block |

### L1: Language Compliance

| Status | Criteria |
|--------|----------|
| PASS | `description:` is in Japanese; body (headings, text, tables, bullets) is in English; `references/` files are English only |
| PARTIAL | Body contains minimal Japanese (1-3 instances) outside technical terms/proper nouns |
| FAIL | Body contains significant Japanese text (4+ instances) or `description:` is in English |

**Detection:** Scan body for hiragana/katakana/kanji characters. Whitelist: agent names, tool names, technical terms, proper nouns.

### H1: CAPABILITIES_SUMMARY

| Status | Criteria |
|--------|----------|
| PASS | HTML comment block contains `CAPABILITIES_SUMMARY:` with 3+ capability entries in `key: description` format |
| PARTIAL | Block present but fewer than 3 entries or malformed format |
| FAIL | No CAPABILITIES_SUMMARY block |

### H2: COLLABORATION_PATTERNS

| Status | Criteria |
|--------|----------|
| PASS | HTML comment block contains `COLLABORATION_PATTERNS:` with directed patterns (Agent -> Agent: description) |
| PARTIAL | Block present but missing direction arrows or descriptions |
| FAIL | No COLLABORATION_PATTERNS block |

### H3: PROJECT_AFFINITY

| Status | Criteria |
|--------|----------|
| PASS | HTML comment block contains `PROJECT_AFFINITY:` with domain ratings (e.g., `Game(H) SaaS(L)`) or `universal` |
| PARTIAL | Block present but no ratings or `universal` keyword |
| FAIL | No PROJECT_AFFINITY line |

### S1: Trigger Guidance

| Status | Criteria |
|--------|----------|
| PASS | Section with heading containing "Trigger" exists; includes "Use X when" list AND "Route elsewhere" list |
| PARTIAL | Section exists but missing one of the two lists |
| FAIL | No Trigger Guidance section |

### S2: Core Contract

| Status | Criteria |
|--------|----------|
| PASS | Section with heading containing "Core Contract" or "Contract" exists; contains 3+ commitment items |
| PARTIAL | Section exists but fewer than 3 items |
| FAIL | No Core Contract section |

### S3: Boundaries

| Status | Criteria |
|--------|----------|
| PASS | Section contains all three subsections: Always, Ask (First), Never |
| PARTIAL | Section exists but missing one subsection |
| FAIL | No Boundaries section or missing two+ subsections |

### S4: Workflow

| Status | Criteria |
|--------|----------|
| PASS | Section contains phase pipeline (e.g., `PHASE1 → PHASE2 → ...`) AND phase table with columns |
| PARTIAL | Pipeline or table present but not both |
| FAIL | No Workflow section |

### S5: Output Routing

| Status | Criteria |
|--------|----------|
| PASS | Section contains signal-to-output routing table with columns: Signal, Approach/output, Read next |
| PARTIAL | Table present but missing columns or fewer than 3 routes |
| FAIL | No Output Routing section |

### S6: Output Requirements

| Status | Criteria |
|--------|----------|
| PASS | Section lists 3+ required elements for every deliverable |
| PARTIAL | Section exists but fewer than 3 elements |
| FAIL | No Output Requirements section |

### S7: Collaboration

| Status | Criteria |
|--------|----------|
| PASS | Section contains both `Receives:` and `Sends:` entries with agent names and context |
| PARTIAL | Section exists but missing Receives or Sends |
| FAIL | No Collaboration section |

### S8: Reference Map

| Status | Criteria |
|--------|----------|
| PASS | Section contains reference table with `Reference` and `Read this when` columns; all listed files exist in `references/` |
| PARTIAL | Table present but incomplete or references missing |
| FAIL | No Reference Map section (N/A acceptable if agent has no `references/` directory) |

### S9: Operational

| Status | Criteria |
|--------|----------|
| PASS | Section mentions journal file (`.agents/{name}.md`), PROJECT.md logging, and links to `_common/OPERATIONAL.md` |
| PARTIAL | Section exists but missing one of the three elements |
| FAIL | No Operational section |

### A1: AUTORUN Support (_STEP_COMPLETE)

| Status | Criteria |
|--------|----------|
| PASS | Section describes `_AGENT_CONTEXT` parsing AND contains `_STEP_COMPLETE` YAML block with Status, Output, Next, Reason fields |
| PARTIAL | Section exists but _STEP_COMPLETE block is incomplete |
| FAIL | No AUTORUN Support section |

### A2: Nexus Hub Mode (NEXUS_HANDOFF)

| Status | Criteria |
|--------|----------|
| PASS | Section describes `NEXUS_ROUTING` detection AND contains `NEXUS_HANDOFF` block with Step, Agent, Summary, Next action fields |
| PARTIAL | Section exists but NEXUS_HANDOFF block is incomplete |
| FAIL | No Nexus Hub Mode section |

---

## Priority Classification

| Priority | Scope | Items | Rationale |
|----------|-------|-------|-----------|
| **P0 (Critical)** | Ecosystem integration | A1, A2, S7 | Required for AUTORUN execution and agent collaboration |
| **P1 (High)** | Quality and consistency | S2, S3, S4, L1 | Core behavioral contract, safety boundaries, workflow definition, language standard |
| **P2 (Medium)** | Discoverability and routing | S1, S5, S6 | Trigger conditions, output routing, deliverable requirements |
| **P3 (Low)** | Metadata and reference | F1, H1, H2, H3, S8, S9 | Frontmatter, machine-readable metadata, reference pointers, operational logging |

### Priority-Based Fix Order

1. Fix all P0 violations first (ecosystem cannot function without these).
2. Fix P1 violations (quality and behavioral consistency).
3. Fix P2 violations (improve discoverability and routing accuracy).
4. Fix P3 violations (metadata completeness).

---

## Quest Exemplar Reference

Quest (`quest/SKILL.md`) is the reference standard for all 16 items. When generating fix snippets, cite the corresponding Quest section as exemplar:

| Item | Quest Section Reference |
|------|------------------------|
| F1 | Lines 1-4 (YAML frontmatter) |
| L1 | Full file (Japanese description, English body) |
| H1 | Lines 6-17 (CAPABILITIES_SUMMARY) |
| H2 | Lines 19-32 (COLLABORATION_PATTERNS) |
| H3 | Lines 38 (PROJECT_AFFINITY) |
| S1 | "Trigger Guidance" section |
| S2 | "Core Contract" section |
| S3 | "Boundaries" section (Always/Ask First/Never) |
| S4 | "Workflow" section (pipeline + phase table) |
| S5 | "Output Routing" section (signal table + routing rules) |
| S6 | "Output Requirements" section |
| S7 | "Collaboration" section (Receives/Sends + overlap boundaries) |
| S8 | "Reference Map" section |
| S9 | "Operational" section |
| A1 | "AUTORUN Support" section (_STEP_COMPLETE YAML) |
| A2 | "Nexus Hub Mode" section (NEXUS_HANDOFF block) |
