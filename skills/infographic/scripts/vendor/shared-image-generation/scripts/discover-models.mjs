import fs from "node:fs";
import path from "node:path";
import process from "node:process";
import { fileURLToPath } from "node:url";

// Prefer raw markdown endpoints (if exposed), but fall back to HTML pages.
const TEXT_TO_IMAGE_DOC_URLS = [
  "https://docs.weryai.com/api-reference/image-generation/submit-text-to-image-task.md",
  "https://docs.weryai.com/api-reference/image-generation/submit-text-to-image-task",
];
const IMAGE_TO_IMAGE_DOC_URLS = [
  "https://docs.weryai.com/api-reference/image-generation/submit-image-to-image-task.md",
  "https://docs.weryai.com/api-reference/image-generation/submit-image-to-image-task",
];

function printHelp() {
  console.log(`Usage:
  node scripts/discover-models.mjs [options]

Options:
  --out <path>   Write discovered MODELS.json starter to this path
  --force        Overwrite existing output file
  --json         Print JSON output
  -h, --help     Show help`);
}

function parseArgs(argv) {
  const args = {
    out: null,
    force: false,
    json: false,
  };

  for (let i = 0; i < argv.length; i++) {
    const current = argv[i];
    if (current === "--out") args.out = argv[++i] ?? null;
    else if (current === "--force") args.force = true;
    else if (current === "--json") args.json = true;
    else if (current === "--help" || current === "-h") {
      printHelp();
      process.exit(0);
    }
  }
  return args;
}

export function parseModelTable(markdown) {
  const lines = markdown.replace(/\r\n/g, "\n").split("\n");

  // Locate a markdown table containing a model list. Docs may tweak header wording,
  // so we infer the "key" column and the "name" column rather than matching a fixed string.
  const normalize = (s) => s.toLowerCase().replace(/\s+/g, " ").trim();
  const isTableLine = (line) => /^\s*\|/.test(line);
  const splitRow = (line) => line.split("|").map((part) => part.trim()).filter(Boolean);

  let headerIndex = -1;
  for (let i = 0; i < lines.length; i++) {
    const line = lines[i];
    if (!isTableLine(line)) continue;
    const parts = splitRow(line);
    if (parts.length < 2) continue;
    const joined = normalize(parts.join(" | "));
    if (joined.includes("model") && joined.includes("key")) {
      headerIndex = i;
      break;
    }
  }
  if (headerIndex === -1) return [];

  const headerParts = splitRow(lines[headerIndex] ?? "").map(normalize);
  const keyIdx = headerParts.findIndex((h) => h === "model key" || h.endsWith(" key") || h.includes("key"));
  const nameIdx = headerParts.findIndex((h) => h === "model" || h === "model name" || h.includes("name"));
  const resolvedKeyIdx = keyIdx >= 0 ? keyIdx : 1;
  const resolvedNameIdx = nameIdx >= 0 ? nameIdx : 0;

  const rows = [];
  for (let i = headerIndex + 2; i < lines.length; i++) {
    const line = lines[i].trim();
    if (!line.startsWith("|")) break;
    const parts = line.split("|").map((part) => part.trim()).filter(Boolean);
    if (parts.length < 2) continue;
    const label = parts[resolvedNameIdx] ?? parts[0];
    const id = parts[resolvedKeyIdx] ?? parts[1];
    if (!label || !id) continue;
    rows.push({ label, id });
  }
  return rows;
}

function inferQualityTier(label, id) {
  const value = `${label} ${id}`.toLowerCase();
  if (/(mini|lite|flash)/i.test(value)) return "fast";
  if (/(pro|2\.0|5\.0|v4 pro)/i.test(value)) return "high";
  return "balanced";
}

function inferStrengths(label, id) {
  const value = `${label} ${id}`.toLowerCase();
  const strengths = new Set();
  if (/(vector|recraft)/i.test(value)) {
    strengths.add("flat-illustration");
    strengths.add("infographic");
    strengths.add("diagram");
  }
  if (/(qwen|gpt)/i.test(value)) {
    strengths.add("editorial");
    strengths.add("labels");
  }
  if (/(dreamina|flux|imagen|weryai image|seedream|wan)/i.test(value)) {
    strengths.add("editorial");
    strengths.add("photoreal");
    strengths.add("cover");
  }
  if (/(gemini|banana)/i.test(value)) {
    strengths.add("editorial");
    strengths.add("speed");
  }
  if (strengths.size === 0) strengths.add("editorial");
  return [...strengths];
}

function inferSupportsTextHeavy(label, id) {
  return /(vector|recraft|qwen|gpt|labels)/i.test(`${label} ${id}`);
}

function inferBestFor(strengths, supportsTextHeavy) {
  const bestFor = new Set(["general", "cover"]);
  if (supportsTextHeavy || strengths.includes("infographic") || strengths.includes("diagram")) {
    bestFor.add("infographic");
    bestFor.add("article");
    bestFor.add("rednote");
  }
  return [...bestFor];
}

