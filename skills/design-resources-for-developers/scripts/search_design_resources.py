#!/usr/bin/env python3
"""Search the integrated design resources catalogue."""

from __future__ import annotations

import argparse
import json
from pathlib import Path


ROOT = Path(__file__).resolve().parents[1]
DATA_PATH = ROOT / "references" / "resources.json"


def load_resources() -> list[dict[str, str]]:
    with DATA_PATH.open("r", encoding="utf-8") as handle:
        return json.load(handle)


def matches(item: dict[str, str], query: str | None, category: str | None) -> bool:
    if category and item["category"].lower() != category.lower():
        return False
    if not query:
        return True
    haystack = " ".join(
        [
            item.get("category", ""),
            item.get("name", ""),
            item.get("description", ""),
            item.get("url", ""),
        ]
    ).lower()
    return all(part in haystack for part in query.lower().split())


def format_markdown(items: list[dict[str, str]]) -> str:
    if not items:
        return "No resources matched."

    lines: list[str] = []
    current_category = None
    for item in items:
        if item["category"] != current_category:
            current_category = item["category"]
            lines.append(f"\n## {current_category}")
        description = item.get("description", "").strip()
        suffix = f" - {description}" if description else ""
        lines.append(f"- [{item['name']}]({item['url']}){suffix}")
    return "\n".join(lines).strip()


def main() -> int:
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument("--query", help="Search terms matched across category, name, URL, and description.")
    parser.add_argument("--category", help="Exact category name to filter.")
    parser.add_argument("--limit", type=int, default=25, help="Maximum resources to return.")
    parser.add_argument("--format", choices=["json", "md"], default="md", help="Output format.")
    args = parser.parse_args()

    resources = [item for item in load_resources() if matches(item, args.query, args.category)]
    resources = resources[: max(args.limit, 0)]

    if args.format == "json":
        print(json.dumps(resources, ensure_ascii=False, indent=2))
    else:
        print(format_markdown(resources))
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
