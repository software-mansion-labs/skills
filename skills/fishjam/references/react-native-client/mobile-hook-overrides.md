# Mobile Hook Overrides

The mobile SDK overrides several hooks from `@fishjam-cloud/react-client` with native-aware implementations. The **API signatures are the same** — call them as you would on web. The behavioral deltas below cover what changes.

> For the call-shape of each hook, read the corresponding reference in `../react-client/`. This file is only the mobile-specific notes.

## `useCamera`

Same API as web (`../react-client/devices.md`). Mobile-specific notes:

- **Permissions must be granted first.** Call `useCameraPermissions().request()` before `useInitializeDevices` (see `permissions.md`).
- **Prefer `selectCamera(deviceId)` to switch cameras.** Constraint-based device selection via the provider may behave differently than on web; the hook's imperative `selectCamera` is the supported path. (Verify against the override implementation in `packages/mobile-client/src/overrides/hooks.ts` for your SDK version.)
- **Front vs back camera.** `cameraDevices` lists both as separate `DeviceItem`s; `selectCamera(deviceId)` is the only public path to swap.
- **`cameraStream` renders only via `RTCView`** (not `<video>`). See `rtcview.md`.

## `useMicrophone`

- **iOS background audio** requires the `audio` background mode in `Info.plist`. The Expo plugin adds `audio` when `ios.supportsPictureInPicture: true`; `ios.enableVoIPBackgroundMode: true` adds `voip` (not `audio`); `ios.enableScreensharing: true` adds neither. There is no dedicated flag for audio-only persistence — enable PiP support (writes `audio` as a side-effect) or add `audio` to `UIBackgroundModes` manually in a bare workflow.
- **`microphoneStream` is not directly playable** the way an HTML `<audio>` element plays a stream. The SDK handles playback of remote audio internally; for local mic, you typically just visualize the levels, not play it back.

## `useScreenShare`

Mobile-only flow detailed in `screen-sharing.md`. Key API delta over the web hook:

- Adds `presentBroadcastPicker: () => Promise<void>` to the destructure — iOS-only; opens the system Stop sheet for an active broadcast (use instead of `stopStreaming` to avoid the error dialog the platform pops on force-close).
- On iOS, `startStreaming()` internally presents the system broadcast picker — no separate component to mount.
- On Android, `startStreaming` triggers the MediaProjection consent prompt and requires `useForegroundService({ enableScreenSharing: true })` to be active.
- `setTracksMiddleware` is supported; the `TracksMiddleware` shape is mobile-specific (operates on `@fishjam-cloud/react-native-webrtc` `MediaStreamTrack`s — see `packages/mobile-client/src/overrides/types.ts`).

## `useCustomSource`

- Same `(sourceId)` signature. The mobile version accepts `MediaStream` from any source that produces native streams (e.g. canvas streams from a React Native canvas library, or custom native sources).
- No browser `<canvas>.captureStream()` — to produce a custom stream, use a library like `react-native-skia` or a custom native module that exposes a `MediaStream`.

## `useInitializeDevices`

- `InitializeDevicesSettings` is `{ enableVideo?: boolean; enableAudio?: boolean }`. There is no `enableSimulcast` field — simulcast is configured at provider level via `<FishjamProvider videoConfig={{ sentQualities: [...] | false }}>` (or `audioConfig`). `useCamera()` does not expose a `startStreaming` method, so don't try to pass `StreamConfig` there. The mobile encoder may decline simulcast if hardware doesn't support it.
- Calling `useInitializeDevices` triggers OS-level permission prompts if you haven't already requested them. Best practice: request permissions explicitly first, then init.

## `useLivestreamStreamer`

- Same API. The streamer's video / audio inputs are `MediaStream` objects you destructure from the device hooks: `useCamera().cameraStream`, `useMicrophone().microphoneStream`, or `useScreenShare().stream`.
- iOS background streaming requires CallKit-style background modes or the broadcast-extension flow for screen-share livestreams.

## `useLivestreamViewer`

- Same API. The received `stream` is a `MediaStream` you render with `RTCView`:

```tsx
<RTCView mediaStream={stream} style={{ flex: 1 }} />
```

- WHEP playback works the same as web — viewer tokens, public livestreams, etc.

## `usePeers`

- Same `PeerWithTracks` shape. Each `Track` (and `RemoteTrack`) has `stream: MediaStream | null`.
- Render with `RTCView` instead of `<video>`. Mirror the local peer's camera tile for usability.

## `useCallKit`, `useCallKitEvent`, `useCallKitService`

These have **no web equivalent** — see `callkit.md`.

## `useAudioOutput`

Mobile-only — see `audio-output.md`.

## `useCameraPermissions`, `useMicrophonePermissions`

Mobile-only — see `permissions.md`. Note the API is `[query, request]` **tuple**, unlike the imperative `request().status` shape some other libraries use.

## `useForegroundService`

Mobile-only (Android-only effective) — see `foreground-service.md`.

## Things the mobile SDK does NOT expose

- `persistLastDevice` — the mobile `FishjamProvider` omits this prop.
- `fishjamClient` override — the mobile `FishjamProvider` constructs its own client with `clientType: 'mobile'` and doesn't accept an override.
- `useStatistics`, `useLocalVAD` — not in the public surface on web or mobile.

## Sources

- `packages/mobile-client/src/index.ts` (exports + re-exports)
- `packages/mobile-client/src/overrides/hooks.ts`
- `packages/mobile-client/src/overrides/types.ts`
- For each hook's full API: corresponding reference in `../react-client/`.
