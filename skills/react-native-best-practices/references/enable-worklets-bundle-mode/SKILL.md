---
name: enable-worklets-bundle-mode
description:
  Enable react-native-worklets Bundle Mode (imports inside worklets, third party libraries on worklet runtimes) in an Expo, RN
  CLI or brownfield React Native app, including the mandatory metro/metro-runtime patches. Use when the user asks to enable
  Bundle Mode, wants to import npm libraries inside worklets or offload JS code to worklet runtimes, or hits "Failed to get the
  SHA-1" errors in a worklets Bundle Mode project / missing Fast Refresh for worklet code. Covers enabling only — not disabling;
  upgrades react-native-worklets to a compatible version when needed, but never installs it from scratch and never touches
  Reanimated.
compatibility: Requires network access to raw.githubusercontent.com and api.github.com to fetch metro patches
metadata:
  author: Tomasz Żelawski
  version: 1.0.0
---

# Enable Worklets Bundle Mode

Bundle Mode lets worklets access the entire JS bundle (imports inside worklets, etc.). Stable since `react-native-worklets`
0.10.0. Official setup docs: https://docs.swmansion.com/react-native-worklets/docs/bundleMode/setup

Bundle Mode has **three** parts: the babel plugin option, the metro config helper, and the metro + metro-runtime patches. The
patches are mandatory - the docs label them "recommended", but without them DX is very poor.

Reference implementation: https://github.com/software-mansion-labs/Bundle-Mode-showcase-app

## 0. Detect current state and environment

Bundle Mode is already ON when all three hold:

- `babel.config.js`: `react-native-worklets/plugin` has `bundleMode: true`.
- `metro.config.js`: imports from `react-native-worklets/bundleMode` (`bundleModeMetroConfig` or `getBundleModeMetroConfig`).
- Patches applied — layout-aware check (hoisted, pnpm isolated, and workspace layouts; run at the workspace root):
  `find node_modules -path '*metro/src/node-haste/DependencyGraph.js' -exec grep -l react-native-worklets {} +` and
  `find node_modules -path '*metro-runtime/src/modules/HMRClient.js' -exec grep -l __workletsModuleProxy {} +` each match at
  least one file. The copy that matters is the one `react-native` resolves — nested under `node_modules/react-native/` if
  present there.

If all three hold, report that Bundle Mode is already enabled and stop. If only some hold, that's a partial (broken) setup —
report which parts are missing, then complete them with the steps below.

Environment facts you need:

- Confirm `react-native-worklets` is a dependency. If it's missing entirely, stop and tell the user — this skill upgrades
  worklets when needed but does not introduce it to an app.
- Bundle Mode needs worklets >= 0.10.0. If the installed version is older, upgrade it to the newest applicable version:
  1. Read the installed `react-native` and (if present) `react-native-reanimated` versions.
  2. Fetch both compatibility tables and pick the HIGHEST worklets minor that is (a) compatible with the app's RN version per
     https://docs.swmansion.com/react-native-worklets/docs/guides/compatibility and (b) accepted by the installed reanimated
     minor per https://docs.swmansion.com/react-native-reanimated/docs/guides/compatibility/ — reanimated pins narrow worklets
     ranges (e.g. 4.5.x accepts 0.10.x–0.11.x while 4.3.x only accepts 0.8.x), so both tables must agree on the pick.
  3. If the best version satisfying both is still < 0.10.0 (old reanimated), stop and explain: Bundle Mode requires a
     reanimated upgrade first, and that is the user's decision — never upgrade reanimated yourself.
  4. Install the picked version with the detected package manager (`npx expo install react-native-worklets@{ver}` on Expo,
     plain add elsewhere), then re-check the resolved version before continuing.
