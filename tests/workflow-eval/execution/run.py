"""Run harmless workflow execution cases against the public CLI and the tracked source checkout."""

import argparse
from pathlib import Path
from tempfile import TemporaryDirectory

import debug_regression
import reconciliation
from router import Router


def main():
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument("--output", type=Path, help="Preserve plans, source hashes, generated files, observed evidence and review ledgers")
    args = parser.parse_args()
    scenarios = [("reconciliation", reconciliation), ("debug-regression", debug_regression)]
    names = {source[0] for _, scenario in scenarios for stage in scenario.STAGES for source in stage["sources"]}
    with TemporaryDirectory(prefix="skill-workflow-execution-") as scratch:
        root = Path(scratch)
        out = (args.output or root / "results").resolve()
        out.mkdir(parents=True, exist_ok=True)
        if any((out / name).exists() for name, _ in scenarios):
            parser.error("Output already contains a scenario directory. Choose a fresh output destination.")
        router = Router(root, names)
        for name, scenario in scenarios:
            scenario.run(router, out / name)
        print(f"Two execution cases verified against {len(names)} canonical bundles at {router.revision}.")
        if args.output:
            print(f"Evidence and artifacts: {out}")


if __name__ == "__main__":
    main()
