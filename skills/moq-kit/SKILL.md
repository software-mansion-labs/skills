---
name: moq-kit
description: "Software Mansion's moq-kit — native Swift (iOS) and Kotlin (Android) SDKs for Media over QUIC (moq-lite) live streaming: sub-second-latency playback, camera/microphone/screen publishing, and realtime data tracks over a MoQ relay. MUST USE before writing, reviewing, or debugging ANY code that imports MoQKit (Swift) or com.swmansion.moqkit (Kotlin); for React Native apps use the react-native-moq skill instead. Trigger on: 'moq-kit', 'MoQKit', 'Media over QUIC', 'moq-lite', 'MoQ relay', 'moq-ffi', 'BroadcastSubscription', 'DataTrackEmitter', 'CameraCapture', 'MicrophoneCapture', 'MultiCameraCapture', 'ScreenCapture', 'TrackSubscription', 'AudioDataStream', 'MoQReplayKitBroadcastSampleHandler', or a Session/Player/Publisher/Catalog from moq-kit."
license: Apache-2.0
---

# moq-kit

Native Swift and Kotlin SDKs for Media over QUIC live streaming: connect to a relay, discover broadcasts, play catalog-described streams with sub-second latency, publish camera/microphone/screen tracks, and exchange raw data tracks. Wraps the UniFFI bindings generated from `moq-ffi` (Luke Curley's `moq-dev/moq`); targets the `moq-lite` protocol, not the IETF `moq-transport` draft.

**Requirements:** iOS 16+ / macOS 13+ (Xcode 16+ to build) / Android minSdk 29, compileSdk 35. A moq-lite relay to connect to. The references cover iOS and Android; macOS is declared by the Swift package but untested here.

moq-kit is an **active preview**: public APIs, packaging, and codec coverage can still change between releases. Pin an exact version and check signatures against the installed artifact rather than assuming a shape.

**Install — iOS (SPM):**

```swift
.package(url: "https://github.com/software-mansion-labs/moq-kit", exact: "<latest release>")
// product: .product(name: "MoQKit", package: "moq-kit")
```

**Install — Android (Maven Central):**

```kotlin
dependencies { implementation("com.swmansion.moqkit:moqkit:<latest release>") }
```

Fill in `<latest release>` from https://github.com/software-mansion-labs/moq-kit/releases. Both pins are exact on purpose — loosen the SPM one (`.upToNextMinor(from:)`) only if you mean to track preview releases as they land.

## Platform detection

Reference files are split per platform — read only the ones for the project at hand:

- `Package.swift`, `*.xcodeproj`, or a `Podfile` → the `*-ios.md` references (Swift).
- `build.gradle`/`build.gradle.kts` with `com.swmansion.moqkit` → the `*-android.md` references (Kotlin).
- Cross-platform work (feature parity, bindings on top of both SDKs) → read both; the APIs are deliberately kept equivalent.

## Mental model

- A **`Session`** owns one QUIC connection to a relay and serves both subscribing and publishing. `connect()` once; a session is **one-shot** — after `close()` (or a connection error) create a new `Session`. Swift: `actor`, `connect() async throws`. Kotlin: optionally takes a `parentScope` (defaults to an internal IO scope; pass `viewModelScope`/`lifecycleScope` so cancelling it tears the session down); `connect()` is the only suspend call.
- **Subscribing:** `session.subscribe(prefix:)` → stream of `Broadcast` → `broadcast.catalogs()` → pick playable tracks → `Player(catalog:videoTrackName:audioTrackName:targetBuffering:)` → `play()`. Rendering is platform-specific: iOS adds `player.videoLayer` (an `AVSampleBufferDisplayLayer`) to the view hierarchy; Android calls `player.setSurface(surface)` with a `Surface` you own.
- **Publishing:** create capture sources (`CameraCapture`, `MicrophoneCapture`, `MultiCameraCapture`, `ScreenCapture`) and **start them yourself** → `Publisher()` + `addVideoTrack`/`addAudioTrack`/`addDataTrack` → `session.publish(path, publisher)` → `publisher.start()`. A publisher is single-use; add all tracks before `start()`.
- **Data tracks** (`DataTrackEmitter` → `addDataTrack`; consume with `broadcast.subscribeTrack(name:)`) bypass the media catalog — publisher and subscriber agree on the track name out of band.
- Everything is observable: `session.state`, `publisher.state`/`events`, player events/stats. Kotlin uses `StateFlow`/`SharedFlow`/cold `Flow`; Swift uses `AsyncStream` for session/publisher state and diagnostics, but callback subscriptions (`subscribeEvents`/`subscribeStats`) for player events/stats.

## Quick start — watch

**iOS:**

```swift
let session = Session(url: "http://localhost:4443/anon")
try await session.connect()
let subscription = try await session.subscribe(prefix: "live")
for await broadcast in subscription.broadcasts {
    for await catalog in broadcast.catalogs() {
        let video = catalog.playableVideoTracks.first?.name
        let audio = catalog.playableAudioTracks.first?.name
        guard video != nil || audio != nil else { continue }
        let player = try await MainActor.run {
            try Player(catalog: catalog, videoTrackName: video, audioTrackName: audio)
        }
        try await player.play() // render via player.videoLayer — see playback-ios.md
    }
}
```

**Android:**

```kotlin
scope.launch {
    val session = Session(url = "http://localhost:4443/anon", parentScope = scope)
    session.connect()
    session.subscribe(prefix = "live").broadcasts.collect { broadcast ->
        broadcast.catalogs().collect { catalog ->
            val video = catalog.playableVideoTracks.firstOrNull()?.name
            val audio = catalog.playableAudioTracks.firstOrNull()?.name
            if (video == null && audio == null) return@collect
            val player = Player(catalog, video, audio, parentScope = scope)
            player.setSurface(surfaceView.holder.surface)
            player.play()
        }
    }
}
```

Minimal on purpose: the nested loops handle one broadcast at a time — spawn a `Task`/`launch` per broadcast in real code — and each catalog update should tear down the previous player before creating the new one (players hold native handles).

## Quick start — go live

Same four steps on both platforms, and the order is load-bearing (rule 3):

**iOS:**

```swift
try await camera.start()                                   // 1. captures — publisher.start() won't do this
let publisher = try Publisher()                            // 2. every track before start()
publisher.addVideoTrack(name: "camera", source: camera, config: VideoEncoderConfig())
try await session.publish(path: "live/ios", publisher: publisher)  // 3. register
try await publisher.start()                                // 4. go live
```

**Android:**

```kotlin
camera.start(context, lifecycleOwner)
val publisher = Publisher()
publisher.addVideoTrack(name = "camera", source = camera, config = VideoEncoderConfig())
session.publish(path = "live/android", publisher = publisher)
publisher.start()
```

Captures, previews, encoder configs, and permissions are in the publishing references.

## Shared concepts (both platforms)

- **Catalogs:** `broadcast.catalogs()` streams catalog updates — each element **replaces** the previous; the stream ending means the broadcast went away (mark it offline). It often ends by **throwing** rather than completing — a producer that disappears surfaces as `MoqException$Mux: json: remote error: code=…` out of `collect`, so handle the exception and the clean completion identically (both demo apps and react-native-moq do). Always select from `playableVideoTracks` / `playableAudioTracks` (device-filtered on Android and for iOS audio; iOS video playability is codec-family recognition only), not the raw lists; per-track `isPlayable`/`unsupportedReason` explain exclusions (they live on the concrete `VideoTrackInfo`/`AudioTrackInfo`, not the base `TrackInfo`). Video codec strings: `"avc1"`, `"hev1"`, `"av01"`; audio: `"mp4a.40.2"`/`"aac"` (Android-published tracks advertise `"aac"`), `"opus"`. AV1 is **playback-only** — neither platform publishes it, and Apple-side decode is iPhone 15 Pro-class hardware, not a general capability.
- **Broadcast paths are prefix-relative:** subscribing with prefix `"live"` reports broadcast `live/game` as `path: "game"` — subscribe with `""` to see full paths.
- **Three observation channels on `Player`** — don't mix them up: **events** (lifecycle milestones — `player.init`, `playback.start`, `track.ready`, `track.switch`, `track.stall.start`, `rebuffer.end`, `decode.error` — wire values; Kotlin matches the `PlayerEventName` enum; not replayed), **stats** (`PlaybackStats` at ~1 s cadence — latency, fps, bitrates, `StallStats`, `timeToFirst`, buffer depths, drops, switches), and **diagnostics** (typed per-frame `PipelineEvent` — `frameDropped` with `DropStage`/`DropReason`, `stallStarted`, `latencySample`, … — bounded at 256 drop-oldest — Android hands every caller the same shared flow — and never backpressures playback; for logging/telemetry only).
- **Publisher states:** `Idle / Publishing / Stopped / Error(message)` — but Android never actually enters `Error`, and it does not react to the session dying either: an Android publisher sits in `Publishing` indefinitely after its session is gone, so the session's own state is the only signal that the broadcast stopped. Wire failure UI to **events**: `TrackStarted` / `TrackStopped` / a per-track error (Kotlin `TrackError(name, message)`, Swift `.error(name, message)`); per-track `PublishedTrack.state` (`Idle / Starting / Active / Stopped`) with `PublishedTrack.stop()` to end one track. `session.unpublish(path)` calls `publisher.stop()`.
- **`TrackDelivery`** for raw track subscription: `monotonic` (default) skips late groups — right for newest-value-wins live state; `arrival` delivers every group in order — right for chat and anything lossless.

## Critical rules

1. **Sessions, publishers, and players are single-use.** A `Session` connects once (`alreadyConnected`/`IllegalStateException` on reuse); after `close()` make a new one. A stopped `Publisher` cannot be restarted. Never restart a torn-down `Player` on either platform — create a fresh one.
2. **Media track names are local labels, catalog names are muxer-generated.** The `name:` passed to `addVideoTrack`/`addAudioTrack` (defaults `"video"`/`"audio"`) only labels `PublishedTrack` and publisher events; subscribers must discover real track names from `Catalog.videoTracks`/`audioTracks`. **Exception: data tracks** — `addDataTrack(name:)` (default `"data"`) IS the name subscribers pass to `subscribeTrack(name:)`.
3. **`publisher.start()` does not start capture.** Start `CameraCapture`/`MicrophoneCapture`/etc. yourself, in this order: start captures → add tracks → `session.publish(path, publisher)` → `publisher.start()`.
4. **Gate codecs before use.** `VideoEncoderConfig.supportedCodecs()` / `AudioEncoderConfig.supportedCodecs()` / `.isSupported` on the publish side; `catalog.playableVideoTracks` / `playableAudioTracks` / per-track `isPlayable` on the playback side. Unsupported configs throw (`SessionError.unsupportedCodec` / `UnsupportedCodecException`). H.264 + AAC/Opus are the best-tested paths; default to H.264. Opus wants 48 kHz — on Android keep `MicrophoneCapture(sampleRate = …)` equal to the encoder config; iOS has no mic rate knob and resamples internally.
5. **Player startup events are not replayed.** Subscribe before `play()`: iOS `player.subscribeEvents { ... }` (and **retain** the returned subscription — it cancels on deinit); Android collect `player.events` with `CoroutineStart.UNDISPATCHED`. On both platforms `player.init` fires during construction and can never be observed — the rule protects `track.ready`/`playback.start`.
6. **Permissions and audio sessions are the app's job.** iOS: `NSCameraUsageDescription`, `NSMicrophoneUsageDescription`, `NSLocalNetworkUsageDescription` (local relays), and configure `AVAudioSession` yourself (`.playAndRecord` to publish, `.playback` to watch) — `MicrophoneCapture` will not. Android: declare + request `CAMERA`/`RECORD_AUDIO` (`MicrophoneCapture.start()` silently no-ops without it), `INTERNET`, and the MediaProjection foreground-service setup for screen capture.
7. **Own your references and clean up.** iOS: keep strong refs to `Session`, `BroadcastSubscription`, `Publisher`, `DataTrackEmitter`, `Player`, and every event/stats subscription; hold the `Task` handles that drive `for await` loops (the `AsyncStream`s are single-consumer, too). Android: nearly everything is `AutoCloseable` and holds ref-counted native handles — `close()` each `Broadcast`, `BroadcastSubscription`, `TrackSubscription`, `MediaTrack`, `AudioDataStream`, and `Player` (leaking one keeps the broadcast alive). `broadcasts`, `objects`, `frames`, and `audio` flows allow a **single collector**.
8. **iOS `Player` is `@MainActor`** — construct and drive it on the main actor. **iOS `Publisher.stop()` flushes encoders synchronously** — call it off the main thread (`Task.detached`; note `Publisher` isn't `Sendable`, so under strict Swift 6 concurrency this needs `nonisolated(unsafe)` or a Swift 5-mode target).
9. **When one session both publishes and subscribes the same prefix, it sees its own broadcast** — filter it out, remembering reported paths are prefix-relative: compare against your publish path with the subscribe prefix stripped. Two `Session`s in one process (one publishing, one subscribing) is also supported — and since sessions are one-shot, separate ones keep a transport failure or teardown on the publish side from taking playback down with it.
10. **Latency is tunable live**: `targetBuffering` defaults to 100 ms; `player.updateTargetLatency(_:)` adjusts during playback. A second consumer of the same media track (another `Player`, an `AudioDataStream`) reuses the first subscriber's upstream buffering — it can't independently lower it.
11. **Custom sources must share the publisher's clock domain.** A `Publisher` stamps all of its tracks against one epoch, set by the first frame of *any* track, so a source emitting its own zero-based timeline drifts against the mic or leaves video scheduled far in the future. iOS: PTS on the host clock (`CMClockGetTime(CMClockGetHostTimeClock())`) — the domain AVFoundation capture already uses. Android: `SystemClock.elapsedRealtimeNanos() / 1_000` for PCM (what `MicrophoneCapture` uses) and the same clock behind the encoder surface's presentation timestamps.

## Symptom → cause

| Symptom | Likely cause |
| --- | --- |
| No broadcasts ever arrive | Relay URL missing the namespace path (`/anon` locally) — the relay answers `NotFound` (`code=13`) and nothing surfaces; or `subscribe()` ran before the session reached connected |
| Publishing "works" but subscribers get no media | Captures were never started — `publisher.start()` does not start them (rule 3) |
| No microphone audio | iOS: `AVAudioSession` isn't on a record-capable category. Android: `RECORD_AUDIO` is missing and `MicrophoneCapture.start()` no-ops silently — no error either way |
| Android: audio plays, video is black | No surface set, or a stale one after `surfaceDestroyed` — null it there and only re-apply a surface that `isValid` |
| iOS: audio plays, video is black | `player.videoLayer` was never added to a layer hierarchy or never sized — nothing renders until it's in the tree and `play()` has run |
| `track.ready` / `playback.start` never fire | Subscribed after `play()`, or (iOS) the returned `PlayerEventSubscription` wasn't retained |
| Data-track sends disappear | Sent before `publisher.start()` or after `stop()` — dropped silently on both platforms; gate sends on publisher state |
| Nothing recovers after a network blip | Sessions are one-shot; on Android a `Closed` you didn't ask for **is** the failure signal. Build a new `Session` — but log `Error(message)` if you see it, it carries the only description of what went wrong |
| A remote error carries a bare `code=N` | Codes are `moq_net::Error` discriminants: `13` NotFound, `24` Dropped "producer dropped without finishing", `25` Closed, `26` Lagged; full list in moq-net's `error.rs` (`vendor/moq/rs/moq-net/src/error.rs` in a moq-kit checkout) |
| iOS: UI hitches when going off-air | `Publisher.stop()` flushes encoders synchronously — call it off the main thread |
| A publisher's own broadcast shows up in its subscriber list | Same session publishing and subscribing the same prefix (rule 9) — filter by prefix-relative path |
| Custom source: A/V drifts apart, or video freezes while audio keeps playing | Frame timestamps aren't in the same clock domain as the publisher's other tracks (rule 11) |
| Android: `play()` throws `IllegalArgumentException` | It's an `UnsupportedCodecException` (a subclass) for an undecodable selected track — catch that first |
| Android: screen capture starts, then dies | The `mediaProjection` foreground service wasn't running yet, or was started outside the consent-result callback (Android 14+ rejects that) |
| Device or emulator can't reach a local relay | `localhost` is the device itself — use `10.0.2.2` from the Android emulator, the host's LAN IP from physical devices |

## Debugging

Relay, subscription, and publish failures usually surface in the native transport log before they show up as state. Raise the level before touching any other API:

```swift
KitLogger.setNativeLogLevel("debug")   // "error"/"warn"/"info"/"debug"/"trace"; only the first call takes effect
```

```kotlin
NativeLogging.setLogLevel("debug")     // same levels; invalid values are ignored with a warning
```

Android native logs land in **logcat under the tag `MoQNative`**, but only when the native library was built with the `android-logcat` cargo feature (`moq-ffi` → `moq-native` → `tracing-android`); without it the `tracing` subscriber writes to stderr, which Android discards. A successful `setLogLevel` call with no `MoQNative` lines afterwards means the feature is off in the artifact you are using — fall back to session/publisher state and player events there.

iOS also logs through `os_log` under subsystem `com.swmansion.MoQKit` (categories `session`, `transport`, `catalog`, `media`, `player`, `publish`) — filter Console.app on it. Per-frame drop/stall detail is on the player's diagnostics channel, not in these logs.

## Testing against a relay

From a moq-kit checkout (needs Rust + [mise](https://mise.jdx.dev)): `mise run relay:run` starts a local moq-lite relay at `http://localhost:4443/anon`; `mise run media:to-fmp4` converts a video to CMAF fMP4 and `mise run stream:file --input file.mp4` loops it into the relay as a broadcast. Use the machine's LAN IP instead of `localhost` from physical devices, and `10.0.2.2` from the Android emulator (the emulator's fixed alias for the host's loopback — its own `localhost` is itself; the iOS simulator shares the host's). The relay URL must include the namespace path (`/anon` locally) — with it wrong or missing, broadcasts simply never appear (relay `NotFound`, `code=13`). The demo apps (`examples/ios/demo/MoQDemo`, `examples/android/demo/MoQDemo`) are the canonical integration references.

## References

| Task | iOS (Swift) | Android (Kotlin) |
| --- | --- | --- |
| Playback: session, discovery, Player, rendering, events/stats wiring | [references/playback-ios.md](references/playback-ios.md) | [references/playback-android.md](references/playback-android.md) |
| Publishing: camera, multi-camera, microphone, previews, encoder configs, permissions | [references/publishing-ios.md](references/publishing-ios.md) | [references/publishing-android.md](references/publishing-android.md) |
| Screen sharing | [references/screen-capture-ios.md](references/screen-capture-ios.md) | [references/screen-capture-android.md](references/screen-capture-android.md) |
| Data tracks, raw compressed media (MediaTrack), decoded PCM (AudioDataStream) | [references/data-tracks-ios.md](references/data-tracks-ios.md) | [references/data-tracks-android.md](references/data-tracks-android.md) |
