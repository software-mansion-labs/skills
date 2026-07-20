# Bundle Mode metro patches — bun

Applies the metro + metro-runtime Bundle Mode patches with bun's native patching (works with monorepos and
`linker=isolated`). Uses the **patch-package-format** patch files.

```bash
BASE=https://raw.githubusercontent.com/software-mansion/react-native-reanimated/main/packages/react-native-worklets/bundleMode/patches
```

In non-interactive shells `bun` may not be on PATH (it typically lives at `~/.bun/bin/bun`) — a patch attempt can then
silently no-op. Verify `command -v bun` (or use the full path) before starting.

Find your installed versions with `bun why metro --top` and `bun why metro-runtime --top`, then for each package (substitute
the version; `%2B` is the URL-encoded `+` in the filename):

```bash
bun patch metro
curl -fsSL "$BASE/patch-package/metro/metro%2B0.84.4.patch" | git apply
bun patch --commit 'node_modules/metro'
```

```bash
bun patch metro-runtime
curl -fsSL "$BASE/patch-package/metro-runtime/metro-runtime%2B0.84.4.patch" | git apply
bun patch --commit 'node_modules/metro-runtime'
```

`bun patch --commit` records the patch in `patchedDependencies` in `package.json` so it re-applies on installs.
