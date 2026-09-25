# Shader Performance

Judge a change by the compiler's statistics (registers, spills, instruction count) or GPU timing on the target device, not by how the source looks. Live timestamps are noisy (see `references/timing.md`).

## Leave to the compiler

Shader compilers already do SSA, inlining, code motion of pure loads and math, and dead-code and dead-field elimination. These choices usually compile to the same code, so don't restructure for them:
- declaration order of large values;
- struct copies, and field-by-field vs whole-struct writes;
- early-outs placed before *pure* loads;
- hoisting `uniform.$.x` out of loops.

Only non-movable work (derivative sampling, atomics, writes) is worth moving behind a branch.

## Do by hand

- **Hoist loop-invariant algebra.** Rewrite the maths so invariant work leaves the loop, e.g. transform a ray once instead of every marched point.
- **Limit simultaneously live state.** Many independent accumulators, large local arrays and very long unrolled bodies are what spill.
- **Avoid integer `%` and division by runtime values in hot loops.** For unsigned values and a power-of-two `n`, use `x & (n - 1)` and `x >> log2(n)`.

## Branches and divergence

- A runtime ternary or `std.select` evaluates both sides. With a uniform condition (a mode, a uniform flag), `if`/`else` skips the dead side, which matters when the branches are expensive.
- A divergent `break` or branch costs as much as the slowest lane in its SIMD group, so early exits pay off only when neighbouring threads exit together.

## Loops and unrolling

- Compilers auto-unroll short loops on their own, so `tgpu.unroll` changes nothing there.
- Unrolling longer dependent chains (FBM octaves, fixed-step marches) often helps well beyond the counts compilers unroll themselves.
- Loops that index a small local array gain the most, because constant indices let the array live in registers.
- The limit is code size (instruction cache), not just registers. For long loops, unroll an inner chunk of a few iterations inside a rolled loop:

```ts
for (const i of std.range(N / 8)) {
  for (const j of tgpu.unroll(std.range(8))) {
    acc = step(acc, d.f32(i * 8 + j));
  }
}
```

- Don't replace a data-dependent `break` with unrolled, guarded iterations.

## Memory layout

- **SoA vs AoS by access pattern.** Use SoA when kernels read a subset of fields and AoS when they read whole records.
- **Texture or storage buffer:**
  - Random gathers and lookup tables: a packed storage buffer (e.g. rgba8 in a `u32`, decoded with `std.unpack4x8unorm`) is often faster than a texture, even with manual interpolation.
  - Coherent filtered reads (blurs, warps, primary-ray volume marching): use a texture. Hardware filtering is nearly free there, while manual filtering multiplies the loads.
- **Format size.** Texture cost grows with bytes per texel, so use the smallest format that holds the data. Wide float formats (rgba32float) are the worst case for `textureLoad`; a plain `array<vec4f>` is usually cheaper.

## Workgroups and precision

- Size compute workgroups to at least the SIMD width; smaller groups leave lanes idle. A multiple of 64 is safe across vendors.
- Don't assume `f16` doubles throughput; the speedup is GPU-dependent. Its dependable gains are fewer registers and less bandwidth.
