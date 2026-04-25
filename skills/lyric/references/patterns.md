# Patterns & Anti-Patterns

よくあるミスと推奨パターン。

## Anti-Patterns (Avoid)

### AP-1: Cliché Overload
```
BAD:
You are my sunshine in the dark
You light the fire in my heart
You are the star up in the sky

GOOD:
You left your reading glasses on my shelf
I wear them sometimes just to see the world
The way you did before you disappeared
```
- 抽象的な決まり文句 → 具体的な個人的ディテール

### AP-2: Tell, Don't Show
```
BAD:
I'm so sad and lonely tonight
My heart is broken and I cry

GOOD:
The clock reads 3:47 again
Same stain on the ceiling, same empty bed
I count the rings the coffee left
On the table where you used to read
```
- 感情を直接述べる → 状況で感情を想起させる

### AP-3: Over-Stuffed Sections
```
BAD:
[Verse 1]
Line 1
Line 2
Line 3
Line 4
Line 5
Line 6
Line 7
Line 8
Line 9
Line 10

GOOD:
[Verse 1]
Line 1
Line 2
Line 3
Line 4
```
- 1セクション2-6行。長すぎるとボーカルがずれる

### AP-4: Missing Structure Tags
```
BAD:
Walking down the street
Thinking about you
I love you so much
You're my everything

GOOD:
[Verse 1]
Walking down the street at dawn
Thinking of the words you said

[Chorus]
You're the echo in my mind
Coming back like the morning light
```
- タグなし → Suno が構造を誤認する

### AP-5: Repeat Chorus Shortcut
```
BAD:
[Chorus]
(repeat chorus)

GOOD:
[Chorus]
I keep coming back to you
Like the tide keeps finding shore
Every time I think we're through
I just want you even more
```
- 毎回フルテキストで繰り返す

### AP-6: Custom Tags
```
BAD:
[My Epic Section]
[Emotional Breakdown]
[The Good Part]

GOOD:
[Bridge | emotional | soft]
[Build-Up]
[Final Chorus | high energy]
```
- 認識される標準タグのみ使用

### AP-7: Style in Lyrics
```
BAD:
[Verse 1]
*guitar strums softly*
The wind blows through the trees
(play this part slowly)

GOOD:
[Verse 1 | soft | acoustic guitar]
The wind blows through the trees
A whisper in the leaves
```
- 歌詞内にスタイル指示やSFX記述を入れない

### AP-8: Exclamation Overflow
```
BAD:
I won't stop fighting!
I'll never give up!
This is my moment!
I'm taking it all!

GOOD:
I won't stop fighting!
The ground beneath me shakes
But I stand tall, I stand tall
This is where I break through
```
- `!` は控えめに。次行にアグレッシブさが伝播する

### AP-9: Contradictory Style Prompts
```
BAD:
aggressive, calm, fast-tempo, slow-tempo, heavy bass, minimal

GOOD:
aggressive, heavy bass, fast-tempo, distorted guitar, 160 BPM
```
- 矛盾するタグは相殺される

### AP-10: Pipe Overload
```
BAD:
[Verse | soft | whisper | acoustic | sad | rainy | nostalgic | female | low]

GOOD:
[Verse | soft | whisper | acoustic]
```
- パイプスタッキングは2-4要素に抑える

---

## Recommended Patterns

### RP-1: Scene Anchoring
曲全体を1つのシーン/場所に固定する:
```
Theme: "真夜中のプラットフォーム"
→ 全ての描写がこの場所から派生する
→ 雨、電車の音、ベンチ、時刻表
```

### RP-2: Bookend Structure
最初と最後に同じフレーズ/イメージを配置:
```
[Verse 1] 最初の行: "The porch light flickers on"
...
[Outro] 最後の行: "The porch light flickers off"
```

### RP-3: Progressive Revelation
Verseごとに情報を段階的に開示:
```
Verse 1: 何かがおかしい（曖昧なヒント）
Verse 2: 具体的な手がかり
Bridge: 真実の暴露
Final Chorus: 新しい理解でのリフレイン
```

### RP-4: Perspective Shift
ブリッジで視点を切り替える:
```
Verses/Chorus: 1人称 "I"
Bridge: 2人称 "You" または 3人称で自分を俯瞰
```

### RP-5: Dynamic Metadata
セクションごとにメタデータを変化させる:
```
[Verse 1 | soft | whisper]
[Chorus | powerful | full mix]
[Bridge | vulnerable | piano only]
[Final Chorus | belting | high energy]
```

### RP-6: Placeholder-First Workflow
1. まずプレースホルダー歌詞で構造をテスト
2. Sunoで生成してメロディ/リズムを確認
3. フックとコーラスを先に完成させる
4. Verseを最後に仕上げる
5. 6回以上イテレーションする

### RP-7: Title = Hook Alignment
曲タイトルとコーラスのフックを一致させる:
```
Title: "Coming Back to You"
Chorus first line: "I keep coming back to you"
```

### RP-8: Instrumental Breathing Room
歌詞だけでなく、インスト部分で曲に呼吸を与える:
```
[Verse 1]
...
[Instrumental Break]
[Verse 2]
...
```

### RP-9: Emotion Contrast in Bridge
ブリッジで曲全体と対比する感情を入れる:
```
Song mood: energetic, upbeat
Bridge: [Bridge | quiet | reflective | piano only]
→ 静かな瞬間が最後のコーラスをより印象的にする
```

### RP-10: Strategic Ad-libs
アドリブを感情のピークに配置:
```
[Chorus]
I keep coming back to you (ooh)    ← ここ
Like the tide keeps finding shore
Every time I think we're through
I just want you even more (yeah)   ← ここ
```
