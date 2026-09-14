---
name: fishjam-composition
description: "Fishjam compositions: server-side, real-time mixing of several live streams (room peers, WHIP, RTMP, MP4) into new output streams with layouts, overlays, captions, and audio mixes, built on Smelter. Covers lifecycle and billing, static scenes, React templates, composing a room, recording, and the raw Composition API. Use when building a grid or picture-in-picture of a call, restreaming a room to a livestream or to YouTube / Twitch over RTMP, adding overlays or captions driven from a backend, or recording a composed layout. Trigger on: 'Fishjam composition', 'compose a room', 'CompositionClient', 'createComposition', 'create_composition', 'registerWhipOutput', 'register_whip_output', 'registerTemplateOutput', 'register_template_output', 'updateOutput', 'update_output', 'forwardRoomTracks', 'forward_room_tracks', 'track_forwardings', 'sendEvent', 'send_event', '@fishjam-cloud/composition', '@fishjam-cloud/composition-cli', 'usePeers', 'eventBus', 'rtc.fishjam.io', 'createRecording', 'Smelter'."
license: MIT
---

# Fishjam Compositions

A composition takes several live streams, lays them out in a scene, and pushes the result out as new streams in real time. Rendering happens on Fishjam's servers, built on Software Mansion's Smelter.

Because the layout is rendered once on the server, every viewer receives one ready-made stream: it looks the same on every device, and a low-end phone only has to play a single video instead of decoding every participant.

> **Read `../platform/SKILL.md` first.** It defines rooms, peers, tracks, management tokens, livestream rooms, and notifications that this skill builds on.

## What you can build

- **A video call turned into a broadcast.** The room's participants in a grid or speaker layout, with names, speaking highlights, and camera-off placeholders, streamed to a livestream audience in your app.
- **Restreaming.** A call, an OBS feed, or any other input sent to YouTube, Twitch, or another RTMP platform.
- **Broadcast graphics driven by your backend.** A LIVE badge, captions, lower thirds, sponsor banners, or a score bug, toggled as things happen in your product.
- **A standby or placeholder channel.** Graphics only, or a looping video, switched to the live layout when the show starts. An output shows one scene at a time; events or updates switch it.
- **The show as a file.** A recording captures one output exactly as viewers see it, including every layout change, so recording a live show takes one extra call.

## Mental model

```
inputs                       composition                 outputs
room tracks (forwarded) ──┐                            ┌─▶ WHIP → Fishjam livestream → viewers
WHIP publisher ───────────┼─▶ scene per output ────────┼─▶ RTMP → YouTube / Twitch / other
RTMP encoder ─────────────┤   video tree + audio mix   └─▶ recording of an output
MP4 file URL ─────────────┘
```

| Concept | What it is |
|---|---|
| **Composition** | One billed rendering session on the Composition API (`https://rtc.fishjam.io`), authenticated with the management token (`../platform/auth-model.md`). Holds inputs, outputs, and assets. |
| **Input** | A stream coming in, addressed by an `input_id` you choose. Registered by you (WHIP, WHEP, RTMP, MP4) or created for you when a Fishjam room is forwarded. |
| **Output** | A stream going out, addressed by an `output_id` you choose. Pushes over WHIP or RTMP. Each output renders its own scene. |
| **Template** | A React component, built into a JS bundle and uploaded with an output, that renders the output's scene and re-renders when room peers change or your backend sends an event. |
| **Static scene** | The alternative to a template: a JSON component tree (`view`, `tiles`, `rescaler`, `input_stream`, `text`, `image`) plus an audio mix, sent by your backend. |
| **Assets** | Images (registered by URL) and fonts (uploaded), usable from templates and scenes. |
| **Recording** | A file capture of one output, managed through the Fishjam Server API, not the Composition API. |

**A composition produces video; it does not serve it.** There is no URL a player can open. Every output must push to a destination that viewers can reach.

## Designing a composition

Answer three questions: where the streams come from, where the result goes, and what draws the layout.

### Where do the streams come from?

| Source | How it gets in | Read |
|---|---|---|
| Peers in a Fishjam room (`../platform/room-types.md`) | Forward the whole room; inputs appear per track | `room-composition.md` |
| OBS, a browser, or any WHIP publisher | WHIP input; the publisher pushes to a URL you get back | `inputs-and-outputs.md` |
| A hardware or software RTMP encoder | RTMP input; the encoder pushes to a URL you get back | `inputs-and-outputs.md` |
| A remote WHEP stream | WHEP input; the composition pulls it | `inputs-and-outputs.md` |
| A video file | MP4 input by URL, optionally looped | `inputs-and-outputs.md` |

### Where does the result go?

| Destination | Output | Read |
|---|---|---|
| Your own web or mobile app | WHIP output into a Fishjam livestream room. Your app's users join that room as livestream viewers with the client SDK (`../react-client/livestream.md`, same hooks in React Native). | `inputs-and-outputs.md` |
| YouTube, Twitch, or another streaming platform | RTMP output to the platform's ingest URL and stream key | `inputs-and-outputs.md` |
| A file to download later | Recording of any output | `recording.md` |

### What draws the layout?

**Use a template unless the layout is fixed.** Most real compositions react to something: peers joining a call, someone speaking, a caption or "live" badge toggled from the backend. With static scenes your backend has to track all of that state and rebuild and resend the scene on every change; a template does it inside the composition.

| | Template (recommended) | Static scene |
|---|---|---|
| Layout logic | React component running next to the composition | JSON tree built by your backend |
| Room peers join or leave | Re-renders by itself (`usePeers`) | Backend rebuilds and resends the whole scene |
| Backend-driven change | `sendEvent` with a small payload | `updateOutput` with the whole scene |
| Extra setup | Build step with `@fishjam-cloud/composition-cli` | None |
| Pick it for | Room composition, overlays, captions, anything dynamic | A few known inputs in a fixed layout, quick prototypes |

