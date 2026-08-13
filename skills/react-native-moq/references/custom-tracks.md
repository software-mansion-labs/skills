# Data tracks & custom media sources

Data tracks are **not listed in the broadcast catalog** (`BroadcastInfo` only has video/audio tracks) — publisher and subscriber agree on the track name out of band. Default name: `'data'`.

## Publish data: `useDataTrack`

```ts
const dataTrack = useDataTrack({ name: 'subtitles' });
publisher.publish({ path, tracks: [camera, dataTrack] });
dataTrack.send(JSON.stringify({ text: 'hello' })); // UTF-8 string; no-op until publishing
```

Binary payloads: encode to base64 yourself.

## Consume data: `useDataMessages`

```ts
const sub = useDataMessages(broadcast, (m) => setCaption(m.payload), {
  trackName: 'subtitles', // default 'data'
  autoStart: true,
});
```

- `DataMessage`: `{ payload: string, trackName, groupSequence, objectIndex }`.
- Returns `ChunkSubscription`: `{ isActive, start(), stop() }`. `stop()` closes the track and stops pulling it over the network — stop when idle.
- The callback is kept in a ref: changing it does not re-subscribe.

## Consume audio in JS: `useAudioChunks`

For processing, visualization, or transcription (vs `useAudioPlayer` which just plays).

```ts
useAudioChunks(broadcast, (chunk) => {
  const pcm = new Float32Array(chunk.data);
  meter(pcm); // or feed a playback/ASR pipeline
}, { format: 'pcm-f32' });
```

- `format`: `'encoded'` (default; raw Opus/AAC packets) | `'pcm-f32'` | `'pcm-i16'` (decoded, interleaved).
- `AudioChunk`: `{ data: ArrayBuffer, format, trackName, codec, sampleRate, channelCount? }` + PCM-only `frameCount`, `timestampUs` + encoded-only `groupSequence`, `objectIndex`.
- Changing `format` re-subscribes; changing the callback does not. Returns the same `ChunkSubscription` as above.
- Downstream consumers often need resampling (e.g. 48 kHz → 16 kHz mono for Whisper).

## Publish your own audio: `useAudioSource`

Push PCM you generate (TTS, synthesized audio, decoded files).

```ts
const tts = useAudioSource({ name: 'audio', audioCodec: 'opus', sampleRate: 48000, channels: 1 });
publisher.publish({ path, tracks: [tts] });
tts.send(float32Pcm); // ArrayBuffer | Int16Array | Float32Array (floats in [-1, 1], interleaved)
```

- Send ahead of realtime freely — native paces output and fills gaps with silence. Slicing an utterance into ~120 ms chunks works well.
- `sampleRate`/`channels` are fixed for the source's lifetime — remount with a `key` to change; resample if your generator differs (e.g. 24 kHz TTS → 48 kHz for Opus).

## Publish your own video: `useVideoSource`

Frame-buffer pool you fill and push (game screens, canvases, filters).

```ts
const video = useVideoSource({ width: 640, height: 480, framerate: 24, poolSize: 3 });
publisher.publish({ path, tracks: [video] });

const bufferIndex = frame % video.buffers.length;
video.fillTestPattern(bufferIndex, frame); // demo filler
video.pushFrame({ bufferIndex }); // no timestampNs → stamped at push time
```

- `options` is required (`width`/`height`). Fixed for the lifetime — remount with a `key` to change.
- `buffers: CustomVideoBufferDescriptor[]` — `{ index, surfaceHandle, width, height }`.
- **iOS:** zero-copy — `surfaceHandle` is an `IOSurfaceRef` (decimal string) you can write into natively; `timestampNs` and `fence` (`MTLSharedEvent`) supported. Timestamps must be host-clock real time.
- **Android:** no JS-importable surface handle yet (`surfaceHandle: '0'`) — only natively-filled frames (like `fillTestPattern`) can be pushed; `timestampNs`/`fence` ignored.

## Realtime command channel pattern

Data tracks are fast enough for input streaming (e.g. game controls):

```ts
dataTrack.send(JSON.stringify({ type: 'buttons', buttons: held }));
```

Subscribe on the other side with `useDataMessages(broadcast, onMessage, { trackName: 'command' })` and parse the JSON.
