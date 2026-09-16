# Properties

Question 3 of the walk.

CSS animates a property only when Reanimated has an interpolator for it. Check the [supported properties](https://docs.swmansion.com/react-native-reanimated/docs/guides/supported-properties) table, never memory; it describes the latest release, so on an older installed version also check the feature floors in `../animations/animations.md` (filter 4.2.0, SVG 4.4.0, pseudo-selectors 4.5.0, callbacks 4.6.0).

- If React Native never renders the property on a platform (`shadowOffset`, `shadowOpacity` and `shadowRadius` on Android, `elevation` on iOS), CSS not animating it there changes nothing.
- If React Native renders it and CSS cannot animate it: Keep on shared values.
- Discrete (keyword) properties change at the transition midpoint only with `transitionBehavior: 'allow-discrete'` (`../animations/animations.md`, Discrete properties). Route by where the shared value version flipped the keyword:

| The keyword flipped | Verdict |
|---|---|
| at a state change (a boolean driver, a ternary on state) | render the keyword conditionally and leave it out of `transitionProperty`; Migrate |
| at `> 0.5` of a 0..1 numeric driver | propose `transitionBehavior: 'allow-discrete'`, which flips at the midpoint; Needs approval |
| at any other threshold, or where the flip breaks layout | Keep on shared values |

## Example

```tsx
// Before
const style = useAnimatedStyle(() => ({ display: progress.value > 0.5 ? 'flex' : 'none', opacity: progress.value }));
```

```tsx
// After (Needs approval: display flips at the midpoint of the opacity transition)
{ display: open ? 'flex' : 'none', opacity: open ? 1 : 0, transitionProperty: ['display', 'opacity'], transitionBehavior: 'allow-discrete', transitionDuration: reduced ? 1 : 300, transitionTimingFunction: 'ease-in-out' }
```
