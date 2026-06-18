"""Shared harness for skill-router + registry characterization tests.

These tests pin the *current, pre-refactor* behaviour of the Universal AI
Skills Library so the in-place consolidation (router decoupling, manus merge,
registry unification) can be proven not to regress existing behaviour.

Design goals
------------
* **Black-box.** The router is exercised as a built CLI binary, not via Go
  internals, so the suite keeps working after the router source is moved into
  a different package/directory.
* **Refactor-robust discovery.** The Go module is located by content, not by a
  hard-coded path, and an explicit ``SKILL_ROUTER_BIN`` override short-circuits
  the build entirely (used by CI to build once).
* **Deterministic.** Only commands documented to avoid network / model APIs are
  used (``skills validate-manifest --json`` and ``preflight --json``).
* **Zero third-party dependencies.** Standard-library ``unittest`` only.

This module is imported by both the test cases and ``update_baseline.py`` so the
prompt battery and discovery logic have a single source of truth.
"""

from __future__ import annotations

import functools
import hashlib
import json
import os
import shutil
import subprocess
import sys
import tempfile
from typing import Dict, List, Optional, Tuple

# --- Fixture --------------------------------------------------------------
# Characterization tests run the router against a *frozen* fixture skills
# library, NOT the live ``skills/`` tree, so that Goal 2/3 consolidation edits
# (manus merge, registry unification) cannot break the routing/lookup/alias
# pins. The router is pointed at the fixture via ``SKILL_ROUTER_REPO_DIR`` (the
# same hook the Go ``configurePreflightTest`` helper uses), and HOME/config are
# redirected to throwaway dirs so no live *external* skill roots leak in.
#
# The fixture is the SHARED, single source of truth co-authored/verified with
# the Go routing tests (Builder 2 owns its content; only ``go test ./cmd/skills``
# can verify the scoring-sensitive routes). Resolution order:
#   1. ``CHAR_FIXTURE_DIR`` env override (explicit).
#   2. the Go-conventional shared location ``skill-router-cli/testdata/route-fixture``.
#   3. a local ``tests/fixtures/skills-lib`` fallback (transition / standalone).
def _resolve_fixture_dir() -> str:
    override = os.environ.get("CHAR_FIXTURE_DIR")
    if override:
        return os.path.abspath(override)
    here = os.path.dirname(os.path.abspath(__file__))
    candidates = [
        # Shared fixture authored + routing-verified by Builder 2's Go tests.
        os.path.join(here, "..", "..", "skill-router-cli", "cmd", "skills", "testdata", "route-fixture"),
        os.path.join(here, "..", "fixtures", "skills-lib"),
    ]
    for c in candidates:
        c = os.path.abspath(c)
        if os.path.isfile(os.path.join(c, "manifest.json")):
            return c
    # Default to the shared location even if absent, so errors point there.
    return os.path.abspath(candidates[0])


FIXTURE_DIR = _resolve_fixture_dir()

# --- Routing battery -------------------------------------------------------
# Each prompt pins a routing *outcome* (decision + chosen skill). Scores are
# recorded in the golden for diagnostics but intentionally NOT asserted: scoring
# internals may legitimately drift during the refactor, whereas the selected
# skill is the behavioural contract that must be preserved. All prompts resolve
# deterministically within the fixture (no external skill roots involved).
PROMPTS: List[Dict[str, object]] = [
    {
        "id": "card-creator-mothers-day",
        "prompt": "use the universal AI skills card creator skill to create a beautiful mothers day card",
    },
    {
        "id": "chat-summarizer-handoff",
        "prompt": "summarize this chat session into a handoff document",
    },
    {
        "id": "router-maintenance-meta",
        "prompt": "improve the skill router automatic routing accuracy",
    },
    {
        "id": "python-testing-patterns",
        "prompt": "write pytest fixtures and mocking tests for a python module",
    },
    {
        "id": "universal-ai-setup",
        "prompt": (
            "please make sure the universal AI tools is cleanly and universally installed "
            "not redundant and clean everything updated on GitHub and in all the different "
            "AI services on my computer"
        ),
    },
    {
        "id": "generic-no-route",
        "prompt": "what is the capital of France",
    },
]


class RouterUnavailable(Exception):
    """Raised when the Go toolchain is missing so tests can SKIP, not fail."""


@functools.lru_cache(maxsize=1)
def repo_root() -> str:
    """Return the canonical repository root.

    Resolved from this file's location (``<repo>/tests/characterization``) and
    validated by the presence of ``manifest.json``. Falls back to
    ``git rev-parse --show-toplevel`` so the suite is portable if relocated.
    """
    here = os.path.dirname(os.path.abspath(__file__))
    candidate = os.path.dirname(os.path.dirname(here))
    if os.path.exists(os.path.join(candidate, "manifest.json")):
        return candidate
    try:
        top = subprocess.run(
            ["git", "rev-parse", "--show-toplevel"],
            cwd=here, capture_output=True, text=True, timeout=30,
        ).stdout.strip()
        if top and os.path.exists(os.path.join(top, "manifest.json")):
            return top
    except (OSError, subprocess.SubprocessError):
        pass
    return candidate


