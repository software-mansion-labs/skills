# Upstream Resources

When the skill references aren't enough or have drifted out of date, go to these sources. Documentation is the **source of truth** — examples and SDK source are supplementary.

## Documentation

- **Docs site:** <https://fishjam.swmansion.com/docs>
- **llms.txt:** <https://fishjam.swmansion.com/docs/llms.txt> — machine-readable index of all doc pages
- **llms-full.txt:** <https://fishjam.swmansion.com/docs/llms-full.txt> — same index plus the full page bodies inlined

Useful sub-trees of the docs:

- `/docs/explanation/` — concept-level pages (what each thing is and why)
- `/docs/tutorials/` — guided walkthroughs (`backend-quick-start`, `react-quick-start`, `react-native-quick-start`, `livestreaming`, `agents`, `moq`)
- `/docs/how-to/backend/` — recipes for the server side (`server-setup`, `fastapi-example`, `fastify-example`, `production-deployment`, `selective-subscriptions`, `whip-whep`, etc.)
- `/docs/how-to/client/` — recipes for the client side (`installation`, `connecting`, `managing-devices`, `screensharing`, `picture-in-picture`, `background-streaming`, `reconnection-handling`, `stream-middleware`, etc.)
- `/docs/api/` — generated TypeDoc / pdoc API references
- `/docs/integrations/` — third-party integrations (`gemini-live-integration`, `vapi-integration`, etc.)

## Dashboard

- **App:** <https://fishjam.io/app> — view and manage your Fishjam instances; get Fishjam ID + management token; provision sandbox URLs; manage billing.
- **Sandbox tab:** <https://fishjam.io/app/sandbox> — enable / disable / rotate the Sandbox API URL.

## API specifications

- **OpenAPI YAML (REST):**
  - Source (browsable): <https://github.com/fishjam-cloud/documentation/blob/main/static/api/fishjam-server-openapi.yaml>
  - Source (raw YAML): <https://raw.githubusercontent.com/fishjam-cloud/documentation/main/static/api/fishjam-server-openapi.yaml>
  - Served: <https://fishjam.swmansion.com/docs/api/fishjam-server-openapi.yaml>
  - Rendered (Scalar UI): <https://fishjam.swmansion.com/docs/api/rest>

- **Protobuf (WebSocket + webhook events):**
  - Source-of-truth repo: <https://github.com/fishjam-cloud/protos>
  - `server_notifications.proto` (browsable): <https://github.com/fishjam-cloud/protos/blob/main/fishjam/server_notifications.proto>
  - `agent_notifications.proto`: not in `fishjam-cloud/protos`; shipped only via the documentation repo at <https://github.com/fishjam-cloud/documentation/blob/main/static/api/protobuf/agent_notifications.proto>

## SDK source repositories

- **JS server SDK:** <https://github.com/fishjam-cloud/js-server-sdk> (package: `@fishjam-cloud/js-server-sdk`)
- **Python server SDK:** <https://github.com/fishjam-cloud/python-server-sdk> (package: `fishjam-server-sdk` on PyPI; `import fishjam`)
- **Web + Mobile client SDK monorepo:** <https://github.com/fishjam-cloud/web-client-sdk>
  - `@fishjam-cloud/ts-client` — core TS WebRTC client
  - `@fishjam-cloud/react-client` — React hooks
  - `@fishjam-cloud/react-native-client` — React Native / Expo wrapper (source at `packages/mobile-client/`)
  - `@fishjam-cloud/webrtc-client` — low-level WebRTC primitives
  - `@fishjam-cloud/react-native-webrtc` — fork of `react-native-webrtc` with native modules

## When to use which

| If you need… | Look at |
| --- | --- |
| What something is, conceptually | docs `/explanation/*` |
| How to build something end-to-end | docs `/tutorials/*` |
| How to do one specific thing | docs `/how-to/*` |
| Exact method signature, types | docs `/api/*` (or the SDK source) |
| Wire format, raw HTTP behavior | OpenAPI YAML, protobuf file |
| What changed in a recent release | upstream SDK GitHub releases |
| Working example app | Each SDK's own `examples/` directory in its GitHub repo |
