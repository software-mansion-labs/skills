# Inputs and Outputs

Inputs bring streams into a composition; outputs push the composed result out. Both are addressed by IDs you choose, unique within the composition. Every input and output you register is billed per minute while registered, so unregister the ones you no longer need. Inputs Fishjam creates when forwarding a room are not billed (`room-composition.md`).

## Inputs

| Type | Who connects | JS | Python | Returns |
|---|---|---|---|---|
| WHIP | A WHIP publisher pushes to the composition | `registerWhipInput` | `register_whip_input` | Publish `url` and `bearerToken` |
| RTMP | An RTMP encoder pushes to the composition | `registerRtmpInput` | `register_rtmp_input` | Publish URL |
| WHEP | The composition pulls from a WHEP server | `registerWhepInput` | `register_whep_input` | Nothing |
| MP4 | The composition downloads a file by URL | `registerMp4Input` | `register_mp4_input` | Video and audio duration |
| Room tracks | Fishjam forwards every peer's tracks | `fishjamClient.forwardRoomTracks` | `fishjam_client.forward_room_tracks` | Nothing; inputs appear per peer (`room-composition.md`) |

Remove any input with `unregisterInput(compositionId, inputId)` / `unregister_input`.

### WHIP input

```ts
const { url, bearerToken } = await compositionClient.registerWhipInput(compositionId, 'camera' as InputId);
```

- **Hand `url` and `bearerToken` to the publisher**; any WHIP client can publish with them.
- **Only the returned token works on that URL.** The management token is rejected there, which makes the pair safe to give to a streaming tool. Pass your own `bearerToken` when registering if you want to choose it.
- **One publisher at a time.** A second publisher is refused while the first is connected; it can connect after the first leaves.
- **Video is H.264.** Register with `video: false` for an audio-only input.

### RTMP input

```ts
const publishUrl = await compositionClient.registerRtmpInput(compositionId, 'encoder' as InputId, {
  streamKey: 'my-secret-key',
});
```

- **Hand `publishUrl` to the encoder** and use it exactly as returned; do not build it yourself.
- **RTMPS is highly recommended.** It is encrypted; plain RTMP sends the stream key and the media unencrypted. Use an encoder that supports RTMPS.

### WHEP input

```ts
await compositionClient.registerWhepInput(compositionId, 'remote' as InputId, {
  endpointUrl: 'https://fishjam.io/api/v1/live/api/whep',
  bearerToken: viewerToken,
});
```

- **Pulls from any WHEP server**, including another Fishjam livestream: use the livestream WHEP URL with a viewer token from `createLivestreamViewerToken` (`../js-server-sdk/livestream-and-moq.md`). Python has `fishjam_client.livestream_whep_url()`.
- **Connects during registration.** The source must already be streaming.
- **If the source stops, register the input again** once it is back.

### MP4 input

```ts
const { videoDurationMs, audioDurationMs } = await compositionClient.registerMp4Input(compositionId, 'intro' as InputId, {
  url: 'https://smelter.dev/videos/template-scene-race.mp4',
  loop: true,
});
```

- **The file is downloaded during registration**, so an unreachable URL fails the call.
- **Without `loop`, the input finishes at the end of the file.** A finished input carries no media, which counts towards idle cleanup (`SKILL.md`, Lifecycle and billing).

## Outputs

| Type | Pushes to | JS | Python |
|---|---|---|---|
| WHIP | A Fishjam livestream room, or any WHIP server | `registerWhipOutput` | `register_whip_output` |
| RTMP | YouTube, Twitch, or any RTMP / RTMPS ingest | `registerRtmpOutput` | `register_rtmp_output` |
| Template | Either of the above, rendered by a template | `registerTemplateOutput` | `register_template_output` |

Every output takes:

- `video`: `{ resolution: { width, height }, initial }`, with both sides even.
- `audio`: `{ initial }`, optionally `channels: 'mono' | 'stereo'`.
- At least one of the two. Registering only `video` gives an output without an audio track.

Scenes (`initial`) are described in `scenes.md`; templates in `templates.md`.

- **The destination must accept the connection during registration.** Create the livestream room or start the platform's stream first.
- **One composition can have several outputs**, each with its own scene and destination: for example the same show to your app and to YouTube.
- **End an output when its inputs end** with `video.sendEosWhen` / `audio.sendEosWhen` (`anyOf` / `allOf` a list of input IDs, or `anyInput` / `allInputs`); useful when an output shows a file that should not loop.
- **`requestKeyframe(compositionId, outputId)`** asks the output for a fresh keyframe, for a destination that needs one to start decoding.
- **Unregister** with `unregisterOutput(compositionId, outputId)`.

### To your app, through a Fishjam livestream

```ts
const livestream = await fishjamClient.createRoom({ roomType: 'livestream' });
const { token } = await fishjamClient.createLivestreamStreamerToken(livestream.id);

await compositionClient.registerWhipOutput(compositionId, 'main' as OutputId, {
  endpointUrl: fishjamClient.livestreamWhipUrl(),
  bearerToken: token,
  video: { resolution: { width: 1280, height: 720 }, initial: scene },
  audio: { initial: audioScene },
});
```

- The composition is the livestream's streamer; **a livestream room takes one output**. Room types and public vs private livestreams: `../platform/room-types.md`.
- Viewers join `livestream.id` in your app with `useLivestreamViewer` (`../react-client/livestream.md`), using a viewer token for a private livestream.

### To YouTube, Twitch, or another platform

```ts
await compositionClient.registerRtmpOutput(compositionId, 'restream' as OutputId, {
  url: ingestUrlWithStreamKey,
  video: { resolution: { width: 1920, height: 1080 }, initial: scene },
  audio: { initial: audioScene },
});
```

- **`url` is the full ingest URL including the stream key**, as the platform provides it. Both `rtmp://` and `rtmps://` work; **use `rtmps://` whenever the platform offers it**, since it is encrypted and plain RTMP sends the stream key unencrypted.
