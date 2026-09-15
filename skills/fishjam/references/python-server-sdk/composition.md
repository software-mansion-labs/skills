# CompositionClient (Python)

Client for Fishjam compositions. **Synchronous**, like `FishjamClient`. Source: `fishjam/api/_composition_client.py` in `fishjam-server-sdk`.

> Concepts, scenes, templates, and room composition are covered in `../composition/SKILL.md`. This file is the Python surface only. Templates themselves are always TypeScript; Python deploys the built bundle.

## Constructor

```python
import os

from fishjam import CompositionClient

composition_client = CompositionClient(management_token=os.environ["FISHJAM_MANAGEMENT_TOKEN"])
```

- Same package and same management token as `FishjamClient`; no Fishjam ID.
- Create one per process and share it. From async code, call it the same way as `FishjamClient` (`client.md`).

## Minimal flow

A looping MP4 composed into a Fishjam livestream room, the Python version of the flow in `../composition/SKILL.md`.

```python
import os

from fishjam import CompositionClient, FishjamClient, RoomOptions
from fishjam.composition import (
    AudioScene,
    AudioSceneInput,
    InputStream,
    InputStreamType,
    OutputWhipAudioOptions,
    OutputWhipVideoOptions,
    Rescaler,
    RescalerType,
    Resolution,
    VideoScene,
)

management_token = os.environ["FISHJAM_MANAGEMENT_TOKEN"]
fishjam_client = FishjamClient(fishjam_id=os.environ["FISHJAM_ID"], management_token=management_token)
composition_client = CompositionClient(management_token=management_token)

livestream = fishjam_client.create_room(RoomOptions(room_type="livestream"))
streamer_token = fishjam_client.create_livestream_streamer_token(livestream.id)

composition_id = composition_client.create_composition().composition_id
composition_client.register_mp4_input(
    composition_id, "movie", url="https://smelter.dev/videos/template-scene-race.mp4", loop=True
)
composition_client.register_whip_output(
    composition_id,
    "main",
    endpoint_url=fishjam_client.livestream_whip_url(),
    bearer_token=streamer_token,
    video=OutputWhipVideoOptions(
        resolution=Resolution(width=1280, height=720),
        initial=VideoScene(
            root=Rescaler(
                type_=RescalerType.RESCALER,
                child=InputStream(type_=InputStreamType.INPUT_STREAM, input_id="movie"),
            )
        ),
    ),
    audio=OutputWhipAudioOptions(initial=AudioScene(inputs=[AudioSceneInput(input_id="movie")])),
)

composition_client.delete_composition(composition_id)
fishjam_client.delete_room(livestream.id)
```

## Scenes in Python

Scene and config models live in `fishjam.composition`. They mirror the JSON in `../composition/scenes.md` with snake_case fields, plus two Python specifics:

- **Every component and output takes `type_`** with its enum: `View(type_=ViewType.VIEW, ...)`, `Tiles(type_=TilesType.TILES, ...)`, `Rescaler(type_=RescalerType.RESCALER, ...)`, `InputStream(type_=InputStreamType.INPUT_STREAM, ...)`, `Text(type_=TextType.TEXT, ...)`, `Image(type_=ImageType.IMAGE, ...)`, `WhipOutput(type_=WhipOutputType.WHIP_CLIENT, ...)`, `RtmpOutput(type_=RtmpOutputType.RTMP_CLIENT, ...)`.
- **Other enum fields use enums too**: `mode=RescaleMode.FILL`, `asset_type=ImageSpecAutoAssetType.AUTO`.

```python
from fishjam.composition import InputStream, InputStreamType, RescaleMode, Rescaler, RescalerType, VideoScene, View, ViewType


def picture_in_picture(main: str, corner: str) -> VideoScene:
    return VideoScene(
        root=View(
            type_=ViewType.VIEW,
            children=[
                Rescaler(type_=RescalerType.RESCALER, child=InputStream(type_=InputStreamType.INPUT_STREAM, input_id=main)),
                Rescaler(
                    type_=RescalerType.RESCALER,
                    mode=RescaleMode.FILL,
                    top=24,
                    right=24,
                    width=320,
                    height=180,
                    child=InputStream(type_=InputStreamType.INPUT_STREAM, input_id=corner),
                ),
            ],
        )
    )
```

## Compositions

```python
from fishjam.composition import CreateCompositionRequest

composition_id = composition_client.create_composition().composition_id
composition_client.create_composition(CreateCompositionRequest(autostart=False, cleanup_without_inputs=False))

composition_client.start_composition(composition_id)
composition_client.delete_composition(composition_id)

url = composition_client.composition_url(composition_id)
```

`composition_url` makes no request; pass its result to `fishjam_client.forward_room_tracks`.

## Inputs

```python
whip = composition_client.register_whip_input(composition_id, "camera")
composition_client.register_whip_input(composition_id, "camera", bearer_token="chosen-token", video=False)

publish_url = composition_client.register_rtmp_input(composition_id, "encoder", stream_key="my-secret-key")

composition_client.register_whep_input(
    composition_id, "remote", endpoint_url=fishjam_client.livestream_whep_url(), bearer_token=viewer_token
)

durations = composition_client.register_mp4_input(
    composition_id, "intro", url="https://smelter.dev/videos/template-scene-race.mp4", loop=True
)

composition_client.unregister_input(composition_id, "camera")
```

