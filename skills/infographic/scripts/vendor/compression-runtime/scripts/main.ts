#!/usr/bin/env bun
import { spawn } from "node:child_process";
import {
  existsSync,
  mkdirSync,
  readdirSync,
  renameSync,
  statSync,
} from "node:fs";
import {
  basename,
  dirname,
  extname,
  join,
  relative,
  resolve,
} from "node:path";

type Compressor = "cwebp" | "magick" | "convert" | "sips";
type Format = "webp" | "png" | "jpeg";

interface CliArgs {
  input: string;
  output?: string;
  format: Format;
  quality: number;
  keep: boolean;
  recursive: boolean;
  json: boolean;
}

interface CompressResult {
  input: string;
  output: string;
  inputSize: number;
  outputSize: number;
  ratio: number;
  reductionPercent: number;
  compressor: Compressor;
}

interface BatchSummary {
  totalFiles: number;
  processedFiles: number;
  failedFiles: number;
  totalInputSize: number;
  totalOutputSize: number;
  ratio: number;
  reductionPercent: number;
  compressor: Compressor;
}

const SUPPORTED_EXTS = new Set([".png", ".jpg", ".jpeg", ".webp", ".gif", ".tif", ".tiff"]);
const ALLOWED_COMMANDS = new Set(["which", "cwebp", "magick", "convert", "sips"]);


function printHelp(): void {
  console.log(`Usage: bun scripts/main.ts <input> [options]

Options:
  -o, --output <path>   Output file or directory
  -f, --format <fmt>    Output format: webp, png, jpeg (default: webp)
  -q, --quality <n>     Quality 0-100 (default: 80)
  -k, --keep            Keep original file when output extension matches input
  -r, --recursive       Process directories recursively
      --json            JSON output
  -h, --help            Show help`);
}

function parseArgs(argv: string[]): CliArgs | null {
  const args: CliArgs = {
    input: "",
    format: "webp",
    quality: 80,
    keep: false,
    recursive: false,
    json: false,
  };

  for (let i = 0; i < argv.length; i++) {
    const arg = argv[i];
    if (arg === "-h" || arg === "--help") {
      printHelp();
      process.exit(0);
    }
    if (arg === "-o" || arg === "--output") {
      args.output = argv[++i];
      continue;
    }
    if (arg === "-f" || arg === "--format") {
      const value = argv[++i]?.toLowerCase();
      if (value === "webp" || value === "png" || value === "jpeg" || value === "jpg") {
        args.format = value === "jpg" ? "jpeg" : value;
        continue;
      }
      console.error(`Invalid format: ${value}`);
      return null;
    }
    if (arg === "-q" || arg === "--quality") {
      const value = Number.parseInt(argv[++i] ?? "", 10);
      if (!Number.isFinite(value) || value < 0 || value > 100) {
        console.error(`Invalid quality: ${argv[i]}`);
        return null;
      }
      args.quality = value;
      continue;
    }
    if (arg === "-k" || arg === "--keep") {
      args.keep = true;
      continue;
    }
    if (arg === "-r" || arg === "--recursive") {
      args.recursive = true;
      continue;
    }
    if (arg === "--json") {
      args.json = true;
      continue;
    }
    if (!arg.startsWith("-") && !args.input) {
      args.input = arg;
      continue;
    }
  }

  if (!args.input) {
    console.error("Input file or directory required.");
    printHelp();
    return null;
  }

  return args;
}

function extForFormat(format: Format): string {
  return format === "jpeg" ? ".jpg" : `.${format}`;
}

function formatSize(bytes: number): string {
  if (bytes < 1024) return `${bytes}B`;
  if (bytes < 1024 * 1024) return `${Math.round(bytes / 1024)}KB`;
  return `${(bytes / (1024 * 1024)).toFixed(1)}MB`;
}

function formatChange(percent: number): string {
  if (percent >= 0) {
    return `${percent}% reduction`;
  }
  return `${Math.abs(percent)}% larger`;
}

function isSupportedImage(pathname: string): boolean {
  return SUPPORTED_EXTS.has(extname(pathname).toLowerCase());
}

function commandExists(command: string): Promise<boolean> {
  return new Promise((resolvePromise) => {
    const child = spawn("which", [command], { stdio: ["ignore", "ignore", "ignore"] });
    child.on("close", (code) => resolvePromise(code === 0));
    child.on("error", () => resolvePromise(false));
  });
}

