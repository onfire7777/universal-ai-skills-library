---
name: printable-cards
description: Create beautiful printable foldable greeting cards (birthday, thank you, etc.) as PDFs. Use when the user asks to make, design, or create a printable card, greeting card, birthday card, or foldable card for someone.
---

# Printable Cards

Create high-quality printable foldable greeting cards rendered as 2-page PDFs via HTML/CSS + WeasyPrint.

## Card Layout

US Letter landscape (11 x 8.5 in), 2 pages, designed for double-sided printing (flip on short edge), fold in half.

| Page | Left Half | Right Half |
|------|-----------|------------|
| 1 (outside) | Back cover (accent bg, border, symbol, label) | Front cover (full-bleed image) |
| 2 (inside) | Open space for handwriting | Header + message + closing |

## Workflow

1. **Gather info**: recipient name(s), relationship/context, occasion, gender/style preference, any personal details for the message
2. **Generate front cover image**: use `generate_image` at 3:4 aspect ratio. See `references/cover-prompts.md` for prompt templates. Feminine = floral/botanical wreath. Masculine = textured background with gold calligraphy.
3. **Clean cover image**: run `scripts/clean_cover_image.py <raw.png> <clean.png>` to remove AI artifacts (shadows, table surfaces, grey borders). Verify visually.
4. **Choose color palette**: see `references/color-palettes.md` for tested palettes. Match palette to the front cover style. Parameters: `accent`, `accent_dark`, `light_bg`, `dark_text`.
5. **Write the message**: write in HTML using `<p>` tags. Save to a `.html` file. Key principles:
   - Write from the correct perspective (who is giving the card)
   - Include specific personal details, not generic compliments
   - Sound human, not like a greeting card company
   - High emotional intelligence: loving, positive, honest
   - 3 paragraphs is the sweet spot
6. **Generate the PDF**: run `scripts/generate_card.py` with all parameters. Use `--masculine` flag for bold sans-serif styling.
7. **Verify**: view both pages of the PDF. Check front cover fills the right half with no gaps, inside text is centered and readable, colors are cohesive between outside and inside.
8. **Deliver**: attach the PDF. Include printing instructions: double-sided, flip on short edge, fold in half.

## Script Usage

```bash
python scripts/generate_card.py \
  --name "Chloe" \
  --front-image /path/to/clean_cover.png \
  --message-file /path/to/message.html \
  --closing "All my love, your little brother" \
  --accent "#C4889B" \
  --accent-dark "#9E6B7E" \
  --light-bg "#FDF2F5" \
  --dark-text "#5A2D45" \
  --output /path/to/Card.pdf \
  [--masculine]
```

## Message Writing Guidelines

The message is the hardest part. Follow these rules:

- **Ask for context first.** Who is giving the card? What is their relationship? What do they love about this person? Any inside jokes or specific things to mention?
- **Write from the giver's voice**, not a generic third person.
- **Be specific.** "The way you always sit with me and really talk" beats "you are always there for me."
- **Stay honest.** Do not over-idealize. If siblings tease each other, acknowledge it warmly.
- **3 paragraphs**: (1) opening with personal connection, (2) what you admire/appreciate with specifics, (3) birthday wish + love.
- **Closing line**: short, warm, identifies the relationship (e.g. "Your favorite brother", "Love always, Mom").

## Dependencies

Ensure `weasyprint` is installed: `pip install weasyprint`

## Common Issues

| Issue | Fix |
|-------|-----|
| Front cover has grey edges/shadows | Run `clean_cover_image.py` more aggressively, or regenerate with "no shadows, no borders, no 3D effects" in prompt |
| Inside background doesn't match outside | Ensure `light_bg` color complements the front cover palette |
| Text looks generic/AI | Rewrite with specific personal details from the user |
| Card looks too feminine/masculine | Adjust palette and toggle `--masculine` flag |
