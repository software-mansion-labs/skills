---
name: on-device-ai
description: Integrate on-device AI into React Native apps using React Native ExecuTorch, which provides APIs for LLMs, computer vision, OCR, audio processing, and embeddings without cloud dependencies, as well as a variety of pre-exported models for common use cases. Use when user asks to build AI features into mobile apps - AI chatbots, image classification, object detection, style transfer, OCR, document parsing, speech processing, or semantic search - all running locally without cloud dependencies. Use when user mentions offline support, privacy, latency or cost concerns in AI-based applications.
---

## When to Use This Skill

Use this skill when you need to:

- **Build AI features directly into mobile apps** without cloud infrastructure
- **Deploy LLMs locally** for text generation, chat, or function calling
- **Add computer vision** (image classification, object detection, OCR)
- **Process audio** (speech-to-text, text-to-speech, voice activity detection)
- **Implement semantic search** with text embeddings
- **Ensure privacy** by keeping all AI processing on-device
- **Reduce latency** by eliminating cloud API calls
- **Work offline** once models are downloaded

## Overview

[React Native Executorch](https://github.com/software-mansion/react-native-executorch) is a library developed by [Software Mansion](https://swmansion.com/) that enables on-device AI model execution in React Native applications. It provides APIs for running machine learning models directly on mobile devices without requiring cloud infrastructure or internet connectivity (after initial model download). React Native Executorch provides APIs for LLMs, computer vision, OCR, audio processing and embeddings without cloud dependencies, as well as a variety of pre-exported models for common use cases. React Native Executorch is a way of bringing ExecuTorch into the React Native world.

## Key Use Cases

### Use Case 1: Mobile Chatbot/Assistant

**Trigger:** User asks to build a chat interface, create a conversational AI, or add an AI assistant to their app

**Steps:**

1. Choose appropriate LLM based on device memory constraints
2. Load model using ExecuTorch hooks
3. Implement message handling and conversation history
4. Optionally add system prompts, tool calling, or structured output

**Result:** Functional chat interface with on-device AI responding without cloud dependency

**Reference:** [./references/reference-llms.md](./references/reference-llms.md)

---

### Use Case 2: Computer Vision

**Trigger:** User needs to classify images, detect objects, or recognize content in photos

**Steps:**

1. Select vision model (classification, detection, or segmentation)
2. Load model for image processing task
3. Pass image and process results
4. Display detections or classifications in app UI

**Result:** App that understands image content without sending data to servers

**Reference:** [./references/reference-cv.md](./references/reference-cv.md)

---

### Use Case 3: Document/Receipt Scanning

**Trigger:** User wants to extract text from photos (receipts, documents, business cards)

**Steps:**

1. Choose OCR model matching target language
2. Load appropriate recognizer for alphabet/language
3. Capture or load image
4. Extract text regions with bounding boxes
5. Post-process results for application

**Result:** OCR-enabled app that reads text directly from device camera

**Reference:** [./references/reference-ocr.md](./references/reference-ocr.md)

---

### Use Case 4: Voice Interface

**Trigger:** User wants to add voice commands, transcription, or voice output to app

**Steps:**

- **For voice input:** Capture audio at correct sample rate → transcribe with STT model
- **For voice output:** Generate speech from text → play through audio context
- Handle audio format/sample rate conversion

**Result:** App with hands-free voice interaction

**Reference:** [./references/reference-audio.md](./references/reference-audio.md)

---

### Use Case 5: Semantic Search

**Trigger:** User needs intelligent search, similarity matching, or content recommendations

**Steps:**

1. Load text or image embeddings model
2. Generate embeddings for searchable content
3. Compute similarity scores between queries and content
4. Rank and return results

**Result:** Smart search that understands meaning, not just keywords

**Reference:** [./references/reference-nlp.md](./references/reference-nlp.md)

---

## Core Capabilities by Category

### Large Language Models (LLMs)

Run text generation, chat, function calling, and structured output generation locally on-device.

**Supported features:**

- Text generation and chat completions
- Function/tool calling
- Structured output with JSON schema validation
- Streaming responses
- Multiple model families (Llama 3.2, Qwen 3, Hammer 2.1, SmolLM2, Phi 4)

**Reference:** See [./references/reference-llms.md](./references/reference-llms.md)

---

### Computer Vision

Perform image understanding and manipulation tasks entirely on-device.

**Supported tasks:**

- **Image Classification** - Categorize images into predefined classes
- **Object Detection** - Locate and identify objects with bounding boxes
- **Image Segmentation** - Pixel-level classification
- **Style Transfer** - Apply artistic styles to images
- **Text-to-Image** - Generate images from text descriptions
- **Image Embeddings** - Convert images to numerical vectors for similarity/search

**Reference:** See [./references/reference-cv.md](./references/reference-cv.md) and [./references/reference-cv-2.md](./references/reference-cv-2.md)

---

### Optical Character Recognition (OCR)

Extract and recognize text from images with support for multiple languages and text orientations.

**Supported features:**

- Text detection in images
- Text recognition across different alphabets
- Horizontal text (standard documents, receipts)
- Vertical text support (experimental, for CJK languages)
- Multi-language support with language-specific recognizers

**Reference:** See [./references/reference-ocr.md](./references/reference-ocr.md)

---

### Audio Processing

Convert between speech and text, and detect speech activity in audio.

**Supported tasks:**

- **Speech-to-Text** - Transcribe audio to text (supports multiple languages including English)
- **Text-to-Speech** - Generate natural-sounding speech from text
- **Voice Activity Detection** - Detect speech segments in audio

**Reference:** See [./references/reference-audio.md](./references/reference-audio.md)

---

### Natural Language Processing

Convert text to numerical representations for semantic understanding and search.

**Supported tasks:**

- **Text Embeddings** - Convert text to vectors for similarity/search
- **Tokenization** - Convert text to tokens and vice versa

**Reference:** See [./references/reference-nlp.md](./references/reference-nlp.md)

---

## Understanding Model Loading

Before using any AI model, you need to load it. Models can be loaded from three sources:

**1. Bundled with app (assets folder)**

- Best for small models (< 512MB)
- Available immediately without download
- Increases app installation size

**2. Remote URL (downloaded on first use)**

- Best for large models (> 512MB)
- Downloaded once and cached locally
- Keeps app size small
- Requires internet on first use

**3. Local file system**

- Maximum flexibility for user-managed models
- Requires custom download/file management UI

**Model selection strategy:**

1. Small models (< 512MB) → Bundle with app or download from URL
2. Large models (> 512MB) → Download from URL on first use with progress tracking
3. Quantized models → Preferred for lower-end devices to save memory

**Reference:** [./references/reference-models.md](./references/reference-models.md) - Loading Models section

---

## Device Constraints and Model Selection

Not all models work on all devices. Consider these constraints:

**Memory limitations:**

- Low-end devices: Use smaller models (135M-1.7B parameters) and quantized variants
- High-end devices: Can run larger models (3B-4B parameters)

**Processing power:**

- Lower-end devices: Expect longer inference times
- Audio processing requires specific sample rates (16kHz for STT and VAD, 24kHz for TTS output)

**Storage:**

- Large models require significant disk space
- Implement cleanup mechanisms to remove unused models
- Monitor total downloaded model size

**Guidance:**

- Always check model memory requirements before recommending models
- Prefer quantized model variants on lower-end devices
- Show download progress for models > 512MB
- Test on target devices before release

**Reference:** [./references/reference-models.md](./references/reference-models.md)

---

## Important Technical Requirements

### Audio Processing

Audio must be in correct sample rate for processing:

- **Speech-to-Text or VAD input:** 16kHz sample rate
- **Text-to-Speech output:** 24kHz sample rate
- Always decode/resample audio to correct rate before processing

**Reference:** [./references/reference-audio.md](./references/reference-audio.md)

### Image Processing

Images can be provided as:

- Remote URLs (http/https) - automatically cached
- Local file URIs (file://)
- Base64-encoded strings

Image preprocessing (resizing, normalization) is handled automatically by most hooks.

**Reference:** [./references/reference-cv.md](./references/reference-cv.md) and [./references/reference-cv-2.md](./references/reference-cv-2.md)

### Text Tokens

Text embeddings and LLMs have maximum token limits. Text exceeding these limits will be truncated. Use `useTokenizer` to count tokens before processing.

**Reference:** [./references/reference-nlp.md](./references/reference-nlp.md)

---

## Core Utilities and Error Handling

The library provides core utilities for managing models and handling errors:

**ResourceFetcher:** Manage model downloads with pause/resume capabilities, storage cleanup, and progress tracking.

**Error Handling:** Use `RnExecutorchError` and error codes for robust error handling and user feedback.

**useExecutorchModule:** Low-level API for custom models not covered by dedicated hooks.

**Reference:** [./references/core-utilities.md](./references/core-utilities.md)

---

## Common Troubleshooting

**Model not loading:** Check model source URL/path validity and sufficient device storage

**Out of memory errors:** Switch to smaller model or quantized variant

**Poor LLM quality:** Adjust temperature/top-p parameters or improve system prompt

**Audio issues:** Verify correct sample rate (16kHz for STT and VAD, 24kHz output for TTS)

**Download failures:** Implement retry logic and check network connectivity

**Reference:** [./references/core-utilities.md](./references/core-utilities.md) for error handling details, or specific reference file for your use case

---

## Quick Reference by Hook

| Hook                   | Purpose                                   | Reference                                             |
| ---------------------- | ----------------------------------------- | ----------------------------------------------------- |
| `useLLM`               | Text generation, chat, function calling   | [reference-llms.md](./references/reference-llms.md)   |
| `useClassification`    | Image categorization                      | [reference-cv.md](./references/reference-cv.md)       |
| `useObjectDetection`   | Object localization                       | [reference-cv.md](./references/reference-cv.md)       |
| `useImageSegmentation` | Pixel-level classification                | [reference-cv.md](./references/reference-cv.md)       |
| `useStyleTransfer`     | Artistic image filters                    | [reference-cv-2.md](./references/reference-cv-2.md)   |
| `useTextToImage`       | Image generation                          | [reference-cv-2.md](./references/reference-cv-2.md)   |
| `useImageEmbeddings`   | Image similarity/search                   | [reference-cv-2.md](./references/reference-cv-2.md)   |
| `useOCR`               | Text recognition (horizontal)             | [reference-ocr.md](./references/reference-ocr.md)     |
| `useVerticalOCR`       | Text recognition (vertical, experimental) | [reference-ocr.md](./references/reference-ocr.md)     |
| `useSpeechToText`      | Audio transcription                       | [reference-audio.md](./references/reference-audio.md) |
| `useTextToSpeech`      | Voice synthesis                           | [reference-audio.md](./references/reference-audio.md) |
| `useVAD`               | Voice activity detection                  | [reference-audio.md](./references/reference-audio.md) |
| `useTextEmbeddings`    | Text similarity/search                    | [reference-nlp.md](./references/reference-nlp.md)     |
| `useTokenizer`         | Text to tokens conversion                 | [reference-nlp.md](./references/reference-nlp.md)     |
| `useExecutorchModule`  | Custom model inference (advanced)         | [core-utilities.md](./references/core-utilities.md)   |

---

## Quick Checklist for Implementation

Use this when building AI features with ExecuTorch:

**Planning Phase**

- Identified what AI task you need (chat, vision, audio, search)
- Considered device memory constraints and target devices
- Chose appropriate model from available options
- Determined if cloud backup fallback is needed

**Development Phase**

- Selected correct hook for your task
- Configured model loading (bundled, remote URL, or local)
- Implemented proper error handling
- Added loading states for model operations
- Tested audio sample rates (if audio task)
- Set up resource management for large models

**Testing Phase**

- Tested on target minimum device
- Verified offline functionality works
- Checked memory usage doesn't exceed device limits
- Tested error handling (network, memory, invalid inputs)
- Measured inference time for acceptable UX

**Deployment Phase**

- Model bundling strategy decided (size/download tradeoff)
- Download progress UI implemented (if remote models)
- Version management plan for model updates
- User feedback mechanism for quality issues

---

## Best Practices

For detailed best practices on model selection, error handling, UX, resource management, and performance optimization, see [./references/core-utilities.md](./references/core-utilities.md).

---

## External Resources

- **Official Documentation:** https://docs.swmansion.com/react-native-executorch
- **HuggingFace Models:** https://huggingface.co/software-mansion/collections
- **GitHub Repository:** https://github.com/software-mansion/react-native-executorch
- **API Reference:** https://docs.swmansion.com/react-native-executorch/docs/api-reference
- **Software Mansion:** https://swmansion.com/
