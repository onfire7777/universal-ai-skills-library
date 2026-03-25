# Lip Sync and Expression

Purpose: Read this when converting Japanese phonemes into avatar mouth movement, extracting timing from `VOICEVOX`, or compositing lip sync with emotion layers.

## Contents

- [Japanese phoneme to viseme mapping](#japanese-phoneme-to-viseme-mapping)
- [VOICEVOX phoneme timing extraction](#voicevox-phoneme-timing-extraction)
- [Lip sync playback engine](#lip-sync-playback-engine)
- [Live2D application](#live2d-application)
- [Emotion transition](#emotion-transition)
- [Compositing rule](#compositing-rule)

## Japanese Phoneme to Viseme Mapping

Japanese lip sync is vowel-dominant, so the vowel viseme is usually the main control signal.

### Vowel Viseme Table

| Vowel | Phoneme | Viseme | Live2D MouthOpenY | Live2D MouthForm | VRM BlendShape |
|-------|---------|--------|-------------------|------------------|----------------|
| あ (a) | `a` | Open wide | `0.8` | `0.0` | `aa: 1.0` |
| い (i) | `i` | Wide narrow | `0.2` | `0.6` | `ih: 1.0` |
| う (u) | `u` | Rounded small | `0.3` | `-0.3` | `ou: 1.0` |
| え (e) | `e` | Half open | `0.5` | `0.3` | `ee: 1.0` |
| お (o) | `o` | Rounded open | `0.6` | `-0.5` | `oh: 1.0` |
| ん (N) | `N` | Closed | `0.0` | `0.0` | all `0` |
| 促音 (cl) | `cl` | Closed (pause) | `0.0` | `0.0` | all `0` |
| 無音 (pau) | `pau` | Neutral | `0.0` | `0.0` | all `0` |

### Consonant Influence

| Consonant group | Effect | MouthOpenY adjustment | Example |
|-----------------|--------|-----------------------|---------|
| `k`, `g` | Back-tongue | `-0.05` | か, が |
| `s`, `z`, `sh` | Dental fricative | `-0.1`, `MouthForm +0.2` | さ, ざ, しゃ |
| `t`, `d`, `n` | Alveolar | `-0.05` | た, だ, な |
| `h`, `f` | Lip / glottal | `-0.05` | は, ふ |
| `m`, `b`, `p` | Bilabial | `-0.2` (momentary close) | ま, ば, ぱ |
| `r` | Tap | `0` | ら |
| `y` | Semivowel | `0` | や |
| `w` | Semivowel | `-0.1`, `MouthForm -0.2` | わ |

## VOICEVOX Phoneme Timing Extraction

Use the `audio_query` response to derive lip-sync timing before synthesis starts.

### `audio_query` Response Structure

```json
{
  "accent_phrases": [
    {
      "moras": [
        {
          "text": "こ",
          "consonant": "k",
          "consonant_length": 0.065,
          "vowel": "o",
          "vowel_length": 0.120,
          "pitch": 5.5
        },
        {
          "text": "ん",
          "consonant": null,
          "consonant_length": null,
          "vowel": "N",
          "vowel_length": 0.100,
          "pitch": 5.2
        },
        {
          "text": "に",
          "consonant": "n",
          "consonant_length": 0.055,
          "vowel": "i",
          "vowel_length": 0.110,
          "pitch": 5.0
        },
        {
          "text": "ち",
          "consonant": "ch",
          "consonant_length": 0.070,
          "vowel": "i",
          "vowel_length": 0.100,
          "pitch": 4.8
        },
        {
          "text": "は",
          "consonant": "w",
          "consonant_length": 0.050,
          "vowel": "a",
          "vowel_length": 0.150,
          "pitch": 4.5
        }
      ],
      "accent": 3,
      "pause_mora": null
    }
  ],
  "speedScale": 1.0,
  "pitchScale": 0.0,
  "intonationScale": 1.0,
  "volumeScale": 1.0
}
```

### Viseme Timeline Generation

```typescript
interface VisemeKeyframe {
  time: number;       // ms from audio start
  viseme: string;     // 'a' | 'i' | 'u' | 'e' | 'o' | 'N' | 'cl' | 'pau'
  duration: number;   // ms
  intensity: number;  // 0.0 - 1.0
}

function generateVisemeTimeline(query: VOICEVOXQuery): VisemeKeyframe[] {
  const keyframes: VisemeKeyframe[] = [];
  let currentTime = 0;

  for (const phrase of query.accent_phrases) {
    for (const mora of phrase.moras) {
      if (mora.consonant && mora.consonant_length) {
        const duration = mora.consonant_length * 1000 / query.speedScale;
        keyframes.push({
          time: currentTime,
          viseme: getConsonantViseme(mora.consonant, mora.vowel),
          duration,
          intensity: 0.5,
        });
        currentTime += duration;
      }

      const vowelDuration = mora.vowel_length * 1000 / query.speedScale;
      keyframes.push({
        time: currentTime,
        viseme: mora.vowel,
        duration: vowelDuration,
        intensity: 1.0,
      });
      currentTime += vowelDuration;
    }

    if (phrase.pause_mora) {
      const pauseDuration = phrase.pause_mora.vowel_length * 1000 / query.speedScale;
      keyframes.push({
        time: currentTime,
        viseme: 'pau',
        duration: pauseDuration,
        intensity: 0,
      });
      currentTime += pauseDuration;
    }
  }

  return keyframes;
}

function getConsonantViseme(consonant: string, followingVowel: string): string {
  if (['m', 'b', 'p'].includes(consonant)) return 'cl';
  return followingVowel;
}
```

## Lip Sync Playback Engine

```typescript
class LipSyncPlayer {
  private timeline: VisemeKeyframe[] = [];
  private startTime = 0;
  private currentIndex = 0;

  load(timeline: VisemeKeyframe[]): void {
    this.timeline = timeline;
    this.currentIndex = 0;
  }

  start(): void {
    this.startTime = performance.now();
    this.currentIndex = 0;
  }

  getCurrentViseme(): { viseme: string; intensity: number } | null {
    if (this.timeline.length === 0) return null;

    const elapsed = performance.now() - this.startTime;

    while (
      this.currentIndex < this.timeline.length - 1 &&
      this.timeline[this.currentIndex + 1].time <= elapsed
    ) {
      this.currentIndex++;
    }

    const kf = this.timeline[this.currentIndex];
    if (!kf || elapsed > kf.time + kf.duration) {
      return { viseme: 'pau', intensity: 0 };
    }

    const progress = (elapsed - kf.time) / kf.duration;
    const easedIntensity = kf.intensity * this.easeInOut(progress);

    return { viseme: kf.viseme, intensity: easedIntensity };
  }

  private easeInOut(t: number): number {
    if (t < 0.2) return t / 0.2;
    if (t > 0.7) return (1 - t) / 0.3;
    return 1.0;
  }
}
```

## Live2D Application

```typescript
function applyVisemeToLive2D(
  model: Live2DModel,
  viseme: { viseme: string; intensity: number }
): void {
  const VISEME_MAP: Record<string, { openY: number; form: number }> = {
    a:   { openY: 0.8, form: 0.0 },
    i:   { openY: 0.2, form: 0.6 },
    u:   { openY: 0.3, form: -0.3 },
    e:   { openY: 0.5, form: 0.3 },
    o:   { openY: 0.6, form: -0.5 },
    N:   { openY: 0.0, form: 0.0 },
    cl:  { openY: 0.0, form: 0.0 },
    pau: { openY: 0.0, form: 0.0 },
  };

  const params = VISEME_MAP[viseme.viseme] ?? VISEME_MAP.pau;
  model.setParameterValueById('ParamMouthOpenY', params.openY * viseme.intensity);
  model.setParameterValueById('ParamMouthForm', params.form * viseme.intensity);
}
```

## Emotion Transition

### Emotion Detection Sources

| Source | Method | Latency | Accuracy |
|--------|--------|---------|----------|
| LLM response metadata | Emotion tag in prompt | `0ms` | High |
| Chat sentiment | Keyword / emoji analysis | `< 10ms` | Medium |
| Separate LLM analysis | Dedicated classifier | `200-500ms` | High |

### Transition Algorithm

```typescript
interface EmotionState {
  emotion: string;
  intensity: number;
  timestamp: number;
}

class EmotionController {
  private currentEmotion: EmotionState = { emotion: 'neutral', intensity: 1.0, timestamp: 0 };
  private transitionDuration = 500;

  setEmotion(emotion: string, intensity: number): void {
    this.currentEmotion = { emotion, intensity, timestamp: performance.now() };
  }

  getExpressionParams(now: number): Record<string, number> {
    const elapsed = now - this.currentEmotion.timestamp;
    const progress = Math.min(elapsed / this.transitionDuration, 1.0);
    const eased = this.easeOutQuad(progress);

    const targetParams =
      EMOTION_EXPRESSION_MAP[this.currentEmotion.emotion] ?? EMOTION_EXPRESSION_MAP.neutral;
    const result: Record<string, number> = {};

    for (const [key, value] of Object.entries(targetParams)) {
      result[key] = value * this.currentEmotion.intensity * eased;
    }

    return result;
  }

  private easeOutQuad(t: number): number {
    return 1 - (1 - t) * (1 - t);
  }
}

const EMOTION_EXPRESSION_MAP: Record<string, Record<string, number>> = {
  neutral:   { happy: 0, angry: 0, sad: 0, surprised: 0, relaxed: 0.3 },
  joy:       { happy: 0.8, angry: 0, sad: 0, surprised: 0, relaxed: 0.2 },
  sad:       { happy: 0, angry: 0, sad: 0.7, surprised: 0, relaxed: 0 },
  angry:     { happy: 0, angry: 0.8, sad: 0, surprised: 0, relaxed: 0 },
  surprised: { happy: 0, angry: 0, sad: 0, surprised: 0.9, relaxed: 0 },
  thinking:  { happy: 0, angry: 0, sad: 0, surprised: 0, relaxed: 0.5 },
};
```

## Compositing Rule

Compose animation layers in this order:

```text
Layer 1: Idle animation
  ↓ blend
Layer 2: Emotion expression
  ↓ blend
Layer 3: Lip sync
  ↓
Final avatar parameters
```

Rule: lip-sync mouth parameters override emotion mouth parameters. Eyes, brows, and other non-mouth facial controls stay driven by the emotion layer.
