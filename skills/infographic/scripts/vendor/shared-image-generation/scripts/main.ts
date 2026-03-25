#!/usr/bin/env bun
/**
 * Single-gateway image generation CLI (currently: WeryAI).
 * POST /v1/generation/text-to-image, image-to-image; GET .../status.
 */
import { access, mkdir, readFile, writeFile } from "node:fs/promises";
import path from "node:path";
import process from "node:process";
import { homedir } from "node:os";
import { existsSync, readFileSync, writeFileSync, mkdirSync } from "node:fs";
import { extname, resolve, dirname } from "node:path";

const MAX_ATTEMPTS = 3;
const DEFAULT_MAX_WORKERS = 10;
const POLL_WAIT_MS = 250;

// --- Gateway HTTP（当前: WeryAI；变更以 https://docs.weryai.com 为准）---
const BASE_URL = "https://api.weryai.com";

/** OpenAPI `maxLength` for text-to-image / image-to-image request schemas. */
const WERYAI_PROMPT_MAX_LEN = 2000;
const WERYAI_NEGATIVE_PROMPT_MAX_LEN = 1000;

/** Result URLs are gateway CDN links — extra hop vs direct bytes; harden download path. */
const DOWNLOAD_MAX_ATTEMPTS = 4;
const DOWNLOAD_TIMEOUT_MS = 120_000;
const DOWNLOAD_BACKOFF_MS = 700;
const DOWNLOAD_MIN_BYTES = 32;

const STYLE_PRESETS = {
  photoreal: "photorealistic, natural lighting, realistic texture, high detail",
  cinematic: "cinematic composition, dramatic lighting, film still, moody atmosphere",
  anime: "anime style, clean linework, vibrant colors, expressive character design",
  manga: "manga style, dynamic ink lines, screentone shading, black and white emphasis",
  watercolor: "watercolor illustration, soft brush texture, paper grain, gentle palette",
  "flat-illustration": "flat illustration, simplified shapes, clean vector feel, limited palette",
  "3d-render": "3d render, volumetric lighting, material detail, polished surfaces",
  poster: "poster design, bold focal point, strong typography space, high contrast",
  editorial: "editorial illustration, refined composition, premium layout feel, restrained palette",
  infographic: "infographic style, clear hierarchy, icon-friendly layout, data visualization feel",
  chalk: "chalk drawing, textured strokes, blackboard feel, hand-drawn annotations",
  "ink-brush": "ink brush painting, expressive brushwork, elegant negative space, East Asian aesthetic",
} as const;

type StylePreset = keyof typeof STYLE_PRESETS;

const SKILL_NAMESPACE = process.env.IMAGE_SKILL_NAMESPACE?.trim() || "image-generation";
const SKILL_LABEL = process.env.IMAGE_SKILL_LABEL?.trim() || SKILL_NAMESPACE;

function getApiKey(): string {
  const k = process.env.IMAGE_GEN_API_KEY || process.env.IMAGE_GEN_API_KEY;
  if (!k?.trim()) {
    console.error(
      `Missing IMAGE_GEN_API_KEY. Get a key at https://weryai.com/api/keys, then add it to the environment or project .image-skills/${SKILL_NAMESPACE}/.env. IMAGE_GEN_API_KEY is still accepted as a compatibility alias. If you are setting up this skill for the first time, run \`npm run doctor -- --project <your-project>\` from the ${SKILL_LABEL} skill directory. Doctor is read-only: it prints readiness, exact next commands, and does not create files.`
    );
    process.exit(1);
  }
  return k.trim();
}

type TaskStatus = "waiting" | "processing" | "succeed" | "failed";

interface ApiEnvelope<T> {
  status: number;
  success?: boolean;
  desc?: string;
  message?: string;
  data?: T;
}

function isApiSuccess(json: ApiEnvelope<unknown>): boolean {
  if (typeof json.success === "boolean") return json.success;
  return json.status === 200 || json.status === 0;
}

type RequestPreview = {
  endpoint: string;
  method: "POST";
  body: Record<string, unknown>;
};

interface SubmitResult {
  batch_id?: number;
  task_ids?: string[];
}

interface TaskDetail {
  task_id?: string;
  task_status?: TaskStatus;
  msg?: string;
  images?: string[];
}

function weryaiBusinessError(json: ApiEnvelope<unknown>): Error {
  const code = json.status;
  const base = json.message ?? json.desc ?? `WeryAI business status ${code}`;
  const hints: Record<number, string> = {
    1001: "parameter error — check model, prompt length, aspect_ratio, images, resolution per docs",
    1002: "authentication — verify IMAGE_GEN_API_KEY",
    1003: "resource not found (e.g. invalid task id)",
  };
  const extra = hints[code] ? ` (${hints[code]})` : "";
  return new Error(`${base}${extra}`);
}

async function apiPost<T>(pathSuffix: string, body: Record<string, unknown>): Promise<ApiEnvelope<T>> {
  const res = await fetch(`${BASE_URL}${pathSuffix}`, {
    method: "POST",
    headers: {
      Authorization: `Bearer ${getApiKey()}`,
      "Content-Type": "application/json",
    },
    body: JSON.stringify(body),
  });
  const json = (await res.json()) as ApiEnvelope<T>;
  if (!res.ok) {
    const msg = json?.message ?? json?.desc ?? res.statusText;
    throw new Error(`HTTP ${res.status}: ${msg}`);
  }
  if (!isApiSuccess(json)) throw weryaiBusinessError(json);
  return json;
}

async function apiGet<T>(pathSuffix: string): Promise<ApiEnvelope<T>> {
  const res = await fetch(`${BASE_URL}${pathSuffix}`, {
    headers: { Authorization: `Bearer ${getApiKey()}` },
  });
  const json = (await res.json()) as ApiEnvelope<T>;
  if (!res.ok) {
    const msg = json?.message ?? json?.desc ?? res.statusText;
    throw new Error(`HTTP ${res.status}: ${msg}`);
  }
  if (!isApiSuccess(json)) throw weryaiBusinessError(json);
  return json;
}

async function pollTask(
  taskId: string,
  intervalMs: number,
  timeoutMs: number,
  quiet: boolean
): Promise<TaskDetail> {
  const start = Date.now();
  let last: TaskDetail | undefined;
  while (Date.now() - start < timeoutMs) {
    const r = await apiGet<TaskDetail>(`/v1/generation/${encodeURIComponent(taskId)}/status`);
    last = r.data;
    const st = r.data?.task_status;
    if (!quiet && st) console.error(`[poll] ${taskId} → ${st}`);
    if (st === "succeed" || st === "failed") return r.data ?? {};
    await new Promise((r2) => setTimeout(r2, intervalMs));
  }
  throw new Error(`Timeout after ${timeoutMs}ms waiting for task ${taskId}. Last: ${JSON.stringify(last)}`);
}

type FetchResult = { bytes: Uint8Array; contentType: string | null };

/**
 * Pull image bytes from gateway result URL (often CDN). Mitigates transient failures, cold CDN,
 * and occasional link quirks by: timeout, retries with backoff, unauthenticated then Bearer retry.
 */
