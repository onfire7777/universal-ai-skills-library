/**
 * seed-config.mjs — one-time bootstrap for the unified registry source of truth.
 *
 * Historically the corpus carried TWO hand-maintained registries that could (and
 * did) drift: manifest.json (router catalog) and docs/build_manifest.json
 * (provenance), plus a duplicated marketplace.json. This script captures the
 * *non-derivable* curated data from those legacy files exactly once and writes
 * it to registry.config.json — the new single source of truth. From then on the
 * catalog itself is derived from skills/ on disk by generate-registry.mjs, and
 * only curated metadata (core ordering, aliases, description overrides, plugin
 * groupings) lives in the config.
 *
 * Curated (kept in config, not derivable):  core ordering, per-skill aliases,
 *   curated description overrides, package/marketplace metadata, the build-
 *   manifest provenance blocks (merged legacy dirs, disabled aliases, policy).
 * Derived (NOT stored in config):  the skill list, descriptions that match
 *   SKILL.md, has_scripts/scripts, total counts.
 *
 * Usage:  node scripts/registry/seed-config.mjs [--force]
 * Idempotent in spirit: re-running reproduces the same config from the same
 * legacy inputs. Refuses to overwrite an existing config unless --force.
 */
import fs from "node:fs";
import path from "node:path";
import { fileURLToPath } from "node:url";
import { parseFrontMatter, collapseWhitespace } from "./lib/frontmatter.mjs";

const HERE = path.dirname(fileURLToPath(import.meta.url));
const REPO_ROOT = path.resolve(HERE, "..", "..");
const CONFIG_PATH = path.join(HERE, "registry.config.json");

function readJson(rel) {
  return JSON.parse(fs.readFileSync(path.join(REPO_ROOT, rel), "utf8"));
}

function main() {
  const force = process.argv.includes("--force");
  if (fs.existsSync(CONFIG_PATH) && !force) {
    console.error(`refusing to overwrite ${path.relative(REPO_ROOT, CONFIG_PATH)} (use --force)`);
    process.exit(1);
  }

  const manifest = readJson("manifest.json");
  const marketplace = readJson("marketplace.json");
  const build = readJson("docs/build_manifest.json");

  const allEntries = [...manifest.core_skills, ...manifest.library_skills];

  // --- core ordering (curated; library is derived by alpha sort) ---
  const coreSkills = manifest.core_skills.map((s) => s.name);

  // --- per-skill aliases (curated) ---
  const aliases = {};
  for (const entry of allEntries) {
    if (Array.isArray(entry.aliases) && entry.aliases.length > 0) {
      aliases[entry.name] = entry.aliases;
    }
  }

  // --- description overrides: only where the manifest deliberately diverges
  //     from the skill's own SKILL.md front matter (router-tuned wording) ---
  const descriptionOverrides = {};
  for (const entry of allEntries) {
    const skillMd = path.join(REPO_ROOT, entry.directory, "SKILL.md");
    let derived;
    try {
      derived = parseFrontMatter(fs.readFileSync(skillMd, "utf8")).description;
    } catch {
      derived = undefined;
    }
    const normalizedDerived = derived == null ? "" : collapseWhitespace(derived);
    if (normalizedDerived !== entry.description) {
      descriptionOverrides[entry.name] = entry.description;
    }
  }

  const config = {
    $schema: "./registry.config.schema.json",
    $doc: "Single source of truth for the unified skill registry. The skill CATALOG is derived from skills/ on disk; this file holds only the curated, non-derivable metadata. Generated artifacts (manifest.json, marketplace.json, docs/build_manifest.json) are produced by generate-registry.mjs and must never be hand-edited.",
    schemaVersion: 1,
    manifest: {
      version: manifest.version,
      generated: manifest.generated,
      description: manifest.description,
      canonical_id_policy: manifest.canonical_id_policy,
      // alias_count is a legacy informational figure that is not reproducible
      // from the catalog alone; carried verbatim for byte-faithful output and
      // recomputed deterministically in --optimize mode.
      alias_count: manifest.alias_count,
      routing: manifest.routing,
    },
    coreSkills,
    aliases,
    descriptionOverrides,
    marketplace,
    buildManifest: {
      schema: build.schema,
      generated_at: build.generated_at,
      primary_binary: build.primary_binary,
      compatibility_binary_aliases: build.compatibility_binary_aliases || [],
      // Legacy Windows-only paths preserved for provenance; the generator emits
      // the portable relative forms below instead.
      legacy_source_of_truth: build.source_of_truth,
      legacy_router_source: build.router_source,
      source_of_truth: "skills",
      router_source: "skill-router-cli",
      alias_count: build.alias_count,
      merged_legacy_directories: build.merged_legacy_directories,
      disabled_colliding_aliases: build.disabled_colliding_aliases,
      compatibility_policy: build.compatibility_policy,
    },
    // Themed plugin groupings (the former marketplace's unique value). Populated
    // from Builder 3's grouping map; each: { name, description, members: [] }.
    groupings: [],
  };

  fs.writeFileSync(CONFIG_PATH, JSON.stringify(config, null, 2) + "\n");
  console.log(
    `wrote ${path.relative(REPO_ROOT, CONFIG_PATH)}: ` +
      `${coreSkills.length} core, ${Object.keys(aliases).length} aliased, ` +
      `${Object.keys(descriptionOverrides).length} description overrides, ` +
      `${config.buildManifest.merged_legacy_directories.length} merged legacy dirs`
  );
}

main();
