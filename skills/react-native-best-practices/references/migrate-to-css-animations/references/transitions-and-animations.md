# Transitions and animations

Step 3 of the walk: what a `with*` composition becomes.

## Transition or animation

Choose per property and per edge, never once for the whole site; one element may need both (`animationName` for the mount animation, `transitionProperty` for the driver-changed properties, never one property in both).

- **Transition**: an outside trigger drives every step (press, toggle, expand, theme change). The driver becomes React state.
- **Animation**: the steps play themselves once started (mount, loop, `withSequence`, looping `withRepeat`).
- An element conditionally rendered on the driver (`if (!open) return null`), or a mount `useEffect` that writes a value different from the initial one, is an animation on that edge (`../animations/animations.md`, Mount animations).
- Shared phase: a CSS animation restarts on every mount of its element, and N CSS animations share a phase only when mounted in the same commit. A shared value owned by an ancestor that stays mounted kept its phase across child remounts and ticked as one clock for every reader. Such sites: Needs approval, proposing to move the animation to the ancestor.

## `with*`

| Source | CSS | Notes |
|---|---|---|
| `withTiming(v, { duration, easing })` | `transitionDuration` + `transitionTimingFunction`, or a two-keyframe animation | mechanism per the section above; easing per `references/easing.md` |
| `withDelay(ms, anim)` at the top level | `transitionDelay` or `animationDelay` | retargeting mid-flight differs: the hook keeps the running animation going during the delay, a delayed transition cancels it and holds the current value; a site that retargets while running: Needs approval |
| `withDelay(ms, x)` inside `withSequence` or `withRepeat` | a hold before `x`: repeat the previous keyframe value at `x`'s offset, add `ms` to the total; inside `withRepeat` the hold starts every cycle, never `animationDelay` | not a leading delay |
| `withTiming(v, { duration: 0 })` inside `withSequence` | an instant step: the new value on the next offset, `'50%': { x: a }, '50.01%': { x: b }` (keyframes sharing an offset merge) | a sequence of holds ending in instant changes: one keyframe per hold with `animationTimingFunction: steps(1, 'jump-end')` on it |
| `withRepeat(anim)`, count omitted | `animationIterationCount: 2` | `numberOfReps` defaults to 2, `animationIterationCount` to 1: always write the count |
| `withRepeat(anim, n)`, `n <= 0` | `animationIterationCount: 'infinite'` | |
| `withRepeat(withTiming(...), n, true)` | plus `animationDirection: 'alternate'`, exact only for a symmetric easing (`linear`, `inOut(f)`, `'ease-in-out'`); CSS mirrors the easing on the return leg, the hook plays it forward. Asymmetric easing: write the ping-pong as `0%`/`50%`/`100%` keyframes with the easing on both intervals, double the duration, halve the count (a fractional count is accepted and ends on the `50%` keyframe; see the fill-mode caveat in `../animations/animations.md`) | an even `n` rests at the start value, odd at the target |
| `withRepeat(withSequence(...), n, true)` | `animationDirection: 'normal'` when the sequence ends where it starts; otherwise Needs approval, with keyframes for cycles 2+ that start at the sequence's last value | the hook does not reverse a sequence, but each later cycle continues from the sequence's last value instead of jumping back to the start |
| `withRepeat(withDelay(ms, withTiming(v)), n, true)` | none | the hook holds `v` from the second cycle on (the delay ignores the reversed target); Needs approval, or Keep on hooks |
| `withSequence(a, b, c)` | one animation, `animationDuration` = sum of the children, each keyframe at its cumulative fraction | 100ms then 300ms: stops at `0%`, `25%`, `100%`. Finite animations need `animationFillMode: 'forwards'` unless the static style already equals the resting value |

Two loops at different speeds need the array form (`animationName: [spin, pulse]`); never split one animation across entries. A sequence with a different easing per step sets `animationTimingFunction` on every keyframe whose outgoing interval needs it (`../animations/animations.md`, CSS Animations, Rules).

## Example: infinite loop

```tsx
// Before
function Spinner() {
  const rotation = useSharedValue(0);
  useEffect(() => {
    rotation.value = withRepeat(withTiming(360, { duration: 2000, easing: Easing.linear }), -1);
    return () => cancelAnimation(rotation);
  }, []);
  const style = useAnimatedStyle(() => ({ transform: [{ rotateZ: `${rotation.value}deg` }] }));
  return <Animated.View style={[styles.box, style]} />;
}
```

```tsx
// After
const rotate: CSSAnimationKeyframes = {
  from: { transform: [{ rotateZ: '0deg' }] },
  to: { transform: [{ rotateZ: '360deg' }] },
};

function Spinner() {
  const reduced = useReducedMotion();
  return (
    <Animated.View
      style={[
        styles.box,
        {
          animationName: rotate,
          animationDuration: reduced ? 1 : 2000,
          animationIterationCount: reduced ? 1 : 'infinite',
          animationTimingFunction: 'linear',
        },
      ]}
    />
  );
}
```

A JS-thread write in `useEffect` and `cancelAnimation` only in the cleanup, so the walk reaches Migrate. Under reduced motion the loop runs once for 1ms and snaps back to the static style; a non-`reverse` loop rests at its target, which here renders the same as the start (`360deg` is `0deg`), so `animationFillMode` can stay `'none'`.

## Example: play once on mount

```tsx
// Before
useEffect(() => { opacity.value = withTiming(1, { duration: 300 }); }, []);
```

```tsx
// After: keep opacity: 0 so the first painted frame matches the hook
const fadeIn: CSSAnimationKeyframes = { from: { opacity: 0 }, to: { opacity: 1 } };
// on the element
{ opacity: 0, animationName: fadeIn, animationDuration: reduced ? 1 : 300, animationTimingFunction: 'ease-in-out', animationFillMode: 'forwards' }
```

Under reduced motion the 1ms run lands on `opacity: 1` through the fill mode, as `withTiming` jumps to its target. The row records `inOut(quad)` to `'ease-in-out'`, max error 0.012, or the sampled `linear()` form if the user chose it.