async function detectCompressor(format: Format): Promise<Compressor> {
  if (format === "webp" && (await commandExists("cwebp"))) {
    return "cwebp";
  }
  if (format === "webp" && (await commandExists("magick"))) {
    return "magick";
  }
  if (format === "webp" && (await commandExists("convert"))) {
    return "convert";
  }
  if (format === "webp") {
    throw new Error("WebP output requires cwebp or ImageMagick on this machine.");
  }
  if (await commandExists("magick")) {
    return "magick";
  }
  if (await commandExists("convert")) {
    return "convert";
  }
  if (await commandExists("sips")) {
    return "sips";
  }
  throw new Error("No supported compressor found. Install cwebp, ImageMagick, or use macOS sips.");
}

function runCommand(command: string, args: string[]): Promise<void> {
  if (!ALLOWED_COMMANDS.has(command)) {
    return Promise.reject(new Error(`Blocked command: ${command}`));
  }
  return new Promise((resolvePromise, rejectPromise) => {
    const child = spawn(command, args, {
      stdio: ["ignore", "ignore", "pipe"],
    });
    let stderr = "";
    child.stderr.on("data", (chunk) => {
      stderr += chunk.toString();
    });
    child.on("close", (code) => {
      if (code === 0) {
        resolvePromise();
        return;
      }
      rejectPromise(new Error(stderr.trim() || `${command} exited with code ${code ?? 1}`));
    });
    child.on("error", (error) => rejectPromise(error));
  });
}

async function compressWithCwebp(input: string, output: string, quality: number): Promise<void> {
  await runCommand("cwebp", ["-q", String(quality), input, "-o", output]);
}

async function compressWithMagick(command: "magick" | "convert", input: string, output: string, quality: number): Promise<void> {
  if (command === "magick") {
    await runCommand("magick", [input, "-quality", String(quality), output]);
    return;
  }
  await runCommand("convert", [input, "-quality", String(quality), output]);
}

async function compressWithSips(input: string, output: string, format: Format, quality: number): Promise<void> {
  const formatArg = format === "jpeg" ? "jpeg" : format;
  await runCommand("sips", [
    "-s",
    "format",
    formatArg,
    "-s",
    "formatOptions",
    String(quality),
    input,
    "--out",
    output,
  ]);
}

async function compressFile(
  compressor: Compressor,
  input: string,
  output: string,
  format: Format,
  quality: number
): Promise<void> {
  if (compressor === "cwebp") {
    if (format === "webp") {
      await compressWithCwebp(input, output, quality);
      return;
    }
    const fallback = (await commandExists("magick"))
      ? "magick"
      : (await commandExists("convert"))
        ? "convert"
        : (await commandExists("sips"))
          ? "sips"
          : null;
    if (!fallback) {
      throw new Error(`No fallback encoder found for ${format}.`);
    }
    if (fallback === "magick" || fallback === "convert") {
      await compressWithMagick(fallback, input, output, quality);
      return;
    }
    await compressWithSips(input, output, format, quality);
    return;
  }
  if (compressor === "magick" || compressor === "convert") {
    await compressWithMagick(compressor, input, output, quality);
    return;
  }
  await compressWithSips(input, output, format, quality);
}

function collectFiles(root: string, recursive: boolean): string[] {
  const items: string[] = [];
  for (const entry of readdirSync(root, { withFileTypes: true })) {
    const fullPath = join(root, entry.name);
    if (entry.isDirectory()) {
      if (recursive) {
        items.push(...collectFiles(fullPath, recursive));
      }
      continue;
    }
    if (entry.isFile() && isSupportedImage(fullPath)) {
      items.push(fullPath);
    }
  }
  return items;
}

function resolveOutputPath(inputFile: string, opts: CliArgs, sourceRoot?: string): string {
  const inputExt = extname(inputFile).toLowerCase();
  const outputExt = extForFormat(opts.format);

  if (sourceRoot && opts.output) {
    const outputRoot = resolve(opts.output);
    const relativePath = relative(sourceRoot, inputFile);
    const relativeDir = dirname(relativePath);
    const base = basename(inputFile, extname(inputFile));
    return join(outputRoot, relativeDir, `${base}${outputExt}`);
  }

  if (opts.output) {
    return resolve(opts.output);
  }

  const base = basename(inputFile, extname(inputFile));
  const sameExt = inputExt === outputExt;
  if (sameExt && opts.keep) {
    return join(dirname(inputFile), `${base}-compressed${outputExt}`);
  }
  if (sameExt) {
    return inputFile;
  }
  return join(dirname(inputFile), `${base}${outputExt}`);
}

