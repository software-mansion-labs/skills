# TypeGPU Project Setup

## Fastest path: TypeGPU CLI

For a **new project** or when **adding TypeGPU to an existing one**, prefer the CLI over the manual steps below - it handles the install, build plugin, and types in one go:

```sh
npx typegpu@latest                 # scaffold a new project (interactive)
npx typegpu@latest my-app --yes    # non-interactive, defaults
npx typegpu@latest my-app --yes --template vite-react --addons @typegpu/sdf,@typegpu/noise
npx typegpu@latest --enhance       # retrofit TypeGPU into the current project
```

Templates: `vite-bare`, `vite-complex`, `vite-react`, `expo-bare`. `--enhance` installs `typegpu`, wires up the build plugin, adds `@webgpu/types`, and can install the TypeGPU AI skill.

The manual steps below are useful for understanding what the CLI sets up, and for repairing a broken setup.

## 1. Install TypeGPU

```sh
npm install typegpu
# pnpm add typegpu / yarn add typegpu
```

## 2. WebGPU types

TypeScript doesn't ship WebGPU types by default:

```sh
npm install --save-dev @webgpu/types
```

Add to `tsconfig.json`:

```json
{
  "compilerOptions": {
    "types": ["@webgpu/types"]
  }
}
```

---

## GPU features (`tgpu.init` options)

Optional device capabilities are requested at init; `d.f16`/`vec*h` need `shader-f16`, subgroup ops need `subgroups`, GPU timing needs `timestamp-query`:

```ts
const root = await tgpu.init({
  adapter: { powerPreference: 'high-performance' },  // GPURequestAdapterOptions
  device: {
    requiredFeatures: ['shader-f16'],       // init throws if unavailable
    optionalFeatures: ['timestamp-query'],  // requested when available
  },
});

root.enabledFeatures.has('timestamp-query'); // ReadonlySet<GPUFeatureName>
```

**Every `requiredFeatures` entry shrinks the set of devices the app runs on** — init fails outright on hardware without it. Require a feature only when that trade-off is deliberate. Otherwise request it via `optionalFeatures` and write both paths, branching on a captured `root.enabledFeatures.has(...)` result — it's comptime-known, so branch pruning emits only the taken path in the generated WGSL:

```ts
const hasF16 = root.enabledFeatures.has('shader-f16');

const process = (x: number) => {
  'use gpu';
  if (hasF16) {
    return fastF16Path(x);   // only the taken branch survives in WGSL
  } else {
    return f32Path(x);       // keep the fallback in `else` - code after the
  }                          // `if` is not part of the pruned branch and
};                           // would be emitted (unreachably) either way
```

---

## 3. Build plugin - `unplugin-typegpu` (required for `'use gpu'`)

The `'use gpu'` directive and JS/TS shader functions need the build plugin. Without it, TypeGPU functions implemented in TypeScript won't work.

```sh
npm install --save-dev unplugin-typegpu
```

### Vite

```js title="vite.config.js"
import { defineConfig } from 'vite';
import typegpu from 'unplugin-typegpu/vite';

export default defineConfig({
  plugins: [typegpu()],
});
```

### Babel (React Native / Expo)

```js title="babel.config.js"
module.exports = (api) => {
  api.cache(true);
  return {
    presets: ['babel-preset-expo'],
    plugins: ['unplugin-typegpu/babel'],
  };
};
```

Other bundlers supported: esbuild, rollup, rolldown, rspack, webpack, farm. Vite and Babel are actively maintained by the TypeGPU team.

### Plugin options

```ts
typegpu({
  include?: FilterPattern,       // default: [/\.m?[jt]sx?$/]
  exclude?: FilterPattern,
  enforce?: 'pre' | 'post',
  autoNamingEnabled?: boolean,   // default: true - names resources from variable names
  earlyPruning?: boolean,        // default: true - skips files without typegpu/tgpu/'use gpu'
  forceTgpuAlias?: string,       // only if tgpu import is aliased unusually
})
```

The plugin also auto-names TypeGPU resources from variable names, improving debugging without manual `.$name()` calls.

---

## 4. Operator overloading - `tsover` (the default TypeGPU way)

