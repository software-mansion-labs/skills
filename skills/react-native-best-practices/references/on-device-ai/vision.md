# Computer Vision

Production patterns for on-device computer vision in React Native using React Native ExecuTorch. For hook API signatures and model constants, webfetch the relevant page from the [official docs](https://docs.swmansion.com/react-native-executorch/docs/).

For model loading and resource fetcher setup, see **`setup.md`**.

---

## Hook Overview

All vision hooks share a common interface pattern:

| Hook | Task | Input | Output |
|---|---|---|---|
| [useClassification](https://docs.swmansion.com/react-native-executorch/docs/hooks/computer-vision/useClassification) | Label an image | Image URI | `{ label: probability }` object |
| [useObjectDetection](https://docs.swmansion.com/react-native-executorch/docs/hooks/computer-vision/useObjectDetection) | Locate objects | Image URI | `Detection[]` with bbox, label, score |
| [useOCR](https://docs.swmansion.com/react-native-executorch/docs/hooks/computer-vision/useOCR) | Read horizontal text | Image URI | `OCRDetection[]` with bbox, text, score |
| [useVerticalOCR](https://docs.swmansion.com/react-native-executorch/docs/hooks/computer-vision/useVerticalOCR) | Read vertical text | Image URI | `OCRDetection[]` with bbox, text, score |
| [useImageSegmentation](https://docs.swmansion.com/react-native-executorch/docs/hooks/computer-vision/useImageSegmentation) | Pixel-level labels | Image URI | Segmentation mask dictionary |
| [useStyleTransfer](https://docs.swmansion.com/react-native-executorch/docs/hooks/computer-vision/useStyleTransfer) | Apply art style | Image URI | Base64-encoded image URL |
| [useTextToImage](https://docs.swmansion.com/react-native-executorch/docs/hooks/computer-vision/useTextToImage) | Generate image from text | Text prompt | Base64 PNG |
| [useImageEmbeddings](https://docs.swmansion.com/react-native-executorch/docs/hooks/computer-vision/useImageEmbeddings) | Image to vector | Image URI | `Float32Array` embedding |
| [useTextEmbeddings](https://docs.swmansion.com/react-native-executorch/docs/hooks/natural-language-processing/useTextEmbeddings) | Text to vector | String | `Float32Array` embedding |

### Common interface

Every vision hook returns an object with:

- `isReady` -- model is loaded and ready for inference
- `isGenerating` -- inference is in progress
- `error` -- error object if loading or inference failed
- `downloadProgress` -- 0 to 1 during model download
- `forward(input)` -- run inference on a single image (returns a Promise)

Image inputs accept: remote URLs (`https://...`), local file URIs (`file:///...`), or base64-encoded strings.

---

## Image Classification

Assigns a label to an image. Returns an object mapping ImageNet1k class names (1000 classes) to probabilities:

```tsx
import { useClassification, EFFICIENTNET_V2_S } from 'react-native-executorch';

const model = useClassification({ model: EFFICIENTNET_V2_S });

const classify = async (imageUri: string) => {
  const scores = await model.forward(imageUri);

  // Get top 3 predictions
  const top3 = Object.entries(scores)
    .sort(([, a], [, b]) => b - a)
    .slice(0, 3)
    .map(([label, score]) => ({ label, score }));

  return top3;
};
```

If multiple classes have similar probabilities, the model is not confident in its prediction.

---

## Object Detection

Returns a list of detected objects with bounding boxes, labels, and confidence scores:

```tsx
import { useObjectDetection, SSDLITE_320_MOBILENET_V3_LARGE } from 'react-native-executorch';

const model = useObjectDetection({ model: SSDLITE_320_MOBILENET_V3_LARGE });

const detect = async (imageUri: string) => {
  const detections = await model.forward(imageUri);

  for (const det of detections) {
    console.log(det.label, det.score, det.bbox); // bbox: { x1, y1, x2, y2 }
  }
};
```

Bounding box coordinates are bottom-left (`x1`, `y1`) and top-right (`x2`, `y2`) in the original image's pixel space. The `label` corresponds to one of 91 COCO class labels.

---

## OCR (Optical Character Recognition)

`useOCR` reads horizontal text. `useVerticalOCR` reads vertical text (e.g., Japanese, Chinese). Both return `OCRDetection[]` with bounding boxes, text, and confidence scores.

For hook APIs, webfetch [useOCR](https://docs.swmansion.com/react-native-executorch/docs/hooks/computer-vision/useOCR) and [useVerticalOCR](https://docs.swmansion.com/react-native-executorch/docs/hooks/computer-vision/useVerticalOCR).

```tsx
import { useOCR, OCR_ENGLISH } from 'react-native-executorch';

const model = useOCR({ model: OCR_ENGLISH });

const readText = async (imageUri: string) => {
  const results = await model.forward(imageUri);
  for (const det of results) {
    console.log(det.text, det.score, det.bbox); // bbox: Point[] (4 corners)
  }
};
```

### Language support

Each supported alphabet requires its own recognizer model. The simplified language constants (e.g., `OCR_ENGLISH`, `OCR_RUSSIAN`, `OCR_JAPANESE`) bundle the correct detector and recognizer automatically. For the full list of 84 supported languages, webfetch [OCR Supported Alphabets](https://docs.swmansion.com/react-native-executorch/docs/api-reference#ocr-supported-alphabets).

When using custom recognizers, ensure the recognizer alphabet matches the language:
- `RECOGNIZER_LATIN_CRNN` for Latin-alphabet languages (Polish, German, etc.)
- `RECOGNIZER_CYRILLIC_CRNN` for Cyrillic-alphabet languages (Russian, Ukrainian, etc.)

The detector model is CRAFT (text detection); recognizers are CRNN (text recognition).

---

## Image Segmentation

Assigns a class label to every pixel in an image. Useful for background removal, scene understanding, and portrait effects. For the full API, webfetch [useImageSegmentation](https://docs.swmansion.com/react-native-executorch/docs/hooks/computer-vision/useImageSegmentation).

```tsx
import { useImageSegmentation, DEEPLAB_V3_RESNET50, DeeplabLabel } from 'react-native-executorch';

const model = useImageSegmentation({ model: DEEPLAB_V3_RESNET50 });

const segment = async (imageUri: string) => {
  // forward(imageUri, classesOfInterest?, resize?)
  const outputDict = await model.forward(imageUri, [DeeplabLabel.CAT], true);

  // outputDict[DeeplabLabel.ARGMAX]: per-pixel class index (always present)
  // outputDict[DeeplabLabel.CAT]: per-pixel probability for CAT class
};
```

- `classesOfInterest` (optional): `DeeplabLabel[]` specifying which classes to return full probability arrays for. Default is empty (only argmax returned).
- `resize` (optional): if `true`, output is rescaled to original image dimensions. Default is `false` (224x224 internal resolution). Setting to `true` makes `forward` slower.

The model supports 21 `DeeplabLabel` classes.

---

## Style Transfer

Applies an artistic style to an image. Returns a base64-encoded image URL. For the full API, webfetch [useStyleTransfer](https://docs.swmansion.com/react-native-executorch/docs/hooks/computer-vision/useStyleTransfer).

Available models: `STYLE_TRANSFER_CANDY`, `STYLE_TRANSFER_MOSAIC`, `STYLE_TRANSFER_UDNIE`, `STYLE_TRANSFER_RAIN_PRINCESS`.

---

## Text-to-Image

Generates an image from a text prompt using a compressed Stable Diffusion model. For the full API, webfetch [useTextToImage](https://docs.swmansion.com/react-native-executorch/docs/hooks/computer-vision/useTextToImage).

```tsx
import { useTextToImage, BK_SDM_TINY_VPRED_512 } from 'react-native-executorch';

const model = useTextToImage({
  ...BK_SDM_TINY_VPRED_512,
  inferenceCallback: (progress) => console.log(`Step: ${progress}`),
});

const generate = async () => {
  // generate(prompt, imageSize, numSteps, seed?)
  const base64Png = await model.generate('a cat sitting on a couch', 512, 25);
};
```

- `imageSize` must be a multiple of 32 (256 or 512 supported)
- `numSteps`: number of denoising iterations
- `seed` (optional): for reproducible results

Available models: `BK_SDM_TINY_VPRED_256`, `BK_SDM_TINY_VPRED_512`.

---

## Embeddings

Convert images or text into vector representations for similarity search, retrieval-augmented generation (RAG), or clustering.

```tsx
import { useImageEmbeddings, CLIP_VIT_BASE_PATCH32_IMAGE } from 'react-native-executorch';

const model = useImageEmbeddings({ model: CLIP_VIT_BASE_PATCH32_IMAGE });

const embedding = await model.forward(imageUri);
// embedding is a Float32Array (512 dimensions, normalized)
```

Text embeddings follow the same pattern with `useTextEmbeddings`. Available text embedding models:

| Model | Dimensions | Max tokens | Best for |
|---|---|---|---|
| `ALL_MINILM_L6_V2` | 384 | 254 | General purpose |
| `ALL_MPNET_BASE_V2` | 768 | 382 | General purpose (higher quality) |
| `MULTI_QA_MINILM_L6_COS_V1` | 384 | 509 | Search / QA |
| `MULTI_QA_MPNET_BASE_DOT_V1` | 768 | 510 | Search / QA (higher quality) |
| `CLIP_VIT_BASE_PATCH32_TEXT` | 512 | 74 | Cross-modal search with CLIP images |

Combine image and text embeddings from the same model family (CLIP) for cross-modal search. Embeddings are normalized, so cosine similarity equals dot product.

Use `useTokenizer` to count tokens before processing variable-length input. Text exceeding the model's token limit is truncated silently.

---

## VisionCamera Real-Time Frame Processing

Vision hooks that support `runOnFrame` can process camera frames in real time using VisionCamera v5.

### Supported hooks

`useClassification`, `useObjectDetection`, `useOCR`, `useVerticalOCR`, `useImageEmbeddings`, `useImageSegmentation`, `useStyleTransfer`.

### runOnFrame vs forward

| | `runOnFrame` | `forward` |
|---|---|---|
| Thread | JS worklet thread (synchronous) | Background thread (async) |
| Input | VisionCamera `Frame` | Image URI |
| Use case | Real-time camera | Single image |

### Setup

Requires `react-native-vision-camera` v5 and `react-native-worklets`.

```tsx
import { useState, useCallback } from 'react';
import { Text, StyleSheet } from 'react-native';
import {
  Camera,
  Frame,
  useCameraDevices,
  useCameraPermission,
  useFrameOutput,
} from 'react-native-vision-camera';
import { scheduleOnRN } from 'react-native-worklets';
import { useClassification, EFFICIENTNET_V2_S } from 'react-native-executorch';

function LiveClassifier() {
  const { hasPermission, requestPermission } = useCameraPermission();
  const devices = useCameraDevices();
  const device = devices.find((d) => d.position === 'back');
  const model = useClassification({ model: EFFICIENTNET_V2_S });
  const [topLabel, setTopLabel] = useState('');

  const runOnFrame = model.runOnFrame;

  const frameOutput = useFrameOutput({
    pixelFormat: 'rgb',           // Required: must be 'rgb'
    dropFramesWhileBusy: true,    // Skip frames during slow inference
    onFrame: useCallback(
      (frame: Frame) => {
        'worklet';
        if (!runOnFrame) return;
        try {
          const scores = runOnFrame(frame);
          if (scores) {
            let best = '';
            let bestScore = -1;
            for (const [label, score] of Object.entries(scores)) {
              if ((score as number) > bestScore) {
                bestScore = score as number;
                best = label;
              }
            }
            scheduleOnRN(setTopLabel, best);
          }
        } finally {
          frame.dispose(); // Always dispose to avoid memory leaks
        }
      },
      [runOnFrame]
    ),
  });

  if (!hasPermission) {
    requestPermission();
    return null;
  }
  if (!device) return null;

  return (
    <>
      <Camera style={styles.camera} device={device} outputs={[frameOutput]} isActive />
      <Text style={styles.label}>{topLabel}</Text>
    </>
  );
}
```

### Gotchas

- **`pixelFormat: 'rgb'` is mandatory.** The default VisionCamera format is `yuv`, which produces scrambled results.
- **Always call `frame.dispose()` in a `finally` block.** Skipping this leaks memory and crashes after processing many frames.
- **Guard `runOnFrame` for null.** It is `null` until the model finishes loading. Check `if (!runOnFrame) return` inside `onFrame`.
- **Use `dropFramesWhileBusy: true`** for models with longer inference times. Without it, the camera pipeline blocks.
- **`runOnFrame` is synchronous and runs on the JS worklet thread.** For models that take longer than a frame interval, consider VisionCamera's [async frame processing](https://react-native-vision-camera-v5-docs.vercel.app/docs/async-frame-processing).

### Module API with VisionCamera

When using the TypeScript Module API (e.g., `ClassificationModule`) instead of hooks, `runOnFrame` is a worklet function. React would invoke it as a state initializer if passed directly to `useState`. Use the functional updater form:

```tsx
const [module] = useState(() => new ClassificationModule());
const [runOnFrame, setRunOnFrame] = useState<typeof module.runOnFrame | null>(null);

useEffect(() => {
  module.load(EFFICIENTNET_V2_S).then(() => {
    // () => module.runOnFrame prevents React from calling it as initializer
    setRunOnFrame(() => module.runOnFrame);
  });
}, [module]);
```
