# ORM Performance Pitfalls

Purpose: Use this file when the bottleneck sits in ORM-generated SQL, hydration, serialization, or ORM choice.

Contents:
- `OP-01..08`
- ORM comparison
- raw-SQL switch criteria
- review checklist

## ORM Pitfalls

| ID | Pitfall | Typical ORM | Response |
|----|---------|-------------|----------|
| `OP-01` | object hydration cost | TypeORM, Sequelize | use raw SQL or lighter mapping for large result sets |
| `OP-02` | cold-start delay | Prisma | measure serverless overhead, often around `300ms` |
| `OP-03` | serialization overhead | Prisma | re-measure on high-throughput endpoints |
| `OP-04` | lazy-loading trap | all ORMs | use explicit eager loading |
| `OP-05` | production `synchronize: true` | TypeORM | forbid in production |
| `OP-06` | over-fetching | all ORMs | use explicit `select` |
| `OP-07` | over-wide transaction scope | all ORMs | shrink transaction boundaries |
| `OP-08` | false type-safety assumptions | TypeORM and similar | validate runtime behavior |

## ORM Comparison

| ORM | Strength | Risk |
|-----|----------|------|
| Drizzle | lightweight, predictable SQL, good for serverless | fewer guardrails against N+1 |
| Prisma | excellent DX and type safety | `~300ms` cold start, serialization cost, SQL may need review |
| TypeORM | flexible patterns | hydration cost, weaker type guarantees, lower maintenance confidence |
| Sequelize | legacy reach | heaviest hydration, weaker TS support |

General performance order for CRUD-heavy paths:

```text
Drizzle ≈ raw SQL > Prisma > TypeORM > Sequelize
```

## When to Switch to Raw SQL

- analytical queries with many joins, grouping, or window functions
- bulk operations on `10,000+` rows
- DB-specific features the ORM cannot express well
- serverless paths where Prisma-like cold start matters
- high-throughput APIs around `1000+ RPS`

Recommended hybrid:
- CRUD -> ORM
- reporting and analytics -> raw SQL
- bulk operations -> raw SQL or DB-native bulk tools

## Review Checklist

- `findMany()` or equivalent uses explicit `select`
- relations use eager loading where appropriate
- no DB query runs inside loops
- transaction scope is minimal
- ORM bulk writes are not used for `10,000+` row operations
- production does not use `synchronize: true`
- generated SQL is reviewed with `EXPLAIN ANALYZE`
