# Static Scenes

A scene is a description of what an output shows. With a static scene, your backend builds it as JSON, sends it when the output is registered, and replaces it with `updateOutput` whenever the layout should change. Use it for a fixed layout of inputs you know up front.

A template produces the same kind of scene: every time it renders, the output's scene is replaced with the result. The only difference is who builds the scene, your backend or the template. When the layout depends on who is in a room or reacts to events, use a template (`templates.md`).

## Scene shape

Each output has a **video scene** and an **audio scene**:

```ts
const video = { root: component };
const audio = { inputs: [{ inputId: 'host' }, { inputId: 'music', volume: 0.2 }] };
```

- The video scene is a tree with a single `root` component.
- The audio scene is independent of the video scene. **An input is heard only if it is listed in `inputs`**, whether or not its video appears in the video scene. `volume` defaults to `1`.
- The SDKs take camelCase (`inputId`, `borderRadius`); the REST API takes the same fields in snake_case (`input_id`, `border_radius`), see `rest-endpoints.md`.

## Components

| `type` | Draws | Key fields |
|---|---|---|
| `input_stream` | One input's video | `inputId` |
| `view` | A container; children laid out in a `row` or `column`, later children drawn on top | `children`, `direction`, `width`, `height`, `top` / `left` / `right` / `bottom` for absolute position, `backgroundColor`, `borderRadius` |
| `tiles` | Children in an automatic grid | `children`, `padding`, `margin`, `tileAspectRatio` |
| `rescaler` | One child scaled to fit its area, keeping aspect ratio | `child`, `mode: 'fit'` (letterbox, default) or `'fill'` (crop) |
| `text` | Text | `text`, `fontSize`, `fontFamily`, `color`, `width` |
| `image` | A registered image | `imageId`, `width` or `height` |

Colors are `#RRGGBBAA` strings. For every field and layout rule, see **Smelter reference** below.

## Example: picture in picture

A main input filling the frame, a second input in a rounded corner box, a logo, and the main input's audio ducked under the host's voice. Calling the same function with the inputs swapped, then `updateOutput`, swaps them with an animation.

```ts
import type { AudioScene, InputId, OutputId, RendererId, VideoScene } from '@fishjam-cloud/js-server-sdk';

function pictureInPicture(main: string, corner: string): VideoScene {
  return {
    root: {
      type: 'view',
      backgroundColor: '#000000ff',
      children: [
        { type: 'rescaler', child: { type: 'input_stream', inputId: main } },
        {
          type: 'rescaler',
          id: 'corner',
          mode: 'fill',
          top: 24,
          right: 24,
          width: 320,
          height: 180,
          borderRadius: 16,
          transition: { durationMs: 500 },
          child: { type: 'input_stream', inputId: corner },
        },
        {
          type: 'view',
          top: 24,
          left: 24,
          width: 120,
          height: 120,
          children: [{ type: 'image', imageId: 'logo', width: 120 }],
        },
      ],
    },
  };
}

await compositionClient.registerImage(compositionId, 'logo' as RendererId, {
  assetType: 'auto',
  url: 'https://fishjam.swmansion.com/docs/img/logo.svg',
});
await compositionClient.registerMp4Input(compositionId, 'race' as InputId, {
  url: 'https://smelter.dev/videos/template-scene-race.mp4',
  loop: true,
});
await compositionClient.registerMp4Input(compositionId, 'host' as InputId, {
  url: 'https://smelter.dev/videos/template-scene-streamer.mp4',
  loop: true,
});

const audio: AudioScene = { inputs: [{ inputId: 'race', volume: 0.2 }, { inputId: 'host' }] };

await compositionClient.registerWhipOutput(compositionId, 'main' as OutputId, {
  endpointUrl: fishjamClient.livestreamWhipUrl(),
  bearerToken: streamerToken,
  video: { resolution: { width: 1280, height: 720 }, initial: pictureInPicture('race', 'host') },
  audio: { initial: audio },
});

await compositionClient.updateOutput(compositionId, 'main' as OutputId, {
  video: pictureInPicture('host', 'race'),
  audio,
});
```

Python equivalent: `../python-server-sdk/composition.md`.

## Registering a scene

Scenes are first given at output registration, under `video.initial` and `audio.initial`, together with the output's `video.resolution`. Everything a scene references must already be registered: inputs by `inputId`, images by `imageId`, fonts by family name.

- **Images:** `registerImage(compositionId, imageId, { assetType, url })`. `assetType` is `'png'`, `'jpeg'`, `'svg'`, `'gif'`, or `'auto'` to detect it from the URL. Remove with `unregisterImage`.
- **Fonts:** `registerFont(compositionId, pathOrBlob)` uploads a font file. `text` components select it with `fontFamily` set to the family name stored in the file (for example `'Inter'`).

## Updating a scene

`updateOutput(compositionId, outputId, { video, audio })` replaces the output's scene while it streams.

- **Pass whole scenes.** An update is not a patch; send the complete `{ root }` and the complete audio `inputs`. There is no `initial` wrapper here.
- **Mirror the registration.** An output registered with both video and audio needs both in every update; one registered with only video accepts only video. A mismatch is rejected with 400.
- **Animate with `id` + `transition`.** When a `view` or `rescaler` keeps the same `id` between the old and new scene and has a `transition`, its size and position change smoothly over `durationMs`. See <https://smelter.dev/http-api/guides/transitions>.
- **Schedule instead of timing it yourself.** `scheduleTimeMs` applies the update at a moment measured in milliseconds from when the composition started, so several changes can be lined up in one go. `unregisterInput`, `unregisterOutput`, and `unregisterImage` accept the same option.
- **Template outputs are not updated this way.** A template owns its scene; see `templates.md`.

## Smelter reference

Scenes are rendered by Smelter. Look fields and behavior up here rather than guessing:

| Need | Source |
|---|---|
| Fields the Fishjam API accepts in a static scene | <https://fishjam.swmansion.com/docs/api/compositions> |
| What each component field does | <https://github.com/software-mansion/smelter/blob/fishjam-v2/smelter-api/src/video/component.rs> |
| Sizing, absolute positioning, how children share space | <https://smelter.dev/http-api/guides/basic-layouts> |
| Transitions | <https://smelter.dev/http-api/guides/transitions> |
