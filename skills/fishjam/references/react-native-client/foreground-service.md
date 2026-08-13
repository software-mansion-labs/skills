# Foreground Service (Android)

Android suspends WebRTC peer connections when the app is backgrounded — unless the app runs a **foreground service** with an ongoing notification. The SDK provides `useForegroundService` to manage this.

> This hook is a no-op on iOS. On iOS, backgrounded streaming relies on **background modes** declared in the Expo plugin (`enableVoIPBackgroundMode`, etc.) — see `callkit.md` for the VoIP scenario.

Source: `packages/mobile-client/src/useForegroundService.ts`.

## Plugin prerequisite

```json
{
  "android": { "enableForegroundService": true }
}
```

This adds the necessary `<service>` declarations to `AndroidManifest.xml` and the right `FOREGROUND_SERVICE_*` permissions.

## Hook usage

```tsx
import { useForegroundService } from '@fishjam-cloud/react-native-client';

function CallScreen() {
  useForegroundService({
    enableCamera: true,
    enableMicrophone: true,
    enableScreenSharing: false,
    channelId: 'fishjam-call',
    channelName: 'Active call',
    notificationTitle: 'On a call',
    notificationContent: 'Tap to return to MyApp.',
  });

  return /* call UI */;
}
```

### `ForegroundServiceConfig`

```ts
type ForegroundServiceConfig = {
  enableCamera?: boolean;
  enableMicrophone?: boolean;
  enableScreenSharing?: boolean;
  channelId?: string;
  channelName?: string;
  notificationTitle?: string;
  notificationContent?: string;
};
```

| Field | Meaning |
|---|---|
| `enableCamera` / `enableMicrophone` / `enableScreenSharing` | Toggle each foreground-service capability. Match what your app is actually using. Without the right capability flag, Android may suspend the corresponding track. |
| `channelId` | A unique-per-package notification channel ID. Reused across calls — pick a stable string. |
| `channelName` | User-visible channel name in Settings → Apps → Notifications. |
| `notificationTitle` | First line of the ongoing notification. |
| `notificationContent` | Second line. |

The hook runs the service as long as it's mounted. Unmount → service stops, notification disappears.

## Pattern: foreground service tied to call lifecycle

Mount on the `CallScreen` component; unmount when the user leaves the screen.

```tsx
function App() {
  const [inCall, setInCall] = useState(false);
  return inCall ? <CallScreen onLeave={() => setInCall(false)} /> : <Lobby onJoin={() => setInCall(true)} />;
}

function CallScreen({ onLeave }: { onLeave: () => void }) {
  useForegroundService({
    enableCamera: true,
    enableMicrophone: true,
    channelId: 'fishjam-call',
    channelName: 'Call',
    notificationTitle: 'Call in progress',
    notificationContent: 'Tap to return.',
  });
  // ...
}
```

## Without the foreground service

If you stream on Android **without** `useForegroundService`:

- When the user opens the home screen / switches apps, the OS kills the WebRTC peer connection within seconds.
- Other peers see `peerDisconnected` for your peer.
- When the user returns, the SDK auto-reconnects (if `reconnect` is enabled), but the gap is jarring for everyone in the call.

## iOS background streaming

iOS doesn't use the foreground-service concept. Instead:

- Add `background_modes` to your `Info.plist` (the Expo plugin does this for the modes you opt in to).
- `ios.supportsPictureInPicture: true` adds `audio` to `UIBackgroundModes`; `ios.enableVoIPBackgroundMode: true` adds `voip`. No dedicated flag exists for audio-only persistence — enable PiP support (which writes `audio` as a side-effect) or add `audio` to `UIBackgroundModes` manually in a bare workflow.

Without these, iOS suspends your app and the WebRTC connection within ~30 seconds of backgrounding.

## Permissions required

Plugin handles these for Expo CNG. In bare workflow add to `AndroidManifest.xml`:

```xml
<uses-permission android:name="android.permission.FOREGROUND_SERVICE" />
<uses-permission android:name="android.permission.FOREGROUND_SERVICE_CAMERA" />
<uses-permission android:name="android.permission.FOREGROUND_SERVICE_MICROPHONE" />
<uses-permission android:name="android.permission.FOREGROUND_SERVICE_MEDIA_PROJECTION" />
```

Android 14+ enforces these strictly — the service won't start without the matching permission for the type you're using.

## Sources

- `packages/mobile-client/src/useForegroundService.ts`
- `@fishjam-cloud/react-native-webrtc` (underlying foreground-service implementation)
- <https://fishjam.swmansion.com/docs/how-to/client/background-streaming>
