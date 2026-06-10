# Room Types

Fishjam exposes four room types — `conference`, `audio_only`, `livestream`, and `audio_only_livestream`. The room type is set when the backend calls `createRoom` and **cannot be changed** afterwards — pick before you create.

## At a glance

| Type | Direction | Cost vs conference | Max tracks per peer | Use case |
|---|---|---|---|---|
| `conference` (default) | many ↔ many bidirectional | base (`$1` per 1000 participant-minutes) | many (camera + screen + mic + custom) | Meetings, classrooms, interactive webinars |
| `audio_only` | many ↔ many, audio only | **75% off** | 1 audio (video silently dropped) | Voice rooms, audio podcasts, town halls |
| `livestream` | one streamer → many viewers | **20% off** | 1 video + 1 audio (streamer); 0 (viewers) | Live events, broadcasts, sports |
| `audio_only_livestream` | one streamer → many viewers, audio only | **75% off livestream** | 1 audio (streamer); 0 (viewers) | Audio-only livestreaming, large radio-style broadcasts |

Pricing details: <https://fishjam.io/#pricing>.

## Conference rooms

The default. All client SDK features work as documented:

- Multiple participants, all bidirectional.
- Camera + microphone + screen share + custom sources, simultaneously.
- Track-level metadata, simulcast, data channels, VAD, etc.

Use this unless you have a specific reason for one of the alternatives.

```ts
const room = await fishjamClient.createRoom();
// or, explicit:
const room = await fishjamClient.createRoom({ roomType: 'conference' });
```

```python
room = fishjam_client.create_room()  # default = conference
```

## Audio-only rooms

`roomType: 'audio_only'`. Same semantics as conference. The client SDK refuses to add video tracks to an `audio_only` room (throws `TrackTypeError`); the React/React Native hooks downgrade this to a console warning. As a defensive fallback, if a video track somehow reached the server, it is dropped. Useful when you want to lock down the call to voice only and benefit from the 75% discount.

```ts
const room = await fishjamClient.createRoom({ roomType: 'audio_only' });
```

```python
from fishjam import RoomOptions
room = fishjam_client.create_room(RoomOptions(room_type='audio_only'))
```

Cost example: 2 peers × 30 minutes = 60 participant-minutes = `$0.06` in conference, `$0.015` in audio-only.

## Livestream rooms

`roomType: 'livestream'`. One streamer publishes, unlimited viewers consume. Backend mints two separate token types:

- **Streamer token** — `createLivestreamStreamerToken(roomId)` / `create_livestream_streamer_token(room_id)`. Hand to the publishing peer's client.
- **Viewer token** — `createLivestreamViewerToken(roomId)` / `create_livestream_viewer_token(room_id)`. Hand to each viewer's client (only required if room is private).

### Public vs private livestreams

```ts
// Default — viewers need a viewer token
const privateRoom = await fishjamClient.createRoom({ roomType: 'livestream' });

// Public — viewers connect using just the room ID
const publicRoom = await fishjamClient.createRoom({ roomType: 'livestream', public: true });
```

Use `private` for paywalled / authenticated viewers. Use `public` for open broadcasts. Public livestreams still bill per viewer, so don't expose them if you need to gate access.

### Constraints

- **Exactly one video track and one audio track per streamer.** Additional tracks are dropped.
- Viewer-side rendering uses [WHEP](https://blog.swmansion.com/building-interactive-streaming-apps-webrtc-whip-whep-explained-d38f4825ec90); any WHEP-compatible player works.

### `audio_only_livestream`

`roomType: 'audio_only_livestream'`. Livestream variant with no video. 75% off the livestream baseline price (equivalently, 80% off conference) — the cheapest mode for one-to-many audio.

## Video codecs

Conference and livestream rooms use H.264 by default. Override per room via `createRoom({ videoCodec: 'h264' | 'vp8' })`.

- **H.264** — hardware-accelerated on most devices, better perf. Default.
- **VP8** — software-only, works everywhere including Android emulators. Default when using the Sandbox API (specifically to support Android emulators during development).

Use H.264 in production unless you have a specific compatibility need. Use VP8 if you're testing on Android emulators outside the Sandbox.

## Decision matrix

| Use case | Type | Why |
|---|---|---|
| Video meetings | `conference` | Multiple senders, bidirectional, full feature set |
| Webinars with Q&A | `conference` | Audience can speak when invited |
| Voice rooms / clubs | `audio_only` | Cheapest, ignores video |
| Live podcasts (multi-speaker) | `audio_only` | Cheapest, voice-only |
| Sports/event streaming | `livestream` | One-to-many, scales cheaply |
| Internal company broadcasts | `livestream` (private) | One presenter, gated viewers |
| Talk-radio-style broadcast | `audio_only_livestream` | Cheapest broadcast mode |

## Source

- <https://fishjam.swmansion.com/docs/explanation/rooms>
- <https://fishjam.swmansion.com/docs/explanation/livestreams>
- <https://fishjam.swmansion.com/docs/how-to/backend/audio-only-calls>
- OpenAPI `RoomType` enum: `fishjam-server-openapi.yaml` in the `documentation` repo
