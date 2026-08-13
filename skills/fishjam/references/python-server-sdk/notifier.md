# FishjamNotifier (Async)

Long-lived WebSocket subscription to every server-side event. Async-first. Source: `fishjam/_ws_notifier.py` in `fishjam-server-sdk`.

> Event-payload semantics live in `../platform/notifications-taxonomy.md`. This file covers wiring, lifecycle, and the decorator pattern.

## Constructor + handler

```python
from fishjam import FishjamNotifier
from fishjam.events import (
    ServerMessageRoomCreated,
    ServerMessagePeerAdded,
    ServerMessagePeerConnected,
    ServerMessageTrackAdded,
    # ... more event classes — see the union below
)
from fishjam.events.allowed_notifications import AllowedNotification

notifier = FishjamNotifier(
    fishjam_id=os.environ['FISHJAM_ID'],
    management_token=os.environ['FISHJAM_MANAGEMENT_TOKEN'],
)

@notifier.on_server_notification
def handle(notification: AllowedNotification) -> None:
    match notification:
        case ServerMessageRoomCreated() as msg:
            print('room created', msg.room_id)
        case ServerMessagePeerAdded() as msg:
            print('peer added', msg.room_id, msg.peer_id, msg.peer_type)
        case ServerMessagePeerConnected() as msg:
            ...
        case ServerMessageTrackAdded() as msg:
            ...
        case _:
            pass  # other events
```

**Exactly one handler.** Decorating a second function replaces the first. Multi-dispatch inside the one handler via `match`.

The handler can be a `def` (sync) or `async def`. If async, `connect()` awaits it.

## The 18 event classes

All under `fishjam.events`:

```text
ServerMessageRoomCreated
ServerMessageRoomDeleted
ServerMessageRoomCrashed

ServerMessagePeerAdded
ServerMessagePeerDeleted
ServerMessagePeerConnected
ServerMessagePeerDisconnected
ServerMessagePeerMetadataUpdated
ServerMessagePeerCrashed

ServerMessageStreamerConnected
ServerMessageStreamerDisconnected
ServerMessageViewerConnected
ServerMessageViewerDisconnected

ServerMessageTrackAdded
ServerMessageTrackRemoved
ServerMessageTrackMetadataUpdated

ServerMessageChannelAdded
ServerMessageChannelRemoved
```

Each is a `betterproto` dataclass. Fields are snake_case (`room_id`, `peer_id`, `peer_type`, `track`, etc.). The `AllowedNotification` union type imported from `fishjam.events.allowed_notifications` is the type alias for "any of these 18".

**Field-name gotcha — streamer / viewer events do NOT carry `room_id`.** They use `stream_id` (the livestream ID) plus `streamer_id` / `viewer_id`:

- `ServerMessageStreamerConnected { stream_id, streamer_id }`
- `ServerMessageStreamerDisconnected { stream_id, streamer_id }`
- `ServerMessageViewerConnected { stream_id, viewer_id }`
- `ServerMessageViewerDisconnected { stream_id, viewer_id }`

**`ServerMessageChannelAdded` / `Removed`** carry `room_id`, `channel_id` (required), and exactly one of `peer_id` / `component_id` (the endpoint that owns the channel).

**Track events** (`ServerMessageTrackAdded` / `Removed` / `MetadataUpdated`) carry `room_id`, optional `peer_id`, optional `component_id` (exactly one of those two is set), and optional `track`.

## Running

```python
import asyncio

async def main():
    notifier_task = asyncio.create_task(notifier.connect())
    await notifier.wait_ready()        # subscribe is in
    # ...do other startup work...
    await notifier_task                # run until WS closes

asyncio.run(main())
```

`connect()` runs as long as the WebSocket is open. Once it closes (for any reason), `connect()` returns — there's no auto-reconnect.

`wait_ready()` resolves once the subscribe is acknowledged. Useful when other startup tasks want to wait until the notifier is actually receiving events.

## Reconnect pattern

```python
async def run_notifier_forever():
    while True:
        try:
            await notifier.connect()
        except Exception:
            logger.exception('notifier dropped, restarting in 1s')
        await asyncio.sleep(1)
```

Add exponential backoff if your event loop tolerates it. Or switch to webhooks if your deployment is multi-process / serverless (see `webhooks.md`).

## Sync handler

```python
@notifier.on_server_notification
def handle(notification: AllowedNotification) -> None:
    # CPU-light: print, push to a queue, etc.
    queue.put_nowait(notification)
```

Watch out for blocking I/O — the handler runs on the event loop. For real work (DB writes, HTTP calls), prefer async:

```python
@notifier.on_server_notification
async def handle(notification: AllowedNotification) -> None:
    match notification:
        case ServerMessagePeerConnected() as msg:
            await record_peer_connected(msg.room_id, msg.peer_id)
```

## Auth failure detection

The notifier raises `RuntimeError('Invalid management token')` if Fishjam rejects the management token during the handshake. Catch this at startup — it usually means a misconfigured env var.

## When to switch to webhooks

Switch to webhooks (per-room `webhook_url`) if any apply:

- You have multiple Python workers and don't want to dedicate one to event handling.
- You're on serverless (Cloud Run, AWS Lambda) and can't keep a long-lived task running.
- You can't tolerate event loss across deploys.

See `webhooks.md`.

## Sources

- `fishjam/_ws_notifier.py` in `fishjam-server-sdk`
- `fishjam/events/allowed_notifications.py` in `fishjam-server-sdk`
- `static/api/protobuf/server_notifications.proto` in the documentation repo
- Cross-references: `../platform/notifications-taxonomy.md`, `../platform/notifier-vs-webhook.md`
- <https://fishjam.swmansion.com/docs/how-to/backend/server-setup>
