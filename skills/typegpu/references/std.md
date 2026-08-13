# TypeGPU `std` Library Reference

`std` wraps WGSL built-in functions ([WGSL spec section 16](https://www.w3.org/TR/WGSL/#builtin-functions)) plus some TypeGPU additions. The WGSL documentation applies.

```ts
import { std } from 'typegpu';
```

The listing below is the commonly used subset — `std` exposes the WGSL builtin set under the same names (e.g. `fma`, `saturate`, `modf`, `frexp`, `ldexp`, `quantizeToF16`, `countOneBits`, `reverseBits`, `extractBits`, `insertBits`, `firstLeadingBit`, `firstTrailingBit` all exist).

**Math**
```ts
std.abs   std.sign   std.floor   std.ceil   std.round   std.fract   std.trunc
std.sqrt  std.inverseSqrt
std.exp   std.exp2   std.log   std.log2   std.pow
std.min   std.max   std.clamp(x, lo, hi)
std.mod(a, b)
std.mix(a, b, t)              // linear interpolation
std.smoothstep(edge0, edge1, x)
std.step(edge, x)
std.select(falseVal, trueVal, cond)
std.copy(x)                   // schema-agnostic deep copy (works when exact schema is generic)
std.bitcast(d.u32, d.f32)(x)  // reinterpret bits; any scalar/vector pair of equal size
```

**Comparison / boolean** — componentwise on vectors, returning bool vectors (`d.vec2b`/`d.vec3b`/`d.vec4b`); combine with `std.select` for vectorized branching
```ts
std.eq  std.ne  std.lt  std.le  std.gt  std.ge   // componentwise -> vecNb
std.all(bv)  std.any(bv)  std.allEq(a, b)
std.and(bv, bv)  std.or(bv, bv)  std.not(bv)
```

**Trig**
```ts
std.sin  std.cos  std.tan  std.asin  std.acos  std.atan  std.atan2
std.sinh std.cosh std.tanh  std.degrees  std.radians
```

**Vector / matrix**
```ts
std.dot(a, b)         std.cross(a, b)      std.length(v)
std.normalize(v)      std.distance(a, b)
std.reflect(i, n)     std.refract(i, n, eta)
std.faceForward(n, i, nRef)
std.mul(mat, vec)     // matrix-vector multiply
std.transpose(m)      std.determinant(m)
std.identity2/3/4()   std.translation4(v)  std.scaling4(v)   // build matrices in-shader
std.rotationX4(rad)   std.rotationY4(rad)  std.rotationZ4(rad)
```

**Arrays**
```ts
std.arrayLength(arr)  // runtime-sized storage array -> WGSL arrayLength;
                      // fixed-size array -> folds to the constant
```

**Texture** — stage/uniformity rules and the sampling reference table live in `references/textures.md`. Short version: `textureSample` is fragment-only and must be in uniform control flow; `textureSampleLevel` works in any stage; `textureSampleGrad` gives auto-LOD in compute/non-uniform branches.
```ts
std.textureSample(view.$, sampler.$, uv)
std.textureSampleLevel(view.$, sampler.$, uv, mipLevel)
std.textureSampleGrad(view.$, sampler.$, uv, ddx, ddy)
std.textureSampleBias(view.$, sampler.$, uv, bias)
std.textureSampleBaseClampToEdge(view.$, sampler.$, uv)
std.textureSampleCompare(depthView.$, comparisonSampler.$, uv, ref)  // + CompareLevel
std.textureGather(component, view.$, sampler.$, uv)
std.textureLoad(view.$, coords, mipLevel)
std.textureStore(storageView.$, coords, value)
std.textureDimensions(view.$)
```

**Fragment control**
```ts
std.discard()   // fragment-only; discards the fragment (WGSL `discard`)
```

**Derivatives** (fragment only)
```ts
std.dpdx(v)  std.dpdy(v)  std.fwidth(v)   // + Coarse/Fine variants of each
```

**Synchronization** (compute only)
```ts
std.workgroupBarrier()   std.storageBarrier()
```

**Atomic**
```ts
std.atomicLoad(ptr)         std.atomicStore(ptr, val)
std.atomicAdd(ptr, val)     std.atomicSub(ptr, val)
std.atomicMin(ptr, val)     std.atomicMax(ptr, val)
std.atomicAnd(ptr, val)     std.atomicOr(ptr, val)     std.atomicXor(ptr, val)
// atomicExchange / atomicCompareExchangeWeak are NOT exposed
```

**Packing** (the full WGSL set is not exposed — only these four)
```ts
std.pack4x8unorm(v)   std.unpack4x8unorm(x)
std.pack2x16float(v)  std.unpack2x16float(x)
```

**Subgroups** (require the `subgroups` device feature — see `references/setup.md`)
```ts
std.subgroupAdd / Mul / Min / Max / And / Or / Xor (value)
std.subgroupExclusiveAdd / ExclusiveMul / InclusiveAdd / InclusiveMul (value)
std.subgroupAll / Any (bool)    std.subgroupBallot(bool)   std.subgroupElect()
std.subgroupBroadcast(value, lane)   std.subgroupBroadcastFirst(value)
std.subgroupShuffle / ShuffleUp / ShuffleDown / ShuffleXor (value, x)
// builtins: d.builtin.subgroupId, subgroupSize, subgroupInvocationId, numSubgroups
```

**Environment probes** - branch on where the code is running (e.g. CPU fallback path vs generated shader): `std.isBeingTranspiled()`, `std.getTargetShaderLanguage()` (`'wgsl'` during generation, `undefined` otherwise), `std.getShaderStage()` (`'vertex' | 'fragment' | 'compute' | undefined`).
