import process from "node:process";

// Keep this menu aligned with image-generation/scripts/main.ts STYLE_PRESETS.
const STYLES = [
  { id: "photoreal", en: "Photorealistic / natural photography", zh: "写实 / 真实摄影质感" },
  { id: "cinematic", en: "Cinematic / film still mood", zh: "电影感 / 剧照氛围" },
  { id: "anime", en: "Anime / colorful cartoon", zh: "日系动漫 / 彩色" },
  { id: "manga", en: "Manga / B&W screentone", zh: "漫画 / 黑白网点" },
  { id: "watercolor", en: "Watercolor illustration", zh: "水彩插画" },
  { id: "flat-illustration", en: "Flat illustration / knowledge graphic", zh: "扁平插画 / 知识图" },
  { id: "3d-render", en: "3D render", zh: "3D 渲染" },
  { id: "poster", en: "Poster design", zh: "海报设计" },
  { id: "editorial", en: "Editorial / magazine layout feel", zh: "杂志插画 / 高级排版感" },
  { id: "infographic", en: "Infographic / chart style", zh: "信息图 / 图表风" },
  { id: "chalk", en: "Chalkboard / hand-drawn", zh: "粉笔黑板" },
  { id: "ink-brush", en: "East Asian ink brush", zh: "国风水墨" },
];

function printHelp() {
  console.log(`Usage:
  node scripts/style-menu.mjs [--lang en|zh]

Options:
  --lang <code>  Output language: en (default) or zh
  -h, --help     Show help

This prints the available --style presets and example commands.`);
}

function parseLang(argv) {
  const idx = argv.indexOf("--lang");
  if (idx >= 0 && argv[idx + 1]) return argv[idx + 1];
  return "en";
}

if (process.argv.includes("-h") || process.argv.includes("--help")) {
  printHelp();
  process.exit(0);
}

const lang = parseLang(process.argv.slice(2));
const labelKey = lang === "zh" ? "zh" : "en";

console.log("Available styles (--style):");
for (const s of STYLES) {
  console.log(`- ${s.id.padEnd(16)} ${s[labelKey]}`);
}

console.log("\nExample:");
console.log("  M=<model-key>");
console.log("  npx -y bun scripts/main.ts \\");
console.log("    --prompt \"a cat wearing sunglasses, neon city at night\" \\");
console.log("    --style cinematic --ar 3:4 --image out.png -m \"$M\"");
