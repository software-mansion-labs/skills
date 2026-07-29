# Playback

## Session

```ts
const session = useSession(url, (s) => s.connect()); // setup runs once per session
```

- `Session`: `id`, `url`, `state`, `connect(targetLatencyMs?)` (default 200 ms), `disconnect()`, `addListener('stateChange', ...)`.
- `SessionState`: `'idle' | 'connecting' | 'connected' | 'closed' | 'error:${string}'`.
- No auto-connect — call `connect()` in `setup` or on user action. `url` is read at `connect()` time — changing it doesn't recreate the session, it applies on the next `connect()`.

## Discovering broadcasts

```ts
const broadcasts = useBroadcasts(session, 'live'); // prefix; '' matches all
```

- Returns `BroadcastInfo[]`: `{ sessionId, path, videoTracks, audioTracks, player }` — `player` is the broadcast's shared `PlayerHandle`, which `useVideoPlayer`/`createVideoPlayer` wrap. Empty until `state === 'connected'`, repopulates on reconnect, `[]` on disconnect.
- `path` is **relative to the prefix**: broadcast `live/test` → `path: 'test'` under prefix `'live'`, `''` under prefix `'live/test'`.
- Track metadata: `VideoTrackInfo { name, codec, width?, height?, framerate?, bitrate? }`, `AudioTrackInfo { name, codec, sampleRate, channelCount?, bitrate? }`.
- Overlapping prefixes matching the same broadcast are unsupported.

## Video playback

```tsx
const player = useVideoPlayer(broadcast, (p) => p.play());
return <VideoView player={player} style={{ width: '100%', aspectRatio: 16 / 9 }} />;
```

`Player` API:

- State: `isPlaying`, `playbackStats` (null until first update), `currentVideoTrackName`, `currentAudioTrackName`, `volume` (0–1, default 1).
- Methods: `play()`, `pause()`, `stop()`, `setVolume(v)`, `updateTargetLatency(ms)`, `switchVideoTrack(name)`, `switchAudioTrack(name)`.
- Events: `playingChange` (deduped), `trackSwitched` (`{ trackKind, trackName }`), `trackStopped`, `statsUpdate` (periodic).

`<VideoView player={player} style={...} />`:

- **No children** — the Android native view is not a ViewGroup. Overlays (captions, badges) go in absolutely-positioned sibling views.
- Size it yourself (e.g. `aspectRatio`); the native surface fills the view.

## Audio-only playback

```ts
const player = useAudioPlayer(broadcast, (p) => p.play());
```

- `AudioPlayer` = `Player` minus video fields. Uses a separate native player (never subscribes the video track) — cheaper than a hidden `<VideoView>`.
- Type error by design if passed to `<VideoView>` / `<VideoPlayerView>`.

Choosing an audio API:

- Just hear it → `useAudioPlayer`.
- Process/visualize/transcribe in JS → `useAudioChunks` (see custom-tracks.md).

## Rendition switching & latency

```ts
broadcast.videoTracks.map((t) => t.name); // e.g. ['hd', 'sd']
player.switchVideoTrack('sd');
player.updateTargetLatency(500); // higher = more buffer, fewer stalls
```

## Stats

`PlaybackStats`: `videoLatencyMs`, `audioLatencyMs`, `videoBitrateKbps`, `audioBitrateKbps`, `videoFps`, `videoJitterBufferMs`, `audioRingBufferMs`, `timeToFirstVideoFrameMs`, `timeToFirstAudioFrameMs`, `videoFramesDropped`, `audioFramesDropped`, `videoStalls` / `audioStalls` (`{ count, totalDurationMs, rebufferingRatio }`). All optional.

## Events without re-rendering

```ts
const state = useEvent(session, 'stateChange', { state: session.state }); // re-renders
useEventListener(player, 'statsUpdate', (stats) => log(stats)); // side effect only; listener kept in a ref
const sub = player.addListener('playingChange', cb); sub.remove(); // manual
```

## react-native-moq-ui

```tsx
<VideoPlayerView player={player} videoAspectRatio={16 / 9}>
  {/* children render as an overlay */}
</VideoPlayerView>
```

- `VideoPlayerView` props: `player` (required), `videoAspectRatio` (default 16/9 — **pass it on Android**, fullscreen `SurfaceView` stretches otherwise), `controls` / `miniControls` (`boolean | ReactNode`, default true), `onFullscreenEnter/Exit`. Ref: `enterFullscreen()`, `exitFullscreen()`.
- Custom chrome: render your own node as `controls` and use `useFullscreenControls()` → `{ player, exit, show, visible }` (fullscreen) or `useMiniPlayerControls()` → `{ player, enterFullscreen, show, visible }` (inline). Both throw outside their context.
- Extras: `<VolumeSlider player={...} width? theme?>` (works with `AudioPlayer` too), `<SpeakerGlyph size? volume? color?>`.

## Troubleshooting

- `useBroadcasts` empty → session not `connected`, or nothing published under that prefix.
- Listed but won't play / relay reset `code=13` (NotFound) → path/prefix mismatch; subscribe with the exact published sub-path or prefix `''`.
- Stretched fullscreen video on Android → missing `videoAspectRatio`.
