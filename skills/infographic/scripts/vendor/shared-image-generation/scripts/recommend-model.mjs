import fs from "node:fs";
import path from "node:path";
import process from "node:process";
import { fileURLToPath } from "node:url";
import { discoverImageModels } from "./discover-models.mjs";

const homeDir = process.env.HOME ?? "";
const skillDir = path.resolve(path.dirname(fileURLToPath(import.meta.url)), "..");
const bundledRegistryPath = path.join(skillDir, "references", "config", "starter-models.json");
const skillNamespace = process.env.IMAGE_SKILL_NAMESPACE?.trim() || "image-generation";
const VALID_WORKFLOWS = new Set(["general", "cover", "rednote", "infographic", "comic", "article"]);
const VALID_QUALITIES = new Set(["fast", "balanced", "high"]);
const AUTO_DISCOVER = Symbol("auto-discover");
const ROLE_PROFILES = {
  "cover-hero": {
    workflow: "cover",
    needRef: false,
    textHeavy: false,
    preferredQuality: "high",
    preferredStyles: ["editorial", "cinematic", "photoreal"],
  },
  "cover-typography-space": {
    workflow: "cover",
    needRef: false,
    textHeavy: true,
    preferredQuality: "high",
    preferredStyles: ["editorial", "poster"],
  },
  "rednote-cover": {
    workflow: "rednote",
    needRef: true,
    textHeavy: true,
    preferredQuality: "balanced",
    preferredStyles: ["editorial", "flat-illustration", "anime"],
  },
  "rednote-text-card": {
    workflow: "rednote",
    needRef: true,
    textHeavy: true,
    preferredQuality: "balanced",
    preferredStyles: ["flat-illustration", "editorial", "infographic"],
  },
  "infographic-dense": {
    workflow: "infographic",
    needRef: false,
    textHeavy: true,
    preferredQuality: "high",
    preferredStyles: ["diagram", "infographic"],
  },
  "infographic-minimal": {
    workflow: "infographic",
    needRef: false,
    textHeavy: false,
    preferredQuality: "balanced",
    preferredStyles: ["editorial", "infographic"],
  },
  "comic-cover": {
    workflow: "comic",
    needRef: true,
    textHeavy: false,
    preferredQuality: "high",
    preferredStyles: ["comic", "manga", "cinematic"],
  },
  "comic-character-sheet": {
    workflow: "comic",
    needRef: true,
    textHeavy: false,
    preferredQuality: "high",
    preferredStyles: ["comic", "manga", "flat-illustration"],
  },
  "comic-page": {
    workflow: "comic",
    needRef: true,
    textHeavy: false,
    preferredQuality: "high",
    preferredStyles: ["comic", "manga", "flat-illustration"],
  },
  "article-framework": {
    workflow: "article",
    needRef: true,
    textHeavy: true,
    preferredQuality: "balanced",
    preferredStyles: ["diagram", "infographic", "editorial"],
  },
  "article-scene": {
    workflow: "article",
    needRef: true,
    textHeavy: false,
    preferredQuality: "high",
    preferredStyles: ["editorial", "watercolor", "photoreal"],
  },
  "article-flowchart": {
    workflow: "article",
    needRef: true,
    textHeavy: true,
    preferredQuality: "balanced",
    preferredStyles: ["diagram", "infographic", "flat-illustration"],
  },
};
const VALID_ROLES = new Set(Object.keys(ROLE_PROFILES));

function defaultArgs() {
  return {
    workflow: "general",
    role: null,
    needRef: false,
    textHeavy: false,
    quality: "balanced",
    style: null,
    modality: null,
    registryPath: null,
    json: false,
  };
}

