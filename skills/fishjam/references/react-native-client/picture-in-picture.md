# Picture-in-Picture

Picture-in-Picture lets your video call shrink into a floating window when the user backgrounds the app. iOS and Android both support it; the SDK exposes one component (`RTCPIPView`) and two helper functions (`startPIP`, `stopPIP`).

## Plugin prerequisite

```json
{
  "ios": {
    "supportsPictureInPicture": true
  },
  "android": {
    "supportsPictureInPicture": true
  }
}
```

For iOS, `ios.supportsPictureInPicture: true` adds `audio` to `UIBackgroundModes` (and nothing else). You do **not** need `enableVoIPBackgroundMode` for PIP.

## Enabling PIP on a tile (preferred)

Pass a `pip` options object to `RTCView`. The OS auto-promotes the view to the floating window when the app backgrounds — no ref, no `AppState` listener needed.

```tsx
import { RTCView } from '@fishjam-cloud/react-native-client';

<RTCView
  mediaStream={remotePeer.cameraTrack.stream}
  pip={{
    enabled: true,
    startAutomatically: true,
    stopAutomatically: true,
    allowsCameraInBackground: true,
    preferredSize: { width: 16, height: 9 },
  }}
/>
```

`RTCPIPOptions` defaults: `enabled: true`, `startAutomatically: true`, `stopAutomatically: true`, `allowsCameraInBackground: false`. If `pip` is omitted, PIP is disabled for that tile. `allowsCameraInBackground` is iOS 16+.

The canonical example app (`examples/mobile-react-native/fishjam-chat/components/VideosGrid.tsx`) uses exactly this pattern.

## `RTCPIPView` + manual `startPIP` / `stopPIP` (escape hatch)

Use this only if you need imperative control (e.g. trigger PIP from a button mid-foreground). `RTCPIPView` is a `forwardRef` — attach a ref so you can later call `startPIP` / `stopPIP`:

```tsx
import { RTCPIPView, startPIP, stopPIP } from '@fishjam-cloud/react-native-client';
import { useRef } from 'react';
import type { ComponentRef } from 'react';

const pipRef = useRef<ComponentRef<typeof RTCPIPView>>(null);

<RTCPIPView ref={pipRef} mediaStream={remotePeer.cameraTrack.stream} style={{ flex: 1 }} />;

startPIP(pipRef);
stopPIP(pipRef);
```

Only the most recently mounted `RTCPIPView` is eligible for PIP at a time. Pair with `AppState` if you want to auto-trigger on background — but prefer the declarative `pip` prop above.

## Android

Android handles PIP automatically when:

- The activity declares `android:supportsPictureInPicture="true"` (the plugin adds this).
- The view is mounted with `pip={{ enabled: true }}` (or you call `startPIP` explicitly). Without the `pip` prop, no PIP — defaults to off.
- The user presses Home / backgrounds the app.

The declarative `pip={{ enabled: true, startAutomatically: true }}` form works on both platforms.

## Foreground service pairing (Android)

For PIP to keep streaming, the peer connection still needs to be alive — which means a foreground service. Combine `useForegroundService` (see `foreground-service.md`) with `RTCPIPView`.

## VoIP background mode (iOS)

If your app uses CallKit (`useCallKit`), enable `ios.enableVoIPBackgroundMode: true` in the plugin. This lets the WebRTC connection persist while PIP is showing.

## Sources

- `@fishjam-cloud/react-native-webrtc` → `lib/typescript/RTCPIPView.d.ts` (`RTCPIPView` as `forwardRef`, `startPIP(ref)` and `stopPIP(ref)` — synchronous, return `void`)
- `packages/mobile-client/src/index.ts` (re-exports)
- <https://fishjam.swmansion.com/docs/how-to/client/picture-in-picture>
- Cross-reference: `foreground-service.md`, `callkit.md`
