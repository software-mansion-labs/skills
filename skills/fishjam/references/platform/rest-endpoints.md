# Fishjam REST API Endpoints

The full HTTP surface exposed by Fishjam's server. The server SDKs wrap these — prefer the SDK whenever possible. Use this reference when calling Fishjam from a language without an SDK, debugging a wire problem, or inspecting OpenAPI / generated clients directly.

Authoritative spec: `fishjam-server-openapi.yaml` in the `documentation` repo (upstream: <https://github.com/fishjam-cloud/documentation/blob/main/static/api/fishjam-server-openapi.yaml>).

## Auth

Every request requires:

```http
Authorization: Bearer <management-token>
```

(Token issued from the Fishjam Dashboard. See `auth-model.md`.)

## Base URL

For Fishjam Cloud:

```
https://fishjam.io/api/v1/connect/<your-fishjam-id>
```

The Fishjam ID is the hex string shown in the [Dashboard](https://fishjam.io/app). All paths below are appended to this base. For self-hosted Fishjam, pass the full instance URL instead of a bare ID. The SDKs handle this for you via the `fishjamId` config.

## Rooms

| Method | Path | Purpose |
|---|---|---|
| `POST` | `/room` | Create a room. Body (`RoomConfig`): `{ roomType?, public?, webhookUrl?, videoCodec?, maxPeers?, batchWebhookNotifications? }`. Returns 201 with `{ data: { room, fishjam_url } }` (every success response in this API is wrapped in a top-level `data` envelope). |
| `GET` | `/room` | List all rooms in the instance. |
| `GET` | `/room/{room_id}` | Get one room with its peers and tracks. |
| `DELETE` | `/room/{room_id}` | Delete a room. All peers are disconnected. |

## Peers

| Method | Path | Purpose |
|---|---|---|
| `POST` | `/room/{room_id}/peer` | Create a peer in the room. Body: `{ type: 'webrtc' \| 'agent' \| 'vapi', options: PeerOptionsWebRTC \| PeerOptionsAgent \| PeerOptionsVapi }` (the `type` discriminator selects which options shape is valid). Returns `{ peer, token, websocket_url }`. |
| `DELETE` | `/room/{room_id}/peer/{id}` | Remove a peer. |
| `POST` | `/room/{room_id}/peer/{id}/refresh_token` | Issue a new peer token (same 24h validity — see `auth-model.md`). Returns `{ token }`. |
| `POST` | `/room/{room_id}/peer/{id}/subscribe_peer?peer_id={publisher_peer_id}` | In manual subscribe mode: subscribe `id` to all tracks published by the peer named in the `peer_id` query string. Empty body. |
| `POST` | `/room/{room_id}/peer/{id}/subscribe_tracks` | In manual subscribe mode: subscribe `id` to specific tracks. Body: `{ "track_ids": ["...", "..."] }` (required). |

The `subscribePeer` and `subscribeTracks` endpoints only matter for peers created with `subscribeMode: 'manual'` (set on the peer's `PeerOptionsWebRTC` / `PeerOptionsAgent` / `PeerOptionsVapi`, **not** on the room). In the default `auto` mode each peer is subscribed to every other peer's tracks automatically.

## Livestream — direct management

For lower-level livestream management (independent of the room):

| Method | Path | Purpose |
|---|---|---|
| `GET` | `/livestream` | List livestreams. |
| `POST` | `/livestream` | Create a livestream directly. |
| `GET` | `/livestream/{stream_id}` | Get one livestream. |
| `DELETE` | `/livestream/{stream_id}` | Delete a livestream. |
| `POST` | `/livestream/{stream_id}/streamer` | Add a streamer to an existing livestream. |
| `DELETE` | `/livestream/{stream_id}/streamer/{streamer_id}` | Remove a streamer. |
| `POST` | `/livestream/{stream_id}/viewer` | Create a viewer for the stream and return its credentials. |
| `DELETE` | `/livestream/{stream_id}/viewer/{viewer_id}` | Revoke a viewer (kicks the WHEP session). |

Most apps interact with livestreams via the room endpoints; the `/livestream` family is for advanced flows.

## Track Forwarding

| Method | Path | Purpose |
|---|---|---|
| `POST` | `/room/{room_id}/track_forwardings` | Forward a track from this room to another peer or room. Used for relay / fan-out patterns. |

## MoQ

| Method | Path | Purpose |
|---|---|---|
| `POST` | `/moq/token` | Issue a MoQ JWT. Body (`MoqTokenConfig`): `{ publishPath?, subscribePath? }`. Returns `{ token }`. See `room-types.md` for MoQ context. |

## Response shape

JSON for all responses except webhook deliveries. Errors come back as HTTP 4xx/5xx with a JSON body — the SDKs map these to typed exceptions (`exceptions/` directory in the JS SDK, `errors.py` in the Python SDK).

## Deprecation header

Responses may carry `x-fishjam-api-deprecated: <message>` when you call a deprecated endpoint or use a deprecated request shape. The JS SDK logs this once per process and continues; treat it as a TODO to read the changelog and migrate.

## OpenAPI spec

The full OpenAPI YAML lives at:

- Upstream: <https://github.com/fishjam-cloud/documentation/blob/main/static/api/fishjam-server-openapi.yaml>
- Served by the docs site: <https://fishjam.swmansion.com/docs/api/fishjam-server-openapi.yaml>

Use it to:

- Generate clients for languages without an official SDK.
- Inspect exact request/response schemas including optional fields.
- Verify deprecation status before you commit to a flow.
