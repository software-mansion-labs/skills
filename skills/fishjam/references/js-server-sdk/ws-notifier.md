# FishjamWSNotifier

Long-lived WebSocket subscription to **all server-side events** from your Fishjam instance. Source: `packages/js-server-sdk/src/ws_notifier.ts` in the `js-server-sdk` repo (<https://github.com/fishjam-cloud/js-server-sdk>).

> For the **event payload taxonomy**, read `../platform/notifications-taxonomy.md`. This file covers wiring, lifecycle, and reconnection — not the payload field semantics.

## Constructor

```ts
import { FishjamWSNotifier } from '@fishjam-cloud/js-server-sdk';

const notifier = new FishjamWSNotifier(
  { fishjamId, managementToken },                 // FishjamConfig
  (errEvent) => console.error(errEvent),           // onError
  (code, reason) => console.warn('ws closed', code, reason), // onClose
);
```

- Both `onError` and `onClose` are **required** positional callbacks — the constructor signature is `constructor(config: FishjamConfig, onError: (msg: Event) => void, onClose: (code: number, reason: string) => void)`.
- Opens one WebSocket at `<fishjam-url>/socket/server/websocket`.
- Sends an `authRequest` and a `subscribeRequest(eventType: SERVER_NOTIFICATION)` on open.
- Decodes binary `ServerMessage` per frame, finds the populated oneof field, and emits an event named after that field.
- Class extends `TypedEmitter<NotificationEvents>` (`typed-emitter` package) so `.on(...)` is statically type-checked.

## Listening

```ts
notifier.on('roomCreated', ({ roomId }) => { ... });
notifier.on('peerAdded',   ({ roomId, peerId, peerType }) => { ... });
notifier.on('peerConnected', ({ roomId, peerId, peerType }) => { ... });
notifier.on('peerDisconnected', ({ roomId, peerId, peerType }) => { ... });
notifier.on('peerDeleted',  ({ roomId, peerId, peerType }) => { ... });
notifier.on('peerMetadataUpdated', ({ roomId, peerId, peerType, metadata }) => {
  // `metadata` is a JSON-encoded string; parse it if the writer used an object.
});
notifier.on('peerCrashed',  ({ roomId, peerId, peerType, reason }) => { ... }); // `reason` is a string
notifier.on('trackAdded',   ({ roomId, peerId, componentId, track }) => {
  // `peerId` is undefined when the publisher is a component; `componentId` is set instead.
  // `track` may be undefined; when set it's `{ id, type, metadata }`.
});
notifier.on('trackRemoved', ({ roomId, peerId, componentId, track }) => { ... });
notifier.on('trackMetadataUpdated', ({ roomId, peerId, componentId, track }) => { ... });
notifier.on('streamerConnected',     ({ streamId, streamerId }) => { ... });
notifier.on('streamerDisconnected',  ({ streamId, streamerId }) => { ... });
notifier.on('viewerConnected',       ({ streamId, viewerId }) => { ... });
notifier.on('viewerDisconnected',    ({ streamId, viewerId }) => { ... });
notifier.on('channelAdded',   ({ roomId, channelId, peerId, componentId }) => { ... });
notifier.on('channelRemoved', ({ roomId, channelId, peerId, componentId }) => { ... });
notifier.on('roomDeleted',  ({ roomId }) => { ... });
notifier.on('roomCrashed',  ({ roomId }) => { ... });
```

Note that streamer / viewer events use `streamId` (the livestream ID), **not** `roomId`. (The registry's "with `roomId` override" description for these events is misleading — the `Override<T, M>` mechanism only substitutes keys that already exist on `T`, and streamer/viewer proto messages have no `roomId` to override.) That's all 18. Detail of each event's trigger and exact field set: `../platform/notifications-taxonomy.md`.

## Type imports

If you write a handler in a separate function, import the event-payload types:

```ts
import type {
  PeerConnected,
  TrackAdded,
  RoomDeleted,
  NotificationEvents,
} from '@fishjam-cloud/js-server-sdk';

function handlePeerConnected(event: PeerConnected) { ... }
notifier.on('peerConnected', handlePeerConnected);
```

`NotificationEvents` is the full mapped type — keys are event names, values are handler signatures.

## Notifier does NOT auto-reconnect

This is the most-bitten-by gotcha. The constructor opens **one** WebSocket. When it closes — for any reason — your `onClose` is called once, and the notifier becomes inert. You must instantiate a fresh `FishjamWSNotifier` to continue receiving events.

Patterns:

```ts
function startNotifier(): FishjamWSNotifier {
  return new FishjamWSNotifier(
    { fishjamId, managementToken },
    (err) => console.error('fishjam ws error', err),
    (code, reason) => {
      console.warn('fishjam ws closed', code, reason);
      setTimeout(() => {
        currentNotifier = startNotifier();
        attachHandlers(currentNotifier);
      }, 1000);
    },
  );
}

let currentNotifier = startNotifier();
attachHandlers(currentNotifier);
```

Add exponential backoff and a circuit breaker if you're paranoid about flapping. Or **use webhooks** if your deployment isn't a single long-lived process.

## Per-instance, not per-room

The notifier delivers events for **every room** in your Fishjam instance — there's no way to scope it to one room. Filter at the consumer.

If you need per-room scoping with restart-safety, prefer webhooks: each room gets its own `webhookUrl` and the URL itself is the routing primitive.

## When to switch to webhooks

You probably want webhooks instead of the WS notifier if **any** of these apply:

- You have more than one Node.js process behind a load balancer.
- You're on a serverless platform (Lambda, Vercel, Cloud Run, Cloud Functions).
- You can't afford to lose events during a deploy or restart.
- You want each room's events on a different endpoint (per-tenant routing).

See `webhooks.md`.

## Sources

- `js-server-sdk` repo: `packages/js-server-sdk/src/ws_notifier.ts`, `packages/js-server-sdk/src/notifications.ts`
- `packages/fishjam-proto/src/fishjam/server_notifications.ts` (raw proto field shapes for streamer / viewer events)
- Cross-references: `../platform/notifications-taxonomy.md`, `../platform/notifier-vs-webhook.md`
- <https://fishjam.swmansion.com/docs/how-to/backend/server-setup> (Listening to events)
