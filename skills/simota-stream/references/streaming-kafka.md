# Streaming Architecture — Kafka Design

Purpose: topic, consumer, schema, and delivery guidance for event-driven and Kafka-based pipelines.

## Contents

1. Topic design
2. Consumer-group rules
3. Event schema
4. Delivery guarantees
5. Outbox and processing patterns

## Topic Design

```yaml
topic_naming_convention:
  pattern: "{domain}.{entity}.{event}"
  examples:
    - "orders.order.created"
    - "users.profile.updated"
    - "payments.transaction.completed"
```

### Topic Configuration Rules

| Setting | Rule |
|---------|------|
| Partitions | `10x` expected peak throughput in MB/s |
| Minimum partitions | `3` |
| Practical maximum | `100` per topic |
| Production retention | `7d` by default, extend for replay needs |
| Compacted retention | effectively infinite for state topics |
| Replication factor | prod `3`, staging `2`, dev `1` |

Anti-pattern guards:
- do not create a god topic
- do not use `1` partition for scalable workloads
- do not jump to `1000` partitions without a hard need

## Consumer Group Rules

```yaml
consumer_groups:
  naming: "{service}-{purpose}"
  examples:
    - "analytics-aggregator"
    - "notification-sender"
    - "audit-logger"
```

Rules:
- `consumers <= partitions`
- scale off consumer lag, not guesswork
- separate groups by purpose, not convenience

## Event Schema

Required metadata:
- `event_id`
- `event_type`
- `timestamp`
- `version`
- `source`
- `correlation_id`
- `data`

Compatibility rules:
- allow additive fields
- do not remove fields or change types without versioning
- use Schema Registry for shared or critical streams

Recommended event size: `1KB-10KB`. If payloads grow beyond that, prefer reference patterns.

## Delivery Guarantees

| Level | Use When | Notes |
|-------|----------|-------|
| At-most-once | loss is acceptable | avoid for business-critical data |
| At-least-once | default | combine with idempotent consumers |
| Exactly-once | narrow Kafka-to-Kafka cases | external systems still need idempotent sinks |

Prefer "effectively once": `acks=all` + manual commit after processing + idempotent sink or deduplication.

## Outbox And Processing Patterns

Use Outbox when DB writes and event publication must remain atomic:
1. write business data and outbox record in one DB transaction
2. capture outbox changes with CDC
3. publish to Kafka from the CDC pipeline

Common processing shapes:
- stateless transform
- windowed aggregation
- stream-table join

Operational guards:
- use manual commit after successful processing
- send poison pills to a DLT/DLQ
- track consumer lag, throughput, error rate, rebalance frequency, disk usage, and under-replicated partitions