- `register_whip_input` returns `WhipInputTarget(url, bearer_token)` with the full publish URL.
- `register_mp4_input` returns `Mp4InputDurations(video_duration_ms, audio_duration_ms)`.
- `register_input(composition_id, input_id, input_)` is the generic form taking `WhipInput`, `RtmpInput`, `WhepInput`, or `Mp4Input`.
- Scheduled removal: `unregister_input(composition_id, input_id, UnregisterInput(schedule_time_ms=10_000))`; the same pattern applies to outputs and images.

## Outputs

```python
from fishjam.composition import (
    AudioScene,
    OutputRtmpClientAudioOptions,
    OutputRtmpClientVideoOptions,
    OutputWhipAudioOptions,
    OutputWhipVideoOptions,
    Resolution,
    VideoScene,
    View,
    ViewType,
    WhipOutput,
    WhipOutputType,
)

composition_client.register_whip_output(
    composition_id,
    "main",
    endpoint_url=fishjam_client.livestream_whip_url(),
    bearer_token=streamer_token,
    video=OutputWhipVideoOptions(resolution=Resolution(width=1280, height=720), initial=picture_in_picture("race", "host")),
    audio=OutputWhipAudioOptions(initial=audio_scene),
)

composition_client.register_rtmp_output(
    composition_id,
    "restream",
    url=ingest_url_with_stream_key,
    video=OutputRtmpClientVideoOptions(resolution=Resolution(width=1920, height=1080), initial=video_scene),
    audio=OutputRtmpClientAudioOptions(initial=audio_scene),
)

composition_client.register_template_output(
    composition_id,
    "show",
    WhipOutput(
        type_=WhipOutputType.WHIP_CLIENT,
        endpoint_url=fishjam_client.livestream_whip_url(),
        bearer_token=show_streamer_token,
        video=OutputWhipVideoOptions(resolution=Resolution(width=1280, height=720), initial=VideoScene(root=View(type_=ViewType.VIEW))),
        audio=OutputWhipAudioOptions(initial=AudioScene(inputs=[])),
    ),
    "./my-template/dist/App.js",
)

composition_client.update_output(composition_id, "main", video=picture_in_picture("host", "race"), audio=audio_scene)
composition_client.request_keyframe(composition_id, "main")
composition_client.unregister_output(composition_id, "main")
```

- Each output pushes to its own destination, so `show` uses a streamer token of a second livestream room. `update_output` applies to static outputs like `main`, never to a template output.
- RTMP outputs use `OutputRtmpClientVideoOptions` / `OutputRtmpClientAudioOptions`; WHIP outputs use the `OutputWhip...` pair.
- `register_template_output` takes the template as a path (`str` or `Path`) or `bytes`. A path is read from disk, so never pass one taken from user input.
- `update_output` also accepts an `UpdateOutputRequest(video=..., audio=..., schedule_time_ms=...)` instead of keywords.
- `register_output(composition_id, output_id, output)` is the generic form taking `WhipOutput` or `RtmpOutput`.

## Assets and events

```python
from fishjam.composition import ImageSpecAuto, ImageSpecAutoAssetType

composition_client.register_image(
    composition_id,
    "logo",
    ImageSpecAuto(asset_type=ImageSpecAutoAssetType.AUTO, url="https://fishjam.swmansion.com/docs/img/logo.svg"),
)
composition_client.unregister_image(composition_id, "logo")

composition_client.register_font(composition_id, "./fonts/Inter.ttf")

composition_client.send_event(composition_id, "SET_CAPTION", {"text": "Welcome!"})
```

`register_font` takes a path or `bytes`, like `register_template_output`.

## Related `FishjamClient` methods

| Method | Use |
|---|---|
| `forward_room_tracks(room_id, composition_url)` | Forward a room into a composition (`../composition/room-composition.md`). |
| `livestream_whip_url()` | `endpoint_url` for a WHIP output into a Fishjam livestream room. |
| `create_livestream_streamer_token(room_id)` | `bearer_token` for that output (`livestream-and-moq.md`). |
| `livestream_whep_url()` + `create_livestream_viewer_token(room_id)` | `endpoint_url` and `bearer_token` for a WHEP input pulling a private Fishjam livestream. |

## Exception types

Composition calls raise the same `fishjam.errors.HTTPError` hierarchy as `FishjamClient` (`client.md`), plus not-found types for composition resources:

| Class | When |
|---|---|
| `CompositionNotFoundError` | The composition does not exist, for example after idle cleanup or deletion. |
| `InputNotFoundError`, `OutputNotFoundError`, `RendererNotFoundError` | The input, output, or image does not exist. All subclass `NotFoundError`. |
| `BadRequestError` | Invalid body, such as a mismatched `update_output` or an odd resolution. |
| `UnauthorizedError` | Wrong management token. |
| `ServiceUnavailableError` | Temporarily unavailable; retry with backoff. |

## Sources

- <https://fishjam.swmansion.com/docs/explanation/compositions>
- <https://fishjam.swmansion.com/docs/api/compositions>
