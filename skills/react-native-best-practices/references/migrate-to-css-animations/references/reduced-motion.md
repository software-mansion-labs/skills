# Reduced motion

Every migrated site gets a guard. `useReducedMotion()` returns the system setting read once when the app starts and ignores `<ReducedMotionConfig>`; `with*` read the config live. Every `with*` without `ReduceMotion.Never` jumps under reduced motion, so an unguarded migration is a delta. Guard every animation and every transition with `useReducedMotion()`; only an `opacity` or color transition under 300ms may go unguarded, reported as `reduced-motion guard dropped`.

| Kind | Form |
|---|---|
| Transition | `transitionDuration: reduced ? 1 : D`, never `0`, and `transitionDelay: reduced ? 0 : ms` when the source had a `withDelay` (the hook drops the delay under reduced motion) (`../animations/animations.md`, Reduced motion) |
| Animation | `animationDuration: reduced ? 1 : D`, `animationIterationCount: reduced ? 1 : N` and `animationDelay: reduced ? 0 : ms`; keep `animationName` so the fill mode and callbacks survive, and pick `animationFillMode` so it rests where the hook rests: `'forwards'` when the hook rests at the target, `'none'` (snapping back to the static style) when it rests at the start, `reduced ? 'none' : 'forwards'` when the two rests differ (a reversed even-count repeat of a sequence) (`../animations/animations.md`, Reduced motion) |
| `ReduceMotion.Never` in the source, or `<ReducedMotionConfig mode={ReduceMotion.Never}>` | no guard |
| `ReduceMotion.Always`, or `<ReducedMotionConfig mode={ReduceMotion.Always}>` | the reduced form for everyone, no guard |
| no config, or `<ReducedMotionConfig mode={ReduceMotion.System}>` | the guard above |
| `<ReducedMotionConfig>` whose `mode` changes at runtime | Keep on hooks: CSS has no live reduced-motion source |

Where the hook rests: `withTiming(TO)`, `withSequence` and a non-reverse `withRepeat` rest at `TO`; `withRepeat(anim, n, true)` with odd `n` runs once under reduced motion and rests at `TO`, with `n <= 0` or even `n` it rests at the start whatever `anim` is (a reversed sequence included, which otherwise rests at its last value). Evaluate the style body at that driver value to know which end the fill mode must hold.

## Trap: the resting value

```tsx
const progress = useSharedValue(0);
useEffect(() => { progress.value = withRepeat(withTiming(1, { duration: 600 }), -1, true); }, []);
const style = useAnimatedStyle(() => ({ transform: [{ scaleY: interpolate(progress.value, [0, 1], [0.35, 1]) }] }));
```

Reduced motion never advances `progress`, so the hook rests at `interpolate(0, ...)` = `0.35`, the start. The shortened animation must rest there too: static style `scaleY: 0.35` and `animationFillMode: 'none'`, so the 1ms run snaps back to it. `'forwards'` would park the element at the `to` keyframe, nearly three times too tall. Evaluate the style body at the driver's resting value, never read it off the keyframes you wrote.