def find_router_module(repo: str) -> Optional[str]:
    """Locate the skill-router Go module directory (content-based, move-safe)."""
    preferred = ["skill-router-cli", os.path.join("packages", "skill-router")]
    for rel in preferred:
        d = os.path.join(repo, rel)
        if os.path.isfile(os.path.join(d, "go.mod")) and os.path.isfile(os.path.join(d, "main.go")):
            return d
    for root, dirs, files in os.walk(repo):
        # Prune noise; never descend into VCS or dependency trees.
        dirs[:] = [d for d in dirs if d not in (".git", "node_modules", "vendor")]
        if "go.mod" in files and "main.go" in files:
            try:
                if "skill-router" in open(os.path.join(root, "go.mod"), encoding="utf-8").read():
                    return root
            except OSError:
                continue
    return None


@functools.lru_cache(maxsize=1)
def router_binary() -> str:
    """Return a path to a runnable skill-router binary, building once if needed.

    Honours ``SKILL_ROUTER_BIN`` (CI builds once and exports it). Otherwise
    builds the discovered Go module to a temp dir *outside* the repo so no build
    artefact lands in another agent's owned directory.

    Raises :class:`RouterUnavailable` when the Go toolchain is absent so callers
    can skip rather than fail.
    """
    override = os.environ.get("SKILL_ROUTER_BIN")
    if override and os.path.isfile(override) and os.access(override, os.X_OK):
        return override

    if shutil.which("go") is None:
        raise RouterUnavailable("go toolchain not available; cannot build skill-router")

    repo = repo_root()
    module = find_router_module(repo)
    if not module:
        raise RouterUnavailable("could not locate the skill-router Go module")

    out_dir = tempfile.mkdtemp(prefix="skill-router-char-")
    out_bin = os.path.join(out_dir, "skill-router")
    proc = subprocess.run(
        ["go", "build", "-o", out_bin, "."],
        cwd=module, capture_output=True, text=True, timeout=600,
        env={**os.environ, "GOFLAGS": "-mod=readonly"},
    )
    if proc.returncode != 0:
        # A real compile failure is a genuine problem — surface it.
        raise AssertionError(f"skill-router build failed:\n{proc.stderr or proc.stdout}")
    return out_bin


@functools.lru_cache(maxsize=1)
def fixture_env() -> Dict[str, str]:
    """Environment that pins the router to the frozen fixture in isolation.

    * ``SKILL_ROUTER_REPO_DIR`` -> the fixture (manifest.json + skills/).
    * ``HOME`` / ``USERPROFILE`` / config dir -> throwaway temp dirs so the
      router never discovers live external skill roots, keeping routing
      deterministic and independent of the host.
    """
    if not os.path.isdir(FIXTURE_DIR):
        raise AssertionError(f"characterization fixture missing: {FIXTURE_DIR}")
    iso = tempfile.mkdtemp(prefix="skill-router-fixture-home-")
    cfg = os.path.join(iso, "cfg")
    os.makedirs(cfg, exist_ok=True)
    # Mirrors Builder 2's hermetic Go-test setup + Scout 1's golden recipe:
    # pin repo + skills dir to the fixture, isolate HOME/config, and disable
    # colour so text output is stable (the router uses fatih/color).
    return {
        "SKILL_ROUTER_REPO_DIR": FIXTURE_DIR,
        "SKILL_ROUTER_SKILLS_DIR": os.path.join(FIXTURE_DIR, "skills"),
        "SKILL_ROUTER_CONFIG_DIR": cfg,
        "HOME": iso,
        "USERPROFILE": iso,
        "NO_COLOR": "1",
        "CLICOLOR": "0",
    }


def run_router(
    args: List[str], timeout: int = 120, env: Optional[Dict[str, str]] = None,
    cwd: Optional[str] = None,
) -> subprocess.CompletedProcess:
    """Run the router. ``env`` overlays the process environment; ``cwd`` defaults
    to the repo root (the router's repo detection is CWD-based)."""
    full_env = {**os.environ, **(env or {})}
    return subprocess.run(
        [router_binary(), *args],
        cwd=cwd or repo_root(), capture_output=True, text=True, timeout=timeout,
        env=full_env,
    )


def preflight(prompt: str, fixture: bool = True) -> Dict[str, object]:
    """Return parsed ``preflight --json`` output for a prompt.

    By default the router is pinned to the frozen fixture (``fixture=True``).
    """
    env = fixture_env() if fixture else None
    cwd = FIXTURE_DIR if fixture else None
    proc = run_router(["preflight", "--json", prompt], env=env, cwd=cwd)
    if not proc.stdout.strip():
        raise AssertionError(f"empty preflight output for {prompt!r}: {proc.stderr}")
    return json.loads(proc.stdout)


