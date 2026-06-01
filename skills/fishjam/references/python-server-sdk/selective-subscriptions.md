# Selective Subscriptions (Python)

By default, every peer is subscribed to every other peer (`subscribe_mode='auto'` on `PeerOptions`). For large or scene-controlled rooms, enable **manual** subscriptions: peers receive nothing until your backend explicitly subscribes them.

Use cases:

- Audience-style rooms (50 listeners + 1 speaker): subscribe everyone to the speaker, no one to the audience members.
- Spotlight / presenter mode.
- Server-side AI agents that should only receive the active speaker's audio.
- Privacy isolation between sub-groups.

## Enabling

Manual mode is configured **per peer** in this SDK:

```python
from fishjam import PeerOptions

peer, token = fishjam_client.create_peer(
    room.id,
    PeerOptions(subscribe_mode='manual'),
)
```

Once a peer is created with `subscribe_mode='manual'`, it receives nothing until you explicitly subscribe it. Mix peers with different subscribe modes in the same room — each peer's policy is its own.

## Subscribing a peer to all of another peer's tracks

```python
fishjam_client.subscribe_peer(
    room_id=room.id,
    peer_id=subscriber_peer_id,
    target_peer_id=publisher_peer_id,
)
```

The subscriber receives every current and future track from the target. No unsubscribe API — to "unsubscribe", remove the publisher's tracks or design with explicit re-subscribe.

## Subscribing to specific tracks only

```python
fishjam_client.subscribe_tracks(
    room_id=room.id,
    peer_id=subscriber_peer_id,
    track_ids=[track_id_1, track_id_2],
)
```

Use when you want a peer's audio but not their video, the screen share but not the camera, etc.

## Typical pattern: spotlight

```python
# When a new peer connects, subscribe them to the current presenter only.
# (Handled in a webhook or notifier.)

async def on_peer_connected(room_id: str, peer_id: str, peer_type: str):
    if peer_id == presenter_peer_id:
        return  # presenter doesn't subscribe to listeners
    fishjam_client.subscribe_peer(
        room_id=room_id,
        peer_id=peer_id,
        target_peer_id=presenter_peer_id,
    )
```

## Typical pattern: AI agent subscribed only to active speaker

```python
# When VAD reports a new active speaker, switch the agent's subscription.

def switch_agent_to_speaker(new_speaker_peer_id: str):
    # `t.type_` is a TrackType enum from the generated OpenAPI client.
    # No public re-export exists today, so either compare against the string
    # value, or import from the private path (subject to change between SDK
    # releases):
    #   from fishjam._openapi_client.models.track_type import TrackType

    room = fishjam_client.get_room(room_id)
    speaker = next(p for p in room.peers if p.id == new_speaker_peer_id)
    audio_track = next(t for t in speaker.tracks if t.type_ == "audio")

    fishjam_client.subscribe_tracks(
        room_id=room_id,
        peer_id=agent_peer_id,
        track_ids=[audio_track.id],
    )
```

## Sources

- `fishjam/api/_fishjam_client.py` in `fishjam-server-sdk` (`subscribe_peer`, `subscribe_tracks`)
- `examples/selective_subscription/` in `fishjam-server-sdk`
- <https://fishjam.swmansion.com/docs/how-to/backend/selective-subscriptions>
