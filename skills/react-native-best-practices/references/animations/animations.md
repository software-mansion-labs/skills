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
        │   └── NO  → Does it need a defined keyframe sequence, or play on mount?
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
| CSS on `react-native-svg` components, iOS and Android (declarations go in `animatedProps`, see `svg-animations.md`). Enabled by default from 4.4.0 and still labeled experimental (`EXPERIMENTAL_CSS_ANIMATIONS_FOR_SVG_COMPONENTS`, opt-in on 4.1.0-4.3.x) | 4.4.0 |
| CSS on `react-native-svg` components, web | 4.5.0 |
| Pseudo-selectors (`:hover`, `:active`, `:active-deepest`, `:focus`, `:focus-within`) | 4.5.0 |
| CSS animation and transition callbacks (`onCSSAnimation*`, `onCSSTransition*`) | 4.6.0 |

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

`transitionDuration: 300` is 300ms: bare numbers are milliseconds in every `transition*` and `animation*` duration or delay, and strings such as `'300ms'` or `'0.3s'` work too.

A transition runs when the property's value differs from the previously rendered one: the driver is React state, a prop, or from 4.5.0 a pseudo-selector. It never runs on mount, and a shared value written on the UI thread does not re-render, so it never triggers one.

When using arrays, the order must match the `transitionProperty` array:

```tsx
transitionProperty: ['width', 'opacity', 'backgroundColor'],
transitionDuration: [300, 200, 150],
transitionTimingFunction: ['ease-out', 'linear', 'ease-in-out'],
```

### Simple gesture feedback

Press feedback is a transition too. Two cases:

**The pressed element styles itself.** From 4.5.0 write the pressed value inline with the `:active` pseudo-selector. Pseudo-selectors work on any `Animated` component (and on `react-native-svg` elements from 4.6.0); the `Pressable` here only provides `onPress`. Nothing re-renders.

```tsx
import { Pressable } from 'react-native-gesture-handler';
import Animated from 'react-native-reanimated';

const AnimatedPressable = Animated.createAnimatedComponent(Pressable);

<AnimatedPressable
  onPress={onPress}
  style={{
    transform: { default: [{ scale: 1 }], ':active': [{ scale: 0.96 }] },
    transitionProperty: 'transform',
    transitionDuration: 80,
  }}
/>
```

`:active` fires on the element under the finger and on every ancestor that declares `:active`; put `:active-deepest` on an ancestor that must stay quiet while a pressed descendant handles the feedback.

**A child of the pressed element styles itself.** A child the finger does not land on never matches `:active`, so drive it from `Pressable`'s render prop instead. This form also works below 4.5.0:

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
            transform: pressed ? [{ scale: 0.96 }] : [{ scale: 1 }],
            transitionProperty: 'transform',
            transitionDuration: 80,
          }}>
          <Text>{label}</Text>
        </Animated.View>
      )}
    </Pressable>
  );
}
```

Below 4.5.0, when the `Pressable` itself carries the style, keep a `useState` set from `onPressIn`/`onPressOut` on the animated `Pressable`.

Reserve shared value animations for continuous gesture tracking (pan, pinch, scroll-driven) where the animation must follow finger position on every frame without a JS thread round-trip.

### Discrete properties

Properties like `flexDirection`, `justifyContent`, and `alignItems` cannot be smoothly animated. In a transition they change instantly by default. To make them flip at the transition midpoint instead, set:

```tsx
transitionBehavior: 'allow-discrete',
```

In a CSS animation they always flip halfway between the two keyframes. The `display` property is special-cased around `none`: leaving `none` it flips at the start, and going to `none` it holds the visible value until the end, which is what makes it usable for enter/exit. To animate the layout change a keyword flip causes, put a layout transition on the affected views instead (`layout-animations.md`).

### Rules

- `transitionProperty` defaults to `'all'` when omitted, which transitions every property that changes. List the properties explicitly when only some of them should animate.
- Always set `transitionDuration`. The default is `0`, which discards the motion. The default timing function is `'ease'`.
- Both values must be the same kind: `height: open ? 300 : 'auto'` cannot animate between a number and a keyword, so it jumps to the target. Declare the property in both states with the same kind of value.
- Reversing a running transition (a press released while the press-in transition is still running) returns over the remaining distance in proportionally less time, like a browser: reversed 100ms into a 300ms linear transition, the way back takes about 100ms. `withTiming` would restart at the full 300ms.
- Colors interpolate as straight sRGB. `withTiming` and `interpolateColor` gamma-correct, so the same two endpoints produce a visibly different midpoint on wide swings (black to white, red to cyan). Alpha and `opacity` fades match exactly.
- Negative delays start the transition partway through (e.g., `'-5s'` on a 10s transition starts at 50%).

### Reduced motion

CSS transitions do not react to the system reduce-motion setting on their own (`withTiming` does). Read `useReducedMotion()` and shorten the transition rather than removing it, so the end state and any callbacks still arrive:

```tsx
const reduced = useReducedMotion();

