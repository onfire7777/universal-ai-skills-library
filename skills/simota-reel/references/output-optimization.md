# Output Format Optimization

Reference for choosing and optimizing terminal recording output formats.

## Format Comparison

| Format | Pros | Cons | Best For | Typical Size (30s) |
|--------|------|------|----------|-------------------|
| GIF | Universal support, auto-plays everywhere | Large files, 256 color limit, no audio | README embeds, quick demos | 2-10 MB |
| MP4 (H.264) | Small files, excellent quality, wide support | No auto-play on GitHub, requires player | Documentation sites, marketing | 0.5-2 MB |
| WebM (VP9) | Smaller than MP4, open format, good quality | Limited browser support (no Safari <16) | Web documentation, CI previews | 0.3-1.5 MB |
| SVG | Infinite scalability, searchable text, tiny size | Requires svg-term-cli, no native animation | Technical docs, presentations | 50-200 KB |
| PNG sequence | Lossless quality, individual frame access | Many files, requires assembly step | Post-processing pipelines | 5-20 MB |

## GIF Optimization

### Using gifsicle

**Color reduction** - Reduce palette to lower file size:

```bash
gifsicle --colors 256 input.gif -o output.gif    # max colors (default)
gifsicle --colors 128 input.gif -o output.gif    # good balance
gifsicle --colors 64 input.gif -o output.gif     # aggressive, may show banding
```

**Lossy compression** - Allow minor visual artifacts for significant size savings:

```bash
gifsicle --lossy=80 input.gif -o output.gif      # moderate compression
gifsicle --lossy=120 input.gif -o output.gif     # aggressive compression
gifsicle --lossy=200 input.gif -o output.gif     # maximum (visible artifacts)
```

**Frame optimization** - Remove redundant pixel data between frames:

```bash
gifsicle -O1 input.gif -o output.gif             # store changed portions
gifsicle -O2 input.gif -o output.gif             # also use transparency
gifsicle -O3 input.gif -o output.gif             # try multiple methods (slowest)
```

**Resize** - Scale down to reduce dimensions and file size:

```bash
gifsicle --resize-width 800 input.gif -o output.gif
gifsicle --resize-width 600 input.gif -o output.gif   # compact for docs
gifsicle --resize 800x600 input.gif -o output.gif     # exact dimensions
```

**Combined optimization** - Apply multiple techniques together:

```bash
# Balanced (recommended starting point)
gifsicle -O3 --lossy=80 --colors 128 input.gif -o output.gif

# Aggressive (when file size is critical)
gifsicle -O3 --lossy=120 --colors 64 --resize-width 600 input.gif -o output.gif

# Quality-first (when visual fidelity matters)
gifsicle -O3 --lossy=30 --colors 256 input.gif -o output.gif
```

### Target Sizes by Use Case

| Use Case | Target Size | Recommended Settings |
|----------|-------------|---------------------|
| README embed | < 5 MB | `--lossy=80 --colors 128 -O3` |
| Documentation | < 2 MB | `--lossy=80 --colors 64 --resize-width 800 -O3` |
| Social/marketing | < 8 MB | `--lossy=30 --colors 256 -O3` |

## MP4 Optimization

### Using ffmpeg with H.264

**Standard encoding** - CRF controls quality (lower = better, 18-28 typical):

```bash
ffmpeg -i input.mp4 -c:v libx264 -crf 23 output.mp4          # default quality
ffmpeg -i input.mp4 -c:v libx264 -crf 28 output.mp4          # smaller file
ffmpeg -i input.mp4 -c:v libx264 -crf 18 output.mp4          # higher quality
```

**Web-optimized** - Move metadata to front for faster streaming start:

```bash
ffmpeg -i input.mp4 -c:v libx264 -crf 23 -movflags +faststart output.mp4
```

**Resolution and bitrate control**:

