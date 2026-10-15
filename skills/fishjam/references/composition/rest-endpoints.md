# Composition REST API Endpoints

The HTTP surface of the Composition API. Both server SDKs wrap it in `CompositionClient`; prefer the SDK. Use this reference when calling compositions from a language without an SDK or debugging a wire problem.

Authoritative spec: <https://fishjam.swmansion.com/docs/api/composition-openapi.json> (rendered at <https://fishjam.swmansion.com/docs/api/compositions>).

## Auth

Every request requires the same management token as the Fishjam Server API (`../platform/auth-model.md`):

```http
Authorization: Bearer <management-token>
```

The one exception is the WHIP publish endpoint, which accepts only the token returned when its input was registered.

## Base URL

```
https://rtc.fishjam.io/api/composition
```

The Composition API is a different host from the Server API (`https://fishjam.io/api/v1/connect/<fishjam-id>`) and needs no Fishjam ID. A composition's URL, `https://rtc.fishjam.io/api/composition/<composition_id>`, is also what you pass to the Server API when forwarding a room.

## Conventions

- **Bodies are snake_case** (`input_id`, `endpoint_url`, `cleanup_without_inputs`). The Server API is camelCase; do not mix them.
- **Bodies are closed.** A field the endpoint does not define is rejected with a message naming it.
- **Success responses** are `{}` unless listed below. Errors carry `{ "message": "..." }`.
- **IDs in paths are yours**, except `composition_id`, which `POST /api/composition` returns.

## Compositions

| Method | Path | Purpose |
|---|---|---|
| `POST` | `/api/composition` | Create. Body: `{ autostart?, cleanup_without_inputs? }`, `{}` for defaults. Returns 201 `{ composition_id, api_url }`. |
| `POST` | `/api/composition/{composition_id}/start` | Start a composition created with `autostart: false`. |
| `DELETE` | `/api/composition/{composition_id}` | Delete, with all its inputs and outputs. |

## Inputs

| Method | Path | Purpose |
|---|---|---|
| `POST` | `/api/composition/{composition_id}/input/{input_id}/register` | Register an input. Body is one of the shapes below. |
| `POST` | `/api/composition/{composition_id}/input/{input_id}/unregister` | Remove an input. Optional body `{ schedule_time_ms }`. |

| `type` | Body | Response |
|---|---|---|
| `whip_server` | `{ type, bearer_token?, video? }` | `{ bearer_token, endpoint_route }`. Publish URL: composition URL + `endpoint_route`. |
| `rtmp_server` | `{ type, stream_key }` | `{ publish_url }` |
| `whep_client` | `{ type, endpoint_url, bearer_token?, video? }` | `{}` |
| `mp4` | `{ type, url, loop? }` | `{ video_duration_ms, audio_duration_ms }` |

## Outputs

| Method | Path | Purpose |
|---|---|---|
| `POST` | `/api/composition/{composition_id}/output/{output_id}/register` | Register an output with a static scene. Body is one of the shapes below. |
| `POST` | `/api/composition/{composition_id}/output/{output_id}/template` | Register a template output. Multipart: `config` (the same JSON as `register`, as `application/json`) and `template` (the built bundle). |
| `POST` | `/api/composition/{composition_id}/output/{output_id}/update` | Replace a static scene. Body: `{ video?, audio?, schedule_time_ms? }`, where `video` is `{ root }` and `audio` is `{ inputs }`. |
| `POST` | `/api/composition/{composition_id}/output/{output_id}/request_keyframe` | Ask the output for a fresh keyframe. |
| `POST` | `/api/composition/{composition_id}/output/{output_id}/unregister` | Remove an output. Optional body `{ schedule_time_ms }`. |

| `type` | Body |
|---|---|
| `whip_client` | `{ type, endpoint_url, bearer_token?, video?, audio? }` |
| `rtmp_client` | `{ type, url, video?, audio? }` |

`video` is `{ resolution: { width, height }, initial: { root } }`; `audio` is `{ initial: { inputs: [{ input_id, volume? }] } }`. Scene components: `scenes.md`.

## Assets and events

| Method | Path | Purpose |
|---|---|---|
| `POST` | `/api/composition/{composition_id}/image/{image_id}/register` | Register an image. Body: `{ asset_type, url }`, `asset_type` one of `png`, `jpeg`, `svg`, `gif`, `auto`. |
| `POST` | `/api/composition/{composition_id}/image/{image_id}/unregister` | Remove an image. Optional body `{ schedule_time_ms }`. |
| `POST` | `/api/composition/{composition_id}/font/register` | Upload a font. Multipart with a single `font` part. |
| `POST` | `/api/composition/{composition_id}/event` | Send an event to templates. Body: `{ event_name, data? }`. |

## Media

| Method | Path | Purpose |
|---|---|---|
| `POST` | `/api/composition/{composition_id}/whip/{input_id}` | WHIP publish endpoint of a `whip_server` input. Called by the publisher with the input's `bearer_token`, not by your backend. |

## Server API: forwarding a room

Forwarding lives on the Fishjam Server API, with its camelCase body:

| Method | Path | Purpose |
|---|---|---|
| `POST` | `https://fishjam.io/api/v1/connect/<fishjam-id>/room/{room_id}/track_forwardings` | Forward the room into a composition. Body: `{ "compositionURL": "https://rtc.fishjam.io/api/composition/<composition_id>" }`. See `room-composition.md`. |

## Examples

```bash
API=https://rtc.fishjam.io/api/composition
AUTH="Authorization: Bearer $FISHJAM_MANAGEMENT_TOKEN"

curl -X POST "$API" -H "$AUTH" -H 'content-type: application/json' -d '{}'

curl -X POST "$API/$COMPOSITION/input/movie/register" -H "$AUTH" -H 'content-type: application/json' \
  -d '{"type":"mp4","url":"https://smelter.dev/videos/template-scene-race.mp4","loop":true}'

curl -X POST "$API/$COMPOSITION/font/register" -H "$AUTH" -F "font=@Inter.ttf"

curl -X POST "$API/$COMPOSITION/output/main/template" -H "$AUTH" \
  -F "config=$OUTPUT_CONFIG_JSON;type=application/json" \
  -F "template=@my-template/dist/App.js"

curl -X POST "$API/$COMPOSITION/event" -H "$AUTH" -H 'content-type: application/json' \
  -d '{"event_name":"SET_CAPTION","data":{"text":"Welcome!"}}'

curl -X DELETE "$API/$COMPOSITION" -H "$AUTH"
```
