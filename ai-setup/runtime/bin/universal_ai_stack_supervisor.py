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


CREATE_NO_WINDOW = 0x08000000 if os.name == "nt" else 0
ROOT = Path(os.environ.get("UNIVERSAL_AI_STACK_HOME", Path(__file__).resolve().parents[1]))
CONFIG_DIR = ROOT / "config"
LOG_DIR = ROOT / "logs"
STATE_DIR = ROOT / "state"
SECRETS_ENV = ROOT / "secrets" / ".env"


def load_json(path: Path) -> dict[str, Any]:
    with path.open("r", encoding="utf-8") as f:
        return json.load(f)


def load_env(path: Path) -> dict[str, str]:
    env = dict(os.environ)
    env["UNIVERSAL_AI_STACK_HOME"] = str(ROOT)
    if path.exists():
        for raw in path.read_text(encoding="utf-8", errors="replace").splitlines():
            line = raw.strip()
            if not line or line.startswith("#") or "=" not in line:
                continue
            key, value = line.split("=", 1)
            if key.strip():
                env[key.strip()] = value.strip().strip('"')
    return env


def setup_logging() -> None:
    LOG_DIR.mkdir(parents=True, exist_ok=True)
    logging.basicConfig(
        filename=str(LOG_DIR / "supervisor.log"),
        level=logging.INFO,
        format="%(asctime)s %(levelname)s %(message)s",
    )


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
    env = load_env(SECRETS_ENV)
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
    setup_logging()
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
