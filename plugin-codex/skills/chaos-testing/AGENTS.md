# Chaos Testing — Chaos Engineering

## Overview

**Chaos engineering** is the discipline of experimenting on a system in order to build confidence in its ability to withstand turbulent conditions in production. Rather than waiting for failures to happen organically, chaos engineering proactively injects controlled faults to expose weaknesses before they cause real outages.

The discipline is grounded in the **Principles of Chaos Engineering** ([principlesofchaos.org](https://principlesofchaos.org)), which define a rigorous, experiment-driven approach to resilience validation:

1. **Build a hypothesis around steady-state behavior** — define what "normal" looks like in terms of measurable business and system metrics.
2. **Vary real-world events** — inject failures that reflect things that actually happen (server crashes, network partitions, disk full, dependency latency).
3. **Run experiments in production** — only production truly represents production; start with staging, but graduate to production with safeguards.
4. **Automate experiments to run continuously** — one-off experiments decay in value; automated chaos builds lasting confidence.
5. **Minimize blast radius** — design experiments to limit the scope of impact, and have abort conditions ready.

> **Reference**: *Chaos Engineering: System Resiliency in Practice* by Casey Rosenthal & Nora Jones (O'Reilly, 2020) is the definitive guide to the discipline, covering theory, organizational adoption, and practical tooling.

## Core Concepts

### Steady-State Hypothesis

The steady-state hypothesis defines the normal, expected behavior of a system using measurable indicators. Before injecting any fault, you must establish what "healthy" looks like so you can detect whether the system degrades. Common steady-state metrics include request success rate, p99 latency, error rate, throughput, and queue depth. The experiment succeeds if the steady state is maintained despite the injected fault, and fails if the steady state is violated.

**Example hypothesis**: "When 10% of application pods are terminated, the request success rate remains above 99.5% and p99 latency stays below 500ms."

### Blast Radius

Blast radius is the scope of potential impact of a chaos experiment. Effective chaos engineering demands strict control over blast radius — start with the smallest possible scope and expand gradually as confidence grows. Controls include targeting a single host or pod, limiting experiment duration, defining automatic abort conditions (kill switches), running during low-traffic windows initially, and having rollback procedures ready.

### Game Days

A game day is a planned event where teams deliberately inject failures into their systems and practice their incident response. Game days combine chaos experiments with organizational readiness: teams observe how monitoring, alerting, runbooks, and human coordination perform under realistic failure conditions. They are typically scheduled, cross-team events with defined scenarios, observers, and post-game retrospectives.

## Experiment Design

Designing a chaos experiment follows a structured methodology:

1. **Form a hypothesis** — Identify a steady-state metric and predict what will happen when a specific failure is introduced. Example: "If we lose one availability zone, failover completes within 30 seconds and no requests are dropped."

2. **Design the smallest experiment** — Start with the most limited scope that can test the hypothesis. Target a single instance, a single failure mode, and a short duration. Prefer reversible faults (latency injection) over destructive ones (data corruption) when starting out.

3. **Measure and observe** — Instrument the experiment to capture the steady-state metrics, system behavior, and any side effects. Use existing monitoring dashboards and add experiment-specific telemetry if needed.

4. **Analyze results** — Compare observed behavior against the hypothesis. If the steady state held, confidence increases. If it broke, you have discovered a real weakness to fix. Document findings either way.

5. **Expand or remediate** — If the experiment passed, gradually increase blast radius (more instances, longer duration, production environment). If it failed, fix the underlying issue and re-run the experiment.

## Common Failure Modes

| Category | Failure Mode | Example |
|----------|-------------|---------|
| **Infrastructure** | Host/VM termination | Kill a random EC2 instance or Kubernetes pod |
| **Infrastructure** | Availability zone loss | Simulate an entire AZ going offline |
| **Infrastructure** | Disk full | Fill disk to capacity on a target host |
| **Network** | Latency injection | Add 500ms latency to inter-service calls |
| **Network** | Packet loss | Drop 10% of packets between services |
| **Network** | DNS failure | Return NXDOMAIN for a downstream dependency |
| **Network** | Partition | Block traffic between two service groups |
| **Application** | CPU stress | Consume 90% CPU on target hosts |
| **Application** | Memory pressure | Allocate memory until near OOM conditions |
| **Application** | Process crash | Kill the application process (not the host) |
| **Application** | Exception injection | Force specific error codes from a dependency |
| **Data** | Clock skew | Shift system time forward or backward |
| **Data** | Data corruption simulation | Return malformed responses from a dependency |
| **Dependencies** | Third-party outage | Block or timeout calls to external APIs |
| **Dependencies** | Database failover | Trigger a primary-to-replica failover |
| **Dependencies** | Cache eviction | Flush the entire cache (cold cache scenario) |

## Cross-Platform Tools

| Tool | Type | Target | Open Source |
|------|------|--------|-------------|
| **Chaos Monkey** | Random instance termination | AWS VMs / cloud instances | Yes (Netflix OSS) |
| **Gremlin** | SaaS fault injection platform | Multi-layer (infra, network, app, data) | No (commercial SaaS, free tier available) |
| **Litmus** | Kubernetes chaos engineering (CNCF) | Kubernetes + cloud resources + bare metal | Yes (Apache 2.0) |
| **Chaos Mesh** | Kubernetes fault injection | Kubernetes only | Yes (Apache 2.0) |
| **Toxiproxy** | TCP proxy for simulating network faults | Any TCP connection (language-agnostic) | Yes (Shopify, MIT) |
| **Pumba** | Container chaos tool | Docker containers | Yes (Apache 2.0) |
| **AWS Fault Injection Service** | Managed chaos engineering service | AWS resources (EC2, ECS, EKS, RDS, etc.) | No (AWS managed service) |
| **Azure Chaos Studio** | Managed chaos engineering service | Azure resources (VMs, AKS, Cosmos DB, etc.) | No (Azure managed service) |

## Toxiproxy

[Toxiproxy](https://github.com/Shopify/toxiproxy) is a TCP proxy that lets you simulate network conditions between your application and its dependencies. It is lightweight, language-agnostic, and ideal for integration testing environments.

### Setup

```bash
# Install Toxiproxy server
# macOS
brew install toxiproxy

# Linux (download binary)
wget https://github.com/Shopify/toxiproxy/releases/download/v2.9.0/toxiproxy-server-linux-amd64 -O toxiproxy-server
chmod +x toxiproxy-server

# Start the Toxiproxy server
./toxiproxy-server &

# Create a proxy for your database connection
toxiproxy-cli create postgres_primary \
  --listen 127.0.0.1:15432 \
  --upstream 127.0.0.1:5432
```

### Usage Example — Simulating Latency

```bash
# Add 500ms latency to all traffic through the proxy
toxiproxy-cli toxic add postgres_primary \
  --type latency \
  --attribute latency=500 \
  --attribute jitter=100

# Your application connects to 127.0.0.1:15432 instead of :5432
# and now experiences ~500ms (+/- 100ms) of added latency on every query

# Remove the toxic when done
toxiproxy-cli toxic remove postgres_primary --toxicName latency_downstream
```

### Usage Example — Simulating Connection Loss

```bash
# Simulate a complete connection timeout (no data flows)
toxiproxy-cli toxic add postgres_primary \
  --type timeout \
  --attribute timeout=0

# Simulate the downstream service going away entirely
toxiproxy-cli toggle postgres_primary  # disables the proxy

# Re-enable
toxiproxy-cli toggle postgres_primary  # toggles it back on
```

### Programmatic Usage (Node.js)

```javascript
import Toxiproxy from "toxiproxy-node-client";

const toxiproxy = new Toxiproxy("http://localhost:8474");

// Create a proxy
const proxy = await toxiproxy.createProxy({
  name: "redis_cache",
  listen: "127.0.0.1:16379",
  upstream: "127.0.0.1:6379",
});

// Add latency toxic
await proxy.addToxic({
  type: "latency",
  attributes: { latency: 200, jitter: 50 },
});

// Run your test against 127.0.0.1:16379
// ...

// Clean up
await proxy.remove();
```

## Litmus

[Litmus](https://litmuschaos.io/) is a CNCF project that provides a complete chaos engineering platform for Kubernetes. It uses a CRD-based approach where chaos experiments are defined as Kubernetes resources.

### ChaosEngine Example — Pod Delete Experiment

```yaml
apiVersion: litmuschaos.io/v1alpha1
kind: ChaosEngine
metadata:
  name: pod-delete-chaos
  namespace: default
spec:
  appinfo:
    appns: "default"
    applabel: "app=my-api"
    appkind: "deployment"
  engineState: "active"
  chaosServiceAccount: litmus-admin
  experiments:
    - name: pod-delete
      spec:
        components:
          env:
            # Number of pods to target
            - name: TOTAL_CHAOS_DURATION
              value: "30"
            # Duration between each pod deletion (seconds)
            - name: CHAOS_INTERVAL
              value: "10"
            # Forceful deletion (like kill -9)
            - name: FORCE
              value: "false"
            # Number of pods to delete at a time
            - name: PODS_AFFECTED_PERC
              value: "50"
        probe:
          - name: "check-api-health"
            type: "httpProbe"
            httpProbe/inputs:
              url: "http://my-api.default.svc:8080/healthz"
              method:
                get:
                  criteria: "=="
                  responseCode: "200"
            mode: "Continuous"
            runProperties:
              probeTimeout: 5s
              interval: 2s
              retry: 3
```

This experiment deletes 50% of pods matching `app=my-api` every 10 seconds for 30 seconds, while continuously probing the health endpoint to verify the steady-state hypothesis (the API remains responsive).

## Chaos Mesh

[Chaos Mesh](https://chaos-mesh.org/) is a Kubernetes-native chaos engineering platform that supports network, pod, stress, I/O, time, and DNS fault injection through CRDs.

### NetworkChaos Example — Latency and Packet Loss

```yaml
apiVersion: chaos-mesh.org/v1alpha1
kind: NetworkChaos
metadata:
  name: network-delay-payment-service
  namespace: chaos-testing
spec:
  action: delay
  mode: all
  selector:
    namespaces:
      - production
    labelSelectors:
      app: payment-service
  delay:
    latency: "200ms"
    correlation: "75"
    jitter: "50ms"
  direction: to
  target:
    selector:
      namespaces:
        - production
      labelSelectors:
        app: order-service
    mode: all
  duration: "5m"
  scheduler:
    cron: "@every 24h"
```

This experiment injects 200ms of latency (with 50ms jitter and 75% correlation) on all traffic from `payment-service` to `order-service` for 5 minutes, running daily on a cron schedule.

### NetworkChaos Example — Partition

```yaml
apiVersion: chaos-mesh.org/v1alpha1
kind: NetworkChaos
metadata:
  name: network-partition-db
  namespace: chaos-testing
spec:
  action: partition
  mode: all
  selector:
    namespaces:
      - production
    labelSelectors:
      app: api-gateway
  direction: to
  target:
    selector:
      namespaces:
        - production
      labelSelectors:
        app: database-proxy
    mode: all
  duration: "2m"
```

## CI/CD Integration

Chaos experiments should be integrated into your CI/CD pipeline systematically, with increasing scope at each stage:

1. **After integration tests pass** — Chaos tests run only when the system is known to be functionally correct. Running chaos against a broken build wastes time.

2. **Start with staging environments** — Run automated chaos experiments in staging on every deployment to catch resilience regressions early. Use Litmus, Chaos Mesh, or Toxiproxy experiments that are defined as code alongside your application.

3. **Gate production deployments** — Include chaos experiment results as a deployment gate. If the service cannot survive pod termination or latency injection in staging, it should not be promoted to production.

4. **Expand to production with safeguards** — Once experiments consistently pass in staging, graduate them to production with strict controls:
   - Run during low-traffic windows initially
   - Use canary or blue-green deployments to limit blast radius
   - Define automatic abort conditions (if error rate exceeds threshold, halt the experiment)
   - Have a human-in-the-loop approval for destructive experiments

5. **Schedule recurring experiments** — Use cron-based scheduling (Chaos Mesh scheduler, CI cron jobs) to run chaos experiments continuously, not just at deploy time.

```
CI/CD Pipeline with Chaos Integration
  │
  ├── Build + Static Analysis
  │
  ├── Unit Tests
  │
  ├── Integration Tests (with Toxiproxy for dependency faults)
  │
  ├── Deploy to Staging
  │
  ├── Chaos Tests in Staging ──► Litmus pod-delete, NetworkChaos experiments
  │     │
  │     └── Abort deployment if steady-state hypothesis violated
  │
  ├── Deploy to Production (canary)
  │
  └── Chaos Tests in Production ──► Scheduled experiments with kill switches
```

## Game Day Checklist

1. **Define objectives** — Select the systems to test and the failure scenarios to simulate. Align on what you want to learn, not what you want to prove.

2. **Establish steady-state metrics** — Agree on the measurable indicators of system health (error rate, latency, throughput) and the thresholds that constitute a failure.

3. **Prepare the experiments** — Write or configure the chaos experiments, test them in a non-production environment, and verify they can be aborted quickly.

4. **Brief all participants** — Ensure all teams involved (SREs, developers, on-call, support) understand the plan, timeline, and their roles. Designate an experiment lead, observers, and a safety controller.

5. **Verify monitoring and alerting** — Confirm that dashboards are up, alerting is active, and on-call rotations are staffed. You must be able to see the impact in real time.

6. **Confirm rollback procedures** — Validate that you can abort the experiment and restore normal operations within a defined time window. Test the kill switch before starting.

7. **Execute the experiment** — Run the chaos injection according to plan. Observers record system behavior, alerting behavior, and team response actions.

8. **Monitor and respond** — The team responds to the injected failure as if it were a real incident. This exercises runbooks, communication channels, and escalation procedures.

9. **Abort if necessary** — If the blast radius exceeds expectations or real user impact is detected beyond acceptable thresholds, trigger the abort condition immediately.

10. **Conduct a retrospective** — After the game day, hold a blameless retrospective. Document what went well, what surprised the team, what broke, and what actions to take. Create tickets for discovered weaknesses.

11. **Track and remediate** — File issues for every weakness discovered. Re-run the experiment after fixes are applied to verify improvement.

## Best Practices

- **Start small and expand gradually** — Begin with a single pod or instance in a non-production environment. Increase blast radius only as confidence grows.
- **Automate experiments** — Manual chaos is a one-time event; automated chaos is continuous confidence. Define experiments as code and run them in CI/CD.
- **Test in production, but with safeguards** — Production is the only environment that truly represents production. Use kill switches, canary scopes, and low-traffic windows to manage risk.
- **Establish steady state before injecting faults** — Never run a chaos experiment against a system whose baseline health is unknown or degraded.
- **Define abort conditions upfront** — Every experiment must have clear criteria for when to stop (error rate threshold, latency ceiling, manual kill switch).
- **Make experiments observable** — Tag chaos experiments in your monitoring system so you can correlate injected faults with observed impact. Use experiment IDs in logs and dashboards.
- **Run game days regularly** — Schedule game days quarterly at minimum. They exercise not just systems but also people and processes.
- **Involve the whole team** — Chaos engineering is not just for SREs. Developers, product managers, and support teams all benefit from understanding failure modes.
- **Fix what you find** — The value of chaos engineering is in the remediation, not the experiment itself. Every discovered weakness should result in an actionable improvement.
- **Document and share learnings** — Publish experiment results, both successes and failures, so the broader organization builds institutional knowledge about system resilience.
