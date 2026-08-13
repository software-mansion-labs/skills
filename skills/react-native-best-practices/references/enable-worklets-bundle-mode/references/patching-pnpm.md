# Bundle Mode metro patches — pnpm

Applies the metro + metro-runtime Bundle Mode patches with pnpm's native patching (https://pnpm.io/cli/patch). Uses the
**patch-package-format** patch files.

NOT upstream-verified: the worklets repo ships no pnpm instructions (its patches README asks for contributions). This flow
follows pnpm's documented mechanism; the `patch -p3` re-targeting is needed because the patch-package-format files carry
`node_modules/{pkg}/`-prefixed paths while pnpm patch dirs are package-relative (`-p3` strips `a/node_modules/{pkg}/`), and that
re-targeting has been validated locally against the 0.84.4 metro patch. Treat the verification greps in step 4 of the skill as
mandatory.

```bash
BASE=https://raw.githubusercontent.com/software-mansion/react-native-reanimated/main/packages/react-native-worklets/bundleMode/patches
```

Find your installed versions (`pnpm why metro`, `pnpm why metro-runtime`), then for each package (substitute the version; `%2B`
is the URL-encoded `+` in the filename; the `--edit-dir` must not exist beforehand):

```bash
rm -rf .claude/tmp/metro-patch
pnpm patch metro@0.84.4 --edit-dir=.claude/tmp/metro-patch
curl -fsSL "$BASE/patch-package/metro/metro%2B0.84.4.patch" | patch -p3 -d .claude/tmp/metro-patch
pnpm patch-commit .claude/tmp/metro-patch
```

```bash
rm -rf .claude/tmp/metro-runtime-patch
pnpm patch metro-runtime@0.84.4 --edit-dir=.claude/tmp/metro-runtime-patch
curl -fsSL "$BASE/patch-package/metro-runtime/metro-runtime%2B0.84.4.patch" | patch -p3 -d .claude/tmp/metro-runtime-patch
pnpm patch-commit .claude/tmp/metro-runtime-patch
```

`pnpm patch-commit` writes a package-relative patch under `patches/` and registers it in `pnpm.patchedDependencies` in
`package.json`, so it re-applies on every install. Do NOT point `patchedDependencies` at the downloaded patch-package-format
files directly — their `node_modules/{pkg}/` path prefixes don't match pnpm's package-relative patch format. Clean up the
`--edit-dir` temp directories afterwards.
