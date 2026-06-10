# Agent (Python)

Server-side programmatic peer for Python. Async-only. Source: `fishjam/agent/agent.py` in `fishjam-server-sdk`.

> Different shape from the JS `FishjamAgent`. The Python API is built around an **async context manager** and an **async iterator** for incoming messages.

## Lifecycle

```python
from fishjam import AgentOptions, AgentOutputOptions, FishjamClient
from fishjam.agent import OutgoingAudioTrackOptions
from fishjam.events import TrackEncoding

# create_agent returns immediately — connect happens inside the context manager
agent = fishjam_client.create_agent(
    room.id,
    AgentOptions(
        output=AgentOutputOptions(audio_format='pcm16', audio_sample_rate=16000),  # audio_format only 'pcm16'; audio_sample_rate only 16000 | 24000 (pick 24000 for Gemini Live)
        subscribe_mode='auto',
    ),
)

async with agent.connect() as session:
    # session is an AgentSession; can send audio and iterate over incoming
    ...
```

`fishjam_client.create_agent(...)` is a sync REST call that returns an unconnected `Agent` instance. The connection only opens when you `async with agent.connect():`. Authentication happens inside the context manager — `AgentAuthError` is raised if the agent token is bad.

## Adding an outgoing audio track

```python
from fishjam.agent import OutgoingAudioTrackOptions
from fishjam.events import TrackEncoding

async with agent.connect() as session:
    track = await session.add_track(OutgoingAudioTrackOptions(
        encoding=TrackEncoding.TRACK_ENCODING_PCM16,
        sample_rate=24000,
        channels=1,
        metadata={'source': 'tts'},
    ))
```

`OutgoingAudioTrackOptions` fields:

- `encoding: TrackEncoding` — `TRACK_ENCODING_OPUS`, `TRACK_ENCODING_PCM16`.
- `sample_rate: Literal[16000, 24000]`
- `channels: Literal[1, 2]` — mono or stereo
- `metadata: dict | None` — JSON-encodable

Returns an `OutgoingTrack` you push frames to.

## Sending audio

```python
await track.send_chunk(audio_bytes)
```

`audio_bytes` is raw codec frames matching the track's `encoding` and `sample_rate`. Pace yourself — there's no built-in backpressure.

## Interrupting playback

If the agent has been queueing audio for the room and the user starts talking, drop the buffer:

```python
await track.interrupt()
```

Any frames sent but not yet played are discarded. Subsequent `send_chunk` calls play normally.

## Receiving incoming audio + images

```python
from fishjam.events._protos.fishjam import (
    AgentResponseTrackData as IncomingTrackData,
    AgentResponseTrackImage as IncomingTrackImage,
)

async with agent.connect() as session:
    async for msg in session.receive():
        match msg:
            case IncomingTrackData() as data:
                # data.peer_id, data.track (a Track — note `data.track.id` may be empty string for unset wire fields), data.data (bytes)
                process_incoming_audio(data.data)
            case IncomingTrackImage() as image:
                # image.track_id, image.content_type, image.data (bytes)
                process_captured_image(image)
```

What you receive is governed by the `AgentOptions.output` you set on `create_agent` — peer tracks get downmixed into the format Gemini-style consumers expect.

## Capturing an image from a video track

```python
# In one task — request the capture (one-way):
await session.capture_image(track_id='peer-track-id')

# Image arrives via the same receive() iterator as an IncomingTrackImage.
```

The capture is asynchronous — you fire-and-forget the request, then watch the `receive()` stream for the matching `IncomingTrackImage`. Match on `image.track_id` to correlate.

## Disconnecting

The context manager handles it on exit. To disconnect explicitly inside the block:

```python
await session.disconnect()
```

## Full minimal example

```python
import asyncio
from fishjam import AgentOptions, AgentOutputOptions, FishjamClient
from fishjam.agent import OutgoingAudioTrackOptions
from fishjam.events import TrackEncoding
from fishjam.events._protos.fishjam import AgentResponseTrackData as IncomingTrackData

async def run():
    client = FishjamClient(fishjam_id=..., management_token=...)
    room = client.create_room()
    agent = client.create_agent(room.id, AgentOptions())

    async with agent.connect() as session:
        track = await session.add_track(OutgoingAudioTrackOptions(
            encoding=TrackEncoding.TRACK_ENCODING_PCM16, sample_rate=24000, channels=1,
        ))

        async def listen():
            async for msg in session.receive():
                if isinstance(msg, IncomingTrackData):
                    print(f'got {len(msg.data)} bytes from {msg.peer_id}')

        async def push():
            while True:
                await track.send_chunk(generate_audio_frame())
                await asyncio.sleep(0.02)

        await asyncio.gather(listen(), push())

asyncio.run(run())
```

## Vapi agents (managed alternative)

If you don't need raw audio control, prefer `create_vapi_agent` — Vapi runs the agent, no `Agent` / `AgentSession` to manage:

```python
from fishjam import PeerOptionsVapi
peer = fishjam_client.create_vapi_agent(
    room.id,
    PeerOptionsVapi(api_key=os.environ['VAPI_API_KEY'], call_id=vapi_call_id),
)
```

## Sources

- `fishjam/agent/agent.py` in `fishjam-server-sdk`
- `fishjam/agent/errors.py` in `fishjam-server-sdk` (`AgentAuthError`)
- `fishjam/events/_protos/fishjam/__init__.py` (`AgentResponseTrackData`, `AgentResponseTrackImage`)
- <https://fishjam.swmansion.com/docs/tutorials/agents>
- <https://fishjam.swmansion.com/docs/explanation/agent-internals>
