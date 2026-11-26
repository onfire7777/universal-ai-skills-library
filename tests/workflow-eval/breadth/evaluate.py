"""Evaluate workflow coverage and relevance against independently sourced criteria."""
import argparse
import hashlib
import json
import os
from pathlib import Path
from queue import Empty, Queue
from threading import Thread
from time import monotonic
import re
import shutil
import subprocess
import sys
import tempfile

HERE = Path(__file__).resolve().parent
REPO = HERE.parents[2]
CASES = HERE / "cases.json"
FACETS = HERE / "host-facets.json"
BASELINE = HERE / "baseline.json"
POLICY = HERE / "acceptance.json"


def sha(data):
    return hashlib.sha256(data).hexdigest()


class Router:
    """Reuse the same MCP process and its validated corpus snapshot for all cases."""

    def __init__(self, binary, repository, cache):
        actual = shutil.which(str(binary))
        if not actual:
            raise ValueError(f"Router binary not found: {binary}")
        self.binary = Path(actual).resolve()
        self.repo = repository.resolve()
        self.ordinal = 0
        self.env = {
            **os.environ, "SKILL_ROUTER_REPO_DIR": str(self.repo),
            "SKILL_ROUTER_MANIFEST": str(self.repo / "manifest.json"),
            "SKILL_ROUTER_SKILLS_SOURCE_DIR": str(self.repo / "skills"),
            "SKILL_ROUTER_CONFIG_DIR": str(cache), "SKILL_ROUTER_TELEMETRY": "0",
        }
        self.proc = None

    def __enter__(self):
        self.proc = subprocess.Popen([str(self.binary), "serve"], cwd=self.repo, env=self.env,
                                     stdin=subprocess.PIPE, stdout=subprocess.PIPE,
                                     encoding="utf-8")
        self.stdout = Queue()
        def collect():
            for line in self.proc.stdout:
                self.stdout.put(line)
            self.stdout.put(None)
        Thread(target=collect, daemon=True).start()
        self.rpc("initialize", {"protocolVersion": "2024-11-05", "capabilities": {},
                                "clientInfo": {"name": "workflow-breadth-eval", "version": "1"}})
        self.proc.stdin.write(json.dumps({"jsonrpc": "2.0", "method": "notifications/initialized"}) + "\n")
        self.proc.stdin.flush()
        return self

    def rpc(self, method, params):
        self.ordinal += 1
        call = {"jsonrpc": "2.0", "id": self.ordinal, "method": method, "params": params}
        self.proc.stdin.write(json.dumps(call) + "\n")
        self.proc.stdin.flush()
        deadline = monotonic() + 240
        while True:
            try:
                line = self.stdout.get(timeout=max(0, deadline - monotonic()))
            except Empty as err:
                self.proc.kill()
                raise TimeoutError(f"240 seconds elapsed waiting for {method} reply {self.ordinal}") from err
            if not line:
                raise RuntimeError(f"{method}: router exited ({self.proc.poll()})")
            result = json.loads(line)
            if result.get("id") != self.ordinal:
                continue
            if "error" in result:
                raise RuntimeError(result["error"])
            return result["result"]

    def plan(self, request):
        reply = self.rpc("tools/call", {"name": "workflow_plan", "arguments": request})
        if reply.get("isError"):
            raise RuntimeError(reply["content"][0]["text"])
        return json.loads(reply["content"][0]["text"])

    def __exit__(self, *_):
        self.proc.stdin.close()
        try:
            self.proc.wait(timeout=10)
        except subprocess.TimeoutExpired:
            self.proc.kill()
            self.proc.wait()

    def identity(self):
        return {"binarySHA256": sha(self.binary.read_bytes()),
                "routerVersion": subprocess.check_output([str(self.binary), "--version"], text=True).strip(),
                "manifestSHA256": sha((self.repo / "manifest.json").read_bytes())}


def read_source(file):
    raw = file.read_bytes()
    return sha(raw), raw


def selected_source(repo, contribution, e):
    """Only local source bytes in returned inclusive lines can support a match."""
    skill = contribution["skill"]
    if e["skill"] != skill:
        raise ValueError(f"Element belongs to a different skill: {e}")
    bundle = (repo / "skills" / skill).resolve(strict=True)
    if bundle.parent != (repo / "skills").resolve(strict=True):
        raise ValueError(f"Invalid bundle locator: {skill}")
    file = (bundle / e["file"]).resolve(strict=True)
    if not file.is_relative_to(bundle):
        raise ValueError(f"Invalid asset locator: {e['file']}")
    digest, raw = read_source(file)
    if e["digest"] != digest:
        raise ValueError(f"Stale file locator: {skill}/{e['file']}")
    if e.get("binary"):
        if e.get("startLine") or e.get("endLine") or e["sectionDigest"] != digest:
            raise ValueError(f"Invalid binary locator: {file}")
        return ""
    start, end = e["startLine"], e["endLine"]
    lines = raw.decode("utf-8").split("\n")
    if len(lines) > 1 and lines[-1] == "":
        lines.pop()
    if not 1 <= start <= end <= len(lines):
        raise ValueError(f"Invalid source lines: {file}:{start}-{end}")
    text = "\n".join(lines[start - 1:end]).replace("\r\n", "\n")
    if e["sectionDigest"] != sha(text.strip().encode()):
        raise ValueError(f"Stale section locator: {file}:{start}-{end}")
    return text