async function fetchImageBytes(url: string, quiet: boolean): Promise<FetchResult> {
  let lastErr = "";

  const transientStatus = (status: number) =>
    status === 408 ||
    status === 429 ||
    status === 502 ||
    status === 503 ||
    status === 504 ||
    (status >= 520 && status <= 599);

  const transientNetwork = (e: unknown) => {
    const name = e instanceof Error ? e.name : "";
    const msg = e instanceof Error ? e.message : String(e);
    return (
      name === "AbortError" ||
      msg.includes("timeout") ||
      msg.includes("Timeout") ||
      msg.includes("fetch failed") ||
      msg.includes("ECONNRESET") ||
      msg.includes("ETIMEDOUT") ||
      msg.includes("EAI_AGAIN")
    );
  };

  const sleep = (ms: number) => new Promise((r) => setTimeout(r, ms));

  for (let attempt = 1; attempt <= DOWNLOAD_MAX_ATTEMPTS; attempt++) {
    let res: Response | null = null;
    let threw: unknown = null;

    for (const useAuth of [false, true] as const) {
      try {
        const init: RequestInit = { signal: AbortSignal.timeout(DOWNLOAD_TIMEOUT_MS) };
        if (useAuth) init.headers = { Authorization: `Bearer ${getApiKey()}` };
        res = await fetch(url, init);
        threw = null;
        if (res.ok) break;
        lastErr = `HTTP ${res.status}`;
        if (!useAuth && (res.status === 401 || res.status === 403)) continue;
        break;
      } catch (e) {
        threw = e;
        res = null;
        lastErr = e instanceof Error ? e.message : String(e);
        if (!useAuth) continue;
        break;
      }
    }

    if (threw) {
      if (attempt < DOWNLOAD_MAX_ATTEMPTS && transientNetwork(threw)) {
        const wait = DOWNLOAD_BACKOFF_MS * 2 ** (attempt - 1);
        if (!quiet) console.error(`[download] ${lastErr} → retry ${attempt}/${DOWNLOAD_MAX_ATTEMPTS} in ${wait}ms`);
        await sleep(wait);
        continue;
      }
      throw new Error(`Download failed (network): ${lastErr}`);
    }

    if (!res) {
      throw new Error(`Download failed: no response (${lastErr})`);
    }

    if (!res.ok) {
      const st = res.status;
      const retry =
        attempt < DOWNLOAD_MAX_ATTEMPTS &&
        (transientStatus(st) || st === 404 || st === 403);
      if (retry) {
        const wait = DOWNLOAD_BACKOFF_MS * 2 ** (attempt - 1);
        if (!quiet) console.error(`[download] HTTP ${st} → retry ${attempt}/${DOWNLOAD_MAX_ATTEMPTS} in ${wait}ms`);
        await sleep(wait);
        continue;
      }
      throw new Error(`Download HTTP ${st} for ${url.length > 96 ? url.slice(0, 96) + "…" : url}`);
    }

    let buf: Uint8Array;
    try {
      buf = new Uint8Array(await res.arrayBuffer());
    } catch (e) {
      lastErr = e instanceof Error ? e.message : String(e);
      if (attempt < DOWNLOAD_MAX_ATTEMPTS) {
        const wait = DOWNLOAD_BACKOFF_MS * 2 ** (attempt - 1);
        if (!quiet) console.error(`[download] body read error → retry in ${wait}ms`);
        await sleep(wait);
        continue;
      }
      throw e;
    }

    if (buf.byteLength < DOWNLOAD_MIN_BYTES) {
      lastErr = `payload too small (${buf.byteLength} bytes; link may be expired or not an image)`;
      if (attempt < DOWNLOAD_MAX_ATTEMPTS) {
        const wait = DOWNLOAD_BACKOFF_MS * 2 ** (attempt - 1);
        if (!quiet) console.error(`[download] ${lastErr} → retry in ${wait}ms`);
        await sleep(wait);
        continue;
      }
      throw new Error(lastErr);
    }

    return { bytes: buf, contentType: res.headers.get("content-type") };
  }

  throw new Error(`Download exhausted after ${DOWNLOAD_MAX_ATTEMPTS} attempts: ${lastErr}`);
}

function mimeForExt(ext: string): string {
  switch (ext.toLowerCase()) {
    case ".png":
      return "image/png";
    case ".webp":
      return "image/webp";
    case ".jpg":
    case ".jpeg":
      return "image/jpeg";
    default:
      return "application/octet-stream";
  }
}

function extensionForContentType(ct: string): string | null {
  const lower = (ct.split(";")[0] ?? "").trim().toLowerCase();
  switch (lower) {
    case "image/png": return ".png";
    case "image/jpeg":
    case "image/jpg": return ".jpg";
    case "image/webp": return ".webp";
    case "image/gif": return ".gif";
    case "image/avif": return ".avif";
    default: return null;
  }
}

