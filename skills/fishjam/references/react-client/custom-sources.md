# Custom Sources & Middleware

Two related features:

1. **`useCustomSource`** — push a `MediaStream` you created yourself (canvas, file, generated frames) into the room as an additional track.
2. **Track middleware** — wrap an existing track (camera / mic / screen share) with a transform before it goes on the wire.

## useCustomSource

Source: `@fishjam-cloud/react-client` → `packages/react-client/src/hooks/useCustomSource.ts`.

Each call takes a string `sourceId` — that ID is the identity of one custom source. Use different IDs for different concurrent sources.

```tsx
import { useCustomSource } from '@fishjam-cloud/react-client';

const { stream, setStream } = useCustomSource('canvas-overlay');

// To publish:
setStream(myCanvasStream);

// To stop:
setStream(null);
```

`setStream` is `(newStream: MediaStream | null) => Promise<void>` — fire-and-forget is fine for most cases, but you can `await` it if you need to know the previous source's tracks have been removed before continuing.

`stream: MediaStream | undefined` is the stream you previously set, for read-back.

### Generating a stream from a canvas

```tsx
const canvasRef = useRef<HTMLCanvasElement>(null);
const { setStream } = useCustomSource('canvas-art');

useEffect(() => {
  const stream = canvasRef.current?.captureStream(30); // 30fps
  if (stream) setStream(stream);
  return () => setStream(null);
}, []);
```

### From a `<video>` element (file playback)

```tsx
const videoRef = useRef<HTMLVideoElement>(null);
const { setStream } = useCustomSource('intro-clip');

useEffect(() => {
  const v = videoRef.current;
  if (!v) return;
  v.src = '/intro.mp4';
  v.muted = true;
  v.play();
  const stream = (v as any).captureStream?.() ?? (v as any).mozCaptureStream?.();
  setStream(stream);
  return () => setStream(null);
}, []);
```

Custom-source tracks appear in `usePeers()` as `customVideoTracks` / `customAudioTracks` on the local peer (and remote peers, if others use the same hook). Their `metadata.type` is `'customVideo'` or `'customAudio'`.

## Track middleware

A middleware wraps a single `MediaStreamTrack` and returns a transformed one. Used for filters, blur, watermark, color correction.

### Types

```ts
type TrackMiddleware = ((
  track: MediaStreamTrack,
) => MiddlewareResult | Promise<MiddlewareResult>) | null; // { track: MediaStreamTrack, onClear?: () => void }

type TracksMiddleware = (
  videoTrack: MediaStreamTrack,           // always present for screen share
  audioTrack: MediaStreamTrack | null,    // null when the user didn't share system audio
) => TracksMiddlewareResult | Promise<TracksMiddlewareResult>; // { videoTrack, audioTrack, onClear }
```

Both signatures may return a `Promise` — middleware can be async (e.g. waiting for a WebGL context or a worker handshake).

Use `TrackMiddleware` for camera/mic (one track). Use `TracksMiddleware` for screen share (which has a video + optional audio together).

### Camera blur example

```tsx
import { useCamera, type TrackMiddleware } from '@fishjam-cloud/react-client';

function applyBlur(track: MediaStreamTrack) {
  // Use a worker or MediaPipe — here as a sketch:
  const blurred = createBlurredTrack(track);   // your impl
  return {
    track: blurred,
    onClear: () => blurred.stop(),
  };
}

const { setCameraTrackMiddleware } = useCamera();

// Turn blur on:
setCameraTrackMiddleware(applyBlur);

// Turn blur off:
setCameraTrackMiddleware(null);
```

The provider applies the middleware before the track goes on the wire — remote peers see the blurred stream. `onClear` is called when the middleware is removed or the track changes.

### Screen-share middleware

```tsx
const { setTracksMiddleware } = useScreenShare();

setTracksMiddleware((video, audio) => {
  // `video` is always a MediaStreamTrack; `audio` may be null.
  return { videoTrack: video, audioTrack: audio, onClear: () => {} };
});

// Remove:
setTracksMiddleware(null);
```

## Combining: custom source + middleware

These are orthogonal:

- Use `useCustomSource` when the **track is yours** — you generate the frames.
- Use middleware when you want to **transform an existing track** (the SDK's camera/mic/screen share).

For a blur-camera effect, prefer middleware (cleaner; the SDK still owns the track lifecycle). For a "show this animated character instead of my camera", use `useCustomSource` with a canvas.

## Sources

- `@fishjam-cloud/react-client` → `packages/react-client/src/hooks/useCustomSource.ts`
- `packages/react-client/src/hooks/useScreenShare.ts` (`setTracksMiddleware`)
- `packages/react-client/src/types/public.ts` (`TrackMiddleware`, `TracksMiddleware`)
- <https://fishjam.swmansion.com/docs/how-to/client/custom-sources>
- <https://fishjam.swmansion.com/docs/how-to/client/stream-middleware>
- Background-blur example: `examples/mobile-react-native/blur-example/` in the SDK monorepo (RN-specific but same idea on web)
