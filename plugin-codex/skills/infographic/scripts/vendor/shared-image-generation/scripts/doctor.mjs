import fs from "node:fs";
import path from "node:path";
import process from "node:process";
import { spawnSync } from "node:child_process";
import { fileURLToPath } from "node:url";

import { discoverBootstrapTargets } from "./bootstrap.mjs";

const npmCmd = process.platform === "win32" ? "npm.cmd" : "npm";
const npxCmd = process.platform === "win32" ? "npx.cmd" : "npx";
const VALID_WORKFLOWS = new Set(["general", "cover", "rednote", "infographic", "comic", "article", "compress"]);
const skillNamespace = process.env.IMAGE_SKILL_NAMESPACE?.trim() || "image-generation";

function printHelp() {
  console.log(`Usage:
  node scripts/doctor.mjs [options]

Options:
  --dir <path>       Skill root directory (default: current skill root)
  --project <path>   Project directory to inspect for .image-skills config (default: current working directory)
  --workflow <name>  Optional readiness focus: general, cover, rednote, infographic, comic, article, compress
  --json             Print JSON output
  -h, --help         Show help`);
}

function parseArgs(argv) {
  const args = {
    dir: null,
    project: null,
    workflow: null,
    json: false,
  };

  for (let i = 0; i < argv.length; i++) {
    const current = argv[i];
    if (current === "--dir") args.dir = argv[++i] ?? null;
    else if (current === "--project") args.project = argv[++i] ?? null;
    else if (current === "--workflow") args.workflow = argv[++i] ?? null;
    else if (current === "--json") args.json = true;
    else if (current === "--help" || current === "-h") {
      printHelp();
      process.exit(0);
    }
  }
  if (args.workflow && !VALID_WORKFLOWS.has(args.workflow)) {
    throw new Error(`Invalid --workflow: ${args.workflow}`);
  }
  return args;
}

function pathExists(filePath) {
  return !!filePath && fs.existsSync(filePath);
}

function checkCommand(command, args) {
  const result = spawnSync(command, args, { encoding: "utf8" });
  return {
    available: result.status === 0,
    version: result.status === 0 ? (result.stdout || result.stderr).trim().split("\n")[0] : null,
  };
}

function inspectToolchain() {
  const nodeStatus = checkCommand(process.execPath, ["--version"]);
  const npmStatus = checkCommand(npmCmd, ["--version"]);
  const bunStatus = checkCommand("bun", ["--version"]);
  const npxStatus = checkCommand(npxCmd, ["--version"]);
  const sipsStatus = checkCommand("sips", ["--help"]);
  const cwebpStatus = checkCommand("cwebp", ["-version"]);
  const magickStatus = checkCommand("magick", ["-version"]);
  const convertStatus = checkCommand("convert", ["-version"]);
  return [
    { id: "node", label: "Node.js", required: true, ...nodeStatus },
    { id: "npm", label: "npm", required: true, ...npmStatus },
    { id: "bun", label: "Bun", required: false, ...bunStatus },
    { id: "npx", label: "npx", required: false, ...npxStatus },
    { id: "sips", label: "sips", required: false, ...sipsStatus },
    { id: "cwebp", label: "cwebp", required: false, ...cwebpStatus },
    { id: "magick", label: "ImageMagick", required: false, ...magickStatus },
    { id: "convert", label: "convert", required: false, ...convertStatus },
  ];
}

