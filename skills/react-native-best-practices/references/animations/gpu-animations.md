# GPU Shader Animations

GPU-accelerated animations using WebGPU shaders in React Native. Use when animation requirements exceed what Reanimated can handle: massive particle counts, fluid/physics simulations, procedural noise, custom shader effects, SDF rendering, or 3D scenes.

Requires React Native 0.81+ and the New Architecture.

**Version requirements:** `react-native-wgpu` requires `react-native-reanimated >= 4.2.1` and `react-native-worklets >= 0.7.2`. Upgrade both before installing if the project uses older versions:

```sh
npm install react-native-reanimated@latest react-native-worklets@latest
npm install react-native-wgpu
```

---

## When to Use GPU Shaders

GPU shaders are the right choice when:
- The effect requires per-pixel computation (noise fields, distortion, blur, glow)
- Physics simulations drive the animation (fluid, cloth, boids flocking, gravity)
- You need signed distance function (SDF) rendering for procedural shapes
- 3D rendering is required (meshes, lighting, skeletal animation via Three.js)
- Animation state is computed from math that benefits from GPU parallelism

For high element counts without per-pixel computation (hundreds of 2D shapes, sprites, tiles), prefer `react-native-skia` with the Atlas API instead (see `canvas-animations.md`). Skia renders to a single canvas with lower setup cost than WebGPU.

Stick with Reanimated when animating standard UI components (opacity, transforms, layout changes) or responding to gestures. GPU shaders render into a `Canvas` element, separate from the React Native view hierarchy.

---

## Setup

### react-native-wgpu

Provides the WebGPU API in React Native using Dawn as the backend.

```sh
npm install react-native-wgpu
```

With Expo:
```sh
npx create-expo-app@latest -e with-webgpu
```

### TypeGPU (recommended)

Type-safe abstraction over WebGPU. Write shader logic in TypeScript that compiles to WGSL.

```sh
npm install typegpu
npm install --save-dev @webgpu/types unplugin-typegpu
```

Add `@webgpu/types` to `tsconfig.json`:

```json
{ "compilerOptions": { "types": ["@webgpu/types"] } }
```

Add the Babel plugin for `'use gpu'` function syntax:

```js
// babel.config.js
module.exports = (api) => {
  api.cache(true);
  return {
    presets: ['babel-preset-expo'],
    plugins: ['unplugin-typegpu/babel'],
  };
};
```

Clear Metro cache after adding the plugin: `npx expo --clear`.

Expo Go is not supported. Run `npx expo prebuild` if migrating from Expo Go.

On the iOS simulator, disable Metal Validation in Edit Scheme to avoid crashes.

---

## Canvas and Device Setup

The `Canvas` component provides the WebGPU surface. Use `useDevice` and `useCanvasRef` hooks from `react-native-wgpu`:

```tsx
import { Canvas, useDevice, useCanvasRef } from 'react-native-wgpu';
import { tgpu, d } from 'typegpu';

function GPUScene() {
  const { device = null } = useDevice();
  const ref = useCanvasRef();

  useEffect(() => {
    if (!device) return;

    const context = ref.current!.getContext('webgpu')!;
    const root = tgpu.initFromDevice({ device });
    const format = navigator.gpu.getPreferredCanvasFormat();
    context.configure({ device, format, alphaMode: 'premultiplied' });

    // Create pipelines and render here
    // ...
    context.present(); // Required: manually present the frame

    return () => root.destroy();
  }, [device]);

  return <Canvas ref={ref} style={{ flex: 1 }} />;
}
```

### Key differences from Web WebGPU

- **`context.present()`** is required after submitting commands. React Native does not auto-present frames.
- Canvas size uses device pixel ratio: `canvas.width = canvas.clientWidth * PixelRatio.get()`.
- On Android, `alphaMode` is ignored in `context.configure()`. Use the `transparent` prop on `Canvas` instead.

---

## Reanimated Integration

WebGPU objects (`GPUDevice`, `GPUCanvasContext`) are automatically registered for worklet serialization. You can pass them directly to worklets and run GPU rendering on the UI thread:

```tsx
import { scheduleOnUI } from 'react-native-worklets';

const renderFrame = (device: GPUDevice, context: GPUCanvasContext) => {
  'worklet';
  const commandEncoder = device.createCommandEncoder();
  const textureView = context.getCurrentTexture().createView();

  const pass = commandEncoder.beginRenderPass({
    colorAttachments: [{
      view: textureView,
      clearValue: [0, 0, 0, 1],
      loadOp: 'clear',
      storeOp: 'store',
    }],
  });
  pass.setPipeline(pipeline);
  pass.draw(3);
  pass.end();

  device.queue.submit([commandEncoder.finish()]);
  context.present();
};

// Trigger from RN thread
scheduleOnUI(renderFrame, device, context);
```

This requires `react-native-reanimated` and `react-native-worklets` as peer dependencies:

```sh
npm install react-native-reanimated react-native-worklets
```

Use this pattern to tie GPU rendering to Reanimated shared values or `useFrameCallback` for continuous animation loops.

Dawn devices are not thread-safe by default. If the JS thread also touches the device (writes buffers, reads results, submits) while a worklet renders, request the `'implicit-device-synchronization'` feature when creating the device; without it iOS crashes in Metal asserts. The `typegpu` skill shows how with TypeGPU.

---

## Writing TypeGPU Code

For shaders, buffers, pipelines and the `@typegpu/*` packages (`@typegpu/noise` for random numbers and Perlin noise, `@typegpu/sdf` for signed distance fields, `@typegpu/react` for hooks and worklet frame loops), use the `typegpu` skill. It matches the installed TypeGPU version.

## Common Animation Patterns

### Particle system

Keep particle state in a storage buffer of structs, update it in a compute pass each frame, and draw the particles as instanced geometry that reads the same buffer.

### Boids / flocking simulation

A classic GPU compute pattern: each boid reads neighbors from a storage buffer in the compute pass, updates its velocity, then renders as instanced triangles.

### Fluid simulation

Grid-based velocity field stored in a storage buffer or texture. Compute passes handle advection, pressure solving, and diffusion. Render the velocity field as colors or displace geometry.

---

## Three.js and React Three Fiber

For 3D scenes, react-native-wgpu supports Three.js (from r168+) and React Three Fiber:

1. Modify Metro config to resolve Three.js to the WebGPU build.
2. Patch `@react-three/fiber/package.json` to use the WebGPU entry point instead of the React Native bundle:

```diff
-  "react-native": "native/dist/react-three-fiber-native.cjs.js",
+  "react-native": "dist/react-three-fiber.cjs.js",
```

This enables skeletal animation, cloth simulation, particle effects, and complex 3D scene graphs using the Three.js ecosystem.

---

## Performance Notes

- GPU compute runs thousands of threads in parallel with zero JS overhead.
- Frame presentation is manual (`context.present()`). You control timing.
- TypeGPU and raw WebGPU can be mixed freely. Adoption is non-contagious: use TypeGPU where type safety helps and raw WebGPU everywhere else.
