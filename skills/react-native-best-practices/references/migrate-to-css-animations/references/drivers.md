# Drivers: where the values come from, who else touches them

Questions 1, 9 and 10 of the walk.

## 1. Where do the target values come from?

A target is the value a `with*` animates to: `withTiming(1)` animates to 1. A driver is the code or input that changes the target. CSS needs the target to come from React state or props. A render that changes the style starts the CSS transition or animation.

The JS thread runs React and ordinary event handlers. The UI thread handles frame-by-frame interface updates. A worklet is a function prepared to run in the UI runtime.

Classify each target by what sets it.

Keep on shared values for continuous input:

- A scroll offset.
- A moving finger: a gesture's `onUpdate` writes the value, or `onChange` in Gesture Handler 2 does.
- A sensor or `useAnimatedKeyboard`.
- A `useFrameCallback` that computes new targets every frame.
- Any shared value owned by a library, such as a bottom sheet's `animatedIndex`, a carousel's progress, or `useScrollOffset`. Name the library in the reason.

A UI-thread gesture callback may fire once per interaction: `onBegin`, `onStart`, `onEnd`, or `onFinalize`. In the Gesture Handler 3 hooks, `onStart` is named `onActivate` and `onEnd` is named `onDeactivate`.

If such a callback writes a target for a style property that nothing drives every frame, CSS needs that target in React state. A UI-thread callback cannot set React state directly.

Treat all entries in `transform` as one style property:

```tsx
const style = useAnimatedStyle(() => ({
  transform: [
    { translateX: x.value }, // onUpdate writes x every frame
    { scale: scale.value },  // onEnd animates scale on release
  ],
}));
```

Here, `translateX` keeps the whole `transform` array in `useAnimatedStyle`. A CSS `transform` and the hook's `transform` on the same view replace each other. Their entries do not merge. Keep `scale` on shared values too.

For targets that can become state setters, check the gesture's other callbacks:

- If none must stay on the UI thread, set `runOnJS: true` on the gesture (`.runOnJS(true)` on a builder). Call the state setters directly. This is the simpler form. No callback may write a shared value the walk keeps or call `scrollTo` or another UI-only API. An `onUpdate` that only reads can move.
- If one must stay, such as an `onUpdate` tracking the finger, keep the gesture as it is. Call `scheduleOnRN(setTarget, value)` only in the once-per-interaction callback (`../animations/animations.md`, Threading). Setting `runOnJS: true` would move that `onUpdate` to the JS thread too.

Each gesture in a composed gesture has its own `runOnJS`. Decide per gesture. Both forms are Needs approval: the animation now starts after a React render instead of on the next UI frame.

A press pair belongs to question 8, which names its CSS form. This means `onPressIn`/`onPressOut`, or a Tap or LongPress gesture's `onBegin`/`onFinalize`, writing the pressed and rest values of a property on the component that owns the style. A Pan or Pinch activate/deactivate pair is not a press. Classify it here.

Gesture callbacks run the same way in Gesture Handler 2 and 3. They run on the UI thread when the callbacks are worklets and `runOnJS` is neither literal `true` nor a shared value holding `true`. Otherwise, use the JS thread rules below.

Callback form matters:

- A builder (`Gesture.Pan().onStart(...)`) needs every callback to be a worklet. One plain callback moves all callbacks to the JS thread.
- A Gesture Handler 3 hook (`usePanGesture({ onActivate, onUpdate, onDeactivate })`) rejects a mix of plain callbacks and worklets. All its callbacks run on one thread.
- Inline builder callbacks are worklets on every Reanimated 4.x.
- Inline Gesture Handler 3 hook callbacks are worklets from Reanimated 4.2.0. Before that, they are plain functions unless marked `'worklet'`. A hook with those plain callbacks runs on the JS thread.
- A callback defined elsewhere, such as an import or prop, is a worklet only when marked `'worklet'`.

Gesture Handler 3 hooks have no `onChange`; its fields are in `onUpdate`. Builders keep their callback names in both versions.

Legacy handler components (`<PanGestureHandler onGestureEvent>`) call plain JS functions on Reanimated 4, which removed `useAnimatedGestureHandler`. Use the JS thread rules below.

A `runOnUI`/`scheduleOnUI` body is how a write reaches the UI runtime, not its source. Classify the target it writes:

- A literal, state or prop from the calling handler or effect: use the JS thread rules below. Set state where `runOnUI` was called.
- Another shared value: follow it back to its first write (Chains, below).
- A body that also does its own UI-thread work (`scrollTo`, a gesture state), or runs every frame: Keep on shared values.

Targets from the JS thread let the walk continue. These include state, props, `useEffect`, `onPress` and other handlers, timers, network callbacks, and gesture callbacks that run on the JS thread.

This includes a `with*` inside `useAnimatedStyle` whose target is a prop or state:

```tsx
const style = useAnimatedStyle(() => ({ opacity: withTiming(visible ? 1 : 0, { duration: 200 }) }), [visible]);
```

The hook body runs on the UI thread. The target changes only when React renders a new prop or state value. Here, `visible` is the source, and the site needs no shared value.

