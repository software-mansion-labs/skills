# Permissions

iOS and Android both gate camera and microphone behind explicit user grants. The mobile SDK exposes two hooks that wrap the OS APIs.

Source: `packages/mobile-client/src/hooks/usePermissions.ts` in the web-client-sdk monorepo.

## API

```tsx
import {
  useCameraPermissions,
  useMicrophonePermissions,
  type PermissionStatus,
} from '@fishjam-cloud/react-native-client';

const [queryCamera, requestCamera] = useCameraPermissions();
const [queryMic, requestMic] = useMicrophonePermissions();
```

Both hooks return a **tuple**, not an object:

- `query: () => Promise<PermissionStatus>` — non-prompting check
- `request: () => Promise<PermissionStatus>` — triggers the native dialog

```ts
type PermissionStatus = 'granted' | 'denied' | 'prompt';
```

`'prompt'` means the user hasn't been asked yet (or, on Android, the system allows re-asking).

## Pattern: request before init

```tsx
import { useEffect } from 'react';
import {
  FishjamProvider,
  useCameraPermissions,
  useMicrophonePermissions,
  useInitializeDevices,
} from '@fishjam-cloud/react-native-client';

function PermissionGate() {
  const [queryCamera, requestCamera] = useCameraPermissions();
  const [queryMic, requestMic] = useMicrophonePermissions();
  const { initializeDevices } = useInitializeDevices();

  useEffect(() => {
    (async () => {
      let cam = await queryCamera();
      if (cam !== 'granted') cam = await requestCamera();

      let mic = await queryMic();
      if (mic !== 'granted') mic = await requestMic();

      if (cam === 'granted' && mic === 'granted') {
        await initializeDevices({ enableAudio: true, enableVideo: true });
      } else {
        // Show your "please grant access" UI
      }
    })();
  }, []);

  return null;
}
```

**Always** check the result of `request` — the user may deny, in which case `useInitializeDevices` will succeed with no devices and `useCamera().cameraDeviceError` reports the failure.

## Denied state recovery

Once a user has denied, calling `requestCamera()` again may:

- **Android:** show the dialog again with a "Don't ask again" option, until they tick it.
- **iOS:** silently re-return the previous state — the user must change permission in Settings.

UX recommendation: on `denied`, show a screen explaining how to enable in OS Settings, with a deep-link button (`Linking.openSettings()`).

## iOS Info.plist usage descriptions

iOS requires usage-description strings or the permission request silently fails. The Expo plugin adds defaults; override in `app.json`:

```json
{
  "expo": {
    "ios": {
      "infoPlist": {
        "NSCameraUsageDescription": "MyApp uses the camera for video calls.",
        "NSMicrophoneUsageDescription": "MyApp uses the microphone for calls.",
        "NSPhotoLibraryUsageDescription": "MyApp accesses photos for sharing screenshots."
      }
    }
  }
}
```

These show in the OS permission dialog — be honest and specific. Apple rejects review submissions with generic strings like "We need this."

## Foreground service permissions (Android)

If you want to keep media flowing while the app is backgrounded, you also need foreground service permissions and a foreground service notification. See `foreground-service.md`.

## Screen capture permission

Screen capture has its own permission flow on both platforms:

- **Android** — `MediaProjection` consent prompt fires on each capture attempt; no separate permission to request.
- **iOS** — the user picks your app from the system's broadcast picker (`ScreenCapturePickerView`); see `screen-sharing.md`.

`useCameraPermissions` / `useMicrophonePermissions` do not cover screen capture.

## Sources

- `packages/mobile-client/src/hooks/usePermissions.ts` (web-client-sdk monorepo)
- <https://fishjam.swmansion.com/docs/tutorials/react-native-quick-start>
