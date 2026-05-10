# OpenTelemetry Best Practices & Distributed Tracing

> Instrumentation strategy, semantic conventions, collector pipeline, sampling, tracing, telemetry correlation, GenAI observability, cost optimization

---

## 1. Instrumentation Strategy

| # | Practice | Description | Importance |
|---|---------|-------------|------------|
| **OT-01** | **Auto-first, Manual-second** | Start with auto-instrumentation for immediate visibility, then add manual spans for business-critical paths | Required |
| **OT-02** | **SDK initialization first** | Initialize OTel SDK before any application module imports | Required |
| **OT-03** | **Always close spans** | Use async/await + finally blocks to ensure spans are closed | Required |
| **OT-04** | **Add business attributes** | Enrich auto-generated spans with business context (customer_tier, order_value) | Recommended |
| **OT-05** | **Record errors and events** | Log state transitions and errors as span events | Recommended |

```
Initialization order (critical):
  // 1. Initialize OTel SDK first
  const { NodeTracerProvider } = require('@opentelemetry/sdk-trace-node');
  const provider = new NodeTracerProvider();
  provider.register();

  // 2. Then import application modules
  const express = require('express');
```

### Manual Span Creation

```python
@tracer.start_as_current_span("process_payment")
def process_payment(order_id: str, amount: float):
    span = trace.get_current_span()
    span.set_attribute("order.id", order_id)
    span.set_attribute("payment.amount", amount)
    span.set_attribute("payment.currency", "USD")

    try:
        result = charge_card(amount)
        span.set_attribute("payment.status", "success")
        span.set_status(StatusCode.OK)
        return result
    except PaymentError as e:
        span.set_status(StatusCode.ERROR, str(e))
        span.record_exception(e)
        raise
```

---

## 2. Semantic Conventions

### Standard Attributes (must follow)

```
HTTP:       http.method, http.status_code, http.route
DB:         db.system, db.statement, db.name
Messaging:  messaging.system, messaging.destination
RPC:        rpc.system, rpc.service, rpc.method
```

### Application-Specific Attributes

```
- Prefix: app.* for custom attributes
- Naming: snake_case consistently
- Examples: app.customer_tier, app.order_value, app.feature_flag
```

### Attribute Anti-Patterns

```
x  High-cardinality attributes (user_id) on all spans -> use traces not metrics
x  Duplicate attributes on parent/child spans
x  Mixing metric/log data into span attributes
x  Abbreviations (svc -> service)
```

---

## 3. Span Naming Conventions

| Layer | Format | Examples |
|-------|--------|---------|
| **HTTP server** | `HTTP {METHOD} {route}` | `HTTP GET /api/users/:id` |
| **HTTP client** | `HTTP {METHOD}` | `HTTP POST` |
| **Database** | `{db.system} {operation} {table}` | `postgresql SELECT orders` |
| **Message publish** | `{queue} publish` | `orders.created publish` |
| **Message consume** | `{queue} process` | `orders.created process` |
| **Business logic** | `{verb}_{noun}` | `validate_payment`, `calculate_tax` |
| **External service** | `{service}.{operation}` | `stripe.create_charge` |

---

## 4. Context Propagation

| Format | Header | Ecosystem |
|--------|--------|-----------|
| **W3C Trace Context** (recommended) | `traceparent`, `tracestate` | Standard |
| **B3** | `X-B3-TraceId`, `X-B3-SpanId` | Zipkin |
| **Jaeger** | `uber-trace-id` | Jaeger |

### Propagation Checklist

```markdown
- [ ] All HTTP clients inject trace context headers
- [ ] All HTTP servers extract trace context headers
- [ ] Message queues propagate trace context in headers/metadata
- [ ] Async workers link to parent span via context
- [ ] Batch jobs create new root spans with links to triggers
- [ ] Third-party API calls create client spans
```

### Baggage

- Use Baggage for cross-service context (customer_id, tenant_id)
- Caution: Baggage propagates to all downstream services
- Never include sensitive information in Baggage

---

## 5. Collector Deployment Patterns

| Pattern | Configuration | Pros | Cons | Scale |
|---------|--------------|------|------|-------|
| **Agent** | Sidecar per app | Network minimal, app isolation | Config management distributed | Small |
| **Gateway** | Central server | Centralized config, routing | SPOF risk | Medium |
| **Hierarchical** | Agent + Gateway | Optimal reliability/management balance | Complexity | Large (recommended) |

