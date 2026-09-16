# Mapping hook animations to CSS

## Transition or animation

Choose per property and per edge, never once for the whole site; one element may need both (`animationName` for the mount animation, `transitionProperty` for the driver-changed properties, never one property in both).

- **Transition**: an outside trigger drives every step (press, toggle, expand, theme change). The driver becomes React state.
- **Animation**: the steps play themselves once started (mount, loop, `withSequence`, looping `withRepeat`).
- An element conditionally rendered on the driver (`if (!open) return null`), or a mount `useEffect` that writes a value different from the initial one, is an animation on that edge (`../animations/animations.md`, Mount animations).
- A CSS animation restarts on every mount of its element, and N CSS animations share a phase only when mounted in the same commit. A shared value owned by an ancestor that stays mounted kept its phase across child remounts and ticked as one clock for every reader. Such sites: Needs approval, proposing to move the animation to the ancestor.

## `with*`

| Source | CSS | Notes |
|---|---|---|
| `withTiming(v, { duration, easing })` | `transitionDuration` + `transitionTimingFunction`, or a two-keyframe animation | mechanism per the section above |
| `withDelay(ms, anim)` at the top level | `transitionDelay` or `animationDelay` | |
| `withDelay(ms, x)` inside `withSequence` or `withRepeat` | a hold before `x`: repeat the previous keyframe value at `x`'s offset, add `ms` to the total; inside `withRepeat` the hold starts every cycle, never `animationDelay` | not a leading delay |
| `withTiming(v, { duration: 0 })` inside `withSequence` | an instant step: the new value on the next offset, `'50%': { x: a }, '50.01%': { x: b }` (keyframes sharing an offset merge) | `steps(1, 'jump-end')` when every step is instant |
| `withRepeat(anim)`, count omitted | `animationIterationCount: 2` | `numberOfReps` defaults to 2, `animationIterationCount` to 1: always write the count |
| `withRepeat(anim, n)`, `n <= 0` | `animationIterationCount: 'infinite'` | |
| `withRepeat(withTiming(...), n, true)` | plus `animationDirection: 'alternate'`, exact only for a symmetric easing (`linear`, `inOut(f)`, `'ease-in-out'`); CSS mirrors the easing on the return leg, the hook plays it forward. Asymmetric easing: write the ping-pong as `0%`/`50%`/`100%` keyframes with the easing on both intervals, double the duration, halve the count | an even `n` rests at the start value, odd at the target. `withRepeat(withSequence(...), n, true)`: the hook ignores `reverse`, emit `'normal'` |
| `withSequence(a, b, c)` | one animation, `animationDuration` = sum of the children, each keyframe at its cumulative fraction | 100ms then 300ms: stops at `0%`, `25%`, `100%`. Finite animations need `animationFillMode: 'forwards'` unless the static style already equals the resting value |
| completion callback on `withTiming` or `withRepeat` | `onCSSTransitionEnd`/`onCSSAnimationEnd` plus the matching `Cancel`, as props | 4.6.0 only; `cb(true)` becomes End, `cb(false)` Cancel; a `withRepeat` callback fires after the last rep, so an infinite loop only ever reaches Cancel |

Two loops at different speeds need the array form (`animationName: [spin, pulse]`); never split one animation across entries. A sequence with a different easing per step sets `animationTimingFunction` on every keyframe whose outgoing interval needs it (`../animations/animations.md`, CSS Animations, Rules).

## Easing

Map the source expression, never a runtime value. Every row is exact:

