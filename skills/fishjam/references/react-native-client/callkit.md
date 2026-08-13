# CallKit (iOS)

CallKit lets your Fishjam sessions behave as **native phone calls** on iOS — full-screen ringer, lock-screen controls, system call history. No equivalent on Android (the SDK exposes the same hooks but they're no-ops there; build a separate notification UI for Android).

> **Plugin prerequisite:** `ios.enableVoIPBackgroundMode: true`. The plugin also adds the `voip` background mode to `Info.plist`.

## API

The SDK re-exports three hooks and two types from `@fishjam-cloud/react-native-webrtc`:

```ts
import {
  useCallKit,
  useCallKitEvent,
  useCallKitService,
  type CallKitAction,
  type CallKitConfig,
} from '@fishjam-cloud/react-native-client';

// useCallKit() returns:
//   { startCallKitSession(config), endCallKitSession(), getCallKitSessionStatus() }
// useCallKitService(config) starts a session on mount and ends it on unmount.
// useCallKitEvent('<action>', cb) subscribes to a single system action.
```

`CallKitConfig = { displayName: string; isVideo: boolean }`.

`CallKitAction` keys (from `@fishjam-cloud/react-native-webrtc/src/CallKit.ts`):

```ts
type CallKitAction = {
  started?: undefined;   // session began
  ended?: undefined;     // red button / system-initiated hangup
  failed?: string;       // error message
  muted?: boolean;       // true = mute, false = unmute
  held?: boolean;        // true = hold,  false = resume
};
```

## Outbound call pattern (matches the fishjam-chat example)

`useCallKitService` is the recommended pattern — it handles start-on-mount and end-on-unmount automatically.

```tsx
import {
  useCallKitEvent,
  useCallKitService,
  useConnection,
  useMicrophone,
} from '@fishjam-cloud/react-native-client';

function RoomScreen({ peerToken, displayName }) {
  const { joinRoom, leaveRoom } = useConnection();
  const { startMicrophone, stopMicrophone } = useMicrophone();

  useCallKitService({ displayName, isVideo: true });

  useCallKitEvent('ended', () => leaveRoom());
  useCallKitEvent('muted', (isMuted) => {
    if (isMuted === true) stopMicrophone();
    else if (isMuted === false) startMicrophone();
  });
  useCallKitEvent('held', (isHeld) => {
    if (isHeld === true) stopMicrophone();
    else if (isHeld === false) startMicrophone();
  });

  // ...joinRoom({ peerToken }) when you're ready to connect.
}
```

If you need manual control instead of the service hook, call `useCallKit()` and use `startCallKitSession({ displayName, isVideo })` / `endCallKitSession()` (both `Promise<void>`).

**Always handle the `ended` event.** Users expect the red button to actually hang up; ignoring it leaves the system UI out of sync with your app state.

## Inbound calls (with PushKit)

Inbound CallKit calls require a separate VoIP push notification via PushKit, which is not provided by the SDK. The flow:

1. Your backend sends a VoIP push to the user's device with call info.
2. The OS wakes your PushKit handler.
3. Your handler reports the incoming call to CallKit so the system rings.
4. On accept, fetch a peer token from your backend and `joinRoom({ peerToken })`.

`useCallKitService(config)` only takes the `CallKitConfig` shown above (`{ displayName, isVideo }`) — it does not accept lifecycle callbacks in its config. Lifecycle handling is split: subscribe to each `CallKitAction` (`started` / `ended` / `failed` / `muted` / `held`) with one `useCallKitEvent('<action>', cb)` call apiece, exactly as in the outbound-call snippet above. Reporting the *incoming* call to the system still happens through your PushKit handler outside this SDK.

## Plugin flag — `enableVoIPBackgroundMode`

```json
{ "ios": { "enableVoIPBackgroundMode": true } }
```

This adds `voip` to `UIBackgroundModes`. Without it:

- VoIP pushes won't wake your app.
- The peer connection won't survive backgrounding during a call.

## Android

The three hooks are no-ops on Android. For a "call in progress" notification on Android, use `useForegroundService` with appropriate `notificationTitle` / `notificationContent` — see `foreground-service.md`.

## Sources

- `packages/mobile-client/src/index.ts` (re-exports `CallKitAction`, `CallKitConfig`)
- `packages/mobile-client/src/overrides/hooks.ts` (`useCallKit`, `useCallKitEvent<T>`, `useCallKitService` wrappers)
- `@fishjam-cloud/react-native-webrtc` → `lib/typescript/useCallKit.d.ts` (`useCallKit`, `useCallKitService`, `useCallKitEvent`, `UseCallKitResult`) and `lib/typescript/CallKit.d.ts` (`CallKitConfig = { displayName, isVideo }`, `CallKitAction = { started?, ended?, failed?, muted?, held? }`)
- Cross-reference: `foreground-service.md`, `picture-in-picture.md`
