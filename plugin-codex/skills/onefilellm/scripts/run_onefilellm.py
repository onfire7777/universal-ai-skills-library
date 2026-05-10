#!/usr/bin/env python3
"""Run the locally integrated OneFileLLM CLI from its isolated venv."""

from __future__ import annotations

import os
import subprocess
import sys
from pathlib import Path


CLI = Path(r"C:\Users\burni\.onefilellm\venv\Scripts\onefilellm.exe")


def main() -> int:
    if not CLI.exists():
        print(f"OneFileLLM CLI not found: {CLI}", file=sys.stderr)
        return 127

    env = os.environ.copy()
    env.setdefault("PYTHONUTF8", "1")
    return subprocess.call([str(CLI), *sys.argv[1:]], env=env)


if __name__ == "__main__":
    raise SystemExit(main())
