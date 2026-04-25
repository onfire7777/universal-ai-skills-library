# OBS Streaming

Purpose: Read this when controlling OBS through WebSocket, designing scenes, routing audio, choosing `RTMP` vs `SRT`, or automating stream start and stop.

## Contents

- [obs-websocket-js v5](#obs-websocket-js-v5)
- [Scene management](#scene-management)
- [Audio routing](#audio-routing)
- [RTMP vs SRT](#rtmp-vs-srt)
- [Streaming automation](#streaming-automation)
- [Recommended OBS settings](#recommended-obs-settings)

## obs-websocket-js v5

### Connection

```typescript
import OBSWebSocket from 'obs-websocket-js';

const obs = new OBSWebSocket();

await obs.connect('ws://localhost:4455', 'password');
await obs.disconnect();

class OBSConnection {
  private obs = new OBSWebSocket();
  private reconnecting = false;

  async connect(): Promise<void> {
    try {
      await this.obs.connect('ws://localhost:4455', process.env.OBS_WS_PASSWORD);
      this.obs.on('ConnectionClosed', () => this.handleDisconnect());
    } catch (err) {
      await this.handleDisconnect();
    }
  }

  private async handleDisconnect(): Promise<void> {
    if (this.reconnecting) return;
    this.reconnecting = true;
    let delay = 1000;
    while (true) {
      try {
        await new Promise(r => setTimeout(r, delay));
        await this.obs.connect('ws://localhost:4455', process.env.OBS_WS_PASSWORD);
        this.reconnecting = false;
        return;
      } catch {
        delay = Math.min(delay * 2, 30000);
      }
    }
  }
}
```

## Scene Management

### Scene Definitions

| Scene | Purpose | Sources | Transition |
|-------|---------|---------|------------|
| Starting | Pre-stream countdown | Countdown image, BGM | Fade (`1s`) |
| Main | Active streaming | Avatar browser source, chat overlay, capture, BGM | Cut |
| BRB | Break state | BRB animation, BGM, optional chat overlay | Fade (`1s`) |
| Ending | Stream close | Credits, CTA, BGM | Fade (`2s`) |
| Emergency | Technical issues | Static image, alert text | Cut (instant) |

### Scene Switching

```typescript
await obs.call('SetCurrentProgramScene', { sceneName: 'Main' });

const { currentProgramSceneName } = await obs.call('GetCurrentProgramScene');
const { scenes } = await obs.call('GetSceneList');
```

### Source Control

```typescript
await obs.call('SetSceneItemEnabled', {
  sceneName: 'Main',
  sceneItemId: chatOverlayId,
  sceneItemEnabled: true,
});

const { sceneItemId } = await obs.call('GetSceneItemId', {
  sceneName: 'Main',
  sourceName: 'ChatOverlay',
});

await obs.call('SetInputVolume', {
  inputName: 'TTS_Audio',
  inputVolumeDb: -3.0,
});

await obs.call('SetInputMute', {
  inputName: 'BGM',
  inputMuted: true,
});
```

Safety rule: keep scene safety checks explicit so active speech cannot be cut by accidental scene or source toggles.

## Audio Routing

### Audio Source Setup

| Source | Type | Purpose | Default volume |
|--------|------|---------|----------------|
| `TTS_Audio` | Media source / audio pipe | Spoken output | `0 dB` |
| `BGM` | Media source | Background music | `-15 dB` |
| `SFX` | Media source | Alerts and transitions | `-10 dB` |
| `Mic` | Audio input capture | Fallback / manual override | `-∞` (muted) |

### TTS Audio Integration

```text
Option 1: Virtual Audio Cable
  TTS process → Virtual Audio Device → OBS Audio Input Capture
  Pros: Simple, low latency
  Cons: Platform-specific setup

Option 2: Media Source with file
  TTS process → Write WAV file → OBS Media Source (file monitoring)
  Pros: Cross-platform, reliable
  Cons: Small file I/O delay

Option 3: Browser Source with Web Audio API
  TTS process → WebSocket → Browser Source → Web Audio API playback
  Pros: Flexible, integrated with avatar
  Cons: Additional WebSocket overhead
```

Recommendation: prefer Option 3 when avatar and audio already share a browser runtime.

### BGM Ducking

```typescript
class AudioDucker {
  private normalBGMVolume = -15;
  private duckedBGMVolume = -25;
  private transitionTime = 300;

  async onTTSStart(): Promise<void> {
    await obs.call('SetInputVolume', {
      inputName: 'BGM',
      inputVolumeDb: this.duckedBGMVolume,
    });
  }

  async onTTSEnd(): Promise<void> {
    await new Promise(r => setTimeout(r, 200));
    await obs.call('SetInputVolume', {
      inputName: 'BGM',
      inputVolumeDb: this.normalBGMVolume,
    });
  }
}
```

## RTMP vs SRT

| Feature | RTMP | SRT |
|---------|------|-----|
| Protocol | TCP-based | UDP-based |
| Latency | `2-5s` | `0.5-2s` |
| Packet-loss handling | TCP retransmit | ARQ + FEC |
| Encryption | RTMPS (TLS) | AES-128/256 built-in |
| Firewall fit | Port `1935` | Configurable port |
| Platform support | YouTube, Twitch | YouTube, partial Twitch ingest |
| OBS support | Native | Native |
| Recommended for | Maximum compatibility | Lower latency / unstable networks |

### OBS Stream Settings

```text
RTMP:
  Server: rtmp://a.rtmp.youtube.com/live2
  Stream Key: {youtube_stream_key}

SRT:
  Server: srt://a.srt.youtube.com:9710
  Stream Key: (embedded in SRT URL as streamid)

  Full SRT URL format:
  srt://a.srt.youtube.com:9710?streamid={stream_key}&latency=2000000
```

## Streaming Automation

### Stream Lifecycle

```typescript
class StreamLifecycle {
  async startStream(): Promise<void> {
    await obs.call('SetCurrentProgramScene', { sceneName: 'Starting' });
    await obs.call('StartStream');
    await obs.call('StartRecord');
    await new Promise(r => setTimeout(r, 60_000));
    await obs.call('SetCurrentProgramScene', { sceneName: 'Main' });
  }

  async endStream(): Promise<void> {
    await obs.call('SetCurrentProgramScene', { sceneName: 'Ending' });
    await new Promise(r => setTimeout(r, 30_000));
    await obs.call('StopStream');
    await obs.call('StopRecord');
  }

  async emergencyStop(): Promise<void> {
    await obs.call('SetCurrentProgramScene', { sceneName: 'Emergency' });
    // Do NOT stop the stream; preserve a recovery path.
  }
}
```

### Health Monitoring

```typescript
const stats = await obs.call('GetStats');
// stats.cpuUsage, stats.memoryUsage, stats.activeFps,
// stats.renderSkippedFrames, stats.outputSkippedFrames

const streamStatus = await obs.call('GetStreamStatus');
// streamStatus.outputActive, streamStatus.outputDuration,
// streamStatus.outputBytes, streamStatus.outputSkippedFrames

obs.on('StreamStateChanged', (data) => {
  if (data.outputState === 'OBS_WEBSOCKET_OUTPUT_RECONNECTING') {
    // Connection unstable — alert
  }
});
```

## Recommended OBS Settings

### Video Settings

| Setting | Recommended | Notes |
|---------|-------------|-------|
| Base resolution | `1920×1080` | Match monitor |
| Output resolution | `1920×1080` | Scale down to `1280×720` if GPU-limited |
| FPS | `30` | Use `60` only if avatar + encoding load allows it |
| Downscale filter | `Lanczos` | Best visual quality for downscaling |

### Output Settings

| Setting | RTMP | SRT |
|---------|------|-----|
| Encoder | `NVENC H.264` or `x264` | `NVENC H.264` |
| Rate control | `CBR` | `CBR` |
| Bitrate | `4500-6000 kbps` | `4500-6000 kbps` |
| Keyframe interval | `2s` | `2s` |
| Profile | `High` | `High` |
| Tune | — | `zerolatency` |

### Browser Source Settings

| Setting | Value | Notes |
|---------|-------|-------|
| Width | `1920` | Match base resolution |
| Height | `1080` | Match base resolution |
| FPS | `30` | Match OBS FPS |
| Custom CSS | `body { background: transparent; }` | Transparent overlay |
| Shutdown when invisible | Disabled | Keep avatar alive |
| Refresh when active | Disabled | Avoid reload during stream |
