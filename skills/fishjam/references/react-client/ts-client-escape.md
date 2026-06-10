# Dropping to `@fishjam-cloud/ts-client`

`@fishjam-cloud/react-client` is a thin React wrapper over `@fishjam-cloud/ts-client`. If you're not in React (Svelte, Vue, vanilla TS in a Vite app, a web worker, a Node-style headless harness), use `ts-client` directly.

## Install

```bash
npm install @fishjam-cloud/ts-client
```

## Construct a client

```ts
import { FishjamClient } from '@fishjam-cloud/ts-client';

const client = new FishjamClient({
  reconnect: true,    // or false / ReconnectConfig
  debug: false,
});
```

This is the same class that `FishjamProvider` constructs internally. Pass your own instance back into the provider via `fishjamClient={...}` if you also need it from React.

## Connect

```ts
// Pass the host root only — FishjamClient.connect appends `/socket/peer/websocket` itself.
const fishjamUrl = `wss://fishjam.io/api/v1/connect/${FISHJAM_ID}`;
// (or build via the same `httpToWebsocketUrl(resolveFishjamUrl(fishjamId))` logic the React hook uses)

await client.connect({
  url: fishjamUrl,
  token: peerToken,
  peerMetadata: { name: 'alice' },
});
```

Status changes via events:

```ts
client.on('joined', () => console.log('connected'));
client.on('disconnected', () => console.log('disconnected'));
client.on('authError', (reason) => console.error('auth error', reason));
```

## React-equivalent operations

| React hook | ts-client equivalent |
|---|---|
| `useConnection().joinRoom` | `client.connect({ url, token, peerMetadata })` |
| `useConnection().leaveRoom` | `client.disconnect()` |
| `usePeers().remotePeers` | `client.getRemotePeers()` and `'trackReady' / 'trackRemoved'` events |
| `useCamera().cameraStream` | Construct your own `MediaStream`; publish via `client.addTrack(track, metadata, simulcastConfig)` |
| `useScreenShare().startStreaming` | `await navigator.mediaDevices.getDisplayMedia(...)` then `client.addTrack(...)` |
| `useDataChannel().publishData` | `client.publishData(data, options)` after `client.createDataChannels()` |
| `useDataChannel().subscribeData` | `client.subscribeData(callback, options)` |
| `useLivestreamStreamer().connect` | `import { publishLivestream } from '@fishjam-cloud/ts-client'` |
| `useLivestreamViewer().connect` | `import { receiveLivestream } from '@fishjam-cloud/ts-client'` |

The React hooks just translate these primitives into reactive state. The semantics are identical.

## Minimal vanilla TS app

```ts
import { FishjamClient, publishLivestream, receiveLivestream } from '@fishjam-cloud/ts-client';

async function main() {
  const client = new FishjamClient({ reconnect: true });

  client.on('joined', () => console.log('joined'));
  client.on('trackReady', (trackContext) => {
    const video = document.createElement('video');
    video.srcObject = trackContext.stream;
    video.autoplay = true;
    document.body.appendChild(video);
  });

  await client.connect({
    url: `wss://fishjam.io/api/v1/connect/<FISHJAM_ID>`, // no trailing path — client appends /socket/peer/websocket
    token: peerToken,
    peerMetadata: {},
  });

  // Get media:
  const stream = await navigator.mediaDevices.getUserMedia({ video: true, audio: true });
  for (const track of stream.getTracks()) {
    const type = track.kind === 'video' ? 'camera' : 'microphone';
    await client.addTrack(track, { type, paused: false });
  }
}
```

See `examples/web-ts/minimal/` and `simple-app/` in the SDK monorepo for full working examples.

## When NOT to drop to ts-client

If your app is React, **stay in React**. The hooks save you from event-listener spaghetti, device persistence boilerplate, and reactivity glue. The "escape hatch" is for genuinely non-React code paths.

You can also do **both** — use the React hooks for UI and grab the client via `fishjamClient` prop on the provider for one specific advanced flow.

## Sources

- `@fishjam-cloud/ts-client` exports (re-exported types include `AuthErrorReason`, `JoinErrorReason`, `Metadata`, `ReconnectConfig`, `ReconnectionStatus`, `SimulcastConfig`, `SimulcastBandwidthLimit`, `TrackBandwidthLimit`, `Variant`)
- `examples/web-ts/minimal/src/main.ts` in the SDK monorepo (canonical vanilla example)
- `examples/web-ts/simple-app/` in the SDK monorepo (multi-track example)
- The React client source — every hook is a thin wrapper that calls `ts-client` methods
