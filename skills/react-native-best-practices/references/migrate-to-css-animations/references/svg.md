# SVG

Part of question 3, for `react-native-svg` elements driven through `useAnimatedProps`.

CSS reaches SVG attributes from 4.4.0 on iOS and Android (on 4.1.0 to 4.3.x only with the `EXPERIMENTAL_CSS_ANIMATIONS_FOR_SVG_COMPONENTS` static flag, a native rebuild; from 4.4.0 the flag is on unless the app's `package.json` sets it to `false`) and from 4.5.0 on web. Below those floors: Keep on shared values.

The CSS declarations go in `animatedProps`, never `style` (`../animations/svg-animations.md`). The driver becomes state and the attribute a plain prop: `r={grown ? 50 : 20}` beside `animatedProps={{ transitionProperty: 'r', transitionDuration: 300, transitionTimingFunction: 'ease-in-out' }}`. Values from a separate `useAnimatedProps` hook you are not migrating can stay in the same `animatedProps` array.

Which attributes animate varies per component; check the component's row in [Animating SVG](https://docs.swmansion.com/react-native-reanimated/docs/guides/animating-svg) before converting. Keep on shared values: SVG `transform` arrays; `mask`, `filter`, `marker*` and `pointerEvents` (listed in that table but they throw `No interpolator factory found`); `fill`/`stroke` written as `currentColor` or `url(#...)`.

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
