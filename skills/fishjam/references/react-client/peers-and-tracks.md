# Peers & Tracks

`usePeers` is the read-side of the SDK — it gives you everyone in the room and what they're publishing.

Source: `@fishjam-cloud/react-client` → `packages/react-client/src/hooks/usePeers.ts`.

## Basic usage

```tsx
import { usePeers } from '@fishjam-cloud/react-client';

const { localPeer, remotePeers, setReceivedTracksQuality } = usePeers();
```

### `PeerWithTracks` shape

Each peer is a `PeerWithTracks` object:

```ts
type PeerWithTracks<PeerMetadata, ServerMetadata, T extends Track = Track> = {
  id: PeerId;
  metadata?: Metadata<PeerMetadata, ServerMetadata>;
  tracks: T[];
  cameraTrack?: T;
  microphoneTrack?: T;
  screenShareVideoTrack?: T;
  screenShareAudioTrack?: T;
  customVideoTracks: T[];
  customAudioTracks: T[];
};
```

The hook return-type uses two specializations of `PeerWithTracks`:
- `localPeer: PeerWithTracks<PeerMetadata, ServerMetadata> | null` (T = Track)
- `remotePeers: PeerWithTracks<PeerMetadata, ServerMetadata, RemoteTrack>[]`

The SDK pre-bins tracks by type — you don't have to filter `tracks` yourself. `cameraTrack`/`microphoneTrack` etc. are *the* camera/mic/screen share track of that peer, if present.

### `Track` vs `RemoteTrack`

- `localPeer.tracks` are `Track`. Local quality is controlled by `videoConfig` / `sentQualities` on `FishjamProvider` (you decide what *you* publish), not per-track from the consumer side.
- `remotePeers[*].tracks` are `RemoteTrack` — same fields, plus `setReceivedQuality(quality: Variant)` which calls `setTargetTrackEncoding` to request a specific simulcast layer from the SFU for *that one remote track*.

Both shapes:

```ts
type Track = {
  stream: MediaStream | null;
  trackId: TrackId;
  metadata?: TrackMetadata;
  simulcastConfig: SimulcastConfig | null;
  track: MediaStreamTrack | null;
};

type RemoteTrack = Track & {
  setReceivedQuality: (quality: Variant) => void;
};
```

`metadata` is **optional** — always optional-chain (`track.metadata?.type === 'camera'`); the SDK's own `usePeers` implementation does the same.

## Rendering a video grid

```tsx
function VideoGrid() {
  const { localPeer, remotePeers } = usePeers();

  return (
    <div className="grid">
      {localPeer && <PeerTile peer={localPeer} isLocal />}
      {remotePeers.map((peer) => <PeerTile key={peer.id} peer={peer} />)}
    </div>
  );
}

function PeerTile({ peer, isLocal }: { peer: PeerWithTracks<any, any>; isLocal?: boolean }) {
  return (
    <div>
      <div>{peer.metadata?.peer?.displayName ?? peer.id}</div>
      {peer.cameraTrack && (
        <video
          autoPlay
          muted={isLocal}                                    // mute your own video to avoid feedback
          playsInline
          ref={(el) => { if (el) el.srcObject = peer.cameraTrack!.stream; }}
        />
      )}
      {peer.microphoneTrack && !isLocal && (
        <audio
          autoPlay
          ref={(el) => { if (el) el.srcObject = peer.microphoneTrack!.stream; }}
        />
      )}
      {peer.screenShareVideoTrack && (
        <video
          autoPlay
          playsInline
          ref={(el) => { if (el) el.srcObject = peer.screenShareVideoTrack!.stream; }}
        />
      )}
    </div>
  );
}
```

Notes:

- Use `srcObject` (via a `ref`), not `src` — `MediaStream` doesn't go in `src`.
- Mute your own outgoing audio in the rendered `<video>` (set `muted`), or you'll hear yourself.
- The browser may refuse autoplay without `muted` until the user has interacted with the page. The standard pattern is `muted autoplay playsInline` on every render.

## Typed metadata

If you defined a peer-metadata type when you joined:

```tsx
type PeerMetadata = { displayName: string; role: 'host' | 'guest' };
type ServerMetadata = { userId: string };

const { localPeer, remotePeers } = usePeers<PeerMetadata, ServerMetadata>();

// peer.metadata is typed as Metadata<PeerMetadata, ServerMetadata>
remotePeers.forEach((peer) => {
  const display = peer.metadata?.peer.displayName;
  const userId = peer.metadata?.server.userId;
});
```

`PeerMetadata` is what you passed to `joinRoom({ peerMetadata })` (client-set, can be updated via `useUpdatePeerMetadata`). `ServerMetadata` is what your backend passed to `createPeer({ metadata })` (immutable from the client).

## Simulcast: ask for a lower-quality variant

When you don't need full-resolution from every peer (e.g. small tiles in a grid):

```tsx
import { Variant } from '@fishjam-cloud/react-client';

// For one peer's camera track:
peer.cameraTrack?.setReceivedQuality(Variant.VARIANT_LOW);

// In bulk:
const { setReceivedTracksQuality } = usePeers();
setReceivedTracksQuality(remotePeers.flatMap((p) => p.tracks.map((t) => t.trackId)), Variant.VARIANT_LOW);
```

Only works if the publisher enabled simulcast on `videoConfig` (`FishjamProvider`).

## `peers` (deprecated alias)

`usePeers()` also returns a `peers` field which is the same array as `remotePeers`. It's flagged `@deprecated` — use `remotePeers` going forward.

## Sources

- `@fishjam-cloud/react-client` → `packages/react-client/src/hooks/usePeers.ts`
- <https://fishjam.swmansion.com/docs/how-to/client/list-other-peers>
- See `simulcast-and-bandwidth.md` for `Variant` and bandwidth setup.
