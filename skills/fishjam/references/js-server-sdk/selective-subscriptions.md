# Selective Subscriptions

By default, every peer in a Fishjam room is subscribed to every other peer's tracks (`subscribeMode: 'auto'`). This is what you want for most apps — meetings, chat rooms, multi-party calls.

**Manual subscribe mode** changes that: a peer receives nothing until your backend explicitly subscribes them. Useful for:

- **Large rooms / cost control.** A 50-peer "audience" shouldn't receive 49 video streams each. Subscribe only to the active speakers.
- **Privacy / staging.** Hide certain peers from each other until your app permits it (presenter-only mode, breakout sessions, etc.).
- **Selective AI processing.** An agent only needs the audio of the peer who's currently speaking; subscribe based on VAD.
- **Spotlight UX.** Front-row vs back-row visibility.

## Enabling

`subscribeMode` is set **per peer**, at peer creation — it lives on `PeerOptionsWebRTC` / `PeerOptionsAgent` / `PeerOptionsVapi`, **not** on `RoomConfig`. A single room can mix peers with `'auto'` and `'manual'` modes.

```ts
const { peer, peerToken } = await fishjamClient.createPeer(roomId, {
  metadata: { userId: 'user-123' },
  subscribeMode: 'manual',
});

// Agents:
const { agent, peer: agentPeer } = await fishjamClient.createAgent(roomId, {
  output: { audioFormat: 'pcm16', audioSampleRate: 16000 },
  subscribeMode: 'manual',
});
```

You cannot toggle this after the peer is created. Decide at peer creation.

## Subscribing a peer to another peer

`subscribePeer(roomId, subscriberPeerId, publisherPeerId)` — subscriber receives **all current and future tracks** from publisher.

```ts
await fishjamClient.subscribePeer(room.id, viewerPeerId, presenterPeerId);
```

You can call this multiple times to subscribe one peer to many publishers. There's no `unsubscribePeer` — to "unsubscribe", remove the publisher's tracks via `deletePeer`, or design with manual subscribe and explicit re-subscribe semantics for each scene change.

## Subscribing a peer to specific tracks

`subscribeTracks(roomId, subscriberPeerId, trackIds)` — subscriber receives **only the specified tracks**. Use when you want one peer's audio but not their video, or just the screen-share track but not the camera.

```ts
import type { TrackId } from '@fishjam-cloud/js-server-sdk';

// In your trackAdded handler:
notifier.on('trackAdded', async ({ roomId, peerId, track }) => {
  if (track.type === 'audio' && peerId === currentSpeakerId) {
    await fishjamClient.subscribeTracks(
      roomId,
      activeListenerPeerId,
      [track.id as TrackId],
    );
  }
});
```

## Typical patterns

### Subscribe everyone to everyone (manual but blanket)

Use `subscribeMode: 'manual'` on each peer + a `peerConnected` handler that subscribes the new peer to every existing peer. Same effect as auto mode, but with the explicit control point in your code.

> **Precondition:** the **subscriber** peer (second argument) must have been created with `subscribeMode: 'manual'`. The publisher's mode does not matter. If the subscriber is `'auto'`, the call is a silent no-op.

```ts
// Every createPeer call elsewhere in your code:
//   await fishjamClient.createPeer(roomId, { subscribeMode: 'manual' });

notifier.on('peerConnected', async ({ roomId, peerId }) => {
  const room = await fishjamClient.getRoom(roomId);
  for (const other of room.peers) {
    if (other.id === peerId) continue;
    await fishjamClient.subscribePeer(roomId, peerId, other.id);
    await fishjamClient.subscribePeer(roomId, other.id, peerId);
  }
});
```

### Spotlight / single-presenter mode

```ts
// Audience peers MUST be created with subscribeMode: 'manual' for subscribePeer to take effect.
// e.g. await fishjamClient.createPeer(roomId, { subscribeMode: 'manual', metadata: { role: 'audience' } });

// Subscribe everyone to the presenter; the presenter to nobody.
notifier.on('peerConnected', async ({ roomId, peerId }) => {
  if (peerId === presenterPeerId) return;
  await fishjamClient.subscribePeer(roomId, peerId, presenterPeerId);
});
```

### Active-speaker subscription (e.g. for an agent)

```ts
// Use a client-side VAD signal forwarded to the backend, or a server-side VAD on the agent.
// When the speaker changes, subscribe the agent to the new speaker's audio.
async function switchSpeaker(newSpeakerId: PeerId, agentPeerId: PeerId) {
  const room = await fishjamClient.getRoom(roomId);
  const audioTrack = room.peers.find((p) => p.id === newSpeakerId)?.tracks?.find((t) => t.type === 'audio');
  if (audioTrack) await fishjamClient.subscribeTracks(roomId, agentPeerId, [audioTrack.id]);
}
```

## REST endpoints (raw)

```http
POST /room/{room_id}/peer/{id}/subscribe_peer       # subscribe to all of one peer's tracks
POST /room/{room_id}/peer/{id}/subscribe_tracks     # subscribe to specific track IDs
```

The SDK methods wrap these. Hand-rolling raw calls only makes sense if you're not using the SDK.

## Sources

- `js-server-sdk` repo: `packages/js-server-sdk/src/client.ts` (`subscribePeer`, `subscribeTracks`), `packages/fishjam-openapi/src/generated/api.ts` (`PeerOptionsWebRTC.subscribeMode`, etc.)
- Canonical example: `examples/selective-subscription/backend/src/service/fishjam.ts` in the same repo
- <https://fishjam.swmansion.com/docs/how-to/backend/selective-subscriptions>
- Cross-reference: `../platform/rest-endpoints.md`
