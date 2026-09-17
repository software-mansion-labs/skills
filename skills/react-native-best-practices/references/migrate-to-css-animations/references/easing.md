# Easing

Question 5 of the walk. Map the `Easing.*` expression written in the source to a CSS timing function; the `Easing` object itself is not accepted by CSS. The `withTiming` default is `Easing.inOut(Easing.quad)` at 300ms.

## Exact rows

| Source | CSS |
|---|---|
| `Easing.linear`, `Easing.poly(1)`, `Easing.inOut(Easing.linear)` | the `'linear'` keyword (the `linear(...)` function with stops is for sampled curves below) |
| `Easing.ease` | `'ease-in'` (not CSS `'ease'`, which is a different curve) |
| `Easing.quad`, `Easing.poly(2)` | `cubicBezier(1/3, 0, 2/3, 1/3)` |
| `Easing.cubic`, `Easing.poly(3)` | `cubicBezier(1/3, 0, 2/3, 0)` |
| `Easing.bezier(a, b, c, d)`, `Easing.bezierFn(...)` | `cubicBezier(a, b, c, d)` |
| `Easing.back(s)` | `cubicBezier(1/3, 0, 2/3, -s/3)`; default `s = 1.70158` gives `-0.56719` |
| `Easing.steps(n)`, `Easing.steps(n, true)` | `steps(n, 'jump-start')` |
| `Easing.steps(n, false)` | `steps(n, 'jump-end')` |
| `Easing.in(f)` | `Easing.in(f)` is `f` itself: map `f` |
| `Easing.out(f)` | reflect: `cubicBezier(x1, y1, x2, y2)` becomes `cubicBezier(1 - x2, 1 - y2, 1 - x1, 1 - y1)`; `out(Easing.ease)` is `'ease-out'` |

Always write the timing function, and do not write `'ease'` for a curve the source did not have: no `Easing.*` value is CSS `'ease'`, and on 4.0.0 to 4.3.x the native `'ease'`, which is also the default when the timing function is left out, is `cubicBezier(0.25, 0.1, 0.25, 0.1)`, a wrong curve (up to 0.40 off the real one). When the user says the exact curve does not matter, offer `'ease-in-out'` (0.012 off the `withTiming` default) or the CSS default look, written as `cubicBezier(0.25, 0.1, 0.25, 1)` on those versions.

## Every other curve

`Easing.inOut(f)`, `poly(n)` for `n` above 3, `sin`, `circle`, `exp`, `elastic`, `bounce`, and the `in`/`out`/`inOut` forms of those have no cubic-bezier form. Two ways to migrate them. Show both for the first such site (the emitted `linear(...)` with its stop count next to the nearest approximation and its max error) and let the user pick; the choice holds for the run, and either way the site's row records the substitution and its max error:

- Sampled: `linear()` (`react-native-reanimated`, available since 4.0.0) with the function's values at evenly spaced stops, `linear(0, f(1/n), f(2/n), ..., 1)` (n intervals, n + 1 points). Compute the max error for the count you emit. 20 intervals hold `inOut(quad|cubic|sin)`, `sin`, `poly(4)`, `poly(5)`, `back` and `elastic` within 0.006; `exp` needs 50 (0.013 at 20, 0.002 at 50); `bounce` is 0.078 at 20 and 0.024 at 50 with even stops but 0.005 at 20 once stops sit on its three kinks (`t = 1/2.75, 2/2.75, 2.5/2.75`; `linear()` accepts `[value, 'pct%']` positions); `circle` has a vertical tangent at the end (0.080 at 20, 0.035 at 100 even stops), so add stops at 99%, 99.5% and 99.9%.
- Approximate: the nearest named curve or cubic-bezier with its max error, which the user may prefer for readability. `'ease-in-out'` stands in for `inOut(quad)` (the `withTiming` default) with a max error of 0.012, for `inOut(sin)` with 0.019, for `inOut(ease)` with 0.029; `inOut(cubic)` (0.084) and `inOut(circle)` (0.136) are visibly different.

Do not sample a spring: its curve depends on the distance and velocity of each run.

## Example

```tsx
// Before
withTiming(1, { duration: 400, easing: Easing.out(Easing.exp) })
```

```tsx
// After, sampled at 50 intervals (max error 0.002); the approximate form would be the nearest cubic-bezier with its error stated
transitionDuration: 400,
transitionTimingFunction: linear(0, 0.13, 0.24, ..., 1),
```
