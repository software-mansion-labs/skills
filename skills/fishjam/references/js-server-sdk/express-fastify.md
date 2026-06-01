# Express / Fastify Wiring

Production-quality patterns for hosting Fishjam from a Node.js backend. Covers both Express (most common) and Fastify (the example pattern in the official docs).

## Single client per process

Always — one `FishjamClient` shared across all routes / requests. It's an axios wrapper, not a connection pool. Re-instantiating per request leaks.

```ts
const fishjamClient = new FishjamClient({
  fishjamId: process.env.FISHJAM_ID!,
  managementToken: process.env.FISHJAM_MANAGEMENT_TOKEN!,
});
```

## Express — the canonical `/api/join-room` endpoint

```ts
import express from 'express';
import {
  FishjamClient,
  FishjamBaseException,
  FishjamNotFoundException,
  UnauthorizedException,
  BadRequestException,
  ServiceUnavailableException,
} from '@fishjam-cloud/js-server-sdk';

const app = express();
app.use(express.json());

const fishjamClient = new FishjamClient({
  fishjamId: process.env.FISHJAM_ID!,
  managementToken: process.env.FISHJAM_MANAGEMENT_TOKEN!,
});

// Your own auth middleware. NOT Fishjam's — it verifies YOUR user.
app.use(authenticateUser);

app.post('/api/join-room', async (req, res, next) => {
  try {
    const { roomName, roomType = 'conference' } = req.body;
    const user = req.user;

    if (!userCanJoin(user, roomName)) {
      return res.status(403).json({ error: 'Forbidden' });
    }

    // Either look up an existing room by name in your DB, or create a fresh one.
    const room = await getOrCreateRoom(roomName, roomType);

    const { peer, peerToken } = await fishjamClient.createPeer(room.id, {
      metadata: { userId: user.id, name: user.name, role: user.role },
    });

    res.json({ peerToken, roomId: room.id, fishjamId: process.env.FISHJAM_ID });
  } catch (err) {
    next(err); // delegated to the error middleware below
  }
});

// Error middleware — map Fishjam exceptions to HTTP responses
app.use((err: Error, req, res, next) => {
  if (err instanceof FishjamNotFoundException) {
    return res.status(404).json({ error: err.message });
  }
  if (err instanceof BadRequestException) {
    return res.status(400).json({ error: err.message, details: err.details });
  }
  if (err instanceof UnauthorizedException) {
    // Our credentials are bad — surface this as a server error, NOT a user 401.
    console.error('Fishjam auth failure — check FISHJAM_MANAGEMENT_TOKEN', err);
    return res.status(502).json({ error: 'Upstream auth failure' });
  }
  if (err instanceof ServiceUnavailableException) {
    return res.status(503).json({ error: 'Upstream unavailable' });
  }
  if (err instanceof FishjamBaseException) {
    // Catch-all for any other Fishjam exception (incl. 403 surfaced as UnknownException).
    return res.status(502).json({ error: 'Upstream error', status: err.statusCode, details: err.details });
  }
  console.error('Unhandled error', err);
  res.status(500).json({ error: 'Internal server error' });
});

app.listen(3000);
```

Production checklist for this endpoint:

- **CORS** — `cors({ origin: process.env.FRONTEND_ORIGIN, credentials: true })` so the client can call it.
- **Rate limiting** — `express-rate-limit` per IP and per user. This is what protects you from token spam.
- **Idempotency** — pick whether `roomName` resolves to an existing room or always creates a new one. Match your product's semantics.
- **Authorization** — `userCanJoin(user, roomName)` is the policy hook. Fishjam authorizes the token bearer, not your user. Enforce *your* rules here.

## Fastify — official pattern

The Fastify docs show the SDK wired as a plugin so `fastify.fishjam` is decorated everywhere:

```ts
import fastifyPlugin from 'fastify-plugin';
import { FishjamClient } from '@fishjam-cloud/js-server-sdk';

declare module 'fastify' {
  interface FastifyInstance {
    fishjam: FishjamClient;
    config: { FISHJAM_ID: string; FISHJAM_MANAGEMENT_TOKEN: string };
  }
}

export const fishjamPlugin = fastifyPlugin((fastify) => {
  fastify.decorate(
    'fishjam',
    new FishjamClient({
      fishjamId: fastify.config.FISHJAM_ID,
      managementToken: fastify.config.FISHJAM_MANAGEMENT_TOKEN,
    }),
  );
});
```

Routes:

```ts
import type { RoomType } from '@fishjam-cloud/js-server-sdk';

fastify.post<{ Body: { roomName: string; roomType?: RoomType } }>(
  '/api/join-room',
  async (request, reply) => {
    const { roomName, roomType = 'conference' } = request.body;
    const room = await getOrCreateRoom(roomName, roomType);
    const { peer, peerToken } = await fastify.fishjam.createPeer(room.id, {
      metadata: { userId: request.user.id },
    });
    return { peerToken, roomId: room.id };
  },
);

fastify.get('/api/rooms', () => fastify.fishjam.getAllRooms());
```

For Fastify notifier and webhook wiring, see `ws-notifier.md` and `webhooks.md`, plus the official guide at <https://fishjam.swmansion.com/docs/how-to/backend/fastify-example>.

## Combining `/api/join-room` + webhook + notifier

Common production shape:

| Concern | Where it lives |
|---|---|
| `/api/join-room` (mint peer token) | HTTP route — see Express / Fastify above |
| `/api/leave-room` (optional explicit cleanup) | HTTP route → `fishjamClient.deletePeer(roomId, peerId)` |
| `/api/admin/rooms` (debugging) | HTTP route → `fishjamClient.getAllRooms()` |
| `/webhooks/fishjam/:secret` (events) | Raw-body route, `ServerMessage.decode`; see `webhooks.md` |
| `FishjamWSNotifier` (events, alt) | Background plugin / startup hook; see `ws-notifier.md` |
| Token refresh (long-lived peers) | Cron / scheduler calling `fishjamClient.refreshPeerToken(...)`; ship new token via websocket / SSE to the client |

## Sources

- `js-server-sdk` repo: `packages/js-server-sdk/src/client.ts` (constructor JSDoc mentions Fastify pattern), `examples/room-manager/` (Fastify reference implementation). There is no first-party Express example in the SDK repo.
- <https://fishjam.swmansion.com/docs/how-to/backend/fastify-example>
- <https://fishjam.swmansion.com/docs/how-to/backend/production-deployment>