function toImagePayloadValue(pathOrUrl: string): string {
  const t = pathOrUrl.trim();
  if (/^https?:\/\//i.test(t)) return t;
  const abs = resolve(t);
  if (!existsSync(abs)) throw new Error(`Reference image not found: ${abs}`);
  const buf = readFileSync(abs);
  const mime = mimeForExt(extname(abs));
  return `data:${mime};base64,${buf.toString("base64")}`;
}

// --- EXTEND / env ---

type Quality = "normal" | "2k";

type ExtendConfig = {
  version?: number;
  default_model?: string | null;
  default_quality?: Quality | null;
  default_aspect_ratio?: string | null;
  default_resolution?: string | null;
  default_use_web_search?: boolean | null;
  default_caller_id?: number | null;
  batch?: { max_workers?: number | null };
};

function extractYamlFrontMatter(content: string): string | null {
  const match = content.match(/^---\s*\n([\s\S]*?)\n---\s*$/m);
  return match ? match[1] : null;
}

function parseExtendYaml(yaml: string): ExtendConfig {
  const config: ExtendConfig = {};
  const lines = yaml.split("\n");
  let inBatch = false;

  for (const line of lines) {
    const trimmed = line.trim();
    const indent = line.match(/^\s*/)?.[0].length ?? 0;
    if (!trimmed || trimmed.startsWith("#")) continue;

    if (trimmed.includes(":") && !trimmed.startsWith("-")) {
      const colonIdx = trimmed.indexOf(":");
      const key = trimmed.slice(0, colonIdx).trim();
      let value = trimmed.slice(colonIdx + 1).trim().replace(/^["']|["']$/g, "");

      if (key === "version") {
        config.version = value === "null" || value === "" ? 1 : parseInt(value, 10);
      } else if (key === "default_model") {
        if (!value || value === "null") config.default_model = null;
        else config.default_model = value;
        inBatch = false;
      } else if (key === "default_quality") {
        if (value === "normal" || value === "2k") config.default_quality = value;
        else config.default_quality = null;
        inBatch = false;
      } else if (key === "default_aspect_ratio") {
        config.default_aspect_ratio = value === "null" || !value ? null : value;
        inBatch = false;
      } else if (key === "default_resolution") {
        config.default_resolution = value === "null" || !value ? null : value;
        inBatch = false;
      } else if (key === "default_use_web_search") {
        const v = value.toLowerCase();
        if (value === "null" || value === "") config.default_use_web_search = null;
        else config.default_use_web_search = v === "true" || v === "1" || v === "yes";
        inBatch = false;
      } else if (key === "default_caller_id") {
        if (value === "null" || !value) config.default_caller_id = null;
        else {
          const n = parseInt(value, 10);
          config.default_caller_id = Number.isFinite(n) ? n : null;
        }
        inBatch = false;
      } else if (key === "batch") {
        inBatch = true;
        config.batch ??= {};
      } else if (key === "max_workers" && inBatch && indent >= 2) {
        config.batch ??= {};
        config.batch.max_workers = value === "null" || !value ? null : parseInt(value, 10);
      }
    }
  }
  return config;
}

async function loadExtendConfig(): Promise<ExtendConfig> {
  const cwd = process.cwd();
  const home = homedir();
  const paths = [
    path.join(cwd, ".image-skills", SKILL_NAMESPACE, "EXTEND.md"),
    path.join(home, ".image-skills", SKILL_NAMESPACE, "EXTEND.md"),
  ];
  for (const p of paths) {
    try {
      const content = await readFile(p, "utf8");
      const yaml = extractYamlFrontMatter(content);
      if (!yaml) continue;
      return parseExtendYaml(yaml);
    } catch {
      continue;
    }
  }
  return {};
}

async function loadEnvFile(p: string): Promise<Record<string, string>> {
  try {
    const content = await readFile(p, "utf8");
    const env: Record<string, string> = {};
    for (const line of content.split("\n")) {
      const trimmed = line.trim();
      if (!trimmed || trimmed.startsWith("#")) continue;
      const idx = trimmed.indexOf("=");
      if (idx === -1) continue;
      const key = trimmed.slice(0, idx).trim();
      let val = trimmed.slice(idx + 1).trim();
      if ((val.startsWith('"') && val.endsWith('"')) || (val.startsWith("'") && val.endsWith("'"))) {
        val = val.slice(1, -1);
      }
      env[key] = val;
    }
    return env;
  } catch {
    return {};
  }
}

async function loadEnv(): Promise<void> {
  const home = homedir();
  const cwd = process.cwd();
  const mergedFiles = {
    ...(await loadEnvFile(path.join(home, ".image-skills", SKILL_NAMESPACE, ".env"))),
    ...(await loadEnvFile(path.join(cwd, ".image-skills", SKILL_NAMESPACE, ".env"))),
  };
  for (const [k, v] of Object.entries(mergedFiles)) {
    if (!process.env[k]) process.env[k] = v;
  }
}

// --- CLI types ---

type CliArgs = {
  prompt: string | null;
  promptFiles: string[];
  imagePath: string | null;
  batchFile: string | null;
  jobs: number | null;
  model: string | null;
  style: StylePreset | null;
  aspectRatio: string | null;
  size: string | null;
  quality: Quality | null;
  imageSize: string | null;
  referenceImages: string[];
  n: number;
  resolution: string | null;
  negativePrompt: string | null;
  webhookUrl: string | null;
  pollIntervalMs: number;
  pollTimeoutMs: number;
  noDownload: boolean;
  json: boolean;
  help: boolean;
  /** WeryAI `use_web_search` */
  useWebSearch: boolean;
  /** WeryAI `caller_id` */
  callerId: number | null;
  dryRun: boolean;
};

type BatchTaskInput = {
  id?: string;
  prompt?: string | null;
  promptFiles?: string[];
  image?: string;
  model?: string | null;
  style?: StylePreset | null;
  ar?: string | null;
  size?: string | null;
  quality?: Quality | null;
  imageSize?: "1K" | "2K" | "4K" | null;
  ref?: string[];
  n?: number;
  provider?: string | null;
  webhook_url?: string;
  negative_prompt?: string;
  resolution?: string;
  use_web_search?: boolean;
  caller_id?: number;
};

type BatchFile = BatchTaskInput[] | { tasks: BatchTaskInput[]; jobs?: number | null };

function printUsage(): void {
  console.log(`Image generation (gateway: api.weryai.com)
  ${process.argv[1] ?? "main.ts"} --prompt "..." --image out.png [--ar 16:9] [--ref a.png]
  ${process.argv[1] ?? "main.ts"} --promptfiles a.md b.md --image out.png
  ${process.argv[1] ?? "main.ts"} --batchfile batch.json [--jobs 4] [--json]

Legacy subcommands:
  ... text-to-image [options]
  ... image-to-image [options]
  ... status --task-id <id>

Options:
  -p, --prompt <text>       Prompt
  --promptfiles <files...>  Concatenate prompt from files
  --image <path>            Output image path (required for single mode)
  --batchfile <path>        JSON batch (see SKILL.md)
  --jobs <n>                Parallel workers in batch mode
  -m, --model <id>          Required unless EXTEND default_model or IMAGE_GEN_DEFAULT_MODEL (else agent supplies after doc lookup)
  --style <preset>          Style preset (${Object.keys(STYLE_PRESETS).join(", ")})
  --ar <ratio>              Aspect ratio (default: EXTEND or 1:1)
  --size <WxH>              If --ar omitted, ratio may be inferred (e.g. 1024x768 → 4:3)
  --quality normal|2k      Maps to resolution 1k|2k (default: 2k)
  --imageSize 1K|2K|4K      Fine-tune resolution mapping with --quality
  --resolution <s>        Pass through to API (overrides quality mapping)
  --negative-prompt <t>    Negative prompt
  --ref, --reference <paths...>  Reference images (image-to-image)
  --n <count>             image_number (default 1)
  --webhook-url <url>      Gateway webhook_url
  --use-web-search         Gateway use_web_search=true
  --caller-id <n>          Gateway caller_id (int64)
  --poll-interval-ms <n>   Default 2500
  --poll-timeout-ms <n>    Default 600000
  --no-download           Only print URLs (no file write)
  --dry-run               Print final request body (JSON), do not call API
  --json                  JSON summary / batch report
  -h, --help

Env: IMAGE_GEN_API_KEY (preferred); IMAGE_GEN_API_KEY compatible; IMAGE_GEN_DEFAULT_MODEL optional; or EXTEND / --model
Also loaded: <cwd>/.image-skills/${SKILL_NAMESPACE}/.env and ~/.image-skills/${SKILL_NAMESPACE}/.env
Prefs: .image-skills/${SKILL_NAMESPACE}/EXTEND.md (project or home)
`);
}

function takeMany(argv: string[], i: number): { items: string[]; next: number } {
  const items: string[] = [];
  let j = i + 1;
  while (j < argv.length) {
    const v = argv[j]!;
    if (v.startsWith("-")) break;
    items.push(v);
    j++;
  }
  return { items, next: j - 1 };
}

function normalizeStylePreset(value: string): StylePreset | null {
  const normalized = value.trim().toLowerCase().replace(/[_\s]+/g, "-");
  return Object.prototype.hasOwnProperty.call(STYLE_PRESETS, normalized)
    ? (normalized as StylePreset)
    : null;
}

function applyStylePreset(prompt: string, style: StylePreset | null | undefined): string {
  if (!style) return prompt;
  return `${prompt.trim()}\n\nStyle preset (${style}): ${STYLE_PRESETS[style]}.`;
}

function parseArgs(argv: string[]): CliArgs {
  const out: CliArgs = {
    prompt: null,
    promptFiles: [],
    imagePath: null,
    batchFile: null,
    jobs: null,
    model: null,
    style: null,
    aspectRatio: null,
    size: null,
    quality: null,
    imageSize: null,
    referenceImages: [],
    n: 1,
    resolution: null,
    negativePrompt: null,
    webhookUrl: null,
    pollIntervalMs: 2500,
    pollTimeoutMs: 600_000,
    noDownload: false,
    json: false,
    help: false,
    useWebSearch: false,
    callerId: null,
    dryRun: false,
  };
  const positional: string[] = [];

  for (let i = 0; i < argv.length; i++) {
    const a = argv[i]!;
    if (a === "--use-web-search") {
      out.useWebSearch = true;
      continue;
    }
    if (a === "--dry-run") {
      out.dryRun = true;
      continue;
    }
    if (a === "--help" || a === "-h") {
      out.help = true;
      continue;
    }
    if (a === "--json") {
      out.json = true;
      continue;
    }
    if (a === "--no-download") {
      out.noDownload = true;
      continue;
    }
    if (a === "--prompt" || a === "-p") {
      const v = argv[++i];
      if (!v) throw new Error(`Missing value for ${a}`);
      out.prompt = v;
      continue;
    }
    if (a === "--promptfiles") {
      const { items, next } = takeMany(argv, i);
      if (!items.length) throw new Error("Missing files for --promptfiles");
      out.promptFiles.push(...items);
      i = next;
      continue;
    }
    if (a === "--image") {
      const v = argv[++i];
      if (!v) throw new Error("Missing value for --image");
      out.imagePath = v;
      continue;
    }
    if (a === "--output" || a === "-o") {
      const v = argv[++i];
      if (!v) throw new Error("Missing value for --output");
      out.imagePath = v;
      continue;
    }
    if (a === "--batchfile") {
      const v = argv[++i];
      if (!v) throw new Error("Missing value for --batchfile");
      out.batchFile = v;
      continue;
    }
    if (a === "--jobs") {
      const v = argv[++i];
      if (!v) throw new Error("Missing value for --jobs");
      out.jobs = parseInt(v, 10);
      if (Number.isNaN(out.jobs) || out.jobs < 1) throw new Error(`Invalid worker count: ${v}`);
      continue;
    }
    if (a === "--model" || a === "-m") {
      const v = argv[++i];
      if (!v) throw new Error(`Missing value for ${a}`);
      out.model = v;
      continue;
    }
    if (a === "--style") {
      const v = argv[++i];
      if (!v) throw new Error("Missing value for --style");
      const style = normalizeStylePreset(v);
      if (!style) throw new Error(`Invalid --style: ${v}. Available: ${Object.keys(STYLE_PRESETS).join(", ")}`);
      out.style = style;
      continue;
    }
    if (a === "--ar") {
      const v = argv[++i];
      if (!v) throw new Error("Missing value for --ar");
      out.aspectRatio = v;
      continue;
    }
    if (a === "--aspect-ratio") {
      const v = argv[++i];
      if (!v) throw new Error("Missing value for --aspect-ratio");
      out.aspectRatio = v;
      continue;
    }
    if (a === "--size") {
      const v = argv[++i];
      if (!v) throw new Error("Missing value for --size");
      out.size = v;
      continue;
    }
    if (a === "--quality") {
      const v = argv[++i];
      if (v !== "normal" && v !== "2k") throw new Error(`Invalid quality: ${v}`);
      out.quality = v;
      continue;
    }
    if (a === "--imageSize") {
      const v = argv[++i]?.toUpperCase();
      if (v !== "1K" && v !== "2K" && v !== "4K") throw new Error(`Invalid imageSize: ${v}`);
      out.imageSize = v;
      continue;
    }
    if (a === "--ref" || a === "--reference") {
      const { items, next } = takeMany(argv, i);
      if (!items.length) throw new Error(`Missing files for ${a}`);
      out.referenceImages.push(...items);
      i = next;
      continue;
    }
    if (a === "--n") {
      const v = argv[++i];
      if (!v) throw new Error("Missing value for --n");
      out.n = parseInt(v, 10);
      if (Number.isNaN(out.n) || out.n < 1) throw new Error(`Invalid count: ${v}`);
      continue;
    }
    if (a === "--resolution") {
      const v = argv[++i];
      if (!v) throw new Error("Missing value for --resolution");
      out.resolution = v;
      continue;
    }
    if (a === "--negative-prompt") {
      const v = argv[++i];
      if (!v) throw new Error("Missing value for --negative-prompt");
      out.negativePrompt = v;
      continue;
    }
    if (a === "--webhook-url") {
      const v = argv[++i];
      if (!v) throw new Error("Missing value for --webhook-url");
      out.webhookUrl = v;
      continue;
    }
    if (a === "--caller-id") {
      const v = argv[++i];
      if (!v) throw new Error("Missing value for --caller-id");
      const n = parseInt(v, 10);
      if (!Number.isFinite(n)) throw new Error(`Invalid --caller-id: ${v}`);
      out.callerId = n;
      continue;
    }
    if (a === "--poll-interval-ms") {
      const v = argv[++i];
      if (!v) throw new Error("Missing value for --poll-interval-ms");
      out.pollIntervalMs = parseInt(v, 10);
      continue;
    }
    if (a === "--poll-timeout-ms") {
      const v = argv[++i];
      if (!v) throw new Error("Missing value for --poll-timeout-ms");
      out.pollTimeoutMs = parseInt(v, 10);
      continue;
    }
    if (a.startsWith("-")) throw new Error(`Unknown option: ${a}`);
    positional.push(a);
  }

  if (!out.prompt && out.promptFiles.length === 0 && positional.length > 0) {
    out.prompt = positional.join(" ");
  }
  return out;
}

function mergeConfig(args: CliArgs, extend: ExtendConfig): CliArgs {
  const envModel = process.env.IMAGE_GEN_DEFAULT_MODEL?.trim() || null;
  const useWs = Boolean(args.useWebSearch) || extend.default_use_web_search === true;
  const cid = args.callerId ?? extend.default_caller_id ?? null;
  return {
    ...args,
    quality: args.quality ?? extend.default_quality ?? null,
    aspectRatio: args.aspectRatio ?? extend.default_aspect_ratio ?? null,
    resolution: args.resolution ?? extend.default_resolution ?? null,
    model: args.model?.trim() || extend.default_model?.trim() || envModel || null,
    useWebSearch: useWs,
    callerId: cid !== null && cid !== undefined && Number.isFinite(cid) ? cid : null,
  };
}

function gcd(a: number, b: number): number {
  while (b) [a, b] = [b, a % b];
  return Math.abs(a);
}

function inferAspectRatioFromSize(size: string): string | null {
  const m = size.trim().match(/^(\d+)\s*[x*]\s*(\d+)$/i);
  if (!m) return null;
  let w = parseInt(m[1]!, 10);
  let h = parseInt(m[2]!, 10);
  const g = gcd(w, h);
  if (g === 0) return null;
  w /= g;
  h /= g;
  return `${w}:${h}`;
}

/**
 * `resolution` is optional; allowed values depend on the model.
 * Common gateway values: 1k, 2k, 4k. `--resolution` passes through verbatim for exact control.
 */
function resolveResolution(args: CliArgs): string | undefined {
  if (args.resolution) return args.resolution;
  const q = args.quality ?? "2k";
  const is = (args.imageSize ?? "").toUpperCase();
  if (is === "4K") return "4k";
  if (is === "1K") return "1k";
  if (is === "2K") return "2k";
  return q === "normal" ? "1k" : "2k";
}

function resolveAspectRatio(args: CliArgs): string {
  if (args.aspectRatio) return args.aspectRatio;
  if (args.size) {
    const inferred = inferAspectRatioFromSize(args.size);
    if (inferred) return inferred;
  }
  return "1:1";
}

function resolveModel(args: CliArgs): string {
  const m = args.model?.trim();
  if (!m) {
    throw new Error(
      `Missing model: need --model, EXTEND default_model, or IMAGE_GEN_DEFAULT_MODEL. Set default_model in .image-skills/${SKILL_NAMESPACE}/EXTEND.md, export IMAGE_GEN_DEFAULT_MODEL, or rerun with --model <id>. If setup is still incomplete, run \`npm run doctor -- --project <your-project>\` or \`npm run setup -- --project <your-project>\` from the ${SKILL_LABEL} skill directory. Agent should start from the bundled WeryAI starter model list, refresh from docs if needed, rank candidates, ask the user to confirm one, then rerun with --model <id>.`
    );
  }
  return m;
}

function modelLabelForReport(args: CliArgs): string {
  const m = args.model?.trim();
  return m || "(no model configured)";
}

function parsePositiveInt(v: string | undefined): number | null {
  if (!v) return null;
  const n = parseInt(v, 10);
  return Number.isFinite(n) && n > 0 ? n : null;
}

function getMaxWorkers(extend: ExtendConfig): number {
  const env = parsePositiveInt(process.env.BASE_IMAGE_MAX_WORKERS);
  const cfg = extend.batch?.max_workers;
  return Math.max(1, env ?? cfg ?? DEFAULT_MAX_WORKERS);
}

function getConcurrency(): { concurrency: number; startIntervalMs: number } {
  return {
    concurrency: parsePositiveInt(process.env.BASE_IMAGE_GEN_CONCURRENCY) ?? 3,
    startIntervalMs: parsePositiveInt(process.env.BASE_IMAGE_GEN_START_INTERVAL_MS) ?? 1100,
  };
}

async function readPromptFromFiles(files: string[]): Promise<string> {
  const parts: string[] = [];
  for (const f of files) {
    parts.push(await readFile(f, "utf8"));
  }
  return parts.join("\n\n");
}

async function readPromptFromStdin(): Promise<string | null> {
  if (process.stdin.isTTY) return null;
  try {
    const chunks: Buffer[] = [];
    for await (const chunk of process.stdin) {
      chunks.push(Buffer.isBuffer(chunk) ? chunk : Buffer.from(chunk));
    }
    const value = Buffer.concat(chunks).toString("utf8").trim();
    return value.length > 0 ? value : null;
  } catch {
    return null;
  }
}

function normalizeOutputImagePath(p: string): string {
  const full = path.resolve(p);
  const ext = path.extname(full);
  if (ext) return full;
  return `${full}.png`;
}

async function validateReferenceImages(paths: string[]): Promise<void> {
  for (const refPath of paths) {
    const t = refPath.trim();
    if (/^https?:\/\//i.test(t)) continue;
    const fullPath = path.resolve(refPath);
    try {
      await access(fullPath);
    } catch {
      throw new Error(`Reference image not found: ${fullPath}`);
    }
  }
}

function assertWeryaiPromptBounds(prompt: string, negativePrompt: string | null | undefined): void {
  const p = prompt.trim();
  if (!p) throw new Error("prompt is empty (WeryAI requires non-empty prompt)");
  if (p.length > WERYAI_PROMPT_MAX_LEN) {
    throw new Error(
      `prompt length ${p.length} exceeds WeryAI OpenAPI maxLength ${WERYAI_PROMPT_MAX_LEN}`
    );
  }
  if (negativePrompt && negativePrompt.length > WERYAI_NEGATIVE_PROMPT_MAX_LEN) {
    throw new Error(
      `negative_prompt length ${negativePrompt.length} exceeds WeryAI maxLength ${WERYAI_NEGATIVE_PROMPT_MAX_LEN}`
    );
  }
}

async function runAsyncGeneration(input: {
  prompt: string;
  model: string;
  aspectRatio: string;
  resolution?: string;
  negativePrompt?: string | null;
  imageNumber: number;
  referenceImages: string[];
  webhookUrl?: string | null;
  useWebSearch: boolean;
  callerId: number | null;
  pollIntervalMs: number;
  pollTimeoutMs: number;
  quietPoll: boolean;
  dryRun?: boolean;
}): Promise<{
  submit: SubmitResult | undefined;
  detail: TaskDetail | undefined;
  requestPreview?: RequestPreview;
  dryRun?: boolean;
}> {
  const {
    prompt,
    model,
    aspectRatio,
    resolution,
    negativePrompt,
    imageNumber,
    referenceImages,
    webhookUrl,
    useWebSearch,
    callerId,
    pollIntervalMs,
    pollTimeoutMs,
    quietPoll,
  } = input;

  const baseBody: Record<string, unknown> = {
    model,
    prompt,
    aspect_ratio: aspectRatio,
    image_number: imageNumber,
  };
  if (negativePrompt) baseBody.negative_prompt = negativePrompt;
  if (resolution) baseBody.resolution = resolution;
  if (webhookUrl) baseBody.webhook_url = webhookUrl;
  if (useWebSearch) baseBody.use_web_search = true;
  if (callerId !== null && callerId !== undefined && Number.isFinite(callerId)) {
    baseBody.caller_id = Math.trunc(callerId);
  }

  const pathSuffix =
    referenceImages.length > 0 ? "/v1/generation/image-to-image" : "/v1/generation/text-to-image";
  const body: Record<string, unknown> =
    referenceImages.length > 0
      ? { ...baseBody, images: referenceImages.map(toImagePayloadValue) }
      : baseBody;

  if (input.dryRun) {
    const preview = { ...body };
    if (Array.isArray(preview.images)) {
      preview.images = (preview.images as string[]).map((v) =>
        typeof v === "string" && v.length > 256 ? `${v.slice(0, 64)}…[${v.length} chars]` : v
      );
    }
    return {
      submit: undefined,
      detail: undefined,
      dryRun: true,
      requestPreview: {
        endpoint: `${BASE_URL}${pathSuffix}`,
        method: "POST",
        body: preview,
      },
    };
  }

  const submit = await apiPost<SubmitResult>(pathSuffix, body);
  const taskIds = submit.data?.task_ids;
  if (!taskIds?.length) throw new Error("No task_ids in response");
  const detail = await pollTask(taskIds[0]!, pollIntervalMs, pollTimeoutMs, quietPoll);
  return { submit: submit.data, detail };
}

function isRetryable(err: unknown): boolean {
  const msg = err instanceof Error ? err.message : String(err);
  const bad = [
    "Reference image not found",
    "Missing ",
    "Invalid ",
    "HTTP 401",
    "HTTP 400",
    "HTTP 403",
    "WeryAI business status",
    "WeryAI task failed",
    "exceeds WeryAI",
    "parameter error",
    "authentication",
    "API status",
    "Parameter",
    "required",
  ];
  return !bad.some((b) => msg.includes(b));
}

async function saveResultImages(
  images: string[],
  outputPath: string,
  noDownload: boolean,
  quietLogs: boolean
): Promise<string[]> {
  if (noDownload || !images.length) return [];
  const userExt = path.extname(path.resolve(outputPath));
  const resolved = normalizeOutputImagePath(outputPath);
  const dir = path.dirname(resolved);
  const fallbackExt = path.extname(resolved) || ".png";
  const base = path.basename(resolved, fallbackExt);
  const saved: string[] = [];
  for (let i = 0; i < images.length; i++) {
    const { bytes, contentType } = await fetchImageBytes(images[i]!, quietLogs);
    let ext = fallbackExt;
    if (!userExt && contentType) {
      ext = extensionForContentType(contentType) ?? fallbackExt;
    }
    const target =
      images.length === 1 ? path.join(dir, `${base}${ext}`) : path.join(dir, `${base}-${i + 1}${ext}`);
    await mkdir(path.dirname(target), { recursive: true });
    await writeFile(target, bytes);
    saved.push(target);
    if (!quietLogs) console.error(`Saved: ${target}`);
  }
  return saved;
}

type GenResult = {
  success: boolean;
  dryRun?: boolean;
  outputPath: string;
  model: string;
  attempts: number;
  error: string | null;
  images?: string[];
  submit?: SubmitResult;
  detail?: TaskDetail;
  requestPreview?: RequestPreview;
};

async function generateOne(
  label: string,
  prompt: string,
  args: CliArgs,
  quietPoll: boolean
): Promise<GenResult> {
  let outputPath = args.imagePath ? normalizeOutputImagePath(args.imagePath) : "";
  const model = resolveModel(args);
  const styledPrompt = applyStylePreset(prompt, args.style);
  const aspectRatio = resolveAspectRatio(args);
  const resolution = resolveResolution(args);

  console.error(`Using model ${model} (${label})`);
  console.error(`Switch model: --model <id> | EXTEND.md default_model | IMAGE_GEN_DEFAULT_MODEL`);

  let attempts = 0;
  while (attempts < MAX_ATTEMPTS) {
    attempts++;
    try {
      assertWeryaiPromptBounds(styledPrompt, args.negativePrompt);
      if (args.referenceImages.length > 0) await validateReferenceImages(args.referenceImages);
      const { submit, detail, dryRun, requestPreview } = await runAsyncGeneration({
        prompt: styledPrompt,
        model,
        aspectRatio,
        resolution,
        negativePrompt: args.negativePrompt,
        imageNumber: args.n,
        referenceImages: args.referenceImages,
        webhookUrl: args.webhookUrl,
        useWebSearch: args.useWebSearch,
        callerId: args.callerId,
        pollIntervalMs: args.pollIntervalMs,
        pollTimeoutMs: args.pollTimeoutMs,
        quietPoll,
        dryRun: args.dryRun,
      });
      if (dryRun) {
        return {
          success: true,
          dryRun: true,
          outputPath: "",
          model,
          attempts,
          error: null,
          submit,
          detail,
          requestPreview,
        };
      }
      if (!detail) {
        throw new Error("Gateway returned no task detail");
      }
      if (detail.task_status === "failed") {
        throw new Error(
          detail.msg ? `WeryAI task failed: ${detail.msg}` : "WeryAI task failed (see data.msg)"
        );
      }
      const images = detail.images ?? [];
      if (outputPath && !args.noDownload) {
        try {
          const savedPaths = await saveResultImages(images, args.imagePath!, args.noDownload, args.json);
          if (savedPaths.length > 0) outputPath = savedPaths[0]!;
        } catch (downloadErr) {
          const msg = downloadErr instanceof Error ? downloadErr.message : String(downloadErr);
          if (!args.json) {
            console.error(`[${label}] Generation succeeded but download failed — URLs still valid; retry download or use --no-download. ${msg}`);
          }
          return {
            success: false,
            outputPath,
            model,
            attempts,
            error: `Download after succeed: ${msg}`,
            images,
            submit,
            detail,
          };
        }
      }
      return {
        success: true,
        outputPath,
        model,
        attempts,
        error: null,
        images,
        submit,
        detail,
      };
    } catch (e) {
      const message = e instanceof Error ? e.message : String(e);
      if (attempts < MAX_ATTEMPTS && isRetryable(e)) {
        console.error(`[${label}] Attempt ${attempts}/${MAX_ATTEMPTS} failed, retrying... (${message})`);
        continue;
      }
      return { success: false, outputPath, model, attempts, error: message };
    }
  }
  return { success: false, outputPath, model, attempts: MAX_ATTEMPTS, error: "Unknown failure" };
}

// --- Batch ---

async function loadBatchTasks(batchFilePath: string): Promise<{
  tasks: BatchTaskInput[];
  jobs: number | null;
  batchDir: string;
}> {
  const resolved = path.resolve(batchFilePath);
  const content = await readFile(resolved, "utf8");
  const parsed = JSON.parse(content.replace(/^\uFEFF/, "")) as BatchFile;
  const batchDir = path.dirname(resolved);
  if (Array.isArray(parsed)) return { tasks: parsed, jobs: null, batchDir };
  if (parsed && typeof parsed === "object" && Array.isArray(parsed.tasks)) {
    const jobs =
      parsed.jobs !== undefined && parsed.jobs !== null
        ? typeof parsed.jobs === "number"
          ? parsed.jobs
          : parseInt(String(parsed.jobs), 10)
        : null;
    if (parsed.jobs !== undefined && parsed.jobs !== null && (jobs === null || jobs < 1)) {
      throw new Error("Invalid batch file: jobs must be a positive integer.");
    }
    return { tasks: parsed.tasks, jobs: jobs && jobs > 0 ? jobs : null, batchDir };
  }
  throw new Error("Invalid batch file. Expected array or { tasks: [...] }.");
}

function resolveBatchPath(batchDir: string, filePath: string): string {
  return path.isAbsolute(filePath) ? filePath : path.resolve(batchDir, filePath);
}

function createTaskArgs(base: CliArgs, task: BatchTaskInput, batchDir: string): CliArgs {
  const refs = task.ref?.map((p) => resolveBatchPath(batchDir, p)) ?? [];
  return {
    ...base,
    prompt: task.prompt ?? null,
    promptFiles: task.promptFiles?.map((f) => resolveBatchPath(batchDir, f)) ?? [],
    imagePath: task.image ? resolveBatchPath(batchDir, task.image) : null,
    model: task.model ?? base.model,
    style: task.style ?? base.style,
    aspectRatio: task.ar ?? base.aspectRatio,
    size: task.size ?? base.size,
    quality: task.quality ?? base.quality,
    imageSize: task.imageSize ?? base.imageSize,
    referenceImages: refs,
    n: task.n ?? base.n,
    webhookUrl: task.webhook_url ?? base.webhookUrl,
    negativePrompt: task.negative_prompt ?? base.negativePrompt,
    resolution: task.resolution ?? base.resolution,
    useWebSearch: task.use_web_search !== undefined && task.use_web_search !== null ? task.use_web_search : base.useWebSearch,
    callerId: task.caller_id !== undefined && task.caller_id !== null ? task.caller_id : base.callerId,
    batchFile: null,
    jobs: base.jobs,
    json: base.json,
    help: false,
  };
}

function getWorkerCount(taskCount: number, jobs: number | null, maxWorkers: number): number {
  const requested = jobs ?? Math.min(taskCount, maxWorkers);
  return Math.max(1, Math.min(requested, taskCount, maxWorkers));
}

function createGate(limit: { concurrency: number; startIntervalMs: number }) {
  const state = { active: 0, lastStartedAt: 0 };
  return async function acquire(): Promise<() => void> {
    while (true) {
      const now = Date.now();
      const cap = state.active < limit.concurrency;
      const gap = now - state.lastStartedAt >= limit.startIntervalMs;
      if (cap && gap) {
        state.active++;
        state.lastStartedAt = now;
        return () => {
          state.active = Math.max(0, state.active - 1);
        };
      }
      await new Promise((r) => setTimeout(r, POLL_WAIT_MS));
    }
  };
}

// --- Legacy subcommand parse (text-to-image / image-to-image) ---

type LegacyOpts = {
  model: string;
  prompt: string;
  aspectRatio: string;
  negativePrompt?: string;
  resolution?: string;
  imageNumber: number;
  webhookUrl?: string;
  output?: string;
  noDownload: boolean;
  pollIntervalMs: number;
  pollTimeoutMs: number;
  json: boolean;
};

function parseLegacyCommon(args: string[], startIdx: number): { opts: LegacyOpts; next: number } {
  const opts: LegacyOpts = {
    model: "",
    prompt: "",
    aspectRatio: "",
    imageNumber: 1,
    noDownload: false,
    pollIntervalMs: 2500,
    pollTimeoutMs: 600_000,
    json: false,
  };
  let i = startIdx;
  const take = (flag: string): string => {
    const v = args[++i];
    if (!v) throw new Error(`Missing value for ${flag}`);
    return v;
  };
  while (i < args.length) {
    const a = args[i]!;
    if (a === "--json") {
      opts.json = true;
      i++;
      continue;
    }
    if (a === "--no-download") {
      opts.noDownload = true;
      i++;
      continue;
    }
    if (a === "--model" || a === "-m") opts.model = take(a);
    else if (a === "--prompt") opts.prompt = take(a);
    else if (a === "--aspect-ratio") opts.aspectRatio = take(a);
    else if (a === "--negative-prompt") opts.negativePrompt = take(a);
    else if (a === "--resolution") opts.resolution = take(a);
    else if (a === "--image-number") {
      opts.imageNumber = parseInt(take(a), 10);
      if (Number.isNaN(opts.imageNumber) || opts.imageNumber < 1) throw new Error("Invalid --image-number");
    } else if (a === "--webhook-url") opts.webhookUrl = take(a);
    else if (a === "--output" || a === "-o") opts.output = take(a);
    else if (a === "--poll-interval-ms") opts.pollIntervalMs = parseInt(take(a), 10);
    else if (a === "--poll-timeout-ms") opts.pollTimeoutMs = parseInt(take(a), 10);
    else if (a.startsWith("-")) throw new Error(`Unknown option: ${a}`);
    else break;
    i++;
  }
  return { opts, next: i };
}

async function runLegacyTextToImage(args: string[]): Promise<void> {
  const extend = await loadExtendConfig();
  const { opts, next } = parseLegacyCommon(args, 0);
  if (next < args.length) throw new Error(`Unexpected arg: ${args[next]}`);
  if (!opts.prompt) throw new Error("--prompt is required");
  if (!opts.aspectRatio) throw new Error("--aspect-ratio is required");
  const cli = parseArgs([]);
  const merged: CliArgs = mergeConfig(
    {
      ...cli,
      prompt: opts.prompt,
      aspectRatio: opts.aspectRatio,
      model: opts.model || null,
      resolution: opts.resolution ?? null,
      negativePrompt: opts.negativePrompt ?? null,
      webhookUrl: opts.webhookUrl ?? null,
      n: opts.imageNumber,
      imagePath: opts.output ?? null,
      noDownload: opts.noDownload,
      pollIntervalMs: opts.pollIntervalMs,
      pollTimeoutMs: opts.pollTimeoutMs,
      json: opts.json,
      quality: opts.resolution ? null : "2k",
    },
    extend
  );
  const r = await generateOne("legacy-text", opts.prompt, merged, opts.json);
  if (!r.success) throw new Error(r.error ?? "failed");
  if (opts.json)
    console.log(JSON.stringify({ submit: r.submit, detail: r.detail, images: r.images }, null, 2));
  else console.log((r.images ?? []).join("\n"));
}

async function runLegacyImageToImage(args: string[]): Promise<void> {
  const extend = await loadExtendConfig();
  const refs: string[] = [];
  const filtered: string[] = [];
  for (let i = 0; i < args.length; i++) {
    if (args[i] === "--ref") {
      const v = args[++i];
      if (!v) throw new Error("Missing value for --ref");
      refs.push(v);
    } else filtered.push(args[i]!);
  }
  const { opts, next } = parseLegacyCommon(filtered, 0);
  if (next < filtered.length) throw new Error(`Unexpected arg: ${filtered[next]}`);
  if (!opts.prompt) throw new Error("--prompt is required");
  if (!opts.aspectRatio) throw new Error("--aspect-ratio is required");
  if (!refs.length) throw new Error("At least one --ref required");
  const cli = parseArgs([]);
  const merged: CliArgs = mergeConfig(
    {
      ...cli,
      prompt: opts.prompt,
      aspectRatio: opts.aspectRatio,
      model: opts.model || null,
      resolution: opts.resolution ?? null,
      negativePrompt: opts.negativePrompt ?? null,
      webhookUrl: opts.webhookUrl ?? null,
      n: opts.imageNumber,
      imagePath: opts.output ?? null,
      referenceImages: refs,
      noDownload: opts.noDownload,
      pollIntervalMs: opts.pollIntervalMs,
      pollTimeoutMs: opts.pollTimeoutMs,
      json: opts.json,
      quality: opts.resolution ? null : "2k",
    },
    extend
  );
  const r = await generateOne("legacy-i2i", opts.prompt, merged, opts.json);
  if (!r.success) throw new Error(r.error ?? "failed");
  if (opts.json)
    console.log(JSON.stringify({ submit: r.submit, detail: r.detail, images: r.images }, null, 2));
  else console.log((r.images ?? []).join("\n"));
}

async function runLegacyStatus(args: string[]): Promise<void> {
  let taskId = "";
  for (let i = 0; i < args.length; i++) {
    if (args[i] === "--task-id") taskId = args[++i] ?? "";
  }
  if (!taskId) throw new Error("--task-id is required");
  const r = await apiGet<TaskDetail>(`/v1/generation/${encodeURIComponent(taskId)}/status`);
  console.log(JSON.stringify(r, null, 2));
}

// --- main flows ---

async function runSingleMode(args: CliArgs): Promise<void> {
  let prompt: string | null = args.prompt;
  if (!prompt && args.promptFiles.length > 0) {
    prompt = await readPromptFromFiles(args.promptFiles);
  }
  if (!prompt) prompt = await readPromptFromStdin();
  if (!prompt) throw new Error("Prompt is required (--prompt, --promptfiles, or stdin)");
  if (!args.imagePath && !args.dryRun) throw new Error("--image <path> is required for single mode");
  if (args.referenceImages.length > 0) await validateReferenceImages(args.referenceImages);

  const r = await generateOne("single", prompt, args, args.json);
  if (!r.success) throw new Error(r.error ?? "Generation failed");
  if (args.dryRun) {
    console.log(
      JSON.stringify(
        {
          mode: "dry-run",
          model: r.model,
          prompt: prompt.slice(0, 200),
          request: r.requestPreview,
        },
        null,
        2
      )
    );
    return;
  }
  if (args.json) {
    console.log(
      JSON.stringify(
        {
          savedHint: r.outputPath,
          model: r.model,
          attempts: r.attempts,
          images: r.images,
          prompt: prompt.slice(0, 200),
        },
        null,
        2
      )
    );
    return;
  }
  console.log(r.outputPath);
}

async function runBatchMode(args: CliArgs, extend: ExtendConfig): Promise<void> {
  if (!args.batchFile) throw new Error("--batchfile required");
  const { tasks: inputs, jobs: batchJobs, batchDir } = await loadBatchTasks(args.batchFile);
  if (!inputs.length) throw new Error("Batch has no tasks");

  const maxW = getMaxWorkers(extend);
  const workerCount = getWorkerCount(inputs.length, args.jobs ?? batchJobs, maxW);
  const gate = createGate(getConcurrency());

  type Row = {
    id: string;
    success: boolean;
    dryRun?: boolean;
    outputPath: string;
    model: string;
    attempts: number;
    error: string | null;
    requestPreview?: RequestPreview;
  };
  const results: Row[] = new Array(inputs.length);

  if (inputs.length === 1) {
    const t = inputs[0]!;
    if (t.provider && !args.json)
      console.error(`[batch] Ignoring task.provider=${JSON.stringify(t.provider)} (single-gateway mode)`);
    const taskArgs = createTaskArgs(args, t, batchDir);
    const merged = mergeConfig(taskArgs, extend);
    if (!merged.quality) merged.quality = "2k";
    let prompt: string | null = taskArgs.prompt;
    if (!prompt && taskArgs.promptFiles.length)
      prompt = await readPromptFromFiles(taskArgs.promptFiles);
    if (!prompt) throw new Error("Task missing prompt");
    if (!taskArgs.imagePath && !args.dryRun) throw new Error("Task missing image path");
    const r = await generateOne(t.id || "task-01", prompt, merged, args.json);
    results[0] = {
      id: t.id || "task-01",
      success: r.success,
      dryRun: r.dryRun,
      outputPath: r.outputPath,
      model: r.model,
      attempts: r.attempts,
      error: r.error,
      requestPreview: r.requestPreview,
    };
  } else {
    console.error(`Batch: ${inputs.length} tasks, ${workerCount} workers.`);
    let next = 0;
    const workers = Array.from({ length: workerCount }, async () => {
      while (true) {
        const idx = next++;
        if (idx >= inputs.length) return;
        const t = inputs[idx]!;
        if (t.provider && !args.json)
          console.error(`[batch] Ignoring task.provider=${JSON.stringify(t.provider)} (single-gateway mode)`);
        const taskArgs = createTaskArgs(args, t, batchDir);
        const merged = mergeConfig(taskArgs, extend);
        if (!merged.quality) merged.quality = "2k";
        let prompt: string | null = taskArgs.prompt;
        if (!prompt && taskArgs.promptFiles.length) {
          try {
            prompt = await readPromptFromFiles(taskArgs.promptFiles);
          } catch (e) {
            results[idx] = {
              id: t.id || `task-${String(idx + 1).padStart(2, "0")}`,
              success: false,
              outputPath: "",
              model: modelLabelForReport(merged),
              attempts: 0,
              error: (e as Error).message,
            };
            continue;
          }
        }
        if (!prompt) {
          results[idx] = {
            id: t.id || `task-${String(idx + 1).padStart(2, "0")}`,
            success: false,
            outputPath: "",
            model: modelLabelForReport(merged),
            attempts: 0,
            error: "Missing prompt",
          };
          continue;
        }
        if (!taskArgs.imagePath && !args.dryRun) {
          results[idx] = {
            id: t.id || `task-${String(idx + 1).padStart(2, "0")}`,
            success: false,
            outputPath: "",
            model: modelLabelForReport(merged),
            attempts: 0,
            error: "Missing image output path",
          };
          continue;
        }
        const release = await gate();
        try {
          const r = await generateOne(t.id || `task-${String(idx + 1).padStart(2, "0")}`, prompt, merged, args.json);
          results[idx] = {
            id: t.id || `task-${String(idx + 1).padStart(2, "0")}`,
            success: r.success,
            dryRun: r.dryRun,
            outputPath: r.outputPath,
            model: r.model,
            attempts: r.attempts,
            error: r.error,
            requestPreview: r.requestPreview,
          };
        } finally {
          release();
        }
      }
    });
    await Promise.all(workers);
  }

  if (args.dryRun) {
    const failed = results.filter((x) => !x.success).length;
    if (args.json) {
      console.log(
        JSON.stringify(
          {
            mode: "batch-dry-run",
            total: results.length,
            previewed: results.length - failed,
            failed,
            results,
          },
          null,
          2
        )
      );
    } else {
      for (const row of results) {
        if (row.requestPreview) {
          console.log(JSON.stringify({ id: row.id, request: row.requestPreview }, null, 2));
        } else {
          console.error(`[dry-run] ${row.id}: ${row.error}`);
        }
      }
    }
    if (failed > 0) process.exitCode = 1;
    return;
  }

  const ok = results.filter((x) => x?.success).length;
  const fail = results.length - ok;
  console.error("");
  console.error("Batch generation summary:");
  console.error(`- Total: ${results.length}`);
  console.error(`- Succeeded: ${ok}`);
  console.error(`- Failed: ${fail}`);
  if (fail > 0) {
    console.error("Failures:");
    for (const x of results) {
      if (!x.success) console.error(`- ${x.id}: ${x.error}`);
    }
  }
  if (args.json) {
    console.log(JSON.stringify({ mode: "batch", total: results.length, succeeded: ok, failed: fail, results }, null, 2));
  }
  if (fail > 0) process.exitCode = 1;
}

async function main(): Promise<void> {
  const argv = process.argv.slice(2);
  if (!argv.length || argv[0] === "-h" || argv[0] === "--help") {
    printUsage();
    return;
  }

  await loadEnv();

  if (argv[0] === "text-to-image") {
    try {
      await runLegacyTextToImage(argv.slice(1));
    } catch (e) {
      console.error((e as Error).message);
      process.exit(1);
    }
    return;
  }
  if (argv[0] === "image-to-image") {
    try {
      await runLegacyImageToImage(argv.slice(1));
    } catch (e) {
      console.error((e as Error).message);
      process.exit(1);
    }
    return;
  }
  if (argv[0] === "status") {
    try {
      await runLegacyStatus(argv.slice(1));
    } catch (e) {
      console.error((e as Error).message);
      process.exit(1);
    }
    return;
  }

  const extend = await loadExtendConfig();
  let args = parseArgs(argv);
  if (args.help) {
    printUsage();
    return;
  }
  args = mergeConfig(args, extend);
  if (!args.quality) args.quality = "2k";

  try {
    if (args.batchFile) await runBatchMode(args, extend);
    else await runSingleMode(args);
  } catch (e) {
    console.error((e as Error).message);
    process.exit(1);
  }
}

main();