function parseArgs(argv) {
  const args = defaultArgs();
  for (let i = 0; i < argv.length; i++) {
    const current = argv[i];
    if (current === "--workflow") args.workflow = argv[++i] ?? args.workflow;
    else if (current === "--role") args.role = argv[++i] ?? null;
    else if (current === "--need-ref") args.needRef = true;
    else if (current === "--text-heavy") args.textHeavy = true;
    else if (current === "--quality") args.quality = argv[++i] ?? args.quality;
    else if (current === "--style") args.style = argv[++i] ?? null;
    else if (current === "--modality") args.modality = argv[++i] ?? null;
    else if (current === "--registry") args.registryPath = argv[++i] ?? null;
    else if (current === "--json") args.json = true;
    else if (current === "-h" || current === "--help") {
      printHelp();
      process.exit(0);
    }
  }

  if (!VALID_WORKFLOWS.has(args.workflow)) {
    throw new Error(`Invalid workflow: ${args.workflow}`);
  }
  if (args.role && !VALID_ROLES.has(args.role)) {
    throw new Error(`Invalid role: ${args.role}`);
  }
  if (!VALID_QUALITIES.has(args.quality)) {
    throw new Error(`Invalid quality: ${args.quality}`);
  }
  if (!args.modality) {
    args.modality = args.needRef ? "image-to-image" : "text-to-image";
  }
  return args;
}

function printHelp() {
  console.log(`Usage:
  node scripts/recommend-model.mjs --workflow <general|cover|rednote|infographic|comic|article> [options]

Options:
  --workflow <name>   Target workflow (default: general)
  --role <name>       Optional sub-task role such as cover-hero, comic-page, infographic-dense
  --need-ref          Prefer models that support reference images
  --text-heavy        Prefer models that handle text-heavy or label-heavy images better
  --quality <level>   fast | balanced | high (default: balanced)
  --style <name>      Optional style hint such as editorial, diagram, comic, photoreal
  --modality <name>   text-to-image | image-to-image (default inferred from --need-ref)
  --registry <path>   Explicit MODELS.json path
  --json              Print JSON output
  -h, --help          Show help`);
}

export function getRoleProfile(role) {
  return role ? ROLE_PROFILES[role] ?? null : null;
}

export function resolveRecommendationOptions(args) {
  const roleProfile = getRoleProfile(args.role);
  const preferredStyles = new Set();
  if (args.style) preferredStyles.add(args.style);
  for (const style of roleProfile?.preferredStyles ?? []) preferredStyles.add(style);

  return {
    ...args,
    workflow: args.workflow === "general" && roleProfile?.workflow ? roleProfile.workflow : args.workflow,
    needRef: Boolean(args.needRef || roleProfile?.needRef),
    textHeavy: Boolean(args.textHeavy || roleProfile?.textHeavy),
    preferredQuality: roleProfile?.preferredQuality ?? null,
    preferredStyles: [...preferredStyles],
    roleProfile,
  };
}

export function parseSimpleFrontmatter(markdown) {
  const match = markdown.match(/^---\n([\s\S]*?)\n---/);
  if (!match) return {};
  const out = {};
  for (const line of match[1].split("\n")) {
    const idx = line.indexOf(":");
    if (idx === -1) continue;
    const key = line.slice(0, idx).trim();
    const value = line.slice(idx + 1).trim().replace(/^"(.*)"$/, "$1");
    out[key] = value;
  }
  return out;
}

export function candidateRegistryPaths(cwd, overridePath = null) {
  if (overridePath) return [path.resolve(cwd, overridePath)];
  return [
    path.join(cwd, ".image-skills", skillNamespace, "MODELS.json"),
    homeDir ? path.join(homeDir, ".image-skills", skillNamespace, "MODELS.json") : null,
    bundledRegistryPath,
  ].filter(Boolean);
}

function candidateExtendPaths(cwd) {
  return [
    path.join(cwd, ".image-skills", skillNamespace, "EXTEND.md"),
    homeDir ? path.join(homeDir, ".image-skills", skillNamespace, "EXTEND.md") : null,
  ].filter(Boolean);
}

function readFirstExisting(paths) {
  for (const filePath of paths) {
    if (fs.existsSync(filePath)) return filePath;
  }
  return null;
}

