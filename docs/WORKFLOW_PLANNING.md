# Using several skills on one task

For a simple request, preflight identifies a candidate and the agent can read its full skill. For a task with multiple deliverables, constraints, tools or independent checks, use the workflow planner. The planner searches canonical skill sections for every requirement. Each proposed contribution names a skill, its role, the source files and lines, and a revision. Gaps identify work that still needs attention.

Existing `skill-router skills compose --pipeline` and MCP `compose` with `pipeline:true` expose the requirement plan under `workflow` for multi-part requests. Their original `steps`, `edges` and inline token budget remain available. In CLI pipeline mode `--skills a,b`, or MCP `skills:["a","b"]`, requests those named passes in order. Pass the nested `workflow` object to stage loading and review.

## Identify the requirements

A natural prompt works directly from the CLI.

```sh
skill-router skills workflow plan -- 'Refresh the SaaS dashboard so it matches our existing brand and makes common tasks easier for keyboard users' > plan.json
```

For precise coverage, the host agent identifies deliverables and constraints, specifies their dependencies, and supplies them as a JSON request. Separate IDs preserve independent review passes, including repeated reviews by the same explicitly named skill.

```json
{
  "prompt": "Refresh the dashboard using our existing brand, and check keyboard navigation.",
  "requirements": [
    {"id":"brand","kind":"prerequisite","text":"Extract existing brand rules from the project","skills":["brand"]},
    {"id":"dashboard","kind":"deliverable","text":"Improve common dashboard tasks using the approved brand rules","dependsOn":["brand"]},
    {"id":"keyboard","kind":"verification","text":"Test keyboard navigation, visible focus and focus order","skills":["accessibility-testing"],"dependsOn":["dashboard"]}
  ]
}
```

Run `skill-router skills workflow plan --request request.json > plan.json`, or send the same object to the MCP tool `workflow_plan`. Use the optional `available` array for tools and runtimes whose availability the host has established. An empty array records a known empty inventory. When a dependency produces a gap, verify the prerequisite, include the confirmed value in a revised request and plan again. Hosts can request kinds `deliverable`, `constraint`, `prerequisite` and `verification`.

### Specify what the source must cover

Keep the original request in `prompt`. Extract the details that would change the workflow, such as a keyboard test on an Android phone, a native Word template, or a request to preserve duplicate payment records. Define each detail under `facets`. An `anyOf` array gives alternative source phrases for **one** detail. Give independent details separate IDs. Use the same facets in `workflow_search` to inspect individual sections.

For example, a requirement to edit and check an existing Word report can describe the file format, permitted editing approach and two separate acceptance criteria.

```json
{
  "id": "word",
  "kind": "deliverable",
  "text": "Edit the existing Word template and verify its exception table and retained styles",
  "facets": [
    {"id":"format","kind":"format","anyOf":["DOCX","Word document"],"carry":true},
    {"id":"edit-template","kind":"action","anyOf":["edit an existing document","word/document.xml"]}
  ],
  "checks": [
    {"id":"rows","text":"Reopen the file and compare every exception row with the expected amounts","evidenceKind":"test"},
    {"id":"styles","text":"Compare the template styles before and after editing","evidenceKind":"test"}
  ]
}
```

Facet kinds are `object`, `action`, `platform`, `format`, `check` and `exclude`. An `exclude` facet rules out sections that mention a conflicting alternative. For constraints that apply to later stages, use `carry:true` and connect those stages with `dependsOn`. Inspect each stage’s `facetCoverage`: matches cite a contribution, the indexed section, source line and matching phrase. If a required facet lacks a match, inspect the source or refine the requirement before applying it.

An acceptance check specifies what the agent will demonstrate in the project. Give each distinct outcome a check ID and an `evidenceKind` of `test`, `artifact` or `decision`. A source facet of kind `check` can find a verification technique. The corresponding acceptance check still requires an observation from performing the work.

## Refine a requirement

When a contribution misses the project context or an important check is uncovered, search specific sections for that requirement. Keep the original task in `prompt` and the narrower intent in `query`.

```json
{"prompt":"Improve dashboard keyboard navigation","query":"focus order and visible indicators","limit":5}
```

Pass this object to MCP `workflow_search` or CLI `skill-router skills workflow search --request search.json`. Narrow by `skill` if necessary or retrieve the next page with `offset`. Inspect the returned title, path, section, applicability and prerequisites. Pin the applicable sections in that requirement’s `elements` array using the returned `skill`, `file`, `startLine`, `endLine` and `digest`. Then call `workflow_plan` again. Refresh the search when the source digest changes. For an exact set of named skills or pinned sections, set the requirement’s `selection` to `explicit`. Without this setting the planner can also suggest complementary elements.

