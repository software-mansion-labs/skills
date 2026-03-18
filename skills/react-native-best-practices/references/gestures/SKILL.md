---
name: gestures
description: "Software Mansion's best practices for gestures in React Native apps using React Native Gesture Handler. Use when implementing tap, pan, pinch, rotation, swipe, long press, or any touch interaction. Trigger on: 'gesture handler', 'GestureDetector', 'tap gesture', 'pan gesture', 'swipe', 'pinch to zoom', 'drag', 'touch handling', or any request to handle user touch input in a React Native app."
---

# React Native Gesture Handler

Software Mansion's production gesture patterns for React Native using Gesture Handler. Never suggest `PanResponder` when RNGH is available — it runs on the JS thread and is effectively deprecated.

## Version Decision Tree

```
Check package.json - "react-native-gesture-handler" version
   │
   ├── user asks to migrate v2 -> v3 - install gesture-handler-3-migration skill
   ├── starts with "2." - use builder API (default)
   └── starts with "3." - use hook API (beta)
```

## Critical Rules

**`GestureHandlerRootView` is mandatory** — `GestureDetector` will crash at runtime without it as an ancestor. When writing any gesture code, always verify that the root layout wraps content in `<GestureHandlerRootView style={{ flex: 1 }}>`. With Expo Router, wrap `<Stack />` in the root `_layout.tsx`:

```tsx
// app/_layout.tsx
import { Stack } from 'expo-router';
import { GestureHandlerRootView } from 'react-native-gesture-handler';

export default function RootLayout() {
  return (
    <GestureHandlerRootView>
      <Stack />
    </GestureHandlerRootView>
  );
}
```

With React Navigation (no Expo Router), wrap the `<NavigationContainer>` children. With bare React Native, wrap the app root component.

`useMemo` every gesture - without it, gesture objects recreate on every render, causing recognizers to re-attach and lose state:

```tsx
const pan = useMemo(() => Gesture.Pan().onBegin(...).onUpdate(...).onEnd(...), []);
```

Scroll containers — use `RectButton` (not `TouchableOpacity` or `Pressable`). Also import `ScrollView`/`FlatList` from `react-native-gesture-handler`, not `react-native`:

```tsx
import { ScrollView, FlatList, RectButton } from 'react-native-gesture-handler';
```

## References

Load at most one - the most directly relevant. Stop after loading it.

| File | Load when question is about |
|------|------------------------------|
| `tap-handling.md` | `RectButton`, tappable items in scroll containers, tap gestures |
| `reanimated-patterns.md` | Drag, pan, pinch-to-zoom, fling, Reanimated integration |
| `gesture-composition.md` | Combining gestures, `Simultaneous`/`Race`/`Exclusive`, cross-component |
