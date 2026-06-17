#!/usr/bin/env python3
"""Shared runtime helpers for the Universal AI Stack bin/ entrypoints.

Single source of truth for stack path resolution, config/dotenv loading, and
logging configuration used by ``universal_ai_router.py`` and
``universal_ai_stack_supervisor.py``.

This module was extracted to remove duplicated logic (DRY) **without changing
behaviour**: the entrypoints opt into their historical variations via the
``inject_home`` / ``override`` flags on :func:`load_env` and the ``log_file``
argument to :func:`setup_logging`.

It lives alongside the entrypoints in ``runtime/bin`` (a plain scripts directory,
not a package). The entrypoints add this directory to ``sys.path`` before
importing it, so it resolves regardless of the working directory or how the
interpreter was launched.
"""
from __future__ import annotations

import json
import logging
import os
from pathlib import Path
from typing import Any

# Stack home: honour an explicit override, otherwise the ``runtime`` directory
# (the parent of this ``bin`` directory). Computed identically to the legacy
# per-entrypoint definitions so installed/launched behaviour is unchanged.
ROOT = Path(os.environ.get("UNIVERSAL_AI_STACK_HOME", Path(__file__).resolve().parents[1]))
CONFIG_DIR = ROOT / "config"
LOG_DIR = ROOT / "logs"
SECRETS_ENV = ROOT / "secrets" / ".env"


def load_json(path: Path) -> dict[str, Any]:
    """Read and parse a UTF-8 JSON file into a dictionary."""
    with path.open("r", encoding="utf-8") as f:
        return json.load(f)


def load_env(path: Path, *, inject_home: bool = False, override: bool = False) -> dict[str, str]:
    """Build an environment mapping from ``os.environ`` plus a dotenv-style file.

    Lines are ``KEY=value``; blank lines, ``#`` comments, and lines without an
    ``=`` are ignored, and surrounding double quotes are stripped from values.

    Args:
        path: dotenv-style file to read. A missing file is not an error.
        inject_home: when True, set ``UNIVERSAL_AI_STACK_HOME`` to :data:`ROOT`
            before applying the file (matches the supervisor's historical
            behaviour so spawned services inherit the stack home).
        override: when True, file values overwrite existing environment values
            (supervisor semantics); when False, existing environment values take
            precedence and the file only fills in keys that are not already set
            (router semantics).

    Returns:
        A new mapping; ``os.environ`` itself is never mutated.
    """
    env = dict(os.environ)
    if inject_home:
        env["UNIVERSAL_AI_STACK_HOME"] = str(ROOT)
    if not path.exists():
        return env
    for raw in path.read_text(encoding="utf-8", errors="replace").splitlines():
        line = raw.strip()
        if not line or line.startswith("#") or "=" not in line:
            continue
        key, value = line.split("=", 1)
        key = key.strip()
        value = value.strip().strip('"')
        if not key:
            continue
        if override or key not in env:
            env[key] = value
    return env


def setup_logging(log_file: str) -> None:
    """Configure root logging to ``LOG_DIR/<log_file>`` (INFO, timestamped).

    Args:
        log_file: file name (e.g. ``"router.log"``) created under :data:`LOG_DIR`.
    """
    LOG_DIR.mkdir(parents=True, exist_ok=True)
    logging.basicConfig(
        filename=str(LOG_DIR / log_file),
        level=logging.INFO,
        format="%(asctime)s %(levelname)s %(message)s",
    )