### Collector Configuration

```yaml
receivers:
  otlp:
    protocols:
      grpc:
        endpoint: 0.0.0.0:4317

processors:
  memory_limiter:        # OOM prevention (MUST be first)
    check_interval: 1s
    limit_mib: 1000
  batch:                 # Network efficiency
    send_batch_size: 10000
    timeout: 10s

exporters:
  otlp:
    endpoint: observability-backend:4317

service:
  pipelines:
    traces:
      receivers: [otlp]
      processors: [memory_limiter, batch]  # memory_limiter always first
      exporters: [otlp]
```

### Processor Ordering (Critical)

1. **memory_limiter** (prevent crashes)
2. **enrichment** (k8sattributes, resource)
3. **filter/transform** (PII redaction, filtering)
4. **batch** (efficient delivery, always last)

### PII/PHI Filtering

```yaml
processors:
  filter:
    spans:
      include:
        match_type: regexp
        attributes:
          - key: db.statement
            value: "(?i)(?:password|passwd)\\s*=\\s*[^\\s,;]+"
      actions:
        - key: db.statement
          action: update
          value: "REDACTED"
```

### Operational Requirements

- Minimum 4GB node memory for graceful shutdown
- Version-lock Operator, Collector, and Target Allocator together
- Use `nodeAffinity` to prevent deployment on small nodes
- Monitor: `otelcol_receiver_refused_metric_points_total` (non-zero = data loss)

---

## 6. Sampling Strategies

| Strategy | Decision Point | Pros | Cons | Use |
|----------|---------------|------|------|-----|
| **Head Sampling** | Trace start | Simple, low overhead | Misses error traces | Dev environments |
| **Tail Sampling** | Trace completion | Intelligent decisions | Requires buffering | Production (recommended) |
| **Probabilistic** | Random % | Predictable cost | Error miss risk | High traffic |
| **Rate Limiting** | Time-based cap | Spike control | Important trace loss risk | Burst protection |

### Recommended: Composite Sampling Strategy

```yaml
processors:
  tail_sampling:
    decision_wait: 10s
    policies:
      - name: errors
        type: status_code
        status_code: { status_codes: [ERROR] }     # 100% error retention
      - name: slow-requests
        type: latency
        latency: { threshold_ms: 2000 }
      - name: critical-endpoints
        type: string_attribute
        string_attribute:
          key: http.route
          values: ["/api/payments", "/api/auth"]
      - name: baseline
        type: probabilistic
        probabilistic: { sampling_percentage: 5 }   # 5% normal traffic
    decision_cache_size: 50000
```

### Metrics Accuracy Preservation

- Generate metrics BEFORE sampling (spanmetrics processor)
- Use `spanmetrics` processor for automatic RED metric generation
- Use `servicegraph` processor for automatic dependency map generation

```yaml
processors:
  spanmetrics:
    metrics_exporter: prometheus
    dimensions:
      - name: service.name
      - name: http.method
      - name: http.status_code
```

---

## 7. Telemetry Correlation (Three Pillars)

```
Log-Trace correlation:
  - Inject trace_id / span_id into logs automatically
  - Use structured logging (JSON)
  - ERROR/WARN logs must also be recorded as span events

Trace -> Metrics conversion:
  - spanmetrics processor for RED metrics
  - servicegraph processor for dependency maps

Performance tuning:
  BatchSpanProcessor:
    maxQueueSize: 2048
    maxExportBatchSize: 512
    scheduledDelayMillis: 5000
    exportTimeoutMillis: 30000
  - Enable gzip compression (bandwidth reduction)
  - Circuit breaker (telemetry must not affect availability)
```

---

## 8. Trace Analysis Patterns

| Pattern | What to Look For | Action |
|---------|-----------------|--------|
| **Long spans** | Single span > SLO threshold | Optimize or decompose |
| **Wide traces** | Fan-out > 50 spans | Check N+1 queries |
| **Deep traces** | Depth > 10 levels | Simplify call chain |
| **Orphan spans** | Missing parent spans | Fix context propagation |
| **Gap spans** | Time gaps between child spans | Check queuing/scheduling |

---

## 9. Cardinality Management

