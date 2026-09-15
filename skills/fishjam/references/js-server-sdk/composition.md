# CompositionClient

Client for Fishjam compositions. Source: `packages/js-server-sdk/src/composition.ts` in the `js-server-sdk` repo (<https://github.com/fishjam-cloud/js-server-sdk>).

> Concepts, scenes, templates, and room composition are covered in `../composition/SKILL.md`. This file is the TypeScript surface only.

## Constructor

```ts
import { CompositionClient } from '@fishjam-cloud/js-server-sdk';

const compositionClient = new CompositionClient({
  managementToken: process.env.FISHJAM_MANAGEMENT_TOKEN!,
});
```

- Same package and same management token as `FishjamClient`; no Fishjam ID.
- Create one per process and share it.

## IDs

`CompositionId`, `InputId`, `OutputId`, and `RendererId` are branded strings. `createComposition` returns a `CompositionId`; IDs you choose yourself are cast once:

```ts
import type { InputId, OutputId, RendererId } from '@fishjam-cloud/js-server-sdk';

const camera = 'camera' as InputId;
const main = 'main' as OutputId;
const logo = 'logo' as RendererId;
```

Inside scenes, `inputId` and `imageId` are plain strings.

## Compositions

```ts
const { compositionId } = await compositionClient.createComposition();
await compositionClient.createComposition({ autostart: false, cleanupWithoutInputs: false });

await compositionClient.startComposition(compositionId);
await compositionClient.deleteComposition(compositionId);

const url: string = compositionClient.compositionUrl(compositionId);
```

`compositionUrl` is synchronous and makes no request; pass its result to `fishjamClient.forwardRoomTracks`.

## Inputs

```ts
const { url, bearerToken } = await compositionClient.registerWhipInput(compositionId, camera);
await compositionClient.registerWhipInput(compositionId, camera, { bearerToken: 'chosen-token', video: false });

const publishUrl: string = await compositionClient.registerRtmpInput(compositionId, encoder, { streamKey });

await compositionClient.registerWhepInput(compositionId, remote, { endpointUrl: whepUrl, bearerToken: viewerToken });

const { videoDurationMs, audioDurationMs } = await compositionClient.registerMp4Input(compositionId, intro, {
  url: 'https://smelter.dev/videos/template-scene-race.mp4',
  loop: true,
});

await compositionClient.unregisterInput(compositionId, camera);
await compositionClient.unregisterInput(compositionId, camera, { scheduleTimeMs: 10_000 });
```

- `registerWhipInput` returns the full publish URL, ready to hand to a publisher.
- `registerInput(compositionId, inputId, { type, ... })` is the generic form behind these helpers.

## Outputs

```ts
await compositionClient.registerWhipOutput(compositionId, main, {
  endpointUrl: fishjamClient.livestreamWhipUrl(),
  bearerToken: streamerToken,
  video: { resolution: { width: 1280, height: 720 }, initial: videoScene },
  audio: { initial: audioScene },
});

await compositionClient.registerRtmpOutput(compositionId, restream, {
  url: ingestUrlWithStreamKey,
  video: { resolution: { width: 1920, height: 1080 }, initial: videoScene },
  audio: { initial: audioScene },
});

await compositionClient.registerTemplateOutput(
  compositionId,
  'show' as OutputId,
  {
    type: 'whip_client',
    endpointUrl: fishjamClient.livestreamWhipUrl(),
    bearerToken: showStreamerToken,
    video: { resolution: { width: 1280, height: 720 }, initial: { root: { type: 'view' } } },
    audio: { initial: { inputs: [] } },
  },
  './my-template/dist/App.js',
);

await compositionClient.updateOutput(compositionId, main, { video: videoScene, audio: audioScene });
await compositionClient.requestKeyframe(compositionId, main);
await compositionClient.unregisterOutput(compositionId, main);
```

- Each output pushes to its own destination, so `show` uses a streamer token of a second livestream room. `updateOutput` applies to static outputs like `main`, never to a template output.
- `registerTemplateOutput` takes the config with `type` set, and the bundle as a file path or a `Blob`. A path is read from disk, so never pass one taken from user input.
- `registerOutput(compositionId, outputId, { type, ... })` is the generic form behind `registerWhipOutput` and `registerRtmpOutput`.

## Assets and events

```ts
await compositionClient.registerImage(compositionId, logo, { assetType: 'auto', url: 'https://fishjam.swmansion.com/docs/img/logo.svg' });
await compositionClient.unregisterImage(compositionId, logo);

await compositionClient.registerFont(compositionId, './fonts/Inter.ttf');

await compositionClient.sendEvent(compositionId, { eventName: 'SET_CAPTION', data: { text: 'Welcome!' } });
```

`registerFont` takes a file path or a `Blob`, like `registerTemplateOutput`.

## Related `FishjamClient` methods

| Method | Use |
|---|---|
| `forwardRoomTracks(roomId, compositionUrl)` | Forward a room into a composition (`../composition/room-composition.md`). |
| `livestreamWhipUrl()` | `endpointUrl` for a WHIP output into a Fishjam livestream room. |
| `createLivestreamStreamerToken(roomId)` | `bearerToken` for that output (`livestream-and-moq.md`). |
| `createLivestreamViewerToken(roomId)` | `bearerToken` for a WHEP input pulling a private Fishjam livestream. |

## Types

Scene and config types are exported as types only: `VideoScene`, `AudioScene`, `Component`, `View`, `Tiles`, `Rescaler`, `InputStream`, `Text`, `Image`, `Transition`, `RegisterOutput`, `WhipOutput`, `RtmpOutput`, `ImageSpec`, `CreateCompositionRequest`, `UpdateOutputRequest`, and more. Write enum-like fields as string literals (`type: 'rescaler'`, `mode: 'fill'`).

```ts
import type { AudioScene, VideoScene } from '@fishjam-cloud/js-server-sdk';
```

## Exception types

Composition calls throw the same `FishjamBaseException` hierarchy as `FishjamClient` (`client.md`), plus not-found types for composition resources:

| Exception | When |
|---|---|
| `CompositionNotFoundException` | The composition does not exist, for example after idle cleanup or deletion. |
| `InputNotFoundException`, `OutputNotFoundException`, `RendererNotFoundException` | The input, output, or image does not exist. |
| `BadRequestException` | Invalid body, such as a mismatched `updateOutput` or an odd resolution. `details` carries the server message. |
| `UnauthorizedException` | Wrong management token. |
| `ServiceUnavailableException` | Temporarily unavailable; retry with backoff. |

## Sources

- <https://fishjam.swmansion.com/docs/explanation/compositions>
- <https://fishjam.swmansion.com/docs/api/compositions>
