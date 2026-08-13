# Playback — iOS (Swift)

## Session

```swift
let session = Session(url: relayURL)            // actor — all calls await
Task { for await state in session.state { … } } // .idle → .connecting → .connected; .error / .closed
try await session.connect()
// later: await session.close()                 // idempotent; session is dead afterwards
```

- `connect()` is one-shot: reconnecting means a **new** `Session` (throws `SessionError.alreadyConnected`/`alreadyClosed` on reuse). Guard stale async completions across reconnects with a connection token (compare a `UUID` captured before each await).
- Transport failures surface as state (`.error(SessionError)`), not as exceptions from playback calls.

## Broadcast discovery

```swift
let subscription = try await session.subscribe(prefix: "live") // "" = all; throws .alreadySubscribed per duplicate prefix
for await broadcast in subscription.broadcasts { … }     // keep subscription alive; cancel() frees the prefix
```

Requires a connected session; one active subscription per exact prefix. Hold the `Task` driving the `for await` loop — cancelling it stops discovery.

## Catalog track selection

```swift
for await catalog in broadcast.catalogs() {
    let video = catalog.playableVideoTracks
        .max { ($0.config.coded?.height ?? 0) < ($1.config.coded?.height ?? 0) }?.name
    let audio = catalog.playableAudioTracks.first?.name
}
```

Catalog semantics (updates replace, stream end = offline, playable lists) are in SKILL.md. Video-track playability on iOS is based on codec families the renderer recognizes; actual decode support is still decided by AVFoundation at runtime.

## Player

```swift
// @MainActor — construct and drive on the main actor
let player = try Player(catalog: catalog, videoTrackName: video, audioTrackName: audio,
                        targetBuffering: .milliseconds(100), volume: 1.0)
try await player.play()
```

- At least one of `videoTrackName`/`audioTrackName` is required; unknown names throw at init.
- Controls: `pause()` (resumes from live on next `play()`), `stopAll(reason:)` (terminal), `switchTrack(to:)` (video rendition; seamless when both are active), `switchAudioTrack(to:)`, `updateTargetLatency(_:)` (live), `setVolume` / `audioVolume` (clamped 0–1).
- Audio plays via `AVAudioEngine` through the system output. The SDK never configures `AVAudioSession` (its docs claim none is needed) — set a `.playback`-capable category yourself or the mute switch can silence playback.

## Rendering video

`player.videoLayer` is an `AVSampleBufferDisplayLayer`; add it to a layer hierarchy and size it yourself. Nothing renders until it's in the tree and `play()` has run.

```swift
final class VideoContainerView: UIView {
    private var displayLayer: AVSampleBufferDisplayLayer?
    func setDisplayLayer(_ layer: AVSampleBufferDisplayLayer?) {
        guard layer !== displayLayer else { return }
        displayLayer?.removeFromSuperlayer()
        displayLayer = layer
        if let layer { self.layer.addSublayer(layer) }
        setNeedsLayout()
    }
    override func layoutSubviews() { super.layoutSubviews(); displayLayer?.frame = bounds }
}

struct VideoLayerView: UIViewRepresentable {
    let videoLayer: AVSampleBufferDisplayLayer?
    func makeUIView(context: Context) -> VideoContainerView { VideoContainerView() }
    func updateUIView(_ view: VideoContainerView, context: Context) { view.setDisplayLayer(videoLayer) }
    static func dismantleUIView(_ view: VideoContainerView, coordinator: ()) { view.setDisplayLayer(nil) }
}
```

Moving playback to another screen (e.g. fullscreen) reuses the **same** layer — force the representable to re-add it (the demo bumps an `.id(generation)`).

## Events, stats, diagnostics

Channel semantics are in SKILL.md; the iOS wiring:

```swift
let events = player.subscribeEvents { event in … }   // BEFORE play(); RETAIN the subscription
let stats  = player.subscribeStats { stats in … }    // first push deferred until real data
Task { for await event in player.diagnostics() { … } } // bounded (256) per call, non-replayed
```

Both `subscribeEvents` and `subscribeStats` return a `PlayerEventSubscription` that cancels itself on deinit — store it or events silently stop. Listeners are `@MainActor`. A synchronous `player.stats` snapshot property also exists.

## Lifecycle checklist

1. Create `Session` → observe state → `connect()`.
2. `subscribe(prefix)` → per broadcast, observe `catalogs()`; catalog stream ending = broadcast offline.
3. On the main actor: create `Player` → subscribe events → attach `videoLayer` → `play()`.
4. Teardown in order: cancel observation tasks → `await player.stopAll()` → cancel the broadcast subscription → `await session.close()`.

Keep strong references throughout: `Session`, `BroadcastSubscription`, `Player`, and every `PlayerEventSubscription`; the `Task` handles driving `for await` loops are the lifetime of those streams.
