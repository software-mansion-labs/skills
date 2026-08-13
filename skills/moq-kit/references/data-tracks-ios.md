# Data tracks and raw media consumers — iOS (Swift)

Shared semantics (names agreed out of band, `TrackDelivery` modes, self-broadcast filtering, shared upstream latency) are in SKILL.md.

## Publish a data track

```swift
let emitter = DataTrackEmitter()               // retain it — it's your send handle
let publisher = try Publisher()
publisher.addDataTrack(name: "chat", source: emitter)
try await session.publish(path: "chat/room-1", publisher: publisher)
try await publisher.start()
try emitter.send(JSONEncoder().encode(message)) // silent no-op before start / after stop; throws only on a live write failure
```

High-frequency sends are fine — the demo's MoQBoy controller sends held-button JSON every 10 ms.

## Subscribe to a data track

```swift
let track = try broadcast.subscribeTrack(name: "chat", delivery: .arrival)
for try await object in track.objects {        // TrackObject: payload, groupSequence, objectIndex
    let msg = try JSONDecoder().decode(ChatMessage.self, from: object.payload)
}
// track.close() when done (also closes on deinit and stream termination)
```

## MediaTrack — raw compressed frames

Bypass the player and receive encoded media frames (e.g. to feed your own decoder or remux):

```swift
let media = try broadcast.subscribeMedia(
    MediaTrackRequest(name: trackName, container: .cmaf(initializationData: initData),
                      targetBuffering: .milliseconds(100)),
    options: MediaTrackOptions(bufferingPolicy: .bufferingNewest(30)))
for try await frame in media.frames { … }      // MediaFrame: payload, timestampUs, keyframe
```

Containers: `.legacy`, `.cmaf(initializationData:)`, `.loc`. Buffering policy: `.unbounded` or `.bufferingNewest(limit)`. Watch `media.state` (`.idle/.active/.closed/.error(String)`); `close()` when done.

## AudioDataStream — decoded PCM

Decoded audio samples for processing (metering, transcription, effects) — this does **not** play audio:

```swift
let stream = try AudioDataStream(catalog: catalog, trackName: audioTrack,
                                 format: AudioDataFormat(sampleFormat: .float32))
for try await audio in stream.audio { … }      // AudioData: bytes, sampleFormat, timestampUs, sampleRate, channelCount, frameCount
```

Formats: `.float32` / `.int16`, optional resample (`sampleRate`) and channel-count conversion. Bounded buffer — a slow consumer gets the newest data. Close it like every other consumer.
