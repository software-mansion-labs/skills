# FishjamAgent

A server-side, programmatic peer that joins a room and exchanges audio frames with other peers. Used for AI voice integrations, custom server-side audio processing, and managed voice agents.

Source: `packages/js-server-sdk/src/agent.ts` in the `js-server-sdk` repo (<https://github.com/fishjam-cloud/js-server-sdk>).

> For the **Gemini Live** integration that wraps this class, see `gemini-integration.md`. For Vapi (managed voice agents), use `client.createVapiAgent` — covered briefly in `client.md`.

## Lifecycle

```ts
import type { AgentCallbacks } from '@fishjam-cloud/js-server-sdk';

const { agent, peer } = await fishjamClient.createAgent(
  roomId,
  {                                      // PeerOptionsAgent
    output: { audioFormat: 'pcm16', audioSampleRate: 16000 }, // audioFormat only 'pcm16'; audioSampleRate only 16000 | 24000 (pick 24000 for Gemini Live)
  },
  {                                      // AgentCallbacks
    onError: (event) => console.error('agent ws error', event),
    onClose: (code, reason) => console.log('agent ws closed', code, reason),
  } satisfies AgentCallbacks,
);

// ...interact with audio...
agent.disconnect();
```

`agent` is a typed `EventEmitter`, ready as soon as `createAgent` resolves — no need to call `agent.awaitConnected()` yourself (the SDK awaits it inside `createAgent`).

## Outgoing audio

```ts
import type { AudioCodecParameters, TrackId } from '@fishjam-cloud/js-server-sdk';

const codec: AudioCodecParameters = { encoding: 'opus', sampleRate: 48000, channels: 1 };
// or { encoding: 'pcm16', sampleRate: 24000, channels: 1 }
// or { encoding: 'pcm16', sampleRate: 16000, channels: 1 }

const track = agent.createTrack(codec, { tag: 'tts-stream' });
// track: { id: TrackId, type, metadata }

agent.sendData(track.id, audioBytes /* Uint8Array */); // push frames
agent.interruptTrack(track.id);                         // drop any buffered playback
agent.deleteTrack(track.id);                            // remove the track from the room
```

Supported codecs (from the source):

- `opus` → encoded internally as `TRACK_ENCODING_OPUS`.
- `pcm16` → `TRACK_ENCODING_PCM16`. Sample rates: 16000, 24000, 48000.

Frames go onto the WebSocket as binary `AgentRequest.trackData` messages. Pace yourself — there's no built-in backpressure; if you push faster than realtime, you'll inflate latency.

## Incoming audio (`trackData` event)

```ts
import type { IncomingTrackData } from '@fishjam-cloud/js-server-sdk';

agent.on('trackData', (msg: IncomingTrackData) => {
  // msg.peerId (PeerId)
  // msg.track  (Track | undefined — check before use, e.g. msg.track?.id)
  // msg.data   (Uint8Array of audio frames)
});
```

What you receive is governed by `output` you set on `PeerOptionsAgent`. For example, `output: { audioFormat: 'pcm16', audioSampleRate: 16000 }` downmixes every peer's audio into pcm16 at 16kHz so the agent can feed it straight into an STT or LLM.

## Image capture

If you want to "see" what a peer is broadcasting on a video track, ask for a single frame:

```ts
const image = await agent.captureImage(trackId, 5000 /* timeout ms */);
// image: { trackId, contentType: string (e.g. 'image/jpeg'), data: Uint8Array }
```

One pending capture per track at a time. The second concurrent call on the same track returns a **rejected promise** (it does not throw synchronously), with message `captureImage already pending for track <trackId>`. Timeouts also reject, they don't throw — so always `await` inside `try/catch` or chain `.catch()`.

## Disconnecting

```ts
agent.disconnect();
```

Closes the WebSocket. Any pending `captureImage` promises reject. Call this from your room teardown — the agent does not exit on its own.

## Vapi (managed alternative)

If you just want a hosted voice agent and don't need fine-grained audio control, use `createVapiAgent` instead. It returns a peer (no client `agent` object) — Vapi manages the WebSocket and AI inference for you.

```ts
const { peer } = await fishjamClient.createVapiAgent(roomId, vapiOptions);
// no agent object, no track APIs — Vapi handles them
```

## Sources

- `js-server-sdk` repo: `packages/js-server-sdk/src/agent.ts`, `packages/js-server-sdk/src/client.ts` (`createAgent`, `createVapiAgent`)
- `packages/fishjam-proto/src/fishjam/agent_notifications.ts` (`AgentResponse_TrackData`, `AgentResponse_TrackImage`)
- Canonical example: `examples/transcription/src/service/transcription.ts` in the same repo
- <https://fishjam.swmansion.com/docs/tutorials/agents>
- <https://fishjam.swmansion.com/docs/explanation/agent-internals>
