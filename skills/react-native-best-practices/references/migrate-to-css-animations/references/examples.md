# Worked examples

Calibration, not an allow-list: the checks in `SKILL.md` decide.

## Migrate: infinite loop

```tsx
// Before
function Spinner() {
  const rotation = useSharedValue(0);
  useEffect(() => {
    rotation.value = withRepeat(withTiming(360, { duration: 2000, easing: Easing.linear }), -1);
    return () => cancelAnimation(rotation);
  }, []);
  const style = useAnimatedStyle(() => ({ transform: [{ rotateZ: `${rotation.value}deg` }] }));
  return <Animated.View style={[styles.box, style]} />;
}
```

```tsx
// After
const rotate: CSSAnimationKeyframes = {
  from: { transform: [{ rotateZ: '0deg' }] },
  to: { transform: [{ rotateZ: '360deg' }] },
};

function Spinner() {
  const reduced = useReducedMotion();
  return (
    <Animated.View
      style={[
        styles.box,
        reduced
          ? { transform: [{ rotateZ: '360deg' }] }
          : { animationName: rotate, animationDuration: 2000, animationIterationCount: 'infinite', animationTimingFunction: 'linear' },
      ]}
    />
  );
}
```

Permitted by check 1 (one write in `useEffect`) and check 5 (`cancelAnimation` only in the cleanup). The reduced branch renders `360deg`: no `reverse`, so the loop rests at the target.

## Migrate: play once on mount

```tsx
// Before
useEffect(() => { opacity.value = withTiming(1, { duration: 300 }); }, []);
```

```tsx
// After: keep opacity: 0 so the first painted frame matches the hook
const fadeIn: CSSAnimationKeyframes = { from: { opacity: 0 }, to: { opacity: 1 } };
// on the element
reduced
  ? { opacity: 1 }
  : { opacity: 0, animationName: fadeIn, animationDuration: 300, animationTimingFunction: 'ease-in-out', animationFillMode: 'forwards' }
```

The reduced branch renders the resting value. The report notes `inOut(quad)` to `'ease-in-out'`, max error 0.012.

## Migrate: one-way state change

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

Permitted by check 1 (JS effect) and check 7 (an upload completes once, no reversal). Grey to green is the moderate color case: Migrate and state the 36/255 peak gap in the row (the transition is a color, so no reduced-motion guard).

## Needs approval: the toggle that can reverse

Expand/collapse, show/hide, a switch: the user can flip back inside the duration, and check 7 stops there.

```tsx
// Before: leave exactly as is until the user accepts
const [expanded, setExpanded] = useState(false);
const progress = useSharedValue(0);
useEffect(() => { progress.value = withTiming(expanded ? 1 : 0, { duration: 200 }); }, [expanded]);
const style = useAnimatedStyle(() => ({ height: 100 + progress.value * 200 }));
```

Reversed 80ms in, the hook takes the full 200ms back and CSS about 66ms. Propose the transition (`height: expanded ? 300 : 100`, `transitionProperty: 'height'`, `transitionDuration: reduced ? 1 : 200`, `transitionTimingFunction: 'ease-in-out'`) as a Needs approval row.

## Trap: the resting value under reduced motion

```tsx
const progress = useSharedValue(0);
useEffect(() => { progress.value = withRepeat(withTiming(1, { duration: 600 }), -1, true); }, []);
const style = useAnimatedStyle(() => ({ transform: [{ scaleY: interpolate(progress.value, [0, 1], [0.35, 1]) }] }));
```

Reduced motion never advances `progress`, so the style rests at `interpolate(0, ...)` = `0.35`. The reduced branch is `{ transform: [{ scaleY: 0.35 }] }`; copying the `to` keyframe parks the element nearly three times too tall. Evaluate the style body at the driver's resting value, never read it off the keyframes you wrote.
