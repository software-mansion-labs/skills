# Setup and Advanced Topics

---

## Installation

### Expo (SDK 54+)

The Expo starter template includes the Worklets Babel plugin by default since SDK 54. Install and rebuild native dependencies:

```bash
npm install react-native-worklets
npx expo prebuild
```

### React Native Community CLI

Install the package and add the Babel plugin manually:

```bash
npm install react-native-worklets
```

```js
// babel.config.js
module.exports = {
  presets: ['module:@react-native/babel-preset'],
  plugins: [
    'react-native-worklets/plugin',
  ],
};
```

Then clear the Metro cache and install iOS pods (Android needs no extra steps):

```bash
npm start -- --reset-cache
cd ios && pod install && cd ..
```

### Prerequisites

Worklets requires the New Architecture (Fabric). It is untested on the Legacy Architecture (Paper).

Supported platforms: Android, iOS, macOS, tvOS, visionOS, Web.

---

## Babel Plugin

The Worklets Babel plugin transforms functions marked with `'worklet'` into serializable objects that can run on Worklet Runtimes. It also autoworkletizes callbacks passed to Worklets APIs (`scheduleOnUI`, `scheduleOnRuntime`, etc.) and Reanimated/Gesture Handler hooks.

### What can be a worklet

- Function declarations, function expressions, arrow functions, and object methods with `'worklet';` as the first statement.
- Callbacks passed to autoworkletized APIs (no directive needed).
- All top-level functions in files that start with `'worklet';` at the file level.

### Worklet Context Objects

Object methods lose their `this` binding on the UI thread. Worklet Context Objects preserve it:

```tsx
const counter = {
  __workletContextObject: true,
  count: 0,
  increment() {
    this.count += 1; // `this` is preserved
  },
};
```

Changes on the UI thread are visible only on the UI thread. Changes on the JS thread are visible only on the JS thread.

### Worklet Classes

Hermes doesn't support native classes on Worklet Runtimes. Mark a class with `__workletClass = true` to make it instantiable on the UI thread:

```tsx
class Particle {
  __workletClass = true;
  x = 0;
  y = 0;
  update(dt: number) {
    this.x += dt;
  }
}

scheduleOnUI(() => {
  const p = new Particle();
  p.update(16);
});
```

Limitations: no inheritance, no static members, instances cannot be shared between threads.

### Plugin options

Configure by passing an options object to the Babel plugin:

```js
/** @type {import('react-native-worklets/plugin').PluginOptions} */
const workletsPluginOptions = {
  bundleMode: true,
  strictGlobal: true,
  globals: ['myGlobalVar'],
};

// babel.config.js
plugins: [['react-native-worklets/plugin', workletsPluginOptions]];
```

| Option | Default | Purpose |
|--------|---------|---------|
| `bundleMode` | `false` | Enable Bundle Mode (full bundle on all runtimes) |
| `strictGlobal` | `false` | Stricter access to global variables inside worklets. Recommended. |
| `globals` | `[]` | Identifiers that should not be copied to worklet runtimes |
| `importForwarding.moduleNames` | `[]` | Module names whose imports are forwarded into worklets (Bundle Mode) |
| `importForwarding.relativePaths` | `[]` | Module paths whose relative imports are forwarded into worklets (Bundle Mode) |
| `disableWorkletClasses` | `false` | Disable Worklet Classes (needed for Custom Serializables with `new`) |
| `hermesBytecode` | `false` | Compile worklets to Hermes bytecode ahead of time instead of shipping source (Legacy Eval Mode) |
| `getHBCBinary` | `undefined` | Returns the path to the Hermes bytecode compiler (required by `hermesBytecode`) |
| `extraPlugins` | `[]` | Extra Babel plugins applied when transforming worklet code |
| `extraPresets` | `[]` | Extra Babel presets applied when transforming worklet code |
| `disableInlineStylesWarning` | `false` | Suppress warnings about `.value` access on shared values in Reanimated inline styles |
| `disableSourceMaps` | `false` | Turn off source map generation for worklets |
| `relativeSourceLocation` | `false` | Worklet file paths relative to `process.cwd()` for stable test snapshots |
| `omitNativeOnlyData` | `false` | Smaller bundles for Web builds |
| `substituteWebPlatformChecks` | `false` | Helps tree-shaking for Web builds |

### Pitfalls

- **Worklets are not hoisted.** Using a workletized function before its declaration crashes at runtime.
- **Imported functions need explicit `'worklet'` directive.** Autoworkletization only applies within the same file.
- **Conditional expressions bypass autoworkletization.** Add `'worklet';` to each branch manually.
- **Custom hooks are not autoworkletized.** Only registered APIs trigger autoworkletization; callbacks passed to your own hooks need explicit directives.

---

## Bundle Mode

Bundle Mode gives worklets access to the full JavaScript bundle, allowing imports inside worklets and third-party libraries to run on Worklet Runtimes without patching. Stable since react-native-worklets 0.10.0.

### Setup

