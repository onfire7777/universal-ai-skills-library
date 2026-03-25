# Preview Page Template

Use this template to generate the preview.html file. Replace the placeholder content with actual icons and metadata.

## Template

```html
<!DOCTYPE html>
<html lang="en">
<head>
<meta charset="UTF-8">
<meta name="viewport" content="width=device-width, initial-scale=1.0">
<title>{{PROJECT_NAME}} — Icon Set Preview</title>
<style>
  * { margin: 0; padding: 0; box-sizing: border-box; }
  body {
    font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', system-ui, sans-serif;
    padding: 2rem;
    max-width: 1200px;
    margin: 0 auto;
    color: #1a1a1a;
    background: #fafafa;
  }
  h1 { font-size: 1.5rem; font-weight: 600; margin-bottom: 0.25rem; }
  .subtitle { color: #666; font-size: 0.9rem; margin-bottom: 2rem; }
  .spec {
    background: #fff;
    border: 1px solid #e5e5e5;
    border-radius: 8px;
    padding: 1rem 1.25rem;
    margin-bottom: 2rem;
    font-size: 0.85rem;
    color: #444;
    display: flex;
    flex-wrap: wrap;
    gap: 0.5rem 1.5rem;
  }
  .spec span { white-space: nowrap; }
  .spec strong { color: #1a1a1a; }
  h2 {
    font-size: 1.1rem;
    font-weight: 600;
    margin: 2rem 0 1rem;
    padding-bottom: 0.5rem;
    border-bottom: 1px solid #e5e5e5;
  }
  .grid {
    display: grid;
    grid-template-columns: repeat(auto-fill, minmax(100px, 1fr));
    gap: 1rem;
  }
  .icon-card {
    display: flex;
    flex-direction: column;
    align-items: center;
    gap: 0.5rem;
    padding: 1rem 0.5rem;
    border-radius: 8px;
    background: #fff;
    border: 1px solid #e5e5e5;
    transition: border-color 0.15s;
  }
  .icon-card:hover { border-color: #999; }
  .icon-card svg { flex-shrink: 0; }
  .icon-card .label {
    font-size: 0.7rem;
    color: #666;
    text-align: center;
    word-break: break-all;
  }
  .dark-section {
    background: #1a1a1a;
    border-radius: 12px;
    padding: 2rem;
    margin-top: 2rem;
  }
  .dark-section h2 {
    color: #fff;
    border-bottom-color: #333;
  }
  .dark-section .icon-card {
    background: #2a2a2a;
    border-color: #333;
    color: #fff;
  }
  .dark-section .icon-card:hover { border-color: #555; }
  .dark-section .icon-card .label { color: #999; }
  .size-label {
    font-size: 0.75rem;
    color: #999;
    font-weight: 500;
    text-transform: uppercase;
    letter-spacing: 0.05em;
  }
</style>
</head>
<body>

<h1>{{PROJECT_NAME}} Icons</h1>
<p class="subtitle">{{ICON_COUNT}} icons · {{PRESET}} style · Generated {{DATE}}</p>

<div class="spec">
  <span><strong>Grid:</strong> {{GRID}}px</span>
  <span><strong>Stroke:</strong> {{STROKE_WIDTH}}px</span>
  <span><strong>Caps:</strong> {{STROKE_LINECAP}}</span>
  <span><strong>Joins:</strong> {{STROKE_LINEJOIN}}</span>
  <span><strong>Corner radius:</strong> {{CORNER_RADIUS}}px</span>
  <span><strong>Padding:</strong> {{PADDING}}px</span>
</div>

<h2>Native Size <span class="size-label">({{GRID}}px)</span></h2>
<div class="grid">
  <!-- REPEAT FOR EACH ICON -->
  <div class="icon-card">
    {{SVG_AT_NATIVE_SIZE}}
    <span class="label">{{ICON_NAME}}</span>
  </div>
  <!-- END REPEAT -->
</div>

<h2>2× Size <span class="size-label">({{GRID_2X}}px)</span></h2>
<div class="grid" style="grid-template-columns: repeat(auto-fill, minmax(140px, 1fr));">
  <!-- REPEAT FOR EACH ICON -->
  <div class="icon-card">
    {{SVG_AT_2X_SIZE}}
    <span class="label">{{ICON_NAME}}</span>
  </div>
  <!-- END REPEAT -->
</div>

<div class="dark-section">
  <h2>Dark Background <span class="size-label">({{GRID_2X}}px)</span></h2>
  <div class="grid" style="grid-template-columns: repeat(auto-fill, minmax(140px, 1fr));">
    <!-- REPEAT FOR EACH ICON -->
    <div class="icon-card">
      {{SVG_AT_2X_SIZE}}
      <span class="label">{{ICON_NAME}}</span>
    </div>
    <!-- END REPEAT -->
  </div>
</div>

</body>
</html>
```

## Implementation Notes

- For "native size" SVGs: render at the grid size (e.g. `width="24" height="24"`)
- For "2x size" SVGs: double the width/height attributes (e.g. `width="48" height="48"`) but keep the same viewBox
- The dark section uses CSS `color: #fff` on the card, which `currentColor` in the SVGs will inherit — no SVG changes needed
- Inline every SVG directly in the HTML. Don't use `<img>` tags or external references
- Replace all `{{PLACEHOLDER}}` values with actual data from the style spec
