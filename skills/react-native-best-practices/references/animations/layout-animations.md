# Layout Animations

Animations for components entering, exiting, or changing position in the view hierarchy. Reanimated 4, New Architecture required.

For predefined animation lists, modifiers, and parameters, webfetch the linked documentation pages below.

---

## [Entering and Exiting Animations](https://docs.swmansion.com/react-native-reanimated/docs/layout-animations/entering-exiting-animations)

Animate elements when they are added to or removed from the view hierarchy:

```tsx
import Animated, { FadeIn, FadeOut } from 'react-native-reanimated';

{visible && (
  <Animated.View entering={FadeIn} exiting={FadeOut}>
    <Text>Hello</Text>
  </Animated.View>
)}
```

Predefined animation families include Fade, Slide, Zoom, Bounce, Flip, Stretch, Roll, Rotate, LightSpeed, and Pinwheel. Each has directional variants (e.g., `FadeInRight`, `FadeInLeft`, `FadeInUp`, `FadeInDown`).

Chain modifiers on any predefined animation:

```tsx
entering={FadeIn.delay(200).springify().damping(15)}
```

`.easing()` has no effect once `.springify()` is used. A spring is duration-based (`.springify().duration(550).dampingRatio(0.75)`) or physics-based (`.springify().damping(30).stiffness(900)`); `.mass()` applies to both. When both kinds of modifier are present, `duration` and `dampingRatio` win and `damping` and `stiffness` are ignored.

### Gotchas

- **`nativeID` conflict (New Architecture)**: Reanimated uses `nativeID` internally for entering animations. Overwriting it breaks the animation. Wrap animated children in a plain `View` to work around this, especially with `TouchableWithoutFeedback`.
- **View flattening**: Removing a non-animated parent triggers exiting animations in its children, but the parent will not wait for children to finish. Add `collapsable={false}` to the parent to prevent this.
- **Spring-based animations**: Not yet available on the web platform.
- **Performance**: Define animation builders outside of components or wrap with `useMemo`.
- `.energyThreshold()` (default `6e-9`) decides when a spring rests (4.1.0+); `.restDisplacementThreshold()` and `.restSpeedThreshold()` are no-ops since 4.1.0.

Override a preset's start and end state (entering/exiting only, not layout transitions):

```tsx
entering={FadeInDown.withInitialValues({ translateY: 420 }).withTargetValues({ translateY: 0 })}
```

From 4.4.0: flat transform props (`{ translateX: 50 }`) and `withTargetValues`. Below 4.4.0 only `withInitialValues({ transform: [{ translateX: 50 }] })` exists.

---

## [Layout Transitions](https://docs.swmansion.com/react-native-reanimated/docs/layout-animations/layout-transitions)

Smooth animations when a component's position or size changes due to state updates:

```tsx
import Animated, { LinearTransition } from 'react-native-reanimated';

<Animated.View layout={LinearTransition}>
  {items.map((item) => (
    <Item key={item.id} {...item} />
  ))}
</Animated.View>
```

Predefined transitions: `LinearTransition`, `SequencedTransition`, `FadingTransition`, `JumpingTransition`, `CurvedTransition`, `EntryExitTransition`.

The generic `Layout` transition from older Reanimated versions is deprecated. Use `LinearTransition`.

**Spring config modes**: Use either physics-based (`damping`/`stiffness`) or duration-based (`duration`/`dampingRatio`), never both. If both are provided, duration-based overrides.

---

## [Keyframe Animations](https://docs.swmansion.com/react-native-reanimated/docs/layout-animations/keyframe-animations)

For complex multi-step entering/exiting animations beyond what presets offer:

```tsx
import Animated, { Easing, Keyframe } from 'react-native-reanimated';

const enteringAnimation = new Keyframe({
  0: { opacity: 0, transform: [{ scale: 0.5 }, { rotate: '-45deg' }] },
  50: {
    opacity: 1,
    transform: [{ scale: 1.2 }, { rotate: '0deg' }],
    easing: Easing.out(Easing.quad),
  },
  100: { transform: [{ scale: 1 }, { rotate: '0deg' }] },
});

<Animated.View entering={enteringAnimation.duration(600)} />
```

### Rules

- Keyframe `0` (or `from`) is **required**. Provide initial values for all properties you intend to animate.
- Keyframe `100` (or `to`) is optional.
- Do not provide both `0` and `from`, or both `100` and `to`.
- Easing is assigned to the second keyframe in a pair. Never provide easing to keyframe `0`.
- Default easing between keyframes is `Easing.linear`.
- **All properties in the transform array must appear in the same order across all keyframes.**

---

## [List Layout Animations](https://docs.swmansion.com/react-native-reanimated/docs/layout-animations/list-layout-animations)

Animate item layout changes in `FlatList` when items are added, removed, or reordered:

```tsx
<Animated.FlatList
  data={data}
  renderItem={renderItem}
  itemLayoutAnimation={LinearTransition}
/>
```

### Rules

- Only works with single-column `FlatList`. `numColumns` cannot be greater than 1.
- Items must have a `key` or `id` property (or provide a custom `keyExtractor`).
- Set `itemLayoutAnimation` to `undefined` to disable at runtime.
- `skipEnteringExitingAnimations` is a prop on `Animated.FlatList`, not a modifier. Any defined value, including `false`, skips; pass it only when skipping.

---

## [LayoutAnimationConfig](https://docs.swmansion.com/react-native-reanimated/docs/layout-animations/layout-animation-config)

Skip entering/exiting animations for a subtree:

```tsx
import { LayoutAnimationConfig } from 'react-native-reanimated';

<LayoutAnimationConfig skipEntering skipExiting>
  {children}
</LayoutAnimationConfig>
```

Can be nested. For FlatLists, pass the `skipEnteringExitingAnimations` prop on `Animated.FlatList` instead, which applies this wrapper.

---

## [Shared Element Transitions](https://docs.swmansion.com/react-native-reanimated/docs/shared-element-transitions/overview)

**Status: Experimental, off by default. Not recommended for production.**

Requires 4.2.0+ and the `ENABLE_SHARED_ELEMENT_TRANSITIONS` static feature flag (`package.json`: `"reanimated": { "staticFeatureFlags": { "ENABLE_SHARED_ELEMENT_TRANSITIONS": true } }`, then `pod install` and rebuild; not possible in Expo Go). With the flag off, `sharedTransitionTag` is silently ignored. The flag disables the `*_SYNCHRONOUSLY_UPDATE_UI_PROPS` fast path, and the iOS `pod install` or the Android Gradle build fails if both are set.

Animates a view between two screens during navigation:

```tsx
<Animated.Image
  sharedTransitionTag="hero-image"
  sharedTransitionStyle={SharedTransition.duration(550).springify()}
/>
```

- With a navigator, only the React Navigation native stack is supported. Tab navigator and `transparentModal` (iOS) are not.
- 4.5.0+: `SharedTransitionBoundary` drops the navigator requirement. Wrap each side and toggle `isActive`: `<SharedTransitionBoundary isActive={activeId === 0}><Animated.View sharedTransitionTag="tag" /></SharedTransitionBoundary>`.
- Tags must be unique per screen. Add the same tag to matching components on both screens.
- Default duration: 500ms. Animates width, height, position, transform, backgroundColor, opacity.
- iOS supports progress-based (swipe gesture) transitions. Android uses timing-based transitions only.
- Custom animation functions are not yet supported.