def evidence(text, alternatives):
    # Each alternative is a conjunction supported within a local span. This
    # prevents unrelated terms from opposite ends of a long section qualifying.
    for offset in range(0, len(text), 700):
        window = text[offset:offset + 1400]
        for alternative in alternatives:
            matches = [re.search(rx, window, re.I | re.M) for rx in alternative]
            if all(matches):
                first = min(m.start() for m in matches)
                return text.count("\n", 0, offset + first)
    return None


def preregistration(cases, repo):
    if len({c["id"] for c in cases}) != len(cases):
        raise ValueError("Repeated case ID in preregistration")
    for case in cases:
        ids = set()
        for r in case["requirements"]:
            if not r["source"] or any(not alt or not all(alt) for alt in r["source"]):
                raise ValueError(f"{case['id']}/{r['id']}: missing actionable source criteria")
            if r["id"] in ids:
                raise ValueError(f"{case['id']}: repeated requirement {r['id']}")
            ids.add(r["id"])
            if r.get("expected") == "gap":
                if not r["source"]:
                    raise ValueError("A gap must still describe what would count as real support.")
                continue
            path = repo / "skills" / r["witness"]
            text = path.read_text(encoding="utf-8")
            if not any(all(re.search(rx, text, re.I | re.M) for rx in alt) for alt in r["source"]):
                raise ValueError(f"{case['id']}/{r['id']}: registered witness no longer supports its criterion: {path}")


def request_for(case, mode, annotations=None):
    request = {"prompt": case["prompt"]}
    if mode in ("structured", "faceted"):
        reqs = []
        for r in case["requirements"]:
            reqs.append({"id": r["id"], "kind": r["kind"], "text": r["text"],
                         "dependsOn": [reqs[-1]["id"]] if reqs else []})
        if mode == "faceted":
            for req in reqs:
                req["facets"] = annotations[case["id"]][req["id"]]
        request["requirements"] = reqs
    return request


def verify_facet_quotes(stage, repo):
    elements = {}
    for c in stage.get("contributions") or []:
        for e in c["elements"]:
            address = f"{e['file']}#L{e.get('startLine', 0)}-L{e.get('endLine', 0)}"
            elements[(c["id"], address)] = (e, selected_source(repo, c, e))
    proven = 0
    coverage = stage.get("facetCoverage") or []
    expected = {f["id"] for f in stage["requirement"].get("facets", [])}
    reported = {c["facet"]["id"] for c in coverage}
    if stage.get("contributions") and (not expected or not expected <= reported):
        raise ValueError(f"Missing facet receipts for stage {stage['requirement']['id']}: {expected - reported}")
    for facet in coverage:
        matches = facet.get("matches") or []
        if facet["satisfied"] and facet["facet"]["kind"] != "exclude" and not matches:
            raise ValueError(f"Facet has no supporting span: {facet}")
        for m in matches:
            entry = elements.get((m["contribution"], m["element"]))
            if not entry:
                raise ValueError(f"Facet points outside the selected stage: {m}")
            e, body = entry
            lines = body.split("\n")
            offset = m["line"] - e["startLine"]
            # A multiword match can straddle a line break. Its first byte
            # must still actually occur on the reported source line.
            quote = m["term"].replace("\r\n", "\n")
            quoted_at = "\n".join(lines[offset:]).find(quote) if 0 <= offset < len(lines) and quote else -1
            if quoted_at < 0 or quoted_at > len(lines[offset]):
                raise ValueError(f"Facet quote does not occur at its reported line: {m}")
        proven += int(facet["satisfied"])
    return proven, len(coverage)


