# Composition Templates

A template is a React component that renders one output's scene. It runs inside the composition, next to the renderer, and re-renders when room peers change or your backend sends an event. The backend only deploys it and sends events; it never rebuilds scenes.

Templates are always TypeScript and React. A Python backend deploys the same built bundle as a Node backend.

## Scaffold and build

```bash
npx @fishjam-cloud/composition-cli init my-template
cd my-template
npm install
npm run typecheck
npm run build
```

`npm run build` bundles `src/App.tsx` into `dist/App.js`, the file you deploy. `init` creates `package.json`, `tsconfig.json`, and `src/App.tsx` with matching versions of these packages:

| Package | Provides |
|---|---|
| `@swmansion/smelter` | Layout components (`View`, `Tiles`, `Rescaler`, `InputStream`, `Text`, `Image`) and hooks (`useInputStreams`, `useAudioInput`). |
| `@fishjam-cloud/composition` | Room hooks (`usePeers`, `usePeer`, `useRoom`, `useSpeakingState`) and `eventBus`. |
| `react` | React itself. |

To build a different entry file: `npx composition-cli build src/Other.tsx --out dist/Other.js`.

What makes a valid bundle:

- **Default-export a React component.** That component is the output's root.
- **`react`, `@swmansion/smelter`, and `@fishjam-cloud/composition` are provided by the composition.** Everything else you import is bundled into the file; `composition-cli build` validates the result and fails with a list of problems.
- **Outside data comes from your backend.** A template gets data only from `sendEvent`, the room hooks, and images and fonts registered on the composition. Do not plan on fetching anything from the template.

## Layout components

Templates build layouts from `View`, `Tiles`, `Rescaler`, `InputStream`, `Text`, and `Image`, the React form of the components in `scenes.md`. Look props up rather than guessing:

- Component props and hooks, as Fishjam runs them: <https://github.com/software-mansion/smelter/tree/fishjam-v2/ts/smelter/src> (`components/`, `hooks.ts`)
- Layout rules and transitions: <https://smelter.dev/ts-sdk/guides/basic-layouts>, <https://smelter.dev/ts-sdk/guides/transitions>

Fishjam specifics:

- `<InputStream inputId>` takes an input ID: one you registered, or `stream.inputId` from the room hooks.
- `<Image imageId>` takes an image registered with `registerImage` / `register_image`.
- `<Text style={{ fontFamily }}>` names a font you uploaded with `registerFont` / `register_font`, by the family name stored in the font file (for example `'Inter'`). The upload takes no name. Register a font whenever the template shows text.
- For inputs you registered yourself (not from a room), `useInputStreams()` lists them with their state, for example to show only inputs whose `videoState` is `'playing'`.

## Room hooks

Available once a Fishjam room is forwarded into the composition (`room-composition.md`).

| Hook | Returns |
|---|---|
| `usePeers<PeerMetadata, ServerMetadata>()` | Every peer in the room as `PeerWithStreams[]`. |
| `usePeer<PeerMetadata, ServerMetadata>(peerId)` | One peer, or `undefined`. |
| `useRoom()` | `{ id }` of the forwarded room, or `undefined` before a room is forwarded. |
| `useSpeakingState(peerId)` | `'speech'` or `'silence'` for that peer's audio. |

```ts
type PeerWithStreams<PeerMetadata, ServerMetadata> = {
  id: string;
  metadata: { peer: PeerMetadata; server: ServerMetadata };
  streams: Stream[];
  cameraStream?: Stream;
  screenShareStream?: Stream;
  customStreams: Stream[];
};

type Stream = {
  inputId: string;
  video?: { type: 'video'; id: string; paused: boolean; metadata: Record<string, unknown> };
  audio?: { type: 'audio'; id: string; paused: boolean; metadata: Record<string, unknown> };
};
```

- `metadata.peer` is what the client set when joining; `metadata.server` is what your backend set in `createPeer` (`../js-server-sdk/client.md`).
- `cameraStream` holds the camera and microphone, `screenShareStream` the screen share with its audio, `customStreams` everything else.
- `paused` is `true` while the peer has that track muted, for example camera turned off.
- A peer is listed as soon as it is in the room; its `streams` fill in once its media reaches the composition.
- These hooks are unrelated to `usePeers` in `@fishjam-cloud/react-client`; a template never imports client SDK packages.

## Audio

A template output's audio mix is built from the template, not from the output config:

- **Register the output with an `audio` config**, even an empty one (`audio: { initial: { inputs: [] } }`). Without it the output has no audio track.
- **Every rendered `<InputStream>` is mixed in**, at its `volume` (default 1) or silent when `muted`.
- **`useAudioInput(inputId, { volume })` adds an input to the mix without drawing it.**
- **An input that is neither rendered nor passed to `useAudioInput` is silent.** A tile that shows a name instead of video, because the camera is off, would otherwise silence that peer's microphone.

The robust pattern: draw video with `muted`, and add every stream's audio with `useAudioInput`, independent of what is on screen. The example below does this.

## Example: room grid with captions

A grid of peers with a speaking highlight, a name when the camera is off, audio for everyone, and a caption bar controlled by the backend.

