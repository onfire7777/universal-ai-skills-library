
import fs from "node:fs";
import path from "node:path";
import process from "node:process";
import { fileURLToPath } from "node:url";

const defaultSkillDir = path.resolve(path.dirname(fileURLToPath(import.meta.url)), "..");
const bundledModelsPath = path.join(defaultSkillDir, "references", "config", "starter-models.json");
const skillNamespace = process.env.IMAGE_SKILL_NAMESPACE?.trim() || "image-generation";
const defaultConfigDir = `.image-skills/${skillNamespace}`;
const defaultEnvPath = `${defaultConfigDir}/.env`;
const defaultExtendPath = `${defaultConfigDir}/EXTEND.md`;
const defaultStarterModel = {
  id: "GEMINI_3_1_FLASH_IMAGE",
  label: "Nano Banana 2",
};

function printHelp() {
  console.log(`Usage:
  node scripts/setup.mjs [options]

Options:
  --project <path>      Project directory to write/read .image-skills config (default: cwd)
  --workflow <name>     Optional focus: general|cover|rednote|infographic|comic|article|compress
  --force-models        Overwrite MODELS.json if it already exists
  --api-key <secret>    Write IMAGE_GEN_API_KEY into local skill config
  --persist-api-key     Persist current IMAGE_GEN_API_KEY from env into local skill config
  --scope <target>      Where to write the local .env: project or home (default: project)
  --json                Print JSON output
  -h, --help            Show help

What this does:
  1) Optionally writes IMAGE_GEN_API_KEY into the local skill .env
  2) Runs doctor (read-only)
  3) Ensures a MODELS.json exists (writes the bundled starter if missing)
  4) Initializes EXTEND.md with Nano Banana 2 if no default model exists yet
  5) Prints a recommend-model command for your workflow

The API key is never echoed back in output. When available, this script can persist it to ${defaultEnvPath}.`);
}

function parseArgs(argv) {
  const args = {
    project: process.cwd(),
    workflow: "general",
    forceModels: false,
    apiKey: null,
    persistApiKey: false,
    scope: "project",
    json: false,
  };
  for (let i = 0; i < argv.length; i++) {
    const current = argv[i];
    if (current === "--project") args.project = argv[++i] ?? args.project;
    else if (current === "--workflow") args.workflow = argv[++i] ?? args.workflow;
    else if (current === "--force-models") args.forceModels = true;
    else if (current === "--api-key") args.apiKey = argv[++i] ?? null;
    else if (current === "--persist-api-key") args.persistApiKey = true;
    else if (current === "--scope") {
      const value = argv[++i] ?? args.scope;
      if (value !== "project" && value !== "home") throw new Error(`Invalid --scope: ${value}`);
      args.scope = value;
    } else if (current === "--json") args.json = true;
    else if (current === "--help" || current === "-h") {
      printHelp();
      process.exit(0);
    }
  }
  return args;
}

function resolveProvidedApiKey({ apiKey, persistApiKey, env }) {
  if (typeof apiKey === "string" && apiKey.trim()) {
    return { value: apiKey.trim(), source: "arg", requested: true };
  }
  if (persistApiKey) {
    const value = env.IMAGE_GEN_API_KEY?.trim() ?? env.IMAGE_GEN_API_KEY?.trim() ?? null;
    return { value, source: value ? "env" : "missing-env", requested: true };
  }
  return {
    value: env.IMAGE_GEN_API_KEY?.trim() ?? env.IMAGE_GEN_API_KEY?.trim() ?? null,
    source: env.IMAGE_GEN_API_KEY || env.IMAGE_GEN_API_KEY ? "env" : "none",
    requested: false,
  };
}

function serializeEnvValue(value) {
  return JSON.stringify(value);
}

function upsertEnvVar(content, key, value) {
  const lines = content ? content.split(/\r?\n/) : [];
  const pattern = new RegExp(`^\\s*${key}\\s*=`);
  let found = false;
  const nextLines = [];
  for (const line of lines) {
    if (!found && pattern.test(line) && !line.trim().startsWith("#")) {
      nextLines.push(`${key}=${serializeEnvValue(value)}`);
      found = true;
    } else if (line.length > 0) {
      nextLines.push(line);
    }
  }
  if (!found) nextLines.push(`${key}=${serializeEnvValue(value)}`);
  return `${nextLines.join("\n")}\n`;
}

function resolveEnvPath({ project, homeDir, scope }) {
  const targetRoot = scope === "home" ? homeDir : project;
  if (!targetRoot) return null;
  return path.resolve(targetRoot, defaultEnvPath);
}

function writeApiKeyConfig({ project, homeDir, scope, apiKeyInput }) {
  const envPath = resolveEnvPath({ project, homeDir, scope });
  if (!apiKeyInput.requested) {
    return { ok: true, wrote: false, path: envPath, status: "not-requested", source: apiKeyInput.source };
  }
  if (!envPath) {
    return { ok: false, wrote: false, path: null, status: "missing-home-dir", source: apiKeyInput.source };
  }
  if (!apiKeyInput.value) {
    return { ok: false, wrote: false, path: envPath, status: "missing-api-key", source: apiKeyInput.source };
  }

  const existing = fs.existsSync(envPath) ? fs.readFileSync(envPath, "utf8") : "";
  const updated = upsertEnvVar(existing, "IMAGE_GEN_API_KEY", apiKeyInput.value);
  const wrote = updated !== existing;
  fs.mkdirSync(path.dirname(envPath), { recursive: true });
  if (wrote) fs.writeFileSync(envPath, updated, "utf8");
  return { ok: true, wrote, path: envPath, status: wrote ? "wrote" : "kept", source: apiKeyInput.source };
}

