import fs from "node:fs";
import path from "node:path";
import process from "node:process";
import { fileURLToPath } from "node:url";

const VALID_WORKFLOWS = new Set(["general", "cover", "rednote", "infographic", "comic", "article"]);

const FIELD_OPTIONS = {
  styleFamily: ["photoreal", "editorial", "flat-illustration", "poster", "manga", "watercolor", "infographic"],
  palette: ["cool", "warm", "neutral", "high-contrast", "pastel", "dark-cinematic"],
  lighting: ["soft", "natural", "dramatic", "studio", "flat", "high-key"],
  composition: ["single-hero", "balanced", "layered-depth", "grid", "bento", "panel-storytelling"],
  shot: ["close-up", "medium", "wide", "top-down", "isometric", "straight-on"],
  textDensity: ["none", "title-only", "light-labels", "text-rich"],
  aspect: ["1:1", "3:4", "4:3", "16:9", "9:16"],
  referenceMode: ["none", "style-reference", "character-reference", "layout-reference"],
};

const WORKFLOW_DEFAULTS = {
  general: {
    styleFamily: "editorial",
    composition: "balanced",
    shot: "medium",
    aspect: "1:1",
    textDensity: "none",
  },
  cover: {
    styleFamily: "editorial",
    composition: "single-hero",
    shot: "medium",
    aspect: "16:9",
    textDensity: "title-only",
  },
  rednote: {
    styleFamily: "flat-illustration",
    composition: "bento",
    shot: "straight-on",
    aspect: "3:4",
    textDensity: "light-labels",
  },
  infographic: {
    styleFamily: "infographic",
    composition: "grid",
    shot: "straight-on",
    aspect: "3:4",
    textDensity: "text-rich",
  },
  comic: {
    styleFamily: "manga",
    composition: "panel-storytelling",
    shot: "medium",
    aspect: "3:4",
    textDensity: "light-labels",
  },
  article: {
    styleFamily: "editorial",
    composition: "balanced",
    shot: "wide",
    aspect: "16:9",
    textDensity: "light-labels",
  },
};

function printHelp() {
  console.log(`Usage:
  node scripts/build-visual-brief.mjs --workflow <name> --topic <text> [options]

Options:
  --workflow <name>       general | cover | rednote | infographic | comic | article
  --topic <text>          Core topic or request
  --subject <text>        Main subject
  --goal <text>           Intended effect or communication goal
  --style-family <name>   Visual family
  --palette <name>        Color temperature / palette direction
  --lighting <name>       Lighting direction
  --composition <name>    Composition pattern
  --shot <name>           Shot language
  --aspect <ratio>        Aspect ratio
  --lang <code>           On-image text language when needed
  --text-density <name>   none | title-only | light-labels | text-rich
  --reference-mode <name> none | style-reference | character-reference | layout-reference
  --output <path>         Write a markdown brief to this path
  --json                  Print JSON output
  -h, --help              Show help`);
}

function parseArgs(argv) {
  const args = {
    workflow: "general",
    topic: null,
    subject: null,
    goal: null,
    styleFamily: null,
    palette: null,
    lighting: null,
    composition: null,
    shot: null,
    aspect: null,
    language: "en",
    textDensity: null,
    referenceMode: null,
    output: null,
    json: false,
  };

  for (let i = 0; i < argv.length; i++) {
    const current = argv[i];
    if (current === "--workflow") args.workflow = argv[++i] ?? args.workflow;
    else if (current === "--topic") args.topic = argv[++i] ?? null;
    else if (current === "--subject") args.subject = argv[++i] ?? null;
    else if (current === "--goal") args.goal = argv[++i] ?? null;
    else if (current === "--style-family") args.styleFamily = argv[++i] ?? null;
    else if (current === "--palette") args.palette = argv[++i] ?? null;
    else if (current === "--lighting") args.lighting = argv[++i] ?? null;
    else if (current === "--composition") args.composition = argv[++i] ?? null;
    else if (current === "--shot") args.shot = argv[++i] ?? null;
    else if (current === "--aspect") args.aspect = argv[++i] ?? null;
    else if (current === "--lang") args.language = argv[++i] ?? args.language;
    else if (current === "--text-density") args.textDensity = argv[++i] ?? null;
    else if (current === "--reference-mode") args.referenceMode = argv[++i] ?? null;
    else if (current === "--output") args.output = argv[++i] ?? null;
    else if (current === "--json") args.json = true;
    else if (current === "--help" || current === "-h") {
      printHelp();
      process.exit(0);
    }
  }

  if (!VALID_WORKFLOWS.has(args.workflow)) throw new Error(`Invalid workflow: ${args.workflow}`);
  if (!args.topic?.trim()) throw new Error("--topic is required");
  return args;
}

