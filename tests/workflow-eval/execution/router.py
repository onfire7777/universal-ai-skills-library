"""Exercise the public workflow CLI with unmodified canonical skill bundles."""

import hashlib
import json
import os
from pathlib import Path
import shutil
import subprocess

REPO = Path(__file__).resolve().parents[3]


def digest(data):
    return hashlib.sha256(data).hexdigest()


def write_json(path, value):
    path.write_text(json.dumps(value, indent=2, ensure_ascii=False) + "\n", encoding="utf-8")


class Router:
    def __init__(self, temp, names):
        executable = os.environ.get("SKILL_ROUTER_BIN") or shutil.which("skill-router")
        if not executable:
            raise RuntimeError("Build the router and set SKILL_ROUTER_BIN")
        self.bin = str(Path(executable).resolve())
        checkout = temp / "catalog"
        checkout.mkdir()
        manifest = json.loads((REPO / "manifest.json").read_text())
        for key in ("core_skills", "library_skills"):
            manifest[key] = [item for item in manifest[key] if item["name"] in names]
        named = manifest["core_skills"] + manifest["library_skills"]
        assert {item["name"] for item in named} == set(names)
        manifest["total_skills"] = len(named)
        manifest["alias_count"] = sum(len(item.get("aliases", [])) for item in named)
        write_json(checkout / "manifest.json", manifest)
        for item in named:
            relative = Path(item["directory"])
            shutil.copytree(REPO / relative, checkout / relative)
        self.env = {key: val for key, val in os.environ.items() if not key.startswith("SKILL_ROUTER_")}
        self.env.update(SKILL_ROUTER_REPO_DIR=str(checkout), SKILL_ROUTER_CONFIG_DIR=str(temp / "cache"))
        self.revision = subprocess.check_output(["git", "-C", str(REPO), "rev-parse", "HEAD"], text=True).strip()

    def call(self, verb, request):
        command = [self.bin, "skills", "workflow", verb, "--request", "-"]
        result = subprocess.run(command, input=json.dumps(request), text=True, capture_output=True,
                                cwd=REPO, env=self.env, timeout=90)
        if result.returncode:
            raise RuntimeError(f"workflow {verb} failed: {result.stderr.strip()}")
        cli_result = json.loads(result.stdout)
        # MCP must accept exactly the same structured request and deliver the same plan,
        # hydrated assets and receipts, or criterion-level review, as the CLI.
        rpc = {"jsonrpc": "2.0", "id": 1, "method": "tools/call",
               "params": {"name": "workflow_" + verb, "arguments": request}}
        server = subprocess.run([self.bin, "serve"], input=json.dumps(rpc) + "\n", capture_output=True,
                                text=True, cwd=REPO, env=self.env, timeout=90)
        if server.returncode:
            raise RuntimeError(f"MCP workflow_{verb} failed: {server.stderr.strip()}")
        packet = json.loads(server.stdout)
        assert "error" not in packet, packet.get("error")
        response = packet["result"]
        assert not response.get("isError"), response["content"]
        mcp_result = json.loads(response["content"][0]["text"])
        assert mcp_result == cli_result, f"CLI and MCP disagree for workflow_{verb}"
        return cli_result

    def plan(self, out, prompt, stages):
        requirements = []
        for spec in stages:
            req = {key: val for key, val in spec.items() if key != "sources"}
            req["selection"] = "explicit"
            req["elements"] = []
            for locator in spec["sources"]:
                skill, title, phrase = locator[:3]
                file = locator[3] if len(locator) == 4 else "SKILL.md"
                query = f"{Path(file).name} {title} {phrase}" if file != "SKILL.md" else phrase
                search = self.call("search", {"prompt": prompt, "query": query, "skill": skill, "limit": 100})
                sections = [hit["element"] for hit in search["hits"]
                            if hit["element"]["file"] == file and hit["element"]["title"] == title]
                assert len(sections) == 1, f"Expected {skill}/{title} in search: {search.get('reason')}"
                section = sections[0]
                raw = (REPO / "skills" / skill / section["file"]).read_bytes()
                assert digest(raw) == section["digest"]
                lines = raw.decode().splitlines()[section["startLine"] - 1:section["endLine"]]
                assert phrase.casefold() in "\n".join(lines).casefold(), (skill, title, lines)
                req["elements"].append({k: section[k] for k in ("skill", "file", "startLine", "endLine", "digest")})
            requirements.append(req)
        plan = self.call("plan", {"prompt": prompt, "requirements": requirements})
        assert not plan.get("gaps"), plan.get("gaps")
        assert [s["requirement"]["id"] for s in plan["stages"]] == [s["id"] for s in stages]
        for stage in plan["stages"]:
            facets = stage["requirement"].get("facets", [])
            coverage = stage.get("facetCoverage", [])
            assert {f["id"] for f in facets} == {item["facet"]["id"] for item in coverage}
            for item in coverage:
                assert item["satisfied"] and item["matches"], item
                for match in item["matches"]:
                    contribution = next(c for c in stage["contributions"] if c["id"] == match["contribution"])
                    source = next(e for e in contribution["elements"] if f"{e['file']}#L{e['startLine']}-L{e['endLine']}" == match["element"])
                    assert source["startLine"] <= match["line"] <= source["endLine"]
                    lines = (REPO / "skills" / source["skill"] / source["file"]).read_text().splitlines()
                    assert match["term"].casefold() in lines[match["line"] - 1].casefold()
        write_json(out / "plan.json", plan)
        return plan

    def load(self, out, plan, stage_id):
        stage = next(s for s in plan["stages"] if s["requirement"]["id"] == stage_id)
        selected = {c["id"]: c for c in stage["contributions"]}
        load = self.call("load", {"plan": plan, "stage": stage_id})
        assert not load.get("gaps") and not load.get("blocked"), load.get("gaps")
        assert {c["id"] for c in load["contributions"]} == set(selected)
        records, sources = [], []
        for contribution in load["contributions"]:
            skill = contribution["skill"]
            base = REPO / "skills" / skill
            entry = (base / "SKILL.md").read_bytes()
            assert entry.decode() == contribution["body"]
            assert entry == subprocess.check_output(["git", "-C", str(REPO), "show", f"{self.revision}:skills/{skill}/SKILL.md"]), skill
            read_files = []
            for asset in contribution["assets"]:
                content = (base / asset["file"]).read_bytes()
                assert digest(content) == asset["digest"] and len(content) == asset["bytes"]
                revision_file = f"{self.revision}:skills/{skill}/{asset['file']}"
                assert subprocess.check_output(["git", "-C", str(REPO), "show", revision_file]) == content
                if asset.get("path"):
                    delivered = Path(asset["path"]).read_bytes()
                    read_files.append(asset["file"])
                else:
                    delivered = asset.get("body", "").encode("utf-8")
                assert delivered == content
                saved = out / "sources" / skill / asset["file"]
                saved.parent.mkdir(parents=True, exist_ok=True)
                saved.write_bytes(delivered)
            elements = selected[contribution["id"]]["elements"]
            addresses = [f"{el['file']}#L{el['startLine']}-L{el['endLine']}" for el in elements]
            assert contribution["receipt"] and addresses
            records.append({"contribution": contribution["id"], "elements": addresses,
                            "readFiles": read_files, "receipt": contribution["receipt"]})
            sources.append({"skill": skill, "checkoutRevision": self.revision, "entrypointSHA256": digest(entry),
                            "elements": addresses, "loadedAssetSHA256": {a["file"]: a["digest"] for a in contribution["assets"]}})
        write_json(out / f"{stage_id}-sources.json", sources)
        return records


def record(contributions, status, observations):
    return [{**contribution, "status": status, "evidence": observations} for contribution in contributions]


def evidence(check, kind, reference, observation, outcome=None):
    e = {"checks": [check], "kind": kind, "reference": reference, "observation": observation}
    if outcome:
        e["outcome"] = outcome
    return e


def check(id_, text, kind):
    return {"id": id_, "text": text, "evidenceKind": kind}


def review(router, out, plan, records, name, missing_checks=()):
    result = router.call("review", {"plan": plan, "records": records})
    write_json(out / f"review-{name}.json", result)
    missing = {c["id"] for stage in result["stages"] for c in stage["checks"] if not c["satisfied"]}
    assert missing == set(missing_checks), (name, missing, missing_checks)
    assert result["complete"] == (not missing_checks), (name, result)
    if not missing_checks:
        for stage in result["stages"]:
            assert all(c["satisfied"] and c["evidence"] for c in stage["checks"]), stage["checks"]
    return result
