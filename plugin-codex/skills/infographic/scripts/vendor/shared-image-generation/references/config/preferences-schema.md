# `EXTEND.md` Schema (`image-generation`)

Configuration file paths, in priority order:

- `<project>/.image-skills/image-generation/EXTEND.md`
- `~/.image-skills/image-generation/EXTEND.md`

Optional model registry paths:

- `<project>/.image-skills/image-generation/MODELS.json`
- `~/.image-skills/image-generation/MODELS.json`

Field-to-gateway mapping is described in [weryai-platform.md](../weryai-platform.md).

The file is Markdown with YAML front matter at the top.

## Example

```yaml
---
version: 1
default_model: <model key chosen from the gateway docs>
default_quality: 2k
default_aspect_ratio: "16:9"
# default_resolution: 2k
# default_use_web_search: false
# default_caller_id: 12345
batch:
  max_workers: 10
---
```

## Fields

| Field | Type | Description |
| --- | --- | --- |
| `version` | int | reserved, currently `1` |
| `default_model` | string | gateway `model` key; initialization should default to `GEMINI_3_1_FLASH_IMAGE` (Nano Banana 2) unless the user has already chosen another model |
| `default_quality` | `normal` \| `2k` | maps to `resolution`, unless overridden by `--resolution` |
| `default_aspect_ratio` | string | such as `1:1`; overridden by `--ar` |
| `default_resolution` | string | such as `2k`; common values are `1k`, `2k`, and `4k` |
| `default_use_web_search` | bool | default `use_web_search` behavior; can still be forced on by `--use-web-search` |
| `default_caller_id` | int | default `caller_id`; can be overridden by `--caller-id` |
| `batch.max_workers` | int | maximum batch concurrency; can be overridden by `BASE_IMAGE_MAX_WORKERS` |

`MODELS.json` is separate from `EXTEND.md`. Use it when this skill should rank multiple valid model keys by workflow or by narrower roles such as `comic-page` or `article-framework`. See [model-registry-schema.md](model-registry-schema.md).
