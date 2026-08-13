# Livestream Hooks

Two hooks, both used outside (not in addition to) the normal `useConnection` flow:

- `useLivestreamStreamer` — broadcast into a `livestream` room.
- `useLivestreamViewer` — consume a livestream via WHEP.

> Tokens come from the backend (`../js-server-sdk/livestream-and-moq.md` → `createLivestreamStreamerToken` / `createLivestreamViewerToken`; same for Python in `../python-server-sdk/livestream-and-moq.md`), or from the Sandbox API via `useSandbox().getSandboxLivestream` / `getSandboxViewerToken`.

## `useLivestreamStreamer` (broadcaster)

Source: `@fishjam-cloud/react-client` → `packages/react-client/src/hooks/useLivestreamStreamer.ts`.

```tsx
import { useLivestreamStreamer, useCamera, useMicrophone } from '@fishjam-cloud/react-client';

const { connect, disconnect, error, isConnected } = useLivestreamStreamer();
const { cameraStream } = useCamera();
const { microphoneStream } = useMicrophone();

const handleStart = async () => {
  await connect({
    inputs: {
      video: cameraStream!,        // MediaStream
      audio: microphoneStream,     // MediaStream | null
    },
    token: streamerToken,           // from backend
  });
};

const handleStop = () => disconnect();
```

### Inputs

```ts
type StreamerInputs =
  | { video: MediaStream; audio?: MediaStream | null }
  | { video?: null; audio: MediaStream };

type ConnectStreamerConfig = {
  inputs: StreamerInputs;
  token: string;
};
```

- **Exactly one** video track and **exactly one** audio track go on the wire — Fishjam picks the first of each from the streams you pass. Extra tracks are silently ignored.
- Audio-only stream: pass `{ audio: micStream }` without video.
- You can also publish screen share as the video source: `{ video: screenShareStream }`.

### Re-publish

Calling `connect` again **replaces** the current publication with the new inputs — useful for switching from camera to screen share without tearing down the connection. Only the latest `connect` call wins.

### Streamer error states

`error: LivestreamError | null`. The enum (from `@fishjam-cloud/ts-client`):

- `LivestreamError.UNAUTHORIZED` — bad / expired / wrong-role streamer token.
- `LivestreamError.STREAM_NOT_FOUND` — the room/stream doesn't exist.
- `LivestreamError.STREAMER_ALREADY_CONNECTED` — another streamer is currently publishing into this livestream room.
- `LivestreamError.UNKNOWN_ERROR` — fallback for anything else (network, TURN, server).

`isConnected` flips `true` once the WebRTC connection state hits `connected`.

## `useLivestreamViewer` (audience)

Source: `@fishjam-cloud/react-client` → `packages/react-client/src/hooks/useLivestreamViewer.ts`.

```tsx
import { useLivestreamViewer } from '@fishjam-cloud/react-client';

const { stream, connect, disconnect, error, isConnected, getStatistics } = useLivestreamViewer();

// Private livestream — backend-issued viewer token:
await connect({ token: viewerToken });

// Public livestream — just the room ID:
await connect({ streamId: 'room-abc-123' });
```

### `ConnectViewerConfig`

```ts
type ConnectViewerConfig =
  | { token: string;  streamId?: never }
  | { streamId: string; token?: never };
```

Pass exactly one — TypeScript will complain if you mix.

### Rendering

```tsx
<video
  ref={(el) => { if (el) el.srcObject = stream; }}
  autoPlay
  playsInline
/>
```

`stream` populates once the WHEP handshake completes. Until then it's `null`.

### Statistics

```ts
const stats: RTCStatsReport | undefined = await getStatistics();
```

Returns a raw `RTCStatsReport` from the underlying `RTCPeerConnection` — use for QoS overlays.

### Viewer error states

Same `LivestreamError` enum as the streamer. Public-livestream-only errors come through `error` too (e.g. `streamId` doesn't exist).

## Single-track constraint

Livestream rooms accept **exactly one video track** and **one audio track from the streamer**. Multiple-camera scenarios aren't supported on livestream — for that, use a `conference` room.

## Public vs private — choosing

Decision lives on the **room** at creation (`fishjam_client.create_room(RoomOptions(room_type='livestream', public=True))`). The client side just adapts:

| Room is | Streamer | Viewer |
|---|---|---|
| Private (default) | streamer token | viewer token |
| Public | streamer token | room ID |

The streamer flow is identical — only the viewer differs.

## WHEP / WHIP

Under the hood the streamer uses WHIP (`buildLivestreamWhipUrl(fishjamId)`) and the viewer uses WHEP (`buildLivestreamWhepUrl(fishjamId)`). You don't need to know these URLs unless you're integrating with a non-Fishjam WHEP client — in which case the URL is `https://<fishjam-host>/api/v1/live/api/whep[/<streamId>]` (or `.../whip` for the streamer). See `packages/react-client/src/utils/fishjamUrl.ts` for the builder.

## Sources

- `@fishjam-cloud/react-client` → `packages/react-client/src/hooks/useLivestreamStreamer.ts`
- `packages/react-client/src/hooks/useLivestreamViewer.ts`
- `packages/ts-client/src/livestream.ts` (`LivestreamError` enum)
- <https://fishjam.swmansion.com/docs/tutorials/livestreaming>
- <https://fishjam.swmansion.com/docs/how-to/backend/whip-whep>
- `examples/web-react/livestreaming/` in the SDK monorepo
