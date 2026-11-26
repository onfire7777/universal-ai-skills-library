# Workflow breadth and specificity evaluation

The suite registers 24 prompts across creative production, technical work, analysis, writing, accessibility, operations, productivity, education and research. It has 55 positively sourced requirements and two expected catalog gaps. Four prompts require four distinct operations across skill families. The source witnesses and near-misses were recorded before the first routing comparison. Witness paths establish that the catalog describes the technique. Runtime eligibility applies to every selection, and any eligible section supporting the registered method can earn credit.

Run from the repository root with a router built from the checkout under test.

```sh
SKILL_ROUTER_BIN=/path/to/skill-router python3 tests/workflow-eval/breadth/evaluate.py
```

Use `--output /tmp/workflow-breadth.json` for case-level results and exact evidence locations. Use `--case diagram` to inspect one case. Use `--baseline /path/to/previous/skill-router --baseline-repo /path/to/previous/checkout` to compare against a frozen release on its matching corpus. A repeated local run can reuse `--cache /tmp/workflow-breadth-cache`. The evaluator keeps one MCP process per binary and evaluates three input modes. Frozen-version comparisons use natural and structured inputs, which are shared by the reference version.

## What the scores mean

- **Natural prompt coverage** counts registered requirements supported by a section selected for a matching stage in the original user request.
- **Structured coverage** measures the same requirements when a host supplies the separate deliverables and dependencies. Inputs contain requirement text and a dependency on the preceding requested operation to exercise stage handoffs. These cases use a linear host-supplied order. Skill names, source witnesses and section pins stay in the evaluator.
- **Faceted coverage** uses the structured requirements with separate object, action, format, platform and verification annotations. The annotations in `host-facets.json` express the request in ordinary language. The evaluator verifies every reported positive facet quote, element address and source line. It still awards requirement credit only for a registered source method.
- **Grounded contribution precision** is the number of selected contributions whose sections meet at least one applicable source criterion, divided by all selected contributions. A contribution that receives no credit may contain incidental advice. The separate near-miss field identifies specific, preregistered substitutions.
- **Section precision** applies the same test to individual selected elements. A mention in the plan or a relevant bundle name cannot satisfy a source criterion. The evaluator checks inclusive line ranges, raw file SHA-256 and section SHA-256 against the checkout, then searches a bounded section excerpt for evidence of the method. Legitimate equivalent sources can earn credit.
- **Catalog gaps** earn credit for an explicit gap with no substitute contribution, or newly available source evidence meeting the registered technical criterion.
- **Deep completion** requires coverage of every requirement and methods assigned across at least three separate stages on a four-operation task.

Use `--gate` for the maintained release gate. The floors in `acceptance.json` require at least 16 of the 18 single-family workflows to have a grounded operation, two of four deep workflows to cover all operations over multiple executable stages, and the two diagnostic gap cases to receive appropriate handling. Each mode must ground at least half of its contributions. The per-mode recall floors are 33, 36 and 39 of 55 positively sourced needs, respectively. Natural and structured results must also retain the aggregate coverage and precision of the frozen reference. `baseline.json` includes the immutable source commit, binary and manifest fingerprints, requirement-level results and unexpected selections.

The source predicates recognize specified techniques and equivalents, and have finite coverage. The email layout witness, for example, describes a centered table layout in its frontmatter and a centered single-column layout in its workflow. A correctly rendered table can use different syntax from the registered patterns. Slide layouts and catalog CSVs may also reuse those terms. The frozen Word structured result also credits a web viewport inspection method for its render requirement. Its two Word-native layout methods have direct DOCX evidence. Use the reported source ranges to check whether the technique actually operates on the requested artifact.

The optional `--strict` gate requires all registered needs and expected gaps, complete coverage of deep workflows and zero named near misses in all measured modes. An evaluation score measures plan quality. Actual execution, artifact validation and user-task outcomes require separate observations.
