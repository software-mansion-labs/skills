# Connection & Sandbox

## useConnection

The hook for joining/leaving rooms and reading status. Source: `@fishjam-cloud/react-client` → `packages/react-client/src/hooks/useConnection.ts`.

```tsx
import { useConnection } from '@fishjam-cloud/react-client';

function Room() {
  const { joinRoom, leaveRoom, peerStatus, reconnectionStatus } = useConnection();
  // ...
}
```

### `joinRoom({ peerToken, peerMetadata? })`

```tsx
await joinRoom({ peerToken: 'JWT_FROM_YOUR_BACKEND' });

// Typed metadata flowing to other peers:
await joinRoom<MyPeerMetadata>({
  peerToken,
  peerMetadata: { displayName: 'Alice', role: 'host' },
});
```

`JoinRoomConfig`:

- `peerToken: string` — required. From your backend (production) or `useSandbox().getSandboxPeerToken(...)` (dev).
- `peerMetadata?: PeerMetadata` — optional. Visible to other peers via `usePeers()`. Different from the server-side metadata set in `createPeer` (which is immutable; this client-side one can be updated later via `useUpdatePeerMetadata`).

Throws if `fishjamId` is missing from the provider.

### `leaveRoom()`

```tsx
leaveRoom();
```

Disconnects. Devices stay initialized — you can `joinRoom` again without re-initializing.

### `peerStatus`

Current peer connection state. `PeerStatus` is one of (from `@fishjam-cloud/react-client`):

- `'idle'` — not connected
- `'connecting'`
- `'connected'`
- `'error'` — connection failed (auth failure, network, etc.); inspect underlying client events for the reason

```tsx
if (peerStatus === 'connected') ...
if (peerStatus === 'error') {
  // connection failed — likely a bad/expired token or network issue.
  // Listen on the underlying client (`fishjamClient` prop) for granular reasons.
}
```

### `reconnectionStatus`

Status of the SDK's auto-reconnect attempts:

- `'idle'`
- `'reconnecting'`
- `'idle'` again after success
- `'error'` after `maxAttempts` exhausted

Drive UI off this — e.g. show a "Reconnecting…" banner while `reconnectionStatus === 'reconnecting'`.

### Error reasons

Connection failures surface as `peerStatus === 'error'`. The reason is delivered via internal events; import the union types:

```tsx
import type { AuthErrorReason, JoinErrorReason } from '@fishjam-cloud/react-client';

// AuthErrorReason: auth-time failure reasons (token expired/invalid, room not found, ...)
// JoinErrorReason: connection-time errors
```

(Listen via the underlying `FishjamClient` events — `authError`, `joinError`, etc. — if you need granular handling. See `ts-client-escape.md`.)

## useSandbox (dev only)

For prototyping a frontend before you've stood up a backend. Calls Fishjam's no-auth Sandbox API. **Never ship this in production.**

Source: `@fishjam-cloud/react-client` → `packages/react-client/src/hooks/useSandbox.ts`.

```tsx
import { useSandbox } from '@fishjam-cloud/react-client';

const {
  getSandboxPeerToken,
  getSandboxViewerToken,
  getSandboxLivestream,
} = useSandbox({ sandboxApiUrl: import.meta.env.VITE_SANDBOX_API_URL });
```

### `getSandboxPeerToken(roomName, peerName, roomType?)`

```tsx
const peerToken = await getSandboxPeerToken('test-room', 'alice');
const peerToken = await getSandboxPeerToken('voice-room', 'alice', 'audio_only');
const peerToken = await getSandboxPeerToken('stream-room', 'alice', 'livestream');

await joinRoom({ peerToken });
```

`roomType: 'conference' | 'livestream' | 'audio_only'` — defaults to `'conference'`.

### `getSandboxViewerToken(roomName)`

For private livestream rooms — issue a viewer token:

```tsx
const viewerToken = await getSandboxViewerToken('my-livestream');
const { connect } = useLivestreamViewer();
await connect({ token: viewerToken });
```

Throws if the livestream room doesn't exist.

### `getSandboxLivestream(roomName, isPublic?)`

Creates a livestream room and returns the **streamer token**:

```tsx
const { streamerToken, room } = await getSandboxLivestream('my-livestream');
// or for a public livestream:
const { streamerToken, room } = await getSandboxLivestream('my-livestream', true);

const { connect } = useLivestreamStreamer();
await connect({ token: streamerToken });
```

`room: { id, name }` lets you share the room ID with viewers (public livestreams need only the ID).

### Migration off the sandbox

Replace each `useSandbox` call with a `fetch` to your own backend's `/api/join-room` endpoint. Detail: `../platform/sandbox-vs-production.md`.

## Sources

- `@fishjam-cloud/react-client` → `packages/react-client/src/hooks/useConnection.ts`
- `packages/react-client/src/hooks/useSandbox.ts`
- `packages/react-client/src/types/public.ts` (`PeerStatus`)
- <https://fishjam.swmansion.com/docs/how-to/client/connecting>
- <https://fishjam.swmansion.com/docs/how-to/client/reconnection-handling>
- Cross-reference: `../platform/auth-model.md`, `../platform/sandbox-vs-production.md`