def score(case, mode, plan, repo):
    results = {r["id"]: {"status": "missing", "evidence": []} for r in case["requirements"]}
    contributions = []
    section_count = credited_sections = satisfied_facets = facet_count = 0
    supported_stages = set()
    stages = plan.get("stages") or []
    previous = set()
    for si, stage in enumerate(stages):
        requirement = stage["requirement"]
        if requirement["id"] in previous or not set(requirement.get("dependsOn") or []) <= previous:
            raise ValueError(f"Invalid workflow dependency order for {case['id']}: {requirement}")
        previous.add(requirement["id"])
        if mode == "faceted":
            proven, count = verify_facet_quotes(stage, repo)
            satisfied_facets += proven
            facet_count += count
        applicable = [r for r in case["requirements"] if
                      (requirement["id"] == r["id"] if mode in ("structured", "faceted") else
                       re.search(r["stage"], requirement["text"], re.I))]
        for c in stage.get("contributions") or []:
            credited = set()
            provenances = []
            for e in c["elements"]:
                text = selected_source(repo, c, e)
                section_count += 1
                element_checks = []
                for r in applicable:
                    offset = evidence(text, r["source"])
                    if offset is None or c.get("status") == "blocked":
                        continue
                    where = f"{c['skill']}/{e['file']}:{e['startLine'] + offset}"
                    results[r["id"]]["status"] = "supported"
                    if where not in results[r["id"]]["evidence"]:
                        results[r["id"]]["evidence"].append(where)
                    element_checks.append(r["id"])
                    credited.add(r["id"])
                    supported_stages.add(si)
                if element_checks:
                    credited_sections += 1
                provenances.append((e["file"], e["title"]))
            labels = c["skill"] + "/" + " ".join(f"{p} {t}" for p, t in provenances)
            contributions.append({"stage": requirement["id"], "skill": c["skill"], "checks": sorted(credited),
                                  "nearMiss": not credited and bool(re.search(case["nearMiss"]["candidate"], labels, re.I)),
                                  "blocked": c.get("status") == "blocked"})
        if stage.get("gaps"):
            for r in applicable:
                if r.get("expected") == "gap" and not stage.get("contributions"):
                    results[r["id"]]["status"] = "declared-gap"
    if not stages and plan.get("gaps"):
        for r in case["requirements"]:
            if r.get("expected") == "gap":
                results[r["id"]]["status"] = "declared-gap"
    supported = sum(v["status"] == "supported" for v in results.values())
    gap_count = sum(v["status"] == "declared-gap" for v in results.values())
    needed = sum(r.get("expected") != "gap" for r in case["requirements"])
    positive = sum(results[r["id"]]["status"] == "supported" for r in case["requirements"] if r.get("expected") != "gap")
    no_gap = all(result["status"] in ("supported", "declared-gap") for result in results.values())
    executable = bool(stages) and not plan.get("gaps") and not any(s.get("gaps") for s in stages)
    return {"id": case["id"], "family": case["family"], "mode": mode, "requirements": results,
            "contributions": contributions, "supported": supported, "positiveCovered": positive, "needed": needed,
            "expectedGaps": len(results) - needed, "declaredGaps": gap_count,
            "gapCorrect": sum(results[r['id']]['status'] in ('supported', 'declared-gap')
                              for r in case['requirements'] if r.get('expected') == 'gap'),
            "count": len(contributions), "credited": sum(bool(c["checks"]) for c in contributions),
            "nearMisses": sum(c["nearMiss"] for c in contributions), "elements": section_count,
            "creditedElements": credited_sections, "stages": len(stages), "groundedStages": len(supported_stages),
            "deep": len(case["requirements"]) >= 4, "facetsSatisfied": satisfied_facets, "facetsReported": facet_count,
            "executable": executable, "deepComplete": no_gap and executable and len(supported_stages) >= min(3, needed),
            "allCriteria": no_gap, "planGaps": plan.get("gaps") or []}


def totals(rows, mode):
    data = [r for r in rows if r["mode"] == mode]
    out = {k: sum(r[k] for r in data) for k in ("positiveCovered", "needed", "gapCorrect", "expectedGaps",
            "credited", "count", "nearMisses", "creditedElements", "elements", "allCriteria", "facetsSatisfied", "facetsReported")}
    out["cases"] = len(data)
    out["plansWithoutGaps"] = sum(r.get("executable", False) for r in data)
    singles = [r for r in data if not r["deep"] and r["needed"]]
    out["singleFamilies"] = len({r["family"] for r in singles})
    out["groundedFamilies"] = len({r["family"] for r in singles if r["positiveCovered"] > 0})
    out["deepComplete"] = sum(r["deepComplete"] for r in data if r["deep"])
    out["deepCases"] = sum(r["deep"] for r in data)
    return out


def evaluate(binary, repo, cache, cases, label, modes, annotations):
    with Router(binary, repo, cache) as router:
        rows = []
        for case in cases:
            for mode in modes:
                rows.append(score(case, mode, router.plan(request_for(case, mode, annotations)), router.repo))
            print(f"{label} {case['id']}: " + " | ".join(
                f"{r['mode']} {r['positiveCovered']}/{r['needed']} needs, {r['gapCorrect']}/{r['expectedGaps']} gaps, "
                f"{r['credited']}/{r['count']} grounded contributions" for r in rows[-len(modes):]), flush=True)
        return {"identity": router.identity(), "modes": {m: totals(rows, m) for m in modes}, "cases": rows}


