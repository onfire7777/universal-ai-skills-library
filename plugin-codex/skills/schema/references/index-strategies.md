# Index Strategy Guide

Purpose: Use this file when selecting index type, composite order, or monitoring strategy.

Contents:
1. Index type selection
2. Composite index rules
3. Partial and covering indexes
4. Monitoring queries

## Index Type Selection

| Query pattern | PostgreSQL | MySQL |
|--------------|------------|-------|
| Exact match (`=`) | `B-tree` | `B-tree` |
| Range (`>`, `<`, `BETWEEN`) | `B-tree` | `B-tree` |
| Full-text search | `GIN (tsvector)` | `FULLTEXT` |
| JSON field lookup | `GIN (jsonb)` | Virtual column + `B-tree` |
| Array membership | `GIN` | N/A |
| Geospatial | `GiST` | `SPATIAL` |

## Composite Index Rules

- Order columns as `Equality -> Range -> Sort`.
- Respect the leftmost-prefix rule.
- Prefer the actual predicate order over abstract “important columns”.

Example:

```sql
WHERE status = 'active' AND created_at > '2024-01-01' ORDER BY name
```

Preferred index:

```sql
(status, created_at, name)
```

## Partial And Covering Indexes

| Pattern | Use when | Example |
|--------|----------|---------|
| Partial index | Query targets a stable subset | `WHERE deleted_at IS NULL` |
| Covering index | Heap fetch cost dominates | `INCLUDE (name, email)` |

```sql
CREATE INDEX idx_active_users
  ON users(email)
  WHERE deleted_at IS NULL;

CREATE INDEX idx_covering_users
  ON users(status)
  INCLUDE (name, email);
```

## Monitoring Queries

```sql
-- Unused indexes (PostgreSQL)
SELECT indexrelname, idx_scan, pg_size_pretty(pg_relation_size(indexrelid))
FROM pg_stat_user_indexes
WHERE NOT indisunique AND idx_scan < 50;

-- Missing-index hints
SELECT relname, seq_scan - idx_scan AS too_much_seq
FROM pg_stat_user_tables
WHERE seq_scan - idx_scan > 100;
```