- Detect the package manager from the lockfile: `yarn.lock` (Yarn 2+ berry when `packageManager: yarn@2+` is set or `.yarnrc.yml`
  with `yarnPath`/`nodeLinker` exists, otherwise Yarn 1 classic), `package-lock.json` (npm), `bun.lockb` / `bun.lock` (bun),
  `pnpm-lock.yaml` (pnpm). In a workspace repo look for the lockfile upward from the app dir — patch registration happens at the
  workspace root.
- Get the **installed metro version** — patches are version-specific:
  ```bash
  node -e "console.log(require('metro/package.json').version)"
  node -e "console.log(require('metro-runtime/package.json').version)"
  ```
  Under pnpm's isolated layout these `require`s fail with MODULE_NOT_FOUND — use `pnpm why metro` / `pnpm why metro-runtime`, or
  read the `version` field of the `package.json` next to the files located by the find commands above.
- Is it Expo or RN community CLI? (Different metro helper — see step 2.) An Expo app may have no `babel.config.js` /
  `metro.config.js` at all — that's normal, not a broken state; steps 1–2 generate them.
- Does another library remap the bare `react-native` specifier in `metro.config.js`? Check for uniwind
  (`grep -ns "uniwind" metro.config.js package.json`); NativeWind is reported to do the same. If present, the plain step 2 setup
  crashes the app at startup — use [references/uniwind-remap-workaround.md](references/uniwind-remap-workaround.md) in step 2
  instead.

## 1. Babel plugin

In `babel.config.js`, add the worklets plugin with `bundleMode: true`. `strictGlobal: true` is optional but recommended by the
docs.

If the worklets plugin is already present, just add the options to it — do not add a second copy. The plugin should stay last in
the `plugins` array, and the options go on the plugin entry, never on a preset.

Always use the typed-const form shown below, including when merging into an existing config: declare `workletsPluginOptions`
with the `/** @type {import('react-native-worklets/plugin').PluginOptions} */` JSDoc annotation and reference it from the plugin
entry — do not inline an untyped options object. Keep the JSDoc line even in codebases with a no-comments convention: it is a
type annotation (editor completion + typo checking for the options), not a prose comment.

RN community CLI (`babel.config.js`):

```js
/** @type {import('react-native-worklets/plugin').PluginOptions} */
const workletsPluginOptions = {
  bundleMode: true,
  strictGlobal: true,
};

module.exports = {
  presets: ['module:@react-native/babel-preset'],
  plugins: [['react-native-worklets/plugin', workletsPluginOptions]],
};
```

