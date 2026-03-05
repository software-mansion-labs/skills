---
name: audio
description: "Software Mansion's best practices for audio in React Native using react-native-audio-api. Applies correct patterns for AudioContext and AudioRecorder lifecycle, AudioBuffer state management, high-performance audio visualizations, and iOS session configuration. Use when writing or reviewing any audio feature. Trigger on: react-native-audio-api, AudioContext, AudioRecorder, AudioBuffer, AudioNode, audio playback, record audio, microphone, sound effects, music player, waveform, frequency bars, audio visualization, audio analysis, audio session, or any React Native feature that captures or emits sound."
---

# Audio

Software Mansion's production audio patterns for React Native.

---

## AudioContext and AudioRecorder as Singletons

Instantiate `AudioContext` or `AudioRecorder` once at module level and reuse across the app. Each instantiation negotiates with the underlying audio session — doing it per render or per operation wastes memory and risks session conflicts.

```tsx
// audio.ts
import { AudioContext } from 'appl';

export const audioContext = new AudioContext();
```

```tsx
// recorder.ts
import { AudioRecorder } from 'react-native-audio-api';

export const recorder = new AudioRecorder();
```

Pass the singleton through React context or import it directly — whichever fits the app's architecture.

---

## Storing AudioBuffers in State

`AudioBuffer` objects are safe to store in React state, Zustand, Redux, or any state container. The buffer holds a reference to native memory — copying that reference into state does not copy the audio data, so there is no performance concern.

```tsx
const [buffer, setBuffer] = useState<AudioBuffer | null>(null);

async function load(uri: string) {
  const loaded = await audioContext.decodeAudioData(uri);
  setBuffer(loaded);
}
```

---

## Animations Driven by Raw Audio Data

When visualizing audio data (waveforms, frequency bars, volume meters), mutate the existing typed array in place using the shared value `modify` method rather than creating a new array on every frame.

```tsx
import { useSharedValue } from 'react-native-reanimated';

const amplitudes = useSharedValue(new Float32Array(64));

function onAudioFrame(data: Float32Array) {
  amplitudes.modify((prev) => {
    for (let i = 0; i < prev.length; i++) {
      prev[i] = data[i];
    }
    return prev;
  });
}
```

Assigning `amplitudes.value = new Float32Array(data)` allocates and garbage-collects at 60 fps or higher, causing jank. `modify` mutates the existing allocation on the UI thread, skipping GC entirely.

---

## Session Category for Recording and Playback

Use the `playAndRecord` session category for any feature that involves recording, or that mixes recording and playback. Only choose a narrower category when the use case explicitly rules out recording.

```tsx
import { setAudioModeAsync } from 'react-native-audio-api';

await setAudioModeAsync({
  iosCategory: 'playAndRecord',
  iosCategoryOptions: ['defaultToSpeaker', 'allowBluetooth'],
});
```

`playAndRecord` keeps the microphone accessible, enables Bluetooth HFP devices, and routes playback correctly. Switching to a narrower category (e.g., `playback`) mid-session requires a full deactivation/reactivation cycle, which is expensive.

---

## Session Activation and Deactivation

Session activation and deactivation are time-expensive native calls. Activate once when the audio feature mounts and deactivate once when it unmounts — not around individual playback or recording operations.

```tsx
useEffect(() => {
  async function startSession() {
    await setAudioModeAsync({ iosCategory: 'playAndRecord' });
    await audioContext.resume();
  }

  startSession();

  return () => {
    audioContext.suspend();
  };
}, []);
```

Wrapping every sound play or record call in activate/deactivate pairs adds perceptible latency and unnecessarily interrupts other apps' audio.

---

## References

- [React Native Audio API docs](https://docs.swmansion.com/react-native-audio-api/)