export function loadDefaultModel(cwd) {
  const envModel = process.env.IMAGE_GEN_DEFAULT_MODEL?.trim();
  if (envModel) return { model: envModel, source: "IMAGE_GEN_DEFAULT_MODEL" };

  const extendPath = readFirstExisting(candidateExtendPaths(cwd));
  if (!extendPath) return { model: null, source: null };
  const parsed = parseSimpleFrontmatter(fs.readFileSync(extendPath, "utf8"));
  const model = typeof parsed.default_model === "string" ? parsed.default_model.trim() : "";
  if (!model) return { model: null, source: null };
  return { model, source: path.relative(cwd, extendPath) || "EXTEND.md" };
}

export function loadRegistry(cwd, overridePath = null) {
  const registryPath = readFirstExisting(candidateRegistryPaths(cwd, overridePath));
  if (!registryPath) {
    return { registry: null, path: null };
  }
  const parsed = JSON.parse(fs.readFileSync(registryPath, "utf8"));
  if (!parsed || !Array.isArray(parsed.models)) {
    throw new Error(`Invalid model registry: ${registryPath}`);
  }
  return { registry: parsed, path: registryPath };
}

function registrySourceForPath(registryPath) {
  if (!registryPath) return null;
  return path.resolve(registryPath) === path.resolve(bundledRegistryPath) ? "bundled-starter" : "local-registry";
}

export function scoreModel(model, options, defaultModel) {
  const reasons = [];
  let score = 0;

  const modalities = Array.isArray(model.modalities) ? model.modalities : [];
  if (options.modality && modalities.length > 0) {
    if (modalities.includes(options.modality)) {
      score += 3;
      reasons.push(`supports ${options.modality}`);
    } else {
      score -= 6;
      reasons.push(`does not list ${options.modality}`);
    }
  }

  if (options.workflow && Array.isArray(model.best_for)) {
    if (model.best_for.includes(options.workflow)) {
      score += 5;
      reasons.push(`matched workflow ${options.workflow}`);
    } else if (options.workflow !== "general") {
      score -= 1;
    }
  }

  if (options.role) {
    if (Array.isArray(model.best_for_roles) && model.best_for_roles.includes(options.role)) {
      score += 6;
      reasons.push(`matched role ${options.role}`);
    } else if (options.role) {
      score -= 1;
    }
  }

  if (options.needRef) {
    if (model.supports_ref) {
      score += 4;
      reasons.push("supports reference images");
    } else {
      score -= 8;
      reasons.push("missing reference-image support");
    }
  }

  if (options.textHeavy) {
    if (model.supports_text_heavy) {
      score += 4;
      reasons.push("handles text-heavy layouts");
    } else {
      score -= 3;
      reasons.push("not marked for text-heavy layouts");
    }
  }

  if (options.quality && model.quality_tier === options.quality) {
    score += 3;
    reasons.push(`matched quality tier ${options.quality}`);
  }

  if (Array.isArray(model.strengths)) {
    for (const style of options.preferredStyles ?? []) {
      if (model.strengths.includes(style)) {
        score += style === options.style ? 3 : 2;
        reasons.push(`matched style strength ${style}`);
      }
    }
  }

  if (options.preferredQuality && model.quality_tier === options.preferredQuality && options.quality !== options.preferredQuality) {
    score += 1;
    reasons.push(`fits preferred role quality ${options.preferredQuality}`);
  }

  if (defaultModel && model.id === defaultModel) {
    score += 2;
    reasons.push("matches configured default model");
  }

  return { score, reasons };
}

export function recommendModels(registry, options, defaultModel = null) {
  const resolvedOptions = resolveRecommendationOptions(options);
  return registry.models
    .map((model) => {
      const scored = scoreModel(model, resolvedOptions, defaultModel);
      return {
        id: model.id,
        label: model.label ?? model.id,
        score: scored.score,
        reasons: scored.reasons,
        quality_tier: model.quality_tier ?? null,
        supports_ref: Boolean(model.supports_ref),
        supports_text_heavy: Boolean(model.supports_text_heavy),
        best_for: Array.isArray(model.best_for) ? model.best_for : [],
        best_for_roles: Array.isArray(model.best_for_roles) ? model.best_for_roles : [],
        notes: model.notes ?? null,
      };
    })
    .sort((a, b) => b.score - a.score || a.id.localeCompare(b.id));
}

