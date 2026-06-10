# Audio Output Routing

Mobile devices have multiple audio outputs — earpiece, speaker, wired headset, Bluetooth headset — and the OS routes audio to one of them. The SDK exposes a hook to read and override the routing.

> No web equivalent in the React client — programmatic audio-output routing on the web relies on browser-specific APIs (e.g. `HTMLMediaElement.setSinkId`) with uneven cross-browser support, so the SDK doesn't expose an equivalent hook. This is mobile-only.

## API

```tsx
import {
  useAudioOutput,
  AudioDeviceType,
} from '@fishjam-cloud/react-native-client';
import type {
  AudioDevice,
  AudioOutputChangedInfo,
  UseAudioOutputResult,
} from '@fishjam-cloud/react-native-client';

const result: UseAudioOutputResult = useAudioOutput();
```

The SDK re-exports `useAudioOutput` and `AudioDeviceType` directly from `@fishjam-cloud/react-native-webrtc`, plus the types `AudioDevice`, `AudioOutputChangedInfo`, and `UseAudioOutputResult`.

> **Exact field names** on `UseAudioOutputResult` (list of devices, current selection, selector function) and `AudioDevice` come from the upstream `@fishjam-cloud/react-native-webrtc` types. Use TypeScript intellisense or check those types directly when wiring this up — destructure the hook result and let the compiler guide you.

`AudioDeviceType` enumerates the broad categories: `earpiece`, `speaker`, `bluetooth`, `wiredHeadset`, `usb`, `hdmi`, `airplay`, `carAudio`, `hearingAid`, `lineOut`, `unknown`. Use the enum members (e.g. `AudioDeviceType.speaker`, `AudioDeviceType.bluetooth`) — they're lowercase to match the upstream package — rather than string literals so renames don't silently break your code.

## Default routing

If you never call the selector, the OS picks a sensible default:

- Headset connected → headset.
- Bluetooth connected and the user prefers it → Bluetooth.
- Otherwise → earpiece for voice-only / speaker for video calls (heuristic-based).

Most apps don't need to override. Override when:

- The user explicitly toggles "speakerphone" on or off.
- You want video calls to start in speaker mode by default.
- An external audio router demands a specific output (Bluetooth deck, USB audio interface, etc.).

## API shape

`useAudioOutput()` returns:

```ts
type UseAudioOutputResult = {
  currentAudioOutput: AudioDevice | null;
  availableAudioOutputs: AudioDevice[];
  ios: {
    showAudioRoutePicker(): void;
    overrideAudioOutput(output: 'speaker' | 'none'): Promise<void>;
  };
  android: {
    selectAudioOutput(deviceId: string): Promise<void>;
    // …additional Android-only controls
  };
};
```

The `.ios` / `.android` namespaces throw if you call them on the wrong platform — guard with `Platform.OS` first.

## Forcing speaker mode

On iOS, override the route to the built-in speaker:

```tsx
import { useAudioOutput } from '@fishjam-cloud/react-native-client';
import { Platform } from 'react-native';
import { useEffect } from 'react';

function ForceSpeaker() {
  const audio = useAudioOutput();

  useEffect(() => {
    if (Platform.OS === 'ios') {
      audio.ios.overrideAudioOutput('speaker').catch(console.error);
    }
  }, [audio]);

  return null;
}
```

On Android, pick a specific output device by its id:

```tsx
const speaker = audio.availableAudioOutputs.find((d) => d.type === AudioDeviceType.speaker);
if (speaker && Platform.OS === 'android') {
  await audio.android.selectAudioOutput(speaker.id);
}
```

## Reactive Bluetooth in/out

The hook's device list updates when the user connects or disconnects a Bluetooth or wired headset — re-render off the hook value rather than subscribing to OS audio events separately. When a previously-selected device disconnects, the OS picks a fallback automatically and the hook reports the new selection.

## Pairing with CallKit / foreground service

- During a CallKit session, iOS respects audio routes the same way as during a phone call (e.g. Bluetooth headset, AirPods).
- During an Android foreground-service-backed session, the SDK's audio focus declarations cooperate with the OS routing.

In both cases the selector works as expected — no special handling needed.

## Sources

- `packages/mobile-client/src/index.ts` (re-exports `useAudioOutput`, `AudioDeviceType`, type aliases)
- `@fishjam-cloud/react-native-webrtc` → `lib/typescript/useAudioOutput.d.ts` (`UseAudioOutputResult`) and `lib/typescript/audioOutputManager/AudioOutputManager.d.ts` (`AudioDevice`, `AudioDeviceType`, `AudioOutputChangedInfo`, `ios.showAudioRoutePicker`, `ios.overrideAudioOutput('speaker' | 'none')`, `android.selectAudioOutput(deviceId)`)
