# TTS Engines

Purpose: Read this when choosing a TTS engine, implementing the `TTSAdapter` boundary, tuning voice parameters, or managing the audio queue.

## Contents

- [Engine comparison](#engine-comparison)
- [Use-case recommendations](#use-case-recommendations)
- [TTSAdapter contract](#ttsadapter-contract)
- [VOICEVOX example](#voicevox-example)
- [Audio queue management](#audio-queue-management)
- [Parameter tuning](#parameter-tuning)

## Engine Comparison

| Engine | Type | Latency | Quality | Japanese | GPU required | Cost | Notes |
|--------|------|---------|---------|----------|--------------|------|-------|
| `VOICEVOX` | Local | `200-800ms` | High | Native | Optional | OSS | Strong lip-sync support |
| `Style-Bert-VITS2` | Local | `500-1500ms` | Very high | Native | Recommended | OSS | Best quality and custom voice training |
| `COEIROINK` | Local | `200-600ms` | High | Native | Optional | Free | Lightweight |
| `NIJIVOICE` | Cloud | `300-1000ms` | Very high | Native | No | Paid | Good when GPU is unavailable |
| `VOICEVOX Nemo` | Local | `200-800ms` | High | Native | Optional | OSS | Same API family as VOICEVOX |

## Use-Case Recommendations

| Use case | Recommended engine | Why |
|----------|--------------------|-----|
| v1 / fastest setup | `VOICEVOX` | Lowest setup friction and phoneme timing support |
| Highest naturalness | `Style-Bert-VITS2` | Strongest voice quality |
| Low-resource local run | `COEIROINK` | Lightweight and quick |
| No GPU available | `NIJIVOICE` | Offloads compute |
| Custom trained voice | `Style-Bert-VITS2` | Fine-tuning support |

## TTSAdapter Contract

Keep the engine boundary behind a single interface so the pipeline can switch engines without architectural rewrites.

```typescript
interface TTSAdapter {
  /** Engine name */
  readonly name: string;

  /** Check availability */
  isAvailable(): Promise<boolean>;

  /** Synthesize speech */
  synthesize(params: TTSSynthesizeParams): Promise<TTSResult>;

  /** Get phoneme timing for lip sync when the engine supports it */
  getPhonemeTimings?(text: string, speakerId: number): Promise<PhonemeTiming[]>;

  /** List available speakers */
  getSpeakers(): Promise<TTSSpeaker[]>;

  /** Shut down or release the engine */
  dispose(): Promise<void>;
}

interface TTSSynthesizeParams {
  text: string;
  speakerId: number;
  speed?: number;
  pitch?: number;
  intonation?: number;
  volume?: number;
}

interface TTSResult {
  audio: Buffer;
  format: 'wav' | 'mp3';
  duration: number;
  phonemes?: PhonemeTiming[];
  synthesisTime: number;
}
```

## VOICEVOX Example

Use this when phoneme timing is required for lip sync:

```typescript
class VOICEVOXAdapter implements TTSAdapter {
  readonly name = 'voicevox';
  private baseUrl: string;

  constructor(baseUrl = 'http://localhost:50021') {
    this.baseUrl = baseUrl;
  }

  async isAvailable(): Promise<boolean> {
    try {
      const res = await fetch(`${this.baseUrl}/version`);
      return res.ok;
    } catch {
      return false;
    }
  }

  async synthesize(params: TTSSynthesizeParams): Promise<TTSResult> {
    const start = Date.now();

    const query = await fetch(
      `${this.baseUrl}/audio_query?text=${encodeURIComponent(params.text)}&speaker=${params.speakerId}`,
      { method: 'POST' }
    ).then(r => r.json());

    if (params.speed) query.speedScale = params.speed;
    if (params.pitch) query.pitchScale = params.pitch;
    if (params.intonation) query.intonationScale = params.intonation;
    if (params.volume) query.volumeScale = params.volume;

    const audioRes = await fetch(
      `${this.baseUrl}/synthesis?speaker=${params.speakerId}`,
      { method: 'POST', headers: { 'Content-Type': 'application/json' }, body: JSON.stringify(query) }
    );
    const audio = Buffer.from(await audioRes.arrayBuffer());

    return {
      audio,
      format: 'wav',
      duration: this.calculateDuration(query),
      phonemes: this.extractPhonemes(query),
      synthesisTime: Date.now() - start,
    };
  }

  async getPhonemeTimings(text: string, speakerId: number): Promise<PhonemeTiming[]> {
    const query = await fetch(
      `${this.baseUrl}/audio_query?text=${encodeURIComponent(text)}&speaker=${speakerId}`,
      { method: 'POST' }
    ).then(r => r.json());
    return this.extractPhonemes(query);
  }

  async getSpeakers(): Promise<TTSSpeaker[]> {
    return fetch(`${this.baseUrl}/speakers`).then(r => r.json());
  }

  async dispose(): Promise<void> {
    // VOICEVOX runs as a separate process
  }

  private extractPhonemes(query: any): PhonemeTiming[] {
    const phonemes: PhonemeTiming[] = [];
    let currentTime = 0;
    for (const phrase of query.accent_phrases) {
      for (const mora of phrase.moras) {
        if (mora.consonant) {
          const duration = mora.consonant_length * 1000;
          phonemes.push({ phoneme: mora.consonant, startTime: currentTime, endTime: currentTime + duration });
          currentTime += duration;
        }
        if (mora.vowel) {
          const duration = mora.vowel_length * 1000;
          phonemes.push({ phoneme: mora.vowel, startTime: currentTime, endTime: currentTime + duration });
          currentTime += duration;
        }
      }
      if (phrase.pause_mora) {
        const duration = phrase.pause_mora.vowel_length * 1000;
        phonemes.push({ phoneme: 'pau', startTime: currentTime, endTime: currentTime + duration });
        currentTime += duration;
      }
    }
    return phonemes;
  }

  private calculateDuration(query: any): number {
    let total = 0;
    for (const phrase of query.accent_phrases) {
      for (const mora of phrase.moras) {
        total += (mora.consonant_length ?? 0) + mora.vowel_length;
      }
      if (phrase.pause_mora) total += phrase.pause_mora.vowel_length;
    }
    return total * 1000;
  }
}
```

## Audio Queue Management

Real-time chat produces bursty audio work. The queue must stay bounded.

```typescript
interface AudioQueueItem {
  id: string;
  text: string;
  priority: number;     // 1=highest (superchat), 4=lowest (regular)
  timestamp: number;
  speakerId: number;
  emotion?: string;
}

class AudioQueue {
  private queue: AudioQueueItem[] = [];
  private maxSize = 10;

  enqueue(item: AudioQueueItem): boolean {
    if (this.queue.some(q => q.text === item.text && Date.now() - q.timestamp < 5000)) {
      return false;
    }

    if (this.queue.length >= this.maxSize) {
      const lowestPriority = Math.max(...this.queue.map(q => q.priority));
      if (item.priority >= lowestPriority) return false;
      const dropIndex = this.queue.findLastIndex(q => q.priority === lowestPriority);
      this.queue.splice(dropIndex, 1);
    }

    const insertIndex = this.queue.findIndex(q => q.priority > item.priority);
    if (insertIndex === -1) this.queue.push(item);
    else this.queue.splice(insertIndex, 0, item);
    return true;
  }

  dequeue(): AudioQueueItem | null {
    return this.queue.shift() ?? null;
  }
}
```

### Playback Rule

Use double buffering so the next sentence is synthesized while the current sentence is still playing.

## Parameter Tuning

| Character type | Speed | Pitch | Intonation | Volume |
|----------------|-------|-------|------------|--------|
| Energetic | `1.1-1.2` | `+0.05` | `1.3-1.5` | `1.1` |
| Gentle | `0.85-0.95` | `-0.03` | `0.8-1.0` | `0.9` |
| Cool | `0.95-1.0` | `-0.05` | `0.7-0.9` | `1.0` |
| Tsundere | `1.0-1.1` | `+0.03` | `1.2-1.5` | `1.0-1.2` |
| Calm | `0.9-1.0` | `-0.02` | `0.9-1.1` | `0.95` |

### Emotion Overrides

| Emotion | Speed Δ | Pitch Δ | Intonation Δ | Volume Δ |
|---------|---------|---------|--------------|----------|
| Joy / excited | `+0.1` | `+0.05` | `+0.3` | `+0.1` |
| Sad | `-0.15` | `-0.05` | `-0.2` | `-0.1` |
| Angry | `+0.05` | `+0.03` | `+0.4` | `+0.2` |
| Surprised | `+0.15` | `+0.08` | `+0.5` | `+0.15` |
| Thinking | `-0.1` | `0` | `-0.1` | `-0.05` |
