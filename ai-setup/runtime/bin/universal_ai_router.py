#!/usr/bin/env python3
"""Small OpenAI-compatible failover router for HTTP-capable local AI clients."""

from __future__ import annotations

import json
import logging
import os
import sys
import threading
import time
import urllib.error
import urllib.request
from http.server import BaseHTTPRequestHandler, ThreadingHTTPServer
from pathlib import Path
from typing import Any


ROOT = Path(os.environ.get("UNIVERSAL_AI_STACK_HOME", Path(__file__).resolve().parents[1]))
CONFIG_DIR = ROOT / "config"
LOG_DIR = ROOT / "logs"
SECRETS_ENV = ROOT / "secrets" / ".env"


def load_json(path: Path) -> dict[str, Any]:
    with path.open("r", encoding="utf-8") as f:
        return json.load(f)


def load_env(path: Path) -> dict[str, str]:
    env = dict(os.environ)
    if not path.exists():
        return env
    for raw in path.read_text(encoding="utf-8", errors="replace").splitlines():
        line = raw.strip()
        if not line or line.startswith("#") or "=" not in line:
            continue
        key, value = line.split("=", 1)
        key = key.strip()
        value = value.strip().strip('"')
        if key and key not in env:
            env[key] = value
    return env


def setup_logging() -> None:
    LOG_DIR.mkdir(parents=True, exist_ok=True)
    logging.basicConfig(
        filename=str(LOG_DIR / "router.log"),
        level=logging.INFO,
        format="%(asctime)s %(levelname)s %(message)s",
    )


setup_logging()
REGISTRY = load_json(CONFIG_DIR / "model-registry.json")
POLICY = load_json(CONFIG_DIR / "routing-policy.json")
ENV = load_env(SECRETS_ENV)
CIRCUIT_LOCK = threading.Lock()
CIRCUIT_STATE: dict[str, dict[str, Any]] = {}


def registry_models() -> list[dict[str, Any]]:
    return REGISTRY.get("models", [])


def retry_policy() -> dict[str, Any]:
    return POLICY.get("retry", {})


def circuit_config() -> dict[str, Any]:
    return retry_policy().get("circuitBreaker", {})


def circuit_enabled() -> bool:
    return bool(circuit_config().get("enabled", False))


def circuit_snapshot() -> dict[str, Any]:
    now = time.time()
    with CIRCUIT_LOCK:
        return {
            provider: {
                **state,
                "open": float(state.get("openedUntil", 0)) > now,
                "remainingSeconds": max(0, int(float(state.get("openedUntil", 0)) - now)),
            }
            for provider, state in CIRCUIT_STATE.items()
        }


def circuit_is_open(provider_id: str) -> bool:
    if not circuit_enabled():
        return False
    now = time.time()
    with CIRCUIT_LOCK:
        state = CIRCUIT_STATE.get(provider_id)
        if not state:
            return False
        opened_until = float(state.get("openedUntil", 0))
        if opened_until > now:
            return True
        if opened_until:
            state["openedUntil"] = 0
        return False


def record_provider_success(provider_id: str) -> None:
    if not circuit_enabled():
        return
    with CIRCUIT_LOCK:
        CIRCUIT_STATE.pop(provider_id, None)


def record_provider_failure(provider_id: str) -> None:
    if not circuit_enabled():
        return
    cfg = circuit_config()
    threshold = int(cfg.get("failureThreshold", 2))
    cooldown = int(cfg.get("cooldownSeconds", 600))
    now = time.time()
    with CIRCUIT_LOCK:
        state = CIRCUIT_STATE.setdefault(provider_id, {"failures": 0, "openedUntil": 0})
        state["failures"] = int(state.get("failures", 0)) + 1
        state["lastFailureAt"] = int(now)
        if state["failures"] >= threshold:
            state["openedUntil"] = now + cooldown
            logging.warning("provider circuit opened provider=%s cooldown=%ss", provider_id, cooldown)


def model_by_id(model_id: str) -> dict[str, Any] | None:
    aliases = REGISTRY.get("canonicalAliases", {})
    actual = aliases.get(model_id, model_id)
    for model in registry_models():
        if model.get("id") == actual or model.get("model") == actual:
            return model
    return None


def routeable_models() -> list[dict[str, Any]]:
    ordered = []
    for model_id in POLICY.get("fallbackOrder", []):
        model = model_by_id(model_id)
        if not model or not model.get("enabled", False):
            continue
        if model.get("routeKind") != "openai-compatible-http":
            continue
        ordered.append(model)
    return ordered


def response_json(handler: BaseHTTPRequestHandler, status: int, payload: dict[str, Any]) -> None:
    data = json.dumps(payload, ensure_ascii=False).encode("utf-8")
    handler.send_response(status)
    handler.send_header("Content-Type", "application/json")
    handler.send_header("Content-Length", str(len(data)))
    handler.send_header("Access-Control-Allow-Origin", "http://127.0.0.1")
    handler.send_header("Access-Control-Allow-Headers", "authorization,content-type,x-api-key")
    handler.end_headers()
    handler.wfile.write(data)


