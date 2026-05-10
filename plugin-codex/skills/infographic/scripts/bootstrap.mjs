import fs from "node:fs";
import path from "node:path";
import process from "node:process";
import { spawnSync } from "node:child_process";
import { fileURLToPath } from "node:url";

const npmCmd = process.platform === "win32" ? "npm.cmd" : "npm";
const SKIP_DIRS = new Set(["node_modules", ".git", "dist"]);

function printHelp() {
  console.log(`Usage:
  node scripts/bootstrap.mjs [options]

Options:
  --dir <path>     Skill root directory (default: current skill root)
  --dry-run        Print install targets without running npm install
  --json           Print JSON output
  -h, --help       Show help`);
}

function parseArgs(argv) {
  const args = {
    dir: null,
    dryRun: false,
    json: false,
  };

  for (let i = 0; i < argv.length; i++) {
    const current = argv[i];
    if (current === "--dir") args.dir = argv[++i] ?? null;
    else if (current === "--dry-run") args.dryRun = true;
    else if (current === "--json") args.json = true;
    else if (current === "--help" || current === "-h") {
      printHelp();
      process.exit(0);
    }
  }
  return args;
}

function readJson(filePath) {
  return JSON.parse(fs.readFileSync(filePath, "utf8"));
}

function countInstallableDependencies(pkg) {
  const sections = ["dependencies", "optionalDependencies", "devDependencies"];
  return sections.reduce((total, field) => total + Object.keys(pkg[field] ?? {}).length, 0);
}

function collectPackageJsonPaths(rootDir) {
  const results = [];
  function visit(currentDir) {
    for (const entry of fs.readdirSync(currentDir, { withFileTypes: true })) {
      if (entry.isDirectory()) {
        if (SKIP_DIRS.has(entry.name)) continue;
        visit(path.join(currentDir, entry.name));
        continue;
      }
      if (entry.isFile() && entry.name === "package.json") {
        results.push(path.join(currentDir, entry.name));
      }
    }
  }
  visit(rootDir);
  return results;
}

export function discoverBootstrapTargets(rootDir) {
  const suiteRoot = path.resolve(rootDir);
  return collectPackageJsonPaths(suiteRoot)
    .map((packageJsonPath) => {
      const pkg = readJson(packageJsonPath);
      const dependencyCount = countInstallableDependencies(pkg);
      return {
        dir: path.dirname(packageJsonPath),
        relativeDir: path.relative(suiteRoot, path.dirname(packageJsonPath)) || ".",
        packageName: pkg.name ?? path.basename(path.dirname(packageJsonPath)),
        dependencyCount,
      };
    })
    .filter((target) => target.dependencyCount > 0)
    .sort((a, b) => a.relativeDir.localeCompare(b.relativeDir));
}

function installTarget(target) {
  const result = spawnSync(npmCmd, ["install"], {
    cwd: target.dir,
    stdio: "inherit",
  });
  return {
    ...target,
    ok: result.status === 0,
    exitCode: result.status ?? 1,
  };
}

export function bootstrapSuite({ dir, dryRun = false }) {
  const suiteRoot = path.resolve(dir);
  const targets = discoverBootstrapTargets(suiteRoot);
  const summary = {
    ok: true,
    dir: suiteRoot,
    dryRun,
    targetCount: targets.length,
    targets,
    installed: [],
  };

  if (dryRun || targets.length === 0) {
    return summary;
  }

  for (const target of targets) {
    const installed = installTarget(target);
    summary.installed.push(installed);
    if (!installed.ok) {
      summary.ok = false;
      break;
    }
  }
  return summary;
}

function formatText(summary) {
  const lines = [];
  lines.push(`Skill root: ${summary.dir}`);
  lines.push(`Dry run: ${summary.dryRun ? "yes" : "no"}`);
  lines.push(`Install targets: ${summary.targetCount}`);
  if (summary.targets.length) {
    lines.push("Targets:");
    for (const target of summary.targets) {
      lines.push(`- ${target.relativeDir} (${target.dependencyCount} dependencies)`);
    }
  }
  if (!summary.dryRun && summary.installed.length) {
    lines.push("Install results:");
    for (const target of summary.installed) {
      lines.push(`- ${target.relativeDir}: ${target.ok ? "ok" : `failed (${target.exitCode})`}`);
    }
  }
  return lines.join("\n");
}

const defaultDir = path.resolve(path.dirname(fileURLToPath(import.meta.url)), "..");
const isDirectRun =
  process.argv[1] && path.resolve(process.argv[1]) === path.resolve(fileURLToPath(import.meta.url));

if (isDirectRun) {
  try {
    const args = parseArgs(process.argv.slice(2));
    const summary = bootstrapSuite({ dir: args.dir ?? defaultDir, dryRun: args.dryRun });
    if (args.json) console.log(JSON.stringify(summary, null, 2));
    else console.log(formatText(summary));
    if (!summary.ok) process.exitCode = 1;
  } catch (error) {
    console.error(error instanceof Error ? error.message : String(error));
    process.exit(1);
  }
}
