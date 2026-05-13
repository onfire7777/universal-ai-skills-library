from __future__ import annotations

import argparse
import json
import os
import shutil
import socket
import subprocess
import threading
import time
import urllib.error
import urllib.request
from http.server import BaseHTTPRequestHandler, ThreadingHTTPServer
from pathlib import Path


HOP_BY_HOP_HEADERS = {
    "connection",
    "keep-alive",
    "proxy-authenticate",
    "proxy-authorization",
    "te",
    "trailer",
    "transfer-encoding",
    "upgrade",
}

WIN_PRIORITY_FLAGS = {
    "normal": 0,
    "below-normal": 0x00004000,
    "idle": 0x00000040,
}


def free_ram_gb() -> float | None:
    if os.name != "nt":
        return None
    try:
        import ctypes
        from ctypes import wintypes

        class MemoryStatusEx(ctypes.Structure):
            _fields_ = [
                ("dwLength", wintypes.DWORD),
                ("dwMemoryLoad", wintypes.DWORD),
                ("ullTotalPhys", ctypes.c_ulonglong),
                ("ullAvailPhys", ctypes.c_ulonglong),
                ("ullTotalPageFile", ctypes.c_ulonglong),
                ("ullAvailPageFile", ctypes.c_ulonglong),
                ("ullTotalVirtual", ctypes.c_ulonglong),
                ("ullAvailVirtual", ctypes.c_ulonglong),
                ("ullAvailExtendedVirtual", ctypes.c_ulonglong),
            ]

        status = MemoryStatusEx()
        status.dwLength = ctypes.sizeof(MemoryStatusEx)
        if not ctypes.windll.kernel32.GlobalMemoryStatusEx(ctypes.byref(status)):
            return None
        return status.ullAvailPhys / (1024**3)
    except Exception:
        return None


def free_vram_gb() -> float | None:
    nvidia_smi = shutil.which("nvidia-smi")
    if not nvidia_smi:
        return None
    try:
        result = subprocess.run(
            [
                nvidia_smi,
                "--query-gpu=memory.free",
                "--format=csv,noheader,nounits",
            ],
            check=True,
            capture_output=True,
            text=True,
            timeout=5,
        )
        first = result.stdout.strip().splitlines()[0].strip()
        return float(first) / 1024
    except Exception:
        return None


