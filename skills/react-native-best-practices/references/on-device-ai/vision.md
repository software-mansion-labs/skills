# Computer Vision

Production patterns for on-device computer vision in React Native using React Native ExecuTorch. For hook API signatures and model constants, webfetch the relevant page from the [official docs](https://docs.swmansion.com/react-native-executorch/docs/).

For model loading and resource fetcher setup, see **`setup.md`**.

---

## Hook Overview

All vision hooks share a common interface pattern:

| Hook | Task | Input | Output |
|---|---|---|---|
| [useClassification](https://docs.swmansion.com/react-native-executorch/docs/hooks/computer-vision/useClassification) | Label an image | Image URI / PixelData | `{ label: probability }` object |
| [useObjectDetection](https://docs.swmansion.com/react-native-executorch/docs/hooks/computer-vision/useObjectDetection) | Locate objects | Image URI / PixelData | `Detection[]` with bbox, label, score |
| [useOCR](https://docs.swmansion.com/react-native-executorch/docs/hooks/computer-vision/useOCR) | Read horizontal text | Image URI / PixelData | Recognized text with bounding boxes |
| [useVerticalOCR](https://docs.swmansion.com/react-native-executorch/docs/hooks/computer-vision/useVerticalOCR) | Read vertical text | Image URI / PixelData | Recognized text with bounding boxes |
| [useSemanticSegmentation](https://docs.swmansion.com/react-native-executorch/docs/hooks/computer-vision/useSemanticSegmentation) | Pixel-level labels | Image URI / PixelData | Segmentation mask |
| [useStyleTransfer](https://docs.swmansion.com/react-native-executorch/docs/hooks/computer-vision/useStyleTransfer) | Apply art style | Image URI / PixelData | Styled image URI |
| [useTextToImage](https://docs.swmansion.com/react-native-executorch/docs/hooks/computer-vision/useTextToImage) | Generate image from text | Text prompt | Generated image URI |
| [useImageEmbeddings](https://docs.swmansion.com/react-native-executorch/docs/hooks/computer-vision/useImageEmbeddings) | Image to vector | Image URI / PixelData | `number[]` embedding |
| [useTextEmbeddings](https://docs.swmansion.com/react-native-executorch/docs/hooks/natural-language-processing/useTextEmbeddings) | Text to vector | String | `number[]` embedding |

### Common interface

Every vision hook returns an object with:

- `isReady` -- model is loaded and ready for inference
- `isGenerating` -- inference is in progress
- `error` -- error object if loading or inference failed
- `downloadProgress` -- 0 to 1 during model download
- `forward(input)` -- run inference on a single image (returns a Promise)

Image inputs accept: remote URLs (`https://...`), local file URIs (`file:///...`), base64-encoded strings, or `PixelData` objects (raw RGB buffer).

---

## Image Classification

Assigns a label to an image. Returns an object mapping class names to probabilities:

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
import { useObjectDetection, RF_DETR_NANO } from 'react-native-executorch';

const model = useObjectDetection({ model: RF_DETR_NANO });

const detect = async (imageUri: string) => {
  // detectionThreshold defaults to 0.7; lower it to find more objects
  const detections = await model.forward(imageUri, 0.5);

  for (const det of detections) {
    console.log(det.label, det.score, det.bbox); // bbox: { x1, y1, x2, y2 }
  }
};
```

Bounding box coordinates are in the original image's pixel space. The `label` property is typed to the model's label map (e.g., COCO labels).

---

## OCR (Optical Character Recognition)

`useOCR` reads horizontal text. `useVerticalOCR` reads vertical text (e.g., Japanese, Chinese). Both return recognized text with bounding boxes. For hook APIs, webfetch [useOCR](https://docs.swmansion.com/react-native-executorch/docs/hooks/computer-vision/useOCR) and [useVerticalOCR](https://docs.swmansion.com/react-native-executorch/docs/hooks/computer-vision/useVerticalOCR).

```tsx
import { useOCR, EASYOCR_DETECTION_MODEL } from 'react-native-executorch';

const model = useOCR({ model: EASYOCR_DETECTION_MODEL });

const readText = async (imageUri: string) => {
  const results = await model.forward(imageUri);
  // results contain recognized text and bounding boxes
};
```

For multilingual OCR, pass a `language` option. Check the docs for supported languages per model.

---

## Semantic Segmentation

Assigns a class label to every pixel in an image. Useful for background removal, scene understanding, and portrait effects. For the full API, webfetch [useSemanticSegmentation](https://docs.swmansion.com/react-native-executorch/docs/hooks/computer-vision/useSemanticSegmentation).

---

## Style Transfer

Applies an artistic style to an image. Returns a URI to the styled output image. For the full API, webfetch [useStyleTransfer](https://docs.swmansion.com/react-native-executorch/docs/hooks/computer-vision/useStyleTransfer).

---

## Text-to-Image

Generates an image from a text prompt. For the full API, webfetch [useTextToImage](https://docs.swmansion.com/react-native-executorch/docs/hooks/computer-vision/useTextToImage).

---

## Embeddings

Convert images or text into vector representations for similarity search, retrieval-augmented generation (RAG), or clustering.

```tsx
import { useImageEmbeddings, CLIP_VIT_BASE_PATCH16 } from 'react-native-executorch';

const model = useImageEmbeddings({ model: CLIP_VIT_BASE_PATCH16 });

const embedding = await model.forward(imageUri);
// embedding is a number[] vector
```

Text embeddings follow the same pattern with `useTextEmbeddings`. Combine image and text embeddings (using the same model family like CLIP) for cross-modal search.

Use `useTokenizer` to count tokens before processing variable-length input. Text exceeding the model's token limit is truncated silently.

---

## VisionCamera Real-Time Frame Processing

Vision hooks that support `runOnFrame` can process camera frames in real time using VisionCamera v5.

### Supported hooks

`useClassification`, `useObjectDetection`, `useOCR`, `useVerticalOCR`, `useImageEmbeddings`, `useSemanticSegmentation`, `useStyleTransfer`.

### runOnFrame vs forward

| | `runOnFrame` | `forward` |
|---|---|---|
| Thread | JS worklet thread (synchronous) | Background thread (async) |
| Input | VisionCamera `Frame` | Image URI / PixelData |
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
