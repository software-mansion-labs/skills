# Migrating from expo-haptics to Pulsar

Use this reference only for a React Native or Expo migration from `expo-haptics` to
`react-native-pulsar`. Confirm package versions and callable symbols against resolved
project artifacts before changing code.

## System-feedback mapping

| expo-haptics                                                  | Pulsar system preset                   | Migration note                                     |
| ------------------------------------------------------------- | -------------------------------------- | -------------------------------------------------- |
| `Haptics.impactAsync(ImpactFeedbackStyle.Light)`              | `Presets.System.impactLight()`         | Direct system preset                               |
| `Haptics.impactAsync(ImpactFeedbackStyle.Medium)`             | `Presets.System.impactMedium()`        | Direct system preset                               |
| `Haptics.impactAsync(ImpactFeedbackStyle.Heavy)`              | `Presets.System.impactHeavy()`         | Direct system preset                               |
| `Haptics.impactAsync(ImpactFeedbackStyle.Soft)`               | `Presets.System.impactSoft()`          | expo-haptics Android waveform matches Light        |
| `Haptics.impactAsync(ImpactFeedbackStyle.Rigid)`              | `Presets.System.impactRigid()`         | expo-haptics Android waveform matches Medium       |
| `Haptics.notificationAsync(NotificationFeedbackType.Success)` | `Presets.System.notificationSuccess()` | Direct system preset                               |
| `Haptics.notificationAsync(NotificationFeedbackType.Warning)` | `Presets.System.notificationWarning()` | Direct system preset                               |
| `Haptics.notificationAsync(NotificationFeedbackType.Error)`   | `Presets.System.notificationError()`   | Direct system preset                               |
| `Haptics.selectionAsync()`                                    | `Presets.System.selection()`           | expo-haptics Android waveform matches Light impact |

Pulsar system presets use Android system effects, respecting device haptic tuning. Do not
promise identical Android waveform output: expo-haptics uses raw waveforms and collapses
several styles. If exact Android effect selection matters, verify matching-version
`Presets.System.Android.*` exports before using platform-specific presets.

## Code shape

```tsx
// Before
import * as Haptics from "expo-haptics";

function handleSuccess() {
  Haptics.notificationAsync(Haptics.NotificationFeedbackType.Success);
}

// After
import { Presets } from "react-native-pulsar";

function handleSuccess() {
  Presets.System.notificationSuccess();
}
```

Pulsar functions are synchronous and worklet-compatible in upstream guidance: do not
carry `await` from `expo-haptics`. After preserving system semantics, consider a named
Pulsar preset only when richer feedback fits event meaning. Verify support, setup,
cleanup, and physical-device feel before reporting migration complete.
