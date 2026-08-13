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

// Map Fishjam exceptions to HTTP responses. Called explicitly from each route's catch
function respondFishjamError(err: unknown, res: express.Response) {
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
  return res.status(500).json({ error: 'Internal server error' });
}

// Your own auth middleware. NOT Fishjam's — it verifies YOUR user.
app.use(authenticateUser);

app.post('/api/join-room', async (req, res) => {
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
    respondFishjamError(err, res);
  }
});

app.listen(3000);
```

If the app already centralizes errors in a four-argument Express error-handling middleware, `next(err)` from the catch block reaches it the same way (that's the pattern Express's own error-handling guide documents) — but don't introduce one just for Fishjam; the explicit helper is simpler and framework-portable.

Fishjam-specific decisions for this endpoint:

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

- `js-server-sdk` repo: `packages/js-server-sdk/src/client.ts` (constructor JSDoc mentions Fastify pattern), `examples/room-manager/` (Fastify reference implementation; `src/errors.ts` has the `parseError` helper the error mapping above mirrors). There is no first-party Express example in the SDK repo.
- <https://fishjam.swmansion.com/docs/how-to/backend/fastify-example>
- <https://fishjam.swmansion.com/docs/how-to/backend/production-deployment>
- <https://expressjs.com/en/guide/error-handling.html> (the `next(err)` / error-middleware alternative)