`tsover` is a drop-in TypeScript replacement adding operator overloading (`+ - * / %` on vectors and matrices). It is technically optional — without it the code still compiles and runs, but the IDE treats `d.vec3f() * 2` as a type error. **Treat it as part of the standard setup**: operators are the idiomatic TypeGPU style, so install `tsover` unless the user explicitly declines or the project has concrete TypeScript language-server constraints that rule it out; only then fall back to infix methods (`.add()`, `.mul()`).

**Note:** `unplugin-typegpu` already handles runtime operator overloads inside `'use gpu'` functions — no bundler plugin needed for shader code. `tsover` adds IDE type-checking support and enables operators outside `'use gpu'` blocks (CPU-side vector math).

### Install

Replace `typescript` with `tsover` in `package.json`:

```json
{
  "devDependencies": {
    "typescript": "npm:tsover@latest"
  }
}
```

For monorepos, add overrides:

```json
// npm / pnpm
{ "overrides": { "typescript": "npm:tsover@latest" } }

// Yarn
{ "resolutions": { "typescript": "npm:tsover@latest" } }
```

Match major.minor: if your project uses `typescript@5.8.x`, use `tsover@5.8.x`.

### Enable in code

No `tsconfig.json` changes. Operators typecheck inside `'use gpu'` functions as-is once tsover is the project's TypeScript. For CPU-side code outside `'use gpu'`, add a `'use tsover'` directive at file or function scope:

```ts
'use tsover'; // file-level; or place inside a single function
const c = a + b; // d.v2f + d.v2f
```

### Bundler plugin (only for CPU-side operators)

Only needed if you use operators outside `'use gpu'` blocks. Most projects can skip this.

```js title="vite.config.js"
import tsoverPlugin from 'tsover/plugin/vite';
// add to plugins array alongside typegpuPlugin()
```

Also available: `tsover/plugin/rollup`, `tsover/plugin/rolldown`.

### IDE: use workspace TypeScript version

VS Code / Cursor / Windsurf: command palette -> **TypeScript: Select TypeScript Version** -> **Use Workspace Version**.

Or in `.vscode/settings.json`:

```json
{ "typescript.tsdk": "node_modules/typescript/lib" }
```

Zed: set `tsdk` in `.zed/settings.json` for `vtsls` or `typescript-language-server`.

---

## Troubleshooting

**`'use gpu'` silently does nothing / `ResolutionError` about an untranspiled function** - the build plugin isn't running on that file. This is the #1 setup failure. Check, in order:
1. The plugin is in the right config for your bundler (`unplugin-typegpu/vite` in `vite.config.js` vs `'unplugin-typegpu/babel'` in `babel.config.js`).
2. The file matches the plugin's `include` pattern (default `[/\.m?[jt]sx?$/]`).
3. The dev server was restarted after adding the plugin; on React Native, clear the Metro cache (`npx expo start --clear`).
4. If the `tgpu` import is aliased unusually, set `forceTgpuAlias`.

**IDE flags `vec * scalar` as a type error while the code runs fine** - `tsover` isn't active in the editor: select the workspace TypeScript version (see above); for CPU-side code also check the `'use tsover'` directive is present.

---

## Lint plugin - `eslint-plugin-typegpu`

Highlights common pitfalls and unsupported syntax in `'use gpu'` functions. Optional but highly recommended - use it unless told otherwise or the project's linter is incompatible, and include a `lint` script in `package.json`.

```sh
npm install --save-dev eslint-plugin-typegpu
# pnpm add -D eslint-plugin-typegpu / yarn add -D eslint-plugin-typegpu
```

### ESLint (`eslint.config.js`)

```ts
import { defineConfig } from 'eslint/config';
import typegpu from 'eslint-plugin-typegpu';

export default defineConfig([
  {
    ...typegpu.configs.recommended,
    files: ['**/*.{js,mjs,ts,jsx,tsx}'],
  },
]);
```

### Oxlint (`oxlint.config.ts`)

```ts
import { defineConfig } from 'oxlint';
import typegpu from 'eslint-plugin-typegpu';

export default defineConfig({
  jsPlugins: ['eslint-plugin-typegpu'],
  rules: {
    ...typegpu.configs.recommended.rules,
  },
  ignorePatterns: ['node_modules'],
});
```