def request_body(handler: BaseHTTPRequestHandler) -> dict[str, Any]:
    length = int(handler.headers.get("Content-Length", "0") or "0")
    raw = handler.rfile.read(length) if length else b"{}"
    return json.loads(raw.decode("utf-8") or "{}")


def authorized(handler: BaseHTTPRequestHandler) -> bool:
    router_cfg = POLICY.get("httpRouter", {})
    if not router_cfg.get("requireApiKey", True):
        return True
    expected = ENV.get(router_cfg.get("apiKeyEnv", "UNIVERSAL_AI_STACK_API_KEY")) or ENV.get("API_SERVER_KEY")
    if not expected:
        return True
    auth = handler.headers.get("Authorization", "")
    token = auth.removeprefix("Bearer ").strip() if auth.lower().startswith("bearer ") else ""
    api_key = handler.headers.get("x-api-key", "").strip()
    return token == expected or api_key == expected


def provider_health(model: dict[str, Any]) -> dict[str, Any]:
    item = {
        "id": model.get("id"),
        "provider": model.get("provider"),
        "routeKind": model.get("routeKind"),
        "enabled": model.get("enabled", False),
        "routeable": model.get("routeKind") == "openai-compatible-http",
        "available": False,
    }
    if not model.get("enabled", False):
        item["available"] = False
        item["status"] = "disabled"
        return item
    if model.get("apiKeyEnv"):
        item["apiKeyPresent"] = bool(ENV.get(model["apiKeyEnv"]))
        item["available"] = bool(ENV.get(model["apiKeyEnv"]))
    health_url = model.get("healthUrl")
    if health_url:
        try:
            with urllib.request.urlopen(health_url, timeout=2) as response:
                item["available"] = 200 <= response.status < 300
                item["status"] = response.status
        except Exception as exc:  # noqa: BLE001 - diagnostics only
            item["available"] = False
            item["status"] = type(exc).__name__
    if model.get("routeKind") == "host-cli-session":
        item["available"] = True
        item["note"] = "Available only through host CLI/session integrations, not this HTTP router."
    return item


def openai_error(message: str, code: str, status: int = 503, details: list[dict[str, Any]] | None = None) -> dict[str, Any]:
    return {
        "error": {
            "message": message,
            "type": "universal_ai_stack_error",
            "code": code,
            "details": details or [],
            "status": status,
        }
    }


def normalize_provider_payload(model: dict[str, Any], payload: dict[str, Any]) -> dict[str, Any]:
    provider_payload = dict(payload)
    provider_payload["model"] = model.get("model", model.get("id"))

    # Kimi/Moonshot thinking/code models reject the usual low-temperature
    # coding defaults. Normalize here so clients can use one shared config.
    if model.get("provider") == "moonshot-ai":
        provider_payload["temperature"] = 1
        provider_payload["top_p"] = 0.95
        if int(provider_payload.get("max_tokens") or 0) < 256:
            provider_payload["max_tokens"] = 256
        if int(provider_payload.get("max_completion_tokens") or 0) < 256:
            provider_payload.pop("max_completion_tokens", None)
        for unsupported in ("top_k", "repeat_penalty", "repetition_penalty"):
            provider_payload.pop(unsupported, None)

    return provider_payload


def call_provider(model: dict[str, Any], payload: dict[str, Any], timeout_override: float | None = None) -> tuple[int, bytes, str]:
    provider_payload = normalize_provider_payload(model, payload)
    base_url = model.get("baseUrl", "").rstrip("/")
    url = f"{base_url}/chat/completions"
    headers = {"Content-Type": "application/json"}
    api_key_env = model.get("apiKeyEnv")
    if api_key_env:
        api_key = ENV.get(api_key_env)
        if not api_key:
            raise RuntimeError(f"{api_key_env} is not configured")
        headers["Authorization"] = f"Bearer {api_key}"
    data = json.dumps(provider_payload).encode("utf-8")
    request = urllib.request.Request(url, data=data, headers=headers, method="POST")
    timeout = int(retry_policy().get("timeoutSeconds", 180))
    if timeout_override is not None:
        timeout = max(1, min(timeout, int(timeout_override)))
    with urllib.request.urlopen(request, timeout=timeout) as response:
        return response.status, response.read(), response.headers.get("Content-Type", "application/json")


def choose_candidates(requested_model: str | None) -> list[dict[str, Any]]:
    if not requested_model or requested_model in {"auto", "auto-coding", "primary-coding"}:
        return routeable_models()
    if requested_model == "primary-api":
        model = model_by_id("kimi-k2.6-thinking")
        return [model] if model else []
    model = model_by_id(requested_model)
    if model and model.get("routeKind") == "openai-compatible-http":
        return [model]
    return routeable_models()


