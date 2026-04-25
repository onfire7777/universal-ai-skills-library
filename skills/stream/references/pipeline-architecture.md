# Pipeline Architecture

Purpose: choose the right pipeline mode and core architecture before designing detailed components.

## Contents

1. Batch vs streaming matrix
2. Decision rules
3. Architecture patterns
4. ETL vs ELT
5. Minimal orchestration template

## Batch vs Streaming Decision Matrix

| Factor | Batch | Streaming | Hybrid |
|--------|-------|-----------|--------|
| Latency requirement | hours to days | seconds to minutes | mixed |
| Data volume | large and bounded | continuous | both |
| Processing complexity | complex joins and aggregations | simpler incremental transforms | both |
| Cost sensitivity | high | moderate | balanced |
| Replay requirement | easy | harder | depends |

## Decision Rules

- `latency < 1 minute` -> streaming candidate
- `volume > 10K events/sec` + low latency -> Kafka + Flink/Spark
- complex analytics without sub-minute latency -> batch
- mixed operational + analytical outputs -> hybrid

## Architecture Patterns

### Lambda Architecture

Use only when legacy constraints require separate batch and speed layers. Avoid for new systems when Kappa or hybrid designs are sufficient.

### Kappa Architecture

Use when stream processing can be unified around an immutable log and replay is a first-class requirement.

### Medallion Architecture

Use for lakehouse-style pipelines with progressive refinement:
- `Bronze`: raw
- `Silver`: cleaned and standardized
- `Gold`: business-ready outputs

## ETL vs ELT

| Aspect | ETL | ELT |
|--------|-----|-----|
| Transform location | before load | after load |
| Best for | constrained destinations or legacy systems | powerful warehouses |
| Data volume | small to medium | medium to large |
| Flexibility | lower | higher |
| Common stack | Airflow + Python | dbt + Snowflake/BigQuery |

Rules:
- cloud warehouses usually favor ELT
- constrained operational databases usually favor ETL

## Minimal Orchestration Template

```python
with DAG("etl_orders_daily", default_args=default_args, catchup=False) as dag:
    extract = PythonOperator(task_id="extract_orders", python_callable=extract_from_source)
    validate_source = PythonOperator(task_id="validate_source", python_callable=run_quality_checks)
    transform = PythonOperator(task_id="transform_orders", python_callable=apply_business_logic)
    validate_transform = PythonOperator(task_id="validate_transform", python_callable=run_quality_checks)
    load = PostgresOperator(task_id="load_to_warehouse", postgres_conn_id="warehouse", sql="sql/load_orders.sql")

    extract >> validate_source >> transform >> validate_transform >> load
```
