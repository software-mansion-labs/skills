# Animations

Production-quality animation patterns for React Native apps using Reanimated 4 on the New Architecture.

For animation function APIs and core hooks, see **`animation-functions.md`**.
For entering/exiting and layout transition animations, see **`layout-animations.md`**.
For scroll-driven animations and event-based patterns, see **`scroll-and-events.md`**.
For canvas animations with Skia (high element counts, sprites, path morphing), see **`canvas-animations.md`**.
For GPU shader animations (particles, noise, SDF, physics, 3D), see **`gpu-animations.md`**.
For performance tuning and feature flags, see **`animations-performance.md`**.

---

## Decision Tree

Pick the animation type based on what drives the animation and what it needs to compute.

```
Does the effect require per-pixel GPU computation?
(Particle systems, fluid/physics sims, procedural noise, SDF shapes, 3D scenes)
├── YES → Use GPU Shaders (react-native-wgpu + TypeGPU)   → see gpu-animations.md
└── NO  → Does it animate more than ~100 elements (low-end Android) or ~500 (iOS)?
    ├── YES → Use Reanimated + react-native-skia           → see canvas-animations.md
    └── NO  → Is the animation driven by a state change (not a gesture or continuous input)?
        ├── YES → Can it be expressed as a simple A→B property transition?
        │   ├── YES → Use CSS Transition (transitionProperty)
        │   └── NO  → Does it need a defined keyframe sequence?
        │       ├── YES → Use CSS Animation (animationName + keyframes)
        │       └── NO  → Use CSS Transition with multiple properties
        └── NO  → Is it gesture-driven, or does it need math / trig / layout reads?
            ├── Simple feedback (press/release, toggle)?
            │   └── YES → Use CSS Transition with `:active` (4.5.0+), else Pressable + React state
            └── Continuous tracking, math, or layout reads?
                └── YES → Use Shared Value Animation (useSharedValue + useAnimatedStyle)
```

Default to CSS transitions and CSS animations. They are declarative, easier to read, and remove the overhead of worklet execution. This includes simple gesture feedback like button presses: from 4.5.0 use a CSS transition with the `:active` pseudo-selector, and below that a CSS transition driven by `Pressable` + React state; either way you avoid shared values, worklets and thread bridging. Reach for shared values when the animation requires continuous tracking (pan, pinch, scroll), per-frame math, or layout reads. When the scene animates more than ~100 elements on low-end Android or ~500 on iOS, switch to Reanimated + `react-native-skia`, which renders to a single canvas and avoids per-view overhead. Reach for GPU shaders (`react-native-wgpu` + TypeGPU) when the animation involves per-pixel computation, physics simulations, or 3D rendering that operates outside the React Native view hierarchy.

---

## CSS feature availability

Check the installed version first (see `SKILL.md`). Everything below works from Reanimated 4.0.0 unless a row says otherwise; a feature used on an older version is silently ignored or throws.

| Feature | From |
|---|---|
| CSS transitions and CSS animations: all `transition*` and `animation*` properties, keyframes, and every timing function (the named ones like `'ease-in-out'`, plus `cubicBezier()`, `steps()` and `linear()`) | 4.0.0 |
| `filter` and its functions (`blur`, `brightness`, `dropShadow`, ...) on iOS and Android; web has it from 4.0.0 | 4.2.0 |
| CSS on `react-native-svg` components, iOS and Android (declarations go in `animatedProps`, see `svg-animations.md`; on 4.1.0-4.3.x only with the `EXPERIMENTAL_CSS_ANIMATIONS_FOR_SVG_COMPONENTS` flag) | 4.4.0 |
| CSS on `react-native-svg` components, web | 4.5.0 |
| Pseudo-selectors (`:hover`, `:active`, `:active-deepest`, `:focus`, `:focus-within`) | 4.5.0 |
| CSS animation and transition callbacks (`onCSSAnimation*`, `onCSSTransition*`) | 4.6.0 |

---

## Shared rules for CSS transitions and animations

