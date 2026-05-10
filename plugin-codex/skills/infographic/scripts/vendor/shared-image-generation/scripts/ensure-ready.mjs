import fs from "node:fs";
import path from "node:path";
import process from "node:process";
import { fileURLToPath } from "node:url";

import { bootstrapSuite } from "./bootstrap.mjs";
import { doctorSuite } from "./doctor.mjs";

const DEFAULT_MODEL = "GEMINI_3_1_FLASH_IMAGE";

function ensureDefaultModelConfig({ suiteDir, projectDir, dryRun }) {
  if (process.env.IMAGE_GEN_DEFAULT_MODEL) return { wrote: false, reason: "env" };
  const namespace = path.basename(path.resolve(suiteDir));
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

function printHelp() {
  console.log(`Usage:
  node scripts/ensure-ready.mjs [options]

Options:
  --dir <path>       Skill root directory (default: current image-generation root)
  --project <path>   Project directory to inspect for .image-skills config (default: current working directory)
  --workflow <name>  Optional readiness focus: general, cover, rednote, infographic, comic, article, compress
  --dry-run          Report what would be installed without running bootstrap
  --json             Print JSON output
  -h, --help         Show help`);
}

function parseArgs(argv) {
  const args = {
    dir: null,
    project: null,
    workflow: null,
    dryRun: false,
    json: false,
  };

  for (let i = 0; i < argv.length; i++) {
    const current = argv[i];
    if (current === "--dir") args.dir = argv[++i] ?? null;
    else if (current === "--project") args.project = argv[++i] ?? null;
    else if (current === "--workflow") args.workflow = argv[++i] ?? null;
    else if (current === "--dry-run") args.dryRun = true;
    else if (current === "--json") args.json = true;
    else if (current === "--help" || current === "-h") {
      printHelp();
      process.exit(0);
    }
  }

  return args;
}

function formatText(summary) {
  const lines = [];
  lines.push(`Skill root: ${summary.suiteDir}`);
  lines.push(`Project dir: ${summary.projectDir}`);
  lines.push(`Workflow: ${summary.workflow ?? "general"}`);
  if (summary.defaultModelConfig?.wrote) lines.push(`Default model: initialized to GEMINI_3_1_FLASH_IMAGE`);
  lines.push(`Bootstrap needed: ${summary.initial.bootstrap.pendingCount > 0 ? "yes" : "no"}`);
  lines.push(`Bootstrap attempted: ${summary.bootstrap.attempted ? "yes" : "no"}`);
  lines.push(`Bootstrap changed files: ${summary.bootstrap.changed ? "yes" : "no"}`);
  lines.push(`Ready: ${summary.final.ok ? "yes" : "no"}`);

  if (summary.bootstrap.targets.length) {
    lines.push("Bootstrap targets:");
    for (const target of summary.bootstrap.targets) {
      lines.push(`- ${target.relativeDir}`);
    }
  }

  if (summary.final.errors.length) {
    lines.push("Errors:");
    for (const error of summary.final.errors) lines.push(`- ${error}`);
  }
  if (summary.final.warnings.length) {
    lines.push("Warnings:");
    for (const warning of summary.final.warnings) lines.push(`- ${warning}`);
  }
  if (summary.final.suggestions.length) {
    lines.push("Suggestions:");
    for (const suggestion of summary.final.suggestions) lines.push(`- ${suggestion}`);
  }
  return lines.join("\n");
}

export function ensureReady({
  suiteDir,
  projectDir = process.cwd(),
  homeDir = process.env.HOME ?? null,
  env = process.env,
  workflow = null,
  dryRun = false,
  doctorImpl = doctorSuite,
  bootstrapImpl = bootstrapSuite,
}) {
  const defaultModelConfig = ensureDefaultModelConfig({ suiteDir, projectDir, dryRun });

  const initial = doctorImpl({
    suiteDir,
    projectDir,
    homeDir,
    env,
    workflow,
  });
  const pendingTargets = initial.bootstrap.targets.filter((target) => !target.installed);

  let bootstrapSummary = {
    attempted: false,
    changed: false,
    dryRun,
    ok: true,
    targets: pendingTargets,
    result: null,
  };

  if (pendingTargets.length > 0) {
    bootstrapSummary = {
      attempted: true,
      changed: !dryRun,
      dryRun,
      ok: true,
      targets: pendingTargets,
      result: bootstrapImpl({ dir: suiteDir, dryRun }),
    };
    if (bootstrapSummary.result && bootstrapSummary.result.ok === false) {
      bootstrapSummary.ok = false;
    }
  }

  const final = doctorImpl({
    suiteDir,
    projectDir,
    homeDir,
    env,
    workflow,
  });

  return {
    ok: final.ok && bootstrapSummary.ok,
    suiteDir: path.resolve(suiteDir),
    projectDir: path.resolve(projectDir),
    workflow,
    initial,
    bootstrap: bootstrapSummary,
    final,
    defaultModelConfig,
  };
}

const defaultSuiteDir = path.resolve(path.dirname(fileURLToPath(import.meta.url)), "..");
const isDirectRun =
  process.argv[1] && path.resolve(process.argv[1]) === path.resolve(fileURLToPath(import.meta.url));

if (isDirectRun) {
  try {
    const args = parseArgs(process.argv.slice(2));
    const summary = ensureReady({
      suiteDir: args.dir ?? defaultSuiteDir,
      projectDir: args.project ?? process.cwd(),
      workflow: args.workflow,
      dryRun: args.dryRun,
    });
    if (args.json) console.log(JSON.stringify(summary, null, 2));
    else console.log(formatText(summary));
    if (!summary.ok) process.exitCode = 1;
  } catch (error) {
    console.error(error instanceof Error ? error.message : String(error));
    process.exit(1);
  }
}
