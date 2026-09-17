# Reduced motion

`with*` animations follow the system Reduce Motion setting on their own: every `with*` without `ReduceMotion.Never` jumps to its target when the setting is on. CSS transitions and animations ignore the setting (`../animations/animations.md`, Reduced motion), so a migration silently drops that behavior. Step 1 asks the user once whether to keep it or drop it; this file is the keep form.

`useReducedMotion()` returns the system setting read once when the app starts and ignores `<ReducedMotionConfig>`; `with*` read the config live.

| Source | Form |
|---|---|
| a transition, reduced motion kept | `transitionDuration: reduced ? 1 : D`, never `0`, and `transitionDelay: reduced ? 0 : ms` when the source had a `withDelay` (the shared value dropped the delay under reduced motion) |
| an animation, reduced motion kept | `animationDuration: reduced ? 1 : D`, `animationIterationCount: reduced ? 1 : N` and `animationDelay: reduced ? 0 : ms`; keep `animationName` so the fill mode and callbacks survive, and pick `animationFillMode` so it rests where the shared value rested: `'forwards'` when it rested at the target, `'none'` (snapping back to the static style) when it rested at the start, `reduced ? 'none' : 'forwards'` when the two rests differ (a reversed even-count repeat of a sequence) |
| `ReduceMotion.Never` in the source, or `<ReducedMotionConfig mode={ReduceMotion.Never}>` | no guard, whatever the user answered |
| `ReduceMotion.Always`, or `<ReducedMotionConfig mode={ReduceMotion.Always}>` | the reduced form for everyone, no guard |
| `<ReducedMotionConfig mode={ReduceMotion.System}>` | the same as no config |
| `<ReducedMotionConfig>` whose `mode` changes at runtime | Keep on shared values: CSS has no live reduced-motion source |
| reduced motion dropped | no guard; say once in the report that users with Reduce Motion on now see the animations |

Where the shared value rested: `withTiming(TO)`, `withSequence` and a non-reverse `withRepeat` rest at `TO`; `withRepeat(anim, n, true)` with odd `n` runs once under reduced motion and rests at `TO`, with `n <= 0` or even `n` it rests at the start whatever `anim` is (a reversed sequence included, which otherwise rests at its last value). Evaluate the style body at that driver value to know which end the fill mode must hold.

## Trap: the resting value

```tsx
const progress = useSharedValue(0);
useEffect(() => { progress.value = withRepeat(withTiming(1, { duration: 600 }), -1, true); }, []);
const style = useAnimatedStyle(() => ({ transform: [{ scaleY: interpolate(progress.value, [0, 1], [0.35, 1]) }] }));
```

Reduced motion never advances `progress`, so the shared value rests at `interpolate(0, ...)` = `0.35`, the start. The shortened animation must rest there too: static style `scaleY: 0.35` and `animationFillMode: 'none'`, so the 1ms run snaps back to it. `'forwards'` would park the element at the `to` keyframe, nearly three times too tall. Evaluate the style body at the driver's resting value, never read it off the keyframes you wrote.
