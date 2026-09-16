# Animations Performance and Accessibility

Reanimated 4 requires the New Architecture (Fabric). All guidance here assumes that.

---

## 120fps Support

Enable ProMotion display support on iOS by adding to `Info.plist`:

```xml
<key>CADisableMinimumFrameDurationOnPhone</key>
<true/>
```

Without this flag, iOS caps animations at 60fps even on ProMotion devices. The React Native app template sets it from React Native 0.82; check `Info.plist` before adding it.

---

## Feature Flags

Reanimated 4 exposes feature flags to opt into fixes for known New Architecture issues. Every flag below is *static*: it resolves at compile time and cannot be changed at runtime. Set it in your app's `package.json`, run `pod install` (iOS), and rebuild the native app.

```json
{
  "reanimated": {
    "staticFeatureFlags": {
      "DISABLE_COMMIT_PAUSING_MECHANISM": true
    }
  }
}
```

Static flags cannot be changed where Reanimated ships prebuilt (Expo Go, RNRepo). Read one back with `getStaticFeatureFlag(name)` from `react-native-reanimated`.

### Platform-Driven CSS Transitions

Run CSS transitions on the platform's own animation API instead of Reanimated's animation loop, so Reanimated stops recomputing and committing the value every frame.

- `IOS_CSS_CORE_ANIMATION` (4.4.0): a Core Animation animation on the view's layer
- `ANDROID_CSS_PLATFORM_TRANSITIONS` (4.6.0): an `ObjectAnimator` writing straight to the platform view

Both experimental, both default `false`. CSS *animations* always stay on the loop.

Routing is per property. On 4.4.x iOS routes `opacity` only; from 4.5.0 it routes `opacity`, `backgroundColor`, `borderColor`, `borderRadius`, `borderWidth`, `shadowColor`, `shadowOffset`, `shadowOpacity`, `shadowRadius`. Android routes `opacity` only. A property falls back to the loop if the component uses any `onCSSTransition*` prop (4.6.0), or, on iOS, if its timing function is `steps` or `linear()` with stops.

iOS caveat: `backgroundColor`, `borderColor`, `borderWidth` and `borderRadius` are routed even when React Native draws them on separate layers, where the routed animation never arrives and the value jumps. That happens on any view with a visible border and the default `overflow`, or with per-side border or per-corner radius differences.

### Animation Backend

`USE_ANIMATION_BACKEND` (4.4.0, default `false`) hands applying animated changes to React Native's Animation Backend. Requires React Native 0.85.2+ with its `useSharedAnimatedBackend` flag on, overridden per app like `preventShadowTreeCommitExhaustion` below.

It cannot be enabled alongside `FORCE_REACT_RENDER_FOR_SETTLED_ANIMATIONS`, which is on by default since 4.3.0; turn that off in the same block:

```json
{
  "reanimated": {
    "staticFeatureFlags": {
      "USE_ANIMATION_BACKEND": true,
      "FORCE_REACT_RENDER_FOR_SETTLED_ANIMATIONS": false
    }
  }
}
```

### Flickering / Jittering While Scrolling

Animated components like sticky headers flicker during `FlatList` or `ScrollView` scrolling on the New Architecture.

**Fix:** Upgrade to React Native 0.81+ and enable:
- `preventShadowTreeCommitExhaustion` (React Native). Override it per app through React Native's feature-flag override API before the host starts (Android: `ReactNativeFeatureFlags.override(object : ReactNativeFeatureFlagsDefaults() { override fun preventShadowTreeCommitExhaustion() = true })`; iOS: `ReactNativeFeatureFlags::override(...)` in `AppDelegate.mm`) rather than switching RN to the experimental release level, which enables unrelated flags too.
- `DISABLE_COMMIT_PAUSING_MECHANISM` (Reanimated feature flag)

### FPS Drops During Scrolling

FPS drops when many animated components are visible during scroll.

**Fix:** Upgrade to React Native 0.80+. `USE_COMMIT_HOOK_ONLY_FOR_REACT_COMMITS` (added 4.2.0) is on by default since 4.3.0; set it explicitly on 4.2.x.

Alternative: `enableCppPropsIteratorSetter`, a React Native flag. Experimental, and it requires patching React Native's source files and building React Native from source.

### Low FPS with Many Simultaneous Animations

**Fix:** Enable platform-specific synchronous UI update flags:
- `ANDROID_SYNCHRONOUSLY_UPDATE_UI_PROPS` (available since 4.0.0)
- `IOS_SYNCHRONOUSLY_UPDATE_UI_PROPS` (available since 4.2.0)

Note: these flags may interfere with touch detection on animated `transform` elements. Prefer `Pressable` from `react-native-gesture-handler` over the core `Pressable` when using these flags.

---

## Simultaneous Animation Limits

Reanimated can handle many animated components, but performance degrades at scale:

| Platform       | Practical limit |
|----------------|-----------------|
| iOS            | ~500 components |
| Low-end Android | ~100 components |

For lists with many animated items, consider reducing animation complexity on low-end devices (by device class, not `useReducedMotion`, which reflects the accessibility setting). For highly complex animation scenes (hundreds of elements), consider Reanimated + `react-native-skia` instead of animating native views.

---

## Prefer Non-Layout Properties

Animating layout properties (`top`, `left`, `width`, `height`, `margin`, `padding`) forces a layout pass on every frame.

Prefer non-layout properties: `transform` (all transforms), `opacity`, `backgroundColor`.

The synchronous fast path is a separate mechanism that exists only with `ANDROID_SYNCHRONOUSLY_UPDATE_UI_PROPS` / `IOS_SYNCHRONOUSLY_UPDATE_UI_PROPS` enabled (incompatible with `ENABLE_SHARED_ELEMENT_TRANSITIONS`, `layout-animations.md`): it carries `opacity`, `transform`, `zIndex`, `elevation`, `borderRadius`, `outline*` and the color props (`backgroundColor`, `borderColor`, `shadowColor`, `tintColor`, `placeholderTextColor`; not `PlatformColor` values on Android), plus `shadowOffset`/`shadowOpacity`/`shadowRadius` on iOS. Everything else goes through a shadow tree commit.

If a design requires a size change, consider `scale` transforms for the same visual effect without triggering layout.

---

## Avoid Reading Shared Values on the JS Thread

Reading `sv.value` inside React render, event handlers, or `useEffect` triggers a synchronization from the UI thread to the JS thread, which can block the JS thread.

Instead, use `useDerivedValue` to transform shared values and `useAnimatedStyle` to consume them — both run on the UI thread.

---

## Memoize Callbacks and Gesture Objects

Frame callbacks and gesture objects are re-created on every render by default. Wrap them:

```tsx
const frameCallback = useFrameCallback(
  useCallback((frameInfo) => {
    // runs on UI thread every frame
  }, [])
);

const gesture = useMemo(() =>
  Gesture.Pan().onUpdate((e) => {
    offset.value = e.translationX;
  }),
  []
);
```

If React Compiler is available, it handles memoization automatically.

---

## Worklet Closure Optimization

Worklets capture variables from their surrounding scope. Capturing large objects causes performance issues due to serialization overhead.

```tsx
// Bad: captures entire theme object
const theme = useTheme();
const style = useAnimatedStyle(() => ({
  backgroundColor: theme.colors.primary,
}));

// Good: extract only what is needed
const primaryColor = useTheme().colors.primary;
const style = useAnimatedStyle(() => ({
  backgroundColor: primaryColor,
}));
```

Extract specific properties before the worklet to minimize the closure payload.

Functions marked with `'worklet'` are not hoisted. They must be defined before they are referenced in other worklets.

---

## Debug vs. Release Builds

Always profile animations in a release build. Debug builds add significant JS overhead (Metro bundler, Hermes debug mode, dev warnings) that makes animations appear slower than they are in production.

```
npx react-native run-android --mode=release
```

On Android, use the `debugOptimized` build variant (React Native 0.82, backported to 0.81.2) for a better dev experience with closer-to-production performance.

---

## Accessibility

### useReducedMotion

```tsx
const reduceMotion = useReducedMotion();
```

Returns `true` if the device has reduced motion enabled **at app start**. Does not update at runtime if the user changes the setting.

Use it to conditionally render animations or pick simpler alternatives:

```tsx
<Animated.View entering={reduceMotion ? undefined : FadeIn} />
```

### ReducedMotionConfig

Sets global animation behavior for the entire app:

```tsx
import { ReducedMotionConfig, ReduceMotion } from 'react-native-reanimated';

// Place near app root
<ReducedMotionConfig mode={ReduceMotion.System} />
```

Modes:
- `ReduceMotion.System` (default) — follow device setting
- `ReduceMotion.Always` — always disable animations
- `ReduceMotion.Never` — always enable animations

### Behavior per animation type when reduced motion is enabled

| Animation | Behavior |
|-----------|----------|
| `withSpring`, `withTiming` | Jump to `toValue` immediately |
| `withDecay` | Return current value (respecting clamp) |
| `withDelay` | Start next animation immediately |
| `withRepeat` (infinite or even + reversed) | Do not start |
| `withRepeat` (other) | Run once |
| `withSequence` | Only start children with `reduceMotion: Never` |
| Entering / keyframe / layout animations | Jump to endpoint immediately |
| Exiting / shared element transitions | Omitted entirely |
| CSS transitions and CSS animations | Not affected: they ignore the setting. Shorten them yourself from `useReducedMotion()` (`animations.md`, Reduced motion) |

Higher-order animations pass their `reduceMotion` config to children unless a child has its own explicit config.
