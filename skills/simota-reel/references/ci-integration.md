# Reel CI/CD Integration

Automated terminal recording workflows using VHS in CI pipelines.

---

## GitHub Actions: VHS Recording Workflow

### Using charmbracelet/vhs-action

```yaml
# .github/workflows/record-demos.yml
name: Record Terminal Demos

on:
  push:
    branches: [main]
    paths:
      - 'src/cli/**'
      - '*.tape'
      - '.github/workflows/record-demos.yml'
  pull_request:
    paths:
      - 'src/cli/**'
      - '*.tape'

# ...
```

### Manual VHS Installation

```yaml
      - name: Install VHS (manual)
        run: |
          VHS_VERSION="${{ env.VHS_VERSION || '0.8.0' }}"
          curl -fsSL "https://github.com/charmbracelet/vhs/releases/download/v${VHS_VERSION}/vhs_${VHS_VERSION}_linux_amd64.tar.gz" \
            | tar xz -C /usr/local/bin vhs
          vhs --version

      - name: Install recording dependencies
        run: |
          sudo apt-get update
          sudo apt-get install -y ffmpeg ttyd
```

---

## PR Preview Comments

Post recording previews as PR comments so reviewers can see CLI output changes.

```yaml
  preview:
    runs-on: ubuntu-latest
    if: github.event_name == 'pull_request'
    needs: record
    permissions:
      pull-requests: write
    steps:
      - uses: actions/download-artifact@v4
        with:
          name: terminal-recordings
          path: recordings/

      - name: Upload GIFs to temporary hosting
        id: upload
        run: |
# ...
```

### Before/After Comparison

```yaml
      - name: Download base branch recordings
        run: |
          git fetch origin ${{ github.base_ref }}
          git checkout origin/${{ github.base_ref }} -- docs/recordings/ 2>/dev/null || true
          mv docs/recordings/ base-recordings/ 2>/dev/null || mkdir -p base-recordings/

      - name: Build comparison comment
        id: compare
        run: |
          body="## CLI Output Comparison\n\n"
          for gif in recordings/*.gif; do
            name=$(basename "$gif")
            if [ -f "base-recordings/$name" ]; then
              body+="### ${name}\n| Before | After |\n|--------|-------|\n"
              body+="| ![before](base-recordings/${name}) | ![after](recordings/${name}) |\n\n"
# ...
```

---

## Automated Demo Updates

Auto-commit updated GIFs when CLI source changes on the main branch.

```yaml
  update-docs:
    runs-on: ubuntu-latest
    if: github.ref == 'refs/heads/main'
    needs: record
    permissions:
      contents: write
    steps:
      - uses: actions/checkout@v4
        with:
          token: ${{ secrets.GITHUB_TOKEN }}

      - uses: actions/download-artifact@v4
        with:
          name: terminal-recordings
          path: docs/recordings/
# ...
```

### Git LFS for Large GIFs

```bash
# .gitattributes
docs/recordings/*.gif filter=lfs diff=lfs merge=lfs -text
docs/recordings/*.mp4 filter=lfs diff=lfs merge=lfs -text
```

```yaml
      - name: Setup Git LFS
        run: |
          git lfs install
          git lfs track "docs/recordings/*.gif"
          git lfs track "docs/recordings/*.mp4"
```

---

## Caching Strategies

### VHS Binary Cache

```yaml
      - name: Cache VHS binary
        id: cache-vhs
        uses: actions/cache@v4
        with:
          path: |
            /usr/local/bin/vhs
            ~/.local/share/vhs/
          key: vhs-${{ runner.os }}-v0.8.0

      - name: Install VHS
        if: steps.cache-vhs.outputs.cache-hit != 'true'
        run: |
          curl -fsSL "https://github.com/charmbracelet/vhs/releases/download/v0.8.0/vhs_0.8.0_linux_amd64.tar.gz" \
            | tar xz -C /usr/local/bin vhs
```

### Font Cache

```yaml
      - name: Cache fonts
        uses: actions/cache@v4
        with:
          path: ~/.local/share/fonts
          key: fonts-${{ runner.os }}-${{ hashFiles('fonts.txt') }}

      - name: Install fonts
        if: steps.cache-fonts.outputs.cache-hit != 'true'
        run: |
          mkdir -p ~/.local/share/fonts
          curl -fsSL -o JetBrainsMono.zip \
            "https://github.com/JetBrains/JetBrainsMono/releases/download/v2.304/JetBrainsMono-2.304.zip"
          unzip -o JetBrainsMono.zip -d ~/.local/share/fonts/
          fc-cache -fv
```

### Hash-Based Output Cache

Skip re-recording when neither the tape file nor the CLI binary has changed.

```yaml
      - name: Compute source hash
        id: source-hash
        run: |
          hash=$(cat demos/*.tape src/cli/**/*.ts | sha256sum | cut -d' ' -f1)
          echo "hash=$hash" >> "$GITHUB_OUTPUT"

      - name: Cache recordings
        id: cache-recordings
        uses: actions/cache@v4
        with:
          path: docs/recordings/
          key: recordings-${{ steps.source-hash.outputs.hash }}

      - name: Record demos
        if: steps.cache-recordings.outputs.cache-hit != 'true'
# ...
```

---

## Matrix Recording

Record demos across different shells and themes for comprehensive documentation.

```yaml
  matrix-record:
    runs-on: ubuntu-latest
    strategy:
      matrix:
        shell: [bash, zsh, fish]
        theme: [light, dark]
    steps:
      - uses: actions/checkout@v4

      - name: Install shells
        run: |
          sudo apt-get update
          sudo apt-get install -y ${{ matrix.shell }}

      - name: Install VHS
# ...
```

---

## Complete Example Workflow

A production-ready workflow combining all strategies.

```yaml
# .github/workflows/record-demos.yml
name: Record Terminal Demos

on:
  push:
    branches: [main]
    paths: ['src/cli/**', 'demos/**/*.tape']
  pull_request:
    paths: ['src/cli/**', 'demos/**/*.tape']

env:
  VHS_VERSION: '0.8.0'

jobs:
  record:
# ...
```
