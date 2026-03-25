# Handoffs

Lyric エージェントと他エージェントの連携パターン。

## Lyric → Tone

### Purpose
完成した歌詞 + スタイルプロンプトを Tone に渡し、Suno API コード生成・オーディオ処理を依頼する。

### Handoff Format
```
## LYRIC_TO_TONE_HANDOFF

### Song Metadata
- Title: [曲タイトル]
- Genre: [ジャンル]
- Language: [言語]
- Duration Target: [目標尺]

### Style Prompt
[200文字以内のスタイルプロンプト]

### Lyrics
[フォーマット済み歌詞全文]

### Technical Notes
- Character count: [X] / 3,000
- Line count: [X]
- Structure: [構成サマリー]
- Special tags: [使用した特殊タグ]
- Recommended iterations: [X回]
```

### What Tone Does Next
1. Suno API 呼び出しコード生成
2. LUFS正規化設定
3. 出力フォーマット最適化
4. 複数バリエーション生成スクリプト

---

## Tone → Lyric

### Purpose
Suno 生成結果のフィードバックを Lyric に返し、歌詞の修正を依頼する。

### Feedback Triggers
- ボーカルのフレージングがずれる → 音節数の調整
- 特定セクションが期待と異なる → メタタグ修正
- メロディが単調 → 構造の変更
- ボーカルスタイルが合わない → ボーカルタグ追加

### Feedback Format
```
## TONE_TO_LYRIC_FEEDBACK

### Issue
[具体的な問題]

### Affected Section
[問題のあるセクション]

### Suggested Fix
[修正方向の提案]

### Audio Reference
[生成結果のファイルパスまたはID]
```

---

## Quest → Lyric

### Purpose
ゲームのナラティブブリーフに基づいて、ゲーム内楽曲の歌詞を作成する。

### Handoff Format
```
## QUEST_TO_LYRIC_HANDOFF

### Game Context
- Title: [ゲームタイトル]
- Scene: [使用シーン]
- Mood: [ムード]
- Narrative Context: [物語上の位置づけ]

### Musical Direction
- Genre: [ジャンル]
- Tempo: [テンポ感]
- Duration: [尺]
- Vocal: [ボーカルの方向性]

### Lyric Requirements
- Theme: [テーマ]
- Keywords: [必須キーワード]
- Tone: [トーン]
- Language: [言語]
- Constraints: [制約事項]
```

---

## Lyric → Oracle

### Purpose
Suno AI向けプロンプトの最適化について Oracle に相談する。

### Consultation Triggers
- スタイルプロンプトが200文字に収まらない
- 複数のスタイル方向性で迷っている
- 生成品質が安定しない
- コスト最適化（イテレーション回数削減）

### Handoff Format
```
## LYRIC_TO_ORACLE_HANDOFF

### Current Prompt
[現在のスタイルプロンプト]

### Issue
[具体的な課題]

### Attempted Variations
[試したバリエーション]

### Expected vs Actual
[期待した結果 vs 実際の結果]
```

---

## Nexus Integration

### AUTORUN Entry
```
_AGENT_CONTEXT:
  Agent: Lyric
  Task: [composition | formatting | styling | refinement]
  Input: [song request or lyrics to process]
  Constraints: [genre, language, mood, etc.]
```

### AUTORUN Exit
```
_STEP_COMPLETE:
  Agent: Lyric
  Task_Type: [composition | formatting | styling | refinement]
  Status: [done | needs_review | needs_iteration]
  Output: [lyrics + style prompt]
  Handoff: [LYRIC_TO_TONE_HANDOFF | none]
  Next: [Tone for API code | User for review | Oracle for optimization]
  Reason: [brief explanation]
```

### Hub Mode
When receiving `## NEXUS_ROUTING`:
```
## NEXUS_HANDOFF
- Step: [current step number]
- Agent: Lyric
- Summary: [what was done]
- Key findings / decisions: [key choices made]
- Artifacts: [lyrics file, style prompt]
- Risks / trade-offs: [genre ambiguity, iteration needs]
- Open questions: [unresolved items]
- Pending Confirmations: [items needing user input]
- User Confirmations: [items already confirmed]
- Suggested next agent: [Tone | Oracle | User]
- Next action: [what should happen next]
```
