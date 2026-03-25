#!/usr/bin/env bun
import path from "node:path";
import process from "node:process";
import { mkdir, readFile, writeFile } from "node:fs/promises";

type CliArgs = {
  structuredContentPath: string | null;
  outputPath: string | null;
  audience: string;
  style: string;
  aspect: string;
  language: string;
  help: boolean;
};

type Section = {
  heading: string;
  keyPoints: string[];
  labels: string[];
  visualHints: string[];
};

function printUsage(): void {
  console.log(`Usage:
  npx -y bun scripts/build-prompt.ts --structured-content structured-content.md --output prompts/infographic.md [options]

Options:
  --structured-content <path>  Path to structured-content.md
  --output <path>              Output prompt path
  --audience <text>            Target audience (default: target audience)
  --style <name>               Style direction (default: infographic)
  --aspect <ratio>             Aspect ratio (default: 3:4)
  --lang <code>                On-image text language when needed (default placeholder: en)
  -h, --help                   Show help`);
}

function parseArgs(argv: string[]): CliArgs {
  const args: CliArgs = {
    structuredContentPath: null,
    outputPath: null,
    audience: "target audience",
    style: "infographic",
    aspect: "3:4",
    language: "en",
    help: false,
  };

  for (let i = 0; i < argv.length; i++) {
    const current = argv[i]!;
    if (current === "--structured-content") args.structuredContentPath = argv[++i] ?? null;
    else if (current === "--output") args.outputPath = argv[++i] ?? null;
    else if (current === "--audience") args.audience = argv[++i] ?? args.audience;
    else if (current === "--style") args.style = argv[++i] ?? args.style;
    else if (current === "--aspect") args.aspect = argv[++i] ?? args.aspect;
    else if (current === "--lang") args.language = argv[++i] ?? args.language;
    else if (current === "--help" || current === "-h") args.help = true;
  }

  return args;
}

function extractBlock(content: string, heading: string): string {
  const normalized = content.replace(/\r\n/g, "\n");
  const lines = normalized.split("\n");
  const start = lines.findIndex((line) => line.trim() === `## ${heading}`);
  if (start === -1) return "";

  const body: string[] = [];
  for (let i = start + 1; i < lines.length; i++) {
    if (lines[i]!.startsWith("## ")) break;
    body.push(lines[i]!);
  }
  return body.join("\n").trim();
}

function extractSingleLineValue(block: string, label: string): string {
  const escaped = label.replace(/[.*+?^${}()|[\]\\]/g, "\\$&");
  const match = block.match(new RegExp(`^- ${escaped}:\\s*(.+)$`, "m"));
  return match?.[1]?.trim() ?? "";
}

function extractBulletList(block: string, label: string): string[] {
  const escaped = label.replace(/[.*+?^${}()|[\]\\]/g, "\\$&");
  const match = block.match(new RegExp(`^- ${escaped}:\\s*\\n((?:  - .+\\n?)*)`, "m"));
  if (!match?.[1]) return [];
  return match[1]
    .split("\n")
    .map((line) => line.trim())
    .filter((line) => line.startsWith("- "))
    .map((line) => line.slice(2).trim())
    .filter(Boolean);
}

function parseSections(content: string): Section[] {
  const normalized = content.replace(/\r\n/g, "\n");
  const sectionsBlock = extractBlock(normalized, "Sections");
  const parts = sectionsBlock.split(/^### /m).slice(1);
  return parts.map((part) => {
    const firstBreak = part.indexOf("\n");
    const name = firstBreak === -1 ? part.trim() : part.slice(0, firstBreak).trim();
    const body = firstBreak === -1 ? "" : part.slice(firstBreak + 1);
    return {
      heading: extractSingleLineValue(body, "Heading") || name,
      keyPoints: extractBulletList(body, "Key points"),
      labels: extractBulletList(body, "Labels"),
      visualHints: extractBulletList(body, "Visual hint"),
    };
  }).filter((section) => section.heading);
}

function buildPrompt(content: string, args: CliArgs): string {
  const title = extractBlock(content, "Title") || "Infographic topic";
  const subtitle = extractBlock(content, "Subtitle") || "leave empty";
  const learningObjective = extractBulletList(extractBlock(content, "Learning Objective"), "").join("; ");
  const dataIntegrity = extractBulletList(extractBlock(content, "Data Integrity"), "Must preserve");
  const layoutNotes = extractBlock(content, "Layout Notes");
  const preferredLayout = extractSingleLineValue(layoutNotes, "Preferred layout") || "bento-grid";
  const whitespacePriority = extractSingleLineValue(layoutNotes, "Whitespace priority");
  const sections = parseSections(content);

  const sectionLines = sections.map((section, index) => {
    const keyPoints = section.keyPoints.join("; ");
    return `- Section ${index + 1}: ${section.heading} | key points: ${keyPoints}`;
  }).join("\n");

  const labelHints = sections
    .flatMap((section) => section.labels)
    .filter(Boolean)
    .slice(0, 8)
    .join("; ");
  const visualHints = sections
    .flatMap((section) => section.visualHints)
    .filter(Boolean)
    .slice(0, 6)
    .join("; ");

  const promptLines = [
    `Create a high-density infographic about: ${title}.`,
    "",
    `Target audience: ${args.audience}.`,
    `Layout: ${preferredLayout}.`,
    `Style direction: ${args.style}.`,
    `Aspect ratio: ${args.aspect}.`,
    `Language: ${args.language}.`,
    "",
    `Main title: ${title}.`,
    `Subtitle: ${subtitle || "leave empty"}.`,
    learningObjective ? `Learning objective: ${learningObjective}.` : "",
    "",
    "Sections:",
    sectionLines || "- Section 1: Overview | key points: point 1; point 2; point 3",
    "",
    "Visual requirements:",
    "- clear hierarchy",
    "- readable labels",
    "- consistent icon style",
    "- enough whitespace for clarity",
    "- strong section separation",
    whitespacePriority ? `- whitespace priority: ${whitespacePriority}` : "",
    labelHints ? `- labels to preserve: ${labelHints}` : "",
    visualHints ? `- visual hints: ${visualHints}` : "",
    dataIntegrity.length ? `- preserve exact data: ${dataIntegrity.join("; ")}` : "",
    "",
    "Avoid: cluttered background, unreadable tiny text, broken charts, random decorative noise, watermark.",
  ].filter(Boolean);

  return promptLines.join("\n") + "\n";
}

async function main(): Promise<void> {
  const args = parseArgs(process.argv.slice(2));
  if (args.help) {
    printUsage();
    return;
  }
  if (!args.structuredContentPath) {
    console.error("Error: --structured-content is required");
    process.exit(1);
  }
  if (!args.outputPath) {
    console.error("Error: --output is required");
    process.exit(1);
  }

  const content = await readFile(path.resolve(args.structuredContentPath), "utf8");
  const prompt = buildPrompt(content, args);
  const outputPath = path.resolve(args.outputPath);
  await mkdir(path.dirname(outputPath), { recursive: true });
  await writeFile(outputPath, prompt, "utf8");
  console.log(`Prompt written: ${outputPath}`);
}

main().catch((error) => {
  console.error(error instanceof Error ? error.message : String(error));
  process.exit(1);
});