If a writer is outside the requested scope, widen the scope or note Needs approval.

Chains: a `useDerivedValue`, a `useAnimatedReaction`, or a `SharedValue` passed as a prop is not a source. Follow the chain back to where the value is first written.

When that write is on the JS thread, collapse the whole chain:

- Set state where the shared value was written.
- Remove the derived value and the reaction.
- Let question 4 judge the computation the chain performed.

Keep on shared values if a reaction responds to continuous input or does its own UI-thread work (`scrollTo`, a gesture state). Remove a reaction whose only job is to forward a JS-written value. If you cannot tell what a reaction does for other code, note Needs approval and ask.

A custom hook wrapping `useAnimatedStyle`, such as `useFadeIn()`, is one site. Classify its body. The verdict must hold for every call site, which then receives plain style props.

A measured value from `onLayout`, or from `measure` in an effect or handler, is a JS-thread target. This holds whether it was stored in state, written into a shared value, or passed directly to `withTiming`. Hold it in state and transition to it. Reading `measure` inside the hook every frame is continuous input.

A shared value used directly in `style` or as a prop has no hook: `style={{ opacity: sv }}` or `<AnimatedCircle r={r} />`. Classify its writes the same way.

## 9. Other writers

Every change to a transitioned property animates. A write that previously changed it instantly would now animate too. Check for:

- `sv.value = x` without a `with*`, such as a jump in a handler or on mount.
- A branch in the hook body or a style-array entry that supplies a different value based on a condition other than the animated driver. For example, `opacity: disabled ? 0.5 : opacity.value` can change opacity when `disabled` changes.

To make a change jump in CSS:

1. Leave the property out of `transitionProperty` in the render that sets its new value.
2. Put it back in the next render, when the value no longer changes, so nothing animates.

Set a state flag with the new value and clear it in an effect:

```tsx
const [opacity, setOpacity] = useState(1);
const [jumping, setJumping] = useState(false);
const jumpTo = (value: number) => { setOpacity(value); setJumping(true); };
useEffect(() => { if (jumping) setJumping(false); }, [jumping]);
// on the element, with transform transitioning too:
// { opacity, transitionProperty: jumping ? 'transform' : ['opacity', 'transform'], transitionDuration: 200 }
```

Leave out only the property that jumps. If it is the only transitioned property, `'none'` removes the whole transition config. Adding that config back to a mounted view animates from the previous value only from 4.6.0. Earlier versions animate from the property's default, so such a site stays on shared values on those versions.

Propose this recipe as Needs approval. Name the write it covers.

A reset immediately followed by a `with*` is a replay, not a jump: question 6, restart.

If all writes are instant, the site never animated. Use plain state in the static style, with no `transitionProperty`.

## 10. Other readers

Removing the shared value breaks every other reader. Classify each reader by what it consumes:

- Another hook in scope, in this component or a child receiving the value through a prop: walk that hook as its own site. Question 1's Chains rule collapses both to the same state. For a looping value read by several hooks, see `references/transitions-and-animations.md`, shared phase.
- A JS handler that reads `sv.value` only at rest and cannot run during the animation: note Needs approval and propose a state mirror.
- A reader that consumes the value during the animation: Keep on shared values. Name the reader as the reason. This includes a worklet (`useAnimatedReaction` not removed by Chains, a gesture callback, `scrollTo`), a child or library component outside scope rendering the value from a prop, or a handler that can run during the animation.

## Example: state-driven fade

```tsx
// Before
const opacity = useSharedValue(0);
useEffect(() => { opacity.value = withTiming(visible ? 1 : 0, { duration: 200 }); }, [visible]);
const style = useAnimatedStyle(() => ({ opacity: opacity.value }));
```

```tsx
// After: the state that drove the effect drives the transition
<Animated.View
  style={[styles.box, {
    opacity: visible ? 1 : 0,
    transitionProperty: 'opacity',
    transitionDuration: 200,
    transitionTimingFunction: 'ease-in-out',
  }]}
/>
```

The row records `inOut(quad)` to `'ease-in-out'` (0.012). When `visible` can flip back inside 200ms, it also records the shortened return from question 8.

Here, `visible` starts `false`. If it can start `true`, the effect faded the element in on mount. See `references/transitions-and-animations.md`, section 'Transition or animation': render `0`, then flip the state in a mount effect.

With reduced motion kept, `transitionDuration` becomes `reduced ? 1 : 200`.

## Example: a reaction that only forwards a JS write

```tsx
// Before
const step = useSharedValue(0);
const width = useSharedValue(0);
useAnimatedReaction(() => step.value, (s) => { width.value = withTiming(s * 80, { duration: 250 }); });
const style = useAnimatedStyle(() => ({ width: width.value }));
const next = () => { step.value = step.value + 1; };
```

```tsx
// After: the handler that wrote step now sets state; the reaction and both shared values go
const [step, setStep] = useState(0);
const next = () => setStep((s) => s + 1);
<Animated.View style={[styles.bar, { width: step * 80, transitionProperty: 'width', transitionDuration: 250, transitionTimingFunction: 'ease-in-out' }]} />
```
