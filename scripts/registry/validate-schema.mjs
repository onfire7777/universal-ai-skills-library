#!/usr/bin/env node
/**
 * validate-schema.mjs — gate SKILL.md frontmatter against schemas/skill.schema.json.
 *
 * Phase 0 of ARCHITECTURE_IMPROVEMENT_PLAN.md (§3.2, §9.1). This is the
 * `validate-schema` step the plan calls for, implemented in the project's
 * existing zero-dependency Node registry toolchain (no new Go, no npm deps) so
 * it adds ZERO runtime-routing surface.
 *
 *   node scripts/registry/validate-schema.mjs            # --warn (default): report, exit 0
 *   node scripts/registry/validate-schema.mjs --error     # strict: exit 1 on any violation
 *   node scripts/registry/validate-schema.mjs --json       # machine-readable report
 *   node scripts/registry/validate-schema.mjs --quiet       # only the summary line
 *
 * Design (mirrors lib/frontmatter.mjs's "intentionally limited" philosophy):
 *   - Validation is driven BY the schema JSON (a tiny interpreter for exactly the
 *     keywords the schema uses), so editing schemas/skill.schema.json changes
 *     behavior — the validator cannot drift from the contract.
 *   - Only `name` + `description` are required; every other field is optional and
 *     backfilled incrementally, so --warn is the default and a missing optional
 *     field is NOT a violation — it is tracked in the coverage report instead.
 *   - --warn always exits 0 for schema violations (operational failures — missing
 *     schema, no skills found — still exit non-zero). --error exits 1 if any
 *     skill has a violation. This is the plan's "--warn then --error" ramp.
 */

import { readFileSync, readdirSync, statSync } from "node:fs";
import { fileURLToPath } from "node:url";
import { dirname, join, resolve } from "node:path";
import { extractFrontMatterBlock, collapseWhitespace } from "./lib/frontmatter.mjs";

const HERE = dirname(fileURLToPath(import.meta.url));
const REPO_ROOT = resolve(HERE, "..", "..");
const SCHEMA_PATH = join(REPO_ROOT, "schemas", "skill.schema.json");
const SKILLS_DIR = join(REPO_ROOT, "skills");

// --- args -----------------------------------------------------------------
const args = process.argv.slice(2);
const MODE = args.includes("--error") ? "error" : "warn";
const JSON_OUT = args.includes("--json");
const QUIET = args.includes("--quiet");

