# Example Projects

Runnable Expo apps demonstrating the mobile client SDK. All live in the [`fishjam-cloud/examples`](https://github.com/fishjam-cloud/examples) repo under [`mobile-react-native/`](https://github.com/fishjam-cloud/examples/tree/main/mobile-react-native).

Use these as **reference implementations**. Copy the parts you need (`app.json` plugin config, hook usage patterns, navigation skeleton) — don't try to run the whole tree if you only want one feature.

## `minimal-react-native/`

Smallest possible Fishjam React Native app. Token input box, join/leave button, local + remote video tiles. Good for:

- Verifying your install / Expo plugin config works.
- Seeing the bare-minimum hook set: `FishjamProvider`, `useConnection`, `useCamera`, `useMicrophone`, `usePeers`.
- Confirming `RTCView` rendering for local and remote peers.

Start here if you've never built a Fishjam app before.

## `fishjam-chat/`

Full video-chat app with Expo Router, real navigation, screen share, CallKit integration, foreground service, and runtime Fishjam ID swap. Demonstrates:

- Permission flow via `useCameraPermissions` / `useMicrophonePermissions`.
- Foreground service for backgrounded streaming.
- Screen share on both Android and iOS (with the broadcast extension setup in `app.json`).
- CallKit (`useCallKit`, `useCallKitService`) integration for treating sessions as native phone calls on iOS.
- Picture-in-Picture via `RTCView`'s `pip` prop.
- Multi-screen navigation (lobby → call → leave).
- Dynamic Fishjam ID swap (switch between dev/staging/prod instances at runtime).

This is the production-flavor example. Read its `app.json` for the canonical Expo plugin configuration.

## `blur-example/`

Tiny app focused on the **camera middleware** pattern — applies a background blur effect using a native processing library and exposes a toggle. Demonstrates:

- `setCameraTrackMiddleware` lifecycle.
- Performance considerations (frame-rate cost of GPU processing).

Use as a starting template for any custom video effect (LUTs, virtual backgrounds, filters).

## `text-chat/`

Demonstrates `useDataChannel` — peer-to-peer text messaging without any media. Useful for:

- Understanding the data-channel lifecycle (`initializeDataChannel` → wait for `dataChannelReady` → publish/subscribe).
- Building chat alongside a call.
- Seeing how `DataChannelOptions` map to RTC reliability modes.

## `video-player/`

Livestream **streamer + viewer** combo. Three-mode UI (`'streamer' | 'viewer' | 'none'`) — connects to a `livestream` room either as the broadcaster or as a WHEP viewer. Demonstrates:

- `useLivestreamStreamer` for the broadcasting side (single peer).
- `useLivestreamViewer().connect({ token })` for private livestreams.
- `useLivestreamViewer().connect({ streamId })` for public livestreams.
- `RTCView` for WHEP playback.
- Disconnect / reconnect flow.

Pair with `web-react/livestreaming/` for cross-platform tests.

## Running the examples

```bash
# clone the examples repo:
git clone https://github.com/fishjam-cloud/examples
cd examples/mobile-react-native/minimal-react-native

yarn install
npx expo prebuild --clean
npx expo run:ios       # or :android
```

You'll need an `EXPO_PUBLIC_FISHJAM_ID` (and optionally `EXPO_PUBLIC_FISHJAM_SANDBOX_API_URL` for the no-backend flow) — the examples read these from environment variables defined per-app.

## Sources

- [`fishjam-cloud/examples`](https://github.com/fishjam-cloud/examples) — example apps across platforms (mobile, web, backends); React Native apps live under [`mobile-react-native/`](https://github.com/fishjam-cloud/examples/tree/main/mobile-react-native)
- Each app's own README for its specific setup instructions
