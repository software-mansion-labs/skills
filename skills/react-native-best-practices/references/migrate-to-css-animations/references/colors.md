# Colors

Part of question 4. `withTiming` and `interpolateColor` interpolate gamma-corrected; CSS lerps sRGB (`../animations/animations.md`, Shared rules). The endpoints match, the frames between them do not, and the gap peaks about a quarter in from the darker endpoint.

| Endpoints | Worst channel gap | Verdict |
|---|---|---|
| small swing between bright values (`#eee` to `#ccc`) | 1/255 | Migrate |
| about half the range (grey to mid green) | 36/255 | Migrate, state it in the row |
| a channel crossing most of 0..255 (black to white, red to cyan) | 72/255 | Needs approval |

`interpolateColor` in `'HSV'` or `'LAB'` stays on shared values (`references/value-functions.md`).

## Example: one-way state change

```tsx
// Before
const sv = useSharedValue(false);
const style = useAnimatedStyle(() => ({ backgroundColor: withTiming(sv.value ? 'green' : 'grey', { duration: 200 }) }));
useEffect(() => { if (uploadComplete) sv.value = true; }, [uploadComplete]);
```

```tsx
// After
<Animated.View
  style={[styles.box, {
    backgroundColor: uploadComplete ? 'green' : 'grey',
    transitionProperty: 'backgroundColor',
    transitionTimingFunction: 'ease-in-out',
    transitionDuration: 200,
  }]}
/>
```

A JS effect drives it and an upload completes once, so nothing reverses mid-flight. Grey to green is the moderate case: Migrate and state the 36/255 peak gap in the row, together with `inOut(quad)` to `'ease-in-out'` (0.012).