// --- minimal YAML-subset frontmatter reader -------------------------------
// Handles exactly the shapes that appear in the corpus and the schema:
// top-level scalars (quoted / block `|`/`>`), inline `[a, b]` / `{k: v}`,
// block sequences (`- item`), and one level of nested maps (`requires:`).
function stripQuotes(v) {
  const m = v.match(/^(['"])([\s\S]*)\1$/);
  return m ? m[2] : v;
}
function parseInlineArray(v) {
  // [a, "b", c]  -> ["a","b","c"]   (scalars only; good enough for our fields)
  const inner = v.slice(1, -1).trim();
  if (!inner) return [];
  return inner.split(",").map((s) => stripQuotes(s.trim())).filter((s) => s.length);
}
function parseScalar(v) {
  const t = stripQuotes(v.trim());
  return collapseWhitespace(t);
}

/** Parse a frontmatter body into a plain object (limited YAML). */
function parseFrontMatter(content) {
  const body = extractFrontMatterBlock(content);
  if (body === null) return { __noFrontMatter: true };
  const lines = body.replace(/\r\n/g, "\n").split("\n");
  const obj = {};
  for (let i = 0; i < lines.length; i++) {
    const line = lines[i];
    if (/^\s/.test(line) || line.trim() === "" || line.trimStart().startsWith("#")) continue;
    const m = line.match(/^([A-Za-z][A-Za-z0-9_-]*):(.*)$/);
    if (!m) continue;
    const key = m[1];
    const rest = m[2].replace(/^[ \t]/, "");

    // Block scalar (| or >)
    if (/^[|>][+-]?\d*\s*$/.test(rest)) {
      const collected = [];
      for (let j = i + 1; j < lines.length; j++) {
        const l = lines[j];
        if (l.trim() === "") { collected.push(""); continue; }
        if (!/^\s/.test(l)) break;
        collected.push(l.replace(/^\s+/, ""));
        i = j;
      }
      obj[key] = collapseWhitespace(collected.join(" "));
      continue;
    }
    // Inline array / object
    if (rest.startsWith("[")) { obj[key] = parseInlineArray(rest); continue; }
    if (rest.startsWith("{")) { obj[key] = parseInlineObjectIsh(rest); continue; }

    // Empty value: peek the indented block — sequence or nested map.
    if (rest.trim() === "") {
      const seq = [];
      const map = {};
      let kind = null;
      for (let j = i + 1; j < lines.length; j++) {
        const l = lines[j];
        if (l.trim() === "") { i = j; continue; }
        if (!/^\s/.test(l)) break;
        const item = l.trim();
        const sm = item.match(/^-\s*(.*)$/);
        if (sm) { kind = kind || "seq"; if (kind === "seq") seq.push(stripQuotes(sm[1].trim())); i = j; continue; }
        const km = item.match(/^([A-Za-z][A-Za-z0-9_-]*):(.*)$/);
        if (km) {
          kind = kind || "map";
          if (kind === "map") {
            const sub = km[2].trim();
            map[km[1]] = sub.startsWith("[") ? parseInlineArray(sub) : (sub === "" ? [] : parseScalar(sub));
          }
          i = j; continue;
        }
        break;
      }
      obj[key] = kind === "map" ? map : seq;
      continue;
    }
    // Plain scalar
    obj[key] = parseScalar(rest);
  }
  return obj;
}
function parseInlineObjectIsh(v) {
  // Best-effort {mcp: [], tools: []} — we only need to recognize it as an object.
  const out = {};
  const inner = v.slice(1, v.lastIndexOf("}") >= 0 ? v.lastIndexOf("}") : v.length).replace(/^{/, "");
  const m = inner.matchAll(/([A-Za-z][A-Za-z0-9_-]*):\s*(\[[^\]]*\]|[^,}]*)/g);
  for (const mm of m) {
    const val = mm[2].trim();
    out[mm[1]] = val.startsWith("[") ? parseInlineArray(val) : parseScalar(val);
  }
  return out;
}

// --- minimal JSON-Schema interpreter --------------------------------------
// Supports only the keywords schemas/skill.schema.json actually uses.
function resolveRef(ref, root) {
  if (!ref.startsWith("#/")) throw new Error(`unsupported $ref: ${ref}`);
  return ref.slice(2).split("/").reduce((acc, seg) => acc[seg], root);
}
function typeOf(v) {
  if (Array.isArray(v)) return "array";
  if (v === null) return "null";
  return typeof v === "object" ? "object" : typeof v;
}
function validate(value, schema, root, path, out) {
  if (schema.$ref) schema = { ...resolveRef(schema.$ref, root), ...schema };
  if (schema.type) {
    const t = typeOf(value);
    const ok = schema.type === "integer" ? t === "number" && Number.isInteger(value) : t === schema.type;
    if (!ok) { out.push({ path, msg: `expected ${schema.type}, got ${t}` }); return; }
  }
  if (typeof value === "string") {
    if (schema.minLength != null && value.length < schema.minLength) out.push({ path, msg: `shorter than minLength ${schema.minLength}` });
    if (schema.maxLength != null && value.length > schema.maxLength) out.push({ path, msg: `longer than maxLength ${schema.maxLength}` });
    if (schema.pattern && !new RegExp(schema.pattern).test(value)) out.push({ path, msg: `does not match pattern ${schema.pattern}` });
    if (schema.enum && !schema.enum.includes(value)) out.push({ path, msg: `"${value}" not in allowed set (${schema.enum.join(", ")})` });
  }
  if (Array.isArray(value)) {
    if (schema.uniqueItems && new Set(value.map((x) => JSON.stringify(x))).size !== value.length) out.push({ path, msg: "items not unique" });
    if (schema.items) value.forEach((el, idx) => validate(el, schema.items, root, `${path}[${idx}]`, out));
  }
  if (value && typeof value === "object" && !Array.isArray(value)) {
    for (const req of schema.required || []) if (!(req in value)) out.push({ path: `${path}.${req}`, msg: "required field missing", missingRequired: true });
    for (const [k, v] of Object.entries(value)) {
      const ps = schema.properties && schema.properties[k];
      if (ps) validate(v, ps, root, `${path}.${k}`, out);
      else if (schema.additionalProperties === false) out.push({ path: `${path}.${k}`, msg: "additional property not allowed" });
    }
  }
}

// --- run -------------------------------------------------------------------
function fail(msg) {
  if (JSON_OUT) process.stdout.write(JSON.stringify({ ok: false, error: msg }) + "\n");
  else process.stderr.write(`validate-schema: ${msg}\n`);
  process.exit(2);
}

let schema;
try { schema = JSON.parse(readFileSync(SCHEMA_PATH, "utf8")); }
catch (e) { fail(`cannot load schema ${SCHEMA_PATH}: ${e.message}`); }

let dirs;
try { dirs = readdirSync(SKILLS_DIR).filter((d) => statSync(join(SKILLS_DIR, d)).isDirectory()).sort(); }
catch (e) { fail(`cannot read skills dir ${SKILLS_DIR}: ${e.message}`); }

const trackedProps = Object.keys(schema.properties);
const coverage = Object.fromEntries(trackedProps.map((p) => [p, 0]));
const offenders = [];
let total = 0, withViolations = 0, violationCount = 0, noFrontMatter = 0, nameMismatch = 0;

for (const dir of dirs) {
  const file = join(SKILLS_DIR, dir, "SKILL.md");
  let raw;
  try { raw = readFileSync(file, "utf8"); } catch { continue; }
  total++;
  const fm = parseFrontMatter(raw);
  if (fm.__noFrontMatter) {
    noFrontMatter++;
    const v = [{ path: dir, msg: "no YAML frontmatter block", missingRequired: true }];
    offenders.push({ skill: dir, violations: v });
    withViolations++; violationCount += v.length;
    continue;
  }
  for (const p of trackedProps) if (p in fm) coverage[p]++;

  const out = [];
  validate(fm, schema, schema, dir, out);
  // Cross-file invariant the JSON Schema cannot express: name MUST equal dir id.
  if (typeof fm.name === "string" && fm.name !== dir) {
    out.push({ path: `${dir}.name`, msg: `name "${fm.name}" != directory id "${dir}"` });
    nameMismatch++;
  }
  if (out.length) {
    offenders.push({ skill: dir, violations: out });
    withViolations++; violationCount += out.length;
  }
}

const report = {
  ok: violationCount === 0,
  mode: MODE,
  schema: "schemas/skill.schema.json",
  totals: { skills: total, with_violations: withViolations, violations: violationCount, no_frontmatter: noFrontMatter, name_dir_mismatch: nameMismatch },
  coverage: Object.fromEntries(trackedProps.map((p) => [p, { present: coverage[p], pct: total ? +(100 * coverage[p] / total).toFixed(1) : 0 }])),
  offenders: offenders.slice(0, 50),
  offenders_truncated: Math.max(0, offenders.length - 50),
};

if (JSON_OUT) {
  process.stdout.write(JSON.stringify(report, null, 2) + "\n");
} else {
  const C = (s) => s; // keep dependency-free, no color
  console.log(`\nSKILL.md schema validation (${MODE} mode) — ${total} skills`);
  console.log("─".repeat(64));
  console.log("Field coverage (backfill tracker):");
  for (const p of trackedProps) {
    const c = report.coverage[p];
    const req = (schema.required || []).includes(p) ? " [required]" : "";
    console.log(`  ${p.padEnd(14)} ${String(c.present).padStart(5)}/${total}  ${String(c.pct).padStart(5)}%${req}`);
  }
  console.log("─".repeat(64));
  if (!QUIET && offenders.length) {
    console.log(`Violations in ${withViolations} skill(s):`);
    for (const o of offenders.slice(0, 50)) {
      for (const v of o.violations) console.log(`  ${v.missingRequired ? "REQUIRED" : "invalid "}  ${v.path}: ${v.msg}`);
    }
    if (report.offenders_truncated) console.log(`  … and ${report.offenders_truncated} more skill(s) with violations`);
    console.log("─".repeat(64));
  }
  console.log(`Result: ${violationCount} violation(s) across ${withViolations}/${total} skills`
    + (MODE === "warn" ? "  (warn mode — not failing the build)" : ""));
}

// Exit policy: --warn never fails on schema violations; --error does.
process.exit(MODE === "error" && violationCount > 0 ? 1 : 0);
