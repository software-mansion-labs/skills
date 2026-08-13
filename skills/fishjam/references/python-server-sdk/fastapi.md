# FastAPI / Flask Wiring

Production patterns for hosting Fishjam from a Python backend. The official docs show FastAPI; the reference example app uses Flask. Both work.

## FastAPI

### Singleton client via dependency

```python
import os
from functools import lru_cache
from typing import Annotated
from fastapi import Depends, FastAPI, HTTPException

from fishjam import FishjamClient, PeerOptions, RoomOptions
from fishjam.errors import NotFoundError, UnauthorizedError, HTTPError

@lru_cache
def get_fishjam() -> FishjamClient:
    return FishjamClient(
        fishjam_id=os.environ['FISHJAM_ID'],
        management_token=os.environ['FISHJAM_MANAGEMENT_TOKEN'],
    )

FishjamDep = Annotated[FishjamClient, Depends(get_fishjam)]

app = FastAPI()
```

### `/api/join-room`

```python
from pydantic import BaseModel

class JoinRoomRequest(BaseModel):
    room_name: str
    room_type: str = 'conference'

class JoinRoomResponse(BaseModel):
    peer_token: str
    room_id: str
    fishjam_id: str

@app.post('/api/join-room', response_model=JoinRoomResponse)
async def join_room(
    req: JoinRoomRequest,
    fishjam: FishjamDep,
    user: Annotated[User, Depends(authenticate_user)],  # your auth dependency
):
    if not user_can_join(user, req.room_name):
        raise HTTPException(403)

    try:
        # Look up an existing room id, or create one
        room = get_or_create_room(req.room_name, req.room_type, fishjam)

        peer, peer_token = fishjam.create_peer(
            room.id,
            PeerOptions(metadata={
                'userId': user.id, 'name': user.name, 'role': user.role,
            }),
        )
    except UnauthorizedError as e:
        # OUR Fishjam credentials — surface as a server error, NOT user 401
        raise HTTPException(502, detail='Upstream auth failure') from e
    except NotFoundError as e:
        raise HTTPException(404, detail='Room not found') from e
    except HTTPError as e:
        raise HTTPException(500, detail=str(e)) from e

    return JoinRoomResponse(
        peer_token=peer_token,
        room_id=room.id,
        fishjam_id=os.environ['FISHJAM_ID'],
    )
```

Both Fishjam HTTP calls (`get_or_create_room` and `create_peer`) are sync — they block briefly inside an async route. For high-throughput services, wrap with `asyncio.to_thread(...)`.

### Webhook receiver (FastAPI)

```python
from fastapi import Request, HTTPException
from fishjam import receive_binary
from fishjam.events import ServerMessagePeerConnected, ServerMessagePeerDisconnected, ServerMessageRoomDeleted

@app.post('/webhooks/fishjam/{secret}')
async def fishjam_webhook(secret: str, request: Request):
    if secret != os.environ['FISHJAM_WEBHOOK_SECRET']:
        raise HTTPException(401)

    notification = receive_binary(await request.body())
    if notification is None:
        return {'ok': True}

    match notification:
        case ServerMessagePeerConnected() as msg:
            await record_peer_connected(msg.room_id, msg.peer_id)
        case ServerMessagePeerDisconnected() as msg:
            await record_peer_disconnected(msg.room_id, msg.peer_id)
        case ServerMessageRoomDeleted() as msg:
            await cleanup_room(msg.room_id)

    return {'ok': True}
```

### Background notifier on startup

If you'd rather use the WebSocket notifier than webhooks (single-process backend):

```python
from contextlib import asynccontextmanager
from fishjam import FishjamNotifier
from fishjam.events import ServerMessagePeerConnected
import asyncio

notifier_task: asyncio.Task | None = None

@asynccontextmanager
async def lifespan(app: FastAPI):
    global notifier_task
    notifier = FishjamNotifier(
        fishjam_id=os.environ['FISHJAM_ID'],
        management_token=os.environ['FISHJAM_MANAGEMENT_TOKEN'],
    )

    @notifier.on_server_notification
    async def handle(msg):
        match msg:
            case ServerMessagePeerConnected():
                ...

    async def supervise():
        while True:
            try:
                await notifier.connect()
            except Exception:
                logger.exception('notifier dropped, restarting')
            await asyncio.sleep(1)

    notifier_task = asyncio.create_task(supervise())
    await notifier.wait_ready()
    yield
    notifier_task.cancel()

app = FastAPI(lifespan=lifespan)
```

Use **either** webhooks **or** the notifier, not both — picking duplicates events.

### CORS

```python
from fastapi.middleware.cors import CORSMiddleware

app.add_middleware(
    CORSMiddleware,
    allow_origins=[os.environ['FRONTEND_ORIGIN']],
    allow_credentials=True,
    allow_methods=['POST', 'GET'],
    allow_headers=['*'],
)
```

## Flask (reference example pattern)

The official Python example (`examples/room_manager/` in `fishjam-server-sdk`) uses Flask. Worth knowing if you inherit a Flask codebase:

```python
from flask import Flask, request, jsonify, abort
from flask_cors import CORS
from fishjam import FishjamClient, PeerOptions, RoomOptions, receive_binary

app = Flask(__name__)
CORS(app)

fishjam_client = FishjamClient(
    fishjam_id=os.environ['FISHJAM_ID'],
    management_token=os.environ['FISHJAM_MANAGEMENT_TOKEN'],
)

@app.get('/api/rooms')
def join_room():
    room_name = request.args.get('roomName')
    peer_name = request.args.get('peerName')
    if not room_name or not peer_name:
        return abort(400)

    room = get_or_create_room(room_name)
    peer, token = fishjam_client.create_peer(
        room.id, PeerOptions(metadata={'name': peer_name}),
    )
    return jsonify({'peerToken': token, 'roomId': room.id})

@app.post('/api/rooms/webhook')
def webhook():
    notification = receive_binary(request.data)
    if notification:
        handle(notification)
    return ('Webhook Notification Received', 200)
```

Flask doesn't run the async notifier well — for Flask, prefer webhooks. If you really need the notifier alongside Flask, run it in a separate process or thread (`threading.Thread(target=lambda: asyncio.run(supervise()))`).

## Sources

- `examples/room_manager/main.py`, `routes.py`, `room_service.py` in `fishjam-server-sdk`
- <https://fishjam.swmansion.com/docs/how-to/backend/fastapi-example>
- <https://fishjam.swmansion.com/docs/how-to/backend/production-deployment>
