/**
 * Characterization + contract tests for the unified registry generator.
 *
 * Run: node --test scripts/registry/   (or: npm test if wired by Foundation)
 *
 * These tests are the refactor-only proof: the generator reproduces today's
 * manifest.json and build provenance, and the generated output obeys the exact
 * contract enforced by skill-router-cli validate_manifest.go, so the decoupled
 * router reads it unchanged.
 */
import test from "node:test";
import assert from "node:assert/strict";
import fs from "node:fs";
import path from "node:path";

import {
  REPO_ROOT,
  CONFIG_PATH,
  STALE_REGISTRIES,
  serialize,
  listSkillScripts,
  scanSkills,
  buildManifest,
  buildBuildManifest,
  normalizeForCompare,
} from "./generate-registry.mjs";

const config = JSON.parse(fs.readFileSync(CONFIG_PATH, "utf8"));
const skills = scanSkills(config);
const readCommitted = (rel) => fs.readFileSync(path.join(REPO_ROOT, rel), "utf8");
const allEntries = (m) => [...m.core_skills, ...m.library_skills];
const sorted = (a) => (a || []).slice().sort();

test("faithful manifest is semantically identical to the committed manifest", () => {
  const committed = JSON.parse(readCommitted("manifest.json"));
  const generated = buildManifest(config, skills, { optimize: false });

  // top-level metadata + partition sizes. alias_count is intentionally excluded:
  // it is a mode-dependent informational figure (the legacy 1917 was not
  // reproducible from the catalog; --optimize recomputes it deterministically).
  for (const k of ["version", "generated", "description", "canonical_id_policy", "total_skills"]) {
    assert.deepEqual(generated[k], committed[k], `metadata ${k}`);
  }
  assert.deepEqual(generated.routing, committed.routing, "routing block");
  assert.equal(generated.core_skills.length, committed.core_skills.length, "core count");
  assert.equal(generated.library_skills.length, committed.library_skills.length, "library count");

  // core ordering and library ordering are preserved exactly
  assert.deepEqual(
    generated.core_skills.map((s) => s.name),
    committed.core_skills.map((s) => s.name),
    "core ordering"
  );
  assert.deepEqual(
    generated.library_skills.map((s) => s.name),
    committed.library_skills.map((s) => s.name),
    "library ordering"
  );

  // every entry equal field-by-field (scripts compared as a set: the Go reader
  // sorts them, so order is not semantically meaningful)
  const G = new Map(allEntries(generated).map((s) => [s.name, s]));
  for (const c of allEntries(committed)) {
    const g = G.get(c.name);
    assert.ok(g, `generated missing ${c.name}`);
    assert.equal(g.directory, c.directory, `${c.name} directory`);
    assert.equal(g.description, c.description, `${c.name} description`);
    // Boolean-normalize: the optimized form omits has_scripts:false (omitempty).
    assert.equal(Boolean(g.has_scripts), Boolean(c.has_scripts), `${c.name} has_scripts`);
    assert.deepEqual(g.aliases || null, c.aliases || null, `${c.name} aliases`);
    assert.deepEqual(sorted(g.scripts), sorted(c.scripts), `${c.name} scripts (set)`);
  }
});

test("optimize is behaviour-neutral: faithful and optimized manifests are semantically equal", () => {
  const f = buildManifest(config, skills, { optimize: false });
  const o = buildManifest(config, skills, { optimize: true });
  const norm = (m) =>
    allEntries(m)
      .map((e) => ({
        name: e.name,
        directory: e.directory,
        description: e.description,
        aliases: e.aliases && e.aliases.length ? e.aliases : [],
        has_scripts: Boolean(e.has_scripts) || (e.scripts || []).length > 0,
        scripts: sorted(e.scripts),
      }))
      .sort((a, b) => (a.name < b.name ? -1 : 1));
  assert.deepEqual(norm(f), norm(o), "faithful vs optimize entries");
  assert.equal(f.total_skills, o.total_skills, "total_skills");
  assert.deepEqual(f.routing, o.routing, "routing");
});

test("generated manifest obeys the validate-manifest contract (faithful + optimize)", () => {
  for (const optimize of [false, true]) {
    const m = buildManifest(config, skills, { optimize });
    const entries = allEntries(m);
    const names = new Set();
    const dirs = new Set();
    for (const e of entries) {
      // unique names
      assert.ok(!names.has(e.name.toLowerCase()), `duplicate name ${e.name}`);
      names.add(e.name.toLowerCase());
      // unique, safe directories shaped skills/<name>
      assert.equal(e.directory, `skills/${e.name}`, `${e.name} directory shape`);
      assert.ok(!dirs.has(e.directory), `duplicate dir ${e.directory}`);
      dirs.add(e.directory);
      assert.ok(!e.directory.includes(".."), `unsafe dir ${e.directory}`);
      assert.ok(!path.isAbsolute(e.directory), `absolute dir ${e.directory}`);
      // SKILL.md exists
      assert.ok(
        fs.existsSync(path.join(REPO_ROOT, e.directory, "SKILL.md")),
        `${e.name} missing SKILL.md`
      );
      // has_scripts / scripts match the Go validator's listSkillScripts exactly
      const actual = listSkillScripts(path.join(REPO_ROOT, e.directory));
      const listed = sorted(e.scripts);
      assert.deepEqual(listed, actual, `${e.name} scripts match disk`);
      const hasScripts = optimize ? e.has_scripts === true : e.has_scripts;
      assert.equal(Boolean(hasScripts), actual.length > 0, `${e.name} has_scripts correctness`);
    }
  }
});

