# SVG Examples Reference

These examples show how to construct icons correctly using the Clean preset (24px grid, 1.5px stroke, round caps/joins, 2px padding). Study the patterns before generating your set.

## Key Principles Demonstrated

1. All coordinates stay within the padding zone (2–22 for 24px grid with 2px padding)
2. Coordinates use at most 2 decimal places, preferring whole and half numbers
3. Paths are minimal — no unnecessary points or elements
4. Icons are visually centred, not just mathematically centred

---

## Simple Icon: Chevron Right

A basic directional indicator. Note the optical overshoot — the chevron extends slightly beyond what pure math would suggest, so it feels the same visual weight as boxier icons.

```xml
<svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round">
  <polyline points="9 6 15 12 9 18"/>
</svg>
```

**Why it works**: Single element. Points are clean whole numbers. The chevron occupies 6px horizontally (9→15) and 12px vertically (6→18) — slightly taller than wide, which is correct for a chevron that needs to feel balanced.

---

## Simple Icon: Close / X

Two diagonal lines crossing. Symmetrical and centred.

```xml
<svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round">
  <line x1="6" y1="6" x2="18" y2="18"/>
  <line x1="18" y1="6" x2="6" y2="18"/>
</svg>
```

**Why it works**: Perfectly symmetrical. Uses `<line>` elements because that's the simplest representation. Coordinates are all whole numbers. The X spans 12px in each direction (6→18), leaving 6px of padding on each side — generous, which gives the icon breathing room.

---

## Medium Icon: Home

A house shape combining a roof (triangle) and body (rectangle). Demonstrates combining multiple elements.

```xml
<svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round">
  <path d="M3 10.5L12 3l9 7.5"/>
  <path d="M5 9.5v10a1 1 0 001 1h12a1 1 0 001-1v-10"/>
  <path d="M9.5 20.5v-6h5v6"/>
</svg>
```

**Why it works**: Three paths — roof line, house body, door. The roof peak (12, 3) is at the top padding boundary. The house bottom (20.5) leaves room for the baseline. The door is centred horizontally. Corner radius on the house body (`a1 1 0 001 1`) matches the preset's cornerRadius of 2px (approximated as 1 in the arc for this scale).

---

## Medium Icon: Email / Mail

An envelope shape. Shows how to handle a recognisable real-world object.

```xml
<svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round">
  <rect x="3" y="5" width="18" height="14" rx="2"/>
  <polyline points="3 5 12 13 21 5"/>
</svg>
```

**Why it works**: Two elements — rectangle body and V-shaped flap line. The `rx="2"` on the rect matches the preset's corner radius. The envelope flap's peak (12, 13) is slightly above centre, which matches how real envelopes look. The rect spans the full usable width (3–21) because envelopes are wide.

---

## Medium Icon: Phone

A phone handset. Demonstrates curved paths and the curved stroke compensation principle.

```xml
<svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round">
  <path d="M22 16.92v3a2 2 0 01-2.18 2 19.79 19.79 0 01-8.63-3.07 19.5 19.5 0 01-6-6 19.79 19.79 0 01-3.07-8.67A2 2 0 014.11 2h3a2 2 0 012 1.72c.13.81.36 1.6.68 2.35a2 2 0 01-.45 2.11L8.09 9.91a16 16 0 006 6l1.27-1.27a2 2 0 012.11-.45c.75.32 1.54.55 2.35.68a2 2 0 011.72 2.03z"/>
</svg>
```

**Why it works**: Single path for the complete handset shape. The curves create a recognisable phone silhouette. Because this icon is almost entirely curved, the path is sized slightly more generously within the grid — it extends close to the edges to compensate for the visual thinning effect of curves.

---

## Complex Icon: Shield Check

A shield with a checkmark inside. Shows layering meaning (protection + verification).

```xml
<svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round">
  <path d="M12 2l7.5 3.5v5c0 5.25-3.19 8.69-7.5 11.5-4.31-2.81-7.5-6.25-7.5-11.5v-5L12 2z"/>
  <polyline points="9 12 11 14 15 10"/>
</svg>
```

**Why it works**: Two elements — shield outline and checkmark. The shield is vertically asymmetric (taller below centre than above) which matches how real shields look. The checkmark is positioned slightly above the shield's visual centre. The shield's top point (12, 2) is at the padding boundary, and its bottom (implicit from the curve, ~22) fills the vertical space.

---

## Complex Icon: Star

A five-pointed star. Demonstrates precise geometric construction.

```xml
<svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round">
  <polygon points="12 2 15.09 8.26 22 9.27 17 14.14 18.18 21.02 12 17.77 5.82 21.02 7 14.14 2 9.27 8.91 8.26"/>
</svg>
```

**Why it works**: Single `<polygon>` element. The star is calculated from proper geometry (inner/outer radius ratios) so it looks regular. The top point (12, 2) is at the padding boundary. The star is wider than it is tall, which is correct for five-pointed stars. Points use max 2 decimal places.

---

## Icon with Fill: Location Pin

Some icons need a filled element alongside stroked elements. Use `fill="currentColor"` on the specific element, not globally.

```xml
<svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round">
  <path d="M12 2C8.13 2 5 5.13 5 9c0 5.25 7 13 7 13s7-7.75 7-13c0-3.87-3.13-7-7-7z"/>
  <circle cx="12" cy="9" r="2.5" fill="currentColor"/>
</svg>
```

**Why it works**: The pin body is stroked (outline). The inner dot uses `fill="currentColor"` to appear solid while still respecting the `currentColor` system. The pin extends from y=2 (top) to an implied y=22 (bottom of the teardrop), using the full vertical space because pins are naturally tall and narrow.

---

## What to Avoid

### Bad: Hardcoded colours
```xml
<!-- WRONG -->
<circle cx="12" cy="12" r="8" stroke="#333" fill="#eee"/>
```

### Bad: Excessive precision
```xml
<!-- WRONG -->
<line x1="5.333333" y1="7.142857" x2="18.666667" y2="16.857143"/>
<!-- RIGHT -->
<line x1="5.5" y1="7" x2="18.5" y2="17"/>
```

### Bad: Transform instead of coordinates
```xml
<!-- WRONG -->
<g transform="translate(2, 3) rotate(45)">
  <rect x="0" y="0" width="10" height="10"/>
</g>
<!-- RIGHT: Bake the transform into coordinates -->
<rect x="7" y="5" width="10" height="10"/>
```

### Bad: Inconsistent padding
```xml
<!-- WRONG: One icon uses full bleed, another has generous padding -->
<!-- Icon A: path starts at x=1 (too close to edge) -->
<!-- Icon B: path starts at x=4 (too much padding) -->
<!-- RIGHT: Both should start at x=2 (or wherever the preset padding specifies) -->
```
