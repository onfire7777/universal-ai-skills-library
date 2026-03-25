# dbt Model Design

Purpose: canonical dbt layer structure, naming, materialization defaults, and test expectations for Stream output.

## Contents

1. Layer structure
2. Naming rules
3. Materialization defaults
4. Minimal templates
5. Core tests

## Layer Structure

```text
models/
├── staging/
├── intermediate/
├── marts/
└── exposures/
```

### Roles

| Layer | Purpose | Rules |
|-------|---------|-------|
| `staging` | source-aligned cleanup | `1:1` with source tables; no joins |
| `intermediate` | business logic and enrichment | small, reusable, single-responsibility transforms |
| `marts` | business-facing outputs | `dim_*`, `fct_*`, `rpt_*` |
| `exposures` | BI or API dependencies | document downstream usage |

## Naming Rules

- `stg_{source}__{entity}`
- `int_{entity}_{verb}`
- `dim_{entity}`
- `fct_{event}`
- `rpt_{report}`

## Materialization Defaults

| Layer | Default |
|-------|---------|
| `stg_*` | `view` |
| `int_*` | `ephemeral` or `table` |
| `dim_*` | `table` |
| `fct_*` | `table` or `incremental` |
| `rpt_*` | `table` |

## Minimal Templates

### Staging

```sql
{{ config(materialized='view', tags=['staging']) }}

with source as (
  select * from {{ source('raw', 'orders') }}
)
select
  id as order_id,
  customer_id,
  total_amount,
  created_at,
  updated_at
from source
```

### Intermediate

```sql
{{ config(materialized='table', tags=['intermediate']) }}

with orders as (
  select * from {{ ref('stg_orders') }}
)
select
  order_id,
  customer_id,
  total_amount,
  date_trunc('day', created_at) as order_date
from orders
```

## Core Tests

Always define:
- `unique` and `not_null` on primary keys
- `relationships` on foreign keys
- range or business-rule tests for critical measures
