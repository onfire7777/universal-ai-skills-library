# Gateway Alignment Notes (Current Gateway: WeryAI)

This skill follows a single-gateway architecture, and the current gateway happens to be WeryAI. The HTTP contract below follows the public **Image Generation** documentation. If the gateway changes later, this file must be updated as well.

## Base URL and Authentication

- **Base URL**: `https://api.weryai.com`
- **Header**: `Authorization: Bearer <API_KEY>` using environment variable `IMAGE_GEN_API_KEY`

## Async Endpoints

| Capability | Method | Path |
| --- | --- | --- |
| text-to-image | POST | `/v1/generation/text-to-image` |
| image-to-image | POST | `/v1/generation/image-to-image` |
| task status | GET | `/v1/generation/{taskId}/status` |

Successful business responses return **`success: true`** (or `status: 200`). Failed business responses return `success: false` with a business status code such as `1002`. On success, `data.task_ids` is a string array and polling uses the returned `task_id`.

## `task_status` Values

Follow the lowercase OpenAPI enum:

- `waiting`
- `processing`
- `succeed` -> `data.images` contains image URLs
- `failed` -> `data.msg` explains the failure

Do not rely on uppercase aliases such as `SUCCESS` or `FAILED`. If examples and OpenAPI disagree, trust the OpenAPI definition.

## Result Image Download

After `succeed`, the platform usually returns **download URLs** (often behind a CDN). Compared with APIs that return raw bytes directly, this introduces an extra network hop and can fail because the link is cold, temporarily unavailable, or subject to different authentication behavior.

The CLI handles each result image with:

- a request timeout (default: 120s)
- up to 4 attempts with exponential backoff
- an unauthenticated `GET` first, then a retry with the same Bearer token on `401/403`
- limited retries for `5xx`, `429`, `408`, and even `404` when the asset is not ready yet
- minimum payload validation to avoid saving empty responses or error pages as images

If download still fails, use the returned image URL manually and inspect the request with `task_id` plus [call history](https://weryai.com/api/history).

**Billing behavior**: once the platform-side task is already `succeed`, the CLI will **not** resubmit the generation request just because the file download failed. It only retries the HTTP file download layer, which avoids accidental double billing.

## Text-to-Image Request Body

Important fields:

| Field | Notes |
| --- | --- |
| `model` | required; see the model table for that endpoint |
| `prompt` | required; positive prompt, subject to documented length limits |
| `aspect_ratio` | required; e.g. `16:9`, `1:1` |
| `negative_prompt` | optional |
| `image_number` | optional, default 1 |
| `resolution` | optional; common values include `1k`, `2k`, `4k`, depending on model support |
| `use_web_search` | optional, default `false` |
| `webhook_url` | optional |
| `caller_id` | optional, int64, for business correlation |

## Image-to-Image Request Body

The image-to-image endpoint adds:

| Field | Notes |
| --- | --- |
| `images` | required string array. The docs show public HTTPS URLs, but the CLI can encode local files as `data:image/...;base64,...` and send them through the same field. |

The supported `model` set for image-to-image may differ from text-to-image. Always check the corresponding endpoint documentation.

## Common Business Status Codes

Examples from the documentation:

- `0` -> success
- `1001` -> parameter error
- `1002` -> authentication failure
- `1003` -> resource not found, such as an invalid task id

The script raises an error when the business response is not successful. For troubleshooting, compare the request body and response with [call history](https://weryai.com/api/history).

## Official Documentation Entry Points

- [API introduction](https://docs.weryai.com/en)
- [text-to-image](https://docs.weryai.com/api-reference/image-generation/submit-text-to-image-task)
- [image-to-image](https://docs.weryai.com/api-reference/image-generation/submit-image-to-image-task)
- [query task status](https://docs.weryai.com/api-reference/tasks/query-task-details)
- [documentation index `llms.txt`](https://docs.weryai.com/llms.txt)

The script validates `prompt` and `negative_prompt` lengths before making requests, following the OpenAPI `maxLength` values. The current implementation uses `prompt <= 2000` and `negative_prompt <= 1000`; update `scripts/main.ts` if the upstream contract changes.