Templates are written in TypeScript and React whatever the backend language; a Python backend uploads the built bundle the same way. Read `templates.md` for templates and `scenes.md` for the component tree both approaches share.

## Flow and SDK calls

Both server SDKs ship a `CompositionClient`, authenticated with the same management token as `FishjamClient` (no Fishjam ID needed).

| Step | JS (`@fishjam-cloud/js-server-sdk`) | Python (`fishjam`) |
|---|---|---|
| 1. Create | `createComposition()` | `create_composition()` |
| 2. Register assets | `registerImage`, `registerFont` | `register_image`, `register_font` |
| 3. Register inputs, or forward a room | `registerMp4Input`, `registerWhipInput`, …, or `fishjamClient.forwardRoomTracks` | `register_mp4_input`, `register_whip_input`, …, or `fishjam_client.forward_room_tracks` |
| 4. Register outputs | `registerTemplateOutput`, `registerWhipOutput`, `registerRtmpOutput` | `register_template_output`, `register_whip_output`, `register_rtmp_output` |
| 5. Change while live | `sendEvent`, `updateOutput` | `send_event`, `update_output` |
| 6. Tear down | `deleteComposition` | `delete_composition` |

Register in dependency order: assets and inputs before the scenes that reference them, and the destination (for example a livestream room) before the output that pushes to it.

## Minimal flow

A looping MP4 composed into a Fishjam livestream room. It uses a static scene to stay self-contained with no build step; for a template version see `templates.md`. Python equivalent: `../python-server-sdk/composition.md`.

```ts
import { CompositionClient, FishjamClient } from '@fishjam-cloud/js-server-sdk';
import type { InputId, OutputId } from '@fishjam-cloud/js-server-sdk';

const managementToken = process.env.FISHJAM_MANAGEMENT_TOKEN!;
const fishjamClient = new FishjamClient({ fishjamId: process.env.FISHJAM_ID!, managementToken });
const compositionClient = new CompositionClient({ managementToken });

const livestream = await fishjamClient.createRoom({ roomType: 'livestream' });
const { token } = await fishjamClient.createLivestreamStreamerToken(livestream.id);

const { compositionId } = await compositionClient.createComposition();
await compositionClient.registerMp4Input(compositionId, 'movie' as InputId, {
  url: 'https://smelter.dev/videos/template-scene-race.mp4',
  loop: true,
});
await compositionClient.registerWhipOutput(compositionId, 'main' as OutputId, {
  endpointUrl: fishjamClient.livestreamWhipUrl(),
  bearerToken: token,
  video: {
    resolution: { width: 1280, height: 720 },
    initial: { root: { type: 'rescaler', child: { type: 'input_stream', inputId: 'movie' } } },
  },
  audio: { initial: { inputs: [{ inputId: 'movie' }] } },
});

await compositionClient.deleteComposition(compositionId);
await fishjamClient.deleteRoom(livestream.id);
```

Between registering the output and deleting, viewers watch `livestream.id` like any Fishjam livestream. Full signatures: `../js-server-sdk/composition.md` and `../python-server-sdk/composition.md`.

## Lifecycle and billing

- **Billed per minute for each registered input and output**, whether media flows or anyone watches. Inputs created by forwarding a Fishjam room are not billed. Recordings are billed separately. Rates: <https://fishjam.swmansion.com/pricing>.
- **Unregister inputs and outputs you no longer need**; they bill for as long as they stay registered.
- **Starts automatically** by default. Create with `autostart: false` when inputs arrive later (room forwarding), then call start.
- **Idle cleanup:** by default a composition is removed after about 5 minutes in which none of its inputs carry media. A composition with no inputs at all counts as idle.
- **`cleanup_without_inputs: false`** makes cleanup wait for inputs *and* outputs to go silent. Needed for late-joining room peers or input-less templates (captions only). An output that keeps publishing is then never cleaned up for you.
- **Delete explicitly** when finished, even after a test. Deleting removes all inputs and outputs and finalizes recordings.

## Key rules

- **Destinations must exist before you register.** An output connects to its endpoint during registration, and an MP4 URL is fetched during registration. Unreachable targets fail the register call, not later.
- **One output per livestream room.** A Fishjam livestream accepts one streamer (`../platform/room-types.md`); a second output into the same room fails with 400.
- **`updateOutput` mirrors registration.** An output registered with video and audio needs both in every update; one registered with only video accepts only video.
- **Template outputs sound like what they render.** Give the output an `audio` config, or it has no audio track. The mix is every `<InputStream>` currently rendered; an input that is not rendered, such as a peer with no tile, is silent. Detail: `templates.md`.
- **Resolution** must be even on both sides.
- **Composed rooms use H.264.** Create them with `videoCodec: "h264"` explicitly, even though it is the default, so a changed default cannot break the composition.
- **Casing differs by API.** Composition API bodies are snake_case and closed (unknown fields return 422). Server API bodies are camelCase (`compositionURL`, `outputId`). The SDKs handle both.

## References

| File | When to read |
|---|---|
| `templates.md` | `composition-cli` scaffold and build, `@fishjam-cloud/composition` hooks, `eventBus`, deploy and redeploy. |
| `scenes.md` | Component tree, audio mix, updating static scenes live, transitions, scheduled changes. |
| `inputs-and-outputs.md` | Every input and output type, publishing into WHIP / RTMP inputs, reaching viewers. |
| `room-composition.md` | Composing a Fishjam room end to end with `forwardRoomTracks`. |
| `recording.md` | Recording an output, statuses, fetching files. |
| `rest-endpoints.md` | Raw Composition API over HTTP, without an SDK. |
