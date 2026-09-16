# Worked examples

Calibration, not an allow-list: the decision tree in `SKILL.md` decides.

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
        {
          animationName: rotate,
          animationDuration: reduced ? 1 : 2000,
          animationIterationCount: reduced ? 1 : 'infinite',
          animationTimingFunction: 'linear',
        },
      ]}
    />
  );
}
```

A JS-thread write in `useEffect`, and `cancelAnimation` only in the cleanup, so the tree reaches Migrate. Under reduced motion the loop runs once for 1ms and snaps back to the static style; a non-`reverse` loop rests at its target, which here renders the same as the start (`360deg` is `0deg`), so `animationFillMode` can stay `'none'`.

## Migrate: play once on mount

```tsx
// Before
useEffect(() => { opacity.value = withTiming(1, { duration: 300 }); }, []);
```

```tsx
// After: keep opacity: 0 so the first painted frame matches the hook
const fadeIn: CSSAnimationKeyframes = { from: { opacity: 0 }, to: { opacity: 1 } };
// on the element
{ opacity: 0, animationName: fadeIn, animationDuration: reduced ? 1 : 300, animationTimingFunction: 'ease-in-out', animationFillMode: 'forwards' }
```

Under reduced motion the 1ms run lands on `opacity: 1` through the fill mode, as `withTiming` jumps to its target. The easing is the `withTiming` default: the report notes `inOut(quad)` to `'ease-in-out'`, max error 0.012, or the exact `linear()` form if the user chose it.

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

A JS effect drives it and an upload completes once, so nothing reverses mid-flight. Grey to green is the moderate color case: Migrate and state the 36/255 peak gap in the row; the transition is a color under 300ms, so the row also says `reduced-motion guard dropped`.

## Needs approval: the toggle that can reverse

Expand/collapse, show/hide, a switch: the user can flip back inside the duration, so the reversal question notes Needs approval.

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

Reduced motion never advances `progress`, so the hook rests at `interpolate(0, ...)` = `0.35`, the start. The shortened animation must rest there too: static style `scaleY: 0.35` and `animationFillMode: 'none'`, so the 1ms run snaps back to it. `'forwards'` would park the element at the `to` keyframe, nearly three times too tall. Evaluate the style body at the driver's resting value, never read it off the keyframes you wrote.
