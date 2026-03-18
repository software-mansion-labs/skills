# Gestures

Choosing the right gesture type, understanding the callback lifecycle, and handling threading for React Native Gesture Handler.

For tap-specific patterns, see **`tap-handling.md`**.
For continuous gesture patterns (pan, pinch, rotation), see **`continuous-gestures.md`**.
For combining gestures, see **`gesture-composition.md`**.
For pre-built components (Swipeable, Drawer), see **`swipeable-and-drawer.md`**.

---

## Decision Tree

Pick the gesture type based on the interaction you need.

```
What touch interaction does the user perform?
│
├── Tap / press on a UI element?
│   ├── Inside a ScrollView/FlatList → RectButton (v2) / LegacyRectButton (v3)
│   ├── Simple press outside scroll containers → RNGH Pressable
│   ├── Needs UI-thread animation on press → GestureDetector + Tap gesture
│   └── Double-tap / multi-tap → Tap gesture with numberOfTaps
│
├── Drag / pan movement?
│   ├── Simple drag-and-drop → Pan gesture
│   ├── Pan with momentum/fling → Pan gesture + withDecay
│   └── Pan inside ScrollView → Pan + Native gesture (simultaneousWith)
│
├── Pinch / zoom?
│   └── Pinch gesture (combine with Pan for photo viewer)
│
├── Rotation?
│   └── Rotation gesture (combine with Pinch for photo viewer)
│
├── Quick flick in a direction?
│   └── Fling gesture (configure direction)
│
├── Long press?
│   ├── Context menu / haptic → LongPress gesture
│   └── Long press then drag → Pan gesture with activateAfterLongPress
│
├── Mouse/stylus hover?
│   └── Hover gesture (iOS/Android/Web)
│
├── Swipeable list item?
│   └── ReanimatedSwipeable component   → see swipeable-and-drawer.md
│
├── Side drawer / navigation?
│   └── ReanimatedDrawerLayout component → see swipeable-and-drawer.md
│
└── Completely custom recognition logic?
    └── Manual gesture + GestureStateManager
```

---

## Callback Lifecycle

All gestures follow the same state-driven callback flow. Guaranteed pairs: if `onBegin` fires, `onFinalize` will fire; if `onActivate` fires, `onDeactivate` will fire.

```
Touch starts
  │
  ├── onBegin ─── gesture begins receiving touches
  │
  ├── onActivate ─── activation criteria met (finger moved enough, held long enough, etc.)
  │   │
  │   ├── onUpdate ─── (continuous gestures only) called on each pointer move
  │   │
  │   └── onDeactivate(didSucceed) ─── gesture ends
  │       ├── didSucceed = true → user completed the gesture normally (END)
  │       └── didSucceed = false → gesture was cancelled or failed
  │
  └── onFinalize(didSucceed) ─── handler stops recognizing (always called)
```

For touch-level events (tracking individual fingers), use `onTouchesDown`, `onTouchesMove`, `onTouchesUp`, `onTouchesCancel`. Track touches by `id`, not array index -- the order in `changedTouches`/`allTouches` can change during a gesture.

### v2 Builder API callback names

| Lifecycle | v2 name | v3 name |
|-----------|---------|---------|
| Begin | `.onBegin()` | `onBegin` |
| Activate | `.onStart()` | `onActivate` |
| Update | `.onUpdate()` / `.onChange()` | `onUpdate` (includes `change*` properties) |
| Deactivate | `.onEnd()` | `onDeactivate` |
| Finalize | `.onFinalize()` | `onFinalize` |

---

## Threading

When `react-native-reanimated` is installed, gesture callbacks run on the UI thread by default. This enables 60/120fps gesture-driven animations without JS thread involvement.

**Auto-workletization**: Callbacks passed directly (inline) to gesture hooks/builders are auto-workletized by the Babel plugin. Callbacks passed as variable references or wrapped in `useCallback`/`useMemo` require an explicit `'worklet';` directive:

```tsx
// Auto-workletized (inline)
const gesture = usePanGesture({
  onUpdate: (e) => { offset.value = e.translationX; },
});

// Needs explicit 'worklet' (variable reference)
const handleUpdate = (e: GestureEvent) => {
  'worklet';
  offset.value = e.translationX;
};
const gesture = usePanGesture({ onUpdate: handleUpdate });
```

**Calling JS-thread functions from gesture callbacks**: Use `scheduleOnRN` from `react-native-worklets` to call React state setters or navigation from a UI-thread gesture callback. `runOnJS` is deprecated in Reanimated 4:

```tsx
import { scheduleOnRN } from 'react-native-worklets';

// v3
const pan = usePanGesture({
  onDeactivate: () => { scheduleOnRN(setPosition, offset.value); },
});

// v2
const pan = useMemo(() =>
  Gesture.Pan().onEnd(() => { scheduleOnRN(setPosition, offset.value); }),
[]);
```

**Disabling Reanimated**: Set `disableReanimated: true` (v3) to run all callbacks on the JS thread without Reanimated overhead. Useful for gestures that only trigger JS-side logic and do not animate.

**`runOnJS` property** (v3 only): Dynamically control per-gesture whether callbacks run on the JS or UI thread. Accepts a `SharedValue<boolean>` for runtime switching.

---

