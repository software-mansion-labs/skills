# Gemini Live Integration

Fishjam ships an opinionated helper for wiring **Google's Gemini Live** (the multimodal `@google/genai` API) into a `FishjamAgent`. Source: `packages/js-server-sdk/src/integrations/gemini.ts` in the `js-server-sdk` repo (<https://github.com/fishjam-cloud/js-server-sdk>). Subpath import:

```ts
import {
  createClient,
  geminiInputAudioSettings,
  geminiOutputAudioSettings,
  inputMimeType,
} from '@fishjam-cloud/js-server-sdk/gemini';
```

`@google/genai` is an optional peer dependency - you only need to install it if you import from the `'@fishjam-cloud/js-server-sdk/gemini'` path.

## What the helper provides

| Export | What it is |
|---|---|
| `createClient(options)` | A thin wrapper around `new GoogleGenAI(options)` |
| `geminiInputAudioSettings` | `{ audioFormat: 'pcm16', audioSampleRate: 16000 }` — pass to `createAgent` as `output` so other peers' audio is downmixed into the format Gemini expects. |
| `geminiOutputAudioSettings` | `{ encoding: 'pcm16', sampleRate: 24000, channels: 1 }` — pass to `agent.createTrack(...)` so Gemini's TTS output flows back into the room. |
| `inputMimeType` | `'audio/pcm;rate=16000'` — the MIME string to attach to audio chunks sent to Gemini. |

## End-to-end pattern

```ts
import { FishjamClient } from '@fishjam-cloud/js-server-sdk';
import {
  createClient,
  geminiInputAudioSettings,
  geminiOutputAudioSettings,
  inputMimeType,
} from '@fishjam-cloud/js-server-sdk/gemini';
import { Modality } from '@google/genai';

const fishjamClient = new FishjamClient({
  fishjamId: process.env.FISHJAM_ID!,
  managementToken: process.env.FISHJAM_MANAGEMENT_TOKEN!,
});

const genai = createClient({ apiKey: process.env.GEMINI_API_KEY! });

async function spawnGeminiAgent(roomId: RoomId) {
  const { agent, peer } = await fishjamClient.createAgent(
    roomId,
    {
      output: geminiInputAudioSettings,
    },
  );

  // Create the Fishjam output track FIRST so the closure below captures a real reference.
  const geminiTrack = agent.createTrack(geminiOutputAudioSettings, { source: 'gemini' });

  const session = await genai.live.connect({
    model: 'gemini-2.0-flash-exp',
    config: {
      responseModalities: [Modality.AUDIO],
      systemInstruction: 'You are a helpful assistant.',
      outputAudioTranscription: {},
    },
    callbacks: {
      onmessage: (event) => {
        // Gemini → Fishjam: forward inbound audio to the room
        if (event.data) agent.sendData(geminiTrack.id, Buffer.from(event.data, 'base64'));
      },
      onerror: (e) => console.error('gemini error', e),
      onclose: () => agent.disconnect(),
    },
  });

  // Fishjam → Gemini: forward each peer's audio chunk
  agent.on('trackData', (msg) => {
    session.sendRealtimeInput({
      audio: { data: Buffer.from(msg.data).toString('base64'), mimeType: inputMimeType },
    });
  });

  return { agent, peer, session };
}
```

The full reference implementation lives in the `js-server-sdk` repo under `examples/multimodal/` (full duplex audio with `outputAudioTranscription` and `createTrack`) and `examples/gemini-demo/backend/` (production-shape end-to-end demo). `examples/transcription/` is text-only — no `createTrack`, no audio echo back — so it won't match this pattern.

## When to use the helper vs raw

- **Use the helper** for the canonical 16kHz-in / 24kHz-out Gemini Live flow. It saves you from picking codec parameters and getting them wrong.
- **Skip the helper** and import `@google/genai` directly if you need a different model with different audio formats, or if you're using Gemini's text-only modes (the helper is only ergonomic when you actually want voice).

## Image capture for visual prompting

`agent.captureImage(trackId)` returns a still frame from a video track — useful for "look at what I'm showing and answer". Pair with Gemini's multimodal prompts.

## Sources

- `js-server-sdk` repo: `packages/js-server-sdk/src/integrations/gemini.ts`
- Canonical reference implementations: `examples/multimodal/src/service/multimodal.ts` and `examples/gemini-demo/backend/src/agents.ts` in the same repo
- <https://fishjam.swmansion.com/docs/integrations/gemini-live-integration>
