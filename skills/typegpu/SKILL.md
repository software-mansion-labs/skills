---
name: typegpu
description: >-
  Use when writing, debugging, reviewing or optimizing code that imports typegpu or a @typegpu/* package, or uses tgpu.*, d.* schemas, std.* or 'use gpu' functions: shaders, buffers, textures, bind groups, compute and render pipelines, vertex layouts, command encoders and render bundles, slots and accessors, React and React Native rendering (@typegpu/react: useRoot, useFrame, useUniform, useConfigureContext), random numbers and noise (@typegpu/noise), signed distance fields (@typegpu/sdf), colour and tonemapping (@typegpu/color), 2D global illumination (@typegpu/radiance-cascades). Covers the CPU side and the shader side together, so use it even when the question names only one. TypeGPU is type-safe WebGPU in TypeScript. Not for raw WebGPU written without tgpu, standalone WGSL, or WebGL.
---

# TypeGPU

The full TypeGPU skill ships inside the `typegpu` package, so it always matches the installed version. This skill only routes to it.

1. Read `node_modules/typegpu/skills/typegpu/SKILL.md` whole and follow it. In a workspace, `typegpu` may be hoisted to the `node_modules` of a parent folder.
2. If `typegpu` is installed but that file is missing, the installed release predates the bundled skill. Suggest upgrading `typegpu`. Until then, the `.d.ts` files in `node_modules/typegpu` are the exact API, and the docs cover the rest: https://docs.swmansion.com/TypeGPU/llms.txt lists the guides, and every page has a Markdown copy (replace the trailing `/` of its URL with `.md`).
3. If the project has no `typegpu` yet, add it with `npx typegpu@latest --enhance` (or create a project with `npx typegpu@latest`), then go to step 1.
