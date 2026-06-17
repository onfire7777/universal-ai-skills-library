#!/usr/bin/env python3
"""Regenerate the characterization golden baselines.

Run this **deliberately** to refresh the pinned baselines after an intended,
reviewed behaviour change:

    python3 tests/characterization/update_baseline.py

It writes two artifacts consumed by the test suite:

* ``baseline/routing_golden.json`` — the router's decision + chosen skill for
  every prompt in :data:`harness.PROMPTS`, evaluated against the frozen fixture.
* ``baseline/registry_baseline.json`` — a fingerprint of the *live* skills
  corpus (counts + sorted ``(name, directory)`` pairs) used by the
  "no existing skill is ever lost" guard.

It intentionally does **not** touch ``baseline/go_known_failures.json``: that
file is hand-curated so the known-red Go baseline can only shrink on purpose,
never silently absorb a newly introduced failure.
"""

from __future__ import annotations

import json
import os
import sys

sys.path.insert(0, os.path.dirname(os.path.abspath(__file__)))
import harness  # noqa: E402


def build_routing_golden() -> dict:
    entries = []
    for spec in harness.PROMPTS:
        out = harness.preflight(spec["prompt"], fixture=True)
        best = out.get("best") or {}
        entries.append({
            "id": spec["id"],
            "prompt": spec["prompt"],
            "decision": out.get("decision"),
            "best_name": best.get("name"),
            "best_source": best.get("source"),
            "best_score": best.get("score"),  # diagnostics only; not asserted
        })
    return {
        "_comment": (
            "Pinned router outcomes against tests/fixtures/skills-lib. Tests assert "
            "decision + best_name + best_source; best_score is diagnostic only."
        ),
        "fixture": os.path.relpath(harness.FIXTURE_DIR, harness.repo_root()),
        "entries": entries,
    }


def build_registry_baseline() -> dict:
    manifest = harness.load_manifest()
    pairs = harness.manifest_skill_pairs(manifest)
    vm = harness.validate_manifest(fixture=False)
    return {
        "_comment": (
            "Fingerprint of the live skills corpus at baseline. The registry guard "
            "requires the post-refactor corpus to remain a SUPERSET of these skills "
            "(no existing skill may be lost) and to preserve each skill's directory."
        ),
        "manifest_version": manifest.get("version"),
        "counts": {
            "core": vm.get("coreSkills"),
            "library": vm.get("librarySkills"),
            "total": vm.get("totalSkills"),
        },
        "validate_manifest_ok": vm.get("ok"),
        "fingerprint_sha256": harness.corpus_fingerprint(pairs),
        "skills": [[name, directory] for name, directory in pairs],
    }


def main() -> int:
    os.makedirs(harness.BASELINE_DIR, exist_ok=True)

    routing = build_routing_golden()
    with open(harness.ROUTING_GOLDEN, "w", encoding="utf-8") as fh:
        json.dump(routing, fh, indent=2)
        fh.write("\n")
    print(f"wrote {harness.ROUTING_GOLDEN} ({len(routing['entries'])} prompts)")

    registry = build_registry_baseline()
    with open(harness.REGISTRY_BASELINE, "w", encoding="utf-8") as fh:
        json.dump(registry, fh, indent=2)
        fh.write("\n")
    print(f"wrote {harness.REGISTRY_BASELINE} "
          f"(total={registry['counts']['total']}, fp={registry['fingerprint_sha256'][:12]})")
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
