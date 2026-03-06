# Tap Handling

## Decision

| What to use | When to use |
|------|------------------------------|
| RectButton (icons: BorderlessButton) | Inside ScrollView / FlatList |
| Gesture.Tap() + useMemo | Custom animation or multi-tap |
| Pressable | Other cases |

## RectButton in Scroll Containers

Import both `RectButton` and the scroll container from `react-native-gesture-handler` - mixing RNGH buttons with React Native's `ScrollView`/`FlatList` causes gesture conflicts:

```tsx
import { FlatList, RectButton, BorderlessButton } from 'react-native-gesture-handler';

<RectButton onPress={handlePress} style={styles.row}>
  <Text>{item.title}</Text>
  <BorderlessButton onPress={handleDelete}>
    <DeleteIcon />
  </BorderlessButton>
</RectButton>
```

## Gesture.Tap() - Custom Feedback / Multi-tap

```tsx
const tap = useMemo(() =>
  Gesture.Tap()
    .onBegin(() => { scale.value = withTiming(0.95); })
    .onFinalize(() => { scale.value = withTiming(1); }),
[]);
```

Double-tap with single-tap fallback - `Gesture.Exclusive` gives doubleTap priority:

```tsx
const doubleTap = useMemo(() => Gesture.Tap().numberOfTaps(2).onEnd(() => { runOnJS(onDouble)(); }), []);
const singleTap = useMemo(() => Gesture.Tap().numberOfTaps(1).onEnd(() => { runOnJS(onSingle)(); }), []);
const composed = useMemo(() => Gesture.Exclusive(doubleTap, singleTap), [doubleTap, singleTap]);
```