## SharedValue in Gesture Config (v3)

v3 gesture hooks accept `SharedValue` objects for configuration properties, allowing gesture parameters to change without re-renders:

```tsx
const numberOfTaps = useSharedValue(2);

const tap = useTapGesture({
  numberOfTaps,
  onDeactivate: () => {
    numberOfTaps.value += 1; // Next time requires one more tap
  },
});
```

Works for `enabled`, `numberOfTaps`, `minDistance`, `hitSlop`, and most configuration properties.

---

## GestureStateManager

For manual control of gesture activation. Use when automatic activation criteria do not match your needs.

```tsx
import { GestureStateManager } from 'react-native-gesture-handler';

// v3 -- skip long press wait, activate immediately on touch
const longPress = useLongPressGesture({
  onTouchesDown: (e) => {
    GestureStateManager.activate(e.handlerTag);
  },
  onActivate: () => { /* fires immediately */ },
});

// v3 -- activate one gesture from another
const pan = usePanGesture({ manualActivation: true, onActivate: () => { ... } });
const longPress = useLongPressGesture({
  onActivate: () => { GestureStateManager.activate(pan.handlerTag); },
});
```

Methods: `begin(tag)`, `activate(tag)`, `deactivate(tag)`, `fail(tag)`.

`manualActivation: true` prevents the gesture from activating on its own -- you must call `GestureStateManager.activate()` explicitly.

v2 equivalent: the `stateManager` object was passed as a second argument to touch callbacks. In v3, use the global `GestureStateManager` with `e.handlerTag`.

---

## API Reference

For up-to-date API signatures, configuration options, and event data, webfetch these pages:

**Gestures (v3 hook API):**

| Gesture | Documentation |
|---------|--------------|
| Pan | [usePanGesture](https://docs.swmansion.com/react-native-gesture-handler/docs/gestures/use-pan-gesture) |
| Tap | [useTapGesture](https://docs.swmansion.com/react-native-gesture-handler/docs/gestures/use-tap-gesture) |
| Pinch | [usePinchGesture](https://docs.swmansion.com/react-native-gesture-handler/docs/gestures/use-pinch-gesture) |
| Rotation | [useRotationGesture](https://docs.swmansion.com/react-native-gesture-handler/docs/gestures/use-rotation-gesture) |
| Fling | [useFlingGesture](https://docs.swmansion.com/react-native-gesture-handler/docs/gestures/use-fling-gesture) |
| Long Press | [useLongPressGesture](https://docs.swmansion.com/react-native-gesture-handler/docs/gestures/use-long-press-gesture) |
| Hover | [useHoverGesture](https://docs.swmansion.com/react-native-gesture-handler/docs/gestures/use-hover-gesture) |
| Native | [useNativeGesture](https://docs.swmansion.com/react-native-gesture-handler/docs/gestures/use-native-gesture) |
| Manual | [useManualGesture](https://docs.swmansion.com/react-native-gesture-handler/docs/gestures/use-manual-gesture) |

**Gestures (v2 builder API):**

| Gesture | Documentation |
|---------|--------------|
| Pan | [Gesture.Pan](https://docs.swmansion.com/react-native-gesture-handler/docs/legacy-gestures/pan-gesture) |
| Tap | [Gesture.Tap](https://docs.swmansion.com/react-native-gesture-handler/docs/legacy-gestures/tap-gesture) |
| Pinch | [Gesture.Pinch](https://docs.swmansion.com/react-native-gesture-handler/docs/legacy-gestures/pinch-gesture) |
| Rotation | [Gesture.Rotation](https://docs.swmansion.com/react-native-gesture-handler/docs/legacy-gestures/rotation-gesture) |
| Fling | [Gesture.Fling](https://docs.swmansion.com/react-native-gesture-handler/docs/legacy-gestures/fling-gesture) |
| Long Press | [Gesture.LongPress](https://docs.swmansion.com/react-native-gesture-handler/docs/legacy-gestures/long-press-gesture) |
| Hover | [Gesture.Hover](https://docs.swmansion.com/react-native-gesture-handler/docs/legacy-gestures/hover-gesture) |
| Native | [Gesture.Native](https://docs.swmansion.com/react-native-gesture-handler/docs/legacy-gestures/native-gesture) |
| Manual | [Gesture.Manual](https://docs.swmansion.com/react-native-gesture-handler/docs/legacy-gestures/manual-gesture) |

**Fundamentals:**

- [GestureHandlerRootView](https://docs.swmansion.com/react-native-gesture-handler/docs/fundamentals/root-view)
- [Gesture Detectors](https://docs.swmansion.com/react-native-gesture-handler/docs/fundamentals/gesture-detectors)
- [Callbacks & Events](https://docs.swmansion.com/react-native-gesture-handler/docs/fundamentals/callbacks-events)
- [Gesture Composition](https://docs.swmansion.com/react-native-gesture-handler/docs/fundamentals/gesture-composition)
- [State Manager](https://docs.swmansion.com/react-native-gesture-handler/docs/fundamentals/state-manager)
- [Reanimated Integration](https://docs.swmansion.com/react-native-gesture-handler/docs/fundamentals/reanimated-interactions)
- [State Machine](https://docs.swmansion.com/react-native-gesture-handler/docs/under-the-hood/state)
