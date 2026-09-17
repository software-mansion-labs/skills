# Springs, decay and clamp

Question 2 of the walk.

`withSpring` and `withDecay`, anywhere in the composition (inside `withSequence`, `withDelay` or `withRepeat` too): Keep on shared values. CSS has no spring and no velocity-driven timeline, and a spring must never be hand-sampled into `linear()`: its curve depends on the distance and velocity of each run.

`withClamp({ min, max }, anim)` clamps the inner value every frame, so it is a no-op only when every value stays inside `[min, max]`:

| Inner animation | Verdict |
|---|---|
| `withSpring` or `withDecay` | Keep on shared values |
| timing animations whose endpoints lie inside the bounds and whose every frame stays inside too: an easing within 0..1 (`Easing.bounce` included, it never leaves 0..1), or an overshooting easing whose excursion stays inside the bounds | the clamp never triggers: drop it and continue the walk with the inner animation |
| timing animations with an endpoint outside the bounds, or an overshooting easing whose excursion crosses a bound | the clamp triggers: Keep on shared values, CSS cannot clamp |

The excursions: `Easing.back(s)` undershoots the start by `4 s^3 / (27 (s + 1)^2)` of the distance, 0.100 at the default `s = 1.70158`; `Easing.elastic()` overshoots the end by 0.066 of the distance; a cubic bezier by however far its y control points take the curve outside 0..1. So `withClamp({ min: 0, max: 1 }, withTiming(1, { easing: Easing.back() }))` from 0 triggers (the curve reaches -0.100), while the same clamp around a run from 0.2 to 0.8 never does.

A clamp spelled differently is the same question: `Extrapolation.CLAMP` in an `interpolate` (the default), `clamp()`, `Math.min`/`Math.max` around the driver, or `interpolateColor` (which clamps its input) hold the value at a bound whenever the driver's easing overshoots the input range. With an overshooting easing they trigger: Keep on shared values; with an easing within 0..1 they never trigger and drop.

## Example

```tsx
// Before: the clamp never triggers (0..1 inside 0..1, linear easing)
width.value = withClamp({ min: 0, max: 1 }, withTiming(open ? 1 : 0, { duration: 200, easing: Easing.linear }));
```

```tsx
// After: the plain transition
{ width: open ? 1 : 0, transitionProperty: 'width', transitionDuration: reduced ? 1 : 200, transitionTimingFunction: 'linear' }
```
