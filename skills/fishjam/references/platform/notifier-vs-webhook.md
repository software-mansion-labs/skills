# WebSocket Notifier vs Webhook

Fishjam's SDKs surface **18 server-side events** covering the room / peer / track / streamer / viewer / channel lifecycle. (The underlying `ServerMessage` protobuf carries additional oneof fields that the SDKs intentionally drop; if you decode the wire protocol yourself you'll see them.) Two transports deliver the same set with the same payload shape — you choose by deployment topology.

| | WebSocket notifier | Webhook |
|---|---|---|
| Transport | Single long-lived WS (`/socket/server/websocket`) | HTTP POST per event |
| Encoding | Binary protobuf (`ServerMessage`) | Binary protobuf (`Content-Type: application/x-protobuf`) |
| Server class (JS) | `FishjamWSNotifier` | (decode `ServerMessage` from raw body) |
| Server class (Python) | `FishjamNotifier` | `receive_binary(raw_body)` |
| Subscription scope | Whole Fishjam instance | One room (set on `createRoom({ webhookUrl })`) |
| Auth | Management token in WS handshake | Whatever you build into the URL path |
| State | Notifier reconnect on its own? **No** — supervise yourself | Stateless |
| Best for | Single dedicated worker, integration test suites, dev tools | Production backends, serverless, multi-instance, K8s |

## Pick WebSocket if…

- You have a **single long-lived process** that always needs to see every event.
- You're running in dev / CI / a test fixture.
- You want low integration overhead — `notifier.on('peerAdded', …)` and you're done.

Limitations to plan for:

- One WS subscription = one process. If you scale out the worker, you'll get duplicate event handling.
- The notifier does **not** auto-reconnect. Wrap it in a supervisor / restart loop that calls back into your code if the socket drops.
- Auth is the management token, so the WS process needs it. Same blast-radius rules apply (`auth-model.md`).

JS:

```ts
import { FishjamWSNotifier } from '@fishjam-cloud/js-server-sdk';

const notifier = new FishjamWSNotifier(
  { fishjamId, managementToken },
  (err) => { console.error(err); /* restart */ },
  (code, reason) => { console.warn('WS closed', code, reason); /* restart */ },
);

notifier.on('peerConnected', ({ roomId, peerId, peerType }) => {
  // ...
});
```

Python:

```python
import asyncio
from fishjam import FishjamNotifier
from fishjam.events import ServerMessagePeerConnected

notifier = FishjamNotifier(fishjam_id, management_token)

@notifier.on_server_notification
def handle(notification):
    match notification:
        case ServerMessagePeerConnected():
            ...
        case _:
            pass

async def main():
    task = asyncio.create_task(notifier.connect())
    await notifier.wait_ready()
    await task  # returns when the WS drops — wrap in `while True: try: await notifier.connect() except ...` to auto-reconnect
```

## Pick webhooks if…

- You're running on **serverless** (Lambda, Cloud Functions, Cloud Run, Vercel functions, etc.).
- Your backend is horizontally scaled across multiple instances behind a load balancer.
- You want each room's events to land on the closest endpoint (you can set a per-room `webhookUrl`).
- You can persist POST bodies on your side for replay / debugging (Fishjam itself has no replay store).

Configure once at room creation:

```ts
await fishjamClient.createRoom({
  webhookUrl: 'https://api.your-app.com/webhooks/fishjam/<shared-secret>',
});
```

```python
fishjam_client.create_room(RoomOptions(webhook_url='https://api.your-app.com/webhooks/fishjam/<shared-secret>'))
```

Decode bodies:

```ts
// Express — raw body middleware required
import express from 'express';
import { ServerMessage } from '@fishjam-cloud/js-server-sdk';

app.post(
  '/webhooks/fishjam/:secret',
  express.raw({ type: 'application/x-protobuf' }),
  (req, res) => {
    if (req.params.secret !== process.env.FISHJAM_WEBHOOK_SECRET) return res.sendStatus(401);
    const msg = ServerMessage.decode(new Uint8Array(req.body));
    // msg has one populated oneof field, e.g. msg.peerConnected, msg.trackAdded, …
    res.sendStatus(200);
  },
);
```

```python
# FastAPI
from fastapi import APIRouter, Request, HTTPException
from fishjam import receive_binary

@router.post('/webhooks/fishjam/{secret}')
async def fishjam_webhook(secret: str, req: Request):
    if secret != os.environ['FISHJAM_WEBHOOK_SECRET']:
        raise HTTPException(401)
    notification = receive_binary(await req.body())
    if notification is None:
        return {'ok': True}
    # match on notification type — see notifications-taxonomy.md
    return {'ok': True}
```

### Webhook security pattern

Fishjam doesn't sign webhook bodies. Common pattern: include a long random secret in the URL path (`/webhooks/fishjam/<secret>`) and reject mismatched secrets.

## Hybrid

Nothing stops you from using both for different rooms — for example, webhooks for production rooms and the notifier for a dev / debug process attached to a sandbox room. Just remember the WS notifier sees events for **every** room in the instance, while webhook URLs are per-room.

## Event set

Both transports deliver the same 18 SDK-surfaced events. Detail: `notifications-taxonomy.md`.

## Source

- JS: `packages/js-server-sdk/src/ws_notifier.ts` and `notifications.ts` in <https://github.com/fishjam-cloud/js-server-sdk>
- Python: `fishjam/_ws_notifier.py` and `_webhook_notifier.py` in <https://github.com/fishjam-cloud/python-server-sdk>
- Docs: <https://fishjam.swmansion.com/docs/how-to/backend/server-setup> (Listening to events)
