# Value functions

Question 4 of the walk. A transition tweens a property in a straight line between two endpoints, so the style must be an affine function of one driver between them, and both endpoints must be the same kind of value. Anything else can still become a keyframe animation that samples the curve; that is a proposal the user approves.

| The style body is | Verdict |
|---|---|
| `a * driver + b`, or `interpolate` with no input stop strictly between the two endpoint values | continue, unless an explicit `Extrapolation.CLAMP`, `clamp()`, `Math.min`/`Math.max` or `interpolateColor` meets an overshooting easing or a driver that leaves the input range: that clamp triggers, `references/springs-and-clamp.md` |
| `interpolate` with a stop between the endpoints (`interpolate(p, [0, 0.5, 1], [0, 100, 0])` passes through 100 where a transition from 0 to 0 stays put) | a keyframe animation with a keyframe at every stop, each keyframe at the time the driver's easing reaches that stop (`E^-1(stop)`, so a stop at 0.5 under an `inOut` easing sits at `50%`) and, as its `animationTimingFunction`, the driver's easing restricted to the outgoing interval and renormalized to 0..1: exact for `linear`, for any exact bezier row, and for `inOut(f)` with a stop at 0.5 when `f` has an exact row (the halves are `f` then `out(f)`: the `withTiming` default gives `cubicBezier(1/3, 0, 2/3, 1/3)` then `cubicBezier(1/3, 2/3, 2/3, 1)`); otherwise Needs approval with the error stated. For a state-driven toggle emit one rule per direction, chosen by state and attached only once the state has changed (second example), with `animationFillMode: 'forwards'` and the static style at the resting value, and say that a flip mid-flight restarts the incoming rule from its first keyframe instead of retargeting |
| any other function of the driver (`Math.sin(driver)` is one example) | sample it every 5 to 10 percent into keyframes and show the result: Needs approval with the step and the error stated, and the same retrigger note for a state-driven site |
| two drivers feeding entries of one compound property, written together with one config (`translateX: x.value, translateY: y.value` from one handler, `shadowOffset`, `boxShadow`, `filter`) | one state object holds both entries: continue |
| entries of one compound property written with different configs (`scale` at 120ms, `translateY` at 300ms) | Needs approval proposing one config and naming the entry whose timing changes, or Keep on shared values |
| two drivers combined into one value (`a.value * b.value`, `base.value + offset.value`) | Keep on shared values: no single state value to transition |
| a number and a keyword (`300` to `'auto'`) | Keep on shared values: nothing tweens, the value jumps to the target when the transition starts (at the midpoint only with `transitionBehavior: 'allow-discrete'`; `../animations/animations.md`, Shared rules) |
| a percentage string at both ends (`'0%'` to `'100%'`) | continue: both tween |
| a number and a percentage (`50` to `'100%'`, `'50%'` to `300`) | Needs approval: CSS resolves the percentage (against the parent; `translateX`/`translateY`, border radii, gaps and transform origin against the view itself) and tweens between the resolved lengths (it switches at the midpoint only when the parent cannot be measured). The shared value did not tween: from a number start it produced `NaN` every frame (the property rendered unset) and jumped to the target at the end, so the migration fixes the original and only the frames differ; from a percentage start it kept the `%` and animated the number, so `'50%'` to `300` ended at `'300%'` and the end state changes. Say which |
| `interpolateColor` in `'HSV'` or `'LAB'` | Keep on shared values: CSS lerps sRGB |
| a ternary on the driver whose result is not wrapped in a `with*` | a step, not a tween: render it conditionally and leave it out of `transitionProperty`. That moves the step from the threshold crossing to the state change, identical for a boolean driver or a threshold at the start value; any other threshold: Needs approval with the timing shift for a numeric step; a keyword follows `references/properties.md`, which keeps other thresholds on shared values |

## Example

```tsx
// Before: affine in the driver, so a transition between 100 and 300
const style = useAnimatedStyle(() => ({ height: 100 + progress.value * 200 }));
useEffect(() => { progress.value = withTiming(expanded ? 1 : 0, { duration: 200 }); }, [expanded]);
```

```tsx
// After
{ height: expanded ? 300 : 100, transitionProperty: 'height', transitionDuration: reduced ? 1 : 200, transitionTimingFunction: 'ease-in-out' }
```

A toggle the user can flip back inside 200ms, so question 8 notes the shortened return: reversed 80ms in, the shared value took the 200ms the effect gave it and CSS takes about 66ms with `'ease-in-out'`.

## Example: a multi-stop toggle

```tsx
// Before: opening passes through 200 on its way to 100, the withTiming default easing
const style = useAnimatedStyle(() => ({ height: interpolate(progress.value, [0, 0.5, 1], [0, 200, 100]) }));
useEffect(() => { progress.value = withTiming(open ? 1 : 0, { duration: 400 }); }, [open]);
```

```tsx
// After: one rule per direction, module scope; the stop at 0.5 sits at 50% under inOut(quad),
// each interval gets its half of the easing
const openRule = css.keyframes({
  from: { height: 0, animationTimingFunction: cubicBezier(1 / 3, 0, 2 / 3, 1 / 3) },
  '50%': { height: 200, animationTimingFunction: cubicBezier(1 / 3, 2 / 3, 2 / 3, 1) },
  to: { height: 100 },
});
const closeRule = css.keyframes({
  from: { height: 100, animationTimingFunction: cubicBezier(1 / 3, 0, 2 / 3, 1 / 3) },
  '50%': { height: 200, animationTimingFunction: cubicBezier(1 / 3, 2 / 3, 2 / 3, 1) },
  to: { height: 0 },
});
// in the component: no animation until the first toggle, so mount stays static like the shared value did
const [rule, setRule] = useState<ReturnType<typeof css.keyframes>>();
const toggle = () => { setRule(open ? closeRule : openRule); setOpen(!open); };
// on the element: { height: open ? 100 : 0, ...(rule && { animationName: rule, animationDuration: 400, animationFillMode: 'forwards' }) }
```

Needs approval row: a toggle flipped while a run is in flight restarts the incoming rule from its first keyframe (100 or 0) where the shared value turned back from the current height; the keyframes and easings are exact.
