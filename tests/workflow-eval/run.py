"""Score held-out natural requests against the previous pipeline and section plans."""
import hashlib
import json
import os
from pathlib import Path
import re
import subprocess
import tempfile

HERE = Path(__file__).resolve().parent
REPO = HERE.parents[1]
CASES = json.loads((HERE / "cases.json").read_text())
ROUTER = os.environ.get("SKILL_ROUTER_BIN", "skill-router")


def source_text(skill, file, start=1, end=None):
    path = REPO / "skills" / skill / file
    body = path.read_text(encoding="utf-8").splitlines()
    return "\n".join(body[start - 1:end])


def covered(check, entries):
    return any((not check.get("skills") or skill in check["skills"])
               and (not check.get("kind") or kind == check["kind"])
               and (not check.get("inapplicable") or not re.search(check["inapplicable"], body, re.I))
               and re.search(check["source"], body, re.I)
               for skill, kind, body in entries)


def run(args, env):
    p = subprocess.run([ROUTER, *args], cwd=REPO, env=env, capture_output=True,
                       text=True, check=True, timeout=240)
    return json.loads(p.stdout)


def verify_inventory(env):
    result = run(["skills", "workflow", "index", "--inventory"], env)
    inventory = result["inventory"]
    manifest = json.loads((REPO / "manifest.json").read_text())
    names = {s["name"] for s in manifest["core_skills"] + manifest["library_skills"]}
    assert set(inventory) == names, "indexed bundles differ from canonical manifest"
    files = subprocess.check_output(["git", "ls-files", "-z", "--cached", "--others",
                                     "--exclude-standard", "--", "skills"], cwd=REPO)
    tracked = {os.fsdecode(f) for f in files.split(b"\0") if f}
    expected = {p for p in tracked if len(p.split("/")) > 2 and p.split("/")[1] in names}
    found = {f"skills/{name}/{file}": digest for name, items in inventory.items()
             for file, digest in items.items()}
    assert set(found) == expected, "indexed files differ from canonical source tree"
    for file, digest in found.items():
        path = (REPO / file).resolve()
        path.relative_to((REPO / "skills" / file.split("/")[1]).resolve())
        assert hashlib.sha256(path.read_bytes()).hexdigest() == digest, file
    print(f"Verified {len(names)} bundles and {len(found)} file digests", flush=True)


def main():
    totals = [0, 0, 0]
    failed = []
    with tempfile.TemporaryDirectory() as tmp:
        env = {**os.environ, "SKILL_ROUTER_REPO_DIR": str(REPO),
               "SKILL_ROUTER_SKILLS_DIR": str(REPO / "skills"), "HOME": tmp,
               "USERPROFILE": tmp, "SKILL_ROUTER_CONFIG_DIR": tmp,
               "XDG_CACHE_HOME": tmp, "SKILL_ROUTER_TELEMETRY": "0"}
        verify_inventory(env)
        for case in CASES:
            prompt = case["prompt"]
            old = run(["skills", "compose", "--pipeline", "--json", "--", prompt], env)
            legacy = [(s["skill"], "verification" if "test.author" in s["capabilities"] else "",
                       source_text(s["skill"], "SKILL.md"))
                      for s in old["steps"] if s.get("skill") and s["decision"] == "route"]
            new = run(["skills", "workflow", "plan", "--", prompt], env)
            planned = [(c["skill"], stage["requirement"]["kind"], source_text(
                           c["skill"], e["file"], e["startLine"], e["endLine"]))
                       for stage in new["stages"] for c in stage.get("contributions") or []
                       for e in c["elements"] if not e.get("binary")]
            scores = [sum(covered(c, entries) for c in case["checks"]) for entries in (legacy, planned)]
            scores.append(len(case["checks"]))
            totals = [a + b for a, b in zip(totals, scores)]
            missing = [c["purpose"] for c in case["checks"] if not covered(c, planned)]
            if missing:
                failed.extend(f"{case['id']}: {purpose}" for purpose in missing)
            print(f"{case['id']}: pipeline {scores[0]}/{scores[2]}, workflow {scores[1]}/{scores[2]}", flush=True)
    print(f"All cases: pipeline {totals[0]}/{totals[2]}, workflow {totals[1]}/{totals[2]}")
    if failed:
        raise SystemExit("Missing required skill evidence: " + ", ".join(failed))


if __name__ == "__main__":
    main()
