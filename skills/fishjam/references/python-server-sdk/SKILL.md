---
name: fishjam-python-server-sdk
description: "Python server SDK for Fishjam — backends that create rooms, mint peer tokens, receive server notifications, and run voice agents. Use when writing a Python backend (FastAPI, Flask, Starlette, aiohttp) that talks to Fishjam, decorates a notification handler, decodes a Fishjam webhook, or builds an AI voice agent in Python. Trigger on: 'fishjam-server-sdk', 'pip install fishjam-server-sdk', 'from fishjam import', 'fishjam.FishjamClient', 'FishjamNotifier', 'on_server_notification', 'receive_binary', 'fishjam Agent', 'AgentSession', 'PeerOptions', 'RoomOptions', 'AgentOptions', 'AgentOutputOptions', 'OutgoingAudioTrackOptions', 'create_room', 'create_peer', 'create_agent', 'create_vapi_agent', 'create_livestream_streamer_token', 'create_moq_token', 'subscribe_peer', 'fastapi fishjam', 'flask fishjam', 'fishjam python', 'gemini fishjam python'. Python 3.10+. The REST client is synchronous; notifier and agent are async."
license: Apache-2.0
---

# Fishjam Python Server SDK

`fishjam-server-sdk` on PyPI; `import fishjam`. Server-side Python SDK for Fishjam.

> **Read `../platform/SKILL.md` first.** It defines rooms, peers, tracks, the two-tier token model, and the WS-vs-webhook tradeoff that the Python SDK is built on.

## Sync REST + async events

The architecture differs from the JS SDK in one important way:

- `FishjamClient` is **synchronous**. Methods like `create_room`, `create_peer`, `delete_peer` are plain `def`, not `async def`. Call them from sync routes or via `run_in_executor` from async contexts.
- `FishjamNotifier` is **async** (`asyncio` + `websockets`). Use `@notifier.on_server_notification` and `await notifier.connect()`.
- `Agent` is **async** (async context manager — `async with agent.connect() as session: ...`).
- `receive_binary(raw_bytes)` is **sync** — call it directly inside an async route handler or a sync one.

## Install + minimal flow

```bash
pip install fishjam-server-sdk
# or: uv add fishjam-server-sdk
# or: poetry add fishjam-server-sdk
```

```python
import os
from fishjam import FishjamClient, PeerOptions

fishjam_client = FishjamClient(
    fishjam_id=os.environ['FISHJAM_ID'],
    management_token=os.environ['FISHJAM_MANAGEMENT_TOKEN'],
)

room = fishjam_client.create_room()
peer, peer_token = fishjam_client.create_peer(
    room.id,
    PeerOptions(metadata={'userId': 'user-123'}),
)
# ship `peer_token` to your client
```

## Components

| Component | What it does | Reference |
|---|---|---|
| `FishjamClient` (sync) | REST — rooms, peers, agents, livestream/MoQ tokens, subscribe modes. | `client.md` |
| `FishjamNotifier` (async) | Persistent WebSocket; one handler via decorator gets every event. | `notifier.md` |
| `receive_binary(bytes)` (sync) | Decode an HTTP webhook body to a typed event. | `webhooks.md` |
| `Agent` + `AgentSession` (async) | Server-side voice agent — async iter over incoming, `await track.send_chunk(...)`. | `agent.md` |

## When to read which reference

| Task | Reference |
|---|---|
| Any REST method on `FishjamClient` | `client.md` |
| Background WS task receiving events | `notifier.md` |
| HTTP route receiving webhook payloads | `webhooks.md` |
| Build a voice agent in Python | `agent.md` |
| Wire Gemini Live to an agent | `gemini-integration.md` |
| FastAPI / Flask wiring patterns | `fastapi.md` |
| Issue livestream / MoQ tokens | `livestream-and-moq.md` |
| Use manual subscribe mode | `selective-subscriptions.md` |

## Key rules

- **Python 3.10+** required. SDK uses `match`, structural typing, PEP 604 unions.
- **`FishjamClient` is sync.** Calling it from `async def` works (it's a quick HTTP round-trip) but if you're in a hot async path consider `await asyncio.to_thread(fishjam_client.create_peer, ...)`.
- **`FishjamNotifier` does NOT auto-reconnect.** Its `connect()` coroutine runs until the WS closes. Wrap it in a supervisor loop.
- **One handler per notifier.** `@notifier.on_server_notification` only stores the *last* function you decorate. Multi-dispatch in your one handler via `match` on the message type.
- **Notifier handler can be sync or async.** Both work — `await result if inspect.isawaitable(result)`.
- **Webhook decoding is sync.** `receive_binary(await req.body())` is fine inside an async FastAPI route.
- **Agent is async-only.** `async with agent.connect() as session:` is the only entry point.
- **Management token stays on the server.** Standard rule — never leak.

## References

| File | When to read |
|---|---|
| `client.md` | `FishjamClient` — every method, dataclass, error type. |
| `notifier.md` | `FishjamNotifier`, `wait_ready`, the decorator pattern, pattern-matching on event types. |
| `webhooks.md` | `receive_binary`, FastAPI / Flask receivers, security patterns. |
| `agent.md` | `Agent`, `AgentSession`, `OutgoingAudioTrackOptions`, `receive()` async iterator. |
| `gemini-integration.md` | `fishjam.integrations.gemini`, the `GeminiIntegration` singleton. |
| `fastapi.md` | FastAPI patterns (`Depends`-injected client, async webhook route, background notifier). |
| `livestream-and-moq.md` | `create_livestream_streamer_token`, `create_livestream_viewer_token`, `create_moq_token`. |
| `selective-subscriptions.md` | `subscribe_peer` / `subscribe_tracks` with `subscribe_mode='manual'` rooms. |
