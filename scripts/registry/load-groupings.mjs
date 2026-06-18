/**
 * load-groupings.mjs — fold themed plugin groupings into the registry config.
 *
 * The former marketplace's only non-duplicate value was its 14 themed plugin
 * groupings (security, data-analysis, ...). Every member skill already exists in
 * skills/, so the groupings are pure universal category metadata over existing kebab ids.
 * This script imports that metadata into registry.config.json (self-contained,
 * no external file dependency) after validating every member resolves to a real
 * skills/<id>/ directory. The generator then surfaces them in marketplace.json
 * under --optimize.
 *
 * Usage: node scripts/registry/load-groupings.mjs <groupings.json> [--write]
 *   <groupings.json>  { source_repo, plugin_count, groupings: [{ name|plugin,
 *                       description, members: [kebab-id] }] }
 *   --write           persist into registry.config.json (otherwise dry-run)
 */
import fs from "node:fs";
import path from "node:path";
import { fileURLToPath } from "node:url";

const HERE = path.dirname(fileURLToPath(import.meta.url));
const REPO_ROOT = path.resolve(HERE, "..", "..");
const CONFIG_PATH = path.join(HERE, "registry.config.json");

function main() {
  const src = process.argv[2];
  const write = process.argv.includes("--write");
  if (!src) {
    console.error("usage: node scripts/registry/load-groupings.mjs <groupings.json> [--write]");
    process.exit(2);
  }

  const raw = JSON.parse(fs.readFileSync(src, "utf8"));
  const input = Array.isArray(raw) ? raw : raw.groupings;
  if (!Array.isArray(input)) {
    console.error("groupings file has no 'groupings' array");
    process.exit(2);
  }

  const skillExists = (id) =>
    fs.existsSync(path.join(REPO_ROOT, "skills", id, "SKILL.md"));

  const groupings = [];
  const missing = [];
  for (const g of input) {
    const name = g.name || g.plugin;
    const members = g.members || g.skills || [];
    const resolved = [];
    for (const id of members) {
      if (skillExists(id)) resolved.push(id);
      else missing.push(`${name}:${id}`);
    }
    groupings.push({
      name,
      description: g.description || "",
      members: resolved.sort(),
    });
  }

  const totalMembers = groupings.reduce((n, g) => n + g.members.length, 0);
  console.log(
    `${groupings.length} groupings, ${totalMembers} resolved members, ${missing.length} unresolved`
  );
  if (missing.length > 0) {
    console.error("UNRESOLVED members (not in skills/): " + missing.slice(0, 20).join(", "));
    // Unresolved members would dangle in the registry — fail rather than ship them.
    process.exit(1);
  }

  if (!write) {
    console.log("dry-run (pass --write to persist into registry.config.json)");
    return;
  }

  const config = JSON.parse(fs.readFileSync(CONFIG_PATH, "utf8"));
  config.groupings = groupings.sort((a, b) =>
    a.name < b.name ? -1 : a.name > b.name ? 1 : 0
  );
  fs.writeFileSync(CONFIG_PATH, JSON.stringify(config, null, 2) + "\n");
  console.log(`wrote ${groupings.length} groupings into ${path.relative(REPO_ROOT, CONFIG_PATH)}`);
}

main();
