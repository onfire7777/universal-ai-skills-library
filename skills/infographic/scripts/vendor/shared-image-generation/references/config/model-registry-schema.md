# `MODELS.json` Schema (`image-generation`)

This optional registry helps the skill recommend a model based on workflow needs instead of treating model choice as an unstructured string.

Recommended file paths, in priority order:

- `<project>/.image-skills/image-generation/MODELS.json`
- `~/.image-skills/image-generation/MODELS.json`
- `<skill>/references/config/starter-models.json` (bundled starter snapshot from current WeryAI docs)

## Example

```json
{
  "version": 1,
  "models": [
    {
      "id": "gateway-model-fast",
      "label": "Fast general model",
      "modalities": ["text-to-image"],
      "best_for": ["general", "cover", "rednote"],
      "best_for_roles": ["rednote-cover"],
      "supports_ref": false,
      "supports_text_heavy": false,
      "quality_tier": "fast",
      "strengths": ["editorial", "flat-illustration", "speed"],
      "notes": "Good for quick visual exploration."
    },
    {
      "id": "gateway-model-pro",
      "label": "High quality ref-capable model",
      "modalities": ["text-to-image", "image-to-image"],
      "best_for": ["cover", "comic", "article"],
      "best_for_roles": ["cover-hero", "comic-page", "comic-character-sheet"],
      "supports_ref": true,
      "supports_text_heavy": false,
      "quality_tier": "high",
      "strengths": ["photoreal", "editorial", "comic"],
      "notes": "Best when consistency and image-to-image control matter."
    },
    {
      "id": "gateway-model-diagram",
      "label": "Diagram and infographic model",
      "modalities": ["text-to-image"],
      "best_for": ["infographic", "article"],
      "best_for_roles": ["infographic-dense", "article-framework", "article-flowchart"],
      "supports_ref": false,
      "supports_text_heavy": true,
      "quality_tier": "balanced",
      "strengths": ["diagram", "infographic", "labels"],
      "notes": "Prefer this for dense, label-heavy visuals."
    }
  ]
}
```

## Fields

| Field | Type | Description |
| --- | --- | --- |
| `version` | int | reserved, currently `1` |
| `models[]` | array | list of locally approved gateway model entries |
| `models[].id` | string | exact gateway `model` key |
| `models[].label` | string | human-readable short label |
| `models[].modalities` | string[] | e.g. `text-to-image`, `image-to-image` |
| `models[].best_for` | string[] | workflow hints such as `cover`, `rednote`, `infographic`, `comic`, `article`, `general` |
| `models[].best_for_roles` | string[] | optional sub-task hints such as `cover-hero`, `cover-typography-space`, `rednote-text-card`, `infographic-dense`, `comic-page`, `comic-character-sheet`, `article-framework`, `article-flowchart` |
| `models[].supports_ref` | bool | whether the model is dependable for reference-image workflows |
| `models[].supports_text_heavy` | bool | whether the model is dependable for label-heavy or diagram-heavy outputs |
| `models[].quality_tier` | string | `fast`, `balanced`, or `high` |
| `models[].strengths` | string[] | optional style or task hints such as `editorial`, `diagram`, `comic`, `photoreal` |
| `models[].notes` | string | optional short operational note |

## Recommendation CLI

After creating the registry, run:

```bash
npm run recommend-model -- --workflow comic --need-ref --quality high --style comic
npm run recommend-model -- --workflow comic --role comic-page
```

Use `--json` when the output should be consumed by another script.

This skill already ships with a bundled starter registry. To refresh it for a project from live docs:

```bash
npm run discover-image-models -- --out .image-skills/image-generation/MODELS.json
```
