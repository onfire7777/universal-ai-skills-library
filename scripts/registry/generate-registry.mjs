#!/usr/bin/env node
/**
 * generate-registry.mjs — the (now SECONDARY) registry emitter / parity oracle.
 *
 * STATUS — Node→Go cut-over, Stage 3 (see docs/MIGRATION_NODE_TO_GO.md §5):
 * `skill-router registry build` (Go) is now the AUTHORITATIVE generator and the
 * blocking `--check` drift gate (byte-parity proven across all 4 artifacts in
 * both modes). This script is retained ONLY as an independent, non-blocking
 * parity ORACLE and is slated for removal at Stage 5 after one clean release.
 * Keep its output byte-identical to the Go generator until then.
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
 * Modes (the committed registries ARE the optimized output, so optimize is the
 * default; --faithful selects the legacy characterization mode):
 *   --check            generate in memory, diff against the committed files,
 *                      exit 1 on drift. No writes. Default covers ALL four
 *                      artifacts + the stale-duplicate guard. (Use in CI.)
 *   --write            write the generated artifacts to disk.
 *   --print <artifact> print one artifact to stdout (manifest|marketplace|codex-marketplace|build-manifest).
 *   --faithful         reproduce the legacy manifest.json/marketplace.json byte-
 *                      for-byte instead of the optimized form (characterization).
 *   --optimize         explicit form of the default optimized output.
 *   --only <list>      restrict to a comma list of artifacts.
 *
 * Optimize (default) vs --faithful:
 *   faithful   reproduces the legacy manifest.json/marketplace.json byte-for-byte
 *              (the refactor-only proof). build_manifest is not reproduced in this
 *              mode because it is intentionally restructured by optimize.
 *   optimize   - manifest entries omit empty optional fields (has_scripts:false,
 *                empty scripts/aliases) relying on the Go reader's `omitempty`
 *              - build_manifest slimmed to provenance only (drops the duplicated
 *                ~750KB catalog; nothing consumes its skills[])
 *              - build_manifest paths made portable (relative, not %USERPROFILE%)
 *              - counts recomputed from the live tree (kills the 1812/1811 drift)
 *              - marketplace gains the 14 themed groupings + live skill count
 *
 * Invariants preserved in BOTH modes (breaking these is CHANGES_REQUESTED):
 *   - manifest.routing.compatibility_access records legacy command aliases
 *   - build_manifest.compatibility_binary_aliases records legacy binary names
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

// Registries that were collapsed into the canonical single source and must NOT
// reappear. plugin/marketplace.json was a stray byte-duplicate of the root
// aggregate (it points at ./plugin, which only resolves from the repo root; the
// plugin is self-described by plugin/plugin.json). --check fails if any return.
export const STALE_REGISTRIES = ["plugin/marketplace.json"];

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
    routing: config.manifest.routing,
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
    // Preserve the former marketplace's only unique value: themed collections.
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
    compatibility_binary_aliases: b.compatibility_binary_aliases,
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
// drift comparison — compare by MEANING, not bytes. This is intentionally
// invariant to the faithful-vs-optimize formatting choice so the guard reports
// real drift (skills added/removed, descriptions/aliases/script-sets changed)
// regardless of which form is committed:
//   - drop volatile / mode-dependent fields (generated, alias_count, the
//     marketplace skill-count token, groupings)
//   - canonicalize each entry: aliases default [], has_scripts := scripts
//     non-empty, scripts sorted (the Go validator sorts both sides anyway)
// ---------------------------------------------------------------------------
function canonicalEntry(e) {
  const scripts = (e.scripts || []).slice().sort();
  return {
    name: e.name,
    directory: e.directory,
    description: e.description,
    aliases: e.aliases && e.aliases.length ? e.aliases : [],
    has_scripts: scripts.length > 0,
    scripts,
  };
}

export function normalizeForCompare(key, text) {
  let obj;
  try {
    obj = JSON.parse(text);
  } catch {
    return text;
  }
  if (key === "manifest") {
    return serialize({
      version: obj.version,
      description: obj.description,
      canonical_id_policy: obj.canonical_id_policy,
      routing: obj.routing,
      total_skills: obj.total_skills,
      core_skills: (obj.core_skills || []).map(canonicalEntry),
      library_skills: (obj.library_skills || []).map(canonicalEntry),
    });
  }
  if (key === "marketplace") {
    return serialize({
      name: obj.name,
      owner: obj.owner,
      plugins: (obj.plugins || []).map((p) => ({
        name: p.name,
        source: p.source,
        version: p.version,
        author: p.author,
      })),
    });
  }
  if (key === "build-manifest") {
    delete obj.generated_at;
    (obj.skills || []).forEach((e) => {
      if (Array.isArray(e.scripts)) e.scripts = e.scripts.slice().sort();
    });
    return serialize(obj);
  }
  return text;
}

// ---------------------------------------------------------------------------
// main
// ---------------------------------------------------------------------------
function parseArgs(argv) {
  // The committed registries ARE the optimized output, so optimize is the
  // default; --faithful selects the legacy-reproduction (characterization) mode.
  const args = { mode: "check", optimize: true, only: null, printArtifact: null };
  for (let i = 0; i < argv.length; i++) {
    const a = argv[i];
    if (a === "--write") args.mode = "write";
    else if (a === "--check") args.mode = "check";
    else if (a === "--optimize") args.optimize = true;
    else if (a === "--faithful") args.optimize = false;
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

  // Default (optimize): check ALL four artifacts byte-for-byte against the
  // committed tree. --faithful only reproduces manifest+marketplace (the legacy
  // characterization set); build-manifest/codex have no faithful committed form.
  const FAITHFUL_CHECK = ["manifest", "marketplace"];
  let selected;
  if (args.only) selected = args.only.filter((k) => k in ARTIFACTS);
  else if (args.mode === "check" && !args.optimize) selected = FAITHFUL_CHECK;
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
  // Stale-duplicate guard: collapsed registries must not reappear.
  for (const rel of STALE_REGISTRIES) {
    if (fs.existsSync(path.join(REPO_ROOT, rel))) {
      console.error(`DRIFT: ${rel} is a stale duplicate registry (collapsed into the canonical marketplace.json — delete it)`);
      drift++;
    } else {
      console.log(`ok: ${rel} absent (collapsed)`);
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
