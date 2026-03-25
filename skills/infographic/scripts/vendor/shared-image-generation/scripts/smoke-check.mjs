import process from "node:process";
import { spawnSync } from "node:child_process";
import path from "node:path";
import { fileURLToPath } from "node:url";

const rootDir = path.resolve(path.dirname(fileURLToPath(import.meta.url)), "..");
const npxCmd = process.platform === "win32" ? "npx.cmd" : "npx";

function parseArgs(argv) {
  const out = { model: null };
  for (let i = 0; i < argv.length; i++) {
    if (argv[i] === "--model") out.model = argv[++i] ?? null;
  }
  return out;
}

function runBunScript(relativeScript, args) {
  const result = spawnSync(npxCmd, ["-y", "bun", relativeScript, ...args], {
    cwd: rootDir,
    encoding: "utf8",
  });
  if (result.status !== 0) {
    throw new Error(result.stderr || result.stdout || `Command failed: ${relativeScript}`);
  }
  return result.stdout.trim();
}

const cli = parseArgs(process.argv.slice(2));
const model = cli.model ?? process.env.IMAGE_GEN_DEFAULT_MODEL ?? "smoke-model-placeholder";
const payload = runBunScript("scripts/main.ts", [
  "--prompt",
  "image-generation smoke preview",
  "--ar",
  "1:1",
  "--model",
  model,
  "--dry-run",
  "--json",
]);
console.log(payload);
