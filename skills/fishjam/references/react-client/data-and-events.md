# Data Channels, VAD, and Metadata

Three hooks for non-media interaction with the room.

## useDataChannel

WebRTC data-channel messaging. Send `Uint8Array` payloads to all peers in the room.

Source: `@fishjam-cloud/react-client` → `packages/react-client/src/hooks/useDataChannel.ts`.

```tsx
import { useDataChannel } from '@fishjam-cloud/react-client';

const {
  initializeDataChannel,
  publishData,
  subscribeData,
  dataChannelReady,
  dataChannelLoading,
  dataChannelError,
} = useDataChannel();
```

### Lifecycle

```tsx
useEffect(() => {
  if (peerStatus === 'connected') initializeDataChannel();
}, [peerStatus]);
```

You **must** call `initializeDataChannel()` after the peer is connected. Don't call before — it'll error with `"Peer is not connected"`.

`dataChannelReady` flips `true` once the channel is open. `dataChannelLoading` is true during init. `dataChannelError` holds any error.

### Sending data

```tsx
import type { DataChannelOptions } from '@fishjam-cloud/react-client';

const message = new TextEncoder().encode(JSON.stringify({ text: 'hello' }));

const options: DataChannelOptions = {
  // Fishjam exposes a single boolean. Two underlying channels exist —
  // one reliable (ordered, guaranteed) and one lossy (unordered, low latency).
  reliable: true,
};

publishData(message, options);
```

`DataChannelOptions` is just `{ reliable: boolean }`:

- `reliable: true` — uses the reliable channel: ordered, guaranteed delivery (TCP-like). Use for chat, presence, signalling.
- `reliable: false` — uses the lossy channel: unordered, low latency, may drop (UDP-like). Use for game state, cursor positions, anything where freshness beats completeness.

### Receiving data

The receive callback gets only the binary payload — there is no `peerId` argument. If you need sender identity, embed it in the encoded message yourself (e.g. JSON with a `from` field) or read it from `usePeers()` context.

```tsx
useEffect(() => {
  if (!dataChannelReady) return;

  const unsubscribe = subscribeData(
    (data) => {
      const message = JSON.parse(new TextDecoder().decode(data));
      console.log('received', message);
    },
    { reliable: true },
  );

  return () => unsubscribe();
}, [dataChannelReady, subscribeData]);
```

`subscribeData` returns an unsubscribe function. The `options` you pass select which underlying channel you listen on — subscribe to the same reliability mode you publish on.

### Text chat pattern

```tsx
function TextChat() {
  const { publishData, subscribeData, dataChannelReady, initializeDataChannel } = useDataChannel();
  const { localPeer } = usePeers();
  const [messages, setMessages] = useState<Array<{ from: string; text: string }>>([]);
  const { peerStatus } = useConnection();

  useEffect(() => {
    if (peerStatus === 'connected') initializeDataChannel();
  }, [peerStatus]);

  useEffect(() => {
    if (!dataChannelReady) return;
    const unsubscribe = subscribeData(
      (data) => {
        const { from, text } = JSON.parse(new TextDecoder().decode(data));
        setMessages((m) => [...m, { from, text }]);
      },
      { reliable: true },
    );
    return unsubscribe;
  }, [dataChannelReady, subscribeData]);

  const send = (text: string) => {
    const payload = JSON.stringify({ from: localPeer?.id, text });
    publishData(new TextEncoder().encode(payload), { reliable: true });
  };

  return <ChatUI messages={messages} onSend={send} />;
}
```

## useVAD (voice activity detection)

Tells you whether each peer in a watch-list is currently speaking. Useful for active-speaker UI, agent activation, transcription gating.

```tsx
import { useVAD, usePeers, type PeerId } from '@fishjam-cloud/react-client';

const { localPeer, remotePeers } = usePeers();
const peerIds = [localPeer?.id, ...remotePeers.map((p) => p.id)].filter(Boolean) as PeerId[];

const vadStatuses = useVAD({ peerIds });
// vadStatuses: Record<PeerId, boolean>  — true means voice activity right now

const speakingPeerIds = (Object.keys(vadStatuses) as PeerId[]).filter((id) => vadStatuses[id]);
```

Signature: `useVAD({ peerIds: ReadonlyArray<PeerId> }) -> Record<PeerId, boolean>`.

- Remote-peer VAD is server-driven (`vadNotification` from the backend).
- Including the **local peer's** id activates client-side polling of the microphone's audio level (handled internally; you don't see `useLocalVAD` directly in the public surface).
- Returned record updates reactively as voice activity changes.

## useUpdatePeerMetadata

Update the client-side `peerMetadata` you passed when calling `joinRoom`. Different from the **server-side** metadata set in `createPeer` — that one is immutable.

```tsx
import { useUpdatePeerMetadata } from '@fishjam-cloud/react-client';

const { updatePeerMetadata } = useUpdatePeerMetadata<{ hand: 'raised' | 'lowered' }>();

const raiseHand = () => updatePeerMetadata({ hand: 'raised' });
```

The update is broadcast to all peers and surfaces on their side via `usePeers()`. The combined client + server metadata is exposed under `peer.metadata` (typed as `Metadata<PeerMetadata, ServerMetadata>` from `@fishjam-cloud/ts-client`); the client-side update only mutates the client-set portion — server-set metadata from `createPeer` is immutable.

Use cases:

- "Hand raised" indicator.
- Typing-in-chat status.
- Sub-room / breakout assignment.
- Mic-muted state (separate from the actual track being on/off).

## Sources

- `@fishjam-cloud/react-client` → `packages/react-client/src/hooks/useDataChannel.ts`
- `packages/react-client/src/hooks/useVAD.ts`
- `packages/react-client/src/hooks/useUpdatePeerMetadata.ts`
- `packages/webrtc-client/src/types.ts` (`DataChannelOptions`, `DataCallback`)
- <https://fishjam.swmansion.com/docs/how-to/client/text-chat>
- <https://fishjam.swmansion.com/docs/how-to/client/metadata>
- <https://fishjam.swmansion.com/docs/explanation/data-channels>
