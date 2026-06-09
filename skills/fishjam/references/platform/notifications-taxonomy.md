# Server Notifications Taxonomy

18 events surfaced by the server SDKs, grouped by domain. Same set across both transports (WebSocket notifier and webhook). The underlying `ServerMessage` protobuf carries more oneof fields than these — the SDKs intentionally surface only this subset. JS sources: `packages/js-server-sdk/src/notifications.ts` in <https://github.com/fishjam-cloud/js-server-sdk>. Python sources: `fishjam/events/allowed_notifications.py` in <https://github.com/fishjam-cloud/python-server-sdk>. Wire format: `static/api/protobuf/server_notifications.proto` in <https://github.com/fishjam-cloud/documentation>.

## Common fields

- **Room-scoped events** carry **`roomId`** (the room they relate to). **Exception:** the four livestream streamer/viewer events use **`streamId`** (the livestream ID) instead — see the per-event tables below.
- Peer-scoped events carry **`peerId`**.
- Peer-lifecycle events also carry **`peerType`**: `'webrtc' | 'agent' | 'vapi` (mapped from protobuf enum in the JS SDK; raw enum in Python).
- Track-scoped events carry **`track: { id, type, metadata }`** where `type` is `'video' | 'audio'`.

## Room lifecycle (3)

| Event | Trigger | Payload (beyond `roomId`) |
|---|---|---|
| `roomCreated` | Backend successfully calls `createRoom`. | `room` (full Room object). |
| `roomDeleted` | Room is deleted (explicitly via `deleteRoom`, or by Fishjam when empty for long). | — |
| `roomCrashed` | Internal Fishjam fault terminated the room. | — |

## Peer lifecycle (6)

| Event | Trigger | Payload |
|---|---|---|
| `peerAdded` | Your backend called `createPeer` (peer record created; peer has not connected yet). | `peerId`, `peerType`. Metadata is **not** on this event — it arrives via `peerMetadataUpdated`. |
| `peerDeleted` | Your backend called `deletePeer`, or token expired and peer was reaped. | `peerId`, `peerType`. |
| `peerConnected` | Client established a WebRTC connection using its peer token. | `peerId`, `peerType`. |
| `peerDisconnected` | Client dropped the WebRTC connection (may reconnect; if not, eventually `peerDeleted`). | `peerId`, `peerType`. |
| `peerMetadataUpdated` | Peer metadata changed (also fires on initial set after `peerAdded`). | `peerId`, `peerType`, new `metadata`. |
| `peerCrashed` | Fishjam-side fault for that peer. | `peerId`, `peerType`, `reason: string`. |

Pattern: `peerAdded` precedes `peerConnected`. `peerDisconnected` may or may not be followed by `peerDeleted` (depending on whether the client reconnects).

## Track lifecycle (3)

| Event | Trigger | Payload |
|---|---|---|
| `trackAdded` | A peer or component published a new track (camera, mic, screen share, custom source). | exactly one of `peerId` / `componentId`, optional `track` (`{ id, type, metadata }` — may be undefined). |
| `trackRemoved` | Peer or component unpublished a track. | Same shape as `trackAdded`. |
| `trackMetadataUpdated` | Peer or component changed track-level metadata. | Same shape as `trackAdded`, with new metadata on `track`. |

## Livestream — streamer (2)

For `livestream` (and `audio_only_livestream`) rooms only.

| Event | Trigger |
|---|---|
| `streamerConnected` | Streamer client connected with a streamer token. |
| `streamerDisconnected` | Streamer client disconnected. |

## Livestream — viewer (2)

| Event | Trigger |
|---|---|
| `viewerConnected` | A viewer began consuming the livestream. |
| `viewerDisconnected` | A viewer left. |

For public livestreams, viewer events still fire — each WHEP session is tracked.

## Channel lifecycle (2)

| Event | Trigger |
|---|---|
| `channelAdded` | A data channel was opened in the room. Payload: `roomId`, `channelId` (required), and exactly one of `peerId` / `componentId` (the endpoint that owns the channel). |
| `channelRemoved` | A data channel was closed. Same payload shape as `channelAdded`. |

## Intentionally NOT emitted

The protobuf `ServerMessage` union includes other oneof fields that the SDKs deliberately **filter out**. Don't expect to see these:

- `authenticated`, `authRequest`, `subscribeRequest`, `subscribeResponse` — handshake / request-side messages, never inbound notifications.
- `trackForwarding`, `trackForwardingRemoved` — feature exists but no consumer demand to surface.
- `notificationBatch` — webhook-only transport wrapper, the WS notifier never receives it.
- `streamConnected`, `streamDisconnected`, `hlsPlayable`, `hlsUploaded`, `hlsUploadCrashed`, `componentCrashed` — deprecated.

JS source for this partition: the `ignoredEventsList` constant in `notifications.ts`. Python mirrors it.

## Idempotency and ordering

- Events are delivered **in order per room**, but **not** ordered across rooms.
- Fishjam retries 5xx and transport errors with randomized exponential backoff (up to 8 attempts, ~30 s total budget). 4xx is not retried. Design handlers to be idempotent — repeats are expected. VAD notifications use a separate "lossy" path (500 ms timeout, no retry). Treat each event as a hint, not a transactional fact: re-derive state from `GET /room/{id}` when accuracy matters.
- The WebSocket notifier has no replay buffer; restarting the process means losing the events that fired while you were down. Wrap critical state changes (e.g. billing) in a backend-side `getRoom` cross-check.

## Same set, both transports

Whether you receive an event via `FishjamWSNotifier.on('peerConnected', …)` or by decoding a webhook body, you get the **same event set with the same semantic content** — but field representations differ by SDK and transport:

- JS WS notifier: `mapNotification` in `notifications.ts` decodes the protobuf and maps `peerType` / `track.type` enums to friendly strings (`'webrtc'`, `'video'`, ...).
- JS webhook receiver: you call `ServerMessage.decode(...)` yourself and get the **raw protobuf shape** — `peerType` is the enum int (`ServerMessage_PeerType.PEER_TYPE_WEBRTC`), `track.type` is the `TrackType` enum, and `notificationBatch` may appear if `batchWebhookNotifications` is enabled on the room. `mapNotification` is not re-exported.
- Python (both notifier and `receive_binary`): typed `betterproto` message with raw enum ints — no string mapping in either transport.

## Source

- `packages/js-server-sdk/src/notifications.ts` in <https://github.com/fishjam-cloud/js-server-sdk>
- `fishjam/events/allowed_notifications.py` in <https://github.com/fishjam-cloud/python-server-sdk>
- `static/api/protobuf/server_notifications.proto` in <https://github.com/fishjam-cloud/documentation>
- <https://fishjam.swmansion.com/docs/how-to/backend/server-setup>
