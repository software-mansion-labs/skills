# Springs, decay and clamp

Question 2 of the walk.

`withSpring` and `withDecay`, anywhere in the composition (inside `withSequence`, `withDelay` or `withRepeat` too): Keep on shared values. CSS has no spring and no velocity-driven timeline, and a spring must never be hand-sampled into `linear()`: its curve depends on the distance and velocity of each run.

`withClamp({ min, max }, anim)` clamps the inner value every frame, so it is a no-op only when every value stays inside `[min, max]`:

| Inner animation | Verdict |
|---|---|
| `withSpring` or `withDecay` | Keep on shared values |
| timing animations whose endpoints lie inside the bounds and whose easing stays within 0..1 (not `back` or `elastic`, not a bezier with y outside 0..1) | the clamp never triggers: drop it and continue the walk with the inner animation |
| timing animations with an endpoint outside the bounds, or an overshooting easing | the clamp triggers: Keep on shared values, CSS cannot clamp |

## Example

```tsx
// Before: the clamp never triggers (0..1 inside 0..1, linear easing)
width.value = withClamp({ min: 0, max: 1 }, withTiming(open ? 1 : 0, { duration: 200, easing: Easing.linear }));
```

```tsx
// After: the plain transition
{ width: open ? 1 : 0, transitionProperty: 'width', transitionDuration: reduced ? 1 : 200, transitionTimingFunction: 'linear' }
```
