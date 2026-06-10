# Livestream & MoQ Tokens

Three additional token types live alongside peer tokens. Issued by `FishjamClient`, consumed by client SDKs:

| Token | Method | Used by | Required for |
|---|---|---|---|
| Streamer token | `createLivestreamStreamerToken(roomId)` | `useLivestreamStreamer` (web/RN) | `livestream` rooms |
| Viewer token | `createLivestreamViewerToken(roomId)` | `useLivestreamViewer` | **Private** livestreams only |
| MoQ token | `createMoqToken({ publishPath?, subscribePath? })` | MoQ relay clients | MoQ streaming |

## Livestream rooms — the flow

```ts
// 1. Create a livestream room
const room = await fishjamClient.createRoom({ roomType: 'livestream' });
//   → for public viewing without tokens:
//     await fishjamClient.createRoom({ roomType: 'livestream', public: true });

// 2. Mint a streamer token for the one peer that broadcasts
const { token: streamerToken } = await fishjamClient.createLivestreamStreamerToken(room.id);
//   → ship to the publishing client; client uses `useLivestreamStreamer({ token: streamerToken })`

// 3. Mint viewer tokens (private livestream)
const { token: viewerToken } = await fishjamClient.createLivestreamViewerToken(room.id);
//   → ship to each viewer; client uses `useLivestreamViewer().connect({ token: viewerToken })`

// 3a. For public livestreams, viewers connect with just the room ID:
//     useLivestreamViewer().connect({ streamId: room.id })
```

Streamer + viewer tokens have their own lifetimes; check the OpenAPI response shape if you need the exact `expiresAt`.

Constraints:

- **One streamer per room.** Issuing a second streamer token doesn't add a streamer — the room only ever has one active broadcasting peer.
- **One video + one audio track on the streamer.** Additional tracks are silently dropped by Fishjam.

## Public vs private livestreams

| Mode | `public` | Viewer token? | Use when |
|---|---|---|---|
| Private (default) | `false` (or omitted) | required per viewer | Paywalled, authenticated, internal broadcasts |
| Public | `true` | none — viewers connect with room ID | Open broadcasts; remember you still pay per viewer-minute |

`public: true` is irreversible for that room — to switch, create a new room.

## Revocation

To kick a viewer or streamer, hit the REST API directly (no high-level method right now):

```http
DELETE /livestream/{stream_id}/viewer/{viewer_id}
DELETE /livestream/{stream_id}/streamer/{streamer_id}
```

The `stream_id` here is the livestream ID surfaced when the streamer/viewer joins. Use the room ID lookup or events to obtain it.

## MoQ tokens

Media-over-QUIC streaming uses a separate auth model — short-lived JWTs scoped to a path. See `../platform/room-types.md` for the MoQ conceptual overview and the upstream `docs/explanation/moq-streaming.mdx`.

```ts
import type { MoqTokenConfig } from '@fishjam-cloud/js-server-sdk';

// Publisher token — broad path: client can name themselves under it
const { token: pubToken } = await fishjamClient.createMoqToken({
  publishPath: 'my-room',
});

// Publisher token — specific path: identity is server-pinned
const { token: pubTokenAlice } = await fishjamClient.createMoqToken({
  publishPath: 'my-room/alice',
});

// Subscriber token — sees all publishers in the namespace
const { token: subToken } = await fishjamClient.createMoqToken({
  subscribePath: 'my-room',
});
```

Path scoping is **prefix-based**:

- `publishPath: 'stream-name'` permits broadcasts at `stream-name/*`.
- `publishPath: 'stream-name/alice'` permits only the exact path.
- Same for `subscribePath`.

Your Fishjam ID is automatically the root namespace — never include it in the paths.

The client attaches the token as `?jwt=<token>` when connecting to `https://relay.fishjam.io/<fishjam-id>?jwt=<token>` (HTTPS bootstrap — MoQ runs over WebTransport/QUIC, not WebSocket). The Fishjam ID is the only URL path segment; the publish/subscribe path is encoded in the JWT.

## Client-side counterpart

- Web: `../react-client/livestream.md` (for livestream WHEP-style flow).
- React Native: `../react-native-client/SKILL.md` (re-exports `useLivestreamStreamer` / `useLivestreamViewer` from the web client).
- MoQ: no first-party React hook — connect with a MoQ-compatible JS library (e.g. `@moq/publish` / `@moq/watch`); see the MoQ examples in the `js-server-sdk` repo and on the Fishjam docs site.

## Sources

- `js-server-sdk` repo: `packages/js-server-sdk/src/client.ts` (`createLivestream*Token`, `createMoqToken`)
- <https://fishjam.swmansion.com/docs/tutorials/livestreaming>
- <https://fishjam.swmansion.com/docs/tutorials/moq>
- <https://fishjam.swmansion.com/docs/explanation/moq-streaming>
