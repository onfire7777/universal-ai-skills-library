# Style Presets Reference

Each preset defines a complete set of SVG attributes that ensure visual consistency across an icon set. Pick a preset as a starting point, then adjust individual parameters if the project calls for it.

## Preset: Clean

The safe default. Works for most business websites, SaaS dashboards, and professional services.

```json
{
  "grid": 24,
  "strokeWidth": 1.5,
  "strokeLinecap": "round",
  "strokeLinejoin": "round",
  "cornerRadius": 2,
  "padding": 2,
  "opticalBalance": true
}
```

**Character**: Balanced, professional, approachable without being casual. The round caps and joins soften the icons without making them feel childish. The 1.5px stroke is the sweet spot — visible without being heavy.

**Use when**: You don't have strong reasons to choose something else. This is the "can't go wrong" option.

**Drawing range**: With 24px grid and 2px padding, draw within coordinates 2–22.

---

## Preset: Sharp

Corporate, technical, precise. Good for law firms, engineering companies, fintech.

```json
{
  "grid": 24,
  "strokeWidth": 1.5,
  "strokeLinecap": "square",
  "strokeLinejoin": "miter",
  "cornerRadius": 0,
  "padding": 2,
  "opticalBalance": true
}
```

**Character**: Authoritative and precise. Square caps create definitive endpoints. Miter joins make clean 90° corners. Zero corner radius means rectangles are sharp.

**Use when**: The brand identity is formal, technical, or conveys precision and authority.

**Watch out for**: Miter joins can create spiky artifacts at acute angles. If a path has angles less than ~30°, consider using `stroke-linejoin="round"` on that specific element only.

---

## Preset: Soft

Friendly, warm, approachable. Good for childcare, health & wellness, food, community orgs.

```json
{
  "grid": 24,
  "strokeWidth": 2,
  "strokeLinecap": "round",
  "strokeLinejoin": "round",
  "cornerRadius": 4,
  "padding": 2.5,
  "opticalBalance": true
}
```

**Character**: Warm and inviting. The thicker 2px stroke gives more visual presence. Large corner radius makes everything feel rounded and friendly. Extra padding (2.5px) gives breathing room.

**Use when**: The project needs to feel approachable, non-intimidating, or playful.

**Drawing range**: With 2.5px padding, draw within coordinates 2.5–21.5.

---

## Preset: Minimal

Elegant, restrained, editorial. Good for luxury brands, design studios, photography portfolios.

```json
{
  "grid": 24,
  "strokeWidth": 1,
  "strokeLinecap": "round",
  "strokeLinejoin": "round",
  "cornerRadius": 0,
  "padding": 2,
  "opticalBalance": true
}
```

**Character**: Delicate and refined. The 1px stroke is thin — it looks elegant at larger sizes (32px+) but needs careful testing at smaller sizes. Zero corner radius keeps geometry pure.

**Use when**: The design is high-end, editorial, or minimalist. Works best when icons will be displayed at 28px+ size.

**Watch out for**: At 16-20px rendering, 1px strokes can look faint on standard-DPI screens. Recommend this preset only when you know the icons will be used at adequate sizes, or when the site targets retina displays.

---

## Preset: Bold

High impact, accessible. Good for outdoor brands, construction, emergency services, signage.

```json
{
  "grid": 24,
  "strokeWidth": 2.5,
  "strokeLinecap": "round",
  "strokeLinejoin": "round",
  "cornerRadius": 2,
  "padding": 2.5,
  "opticalBalance": true
}
```

**Character**: Strong and unmistakable. The 2.5px stroke is thick enough to read at small sizes and low contrast. Ideal for accessibility-focused projects or brands that want to feel strong and reliable.

**Use when**: Icons need to be legible in challenging conditions — small sizes, low contrast, outdoor signage, or accessibility requirements.

**Drawing range**: With 2.5px padding and thick strokes, the usable drawing area is tighter. Keep paths within 3–21 to avoid strokes bleeding into the padding.

---

## Customising Presets

You can mix parameters. Common adjustments:

- **Clean preset but with square caps**: More structured feel while keeping round joins
- **Soft preset but with 1.5px stroke**: Friendly corners but less heavy
- **Minimal preset but with 1.25px stroke**: Slightly more visible while staying refined
- **Bold preset but with miter joins**: Industrial/construction feel

Always document the customisation in the style-spec.json so the user knows exactly what was used.

---

## Grid Size Notes

24px is the standard and works for almost everything. Consider alternatives only if:

- **20px**: The design system is built on a 20px base unit (uncommon)
- **32px**: Icons will primarily display at large sizes (hero sections, feature lists)
- **16px**: Icons are exclusively for very small UI elements (rare — usually just use 24px and scale down)

Stick with 24px unless there's a specific reason not to.
