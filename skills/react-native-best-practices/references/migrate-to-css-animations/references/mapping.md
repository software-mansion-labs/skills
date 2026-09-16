# Mapping hook animations to CSS

## Transition or animation

Choose per property and per edge, never once for the whole site; one element may need both (`animationName` for the mount animation, `transitionProperty` for the driver-changed properties, never one property in both).

- **Transition**: an outside trigger drives every step (press, toggle, expand, theme change). The driver becomes React state.
- **Animation**: the steps play themselves once started (mount, loop, `withSequence`, looping `withRepeat`).
- An element conditionally rendered on the driver (`if (!open) return null`), or a mount `useEffect` that writes a value different from the initial one, is an animation on that edge (`../animations/animations.md`, Mount animations).
- A CSS animation restarts on every mount of its element, and N CSS animations share a phase only when mounted in the same commit. A shared value owned by an ancestor that stays mounted kept its phase across child remounts and ticked as one clock for every reader. Such sites: Needs approval, proposing to move the animation to the ancestor.

## Properties

CSS animates a property only when Reanimated has an interpolator for it. Check the [supported properties](https://docs.swmansion.com/react-native-reanimated/docs/guides/supported-properties) table, never memory; it describes the latest release, so on an older installed version also check the feature floors in `../animations/animations.md` (filter 4.2.0, SVG 4.4.0, pseudo-selectors 4.5.0, callbacks 4.6.0). If React Native never renders the property on a platform (`shadowOffset`, `shadowOpacity` and `shadowRadius` on Android, `elevation` on iOS), CSS not animating it there changes nothing. If React Native renders it and CSS cannot animate it: Keep on hooks. Discrete (keyword) properties change at the transition midpoint only with `transitionBehavior: 'allow-discrete'` (`../animations/animations.md`, Discrete properties). Route by where the hook flipped the value: a boolean or state driver flipped at the state change, so render the keyword conditionally and leave it out of `transitionProperty`; a `> 0.5` threshold on a 0..1 numeric driver matches the midpoint, so propose `allow-discrete` as Needs approval; any other threshold, or a flip that breaks layout: Keep on hooks.

## Value functions

A transition tweens the property between two endpoints, so between those endpoints the hook's style must be an affine function of the driver: `a * driver + b`, or `interpolate` with no input stop strictly between the two endpoint values (`interpolate(p, [0, 0.5, 1], [0, 100, 0])` passes through 100 where a transition from 0 to 0 stays put). Anything else gives different intermediate frames (`Math.sin(driver)` is one example): Keep on hooks for a transition. An animation may instead put a keyframe at every `interpolate` stop, exact when the driver's own easing is linear, or sample any other function every 5 to 10 percent; with an eased driver both are approximations, so Needs approval with the step and the error stated. `interpolateColor` in `'HSV'` or `'LAB'`: Keep on hooks, CSS lerps sRGB. A ternary on the driver whose result is not wrapped in a `with*` is a step, not a tween: render it conditionally and leave it out of `transitionProperty`. That moves the step from the threshold crossing to the state change, identical for a boolean driver or a threshold at the start value; any other threshold: Needs approval with the timing shift (keywords: see Properties above).

## `with*`

| Source | CSS | Notes |
|---|---|---|
| `withTiming(v, { duration, easing })` | `transitionDuration` + `transitionTimingFunction`, or a two-keyframe animation | mechanism per the section above |
| `withDelay(ms, anim)` at the top level | `transitionDelay` or `animationDelay` | retargeting mid-flight differs: the hook keeps the running animation going during the delay, a delayed transition cancels it and holds the current value; a site that retargets while running: Needs approval |
| `withDelay(ms, x)` inside `withSequence` or `withRepeat` | a hold before `x`: repeat the previous keyframe value at `x`'s offset, add `ms` to the total; inside `withRepeat` the hold starts every cycle, never `animationDelay` | not a leading delay |
| `withTiming(v, { duration: 0 })` inside `withSequence` | an instant step: the new value on the next offset, `'50%': { x: a }, '50.01%': { x: b }` (keyframes sharing an offset merge) | a sequence of holds ending in instant changes: one keyframe per hold with `animationTimingFunction: steps(1, 'jump-end')` on it |
| `withRepeat(anim)`, count omitted | `animationIterationCount: 2` | `numberOfReps` defaults to 2, `animationIterationCount` to 1: always write the count |
| `withRepeat(anim, n)`, `n <= 0` | `animationIterationCount: 'infinite'` | |
| `withRepeat(withTiming(...), n, true)` | plus `animationDirection: 'alternate'`, exact only for a symmetric easing (`linear`, `inOut(f)`, `'ease-in-out'`); CSS mirrors the easing on the return leg, the hook plays it forward. Asymmetric easing: write the ping-pong as `0%`/`50%`/`100%` keyframes with the easing on both intervals, double the duration, halve the count (a fractional count is accepted and ends on the `50%` keyframe; see the fill-mode caveat in `../animations/animations.md`) | an even `n` rests at the start value, odd at the target |
| `withRepeat(withSequence(...), n, true)` | `animationDirection: 'normal'` when the sequence ends where it starts; otherwise Needs approval, with keyframes for cycles 2+ that start at the sequence's last value | the hook does not reverse a sequence, but each later cycle continues from the sequence's last value instead of jumping back to the start |
| `withRepeat(withDelay(ms, withTiming(v)), n, true)` | none | the hook holds `v` from the second cycle on (the delay ignores the reversed target); Needs approval, or Keep on hooks |
| `withSequence(a, b, c)` | one animation, `animationDuration` = sum of the children, each keyframe at its cumulative fraction | 100ms then 300ms: stops at `0%`, `25%`, `100%`. Finite animations need `animationFillMode: 'forwards'` unless the static style already equals the resting value |
| `withClamp({ min, max }, anim)` around timing animations | drop the clamp when every value stays inside `[min, max]`: both endpoints inside and an easing within 0..1 (not `back` or `elastic`, not a bezier with y outside 0..1) | otherwise the clamp triggers: Keep on hooks, CSS cannot clamp |
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

Never emit CSS `'ease'` and never leave the timing function out: on 4.0.0 to 4.3.x the native `'ease'`, which is also the default, is `cubicBezier(0.25, 0.1, 0.25, 0.1)`, a wrong curve (up to 0.40 off the real one); write `cubicBezier(0.25, 0.1, 0.25, 1)` when that curve is wanted.

Every other curve (`Easing.inOut(f)`, `poly(n)` for `n` above 3, `sin`, `circle`, `exp`, `elastic`, `bounce`, and the `in`/`out`/`inOut` forms of those) has no cubic-bezier form. Two ways to migrate it; the single question after the inventory picks one for the whole run:

- Sampled: `linear()` (`react-native-reanimated`, available since 4.0.0) with the function's values at evenly spaced stops, `linear(0, f(1/n), f(2/n), ..., 1)` (n intervals, n + 1 points). Compute the max error for the count you emit and state it. 20 intervals hold `inOut(quad|cubic|sin)`, `sin`, `poly(4)`, `back` and `elastic` within 0.006; `exp` and `poly(5)` need 50; `bounce` is 0.078 at 20 and 0.024 at 50 with even stops but 0.005 at 20 once stops sit on its three kinks (`t = 1/2.75, 2/2.75, 2.5/2.75`; `linear()` accepts `[value, 'pct%']` positions); `circle` has a vertical tangent at the end (0.080 at 20, 0.035 at 100 even stops), so add stops at 99%, 99.5% and 99.9%.
- Approximate: the nearest named curve or cubic-bezier with its max error, which the user may prefer for readability. `'ease-in-out'` stands in for `inOut(quad)` (the `withTiming` default) with a max error of 0.012, for `inOut(sin)` with 0.019, for `inOut(ease)` with 0.029; `inOut(cubic)` (0.084) and `inOut(circle)` (0.136) are visibly different.

Either way the site's row records the substitution and its max error. Do not sample a spring: its curve depends on the distance and velocity of each run.

## Imperative control

| Source | CSS | Verdict |
|---|---|---|
| `cancelAnimation(sv)` in an unmount cleanup | nothing; CSS cleans up on unmount | Migrate |
| `cancelAnimation(sv)` then `sv.value = withTiming(other)` | a new target; the transition retargets from the current value | Migrate |
| `cancelAnimation(sv)` then `sv.value = withTiming(sameTarget)` | the hook restarts at full duration; a transition to a target already in flight changes nothing | Needs approval |
| pausing and resuming (the hook has no pause; sites emulate it with `cancelAnimation` and a later `with*` from the current value) | `animationPlayState: 'paused'` / `'running'` on an animation | Needs approval |
| reversing mid-flight from code (`cancelAnimation` then `withTiming` back to the start) | a transition back to the start value, which takes the shortened return leg (`../animations/animations.md`, CSS Transitions, Rules) | Needs approval |
| stopping a running loop and tweening to rest (`cancelAnimation(pulse); pulse.value = withTiming(1)` when loading ends) | removing `animationName` snaps to the static style (`../animations/animations.md`, Mount animations); there is no tween from the current animated value | Needs approval, stating the snap, or Keep on hooks |
| restarting a finished animation (`sv.value = 0; sv.value = withTiming(1)`) | a new keyframes rule restarts the animation: `animationName: useMemo(() => css.keyframes(frames), [replayCount])` (`../animations/animations.md`, Defining keyframes); a remount or `key` change also works but changes the element tree | Needs approval |

## Callbacks

`with*` completion callbacks map onto the `onCSS*` props from 4.6.0 (`../animations/animations.md`, Callbacks); below 4.6.0 an observable callback keeps the site on hooks. A log-only callback is dropped.

| Source | CSS |
|---|---|
| `withTiming(v, cfg, cb)` on a transitioned property | `onCSSTransitionEnd` for `cb(true)`, `onCSSTransitionCancel` for `cb(false)`; the payload carries `propertyName`, so one handler serves several properties |
| `withTiming(v, cfg, cb)` on a keyframe animation, `withSequence(...)` completion | `onCSSAnimationEnd` for `cb(true)`, `onCSSAnimationCancel` for `cb(false)` |
| `withRepeat(withTiming(v, cfg, innerCb), n)` | `innerCb` fires after every repetition: `onCSSAnimationIteration` (not after the last one, which is `onCSSAnimationEnd`) |
| `withRepeat(anim, n, reverse, outerCb)` | `outerCb` fires once after the last repetition: `onCSSAnimationEnd`; with `n <= 0` it fires `false` on cancel only (`onCSSAnimationCancel`), and `true` under reduced motion, which the `animationIterationCount: reduced ? 1 : 'infinite'` form also reaches as `onCSSAnimationEnd` |
| a callback on a step inside `withSequence` | no per-keyframe event; Needs approval, proposing a timer at the step's offset, else Keep on hooks |
| a callback that relies on `finished: true` when the target already equals the current value | CSS fires nothing when nothing changes; call the handler directly in that branch (`if (next === current) onDone(); else setValue(next)`) and say so in the row |
| the same target written again mid-flight | the hook fired the first callback with `false` and the second with `true`; CSS leaves the running transition alone and fires one End: Needs approval when the site counts on the `false` call |
| a callback that writes another shared value | the `onCSS*` handler runs on the JS thread: set state when the chained site migrates too, else assign the shared value from JS (one frame later) |

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
| Transition | `transitionDuration: reduced ? 1 : D`, never `0`, and `transitionDelay: reduced ? 0 : ms` when the source had a `withDelay` (the hook drops the delay under reduced motion) (`../animations/animations.md`, Reduced motion) |
| Animation | `animationDuration: reduced ? 1 : D`, `animationIterationCount: reduced ? 1 : N` and `animationDelay: reduced ? 0 : ms`; keep `animationName` so the fill mode and callbacks survive, and pick `animationFillMode` so it rests where the hook rests: `'forwards'` when the hook rests at the target, `'none'` (snapping back to the static style) when it rests at the start, `reduced ? 'none' : 'forwards'` when the two rests differ (a reversed even-count repeat of a sequence) (`../animations/animations.md`, Reduced motion) |
| `ReduceMotion.Never` in the source, or `<ReducedMotionConfig mode={ReduceMotion.Never}>` | no guard |
| `ReduceMotion.Always`, or `<ReducedMotionConfig mode={ReduceMotion.Always}>` | the reduced form for everyone, no guard |
| `<ReducedMotionConfig>` whose `mode` changes at runtime | Keep on hooks |

Where the hook rests: `withTiming(TO)`, `withSequence` and a non-reverse `withRepeat` rest at `TO`; `withRepeat(anim, n, true)` with odd `n` runs once under reduced motion and rests at `TO`, with `n <= 0` or even `n` it rests at the start whatever `anim` is (a reversed sequence included, which otherwise rests at its last value). Evaluate the style body at that driver value to know which end the fill mode must hold.

## Colors

`withTiming` and `interpolateColor` interpolate gamma-corrected; CSS lerps sRGB (`../animations/animations.md`, Shared rules). The gap peaks about a quarter in from the darker endpoint.

| Endpoints | Worst channel gap | Verdict |
|---|---|---|
| small swing between bright values (`#eee` to `#ccc`) | 1/255 | Migrate |
| about half the range (grey to mid green) | 36/255 | Migrate, state it |
| a channel crossing most of 0..255 (black to white, red to cyan) | 72/255 | Needs approval |

## SVG

From 4.4.0 (web 4.5.0; on 4.1.0 to 4.3.x only with the `EXPERIMENTAL_CSS_ANIMATIONS_FOR_SVG_COMPONENTS` static flag, a native rebuild; from 4.4.0 the flag is on unless the app's `package.json` sets it to `false`) CSS declarations go in `animatedProps`, never `style` (`../animations/svg-animations.md`). The driver becomes state and the attribute a plain prop: `r={grown ? 50 : 20}` beside `animatedProps={{ transitionProperty: 'r', transitionDuration: 300 }}`. Values from a separate `useAnimatedProps` hook you are not migrating can stay in the same `animatedProps` array. Which attributes animate varies per component; check the component's row in [Animating SVG](https://docs.swmansion.com/react-native-reanimated/docs/guides/animating-svg) before converting. Keep on hooks: SVG below the floors above, SVG `transform` arrays, `mask`, `filter`, `marker*` and `pointerEvents` (listed in that table but they throw `No interpolator factory found`), and `fill`/`stroke` written as `currentColor` or `url(#...)`.