## Read and apply each stage

Send the unedited plan object and a stage ID to MCP `workflow_load`. The CLI accepts the same JSON. For example, with jq installed:

```sh
jq '{plan: ., stage: "keyboard"}' plan.json | skill-router skills workflow load > keyboard-load.json
```

Read each contribution's full `body` and its relevant `assets`. Text assets carry their source text when available. Open a path-only asset locally if the selected step calls for it. Use each skill's referenced workflow, installed native CLI, or script as appropriate. Keep the receipt returned by the load alongside each contribution's ID. A shared action can use one observation for multiple elements when the action genuinely addressed all of them. Requested independent checks each need their own observation.

When a stage is large, ask `workflow_load` for `"mode":"manifest"`. Its `assets` list includes the full entrypoint and each relevant file, with its checked absolute path, byte count and SHA-256. Manifest responses contain no load receipts. Blocked contributions and their reasons remain visible alongside any contributions ready to read.

Load one contribution at a time by adding `"contribution":"<id from the plan>"` to the stage request. If the MCP client truncates full output, save the complete CLI response in the task workspace. For example:

```sh
jq --arg id "$contribution_id" '{plan: ., stage: "keyboard", contribution:$id}' plan.json |
  skill-router skills workflow load > contribution-load.json
jq -j '.contributions[0].body' contribution-load.json > contribution-skill.md
```

Read the complete entrypoint and relevant assets in manageable ranges. Compare sizes and digests when opening files by path. Retain the full response receipt after reading the inputs. If only metadata or a truncated preview is available, record the outstanding read as a gap.

The loader checks the plan and requested bundle against the current checkout. If the checkout or requirements change, generate a fresh plan before continuing.

## Reconcile observed evidence

Maintain an array of records as work progresses. Use status `applied`, `verified`, `skipped`, or `blocked`. Applied and verified records include the contribution ID, the receipt from its load response, and the exact elements addressed as `file#Lstart-Lend`. A record for a single-element contribution can omit `elements`. Include bundle-relative paths in `readFiles` for path-only assets that were consulted. Record the action's actual outcome as evidence:

| Evidence kind | Reference | Observation |
| --- | --- | --- |
| `decision` | Project rule or recorded design decision | The specific choice the agent made and where it applied |
| `artifact` | Path and location of the change | What the agent produced or inspected |
| `test` | Test command, run or output file | What happened and the outcome `passed` or `failed` |

When a requirement declares acceptance checks, include their IDs in each observation’s `checks` list. An observation can address several related checks when its reference and result genuinely cover all of them. A failed test or a missing check keeps that criterion open. For example:

```json
{"kind":"test","checks":["rows"],"reference":"verify-report rows reconciliation.docx","observation":"Reopened the Word file and matched all three exception rows, including currency deltas.","outcome":"passed"}
```

Inspect `stages[].checks` in the review result for each criterion’s status and evidence. In this example the style comparison remains open until its observed result is recorded.

Skipped and blocked records require a reason. Verification requires an observed passing check. A planned check or a load receipt alone leaves the verification outstanding. The router evaluates the evidence supplied by the host. The host performs the task and runs the checks.

```sh
jq --slurpfile records records.json '{plan: ., records: $records[0]}' plan.json | skill-router skills workflow review
```

The same request object works with MCP `workflow_review`. Keep records for unchanged contribution IDs when revising a plan. Review reports obsolete records and outstanding work. Task-local plan and record files allow an agent to resume long workflows.

## Index coverage

Run `skill-router skills workflow index` to see indexed bundle, file and element counts and how many bundles were rebuilt or reused. Add `--inventory` for a complete map from skill name to each relative file path and its SHA-256. The index uses the configured canonical checkout and keeps its derived cache locally.

## Enable the Codex adapter

From the configured checkout run `skill-router sync codex`. When the Codex CLI is available, sync also registers `universal-ai-skills` with that executable and checkout. An existing registration retains its launch command and custom settings. Use `--mcp=false` for CLI-only installation. This installs the compact router skill, updates the router section of the active Codex home’s `AGENTS.md`, and registers its user-prompt hook alongside existing hooks such as Mnemosyne. Codex owns hook trust. On a fresh profile use Codex's hook review controls to enable the new registration. The global instructions provide a fallback when a trusted hook is unavailable.

To configure MCP separately, use the native Codex CLI from the library checkout:

```sh
codex mcp add --env "SKILL_ROUTER_REPO_DIR=$PWD" universal-ai-skills -- "$(command -v skill-router)" serve
codex mcp get universal-ai-skills
```

After configuration changes, open a fresh Codex session. On multi-part work, Codex uses planning even when the single-skill preflight only covers one requirement. Simple work retains direct skill loading.