async function runDoctor({ project, workflow, env }) {
  const { doctorSuite } = await import("./doctor.mjs");
  return doctorSuite({ suiteDir: defaultSkillDir, projectDir: project, workflow, env });
}

async function ensureModelsRegistry({ project, force }) {
  const modelsPath = path.resolve(project, ".image-skills", skillNamespace, "MODELS.json");
  if (fs.existsSync(modelsPath) && !force) {
    return { ok: true, path: modelsPath, wrote: false, reason: "exists" };
  }
  const registry = fs.existsSync(bundledModelsPath)
    ? JSON.parse(fs.readFileSync(bundledModelsPath, "utf8"))
    : await import("./discover-models.mjs").then((m) => m.discoverImageModels());
  fs.mkdirSync(path.dirname(modelsPath), { recursive: true });
  fs.writeFileSync(modelsPath, `${JSON.stringify(registry, null, 2)}\n`, "utf8");
  return {
    ok: true,
    path: modelsPath,
    wrote: true,
    reason: fs.existsSync(bundledModelsPath) ? "bundled-starter" : "generated-from-docs",
  };
}

function ensureDefaultExtend({ project }) {
  const extendPath = path.resolve(project, defaultExtendPath);
  if (fs.existsSync(extendPath)) {
    return {
      ok: true,
      path: extendPath,
      wrote: false,
      defaultModel: null,
    };
  }

  const content = `---
version: 1
default_model: "${defaultStarterModel.id}"
default_quality: 2k
default_aspect_ratio: "1:1"
batch:
  max_workers: 4
---
`;
  fs.mkdirSync(path.dirname(extendPath), { recursive: true });
  fs.writeFileSync(extendPath, content, "utf8");
  return {
    ok: true,
    path: extendPath,
    wrote: true,
    defaultModel: defaultStarterModel,
  };
}

function formatApiKeyLine(apiKey) {
  switch (apiKey.status) {
    case "wrote":
      return `API key config: ${apiKey.path} (wrote)`;
    case "kept":
      return `API key config: ${apiKey.path} (already up to date)`;
    case "missing-api-key":
      return `API key config: no key value was available to write to ${apiKey.path}`;
    case "missing-home-dir":
      return "API key config: HOME is unavailable, so --scope home could not be used";
    default:
      return `API key config: no local write requested (target would be ${apiKey.path ?? defaultEnvPath})`;
  }
}

function formatText(summary) {
  const lines = [];
  lines.push(`Project: ${summary.project}`);
  lines.push(formatApiKeyLine(summary.apiKey));
  lines.push(`Doctor ok: ${summary.doctor.ok ? "yes" : "no"}`);
  if (summary.doctor.errors?.length) {
    lines.push("Doctor errors:");
    for (const e of summary.doctor.errors) lines.push(`- ${e}`);
  }
  if (summary.doctor.warnings?.length) {
    lines.push("Doctor warnings:");
    for (const w of summary.doctor.warnings) lines.push(`- ${w}`);
  }
  lines.push(`Models registry: ${summary.models.path} (${summary.models.wrote ? "wrote" : "kept"})`);
  lines.push(`Default config: ${summary.extend.path} (${summary.extend.wrote ? "wrote" : "kept"})`);
  if (summary.extend.defaultModel) {
    lines.push(
      `Initialized default model: ${summary.extend.defaultModel.label} (${summary.extend.defaultModel.id}). You can switch anytime later.`
    );
  }
  lines.push("Next:");
  lines.push(`- npm run recommend-model -- --workflow ${summary.workflow} ${summary.workflow === "comic" ? "--role comic-page" : ""}`.trim());
  lines.push(`- Current default is ${defaultStarterModel.label} (${defaultStarterModel.id}); change ${defaultExtendPath} or IMAGE_GEN_DEFAULT_MODEL any time if needed`);
  if (!summary.apiKey.ok) {
    lines.push(`- API key write failed: ${summary.apiKey.status}`);
  } else if (!summary.doctor.config.apiKeyInEnv && !summary.doctor.config.apiKeyInProjectEnv && !summary.apiKey.wrote) {
    lines.push(`- Provide IMAGE_GEN_API_KEY and rerun setup with --persist-api-key, or let the agent write ${defaultEnvPath} for you`);
  }
  return lines.join("\n");
}

const isDirectRun = process.argv[1] && path.resolve(process.argv[1]).endsWith("setup.mjs");

if (isDirectRun) {
  try {
    const args = parseArgs(process.argv.slice(2));
    const apiKeyInput = resolveProvidedApiKey({ apiKey: args.apiKey, persistApiKey: args.persistApiKey, env: process.env });
    const apiKey = writeApiKeyConfig({
      project: args.project,
      homeDir: process.env.HOME ?? null,
      scope: args.scope,
      apiKeyInput,
    });
    const doctorEnv = apiKeyInput.value
      ? { ...process.env, IMAGE_GEN_API_KEY: apiKeyInput.value, IMAGE_GEN_API_KEY: apiKeyInput.value }
      : process.env;
    const doctor = await runDoctor({ project: args.project, workflow: args.workflow, env: doctorEnv });
    const models = await ensureModelsRegistry({ project: args.project, force: args.forceModels });
    const extend = ensureDefaultExtend({ project: args.project });
    const summary = {
      ok: doctor.ok && models.ok && extend.ok && apiKey.ok,
      project: path.resolve(args.project),
      workflow: args.workflow,
      apiKey,
      doctor,
      models,
      extend,
    };
    if (args.json) console.log(JSON.stringify(summary, null, 2));
    else console.log(formatText(summary));
    if (!summary.ok) process.exitCode = 1;
  } catch (error) {
    console.error(error instanceof Error ? error.message : String(error));
    process.exit(1);
  }
}
