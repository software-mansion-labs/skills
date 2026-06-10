---
name: fishjam-react-native-client
description: "React Native / Expo SDK for Fishjam — video/audio streaming on iOS and Android. Use when writing a React Native or Expo app that calls Fishjam, configures the Fishjam Expo plugin, sets up permissions, runs background streaming, integrates CallKit, or renders RTCView. Trigger on: '@fishjam-cloud/react-native-client', 'fishjam expo plugin', 'FishjamProvider mobile', 'useCameraPermissions', 'useMicrophonePermissions', 'useForegroundService', 'useCallKit', 'useCallKitEvent', 'useCallKitService', 'RTCView', 'RTCPIPView', 'ScreenCapturePickerView', 'startPIP', 'stopPIP', 'AudioDeviceType', 'useAudioOutput', '@fishjam-cloud/react-native-webrtc', 'fishjam react native', 'expo fishjam', 'fishjam ios', 'fishjam android', 'broadcast extension'. Re-exports @fishjam-cloud/react-client hooks plus mobile-only: permissions, foreground service, iOS broadcast extension, audio routing, CallKit, Expo config plugin."
license: MIT
---

# Fishjam React Native Client

`@fishjam-cloud/react-native-client` — Fishjam's SDK for **React Native / Expo apps** on iOS and Android. Re-exports the web `@fishjam-cloud/react-client` hook surface, overrides some hooks with mobile-aware versions, and adds mobile-only APIs.

> **Read `../platform/SKILL.md` first** for the domain model.
>
> **Read `../react-client/SKILL.md` for the shared hook catalog.** This skill is the *mobile delta* — permissions, foreground service, screen-share broadcast extension, CallKit, audio routing, the Expo plugin. The hooks that work identically to the web SDK aren't re-documented here.

## Minimal app

```bash
npm install @fishjam-cloud/react-native-client @fishjam-cloud/react-native-webrtc
# @fishjam-cloud/react-native-webrtc is a peerDependency — install it explicitly.
# Do NOT install upstream `react-native-webrtc` (the unforked package) — it will collide.
```

`app.json` (Expo):

```json
{
  "expo": {
    "plugins": [
      [
        "@fishjam-cloud/react-native-client",
        {
          "android": {
            "enableForegroundService": true,
            "enableScreensharing": true
          },
          "ios": {
            "enableScreensharing": true,
            "broadcastExtensionTargetName": "MyAppScreenSharing",
            "broadcastExtensionDisplayName": "MyApp Screen Sharing",
            "appGroupContainerId": "group.com.myapp.screensharing",
            "mainTargetName": "MyApp"
          }
        }
      ]
    ]
  }
}
```

Then `npx expo prebuild` to generate native code.

```tsx
import { FishjamProvider } from '@fishjam-cloud/react-native-client';

export default function App() {
  return (
    <FishjamProvider fishjamId={process.env.EXPO_PUBLIC_FISHJAM_ID!}>
      <Room />
    </FishjamProvider>
  );
}
```

## Hook delta vs web (`../react-client/`)

### Re-exported unchanged

These hooks behave identically to the web SDK — see `../react-client/` references.

`useConnection`, `useDataChannel`, `useSandbox`, `useUpdatePeerMetadata`, `useVAD`, `Variant`, `InitializeDevicesSettings`.

> **`useSandbox` `RoomType` divergence.** The web SDK accepts `'audio_only'` (underscore); the React Native SDK accepts `'audio-only'` (hyphen). Passing the wrong form on either platform makes the Sandbox API reject the request. (Web wire format is `audio_only`; the RN hook silently sends `'audio-only'` and the server validation 400s.)

### Overridden (mobile-specific implementations of same API)

Same call signatures as web; internals adapted to React Native. Read `../react-client/` references for the API, then this skill's `mobile-hook-overrides.md` for behavioral deltas.

`useCamera`, `useMicrophone`, `useScreenShare`, `useCustomSource`, `useInitializeDevices`, `useLivestreamStreamer`, `useLivestreamViewer`, `usePeers`.

