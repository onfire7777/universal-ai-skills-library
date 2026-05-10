#!/usr/bin/env bun
import process from "node:process";
import { mkdir, writeFile, access } from "node:fs/promises";
import { constants } from "node:fs";
import path from "node:path";

type CliArgs = {
  outputDir: string | null;
  topic: string;
  language: string;
  force: boolean;
  help: boolean;
};

function printUsage(): void {
  console.log(`Usage:
  npx -y bun scripts/scaffold.ts --output-dir infographic/topic-slug --topic "Topic name" [options]

Options:
  --output-dir <path>   Target infographic working directory
  --topic <text>        Main topic placeholder
  --lang <code>         On-image text language when needed (default placeholder: en)
  --force               Overwrite existing files
  -h, --help            Show help`);
}

function parseArgs(argv: string[]): CliArgs {
  const args: CliArgs = {
    outputDir: null,
    topic: "Topic",
    language: "en",
    force: false,
    help: false,
  };

  for (let i = 0; i < argv.length; i++) {
    const current = argv[i]!;
    if (current === "--output-dir") args.outputDir = argv[++i] ?? null;
    else if (current === "--topic") args.topic = argv[++i] ?? args.topic;
    else if (current === "--lang") args.language = argv[++i] ?? args.language;
    else if (current === "--force") args.force = true;
    else if (current === "--help" || current === "-h") args.help = true;
  }

  return args;
}

async function exists(filePath: string): Promise<boolean> {
  try {
    await access(filePath, constants.F_OK);
    return true;
  } catch {
    return false;
  }
}

async function writeScaffoldFile(filePath: string, content: string, force: boolean): Promise<void> {
  if (!force && (await exists(filePath))) {
    throw new Error(`File already exists: ${filePath}. Use --force to overwrite.`);
  }
  await writeFile(filePath, content, "utf8");
}

function analysisTemplate(topic: string, language: string): string {
  return `# Infographic Analysis

## Topic
- Main topic: ${topic}
- Audience: <target audience>
- Language: ${language}

## Content Shape
- Density: <low|medium|high>
- Primary structure: <list|comparison|timeline|hierarchy|process|mixed>
- Data presence: <none|light|heavy>
- Text-in-image needed: <yes|no>

## Core Message
- Main message: <one-sentence summary>
- Supporting points:
  - <point 1>
  - <point 2>
  - <point 3>

## Visual Implications
- Best layout candidates:
  - <layout 1>
  - <layout 2>
- Best style candidates:
  - <style 1>
  - <style 2>
- Recommended aspect:
  - <3:4|4:3|16:9>

## Risks
- <too much density, too much image text, incomplete data, or similar>
`;
}

function structuredContentTemplate(topic: string): string {
  return `# Structured Content

## Title
${topic}

## Subtitle
<subtitle, or leave empty>

## Learning Objective
- <what the reader should understand after seeing the infographic>

## Sections

### Section 1
- Heading: <heading>
- Purpose: <why this section exists>
- Key points:
  - <point 1>
  - <point 2>
  - <point 3>
- Labels:
  - <short labels, terms, or numbers>
- Visual hint:
  - <icons, arrows, modules, timelines, contrast boxes, and so on>

### Section 2
- Heading: <heading>
- Purpose: <purpose>
- Key points:
  - <point 1>
  - <point 2>
- Labels:
  - <short labels>
- Visual hint:
  - <visual hint>

## Data Integrity
- Must preserve:
  - <year>
  - <number>
  - <term>

## Layout Notes
- Preferred layout: <layout>
- Preferred reading order: <left-to-right|top-to-bottom|center-out>
- Whitespace priority: <low|medium|high>
`;
}

function promptTemplate(topic: string, language: string): string {
  return `Create a high-density infographic about: ${topic}.

Target audience: <target audience>.
Layout: <layout>.
Style direction: <style>.
Aspect ratio: <aspect>.
Language: ${language}.

Main title: <title>.
Subtitle: <subtitle or leave empty>.

Sections:
- Section 1: Overview | key points: point 1; point 2; point 3
- Section 2: Breakdown | key points: point 1; point 2; point 3
- Section 3: Checklist | key points: point 1; point 2; point 3

Visual requirements:
- clear hierarchy
- readable labels
- consistent icon style
- enough whitespace for clarity
- strong section separation

Avoid: cluttered background, unreadable tiny text, broken charts, random decorative noise, watermark.
`;
}

async function main(): Promise<void> {
  const args = parseArgs(process.argv.slice(2));
  if (args.help) {
    printUsage();
    return;
  }
  if (!args.outputDir) {
    console.error("Error: --output-dir is required");
    process.exit(1);
  }

  const outputDir = path.resolve(args.outputDir);
  const promptsDir = path.join(outputDir, "prompts");
  await mkdir(promptsDir, { recursive: true });

  await writeScaffoldFile(path.join(outputDir, "analysis.md"), analysisTemplate(args.topic, args.language), args.force);
  await writeScaffoldFile(path.join(outputDir, "structured-content.md"), structuredContentTemplate(args.topic), args.force);
  await writeScaffoldFile(path.join(promptsDir, "infographic.md"), promptTemplate(args.topic, args.language), args.force);

  console.log(`Scaffold created in: ${outputDir}`);
  console.log("Files:");
  console.log("- analysis.md");
  console.log("- structured-content.md");
  console.log("- prompts/infographic.md");
}

main().catch((error) => {
  console.error(error instanceof Error ? error.message : String(error));
  process.exit(1);
});
