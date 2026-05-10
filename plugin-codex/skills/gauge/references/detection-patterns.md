# Detection Patterns

**Purpose:** Structural detection rules for each of the 16 checklist items.
**Read when:** Executing the CLASSIFY phase of the audit workflow.

---

## F1: YAML Frontmatter

**Detection:**
1. Check first line is exactly `---`
2. Scan for closing `---` (second occurrence)
3. Between delimiters, check for `name:` field (non-empty)
4. Between delimiters, check for `description:` field (non-empty)

**PARTIAL trigger:** Delimiters present but `name:` or `description:` missing/empty.

---

## L1: Language Compliance

**Detection:**
1. Extract `description:` value from frontmatter → check for Japanese characters (hiragana: `\u3040-\u309F`, katakana: `\u30A0-\u30FF`, kanji: `\u4E00-\u9FFF`)
2. Extract body text (everything after frontmatter closing `---`)
3. Scan body for Japanese character ranges (same as above)
4. Apply whitelist exclusions:
   - Agent names (e.g., `Nexus`, `Architect`)
   - Technical identifiers (e.g., `_STEP_COMPLETE`, `NEXUS_HANDOFF`)
   - Quoted proper nouns
5. Scan `references/` files for Japanese characters (none allowed)

**PARTIAL trigger:** 1-3 Japanese character instances in body after whitelist filtering.
**FAIL trigger:** 4+ instances in body, OR `description:` is English-only, OR `references/` files contain Japanese.

---

## H1: CAPABILITIES_SUMMARY

**Detection:**
1. Scan for HTML comment block: `<!--` ... `-->`
2. Within comment, search for `CAPABILITIES_SUMMARY:` keyword
3. After keyword, count entries matching pattern: `- key: description`
4. Verify minimum 3 entries

**PARTIAL trigger:** Keyword found but fewer than 3 entries or entries lack `key: description` format.

---

## H2: COLLABORATION_PATTERNS

**Detection:**
1. Within HTML comment block, search for `COLLABORATION_PATTERNS:` keyword
2. After keyword, count entries matching pattern: `- Agent -> Agent: description` or `- Agent → Agent: description`
3. Also check for `BIDIRECTIONAL_PARTNERS:` subsection with `INPUT:` and `OUTPUT:` lines

**PARTIAL trigger:** Keyword found but missing direction arrows, agent names, or descriptions.

---

## H3: PROJECT_AFFINITY

**Detection:**
1. Within HTML comment block, search for `PROJECT_AFFINITY:` keyword
2. After keyword, check for either:
   - `universal` keyword, OR
   - Domain ratings in format `Domain(H|M|L)` (minimum 1)

**PARTIAL trigger:** Keyword found but no ratings and no `universal` keyword.

---

## S1: Trigger Guidance

**Detection:**
1. Search for heading matching: `## Trigger Guidance` or `## Trigger` (case-insensitive)
2. Under heading, search for "Use" + agent name + "when" pattern (positive triggers)
3. Search for "Route elsewhere" or "Route to" pattern (negative routing)

**PARTIAL trigger:** Heading exists but missing positive triggers or negative routing.

---

## S2: Core Contract

**Detection:**
1. Search for heading matching: `## Core Contract` or `## Contract` (case-insensitive)
2. Count bullet points or numbered items under heading
3. Verify minimum 3 items

**PARTIAL trigger:** Heading exists but fewer than 3 contract items.

---

## S3: Boundaries

**Detection:**
1. Search for heading matching: `## Boundaries` (case-insensitive)
2. Under heading, search for three subsections:
   - `### Always` or `**Always**` or `**Always:**`
   - `### Ask` or `**Ask**` or `**Ask First**` or `**Ask First:**`
   - `### Never` or `**Never**` or `**Never:**`
3. Verify each subsection has at least 1 item

**PARTIAL trigger:** Boundaries heading exists but one subsection missing or empty.
**FAIL trigger:** No Boundaries heading, or 2+ subsections missing.

