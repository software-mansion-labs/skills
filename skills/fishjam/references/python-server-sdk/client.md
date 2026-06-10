# FishjamClient (Python)

The REST client. **Synchronous** — no `await`. Source: `fishjam/api/_fishjam_client.py` in `fishjam-server-sdk`.

## Constructor

```python
from fishjam import FishjamClient

fishjam_client = FishjamClient(
    fishjam_id=os.environ['FISHJAM_ID'],
    management_token=os.environ['FISHJAM_MANAGEMENT_TOKEN'],
)
```

Use one instance per process. The constructor accepts both positional and keyword args; pin yourself to keyword for clarity.

## Rooms

```python
from fishjam import FishjamClient, RoomOptions

room = fishjam_client.create_room()  # defaults: conference room_type, no webhook, server-default video codec
room = fishjam_client.create_room(RoomOptions(
    room_type='audio_only',
))
room = fishjam_client.create_room(RoomOptions(
    room_type='livestream',
    public=True,
))
room = fishjam_client.create_room(RoomOptions(
    webhook_url='https://api.example.com/webhooks/fishjam/<secret>',
    video_codec='vp8',
    max_peers=20,
))

room = fishjam_client.get_room(room_id)
rooms = fishjam_client.get_all_rooms()
fishjam_client.delete_room(room_id)
```

`RoomOptions` dataclass fields:

- `max_peers: int | None` — cap concurrent peers
- `video_codec: Literal['h264', 'vp8'] | None` — default `None` (server picks)
- `webhook_url: str | None`
- `room_type: Literal['conference', 'audio_only', 'livestream', 'audio_only_livestream', 'full_feature']` — default `'conference'`
- `public: bool` — livestream only

Returns `Room(id, config, peers)` (a dataclass).

## Peers (WebRTC clients)

```python
from fishjam import PeerOptions

peer, peer_token = fishjam_client.create_peer(
    room.id,
    PeerOptions(metadata={'userId': 'user-123', 'name': 'Jane'}),
)

fishjam_client.delete_peer(room.id, peer.id)
new_token = fishjam_client.refresh_peer_token(room.id, peer.id)
```

`PeerOptions` fields:

- `metadata: dict[str, Any] | None` — initial peer metadata set at creation. The backend cannot mutate it afterwards, but the client owning the peer can update it mid-session (web: `useUpdatePeerMetadata`), surfacing as `peerMetadataUpdated` notifications.
- `subscribe_mode: Literal['auto', 'manual']` — default `'auto'`

Return: `tuple[Peer, str]` — the peer object and the peer token. The token is valid for 24 hours from creation; the initial connect consumes it and established sessions stay alive on their own. For peers that haven't connected (or need to reconnect after 24h), call `refresh_peer_token` for a fresh one.

## Agents (server-side programmatic peers)

```python
from fishjam import AgentOptions, AgentOutputOptions

agent = fishjam_client.create_agent(
    room.id,
    AgentOptions(
        output=AgentOutputOptions(audio_format='pcm16', audio_sample_rate=16000),
        subscribe_mode='auto',
    ),
)
# `agent` is a fishjam.agent.Agent — use as async context manager:
async with agent.connect() as session:
    ...  # see agent.md
```

`AgentOutputOptions` controls the format of incoming audio (what other peers' tracks get downmixed into for the agent).

## Vapi agents

```python
from fishjam import PeerOptionsVapi

peer = fishjam_client.create_vapi_agent(
    room.id,
    PeerOptionsVapi(api_key=os.environ['VAPI_API_KEY'], call_id=vapi_call_id),
)
```

Returns just a `Peer` — Vapi runs the agent for you; nothing to await or context-manage. `PeerOptionsVapi` requires both `api_key` and `call_id` (strings); also accepts optional `subscribe_mode`.

## Selective subscriptions

Only meaningful with `subscribe_mode='manual'` rooms — see `selective-subscriptions.md`.

```python
fishjam_client.subscribe_peer(room.id, subscriber_peer_id, target_peer_id)
fishjam_client.subscribe_tracks(room.id, subscriber_peer_id, [track_id1, track_id2])
```

## Livestream tokens

```python
streamer_token = fishjam_client.create_livestream_streamer_token(room.id)
viewer_token = fishjam_client.create_livestream_viewer_token(room.id)
```

Issue only for `livestream` rooms. Public livestreams skip the viewer token.

## MoQ tokens

```python
publisher_token = fishjam_client.create_moq_token(publish_path='my-room/alice')
subscriber_token = fishjam_client.create_moq_token(subscribe_path='my-room')
```

`publish_path` and `subscribe_path` are separate kwargs (unlike the JS SDK which takes a config object).

## Exception types

All HTTP errors raise subclasses of `fishjam.errors.HTTPError`:

| Class | Triggered by |
|---|---|
| `BadRequestError` | 400 |
| `UnauthorizedError` | 401 — bad management token |
| `NotFoundError` | 404 — room or peer missing |
| `ConflictError` | 409 — name collision / state conflict |
| `ServiceUnavailableError` | 503 |
| `InternalServerError` | 500 / any other |

Pattern:

```python
from fishjam.errors import NotFoundError, UnauthorizedError, HTTPError

try:
    peer, token = fishjam_client.create_peer(room_id, options)
except NotFoundError:
    raise HTTPException(404, detail='Room not found')
except UnauthorizedError:
    # OUR creds — surface as a server-side config error, not a user 401
    logger.exception('Fishjam credentials misconfigured')
    raise HTTPException(502)
except HTTPError as e:
    logger.exception('Fishjam call failed: %s', e)
    raise
```

`HTTPError.__init__` stores `errors` from the response body — inspect for debugging.

## Calling sync client from async code

Each method is a single HTTP round-trip, usually under ~50ms. From async code, two valid patterns:

```python
# Pattern 1 — just call it. Fine for most cases.
peer, token = fishjam_client.create_peer(room.id, options)

# Pattern 2 — for hot paths or when you care about not blocking the event loop.
peer, token = await asyncio.to_thread(fishjam_client.create_peer, room.id, options)
```

If you wrap it for async, do so consistently — don't mix patterns in the same module.

## Sources

- `fishjam/api/_fishjam_client.py` in `fishjam-server-sdk`
- `fishjam/__init__.py` in `fishjam-server-sdk`
- `fishjam/errors.py` in `fishjam-server-sdk`
- <https://fishjam.swmansion.com/docs/how-to/backend/server-setup>
- <https://fishjam.swmansion.com/docs/api/server-python>
