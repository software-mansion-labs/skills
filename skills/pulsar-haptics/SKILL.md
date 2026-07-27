---
name: pulsar-haptics
description: >
  Implements, migrates, designs, and troubleshoots haptic feedback with Software
  Mansion's Pulsar SDK across its native, cross-platform, and Web packages. Use when
  the user names Pulsar, react-native-pulsar, com.swmansion:pulsar,
  com.swmansion:pulsar-kmp, pulsar_haptics, pulsar-haptics, or Pulsar's Presets,
  PatternComposer, RealtimeComposer, Settings, migration, installation, or runtime
  problem. Do not use for generic haptics, browser Vibration API, Core Haptics,
  expo-haptics, or UI-polish requests unless the user explicitly wants Pulsar.
---

# Pulsar Haptics

Implement against the project's installed Pulsar package, not remembered API details.
Keep one workflow across platforms; branch only where the installed SDK differs.

## Sources

Treat the [SDK overview](https://docs.swmansion.com/pulsar/sdk/overview/) as the
current platform and version index. Then open the matching official SDK page:

- [React Native and Expo](https://docs.swmansion.com/pulsar/sdk/react-native/)
- [iOS and Swift](https://docs.swmansion.com/pulsar/sdk/ios/)
- [Android and Kotlin](https://docs.swmansion.com/pulsar/sdk/android/)
- [Kotlin Multiplatform](https://docs.swmansion.com/pulsar/sdk/kmp/)
- [Flutter](https://docs.swmansion.com/pulsar/sdk/flutter/)
- [Web](https://docs.swmansion.com/pulsar/sdk/web/)

Use the [official Pulsar source](https://github.com/software-mansion/pulsar) and
the installed package as version-specific fallbacks. Use the
[preset playground](https://docs.swmansion.com/pulsar/presets-playground/) to
compare current native-family presets on real hardware. Do not copy a full preset
catalog, version table, installation command, API signature, support enum, or
device matrix into the answer without checking it for the installed version.

## Workflow

### 1. Inspect before asking

Read the target handler/component and adjacent state, validation, success, failure,
gesture, and cleanup paths. Inspect manifests and package metadata. Infer:

- platform and framework;
- event semantics and whether feedback marks intent, success, failure, or progress;
- user-triggered versus system-triggered behavior;
- repetition rate, reversibility, urgency, and existing visual/audio feedback;
- whether the request needs a preset, a fixed custom timeline, or live modulation.

Ask only unresolved questions that change implementation. Ask at most two together.
If platform or success semantics remain unknown, clarify before emitting code.

### 2. Resolve the installed package and version

Use both declared and resolved evidence:

- React Native or Web: `package.json`, lockfile, installed package metadata, and
  exported TypeScript declarations;
- iOS: `Package.resolved`, `Package.swift`, Xcode project settings, or resolved
  CocoaPods metadata;
- Android or Kotlin Multiplatform: Gradle build files, version catalogs, dependency
  locks, and resolved dependency reports;
- Flutter: `pubspec.yaml`, `pubspec.lock`, and installed package source.

Do not upgrade, add, or replace a dependency unless requested. If Pulsar is absent,
stop implementation and ask whether to add it. If declared and resolved versions
differ, use the resolved version and report the mismatch.

### 3. Match documentation to installed evidence

Open the overview and matching platform page. Compare their current version and
requirements with the project's resolved package.

If they match, verify symbols and setup against installed types or source when cheap.
If they differ, prefer this order:

1. installed declarations, generated interfaces, or package source;
2. vendor source or release tag for the resolved version;
3. resolved registry metadata and lockfiles;
4. current docs only for concepts that are unchanged.

When official docs are unavailable, use the same fallback order. State which source
was unavailable, what local evidence was used, and what remains unverified. Never
silently combine signatures, requirements, presets, or support levels from different
versions.

### 4. Choose preset first

Prefer a named or system preset when it matches the event.

1. Read code and infer event meaning before asking questions.
2. Search the installed preset exports and matching official platform page.
3. Choose one primary preset that fits semantics, intensity, duration, repetition,
   and reversibility.
4. Give paste-ready code using the exact installed-version symbol.
5. Mention at most two focused alternatives, each tied to a concrete condition.

Use system presets when migrating system-style feedback and semantics must remain
native. Do not print a catalog or expose internal tag-filtering steps.

### 5. Use composers only when needed

Use `PatternComposer` when the complete custom timeline is known before playback.
Use `RealtimeComposer` when amplitude or frequency must follow changing gesture or
sensor values.

Before coding, verify the installed platform's data model, method names, sync/async
behavior, ownership, and cleanup rules. Native-family and Web pattern shapes are not
interchangeable.

For live feedback:

- map inputs into the installed API's accepted range;
- reuse the composer for the interaction lifecycle;
- layer discrete feedback only at meaningful snap points or boundaries;
- stop on normal end, cancellation, failure, unmount/disposal, and any relevant
  lifecycle interruption;
- keep visual behavior usable when haptics are unavailable.

### 6. Implement platform preconditions

Follow only requirements confirmed for the installed version. Check, as applicable:

- native dependency linking, generated native projects, development-client rebuild,
  and test mocks for React Native or Expo;
- manifest permissions, activity/context needs, capability tiers, and device/OEM
  behavior for Android-family targets;
- package-manager constraints, engine lifecycle, current playability, and hardware
  support for Apple targets;
- plugin registration and explicit native-handle disposal for Flutter;
- feature detection, user activation, promise behavior, and audio fallback for Web.

Respect user and system haptics settings. Do not force a support level in production.
Pair critical feedback with visible or audible state; haptics remain progressive
enhancement.

### 7. Verify

Run the repository's typecheck, build, lint, and relevant tests. Exercise success,
failure, cancellation, disabled-haptics, and unsupported-capability paths.

Simulator/emulator audio can validate invocation, timing, and rough pattern shape.
It cannot validate tactile feel. Finish on supported physical hardware; for Android,
include representative actuator/OEM coverage when fidelity matters.

Report:

- detected platform, package, and resolved version;
- canonical or installed source used;
- chosen preset/composer and why;
- preconditions and fallback behavior;
- simulator/emulator result and physical-hardware result separately;
- any unverified limitation.

## Stop conditions

- No platform/package evidence: ask for it; do not guess imports or API calls.
- Docs unavailable and installed evidence incomplete: provide only a clearly labeled
  approach or pseudocode, not claimed-compiling code.
- Installed platform/version lacks the requested feature: explain the limitation and
  offer a preset, system feedback, visual/audio fallback, or supported upgrade path.
- Hardware, browser, OS, permission, lifecycle, or user settings prevent playback:
  diagnose that boundary before redesigning the interaction.
