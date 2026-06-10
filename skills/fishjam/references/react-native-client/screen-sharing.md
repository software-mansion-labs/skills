# Screen Sharing

Screen sharing in React Native differs significantly between iOS and Android. Both work through the same `useScreenShare` hook (inherited from `../react-client/`), but the native plumbing is platform-specific.

## API at the React layer

```tsx
import { useScreenShare } from '@fishjam-cloud/react-native-client';

const {
  startStreaming,
  stopStreaming,
  stream,
  presentBroadcastPicker, // iOS-only: opens the system Stop sheet for an active broadcast
} = useScreenShare();
const isSharing = stream != null;
```

`useScreenShare` extends the web hook (see `../react-client/devices.md`) with one extra mobile-only function, `presentBroadcastPicker: () => Promise<void>`. There's no dedicated `isSharing` flag — derive it from `stream != null` (or `videoTrack != null`).

## Android (MediaProjection)

### Android plugin config

```json
{
  "android": {
    "enableScreensharing": true,
    "enableForegroundService": true
  }
}
```

`enableForegroundService` is required because Android only allows MediaProjection when a foreground service is running. Note that in the current plugin (`packages/mobile-client/plugin/src/withFishjamAndroid.ts`), the `FOREGROUND_SERVICE_*` permissions — including `FOREGROUND_SERVICE_MEDIA_PROJECTION` — are gated by `enableForegroundService`, not by `enableScreensharing`. So:

- Enable **both** `enableScreensharing: true` and `enableForegroundService: true` for Android screen share to work.
- `enableScreensharing` adds the `MediaProjection` service declarations to `AndroidManifest.xml`.
- `enableForegroundService` adds the runtime permissions: `FOREGROUND_SERVICE`, `FOREGROUND_SERVICE_CAMERA`, `FOREGROUND_SERVICE_MICROPHONE`, `FOREGROUND_SERVICE_MEDIA_PROJECTION`.

### Runtime — start the foreground service

```tsx
import { useForegroundService, useScreenShare } from '@fishjam-cloud/react-native-client';

useForegroundService({
  enableScreenSharing: true,
  enableCamera: true,
  enableMicrophone: true,
  channelId: 'fishjam-channel',
  channelName: 'Fishjam',
  notificationTitle: 'In a call',
  notificationContent: 'You are sharing your screen.',
});

const { startStreaming } = useScreenShare();

<Button title="Share screen" onPress={() => startStreaming()} />
```

The system shows a consent dialog before capture begins. On user dismissal, `startStreaming` rejects.

### What gets shared on Android

The entire device screen (apps switch, notifications, etc.). Android does not have a per-app or per-window picker like macOS.

## iOS (Broadcast Extension)

iOS only allows screen sharing through a separate **Broadcast Extension** target. The Expo plugin generates this target during `prebuild`.

### iOS plugin config

In your `app.json` / `app.config.js`, plugins use a `[name, options]` tuple form:

```json
{
  "expo": {
    "plugins": [
      ["@fishjam-cloud/react-native-client", {
        "ios": {
          "enableScreensharing": true,
          "broadcastExtensionTargetName": "MyAppScreenSharing",
          "broadcastExtensionDisplayName": "MyApp Screen Sharing",
          "appGroupContainerId": "group.com.myapp.screensharing",
          "mainTargetName": "MyApp"
        }
      }]
    ]
  }
}
```

`appGroupContainerId` is **optional** — it defaults to `group.<your-bundle-identifier>`. Override it only if you need a custom App Group identifier (e.g. shared with other apps).

After `npx expo prebuild`, your Xcode project gains a new target with the configured name.

### `appGroupContainerId` pitfall

This is the **single most-bitten-by misconfiguration** in iOS screen sharing.

- Both the main app and the broadcast extension must belong to the same App Group.
- The plugin sets entitlements for the configured `appGroupContainerId` on both targets.
- If you regenerate prebuild with a different `appGroupContainerId`, the extension's entitlements change but cached signing identities can lag — clean Xcode derived data and re-sign.
- If you manually edit one of the entitlements files (e.g. add a different group), the broadcast extension can't reach the main app's media plumbing and screen sharing silently fails.

Format the value as `group.<reverse-dns>.<suffix>` and pick something you won't regret — changing it later is painful.

### Triggering the iOS picker

On iOS, call `startStreaming()` — the SDK internally presents the system broadcast picker (`BroadcastPickerHelper.presentSystemPickerWithError`) and the user picks your extension from the system UI. To stop a live broadcast cleanly, call `await presentBroadcastPicker()` (returned by `useScreenShare()`); this opens the system Stop sheet without surfacing the "Screen sharing stopped" error dialog that `stopStreaming()` triggers.

```tsx
import { useScreenShare } from '@fishjam-cloud/react-native-client';

const { startStreaming, stopStreaming, presentBroadcastPicker, stream } = useScreenShare();
const isSharing = stream != null;

// Start: system picker pops up automatically.
const onStart = () => startStreaming();

// Stop: open the system Stop sheet (preferred over stopStreaming on iOS).
const onStop = () => presentBroadcastPicker();
```

### What gets shared on iOS

The user's screen across all apps until they tap the stop button (or you call `stopStreaming()`). iOS doesn't allow restricting to your app's window.

## Stopping

```tsx
stopStreaming();
```

Tears down both the WebRTC track and the platform-specific capture pipeline. On Android, the foreground service notification updates. On iOS, the system broadcast bar goes away.

## Cross-platform UI sketch

```tsx
import { Platform, Button } from 'react-native';
import {
  useScreenShare,
  useForegroundService,
} from '@fishjam-cloud/react-native-client';

function ScreenShareButton() {
  const { startStreaming, stopStreaming, presentBroadcastPicker, stream } = useScreenShare();
  const isSharing = stream != null;

  useForegroundService({ enableScreenSharing: true });

  const handlePress = () => {
    if (isSharing) {
      // iOS: prefer the system Stop sheet to avoid the "Screen sharing stopped" error dialog.
      if (Platform.OS === 'ios') presentBroadcastPicker();
      else stopStreaming();
      return;
    }
    startStreaming();
  };

  return (
    <Button
      title={isSharing ? 'Stop sharing' : 'Share screen'}
      onPress={handlePress}
    />
  );
}
```

## Sources

- `packages/mobile-client/src/index.ts` (re-exports)
- `packages/mobile-client/plugin/src/withFishjamIos.ts` (broadcast extension generation)
- `packages/mobile-client/plugin/src/withFishjamAndroid.ts` (MediaProjection wiring)
- <https://fishjam.swmansion.com/docs/how-to/client/screensharing>
- Cross-reference: `foreground-service.md`, `picture-in-picture.md`
