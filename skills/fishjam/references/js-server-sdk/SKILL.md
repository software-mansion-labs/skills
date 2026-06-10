---
name: fishjam-js-server-sdk
description: "Node.js / TypeScript server SDK for Fishjam — backends that create rooms, mint peer tokens, listen to server notifications, and run agents. Use when writing a Node.js / Express / Fastify / Hono / NestJS backend that talks to Fishjam, sets up a webhook receiver, runs an AI agent, or generates livestream tokens. Trigger on: '@fishjam-cloud/js-server-sdk', 'FishjamClient', 'FishjamWSNotifier', 'FishjamAgent', 'createPeer', 'createRoom', 'createAgent', 'createVapiAgent', 'createLivestreamStreamerToken', 'createLivestreamViewerToken', 'refreshPeerToken', 'subscribeTracks', 'subscribePeer', 'createMoqToken', 'ServerMessage', 'roomCreated', 'peerAdded', 'trackAdded', 'fishjam backend Node', 'express fishjam', 'fastify fishjam', 'gemini agent fishjam', 'vapi fishjam'. Covers REST, WebSocket notifier, webhook protobuf decoding, FishjamAgent, and Gemini Live integration."
license: MIT
---

# Fishjam JS Server SDK

`@fishjam-cloud/js-server-sdk` — server-side Node.js / TypeScript SDK for the Fishjam platform.

> **Read `../platform/SKILL.md` first.** It defines rooms, peers, tracks, management tokens, peer tokens, and the WS-vs-webhook tradeoff that this skill builds on.

## Three components

| Component | What it does | Reference |
|---|---|---|
| `FishjamClient` | REST client — rooms, peers, agents, livestream/MoQ tokens, subscribe modes. | `client.md` |
| `FishjamWSNotifier` | Subscribes to all 18 server events over a single WebSocket. Best for long-lived workers. | `ws-notifier.md` |
| `FishjamAgent` | Programmatic peer — sends/receives audio frames over a server-side WS. Used for AI agents (Gemini, custom). | `agent.md` |

Webhook receivers are not a class — you decode `ServerMessage` protobuf yourself; covered in `webhooks.md`.

## Install + minimal flow

```bash
npm install @fishjam-cloud/js-server-sdk
# or yarn add @fishjam-cloud/js-server-sdk
```

```ts
import { FishjamClient } from '@fishjam-cloud/js-server-sdk';

const fishjamClient = new FishjamClient({
  fishjamId: process.env.FISHJAM_ID!,
  managementToken: process.env.FISHJAM_MANAGEMENT_TOKEN!,
});

const room = await fishjamClient.createRoom();
const { peer, peerToken } = await fishjamClient.createPeer(room.id, {
  metadata: { userId: 'user-123' },
});
// ship `peerToken` to your client
```

## When to read which reference

| Task | Reference |
|---|---|
| Call any REST method (rooms, peers, subscribe, refresh) | `client.md` |
| Subscribe to server events from a long-lived worker | `ws-notifier.md` |
| Receive events from a serverless / scaled-out backend | `webhooks.md` |
| Build an AI voice agent (audio in/out, image capture) | `agent.md` |
| Wire a Gemini Live agent into a room | `gemini-integration.md` |
| Wire all of this into Express or Fastify | `express-fastify.md` |
| Issue livestream streamer/viewer tokens or MoQ tokens | `livestream-and-moq.md` |
| Use manual subscribe mode (cost / privacy / staging) | `selective-subscriptions.md` |

## Key rules

- **One `FishjamClient` per process.** It's a thin axios wrapper — share, don't recreate per request.
- **Catch SDK exceptions, not axios errors.** All methods wrap REST calls in `mapException`, throwing typed errors (`BadRequestException` 400, `UnauthorizedException` 401, `RoomNotFoundException` / `PeerNotFoundException` / `FishjamNotFoundException` 404, `ServiceUnavailableException` 503, `UnknownException` for anything else — including 403). All inherit from `FishjamBaseException`, which is useful for a catch-all `instanceof` branch. `ForbiddenException` is exported but never produced by the current mapper. Match on these to map to HTTP status codes. Constructor also throws `MissingFishjamIdException` synchronously if `fishjamId` is falsy.
- **Management token stays in the backend.** Never put it in `NEXT_PUBLIC_*`, `VITE_*`, or any bundled env. If exposed, regenerate from the Dashboard.
- **Deprecation header.** The client logs `x-fishjam-api-deprecated` once per process. If you see it, plan a migration — don't silence it.
- **`FishjamWSNotifier` does NOT auto-reconnect.** The constructor opens one WS; if it closes, you must reinstantiate. Wrap it in a supervisor or use webhooks instead.
- **Peer tokens are valid for 24 hours from creation.** The token is consumed during the initial WS handshake; an established session keeps running on its own and isn't dropped when the token's wall-clock expiry passes. If the peer hasn't connected yet, or needs to *reconnect* after the 24h window, call `refreshPeerToken` to mint a new one and ship it to the client.

## References

| File | When to read |
|---|---|
| `client.md` | The REST surface — every method on `FishjamClient` with signatures and examples. |
| `ws-notifier.md` | `FishjamWSNotifier` wiring, the 18 events, reconnection patterns. |
| `webhooks.md` | Decoding `application/x-protobuf` webhook bodies with `ServerMessage.decode`. |
| `agent.md` | `FishjamAgent` — opus/pcm16, `createTrack`, `sendData`, `captureImage`, `interruptTrack`. |
| `gemini-integration.md` | `@fishjam-cloud/js-server-sdk/gemini` — `createClient`, audio settings, end-to-end demo. |
| `express-fastify.md` | Production wiring — `/api/join-room`, error mapping, webhook routes, Fastify plugin. |
| `livestream-and-moq.md` | `createLivestreamStreamerToken`, `createLivestreamViewerToken`, `createMoqToken`. |
| `selective-subscriptions.md` | `subscribePeer` / `subscribeTracks` with peers created using `subscribeMode: 'manual'`. |
