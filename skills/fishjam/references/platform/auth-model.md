# Auth Model

Fishjam uses **two distinct token types** with very different blast radii. Mixing them up is the #1 source of bugs — the management token must never leave the backend.

## Tier 1 — Management Token

- **Holder:** your backend (server SDK).
- **Issued by:** Fishjam Dashboard (<https://fishjam.io/app>). Manually rotated.
- **Lifetime:** long-lived (until rotated).
- **Scope:** full admin — create / delete rooms, create / delete peers, mint peer tokens, set per-room webhook URLs (via `RoomConfig.webhookUrl` / `RoomOptions.webhook_url` on room creation), set immutable peer metadata.
- **Auth header:** `Authorization: Bearer <management-token>` on every REST call.
- **Risk:** anyone with this token can spin up rooms and peers on your account. If exposed, regenerate immediately from the Dashboard.

The server SDKs accept it once in their constructor:

```ts
const fishjamClient = new FishjamClient({
  fishjamId: process.env.FISHJAM_ID!,
  managementToken: process.env.FISHJAM_MANAGEMENT_TOKEN!,
});
```

```python
from fishjam import FishjamClient
fishjam_client = FishjamClient(
  fishjam_id=os.environ['FISHJAM_ID'],
  management_token=os.environ['FISHJAM_MANAGEMENT_TOKEN'],
)
```

## Tier 2 — Participant (Peer, Viewer, Streamer) Tokens

- **Holder:** the client (browser, mobile app).
- **Issued by:** your backend, by calling `createPeer` / `create_peer` (which internally calls `POST /room/{room_id}/peer`).
- **Lifetime:** 24 hours from creation. After that, calls to Fishjam with the token fail with an auth error.
- **Scope:** one specific peer in one specific room. Cannot create rooms, list rooms, or impersonate other peers.
- **Safe for client transport:** yes — send it back from your `/join-room` endpoint to the frontend.

### Peer token

Generation:

```ts
const { peer, peerToken } = await fishjamClient.createPeer(roomId, {
  metadata: { userId: 'user-123', name: 'Jane Doe' },
});
res.json({ peerToken, roomId });
```

```python
peer, token = fishjam_client.create_peer(
  room_id,
  PeerOptions(metadata={'userId': 'user-123', 'name': 'Jane Doe'}),
)
return {'peerToken': token, 'roomId': room_id}
```

Client usage:

```ts
const { joinRoom } = useConnection();
await joinRoom({ peerToken });
```

#### Token refresh

Peer tokens expire after 24 hours. Two strategies:

1. **Short-lived rooms (recommended).** Treat a peer token as a one-call credential. Tear down the room when the call ends, and create a fresh peer on the next join. This matches how almost all real-world meetings work.
2. **Long-lived peers.** Call `refreshPeerToken(roomId, peerId)` (JS) or `refresh_peer_token(room_id, peer_id)` (Python) from your backend on a schedule, return the new token to the client, and have the client update its connection. Use this only when you genuinely need a peer that outlives 24h.

The HTTP endpoint is `POST /room/{room_id}/peer/{id}/refresh_token`.

### Livestream tokens

Livestream rooms (`roomType: 'livestream'`) use additional token types issued by the backend:

| Token | Method (JS) | Method (Python) | Recipient |
|---|---|---|---|
| Streamer token | `createLivestreamStreamerToken(roomId)` | `create_livestream_streamer_token(room_id)` | The single broadcasting peer |
| Viewer token (private livestream only) | `createLivestreamViewerToken(roomId)` | `create_livestream_viewer_token(room_id)` | Each viewer; public livestreams skip this |

Both endpoints under the hood:

- `POST /room/{room_id}/streamer`
- `POST /room/{room_id}/viewer`

Detail: `room-types.md`.

### MoQ tokens

For [Media-over-QUIC](https://fishjam.swmansion.com/docs/explanation/moq-streaming) streaming Fishjam issues a separate, path-scoped JWT.

- JS: `fishjamClient.createMoqToken({ publishPath })` → publisher token (write to one path / prefix); `createMoqToken({ subscribePath })` → subscriber token (read from one path / prefix).
- Python: `fishjam_client.create_moq_token(publish_path=...)` or `create_moq_token(subscribe_path=...)` — kwargs, not an object.

Client passes it as a query param when connecting to the relay: `https://relay.fishjam.io/<fishjam-id>?jwt=<token>`. The Fishjam ID is the only URL path segment — the publish/subscribe path is encoded in the JWT (and used by the MoQ client to name each broadcast). Path scoping inside the token is prefix-based — `subscribePath: 'my-room'` lets the client consume any broadcast under `my-room/*`.

Endpoint: `POST /moq/token`.

## Full auth flow

1. User calls the backend controlled by you with your app-level auth.
2. Your backend validates the user against your DB/auth system.
3. Your backend calls Fishjam `createPeer` using the **management token**.
4. Fishjam returns `{ peer, peerToken }` to your backend.
5. Your backend returns `{ peerToken, roomId }` to the client.
6. Client joins Fishjam with `peerToken`.

Backend authenticates the user (your auth scheme, not Fishjam's), creates the peer using the management token, returns the peer token. The client never sees the management token; Fishjam never sees your user IDs (unless you put them in `metadata`).

## Rules

- **Never** put the management token in client code, in HTML, in environment variables exposed to the frontend bundler (e.g. `NEXT_PUBLIC_*`, `VITE_*`), or in API responses to clients.
- **Always** authenticate your own users *before* you call `createPeer` on their behalf. Fishjam authorizes the bearer of the peer token, not the user behind it. Your backend is the policy enforcer.
- **Always** rotate the management token if it leaks. The Dashboard makes this one click.
- Peer tokens expire silently — clients will see a connect / reconnect failure. Detect at the client and re-fetch from your backend.

## Source

- <https://fishjam.swmansion.com/docs/explanation/security-tokens>
- <https://fishjam.swmansion.com/docs/how-to/backend/production-deployment> (token refresh section)
- OpenAPI: `fishjam-server-openapi.yaml` in the `documentation` repo
