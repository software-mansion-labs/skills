# Screen capture — iOS

Pick the right variant first: in-app `ScreenCapture` runs in the host app process, so frame delivery can stop once the app backgrounds — don't rely on it surviving an app switch. Full-device capture across apps requires the ReplayKit Broadcast Upload Extension path.

## In-app `ScreenCapture`

Captures the app's own screen via `RPScreenRecorder`:

```swift
let screen = ScreenCapture()
try await screen.start()
publisher.addVideoTrack(name: "screen", source: screen.videoSource)
publisher.addAudioTrack(name: "screen-audio", source: screen.audioSource)
// … session.publish + publisher.start() as usual; later: await screen.stop()
```

## ReplayKit Broadcast Upload Extension (full-device)

The extension runs out-of-process and opens its **own** relay connection; the host app hands it the relay URL and path through App Group `UserDefaults`.

Requirements:

- A Broadcast Upload Extension target (`NSExtensionPointIdentifier = com.apple.broadcast-services-upload`, principal class = your handler).
- A real, provisioned App Group (`com.apple.security.application-groups`) on **both** the host app and the extension — `ReplayKitBroadcastDescriptorStore` throws `.invalidAppGroup` if the container isn't accessible.
- `NSMicrophoneUsageDescription` in the extension too if the mic button is offered.

Extension — the whole thing can be three lines:

```swift
final class SampleHandler: MoQReplayKitBroadcastSampleHandler {
    override var replayKitAppGroupIdentifier: String? { "group.com.example.app" }
}
```

`MoQReplayKitBroadcastSampleHandler` reads the descriptor, runs a `ReplayKitBroadcastPipeline` (connect → publish → encode screen/app-audio/mic samples), and handles pause/resume/finish. Override `makeReplayKitBroadcastConfiguration(setupInfo:)` to customize track names or encoder configs (`ReplayKitBroadcastConfiguration`: `videoTrackName = "screen"`, `appAudioTrackName = "screen-audio"`, `micAudioTrackName = nil` by default).

Host app — write the descriptor **before** the user starts the broadcast, then show the system picker:

```swift
try ReplayKitBroadcastDescriptorStore(appGroupIdentifier: "group.com.example.app")
    .save(ReplayKitBroadcastDescriptor(relayURL: relayURL, broadcastPath: "live/screen"))
```

```swift
// RPSystemBroadcastPickerView wrapped for SwiftUI; the USER starts the broadcast from it
let picker = RPSystemBroadcastPickerView()
picker.preferredExtension = "com.example.app.broadcast" // extension bundle id
picker.showsMicrophoneButton = true
```

Clear the descriptor (`store.clear()`) when the feature is turned off so a stale URL isn't picked up later.
