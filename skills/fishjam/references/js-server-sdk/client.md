# FishjamClient

REST client for the Fishjam server. Source: `packages/js-server-sdk/src/client.ts` in the `js-server-sdk` repo (<https://github.com/fishjam-cloud/js-server-sdk>).

## Constructor

```ts
import { FishjamClient } from '@fishjam-cloud/js-server-sdk';

const fishjamClient = await FishjamClient.create({
  fishjamId: process.env.FISHJAM_ID!,
  managementToken: process.env.FISHJAM_MANAGEMENT_TOKEN!,
});
```

A thin axios wrapper — the management token goes out as an `Authorization: Bearer` header on every request. The async factory `await FishjamClient.create(config)` **does** validate credentials. The synchronous constructor flow does **not** validate credentials on its own unless you explicitly call `client.checkCredentials()`, which throws `InvalidFishjamCredentialsException` when the backend rejects them. If the backend reports the SDK version as deprecated or unsupported, the client logs a one-time console warning/error.

## Rooms

```ts
import type { RoomConfig, Room, RoomId } from '@fishjam-cloud/js-server-sdk';

await fishjamClient.createRoom();                                              // defaults
await fishjamClient.createRoom({ roomType: 'audio_only' });
await fishjamClient.createRoom({ roomType: 'livestream', public: true });
await fishjamClient.createRoom({ webhookUrl: 'https://api/webhook/<secret>' });
await fishjamClient.createRoom({ batchWebhookNotifications: true });
await fishjamClient.createRoom({ videoCodec: 'vp8', maxPeers: 20 });

const room: Room        = await fishjamClient.getRoom(roomId);
const rooms: Room[]     = await fishjamClient.getAllRooms();
await fishjamClient.deleteRoom(roomId);
```

`RoomConfig` (from `@fishjam-cloud/fishjam-openapi`) includes:

- `roomType?: 'conference' | 'audio_only' | 'livestream' | 'audio_only_livestream' | 'full_feature' | 'broadcaster'` — prefer the modern values; `'full_feature'` and `'broadcaster'` are legacy aliases
- `public?: boolean` (livestream-only — viewers can connect without a viewer token)
- `videoCodec?: 'h264' | 'vp8'`
- `webhookUrl?: string`
- `batchWebhookNotifications?: boolean` (coalesce per-room webhook events into a single `NotificationBatch` POST; VAD notifications are unaffected)
- `maxPeers?: number` (cap the number of peers allowed in the room)

Note: `subscribeMode` is **NOT** a `RoomConfig` field — it lives on `PeerOptionsWebRTC` / `PeerOptionsAgent` / `PeerOptionsVapi`. See "Selective subscriptions" below.

Returns the `Room` shape — `{ id: RoomId, peers: Peer[], config: RoomConfig }`. `RoomId` and `PeerId` are *branded* strings — pass them around as-is, don't `String(...)` cast.

## Peers (WebRTC clients)

```ts
import type { PeerOptionsWebRTC, Peer, PeerId } from '@fishjam-cloud/js-server-sdk';

const { peer, peerToken } = await fishjamClient.createPeer(roomId, {
  metadata: { userId: 'user-123', role: 'host' },
});

await fishjamClient.deletePeer(roomId, peer.id);

const newToken: string = await fishjamClient.refreshPeerToken(roomId, peer.id);
```

`PeerOptionsWebRTC.metadata` is the peer's **server metadata** — set by your backend at creation, fixed for the peer's lifetime, and read by clients as `peer.metadata.server`. Nobody can change it afterwards: clients have no API for it, and there is no REST endpoint to update a peer.

Don't confuse it with **peer metadata** (`peer.metadata.peer` on clients) — a separate, client-owned field set at `joinRoom({ peerMetadata })` and updated mid-session via `useUpdatePeerMetadata` / `updatePeerMetadata`; those updates surface as `peerMetadataUpdated` notifications. Rule of thumb: server metadata for trusted facts your backend vouches for (user id, role), peer metadata for client-mutable state (display name).

`peerToken` is the credential you hand to the client. It's valid for 24 hours from creation. The initial WS handshake consumes it; once the peer has joined, the established session stays alive regardless of token wall-clock expiry. For peers that haven't yet connected, or need to *reconnect* after the 24h window, mint a fresh token with `refreshPeerToken`.

`PeerOptionsWebRTC` also accepts `subscribeMode?: 'auto' | 'manual'` — see "Selective subscriptions" below.

## Agents (server-side programmatic peers)

```ts
import type { PeerOptionsAgent, AgentCallbacks } from '@fishjam-cloud/js-server-sdk';

const { agent, peer } = await fishjamClient.createAgent(roomId, {
  output: { audioFormat: 'pcm16', audioSampleRate: 16000 }, // audioFormat only 'pcm16'; audioSampleRate only 16000 | 24000 (pick 24000 for Gemini Live)
}, {
  onError: (event) => console.error('agent ws error', event), // event: DOM Event, not Error
  onClose: (code, reason) => console.log('agent ws closed', code, reason),
});

// agent is a connected FishjamAgent — see agent.md
```

`PeerOptionsAgent.output` is a single `AgentOutput` (NOT an array) that controls the format of incoming audio — what other peers' tracks are downmixed into for the agent to consume. `PeerOptionsAgent` also accepts `subscribeMode?: 'auto' | 'manual'`.

## Vapi agents (managed voice agents)

```ts
import type { PeerOptionsVapi } from '@fishjam-cloud/js-server-sdk';

const { peer } = await fishjamClient.createVapiAgent(roomId, {
  apiKey: process.env.VAPI_API_KEY!,
  callId: vapiCallId,
});
```

No client-side agent object — Vapi runs the agent for you. `PeerOptionsVapi` requires both `apiKey` and `callId` (strings); also accepts optional `subscribeMode`.

## Selective subscriptions

Only meaningful for peers created with `subscribeMode: 'manual'` (set per peer in `PeerOptionsWebRTC` / `PeerOptionsAgent` / `PeerOptionsVapi` — **not** on the room). See `selective-subscriptions.md`.

```ts
import type { TrackId } from '@fishjam-cloud/js-server-sdk';

// At peer creation:
const { peer, peerToken } = await fishjamClient.createPeer(roomId, {
  metadata: { userId: 'user-123' },
  subscribeMode: 'manual',
});

await fishjamClient.subscribePeer(roomId, subscriberPeerId, publisherPeerId);
await fishjamClient.subscribeTracks(roomId, subscriberPeerId, [trackId1, trackId2] as TrackId[]);
```

## Livestream tokens

```ts
const { token: streamerToken } = await fishjamClient.createLivestreamStreamerToken(roomId);
const { token: viewerToken }   = await fishjamClient.createLivestreamViewerToken(roomId);
```

Issue these only for `livestream`-type rooms. Public livestreams don't need viewer tokens. See `livestream-and-moq.md`.

## MoQ tokens

```ts
import type { MoqTokenConfig } from '@fishjam-cloud/js-server-sdk';

const { token: pubToken } = await fishjamClient.createMoqToken({ publishPath: 'my-room/alice' });
const { token: subToken } = await fishjamClient.createMoqToken({ subscribePath: 'my-room' });
```

## Exception types

Every method wraps axios errors in typed exceptions via `mapException`. Import from `@fishjam-cloud/js-server-sdk`:

| Class | Triggers |
|---|---|
| `MissingFishjamIdException` | `fishjamId` blank in constructor (sync throw) |
| `InvalidFishjamCredentialsException` | credentials rejected (401/404) during synchronous constructor credential check / explicit `checkCredentials()` |
| `BadRequestException` | 400 — invalid request body or path params |
| `UnauthorizedException` | 401 — bad management token |
| `FishjamNotFoundException` | 404 — generic not found (instance / endpoint) |
| `RoomNotFoundException` | 404 on a room-scoped endpoint |
| `PeerNotFoundException` | 404 on a peer-scoped endpoint |
| `ServiceUnavailableException` | 503 — Fishjam temporarily unavailable |
| `UnknownException` | anything else — including 403 (the current mapper has no `case 403` so it falls through here) |

`ForbiddenException` is also exported, but the current `mapException` never produces it. Catch `UnknownException` if you need to handle 403 responses.

All extend `FishjamBaseException` with `.statusCode`, `.axiosCode`, `.details` fields.

Pattern for mapping them to HTTP responses — call explicitly from each route's `catch` (see `express-fastify.md`):

```ts
import {
  RoomNotFoundException,
  PeerNotFoundException,
  UnauthorizedException,
} from '@fishjam-cloud/js-server-sdk';

function respondFishjamError(err: unknown, res: Response) {
  if (err instanceof RoomNotFoundException || err instanceof PeerNotFoundException) {
    return res.status(404).json({ error: err.message });
  }
  if (err instanceof UnauthorizedException) {
    return res.status(502).json({ error: 'Fishjam credentials misconfigured' }); // not the user's auth!
  }
  return res.status(500).json({ error: 'Internal' });
}
```

## Sources

- `js-server-sdk` repo (<https://github.com/fishjam-cloud/js-server-sdk>): `packages/js-server-sdk/src/client.ts`, `packages/js-server-sdk/src/index.ts`, `packages/js-server-sdk/src/exceptions/index.ts`, `packages/js-server-sdk/src/types.ts`
- `packages/fishjam-openapi/src/generated/api.ts` (`RoomConfig`, `PeerOptionsWebRTC` / `PeerOptionsAgent` / `PeerOptionsVapi`, `RoomType`)
- <https://fishjam.swmansion.com/docs/how-to/backend/server-setup>
- <https://fishjam.swmansion.com/docs/how-to/client/metadata> (server metadata vs peer metadata)
- <https://fishjam.swmansion.com/docs/api/server> (generated TypeDoc)
