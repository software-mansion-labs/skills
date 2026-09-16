# Drivers: who writes and reads the shared value

Questions 1, 9 and 10 of the walk.

## 1. Where is the value written?

Per frame or on the UI thread, so it stays on hooks: a scroll handler, a gesture `onUpdate`/`onChange`, gesture callbacks that run as worklets (every callback inline or marked `'worklet'` and no `.runOnJS(true)`; Gesture Handler 2 and 3 alike), `useAnimatedReaction`, `useFrameCallback`, a sensor hook, `useAnimatedKeyboard`, a `scheduleOnUI`/`runOnUI` body, a JS loop assigning every frame, and any shared value a library owns (a bottom sheet's `animatedIndex`, a carousel's progress, `useScrollOffset`): name the library in the reason.

On the JS thread, so the walk continues: `useEffect`, `onPress`/`onChange`, a timer, a network callback, a gesture `onEnd`/`onFinalize` with `.runOnJS(true)`, or no shared value at all because the hook body reads a prop or state with the `with*` inline, where the render is the driver:

```tsx
const style = useAnimatedStyle(() => ({ opacity: withTiming(visible ? 1 : 0, { duration: 200 }) }), [visible]);
```

A writer outside the requested scope: widen the scope, or note Needs approval.

Indirection: a `useDerivedValue`, or a `SharedValue` passed as a prop, is not a driver. Follow it to the shared value that is written and classify that write; the derived computation joins question 4. A derived value that anything besides the migrated hook consumes (a `useAnimatedReaction`, a gesture, a child through props) stays on hooks.

A custom hook that wraps `useAnimatedStyle` (`useFadeIn()`) is one site: classify the hook body, and the verdict must hold for every call site, which then receives plain style props.

Measured layout (`onLayout`, `measure` in an effect) stored in React state is a plain state driver: a transition to the measured target is fine. `measure` read inside the hook every frame is the per-frame case above.

## 9. Other writers

With a transition on a property, every change of that property animates. A site that also sets the value without a `with*` (`sv.value = 0` before a `withTiming(1)`, a mount jump, a conditional style) turned an instant jump into a glide: Needs approval, saying which write changes. A site whose writes are all instant never animated: plain state in the static style, no `transitionProperty`.

A property that animates on mount and is later retargeted by state is one transition: render the start value, flip the state in a mount `useEffect` (the second render starts the transition), and list the property only in `transitionProperty`, never in both `animationName` and `transitionProperty`.

## 10. Other readers

Removing the shared value breaks every other reader. A worklet that reads it (`useAnimatedReaction`, a gesture, `scrollTo`), a child that receives it as a prop, or JS that reads `sv.value`: Keep on hooks, or Needs approval proposing a state mirror for those readers. Another `useAnimatedStyle` reading the same looping value: see `references/transitions-and-animations.md`, shared phase.

## Example: state-driven fade

```tsx
// Before
const opacity = useSharedValue(0);
useEffect(() => { opacity.value = withTiming(visible ? 1 : 0, { duration: 200 }); }, [visible]);
const style = useAnimatedStyle(() => ({ opacity: opacity.value }));
```

```tsx
// After: the state that drove the effect drives the transition
const reduced = useReducedMotion();
<Animated.View
  style={[styles.box, {
    opacity: visible ? 1 : 0,
    transitionProperty: 'opacity',
    transitionDuration: reduced ? 1 : 200,
    transitionTimingFunction: 'ease-in-out',
  }]}
/>
```

The row records `inOut(quad)` to `'ease-in-out'` (0.012) and, when `visible` can flip back inside 200ms, the shortened return from question 8.
