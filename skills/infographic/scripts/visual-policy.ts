export type WorkflowKind = "cover" | "rednote" | "infographic" | "comic" | "article";

const STYLE_ALIASES: Record<WorkflowKind, Record<string, string>> = {
  cover: {
    editorial: "editorial",
    poster: "poster",
    cinematic: "cinematic",
    watercolor: "watercolor",
    "flat-vector": "flat-illustration",
    "ink-brush": "ink-brush",
    chalk: "chalk",
    manga: "manga",
    anime: "anime",
    photoreal: "photoreal",
    "3d-render": "3d-render",
    infographic: "infographic",
  },
  rednote: {
    cute: "anime",
    fresh: "flat-illustration",
    warm: "watercolor",
    bold: "poster",
    minimal: "editorial",
    retro: "poster",
    notion: "flat-illustration",
    chalkboard: "chalk",
    editorial: "editorial",
  },
  infographic: {
    "craft-handmade": "watercolor",
    chalkboard: "chalk",
    "corporate-memphis": "flat-illustration",
    "technical-schematic": "infographic",
    "bold-graphic": "poster",
    "storybook-watercolor": "watercolor",
    "retro-pop-grid": "poster",
    minimal: "editorial",
    infographic: "infographic",
  },
  comic: {
    "ligne-claire": "flat-illustration",
    manga: "manga",
    realistic: "photoreal",
    "ink-brush": "ink-brush",
    chalk: "chalk",
  },
  article: {
    minimal: "editorial",
    notion: "flat-illustration",
    blueprint: "infographic",
    watercolor: "watercolor",
    elegant: "editorial",
    poster: "poster",
  },
};

const WORKFLOW_NEGATIVE_PROMPTS: Record<WorkflowKind, string[]> = {
  cover: [
    "blurry",
    "muddy contrast",
    "crowded title area",
    "tiny unreadable text",
    "weak focal hierarchy",
    "cluttered composition",
    "watermark",
    "logo",
    "jpeg artifacts",
  ],
  rednote: [
    "blurry",
    "tiny unreadable words",
    "cluttered layout",
    "weak hierarchy",
    "merged card panels",
    "decorative clutter",
    "watermark",
    "jpeg artifacts",
  ],
  infographic: [
    "blurry",
    "dense tiny text",
    "unreadable labels",
    "weak information hierarchy",
    "random decorative icons",
    "confusing arrows",
    "watermark",
    "jpeg artifacts",
  ],
  comic: [
    "blurry",
    "inconsistent character design",
    "broken panel borders",
    "unreadable speech bubbles",
    "extra limbs",
    "panel clutter",
    "watermark",
    "jpeg artifacts",
  ],
  article: [
    "blurry",
    "generic stock-photo feel",
    "weak focal point",
    "unrelated decorative objects",
    "unreadable labels",
    "cluttered composition",
    "watermark",
    "jpeg artifacts",
  ],
};

function normalizeValue(value: string): string {
  return value.trim().toLowerCase().replace(/[_\s]+/g, "-");
}

export function resolveWorkflowStyle(
  workflow: WorkflowKind,
  explicitStyle: string | null,
  domainStyle: string | undefined | null,
): string | null {
  if (explicitStyle?.trim()) return explicitStyle.trim();
  if (!domainStyle?.trim()) return null;
  const normalized = normalizeValue(domainStyle);
  return STYLE_ALIASES[workflow][normalized] ?? null;
}

export function buildWorkflowNegativePrompt(workflow: WorkflowKind): string {
  return WORKFLOW_NEGATIVE_PROMPTS[workflow].join(", ");
}
