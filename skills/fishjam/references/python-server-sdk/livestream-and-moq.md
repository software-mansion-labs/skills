# Livestream & MoQ Tokens (Python)

Three token types issued by `FishjamClient` for non-WebRTC-conference flows. All synchronous.

## Livestream rooms

```python
from fishjam import FishjamClient, RoomOptions

room = fishjam_client.create_room(RoomOptions(room_type='livestream'))
# or public (no viewer tokens needed):
room = fishjam_client.create_room(RoomOptions(room_type='livestream', public=True))

streamer_token = fishjam_client.create_livestream_streamer_token(room.id)
# → ship to the single broadcasting client

viewer_token = fishjam_client.create_livestream_viewer_token(room.id)
# → ship to each viewer (private livestreams only)
```

### Public vs private

| Mode | `public` | Viewer token? |
|---|---|---|
| Private (default) | `False` | required |
| Public | `True` | not used — viewers connect by room ID |

`public` is set at room creation and is immutable for that room.

### Constraints

- **One streamer per room.** Issuing a second streamer token doesn't add a second broadcaster.
- **One video + one audio track on the streamer.** Extra tracks dropped silently.

### Revocation

Livestream streamer and viewer tokens issued via `create_livestream_streamer_token(room_id)` / `create_livestream_viewer_token(room_id)` cannot be revoked individually — there is no public REST endpoint or SDK method for it. To invalidate all outstanding tokens for a livestream room, call `fishjam_client.delete_room(room_id)`.

## MoQ tokens

```python
publisher_token = fishjam_client.create_moq_token(publish_path='my-room/alice')
subscriber_token = fishjam_client.create_moq_token(subscribe_path='my-room')
```

Path scoping is **prefix-based**:

| `publish_path` value | Allows publishing to |
|---|---|
| `'stream-name'` | Any path under `stream-name/*` — client picks identity |
| `'stream-name/alice'` | Only the exact path `stream-name/alice` — identity pinned |

Same for `subscribe_path`. Your Fishjam ID is the implicit root namespace; don't include it in the path.

The token is attached as `?jwt=<token>` to `https://relay.fishjam.io/<fishjam-id>?jwt=<token>` (HTTPS bootstrap — MoQ runs over WebTransport/QUIC, not WebSocket). The Fishjam ID is the only URL path segment; the publish/subscribe path is encoded in the JWT.

See `../platform/room-types.md` for MoQ conceptual overview.

## Client-side counterpart

- React web: `../react-client/livestream.md` (uses `useLivestreamStreamer` / `useLivestreamViewer`).
- React Native: `../react-native-client/SKILL.md` (re-exports those hooks).
- MoQ: no first-party React hook — use a MoQ-compatible JS library (`@moq/publish`, `@moq/watch`).

## Sources

- `fishjam/api/_fishjam_client.py` in `fishjam-server-sdk` (`create_livestream_*_token`, `create_moq_token`)
- <https://fishjam.swmansion.com/docs/tutorials/livestreaming>
- <https://fishjam.swmansion.com/docs/tutorials/moq>
- <https://fishjam.swmansion.com/docs/explanation/moq-streaming>