export function inspectConfigState({ projectDir, homeDir, env }) {
  const projectRoot = path.resolve(projectDir);
  const homeRoot = homeDir ? path.resolve(homeDir) : null;
  const projectSkillDir = path.join(projectRoot, ".image-skills", skillNamespace);
  const homeSkillDir = homeRoot ? path.join(homeRoot, ".image-skills", skillNamespace) : null;

  const projectEnvPath = path.join(projectRoot, ".image-skills", skillNamespace, ".env");
  const projectEnvFile = pathExists(projectEnvPath);
  let apiKeyInProjectEnv = false;
  if (projectEnvFile) {
    try {
      const content = fs.readFileSync(projectEnvPath, "utf8");
      apiKeyInProjectEnv =
        /(^|\n)\s*IMAGE_GEN_API_KEY\s*=\s*[^\n#]+/m.test(content) ||
        /(^|\n)\s*IMAGE_GEN_API_KEY\s*=\s*[^\n#]+/m.test(content);
    } catch {
      apiKeyInProjectEnv = false;
    }
  }

  return {
    projectDir: projectRoot,
    apiKeyInEnv: Boolean(env.IMAGE_GEN_API_KEY || env.IMAGE_GEN_API_KEY),
    apiKeyInProjectEnv,
    defaultModelInEnv: Boolean(env.IMAGE_GEN_DEFAULT_MODEL),
    projectEnvFile,
    projectExtend: pathExists(path.join(projectSkillDir, "EXTEND.md")),
    projectModelRegistry: pathExists(path.join(projectSkillDir, "MODELS.json")),
    homeExtend: homeSkillDir ? pathExists(path.join(homeSkillDir, "EXTEND.md")) : false,
    homeModelRegistry: homeSkillDir ? pathExists(path.join(homeSkillDir, "MODELS.json")) : false,
  };
}

function readJson(filePath) {
  return JSON.parse(fs.readFileSync(filePath, "utf8"));
}

export function inspectBootstrapDependencyState(suiteDir) {
  return discoverBootstrapTargets(suiteDir).map((target) => {
    const pkg = readJson(path.join(target.dir, "package.json"));
    const dependencyNames = [
      ...Object.keys(pkg.dependencies ?? {}),
      ...Object.keys(pkg.optionalDependencies ?? {}),
      ...Object.keys(pkg.devDependencies ?? {}),
    ];
    const missingPackages = dependencyNames.filter(
      (dependencyName) => !pathExists(path.join(target.dir, "node_modules", dependencyName, "package.json"))
    );
    return {
      ...target,
      dependencyNames,
      installed: missingPackages.length === 0,
      missingPackages,
    };
  });
}

function hasDetectedApiKey(config) {
  return config.apiKeyInEnv || config.apiKeyInProjectEnv;
}

function hasDetectedDefaultModel(config) {
  return config.defaultModelInEnv || config.projectExtend || config.homeExtend;
}

function findBootstrapTarget(bootstrapTargets, relativeDir) {
  return bootstrapTargets.find((target) => target.relativeDir === relativeDir) ?? null;
}

function toolAvailable(tools, id) {
  return Boolean(tools.find((tool) => tool.id === id)?.available);
}

function buildRootCommandHints(projectDir, workflow = null) {
  const commands = [
    `cd "${projectDir}"`,
    "npm run doctor -- --project .",
  ];
  if (workflow) commands.push(`npm run doctor -- --project . --workflow ${workflow}`);
  commands.push("npm run setup -- --project .");
  commands.push(`npm run discover-image-models -- --out .image-skills/${skillNamespace}/MODELS.json`);
  commands.push("npm run recommend-model -- --workflow cover --role cover-hero");
  return commands;
}

function defaultRoleForWorkflow(workflow) {
  switch (workflow) {
    case "cover":
      return "cover-hero";
    case "rednote":
      return "rednote-text-card";
    case "infographic":
      return "infographic-dense";
    case "comic":
      return "comic-page";
    case "article":
      return "article-framework";
    default:
      return null;
  }
}

export function inspectWorkflowReadiness({ workflow, tools, config, bootstrapTargets }) {
  if (!workflow) return null;

  const blockers = [];
  const warnings = [];
  const suggestions = [];
  const checks = [];
  const hasApiKey = hasDetectedApiKey(config);
  const hasDefaultModel = hasDetectedDefaultModel(config);
  const hasModelRegistry = config.projectModelRegistry || config.homeModelRegistry;
  const recommendedRole = defaultRoleForWorkflow(workflow);

  const mark = (label, ok, detail = null) => {
    checks.push({ label, ok, detail });
    return ok;
  };

  if (workflow === "compress") {
    const compressorOk =
      toolAvailable(tools, "sips") ||
      toolAvailable(tools, "cwebp") ||
      toolAvailable(tools, "magick") ||
      toolAvailable(tools, "convert");
    if (!mark("local-compressor-tool", compressorOk, "Need one of sips, cwebp, magick, or convert")) {
      blockers.push("No local image compression tool was detected");
      suggestions.push("Install at least one of sips, cwebp, or ImageMagick before using compress-image");
    }
  } else {
    if (!mark("gateway-api-key", hasApiKey, `Need IMAGE_GEN_API_KEY in env or .image-skills/${skillNamespace}/.env`)) {
      blockers.push("Gateway API key is not configured");
      suggestions.push(`Provide IMAGE_GEN_API_KEY, then rerun \`npm run setup -- --project . --persist-api-key\` or write .image-skills/${skillNamespace}/.env`);
    }
    if (!mark("default-model", hasDefaultModel, "Need IMAGE_GEN_DEFAULT_MODEL or EXTEND.md default_model for smooth first use")) {
      warnings.push("No default model is configured yet");
      suggestions.push("Run npm run setup -- --project . to initialize Nano Banana 2 as the default model, or set IMAGE_GEN_DEFAULT_MODEL / EXTEND.md yourself");
    }
    if (!mark("model-registry", hasModelRegistry, "MODELS.json enables workflow and role-aware model ranking")) {
      warnings.push("No MODELS.json registry was detected for workflow-aware model recommendation");
      suggestions.push("Run npm run setup -- --project . to write the bundled starter registry, or refresh from docs before model recommendation");
    }
  }

  if (workflow === "comic") {
    const comicBootstrap = findBootstrapTarget(bootstrapTargets, "comic/scripts");
    const pdfReady = comicBootstrap ? comicBootstrap.installed : true;
    if (!mark("comic-pdf-merge", pdfReady, "merge-to-pdf.ts needs local pdf-lib install")) {
      warnings.push("Comic PDF merge dependency is not installed yet");
      suggestions.push("Run npm run bootstrap in this skill directory before using comic PDF merge");
    }
  }

  const recommendationHint =
    workflow !== "compress" && recommendedRole
      ? `npm run recommend-model -- --workflow ${workflow} --role ${recommendedRole}`
      : null;
  if (recommendationHint) {
    suggestions.push(`Before generation, rank models with: ${recommendationHint}`);
  }

  return {
    workflow,
    ready: blockers.length === 0,
    recommendationHint,
    blockers,
    warnings,
    suggestions,
    checks,
  };
}

export function doctorSuite({
  suiteDir,
  projectDir = process.cwd(),
  homeDir = process.env.HOME ?? null,
  env = process.env,
  tools = null,
  workflow = null,
}) {
  const resolvedSuiteDir = path.resolve(suiteDir);
  const toolchain = tools ?? inspectToolchain();
  const config = inspectConfigState({ projectDir, homeDir, env });
  const bootstrap = inspectBootstrapDependencyState(resolvedSuiteDir);
  const warnings = [];
  const errors = [];
  const suggestions = [];

  const hasNode = toolchain.find((tool) => tool.id === "node")?.available;
  const hasNpm = toolchain.find((tool) => tool.id === "npm")?.available;
  const hasBun = toolchain.find((tool) => tool.id === "bun")?.available;
  const hasNpx = toolchain.find((tool) => tool.id === "npx")?.available;

  if (!hasNode) errors.push("Node.js is required but was not detected");
  if (!hasNpm) errors.push("npm is required but was not detected");
  if (!hasBun && !hasNpx) errors.push("Need Bun or npx so Bun-based workflow scripts can run");

  const hasApiKey = hasDetectedApiKey(config);
  if (!hasApiKey) {
    warnings.push(`IMAGE_GEN_API_KEY was not detected in the environment or project .image-skills/${skillNamespace}/.env`);
    suggestions.push(`Provide IMAGE_GEN_API_KEY, then rerun \`npm run setup -- --project . --persist-api-key\` or write .image-skills/${skillNamespace}/.env before using gateway generation`);
  }

  const hasDefaultModel = hasDetectedDefaultModel(config);
  if (!hasDefaultModel) {
    warnings.push("No default model was detected in IMAGE_GEN_DEFAULT_MODEL or EXTEND.md");
    suggestions.push("Run npm run setup -- --project . to initialize Nano Banana 2 as the default model, or set IMAGE_GEN_DEFAULT_MODEL / EXTEND.md yourself");
  }

  const pendingBootstrap = bootstrap.filter((target) => !target.installed);
  if (pendingBootstrap.length > 0) {
    warnings.push(`Optional local dependencies are not installed for ${pendingBootstrap.map((target) => target.relativeDir).join(", ")}`);
    suggestions.push("Run npm run bootstrap in this skill directory to install optional local dependencies");
  }

  const workflowReadiness = inspectWorkflowReadiness({
    workflow,
    tools: toolchain,
    config,
    bootstrapTargets: bootstrap,
  });

  return {
    ok: errors.length === 0,
    suiteDir: resolvedSuiteDir,
    projectDir: config.projectDir,
    workflow,
    behavior: {
      readOnly: true,
      writesFiles: false,
      generatedArtifacts: [],
    },
    commandHints: buildRootCommandHints(config.projectDir, workflow),
    tools: toolchain,
    config,
    bootstrap: {
      targets: bootstrap,
      pendingCount: pendingBootstrap.length,
    },
    workflowReadiness,
    warnings,
    errors,
    suggestions,
  };
}

function formatText(summary) {
  const lines = [];
  lines.push(`Skill root: ${summary.suiteDir}`);
  lines.push(`Project dir: ${summary.projectDir}`);
  lines.push(`Status: ${summary.ok ? "ok" : "failed"}`);
  lines.push("Doctor behavior: read-only check; it does not create or modify files.");
  if (summary.workflowReadiness) {
    lines.push(`Workflow focus: ${summary.workflowReadiness.workflow}`);
    lines.push(`Workflow ready: ${summary.workflowReadiness.ready ? "yes" : "no"}`);
  }
  lines.push("Toolchain:");
  for (const tool of summary.tools) {
    lines.push(`- ${tool.label}: ${tool.available ? tool.version ?? "ok" : "missing"}`);
  }
  lines.push("Gateway config:");
  lines.push(`- API key: ${summary.config.apiKeyInEnv || summary.config.projectEnvFile ? "detected" : "missing"}`);
  lines.push(`- Default model: ${summary.config.defaultModelInEnv || summary.config.projectExtend || summary.config.homeExtend ? "detected" : "missing"}`);
  lines.push("Optional dependencies:");
  if (summary.bootstrap.targets.length === 0) {
    lines.push("- none");
  } else {
    for (const target of summary.bootstrap.targets) {
      lines.push(`- ${target.relativeDir}: ${target.installed ? "installed" : `missing ${target.missingPackages.join(", ")}`}`);
    }
  }
  if (summary.errors.length) {
    lines.push("Errors:");
    for (const error of summary.errors) lines.push(`- ${error}`);
  }
  if (summary.warnings.length) {
    lines.push("Warnings:");
    for (const warning of summary.warnings) lines.push(`- ${warning}`);
  }
  if (summary.suggestions.length) {
    lines.push("Suggestions:");
    for (const suggestion of summary.suggestions) lines.push(`- ${suggestion}`);
  }
  if (summary.commandHints?.length) {
    lines.push("Suggested commands:");
    for (const command of summary.commandHints) lines.push(`- ${command}`);
  }
  if (summary.workflowReadiness) {
    if (summary.workflowReadiness.blockers.length) {
      lines.push("Workflow blockers:");
      for (const blocker of summary.workflowReadiness.blockers) lines.push(`- ${blocker}`);
    }
    if (summary.workflowReadiness.warnings.length) {
      lines.push("Workflow warnings:");
      for (const warning of summary.workflowReadiness.warnings) lines.push(`- ${warning}`);
    }
    if (summary.workflowReadiness.suggestions.length) {
      lines.push("Workflow suggestions:");
      for (const suggestion of summary.workflowReadiness.suggestions) lines.push(`- ${suggestion}`);
    }
    if (summary.workflowReadiness.recommendationHint) {
      lines.push(`Workflow model hint: ${summary.workflowReadiness.recommendationHint}`);
    }
  }
  return lines.join("\n");
}

const defaultSuiteDir = path.resolve(path.dirname(fileURLToPath(import.meta.url)), "..");
const isDirectRun =
  process.argv[1] && path.resolve(process.argv[1]) === path.resolve(fileURLToPath(import.meta.url));

if (isDirectRun) {
  try {
    const args = parseArgs(process.argv.slice(2));
    const summary = doctorSuite({
      suiteDir: args.dir ?? defaultSuiteDir,
      projectDir: args.project ?? process.cwd(),
      workflow: args.workflow,
    });
    if (args.json) console.log(JSON.stringify(summary, null, 2));
    else console.log(formatText(summary));
    if (!summary.ok) process.exitCode = 1;
  } catch (error) {
    console.error(error instanceof Error ? error.message : String(error));
    process.exit(1);
  }
}
