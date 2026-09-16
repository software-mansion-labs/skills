# Mapping hook animations to CSS

## Transition or animation

Choose per property and per edge, never once for the whole site; one element may need both (`animationName` for the mount animation, `transitionProperty` for the driver-changed properties, never one property in both).

- **Transition**: an outside trigger drives every step (press, toggle, expand, theme change). The driver becomes React state.
- **Animation**: the steps play themselves once started (mount, loop, `withSequence`, looping `withRepeat`).
- An element conditionally rendered on the driver (`if (!open) return null`), or a mount `useEffect` that writes a value different from the initial one, is an animation on that edge (`../animations/animations.md`, Mount animations).
- A CSS animation restarts on every mount of its element, and N CSS animations share a phase only when mounted in the same commit. A shared value owned by an ancestor that stays mounted kept its phase across child remounts and ticked as one clock for every reader. Such sites: Needs approval, proposing to move the animation to the ancestor.

## Properties

CSS animates a property only when Reanimated has an interpolator for it. Check the [supported properties](https://docs.swmansion.com/react-native-reanimated/docs/guides/supported-properties) table, never memory; it describes the latest release, so on an older installed version also check the feature floors in `../animations/animations.md` (filter 4.2.0, SVG 4.4.0, pseudo-selectors 4.5.0, callbacks 4.6.0). React Native does not render the property on a platform either: fine, it was inert there (iOS `shadow*` on Android, `elevation` on iOS). React Native renders it and CSS cannot animate it: Keep on hooks. Discrete (keyword) properties change at the transition midpoint only with `transitionBehavior: 'allow-discrete'` (`../animations/animations.md`, Discrete properties); propose it as Needs approval, and Keep when the flip breaks layout.

## Value functions

A transition tweens the property between two endpoints, so the hook's style must be an affine function of the driver (`a * driver + b`, or `interpolate` with fixed stops evaluated at the endpoints). Anything else (trigonometry, `Math.pow` of the driver, modulo, `interpolate` past its stops without `Extrapolation.CLAMP`) gives different intermediate frames: Keep on hooks for a transition. An animation may instead sample the function into keyframes (every 5 to 10 percent, linear between them); that is an approximation, so Needs approval with the sampling step stated. `interpolateColor` in `'HSV'` or `'LAB'`: Keep on hooks, CSS lerps sRGB. A ternary on the driver whose result is not wrapped in a `with*` is a step, not a tween: render it conditionally and leave it out of `transitionProperty`.

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
| `withClamp({ min, max }, withTiming(v))` | drop the clamp when the start value and `v` both lie inside `[min, max]` (the clamp never triggers) | an endpoint outside the bounds: Keep on hooks, CSS cannot clamp |
| `withSpring`, `withDecay`, `withClamp` around either | none | Keep on hooks: no CSS spring or velocity-driven timeline; never hand-sample a spring into `linear()` |

Two loops at different speeds need the array form (`animationName: [spin, pulse]`); never split one animation across entries. A sequence with a different easing per step sets `animationTimingFunction` on every keyframe whose outgoing interval needs it (`../animations/animations.md`, CSS Animations, Rules).

## Easing

Map the `Easing.*` expression written in the source to a CSS timing function; the `Easing` object itself is not accepted by CSS. Exact rows:

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

Every other curve (`Easing.inOut(f)`, `poly(n)` for `n` above 3, `sin`, `circle`, `exp`, `elastic`, `bounce`, and the `in`/`out`/`inOut` forms of those) has no cubic-bezier form. Two ways to migrate it, both to be offered as Needs approval, exact first:

- Exact: sample the function into `linear()` (`react-native-reanimated`, available since 4.0.0): `linear(0, f(1/n), f(2/n), ..., 1)` with evenly spaced stops. 20 stops hold a smooth curve within 0.005; `bounce` and `elastic` need about 50. Verbose, so state the stop count.
- Approximate: the nearest named curve or cubic-bezier with its max error, which the user may prefer for readability. `'ease-in-out'` stands in for `inOut(quad)` (the `withTiming` default) with a max error of 0.012, for `inOut(sin)` with 0.019, for `inOut(ease)` with 0.029; `inOut(cubic)` (0.084) and `inOut(circle)` (0.136) are visibly different.

Ask which the user wants when the first such site comes up, then apply the answer to every site. Do not sample a spring: its curve depends on the distance and velocity of each run.

## Imperative control

| Source | CSS | Verdict |
|---|---|---|
| `cancelAnimation(sv)` in an unmount cleanup | nothing; CSS cleans up on unmount | Migrate |
| `cancelAnimation(sv)` then `sv.value = withTiming(other)` | a new target; the transition retargets from the current value | Migrate |
| `cancelAnimation(sv)` then `sv.value = withTiming(sameTarget)` | the hook restarts at full duration; a transition to a target already in flight changes nothing | Needs approval |
| pausing and resuming (`animationPlayState` has no hook twin, the hook cancels and re-animates from the current value) | `animationPlayState: 'paused'` on an animation | Needs approval |
| reversing mid-flight from code (`cancelAnimation` then `withTiming` back to the start) | a transition back to the start value, which takes the shortened return leg (`../animations/animations.md`, CSS Transitions, Rules) | Needs approval |
| restarting a finished animation (`sv.value = 0; sv.value = withTiming(1)`) | remount the element or change its `key` | Needs approval: a visible structural change |

## Callbacks

`with*` completion callbacks map onto the `onCSS*` props from 4.6.0 (`../animations/animations.md`, Callbacks); below 4.6.0 an observable callback keeps the site on hooks. A log-only callback is dropped.

| Source | CSS |
|---|---|
| `withTiming(v, cfg, cb)` on a transitioned property | `onCSSTransitionEnd` for `cb(true)`, `onCSSTransitionCancel` for `cb(false)`; the payload carries `propertyName`, so one handler serves several properties |
| `withTiming(v, cfg, cb)` on a keyframe animation, `withSequence(...)` completion | `onCSSAnimationEnd` for `cb(true)`, `onCSSAnimationCancel` for `cb(false)` |
| `withRepeat(withTiming(v, cfg, innerCb), n)` | `innerCb` fires after every repetition: `onCSSAnimationIteration` (not after the last one, which is `onCSSAnimationEnd`) |
| `withRepeat(anim, n, reverse, outerCb)` | `outerCb` fires once after the last repetition: `onCSSAnimationEnd`; with `n <= 0` it never fires, and the CSS animation only ever reaches `onCSSAnimationCancel` |
| a callback on a step inside `withSequence` | no per-keyframe event; Needs approval, proposing a timer at the step's offset, else Keep on hooks |
| a callback that relies on `finished: true` when the target already equals the current value | CSS fires nothing when nothing changes; call the handler directly in that branch (`if (next === current) onDone(); else setValue(next)`) and say so in the row |

Callbacks fire for pseudo-selector driven transitions too, and a transition removed by a zero effective duration (4.3.0+) fires nothing.

## Reversal

A CSS transition reversed mid-flight takes a shortened return leg (`../animations/animations.md`, CSS Transitions, Rules); a hook plays the full duration back. Sites where the user can flip the driver inside the duration (press in/out, a toggle, expand/collapse) therefore change behavior: Needs approval, showing the difference in milliseconds for the site's duration. Offer the version-appropriate shape:

- press feedback on the pressed element: from 4.5.0 `:active` in the style (4.6.0 for SVG elements in `animatedProps`), see `../animations/css-pseudo-selectors.md`; `:active` shortens the return leg too;
- press feedback below 4.5.0, or on a descendant or ancestor of the pressed element: a `Pressable` render prop styling the child, or `useState` from `onPressIn`/`onPressOut` when the `Pressable` itself carries the style (`../animations/animations.md`, Simple gesture feedback);
- other toggles: the transition, with the shortened return stated.

A driver behind a timer longer than the duration, or a rare event (rotation, navigation, network): Migrate and state the difference. Animations have no reversal behavior.

## Reduced motion

`useReducedMotion()` returns the system setting read once when the app starts and ignores `<ReducedMotionConfig>`; `with*` read the config live. Every `with*` without `ReduceMotion.Never` jumps under reduced motion, so an unguarded migration is a delta. Guard every animation and every transition with `useReducedMotion()`; only an `opacity` or color transition under 300ms may go unguarded, reported as `reduced-motion guard dropped`.

| Kind | Form |
|---|---|
| Transition | `transitionDuration: reduced ? 1 : D`, never `0` (`../animations/animations.md`, Reduced motion) |
| Animation | `animationDuration: reduced ? 1 : D` and `animationIterationCount: reduced ? 1 : N`; keep `animationName` so the fill mode and callbacks survive, and pick `animationFillMode` so it rests where the hook rests: `'forwards'` when the hook rests at the target, `'none'` (snapping back to the static style) when it rests at the start (`../animations/animations.md`, Reduced motion) |
| `ReduceMotion.Never` in the source, or `<ReducedMotionConfig mode={ReduceMotion.Never}>` | no guard |
| `ReduceMotion.Always`, or `<ReducedMotionConfig mode={ReduceMotion.Always}>` | the reduced form for everyone, no guard |
| `<ReducedMotionConfig>` whose `mode` changes at runtime | Keep on hooks |

Where the hook rests: `withTiming(TO)`, `withSequence` and a non-reverse `withRepeat` rest at `TO`; `withRepeat(anim, n, true)` with `n <= 0` or even `n` rests at the start, unless `anim` is a `withSequence`, whose `reverse` is ignored, so it rests at the sequence's last value. Evaluate the style body at that driver value to know which end the fill mode must hold.

## Colors

`withTiming` and `interpolateColor` interpolate gamma-corrected; CSS lerps sRGB (`../animations/animations.md`, Shared rules). The gap peaks about a quarter in from the darker endpoint.

| Endpoints | Worst channel gap | Verdict |
|---|---|---|
| small swing between bright values (`#eee` to `#ccc`) | 1/255 | Migrate |
| about half the range (grey to mid green) | 36/255 | Migrate, state it |
| a channel crossing most of 0..255 (black to white, red to cyan) | 72/255 | Needs approval |

## SVG

From 4.4.0 (web 4.5.0; on 4.1.0 to 4.3.x only with the `EXPERIMENTAL_CSS_ANIMATIONS_FOR_SVG_COMPONENTS` static flag, a native rebuild) CSS declarations go in `animatedProps`, never `style` (`../animations/svg-animations.md`). The driver becomes state and the attribute a plain prop: `r={grown ? 50 : 20}` beside `animatedProps={{ transitionProperty: 'r', transitionDuration: 300 }}`. Values from a separate `useAnimatedProps` hook you are not migrating can stay in the same `animatedProps` array. Which attributes animate varies per component; check the component's row in [Animating SVG](https://docs.swmansion.com/react-native-reanimated/docs/guides/animating-svg) before converting. Keep on hooks: SVG below the floors above, SVG `transform` arrays, `mask`, `filter`, `marker*` and `pointerEvents` (listed in that table but they throw `No interpolator factory found`), and `fill`/`stroke` written as `currentColor` or `url(#...)`.