class ProxyState:
    def __init__(self, args: argparse.Namespace) -> None:
        self.args = args
        self.lock = threading.Lock()
        self.backend_proc: subprocess.Popen | None = None
        self.last_used = time.time()
        self.log_dir = Path(args.log_dir)
        self.log_dir.mkdir(parents=True, exist_ok=True)
        self.proxy_log = self.log_dir / f"{args.name}-proxy.log"
        self.backend_log = self.log_dir / f"{args.name}-llama.log"

    def log(self, message: str) -> None:
        stamp = time.strftime("%Y-%m-%d %H:%M:%S")
        with self.proxy_log.open("a", encoding="utf-8") as f:
            f.write(f"{stamp} {message}\n")

    def port_listening(self) -> bool:
        try:
            with socket.create_connection(("127.0.0.1", self.args.backend_port), timeout=1.0):
                return True
        except OSError:
            return False

    def backend_url(self, path: str) -> str:
        return f"http://127.0.0.1:{self.args.backend_port}{path}"

    def check_resource_headroom(self) -> None:
        ram = free_ram_gb()
        if self.args.min_free_ram_gb > 0 and ram is not None and ram < self.args.min_free_ram_gb:
            raise RuntimeError(
                f"refusing to start {self.args.name}: free RAM {ram:.1f}GB is below "
                f"{self.args.min_free_ram_gb:.1f}GB guard"
            )
        vram = free_vram_gb()
        if self.args.min_free_vram_gb > 0 and vram is not None and vram < self.args.min_free_vram_gb:
            raise RuntimeError(
                f"refusing to start {self.args.name}: free VRAM {vram:.1f}GB is below "
                f"{self.args.min_free_vram_gb:.1f}GB guard"
            )
        self.log(
            "resource guard passed "
            f"free_ram_gb={ram if ram is not None else 'unknown'} "
            f"free_vram_gb={vram if vram is not None else 'unknown'}"
        )

    def start_backend(self) -> None:
        if self.port_listening():
            return
        model = Path(self.args.model)
        llama_server = Path(self.args.llama_server)
        if not model.exists():
            raise RuntimeError(f"model file missing: {model}")
        if not llama_server.exists():
            raise RuntimeError(f"llama-server missing: {llama_server}")
        self.check_resource_headroom()

        cmd = [
            str(llama_server),
            "-m",
            str(model),
            "--host",
            "127.0.0.1",
            "--port",
            str(self.args.backend_port),
            "--alias",
            self.args.alias,
            "--ctx-size",
            str(self.args.ctx_size),
            "--n-gpu-layers",
            str(self.args.n_gpu_layers),
            "--flash-attn",
            self.args.flash_attn,
            "--batch-size",
            str(self.args.batch_size),
            "--ubatch-size",
            str(self.args.ubatch_size),
            "--threads",
            str(self.args.threads),
            "--cache-type-k",
            self.args.cache_type_k,
            "--cache-type-v",
            self.args.cache_type_v,
            "--jinja",
            "--parallel",
            str(self.args.parallel),
            "--cont-batching",
            "--no-webui",
            "--no-webui-mcp-proxy",
        ]
        for item in self.args.extra_llama_arg:
            cmd.append(item)

        env = os.environ.copy()
        env["PATH"] = str(llama_server.parent) + os.pathsep + env.get("PATH", "")
        flags = 0
        if os.name == "nt":
            flags = 0x00000008 | 0x00000200 | 0x08000000
            flags |= WIN_PRIORITY_FLAGS.get(self.args.process_priority, 0)

        self.log(f"starting backend on 127.0.0.1:{self.args.backend_port}")
        backend_log_handle = self.backend_log.open("ab", buffering=0)
        try:
            self.backend_proc = subprocess.Popen(
                cmd,
                cwd=str(llama_server.parent),
                env=env,
                stdin=subprocess.DEVNULL,
                stdout=backend_log_handle,
                stderr=backend_log_handle,
                creationflags=flags,
                close_fds=True,
            )
        except Exception:
            backend_log_handle.close()
            raise

    def ensure_backend(self) -> None:
        with self.lock:
            self.last_used = time.time()
            if self.port_listening():
                return
            self.start_backend()

        deadline = time.time() + self.args.startup_timeout
        while time.time() < deadline:
            if self.port_listening() and self.backend_ready():
                self.log("backend ready")
                return
            if self.backend_proc is not None and self.backend_proc.poll() is not None:
                raise RuntimeError(f"backend exited with code {self.backend_proc.returncode}")
            time.sleep(1.0)
        raise TimeoutError(f"backend did not become ready within {self.args.startup_timeout}s")

    def backend_ready(self) -> bool:
        for path in ("/health", "/v1/models"):
            try:
                with urllib.request.urlopen(self.backend_url(path), timeout=2.0) as response:
                    if 200 <= response.status < 500:
                        return True
            except Exception:
                pass
        return False

    def stop_backend_if_idle(self) -> None:
        with self.lock:
            proc = self.backend_proc
            if proc is None:
                return
            if proc.poll() is not None:
                self.backend_proc = None
                return
            if time.time() - self.last_used < self.args.idle_timeout:
                return
            self.log("stopping idle backend")
            proc.terminate()
        try:
            proc.wait(timeout=20)
        except subprocess.TimeoutExpired:
            proc.kill()
            try:
                proc.wait(timeout=10)
            except subprocess.TimeoutExpired:
                pass
        with self.lock:
            if self.backend_proc is proc:
                self.backend_proc = None