const QUESTIONS_I18N = {
  en: {
    subject:       { prompt: "What is the main visual subject?", options: ["person", "product", "scene", "diagram/structure", "abstract concept"], recommended: "scene" },
    goal:          { prompt: "What is the primary purpose of this image?", options: ["attract clicks", "explain a concept", "set the mood", "highlight contrast", "drive action"], recommended: "explain a concept" },
    styleFamily:   { prompt: "More photorealistic or illustration?" },
    palette:       { prompt: "What color temperature / palette direction?" },
    lighting:      { prompt: "What lighting or atmosphere?" },
    composition:   { prompt: "What composition pattern works best?" },
    shot:          { prompt: "What shot language fits best?" },
    aspect:        { prompt: "What aspect ratio?" },
    textDensity:   { prompt: "How much text should appear on the image?" },
    referenceMode: { prompt: "Do you have a reference image?" },
  },
  zh: {
    subject:       { prompt: "主视觉主体是什么？", options: ["人物", "产品", "场景", "图表/结构", "抽象概念"], recommended: "场景" },
    goal:          { prompt: "这张图最重要的任务是什么？", options: ["吸引点击", "解释概念", "建立氛围", "突出对比", "引导行动"], recommended: "解释概念" },
    styleFamily:   { prompt: "更偏写实还是插画？" },
    palette:       { prompt: "色彩温度想要什么方向？" },
    lighting:      { prompt: "光线/气氛更偏哪种？" },
    composition:   { prompt: "构图组织更适合哪种？" },
    shot:          { prompt: "镜头语言更适合哪种？" },
    aspect:        { prompt: "最终比例是什么？" },
    textDensity:   { prompt: "画面里文字密度要多高？" },
    referenceMode: { prompt: "是否有参考图？" },
  },
};

function question(id, prompt, options, recommended) {
  return { id, prompt, options, recommended };
}

function getI18nQuestion(lang, id) {
  const table = QUESTIONS_I18N[lang] ?? QUESTIONS_I18N.en;
  return table[id] ?? QUESTIONS_I18N.en[id];
}