function ensureParentDir(pathname: string): void {
  mkdirSync(dirname(pathname), { recursive: true });
}

async function processOne(
  compressor: Compressor,
  inputFile: string,
  opts: CliArgs,
  sourceRoot?: string
): Promise<CompressResult> {
  const input = resolve(inputFile);
  const output = resolveOutputPath(input, opts, sourceRoot);
  const tempOutput = `${output}.tmp`;
  const inputSize = statSync(input).size;

  ensureParentDir(output);
  await compressFile(compressor, input, tempOutput, opts.format, opts.quality);

  const outputSize = statSync(tempOutput).size;
  if (existsSync(output) && output !== input) {
    renameSync(output, `${output}.backup-${Date.now()}`);
  }
  renameSync(tempOutput, output);

  const ratio = inputSize > 0 ? outputSize / inputSize : 0;
  return {
    input,
    output,
    inputSize,
    outputSize,
    ratio,
    reductionPercent: Math.round((1 - ratio) * 100),
    compressor,
  };
}

async function main(): Promise<void> {
  const args = parseArgs(process.argv.slice(2));
  if (!args) {
    process.exit(1);
  }

  const inputPath = resolve(args.input);
  if (!existsSync(inputPath)) {
    console.error(`Not found: ${inputPath}`);
    process.exit(1);
  }

  const compressor = await detectCompressor(args.format);
  const isDirectory = statSync(inputPath).isDirectory();

  if (!isDirectory) {
    if (!isSupportedImage(inputPath)) {
      console.error(`Unsupported image type: ${inputPath}`);
      process.exit(1);
    }
    const result = await processOne(compressor, inputPath, args);
    if (args.json) {
      console.log(JSON.stringify(result, null, 2));
      return;
    }
    console.log(
      `${result.input} -> ${result.output} (${formatSize(result.inputSize)} -> ${formatSize(result.outputSize)}, ${formatChange(result.reductionPercent)}, ${result.compressor})`
    );
    return;
  }

  const files = collectFiles(inputPath, args.recursive);
  if (files.length === 0) {
    console.error("No supported images found.");
    process.exit(1);
  }

  const results: CompressResult[] = [];
  const failures: Array<{ input: string; error: string }> = [];

  for (const file of files) {
    try {
      const result = await processOne(compressor, file, args, inputPath);
      results.push(result);
      if (!args.json) {
        console.log(
          `${result.input} -> ${result.output} (${formatSize(result.inputSize)} -> ${formatSize(result.outputSize)}, ${formatChange(result.reductionPercent)}, ${result.compressor})`
        );
      }
    } catch (error) {
      failures.push({
        input: file,
        error: error instanceof Error ? error.message : String(error),
      });
      if (!args.json) {
        console.error(`Failed: ${file} (${failures.at(-1)?.error ?? "unknown error"})`);
      }
    }
  }

  const totalInputSize = results.reduce((sum, item) => sum + item.inputSize, 0);
  const totalOutputSize = results.reduce((sum, item) => sum + item.outputSize, 0);
  const ratio = totalInputSize > 0 ? totalOutputSize / totalInputSize : 0;
  const summary: BatchSummary = {
    totalFiles: files.length,
    processedFiles: results.length,
    failedFiles: failures.length,
    totalInputSize,
    totalOutputSize,
    ratio,
    reductionPercent: Math.round((1 - ratio) * 100),
    compressor,
  };

  if (args.json) {
    console.log(JSON.stringify({ files: results, failures, summary }, null, 2));
    return;
  }

  console.log(
    `Processed ${summary.processedFiles}/${summary.totalFiles} files: ${formatSize(summary.totalInputSize)} -> ${formatSize(summary.totalOutputSize)} (${formatChange(summary.reductionPercent)}, ${summary.compressor})`
  );
  if (summary.failedFiles > 0) {
    process.exitCode = 1;
  }
}

main().catch((error) => {
  console.error(error instanceof Error ? error.message : String(error));
  process.exit(1);
});
