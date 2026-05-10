#!/usr/bin/env python3
"""Fetch a YouTube transcript in English, translating to English when available."""

from __future__ import annotations

import argparse
import json
import os
import re
import sys
from pathlib import Path
from urllib.parse import parse_qs, urlparse

from openai import OpenAI
from youtube_transcript_api import YouTubeTranscriptApi
from youtube_transcript_api._transcripts import FetchedTranscript, FetchedTranscriptSnippet
from youtube_transcript_api.formatters import (
    PrettyPrintFormatter,
    SRTFormatter,
    TextFormatter,
    WebVTTFormatter,
)


VIDEO_ID_RE = re.compile(r"^[A-Za-z0-9_-]{11}$")


def extract_video_id(value: str) -> str:
    value = value.strip()
    if VIDEO_ID_RE.match(value):
        return value

    parsed = urlparse(value)
    host = parsed.hostname or ""
    if host.endswith("youtu.be"):
        candidate = parsed.path.strip("/").split("/")[0]
        if VIDEO_ID_RE.match(candidate):
            return candidate
    if host.endswith("youtube.com"):
        query_id = parse_qs(parsed.query).get("v", [""])[0]
        if VIDEO_ID_RE.match(query_id):
            return query_id
        parts = [part for part in parsed.path.split("/") if part]
        for marker in ("embed", "shorts", "live"):
            if marker in parts:
                idx = parts.index(marker) + 1
                if idx < len(parts) and VIDEO_ID_RE.match(parts[idx]):
                    return parts[idx]

    raise SystemExit(f"Could not extract a YouTube video ID from: {value}")


def translation_codes(transcript) -> set[str]:
    codes = set()
    for item in getattr(transcript, "translation_languages", []) or []:
        code = getattr(item, "language_code", None)
        if code:
            codes.add(code)
        elif isinstance(item, dict) and item.get("language_code"):
            codes.add(item["language_code"])
    return codes


def _json_from_model_text(text: str) -> dict:
    text = text.strip()
    if text.startswith("```"):
        lines = text.splitlines()
        if lines and lines[0].strip().startswith("```"):
            lines = lines[1:]
        if lines and lines[-1].strip() == "```":
            lines = lines[:-1]
        text = "\n".join(lines).strip()
    parsed = json.loads(text)
    if not isinstance(parsed, dict):
        raise ValueError("AI translation response must be a JSON object")
    return parsed


def _translate_chunk_with_ai(
    chunk: list[tuple[int, str]],
    *,
    source_language: str,
    model: str,
    base_url: str | None,
    api_key: str,
) -> dict[int, str]:
    client_kwargs = {"api_key": api_key}
    if base_url:
        client_kwargs["base_url"] = base_url.rstrip("/")
    client = OpenAI(**client_kwargs)
    payload = [{"index": index, "text": text} for index, text in chunk]
    response = client.chat.completions.create(
        model=model,
        temperature=0,
        response_format={"type": "json_object"},
        messages=[
            {
                "role": "system",
                "content": (
                    "Translate transcript caption segments into natural English. "
                    "Preserve meaning, names, numbers, code terms, and concise caption style. "
                    "Return only JSON shaped as {\"translations\":[{\"index\":0,\"text\":\"...\"}]}."
                ),
            },
            {
                "role": "user",
                "content": json.dumps(
                    {
                        "source_language": source_language,
                        "segments": payload,
                    },
                    ensure_ascii=False,
                ),
            },
        ],
    )
    content = response.choices[0].message.content or ""
    parsed = _json_from_model_text(content)
    translations = parsed.get("translations")
    if not isinstance(translations, list):
        raise RuntimeError("AI translation response missing translations list")

    expected = {index for index, _ in chunk}
    out: dict[int, str] = {}
    for item in translations:
        if not isinstance(item, dict):
            continue
        index = item.get("index")
        text = item.get("text")
        if isinstance(index, int) and isinstance(text, str):
            out[index] = text.strip()

    if set(out) != expected:
        missing = sorted(expected - set(out))
        extra = sorted(set(out) - expected)
        raise RuntimeError(f"AI translation index mismatch; missing={missing}, extra={extra}")
    return out


def translate_fetched_transcript_to_english(
    fetched: FetchedTranscript,
    *,
    source_language: str,
    model: str,
    base_url: str | None,
    chunk_size: int,
) -> FetchedTranscript:
    api_key = os.getenv("OPENAI_API_KEY")
    if not api_key:
        raise RuntimeError(
            "No direct English or YouTube English translation is available, "
            "and OPENAI_API_KEY is not set for AI translation."
        )
    if chunk_size < 1:
        raise ValueError("--ai-chunk-size must be at least 1")

    snippets = list(fetched)
    translated: dict[int, str] = {}
    indexed = [(index, snippet.text) for index, snippet in enumerate(snippets)]
    for start in range(0, len(indexed), chunk_size):
        chunk = indexed[start : start + chunk_size]
        translated.update(
            _translate_chunk_with_ai(
                chunk,
                source_language=source_language,
                model=model,
                base_url=base_url,
                api_key=api_key,
            )
        )

    translated_snippets = [
        FetchedTranscriptSnippet(
            text=translated[index],
            start=snippet.start,
            duration=snippet.duration,
        )
        for index, snippet in enumerate(snippets)
    ]
    return FetchedTranscript(
        snippets=translated_snippets,
        video_id=fetched.video_id,
        language="English",
        language_code="en",
        is_generated=fetched.is_generated,
    )