test("every skills/ directory with a SKILL.md is indexed (no unindexed top dirs)", () => {
  const m = buildManifest(config, skills, { optimize: false });
  const indexed = new Set(allEntries(m).map((s) => s.directory));
  const onDisk = fs
    .readdirSync(path.join(REPO_ROOT, "skills"), { withFileTypes: true })
    .filter((e) => e.isDirectory() && fs.existsSync(path.join(REPO_ROOT, "skills", e.name, "SKILL.md")))
    .map((e) => `skills/${e.name}`);
  for (const dir of onDisk) assert.ok(indexed.has(dir), `unindexed: ${dir}`);
  assert.equal(indexed.size, onDisk.length, "indexed count equals on-disk count");
});

test("compatibility aliases are preserved in generic fields", () => {
  for (const optimize of [false, true]) {
    const m = buildManifest(config, skills, { optimize });
    assert.ok(Array.isArray(m.routing.compatibility_access), "routing.compatibility_access");
    assert.ok(
      m.routing.compatibility_access.includes("manus skill <name>"),
      "routing compatibility command alias"
    );
    const b = buildBuildManifest(config, skills, { optimize });
    assert.deepEqual(b.compatibility_binary_aliases, ["manus"], "compatibility_binary_aliases");
  }
});

test("optimize trims empty optional fields from manifest entries", () => {
  const m = buildManifest(config, skills, { optimize: true });
  for (const e of allEntries(m)) {
    assert.ok(!("has_scripts" in e) || e.has_scripts === true, `${e.name} has_scripts:false present`);
    assert.ok(!("scripts" in e) || e.scripts.length > 0, `${e.name} empty scripts present`);
    assert.ok(!("aliases" in e) || e.aliases.length > 0, `${e.name} empty aliases present`);
  }
});

test("optimize slims build_manifest to provenance and fixes drift + portability", () => {
  const b = buildBuildManifest(config, skills, { optimize: true });
  assert.ok(!("skills" in b), "skills[] catalog dropped");
  assert.equal(b.catalog_ref, "../manifest.json", "points at canonical catalog");
  assert.equal(b.skill_count, skills.length, "skill_count matches live tree");
  assert.equal(b.directories_total, skills.length, "directories_total matches live tree");
  // portable, not Windows %USERPROFILE%
  assert.ok(!/%USERPROFILE%|\\/.test(b.source_of_truth), "source_of_truth portable");
  assert.ok(!/%USERPROFILE%|\\/.test(b.router_source), "router_source portable");
  // curated provenance preserved
  assert.ok(Array.isArray(b.merged_legacy_directories) && b.merged_legacy_directories.length > 0);
  assert.ok(b.disabled_colliding_aliases && Object.keys(b.disabled_colliding_aliases).length > 0);
  assert.ok(b.compatibility_policy && typeof b.compatibility_policy === "object");
});

test("build_manifest no longer drifts from manifest (single scan, equal counts)", () => {
  const m = buildManifest(config, skills, { optimize: true });
  const b = buildBuildManifest(config, skills, { optimize: true });
  assert.equal(b.skill_count, m.total_skills, "build_manifest count == manifest count");
});

test("collapsed stale duplicate registries do not exist on disk", () => {
  for (const rel of STALE_REGISTRIES) {
    assert.ok(
      !fs.existsSync(path.join(REPO_ROOT, rel)),
      `${rel} is a retired marketplace registry and must be deleted`
    );
  }
});

test("normalizeForCompare treats scripts[] order as insignificant", () => {
  const a = { core_skills: [{ name: "x", directory: "skills/x", description: "d", has_scripts: true, scripts: ["scripts/b.py", "scripts/a.py"] }], library_skills: [] };
  const b = { core_skills: [{ name: "x", directory: "skills/x", description: "d", has_scripts: true, scripts: ["scripts/a.py", "scripts/b.py"] }], library_skills: [] };
  assert.equal(normalizeForCompare("manifest", serialize(a)), normalizeForCompare("manifest", serialize(b)));
});

test("--check (faithful default) reports manifest in sync today", () => {
  // mirrors the CLI default-selection behavior without spawning a subprocess
  for (const [key, rel] of [["manifest", "manifest.json"]]) {
    const built = buildManifest(config, skills, { optimize: false });
    const generated = normalizeForCompare(key, serialize(built));
    const committed = normalizeForCompare(key, readCommitted(rel));
    assert.equal(generated, committed, `${rel} drifted`);
  }
});
