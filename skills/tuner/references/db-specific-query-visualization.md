# Database-Specific Optimizations & Query Plan Visualization

Purpose: Use this file when you need DB-specific tuning baselines or want to visualize a query plan for `Canvas`.

Contents:
- PostgreSQL baseline settings
- MySQL baseline settings
- SQLite baseline settings
- Canvas query-plan visualization

## PostgreSQL Baseline

- `shared_buffers`: `25%` of RAM
- `work_mem`: `RAM / max_connections / 4`
- `effective_cache_size`: `75%` of RAM
- `random_page_cost`: `1.1` for SSD, `4.0` for HDD
- maintain with regular `VACUUM ANALYZE`
- consider partial, covering, BRIN, and expression indexes where appropriate

## MySQL Baseline

- `innodb_buffer_pool_size`: `70%` of RAM
- `innodb_log_file_size`: `256M-2G`
- `query_cache_size`: `0` for MySQL 8.0
- `tmp_table_size` / `max_heap_table_size`: `64M-256M`
- maintain with `ANALYZE TABLE`, `OPTIMIZE TABLE`, and online schema tools when needed

## SQLite Baseline

- `PRAGMA journal_mode = WAL`
- `PRAGMA synchronous = NORMAL`
- `PRAGMA cache_size = -64000`
- `PRAGMA temp_store = MEMORY`
- run `ANALYZE` after bulk changes
- remember limited concurrent writes and no parallel query execution

## Canvas Query-Plan Visualization

Tuner can hand a query plan to `Canvas` as Mermaid input.

Example packet:

```markdown
/Canvas visualize query plan

Query: SELECT ...
Plan nodes:
1. Limit
2. Sort - 500ms
3. Hash Join - 200ms
4. Seq Scan users - 800ms <- bottleneck
5. Index Scan orders - 50ms
```

Use this when a diagram will help explain the bottleneck or compare before/after plans.
