#!/usr/bin/env node
/**
 * generate-registry.mjs — the single emitter for every registry artifact.
 *
 * ONE source of truth → MANY generated artifacts (kept in lockstep, so they can
 * never drift the way the old hand-maintained pair did):
 *
 *   source:   skills/ on disk  +  scripts/registry/registry.config.json
 *   emits:    manifest.json                 (router catalog — read by skill-router)
 *             marketplace.json              (canonical Claude plugin marketplace)
 *             .agents/plugins/marketplace.json (codex variant; shared metadata in lockstep)
 *             docs/build_manifest.json      (provenance / build report)
 *
 * The skill CATALOG is derived from skills/ (every top-level dir containing a
 * SKILL.md). Per-skill descriptions come from SKILL.md front matter, except for
 * a small curated set of router-tuned overrides in the config. has_scripts /
 * scripts are computed with the exact same rules as the Go validator
 * (skill-router-cli validate-manifest → listSkillScripts), so generated output
 * always passes `skill-router validate-manifest`.
 *
 * Modes:
 *   --check            generate in memory, diff against the committed files,
 *                      exit 1 on drift. No writes. (Use in CI.)
 *   --write            write the generated artifacts to disk.
 *   --print <artifact> print one artifact to stdout (manifest|marketplace|codex-marketplace|build-manifest).
 *   --optimize         apply the de-bloat transforms (see below) instead of the
 *                      byte-faithful reproduction of the legacy files.
 *   --only <list>      restrict to a comma list of artifacts (default: all).
 *
 * Faithful (default) vs --optimize:
 *   faithful   reproduces the legacy manifest.json/marketplace.json byte-for-byte
 *              (the refactor-only proof). build_manifest is not reproduced in this
 *              mode because it is intentionally restructured by --optimize.
 *   optimize   - manifest entries omit empty optional fields (has_scripts:false,
 *                empty scripts/aliases) relying on the Go reader's `omitempty`
 *              - build_manifest slimmed to provenance only (drops the duplicated
 *                ~750KB catalog; nothing consumes its skills[])
 *              - build_manifest paths made portable (relative, not %USERPROFILE%)
 *              - counts recomputed from the live tree (kills the 1812/1811 drift)
 *              - marketplace gains the 14 themed groupings + live skill count
 *
 * Invariants preserved in BOTH modes (breaking these is CHANGES_REQUESTED):
 *   - manifest.routing.legacy_access = "manus skill <name>"   (the `manus` alias)
 *   - build_manifest.legacy_binary_alias = "manus"
 *   - merged_legacy_directories / disabled_colliding_aliases / compatibility_policy
 */
import fs from "node:fs";
import path from "node:path";
import { fileURLToPath } from "node:url";
import { parseFrontMatter, collapseWhitespace } from "./lib/frontmatter.mjs";

const HERE = path.dirname(fileURLToPath(import.meta.url));
export const REPO_ROOT = path.resolve(HERE, "..", "..");
export const CONFIG_PATH = path.join(HERE, "registry.config.json");

export const ARTIFACTS = {
  manifest: "manifest.json",
  marketplace: "marketplace.json",
  "codex-marketplace": ".agents/plugins/marketplace.json",
  "build-manifest": "docs/build_manifest.json",
};

// ---------------------------------------------------------------------------
// serialization — matches the legacy files: 2-space indent, trailing newline,
// no HTML escaping of <>& (Node's JSON.stringify already leaves them raw).
// ---------------------------------------------------------------------------
export function serialize(obj) {
  return JSON.stringify(obj, null, 2) + "\n";
}

// ---------------------------------------------------------------------------
// scripts discovery — byte-compatible port of validate_manifest.go listSkillScripts
// ---------------------------------------------------------------------------
export function listSkillScripts(skillDir) {
  const scriptsDir = path.join(skillDir, "scripts");
  if (!fs.existsSync(scriptsDir)) return [];
  const out = [];
  const walk = (dir) => {
    for (const entry of fs.readdirSync(dir, { withFileTypes: true })) {
      if (entry.isDirectory()) {
        if (entry.name === "__pycache__" || entry.name === ".git") continue;
        walk(path.join(dir, entry.name));
      } else if (entry.isFile()) {
        if (entry.name.endsWith(".pyc") || entry.name.endsWith(".pyo")) continue;
        const rel = path.relative(skillDir, path.join(dir, entry.name));
        out.push(rel.split(path.sep).join("/"));
      }
    }
  };
  walk(scriptsDir);
  out.sort();
  return out;
}

