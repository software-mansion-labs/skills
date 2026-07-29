# Publishing — iOS (Swift)

## Go-live order

Shared ordering and publisher states are in SKILL.md. iOS specifics: `Publisher()` and `session.publish` throw; `publisher.start()` is `async throws`; `publisher.stop()` flushes encoders **synchronously** — call it off the main thread (`Task.detached { publisher.stop() }`; `Publisher` isn't `Sendable`, so under strict Swift 6 concurrency this needs `nonisolated(unsafe)` or a Swift 5-mode target).

```swift
try await camera.start(); try await microphone.start()   // 1. captures first
let publisher = try Publisher()                          // 2. add ALL tracks before start
publisher.addVideoTrack(name: "camera", source: camera, config: videoConfig)
publisher.addAudioTrack(name: "mic", source: microphone, config: audioConfig)
try await session.publish(path: "live/ios", publisher: publisher) // 3. register
try await publisher.start()                              // 4. go live
```

Observe `publisher.state` / `publisher.events` (`AsyncStream`s) and each `PublishedTrack.state` in `Task`s you keep handles to.

## Camera and preview

```swift
let camera = CameraCapture(camera: Camera(position: .front, width: 720, height: 1280, orientation: .portrait))
try await camera.start()                       // NSCameraUsageDescription required
try camera.switch(to: Camera(position: .back)) // flip
camera.stop()
```

`CameraCapture` exposes `captureSession: AVCaptureSession` — attach an `AVCaptureVideoPreviewLayer` (e.g. a `UIView` whose `layerClass` is `AVCaptureVideoPreviewLayer`) and reuse the same `CameraCapture` for preview and publishing.

## Multi-camera

```swift
guard MultiCameraCapture.isSupported else { … }   // fails on many devices — always check
let cameras = MultiCameraCapture(front: Camera(position: .front, width: 720, height: 1280),
                                 back: Camera(position: .back, width: 720, height: 1280),
                                 maxFrameRate: 30)
try await cameras.start()
publisher.addVideoTrack(name: "front-camera", source: cameras.frontSource, config: videoConfig)
publisher.addVideoTrack(name: "back-camera", source: cameras.backSource, config: videoConfig)
```

Front and back are two separate sources → two video tracks. Multi-cam preview needs `AVCaptureVideoPreviewLayer(sessionWithNoConnection:)` with manually wired connections.

## Microphone and AVAudioSession

```swift
// Configure AVAudioSession BEFORE start — MicrophoneCapture sets
// usesApplicationAudioSession = true and will NOT configure it for you.
let audioSession = AVAudioSession.sharedInstance()
try audioSession.setCategory(.playAndRecord, mode: .videoRecording,
                             options: [.defaultToSpeaker, .allowBluetoothHFP]) // .allowBluetooth on pre-iOS-26 SDKs
try audioSession.setActive(true)

let microphone = MicrophoneCapture() // NSMicrophoneUsageDescription required
try await microphone.start()
```

Wrong/missing category is the classic "no mic audio" bug; switch back to `.playback` when returning to watch-only mode. `MicrophoneCapture` has no sample-rate knob — the encoder resamples to `AudioEncoderConfig.sampleRate` internally.

## Encoder configs and codec gating

```swift
VideoEncoderConfig(codec: .h264, width: 1920, height: 1080, bitrate: 1_500_000,
                   keyframeInterval: 2.0, maxFrameRate: 30.0)   // defaults shown; also optional profileLevel, naluFormat
AudioEncoderConfig(codec: .opus, sampleRate: 48000, channels: 1, bitrate: 128_000) // iOS default codec: opus
```

Build codec pickers from `VideoEncoderConfig.supportedCodecs()` / `AudioEncoderConfig.supportedCodecs()`; check an exact config with `.isSupported` / `.unsupportedReason`. Video encode is H.264/H.265 only (no AV1); H.265 auto-selects Annex B NALU format, and concurrent H.265 encoders (multi-cam) are unreliable on many devices — default to H.264.

## Custom sources

Any `FrameSource` (class-bound, `Sendable`; a settable `var onFrame: (@Sendable (CMSampleBuffer) -> Bool)?`) works as a track source; `FrameRelay` is the ready-made bridge — call `relay.send(sampleBuffer)` with your own video or audio samples.

- **The callback's `Bool` return is the stop signal.** `false` means the consumer is gone (the track stopped) — stop producing. `stop()` on the built-in captures clears `onFrame` for the same reason.
- **Timestamps must be host-clock.** `Publisher` shares one `PublisherClock` across all its tracks and takes the epoch from the first frame of *any* track, so PTS from a private zero-based timeline drifts against a live mic or schedules video far ahead of the render clock. Stamp buffers with `CMClockGetTime(CMClockGetHostTimeClock())` — the domain `CameraCapture` and `MicrophoneCapture` already produce.

## Permissions

Info.plist: `NSCameraUsageDescription`, `NSMicrophoneUsageDescription`, `NSLocalNetworkUsageDescription` (relays on the local network). The SDK adds none of these.

## Teardown

Cancel observer tasks → `Task.detached { publisher.stop() }` → `camera.stop()` / `microphone.stop()` → restore the audio session category → `await session.close()` if no longer needed.

Captures outlive publishers: `CameraCapture.start()` / `MicrophoneCapture.start()` are idempotent (a second call while running is a no-op, not a reconfigure), and `stop()` only stops the session and clears `onFrame`. Keep one capture instance and hand it to the next `Publisher` — the `Publisher` is what can't be reused.