| Source | CSS |
|---|---|
| `Easing.linear`, `Easing.poly(1)`, `Easing.inOut(Easing.linear)` | `'linear'` |
| `Easing.ease` | `'ease-in'`, i.e. `cubicBezier(0.42, 0, 1, 1)`; never CSS `'ease'`, which is a different curve |
| `Easing.quad`, `Easing.poly(2)` | `cubicBezier(1/3, 0, 2/3, 1/3)` |
| `Easing.cubic`, `Easing.poly(3)` | `cubicBezier(1/3, 0, 2/3, 0)` |
| `Easing.bezier(a, b, c, d)`, `Easing.bezierFn(...)` | `cubicBezier(a, b, c, d)` |
| `Easing.back(s)` | `cubicBezier(1/3, 0, 2/3, -s/3)`; default `s = 1.70158` gives `-0.56719` |
| `Easing.steps(n)`, `Easing.steps(n, true)` | `steps(n, 'jump-start')` |
| `Easing.steps(n, false)` | `steps(n, 'jump-end')` |
| `Easing.in(f)` | `f` unchanged |
| `Easing.out(f)` | reflect: `cubicBezier(x1, y1, x2, y2)` becomes `cubicBezier(1 - x2, 1 - y2, 1 - x1, 1 - y1)`; `out(Easing.ease)` is `'ease-out'` |

`Easing.inOut(f)` has no exact form except for `linear`. Substitute `'ease-in-out'` and state the max error when it is at or under 0.03; above that Keep:

| Source | Max error as `'ease-in-out'` |
|---|---|
| `inOut(quad)`, the `withTiming` default | 0.012 |
| `inOut(sin)` | 0.019 |
| `inOut(ease)` | 0.029 |
| `inOut(cubic)` | 0.084, Keep |
| `inOut(circle)` | 0.136, Keep |

## Reduced motion

Every `with*` without `ReduceMotion.Never` jumps under reduced motion, so an unguarded migration is a delta. Guard every animation and every transition with `useReducedMotion()`; only an `opacity` or color transition under 300ms may go unguarded, reported as `reduced-motion guard dropped`.

| Kind | Form |
|---|---|
| Transition | `transitionDuration: reduced ? 1 : D`, never `0` (`../animations/animations.md`, Reduced motion) |
| Animation | `animationDuration: reduced ? 1 : D` and `animationIterationCount: reduced ? 1 : N`; keep `animationName` so the fill mode and callbacks survive, and pick `animationFillMode` so it rests where the hook rests: `'forwards'` when the hook rests at the target, `'none'` (snapping back to the static style) when it rests at the start (`../animations/animations.md`, Reduced motion) |
| `ReduceMotion.Never` in the source | no guard |
| `ReduceMotion.Always` | the reduced form for everyone, no guard |

Where the hook rests: `withTiming(TO)`, `withSequence` and a non-reverse `withRepeat` rest at `TO`; `withRepeat(anim, n, true)` with `n <= 0` or even `n` rests at the start, unless `anim` is a `withSequence`, whose `reverse` is ignored, so it rests at the sequence's last value. Evaluate the style body at that driver value to know which end the fill mode must hold.

## Colors

`withTiming` and `interpolateColor` interpolate gamma-corrected; CSS lerps sRGB (`../animations/animations.md`, Shared rules). The gap peaks about a quarter in from the darker endpoint.

| Endpoints | Worst channel gap | Verdict |
|---|---|---|
| small swing between bright values (`#eee` to `#ccc`) | 1/255 | Migrate |
| about half the range (grey to mid green) | 36/255 | Migrate, state it |
| a channel crossing most of 0..255 (black to white, red to cyan) | 72/255 | Migrate only if the user chose all colors |

The user's answer maps onto the rows: all colors migrates all three; visually close pairs migrates the first two; keep colors leaves every color property on hooks.

## SVG

From 4.4.0 (web 4.5.0) CSS declarations go in `animatedProps`, never `style` (`../animations/svg-animations.md`). The driver becomes state and the attribute a plain prop: `r={grown ? 50 : 20}` beside `animatedProps={{ transitionProperty: 'r', transitionDuration: 300 }}`. Values from a separate `useAnimatedProps` hook you are not migrating can stay in the same `animatedProps` array. Geometry varies per component; check the component's row in [Animating SVG](https://docs.swmansion.com/react-native-reanimated/docs/guides/animating-svg) before converting (the four props the Keep table names throw despite being listed there).
