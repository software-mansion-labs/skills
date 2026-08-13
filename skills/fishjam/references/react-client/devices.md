# Devices

Four hooks. All require `FishjamProvider` higher in the tree.

## useInitializeDevices

Call **once** on app mount (or on first navigation to a page that needs media). Prompts the browser for camera + mic permissions, enumerates devices, and sets up the persisted last-device.

```tsx
import { useInitializeDevices } from '@fishjam-cloud/react-client';
import { useEffect } from 'react';

function MediaBootstrap() {
  const { initializeDevices } = useInitializeDevices();

  useEffect(() => {
    initializeDevices({
      enableAudio: true,
      enableVideo: true,
    });
  }, []);

  return null;
}
```

`InitializeDevicesSettings`:

- `enableAudio?: boolean` — request mic
- `enableVideo?: boolean` — request camera

Simulcast layers are configured separately via `videoConfig` on `FishjamProvider` — not here. See `simulcast-and-bandwidth.md`.

This hook can batch camera + microphone permission requests into a single browser popup if both are enabled simultaneously.

Returns an `InitializeDevicesResult` describing what was granted; status enum (`InitializeDevicesStatus`) shows progress.

Browser permission prompt timing: this is when the user sees the permission dialog. Trigger off a user gesture if you want predictable UX (don't auto-request on mount in a tab the user just opened — Safari sometimes silently denies).

## useCamera

```tsx
import { useCamera } from '@fishjam-cloud/react-client';

const {
  cameraStream,            // MediaStream | null
  isCameraOn,              // boolean
  currentCamera,           // MediaDeviceInfo | null (selected device)
  cameraDevices,           // MediaDeviceInfo[]
  cameraDeviceError,       // DeviceError | null
  toggleCamera,            // () => void
  startCamera,             // () => void
  stopCamera,              // () => void
  selectCamera,            // (deviceId: string) => void
  currentCameraMiddleware, // TrackMiddleware | null
  setCameraTrackMiddleware,// (mw: TrackMiddleware | null) => void
} = useCamera();
```

Render with:

```tsx
<video
  ref={(el) => { if (el) el.srcObject = cameraStream; }}
  autoPlay
  muted
  playsInline
/>
```

Device switching:

```tsx
<select onChange={(e) => selectCamera(e.target.value)}>
  {cameraDevices.map((d) => <option key={d.deviceId} value={d.deviceId}>{d.label}</option>)}
</select>
```

Middleware lets you wrap the raw `MediaStreamTrack` — for filters, blur, watermarks. See `custom-sources.md`.

## useMicrophone

Same API shape as `useCamera` but for audio:

```tsx
const {
  microphoneStream,
  isMicrophoneOn,
  currentMicrophone,
  microphoneDevices,
  microphoneDeviceError,
  toggleMicrophone,
  startMicrophone,
  stopMicrophone,
  selectMicrophone,
  currentMicrophoneMiddleware,
  setMicrophoneTrackMiddleware,
  toggleMicrophoneMute,           // () => void          — toggles track.enabled, keeps publish active
  isMicrophoneMuted,              // boolean             — current mute state
} = useMicrophone();
```

## useScreenShare

```tsx
import { useScreenShare } from '@fishjam-cloud/react-client';

const {
  startStreaming,             // (constraints?) => Promise<void>
  stopStreaming,              // () => Promise<void>
  stream,                     // MediaStream | null     (combined video + audio)
  videoTrack,                 // MediaStreamTrack | null
  audioTrack,                 // MediaStreamTrack | null (system audio, if user opted in)
  currentTracksMiddleware,    // TracksMiddleware | null
  setTracksMiddleware,        // (mw: TracksMiddleware | null) => Promise<void>
} = useScreenShare();
```

**User gesture required.** Call `startStreaming` inside a click handler — browsers reject `getDisplayMedia` calls that aren't reacting to a user event.

```tsx
<button onClick={() => startStreaming()}>Share screen</button>
```

Pass per-kind constraints to override defaults (e.g. preferring tab capture, including system audio):

```tsx
await startStreaming({
  videoConstraints: { displaySurface: 'browser' },
  audioConstraints: true,     // ask for system audio
});
```

Both `videoConstraints` and `audioConstraints` accept `boolean | MediaTrackConstraints`.

Render the combined stream by assigning it to a `<video>` element's `srcObject` via a ref (same pattern as `useCamera` — `<video ref={(el) => { if (el) el.srcObject = stream; }} autoPlay playsInline />`); JSX won't accept `srcObject` as a prop. Or grab `videoTrack` / `audioTrack` individually if you need to layer them separately. Apply filters / transforms with `setTracksMiddleware(...)` — distinct from camera middleware because screen share carries two tracks at once. See `custom-sources.md`.

## Device persistence

By default `FishjamProvider` stores the last-used camera and mic in `localStorage`. After a page reload, `useCamera().currentCamera` resolves to the previously-used device automatically — but the browser still needs the user to grant the permission.

Override the storage via `persistLastDevice` on the provider (see `provider.md`).

## Common errors

`cameraDeviceError` and `microphoneDeviceError` are `DeviceError | null` objects shaped as `{ name: ... }`. Compare on `.name`:

- `cameraDeviceError?.name === 'NotAllowedError'` — user denied permission. Show UI explaining how to re-grant.
- `cameraDeviceError?.name === 'NotFoundError'` — no camera attached.
- `cameraDeviceError?.name === 'OverconstrainedError'` — requested constraints can't be satisfied.
- `cameraDeviceError?.name === 'UNHANDLED_ERROR'` — fallback for anything else.

Same set for microphone.

## Sources

- `@fishjam-cloud/react-client` → `packages/react-client/src/hooks/devices/useCamera.ts`
- `packages/react-client/src/hooks/devices/useMicrophone.ts`
- `packages/react-client/src/hooks/devices/useInitializeDevices.ts`
- `packages/react-client/src/hooks/useScreenShare.ts`
- `packages/react-client/src/hooks/internal/useScreenshareManager.ts` (`startStreaming` param shape)
- `packages/react-client/src/types/public.ts` (`DeviceError`)
- <https://fishjam.swmansion.com/docs/how-to/client/managing-devices>
- <https://fishjam.swmansion.com/docs/how-to/client/screensharing>