def acceptance(result, historical, policy):
    failures = []
    for mode, s in result["current"]["modes"].items():
        thresholds = policy["modes"][mode]
        for metric, minimum in thresholds.items():
            value = s["credited"] / max(s["count"], 1) if metric == "groundedPrecision" else s[metric]
            if value < minimum:
                failures.append(f"{mode} {metric} {value:.3g} is below {minimum}")
        if mode in historical["modes"]:
            previous = historical["modes"][mode]
            for key in ("positiveCovered", "deepComplete", "groundedFamilies", "gapCorrect"):
                if s[key] < previous[key]:
                    failures.append(f"{mode} {key} regressed from {previous[key]} to {s[key]}")
            if s["credited"] * previous["count"] < previous["credited"] * s["count"]:
                failures.append(f"{mode} grounded contribution precision regressed")
    return failures


def main():
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument("--router", default=os.environ.get("SKILL_ROUTER_BIN", "skill-router"))
    parser.add_argument("--repo", type=Path, default=REPO)
    parser.add_argument("--baseline", help="Frozen reference binary supporting workflow_plan")
    parser.add_argument("--baseline-repo", type=Path, help="Corresponding source snapshot (default: same corpus)")
    parser.add_argument("--output", type=Path, help="Write full scores and supporting source locations")
    parser.add_argument("--mode", choices=("natural", "structured", "faceted"), action="append", dest="modes",
                        help="Limit current evaluation to one mode (repeatable)")
    parser.add_argument("--gate", action="store_true", help="Apply maintained breadth floors and frozen baseline nonregression")
    parser.add_argument("--case", action="append", dest="case_ids", help="Evaluate one case (repeatable)")
    parser.add_argument("--cache", type=Path, help="Reuse derived index between invocations")
    parser.add_argument("--strict", action="store_true", help="Require all registered criteria, deep coverage and zero named near misses")
    args = parser.parse_args()
    cases = json.loads(CASES.read_text())
    annotations = json.loads(FACETS.read_text())
    for case in cases:
        if {r["id"] for r in case["requirements"]} != set(annotations.get(case["id"], {})):
            raise ValueError(f"Incomplete host annotations for {case['id']}")
    modes = tuple(dict.fromkeys(args.modes or ("natural", "structured", "faceted")))
    if args.gate and (args.case_ids or len(modes) != 3):
        parser.error("the maintained release gate uses all cases and all three modes")
    if args.case_ids:
        unknown = set(args.case_ids) - {c["id"] for c in cases}
        if unknown:
            parser.error(f"unknown cases: {sorted(unknown)}")
        cases = [c for c in cases if c["id"] in args.case_ids]
    preregistration(cases, args.repo)
    result = {"suiteSHA256": sha(CASES.read_bytes()), "hostFacetsSHA256": sha(FACETS.read_bytes())}
    with tempfile.TemporaryDirectory(prefix="workflow-breadth-") as tmp:
        cache = args.cache or Path(tmp)
        baseline_modes = tuple(m for m in modes if m != "faceted")
        if args.baseline and baseline_modes:
            result["baseline"] = evaluate(args.baseline, args.baseline_repo or args.repo, cache / "baseline",
                                          cases, "baseline", baseline_modes, annotations)
        result["current"] = evaluate(args.router, args.repo, cache / "current", cases, "current", modes, annotations)
    if args.gate:
        historical = json.loads(BASELINE.read_text())
        if result["suiteSHA256"] != historical["suiteSHA256"]:
            raise ValueError("Frozen source criteria changed; re-establish the baseline explicitly.")
        result["gateFailures"] = acceptance(result, historical, json.loads(POLICY.read_text()))
    if args.output:
        args.output.write_text(json.dumps(result, indent=2, ensure_ascii=False) + "\n")
    for which in ("baseline", "current"):
        for mode, values in result.get(which, {}).get("modes", {}).items():
            print(f"{which} {mode}: {json.dumps(values, sort_keys=True)}")
    if result.get("gateFailures"):
        sys.exit("Workflow breadth gate failed:\n" + "\n".join(result["gateFailures"]))
    current = result["current"]["modes"].values()
    if args.strict and any(s["allCriteria"] != s["cases"] or s["deepComplete"] != s["deepCases"]
                           or s["nearMisses"] for s in current):
        sys.exit("Workflow breadth acceptance failed.")


if __name__ == "__main__":
    main()
