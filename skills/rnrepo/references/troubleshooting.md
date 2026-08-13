---
name: rnrepo-troubleshooting
description: "Troubleshooting RNRepo build issues: C++ debug/release mismatch with react-native-worklets/reanimated, duplicate .so files, empty repository list, FAIL_ON_PROJECT_REPOS Gradle errors, iOS Xcode version mismatch, Expo EAS fingerprint failures, and how to verify the plugin is active."
---

# RNRepo Troubleshooting

## Verifying the plugin is active

### Android

Run with `--info` and search for `[📦 RNRepo]` in the output:
```bash
./gradlew :app:assembleDebug --info
```

Expected log:
```
[📦 RNRepo] Found the following supported prebuilt packages:
📦 react-native-safe-area-context@5.7.0
```

If the plugin is working, you will **not** see `compileDebugKotlin` tasks for the listed packages.

### iOS

During `pod install` or build, look for:
```
[📦 RNRepo] ⬇ Downloaded from Maven...
[📦 RNRepo] • react-native-safe-area-context
```

---

## No supported packages found

Log shows:
```
[📦 RNRepo] Found the following supported prebuilt packages: None
```

Causes:
1. The RNRepo plugin is applied **before** the Maven repository is defined. Move the `maven { url "https://packages.rnrepo.org/releases" }` block to `allprojects` in root `build.gradle` and ensure it comes before the plugin is applied.
2. All matching packages are in the deny list.
3. The installed library or React Native version has no prebuild yet (expected — RNRepo falls back to source silently).

---

## C++ debug/release mismatch (reanimated + worklets)

Symptom — `--info` logs show:
```
[📦 RNRepo] react-native-worklets is supported, checking if all packages depending on it are supported.
[📦 RNRepo] react-native-reanimated depending on react-native-worklets is not available as a prebuild, building react-native-worklets from sources.
```

Cause: One dependency uses a prebuilt (release variant) while another that depends on it builds from source (debug variant). The ABI differs for C++ libraries.

Solutions:
1. Keep build variants consistent — all dependencies should share the same build type.
2. Add the problematic libraries to the `denyList` in `rnrepo.config.json` to force source builds.

---

## Duplicate `.so` file conflict

Error:
```
DuplicateRelativeFileException: 2 files found with path 'lib/arm64-v8a/libworklets.so'
```

Cause: A library built from source includes native code from a provider library (e.g., `react-native-worklets`), causing a collision with the prebuilt version.

The plugin automatically adds `pickFirsts` for known provider libraries:
- `react-native-worklets` → `libworklets.so`
- `react-native-nitro-modules` → `libNitroModules.so`

If another library triggers this, add it to the deny list and report the issue on the [RNRepo GitHub](https://github.com/software-mansion/rnrepo/issues).

---

## iOS: Xcode version mismatch (Xcode < 26)

Error:
```
Undefined symbols for architecture arm64:
  "_OBJC_CLASS_$_UITabAccessory", referenced from: ...
```

Cause: Prebuilt xcframeworks are compiled with Xcode 26. Building with an older Xcode version causes linker failures for packages that use newer APIs.

Identify the culprit from the `.o` filename in the error:
```bash
find ./node_modules -path "*RNSTabsBottomAccessoryHelper*"
# → node_modules/react-native-screens/ios/...
```

Solutions (in order of preference):
1. **Upgrade Xcode** to 26 or newer.
2. **Downgrade the library** to a version prebuilt with an older Xcode (if available in the registry).
3. **Add to deny list** under `ios` in `rnrepo.config.json` to force source compilation.

---

## Expo EAS fingerprint mismatch

Cause: The CocoaPods plugin caches xcframeworks in `node_modules/<pkg>/.rnrepo-cache`, changing the directory hash between local and CI environments.

Fix: Add to `.fingerprintignore`:
```text
**/.rnrepo-cache
```

---

## Plugin running during non-build tasks (test, clean, lint)

The plugin skips automatically for tasks it recognizes as non-build: `test`, `signing`, `clean`, `clear`, `init`, `dependencies`, `tasks`, `projects`, `connected`, `device`, `lint`, `check`, `properties`, `help`.

If it still runs when it shouldn't:
```bash
DISABLE_RNREPO=true ./gradlew clean
DISABLE_RNREPO=true ./gradlew test
```

---

## Always target the app module explicitly

Use `./gradlew :app:assembleDebug`, not `./gradlew assembleDebug`. Running a bare task without the module prefix will invoke all subprojects build tasks and can result in longer build times and unexpected behavior.

```bash
# Correct
./gradlew :app:assembleDebug

# Avoid
./gradlew assembleDebug
```

---

## Getting detailed Gradle logs

```bash
# Info level (shows RNRepo decisions)
./gradlew :app:assembleDebug --info

# Scan (full task graph + timing, uploads to scans.gradle.com)
./gradlew :app:assembleDebug --scan
```

Use `--scan` before reporting a bug — it captures all task executions and timing, making it easy to see which packages were and weren't substituted.
