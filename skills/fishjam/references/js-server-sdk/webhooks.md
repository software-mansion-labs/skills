# Webhook Receiver

Fishjam can POST every server event to a URL of your choice instead of (or alongside) using `FishjamWSNotifier`. Transport is stateless HTTP. Best for serverless or scaled-out backends. Decoding is a binary protobuf parse.

> The 18-event set documented in `../platform/notifications-taxonomy.md` is what the WS notifier filters down to. Webhook bodies carry the **raw** `ServerMessage` oneof, so a receiver may also see `vadNotification`, `trackForwarding` / `trackForwardingRemoved`, and the deprecated `streamConnected` / `streamDisconnected` / `hlsPlayable` / `hlsUploaded` / `hlsUploadCrashed` / `componentCrashed` fields depending on what the server emits. Branch on the `oneof` field explicitly and treat unrecognized cases as no-ops.

> Pick between webhook and WS notifier first — see `../platform/notifier-vs-webhook.md`. For event payload semantics: `../platform/notifications-taxonomy.md`.

## Enabling

Set `webhookUrl` when creating a room:

```ts
await fishjamClient.createRoom({
  webhookUrl: `${process.env.PUBLIC_BASE_URL}/webhooks/fishjam/${process.env.FISHJAM_WEBHOOK_SECRET}`,
});
```

That URL receives every event for that room until the room is deleted OR crashes; after either, the notifier shuts down 15 s later (drain window for in-flight events). There is no global "all rooms" webhook — each room is configured at creation.

### Batched delivery

Set `batchWebhookNotifications: true` on `RoomConfig` to have Fishjam coalesce multiple events into a single `ServerMessage_NotificationBatch` POST instead of one HTTP request per event. The proto wrapper is `ServerMessage.notificationBatch` — a `NotificationBatch` cannot itself be nested inside another `NotificationBatch`. If you enable this, your handler must check `msg.notificationBatch` first and iterate its inner notifications. VAD notifications are always sent individually regardless of this setting.

```ts
await fishjamClient.createRoom({
  webhookUrl: '...',
  batchWebhookNotifications: true,
});
```

## Decoding the body

Body is `Content-Type: application/x-protobuf`, raw bytes of a `ServerMessage`. Re-export from the SDK:

```ts
import { ServerMessage } from '@fishjam-cloud/js-server-sdk';

const msg = ServerMessage.decode(new Uint8Array(rawBody));
```

`ServerMessage` is a protobuf `oneof` — exactly one of its fields is populated per message. Switch on the populated field:

```ts
function handle(msg: ServerMessage) {
  if (msg.roomCreated)         handleRoomCreated(msg.roomCreated);
  else if (msg.peerAdded)      handlePeerAdded(msg.peerAdded);
  else if (msg.peerConnected)  handlePeerConnected(msg.peerConnected);
  else if (msg.peerDisconnected) handlePeerDisconnected(msg.peerDisconnected);
  else if (msg.trackAdded)     handleTrackAdded(msg.trackAdded);
  else if (msg.trackRemoved)   handleTrackRemoved(msg.trackRemoved);
  else if (msg.streamerConnected) ...
  else if (msg.viewerConnected)   ...
  // …all 18 fields, same names as the notifier event keys
}
```

Note that webhook payloads use the raw protobuf shape — `peerType` is the enum int (e.g. `ServerMessage_PeerType.PEER_TYPE_WEBRTC`), and `track.type` is the `TrackType` enum, not the friendly strings. The WS notifier maps these via an internal `mapNotification` helper before emitting; it is **not** re-exported from `@fishjam-cloud/js-server-sdk`, so handle the raw enums directly in your webhook code (compare against `ServerMessage_PeerType.PEER_TYPE_WEBRTC`, etc.).

## Express receiver

```ts
import express from 'express';
import { ServerMessage } from '@fishjam-cloud/js-server-sdk';

const app = express();

app.post(
  '/webhooks/fishjam/:secret',
  express.raw({ type: 'application/x-protobuf' }), // <-- crucial: get raw Buffer, not JSON
  (req, res) => {
    if (req.params.secret !== process.env.FISHJAM_WEBHOOK_SECRET) {
      return res.sendStatus(401);
    }

    let msg: ServerMessage;
    try {
      msg = ServerMessage.decode(new Uint8Array(req.body));
    } catch (e) {
      console.error('Bad protobuf body', e);
      return res.sendStatus(400);
    }

    handle(msg);
    res.sendStatus(200);
  },
);
```

The `express.raw({ type: 'application/x-protobuf' })` middleware is **specific** to the protobuf content type — other routes still get the normal JSON parser. Don't apply raw globally, or you'll break the rest of your API.

## Fastify receiver

```ts
import Fastify, { FastifyRequest } from 'fastify';
import { ServerMessage } from '@fishjam-cloud/js-server-sdk';

const fastify = Fastify();

fastify.addContentTypeParser(
  'application/x-protobuf',
  { parseAs: 'buffer' },
  async (_req: FastifyRequest, body: Buffer) => ServerMessage.decode(new Uint8Array(body)),
);

fastify.post<{ Body: ServerMessage; Params: { secret: string } }>(
  '/webhooks/fishjam/:secret',
  async (request, reply) => {
    if (request.params.secret !== process.env.FISHJAM_WEBHOOK_SECRET) {
      return reply.code(401).send();
    }
    handle(request.body); // already a decoded ServerMessage
    return reply.code(200).send();
  },
);
```

This is the pattern in <https://fishjam.swmansion.com/docs/how-to/backend/fastify-example>.

## Security

Fishjam **does not sign webhook bodies**. Pick one of:

1. **Path secret** (shown above) — a long random string embedded in the URL. Rotate by re-creating rooms with a new `webhookUrl`.
2. **IP allowlist** — restrict your receiver to known Fishjam egress IPs (check the current list with Software Mansion if you go this route).
3. **Front-of-Fishjam proxy** — your own gateway that adds an HMAC header and verifies it downstream.

The path-secret pattern is the easiest and works in serverless. Treat the secret like any production credential.

## Idempotency & retries

Fishjam retries 5xx and transport errors with randomized exponential backoff (up to 8 attempts, ~30 s total budget). 4xx is not retried. Design handlers to be idempotent — repeats are expected. VAD notifications use a separate "lossy" path (500 ms timeout, no retry). Per-room dispatch is serialised, so a slow endpoint stalls subsequent events for that room.

- Return 2xx fast — if you 4xx, the event is dropped.
- **Re-derive state on the next event** rather than treating events as a transactional ledger.
- For accuracy when an event matters (billing, audit), reconcile against `GET /room/{room_id}` on a schedule.

## Sources

- `js-server-sdk` repo: `packages/js-server-sdk/src/notifications.ts`, `packages/fishjam-proto/src/fishjam/server_notifications.ts` (`ServerMessage`, `ServerMessage_NotificationBatch`)
- `packages/fishjam-openapi/src/generated/api.ts` (`RoomConfig.batchWebhookNotifications`)
- <https://fishjam.swmansion.com/docs/how-to/backend/server-setup>
- <https://fishjam.swmansion.com/docs/how-to/backend/fastify-example>
