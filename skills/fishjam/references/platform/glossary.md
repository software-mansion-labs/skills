# Fishjam Glossary

Mirrors the canonical definitions at <https://fishjam.swmansion.com/docs/explanation/glossary>. Linked from every Fishjam SDK skill.

## Room

A collection of peers that exchange audio and video. Identified by a `room.id` (server-assigned). Created by the backend via the server SDK or REST API. Four room types — `conference` (default), `audio_only` (75% cheaper, video tracks silently dropped), `livestream` (one streamer, many viewers, single video + single audio track), `audio_only_livestream` (audio-only livestream variant). Optional `webhookUrl` on creation routes server notifications for that room to your endpoint.

## Peer

A user or program connected to a room. Each peer has:

- An `id` (server-assigned).
- A `peerType`: `webrtc` (client SDK), `agent` (server-side via `createAgent` / `create_agent`), `vapi` (managed voice agent).
- Optional `metadata` (set at peer creation by the backend; the **client** can update its own metadata mid-session via `useUpdatePeerMetadata` / `update_peer_metadata`, emitting `peerMetadataUpdated` to notifier/webhook).
- A `peerToken` issued by the backend that authorizes the client to connect.

A peer publishes zero or more **tracks** and subscribes to other peers' tracks (automatic by default; manual via `subscribePeer` / `subscribeTracks`).

## Track

A single stream of media — one camera, one microphone, one screen share, or one audio output from an agent. A peer typically publishes multiple tracks simultaneously. Track shape:

- `id` — server-assigned per track.
- `type` — `video` | `audio` | `unspecified`.
- `metadata` — opaque string set by the publishing peer; visible to all subscribers.

## Agent

A programmatic peer — a server-side program that joins a room via the server SDK (`FishjamClient.createAgent` / `create_agent`) and sends or receives audio frames directly. Used for AI voice integrations (e.g., Gemini Live) and custom server-side audio processing. Agents do not run a WebRTC negotiation; they use a direct protocol over WebSocket.

## Streamer

Livestream-room-specific peer role: the **single** peer authorized to publish into a `livestream` room. Authorized by a livestream streamer token (`createLivestreamStreamerToken`). One streamer per livestream room.

## Viewer

Livestream-room-specific peer role: a peer that **only consumes** the livestream. Unlimited viewers per room. For private livestreams a viewer needs a `viewerToken` (`createLivestreamViewerToken`); for public livestreams the room ID alone suffices.

## Channel

Server-side **notification subject** that tracks WebRTC data-channel lifecycle inside a room. Server notifications `channelAdded` / `channelRemoved` carry `roomId`, `channelId` (required), and exactly one of `peerId` / `componentId` (the endpoint that owns the channel). Distinct from **Data Channel** below — "Channel" is the notification entity; "Data Channel" is the client-side RTC primitive you use via `useDataChannel`.

## Management Token

Long-lived secret credential issued from the Fishjam Dashboard. Authorizes the **server SDK** to call REST endpoints — create rooms, mint peer tokens, configure webhooks, set immutable peer metadata. **Never** ship this token to a client. If exposed, regenerate from the Dashboard.

## Peer Token

Short-lived (24-hour), peer-scoped credential issued by the backend (`createPeer` / `create_peer` returns it alongside the peer object). Authorizes one client to join one room as one peer. Safe to send to a frontend. Refresh via `refreshPeerToken` / `refresh_peer_token` to extend validity.

Adjacent token types: **livestream streamer token**, **livestream viewer token**, **MoQ publisher/subscriber token** — see `auth-model.md`.

## Fishjam ID

Identifier of your Fishjam instance. Both the backend (`FishjamClient` constructor) and the client (`FishjamProvider` / `useConnection`) need it. Visible in the Fishjam Dashboard. Not a secret on its own — the management token does the gating — but knowing it lets clients connect to your instance, so don't publish it pointlessly.

## Sandbox API

A no-auth HTTP backend hosted by Fishjam that mints peer tokens, livestream tokens, and MoQ tokens for your instance. Intended for prototyping a frontend before you've stood up a real backend. Activated and rotated from the Dashboard's Sandbox tab. **Not for production** — anyone with the URL can create rooms on your behalf. See `sandbox-vs-production.md`.

## Track Forwarding

A track forwarded from a peer in one room to another peer or room (e.g. for relay-style fan-out). Set up via `POST /room/{room_id}/track_forwardings` on the REST API. This is a create-only endpoint — there is no list/delete equivalent in the user API (forwardings auto-expire when the source track is removed).

## Data Channel

WebRTC data channel between peers, used for arbitrary app-defined messaging (text chat, signals, etc.). Distinct from media tracks. See `useDataChannel` on the client SDKs.

## Webhook

HTTP POST callback Fishjam sends to your backend per room-lifecycle event. Set the `webhookUrl` when creating a room. Payload is binary protobuf (`Content-Type: application/x-protobuf`) — decode with the SDK helper. Same event set as the WebSocket notifier (see `notifications-taxonomy.md`). Picking between the two: `notifier-vs-webhook.md`.

## Notifier

Object exposed by the server SDKs — `FishjamWSNotifier` in JS, `FishjamNotifier` in Python — that opens a persistent WebSocket to Fishjam and emits typed events. Drop-in alternative to webhooks for single-process worker deployments.

## MoQ (Media over QUIC)

Newer Fishjam streaming mode for one-to-many delivery using the [MoQ standard](https://datatracker.ietf.org/wg/moq/about/). Distinct from livestream rooms (which use WHEP/WebRTC). Tokens issued via `createMoqToken` / `create_moq_token`; clients connect to `relay.fishjam.io/<fishjam-id>`. See `room-types.md` and the upstream docs (`docs/explanation/moq-streaming`, `docs/tutorials/moq`).
