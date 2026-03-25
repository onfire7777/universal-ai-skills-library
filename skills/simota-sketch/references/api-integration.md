# API Integration Reference

Purpose: Read this when Sketch needs canonical Gemini API integration rules, compatibility notes, error handling, rate guidance, or SynthID handling details.

## Contents
- SDK compatibility
- Authentication
- Request patterns
- Parameter rules
- Response handling
- Error handling
- Rate and cost guidance
- SynthID documentation

## SDK Compatibility

Verified baseline:
- `google-genai SDK v1.38.0`
- Google AI API with API-key auth
- confirmed model: `gemini-2.5-flash-image`

| SDK version | Config pattern | Notes |
| --- | --- | --- |
| `v1.38+` | `GenerateContentConfig(response_modalities=["IMAGE"])` | simple config only; control ratio/style in the prompt |
| `v1.50+` | `GenerateContentConfig(image_generation_config=ImageGenerationConfig(...))` | supports `aspect_ratio` and `person_generation` as parameters |

Default guidance:
- keep docs and code `v1.38+` compatible unless the user explicitly targets `v1.50+`

## Authentication

```python
import os
from google import genai

client = genai.Client(api_key=os.environ["GEMINI_API_KEY"])
```

Setup guidance:

```bash
pip install google-genai
export GEMINI_API_KEY="your-api-key-here"
```

Always include:

```bash
# .env
GEMINI_API_KEY=your-api-key-here

# .gitignore
.env
*.env
.env.*
```

## Model Rules

| Model | ID | API type | Speed | Cost | Use |
| --- | --- | --- | --- | --- | --- |
| Flash | `gemini-2.5-flash-image` | Google AI API | fast | low | default |
| Imagen 3.0 | `imagen-3.0-*` | Vertex AI only | medium | higher | Vertex-only paths |

Rule:
- on Google AI API, use `gemini-2.5-flash-image`
- `imagen-3.0-*` on Google AI API will return `404`

## Request Patterns

### Text-to-image

```python
response = client.models.generate_content(
    model="gemini-2.5-flash-image",
    contents="A modern minimalist workspace, soft natural lighting, widescreen 16:9 composition",
    config=types.GenerateContentConfig(
        response_modalities=["IMAGE"],
    ),
)
```

### Text + image response

Use `response_modalities=["TEXT", "IMAGE"]` when the text explanation matters.

### Reference-based editing

Pass the reference image with `types.Part.from_bytes(...)` plus the edit instruction.

### Iterative editing

Use `client.chats.create(...)` for multi-turn image refinement.

### Style transfer

Pass multiple reference images plus the transformation instruction.

## Parameter Rules

### `response_modalities`

| Value | Behavior | Use |
| --- | --- | --- |
| `["IMAGE"]` | image only | standard generation |
| `["TEXT", "IMAGE"]` | text and image | explanations or edit commentary |

### Prompt-based controls for `v1.38+`

| Control | Prompt instruction example |
| --- | --- |
| Aspect ratio | `"widescreen 16:9 composition"` |
| Style | `"photorealistic, DSLR quality"` |
| Quality | `"8K detail, professional photography"` |
| No people | `"no people, empty scene"` |
| Orientation | `"vertical portrait orientation"` |

### Aspect ratio guide

| Ratio | Prompt instruction | Use |
| --- | --- | --- |
| `1:1` | `"square format, 1:1 aspect ratio"` | social posts, icons, avatars |
| `3:2` | `"landscape 3:2 photography format"` | standard photography |
| `2:3` | `"portrait 2:3 vertical format"` | portrait/mobile |
| `16:9` | `"widescreen 16:9 composition"` | hero images, thumbnails |
| `9:16` | `"vertical 9:16 portrait orientation"` | stories, mobile vertical |
| `21:9` | `"ultra-wide 21:9 panoramic"` | banners |

### Reference-image limits

- maximum `14` reference images per request
- supported formats: PNG, JPEG, WebP, GIF first frame
- recommended maximum size: `4MB` per reference image

## Response Handling

Always:
- save images with timestamped filenames
- capture text responses when `TEXT` is requested
- write `metadata.json`

Suggested metadata fields:

```json
{
  "generated_at": "...",
  "prompt": "...",
  "model": "gemini-2.5-flash-image",
  "files": ["..."],
  "synthid": true
}
```

## Error Handling

Keep a comprehensive handler such as `generate_image_safe(...)` with retries.

| Error | Cause | Recovery |
| --- | --- | --- |
| `ResourceExhausted` | quota or rate limit | exponential backoff, quota check |
| `InvalidArgument` | bad prompt or parameters | fix prompt or params |
| `PermissionDenied` | invalid API key | verify `GEMINI_API_KEY` |
| `NotFound` | wrong model for API type | use `gemini-2.5-flash-image` on Google AI API |
| `ServiceUnavailable` | server issue | retry with backoff |
| empty response | content-policy block | simplify or adjust prompt |
| `DeadlineExceeded` | timeout | retry or simplify the request |

Common pitfalls:

| Pitfall | Symptom | Fix |
| --- | --- | --- |
| `ImageGenerationConfig` missing | `AttributeError` | on `< v1.50`, use simple config |
| `imagen-3.0-*` on Google AI API | `NotFound` | Vertex AI only |
| wrong Gemini model name | `NotFound` | keep the `-image` suffix |

## Rate Limits And Cost

| Tier | RPM (est.) | RPD (est.) | Cost/image (est.) |
| --- | --- | --- | --- |
| Free | `~15` | `~1,500` | `$0` |
| Paid | `~60` | higher | `~$0.02-0.04` |

Rules:
- treat all rate and cost figures as estimates
- tell the user to confirm the latest limits in Google AI Studio before production use
- preview `1-3` images before large batches

## Batch Guidance

For batch generation:
- generate sequentially with delay and retries
- show progress
- keep output paths deterministic
- estimate total cost before the run

## SynthID

All Gemini-generated images contain an invisible `SynthID` watermark.

Document:
- that the image is AI-generated
- the model used
- the generation timestamp
- that SynthID is present

Suggested disclosure snippet:

```text
This image was generated using Google Gemini API.
It contains an invisible SynthID watermark for AI-generated content identification.
Model: gemini-2.5-flash-image
Generated: [timestamp]
```
