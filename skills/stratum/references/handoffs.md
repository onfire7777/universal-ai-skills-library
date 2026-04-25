# Handoff Templates

Purpose: Inter-agent handoff templates.
Read when: You need to confirm data exchange formats with other agents.

## Contents

- Inbound Handoffs (Stratum receives)
- Outbound Handoffs (Stratum sends)
- Handoff Examples

---

## Inbound Handoffs

### ATLAS_TO_STRATUM_HANDOFF

When Atlas passes architecture analysis results to Stratum.

```yaml
ATLAS_TO_STRATUM_HANDOFF:
  source: Atlas
  context: "Architecture analysis results for C4 modeling"
  payload:
    dependency_map:
      modules:
        - name: "[module name]"
          path: "[module path]"
          dependencies: ["[dep1]", "[dep2]"]
      external_dependencies:
        - name: "[external service]"
          type: "API" | "Database" | "MessageQueue" | "FileStorage"
          url: "[connection URL pattern]"
    module_boundaries:
      - boundary: "[boundary name]"
        modules: ["[mod1]", "[mod2]"]
        coupling_score: 0.0-1.0
    tech_stack:
      languages: ["TypeScript", "Python"]
      frameworks: ["Express", "React"]
      databases: ["PostgreSQL", "Redis"]
      infrastructure: ["Docker", "AWS ECS"]
    coupling_metrics:
      afferent: { "[module]": 5 }
      efferent: { "[module]": 3 }
    issues:
      - type: "circular_dependency" | "god_class" | "high_coupling"
        location: "[path]"
        severity: "high" | "medium" | "low"
```

### LENS_TO_STRATUM_HANDOFF

When Lens passes codebase investigation results to Stratum.

```yaml
LENS_TO_STRATUM_HANDOFF:
  source: Lens
  context: "Codebase structure investigation for C4 modeling"
  payload:
    entry_points:
      - path: "[main entry file]"
        type: "web_server" | "cli" | "worker" | "lambda"
        port: "[port if applicable]"
    data_flows:
      - name: "[flow name]"
        steps:
          - from: "[source module]"
            to: "[target module]"
            method: "[function/API call]"
    file_structure:
      root: "[project root]"
      key_directories:
        - path: "[dir]"
          purpose: "[description]"
    deployment_configs:
      - type: "Dockerfile" | "docker-compose" | "k8s" | "terraform"
        path: "[config file path]"
        services: ["[service names]"]
```

### RIPPLE_TO_STRATUM_HANDOFF

When Ripple passes change impact analysis results to Stratum (model update trigger).

```yaml
RIPPLE_TO_STRATUM_HANDOFF:
  source: Ripple
  context: "Change impact analysis triggering C4 model update"
  payload:
    change_description: "[what changed]"
    affected_areas:
      - area: "[module/service name]"
        impact: "high" | "medium" | "low"
        c4_level_affected: [1, 2, 3]
    new_dependencies:
      - from: "[source]"
        to: "[target]"
        type: "API" | "Database" | "Event"
    removed_dependencies:
      - from: "[source]"
        to: "[target]"
    structural_changes:
      - type: "new_container" | "removed_container" | "container_split" | "container_merge"
        detail: "[description]"
```

### SCAFFOLD_TO_STRATUM_HANDOFF

When Scaffold passes infrastructure topology to Stratum (for Deployment diagrams).

```yaml
SCAFFOLD_TO_STRATUM_HANDOFF:
  source: Scaffold
  context: "Infrastructure topology for C4 deployment diagram"
  payload:
    environments:
      - name: "Production" | "Staging" | "Development"
        provider: "AWS" | "GCP" | "Azure" | "On-premises"
        nodes:
          - name: "[node name]"
            technology: "[e.g., ECS, GKE, EC2]"
            instances: 1
            containers_deployed: ["[container names from C4 model]"]
            children:
              - name: "[child node]"
                technology: "[technology]"
        infrastructure_nodes:
          - name: "[e.g., Load Balancer, CDN]"
            technology: "[e.g., ALB, CloudFront]"
```

---

## Outbound Handoffs

### STRATUM_TO_CANVAS_HANDOFF

When Stratum requests Canvas to render C4 diagrams.

```yaml
STRATUM_TO_CANVAS_HANDOFF:
  source: Stratum
  target: Canvas
  context: "Render C4 model diagrams"
  payload:
    model_type: "c4"
    levels_to_render:
      - level: 1
        view_key: "L1-SystemContext"
        title: "[diagram title]"
      - level: 2
        view_key: "L2-Containers"
        title: "[diagram title]"
    structurizr_dsl: |
      [complete DSL code]
    render_format: "mermaid" | "draw.io" | "plantuml"
    style_hints:
      color_scheme: "standard_c4"
      layout_direction: "lr" | "tb"
    additional_views:
      - type: "dynamic" | "deployment"
        view_key: "[view key]"
        title: "[title]"
```

### STRATUM_TO_SCRIBE_HANDOFF

When Stratum requests Scribe to document the C4 model.

```yaml
STRATUM_TO_SCRIBE_HANDOFF:
  source: Stratum
  target: Scribe
  context: "Embed C4 model in HLD/LLD document"
  payload:
    document_type: "HLD" | "LLD"
    section: "architecture"
    c4_model:
      system_name: "[name]"
      system_description: "[description]"
      context_summary: "[L1 narrative description]"
      containers:
        - name: "[container name]"
          technology: "[tech]"
          responsibility: "[description]"
      key_relationships:
        - "[relationship description]"
      architectural_decisions:
        - decision: "[what was decided]"
          rationale: "[why]"
          alternatives: ["[alt1]", "[alt2]"]
    structurizr_dsl: |
      [complete DSL for reference]
    diagrams_available:
      - level: 1
        format: "mermaid"
        code: "[mermaid code]"
      - level: 2
        format: "mermaid"
        code: "[mermaid code]"
```

### STRATUM_TO_ATLAS_HANDOFF

When Stratum requests Atlas to create ADRs for architecture decisions.

```yaml
STRATUM_TO_ATLAS_HANDOFF:
  source: Stratum
  target: Atlas
  context: "Architecture decisions from C4 modeling for ADR creation"
  payload:
    decisions:
      - title: "[decision title]"
        context: "[why this decision was needed]"
        decision: "[what was decided]"
        consequences:
          - "[consequence 1]"
          - "[consequence 2]"
        alternatives_considered:
          - name: "[alternative]"
            pros: ["[pro]"]
            cons: ["[con]"]
    system_name: "[C4 system name]"
    affected_containers: ["[container names]"]
```

---

## Handoff Example: Full Model Build Chain

Typical full model build chain:

```
Step 1: User → Stratum
  "Create a C4 model for this project"

Step 2: Stratum (DISCOVER phase)
  - Directly analyze the codebase
  - Request dependency map from Atlas if needed

Step 3: Atlas → Stratum (ATLAS_TO_STRATUM_HANDOFF)
  - Dependency maps, module boundaries, tech stack info

Step 4: Stratum (MODEL + VERIFY phases)
  - Build C4 model and verify consistency
  - Generate Structurizr DSL

Step 5: Stratum → Canvas (STRATUM_TO_CANVAS_HANDOFF)
  - DSL code + rendering instructions

Step 6: Canvas
  - Generate diagrams in Mermaid/draw.io

Step 7: (Optional) Stratum → Scribe (STRATUM_TO_SCRIBE_HANDOFF)
  - Embed C4 model section into HLD document
```
