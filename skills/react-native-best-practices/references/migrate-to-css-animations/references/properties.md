# Properties

Question 3 of the walk.

CSS animates a property only when Reanimated has an interpolator for it. Check the [supported properties](https://docs.swmansion.com/react-native-reanimated/docs/guides/supported-properties) table, never memory; it describes the latest release, so on an older installed version also check the feature floors in `../animations/animations.md` (filter 4.2.0, SVG 4.4.0, pseudo-selectors 4.5.0, callbacks 4.6.0).

- If React Native never renders the property on a platform (`shadowOffset`, `shadowOpacity` and `shadowRadius` on Android, `elevation` on iOS), CSS not animating it there changes nothing.
- If React Native renders it and CSS cannot animate it: Keep on shared values.
- Discrete (keyword) properties change only with `transitionBehavior: 'allow-discrete'`, at the transition midpoint (`../animations/animations.md`, Discrete properties). `display` to or from `none` is the exception: it flips at the start when showing and at the end when hiding, so the element is visible for the whole transition either way. Route by where the shared value version flipped the keyword:

| The keyword flipped | Verdict |
|---|---|
| at a state change (a boolean driver, a ternary on state) | render the keyword conditionally and leave it out of `transitionProperty`; Migrate |
| at `> 0.5` of a 0..1 numeric driver | propose `transitionBehavior: 'allow-discrete'`, which flips at the midpoint (`display` around `none`: at the start when showing, at the end when hiding); Needs approval |
| at any other threshold, or where the flip breaks layout | Keep on shared values |

A Keep row ends the walk for that property; the other two continue it, the Needs approval row noted.

## Example

```tsx
// Before
const style = useAnimatedStyle(() => ({ display: progress.value > 0.5 ? 'flex' : 'none', opacity: progress.value }));
```

```tsx
// After (Needs approval: display leaves none at the start of the show transition and
// returns to none at the end of the hide transition, where the shared value flipped it at the midpoint)
{ display: open ? 'flex' : 'none', opacity: open ? 1 : 0, transitionProperty: ['display', 'opacity'], transitionBehavior: 'allow-discrete', transitionDuration: reduced ? 1 : 300, transitionTimingFunction: 'ease-in-out' }
```
