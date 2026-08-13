# Bundle Mode metro patches — npm / Yarn 1 classic (patch-package)

Applies the metro + metro-runtime Bundle Mode patches with `patch-package`.
Use the **patch-package-format** patch files (named `metro+{ver}.patch` /
`metro-runtime+{ver}.patch`). With Yarn classic the flow is identical — swap
npm commands for their yarn equivalents (`yarn patch-package` instead of
`npx patch-package`, `yarn install` instead of `npm ci`).

```bash
BASE=https://raw.githubusercontent.com/software-mansion/react-native-reanimated/main/packages/react-native-worklets/bundleMode/patches
```

Download the patch files from `$BASE/patch-package/{metro,metro-runtime}/` into
`./patches`, then run `npx patch-package`. Make sure a `prepare` or
`postinstall` script runs `patch-package` so it re-applies after installs.

**Check where each package physically sits first**
(`ls -d node_modules/metro node_modules/metro-runtime node_modules/*/node_modules/metro-runtime node_modules/@*/*/node_modules/metro-runtime`).
In npm workspace repos `metro-runtime` is often NOT hoisted to root (nested
under `react-native`, `metro`, etc.) and the stock patch fails with "Patch file
found for package metro-runtime which is not present at
node_modules/metro-runtime".

Do NOT reach for `npm dedupe` to force hoisting — a full-tree dedupe
re-resolution can corrupt the lockfile's peer-dep layout (seen: it hoisted
`fdir` away from its `picomatch@^3||^4` peer, making every subsequent `npm ci`
fail validation). Instead patch the nested copy that actually ends up in the app
bundle — react-native's own (`react-native/Libraries/Utilities/HMRClient.js` is
what imports `metro-runtime/src/modules/HMRClient`):

```bash
# apply the stock patch content inside react-native (paths line up with -p1 from there)
cd node_modules/react-native && patch -p1 < ../../patches/metro-runtime+0.84.4.patch && cd ../..
# generate the nested patch (note the dir/package syntax → '++' in the filename)
npx patch-package react-native/metro-runtime
# the root-level metro-runtime patch is replaced by the nested one
rm patches/metro-runtime+0.84.4.patch
npx patch-package   # verify: metro ✔ and react-native/metro-runtime ✔
```

The `metro` package itself is normally hoisted, so `metro+<ver>.patch` applies
at root as-is.
