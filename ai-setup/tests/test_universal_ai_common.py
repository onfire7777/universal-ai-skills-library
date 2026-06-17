#!/usr/bin/env python3
"""Behavioural tests for the ai-setup runtime shared helpers.

These pin the contract of the helpers extracted (DRY) from
``universal_ai_router.py`` and ``universal_ai_stack_supervisor.py`` into
``_universal_ai_common``, and verify that both entrypoints still load and reuse
the single shared implementation.  This is a refactor-only change: the observable
behaviour of each entrypoint must be identical to before the extraction.

Written with stdlib ``unittest`` so the suite runs under both
``python -m unittest`` and ``pytest`` without extra dependencies.
"""
from __future__ import annotations

import importlib
import logging
import os
import sys
import tempfile
import unittest
from pathlib import Path

# runtime/bin is a plain scripts directory (not a package); import modules from it
# by putting it on sys.path, mirroring how the entrypoints are launched.
BIN_DIR = Path(__file__).resolve().parents[1] / "runtime" / "bin"


def _import_from_bin(name: str):
    """Import (or re-import) a module that lives in runtime/bin."""
    if str(BIN_DIR) not in sys.path:
        sys.path.insert(0, str(BIN_DIR))
    sys.modules.pop(name, None)
    return importlib.import_module(name)


def _reset_root_logging() -> None:
    """Detach handlers added by ``setup_logging`` so temp log files can be removed."""
    root = logging.getLogger()
    for handler in list(root.handlers):
        root.removeHandler(handler)
        try:
            handler.close()
        except Exception:  # pragma: no cover - best-effort cleanup
            pass


class LoadJsonTests(unittest.TestCase):
    def test_reads_utf8_json(self) -> None:
        common = _import_from_bin("_universal_ai_common")
        with tempfile.TemporaryDirectory() as d:
            path = Path(d) / "x.json"
            path.write_text('{"a": 1, "b": "café"}', encoding="utf-8")
            self.assertEqual(common.load_json(path), {"a": 1, "b": "café"})


class LoadEnvTests(unittest.TestCase):
    def setUp(self) -> None:
        self.common = _import_from_bin("_universal_ai_common")
        self._saved_env = dict(os.environ)

    def tearDown(self) -> None:
        os.environ.clear()
        os.environ.update(self._saved_env)

    def _write_env(self, text: str) -> Path:
        path = Path(tempfile.mkdtemp()) / ".env"
        path.write_text(text, encoding="utf-8")
        return path

    def test_missing_file_returns_environ_copy(self) -> None:
        os.environ["SAMPLE_KEY"] = "keep"
        env = self.common.load_env(Path(tempfile.mkdtemp()) / "absent.env")
        self.assertEqual(env["SAMPLE_KEY"], "keep")

    def test_comments_blanks_and_malformed_lines_ignored(self) -> None:
        path = self._write_env("# comment\n\nNO_EQUALS_LINE\nGOOD=1\n")
        env = self.common.load_env(path, override=True)
        self.assertEqual(env["GOOD"], "1")
        self.assertNotIn("NO_EQUALS_LINE", env)

    def test_default_does_not_override_existing_env(self) -> None:
        # Router semantics: process environment wins over the dotenv file.
        os.environ["SHARED_KEY"] = "from-env"
        path = self._write_env("SHARED_KEY=from-file\nNEW_KEY=v\n")
        env = self.common.load_env(path)
        self.assertEqual(env["SHARED_KEY"], "from-env")
        self.assertEqual(env["NEW_KEY"], "v")

    def test_override_replaces_existing_env(self) -> None:
        # Supervisor semantics: the dotenv file overrides the process environment.
        os.environ["SHARED_KEY"] = "from-env"
        path = self._write_env("SHARED_KEY=from-file\n")
        env = self.common.load_env(path, override=True)
        self.assertEqual(env["SHARED_KEY"], "from-file")

    def test_strips_surrounding_double_quotes(self) -> None:
        path = self._write_env('TOKEN="abc123"\n')
        env = self.common.load_env(path, override=True)
        self.assertEqual(env["TOKEN"], "abc123")

    def test_inject_home_sets_stack_home_even_when_file_missing(self) -> None:
        os.environ.pop("UNIVERSAL_AI_STACK_HOME", None)
        env = self.common.load_env(
            Path(tempfile.mkdtemp()) / "absent.env", inject_home=True
        )
        self.assertEqual(env["UNIVERSAL_AI_STACK_HOME"], str(self.common.ROOT))


class SetupLoggingTests(unittest.TestCase):
    def tearDown(self) -> None:
        _reset_root_logging()

    def test_creates_log_directory(self) -> None:
        common = _import_from_bin("_universal_ai_common")
        with tempfile.TemporaryDirectory() as d:
            log_dir = Path(d) / "logs"
            original = common.LOG_DIR
            common.LOG_DIR = log_dir
            try:
                common.setup_logging("unit.log")
                self.assertTrue(log_dir.is_dir())
            finally:
                common.LOG_DIR = original


class EntrypointWiringTests(unittest.TestCase):
    """Refactor-only guarantee: both entrypoints import cleanly and reuse the
    single shared helper implementation rather than private copies."""

    def setUp(self) -> None:
        self._saved_env = dict(os.environ)
        self._tracked = (
            "_universal_ai_common",
            "universal_ai_router",
            "universal_ai_stack_supervisor",
        )
        self._saved_modules = {n: sys.modules.get(n) for n in self._tracked}

    def tearDown(self) -> None:
        os.environ.clear()
        os.environ.update(self._saved_env)
        for name, module in self._saved_modules.items():
            if module is None:
                sys.modules.pop(name, None)
            else:
                sys.modules[name] = module
        _reset_root_logging()

    def test_entrypoints_load_and_share_helpers(self) -> None:
        with tempfile.TemporaryDirectory() as d:
            home = Path(d)
            (home / "config").mkdir()
            (home / "config" / "model-registry.json").write_text(
                '{"models": []}', encoding="utf-8"
            )
            (home / "config" / "routing-policy.json").write_text(
                "{}", encoding="utf-8"
            )
            os.environ["UNIVERSAL_AI_STACK_HOME"] = str(home)

            common = _import_from_bin("_universal_ai_common")
            router = _import_from_bin("universal_ai_router")
            supervisor = _import_from_bin("universal_ai_stack_supervisor")

            # Single source of truth: entrypoints reuse the shared callables.
            self.assertIs(router.load_json, common.load_json)
            self.assertIs(router.load_env, common.load_env)
            self.assertIs(router.setup_logging, common.setup_logging)
            self.assertIs(supervisor.load_json, common.load_json)
            self.assertIs(supervisor.load_env, common.load_env)
            self.assertIs(supervisor.setup_logging, common.setup_logging)

            # Path resolution honours the UNIVERSAL_AI_STACK_HOME override.
            self.assertEqual(common.ROOT, home)
            self.assertEqual(router.CONFIG_DIR, home / "config")
            self.assertEqual(supervisor.STATE_DIR, home / "state")


if __name__ == "__main__":
    unittest.main()
