# Drivers: where the values come from, who else touches them

Questions 1, 9 and 10 of the walk.

## 1. Where do the target values come from?

A `with*` call produces the frames between two targets; CSS replaces exactly that part. What decides the verdict is where the targets come from.

Continuous input, so the site stays on shared values: a scroll offset, a moving finger (a gesture whose `onUpdate` drives the value, or `onChange` in Gesture Handler 2), a sensor, `useAnimatedKeyboard`, a `useFrameCallback` that computes new targets every frame, and any shared value a library owns (a bottom sheet's `animatedIndex`, a carousel's progress, `useScrollOffset`): name the library in the reason.

A gesture callback that runs on the UI thread, fires once per interaction (`onBegin`, `onStart`, `onEnd`, `onFinalize`; `onActivate` and `onDeactivate` in the Gesture Handler 3 hooks) and writes a target for a style key that nothing drives per frame: CSS needs that target in React state, which a UI-thread callback cannot set directly. A transform entry counts as the whole `transform` key: when `onUpdate` writes `translateX`, a `scale` written on release stays on shared values with it. Once the discrete writes become state setters, look at the gesture's other callbacks:

- None must stay on the UI thread (none writes a shared value the walk keeps, none calls `scrollTo` or another UI-only API; an `onUpdate` that only reads can move): set `runOnJS: true` on the gesture (`.runOnJS(true)` on a builder) and call the state setters directly. This is the simpler form.
- One must (an `onUpdate` tracking the finger): keep the gesture as it is and call `scheduleOnRN(setTarget, value)` in the discrete callback only (`../animations/animations.md`, Threading); `runOnJS: true` would move that `onUpdate` to the JS thread too.

Each gesture of a composed gesture has its own `runOnJS`, so decide per gesture. Both forms are Needs approval: the animation now starts after a React render instead of on the next UI frame.

