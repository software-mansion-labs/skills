# Setup, Model Loading, and Error Handling

Installation, resource fetcher initialization, model loading strategies, download management, error handling, and custom model integration for React Native ExecuTorch.

For the full getting started guide, webfetch [Getting Started](https://docs.swmansion.com/react-native-executorch/docs/fundamentals/getting-started). For version compatibility, webfetch the [Compatibility table](https://docs.swmansion.com/react-native-executorch/docs/other/compatibility).

---

## Installation

### Expo

```bash
npx expo install react-native-executorch
npx expo install react-native-executorch-expo-resource-fetcher expo-file-system expo-asset
```

Expo Go is not supported. Use a custom development build (`npx expo prebuild`).

### Bare React Native

```bash
npm install react-native-executorch
npm install react-native-executorch-bare-resource-fetcher @dr.pogodin/react-native-fs @kesha-antonov/react-native-background-downloader
cd ios && pod install
```

### Prerequisites

- **New Architecture (Fabric) required.** The old architecture is not supported.
- **iOS release builds require a real device.** Simulator release builds are not supported because ExecuTorch uses Metal APIs unavailable in the simulator.

---

## Resource Fetcher Initialization

Call `initExecutorch` once at app startup, before loading any model:

```tsx
// Expo
import { initExecutorch } from 'react-native-executorch';
import { ExpoResourceFetcher } from 'react-native-executorch-expo-resource-fetcher';

initExecutorch({ resourceFetcher: ExpoResourceFetcher });
```

```tsx
// Bare React Native
import { initExecutorch } from 'react-native-executorch';
import { BareResourceFetcher } from '@react-native-executorch/bare-adapter';

initExecutorch({ resourceFetcher: BareResourceFetcher });
```

Skipping this step causes all model loads to throw `ResourceFetcherAdapterNotInitialized`.

---

## Metro Configuration

To load models bundled as app assets via `require()`, register the `.pte` and `.bin` extensions:

```js
// metro.config.js
const { getDefaultConfig } = require('expo/metro-config');

const config = getDefaultConfig(__dirname);
config.resolver.assetExts.push('pte');
config.resolver.assetExts.push('bin');

module.exports = config;
```

---

## Model Loading Strategies

Models can be loaded from three sources. Choose based on model size and UX requirements:

```
How large is the model?
├── Small (< 512MB) and always needed?
│   └── Bundle with app assets: require('../assets/model.pte')
│       + Available instantly, no internet needed
│       - Increases app download size
│
├── Large (> 512MB) or optional feature?
│   └── Remote URL: 'https://huggingface.co/.../model.pte'
│       + Keeps app small, download on first use
│       - Requires internet on first launch
│       - Show download progress UI via downloadProgress
│
└── User-provided model?
    └── Local file: 'file:///path/to/model.pte'
        + Maximum flexibility
        - Requires file management UI
```

All three sources work with every hook:

```tsx
// Bundled asset
const llm = useLLM({ model: { modelSource: require('../assets/model.pte'), ... } });

// Remote URL (downloaded and cached automatically)
const llm = useLLM({ model: LLAMA3_2_1B }); // constants use remote URLs

// Local file
const llm = useLLM({ model: { modelSource: 'file:///var/mobile/.../model.pte', ... } });
```

Use predefined model constants (e.g., `LLAMA3_2_1B`, `EFFICIENTNET_V2_S`) when available. They point to optimized, pre-exported models from Software Mansion's [HuggingFace repository](https://huggingface.co/software-mansion).

### preventLoad

All hooks accept `preventLoad: true` to defer model loading until you're ready:

```tsx
const llm = useLLM({ model: LLAMA3_2_1B, preventLoad: true });

// Load later when needed
// (hook will load automatically once preventLoad changes to false)
```

### Download progress

Track download progress via the `downloadProgress` property (0 to 1):

```tsx
const llm = useLLM({ model: LLAMA3_2_1B });

<Text>Loading: {Math.round(llm.downloadProgress * 100)}%</Text>
```

---

## ResourceFetcher

For advanced download management (pause, resume, cancel, cleanup), use the `ResourceFetcher` utility. For the full API, webfetch [ResourceFetcher](https://docs.swmansion.com/react-native-executorch/docs/utilities/resource-fetcher).

### Download with progress

```tsx
import { ResourceFetcher } from 'react-native-executorch';

const uris = await ResourceFetcher.fetch(
  (progress) => console.log(`${Math.round(progress * 100)}%`),
  'https://huggingface.co/.../model.pte',
  'https://huggingface.co/.../tokenizer.bin'
);
// uris: string[] of local file paths (without file:// prefix)
```

### Pause, resume, cancel

```tsx
// Pause an active download
await ResourceFetcher.pauseFetching('https://...model.pte');

// Resume a paused download (faster than re-calling fetch)
const uris = await ResourceFetcher.resumeFetching('https://...model.pte');

// Cancel entirely
await ResourceFetcher.cancelFetching('https://...model.pte');
```

### Storage management

```tsx
// List all downloaded models
const models = await ResourceFetcher.listDownloadedModels();

// Get total size of remote files before downloading
const bytes = await ResourceFetcher.getFilesTotalSize('https://...model.pte');

// Delete downloaded resources
await ResourceFetcher.deleteResources('https://...model.pte');
```

Downloaded files are stored in the app's documents directory under `react-native-executorch/`.

---

## Error Handling

All errors inherit from `RnExecutorchError` with a `code` property from `RnExecutorchErrorCode`. For the full error reference, webfetch [Error Handling](https://docs.swmansion.com/react-native-executorch/docs/utilities/error-handling).

### Common errors and recovery

| Error Code | When | Recovery |
|---|---|---|
| `ModuleNotLoaded` | Calling `forward`/`generate` before model is ready | Check `isReady` before calling inference methods |
| `ModelGenerating` | Calling inference while another is running | Wait for completion or call `interrupt()` |
| `InvalidConfig` | Invalid config values (e.g., `topp` > 1) | Validate config parameters |
| `ResourceFetcherAdapterNotInitialized` | Forgot to call `initExecutorch()` | Call `initExecutorch({ resourceFetcher })` at app startup |
| `ResourceFetcherDownloadFailed` | Network error during model download | Retry with exponential backoff |
| `MemoryAllocationFailed` | Model too large for device | Use a smaller or more aggressively quantized model variant |
| `DownloadInterrupted` | Download did not complete | Retry the download |

### Error handling pattern

```tsx
import { RnExecutorchError, RnExecutorchErrorCode } from 'react-native-executorch';

try {
  const result = await model.forward(imageUri);
} catch (err) {
  if (err instanceof RnExecutorchError) {
    switch (err.code) {
      case RnExecutorchErrorCode.ModuleNotLoaded:
        // Model still loading, show loading state
        break;
      case RnExecutorchErrorCode.ModelGenerating:
        // Already processing, wait or interrupt
        break;
      case RnExecutorchErrorCode.MemoryAllocationFailed:
        // Device cannot fit this model, suggest smaller variant
        break;
      default:
        console.error('ExecuTorch error:', err.code, err.message);
    }
  } else {
    throw err;
  }
}
```

---

## Custom Models (useExecutorchModule)

For models not covered by the built-in hooks, use `useExecutorchModule` to run arbitrary `.pte` models. For the full API, webfetch [useExecutorchModule](https://docs.swmansion.com/react-native-executorch/docs/hooks/executorch-bindings/useExecutorchModule).

### Exporting a custom model

1. Export your PyTorch model to `.pte` format using the [ExecuTorch export tutorial](https://pytorch.org/executorch/stable/tutorials/export-to-executorch-tutorial.html)
2. Decide on a backend: XNNPACK (CPU, cross-platform default) or Core ML (iOS, uses ANE)
3. Load the `.pte` file via asset, URL, or local path

### Running a custom model

```tsx
import { useExecutorchModule } from 'react-native-executorch';

const model = useExecutorchModule({
  modelSource: require('../assets/custom_model.pte'),
});

const runInference = async () => {
  // Create input tensor matching your model's expected shape
  const input = {
    data: new Float32Array([1.0, 2.0, 3.0]),
    shape: [1, 3],
    dtype: 'float32',
  };

  const output = await model.forward([input]);
  // output is a TensorPtr[] matching your model's output shape
};
```

Input and output use the `TensorPtr` representation: `{ data: TypedArray, shape: number[], dtype: string }`.

### TypeScript Module API

For non-hook usage (e.g., in services or outside React components), use `ExecutorchModule` directly. For the full API, webfetch [ExecutorchModule](https://docs.swmansion.com/react-native-executorch/docs/typescript-api/executorch-bindings/ExecutorchModule).

---

## Device Constraints

### Memory

| Device tier | Parameter range | Examples |
|---|---|---|
| Low-end | 135M-500M | SmolLM 2 135M/360M |
| Mid-range | 500M-1.7B | Qwen 3 0.6B, SmolLM 2 1.7B, LLaMA 3.2 1B |
| High-end | 1.7B-4B | Qwen 3 4B, Phi 4 Mini, LLaMA 3.2 3B |

For detailed memory usage and inference time per model per device, webfetch [Benchmarks](https://docs.swmansion.com/react-native-executorch/docs/benchmarks/inference-time).

### General guidelines

- Prefer quantized model variants to reduce memory usage and storage requirements.
- Test on the lowest-spec device you plan to support.
- Implement a cloud API fallback for devices that cannot fit the model in memory.
- Monitor total downloaded model size and provide cleanup UI via `ResourceFetcher.deleteResources`.
- Always show loading states. Model loading and inference are long-running operations that take seconds to minutes.
