import fs from "node:fs";
import path from "node:path";
import process from "node:process";
import { fileURLToPath } from "node:url";

import { bootstrapSuite } from "./bootstrap.mjs";
import { ensureReady as ensureBaseReady } from "./vendor/weryai-image/ensure-ready.mjs";

const DEFAULT_WORKFLOW = "infographic";
const skillDir = path.resolve(path.dirname(fileURLToPath(import.meta.url)), "..");

const DEFAULT_MODEL = "GEMINI_3_1_FLASH_IMAGE";

function ensureDefaultModelConfig({ projectDir, dryRun }) {
  // Idempotent: only write EXTEND.md when no default model exists yet.
  if (process.env.IMAGE_GEN_DEFAULT_MODEL) return { wrote: false, reason: "env" };

  const namespace = path.basename(skillDir);
  const extendPath = path.resolve(projectDir, ".image-skills", namespace, "EXTEND.md");
  if (fs.existsSync(extendPath)) return { wrote: false, reason: "exists" };
  if (dryRun) return { wrote: false, reason: "dry-run" };

  const content = `---
version: 1
default_model: "GEMINI_3_1_FLASH_IMAGE"
default_quality: 2k
default_aspect_ratio: "1:1"
batch:
  max_workers: 4
---
`;
  fs.mkdirSync(path.dirname(extendPath), { recursive: true });
  fs.writeFileSync(extendPath, content, "utf8");
  return { wrote: true, reason: "created", path: extendPath };
}

function parseArgs(argv) {
  const args = { project: process.cwd(), workflow: DEFAULT_WORKFLOW, dryRun: false, json: false };
  for (let i = 0; i < argv.length; i++) {
    const current = argv[i];
    if (current === "--project") args.project = argv[++i] ?? args.project;
    else if (current === "--workflow") args.workflow = argv[++i] ?? args.workflow;
    else if (current === "--dry-run") args.dryRun = true;
    else if (current === "--json") args.json = true;
  }
  return args;
}

function formatText(summary) {
  const lines = [];
  lines.push(`Skill dir: ${summary.skillDir}`);
  lines.push(`Workflow: ${summary.workflow}`);
  lines.push(`Local bootstrap ok: ${summary.localBootstrap.ok ? "yes" : "no"}`);
  lines.push(`Base runtime ready: ${summary.base.ok ? "yes" : "no"}`);
  if (summary.defaultModelConfig?.wrote) lines.push(`Default model: initialized to GEMINI_3_1_FLASH_IMAGE`);
  if (summary.localBootstrap.targets.length) {
    lines.push("Local bootstrap targets:");
    for (const target of summary.localBootstrap.targets) lines.push(`- ${target.relativeDir}`);
  }
  if (summary.base.final.errors.length) {
    lines.push("Errors:");
    for (const item of summary.base.final.errors) lines.push(`- ${item}`);
  }
  if (summary.base.final.warnings.length) {
    lines.push("Warnings:");
    for (const item of summary.base.final.warnings) lines.push(`- ${item}`);
  }
  if (summary.base.final.suggestions.length) {
    lines.push("Suggestions:");
    for (const item of summary.base.final.suggestions) lines.push(`- ${item}`);
  }
  return lines.join("\n");
}

const args = parseArgs(process.argv.slice(2));
const localBootstrap = bootstrapSuite({ dir: skillDir, dryRun: args.dryRun });
const base = await ensureBaseReady({
  suiteDir: skillDir,
  projectDir: args.project,
  workflow: args.workflow,
  dryRun: args.dryRun,
  skillNamespace: path.basename(skillDir),
});
const defaultModelConfig = ensureDefaultModelConfig({ projectDir: args.project, dryRun: args.dryRun });

const summary = {
  ok: localBootstrap.ok && base.ok,
  skillDir,
  workflow: args.workflow,
  localBootstrap,
  base,
  defaultModelConfig,
};
if (args.json) console.log(JSON.stringify(summary, null, 2));
else console.log(formatText(summary));
if (!summary.ok) process.exitCode = 1;
