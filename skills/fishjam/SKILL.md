---
name: fishjam
description: "Software Mansion's Fishjam — hosted WebRTC platform for video, audio, and one-to-many livestreaming. MUST USE before writing, reviewing, or debugging ANY code that talks to a Fishjam instance from a backend (Node, Python) or a client (React web, React Native / Expo). Routes to platform fundamentals and the right SDK sub-skill. Trigger on: 'Fishjam', 'fishjam.io', 'fishjam.swmansion.com', 'fishjam dashboard', 'Fishjam ID', 'management token', 'peer token', 'sandbox API', 'livestream room', '@fishjam-cloud/js-server-sdk', 'fishjam-server-sdk', '@fishjam-cloud/react-client', '@fishjam-cloud/react-native-client', '@fishjam-cloud/ts-client', 'FishjamClient', 'FishjamProvider', 'FishjamNotifier', 'FishjamAgent', 'createPeer', 'createRoom', 'createAgent', 'createLivestreamStreamerToken', 'createMoqToken', 'refreshPeerToken', 'webhook', 'gemini fishjam', 'vapi fishjam', 'fishjam expo plugin'."
license: MIT
---

# Fishjam

Software Mansion's hosted WebRTC platform. Routes between platform fundamentals and the four SDKs.

Read the **`references/platform/SKILL.md`** first for the domain model (rooms, peers, tracks, tokens, notifications). Then read the SDK sub-skill that matches what you're building. All sub-skills are in `references/`.

## Which sub-skill do I need?

| You are writing…                                                                  | Sub-skill                                                                                       |
| --------------------------------------------------------------------------------- | ----------------------------------------------------------------------------------------------- |
| Anything Fishjam — start here for glossary, room types, auth, notifications, REST | `references/platform/SKILL.md`                                                                  |
| A Node.js / Express / Fastify / Hono backend                                      | `references/js-server-sdk/SKILL.md`                                                             |
| A Python / FastAPI / Flask backend                                                | `references/python-server-sdk/SKILL.md`                                                         |
| A React app in a browser                                                          | `references/react-client/SKILL.md`                                                              |
| A React Native / Expo app for iOS or Android                                      | `references/react-native-client/SKILL.md`                                                       |
| Vanilla TS / Svelte / Vue / web worker (non-React)                                | `references/react-client/SKILL.md` → `ts-client-escape.md` (drop to `@fishjam-cloud/ts-client`) |
| Just calling REST directly (no SDK)                                               | `references/platform/rest-endpoints.md`                                                         |

A production Fishjam app always has **a backend (server SDK) + a client (client SDK)**. For prototyping you can skip the backend with the Sandbox API — see `references/platform/sandbox-vs-production.md`.

## Critical rules

- Management token **never** leaves the backend. If it leaks, regenerate it from the Dashboard.
- Peer tokens are valid for **24 hours from creation**. The token is consumed during the initial WS handshake; an established session keeps running on its own. If the peer hasn't connected yet or needs to _reconnect_ after 24h, call `refreshPeerToken` / `refresh_peer_token` to mint a new one.
- Sandbox API URLs are dev-only. Anyone holding one can create rooms on your account — never ship in a production client build.
