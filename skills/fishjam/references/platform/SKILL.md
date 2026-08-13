---
name: fishjam-platform
description: "Fishjam platform fundamentals — domain model and auth shared by all SDKs. Covers glossary (room, peer, track, agent, streamer, viewer), the four room types (conference / audio_only / livestream / audio_only_livestream), two-tier auth (management vs peer tokens), Sandbox vs production, WebSocket-vs-webhook notifications and the 18-event taxonomy, REST endpoint list, and the end-to-end auth/join flow. Read this first for any Fishjam task. Trigger on: 'Fishjam', 'fishjam.io', 'fishjam.swmansion.com', 'fishjam dashboard', 'management token', 'peer token', 'sandbox API', 'useSandbox', 'room manager', 'livestream room', 'audio-only room', 'conference room', 'Fishjam ID', 'room types', 'auth flow', 'webhook', 'webhookUrl', 'FishjamWSNotifier', 'FishjamNotifier', 'notifier vs webhook', 'ServerMessage', 'peerAdded', 'peerConnected', 'roomCreated', 'refreshPeerToken', 'rest endpoint', 'OpenAPI', 'llms.txt'."
license: MIT
---

# Fishjam Platform

Software Mansion's hosted WebRTC streaming toolkit. Fishjam takes care of media servers, peer networking, and codec negotiation so applications can build video conferencing, audio chat, and one-to-many livestreaming with sub-second latency.

This sub-skill defines the **domain model, auth, room types, and notifications** that every SDK builds on. Read it first; then drop into the SDK sub-skill you need.

A production Fishjam app always has **a backend (server SDK) + a client (client SDK)**. The backend mints peer tokens; the client uses them to connect. For prototyping you can skip the backend with the Sandbox API — see `sandbox-vs-production.md`.

## Key concepts at a glance

- **Room** — a collection of peers exchanging tracks. Three primary room types: `conference` (default, bidirectional), `audio_only` (75% cheaper, video silently dropped), `livestream` (one streamer, many viewers, <300 ms). `audio_only_livestream` combines the audio-only and livestream discounts.
- **Peer** — a participant in a room. Can be a `webrtc` client, an `agent` (server-side program), or a `vapi` voice agent.
- **Track** — one audio or video stream from a peer (camera, mic, screen share, etc.). One peer can publish multiple tracks.
- **Streamer / Viewer** — livestream-room-specific peer roles (one streamer, many viewers).
- **Fishjam ID** — your instance identifier; goes in both backend and client config.
- **Management Token** — long-lived admin secret; backend only.
- **Peer Token** — 24h, scoped to one peer in one room; safe for clients.

Full definitions: `glossary.md`.

## Room types

| Type         | Direction                | Cost        | Use case                                   |
| ------------ | ------------------------ | ----------- | ------------------------------------------ |
| `conference` | many ↔ many              | base        | Meetings, classrooms, interactive webinars |
| `audio_only` | many ↔ many (audio only) | **75% off** | Voice rooms, podcasts, town halls          |
| `livestream` | one → many               | **20% off** | Live events, broadcasts, sports            |

`audio_only_livestream` is 75% off the livestream rate (equivalently, 80% off conference) — the cheapest mode for one-to-many audio.

Detail + decision matrix: `room-types.md`.

## Auth flow (5-step sketch)

1. **Backend → Fishjam:** `createPeer(roomId)` authenticated with the **management token**.
2. **Fishjam → Backend:** returns `{ peer, peerToken }`.
3. **Backend → Client:** ships `peerToken` to the frontend (never the management token).
4. **Client → Fishjam:** `joinRoom({ peerToken })`.
5. **Fishjam:** completes the WebRTC handshake; the peer is now in the room.

Full diagram + token-refresh + livestream + MoQ tokens: `auth-model.md`.

## Sandbox vs production

- **Sandbox API** — no-auth HTTP backend hosted by Fishjam. For prototyping only. Generated URL must never be shipped in a production build (anyone with it can create rooms on your behalf).
- **Production** — your own backend running a server SDK. Required for any real-world deployment.

Migration recipe: `sandbox-vs-production.md`.

## Notifications: WebSocket or webhook?

Both deliver the **same 18-event set surfaced by the SDKs** (room/peer/track/streamer/viewer/channel lifecycle — the underlying protobuf `ServerMessage` carries additional oneof fields the SDKs intentionally ignore). Pick by deployment shape:

- **Long-lived worker, single process** → `FishjamWSNotifier` / `FishjamNotifier` (persistent socket).
- **Serverless, scaled out, behind a load balancer** → webhooks (set `webhookUrl` per room, decode binary protobuf).

Detail: `notifier-vs-webhook.md` and `notifications-taxonomy.md`.

## Key rules

- Management token **never** leaves the backend. If it leaks, regenerate it.
- Peer tokens expire **24 hours** after creation. Refresh via `refreshPeerToken` / `refresh_peer_token` if you need longer.
- Sandbox API URL is **not** a secret in the cryptographic sense, but anyone holding it can create rooms — treat it like a dev credential and rotate from the Dashboard if exposed.
- Conference is the default room type. Use `audio_only` for voice apps (much cheaper). Use `livestream` only when one peer broadcasts and others only watch.
- Livestream rooms accept exactly **one video track + one audio track**. Additional tracks are silently dropped.

## References

| File                        | When to read                                                                    |
| --------------------------- | ------------------------------------------------------------------------------- |
| `glossary.md`               | Define a term (room, peer, track, agent, etc.) precisely.                       |
| `room-types.md`             | Pick `conference` vs `audio_only` vs `livestream`. Cost matrix. Codec defaults. |
| `auth-model.md`             | Token types, lifetimes, refresh, livestream / MoQ tokens, full auth flow.       |
| `sandbox-vs-production.md`  | Use the Sandbox API for dev; migrate to a real backend.                         |
| `notifier-vs-webhook.md`    | Pick WebSocket notifier or webhook receiver.                                    |
| `notifications-taxonomy.md` | All 18 SDK-surfaced server notification events with payload shapes.             |
| `rest-endpoints.md`         | Raw HTTP surface — flat endpoint list, auth header, deprecation header.         |
| `lifecycle-flow.md`         | End-to-end client↔backend↔Fishjam sequence.                                     |
| `llms-and-docs.md`          | Pointers to upstream docs, llms.txt, OpenAPI spec, protobuf, dashboard.         |
