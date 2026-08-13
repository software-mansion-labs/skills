# Webhook Receiver (Python)

Webhooks deliver the **same 18 events** as `FishjamNotifier`, transported as stateless HTTP POSTs. Each room is configured with a `webhook_url` at creation; Fishjam POSTs binary protobuf to that URL per event.

> Pick between webhook and notifier first — see `../platform/notifier-vs-webhook.md`.

## Enabling

```python
from fishjam import RoomOptions

room = fishjam_client.create_room(RoomOptions(
    webhook_url=f"{PUBLIC_BASE_URL}/webhooks/fishjam/{WEBHOOK_SECRET}",
))
```

The URL receives every event for that room until the room is deleted. There's no global webhook — per-room only.

## Decoding the body

`fishjam.receive_binary(raw: bytes) -> AllowedNotification | None`. Synchronous; returns `None` if the body is a message type not in the allowed set (the SDK filters out handshake / deprecated events, including `NotificationBatch` wrappers — if you enable `batch_webhook_notifications=True` on the room, `receive_binary` filters out the batch wrapper and returns `None`; reach into the raw `ServerMessage` if you need its contents).

`AllowedNotification` lives at `fishjam.events.allowed_notifications` — import it from there:

```python
from fishjam.events.allowed_notifications import AllowedNotification
```

```python
from fishjam import receive_binary

msg = receive_binary(raw_body)
if msg is None:
    return  # silently ignore non-event messages
```

Decoded `msg` is a `betterproto` dataclass — one of the 18 event types — with snake_case fields (`room_id`, `peer_id`, `peer_type`, `track`).

## FastAPI receiver

```python
from fastapi import APIRouter, HTTPException, Request
from fishjam import receive_binary
from fishjam.events import (
    ServerMessagePeerAdded,
    ServerMessagePeerConnected,
    ServerMessagePeerDisconnected,
    ServerMessageTrackAdded,
    ServerMessageRoomDeleted,
)

router = APIRouter()

@router.post('/webhooks/fishjam/{secret}')
async def fishjam_webhook(secret: str, request: Request):
    if secret != WEBHOOK_SECRET:
        raise HTTPException(401)

    raw = await request.body()  # bytes
    notification = receive_binary(raw)
    if notification is None:
        return {'ok': True}

    match notification:
        case ServerMessagePeerConnected() as msg:
            await record_peer_connected(msg.room_id, msg.peer_id)
        case ServerMessagePeerDisconnected() as msg:
            await record_peer_disconnected(msg.room_id, msg.peer_id)
        case ServerMessageTrackAdded() as msg:
            ...
        case ServerMessageRoomDeleted() as msg:
            await cleanup_room(msg.room_id)
        case _:
            pass

    return {'ok': True}
```

FastAPI lets you keep the route sync if you prefer; `receive_binary` is sync either way.

## Flask receiver (official example pattern)

The reference example app uses Flask:

```python
from flask import Flask, request
from fishjam import receive_binary

app = Flask(__name__)

# `handle(notification)` must be defined in your code — see the "Pattern-matching
# strategy" section below for the canonical pattern (and `examples/room_manager/routes.py`
# in the SDK for a full reference).

@app.post('/api/rooms/webhook')
def webhook():
    notification = receive_binary(request.data)  # raw bytes
    if notification is None:
        return ('', 200)
    handle(notification)
    return ('Webhook Notification Received', 200)
```

Pattern: `request.data` gives you the raw body without parsing as JSON (default Flask). No content-type-parser configuration needed.

## Pattern-matching strategy

The event union has 18 variants — list every one your app cares about and ignore the rest:

```python
from fishjam.events.allowed_notifications import AllowedNotification
from fishjam.events import (
    ServerMessageRoomCreated,
    ServerMessageRoomDeleted,
    ServerMessagePeerAdded,
    ServerMessagePeerConnected,
    ServerMessagePeerDisconnected,
    ServerMessageTrackAdded,
    # ...
)

def handle(notification: AllowedNotification):
    match notification:
        case ServerMessageRoomCreated():     ...
        case ServerMessageRoomDeleted():     ...
        case ServerMessagePeerAdded():       ...
        case ServerMessagePeerConnected():   ...
        case ServerMessagePeerDisconnected(): ...
        case ServerMessageTrackAdded():      ...
        # ... etc
        case _:
            pass  # explicitly ignore unrecognized
```

Don't `assert False` on `_` — Fishjam may add new event types in the future and you don't want to crash on them.

## Alternative: `betterproto.which_one_of`

If you'd rather route by string name than `isinstance`, fall back to the raw protobuf:

```python
# Private import — `fishjam.events._protos` is the generated betterproto package
# and is not part of the SDK's public API. Pin your SDK version if you rely on it.
import betterproto
from fishjam.events._protos.fishjam import ServerMessage

msg = ServerMessage().parse(raw_body)
field_name, populated = betterproto.which_one_of(msg, 'content')
# field_name: 'peer_added', 'track_removed', etc.
```

But `receive_binary` already does this and filters — usually you don't need the raw path.

## Security

Fishjam doesn't sign webhook bodies. Use a path-secret pattern (long random string in the URL) as shown above. Alternatives: IP allowlist or front-of-Fishjam HMAC proxy.

## Idempotency

Fishjam retries 5xx and transport errors with randomized exponential backoff (up to 8 attempts, ~30 s total budget). 4xx is not retried. Design handlers to be idempotent — repeats are expected. VAD notifications use a separate "lossy" path (500 ms timeout, no retry). Per-room dispatch is serialised, so a slow endpoint stalls subsequent events for that room. Treat events as hints; reconcile against `fishjam_client.get_room(room_id)` when accuracy matters.

## Sources

- `fishjam/_webhook_notifier.py` in `fishjam-server-sdk`
- `examples/room_manager/routes.py` in `fishjam-server-sdk` (Flask)
- `fishjam/events/allowed_notifications.py` in `fishjam-server-sdk`
- <https://fishjam.swmansion.com/docs/how-to/backend/fastapi-example>