class RouterHandler(BaseHTTPRequestHandler):
    server_version = "UniversalAIStackRouter/1.0"

    def log_message(self, fmt: str, *args: Any) -> None:
        logging.info("%s - %s", self.address_string(), fmt % args)

    def do_OPTIONS(self) -> None:  # noqa: N802
        self.send_response(204)
        self.send_header("Access-Control-Allow-Origin", "http://127.0.0.1")
        self.send_header("Access-Control-Allow-Headers", "authorization,content-type,x-api-key")
        self.send_header("Access-Control-Allow-Methods", "GET,POST,OPTIONS")
        self.end_headers()

    def do_GET(self) -> None:  # noqa: N802
        if self.path.rstrip("/") == "/health":
            response_json(
                self,
                200,
                {
                    "status": "ok",
                    "service": "universal-ai-stack-router",
                    "time": int(time.time()),
                    "models": [provider_health(m) for m in registry_models()],
                    "circuitBreakers": circuit_snapshot(),
                },
            )
            return
        if self.path.rstrip("/") == "/v1/models":
            if not authorized(self):
                response_json(self, 401, openai_error("Invalid API key", "invalid_api_key", 401))
                return
            data = [{"id": "auto-coding", "object": "model", "owned_by": "universal-ai-stack"}]
            for model in routeable_models():
                data.append({"id": model["id"], "object": "model", "owned_by": model.get("provider", "local")})
            response_json(self, 200, {"object": "list", "data": data})
            return
        response_json(self, 404, openai_error("Not found", "not_found", 404))

    def do_POST(self) -> None:  # noqa: N802
        if self.path.rstrip("/") not in {"/v1/chat/completions", "/chat/completions"}:
            response_json(self, 404, openai_error("Not found", "not_found", 404))
            return
        if not authorized(self):
            response_json(self, 401, openai_error("Invalid API key", "invalid_api_key", 401))
            return
        try:
            payload = request_body(self)
        except Exception as exc:  # noqa: BLE001
            response_json(self, 400, openai_error(f"Invalid JSON: {exc}", "invalid_json", 400))
            return
        candidates = choose_candidates(payload.get("model"))
        failures: list[dict[str, Any]] = []
        retry_cfg = retry_policy()
        fallback_status = set(retry_cfg.get("fallbackHttpStatus", []))
        max_attempts = int(retry_cfg.get("totalProviderAttempts", len(candidates) or 1))
        global_timeout = int(retry_cfg.get("globalTimeoutSeconds", retry_cfg.get("timeoutSeconds", 180)))
        deadline = time.monotonic() + max(1, global_timeout)
        attempts = 0
        for model in candidates:
            provider_id = model.get("id", "unknown")
            if attempts >= max_attempts:
                failures.append({"provider": provider_id, "status": "attempt_budget_exhausted", "error": "global provider attempt budget exhausted"})
                break
            remaining = deadline - time.monotonic()
            if remaining <= 0:
                failures.append({"provider": provider_id, "status": "global_timeout", "error": f"global timeout exceeded after {global_timeout}s"})
                break
            if circuit_is_open(provider_id):
                failures.append({"provider": provider_id, "status": "circuit_open", "error": "provider temporarily skipped after repeated failures"})
                logging.warning("provider skipped because circuit is open provider=%s", provider_id)
                continue
            try:
                attempts += 1
                status, data, content_type = call_provider(model, payload, remaining)
                self.send_response(status)
                self.send_header("Content-Type", content_type)
                self.send_header("Content-Length", str(len(data)))
                self.send_header("X-Universal-AI-Provider", model["id"])
                self.end_headers()
                self.wfile.write(data)
                record_provider_success(provider_id)
                logging.info("routed model=%s via provider=%s status=%s", payload.get("model"), model["id"], status)
                return
            except urllib.error.HTTPError as exc:
                body = exc.read().decode("utf-8", errors="replace")[:1200]
                failures.append({"provider": provider_id, "status": exc.code, "error": body})
                record_provider_failure(provider_id)
                logging.warning("provider failed provider=%s status=%s", provider_id, exc.code)
                if exc.code not in fallback_status:
                    break
            except Exception as exc:  # noqa: BLE001
                failures.append({"provider": provider_id, "status": type(exc).__name__, "error": str(exc)[:500]})
                record_provider_failure(provider_id)
                logging.warning("provider failed provider=%s error=%s", provider_id, exc)
        response_json(
            self,
            503,
            openai_error("All routeable providers failed or are unavailable", "all_providers_failed", 503, failures),
        )


def main() -> int:
    router_cfg = POLICY.get("httpRouter", {})
    host = router_cfg.get("listenHost", "127.0.0.1")
    port = int(router_cfg.get("listenPort", 18100))
    httpd = ThreadingHTTPServer((host, port), RouterHandler)
    logging.info("universal router listening on %s:%s", host, port)
    try:
        httpd.serve_forever()
    finally:
        httpd.server_close()
    return 0


if __name__ == "__main__":
    sys.exit(main())
