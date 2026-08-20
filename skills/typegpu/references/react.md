# `@typegpu/react` - React and React Native Bindings

Hooks for creating and managing TypeGPU resources in React components. The same import works on web and React Native — the package's `react-native` export condition selects the RN build automatically.

```sh
npm install @typegpu/react
```

## Minimal complete example

```tsx
import { useMemo } from 'react';
import { d, common } from 'typegpu';
import { useConfigureContext, useFrame, useRoot, useUniform } from '@typegpu/react';

function MyEffect() {
  const root = useRoot();
  const time = useUniform(d.f32);

  const renderPipeline = useMemo(
    () =>
      root.createRenderPipeline({
        vertex: common.fullScreenTriangle,
        fragment: ({ uv }) => {
          'use gpu';
          return d.vec4f((uv * 5 + time.$) % 1, 0, 1);
        },
      }),
    [root, time],
  );

  const { ref, ctxRef } = useConfigureContext();

  useFrame(({ elapsedSeconds }) => {
    if (!ctxRef.current) return;
    time.write(elapsedSeconds);
    renderPipeline.withColorAttachment({ view: ctxRef.current }).draw(3);
  });

  return <canvas ref={ref} />;
}
```

## Hook reference

| Hook | Purpose |
|---|---|
| `useRoot()` | Root from the nearest `<Root>` provider (global root if none). Suspends until initialized; throws on init failure. `useRootOrError` / `useRootWithStatus` for manual handling. |
| `<Root root?>` | Optional context provider; all descendants share one `TgpuRoot`. Creates one on demand if `root` not given. |
| `useConfigureContext(opts?)` | Returns `{ ref, ctxRef }` — pass `ref` to the `<canvas>`, use `ctxRef.current` as the attachment view. |
| `useFrame(cb)` | Runs `cb` every frame (rAF) with `{ deltaSeconds, elapsedSeconds }`. No `useCallback` needed — latest closure values are visible. |
| `useUniform(schema, opts?)` | Uniform buffer binding tied to component lifetime; update via `.write()` outside the React lifecycle. Options: `initial`, `onInit`. |
| `useMutable(schema, opts?)` | Same, `var<storage, read_write>`. |
| `useReadonly(schema, opts?)` | Same, `var<storage, read>`. |
| `useMirroredUniform(schema, value)` | Uniform re-synced to `value` on every React re-render — for values living in the React lifecycle (theme, props). |
| `useBuffer(schema, opts?)` | Raw buffer equivalent of the above. |
| `useBindGroup(layout, entries)` | Bind group; API matches `root.createBindGroup`. |
| `ClientOnly` | SSR guard — render GPU components only on the client. |

Guidance:

- **Memoize pipelines** with `useMemo` and `[root, ...capturedResources]` deps — pipeline creation implies shader resolution.
- **Per-frame values go through `.write()` inside `useFrame`**, never through React state (state updates re-render and can re-create resources).
- `useUniform` vs `useMirroredUniform`: updated-every-frame → `useUniform` + `.write()`; derived from React data → `useMirroredUniform`.

## React Native worklets (UI-thread render loops)

With `react-native-worklets` installed, `useFrame` callbacks run on the UI thread, unaffected by RN-thread load. Setup (after the usual RN/webgpu setup) — babel config:

```js
const workletsPluginOptions = {
  bundleMode: true,
  importForwarding: {
    moduleNames: ['typegpu'],
    relativePaths: ['my-app/components'], // dirs with module-scope shader definitions
  },
};

// plugins: ['unplugin-typegpu/babel', ['react-native-worklets/plugin', workletsPluginOptions]]
```

Clear the Metro cache after changing babel config (`npx expo start --clear`).

No extra imports — `@typegpu/react` detects `react-native-worklets` at runtime. While active:

- **Every `useFrame` callback must start with the `'worklet'` directive** and runs on the UI thread; a plain callback throws. Opt out with `<Root disableWorklets>`.
- Call `ctx.present?.()` after drawing (RN canvas contexts need an explicit present).

```tsx
useFrame(({ elapsedSeconds }) => {
  'worklet';
  const ctx = ctxRef.current;
  if (!ctx) return;
  color.write(d.vec3f(0.5 + Math.sin(elapsedSeconds) * 0.5, 0.447, 0.941));
  pipeline.withColorAttachment({ view: ctx }).draw(3);
  ctx.present?.();
});
```

### Rules of transfer

Resources captured by a worklet are transferred to the UI runtime automatically on first use; both runtimes share the same GPU objects.

- **Transfers**: buffers (incl. `createUniform`/`createMutable`/`createReadonly`), textures, samplers, bind groups + layouts, vertex layouts, query sets, pipelines, roots, slots, accessors, consts.
- **Does NOT transfer**: shader definitions (`tgpu.fn`, entry functions, `tgpu.comptime`) and standalone schemas/vector instances. Keep definitions at module scope in files covered by `importForwarding` — worklets re-import them natively.
- **Pipelines carrying attachments/passes throw on transfer** — transfer the bare pipeline and call `withColorAttachment` on the worklet side (as above).
- **Callbacks crossing runtimes** (e.g. `withPerformanceCallback`) must themselves be `'worklet'`-marked.
- **Pipelines created on the UI thread need an explicit `targets: { format: '...' }`** — `navigator.gpu.getPreferredCanvasFormat()` is unavailable on worklet runtimes.
- Query sets: `resolve()` and `read()` on one runtime only; shader logs print only on the runtime that resolved the pipeline.
