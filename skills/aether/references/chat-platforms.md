# Chat Platforms

Purpose: Read this when selecting `YouTube` / `Twitch`, implementing chat ingestion, normalizing platform events, or defining filtering and command handling.

## Contents

- [Platform comparison](#platform-comparison)
- [YouTube Live Chat API](#youtube-live-chat-api)
- [Twitch chat integration](#twitch-chat-integration)
- [Unified message format](#unified-message-format)
- [Filtering and command handling](#filtering-and-command-handling)

## Platform Comparison

| Feature | YouTube Live | Twitch |
|---------|-------------|--------|
| Chat API | `liveChatMessages` (REST polling) | IRC + EventSub (WebSocket) |
| Real-time path | Polling (`5-10s` interval) | WebSocket (near-instant) |
| Auth | OAuth 2.0 (Google) | OAuth 2.0 (Twitch) |
| Rate limit | `10,000` units/day (quota) | `20 msg / 30s` (moderator baseline) |
| Paid interactions | Super Chat / Super Stickers | Bits / Cheers |
| Commands | No native support | Native `!` commands |
| Emotes | Limited | Extensive (`BTTV`, `FFZ`, `7TV`) |
| Moderation | YouTube moderation APIs | AutoMod + custom bots |

Use `YouTube` when platform reach matters more than ingestion latency. Use `Twitch` when command-driven real-time interaction matters more.

## YouTube Live Chat API

### Polling Architecture

```text
┌────────────┐     GET /liveChatMessages      ┌──────────────┐
│   Aether   │ ──────────────────────────────▶ │ YouTube API  │
│  Chat      │ ◀────────────────────────────── │              │
│  Listener  │     { messages, nextPageToken,  │              │
│            │       pollingIntervalMillis }   │              │
└────────────┘                                  └──────────────┘
     │
     │ pollingIntervalMillis (typically 5000-10000ms)
     ▼
  Next poll with nextPageToken
```

### Key Endpoints

```text
# Get live broadcast
GET https://www.googleapis.com/youtube/v3/liveBroadcasts
  ?part=snippet,status
  &broadcastStatus=active
  &mine=true

# Get live chat ID from broadcast
→ response.items[0].snippet.liveChatId

# Poll chat messages
GET https://www.googleapis.com/youtube/v3/liveChat/messages
  ?liveChatId={liveChatId}
  &part=snippet,authorDetails
  &pageToken={nextPageToken}

# Send chat message
POST https://www.googleapis.com/youtube/v3/liveChat/messages
  ?part=snippet
  Body: { snippet: { liveChatId, type: "textMessageEvent", textMessageDetails: { messageText } } }
```

### Message Types

| Type | Field | Meaning |
|------|-------|---------|
| `textMessageEvent` | `textMessageDetails.messageText` | Standard message |
| `superChatEvent` | `superChatDetails.amountMicros` | Super Chat |
| `superStickerEvent` | `superStickerDetails` | Super Sticker |
| `membershipGiftingEvent` | — | Membership gift |
| `newSponsorEvent` | — | New member |

### OAuth 2.0 Flow

```text
1. Authorization URL:
   https://accounts.google.com/o/oauth2/v2/auth
   ?client_id={CLIENT_ID}
   &redirect_uri=http://localhost:3000/callback
   &response_type=code
   &scope=https://www.googleapis.com/auth/youtube
          https://www.googleapis.com/auth/youtube.readonly

2. Token exchange: POST https://oauth2.googleapis.com/token

3. Store refresh_token securely (encrypted at rest)

4. Auto-refresh: access_token expires in 3600s → use refresh_token
```

### Quota Management

YouTube Data API v3 uses a daily quota budget (default `10,000` units).

| Operation | Cost |
|-----------|------|
| `liveChatMessages.list` | `5` units |
| `liveChat/messages.insert` | `50` units |
| `liveBroadcasts.list` | `100` units |

Budget example:

- Poll every `5s`: `5` units × `12/min` × `60 min` = `3,600 units/hour`
- `2`-hour stream: about `7,200` units for polling only
- Remaining budget: about `2,800` units for chat sends and control calls

Optimization rule: increase the polling interval during low-activity windows, for example `10s` when fewer than `1` message arrives per poll.

## Twitch Chat Integration

### IRC + EventSub Architecture

```text
┌────────────┐     WebSocket (IRC)            ┌──────────────┐
│   Aether   │ ──────────────────────────────▶ │   Twitch     │
│  Chat      │ ◀────────────────────────────── │   IRC        │
│  Listener  │     PRIVMSG #channel :message   │              │
└────────────┘                                  └──────────────┘
     │
     │ Also connects to:
     ▼
┌────────────┐     WebSocket (EventSub)       ┌──────────────┐
│   Aether   │ ──────────────────────────────▶ │   Twitch     │
│  Event     │ ◀────────────────────────────── │   EventSub   │
│  Listener  │     Subscription events         │              │
└────────────┘                                  └──────────────┘
```

### IRC Connection

```text
WebSocket URL: wss://irc-ws.chat.twitch.tv:443

Authentication:
  PASS oauth:{access_token}
  NICK {bot_username}

Join channel:
  JOIN #{channel_name}

Request capabilities (for tags, commands):
  CAP REQ :twitch.tv/tags twitch.tv/commands twitch.tv/membership

Incoming message format:
  @badge-info=;badges=moderator/1;color=#FF0000;display-name=User;emotes=;
  first-msg=0;id=msg-id;mod=1;subscriber=0;turbo=0;user-id=12345;
  user-type=mod :user!user@user.tmi.twitch.tv PRIVMSG #channel :Hello!
```

### EventSub

```text
WebSocket URL: wss://eventsub.wss.twitch.tv/ws

Subscription types:
  channel.cheer
  channel.subscribe
  channel.raid
  channel.channel_points_custom_reward_redemption.add
```

### Twitch OAuth 2.0

```text
1. Authorization:
   https://id.twitch.tv/oauth2/authorize
   ?client_id={CLIENT_ID}
   &redirect_uri=http://localhost:3000/callback
   &response_type=code
   &scope=chat:read+chat:edit+channel:read:subscriptions+bits:read

2. Token exchange: POST https://id.twitch.tv/oauth2/token

3. Validate: GET https://id.twitch.tv/oauth2/validate
   (Validate periodically; Twitch may revoke inactive or invalid tokens)
```

## Unified Message Format

Normalize platform-specific events before they enter the response pipeline.

```typescript
interface UnifiedChatMessage {
  id: string;                    // Platform message ID
  platform: 'youtube' | 'twitch';
  timestamp: number;             // Unix ms

  author: {
    id: string;                  // Platform user ID
    name: string;                // Display name
    isModerator: boolean;
    isSubscriber: boolean;       // YouTube member / Twitch sub
    isOwner: boolean;            // Channel owner
    badges: string[];            // Raw badge strings
  };

  content: {
    text: string;                // Raw message text
    isCommand: boolean;          // Starts with ! or /
    command?: string;            // Parsed command name
    commandArgs?: string[];      // Parsed command args
  };

  monetary?: {
    amount: number;              // Normalized to JPY
    currency: string;            // Original currency code
    rawAmount: number;           // Original amount
    tier?: string;               // Superchat tier color / bits badge
  };

  meta: {
    type: 'message' | 'superchat' | 'bits' | 'subscription' | 'raid' | 'membership';
    priority: number;            // 1-4 (computed from type + monetary)
    isFirstMessage: boolean;     // First-time chatter
    replyTo?: string;            // Parent message ID (if reply)
  };
}
```

### Platform Normalizer

```typescript
function normalizeYouTubeMessage(raw: YouTubeLiveChatMessage): UnifiedChatMessage {
  const isSuperChat = raw.snippet.type === 'superChatEvent';
  return {
    id: raw.id,
    platform: 'youtube',
    timestamp: new Date(raw.snippet.publishedAt).getTime(),
    author: {
      id: raw.authorDetails.channelId,
      name: raw.authorDetails.displayName,
      isModerator: raw.authorDetails.isChatModerator,
      isSubscriber: raw.authorDetails.isChatSponsor,
      isOwner: raw.authorDetails.isChatOwner,
      badges: [],
    },
    content: {
      text: isSuperChat
        ? raw.snippet.superChatDetails?.userComment ?? ''
        : raw.snippet.textMessageDetails?.messageText ?? '',
      isCommand: false,
    },
    monetary: isSuperChat ? {
      amount: Number(raw.snippet.superChatDetails.amountMicros) / 1_000_000,
      currency: raw.snippet.superChatDetails.currency,
      rawAmount: Number(raw.snippet.superChatDetails.amountMicros) / 1_000_000,
      tier: raw.snippet.superChatDetails.tier,
    } : undefined,
    meta: {
      type: isSuperChat ? 'superchat' : 'message',
      priority: isSuperChat ? 1 : 4,
      isFirstMessage: false,
    },
  };
}
```

## Filtering and Command Handling

### Filter Pipeline

```text
Raw Message
  → Spam filter (repeated chars, known spam patterns)
  → Safety filter (toxic content, personal info)
  → Command parser (! prefix detection)
  → Rate limit (per-user cooldown)
  → Unified message output
```

### Command Recognition

| Command | Action | Priority |
|---------|--------|----------|
| `!ask [question]` | Force an LLM answer | `2` |
| `!song [title]` | Request BGM change | `3` |
| `!emotion [name]` | Trigger avatar expression | `3` |
| `!skip` | Skip current TTS playback | `2` (moderator only) |
| `!queue` | Show TTS queue status | `4` |
| `!stats` | Show stream stats | `4` |

### Safety Filtering

```typescript
interface SafetyFilter {
  /** Check message safety. Returns null if safe, reason string if filtered */
  check(message: string): string | null;
}

// Layers:
// 1. Regex blocklist (slurs, known patterns)
// 2. Personal info detection (phone, email, address patterns)
// 3. URL filtering (allow only whitelisted domains)
// 4. LLM-based content classification (optional, adds latency)
```
