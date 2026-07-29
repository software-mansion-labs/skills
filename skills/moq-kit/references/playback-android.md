# Playback — Android (Kotlin)

## Session

```kotlin
val session = Session(url = relayURL, parentScope = viewModelScope)
scope.launch { session.state.collect { … } }  // StateFlow: Idle/Connecting/Connected/Error(message)/Closed
session.connect()                             // the only suspend call
// later: session.close()                     // synchronous, idempotent
```

- `connect()` is one-shot: reconnecting means a **new** `Session`. Guard stale async completions across reconnects with a connection token.
- `parentScope` puts all session work under a child `SupervisorJob` of that scope — cancelling `viewModelScope` tears the session down. On connection loss the session sets `Error` and immediately closes itself, milliseconds apart; treat a `Closed` you didn't request as the failure signal, since the conflated `StateFlow` can skip `Error`. Handle `Error(message)` too rather than ignoring it — a plain `collect` does observe it in practice, and the message (e.g. `Session ended: uniffi.moq.MoqException$Protocol: transport: connection closed`) is the only description of the failure you get anywhere.
- Transport failures surface through session state, never as exceptions from playback calls.

## Broadcast discovery

```kotlin
val subscription = session.subscribe(prefix = "live") // not suspend; requires Connected; "" = all
subscription.broadcasts.collect { broadcast -> … }    // cold, SINGLE collector
```

- One active subscription per exact prefix. `subscribe()`, `publish()`, `play()` are **not** suspend — easy to over-wrap.
- `BroadcastSubscription` and each emitted `Broadcast` are `AutoCloseable` and hold ref-counted native handles — `close()` them (the demo does it in `finally`).

## Catalog track selection

```kotlin
broadcast.catalogs().collect { catalog ->
    val video = catalog.playableVideoTracks.maxByOrNull { it.config.coded?.height ?: 0u }?.name
    val audio = catalog.playableAudioTracks.firstOrNull()?.name
}
```

Catalog semantics (updates replace, stream end = offline, playable lists) are in SKILL.md. Wrap the collect in `try`/`catch`: the stream usually ends by throwing `MoqException$Mux` when the producer goes away, so the offline path needs to run from both `catch` and normal completion — see `observeCatalogs` in the demo's `PlayerDemoViewModel`.

## Player

```kotlin
val player = Player(catalog, video, audio,
                    targetBuffering = Duration.ofMillis(100),
                    parentScope = viewModelScope, volume = 1f) // parentScope has NO default
player.play() // not suspend; fine to call before a surface exists
```

- At least one of `videoTrackName`/`audioTrackName` is required; unknown names throw `IllegalArgumentException` at construction; undecodable selected tracks throw `UnsupportedCodecException` from `play()` — it extends `IllegalArgumentException`, so catch it first.
- Controls: `pause()` (resumes from live on next `play()`), `stopAll(reason = …)` (tears down) and `close()` (terminal — only `close()` blocks a later `play()`; treat both as end-of-life), `switchTrack(trackName)` (video rendition; `null` disables video), `switchAudioTrack(trackName)`, `updateTargetLatency(latency)` (live), `setVolume` (clamped 0–1).
- Decodes via MediaCodec, renders audio via AudioTrack; video needs a `Surface` (below).

## Rendering video

No view class is provided — feed `setSurface(Surface?)` from a `SurfaceView`:

```kotlin
AndroidView(factory = { ctx ->
    SurfaceView(ctx).apply {
        holder.addCallback(object : SurfaceHolder.Callback {
            override fun surfaceCreated(h: SurfaceHolder) = player.setSurface(h.surface)
            override fun surfaceChanged(h: SurfaceHolder, f: Int, w: Int, hh: Int) {}
            override fun surfaceDestroyed(h: SurfaceHolder) = player.setSurface(null)
        })
    }
})
```

Set the surface on `surfaceCreated`, **null it on `surfaceDestroyed`** (a stale surface means no video or crashes), and only re-apply a surface that `isValid`. `setSurface(null)` keeps audio running while video waits; `play()` before any surface is fine.

## Events, stats, diagnostics

Channel semantics are in SKILL.md; the Android wiring:

```kotlin
scope.launch(start = CoroutineStart.UNDISPATCHED) { player.events.collect { … } } // BEFORE play()
scope.launch { player.statsUpdates.collect { stats -> … } }                       // ~1 s cadence; also player.stats
// PlaybackStats fields are nullable until data flows: videoFps: Double?,
// videoLatency/audioLatency: Duration?, videoBitrateKbps: Double?, …
scope.launch { player.diagnostics().collect { … } }                               // same shared flow every call; one 256-slot drop-oldest buffer
```

`events` and `statsUpdates` are `SharedFlow`s with no replay — `UNDISPATCHED` collection before `play()` is how the demo avoids missing `track.ready`/`playback.start` (`player.init` fires during construction and is unobservable).

## Lifecycle checklist

1. Create `Session(url, parentScope)` → observe state → `connect()`.
2. `subscribe(prefix)` → per broadcast, collect `catalogs()`; catalog stream ending = broadcast offline → `broadcast.close()`.
3. Create `Player(…, parentScope)` → launch events/stats collectors (`UNDISPATCHED`) → wire `SurfaceHolder.Callback` → `play()`.
4. Teardown in order: cancel collector jobs → `player.close()` → `close()` broadcasts and the subscription → `session.close()`.

Everything `AutoCloseable` holds a ref-counted native broadcast handle — leaking one keeps the broadcast open.
