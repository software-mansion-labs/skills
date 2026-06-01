# Gemini Live Integration (Python)

Source: `fishjam/integrations/gemini.py` in `fishjam-server-sdk`.

The integration is shipped as an **optional extra** — install with:

```bash
pip install 'fishjam-server-sdk[gemini]'
```

Importing without the extra raises an `ImportError` directing you to install it.

## What the helper provides

A singleton instance `GeminiIntegration`:

```python
from fishjam.integrations.gemini import GeminiIntegration
```

With these members:

| Member | What it is |
|---|---|
| `GeminiIntegration.create_client(api_key=..., http_options=..., ...)` | Returns a configured `google.genai.Client` with an `x-goog-api-client: fishjam-python-server-sdk/<version>` header injected. Mirrors `genai.Client(...)` constructor args. |
| `GeminiIntegration.GEMINI_INPUT_AUDIO_SETTINGS` | `AgentOutputOptions(audio_format='pcm16', audio_sample_rate=16000)` — pass to `AgentOptions.output` so other peers' audio is downmixed for Gemini. |
| `GeminiIntegration.GEMINI_OUTPUT_AUDIO_SETTINGS` | `OutgoingAudioTrackOptions(encoding=TRACK_ENCODING_PCM16, sample_rate=24000, channels=1)` — pass to `session.add_track(...)` so Gemini's TTS output flows back into the room. |
| `GeminiIntegration.GEMINI_AUDIO_MIME_TYPE` | `'audio/pcm;rate=16000'` — attach to audio chunks sent to Gemini. |

## End-to-end pattern

```python
import asyncio
import os
from fishjam import AgentOptions, FishjamClient
from fishjam.integrations.gemini import GeminiIntegration
from fishjam.events._protos.fishjam import AgentResponseTrackData as IncomingTrackData
from google.genai import types

fishjam_client = FishjamClient(fishjam_id=..., management_token=...)
genai_client = GeminiIntegration.create_client(api_key=os.environ['GEMINI_API_KEY'])

async def spawn_gemini_agent(room_id: str):
    agent = fishjam_client.create_agent(
        room_id,
        AgentOptions(output=GeminiIntegration.GEMINI_INPUT_AUDIO_SETTINGS),
    )

    async with agent.connect() as session:
        # Gemini -> Fishjam track
        gemini_track = await session.add_track(GeminiIntegration.GEMINI_OUTPUT_AUDIO_SETTINGS)

        # Open Gemini Live session
        async with genai_client.aio.live.connect(
            model='gemini-2.5-flash-native-audio-preview-12-2025',
            config=types.LiveConnectConfig(
                response_modalities=[types.Modality.AUDIO],
                system_instruction='You are a helpful assistant.',
            ),
        ) as gemini_session:

            async def fishjam_to_gemini():
                async for msg in session.receive():
                    if isinstance(msg, IncomingTrackData):
                        await gemini_session.send_realtime_input(
                            audio=types.Blob(
                                data=msg.data,
                                mime_type=GeminiIntegration.GEMINI_AUDIO_MIME_TYPE,
                            )
                        )

            async def gemini_to_fishjam():
                async for response in gemini_session.receive():
                    if response.data:
                        await gemini_track.send_chunk(response.data)

            await asyncio.gather(fishjam_to_gemini(), gemini_to_fishjam())

asyncio.run(spawn_gemini_agent('room-123'))
```

## When to use the helper

- **Use** for the canonical 16kHz-in / 24kHz-out Gemini Live flow. Saves you from picking codec parameters wrong.
- **Skip** if you want a different model, different audio formats, or only text — import `google.genai` directly.

## Multimodal: image capture for visual prompting

```python
await session.capture_image(track_id=peer_video_track_id)
# response arrives as IncomingTrackImage via session.receive()
# forward to Gemini as a `types.Blob` with mime_type='image/jpeg'
```

## Sources

- `fishjam/integrations/gemini.py` in `fishjam-server-sdk`
- `examples/multimodal/` in `fishjam-server-sdk` (canonical reference)
- <https://fishjam.swmansion.com/docs/integrations/gemini-live-integration>