```
Cardinality explosion example:
  http_requests_total{method, path, status, user_id, client_ip}
  -> method(5) x path(100) x status(10) x user_id(100K) x client_ip(50K)
  -> 2.5 trillion unique time series -> system collapse

Detection:
  1. Index size spikes (RAM/disk monitoring)
  2. remote_write ingestion delays
  3. Aggregation query (sum, avg) latency degradation
  4. Observability platform cost spikes

Control strategy:
  Tier 1: Per-service cardinality limits
    - Business-critical metrics: higher thresholds
    - Infrastructure metrics: strict limits

  Tier 2: Tiered retention policies
    - High resolution (raw data): 24-48 hours
    - Medium resolution (1min aggregation): 30 days
    - Low resolution (1hr aggregation): 13+ months

  Tier 3: Adaptive downsampling
    - Keep high-fidelity data at edge (local)
    - Selective downsampling at central aggregation
    - Prefer automation over manual recording rules
```

---

## 10. Cost Optimization

```
Key cost levers (most to least effective):
  1. Pipeline-level filtering BEFORE storage (Collector processors)
  2. Intelligent sampling (tail-based, composite)
  3. Tiered retention (raw -> aggregated -> archived)
  4. Per-node Target Allocator (prevents 20-40x metric duplication)
  5. Cardinality limits per service

Results benchmark (CNCF case study):
  - 72% cost reduction vs previous vendor
  - 100% APM trace coverage (was 5% sampling)
  - Enabled by: OTel Collector + open-source backends (Loki, Tempo, Mimir)

Observability budget framework:
  - Set per-team telemetry budget (GB/day or cost/month)
  - Monitor telemetry volume per service
  - Alert on budget overruns
  - Quarterly review: optimize top-5 cost contributors

Tool sprawl prevention:
  - Standardize on OTel as the single collection layer
  - Consolidate to single backend per signal type
  - Avoid Prometheus + Datadog + New Relic + custom tools in parallel
```

---

## 11. GenAI / Agent Observability

```
OTel GenAI Semantic Conventions (v1.37+):
  Standard attributes:
    gen_ai.request.model         # Model identifier
    gen_ai.usage.input_tokens    # Token usage
    gen_ai.usage.output_tokens
    gen_ai.provider.name         # Provider (openai, anthropic, etc.)
    gen_ai.request.temperature   # Model parameters
    gen_ai.request.max_tokens

  Agent-specific:
    gen_ai.agent.name            # Agent identifier
    gen_ai.agent.description
    gen_ai.tool.name             # Tool calls
    gen_ai.tool.description

Instrumentation approaches:
  Option 1: Baked-in (framework embeds OTel)
    + Simplified adoption, feature-release control
    - Framework bloat, version lag

  Option 2: External OTel libraries (recommended)
    + Decoupled, community-maintained
    - Fragmentation risk if incompatible packages

Key metrics for LLM observability:
  - Token usage per request/session
  - Latency per model call (time to first token, total)
  - Error rates by model/provider
  - Cost per request (tokens x price)
  - Tool call success/failure rates
```

---

## 12. Beacon Integration

```
Usage by mode:
  1. DESIGN: OTel instrumentation strategy, collector pipeline design
  2. SPECIFY: Collector pipeline specs, sampling configuration
  3. MEASURE: Sampling strategy optimization, cardinality monitoring
  4. Periodic review: Semantic Conventions compliance, cost optimization

Quality gates:
  - OTel SDK initialization is first in app startup (OT-02)
  - memory_limiter processor is first in pipeline
  - Error traces retained at 100%
  - Logs inject trace_id (correlation enabled)
  - PII/PHI filtering in Collector
  - Semantic Conventions compliance in attribute naming
  - New metric addition requires cardinality estimate
  - Telemetry budget per team/service defined
```

**Source:** [OTel Semantic Conventions v1.40](https://opentelemetry.io/docs/specs/semconv/) · [Better Stack: OTel Best Practices](https://betterstack.com/community/guides/observability/opentelemetry-best-practices/) · [CNCF: Cost-Effective OTel](https://www.cncf.io/blog/2025/12/16/how-to-build-a-cost-effective-observability-platform-with-opentelemetry/) · [Dash0: OTel Collector Guide](https://dash0.com/blog/opentelemetry-collector-guide) · [OTel: AI Agent Observability](https://opentelemetry.io/blog/2025/ai-agent-observability/) · [OTel GenAI SemConv](https://opentelemetry.io/docs/specs/semconv/gen-ai/) · [OTel Weaver](https://opentelemetry.io/blog/2025/otel-weaver/)
