# Composing a Fishjam Room

Forward a Fishjam room into a composition and every peer's media becomes composition inputs automatically. A template then lays the peers out with the room hooks (`templates.md`), and an output pushes the result to viewers or a platform (`inputs-and-outputs.md`). Peers keep using the room normally through the client SDKs; nothing changes on their side.

## How forwarding works

```ts
await fishjamClient.forwardRoomTracks(roomId, compositionClient.compositionUrl(compositionId));
```

Python: `fishjam_client.forward_room_tracks(room.id, composition_client.composition_url(composition_id))`.

- **One call covers the whole room.** Fishjam registers an input for each peer's media as it starts publishing and removes it when it stops. Never register room inputs yourself.
- **A peer's camera and microphone arrive as one input**, and a screen share with its audio as another. Templates see them as `cameraStream` and `screenShareStream` (`templates.md`).
- **Forwarded inputs are not billed** as composition inputs; the room is billed as usual.
- **A room and a composition pair one to one.** Create one composition per room you want to compose, and put everything that room needs (for example several outputs) into that composition. Calling `forwardRoomTracks` again with the same pair is harmless.
- **Forwarding lasts as long as the room.**
- **The room must use H.264.** Create it with `videoCodec: 'h264'`.

## Order of calls

1. **Create the room** with `roomType: 'conference'` and `videoCodec: 'h264'`, or use an existing H.264 room.
2. **Create the composition** with `autostart: false` and `cleanupWithoutInputs: false`.
3. **Register assets** the template uses: fonts and images.
4. **Create the destination and register the template output**, with an `audio` config so peers are heard.
5. **Forward the room.**
6. **Start the composition.**
7. **Send the template's current state** as events, if it uses any.

Why these options:

- **`autostart: false`**: nothing is rendered or sent until `startComposition`, so the output, assets, and forwarding are all in place before the first frame.
- **`cleanupWithoutInputs: false`**: a room can sit empty, or peers can join late, for longer than the idle cleanup window. Without it the composition could be removed before anyone publishes. The trade-off is that it is never cleaned up for you: delete it when the session ends.

## Example

A conference room composed with a template and streamed to a Fishjam livestream room.

```ts
import { CompositionClient, FishjamClient } from '@fishjam-cloud/js-server-sdk';
import type { OutputId } from '@fishjam-cloud/js-server-sdk';

const managementToken = process.env.FISHJAM_MANAGEMENT_TOKEN!;
const fishjamClient = new FishjamClient({ fishjamId: process.env.FISHJAM_ID!, managementToken });
const compositionClient = new CompositionClient({ managementToken });

const room = await fishjamClient.createRoom({ roomType: 'conference', videoCodec: 'h264' });
const livestream = await fishjamClient.createRoom({ roomType: 'livestream' });
const { token } = await fishjamClient.createLivestreamStreamerToken(livestream.id);

const { compositionId } = await compositionClient.createComposition({
  autostart: false,
  cleanupWithoutInputs: false,
});
await compositionClient.registerFont(compositionId, './fonts/Inter.ttf');
await compositionClient.registerTemplateOutput(
  compositionId,
  'main' as OutputId,
  {
    type: 'whip_client',
    endpointUrl: fishjamClient.livestreamWhipUrl(),
    bearerToken: token,
    video: { resolution: { width: 1280, height: 720 }, initial: { root: { type: 'view' } } },
    audio: { initial: { inputs: [] } },
  },
  './my-template/dist/App.js',
);
await fishjamClient.forwardRoomTracks(room.id, compositionClient.compositionUrl(compositionId));
await compositionClient.startComposition(compositionId);
await compositionClient.sendEvent(compositionId, { eventName: 'SET_CAPTION', data: { text: 'Welcome!' } });
```

From here, peers join `room.id` with peer tokens from `createPeer` as in any Fishjam app (`../platform/SKILL.md`), and viewers watch `livestream.id` (`inputs-and-outputs.md`). The template from `templates.md` matches this example.

When the session ends:

```ts
await compositionClient.deleteComposition(compositionId);
await fishjamClient.deleteRoom(livestream.id);
await fishjamClient.deleteRoom(room.id);
```

Python equivalent: `../python-server-sdk/composition.md`.