export async function buildRecommendationSummary({ cwd, args, autoDiscoveredRegistry = AUTO_DISCOVER }) {
  const resolvedArgs = resolveRecommendationOptions(args);
  const { registry, path: registryPath } = loadRegistry(cwd, args.registryPath);
  const { model: defaultModel, source: defaultSource } = loadDefaultModel(cwd);
  const registrySource = registrySourceForPath(registryPath);
  const discoveredRegistry = registry
    ? null
    : autoDiscoveredRegistry === AUTO_DISCOVER
      ? await discoverImageModels().catch(() => null)
      : autoDiscoveredRegistry;

  if (!registry && !discoveredRegistry) {
    return {
      ok: false,
      workflow: resolvedArgs.workflow,
      requestedWorkflow: args.workflow,
      role: args.role,
      registryPath: null,
      defaultModel,
      defaultModelSource: defaultSource,
      recommendations: defaultModel
        ? [
            {
              id: defaultModel,
              label: defaultModel,
              score: 1,
              reasons: ["no model registry found; falling back to configured default model"],
            },
          ]
        : [],
      notes: [
        "No MODELS.json registry was found.",
        `Create .image-skills/${skillNamespace}/MODELS.json or run \`npm run discover-image-models -- --out .image-skills/${skillNamespace}/MODELS.json\` to let this skill rank models by capability.`,
      ],
    };
  }

  const activeRegistry = registry ?? discoveredRegistry;
  const activeRegistryPath = registryPath ?? "auto-discovered from WeryAI docs";

  return {
    ok: true,
    workflow: resolvedArgs.workflow,
    requestedWorkflow: args.workflow,
    role: args.role,
    registryPath: activeRegistryPath,
    defaultModel,
    defaultModelSource: defaultSource,
    registrySource: registry ? registrySource : "auto-discovered",
    resolvedNeeds: {
      needRef: resolvedArgs.needRef,
      textHeavy: resolvedArgs.textHeavy,
      preferredQuality: resolvedArgs.preferredQuality,
      preferredStyles: resolvedArgs.preferredStyles,
    },
    recommendations: recommendModels(activeRegistry, resolvedArgs, defaultModel).slice(0, 5),
    notes: [
      resolvedArgs.needRef ? "Reference-image support was weighted heavily." : "Reference-image support was not required.",
      resolvedArgs.textHeavy ? "Text-heavy handling was weighted heavily." : "Text-heavy handling was not required.",
      args.role ? `Role profile ${args.role} refined the recommendation.` : "No sub-task role profile was applied.",
      registrySource === "bundled-starter"
        ? "Using the bundled WeryAI starter model registry shipped with this skill."
        : registry
          ? "Using local MODELS.json registry."
          : "Fell back to auto-discovered WeryAI image model list.",
    ],
  };
}

function formatText(summary) {
  const lines = [];
  lines.push(`Workflow: ${summary.workflow}`);
  if (summary.role) lines.push(`Role: ${summary.role}`);
  if (summary.registryPath) lines.push(`Registry: ${summary.registryPath}`);
  if (summary.defaultModel) lines.push(`Configured default: ${summary.defaultModel} (${summary.defaultModelSource})`);
  if (!summary.recommendations.length) {
    lines.push("No concrete model recommendation is available yet.");
  } else {
    lines.push("Recommended models:");
    for (const item of summary.recommendations) {
      lines.push(`- ${item.id} (score ${item.score})`);
      if (item.reasons?.length) lines.push(`  reasons: ${item.reasons.join("; ")}`);
    }
  }
  for (const note of summary.notes ?? []) lines.push(`note: ${note}`);
  return lines.join("\n");
}

const isDirectRun =
  process.argv[1] && path.resolve(process.argv[1]) === path.resolve(fileURLToPath(import.meta.url));

if (isDirectRun) {
  try {
    const cwd = process.cwd();
    const args = parseArgs(process.argv.slice(2));
    const summary = await buildRecommendationSummary({ cwd, args });
    if (args.json) {
      console.log(JSON.stringify(summary, null, 2));
    } else {
      console.log(formatText(summary));
    }
  } catch (error) {
    console.error(error instanceof Error ? error.message : String(error));
    process.exit(1);
  }
}