def validate_manifest(fixture: bool = False) -> Dict[str, object]:
    """Return parsed ``skills validate-manifest --json`` output.

    Defaults to the *live* tree (``fixture=False``) because the live registry's
    structural integrity is itself a characterization target; pass
    ``fixture=True`` to validate the frozen fixture.
    """
    env = fixture_env() if fixture else None
    cwd = FIXTURE_DIR if fixture else None
    proc = run_router(["skills", "validate-manifest", "--json"], env=env, cwd=cwd)
    if not proc.stdout.strip():
        raise AssertionError(f"empty validate-manifest output: {proc.stderr}")
    return json.loads(proc.stdout)


def run_router_as_alias(
    alias: str, args: List[str], timeout: int = 120, fixture: bool = True,
) -> subprocess.CompletedProcess:
    """Run the same router binary under a different argv[0].

    This supports compatibility characterization for explicitly tested command
    names without advertising historical project-specific binaries as an active
    install contract. Pinned to the frozen fixture by default.
    """
    src = router_binary()
    alias_dir = tempfile.mkdtemp(prefix=f"skill-router-alias-{alias}-")
    alias_path = os.path.join(alias_dir, alias)
    shutil.copy2(src, alias_path)
    os.chmod(alias_path, 0o755)
    env = {**os.environ, **(fixture_env() if fixture else {})}
    cwd = FIXTURE_DIR if fixture else repo_root()
    return subprocess.run(
        [alias_path, *args],
        cwd=cwd, capture_output=True, text=True, timeout=timeout, env=env,
    )


def go_module() -> Optional[str]:
    """Convenience wrapper returning the skill-router Go module dir, if any."""
    return find_router_module(repo_root())


def run_go_unit_tests(timeout: int = 600) -> Tuple[set, set]:
    """Run ``go test ./...`` for the router module and return failing names.

    Returns ``(failing_tests, failing_packages)`` where each failing test is
    formatted ``"<package> :: <Test>"``. Raises :class:`RouterUnavailable` if the
    Go toolchain or module is missing so the caller can skip.
    """
    if shutil.which("go") is None:
        raise RouterUnavailable("go toolchain not available")
    module = go_module()
    if not module:
        raise RouterUnavailable("could not locate the skill-router Go module")
    proc = subprocess.run(
        ["go", "test", "./...", "-json"],
        cwd=module, capture_output=True, text=True, timeout=timeout,
        env={**os.environ, "GOFLAGS": "-mod=readonly"},
    )
    failing_tests: set = set()
    failing_pkgs: set = set()
    for line in proc.stdout.splitlines():
        line = line.strip()
        if not line:
            continue
        try:
            evt = json.loads(line)
        except json.JSONDecodeError:
            continue
        if evt.get("Action") == "fail":
            pkg = evt.get("Package", "")
            test = evt.get("Test")
            if test:
                failing_tests.add(f"{pkg} :: {test}")
            else:
                failing_pkgs.add(pkg)
    return failing_tests, failing_pkgs


def load_manifest() -> Dict[str, object]:
    """Load the raw ``manifest.json`` registry document."""
    with open(os.path.join(repo_root(), "manifest.json"), encoding="utf-8") as fh:
        return json.load(fh)


def load_fixture_manifest() -> Dict[str, object]:
    """Load the frozen fixture's ``manifest.json``."""
    with open(os.path.join(FIXTURE_DIR, "manifest.json"), encoding="utf-8") as fh:
        return json.load(fh)


def fixture_skill_count() -> int:
    """Total skills declared by the fixture manifest (core + library)."""
    m = load_fixture_manifest()
    return len(m.get("core_skills") or []) + len(m.get("library_skills") or [])


def manifest_skill_pairs(manifest: Dict[str, object]) -> List[Tuple[str, str]]:
    """Return sorted ``(name, directory)`` pairs across core + library skills."""
    skills = list(manifest.get("core_skills") or []) + list(manifest.get("library_skills") or [])
    return sorted((s["name"], s["directory"]) for s in skills)


def corpus_fingerprint(pairs: List[Tuple[str, str]]) -> str:
    """Stable SHA-256 over the sorted ``name|directory`` corpus."""
    blob = "\n".join(f"{name}|{directory}" for name, directory in pairs)
    return hashlib.sha256(blob.encode("utf-8")).hexdigest()


BASELINE_DIR = os.path.join(os.path.dirname(os.path.abspath(__file__)), "baseline")
REGISTRY_BASELINE = os.path.join(BASELINE_DIR, "registry_baseline.json")
ROUTING_GOLDEN = os.path.join(BASELINE_DIR, "routing_golden.json")


def load_json(path: str) -> Dict[str, object]:
    with open(path, encoding="utf-8") as fh:
        return json.load(fh)