Expo — if the app has no `babel.config.js`, generate the default one first
(https://docs.expo.dev/versions/latest/config/babel/):

```bash
npx expo customize babel.config.js
```

The generated file is a function returning `{ presets: ['babel-preset-expo'] }`. Keep the preset and add a `plugins` array to
the returned object:

```js
/** @type {import('react-native-worklets/plugin').PluginOptions} */
const workletsPluginOptions = {
  bundleMode: true,
  strictGlobal: true,
};

module.exports = function (api) {
  api.cache(true);
  return {
    presets: ['babel-preset-expo'],
    plugins: [['react-native-worklets/plugin', workletsPluginOptions]],
  };
};
```

## 2. Metro config

`react-native-worklets/bundleMode` exports two helpers — pick by project type:

- **RN community CLI** → `bundleModeMetroConfig` (a plain config object, merge it in).
- **Expo** → `getBundleModeMetroConfig(config)` (a function that takes and returns a config).

RN community CLI (`metro.config.js`):

```js
const { getDefaultConfig, mergeConfig } = require('@react-native/metro-config');
const { bundleModeMetroConfig } = require('react-native-worklets/bundleMode');

const config = {
  // your existing metro config
};

module.exports = mergeConfig(getDefaultConfig(__dirname), bundleModeMetroConfig, config);
```

Order matters: put `bundleModeMetroConfig` before the user `config` so the project can still override. `bundleModeMetroConfig`
installs a custom `resolver.resolveRequest` (RN + TurboModuleRegistry shims, `react-native-worklets/.worklets` resolution) and a
`serializer.createModuleIdFactory`.

If the existing config defines its own `resolver.resolveRequest` (svg transformers, monorepo resolvers), do NOT rely on merge
order — `mergeConfig` would replace Bundle Mode's resolver wholesale and silently disable it. Use the function form instead,
which chains the existing resolver (see `getBundleModeMetroConfig` in `react-native-worklets/bundleMode/index.js`):
`module.exports = getBundleModeMetroConfig(mergeConfig(getDefaultConfig(__dirname), config));` — it works on RN CLI configs too.

Expo — if the app has no `metro.config.js`, generate the default one first
(https://docs.expo.dev/guides/customizing-metro/):

```bash
npx expo customize metro.config.js
```

Then wrap the config with the helper (`metro.config.js`):

```js
const { getDefaultConfig } = require('expo/metro-config');
const { getBundleModeMetroConfig } = require('react-native-worklets/bundleMode');

let config = getDefaultConfig(__dirname);
config = getBundleModeMetroConfig(config);
module.exports = config;
```

If step 0 found uniwind (or any resolver remapping the bare `react-native` specifier), do NOT use the plain setups above — they
form a require cycle between Bundle Mode's `react-native` shim and the remapping library, crashing the app at startup
([issue #9817](https://github.com/software-mansion/react-native-reanimated/issues/9817)). Read
[references/uniwind-remap-workaround.md](references/uniwind-remap-workaround.md) and apply the guarded config from there.

## 3. Metro patches (do not skip)

Two patches against `metro` and `metro-runtime`:

- **metro** (`src/node-haste/DependencyGraph.js`): synchronously indexes the `react-native-worklets/.worklets` virtual modules
  (short-circuits `getOrComputeSha1`) → fixes "Failed to get the SHA-1" errors / repeated-reload requirement for new modules.
- **metro-runtime** (`src/modules/HMRClient.js`): calls `global.__workletsModuleProxy?.propagateModuleUpdate(...)` on HMR inject
  → enables Fast Refresh for worklet runtimes (without it, worklet changes need full app reloads).

Expo SDK 57+ also ships an `@expo/metro` package that looks like a vendored metro fork which would bypass these patches — it
isn't: it's a thin re-export shim (`@expo/metro/metro/index.js` is `module.exports = require("metro")`, and its
`DependencyGraph` entry forwards to `metro/private/node-haste/DependencyGraph`). Patch the real `metro` / `metro-runtime`
packages as usual; do not spend time investigating `@expo/metro`.

Patch files are NOT shipped in the npm package — fetch them from the worklets repo, matching your installed metro version
exactly. Base URL (referred to as `$BASE` in the reference files):
`https://raw.githubusercontent.com/software-mansion/react-native-reanimated/main/packages/react-native-worklets/bundleMode/patches`

Always discover the currently available filenames+hashes via the GitHub contents API first — do not trust the version list
below, it rots:
`https://api.github.com/repos/software-mansion/react-native-reanimated/contents/packages/react-native-worklets/bundleMode/patches/yarn/metro`
If the API returns 404, the patches directory moved (the worklets package is being decoupled from the reanimated repo) — search
the software-mansion GitHub org for the new `react-native-worklets` location before giving up.

Formats under `$BASE/{yarn,patch-package}/...` (known versions as of 2026-07: 0.82.4, 0.82.5, 0.83.2, 0.84.4):

- `$BASE/yarn/metro/metro-npm-{ver}-{hash}.patch`
- `$BASE/yarn/metro-runtime/metro-runtime-npm-{ver}-{hash}.patch`
- `$BASE/patch-package/...` (npm / Yarn classic / bun / pnpm style)

Download patch files into a disposable scratch directory. If the environment blocks the usual scratch locations (e.g. creating
`.claude/tmp` is denied), use a throwaway project-root directory like `.bundle-mode-tmp/` and delete it when done — do not
block on the choice of scratch location.

How to apply depends on the package manager detected in step 0. Read ONLY the matching reference file and follow it:

- Yarn 2+ (berry) → [references/patching-yarn-berry.md](references/patching-yarn-berry.md)
- npm or Yarn 1 classic → [references/patching-patch-package.md](references/patching-patch-package.md)
- bun → [references/patching-bun.md](references/patching-bun.md)
- pnpm → [references/patching-pnpm.md](references/patching-pnpm.md) (not upstream-verified — the step 4 checks are mandatory).

## 4. Verify

```bash
find node_modules -path '*metro/src/node-haste/DependencyGraph.js' -exec grep -l react-native-worklets {} +   # expect >=1 file
find node_modules -path '*metro-runtime/src/modules/HMRClient.js' -exec grep -l __workletsModuleProxy {} +    # expect >=1 file
```

These are layout-aware (hoisted, pnpm isolated, npm workspaces) — run them at the workspace root. The copy that must be patched
is the one `react-native` resolves: if `node_modules/react-native/node_modules/metro-runtime` exists, that copy is the one that
matters, not the root one.

Then start Metro with a clean cache (the babel/metro changes need it):

```bash
yarn start --reset-cache    # npm start -- --reset-cache; Expo: npx expo start --clear
```

Strongest check without booting the app — build a dev bundle and look for the bundle-mode virtual modules in it:

```bash
yarn react-native bundle --platform android --dev true --entry-file index.js \
  --bundle-output .claude/tmp/bundle-mode-test.js --reset-cache
grep -c "react-native-worklets/\.worklets" .claude/tmp/bundle-mode-test.js   # expect >0
```

(Expect >0 in an app that already contains worklets; in a worklet-free app a 0 here is inconclusive, not proof of failure —
add a trivial worklet before treating it as a signal.)

(Expo equivalent: `npx expo export --platform android --dev` and grep the output bundle under `dist/`.)

## Example

User says "enable bundle mode in my Expo app". Step 0 finds: `bun.lock` (bun), metro 0.84.4, no `babel.config.js` /
`metro.config.js` (normal for Expo), no `react-native`-remapping resolver. Actions: generate both configs with
`npx expo customize`, add the plugin options (step 1), wrap with `getBundleModeMetroConfig` (step 2), apply the two patches per
[references/patching-bun.md](references/patching-bun.md) (step 3), verify the greps and restart with `npx expo start --clear`
(step 4). Result: imports work inside worklets, no SHA-1 errors, Fast Refresh reaches worklet runtimes.

## Gotchas

- Wrong metro patch version → `yarn install` fails to apply the patch, or you still see `Failed to get the SHA-1`. Re-check
  `metro` version and pick the matching patch.
- Piping npm/yarn commands (`npm ci 2>&1 | tail`) masks their exit code — a failed install can look successful. Check the real
  exit status before trusting greps on `node_modules` (stale trees from a previous package manager pass verification checks
  misleadingly). Note `${PIPESTATUS[0]}` is bash-only — zsh spells it `$pipestatus` — so an empty result does not mean success;
  re-run unpiped when in doubt.
- Forgetting `--reset-cache` after enabling → stale transform, looks like nothing changed.
- Only the babel option + metro config but no patches → SHA-1 errors / repeated reloads / no Fast Refresh on worklet runtimes.
- After upgrading React Native / metro, the patch hashes change — refetch matching patches.
- Startup crash `RangeError: Maximum call stack size exceeded` + `Invariant Violation: "main" has not been registered` right
  after enabling → another resolver remaps `react-native` (uniwind, NativeWind) and cycles with the Bundle Mode shim — see
  [references/uniwind-remap-workaround.md](references/uniwind-remap-workaround.md).
- Expo + pnpm: `npx expo customize babel.config.js` may crash running `pnpm add --dev babel-preset-expo` (`--dev` is not a
  valid pnpm add flag) — and the preset genuinely won't resolve under pnpm's isolated layout until installed. Finish manually
  with `pnpm add -D babel-preset-expo`.