function inferBestForRoles(bestFor, strengths, supportsRef, supportsTextHeavy) {
  const roles = new Set();
  if (bestFor.includes("cover")) roles.add("cover-hero");
  if (supportsTextHeavy) {
    roles.add("rednote-text-card");
    roles.add("infographic-dense");
    roles.add("article-framework");
    roles.add("article-flowchart");
  }
  if (strengths.includes("infographic")) roles.add("infographic-dense");
  if (supportsRef) {
    roles.add("comic-page");
    roles.add("comic-character-sheet");
  }
  return [...roles];
}

export function buildDiscoveredRegistry({ textToImageModels, imageToImageModels }) {
  const byId = new Map();

  for (const item of textToImageModels) {
    const current = byId.get(item.id) ?? {
      id: item.id,
      label: item.label,
      modalities: [],
    };
    current.modalities = [...new Set([...current.modalities, "text-to-image"])];
    byId.set(item.id, current);
  }

  for (const item of imageToImageModels) {
    const current = byId.get(item.id) ?? {
      id: item.id,
      label: item.label,
      modalities: [],
    };
    current.modalities = [...new Set([...current.modalities, "image-to-image"])];
    byId.set(item.id, current);
  }

  const models = [...byId.values()]
    .map((model) => {
      const strengths = inferStrengths(model.label, model.id);
      const supportsTextHeavy = inferSupportsTextHeavy(model.label, model.id);
      const bestFor = inferBestFor(strengths, supportsTextHeavy);
      return {
        id: model.id,
        label: model.label,
        modalities: model.modalities,
        supports_ref: model.modalities.includes("image-to-image"),
        supports_text_heavy: supportsTextHeavy,
        quality_tier: inferQualityTier(model.label, model.id),
        strengths,
        best_for: bestFor,
        best_for_roles: inferBestForRoles(bestFor, strengths, model.modalities.includes("image-to-image"), supportsTextHeavy),
        notes: "Auto-discovered from current WeryAI image generation documentation.",
      };
    })
    .sort((a, b) => a.id.localeCompare(b.id));

  return {
    version: 1,
    generated_at: new Date().toISOString(),
    source: "weryai-docs-auto-discovery",
    docs: {
      text_to_image: TEXT_TO_IMAGE_DOC_URLS[0],
      image_to_image: IMAGE_TO_IMAGE_DOC_URLS[0],
    },
    models,
  };
}

async function fetchFirstOk(urls, fetchImpl) {
  let last = null;
  for (const url of urls) {
    try {
      const res = await fetchImpl(url);
      if (res.ok) return { res, url };
      last = new Error(`HTTP ${res.status} for ${url}`);
    } catch (e) {
      last = e instanceof Error ? e : new Error(String(e));
    }
  }
  throw last ?? new Error("Failed to fetch docs");
}

export async function discoverImageModels({ fetchImpl = fetch } = {}) {
  const [text, image] = await Promise.all([
    fetchFirstOk(TEXT_TO_IMAGE_DOC_URLS, fetchImpl),
    fetchFirstOk(IMAGE_TO_IMAGE_DOC_URLS, fetchImpl),
  ]);

  const [textMarkdown, imageMarkdown] = await Promise.all([text.res.text(), image.res.text()]);
  const textToImageModels = parseModelTable(textMarkdown);
  const imageToImageModels = parseModelTable(imageMarkdown);
  if (textToImageModels.length === 0 && imageToImageModels.length === 0) {
    throw new Error("No image model table was found in current WeryAI docs");
  }

  return buildDiscoveredRegistry({ textToImageModels, imageToImageModels });
}

function writeRegistry(filePath, registry) {
  const resolved = path.resolve(filePath);
  fs.mkdirSync(path.dirname(resolved), { recursive: true });
  fs.writeFileSync(resolved, `${JSON.stringify(registry, null, 2)}\n`, "utf8");
  return resolved;
}

const isDirectRun =
  process.argv[1] && path.resolve(process.argv[1]) === path.resolve(fileURLToPath(import.meta.url));

if (isDirectRun) {
  try {
    const args = parseArgs(process.argv.slice(2));
    const registry = await discoverImageModels();

    let outputPath = null;
    let wrote = false;
    if (args.out) {
      const resolvedOut = path.resolve(args.out);
      if (fs.existsSync(resolvedOut) && !args.force) {
        outputPath = resolvedOut;
        wrote = false;
      } else {
        outputPath = writeRegistry(resolvedOut, registry);
        wrote = true;
      }
    }

    const summary = {
      ok: true,
      modelCount: registry.models.length,
      source: registry.source,
      docs: registry.docs,
      outputPath,
      wrote,
      registry,
    };
    if (args.json) console.log(JSON.stringify(summary, null, 2));
    else {
      console.log(`Discovered ${summary.modelCount} image models from WeryAI docs.`);
      if (args.out && summary.outputPath) {
        if (!summary.wrote) {
          console.log(`Registry already exists (skipped): ${summary.outputPath}`);
          console.log("Tip: pass --force to overwrite.");
        } else {
          console.log(`Wrote starter registry: ${summary.outputPath}`);
        }
      }
    }
  } catch (error) {
    console.error(error instanceof Error ? error.message : String(error));
    process.exit(1);
  }
}
