# Data tracks and raw media consumers — Android (Kotlin)

Shared semantics (names agreed out of band, `TrackDelivery` modes, self-broadcast filtering, shared upstream latency) are in SKILL.md.

## Publish a data track

```kotlin
val emitter = DataTrackEmitter()
val publisher = Publisher()
publisher.addDataTrack(name = "chat", emitter = emitter)
session.publish(path = "chat/room-1", publisher = publisher)
publisher.start()
emitter.send(json.encodeToString(message).toByteArray()) // ignored before start / after stop
```

`send()` never throws — it silently drops before `start()` and after `stop()`, so gate sends on publisher state if delivery matters.

## Subscribe to a data track

```kotlin
val track = broadcast.subscribeTrack(name = "chat", delivery = TrackDelivery.Arrival)
track.objects.collect { obj -> decode(obj.payload) }   // cold Flow, SINGLE collector
// close() the track (and the broadcast) in a finally block
```

`TrackObject`: `payload: ByteArray`, `groupSequence: ULong`, `objectIndex: ULong`.

## MediaTrack — raw compressed frames

Bypass the player and receive encoded media frames (e.g. to feed your own decoder or remux):

```kotlin
val media = broadcast.subscribeMedia(
    MediaTrackRequest(name = trackName, container = MediaContainer.Cmaf(initData),
                      targetBuffering = Duration.ofMillis(100)),
    MediaTrackOptions(bufferingPolicy = MediaTrackBufferingPolicy.BufferingNewest(30)))
media.frames.collect { frame -> … }            // MediaFrame: payload, timestampUs, keyframe; cold, SINGLE collector
```

Containers: `Legacy`, `Cmaf(initializationData)`, `Loc`. Buffering policy: `Unbounded` or `BufferingNewest(limit)`. Watch `media.state: StateFlow<MediaTrackState>` (`Idle/Active/Closed/Error`); `close()` when done.

## AudioDataStream — decoded PCM

Decoded audio samples for processing (metering, transcription, effects) — this does **not** play audio:

```kotlin
val stream = AudioDataStream(catalog, trackName = audioTrack,
                             format = AudioDataFormat(sampleFormat = AudioSampleFormat.Float32),
                             targetBuffering = Duration.ofMillis(100)) // default
stream.audio.collect { audio -> … }            // AudioData: bytes, timestampUs, sampleRate, channelCount, sampleFormat, frameCount
```

Overloads take either a `trackName` or an `AudioTrackInfo` from the catalog (unknown names throw `IllegalArgumentException`), or a `Broadcast` + `AudioTrackRequest` to bypass the catalog.

Single collector; bounded buffer (drops oldest when the consumer is slow). Formats: `Float32` / `Int16`, optional resample (`sampleRate`) and channel-count conversion. Close it like every other consumer.
