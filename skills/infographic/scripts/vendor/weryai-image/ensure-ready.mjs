import path from "node:path";
import process from "node:process";
import { fileURLToPath } from "node:url";

export async function ensureReady(options = {}) {
  const skillNamespace = options.skillNamespace ?? process.env.IMAGE_SKILL_NAMESPACE ?? path.basename(options.projectDir ?? process.cwd());
  const env = { ...process.env, IMAGE_SKILL_NAMESPACE: skillNamespace };
  const previous = process.env.IMAGE_SKILL_NAMESPACE;
  process.env.IMAGE_SKILL_NAMESPACE = skillNamespace;
  try {
    const { ensureReady: sharedEnsureReady } = await import("../shared-image-generation/scripts/ensure-ready.mjs");
    return sharedEnsureReady({
      suiteDir: options.suiteDir,
      projectDir: options.projectDir,
      workflow: options.workflow,
      dryRun: options.dryRun,
      env,
    });
  } finally {
    if (previous == null) delete process.env.IMAGE_SKILL_NAMESPACE;
    else process.env.IMAGE_SKILL_NAMESPACE = previous;
  }
}

const isDirectRun =
  process.argv[1] && path.resolve(process.argv[1]) === path.resolve(fileURLToPath(import.meta.url));

if (isDirectRun) {
  const result = await ensureReady({
    suiteDir: process.cwd(),
    projectDir: process.cwd(),
  });
  console.log(JSON.stringify(result, null, 2));
  if (!result.ok) process.exitCode = 1;
}
