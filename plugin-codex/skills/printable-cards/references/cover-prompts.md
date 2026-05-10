# Front Cover Image Generation

## Prompt Guidelines

Generate front cover images at **3:4 aspect ratio** using `generate_image`.

### Feminine Cover Prompt Template

> A clean, elegant birthday card front cover design for a girl/woman named {NAME}. Soft {COLOR_DESCRIPTION} watercolor background. A beautiful {ELEMENT} wreath framing elegant hand-lettered calligraphy that reads "Happy Birthday {NAME}" in {TEXT_COLOR} with subtle gold accents. Delicate, airy, refined. Flat 2D print design, no shadows, no borders, no 3D effects, no table surface.

Elements: floral wreath, botanical wreath, peony and rose wreath, eucalyptus and wildflower wreath

### Masculine Cover Prompt Template

> A premium masculine birthday card front cover. Deep {COLOR} textured paper background with a subtle linen weave texture. Elegant gold foil hand-lettered calligraphy reads "Happy Birthday {NAME}" in a refined modern script style. Thin gold foil horizontal lines above and below text. Small olive branch illustration in gold foil beneath the name. Sophisticated, understated, luxury stationery feel. No sparkles, no stars, no geometric patterns. Flat print design, no 3D effects, no shadows, no borders.

## Post-Processing

AI-generated images often include unwanted artifacts (table surfaces, shadows, rounded corners). After generation, run `clean_cover_image.py` to crop to card content only:

```bash
python scripts/clean_cover_image.py raw_cover.png clean_cover.png
```

Verify the cleaned image visually before using it in the card.
