# Publishing

## Publisher

```ts
const publisher = usePublisher(session);
publisher.publish({ path: 'live/test', tracks: [camera, microphone] }); // requires session.state === 'connected'
publisher.stop();
```

- `PublisherState`: `'idle' | 'connecting' | 'publishing' | 'stopped' | 'error:${string}'`. Per-track: `publisher.trackStates[name]` (`'idle' | 'starting' | 'active' | 'stopped'`), plus `lastError`.
- Events: `stateChange`, `trackStateChange` (`{ name, state, error? }`).
- The track list is **snapshotted at `publish()`** — to add/remove tracks, call `publish()` again (restarts the broadcast; MoQ finalizes tracks at start).
- Publishing without `connected` fails immediately with `error:session is not connected`.

## Camera

```tsx
const camera = useCamera({ position: 'front', width: 1280, height: 720, framerate: 30, videoCodec: 'h264', enabled: true });
<PublisherView camera={camera} style={styles.preview} />
```

- `CameraTrack`: `state` (`'idle' | 'starting' | 'active' | 'error:...'`), `lastError`, `position`, `encoder`, `flip()`, `setPosition('front' | 'back')`.
- `<PublisherView>` shows the local preview; mounting it does **not** start/stop capture — the hook's `enabled` does.
- Capture hooks are **refcounted device singletons**: all consumers share one physical device; `flip()` is visible everywhere. Keep hooks unconditional and toggle `enabled`.
- Permissions (camera + mic) are the host app's responsibility. Black preview → missing permission.

## Multi-camera (front + back simultaneously)

```ts
const supported = await isMultiCameraSupported();
const multi = useMultiCamera({ enabled: dualCam }); // portrait defaults 720x1280@30
publisher.publish({ path, tracks: [multi.front, multi.back, microphone] });
```

- `MultiCameraTrack`: `isSupported` (null while probing), `state`, `front` / `back` (published as `front-camera` / `back-camera`). `flip()`/`setPosition()` on the sub-tracks are no-ops.
- Use `h264` — a second concurrent H.265 encoder silently produces no frames on many devices.

## Microphone

```ts
const microphone = useMicrophone({ audioCodec: 'opus', audioSampleRate: 48000, enabled: micOn });
```

- `audioSampleRate` defaults to 48000 (Opus's native rate); it's snapshotted at capture start — remount to change.
- iOS: capturing switches the audio session to `playAndRecord`; an idle capture keeps holding it and can block other audio libraries with `insufficientPriority` (OSStatus 561017449). Gate `enabled` on actually publishing.

## Codec support

```ts
getSupportedVideoCodecs(); // ['h264', 'h265'] — synchronous
getSupportedAudioCodecs(); // ['opus', 'aac']
```

Always gate non-default codecs on these. On Android an unsupported encoder makes publishing start then **silently stop with no error** (reported as a clean stop). Defaults: `h264` / `opus`.

## Screen broadcasting

Out-of-process capture (iOS Broadcast Upload Extension, Android foreground Service). It opens its **own** MoQ connection using the session's URL — the host session only supplies the URL.

```tsx
const screen = useScreenBroadcast(session, {
  path: 'live/screen',
  appGroupIdentifier: 'group.com.example.app', // iOS: required
});
// iOS: the user starts it from the picker (start() always rejects there)
<BroadcastPickerView preferredExtension="com.example.app.broadcast" />;
```

- `ScreenBroadcast`: `state` (`'idle' | 'connecting' | 'broadcasting' | 'stopped' | 'error:...'`), `lastError`, `start(): Promise<void>`, `stop()`.
- **iOS:** `start()` always rejects — the user must tap `<BroadcastPickerView preferredExtension={...} />` (renders a plain View on Android). Requires an App Group entitlement on host + extension and a Broadcast Upload Extension subclassing `MoQReplayKitBroadcastSampleHandler`. Options: `appAudio` (default true), `mic` (default true).
- **Android:** `start()` shows the MediaProjection consent dialog. Requires `FOREGROUND_SERVICE` + `FOREGROUND_SERVICE_MEDIA_PROJECTION` permissions.
- Encoder options (`videoCodec`, `width`, `height`, `framerate`, `audioCodec`, `audioSampleRate`) mirror camera/mic defaults.

## Full go-live pattern

```tsx
const tracks: PublishTrack[] = [];
if (cameraOn) tracks.push(dualCam ? multi.front : camera, ...(dualCam ? [multi.back] : []));
if (micOn) tracks.push(microphone);
if (captions) tracks.push(dataTrack); // see custom-tracks.md
publisher.publish({ path, tracks });
```
