# Suno AI Format Guide

Suno AIの技術仕様・メタタグ・制約の完全リファレンス。

## Technical Constraints

| Field | Limit | Notes |
|-------|-------|-------|
| Style Prompt | ~200 chars | Excess is silently truncated |
| Lyrics | ~3,000 chars | 40-60 lines, 200-300 words equivalent |
| Song Title | 80 chars | Minimal impact on music output |
| Recommended lines | 30-40 | For standard 3-4 min songs |
| Lines per section | 2-6 | Longer sections cause vocal drift |
| Audio output | 1-4 min | Per generation |

- 3,000文字超過: セクションが駆け足、または短い出力になる
- 15行未満: 曲が短縮される傾向
- v4.5/v5: APIレベルで最大5,000文字対応

## Metatag Syntax Rules

1. 角括弧 `[ ]` で囲む（コロン `Verse:` は歌詞として歌われる）
2. 独立した行に配置（歌詞の途中に埋め込まない）
3. タグは1-3語で短く保つ
4. セクション間に空白行を入れる
5. 最初の20-30語のタグが最も強い影響力を持つ

## Pipe Stacking

複数の指示を1つの括弧内で組み合わせる:
```
[Chorus | chill | synth lead | soft vocal]
[Verse | spoken word | low energy]
[Bridge | female vocal | melancholic]
```
- 2-4要素に抑える（5つ以上は非推奨）

## Structure Tags

### Primary Structure
| Tag | Purpose | Notes |
|-----|---------|-------|
| `[Short Instrumental Intro]` | Intro | `[Intro]`より安定して動作 |
| `[Verse]` / `[Verse 1]` | Verse | 各2-6行 |
| `[Pre-Chorus]` | Pre-Chorus | 期待感の醸成、短め |
| `[Chorus]` | Chorus/Hook | 繰り返し可能なフック |
| `[Post-Chorus]` | Post-Chorus | コーラス後の余韻 |
| `[Hook]` | Hook | 短く印象的なフレーズ |
| `[Bridge]` | Bridge | 対比とスペース |
| `[Outro]` | Outro | 曲の締めくくり |

### Dynamic Tags
| Tag | Purpose |
|-----|---------|
| `[Build-Up]` / `[Buildup]` | テンション上昇 |
| `[Drop]` | EDMドロップ |
| `[Break]` / `[Percussion Break]` | ブレイク |
| `[Instrumental]` / `[Instrumental Break]` | インスト部分 |
| `[Interlude]` / `[Melodic Interlude]` | インタールード |
| `[Solo]` / `[Guitar Solo]` | ソロパート |

### Ending Tags
| Tag | Purpose |
|-----|---------|
| `[Final Chorus]` | 最大盛り上がりのコーラス |
| `[Big Finish]` | 大サビ的な終結 |
| `[Refrain]` | リフレイン（`[Outro]`より創造的） |
| `[Fade Out]` / `[Fade to End]` | フェードアウト |
| `[End]` | 即座に終了 |

## Punctuation as Musical Direction

| Punctuation | Effect |
|-------------|--------|
| `,` (comma) | マイクロポーズ（微小な間） |
| `--` (dash) | 息継ぎ、スタッガード・デリバリー |
| `...` (ellipsis) | フレージングのゆらぎ |
| `!` (exclamation) | アグレッシブなアタック（次行に伝播注意） |
| Line break | 音楽的なブレスの位置 |
| ALL CAPS | ボーカルの強調・トーン変化 |

## Ad-libs and Vocal Effects

### Ad-libs (parentheses)
```
Now I know (yeah) now I know (uh-huh)
We're gonna make it (oh yeah!)
```

### Sustained notes (hyphens/repetition)
```
lo-ove
sooo-long
knooowwww
```

## Critical Anti-Patterns

| Anti-Pattern | Fix |
|-------------|-----|
| `[Intro]` tag alone | Use `[Short Instrumental Intro]` |
| `[My Custom Tag]` | Only use recognized standard tags |
| `repeat chorus` or similar | Write full chorus text every time |
| Style text as plain lyrics | Always use `[brackets]` for tags |
| `*sound effects*` in lyrics | No asterisks, SFX descriptions, or sound cues |
| `!` overuse | Use sparingly; aggression bleeds to next line |
| Too many pipe elements | Keep to 2-4 per bracket |

## Style Prompt Priority Order (200 chars)

1. Genre/sub-genre (e.g., indie pop, lo-fi hip hop)
2. Vocal direction (e.g., female mid-range, breathy)
3. Primary mood (e.g., melancholic, uplifting)
4. 1-2 instruments (e.g., acoustic guitar, piano)
5. BPM/tempo (e.g., 90 BPM, mid-tempo)
6. Production quality (e.g., lo-fi, polished)

Rules:
- Drop articles, use comma separation
- 4-8 style tags is the sweet spot
- Adding era changes sound drastically ("80s synth-pop" vs "2020s synth-pop")
- Avoid contradictory tags (aggressive + calm)

## Suno v5 Specific

- **Top anchor**: Put vocal role + BPM + structure summary at prompt start
- **Syllable count**: Specify "Verse lines: 8-10 syllables"
- **Rhyme scheme**: Specify "AABB rhyme scheme"
- **Lyric fidelity**: Add "Do not change any words inside brackets. Sing exactly as written."
- **Pronunciation fix**: Adjust at text level (e.g., "bahss" for "bass")
