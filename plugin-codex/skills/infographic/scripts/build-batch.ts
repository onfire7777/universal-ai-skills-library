#!/usr/bin/env bun
import path from "node:path";
import process from "node:process";
import { readdir, readFile, writeFile } from "node:fs/promises";
import {
  buildWorkflowNegativePrompt,
  resolveWorkflowStyle,
} from "./visual-policy.ts";

type CliArgs = {
  promptsDir: string | null;
  outputPath: string | null;
  imagesDir: string | null;
  model: string | null;
  style: string | null;
  aspectRatio: string;
  quality: string;
  jobs: number | null;
  help: boolean;
};

type PromptEntry = {
  order: number;
  promptPath: string;
  imageFilename: string;
};

type PromptMeta = {
  styleDirection?: string;
  aspect?: string;
};

function printUsage(): void {
  console.log(`Usage:
  npx -y bun scripts/build-batch.ts --prompts prompts --output batch.json --model <model> [options]

Options:
  --prompts <path>       Path to prompts directory
  --output <path>        Path to output batch.json
  --images-dir <path>    Directory for generated images
  --model <id>           Model key for bundled runtime batch tasks
  --style <name>         Explicit bundled runtime style override
  --ar <ratio>           Default aspect ratio when prompt files omit it (default: 3:4)
  --quality <level>      Quality for all tasks (default: 2k)
  --jobs <count>         Suggested worker count metadata (optional)
  -h, --help             Show help`);
}

function parseArgs(argv: string[]): CliArgs {
  const args: CliArgs = {
    promptsDir: null,
    outputPath: null,
    imagesDir: null,
    model: null,
    style: null,
    aspectRatio: "3:4",
    quality: "2k",
    jobs: null,
    help: false,
  };

  for (let i = 0; i < argv.length; i++) {
    const current = argv[i]!;
    if (current === "--prompts") args.promptsDir = argv[++i] ?? null;
    else if (current === "--output") args.outputPath = argv[++i] ?? null;
    else if (current === "--images-dir") args.imagesDir = argv[++i] ?? null;
    else if (current === "--model") args.model = argv[++i] ?? null;
    else if (current === "--style") args.style = argv[++i] ?? null;
    else if (current === "--ar") args.aspectRatio = argv[++i] ?? args.aspectRatio;
    else if (current === "--quality") args.quality = argv[++i] ?? args.quality;
    else if (current === "--jobs") {
      const value = argv[++i];
      args.jobs = value ? parseInt(value, 10) : null;
    } else if (current === "--help" || current === "-h") {
      args.help = true;
    }
  }

  return args;
}

function sortPromptFilenames(a: string, b: string): number {
  const matchA = a.match(/^(\d+)/);
  const matchB = b.match(/^(\d+)/);
  if (matchA && matchB) {
    const order = parseInt(matchA[1]!, 10) - parseInt(matchB[1]!, 10);
    if (order !== 0) return order;
  } else if (matchA) {
    return -1;
  } else if (matchB) {
    return 1;
  }
  return a.localeCompare(b);
}

async function collectPromptEntries(promptsDir: string): Promise<PromptEntry[]> {
  const files = await readdir(promptsDir);
  return files
    .filter((filename) => filename.toLowerCase().endsWith(".md"))
    .sort(sortPromptFilenames)
    .map((filename, index) => {
      const baseName = filename.replace(/\.md$/i, "");
      return {
        order: index + 1,
        promptPath: path.join(promptsDir, filename),
        imageFilename: `${baseName}.png`,
      };
    });
}

function parsePromptMeta(content: string): PromptMeta {
  const normalized = content.replace(/\r\n/g, "\n");
  const styleMatch = normalized.match(/^Style direction:\s*(.+?)(?:\.\s*|$)/im);
  const aspectMatch = normalized.match(/^Aspect ratio:\s*(.+?)(?:\.\s*|$)/im);
  return {
    styleDirection: styleMatch?.[1]?.trim().toLowerCase(),
    aspect: aspectMatch?.[1]?.trim(),
  };
}

function resolveStyle(explicitStyle: string | null, styleDirection: string | undefined): string | null {
  return resolveWorkflowStyle("infographic", explicitStyle, styleDirection);
}

async function main(): Promise<void> {
  const args = parseArgs(process.argv.slice(2));
  if (args.help) {
    printUsage();
    return;
  }

  if (!args.promptsDir) {
    console.error("Error: --prompts is required");
    process.exit(1);
  }
  if (!args.outputPath) {
    console.error("Error: --output is required");
    process.exit(1);
  }
  if (!args.model) {
    console.error("Error: --model is required");
    process.exit(1);
  }

  const entries = await collectPromptEntries(args.promptsDir);
  if (entries.length === 0) {
    console.error("No infographic prompt files found in prompts directory.");
    process.exit(1);
  }

  const imageDir = args.imagesDir ?? path.dirname(args.outputPath);
  const tasks = [];

  for (const entry of entries) {
    const promptContent = await readFile(entry.promptPath, "utf8");
    const meta = parsePromptMeta(promptContent);
    const task: Record<string, unknown> = {
      id: `infographic-${String(entry.order).padStart(2, "0")}`,
      promptFiles: [entry.promptPath],
      image: path.join(imageDir, entry.imageFilename),
      model: args.model,
      ar: meta.aspect ?? args.aspectRatio,
      quality: args.quality,
      negative_prompt: buildWorkflowNegativePrompt("infographic"),
    };
    const style = resolveStyle(args.style, meta.styleDirection);
    if (style) task.style = style;
    tasks.push(task);
  }

  const output: Record<string, unknown> = { tasks };
  if (args.jobs) output.jobs = args.jobs;

  await writeFile(args.outputPath, JSON.stringify(output, null, 2) + "\n");
  console.log(`Batch file written: ${args.outputPath} (${tasks.length} tasks)`);
}

main().catch((error) => {
  console.error(error instanceof Error ? error.message : String(error));
  process.exit(1);
});
