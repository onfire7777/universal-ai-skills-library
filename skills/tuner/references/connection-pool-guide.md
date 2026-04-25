# Connection Pool Tuning

Purpose: Use this file when the issue is connection pressure, pool sizing, or wait-event diagnosis.

Contents:
- sizing formula
- PostgreSQL settings
- monitoring queries

## Pool Size Formula

```text
pool_size = (core_count * 2) + effective_spindle_count
```

Typical range: `10-20` connections for many applications, but validate against workload and DB limits.

## Key PostgreSQL Settings

- `max_connections`
- `work_mem`
- `shared_buffers`
- `effective_cache_size`

## Monitoring

```sql
SELECT count(*) FROM pg_stat_activity WHERE state = 'active';
SELECT count(*) FROM pg_stat_activity WHERE state = 'idle';
SELECT wait_event_type, count(*) FROM pg_stat_activity GROUP BY 1;
```
