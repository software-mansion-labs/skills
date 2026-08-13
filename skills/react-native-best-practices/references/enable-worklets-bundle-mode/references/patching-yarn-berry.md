# Bundle Mode metro patches — Yarn 2+ (berry)

Applies the metro + metro-runtime Bundle Mode patches with Yarn Modern's
builtin `patch:` protocol. Use the **yarn-format** patch files.

```bash
BASE=https://raw.githubusercontent.com/software-mansion/react-native-reanimated/main/packages/react-native-worklets/bundleMode/patches
```

Place the yarn-format patch files in `.yarn/patches/` and reference them via
`resolutions` with the `patch:` protocol, then install. (`.yarn/patches/` is
convention, not a requirement — if writing to `.yarn/` is blocked, a
project-root `patches/` directory with `#~/patches/...` locators works
identically; just make sure it isn't gitignored.) Example for metro 0.84.4:

```bash
mkdir -p .yarn/patches
curl -fsSL -o .yarn/patches/metro-npm-0.84.4-68d21d57b4.patch \
  "$BASE/yarn/metro/metro-npm-0.84.4-68d21d57b4.patch"
curl -fsSL -o .yarn/patches/metro-runtime-npm-0.84.4-9533293c73.patch \
  "$BASE/yarn/metro-runtime/metro-runtime-npm-0.84.4-9533293c73.patch"
```

`package.json` (root of the yarn project — in a workspace repo that's the
monorepo root, not the app package; the `npm%3A{ver}` in the value is
`npm:{ver}` URL-encoded, `~/` = project root). First check the lockfile for
OTHER metro versions in the tree (`grep -n 'metro@npm' yarn.lock`):

- Only one metro version → name-only keys are fine and cover every descriptor:

```json
"resolutions": {
  "metro": "patch:metro@npm%3A0.84.4#~/.yarn/patches/metro-npm-0.84.4-68d21d57b4.patch",
  "metro-runtime": "patch:metro-runtime@npm%3A0.84.4#~/.yarn/patches/metro-runtime-npm-0.84.4-9533293c73.patch"
}
```

- Multiple metro versions (e.g. a web workspace pulling metro 0.80.x) →
  name-only keys would force ALL of them onto the patched version and can break
  the other consumer. Scope the keys to the exact descriptors that resolve to
  your target version (find them in the lockfile entry, e.g.
  `"metro@npm:0.84.4, metro@npm:^0.84.3":`):

```json
"resolutions": {
  "metro@npm:0.84.4": "patch:metro@npm%3A0.84.4#~/.yarn/patches/metro-npm-0.84.4-68d21d57b4.patch",
  "metro@npm:^0.84.3": "patch:metro@npm%3A0.84.4#~/.yarn/patches/metro-npm-0.84.4-68d21d57b4.patch",
  "metro-runtime@npm:0.84.4": "patch:metro-runtime@npm%3A0.84.4#~/.yarn/patches/metro-runtime-npm-0.84.4-9533293c73.patch",
  "metro-runtime@npm:^0.84.3": "patch:metro-runtime@npm%3A0.84.4#~/.yarn/patches/metro-runtime-npm-0.84.4-9533293c73.patch"
}
```

Then `yarn install`. Yarn validates the patch applies — a clean install means
the versions matched.

The `~/` locator form is verified on Yarn 4 (the reanimated monorepo itself
uses it); on Yarn 3, if `~/` is rejected, use a `./`-relative path from the
manifest instead (Yarn 3 acceptance unverified).
