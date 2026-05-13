#!/usr/bin/env python3
"""Run the locally integrated OneFileLLM CLI from its isolated venv."""

from __future__ import annotations

import os
import subprocess
import sys
from pathlib import Path


def default_cli() -> Path:
    override = os.environ.get("ONEFILELLM_CLI")
    if override:
        return Path(override)
    if os.name == "nt":
        return Path.home() / ".onefilellm" / "venv" / "Scripts" / "onefilellm.exe"
    return Path.home() / ".onefilellm" / "venv" / "bin" / "onefilellm"


def main() -> int:
    cli = default_cli()
    if not cli.exists():
        print(f"OneFileLLM CLI not found: {cli}", file=sys.stderr)
        return 127

    env = os.environ.copy()
    env.setdefault("PYTHONUTF8", "1")
    return subprocess.call([str(cli), *sys.argv[1:]], env=env)


if __name__ == "__main__":
    raise SystemExit(main())