- **Units**: bare numbers in every `transition*` and `animation*` duration or delay are milliseconds (`transitionDuration: 300` is 300ms); strings such as `'300ms'` or `'0.3s'` work too.
- **Timing functions** are the same for both: the named ones (`'linear'`, `'ease'`, `'ease-in'`, `'ease-out'`, `'ease-in-out'`) plus `cubicBezier()`, `steps()` and `linear()` imported from `react-native-reanimated`. The `Easing` object used by `withTiming` is not compatible and throws. With `steps()`, pass the modifier you mean (`steps(4, 'jump-start')`): the default is `'jump-end'`, while `Easing.steps` defaults to jump-start, so a port that keeps the count and drops the modifier shifts every step.
- **Values must be the same kind on both sides**: `height: open ? 300 : 'auto'` cannot animate between a number and a keyword, so it jumps to the target. Declare the property in both states, or in both keyframes, with the same kind of value.
- **Colors** interpolate as straight sRGB. `withTiming` and `interpolateColor` gamma-correct, so the same two endpoints produce a visibly different midpoint on wide swings (black to white, red to cyan). Alpha and `opacity` fades match exactly.

### Callbacks

From 4.6.0 the animated component takes lifecycle callbacks as props, never as style keys. Transitions report per transitioning property, animations per animation; the event carries `elapsedTime` in seconds (`transitionDuration: 300` reports `0.3`) plus `propertyName` for a transition or `animationName` for an animation. There is no `finished` flag: `End` is completion, `Cancel` is interruption.

```tsx
<Animated.View
  style={{ opacity: visible ? 1 : 0, transitionProperty: 'opacity', transitionDuration: 300 }}
  onCSSTransitionRun={(e) => console.log('triggered, before any delay', e.propertyName)}
  onCSSTransitionStart={(e) => console.log('started, after the delay', e.propertyName)}
  onCSSTransitionEnd={(e) => console.log('finished', e.propertyName, e.elapsedTime)}
  onCSSTransitionCancel={(e) => console.log('interrupted: retargeted mid-flight or unmounted', e.propertyName)}
/>

<Animated.View
  style={{ animationName: pulse, animationDuration: '1200ms', animationIterationCount: 3 }}
  onCSSAnimationStart={(e) => console.log('started, after animationDelay', e.animationName)}
  onCSSAnimationIteration={(e) => console.log('an iteration ended, except the last', e.animationName)}
  onCSSAnimationEnd={(e) => console.log('finished', e.animationName, e.elapsedTime)}
  onCSSAnimationCancel={(e) => console.log('interrupted or unmounted', e.animationName)}
/>
```

An infinite animation never reaches `onCSSAnimationEnd`; its only terminal event is `onCSSAnimationCancel`. Transition callbacks fire for pseudo-selector-driven transitions too.

### Reduced motion

CSS transitions and animations have no reduced-motion option. Unlike `with*` animations (`withTiming`, `withSpring`, ...), which follow the device setting by default (`ReduceMotion.System`), they run regardless of it. Read `useReducedMotion()` and shorten them yourself. Shorten rather than remove: a 1ms run still reaches its end state, keeps `animationFillMode` and fires the callbacks above, whereas dropping `animationName` discards the fill mode too.

```tsx
const reduced = useReducedMotion();

<Animated.View
  style={{
    opacity: visible ? 1 : 0,
    transitionProperty: 'opacity',
    transitionDuration: reduced ? 1 : 300,
  }}
/>

<Animated.View
  style={{
    animationName: pulse,
    animationDuration: reduced ? 1 : '1200ms',
    animationIterationCount: reduced ? 1 : 'infinite',
  }}
/>
```

Use `1` (1ms), never `0`, when the transition must not be dropped, for example to still receive its events: a transition whose duration plus delay is `0` is removed entirely. Cap `animationIterationCount` at `1` so a loop does not strobe. To reduce only the movement and keep the rest smooth, swap the keyframes or shorten only the moving property instead: `animationName: reduced ? fadeIn : slideIn`, or `transitionDuration: reduced ? [1, 300] : [300, 300]` for `transitionProperty: ['transform', 'opacity']`.

---

## CSS Transitions

