# Native Setup

## Install

```bash
npm install @fishjam-cloud/react-native-client @fishjam-cloud/react-native-webrtc react-native-get-random-values
# or yarn add @fishjam-cloud/react-native-client @fishjam-cloud/react-native-webrtc react-native-get-random-values
```

Both `@fishjam-cloud/react-native-webrtc` and `react-native-get-random-values` are declared as **peer dependencies** of `@fishjam-cloud/react-native-client` — you must install them explicitly alongside the wrapper. They are not transitively pulled in. The wrapper imports the `react-native-get-random-values` polyfill itself at module load (used by the WebRTC stack); if the package or its native module isn't linked the SDK throws a setup error at startup.

The wrapper package ships:

- TypeScript wrapper hooks (`useCamera`, `useScreenShare`, etc).
- An Expo config plugin (`@fishjam-cloud/react-native-client/plugin`).

The native iOS/Android modules live in the separate `@fishjam-cloud/react-native-webrtc` peer dep above.

**Do not** install upstream `react-native-webrtc` (the unforked package) — it will collide with the fork at build time.

## Expo plugin

Register in `app.json`, `app.config.ts`, or `app.config.js`:

```json
{
  "expo": {
    "plugins": [
      [
        "@fishjam-cloud/react-native-client",
        {
          "android": {
            "enableForegroundService": true,
            "enableScreensharing": true,
            "supportsPictureInPicture": true
          },
          "ios": {
            "enableScreensharing": true,
            "supportsPictureInPicture": true,
            "broadcastExtensionTargetName": "MyAppScreenSharing",
            "broadcastExtensionDisplayName": "MyApp Screen Sharing",
            "appGroupContainerId": "group.com.myapp.screensharing",
            "mainTargetName": "MyApp",
            "iphoneDeploymentTarget": "14.0",
            "enableVoIPBackgroundMode": true
          }
        }
      ]
    ]
  }
}
```

Then:

```bash
npx expo prebuild --clean
# (or `npx expo run:ios` / `npx expo run:android` which prebuild as a side effect)
```

### Plugin options

```ts
type FishjamPluginOptions = {
  android?: {
    enableForegroundService?: boolean;        // Adds FOREGROUND_SERVICE permissions; enables useForegroundService
    enableScreensharing?: boolean;            // Adds MediaProjection bits to AndroidManifest
    supportsPictureInPicture?: boolean;       // Enables system-PIP attributes on the main activity
  };
  ios?: {
    enableScreensharing?: boolean;            // Generates a Broadcast Extension iOS target
    supportsPictureInPicture?: boolean;       // Adds background-modes for PIP
    broadcastExtensionTargetName?: string;    // Name of the generated screen-sharing target (Xcode)
    broadcastExtensionDisplayName?: string;   // Display name shown in iOS share sheet
    appGroupContainerId?: string;             // App Group container — MUST match in both main app + extension
    mainTargetName?: string;                  // Your main iOS target's name (default: detected)
    iphoneDeploymentTarget?: string;          // Sets IPHONEOS_DEPLOYMENT_TARGET; default '15.1' (see packages/mobile-client/plugin/src/withFishjamIos.ts)
    enableVoIPBackgroundMode?: boolean;       // Adds 'voip' background mode for CallKit/PushKit
  };
} | undefined;
```

### Required: when to enable each flag

| Need | Set |
|---|---|
| Stream when the user backgrounds the app on Android | `android.enableForegroundService: true` |
| Share screen on Android | `android.enableScreensharing: true` (and use `useForegroundService` at runtime) |
| Share screen on iOS | `ios.enableScreensharing: true` and `broadcastExtensionTargetName`. `appGroupContainerId` is optional (defaults to `group.<bundleIdentifier>`) — override only if you need a custom App Group identifier |
| Picture-in-Picture on iOS | `ios.supportsPictureInPicture: true` |
| Picture-in-Picture on Android | `android.supportsPictureInPicture: true` |
| CallKit (treat sessions as phone calls) | `ios.enableVoIPBackgroundMode: true` |

## iOS Broadcast Extension — the gotcha

`broadcastExtensionTargetName` and `appGroupContainerId` **must match** between the main app and the generated extension. If you change one, change the other and re-prebuild.

`appGroupContainerId` follows the Apple convention `group.<reverse-domain>.<suffix>` — e.g. `group.com.acme.myapp.screensharing`. Both targets need entitlement to the same App Group; the plugin sets this up if you let it.

If you've manually edited the Xcode project after `prebuild`, expect drift on the next prebuild. Either let the plugin own those changes or commit the changes and re-run `prebuild --no-install`.

## iOS Info.plist usage descriptions

You must add these to `expo.ios.infoPlist` yourself — the plugin does **not** inject them:

```json
{
  "expo": {
    "ios": {
      "infoPlist": {
        "NSCameraUsageDescription": "MyApp needs your camera to let you appear on video calls.",
        "NSMicrophoneUsageDescription": "MyApp needs your microphone to let you talk on calls."
      }
    }
  }
}
```

Without these, iOS denies the permission request silently.

## Android permissions

**Added by the Fishjam plugin** when you enable `enableForegroundService` (and `enableScreensharing` for the media-projection variant):

| Feature | Permission |
|---|---|
| Foreground service (audio/video) | `android.permission.FOREGROUND_SERVICE`, `FOREGROUND_SERVICE_CAMERA`, `FOREGROUND_SERVICE_MICROPHONE` |
| Screen sharing | `FOREGROUND_SERVICE_MEDIA_PROJECTION` |

**You must declare these yourself** in `expo.android.permissions` (or `AndroidManifest.xml` in bare workflow) — they come from the webrtc fork's manifest or your project, not from the Fishjam plugin:

| Feature | Permission |
|---|---|
| Camera | `android.permission.CAMERA` |
| Microphone | `android.permission.RECORD_AUDIO` |
| Network | `android.permission.INTERNET`, `ACCESS_NETWORK_STATE` |

If you customize `AndroidManifest.xml` manually (bare workflow), make sure all of these are present for the features you use.

## Bare React Native workflow

If you're not using Expo's prebuild flow, you have two options:

1. **Adopt `expo-modules-core`** — install Expo's autolinking, then register the Fishjam plugin the same as above and run `npx expo prebuild`. This is the path of least resistance.
2. **Manual native steps** — open the iOS/Android projects in Xcode/Android Studio and apply the plugin's effects by hand. See the upstream React Native quick-start (`docs/tutorials/react-native-quick-start`) for the latest exact steps.

## Sources

- `packages/mobile-client/plugin/src/types.ts` (plugin options)
- `packages/mobile-client/plugin/src/withFishjam*.ts` (effects)
- `examples/mobile-react-native/minimal-react-native/app.json` (canonical example) — <https://github.com/fishjam-cloud/examples/tree/main/mobile-react-native/minimal-react-native>
- <https://fishjam.swmansion.com/docs/tutorials/react-native-quick-start>