```bash
# Scale to 720p width, maintain aspect ratio
ffmpeg -i input.mp4 -vf "scale=1280:-2" -c:v libx264 -crf 23 output.mp4

# Limit bitrate for predictable file sizes
ffmpeg -i input.mp4 -c:v libx264 -b:v 1M -maxrate 1.5M -bufsize 2M output.mp4
```

**Recommended preset for terminal recordings**:

```bash
ffmpeg -i input.mp4 \
  -c:v libx264 -crf 23 -preset slow \
  -vf "scale=1280:-2" \
  -movflags +faststart \
  -an output.mp4
```

> Note: `-an` removes audio (terminal recordings rarely need it).

## WebM Optimization

### Using ffmpeg with VP9

**Single-pass encoding** - Simple and fast:

```bash
ffmpeg -i input.mp4 -c:v libvpx-vp9 -crf 30 -b:v 0 output.webm
```

**Two-pass encoding** - Better quality/size ratio for final output:

```bash
# Pass 1: Analysis
ffmpeg -i input.mp4 -c:v libvpx-vp9 -b:v 1M -pass 1 -an -f null /dev/null

# Pass 2: Encode
ffmpeg -i input.mp4 -c:v libvpx-vp9 -b:v 1M -pass 2 -an output.webm
```

**Recommended preset for terminal recordings**:

```bash
ffmpeg -i input.mp4 \
  -c:v libvpx-vp9 -crf 30 -b:v 0 \
  -vf "scale=1280:-2" \
  -an output.webm
```

## SVG Output

### From asciinema with svg-term-cli

**Record and convert**:

```bash
# Record with asciinema
asciinema rec recording.cast

# Convert to SVG
svg-term --in recording.cast --out output.svg

# With options
svg-term --in recording.cast --out output.svg \
  --window \
  --width 80 \
  --height 24 \
  --padding 10 \
  --from 1000 \
  --to 15000
```

**Advantages of SVG output**:

- Infinite scalability without quality loss
- Text remains searchable and selectable
- Extremely small file sizes (50-200 KB for 30s)
- Embeddable directly in HTML/Markdown
- Customizable via CSS (colors, fonts, sizing)

**Limitations**:

- Requires svg-term-cli (`npm install -g svg-term-cli`)
- Only works with asciinema `.cast` recordings
- Animation timing may differ slightly from original
- Some complex terminal output may not render perfectly

## File Size Guidelines

| Use Case | Max Size | Recommended Format | Notes |
|----------|----------|-------------------|-------|
| GitHub README | 5 MB | GIF or SVG | GIF for broad support, SVG for quality |
| Documentation site | 2 MB | WebM with GIF fallback | WebM preferred, GIF as fallback |
| Marketing/social | 8 MB | MP4 or GIF | MP4 for quality, GIF for compatibility |
| CI preview | 1 MB | SVG or compressed GIF | Keep small for fast loading |
| Presentation | 10 MB | MP4 | Quality matters, size less critical |
| Blog post | 3 MB | GIF or WebM | Depends on platform support |

## Quality vs Size Decision Matrix

```
Need universal compatibility?
  YES --> Use GIF
    File too large?
      YES --> Apply gifsicle optimization
        Still too large?
          YES --> Reduce resolution, increase lossy
          NO  --> Done
      NO  --> Done
  NO  --> Need scalable/searchable output?
    YES --> Use SVG (via asciinema + svg-term)
    NO  --> Need smallest file size?
      YES --> Use WebM (VP9)
      NO  --> Use MP4 (H.264, widest video support)
```

### Quick Recommendations

| Priority | Format | Why |
|----------|--------|-----|
| Compatibility first | GIF | Works everywhere, no player needed |
| Quality first | MP4 (H.264) | Best quality-to-size ratio with wide support |
| Size first | WebM (VP9) or SVG | Smallest files, SVG is searchable |
| Accessibility first | SVG | Searchable text, scales to any size |
