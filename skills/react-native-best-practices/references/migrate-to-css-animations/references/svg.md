# SVG

Part of question 3, for `react-native-svg` elements driven through `useAnimatedProps`.

CSS reaches SVG attributes from 4.4.0 on iOS and Android and from 4.5.0 on web. On 4.1.0 to 4.3.x a site is eligible only when the app's `package.json` sets `reanimated.staticFeatureFlags.EXPERIMENTAL_CSS_ANIMATIONS_FOR_SVG_COMPONENTS` to `true` (a native rebuild): 4.1.0 and 4.2.x cover Circle, Ellipse, Line, Path and Rect, 4.3.x adds Image, LinearGradient, Pattern, Polygon, Polyline, RadialGradient and Text (the 4.4.0 set). Without the flag there: Keep on shared values, saying that enabling the flag would allow the site. From 4.4.0 the flag is on unless that field sets it to `false`, which keeps every SVG site on shared values. Below the web floor, a project that targets web keeps its SVG sites.

The CSS declarations go in `animatedProps`, never `style` (`../animations/animations.md`, feature table; the component rows are in the [Animating SVG](https://docs.swmansion.com/react-native-reanimated/docs/guides/animating-svg) docs). Pseudo objects for an SVG element (`references/press-feedback.md`) go in `animatedProps` too: only that path injects the hit-test responder on native, `style` works on web only. The driver becomes state and the attribute a plain prop: `r={grown ? 50 : 20}` beside `animatedProps={{ transitionProperty: 'r', transitionDuration: 300, transitionTimingFunction: 'ease-in-out' }}`. Values from a separate `useAnimatedProps` hook you are not migrating can stay in the same `animatedProps` array.

Which attributes animate varies per component; check the component's row in Animating SVG before converting. Keep on shared values: SVG `transform` arrays; `mask`, `filter`, `marker*` and `pointerEvents` (listed in that table but they throw `No interpolator factory found`); `fill`/`stroke` written as `currentColor` or `url(#...)`.

## Example

```tsx
// Before
const animatedProps = useAnimatedProps(() => ({ r: withTiming(grown ? 50 : 20, { duration: 300 }) }), [grown]);
<AnimatedCircle cx={70} cy={70} fill="#f59e0b" animatedProps={animatedProps} />
```

```tsx
// After (4.4.0+)
<AnimatedCircle
  cx={70}
  cy={70}
  fill="#f59e0b"
  r={grown ? 50 : 20}
  animatedProps={{ transitionProperty: 'r', transitionDuration: 300, transitionTimingFunction: 'ease-in-out' }}
/>
```
