# Publishing — Android (Kotlin)

## Go-live order

Shared ordering and publisher states are in SKILL.md. Android specifics: `Publisher()` doesn't throw; `start()`/`stop()` are **not** suspend; `start()` validates every track's codec first and throws `UnsupportedCodecException` before starting anything; `start()` with zero tracks flips straight to `Stopped` (as does stopping the last track).

```kotlin
camera.start(context, lifecycleOwner); microphone.start() // 1. captures first
val publisher = Publisher()                               // 2. add ALL tracks before start
publisher.addVideoTrack(name = "camera", source = camera, config = videoConfig)
publisher.addAudioTrack(name = "mic", source = microphone, config = audioConfig)
session.publish(path = "live/android", publisher = publisher) // 3. register
publisher.start()                                         // 4. go live
```

Observe `publisher.state` (`StateFlow`), `publisher.events` (`SharedFlow`), and each `PublishedTrack.state` in coroutines on your scope. **Also observe `session.state`** — none of the publisher channels react to the session dying: when the transport drops, `publisher.state` stays `Publishing` and no event fires, while `session.state` goes `Error(message)` → `Closed`. A publish screen that watches only the publisher will keep claiming it is live long after the broadcast is gone.

## Camera and preview

```kotlin
val camera = CameraCapture(position = CameraPosition.Back, width = 1920, height = 1080, frameRate = 30)
camera.start(context, lifecycleOwner) // suspend; CameraX binds to this lifecycle; CAMERA permission required
camera.switchCamera()                 // suspend
camera.stop()
```

Preview: `camera.setPreviewSurface(surface)` with a `SurfaceView`'s surface — same `SurfaceHolder.Callback` pattern as the player (null it on `surfaceDestroyed`).

## Multi-camera

```kotlin
if (!MultiCameraCapture.isFrontBackSupported(context)) return // suspend; real CameraX pair check
val cameras = MultiCameraCapture(
    front = CameraStreamConfig(CameraPosition.Front, 1280, 720, 30),
    back = CameraStreamConfig(CameraPosition.Back, 1280, 720, 30),
)
cameras.start(context, lifecycleOwner)
publisher.addVideoTrack("front-camera", cameras.frontSource, videoConfig)
publisher.addVideoTrack("back-camera", cameras.backSource, videoConfig)
```

Two checks: `isSupported(context)` (cheap `PackageManager` feature check — gate UI on it) and `isFrontBackSupported(context)` (suspend; verifies an actual concurrent front/back pair — gate `start()` on it). They're enforced, not advisory: `start()` throws `IllegalStateException` without a concurrent pair, and the constructor `require`s that `front`/`back` positions match their labels. Front and back are two separate sources → two video tracks. The demo previews multi-cam with two `TextureView`s.

## Microphone

```kotlin
val microphone = MicrophoneCapture(sampleRate = 48_000, channels = 1)
microphone.start() // NOT suspend; silently no-ops if RECORD_AUDIO is missing or AudioRecord fails
microphone.stop()
```

The silent no-op means a missing permission produces a broadcast with no audio and no error — request `RECORD_AUDIO` before starting. `start()` is annotated `@RequiresPermission(RECORD_AUDIO)`, so call sites need the permission proven to lint (a checked request, `@RequiresPermission` on the caller, or an explicit suppression). Keep `sampleRate` equal to `AudioEncoderConfig.sampleRate` (Opus wants 48 kHz).

Unlike iOS, `start()` is **not** idempotent: calling it twice without `stop()` leaves the first `AudioRecord` and its reader thread running while only the second is tracked, so `stop()` releases one of them and duplicate PCM keeps flowing. Gate re-publish on your own running flag. Stopping first, then starting again, is fine — a stopped capture can be handed to the next `Publisher` (only the `Publisher` is single-use).

## Encoder configs and codec gating

```kotlin
VideoEncoderConfig(codec = VideoCodec.H264, width = 1920, height = 1080, bitrate = 1_500_000,
                   keyframeIntervalSeconds = 2, frameRate = 30)  // defaults shown; also optional profile
AudioEncoderConfig(codec = AudioCodec.AAC, sampleRate = 48_000, channels = 1, bitrate = 128_000) // Android default codec: aac
```

Build codec pickers from `VideoEncoderConfig.supportedCodecs()` / `AudioEncoderConfig.supportedCodecs()`; check an exact config with `.isSupported` / `.unsupportedReason`. Encoder support varies per device — H.264 is the safe default; gate H.265 explicitly.

## Custom sources

Implement `VideoFrameSource` (encoder-`Surface`-based: `attachEncoderSurface`/`detachEncoderSurface`/`setPreviewSurface`) or `AudioFrameSource` (`onPcmData: (ByteArray, Int, Long) -> Unit`, 16-bit PCM). Advanced — the built-in captures cover camera/mic/screen.

**Timestamps must share one clock domain.** `Publisher` stamps all of its tracks against a single epoch taken from the first frame of *any* track, so a source with a private zero-based timeline drifts against a live mic or leaves video scheduled far ahead. Use `SystemClock.elapsedRealtimeNanos() / 1_000` for `onPcmData`'s `timestampUs` (exactly what `MicrophoneCapture` does) and drive the encoder surface's presentation timestamps off the same clock.

`Publisher` sets `onPcmData` when the audio track starts and clears it when the track stops — don't set it yourself on a source you hand to a publisher.

## Permissions

Manifest: `INTERNET`, `CAMERA`, `RECORD_AUDIO` — declare **and** request at runtime (`RequestMultiplePermissions`) before starting captures. The library does not add them transitively. Screen capture additionally needs the MediaProjection foreground-service setup (see the screen-capture reference).

## Teardown

Cancel observer jobs → `publisher.stop()` → `camera.stop()` / `microphone.stop()` → `session.close()` if no longer needed. `session.unpublish(path)` is equivalent to `publisher.stop()` for that path.
