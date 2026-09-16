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
    └── NO  → What drives the animation?
        ├── A gesture, scroll offset, sensor, or a value that changes every frame
        │   → Shared Value Animation (useSharedValue + useAnimatedStyle)
        ├── Per-frame math, trig, or layout reads (measure)
        │   → Shared Value Animation
        ├── Press feedback
        │   → CSS Transition with `:active` (4.5.0+), else Pressable + React state
        ├── A React state or prop change on an element that stays mounted
        │   → CSS Transition (transitionProperty)
        └── Plays by itself once started: mount, loop, keyframe sequence
            → CSS Animation (animationName + keyframes)
```

Default to CSS transitions and CSS animations: declarative, no worklets, no thread bridging. A transition needs a previously rendered value, so anything that must move on first render is an animation, whatever triggered it. Shared values are for continuous input, per-frame math and layout reads; Skia renders many elements to a single canvas; GPU shaders run outside the React Native view hierarchy.

---

## CSS feature availability

Check the installed version first (see `SKILL.md`). Everything below works from Reanimated 4.0.0 unless a row says otherwise. Under its floor a feature does not work: the value is dropped, or passed through unparsed and thrown on by Reanimated or React Native. Nothing warns about the version.

| Feature | From |
|---|---|
| CSS transitions and animations, every `animation*`/`transition*` longhand, `cubicBezier`, `steps`, `linear` | 4.0.0 |
| `filter` and its functions (`blur`, `brightness`, `dropShadow`, ...) on iOS and Android; web has it from 4.0.0 | 4.2.0 |
| CSS on `react-native-svg` components, iOS and Android (declarations go in `animatedProps`, see `svg-animations.md`; 4.1.0-4.3.x behind the `EXPERIMENTAL_CSS_ANIMATIONS_FOR_SVG_COMPONENTS` static flag) | 4.4.0 |
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

A transition runs when the property's value differs from the previously rendered one: the driver is React state, a prop, or from 4.5.0 a pseudo-selector. It never runs on mount, and a shared value written on the UI thread does not re-render, so it never triggers one. Bare numbers in every `transition*` and `animation*` duration or delay are milliseconds.

When using arrays, the order must match the `transitionProperty` array:

```tsx
transitionProperty: ['width', 'opacity', 'backgroundColor'],
transitionDuration: [300, 200, 150],
transitionTimingFunction: ['ease-out', 'linear', 'ease-in-out'],
```

### Simple gesture feedback

Press feedback is a transition too. From 4.5.0 write the pressed value inline with the `:active` pseudo-selector; nothing re-renders:

```tsx
<Animated.View
  style={{
    transform: { default: [{ scale: 1 }], ':active': [{ scale: 0.96 }] },
    transitionProperty: 'transform',
    transitionDuration: 80,
  }}
/>
```

Below 4.5.0 drive the same transition from `Pressable`'s render prop:

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

### Discrete properties

Keyword-valued properties such as `flexDirection`, `justifyContent`, `alignItems` and `display` cannot tween. A transition skips them unless you set:

```tsx
transitionBehavior: 'allow-discrete',
```

They then flip at the midpoint. `display` is the exception: leaving `none` flips at the start, going to `none` holds the visible value until the end, so it works for enter/exit. Boolean and enum-like props (`includeFontPadding`, SVG `fillRule`, `strokeLinecap`) always flip at the midpoint, with or without `allow-discrete`. To animate the layout change a keyword flip causes, put a layout transition on the affected views instead (`layout-animations.md`).

### Rules

- Always list `transitionProperty` explicitly. Setting any other `transition*` value without it defaults to `'all'`, which transitions whatever happens to change, including properties you never meant to animate.
- Always set `transitionDuration`. The default is `0`, which discards the motion. The default timing function is `'ease'`.
- Both endpoints must be the same kind of value. `height: open ? 300 : 'auto'` jumps to the target on the first frame (from 4.2.0; 4.0.x and 4.1.x flip at the midpoint), and `allow-discrete` only moves that jump to the midpoint. A property declared in one state only transitions from its default, which is `'auto'` for every dimension and inset (`width`, `height`, `minWidth`, `top`, `left`, ...), `flexBasis` and `aspectRatio`, so declare those in both states.
- Reversing a running transition shortens the return leg in proportion to the eased progress, like a browser; `withTiming` restarts at full duration on a retarget.
- Colors interpolate as straight sRGB. `withTiming` and `interpolateColor` gamma-correct, so the same two endpoints produce a visibly different midpoint on wide swings (black to white, red to cyan). Alpha and `opacity` fades match exactly.
- Negative delays start the transition partway through (e.g., `'-5s'` on a 10s transition starts at 50%).
- Respect reduced motion by shortening: `transitionDuration: reduced ? 1 : 300` from `useReducedMotion()`. Never `0`: from 4.3.0 a transition whose duration plus delay is `<= 0` is removed, so no callback fires.

---

## CSS Animations

Use when the animation follows a predefined keyframe sequence independent of external state: loaders, pulse effects, entrance choreography. For the full property list, webfetch the [CSS Animations docs](https://docs.swmansion.com/react-native-reanimated/docs/category/css-animations).

```tsx
const pulse = {
  '0%':   { opacity: 1 },
  '50%':  { opacity: 0.4 },
  '100%': { opacity: 1 },
};

<Animated.View
  style={{
    animationName: pulse,
    animationDuration: 1200,
    animationIterationCount: 'infinite',
    animationTimingFunction: 'ease-in-out',
  }}
/>
```

Keyframe offsets are percentages, `from`/`to`, or numbers in 0..1. The element's current style is the implicit first keyframe, so you only need to define the frames that differ. At minimum, one keyframe is required.

### Mount animations

CSS attaches after the first paint, so a mount animation needs its start value in the static style too, or the first frame shows the resting value. Add `animationFillMode: 'forwards'` to stay at the end; the default `'none'` snaps back when the animation finishes.

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
    animationDuration: [2500, 5000],
    animationIterationCount: ['infinite', 1],
  }}
/>
```

Every `animation*` setting takes a parallel array, one entry per animation. If two animations target the same property, the later one wins.

### Rules

- Define keyframes outside render. A plain keyframes object is compared by content; a `css.keyframes()` rule (`css` from `react-native-reanimated`) by identity, so one created inside render restarts the animation on every re-render.
- A keyframe's own `animationTimingFunction` governs only the interval that starts there. There is no carry-over: an interval whose keyframe declares none uses the animation-level function. The one on the last keyframe is ignored.
- Keep the `transform` array in the same order across all keyframes. Reordered operations, or a `matrix` entry, silently fall back to interpolating the composed matrices, which paces and paths differently.
- Avoid `animationFillMode: 'forwards'` or `'both'` with fractional `animationIterationCount` and relative units (percentages). If the parent resizes after the animation, the child retains stale dimensions.
- `animationIterationCount: 'infinite'` stops automatically on unmount, no manual cleanup needed. Negative delays start the animation partway through its cycle. Pause and resume with `animationPlayState: 'paused'` / `'running'`.
- Respect reduced motion by removing the animation: drop `animationName` and render the value the animation rests at. Never shorten it; a 1ms infinite animation strobes.

### Timing functions

`cubicBezier`, `steps` and `linear` are exported from `react-native-reanimated` and work in transitions and animations alike. Passing an `Easing.*` value throws. Write the `steps` modifier explicitly, `steps(4, 'jump-start')`: the default is `'jump-end'`, the opposite of `Easing.steps`. Modifiers: `'jump-start'`, `'jump-end'`, `'jump-none'`, `'jump-both'`, `'start'`, `'end'`.

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
