---
name: react-native-moq
description: "Software Mansion's react-native-moq — Media over QUIC (MoQ) live streaming for React Native: sub-second-latency video/audio playback, camera/mic/screen publishing, and realtime data tracks over a MoQ relay. MUST USE before writing, reviewing, or debugging ANY code that uses react-native-moq or react-native-moq-ui. Trigger on: 'react-native-moq', 'MoQ', 'Media over QUIC', 'moq relay', 'moq-kit', 'useSession', 'useBroadcasts', 'useVideoPlayer', 'useAudioPlayer', 'usePublisher', 'useCamera', 'useMultiCamera', 'useMicrophone', 'useScreenBroadcast', 'useDataTrack', 'useDataMessages', 'useAudioChunks', 'useAudioSource', 'useVideoSource', 'VideoView', 'PublisherView', 'VideoPlayerView', 'BroadcastPickerView', 'createSession', 'createPublisher', 'subscribeBroadcasts'."
license: MIT
---

# react-native-moq

React Native bindings for Media over QUIC (MoQ): live video, audio, and data tracks with sub-second latency, published to and consumed from a MoQ relay. Built on Software Mansion's moq-kit native library.

**Requirements:** React Native New Architecture (Fabric + TurboModules), iOS 16+, Android API 30+.

```sh
npm install react-native-moq
cd ios && pod install
```

Optional player chrome (fullscreen modal, controls, volume slider):

```sh
npm install react-native-moq-ui @react-native-vector-icons/material-icons
```

iOS also needs `MaterialIcons.ttf` under `UIAppFonts` in Info.plist, then `pod install` again. Skip this package if composing your own player UI on `<VideoView>`.

## Mental model

- A **Session** is one connection to a MoQ relay (`useSession(url)`). It does **not** auto-connect — call `session.connect()` yourself. One session serves both subscribing and publishing.
- **Subscribing:** `useBroadcasts(session, prefix)` lists live broadcasts under a path prefix → pick a `BroadcastInfo` → `useVideoPlayer` / `useAudioPlayer` → render with `<VideoView player={player}>`. Raw-media alternatives: `useAudioChunks` (encoded or PCM to JS) and `useDataMessages` (string payloads).
- **Publishing:** capture hooks produce **tracks** (`useCamera`, `useMicrophone`, `useMultiCamera`) or you generate them (`useDataTrack`, `useAudioSource`, `useVideoSource`) → `usePublisher(session).publish({ path, tracks })`. Screen sharing (`useScreenBroadcast`) is separate and out-of-process.
- Every hook has a hook-free counterpart — `create*` (returns the object plus `destroy()`) for the stateful ones, `subscribe*` (returns `start()`/`stop()`) for the subscription hooks — see [references/imperative.md](references/imperative.md).
- Objects expose reactive state (`session.state`, `publisher.state`, `player.isPlaying`, …); hooks re-render on changes. Subscribe manually with `addListener` or the `useEvent`/`useEventListener` helpers.

## Quick start — watch

```tsx
function Watch() {
  const session = useSession('https://relay.example.com', (s) => s.connect());
  const broadcasts = useBroadcasts(session, ''); // '' = all paths
  return broadcasts.map((b) => <BroadcastPlayer key={b.path} broadcast={b} />);
}

function BroadcastPlayer({ broadcast }: { broadcast: BroadcastInfo }) {
  const player = useVideoPlayer(broadcast, (p) => p.play());
  return <VideoView player={player} style={{ width: '100%', aspectRatio: 16 / 9 }} />;
}
```

`useVideoPlayer`/`useAudioPlayer` require a non-null broadcast — mount them in a child rendered only once one exists.

## Quick start — go live

```tsx
const session = useSession(url, (s) => s.connect());
const camera = useCamera();
const microphone = useMicrophone();
const publisher = usePublisher(session);

<PublisherView camera={camera} style={styles.preview} />
// when session.state === 'connected':
publisher.publish({ path: 'live/test', tracks: [camera, microphone] });
```

## Critical rules

1. **Errors are strings, not exceptions.** States are unions like `'connected' | 'error:...'`; check `state.startsWith('error:')`. Nothing throws on media failure. Sessions carry the message **only** in the state string (`state.slice(6)`); capture tracks and the publisher additionally expose a `lastError` field. Players have no error surface at all — watch session and publisher state instead.
2. **`publish()` requires `session.state === 'connected'`**, otherwise the publisher goes straight to `error:session is not connected`. Gate the go-live button on session state.
3. **Gate codecs with `getSupportedVideoCodecs()` / `getSupportedAudioCodecs()`.** On Android an unsupported encoder (often `h265`) makes publishing start then silently stop with no error. Default to `h264`/`opus`.
4. **Permissions are the host app's job** — request CAMERA/RECORD_AUDIO on Android, add NSCameraUsageDescription/NSMicrophoneUsageDescription on iOS. A black preview usually means a missing permission.
5. **`<VideoView>` takes no children** (the Android native view is not a ViewGroup). Render overlays as absolutely-positioned siblings.
6. **Capture hooks are refcounted device singletons** — all consumers share the physical camera/mic; `flip()` affects everyone. Use `enabled: false` to keep hardware off while mounted instead of calling hooks conditionally.
7. **Broadcast paths are relative to the subscription prefix** (broadcast `live/test` appears as `path: 'test'` under prefix `'live'`). Subscribe with `''` to see full paths. A listed broadcast that won't play (relay error `code=13` NotFound) means a path/prefix mismatch.
8. **Data tracks are not in the catalog** — publisher and subscriber must agree on the track name out of band (default `'data'`).
9. An idle iOS microphone capture holds the audio session in `playAndRecord` and can block other audio libraries (`insufficientPriority`). Disable the mic (`enabled: false`) when not publishing.

## References

| Task | Read |
| --- | --- |
| Sessions, broadcast discovery, video/audio playback, `<VideoView>`, stats, track switching, events, `react-native-moq-ui` | [references/playback.md](references/playback.md) |
| Camera, multi-camera, microphone, publisher lifecycle, screen broadcasting, codec queries | [references/publishing.md](references/publishing.md) |
| Data tracks & messages, audio chunks to JS, push-your-own audio (PCM) and video (frame buffers) | [references/custom-tracks.md](references/custom-tracks.md) |
| Using without React: `create*` handles, `subscribe*` functions, cleanup | [references/imperative.md](references/imperative.md) |