export function buildVisualBrief(input) {
  const lang = (input.language?.trim() || "en").startsWith("zh") ? "zh" : "en";
  const defaults = WORKFLOW_DEFAULTS[input.workflow] ?? WORKFLOW_DEFAULTS.general;
  const brief = {
    workflow: input.workflow,
    topic: input.topic.trim(),
    subject: input.subject?.trim() || null,
    goal: input.goal?.trim() || null,
    styleFamily: input.styleFamily?.trim() || defaults.styleFamily,
    palette: input.palette?.trim() || null,
    lighting: input.lighting?.trim() || null,
    composition: input.composition?.trim() || defaults.composition,
    shot: input.shot?.trim() || defaults.shot,
    aspect: input.aspect?.trim() || defaults.aspect,
    language: input.language?.trim() || "en",
    textDensity: input.textDensity?.trim() || defaults.textDensity,
    referenceMode: input.referenceMode?.trim() || "none",
  };

  const missingFields = [];
  if (!brief.subject) missingFields.push("subject");
  if (!brief.goal) missingFields.push("goal");
  if (!input.palette) missingFields.push("palette");
  if (!input.lighting) missingFields.push("lighting");

  const q = (id, fallbackOptions, fallbackRecommended) => {
    const i18n = getI18nQuestion(lang, id);
    return question(id, i18n.prompt, i18n.options ?? fallbackOptions, i18n.recommended ?? fallbackRecommended);
  };

  const questions = [];
  if (!brief.subject) questions.push(q("subject"));
  if (!brief.goal) questions.push(q("goal"));
  if (!input.styleFamily) questions.push(q("styleFamily", FIELD_OPTIONS.styleFamily, defaults.styleFamily));
  if (!input.palette) questions.push(q("palette", FIELD_OPTIONS.palette, "neutral"));
  if (!input.lighting) questions.push(q("lighting", FIELD_OPTIONS.lighting, "soft"));
  if (!input.composition) questions.push(q("composition", FIELD_OPTIONS.composition, defaults.composition));
  if (!input.shot) questions.push(q("shot", FIELD_OPTIONS.shot, defaults.shot));
  if (!input.aspect) questions.push(q("aspect", FIELD_OPTIONS.aspect, defaults.aspect));
  if (!input.textDensity) questions.push(q("textDensity", FIELD_OPTIONS.textDensity, defaults.textDensity));
  if (!input.referenceMode) questions.push(q("referenceMode", FIELD_OPTIONS.referenceMode, "none"));

  const promptDirectives = [
    `Main topic: ${brief.topic}.`,
    `Primary subject: ${brief.subject || "<need subject>"}.`,
    `Communication goal: ${brief.goal || "<need goal>"}.`,
    `Style family: ${brief.styleFamily}.`,
    `Palette direction: ${brief.palette || "<need palette>"}.`,
    `Lighting: ${brief.lighting || "<need lighting>"}.`,
    `Composition: ${brief.composition}.`,
    `Shot language: ${brief.shot}.`,
    `Aspect ratio: ${brief.aspect}.`,
    `Language: ${brief.language}.`,
    `Text density: ${brief.textDensity}.`,
    `Reference mode: ${brief.referenceMode}.`,
  ];

  return {
    ok: missingFields.length === 0,
    brief,
    missingFields,
    questions,
    promptDirectives,
  };
}

function renderMarkdown(summary) {
  const { brief, missingFields, questions, promptDirectives } = summary;
  const lines = [
    "# Visual Brief",
    "",
    "## Resolved Brief",
    `- Workflow: ${brief.workflow}`,
    `- Topic: ${brief.topic}`,
    `- Subject: ${brief.subject ?? "<need subject>"}`,
    `- Goal: ${brief.goal ?? "<need goal>"}`,
    `- Style family: ${brief.styleFamily}`,
    `- Palette: ${brief.palette ?? "<need palette>"}`,
    `- Lighting: ${brief.lighting ?? "<need lighting>"}`,
    `- Composition: ${brief.composition}`,
    `- Shot: ${brief.shot}`,
    `- Aspect ratio: ${brief.aspect}`,
    `- Language: ${brief.language}`,
    `- Text density: ${brief.textDensity}`,
    `- Reference mode: ${brief.referenceMode}`,
    "",
    "## Missing Fields",
    ...(missingFields.length ? missingFields.map((field) => `- ${field}`) : ["- none"]),
    "",
    "## Questions To Ask",
    ...(questions.length
      ? questions.flatMap((item) => [
          `- ${item.prompt}`,
          `  options: ${item.options.join(" | ")}`,
          `  recommended: ${item.recommended}`,
        ])
      : ["- none"]),
    "",
    "## Prompt Directives",
    ...promptDirectives.map((line) => `- ${line}`),
    "",
  ];
  return `${lines.join("\n")}`;
}

function writeBrief(outputPath, summary) {
  const resolved = path.resolve(outputPath);
  fs.mkdirSync(path.dirname(resolved), { recursive: true });
  fs.writeFileSync(resolved, renderMarkdown(summary), "utf8");
  return resolved;
}

const isDirectRun =
  process.argv[1] && path.resolve(process.argv[1]) === path.resolve(fileURLToPath(import.meta.url));

if (isDirectRun) {
  try {
    const args = parseArgs(process.argv.slice(2));
    const summary = buildVisualBrief(args);
    const outputPath = args.output ? writeBrief(args.output, summary) : null;
    if (args.json) {
      console.log(JSON.stringify({ ...summary, outputPath }, null, 2));
    } else {
      console.log(renderMarkdown(summary));
      if (outputPath) console.log(`Brief written: ${outputPath}`);
    }
  } catch (error) {
    console.error(error instanceof Error ? error.message : String(error));
    process.exit(1);
  }
}