// ---------------------------------------------------------------------------
// catalog scan — every skills/<name>/ that has a SKILL.md, derived + curated
// ---------------------------------------------------------------------------
export function scanSkills(config) {
  const skillsRoot = path.join(REPO_ROOT, "skills");
  const overrides = config.descriptionOverrides || {};
  const aliases = config.aliases || {};
  const skills = [];
  for (const entry of fs.readdirSync(skillsRoot, { withFileTypes: true })) {
    if (!entry.isDirectory()) continue;
    const name = entry.name;
    const skillDir = path.join(skillsRoot, name);
    const skillMd = path.join(skillDir, "SKILL.md");
    if (!fs.existsSync(skillMd)) continue;

    let description = overrides[name];
    if (description === undefined) {
      const parsed = parseFrontMatter(fs.readFileSync(skillMd, "utf8"));
      description = parsed.description ? collapseWhitespace(parsed.description) : "";
    }
    const scripts = listSkillScripts(skillDir);
    skills.push({
      name,
      directory: `skills/${name}`,
      description,
      aliases: aliases[name],
      has_scripts: scripts.length > 0,
      scripts,
    });
  }
  return skills;
}

/** Order an entry's keys exactly as the legacy manifest did, applying omitempty in optimize mode. */
function manifestEntry(skill, { optimize }) {
  const e = { name: skill.name, directory: skill.directory, description: skill.description };
  if (skill.aliases && skill.aliases.length > 0) e.aliases = skill.aliases;
  if (optimize) {
    // rely on the Go reader's omitempty: only emit when meaningful
    if (skill.has_scripts) e.has_scripts = true;
    if (skill.scripts.length > 0) e.scripts = skill.scripts;
  } else {
    e.has_scripts = skill.has_scripts;
    e.scripts = skill.scripts;
  }
  return e;
}

// ---------------------------------------------------------------------------
// artifact builders
// ---------------------------------------------------------------------------
export function buildManifest(config, skills, { optimize }) {
  const coreSet = new Set(config.coreSkills);
  const byName = new Map(skills.map((s) => [s.name, s]));
  // core in curated order; library alphabetical (matches legacy code-unit sort)
  const core = config.coreSkills.filter((n) => byName.has(n)).map((n) => byName.get(n));
  const library = skills
    .filter((s) => !coreSet.has(s.name))
    .sort((a, b) => (a.name < b.name ? -1 : a.name > b.name ? 1 : 0));
  const total = core.length + library.length;
  const aliasCount = optimize
    ? skills.reduce((n, s) => n + (s.aliases ? s.aliases.length : 0), 0)
    : config.manifest.alias_count;
  return {
    version: config.manifest.version,
    generated: config.manifest.generated,
    description: config.manifest.description,
    canonical_id_policy: config.manifest.canonical_id_policy,
    core_skills: core.map((s) => manifestEntry(s, { optimize })),
    library_skills: library.map((s) => manifestEntry(s, { optimize })),
    total_skills: total,
    alias_count: aliasCount,
    routing: config.manifest.routing, // preserves legacy_access = "manus skill <name>"
  };
}

export function buildMarketplace(config, skills, { optimize }) {
  // Deep clone so we never mutate the loaded config.
  const market = JSON.parse(JSON.stringify(config.marketplace));
  if (optimize) {
    const count = skills.length.toLocaleString("en-US");
    for (const plugin of market.plugins || []) {
      if (typeof plugin.description === "string") {
        plugin.description = plugin.description.replace(/[\d,]+ skills/g, `${count} skills`);
      }
    }
    if (Array.isArray(config.groupings) && config.groupings.length > 0) {
      // Preserve the manus marketplace's only unique value: themed collections.
      market.groupings = config.groupings.map((g) => ({
        name: g.name,
        description: g.description,
        members: g.members,
      }));
    }
  }
  return market;
}

export function buildBuildManifest(config, skills, { optimize }) {
  const b = config.buildManifest;
  const base = {
    schema: b.schema,
    generated_at: b.generated_at,
    source_of_truth: optimize ? b.source_of_truth : b.legacy_source_of_truth,
    router_source: optimize ? b.router_source : b.legacy_router_source,
    primary_binary: b.primary_binary,
    legacy_binary_alias: b.legacy_binary_alias, // "manus"
    skill_count: skills.length,
    directories_total: skills.length,
    missing_skill_md: [],
    alias_count: optimize
      ? skills.reduce((n, s) => n + (s.aliases ? s.aliases.length : 0), 0)
      : b.alias_count,
    merged_legacy_directories: b.merged_legacy_directories,
    disabled_colliding_aliases: b.disabled_colliding_aliases,
    compatibility_policy: b.compatibility_policy,
  };
  if (!optimize) {
    // Faithful provenance form keeps the (legacy) full catalog mirror.
    base.skills = skills
      .slice()
      .sort((a, b2) => (a.name < b2.name ? -1 : a.name > b2.name ? 1 : 0))
      .map((s) => ({
        name: s.name,
        directory: s.directory,
        description: s.description,
        has_scripts: s.has_scripts,
        scripts: s.scripts,
      }));
  } else {
    // Slimmed: the catalog lives in manifest.json; build_manifest is provenance
    // only. Point consumers at the canonical catalog rather than duplicating it.
    base.catalog_ref = "../manifest.json";
  }
  return base;
}

export function buildCodexMarketplace(config) {
  // Lockstep codex variant: shares identity with the canonical marketplace but
  // points at ./plugin-codex with the agents-runtime schema.
  const canonical = config.marketplace;
  return {
    name: canonical.name,
    interface: { displayName: canonical.plugins?.[0] ? "Universal AI Skills" : canonical.name },
    plugins: [
      {
        name: canonical.name,
        source: { source: "local", path: "./plugin-codex" },
        policy: { installation: "INSTALLED_BY_DEFAULT", authentication: "ON_INSTALL" },
        category: "Productivity",
      },
    ],
  };
}

