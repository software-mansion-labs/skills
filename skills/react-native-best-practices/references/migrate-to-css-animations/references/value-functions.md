# Value functions

Question 4 of the walk. A transition tweens a property in a straight line between two endpoints, so the style must be an affine function of one driver between them, and both endpoints must be the same kind of value. Anything else can still become a keyframe animation that samples the curve; that is a proposal the user approves.

| The style body is | Verdict |
|---|---|
| `a * driver + b`, or `interpolate` with no input stop strictly between the two endpoint values | continue |
| `interpolate` with a stop between the endpoints (`interpolate(p, [0, 0.5, 1], [0, 100, 0])` passes through 100 where a transition from 0 to 0 stays put) | a keyframe animation with a keyframe at every stop, each keyframe at the time the driver's easing reaches that stop (`E^-1(stop)`, so a stop at 0.5 under an `inOut` easing sits at `50%`) and with the driver's easing applied per interval: exact when the driver's easing is linear, a bezier row, or an `inOut` curve with a stop at 0.5; otherwise Needs approval with the error stated. For a state-driven site the animation replaces the transition, so say that a retrigger mid-flight restarts instead of retargeting |
| any other function of the driver (`Math.sin(driver)` is one example) | sample it every 5 to 10 percent into keyframes and show the result: Needs approval with the step and the error stated, and the same retrigger note for a state-driven site |
| two drivers feeding entries of one compound property, written together with one config (`translateX: x.value, translateY: y.value` from one handler, `shadowOffset`, `boxShadow`, `filter`) | one state object holds both entries: continue |
| two drivers combined into one value (`a.value * b.value`, `base.value + offset.value`), or written with different configs | Keep on shared values: no single state value to transition |
| a number and a keyword (`300` to `'auto'`) | Keep on shared values: nothing tweens, the value switches at the midpoint (`../animations/animations.md`, Shared rules) |
| `0` and a percentage (`0` to `'100%'`) | continue: `0` is `0%`, CSS tweens |
| a nonzero number and a percentage (`50` to `'100%'`) | Needs approval: CSS resolves the percentage against the parent and tweens between the resolved lengths (it switches at the midpoint only when the parent cannot be measured), while the shared value kept the start value's unit and animated the number, so `50` to `'100%'` ran from 50 to 100 points; the end state changes |
| `interpolateColor` in `'HSV'` or `'LAB'` | Keep on shared values: CSS lerps sRGB |
| a ternary on the driver whose result is not wrapped in a `with*` | a step, not a tween: render it conditionally and leave it out of `transitionProperty`. That moves the step from the threshold crossing to the state change, identical for a boolean driver or a threshold at the start value; any other threshold: Needs approval with the timing shift (keywords: `references/properties.md`) |

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
