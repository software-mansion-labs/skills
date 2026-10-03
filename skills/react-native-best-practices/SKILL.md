---
name: react-native-best-practices
description: "Software Mansion's best practices for production React Native and Expo apps on the New Architecture. Use before writing, reviewing, or debugging code in a project whose package.json depends on react-native, expo, or expo-router, and for any question about React Native UI, graphics, animation, gestures, audio, rich text, on-device AI, native modules (JSI, TurboModules, Fabric), or threading (worklets, Reanimated shared values, runOnJS/scheduleOnRN). Routes to a sub-skill per topic; see the table in the body."
license: MIT
---

# React Native Best Practices

Software Mansion's production patterns for React Native apps on the New Architecture.

Read the relevant sub-skill for the topic at hand. All sub-skills are in `references/`.

## Sub-skills

| Sub-skill | When to use |
|-----------|------------|
| `references/animations/SKILL.md` | CSS transitions, CSS animations, CSS pseudo-selectors and callbacks, shared value animations, GPU shader animations (WebGPU, TypeGPU), layout animations (entering/exiting, transitions, keyframes), scroll-driven animations, animation functions (withSpring, withTiming, withDecay), core hooks (useSharedValue, useAnimatedStyle), interpolation, particle systems, procedural noise, SDF rendering, animation performance, 120fps, accessibility, Reanimated 4 |
| `references/gestures/SKILL.md` | Tap, pan, pinch, rotation, swipe, long press, fling, hover, drag, Pressable, RectButton, Swipeable, DrawerLayout, VirtualGestureDetector, gesture composition, gesture testing -- any touch interaction with Gesture Handler |
| `references/svg/SKILL.md` | Vector graphics, icons, charts, illustrations using React Native SVG |
| `references/on-device-ai/SKILL.md` | On-device AI: LLMs (chat, tool calling, structured output, vision-language models), computer vision (classification, object detection, OCR, semantic/instance segmentation, style transfer, embeddings, text-to-image), speech processing (STT with timestamps, TTS with phonemes, VAD), VisionCamera real-time frame processing, model loading, resource management, custom models with ExecuTorch |
| `references/rich-text/SKILL.md` | Rich text editor, formatted text input, WYSIWYG, mentions, HTML/Markdown rendering, react-native-enriched-html (formerly react-native-enriched), react-native-enriched-markdown |
| `references/multithreading/SKILL.md` | Multithreading, react-native-worklets, background processing, Worker Runtimes, UI thread, scheduleOnUI, scheduleOnRN, Serializable, Synchronizable, offloading computation from the JS thread |
| `references/enable-worklets-bundle-mode/SKILL.md` | Enabling react-native-worklets Bundle Mode (imports inside worklets, third-party npm libraries on worklet runtimes) in an Expo, RN CLI, or brownfield app: babel bundleMode plugin option, bundleModeMetroConfig / getBundleModeMetroConfig, mandatory metro and metro-runtime patches per package manager, "Failed to get the SHA-1" errors, missing Fast Refresh for worklet code, uniwind/NativeWind resolver conflicts |
| `references/audio/SKILL.md` | Audio playback (buffer sources, oscillators, streaming, queued playback), recording (file, data callback, graph processing), audio effects (gain, filters, delay, convolver, panner, waveshaper), real-time analysis and visualization, audio worklets (custom processing, synthesis), system integration (sessions, interruptions, notifications, permissions), testing with mocks -- any audio feature with react-native-audio-api |
| `references/jsi/SKILL.md` | JSI, C++ native modules, jsi::Runtime, jsi::Value, jsi::Object, jsi::Function, jsi::HostObject, jsi::HostFunction, jsi::NativeState, jsi::PropNameID, jsi::ArrayBuffer, jsi::WeakObject, jsi::Scope, jsi::BigInt, JSIException, JSError, JSINativeException, calling JS from C++, calling C++ from JS, HostObject destructor constraints, shared_ptr<jsi::Value>, CallInvoker, invokeAsync, JSI threading safety, zero-copy ArrayBuffer, rt.global(), ISerialization, WithRuntimeDecorator, jsi.h |