---

## S4: Workflow

**Detection:**
1. Search for heading matching: `## Workflow` (case-insensitive)
2. Search for phase pipeline pattern: text containing `→` or `->` with 2+ phase names (typically backtick-wrapped)
3. Search for phase table: markdown table with columns including "Phase" and at least one of "Action", "Rule", "Read"

**PARTIAL trigger:** Pipeline or table present but not both.

---

## S5: Output Routing

**Detection:**
1. Search for heading matching: `## Output Routing` (case-insensitive)
2. Search for routing table: markdown table with columns including "Signal" and at least one of "Approach", "Output", "Read"
3. Count table rows (minimum 3)

**PARTIAL trigger:** Table present but missing key columns or fewer than 3 routes.

---

## S6: Output Requirements

**Detection:**
1. Search for heading matching: `## Output Requirements` (case-insensitive)
2. Count bullet points or numbered items under heading
3. Verify minimum 3 items

**PARTIAL trigger:** Heading exists but fewer than 3 requirement items.

---

## S7: Collaboration

**Detection:**
1. Search for heading matching: `## Collaboration` (case-insensitive)
2. Search for `Receives:` or `**Receives:**` pattern with agent names
3. Search for `Sends:` or `**Sends:**` pattern with agent names

**PARTIAL trigger:** Heading exists but missing either Receives or Sends.

---

## S8: Reference Map

**Detection:**
1. Search for heading matching: `## Reference Map` or `## References` (case-insensitive)
2. Search for table with columns: "Reference" (or "File") and "Read this when" (or "Content", "When")
3. Verify listed files exist in `references/` directory

**Special case:** If the agent has no `references/` directory, explicit "N/A" or "No references" is acceptable as PASS.

**PARTIAL trigger:** Table present but references point to non-existent files.

---

## S9: Operational

**Detection:**
1. Search for heading matching: `## Operational` (case-insensitive)
2. Check for journal mention: `.agents/{agent_name}.md` pattern
3. Check for PROJECT.md mention: `.agents/PROJECT.md` or `PROJECT.md` pattern
4. Check for standard protocols link: `_common/OPERATIONAL.md` pattern

**PARTIAL trigger:** Heading exists but missing one of the three elements.

---

## A1: AUTORUN Support (_STEP_COMPLETE)

**Detection:**
1. Search for heading matching: `## AUTORUN` (case-insensitive)
2. Search for `_AGENT_CONTEXT` keyword in section body
3. Search for `_STEP_COMPLETE` keyword
4. Within _STEP_COMPLETE block, verify presence of:
   - `Agent:` field
   - `Status:` field with `SUCCESS | PARTIAL | BLOCKED | FAILED`
   - `Output:` block
   - `Next:` field
   - `Reason:` field

**PARTIAL trigger:** Section exists but _STEP_COMPLETE block missing required fields.

---

## A2: Nexus Hub Mode (NEXUS_HANDOFF)

**Detection:**
1. Search for heading matching: `## Nexus Hub Mode` or `## Nexus` (case-insensitive, must relate to hub mode)
2. Search for `NEXUS_ROUTING` keyword in section body
3. Search for `NEXUS_HANDOFF` keyword
4. Within NEXUS_HANDOFF block, verify presence of:
   - `Step:` field
   - `Agent:` field
   - `Summary:` field
   - `Next action:` field with `CONTINUE | VERIFY | DONE`

**PARTIAL trigger:** Section exists but NEXUS_HANDOFF block missing required fields.

---

## Compound Detection Rules

### Full HTML Comment Block Validation (H1 + H2 + H3)

All three items (H1, H2, H3) must reside within a single HTML comment block (`<!-- ... -->`). If no HTML comment block exists, all three are FAIL. If the block exists but is missing individual sections, score each independently.

### Language Compliance Cross-Check (L1)

L1 applies across all sections and references. When auditing any section (S1-S9), simultaneously check for Japanese text violations. This avoids a second pass.
