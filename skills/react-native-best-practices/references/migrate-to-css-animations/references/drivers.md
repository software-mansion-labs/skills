# Drivers: where the values come from, who else touches them

Questions 1, 9 and 10 of the walk.

## 1. Where do the target values come from?

A `with*` call produces the frames between two targets; CSS replaces exactly that part. What decides the verdict is where the targets come from.

Continuous input, so the site stays on shared values: a scroll offset, a moving finger (a gesture whose `onUpdate` drives the value, or `onChange` in Gesture Handler 2), a sensor, `useAnimatedKeyboard`, a `useFrameCallback` that computes new targets every frame, and any shared value a library owns (a bottom sheet's `animatedIndex`, a carousel's progress, `useScrollOffset`): name the library in the reason.

A worklet gesture callback that fires once per interaction (`onBegin`, `onStart`, `onEnd`, `onFinalize`; in the Gesture Handler 3 hooks `onActivate` and `onDeactivate`) and writes a target for a property that nothing drives per frame: the callback runs on the UI thread, and the CSS form needs its target in React state. Propose `scheduleOnRN(setTarget, value)`, imported from `react-native-worklets` (Reanimated does not re-export it), inside that callback (`runOnJS(setTarget)(value)` on 4.0.x, which has no `scheduleOnRN`; `runOnJS` is deprecated from 4.1.0 but still exported) and hold the target in state, as Needs approval, saying that the animation then starts after a JS round trip and a render instead of on the next UI frame. Never propose `runOnJS: true` on the gesture (`.runOnJS(true)` on a builder): it moves every callback to the JS thread, `onUpdate` included, so a finger-tracking property in the same gesture would stutter. A press pair (`onPressIn`/`onPressOut`, or a Tap or LongPress gesture's `onBegin`/`onFinalize`) that writes the pressed and the rest value of a property on the component that owns the style is question 8, which names its CSS form (from 4.5.0 `:active`; below it React state set from the handlers, through this move when they are worklets). A Pan or Pinch activate/deactivate pair is not a press: it stays in this paragraph.

Gesture Handler 2 and 3 both apply. The builders (`Gesture.Pan().onStart(...)`, both majors) run their callbacks on the UI thread only when `runOnJS` is not `true` and every defined callback is a worklet: one plain callback, typically a function passed by reference without `'worklet'`, moves all of them to the JS thread, the same effect as `.runOnJS(true)`. The 3.x hooks (`usePanGesture({ onActivate, onUpdate, onDeactivate })`) take the UI-thread path when `runOnJS` is not `true` and at least one callback is a worklet, and a plain callback next to a worklet one is an error there, so in practice every callback is a worklet or none is. Inline callbacks are worklets automatically: for the builders on every 4.x, for the hook configs only from Reanimated 4.2.0 (worklets 0.7.0, the package that ships the Babel plugin); on 4.0.x and 4.1.x an inline hook callback is a plain function and the whole gesture runs on the JS thread, which is the JS thread arm. A callback defined elsewhere and passed by reference is a worklet only when marked `'worklet'`. The 3.x hooks rename `onStart`/`onEnd` to `onActivate`/`onDeactivate` and fold `onChange` into `onUpdate`. The legacy `<PanGestureHandler onGestureEvent onHandlerStateChange>` components fire JS-thread callbacks on Reanimated 4 (`useAnimatedGestureHandler` no longer exists there), so a `with*` inside them is the JS thread arm.

A `runOnUI`/`scheduleOnUI` body is how a write reaches the UI runtime, not a source: classify the target the body writes. A literal, state or prop from the handler or effect that called `runOnUI` is the JS thread arm (set state where `runOnUI` was called; a plain `sv.value = withTiming(x)` in a handler skipped that render too, so nothing new is lost). A body whose target comes from another shared value is a chain: follow it back to its first write (Chains, below). A body that also does UI-thread work of its own (`scrollTo`, a gesture state) or runs per frame stays on shared values.

The JS thread, so the walk continues: state, props, `useEffect`, `onPress`/`onChange`, a timer, a network callback, a gesture callback with `runOnJS: true`, or no shared value at all because the hook body reads a prop or state with the `with*` inline, where the render is the driver:

```tsx
const style = useAnimatedStyle(() => ({ opacity: withTiming(visible ? 1 : 0, { duration: 200 }) }), [visible]);
```

A writer outside the requested scope: widen the scope, or note Needs approval.

Chains: a `useDerivedValue`, a `useAnimatedReaction`, or a `SharedValue` passed as a prop is not a source. Follow the chain back to where the value is first written. When that write is on the JS thread, the whole chain collapses: set state where the shared value was written, drop the derived value and the reaction, and let question 4 judge the computation the chain performed. A reaction that reacts to continuous input, or that does UI-thread work of its own (scrollTo, a gesture state), keeps the site on shared values; a reaction whose only job is to forward a JS-written value is removed with it. When you cannot tell what a reaction does for other code, note Needs approval and ask.

A custom hook that wraps `useAnimatedStyle` (`useFadeIn()`) is one site: classify the hook body, and the verdict must hold for every call site, which then receives plain style props.

A measured value (`onLayout`, `measure` in an effect or a handler) is a JS-thread target whether the code stored it in state, wrote it into the shared value, or passed it to `withTiming` directly: hold it in state and transition to it. `measure` read inside the hook every frame is continuous input.

A shared value passed directly in `style` or as a prop (`style={{ opacity: sv }}`, `<AnimatedCircle r={r} />`) has no hook; classify the writes to that value the same way.

## 9. Other writers

With a transition on a property, every change of that property animates. A site that also sets the value without a `with*` (`sv.value = 0` before a `withTiming(1)`, a mount jump, or a branch in the hook body or a style-array entry that gives the property a different value from a condition that is not the animated driver, `opacity: disabled ? 0.5 : opacity.value`) would turn that jump into a glide. The instant write has a CSS form: render the new value with the transition turned off for that property in that render (`transitionProperty` without it, or, on 4.3.0+, `transitionDuration: 0` with no `transitionDelay`, since a property whose duration plus delay is 0 is dropped from the transition) and turn it back on in the next committed render; a reset followed by a `with*` therefore needs two commits (`flushSync`, or the second state set in an effect keyed on the first). Propose that code as Needs approval, saying which write it covers. A reset that question 6 already folded into a restart (`sv.value = 0` before the `with*` that replays from 0, converted to a new `css.keyframes()` rule) is the first keyframe of that rule, not a jump: do not note it. A site whose writes are all instant never animated: plain state in the static style, no `transitionProperty`.

A property that animates on mount and is later retargeted by state is one transition: render the start value, flip the state in a mount `useEffect` (the second render starts the transition), and list the property only in `transitionProperty`, never in both `animationName` and `transitionProperty`.

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
