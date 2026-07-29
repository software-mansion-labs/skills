# Using without React (imperative API)

Every hook has a `create*` / `subscribe*` counterpart for non-React code (services, state stores, other frameworks). `create*` returns the hook's object shape **plus `destroy()`**; `subscribe*` returns a subscription with `start()`/`stop()` (plus an `autoStart` option, default true). Nothing is cleaned up automatically — pair every `create*` with `destroy()` and every `subscribe*` with `stop()`.

| Hook | Imperative |
| --- | --- |
| `useSession(url)` | `createSession(url): SessionHandle` |
| `useBroadcasts(session, prefix)` | `subscribeBroadcasts(session, prefix, onChange, options?): BroadcastSubscription` |
| `useVideoPlayer(broadcast)` | `createVideoPlayer(broadcast): VideoPlayerHandle` |
| `useAudioPlayer(broadcast)` | `createAudioPlayer(broadcast): AudioPlayerHandle` |
| `useAudioChunks(...)` | `subscribeAudioChunks(broadcast, trackName, onChunk, options?)` |
| `useDataMessages(...)` | `subscribeDataMessages(broadcast, trackName, onMessage, options?)` |
| `usePublisher(session)` | `createPublisher(session): PublisherHandle` |
| `useCamera` / `useMultiCamera` / `useMicrophone` | `createCamera` / `createMultiCamera` / `createMicrophone` |
| `useDataTrack` / `useAudioSource` / `useVideoSource` | `createDataTrack` / `createAudioSource` / `createVideoSource` |
| `useScreenBroadcast(session, options)` | `createScreenBroadcast(session, options)` |

Differences from hooks:

- Capture factories take the hook options **minus `enabled`** — creation always starts capture; toggle by `destroy()` + recreate.
- `createVideoSource` additionally exposes `ready: Promise<CustomVideoBufferDescriptor[]>` — await it before pushing frames.
- State changes surface via `addListener` — but the event differs per handle: capture handles (`createCamera`/`createMicrophone`/`createMultiCamera`/`createScreenBroadcast`) emit `stateChange` with `{ state, lastError }`; `PublisherHandle` emits `stateChange` (`{ state }`) plus `trackStateChange` (`{ name, state, error? }`); players emit `playingChange`/`trackSwitched`/`trackStopped`/`statsUpdate`. The `createDataTrack`/`createAudioSource`/`createVideoSource` handles have **no listeners** — just their send/push methods plus `destroy()`.
- `VideoPlayerHandle.destroy()` only detaches listeners (the native player is shared per broadcast); `AudioPlayerHandle.destroy()` also releases its dedicated native player.
- Views still need React, but accept imperative handles: `<VideoView player={videoPlayerHandle}>`, `<PublisherView camera={cameraHandle}>`.

## Watch

```ts
const session = createSession(url);
session.connect();

const sub = subscribeBroadcasts(session, '', (broadcasts) => {
  if (broadcasts.length && !player) {
    player = createVideoPlayer(broadcasts[0]);
    player.play();
  }
});

// teardown
player?.destroy();
sub.stop();
session.destroy();
```

## Publish custom media

```ts
const session = createSession(url);
session.connect();

const audio = createAudioSource({ sampleRate: 48000, channels: 1 });
const video = createVideoSource({ width: 640, height: 480, framerate: 30 });
const publisher = createPublisher(session);

await video.ready;
session.addListener('stateChange', ({ state }) => {
  if (state === 'connected') publisher.publish({ path: 'live/custom', tracks: [audio, video] });
});

// teardown (reverse order)
publisher.destroy(); // also stops publishing
audio.destroy();
video.destroy();
session.destroy();
```
