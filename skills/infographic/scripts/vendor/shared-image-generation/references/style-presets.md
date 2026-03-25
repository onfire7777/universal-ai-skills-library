# Style Presets (`image-generation`)

These presets are **prompt-level style presets**, not native gateway API enums. Each preset appends a style description to the original prompt so the agent and the user can reuse common visual directions more consistently.

Useful when:

- the user only knows a rough direction such as "photoreal", "manga", "poster", or "infographic"
- the same visual language should be reused across a batch
- `--dry-run` is being used to inspect the final assembled request body

Not ideal when:

- the user already provided a very explicit and sufficiently detailed style prompt
- the workflow must rely on model-specific private style tokens

## Available Presets

| preset | User Might Say | Best For | Injected Direction |
| --- | --- | --- | --- |
| `photoreal` | "realistic", "photo-like", "natural" | realistic people, products, environments | photorealistic, natural lighting, realistic texture, high detail |
| `cinematic` | "cinematic", "movie feel", "moody" | cinematic posters, moody scenes | cinematic composition, dramatic lighting, film still, moody atmosphere |
| `anime` | "anime", "colorful cartoon" | anime characters, light-novel covers | anime style, clean linework, vibrant colors, expressive character design |
| `manga` | "manga", "black and white comic" | comic panels, B&W storytelling | manga style, dynamic ink lines, screentone shading, black and white emphasis |
| `watercolor` | "watercolor", "hand-drawn", "soft" | soft illustration, book-style imagery | watercolor illustration, soft brush texture, paper grain, gentle palette |
| `flat-illustration` | "flat", "minimal", "notion-style", "clean" | product illustration, knowledge cards | flat illustration, simplified shapes, clean vector feel, limited palette |
| `3d-render` | "3D", "rendered" | 3D characters, product rendering | 3d render, volumetric lighting, material detail, polished surfaces |
| `poster` | "poster", "bold", "high impact" | promo images, event pages, covers | poster design, bold focal point, strong typography space, high contrast |
| `editorial` | "magazine", "premium", "elegant layout" | magazine-like, brand-oriented visuals | editorial illustration, refined composition, premium layout feel, restrained palette |
| `infographic` | "infographic", "chart", "data visual" | infographics, knowledge cards | infographic style, clear hierarchy, icon-friendly layout, data visualization feel |
| `chalk` | "chalkboard", "hand-written" | chalkboard, hand-drawn tutorials | chalk drawing, textured strokes, blackboard feel, hand-drawn annotations |
| `ink-brush` | "ink painting", "brush", "East Asian" | ink-painting, East Asian aesthetics | ink brush painting, expressive brushwork, elegant negative space, East Asian aesthetic |

When the user describes a style in natural language, match their description to the closest preset. If unsure between two presets, briefly describe both and ask which feels closer. Run `npm run style-menu` from this skill directory to see the full preset list with localized labels.

## Usage

```bash
${BUN_X} {baseDir}/scripts/main.ts --prompt "a cat" --style photoreal --image out.png -m "<model>"
${BUN_X} {baseDir}/scripts/main.ts --prompt "AI trend summary" --style infographic --image out.png -m "<model>"
${BUN_X} {baseDir}/scripts/main.ts --prompt "girl standing on a rainy night street" --style cinematic --dry-run -m "<model>"
```

## Notes

- `--style` only adds style direction. It does not replace `model`.
- The style preset participates in prompt length validation because it is appended to the final request prompt.
- If the user's original prompt is already clear and detailed, respect that first. The preset is only a lightweight enhancement.
