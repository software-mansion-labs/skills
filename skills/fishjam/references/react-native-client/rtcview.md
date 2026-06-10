# RTCView (and RTCPIPView)

`RTCView` is the React Native equivalent of an HTML `<video>` element — it renders a `MediaStream` natively. On the web you'd use `<video ref={el => { if (el) el.srcObject = stream; }} />`; on RN, use `<RTCView mediaStream={stream} />`.

```tsx
import { RTCView, RTCPIPView } from '@fishjam-cloud/react-native-client';
```

These are re-exported from `@fishjam-cloud/react-native-webrtc` (the fork; installed as a peer dependency). Source for the wrapper: `packages/mobile-client/src/overrides/RTCView.tsx`.

## Basic usage

```tsx
import { View } from 'react-native';
import { RTCView, usePeers } from '@fishjam-cloud/react-native-client';

function CallScreen() {
  const { localPeer, remotePeers } = usePeers();

  return (
    <View style={{ flex: 1 }}>
      {localPeer?.cameraTrack && (
        <RTCView
          mediaStream={localPeer.cameraTrack.stream}
          style={{ width: 200, height: 200 }}
          mirror
        />
      )}
      {remotePeers.map((peer) =>
        peer.cameraTrack ? (
          <RTCView
            key={peer.id}
            mediaStream={peer.cameraTrack.stream}
            style={{ width: 200, height: 200 }}
            objectFit="cover"
          />
        ) : null,
      )}
    </View>
  );
}
```

## Props

```ts
import type { MediaStream } from '@fishjam-cloud/react-native-client';
import type { ViewStyle } from 'react-native';
import type { RTCPIPOptions } from '@fishjam-cloud/react-native-webrtc';

type RTCVideoViewProps = {
  mediaStream: MediaStream;   // REQUIRED — from useCamera/useScreenShare/peer.cameraTrack.stream/useLivestreamViewer.

  // Forwarded to the underlying react-native-webrtc RTCView:
  style?: ViewStyle;
  objectFit?: 'contain' | 'cover';
  mirror?: boolean;            // Flip horizontally — common for local camera preview
  zOrder?: number;             // Android only
  pip?: RTCPIPOptions;         // Declarative Picture-in-Picture — see picture-in-picture.md
  onDimensionsChange?: (dims: { width: number; height: number }) => void; // Standard way to compute aspect ratio
  // (...other RN-WebRTC RTCView props — consult the upstream types)
};
```

The SDK's wrapper takes `mediaStream` as a typed prop; internally it calls `mediaStream.toURL()` and forwards it as `streamURL` to the native component. You do **not** pass `streamURL` yourself.

> **Important:** the `MediaStream` type must come from `@fishjam-cloud/react-native-webrtc` (re-exported through `@fishjam-cloud/react-native-client` — `import type { MediaStream } from '@fishjam-cloud/react-native-client'`). The DOM `MediaStream` lacks `.toURL()` and the wrapper logs an error.

## Local camera preview tips

Mirror the local preview so it matches what the user expects from a mirror:

```tsx
<RTCView mediaStream={localPeer.cameraTrack.stream} mirror />
```

Don't mirror remote peers — you'd see their text appearing backwards.

## Performance: many tiles

Each `RTCView` is a native surface. Rendering 20+ at once can stress the GPU. Strategies:

- Request `Variant.LOW` for non-focused tiles via `usePeers().setReceivedTracksQuality(...)` (see `../react-client/simulcast-and-bandwidth.md`).
- Show only the speakers / spotlighted peers; collapse others to thumbnails or text.

## Audio

`RTCView` renders **video** only. Audio plays automatically once the SDK subscribes to a peer's audio track — no extra component needed. Local audio is muted by default in the SDK (no echo).

## RTCPIPView

Superset of `RTCView`. Same props, plus `pip.fallbackView` — the React subtree to render while the PIP window is detached from the host view. This is the distinguishing reason to use `RTCPIPView` over `RTCView` with a `pip` prop. See `picture-in-picture.md`.

```tsx
import { RTCPIPView } from '@fishjam-cloud/react-native-client';

<RTCPIPView
  mediaStream={presenter.cameraTrack.stream}
  style={{ flex: 1 }}
  pip={{ enabled: true, startAutomatically: true, fallbackView: <PausedPlaceholder /> }}
/>
```

When the user backgrounds the app, this view becomes the floating PIP window (assuming you enabled `supportsPictureInPicture` in the Expo plugin).

## Sources

- `packages/mobile-client/src/overrides/RTCView.tsx`
- `@fishjam-cloud/react-native-webrtc` (RTCView native components)
- [`fishjam-cloud/examples`](https://github.com/fishjam-cloud/examples) repo: `mobile-react-native/minimal-react-native/`
- [`fishjam-cloud/examples`](https://github.com/fishjam-cloud/examples) repo: `mobile-react-native/fishjam-chat/`
