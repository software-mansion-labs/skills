---
name: fishjam-react-client
description: "Browser-only React SDK for Fishjam — joining rooms, capturing camera/microphone/screen, displaying peers, and acting as a livestream streamer or viewer in a React web app. Use whenever the user is writing a React app in a browser that calls Fishjam APIs, sets up FishjamProvider, or uses any Fishjam React hook. Trigger on: '@fishjam-cloud/react-client', 'FishjamProvider', 'useConnection', 'useCamera', 'useMicrophone', 'useScreenShare', 'usePeers', 'useDataChannel', 'useVAD', 'useLivestreamStreamer', 'useLivestreamViewer', 'useCustomSource', 'useInitializeDevices', 'useUpdatePeerMetadata', 'useSandbox', 'PeerWithTracks', 'joinRoom', 'peerToken', 'fishjamId', 'fishjam react', '@fishjam-cloud/ts-client', 'FishjamClient ts-client'. Covers the provider, the full hook catalog, simulcast configuration, custom sources, data channels, VAD, livestream WHEP playback, device persistence, and reconnection. Briefly notes when to drop down to @fishjam-cloud/ts-client for non-React or worker contexts."
license: MIT
---

# Fishjam React Client

`@fishjam-cloud/react-client` — Fishjam's React SDK for **browser** apps. Wrap your app in `FishjamProvider`, then use hooks to drive devices, connection, and peers.

> **Read `../platform/SKILL.md` first.** It defines what rooms, peers, tokens, and the sandbox-vs-prod model are. This skill assumes those.
>
> **For React Native, use `../react-native-client/SKILL.md` instead** — it re-exports this hook surface and adds mobile-only concerns.

## Minimal setup

```bash
npm install @fishjam-cloud/react-client
# or yarn add @fishjam-cloud/react-client
```

```tsx
import { FishjamProvider, useConnection, useSandbox } from '@fishjam-cloud/react-client';

function App() {
  return (
    <FishjamProvider fishjamId={import.meta.env.VITE_FISHJAM_ID}>
      <Room />
    </FishjamProvider>
  );
}

function Room() {
  // DEV — sandbox flow (no backend needed):
  const { getSandboxPeerToken } = useSandbox({ sandboxApiUrl: import.meta.env.VITE_SANDBOX_API_URL });
  const { joinRoom } = useConnection();

  const handleJoin = async () => {
    const peerToken = await getSandboxPeerToken('test-room', 'alice');
    await joinRoom({ peerToken });
  };

  return <button onClick={handleJoin}>Join</button>;
}
```

For **production**, replace `useSandbox` with a call to your own backend that returns a peer token (see `../platform/sandbox-vs-production.md`):

```tsx
const peerToken = await fetchPeerTokenFromYourBackend(roomName);
await joinRoom({ peerToken });
```

## Hook catalog (decision tree)

| What you need | Hook | Reference |
|---|---|---|
| Open / close the connection | `useConnection` | `connection.md` |
| Sandbox tokens for prototyping | `useSandbox` | `connection.md` |
| Camera | `useCamera` | `devices.md` |
| Microphone | `useMicrophone` | `devices.md` |
| Screen share | `useScreenShare` | `devices.md` |
| Initialize devices on app load | `useInitializeDevices` | `devices.md` |
| Render peers + their tracks | `usePeers` (returns `PeerWithTracks`) | `peers-and-tracks.md` |
| Livestream as broadcaster | `useLivestreamStreamer` | `livestream.md` |
| Livestream as viewer | `useLivestreamViewer` | `livestream.md` |
| Push a `MediaStreamTrack` you generated yourself | `useCustomSource` | `custom-sources.md` |
| Peer-to-peer messages (text chat, signals) | `useDataChannel` | `data-and-events.md` |
| Detect when a remote peer is speaking | `useVAD` | `data-and-events.md` |
| Update your peer's metadata mid-session | `useUpdatePeerMetadata` | `data-and-events.md` |
| Get WebRTC `RTCStatsReport` for debugging | `useStatistics` (from `@fishjam-cloud/react-client/debug`) | `provider.md` |
| Configure simulcast / bandwidth | (props on `FishjamProvider`) | `simulcast-and-bandwidth.md` |
| Use Fishjam from non-React code (workers, Svelte, vanilla TS) | `@fishjam-cloud/ts-client` | `ts-client-escape.md` |

## Key rules

- **One `FishjamProvider` per app.** Don't nest. It owns the `FishjamClient`, all device managers, and reconnection state.
- **`useSandbox` is dev-only.** Gate it with `if (import.meta.env.DEV)` (or equivalent) so a sandbox URL can never ship in a production bundle.
- **Browser requires user gesture for `useScreenShare`** — call inside an `onClick`, not on mount.
- **Devices stay mounted across connection state.** Initializing with `useInitializeDevices` once is enough — they remain available after you `leaveRoom`.
- **Livestream rooms have a one-track-per-direction constraint.** One video + one audio on the streamer; viewers don't publish.
- **For non-React contexts** (Svelte, Vue, workers, vanilla TS), use `@fishjam-cloud/ts-client` directly. The React SDK is a thin layer over it. See `ts-client-escape.md`.
- **The provider re-exports `Variant` from `ts-client`** for simulcast quality (`Variant.VARIANT_LOW | Variant.VARIANT_MEDIUM | Variant.VARIANT_HIGH`).

## References

| File | When to read |
|---|---|
| `provider.md` | All `FishjamProvider` props — `fishjamId`, `reconnect`, `constraints`, `persistLastDevice`, `bandwidthLimits`, `videoConfig`, `audioConfig`, `debug`, `fishjamClient`. |
| `connection.md` | `useConnection` (`joinRoom`, `leaveRoom`, `peerStatus`, `reconnectionStatus`) and `useSandbox`. |
| `devices.md` | `useCamera`, `useMicrophone`, `useScreenShare`, `useInitializeDevices`. Device enumeration, switching, persistence. |
| `peers-and-tracks.md` | `usePeers` — `localPeer`, `remotePeers`, the `PeerWithTracks` shape. |
| `livestream.md` | `useLivestreamStreamer`, `useLivestreamViewer`, public vs private flows. |
| `custom-sources.md` | `useCustomSource`, `TrackMiddleware`, `TracksMiddleware`. |
| `data-and-events.md` | `useDataChannel`, `useVAD`, `useUpdatePeerMetadata`. |
| `simulcast-and-bandwidth.md` | `Variant`, `BandwidthLimits`, `SimulcastConfig` via `videoConfig` / `audioConfig`. |
| `ts-client-escape.md` | When and how to drop to `@fishjam-cloud/ts-client`. |
