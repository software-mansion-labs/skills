# Command Encoders, Passes, and Render Bundles

> **Unstable API.** These live on `root['~unstable']` — the surface may change between minor releases. If a call documented here errors or doesn't typecheck, verify the signature against the installed `typegpu` version (its `.d.ts` or docs) before debugging elsewhere.

By default, `pipeline.draw()` and `pipeline.dispatchWorkgroups()` each record their own single-pipeline pass and submit it immediately. Encoders are for the cases that need more control: **several pipelines in one pass** (shared attachments, e.g. scene + lights + sky into one MSAA target) and **several passes in one submission**.

## Typed command encoder and render passes

```ts
const encoder = root['~unstable'].createCommandEncoder();

const pass = encoder.beginRenderPass({
  colorAttachments: [{
    view: msaaTexture,        // TypeGPU texture/view, canvas context, or GPUTextureView
    resolveTarget: context,
  }],
  depthStencilAttachment: {
    view: depthTexture,
  },
});

scenePipeline.with(pass).draw(mesh.vertexCount);
lightPipeline.with(pass).draw(6, lightCount);
skyPipeline.with(pass).draw(3);

pass.end();
encoder.submit();
```

Descriptor conveniences over raw `GPURenderPassDescriptor`:

- Attachment `view` accepts TypeGPU textures, texture views, and canvas contexts, as well as raw `GPUTextureView`s.
- `loadOp` / `storeOp` / `depthClearValue` default to `'clear'` / `'store'` / `1`.
- A single color attachment doesn't need to be wrapped in an array.
- `occlusionQuerySet` / `timestampWrites` accept TypeGPU query sets.

## Two equivalent execution styles

`pipeline.with(pass).draw(...)` keeps the pipeline-centric API (all `with*` methods available). Alternatively the pass mirrors `GPURenderPassEncoder`, accepting TypeGPU resources:

```ts
pass.setPipeline(renderPipeline);
pass.setBindGroup(bindGroup);
pass.setVertexBuffer(vertexLayout, vertexBuffer);
pass.draw(3);
```

Both styles share one pass state, applied lazily at draw time and only when it changed since the previous draw. Precedence (TypeGPU 0.12.6+): state set on the pass (`pass.set*`) wins over state held by the pipeline (`pipeline.with*`), which wins over WebGPU defaults. Pass state persists across draws and pipeline switches until overwritten; pipeline-held state applies only to that pipeline's draws. Footgun: `pipeline.with(pass).draw(...)` sets the pass's current pipeline — a subsequent bare `pass.draw(...)` runs *that* pipeline, not one set earlier via `setPipeline`.

`pipeline.with(pass)` allocates a wrapper on every call. When drawing one pipeline many times, bind it once and vary per-draw values through the pass:

```ts
const bound = scenePipeline.with(pass);
for (const object of objects) {
  pass.setBindGroup(object.bindGroup);
  bound.draw(object.vertexCount);
}
```

## Compute passes

```ts
const pass = encoder.beginComputePass();
computePipeline.with(pass).dispatchWorkgroups(16);
pass.end();
encoder.submit();
```

Caveat: guarded compute pipelines (`createGuardedComputePipeline` / `dispatchThreads`) cannot record into passes or encoders — each `dispatchThreads` submits on its own.

## Immediates (per-draw values without buffers)

> **Unstable and not portable.** `tgpu['~unstable'].immediateVar` (TypeGPU 0.12.6+) needs the `immediate_address_space` WGSL language feature. As of September 2026, only Chromium 149+ exposes it, on Windows, macOS, ChromeOS and Android with capable hardware; Firefox, Safari, Chrome on Linux and react-native-wgpu do not. Always keep a fallback path.

Immediates are WebGPU's push constants: a small value set per draw or dispatch, directly on the pass. An accessor lets the same shader read either an immediate or a buffer:

```ts
const params = tgpu.bindGroupLayout({ tint: { uniform: d.vec4f } });
const tint = tgpu.accessor(d.vec4f, () => params.$.tint); // shader reads tint.$

const tintImmediate = root.enabledWgslLanguageFeatures.has('immediate_address_space')
  ? tgpu['~unstable'].immediateVar(d.vec4f)
  : undefined;
const pipeline = (tintImmediate ? root.with(tint, tintImmediate) : root)
  .createRenderPipeline({ vertex, fragment, targets: { format } });

const bound = pipeline.with(pass);
for (const object of objects) {
  if (tintImmediate) {
    pass.setImmediates(tintImmediate, object.tint); // or a Float32Array, copied verbatim
  } else {
    pass.setBindGroup(object.tintBindGroup); // one uniform buffer + bind group per value
  }
  bound.draw(object.vertexCount);
}
```

- **Where values come from:** `pipeline.with(immediate, value)` holds a value on the pipeline, and `pass.setImmediates` overrides it. The variable's optional default (`immediateVar(schema, default)`) is the last resort. Drawing with no value throws `MissingImmediatesError`.
- **Values are serialized when set.** Mutating the source object afterwards has no effect.
- **Limits:** one immediate variable per shader, so group several values in a struct. No arrays, atomics or booleans. The size must be a multiple of 4 bytes and fit `maxImmediateSize` (64 bytes by default; request more via `tgpu.init({ device: { requiredLimits: { maxImmediateSize: 128 } } })` where supported).
- **Fallback:** a single uniform written with `buffer.write` between draws does not work: queue writes are not recorded between draws, so every draw would see the last value. Use one buffer and bind group per distinct value, as above. When one variable holds either an immediate or a uniform, `isImmediateVar(x)` (exported from `typegpu`) tells them apart.

Test the fallback explicitly: a recent desktop Chrome usually takes the immediate path.

## `submit()` vs `finish()`

- `encoder.submit()` finishes and submits to the device queue; shader `console.log` output and performance callbacks are processed as part of that submission.
- `encoder.finish()` returns the raw `GPUCommandBuffer` for manual `device.queue.submit([...])` batching — TypeGPU never sees that submission, so **logs and performance callbacks are not processed**.

## Render bundles

Pre-record a static draw sequence once, replay it cheaply every frame:

```ts
const bundleEncoder = root['~unstable'].createRenderBundleEncoder({
  colorFormats: ['rgba8unorm'],  // must match the pass it will run in
});
scenePipeline.with(bundleEncoder).draw(vertexCount);
const bundle = bundleEncoder.finish();

// each frame:
pass.executeBundles([bundle]);
```

## Raw WebGPU interop

Pipelines also accept raw WebGPU encoders via `.with(...)`:

- `pipeline.with(gpuCommandEncoder)` — TypeGPU opens and ends the needed pass; you finish/submit the encoder.
- `pipeline.with(gpuRenderPass | gpuComputePass)` — state is applied to the existing pass without ending it; you end it.
- `pipeline.with(gpuRenderBundleEncoder)` — records draws into the bundle; you finish it.

Same logging caveat as `finish()`: TypeGPU can't process shader logs or perf callbacks for submissions it doesn't own.

Escape hatch in the other direction: `root.unwrap(encoder)` / `root.unwrap(pass)` return the raw `GPUCommandEncoder` / pass encoder (e.g. for texture copies). Commands recorded raw are invisible to TypeGPU, so after unwrapping a pass, the next typed draw re-applies its full state.

## Encoder-aware buffer ops

`buffer.clear(encoder)` and `dst.copyFrom(src, encoder)` record into the given encoder instead of submitting immediately — use them to fold buffer maintenance into the same submission as your passes.
