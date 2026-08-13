# Sandbox vs Production

The Sandbox API exists so a developer can prototype a Fishjam frontend without standing up a backend. It is **not** for production — anyone with the Sandbox URL can mint peer tokens against your Fishjam instance.

## What the Sandbox API does

A no-auth HTTP backend hosted by Fishjam. Three flavors of token issuance:

```http
GET <sandbox-url>?roomName=...&peerName=...&roomType=...    → { peerToken, room, peer }
GET <sandbox-url>/livestream?roomName=...&public=<bool>     → { streamerToken, room }
GET <sandbox-url>/<roomName>/livestream-viewer-token        → { token }
GET <sandbox-url>/moq/<streamName>/publisher                → { token }   # MoQ publisher
GET <sandbox-url>/moq/<streamName>/subscriber               → { token }   # MoQ subscriber
```

`<sandbox-url>` is the URL shown on the Sandbox tab — shape:
`https://fishjam.io/api/v1/connect/<room-manager-uuid>/room-manager`. The UUID
in the path is the room-manager's own rotatable ID (not the Fishjam ID),
which is why "regenerate" on the Sandbox tab invalidates the old URL.

`useSandbox` exposes the first three as `getSandboxPeerToken` / `getSandboxLivestream` / `getSandboxViewerToken`. Its `RoomType` is `'conference' | 'livestream' | 'audio_only'` on the web client; **on React Native the SDK uses `'audio-only'` (hyphen) — do not pass `'audio_only'` (underscore) on RN or the Sandbox API will reject the request.** `audio_only_livestream` is reachable via the server SDK (`createRoom({ roomType: 'audio_only_livestream' })`) rather than the Sandbox HTTP surface.

From the client, the `useSandbox` hook (web and React Native) wraps these calls:

```ts
import { useSandbox } from '@fishjam-cloud/react-client';

const { getSandboxPeerToken, getSandboxLivestream, getSandboxViewerToken } =
  useSandbox({ sandboxApiUrl: SANDBOX_API_URL });

const peerToken = await getSandboxPeerToken('test-room', 'alice');
await joinRoom({ peerToken });
```

## Why you cannot use it in production

- **No auth.** No API key, no per-user gating, no rate limiting.
- **Public URL = full access.** If you ship the URL in a built frontend, JS bundle, or env var exposed to the client (`VITE_SANDBOX_API_URL`, `NEXT_PUBLIC_SANDBOX_URL`, etc.), every visitor can create rooms on your account.
- **Identical names = identical tokens.** Two users requesting `(testRoom, alice)` get the same peer token, joining as the same peer.
- **No abuse protection.** Your bill is what they make it.

If you suspect a Sandbox URL has leaked, regenerate it from the Sandbox tab. The old URL stops working.

## Migration recipe

Goal: stop calling `useSandbox` and instead fetch a `peerToken` from your own authenticated backend, then pass it to `joinRoom({ peerToken })` (still obtained from `useConnection()`).

### 1. Stand up a backend with the server SDK

Pick `../js-server-sdk/SKILL.md` (Node) or `../python-server-sdk/SKILL.md` (Python). Add a single authenticated endpoint:

```ts
// POST /api/join-room — Express example
app.post('/api/join-room', authenticateUser, async (req, res) => {
  const { roomName } = req.body;
  const user = req.user; // populated by your auth middleware

  // Apply your authorization rules
  if (!userCanJoin(user, roomName)) return res.status(403).end();

  const room = await getOrCreateRoom(roomName);
  const { peer, peerToken } = await fishjamClient.createPeer(room.id, {
    metadata: { userId: user.id, name: user.name },
  });
  res.json({ peerToken });
});
```

```python
# FastAPI example
@router.post('/api/join-room')
async def join_room(req: JoinRoomRequest, user: User = Depends(authenticate_user)):
    if not user_can_join(user, req.room_name):
        raise HTTPException(403)
    room = await get_or_create_room(req.room_name)
    peer, token = fishjam_client.create_peer(
        room.id,
        PeerOptions(metadata={'userId': user.id, 'name': user.name}),
    )
    return {'peerToken': token}
```

### 2. Swap the client call

Before:

```ts
const { getSandboxPeerToken } = useSandbox({ sandboxApiUrl: SANDBOX_API_URL });
const peerToken = await getSandboxPeerToken(roomName, peerName);
await joinRoom({ peerToken });
```

After:

```ts
const res = await fetch('/api/join-room', {
  method: 'POST',
  credentials: 'include',
  headers: { 'Content-Type': 'application/json' },
  body: JSON.stringify({ roomName }),
});
const { peerToken } = await res.json();
await joinRoom({ peerToken });
```

The `FishjamProvider` still gets the `fishjamId` directly — that hasn't changed.

### 3. Remove the Sandbox URL from your frontend config

Delete `VITE_SANDBOX_API_URL`, `NEXT_PUBLIC_SANDBOX_API_URL`, etc. Then disable the Sandbox API from the Dashboard if you no longer need it for dev. (You can re-enable it later for testing.)

### 4. Production checklist

- Backend env vars: `FISHJAM_ID`, `FISHJAM_MANAGEMENT_TOKEN`, never exposed to the bundler.
- HTTPS on the backend; CORS limited to your real origin.
- Your existing auth (JWT, sessions, OAuth — whatever you use) gates `/api/join-room`.
- Rate limiting on `/api/join-room`.
- Webhook receiver (optional) — see `notifier-vs-webhook.md`.
- Token refresh logic if peers may exist longer than 24h — see `auth-model.md`.
- Fishjam plan: **Regular Jar** (paid) is recommended for production. The free **Mini Jar** plan gates production usage (session/peer/duration caps).

## Useful dev pattern: use the Sandbox alongside production

For early development you can keep `useSandbox` behind a feature flag and switch to `useConnection({ peerToken })` for real auth. This avoids deleting your dev path. Once shipping, gate `useSandbox` behind `if (import.meta.env.DEV)` so a Sandbox URL can never leak into a production bundle.

## Source

- <https://fishjam.swmansion.com/docs/explanation/sandbox-api-concept>
- <https://fishjam.swmansion.com/docs/how-to/backend/sandbox-api-testing>
- <https://fishjam.swmansion.com/docs/how-to/backend/production-deployment>
