# Bundle Mode + uniwind (and other `react-native`-remapping resolvers)

Source: https://github.com/software-mansion/react-native-reanimated/issues/9817 (root cause in the issue body, workaround in
the maintainer comment). Not fixable inside worklets — Metro provides no way to prevent resolution cycles — so the fix lives in
the app's `metro.config.js`.

## The problem

Bundle Mode's resolver redirects every `react-native` import to its own shim (`bundleMode/shims/reactNativeShim.js`), and the
shim re-acquires the real module with `require('react-native')` — the bare specifier, which re-enters the resolver chain. Any
other resolver that remaps that bare specifier intercepts the shim's require: uniwind rewrites `import 'react-native'` →
`uniwind/components`, whose own internal `require('react-native')` Bundle Mode redirects back to the shim. The result is a
`shim ↔ uniwind/components` require cycle. NativeWind is reported to trigger the same class of problem.

Symptoms — startup crash before `AppRegistry.registerComponent` runs:

```
ERROR  [runtime not ready]: RangeError: Maximum call stack size exceeded (native stack depth)
ERROR  [runtime not ready]: Invariant Violation: "main" has not been registered.
```

## The workaround (from the maintainers)

Two parts: apply `withUniwindConfig` BEFORE `getBundleModeMetroConfig` (many broken setups have it the other way around), then
install an outermost resolver guard that pins `react-native` requests originating inside the uniwind package to the real
`react-native` path:

```js
const path = require('path');

config = withUniwindConfig(config, {
  cssEntryFile: './global.css',
  dtsFile: './uniwind-types.d.ts',
});
config = getBundleModeMetroConfig(config);

const uniwindDir = path.dirname(require.resolve('uniwind/package.json')) + path.sep;
const realReactNativePath = require.resolve('react-native');
const wrappedResolveRequest = config.resolver.resolveRequest;
config.resolver.resolveRequest = (context, moduleName, platform) => {
  if (
    moduleName === 'react-native' &&
    typeof context.originModulePath === 'string' &&
    context.originModulePath.startsWith(uniwindDir)
  ) {
    return { type: 'sourceFile', filePath: realReactNativePath };
  }
  return (wrappedResolveRequest || context.resolveRequest)(context, moduleName, platform);
};

module.exports = config;
```

Keep the app's existing `withUniwindConfig` import and options — only the ordering and the guard are new.

Why it works (inferred from the issue's root-cause description): app code's `react-native` still goes shim →
`uniwind/components` (both libraries' rewrites stay active), but `uniwind/components`' internal `require('react-native')` hits
the guard and resolves straight to the real `react-native` file — a concrete path that never re-enters resolution, so the cycle
cannot form.

## Caveats

- Non-standard setups (monorepos): verify `uniwindDir` and `realReactNativePath` point at the copies the app actually uses —
  e.g. `require.resolve('react-native', { paths: [__dirname] })` from the app package.
- Other remapping libraries (e.g. NativeWind): same guard pattern, but scope the origin check to that library's package dir
  instead of uniwind's. This variant is untested — the issue thread only validates the uniwind one.
- Watch the upstream issue: it proposes a sentinel-specifier fix inside worklets; if that ships, this guard becomes
  unnecessary (check `node_modules/react-native-worklets/bundleMode/shims/reactNativeShim.js` — if it no longer contains
  `require('react-native')`, the fix landed).
