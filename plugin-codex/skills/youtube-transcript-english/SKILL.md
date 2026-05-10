---
name: youtube-transcript-english
description: Use when fetching YouTube transcripts that must come back in English. Wraps the locally installed jdepoix/youtube-transcript-api package, requests English first, uses YouTube translation when available, then falls back to AI translation. Use for video transcript extraction, caption retrieval, transcript-to-text, JSON/SRT/VTT export, and transcript ingestion before summarization.
license: MIT
---

# YouTube Transcript English

Use this skill when a YouTube transcript must come back in English.

## Source

- Upstream repository: `https://github.com/jdepoix/youtube-transcript-api`
- Integrated source commit: `85c859450c014d860297ab9fd4f88145c47ff2e2`
- Local checkout: `C:\Users\burni\.youtube-transcript-api\youtube-transcript-api`
- Isolated venv CLI: `C:\Users\burni\.youtube-transcript-api\venv\Scripts\youtube_transcript_api.exe`
- English-enforcing shim: `C:\Users\burni\.youtube-transcript-api\bin\youtube-transcript-english.cmd`

## README-Grounded Usage

The upstream README states that `YouTubeTranscriptApi().fetch(video_id)` defaults to English, accepts `languages=['en']`, supports `TranscriptList`, and can translate transcripts with `transcript.translate('en')` when translation is available.

Use the bundled English wrapper first:

```bash
C:\Users\burni\.youtube-transcript-api\venv\Scripts\python.exe scripts\fetch_english_transcript.py "https://www.youtube.com/watch?v=VIDEO_ID"
C:\Users\burni\.youtube-transcript-api\venv\Scripts\python.exe scripts\fetch_english_transcript.py VIDEO_ID --format json --output transcript.json
C:\Users\burni\.youtube-transcript-api\venv\Scripts\python.exe scripts\fetch_english_transcript.py VIDEO_ID --format srt --output transcript.srt
C:\Users\burni\.youtube-transcript-api\venv\Scripts\python.exe scripts\fetch_english_transcript.py VIDEO_ID --ai-model gpt-4o-mini
```

Equivalent direct command:

```bash
C:\Users\burni\.youtube-transcript-api\bin\youtube-transcript-english.cmd VIDEO_ID --format text
```

Use the upstream CLI when exact upstream behavior is needed:

```bash
C:\Users\burni\.youtube-transcript-api\venv\Scripts\youtube_transcript_api.exe VIDEO_ID --languages en --format json
C:\Users\burni\.youtube-transcript-api\venv\Scripts\youtube_transcript_api.exe --list-transcripts VIDEO_ID
```

## English Output Policy

1. Request direct English captions first.
2. If direct English is unavailable, inspect available transcripts.
3. If a transcript is translatable to English, translate it to English.
4. If YouTube cannot provide English, fetch the best available source transcript and translate caption segments to English with an OpenAI-compatible model.
5. If no transcript exists or AI translation is required but `OPENAI_API_KEY` is not set, fail clearly instead of returning a non-English transcript.

Add `--no-translate` to require direct English captions only. Add `--no-ai-translate` to allow YouTube translation but disable AI translation fallback.

AI translation settings:

- `OPENAI_API_KEY` is required only when AI translation fallback is needed.
- `OPENAI_TRANSLATION_MODEL` or `AI_TRANSLATION_MODEL` can override the default model.
- `OPENAI_BASE_URL` can point at an OpenAI-compatible endpoint.
- `--ai-chunk-size` controls how many caption segments are translated per request.

## References

- `references/source-readme.md` contains the full upstream README.
- `references/requirements-lock.txt` documents installed dependencies.
- `references/source-metadata.json` records the integrated commit and local paths.

## Reliability Notes

- YouTube can block requests from some IPs. The upstream README documents `RequestBlocked` and `IpBlocked` behavior and proxy options.
- The library uses an undocumented YouTube web-client API, so live behavior can change.
- For bulk or repeated fetching, respect rate limits and avoid unnecessary retries.
