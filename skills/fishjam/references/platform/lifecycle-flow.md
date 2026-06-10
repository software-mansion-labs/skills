# Lifecycle Flow

End-to-end picture of a Fishjam session, from the moment a user taps "Join" to the moment the room is gone. Skim this once when wiring a new app; refer back for the exact event order.

## Full sequence

1. **User → Backend:** User clicks **Join**.
2. **Backend → Fishjam:** `createRoom` (idempotent), or reuse an existing room.
     - **Fishjam → Backend:** Returns `Room`.
     - **Fishjam event:** `roomCreated` (to notifier/webhook).
3. **Backend → Fishjam:** `createPeer(roomId, metadata)`.
     - **Fishjam → Backend:** Returns `{ peer, peerToken }`.
     - **Fishjam event:** `peerAdded`.
4. **Backend → User/Client (HTTP response):** Returns `{ peerToken, fishjamId }`.
5. **Client → Fishjam:** `joinRoom({ peerToken })`.
     - WebRTC negotiation starts.
     - **Fishjam event:** `peerConnected`.
6. **Client ↔ Fishjam:** Client publishes local tracks (via SDK hooks).
     - **Fishjam event:** `trackAdded` (repeated for each published track).
7. **Client ↔ Fishjam:** Client receives remote tracks.
8. **Client:** Leaves room or closes browser/app.
     - **Fishjam event:** `peerDisconnected`.
9. **Backend → Fishjam (optional cleanup):** `deletePeer`.
     - **Fishjam event:** `peerDeleted`.
10. **Backend → Fishjam (when last peer leaves):** `deleteRoom`.
      - **Fishjam event:** `roomDeleted`.

## Event order summary

For a single peer:

- `roomCreated` (once per room)
  - `peerAdded` → `peerConnected`
    - `trackAdded` (received for each published track)
    - `trackMetadataUpdated` / `trackRemoved` (during session)
    - `peerMetadataUpdated` (if you allow metadata updates)
  - `peerDisconnected`
  - `peerDeleted` (if explicitly removed or token expired)
- `roomDeleted` (when room is torn down)

`peerDisconnected` may be followed by another `peerConnected` if the client reconnects within the peer-token's validity window. Treat the pair as a connectivity blip, not a session-end.

## Token refresh path

For peers that need to outlive 24h:

1. **Trigger:** Just before the current peer token reaches the 24h limit.
2. **Backend → Fishjam:** Call `refreshPeerToken` (`/room/{roomId}/peer/{peerId}/refresh_token`).
3. **Fishjam → Backend:** Returns the new peer token. SDK shape: `string` (JS `refreshPeerToken` → `Promise<string>`, Python `refresh_peer_token` → `str`). Raw HTTP shape: `{ "data": { "token": "..." } }`.
4. **Backend → Client:** Forward the new peer token.
5. **Client behavior:**
     - Reconnect using the new token, or
     - Store it and use it on the next connection attempt.

In practice, short-lived rooms (one call = one room) are easier than long-lived peers with refresh logic. Use refresh only when the product requires a peer to persist across sessions.

## Reconnection

The client SDKs (web + RN) auto-reconnect with the same peer token as long as it's still valid. Configure with the `reconnect` prop on `FishjamProvider`:

```tsx
// Full ReconnectConfig: { maxAttempts?: number, initialDelay?: number, delay?: number, addTracksOnReconnect?: boolean }
<FishjamProvider
  fishjamId={FISHJAM_ID}
  reconnect={{ maxAttempts: 5, initialDelay: 1000, delay: 500 }}
>
  <App />
</FishjamProvider>
```

On the backend side, you don't normally need to do anything — Fishjam keeps the peer record alive across short disconnects. If the peer is unreachable for long enough, Fishjam may emit `peerDeleted`; budget your reconnection time accordingly.

## Failure modes

| Symptom | Likely cause | Remediation |
| --- | --- | --- |
| Client can't connect, sees auth error | Peer token expired (24h) or revoked | Refresh and retry, or treat as new session |
| Backend gets `peerAdded` but no `peerConnected` follows | Client never actually called `joinRoom`, or NAT/firewall blocked WebRTC | Check client logs |
| `peerConnected` then immediately `peerDisconnected` | Failed ICE / TURN renegotiation | Same; check client network and Fishjam status page |
| Webhook delivery fails (your endpoint 500s) | Fishjam retries 5xx with exponential backoff (up to 8 attempts, capped at 30s total wall-clock). 4xx is not retried. Per-room dispatch is serialised, so a slow endpoint stalls subsequent events for that room. | Return 2xx fast; if you 4xx, the event is dropped — reconcile via `GET /room/{id}`. |
| WS notifier silently stops emitting | Process lost the socket, no auto-reconnect | Supervise the notifier, restart on close |

## Sources

- Backend examples: each server SDK's `examples/` directory in its GitHub repo
- <https://fishjam.swmansion.com/docs/tutorials/backend-quick-start>
- <https://fishjam.swmansion.com/docs/how-to/client/reconnection-handling>
- See `auth-model.md` for token-issuance diagram, `notifications-taxonomy.md` for event payload detail.