```tsx
import { InputStream, Rescaler, Text, Tiles, View, useAudioInput } from '@swmansion/smelter';
import { eventBus, usePeers, useRoom, useSpeakingState } from '@fishjam-cloud/composition';
import type { PeerWithStreams } from '@fishjam-cloud/composition';
import { useEffect, useState } from 'react';

type PeerMetadata = { displayName?: string };

export default function App() {
  const room = useRoom();
  const peers = usePeers<PeerMetadata>().filter((peer) => peer.streams.length > 0);
  const caption = useCaption();

  return (
    <View style={{ backgroundColor: '#0b1020ff', direction: 'column' }}>
      {peers.length > 0 ? (
        <Tiles style={{ padding: 16 }}>
          {peers.map((peer) => (
            <PeerTile key={peer.id} peer={peer} />
          ))}
        </Tiles>
      ) : (
        <Label text={room ? 'Waiting for participants' : 'Starting'} />
      )}
      {caption && (
        <View style={{ height: 96, backgroundColor: '#000000cc' }}>
          <Label text={caption} />
        </View>
      )}
    </View>
  );
}

function PeerTile({ peer }: { peer: PeerWithStreams<PeerMetadata> }) {
  const speaking = useSpeakingState(peer.id) === 'speech';
  const camera = peer.cameraStream;

  return (
    <View style={{ borderWidth: 4, borderColor: speaking ? '#00cc66ff' : '#00000000' }}>
      {camera?.video && !camera.video.paused ? (
        <Rescaler style={{ rescaleMode: 'fill' }}>
          <InputStream inputId={camera.inputId} muted />
        </Rescaler>
      ) : (
        <Label text={peer.metadata.peer?.displayName ?? peer.id} />
      )}
      {peer.streams
        .filter((stream) => stream.audio)
        .map((stream) => (
          <PeerAudio key={stream.inputId} inputId={stream.inputId} />
        ))}
    </View>
  );
}

function PeerAudio({ inputId }: { inputId: string }) {
  useAudioInput(inputId, { volume: 1 });
  return null;
}

function Label({ text }: { text: string }) {
  return <Text style={{ fontFamily: 'Inter', fontSize: 40, color: '#ffffffff' }}>{text}</Text>;
}

function useCaption() {
  const [caption, setCaption] = useState<string | null>(null);
  useEffect(() => eventBus.on<{ text: string | null }>('SET_CAPTION', ({ text }) => setCaption(text)), []);
  return caption;
}
```

## Events

Your backend sends named events with any JSON payload; templates subscribe by name.

| Backend | Call |
|---|---|
| JS | `compositionClient.sendEvent(compositionId, { eventName: 'SET_CAPTION', data: { text: 'Welcome!' } })` |
| Python | `composition_client.send_event(composition_id, "SET_CAPTION", {"text": "Welcome!"})` |

In the template, subscribe with `eventBus.on<Payload>('SET_CAPTION', handler)` inside `useEffect`, and return the unsubscribe function it gives back (see `useCaption` in the example).

- **Event names and payloads are yours.** There is no built-in set.
- **Every template in the composition that listens to that name receives the event.**
- **Events are not stored.** One sent while no template is listening, for example before the template output is registered or while it is being redeployed, is gone.
- **Send full state, not deltas.** Prefer `SET_CAPTION { text }` over `APPEND_CAPTION`, so any single event leaves the template correct, and send the current state again after deploying or redeploying a template.

## Deploy and redeploy

Deploy with one call that uploads the bundle together with the output config:

```ts
import type { OutputId } from '@fishjam-cloud/js-server-sdk';

await compositionClient.registerFont(compositionId, './fonts/Inter.ttf');
await compositionClient.registerTemplateOutput(
  compositionId,
  'main' as OutputId,
  {
    type: 'whip_client',
    endpointUrl: fishjamClient.livestreamWhipUrl(),
    bearerToken: streamerToken,
    video: { resolution: { width: 1280, height: 720 }, initial: { root: { type: 'view' } } },
    audio: { initial: { inputs: [] } },
  },
  './my-template/dist/App.js',
);
```

Python: `composition_client.register_template_output(composition_id, "main", WhipOutput(type_=WhipOutputType.WHIP_CLIENT, ...), "./my-template/dist/App.js")`; see `../python-server-sdk/composition.md`.

- **The config is the same as for a regular output**, with `type` set explicitly (`'whip_client'` or `'rtmp_client'`). `video.initial` and `audio.initial` are required by the schema; pass the placeholders above, the template replaces them.
- **Register the template output directly.** Do not register a plain output with the same ID first.
- **The template argument** is a file path or a `Blob` (Python: path, `str`, or `bytes`). Never pass a path that comes from user input; the file is uploaded.
- **Redeploy** by unregistering the output, then registering the template output again. Inputs stay registered. Resend current state events afterwards.
- **Do not call `updateOutput` on a template output.** The template owns the scene and its next render replaces your update.

For composing a room, including the order of create, deploy, forward, and start: `room-composition.md`.

## Tips

- **Render something for every state:** no room forwarded yet, no peers, camera off. An empty or black frame looks broken to viewers.
- **Filter out peers without streams** (`peer.streams.length > 0`) so tiles do not appear before there is media to show.
- **Call `useSpeakingState` inside the tile component**, not in the parent, so a speaking change re-renders one tile instead of the whole layout.
- **Animate with transitions** rather than rapid state changes. Every change the template renders is sent as a scene update, so a timer changing state many times per second is wasteful; a transition gets smooth motion from a single update.
