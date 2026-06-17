#!/usr/bin/env python3
"""Hidden low-resource supervisor for the local Universal AI Stack."""

from __future__ import annotations

import json
import logging
import os
import subprocess
import sys
import time
import urllib.request
from pathlib import Path
from typing import Any

# Shared runtime helpers live beside this script in runtime/bin (a plain scripts
# directory, not a package). Make that directory importable before importing them
# so resolution does not depend on the working directory or interpreter flags.
sys.path.insert(0, str(Path(__file__).resolve().parent))
from _universal_ai_common import (  # noqa: E402  (path bootstrap must precede import)
    CONFIG_DIR,
    ROOT,
    SECRETS_ENV,
    load_env,
    load_json,
    setup_logging,
)


CREATE_NO_WINDOW = 0x08000000 if os.name == "nt" else 0
STATE_DIR = ROOT / "state"


def healthy(url: str, timeout: float = 8.0) -> bool:
    try:
        with urllib.request.urlopen(url, timeout=timeout) as response:
            return 200 <= response.status < 300
    except Exception:
        return False


def pythonw_path() -> str:
    candidate = Path(sys.executable)
    if candidate.name.lower() == "pythonw.exe":
        return str(candidate)
    sibling = candidate.with_name("pythonw.exe")
    if sibling.exists():
        return str(sibling)
    return str(candidate)


def normalize_command(command: list[str]) -> list[str]:
    if not command:
        return command
    if command[0].lower() == "pythonw":
        return [pythonw_path(), *command[1:]]
    return command


def start_service(service: dict[str, Any], env: dict[str, str]) -> None:
    command = normalize_command(service.get("start", []))
    if not command:
        logging.warning("service %s has no start command", service.get("id"))
        return
    for idx, part in enumerate(command):
        if idx == 0:
            continue
        if isinstance(part, str):
            command[idx] = part.replace("%UNIVERSAL_AI_STACK_HOME%", str(ROOT))
    logging.info("starting service=%s", service.get("id"))
    subprocess.Popen(
        command,
        cwd=str(ROOT),
        env=env,
        stdin=subprocess.DEVNULL,
        stdout=subprocess.DEVNULL,
        stderr=subprocess.DEVNULL,
        creationflags=CREATE_NO_WINDOW,
        close_fds=False,
    )


def wait_for_health(url: str | None, timeout_seconds: float = 45.0, interval_seconds: float = 2.0) -> bool:
    if not url:
        return True
    deadline = time.monotonic() + timeout_seconds
    while time.monotonic() < deadline:
        if healthy(url):
            return True
        time.sleep(interval_seconds)
    return healthy(url)


def write_state(results: list[dict[str, Any]]) -> None:
    STATE_DIR.mkdir(parents=True, exist_ok=True)
    payload = {"time": int(time.time()), "results": results}
    (STATE_DIR / "last-supervisor-check.json").write_text(json.dumps(payload, indent=2), encoding="utf-8")


def check_once() -> list[dict[str, Any]]:
    integrations = load_json(CONFIG_DIR / "integrations.json")
    # Supervisor semantics: inject the stack home and let the dotenv file override
    # the process environment so spawned services see the configured values.
    env = load_env(SECRETS_ENV, inject_home=True, override=True)
    results: list[dict[str, Any]] = []
    for service in integrations.get("services", []):
        health_url = service.get("healthUrl")
        service_id = service.get("id")
        is_healthy = healthy(health_url) if health_url else False
        if not is_healthy:
            start_service(service, env)
            is_healthy = wait_for_health(health_url)
        results.append({"id": service_id, "healthy": is_healthy, "healthUrl": health_url})
    write_state(results)
    return results


def main() -> int:
    setup_logging("supervisor.log")
    policy = load_json(CONFIG_DIR / "routing-policy.json")
    interval = int(policy.get("supervisor", {}).get("checkIntervalSeconds", 600))
    once = "--once" in sys.argv
    logging.info("supervisor started once=%s interval=%s", once, interval)
    while True:
        try:
            results = check_once()
            logging.info("supervisor check results=%s", results)
        except Exception as exc:  # noqa: BLE001
            logging.exception("supervisor check failed: %s", exc)
        if once:
            return 0
        time.sleep(interval)


if __name__ == "__main__":
    sys.exit(main())
