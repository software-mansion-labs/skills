---
name: animations
description: "Production animation patterns for React Native using Reanimated 4, Skia, WebGPU, and TypeGPU. Covers CSS transitions, CSS animations, CSS pseudo-selectors, CSS callbacks, shared value animations, canvas animations with react-native-skia, GPU shader animations, layout animations, scroll-driven animations, particle systems, procedural noise, SDF rendering, performance tuning, and accessibility. Trigger on: Reanimated, useSharedValue, useAnimatedStyle, withSpring, withTiming, withDecay, withRepeat, withSequence, CSS transition, CSS animation, layout animation, FadeIn, SlideIn, ZoomIn, LinearTransition, keyframe, interpolate, scrollTo, useFrameCallback, react-native-skia, Skia Canvas, Atlas, usePathInterpolation, usePathValue, useClock, useTexture, SKSL, interpolateColors, Picture API, canvas animation, sprite animation, WebGPU, react-native-wgpu, TypeGPU, GPU shader, WGSL, particle system, Perlin noise, SDF, Three.js, react-three-fiber, animation performance, or any request to animate UI in React Native."
---

# Animations

Software Mansion's production animation patterns for React Native on Reanimated 4 and the New Architecture.

Load at most one reference file per question; `animations.md` and `css-pseudo-selectors.md` may be read together. For API signatures and config options, webfetch the documentation pages linked in each reference file.

## Version Check

Read the installed Reanimated version before writing animation code: `node_modules/react-native-reanimated/package.json` or the lockfile, not the `package.json` range. Reanimated 3.x has no CSS transitions or animations; use shared values there. On 4.x, a feature used below the version that added it (CSS table in `animations.md`; other APIs carry their floor where they are described) is silently ignored or throws: do not emit it, and if the user asked for that feature, say which version adds it.

## Critical Rules

- **Call JS from a worklet with `scheduleOnRN(fn, ...args)`** (scroll handlers, gesture callbacks, `useAnimatedReaction`, `useFrameCallback`, animation callbacks) and schedule UI work from the JS thread with `scheduleOnUI(fn, ...args)`. Both come from `react-native-worklets` 0.5.0+, the range Reanimated 4.1.0 requires; on Reanimated 4.0.x use `runOnJS`/`runOnUI`, which still exist everywhere as deprecated wrappers. `react-native-reanimated` re-exports `runOnJS` but not `scheduleOnRN`.

## References

| File | When to read |
|------|-------------|
| `animations.md` | Choosing between CSS transitions, CSS animations, and shared value animations; CSS feature availability by Reanimated version; CSS transition and CSS animation patterns and rules; CSS callbacks (`onCSS*`, 4.6.0+); animating text; infinite animation cleanup; `scheduleOnRN` |
| `css-pseudo-selectors.md` | Interaction state without React state: `:hover`, `:active`, `:active-deepest`, `:focus`, `:focus-within` (4.5.0+), selector precedence, the property lock, per-platform press traps |
| `animation-functions.md` | Gotchas and rules for core hooks (`useSharedValue`, `useAnimatedStyle`, `useAnimatedProps`, `useDerivedValue`); `withSpring` config modes and presets (`GentleSpringConfig`); `withRepeat` and `withClamp` caveats; composing animations |
| `layout-animations.md` | Entering/exiting animation gotchas (`nativeID` conflict, view flattening); `withInitialValues`/`withTargetValues`; layout transitions; keyframe animation rules; list item animations (`itemLayoutAnimation`); shared element transitions, `SharedTransitionBoundary` |
| `scroll-and-events.md` | Scroll-driven animation patterns (`useAnimatedScrollHandler`, `scrollTo`, `useScrollOffset`); `useAnimatedReaction` patterns; `useFrameCallback`, `useTimestamp`; `interpolate`, `interpolateColor`; `measure` rules |
| `canvas-animations.md` | Canvas animations with `@shopify/react-native-skia`; Reanimated integration (shared values as direct props); `interpolateColors`; retained vs immediate mode (Picture API); path animations (`usePathInterpolation`, `usePathValue`); `useClock`; gesture integration and element tracking; SKSL runtime shaders and image filters; textures |
| `canvas-atlas.md` | Atlas for batched sprite/tile animation; `useTexture`, `useRSXformBuffer`; RSXform matrix format (`[scos, ssin, tx, ty]`) |
| `gpu-animations.md` | GPU shader animations; `react-native-wgpu` Canvas and device setup; TypeGPU typed pipelines; Reanimated + WebGPU worklet integration; compute pipelines for particle systems, physics, and simulations; `@typegpu/noise` (Perlin noise, PRNG); `@typegpu/sdf` (signed distance shapes); Three.js / React Three Fiber for 3D |
| `svg-animations.md` | Animating SVG elements and paths with Reanimated; `createAnimatedComponent` for SVG; progress arcs; pulsing circles |
| `animations-performance.md` | Performance tuning; 120fps setup; feature flags (platform-driven CSS transitions, `USE_ANIMATION_BACKEND`); FPS drop fixes; simultaneous animation limits; accessibility (`useReducedMotion`, `ReducedMotionConfig`); worklet closure optimization; debug vs release builds |