def make_handler(state: ProxyState):
    class Handler(BaseHTTPRequestHandler):
        server_version = "HermesLocalQwenProxy/1.0"

        def log_message(self, fmt: str, *args) -> None:
            state.log(fmt % args)

        def _json(self, status: int, payload: dict) -> None:
            data = json.dumps(payload).encode("utf-8")
            self.send_response(status)
            self.send_header("Content-Type", "application/json")
            self.send_header("Content-Length", str(len(data)))
            self.end_headers()
            self.wfile.write(data)

        def do_GET(self) -> None:
            if self.path.rstrip("/") == "/health":
                self._json(
                    200,
                    {
                        "status": "ok",
                        "name": state.args.name,
                        "alias": state.args.alias,
                        "backend_port": state.args.backend_port,
                        "backend_running": state.port_listening(),
                        "lazy": True,
                    },
                )
                return
            if self.path.rstrip("/") == "/v1/models":
                self._json(
                    200,
                    {
                        "object": "list",
                        "data": [
                            {
                                "id": state.args.alias,
                                "object": "model",
                                "owned_by": "local",
                            }
                        ],
                    },
                )
                return
            self._forward(start_backend=True)

        def do_POST(self) -> None:
            self._forward(start_backend=True)

        def _forward(self, *, start_backend: bool) -> None:
            try:
                length = int(self.headers.get("Content-Length", "0") or "0")
                if state.args.max_body_mb > 0 and length > state.args.max_body_mb * 1024 * 1024:
                    self._json(
                        413,
                        {
                            "error": {
                                "message": f"request body exceeds {state.args.max_body_mb}MB local-model guard",
                                "type": "local_qwen_proxy_request_too_large",
                            }
                        },
                    )
                    return
                if start_backend:
                    state.ensure_backend()
                state.last_used = time.time()
                body = self.rfile.read(length) if length > 0 else None
                headers = {
                    key: value
                    for key, value in self.headers.items()
                    if key.lower() not in HOP_BY_HOP_HEADERS and key.lower() != "host"
                }
                request = urllib.request.Request(
                    state.backend_url(self.path),
                    data=body,
                    headers=headers,
                    method=self.command,
                )
                with urllib.request.urlopen(request, timeout=state.args.request_timeout) as response:
                    self.send_response(response.status)
                    for key, value in response.headers.items():
                        if key.lower() not in HOP_BY_HOP_HEADERS:
                            self.send_header(key, value)
                    self.end_headers()
                    shutil.copyfileobj(response, self.wfile, length=1024 * 1024)
                state.last_used = time.time()
            except urllib.error.HTTPError as exc:
                payload = exc.read()
                self.send_response(exc.code)
                for key, value in exc.headers.items():
                    if key.lower() not in HOP_BY_HOP_HEADERS:
                        self.send_header(key, value)
                self.end_headers()
                self.wfile.write(payload)
            except Exception as exc:
                state.log(f"request failed: {exc!r}")
                self._json(
                    503,
                    {
                        "error": {
                            "message": str(exc),
                            "type": "local_qwen_proxy_error",
                        }
                    },
                )

    return Handler


def idle_loop(state: ProxyState) -> None:
    while True:
        time.sleep(30)
        state.stop_backend_if_idle()


def parse_args() -> argparse.Namespace:
    parser = argparse.ArgumentParser()
    parser.add_argument("--name", required=True)
    parser.add_argument("--listen-port", type=int, required=True)
    parser.add_argument("--backend-port", type=int, required=True)
    parser.add_argument("--alias", required=True)
    parser.add_argument("--model", required=True)
    parser.add_argument("--ctx-size", type=int, required=True)
    parser.add_argument(
        "--llama-server",
        default=str(Path.home() / ".local-ai" / "runtimes" / "llama.cpp-cuda" / "b9128-cuda12.4" / "llama-server.exe"),
    )
    parser.add_argument("--batch-size", type=int, default=512)
    parser.add_argument("--ubatch-size", type=int, default=256)
    parser.add_argument("--threads", type=int, default=6)
    parser.add_argument("--parallel", type=int, default=1)
    parser.add_argument("--n-gpu-layers", type=int, default=99)
    parser.add_argument("--flash-attn", default="on")
    parser.add_argument("--cache-type-k", default="q8_0")
    parser.add_argument("--cache-type-v", default="q8_0")
    parser.add_argument("--extra-llama-arg", action="append", default=[])
    parser.add_argument(
        "--log-dir",
        default=str(Path(os.environ.get("UNIVERSAL_AI_STACK_HOME", Path.home() / ".universal-ai-stack")) / "logs" / "local-qwen"),
    )
    parser.add_argument("--idle-timeout", type=int, default=600)
    parser.add_argument("--startup-timeout", type=int, default=900)
    parser.add_argument("--request-timeout", type=int, default=1800)
    parser.add_argument("--min-free-vram-gb", type=float, default=0.0)
    parser.add_argument("--min-free-ram-gb", type=float, default=0.0)
    parser.add_argument("--max-body-mb", type=int, default=8)
    parser.add_argument(
        "--process-priority",
        choices=sorted(WIN_PRIORITY_FLAGS.keys()),
        default="below-normal",
    )
    return parser.parse_args()


def main() -> None:
    args = parse_args()
    state = ProxyState(args)
    state.log(f"proxy listening on 127.0.0.1:{args.listen_port}")
    threading.Thread(target=idle_loop, args=(state,), daemon=True).start()
    server = ThreadingHTTPServer(("127.0.0.1", args.listen_port), make_handler(state))
    server.serve_forever()


if __name__ == "__main__":
    main()
