/**
 * Minimal, dependency-free YAML front-matter reader for SKILL.md files.
 *
 * The skill router (skill-router-cli/cmd/skills) treats a skill's identity as:
 *   - canonical id  = the kebab-case directory name (NOT the front-matter `name`)
 *   - description   = the front-matter `description`, whitespace-collapsed to a
 *                     single line (block scalars are folded with single spaces)
 *
 * This reader reproduces that exact normalization so the generated manifest is
 * byte-identical to the hand-authored one for the existing skills whose
 * descriptions are pure front-matter. The remaining curated descriptions are
 * supplied as overrides by the generator (see registry.config.json).
 *
 * Only `name` and `description` are needed by the registry, so the parser is
 * intentionally limited to top-level scalar / block-scalar string fields and
 * does not attempt to be a general YAML implementation.
 */

/** Collapse all runs of whitespace to a single space and trim the ends. */
export function collapseWhitespace(value) {
  return value.replace(/\s+/g, " ").trim();
}

/** Strip a single matching pair of surrounding single or double quotes. */
function stripQuotes(value) {
  const m = value.match(/^(['"])([\s\S]*)\1$/);
  return m ? m[2] : value;
}

/**
 * Extract the leading `---\n ... \n---` front-matter block. Returns the inner
 * body (without the fences) or `null` when no front-matter is present.
 */
export function extractFrontMatterBlock(content) {
  // Tolerate a UTF-8 BOM and both LF / CRLF line endings.
  const text = content.replace(/^﻿/, "").replace(/\r\n/g, "\n");
  const m = text.match(/^---\n([\s\S]*?)\n---(?:\n|$)/);
  return m ? m[1] : null;
}

/**
 * Read a single top-level scalar field (`key:`) from a front-matter body,
 * supporting inline values, quoted values, and `|` / `>` block scalars.
 * Returns the whitespace-collapsed string, or `undefined` if the key is absent.
 */
export function readScalarField(body, key) {
  const lines = body.split("\n");
  for (let i = 0; i < lines.length; i++) {
    const line = lines[i];
    // Top-level keys start at column 0; indented lines belong to a parent key.
    if (/^\s/.test(line)) continue;
    const m = line.match(new RegExp("^" + key + ":(.*)$"));
    if (!m) continue;

    const rest = m[1].replace(/^\s/, "");
    // Block scalar indicator: `|`, `>`, with optional chomping/indent markers.
    if (/^[|>][+-]?\d*\s*$/.test(rest)) {
      const collected = [];
      for (let j = i + 1; j < lines.length; j++) {
        const l = lines[j];
        if (l.trim() === "") {
          collected.push("");
          continue;
        }
        // A non-indented line ends the block scalar.
        if (!/^\s/.test(l)) break;
        collected.push(l.replace(/^\s+/, ""));
      }
      return collapseWhitespace(collected.join(" "));
    }

    return collapseWhitespace(stripQuotes(rest.trim()));
  }
  return undefined;
}

/**
 * Parse a SKILL.md file's front matter.
 * @param {string} content raw file contents
 * @returns {{ name?: string, description?: string }}
 */
export function parseFrontMatter(content) {
  const body = extractFrontMatterBlock(content);
  if (body === null) return {};
  return {
    name: readScalarField(body, "name"),
    description: readScalarField(body, "description"),
  };
}
