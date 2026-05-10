# Migration Patterns

Purpose: Use this file when planning safe schema changes, rollback, or framework-specific migration commands.

Contents:
1. Safe change decision tree
2. Expand-contract pattern
3. Zero-downtime index creation
4. Framework commands
5. Pre-migration checklist

## Safe Change Decision Tree

```text
Schema change needed
├── Adding new?
│   ├── New table -> CREATE TABLE (safe)
│   ├── New nullable column -> ADD COLUMN (safe)
│   ├── New NOT NULL column -> Expand-contract
│   └── New index -> Online / concurrent creation
├── Modifying existing?
│   ├── Rename column -> Expand-contract
│   ├── Change type -> Verify conversion safety first
│   ├── Add constraint -> Validate existing data first
│   └── Change default -> Usually safe
└── Removing?
    ├── Drop column -> Backup first, then phased removal
    ├── Drop table -> Backup required, irreversible
    ├── Drop index -> Safe, but verify query impact
    └── Drop constraint -> Safe, but verify integrity risk
```

## Expand-Contract Pattern

| Phase | Goal | Required actions |
|------|------|------------------|
| Expand | Introduce the new structure safely | Add new column/table, keep it nullable, dual-write if needed |
| Migrate | Backfill and validate | Batch copy data, validate consistency, add `NOT NULL` only after backfill |
| Contract | Remove the old structure | Switch reads, remove sync path, drop old column only after verification |

## Zero-Downtime Index Creation

```sql
-- PostgreSQL
CREATE INDEX CONCURRENTLY idx_name ON table_name(column_name);

-- MySQL 8.0+
ALTER TABLE t
  ADD INDEX idx_name (col),
  ALGORITHM=INPLACE,
  LOCK=NONE;
```

## Framework Migration Commands

| Framework | Create | Run | Rollback |
|-----------|--------|-----|----------|
| Prisma | `prisma migrate dev --name [name]` | `prisma migrate deploy` | Manual |
| TypeORM | `typeorm migration:generate -n [Name]` | `typeorm migration:run` | `typeorm migration:revert` |
| Drizzle | `drizzle-kit generate:pg` | `drizzle-kit push:pg` | Manual |
| Knex | `knex migrate:make [name]` | `knex migrate:latest` | `knex migrate:rollback` |

## Pre-Migration Checklist

- Backup is available and verified
- Migration was tested on production-like data
- Rollback path was tested or explicitly marked backup-required
- Lock duration was estimated
- Disk space and index-build impact were checked
- Low-traffic window is selected if blocking work remains
- Post-migration monitoring is prepared