export function buildAll(config, skills, opts) {
  return {
    manifest: buildManifest(config, skills, opts),
    marketplace: buildMarketplace(config, skills, opts),
    "codex-marketplace": buildCodexMarketplace(config, skills, opts),
    "build-manifest": buildBuildManifest(config, skills, opts),
  };
}

// ---------------------------------------------------------------------------
// drift comparison — compare by MEANING, not bytes:
//   - drop volatile timestamp fields (generated / generated_at)
//   - sort each entry's scripts[] (the Go validator treats scripts as an
//     unordered set: validate_manifest.go normalizeScriptList sorts both sides).
// This makes --check robust against the legacy manifest's non-canonical script
// ordering while still catching real drift (added/removed skills, changed
// descriptions, changed script SETS, etc.).
// ---------------------------------------------------------------------------
function sortEntryScripts(entry) {
  if (Array.isArray(entry.scripts)) entry.scripts = entry.scripts.slice().sort();
  return entry;
}

export function normalizeForCompare(key, text) {
  if (key !== "manifest" && key !== "build-manifest") return text;
  try {
    const obj = JSON.parse(text);
    if (key === "manifest") {
      delete obj.generated;
      (obj.core_skills || []).forEach(sortEntryScripts);
      (obj.library_skills || []).forEach(sortEntryScripts);
    }
    if (key === "build-manifest") {
      delete obj.generated_at;
      (obj.skills || []).forEach(sortEntryScripts);
    }
    return serialize(obj);
  } catch {
    return text;
  }
}

// ---------------------------------------------------------------------------
// main
// ---------------------------------------------------------------------------
function parseArgs(argv) {
  const args = { mode: "check", optimize: false, only: null, printArtifact: null };
  for (let i = 0; i < argv.length; i++) {
    const a = argv[i];
    if (a === "--write") args.mode = "write";
    else if (a === "--check") args.mode = "check";
    else if (a === "--optimize") args.optimize = true;
    else if (a === "--print") {
      args.mode = "print";
      args.printArtifact = argv[++i];
    } else if (a === "--only") args.only = argv[++i].split(",").map((s) => s.trim());
    else {
      console.error(`unknown argument: ${a}`);
      process.exit(2);
    }
  }
  return args;
}

function main() {
  const args = parseArgs(process.argv.slice(2));
  const config = JSON.parse(fs.readFileSync(CONFIG_PATH, "utf8"));
  const skills = scanSkills(config);
  const built = buildAll(config, skills, { optimize: args.optimize });

  // Phase-1 safe default: faithful --check guards the artifacts that are byte-
  // consistent with the committed tree today (manifest, marketplace). The other
  // two (build-manifest, codex-marketplace) are intentionally restructured by
  // --optimize, so they enter the check once optimize/--only is requested.
  const DEFAULT_CHECK = ["manifest", "marketplace"];
  let selected;
  if (args.only) selected = args.only.filter((k) => k in ARTIFACTS);
  else if (args.mode === "check" && !args.optimize) selected = DEFAULT_CHECK;
  else selected = Object.keys(ARTIFACTS);

  if (args.mode === "print") {
    const key = args.printArtifact;
    if (!(key in built)) {
      console.error(`unknown artifact: ${key}`);
      process.exit(2);
    }
    process.stdout.write(serialize(built[key]));
    return;
  }

  if (args.mode === "write") {
    for (const key of selected) {
      const target = path.join(REPO_ROOT, ARTIFACTS[key]);
      fs.mkdirSync(path.dirname(target), { recursive: true });
      fs.writeFileSync(target, serialize(built[key]));
      console.log(`wrote ${ARTIFACTS[key]}`);
    }
    return;
  }

  // mode === "check"
  let drift = 0;
  for (const key of selected) {
    const target = path.join(REPO_ROOT, ARTIFACTS[key]);
    const generated = serialize(built[key]);
    if (!fs.existsSync(target)) {
      console.error(`DRIFT: ${ARTIFACTS[key]} is missing (would be generated)`);
      drift++;
      continue;
    }
    const committed = fs.readFileSync(target, "utf8");
    if (normalizeForCompare(key, committed) !== normalizeForCompare(key, generated)) {
      console.error(`DRIFT: ${ARTIFACTS[key]} differs from generated output`);
      drift++;
    } else {
      console.log(`ok: ${ARTIFACTS[key]} in sync`);
    }
  }
  if (drift > 0) {
    console.error(`\n${drift} registry artifact(s) drifted. Run: node scripts/registry/generate-registry.mjs --write`);
    process.exit(1);
  }
  console.log(`\nall ${selected.length} registry artifact(s) in sync (${skills.length} skills)`);
}

import { pathToFileURL } from "node:url";
if (process.argv[1] && import.meta.url === pathToFileURL(process.argv[1]).href) {
  main();
}
