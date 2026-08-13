---
name: rnrepo-installation
description: "How to install and set up RNRepo in a React Native or Expo project. Covers Expo CNG (config plugin), standard React Native (manual Gradle + Podfile edits), Android-only, and iOS-only setups."
---

# RNRepo Installation

## Prerequisites

- React Native **New Architecture** enabled.
- React Native version: `0.77.3`, `0.78.3`, `0.79.7`, or `>= 0.80.0`.

---

## Expo Prebuild (CNG)

For projects that use `expo prebuild` to generate native directories.

### 1. Install the plugin

```bash
npx expo install @rnrepo/expo-config-plugin
```

In a monorepo with hoisted `node_modules`, also install the build tools explicitly:

```bash
npx expo install @rnrepo/build-tools
```

### 2. Register in `app.config.ts` (or `app.json` / `app.config.js`)

```diff
{
  "expo": {
    "plugins": [
+     "@rnrepo/expo-config-plugin"
    ]
  }
}
```

### 3. Add `.fingerprintignore` (if using Expo Fingerprint / EAS)

```text
**/.rnrepo-cache
```

That's it — run `expo prebuild` and the plugin configures both Android and iOS automatically.

---

## Standard React Native (manual native folders)

### 1. Install build tools

```bash
npm install @rnrepo/build-tools@latest
```

### 2. Android — edit `android/build.gradle`

Add the RNRepo Gradle plugin to the `buildscript.dependencies` block:

```diff
buildscript {
  dependencies {
    // ... existing classpath entries
+   def rnrepoDir = new File(
+     providers.exec {
+       workingDir(rootDir)
+       commandLine("node", "--print", "require.resolve('@rnrepo/build-tools/package.json')")
+     }.standardOutput.asText.get().trim()
+   ).getParentFile().absolutePath
+   classpath fileTree(dir: "${rnrepoDir}/gradle-plugin/build/libs", include: ["prebuilds-plugin.jar"])
  }
}
```

### 3. Android — edit `android/app/build.gradle`

Apply the plugin after the existing `apply plugin` lines:

```diff
apply plugin: "com.android.application"
apply plugin: "org.jetbrains.kotlin.android"
apply plugin: "com.facebook.react"
+ apply plugin: "org.rnrepo.tools.prebuilds-plugin"
```

### 4. iOS — edit `ios/Podfile`

Add the `require` at the very top and a `post_install` hook inside your target:

```diff
+ require Pod::Executable.execute_command('node', ['-p',
+   'require.resolve(
+   "@rnrepo/build-tools/cocoapods-plugin/lib/plugin.rb",
+   {paths: [process.argv[1]]},
+ )', __dir__]).strip

# ... rest of Podfile

target 'YourProjectName' do
  post_install do |installer|
    # ... existing post_install hooks
+   rnrepo_post_install(installer)
  end
end
```

---

## Verifying the setup

### Android

Run `npm run android` and look for RNRepo log lines during Gradle preparation:

```
[📦 RNRepo] Found the following supported prebuilt packages:
📦 react-native-safe-area-context@5.7.0
```

If working: **no** `compileDebugKotlin` tasks for substituted packages.

### iOS

Run a build and look for RNRepo log lines during `pod install`:

```
[📦 RNRepo] Total React Native dependencies detected: 1
[📦 RNRepo] ⬇ Downloaded from Maven...
[📦 RNRepo] • react-native-safe-area-context
[📦 RNRepo] Added build phase to <package-name>
```

If working: Xcode compiles **only** JSI glue code (e.g., `safeareacontextJSI-generated.cpp`), not library source files.