def choose_source_transcript(transcript_list, preserve_formatting: bool):
    transcripts = list(transcript_list)
    if not transcripts:
        raise RuntimeError("No transcripts are available for this video.")
    manual = [transcript for transcript in transcripts if not transcript.is_generated]
    selected = manual[0] if manual else transcripts[0]
    return selected.fetch(preserve_formatting=preserve_formatting), selected


def fetch_english(
    video_id: str,
    preserve_formatting: bool,
    allow_translate: bool,
    allow_ai_translate: bool,
    ai_model: str,
    ai_base_url: str | None,
    ai_chunk_size: int,
):
    api = YouTubeTranscriptApi()
    try:
        fetched = api.fetch(video_id, languages=["en"], preserve_formatting=preserve_formatting)
        return fetched, {"method": "direct", "translated": False, "source_language_code": fetched.language_code}
    except Exception as direct_error:
        if not allow_translate:
            raise RuntimeError(f"No direct English transcript available: {direct_error}") from direct_error

    transcript_list = list(api.list(video_id))
    for transcript in transcript_list:
        if transcript.language_code == "en" or transcript.language_code.startswith("en-"):
            fetched = transcript.fetch(preserve_formatting=preserve_formatting)
            return fetched, {"method": "listed-english", "translated": False, "source_language_code": transcript.language_code}

    for transcript in transcript_list:
        if transcript.is_translatable and "en" in translation_codes(transcript):
            fetched = transcript.translate("en").fetch(preserve_formatting=preserve_formatting)
            return fetched, {
                "method": "translated-to-english",
                "translated": True,
                "source_language_code": transcript.language_code,
            }

    if allow_ai_translate:
        source_fetched, source_transcript = choose_source_transcript(transcript_list, preserve_formatting)
        fetched = translate_fetched_transcript_to_english(
            source_fetched,
            source_language=f"{source_transcript.language_code} ({source_transcript.language})",
            model=ai_model,
            base_url=ai_base_url,
            chunk_size=ai_chunk_size,
        )
        return fetched, {
            "method": "ai-translated-to-english",
            "translated": True,
            "ai_translated": True,
            "source_language_code": source_transcript.language_code,
            "ai_model": ai_model,
        }

    available = [f"{t.language_code} ({t.language})" for t in transcript_list]
    raise RuntimeError(
        "No English transcript or English translation is available. "
        f"Available transcripts: {', '.join(available) if available else 'none'}"
    )


def format_transcript(fetched, output_format: str, metadata: dict) -> str:
    if output_format == "json":
        return json.dumps(
            {
                "metadata": metadata,
                "transcript": fetched.to_raw_data(),
            },
            ensure_ascii=False,
            indent=2,
        )

    formatter = {
        "pretty": PrettyPrintFormatter(),
        "text": TextFormatter(),
        "srt": SRTFormatter(),
        "vtt": WebVTTFormatter(),
    }[output_format]
    return formatter.format_transcript(fetched)


def main() -> int:
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument("video", help="YouTube video ID or URL.")
    parser.add_argument("--format", choices=["text", "json", "pretty", "srt", "vtt"], default="text")
    parser.add_argument("--output", help="Optional output file path.")
    parser.add_argument("--preserve-formatting", action="store_true")
    parser.add_argument("--no-translate", action="store_true", help="Require a direct English transcript; do not translate.")
    parser.add_argument("--no-ai-translate", action="store_true", help="Disable AI translation fallback.")
    parser.add_argument(
        "--ai-model",
        default=os.getenv("OPENAI_TRANSLATION_MODEL") or os.getenv("AI_TRANSLATION_MODEL") or "gpt-4o-mini",
        help="OpenAI-compatible model for AI translation fallback.",
    )
    parser.add_argument(
        "--ai-base-url",
        default=os.getenv("OPENAI_BASE_URL"),
        help="Optional OpenAI-compatible base URL. Defaults to OpenAI when unset.",
    )
    parser.add_argument("--ai-chunk-size", type=int, default=40, help="Caption segments per AI translation request.")
    parser.add_argument("--metadata", action="store_true", help="Print metadata before text output.")
    args = parser.parse_args()

    video_id = extract_video_id(args.video)
    fetched, meta = fetch_english(
        video_id,
        args.preserve_formatting,
        not args.no_translate,
        not args.no_translate and not args.no_ai_translate,
        args.ai_model,
        args.ai_base_url,
        args.ai_chunk_size,
    )
    metadata = {
        "video_id": video_id,
        "language": fetched.language,
        "language_code": fetched.language_code,
        "is_generated": fetched.is_generated,
        **meta,
    }
    rendered = format_transcript(fetched, args.format, metadata)

    if args.metadata and args.format != "json":
        rendered = json.dumps(metadata, ensure_ascii=False) + "\n\n" + rendered

    if args.output:
        Path(args.output).write_text(rendered, encoding="utf-8")
    else:
        sys.stdout.write(rendered)
        if not rendered.endswith("\n"):
            sys.stdout.write("\n")
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
