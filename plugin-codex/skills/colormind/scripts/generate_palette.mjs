#!/usr/bin/env node

function usage() {
  console.error(`Usage:
  generate_palette.mjs --model <name> [--pretty] [--input <5 slots...>]

Slots: "r,g,b" or N
Examples:
  generate_palette.mjs --model default
  generate_palette.mjs --model ui --input "44,43,44" "90,83,82" N N N
`);
  process.exit(2);
}

function rgbToHex([r, g, b]) {
  const clamp = (x) => Math.max(0, Math.min(255, Number(x) || 0));
  const to2 = (n) => clamp(n).toString(16).padStart(2, "0");
  return `#${to2(r)}${to2(g)}${to2(b)}`.toUpperCase();
}

function parseSlot(s) {
  if (!s) return "N";
  if (s === "N" || s === "n") return "N";
  const parts = s.split(",").map((x) => x.trim());
  if (parts.length !== 3) throw new Error(`Bad slot: ${s} (expected r,g,b or N)`);
  const rgb = parts.map((x) => {
    const n = Number.parseInt(x, 10);
    if (!Number.isFinite(n)) throw new Error(`Bad number in slot: ${s}`);
    return n;
  });
  return rgb;
}

const args = process.argv.slice(2);
if (args.length === 0 || args.includes("-h") || args.includes("--help")) usage();

let model = null;
let pretty = false;
let inputSlots = null;

for (let i = 0; i < args.length; i++) {
  const a = args[i];
  if (a === "--model") {
    model = args[i + 1];
    i++;
    continue;
  }
  if (a === "--pretty") {
    pretty = true;
    continue;
  }
  if (a === "--input") {
    const slots = args.slice(i + 1);
    if (slots.length !== 5) {
      console.error(`--input expects exactly 5 slots; got ${slots.length}`);
      usage();
    }
    inputSlots = slots;
    break;
  }
  console.error(`Unknown arg: ${a}`);
  usage();
}

if (!model) usage();

const body = { model };
if (inputSlots) {
  body.input = inputSlots.map(parseSlot);
}

const resp = await fetch("http://colormind.io/api/", {
  method: "POST",
  headers: { "Content-Type": "application/json" },
  body: JSON.stringify(body),
});

if (!resp.ok) {
  const text = await resp.text().catch(() => "");
  console.error(`Colormind api failed (${resp.status}): ${text}`);
  process.exit(1);
}

const data = await resp.json();
console.log(JSON.stringify(data, null, 2));

if (pretty) {
  const pal = data?.result ?? [];
  if (Array.isArray(pal) && pal.length) {
    console.log("\n---\n");
    console.log("## Palette\n");
    for (const rgb of pal) {
      if (!Array.isArray(rgb) || rgb.length !== 3) continue;
      console.log(`- ${rgbToHex(rgb)}  (${rgb.join(", ")})`);
    }
  }
}