Use when a style property should animate whenever a state-driven value changes. For the full property list and timing functions, webfetch the [CSS Transitions docs](https://docs.swmansion.com/react-native-reanimated/docs/category/css-transitions).

```tsx
<Animated.View
  style={{
    width: isExpanded ? 200 : 100,
    transitionProperty: 'width',
    transitionDuration: 300,
    transitionTimingFunction: 'ease-out',
  }}
/>
```

A transition runs when the property's value differs from the previously rendered one: the driver is React state, a prop, or from 4.5.0 a pseudo-selector. It never runs on mount, and a shared value written on the UI thread does not re-render, so it never triggers one.

When using arrays, the order must match the `transitionProperty` array:

```tsx
transitionProperty: ['width', 'opacity', 'backgroundColor'],
transitionDuration: [300, 200, 150],
transitionTimingFunction: ['ease-out', 'linear', 'ease-in-out'],
```

### Simple gesture feedback

Press feedback is a transition too. Which element gets the style decides the mechanism: a pseudo-selector on the pressed element, or `Pressable`'s state for anything else. `:hover`, `:focus` and the selector rules are in `css-pseudo-selectors.md`.

**The pressed element styles itself.** From 4.5.0 write the pressed value inline with the `:active` pseudo-selector. Pseudo-selectors work on any `Animated` component (and on `react-native-svg` elements from 4.6.0); the `Pressable` here only provides `onPress`. Nothing re-renders.

```tsx
import { Pressable } from 'react-native-gesture-handler';
import Animated from 'react-native-reanimated';

const AnimatedPressable = Animated.createAnimatedComponent(Pressable);

<AnimatedPressable
  onPress={onPress}
  style={{
    transform: { default: [{ scale: 1 }], ':active': [{ scale: 0.96 }] },
    boxShadow: {
      default: '0px 6px 10px rgba(0, 0, 0, 0.3)',
      ':active': '0px 1px 2px rgba(0, 0, 0, 0.3)',
    },
    transitionProperty: ['transform', 'boxShadow'],
    transitionDuration: '80ms',
  }}
/>
```

`:active` matches the pressed element and every ancestor that declares `:active`, so a card with `:active` also reacts when a button inside it is pressed. `:active-deepest` matches only the innermost element under the finger that declares a press selector, never an ancestor: put it on a container that should react to presses on its own area but stay still while an inner control declaring `:active` or `:active-deepest` is pressed.

**Descendants of the pressed element get styled.** Pseudo-selectors do not help here: a descendant matches `:active` only when the finger is on it. Use the approach that predates pseudo-selectors, `Pressable`'s render prop, which also covers every version below 4.5.0:

```tsx
import { Text } from 'react-native';
import { Pressable } from 'react-native-gesture-handler';
import Animated from 'react-native-reanimated';

function PressableButton({ label, onPress }: { label: string; onPress: () => void }) {
  return (
    <Pressable onPress={onPress}>
      {({ pressed }) => (
        <Animated.View
          style={{
            transform: pressed ? [{ scale: 0.96 }, { translateY: 4 }] : [{ scale: 1 }, { translateY: 0 }],
            boxShadow: pressed ? '0px 1px 2px rgba(0, 0, 0, 0.3)' : '0px 6px 10px rgba(0, 0, 0, 0.3)',
            transitionProperty: ['transform', 'boxShadow'],
            transitionDuration: '80ms',
          }}>
          <Text>{label}</Text>
        </Animated.View>
      )}
    </Pressable>
  );
}
```

**The `Pressable` itself, or an ancestor, gets styled without pseudo-selectors** (below 4.5.0, or when the pressed state must reach an ancestor): keep the pressed flag in React state set from `onPressIn`/`onPressOut` and drive the same transition from it:

```tsx
import { useState } from 'react';
import { Pressable } from 'react-native-gesture-handler';
import Animated from 'react-native-reanimated';

const AnimatedPressable = Animated.createAnimatedComponent(Pressable);

function PressableCard({ children, onPress }: { children: React.ReactNode; onPress: () => void }) {
  const [pressed, setPressed] = useState(false);

  return (
    <AnimatedPressable
      onPress={onPress}
      onPressIn={() => setPressed(true)}
      onPressOut={() => setPressed(false)}
      style={{
        transform: pressed ? [{ scale: 0.96 }] : [{ scale: 1 }],
        transitionProperty: 'transform',
        transitionDuration: '80ms',
      }}>
      {children}
    </AnimatedPressable>
  );
}
```

Reserve shared value animations for continuous gesture tracking (pan, pinch, scroll-driven) where the animation must follow finger position on every frame without a JS thread round-trip.

### Discrete properties

Properties like `flexDirection`, `justifyContent`, and `alignItems` cannot be smoothly animated. In a transition they change instantly by default. To make them flip at the transition midpoint instead, set:

```tsx
transitionBehavior: 'allow-discrete',
```

In a CSS animation they always flip halfway between the two keyframes. The `display` property is special-cased around `none`: changing from `none` to another value flips at the start, and changing to `none` holds the visible value until the end, which is what makes it usable for enter/exit. To animate the layout change a keyword flip causes, put a layout transition on the affected views instead (`layout-animations.md`).

### Rules

- `transitionProperty` defaults to `'all'` when omitted, which transitions every property that changes. List the properties explicitly when only some of them should animate.
- Always set `transitionDuration`. The default is `0`, which discards the motion. The default timing function is `'ease'`.
- Reversing a running transition (the state flips back while the transition is still running) returns over the remaining distance in proportionally less time, like a browser: reversed 100ms into a 300ms linear transition, the way back takes about 100ms. `withTiming` would restart at the full 300ms.
- Negative delays start the transition partway through (e.g., `'-5s'` on a 10s transition starts at 50%).

---

## CSS Animations

Use when the animation follows a predefined keyframe sequence independent of external state — loaders, pulse effects, entrance choreography. For the full property list, webfetch the [CSS Animations docs](https://docs.swmansion.com/react-native-reanimated/docs/category/css-animations).

```tsx
const pulse = {
  '0%':   { opacity: 1 },
  '50%':  { opacity: 0.4 },
  '100%': { opacity: 1 },
};

<Animated.View
  style={{
    animationName: pulse,
    animationDuration: '1200ms',
    animationIterationCount: 'infinite',
    animationTimingFunction: 'ease-in-out',
  }}
/>
```

Keyframe offsets are percentages, `from`/`to`, or numbers in 0..1. The element's current style is the implicit first keyframe, so you only need to define the frames that differ. At minimum, one keyframe is required.

### Mount animations

CSS attaches after the first paint, so a mount animation needs its start value in the static style too, or the first frame shows the resting value. Add `animationFillMode: 'forwards'` to stay at the end; with the default `'none'` the element snaps back to its static style when the animation finishes.

```tsx
<Animated.View
  style={{
    opacity: 0,
    animationName: { to: { opacity: 1 } },
    animationDuration: 300,
    animationFillMode: 'forwards',
  }}
/>
```

### Multiple animations

```tsx
const fadeInOut = { '0%': { opacity: 0 }, '100%': { opacity: 1 } };
const moveLeft = { '100%': { transform: [{ translateX: -100 }] } };

<Animated.View
  style={{
    animationName: [fadeInOut, moveLeft],
    animationDuration: ['2.5s', '5s'],
    animationIterationCount: ['infinite', 1],
  }}
/>
```

Every `animation*` setting takes a parallel array, one entry per animation. If multiple animations target the same property, the later animation in the array wins.

### Defining keyframes

Prefer `css.keyframes()` (`css` imported from `react-native-reanimated`) called once outside the component: the keyframes are processed once and every component that uses the rule shares it.

```tsx
// Best: processed once, shared by every component that uses it
const pulse = css.keyframes({ '50%': { opacity: 0.4 } });

function Dot() {
  return <Animated.View style={{ animationName: pulse, animationDuration: '1200ms' }} />;
}
```

A plain keyframes object is matched by its content, so it never restarts the animation on re-render, but it is re-checked on every render and nothing is shared. Defining it outside the component changes nothing; these two are the same:

```tsx
const pulse = { '50%': { opacity: 0.4 } };

function Dot() {
  return <Animated.View style={{ animationName: pulse, animationDuration: '1200ms' }} />;
}

function Dot() {
  return <Animated.View style={{ animationName: { '50%': { opacity: 0.4 } }, animationDuration: '1200ms' }} />;
}
```

`css.keyframes()` called inside the component creates a new rule on every render and restarts the animation each time. Use it only to re-trigger the animation on purpose:

```tsx
function Dot() {
  const pulse = css.keyframes({ '50%': { opacity: 0.4 } });

  return <Animated.View style={{ animationName: pulse, animationDuration: '1200ms' }} />;
}
```

### Rules

- `animationTimingFunction` at the top level eases every interval between two consecutive keyframes of a property. A keyframe can carry its own `animationTimingFunction` to override it for the interval that starts there, up to the next keyframe that sets the same property; one on the last keyframe has no interval and is ignored.
- Avoid `animationFillMode: 'forwards'` or `'both'` when a fractional `animationIterationCount` meets keyframes that mix relative (percentage) and absolute units for one property. If the parent resizes after the animation, the child retains stale dimensions. Either alone is fine.
- For infinite CSS animations, set `animationIterationCount: 'infinite'`. The animation stops automatically on unmount — no manual cleanup needed.
- Negative delays start the animation partway through its cycle.
- Pause and resume with `animationPlayState: 'paused'` / `'running'`.

---

## Shared Value Animations

Use when:
- The animation is driven by a gesture or continuous input (scroll position, drag offset)
- It requires math, trigonometric functions, or interpolation between computed values
- It needs to read layout measurements on each frame (`measure`, `useAnimatedRef`)
- Multiple animated values need to be derived from a single source of truth

```tsx
const offset = useSharedValue(0);

const animatedStyle = useAnimatedStyle(() => ({
  transform: [{ translateX: withSpring(offset.value) }],
}));

// Gesture-driven example
const gesture = Gesture.Pan().onUpdate((e) => {
  offset.value = e.translationX;
});
```

Avoid reading `sharedValue.value` on the JS thread inside React render or event handlers — it causes a synchronization that blocks the JS thread. Derive values from shared values using `useDerivedValue` instead.

---

## Animating Text

Avoid updating `Animated.Text` content by changing state — it triggers a full React re-render for every frame.

For animated numeric counters or any frequently-changing text, use `AnimatedTextInput` with `animatedProps`:

```tsx
import Animated, { useAnimatedProps, type SharedValue } from 'react-native-reanimated';
import { TextInput } from 'react-native';

const AnimatedTextInput = Animated.createAnimatedComponent(TextInput);

function Counter({ progress }: { progress: SharedValue<number> }) {
  const animatedProps = useAnimatedProps(() => ({
    text: String(Math.round(progress.value)),
    defaultValue: '0',
  }));

  return (
    <AnimatedTextInput
      animatedProps={animatedProps}
      editable={false}
      style={styles.counter}
    />
  );
}
```

This updates the native text node directly on the UI thread, bypassing React and eliminating re-renders.

---

## Infinite Animations

CSS animations with `animationIterationCount: 'infinite'` clean up automatically on unmount.

For shared value infinite animations, always cancel them in the `useEffect` cleanup:

```tsx
useEffect(() => {
  offset.value = withRepeat(withTiming(1, { duration: 800 }), -1, true);

  return () => {
    cancelAnimation(offset);
  };
}, []);
```

Never start infinite animations outside the component lifecycle (module scope, global timers). They cannot be cleaned up and will leak.

---

## Prefer Non-Layout Properties

Animating layout properties (`top`, `left`, `width`, `height`, `margin`, `padding`) forces a layout pass on every frame, which is expensive and causes jank.

Prefer:
- `transform: [{ translateX }, { translateY }, { scale }, { rotate }]`
- `opacity`
- `backgroundColor`

If a design requires what looks like a size change, consider `scale` transforms — same visual effect without triggering layout.

---

## Supported Style Properties

Most React Native style properties are animatable. Key exceptions and platform notes:

- **`flexBasis`**: Changes are calculated but never applied to the view. Use `flexGrow`/`flexShrink` instead.
- **Shadow properties**: `shadowOffset`, `shadowOpacity`, `shadowRadius` do not work on Android. Use `boxShadow` instead (works on all platforms).
- **Web shadows**: All shadow styles must be specified in every keyframe on Web, or they are lost.
- **`tintColor` on iOS**: Must be present in the initial style when the `Image` component mounts. Adding it later has no effect.
- **Style inheritance**: Not supported. Properties that normally inherit in CSS (e.g., `textDecorationColor` from `color`) must be set explicitly.
- **Mixed-unit margins**: Interpolating between absolute and percentage margins may produce unexpected results when the parent's dimensions are affected by the child's margins.

---

## Threading: scheduleOnRN instead of runOnJS

`runOnJS` is removed in Reanimated 4. Use `scheduleOnRN` to call JS-thread functions from a worklet. Arguments are passed directly, not curried:

```tsx
// Reanimated 3 (removed)
runOnJS(setCount)(newCount);

// Reanimated 4
scheduleOnRN(setCount, newCount);
```

`scheduleOnRN` schedules the call asynchronously on the React Native runtime. Functions passed to `scheduleOnRN` must be defined in JS thread scope (they cannot be created inside worklets or animation callbacks).