Enablement (babel plugin options, Expo/RN CLI metro config helpers, the mandatory metro + metro-runtime patches, verification) is covered end-to-end by the dedicated sub-skill: `../enable-worklets-bundle-mode/SKILL.md`. Follow it instead of setting up manually.

### Imports inside worklets (Import Forwarding)

Runtimes don't share state, so a module-level import used inside a worklet is ambiguous — it could mean the RN Runtime's copy of the module or the Worklet Runtime's. Disambiguate one of three ways:

- **RN Runtime state**: read the value before the worklet and let the closure capture it.
- **Worklet Runtime state, per worklet**: `require('my-library')` inside the worklet body.
- **Worklet Runtime state, app-wide**: list the module in `importForwarding.moduleNames` — its imports are then forwarded into worklet bodies automatically:

```js
const workletsPluginOptions = {
  bundleMode: true,
  strictGlobal: true,
  importForwarding: { moduleNames: ['my-library'] },
};
```

`importForwarding.relativePaths` does the same for relative imports from your own code. The docs describe the `importForwarding` API as temporary, to be replaced with a more robust solution. Docs: https://docs.swmansion.com/react-native-worklets/docs/bundleMode/importForwarding

### Using third-party libraries in worklets

Libraries must be allow-listed via `importForwarding.moduleNames` (the earlier `workletizableModules` option no longer exists in the plugin options).

Libraries that import React Native internals cannot run on Worklet Runtimes (they would load a second RN instance).

### Networking in worklets

Worklet Runtimes offer a simplified `fetch` implementation. Enable it with the `FETCH_PREVIEW_ENABLED` static feature flag; it only takes effect in Bundle Mode and requires the Bundle Mode metro patches.

---

## Feature Flags

Static feature flags go in `package.json` under `worklets.staticFeatureFlags`. Changing one requires `pod install` (iOS) and a native rebuild.

```json
{
  "worklets": {
    "staticFeatureFlags": {
      "FETCH_PREVIEW_ENABLED": true
    }
  }
}
```

| Flag | Default | Purpose |
|------|---------|---------|
| `IOS_DYNAMIC_FRAMERATE_ENABLED` | `true` | Auto-adjust frame rate for expensive animations (falls back from 120fps to 60fps) |
| `FETCH_PREVIEW_ENABLED` | `false` | Enable `fetch` on Worklet Runtimes (only takes effect in Bundle Mode) |
| `ENABLE_CROSS_RUNTIME_STACK_TRACES` | `true` | Stitch stack traces across runtimes in dev builds; can hurt performance in scheduling-heavy code paths |

Static flags are unavailable in Expo Go and in RNRepo prebuilt configurations. Use Expo Prebuild or force source builds instead.

Dynamic flags can be toggled at runtime via `setDynamicFeatureFlag('FLAG_NAME', true)` and read via `getDynamicFeatureFlag('FLAG_NAME')`; no dynamic flags are currently available.

---

## Testing with Jest

### Mock implementation (recommended)

```js
// TypeScript
jest.mock('react-native-worklets', () => require('react-native-worklets/src/mock'));

// JavaScript
jest.mock('react-native-worklets', () => require('react-native-worklets/lib/module/mock'));
```

### Web implementation (v0.8+)

Override the Jest resolver to use the Web implementation instead of native:

```js
// jest.config.js
module.exports = {
  resolver: 'react-native-worklets/jest/resolver',
};
```

---

## Troubleshooting

### "Failed to create a worklet"

The Babel plugin is missing. Add `'react-native-worklets/plugin'` to `babel.config.js` and rebuild.

### "Native part of Worklets doesn't seem to be initialized"

Rebuild the app after installing or upgrading. If using a brownfield app, initialize the native library manually.

### Version mismatch errors

Mismatches between the JS code, the Babel plugin, and the native part all have the same fixes: clear the Metro cache (`npm start -- --reset-cache`, `expo start -c`) and rebuild the app after upgrading. If the issue persists, a dependency bundles worklets transpiled with an older Babel plugin version. On Expo Go, use the exact worklets version bundled with the SDK.

### "TypeError: right operand of 'in' is not an object" / "Cannot read property 'createSerializableString' of undefined"

Expo apps disable `inlineRequires` by default, which breaks Worklets initialization. Enable `inlineRequires` in `metro.config.js`.

### "Tried to modify key of an object which has been converted to a serializable"

The object was captured in a worklet's closure and later mutated. In dev builds, captured objects are frozen to surface this mistake. Solutions:
- Use `useSharedValue` for values that change over time.
- Destructure only the needed properties into local variables before the worklet captures them.

### "Tried to synchronously call a non-worklet function on the UI thread"

The called function lacks a `'worklet';` directive. Either add `'worklet';` to make it run on the UI thread, or wrap the call with `scheduleOnRN(fn)` to run it on the JS thread.

---

## Compatibility

Worklets supports at least the last three minor versions of React Native. Latest patch of each minor; full table: https://docs.swmansion.com/react-native-worklets/docs/guides/compatibility.

When Reanimated is installed it additionally pins a narrow worklets range — check both compatibility tables before upgrading.