A press pair (`onPressIn`/`onPressOut`, or a Tap or LongPress gesture's `onBegin`/`onFinalize`) that writes the pressed and the rest value of a property on the component that owns the style is question 8, which names its CSS form. A Pan or Pinch activate/deactivate pair is not a press: it stays here.

Where gesture callbacks run, the same in Gesture Handler 2 and 3: on the UI thread when `runOnJS` is not `true` (a literal or a shared value holding `true`) and the callbacks are worklets, otherwise on the JS thread (the JS thread paragraph below). A builder (`Gesture.Pan().onStart(...)`) needs every callback to be a worklet: one plain callback moves all of them to the JS thread. A Gesture Handler 3 hook (`usePanGesture({ onActivate, onUpdate, onDeactivate })`) rejects a plain callback next to a worklet one, so all of its callbacks run on one thread. Inline builder callbacks are worklets on every Reanimated 4.x, inline Gesture Handler 3 hook callbacks from Reanimated 4.2.0; before that they are plain functions unless marked `'worklet'`, so such a hook runs on the JS thread. A callback defined elsewhere (imported, passed as a prop) is a worklet only when marked `'worklet'`. The Gesture Handler 3 hooks call the builder's `onStart`/`onEnd` `onActivate`/`onDeactivate` and have no `onChange` (its fields are in `onUpdate`); builders keep their names in both versions. Legacy handler components (`<PanGestureHandler onGestureEvent>`) call plain JS functions on Reanimated 4, which removed `useAnimatedGestureHandler`: the JS thread paragraph below.

A `runOnUI`/`scheduleOnUI` body is how a write reaches the UI runtime, not a source: classify the target the body writes. A literal, state or prop from the handler or effect that called `runOnUI` is the JS thread paragraph below (set state where `runOnUI` was called). A body whose target comes from another shared value is a chain: follow it back to its first write (Chains, below). A body that also does UI-thread work of its own (`scrollTo`, a gesture state) or runs per frame stays on shared values.

Targets from the JS thread, so the walk continues: state, props, `useEffect`, `onPress` and other handlers, a timer, a network callback, a gesture callback that runs on the JS thread. This also covers a `with*` written inside `useAnimatedStyle` whose target is a prop or state: the hook body runs on the UI thread, but its target changes only when React renders a new prop or state value, so that value is the source (`visible` here) and the site needs no shared value at all:

```tsx
const style = useAnimatedStyle(() => ({ opacity: withTiming(visible ? 1 : 0, { duration: 200 }) }), [visible]);
```

A writer outside the requested scope: widen the scope, or note Needs approval.

Chains: a `useDerivedValue`, a `useAnimatedReaction`, or a `SharedValue` passed as a prop is not a source. Follow the chain back to where the value is first written. When that write is on the JS thread, the whole chain collapses: set state where the shared value was written, drop the derived value and the reaction, and let question 4 judge the computation the chain performed. A reaction that reacts to continuous input, or that does UI-thread work of its own (scrollTo, a gesture state), keeps the site on shared values; a reaction whose only job is to forward a JS-written value is removed with it. When you cannot tell what a reaction does for other code, note Needs approval and ask.

A custom hook that wraps `useAnimatedStyle` (`useFadeIn()`) is one site: classify the hook body, and the verdict must hold for every call site, which then receives plain style props.

A measured value (`onLayout`, `measure` in an effect or a handler) is a JS-thread target whether the code stored it in state, wrote it into the shared value, or passed it to `withTiming` directly: hold it in state and transition to it. `measure` read inside the hook every frame is continuous input.

A shared value passed directly in `style` or as a prop (`style={{ opacity: sv }}`, `<AnimatedCircle r={r} />`) has no hook; classify the writes to that value the same way.

## 9. Other writers

With a transition on a property, every change of that property animates, so a write that jumped at once would glide. Such writes: `sv.value = x` without a `with*` (a jump in a handler, a jump on mount), or a branch in the hook body or a style-array entry that gives the property a different value from a condition other than the animated driver (`opacity: disabled ? 0.5 : opacity.value`).

The CSS form of a jump: leave the property out of `transitionProperty` in the render that sets the new value, and put it back in the next render, when the value no longer changes, so nothing animates. A flag in state set with the new value and cleared in an effect does it:

```tsx
const [opacity, setOpacity] = useState(1);
const [jumping, setJumping] = useState(false);
const jumpTo = (value: number) => { setOpacity(value); setJumping(true); };
useEffect(() => { if (jumping) setJumping(false); }, [jumping]);
// on the element, with transform transitioning too:
// { opacity, transitionProperty: jumping ? 'transform' : ['opacity', 'transform'], transitionDuration: 200 }
```

Leave out only the property that jumps. When it is the only transitioned property, `'none'` removes the whole transition config, and adding the config back to a mounted view animates from the previous value only from 4.6.0 (earlier versions animate from the property's default), so below 4.6.0 such a site stays on shared values. Propose the recipe as Needs approval, saying which write it covers. A reset followed at once by a `with*` is a replay, not a jump: question 6, restart. A site whose writes are all instant never animated: plain state in the static style, no `transitionProperty`.

## 10. Other readers

Removing the shared value breaks every other reader. Split them by what they consume. Another hook that is a site in scope (this component's, or a child's through a prop) is walked on its own, and question 1's Chains rule collapses both to the same state; a looping value read by several hooks: `references/transitions-and-animations.md`, shared phase. A reader that uses the value only at rest, JS reading `sv.value` in a handler that cannot run while the value animates: note Needs approval proposing a state mirror for it. A reader that consumes the value while it animates, a worklet (a `useAnimatedReaction` that Chains did not remove, a gesture callback, `scrollTo`), a child or library component outside the scope rendering it from a prop, or a handler that can run mid-flight: Keep on shared values, naming the reader as the reason.

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

The row records `inOut(quad)` to `'ease-in-out'` (0.012) and, when `visible` can flip back inside 200ms, the shortened return from question 8. `visible` starts `false` here; when it can start `true` the effect faded the element in on mount, which is the mount edge above (render `0`, flip the state in a mount effect). With reduced motion kept, `transitionDuration` becomes `reduced ? 1 : 200`.

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