### Mobile-only (no web equivalent)

| Hook / API | Purpose | Reference |
|---|---|---|
| `useCameraPermissions`, `useMicrophonePermissions` | Query / request iOS-Android device permissions | `permissions.md` |
| `useForegroundService` | Keep audio/video/screen-share alive when the app is backgrounded on Android | `foreground-service.md` |
| `useCallKit`, `useCallKitEvent`, `useCallKitService` | iOS CallKit integration (treats Fishjam sessions as native phone calls) | `callkit.md` |
| `useAudioOutput` + `AudioDeviceType` | Switch audio routing (speaker / earpiece / Bluetooth) | `audio-output.md` |
| `RTCView`, `RTCPIPView` | Native components for rendering video tracks (in place of `<video>`) | `rtcview.md` |
| `ScreenCapturePickerView` + `startPIP` / `stopPIP` | Native components / functions for iOS screen picker and Picture-in-Picture | `screen-sharing.md`, `picture-in-picture.md` |

### `FishjamProvider` differences

The mobile provider is a thin wrapper around the React one but **drops two props**:

- `persistLastDevice` — not supported on mobile (no `localStorage`; OS manages device selection).
- `fishjamClient` — internal-only on mobile (the SDK constructs its own with `clientType: 'mobile'`).

All other props (`fishjamId`, `reconnect`, `constraints`, `bandwidthLimits`, `videoConfig`, `audioConfig`, `debug`) work the same as web.

## Key rules

- **Install `@fishjam-cloud/react-native-webrtc` (the fork) — but NOT the upstream `react-native-webrtc`.** The fork is a `peerDependency` of `@fishjam-cloud/react-native-client`, so you must add it explicitly. Installing the upstream (unforked) `react-native-webrtc` will collide with the fork and break the native build.
- **The Expo plugin is required even in bare workflow.** For Expo CNG (managed): register in `app.json`/`app.config.ts`. For bare: still register and run `npx expo prebuild`, or do the manual native steps (`native-setup.md`).
- **Request permissions before `useInitializeDevices`.** iOS in particular will silently fail device init if `NSCameraUsageDescription` / `NSMicrophoneUsageDescription` isn't declared and granted.
- **Foreground service is mandatory for backgrounded streaming on Android.** When the user puts the app in background, Android suspends the WebRTC peer connection unless a foreground service is running. `useForegroundService` solves this.
- **iOS screen sharing requires a Broadcast Extension.** A separate iOS target generated by the Expo plugin. Configure `broadcastExtensionTargetName` + `appGroupContainerId` and the plugin scaffolds it during `prebuild`.
- **iOS PIP needs the plugin flag.** `ios.supportsPictureInPicture: true` and `enableVoIPBackgroundMode: true` for CallKit-paired calls.

## References

| File | When to read |
|---|---|
| `native-setup.md` | Install + Expo plugin options (Android + iOS), bare workflow manual steps, prebuild flow. |
| `permissions.md` | `useCameraPermissions`, `useMicrophonePermissions`, request-before-init pattern. |
| `rtcview.md` | Rendering video — `RTCView` props (`streamURL`, `trackId`, `mediaStream`), `RTCPIPView`, mirror / styling. |
| `screen-sharing.md` | Android MediaProjection vs iOS Broadcast Extension; `ScreenCapturePickerView`; app-group pitfalls. |
| `foreground-service.md` | `ForegroundServiceConfig`, Android permissions, iOS background modes. |
| `picture-in-picture.md` | `RTCPIPView`, `startPIP`, `stopPIP`, plugin flags. |
| `callkit.md` | `useCallKit`, `useCallKitEvent`, `useCallKitService`, `CallKitConfig`, `CallKitAction`. |
| `audio-output.md` | `useAudioOutput`, `AudioDeviceType`, switching speaker / earpiece / Bluetooth. |
| `mobile-hook-overrides.md` | Per-hook behavioral deltas vs the web SDK. |
| `example-projects.md` | Runnable apps in the `fishjam-cloud/examples` repo under `mobile-react-native/`. |