<Animated.View
  style={{
    opacity: visible ? 1 : 0,
    transitionProperty: 'opacity',
    transitionDuration: reduced ? 1 : 300,
  }}
/>
```

Use `1` (1ms), never `0`: a transition whose duration plus delay is `0` is dropped entirely, so nothing fires.

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
// Best: processed once, shared, never restarts on re-render
const pulse = css.keyframes({ '50%': { opacity: 0.4 } });

// Fine: a plain object is matched by content, so it does not restart on re-render,
// but it is re-checked on every render
const pulseInline = { '50%': { opacity: 0.4 } };

function Dot() {
  // Restarts the animation on EVERY render: `css.keyframes()` inside a component
  // creates a new rule each time. Use this only to re-trigger the animation on purpose.
  const restarting = css.keyframes({ '50%': { opacity: 0.4 } });

  return <Animated.View style={{ animationName: pulse, animationDuration: '1200ms' }} />;
}
```

A plain keyframes object behaves the same inline or outside the component: matched by content, no restart. Only an inline `css.keyframes()` call restarts the animation each render.

### Rules

- `animationTimingFunction` at the top level eases every interval between two consecutive keyframes of a property. A keyframe can carry its own `animationTimingFunction` to override it for the interval that starts there, up to the next keyframe that sets the same property; one on the last keyframe has no interval and is ignored.
- Avoid `animationFillMode: 'forwards'` or `'both'` with fractional `animationIterationCount` and relative units (percentages). If the parent resizes after the animation, the child retains stale dimensions.
- `animationIterationCount: 'infinite'` runs until the component unmounts; nothing to clean up.
- Negative delays start the animation partway through its cycle.
- Pause and resume with `animationPlayState: 'paused'` / `'running'`.

### Timing functions

`cubicBezier`, `steps` and `linear` are exported from `react-native-reanimated` and work in transitions and animations alike. Passing an `Easing.*` value throws. Write the `steps` modifier explicitly, `steps(4, 'jump-start')`: the default is `'jump-end'`, the opposite of `Easing.steps`. Modifiers: `'jump-start'`, `'jump-end'`, `'jump-none'`, `'jump-both'`, `'start'`, `'end'`.

### Reduced motion

CSS animations do not react to the system reduce-motion setting on their own either. Shorten the animation instead of removing it, and cap the iteration count so a loop does not strobe; a 1ms run still reaches its end state, keeps `animationFillMode` and fires the end callback:

```tsx
const reduced = useReducedMotion();

<Animated.View
  style={{
    animationName: pulse,
    animationDuration: reduced ? 1 : '1200ms',
    animationIterationCount: reduced ? 1 : 'infinite',
  }}
/>
```

Do not remove `animationName` instead: that discards the fill mode too, so an element whose static style is `opacity: 0` never appears. When the motion carries meaning (a slide-in), replace it rather than shorten it: `animationName: reduced ? fadeIn : slideIn`.

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
