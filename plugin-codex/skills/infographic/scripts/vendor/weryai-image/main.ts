#!/usr/bin/env bun
import path from "node:path";
import process from "node:process";
import { fileURLToPath } from "node:url";
import { spawnSync } from "node:child_process";

const sharedScript = path.resolve(path.dirname(fileURLToPath(import.meta.url)), "../shared-image-generation/scripts/main.ts");

const skillNamespace = process.env.IMAGE_SKILL_NAMESPACE?.trim() || path.basename(process.cwd());
const result = spawnSync(process.execPath, [sharedScript, ...process.argv.slice(2)], {
  stdio: "inherit",
  env: {
    ...process.env,
    IMAGE_SKILL_NAMESPACE: skillNamespace,
    IMAGE_SKILL_LABEL: skillNamespace,
  },
});

process.exit(result.status ?? 1);
