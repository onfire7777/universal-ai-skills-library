# PostgreSQL 17 Performance Features

Purpose: Use this file when tuning PostgreSQL 17 or evaluating whether PG17 planner changes affect a recommendation.

Contents:
- optimizer improvements
- vacuum changes
- bulk-load changes
- recommended settings
- upgrade checks

## Query Optimizer Improvements

| Improvement | Typical effect | Notes |
|------------|----------------|-------|
| Streaming I/O | `~27%` faster seq scans | multiple-buffer reads |
| B-tree `IN` optimization | `~10%` faster | fewer leaf-page scans |
| correlated subquery -> `JOIN` conversion | can be dramatic, even `14000x` in best cases | verify actual plan |
| `UNION` optimization | `~15%` faster | better merge behavior |
| CTE statistics propagation | around `2x` on affected queries | planner sees more accurate stats |
| `IS NULL` on `NOT NULL` columns | huge plan simplification | unnecessary check removed |

## Vacuum And Bulk Operations

- vacuum memory footprint can drop significantly in PG17
- autovacuum becomes easier to run under concurrency
- `COPY` and large-row bulk operations can improve, sometimes near `2x`

## Recommended Settings

| Parameter | Suggested range |
|-----------|-----------------|
| `shared_buffers` | `25%` of RAM |
| `effective_cache_size` | `75%` of RAM |
| `work_mem` | `64MB-256MB` |
| `maintenance_work_mem` | `256MB-1GB` |
| `random_page_cost` | `1.1` for SSD |
| `jit` | on for complex-query workloads |
| `max_parallel_workers_per_gather` | `2-4` |

## Upgrade Validation

1. compare top queries with `EXPLAIN ANALYZE` before and after upgrade
2. compare top workload from `pg_stat_statements`
3. check vacuum memory behavior
4. benchmark `COPY` throughput
5. run `ANALYZE` after upgrade to refresh statistics
