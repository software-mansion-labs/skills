<img width="1100" height="382" alt="skills_header" src="https://github.com/user-attachments/assets/f49ab3d0-159d-4169-9419-7dc34c09baf0" />

Production-ready patterns for React Native development, packaged as a [Claude Code plugin](https://code.claude.com/docs/en/plugins#create-plugins). Maintained by [Software Mansion](https://swmansion.com/).

Optimized for Claude models and tested with [Claude Code](https://claude.ai/code). Install the plugin and your AI coding agent gets up-to-date guidance for animations, gestures, on-device AI, audio, real-time video, and other React Native features.

## Installation

### As a plugin (recommended)

Add the Software Mansion marketplace and install the plugin:

```
/plugin marketplace add software-mansion-labs/skills
/plugin install skills@swmansion
/reload-plugins
```

The skills are available immediately. Run `/plugin marketplace update swmansion` to get the latest version.

### Via `npx`

You can also install the skills using the [`skills` CLI](https://www.npmjs.com/package/skills):

```bash
npx skills add software-mansion-labs/skills
```

## Available Skills

### [react-native-best-practices](./skills/react-native-best-practices/)

Production patterns for React Native apps on the New Architecture, covering:

| Topic              | What it covers                                                                                                                                                                                       |
| ------------------ | ---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| **Animations**     | CSS transitions, CSS animations, shared value animations, canvas animations (Skia), GPU shader animations, layout animations, Reanimated 4, 120fps, performance tuning                               |
| **Gestures**       | Tap, pan, pinch, swipe, long press, drag with Gesture Handler                                                                                                                                        |
| **SVG**            | Vector graphics, icons, charts, illustrations with React Native SVG                                                                                                                                  |
| **On-device AI**   | LLMs, computer vision, OCR, audio processing, embeddings with React Native ExecuTorch                                                                                                                |
| **Rich Text**      | Rich text editing with react-native-enriched and Markdown rendering with react-native-enriched-markdown: formatting toolbar, mentions, links, GFM tables, task lists, LaTeX math                     |
| **Multithreading** | Worker Runtimes, scheduling APIs, shared memory with React Native Worklets                                                                                                                           |
| **Audio**          | Playback, recording, visualization, session management with React Native Audio API                                                                                                                   |
| **JSI**            | C++ JavaScript Interface: `HostObject`, `HostFunction`, `NativeState`, zero-copy `ArrayBuffer`, threading safety, `CallInvoker`, TurboModules vs Nitro Modules, C++ memory patterns, crash debugging |

### [radon-mcp](./skills/radon-mcp/)

Best practices for using Radon IDE's MCP tools when developing, debugging, and inspecting React Native and Expo apps. Covers viewing screenshots, reading logs, inspecting the component tree, debugging network requests, reloading the app, and querying React Native documentation.

### [typegpu](./skills/typegpu/)

Practical guidance for building with TypeGPU, from project setup and schema design through shader authoring, resource management, and pipeline composition. It helps an AI coding agent work confidently across the CPU/GPU boundary while avoiding common typing, memory-layout, and WebGPU integration mistakes.

### [rnrepo](./skills/rnrepo/)

Best practices for integrating RNRepo — Software Mansion's infrastructure for pre-building and distributing React Native library artifacts. Covers installation (Expo CNG and standard React Native), configuration (denyList, disabling the plugin, GPG verification), and troubleshooting.

### [detour](./skills/detour/)

Best practices for setting up and migrating to Detour, Software Mansion's deferred deep linking ecosystem. Covers end-to-end SDK initialization, Universal/App Links registration, and type-safe analytics tracking across iOS, Android, React Native, and Flutter, as well as structural mappings for switching away from Branch or AppsFlyer.

### [fishjam](./skills/fishjam/)

Guidance for building real-time video, audio, and livestreaming apps with [Fishjam](https://fishjam.io), Software Mansion's hosted WebRTC platform. Covers the platform fundamentals (rooms, peers, tracks, two-tier auth, notifications, REST API) and all four SDKs: the Node.js and Python server SDKs (including AI voice agents and Gemini Live integration), the React web client, and the React Native / Expo client (permissions, foreground service, CallKit, screen sharing, Picture-in-Picture).

### [expo-horizon](./skills/expo-horizon/)

Software Mansion's guide for migrating Expo SDK apps to Meta Quest using the [expo-horizon](https://github.com/software-mansion-labs/expo-horizon) packages. Covers build flavors for Quest, panel sizing, device detection, and migrating `expo-location` and `expo-notifications` to their Horizon counterparts, through Meta Horizon Store publishing.

## Development

This project uses [Task](https://taskfile.dev) as a task runner. Install with `brew install go-task`.

```bash
task --list        # show available tasks
task lint          # run the skill linter
task lint:test     # run linter unit tests
task check         # run lint + tests
task eval:grade -- /path/to/workspace   # grade a skill-creator workspace
```

## Repository Structure

```
skills/
├── expo-horizon/
│   └── SKILL.md
├── fishjam/
│   ├── references/
│   │   ├── platform/
│   │   │   ├── SKILL.md
│   │   │   ├── auth-model.md
│   │   │   ├── glossary.md
│   │   │   ├── lifecycle-flow.md
│   │   │   ├── llms-and-docs.md
│   │   │   ├── notifications-taxonomy.md
│   │   │   ├── notifier-vs-webhook.md
│   │   │   ├── rest-endpoints.md
│   │   │   ├── room-types.md
│   │   │   └── sandbox-vs-production.md
│   │   ├── js-server-sdk/
│   │   │   ├── SKILL.md
│   │   │   ├── agent.md
│   │   │   ├── client.md
│   │   │   ├── express-fastify.md
│   │   │   ├── gemini-integration.md
│   │   │   ├── livestream-and-moq.md
│   │   │   ├── selective-subscriptions.md
│   │   │   ├── webhooks.md
│   │   │   └── ws-notifier.md
│   │   ├── python-server-sdk/
│   │   │   ├── SKILL.md
│   │   │   ├── agent.md
│   │   │   ├── client.md
│   │   │   ├── fastapi.md
│   │   │   ├── gemini-integration.md
│   │   │   ├── livestream-and-moq.md
│   │   │   ├── notifier.md
│   │   │   ├── selective-subscriptions.md
│   │   │   └── webhooks.md
│   │   ├── react-client/
│   │   │   ├── SKILL.md
│   │   │   ├── connection.md
│   │   │   ├── custom-sources.md
│   │   │   ├── data-and-events.md
│   │   │   ├── devices.md
│   │   │   ├── livestream.md
│   │   │   ├── peers-and-tracks.md
│   │   │   ├── provider.md
│   │   │   ├── simulcast-and-bandwidth.md
│   │   │   └── ts-client-escape.md
│   │   └── react-native-client/
│   │       ├── SKILL.md
│   │       ├── audio-output.md
│   │       ├── callkit.md
│   │       ├── example-projects.md
│   │       ├── foreground-service.md
│   │       ├── mobile-hook-overrides.md
│   │       ├── native-setup.md
│   │       ├── permissions.md
│   │       ├── picture-in-picture.md
│   │       ├── rtcview.md
│   │       └── screen-sharing.md
│   └── SKILL.md
├── radon-mcp/
│   ├── references/
│   │   ├── get-library-description.md
│   │   ├── query-documentation.md
│   │   ├── reload-application.md
│   │   ├── view-application-logs.md
│   │   ├── view-component-tree.md
│   │   ├── view-network-logs.md
│   │   ├── view-network-request-details.md
│   │   └── view-screenshot.md
│   └── SKILL.md
├── rnrepo/
│   ├── references/
│   │   ├── configuration.md
│   │   ├── installation.md
│   │   └── troubleshooting.md
│   └── SKILL.md
├── typegpu/
│   ├── references/
│   │   ├── advanced.md
│   │   ├── matrices.md
│   │   ├── noise.md
│   │   ├── pipelines.md
│   │   ├── sdf.md
│   │   ├── setup.md
│   │   ├── shaders.md
│   │   ├── textures.md
│   │   └── types.md
│   └── SKILL.md
├── detour/
│   ├── README.md
│   ├── migrate-to-detour/
│   │   ├── SKILL.md
│   │   └── references/
│   │       ├── android.md
│   │       ├── ios.md
│   │       ├── react-native.md
│   │       └── flutter.md
│   └── detour-onboarding/
│       ├── SKILL.md
│       └── references/
│           ├── android.md
│           ├── ios.md
│           ├── react-native.md
│           └── flutter.md
└── react-native-best-practices/
    ├── references/
    │   ├── animations/
    │   │   ├── SKILL.md
    │   │   ├── animation-functions.md
    │   │   ├── animations-performance.md
    │   │   ├── animations.md
    │   │   ├── canvas-animations.md
    │   │   ├── canvas-atlas.md
    │   │   ├── gpu-animations.md
    │   │   ├── layout-animations.md
    │   │   ├── scroll-and-events.md
    │   │   └── svg-animations.md
    │   ├── audio/
    │   │   ├── SKILL.md
    │   │   ├── audio.md
    │   │   ├── effects-and-analysis.md
    │   │   ├── playback.md
    │   │   ├── recording.md
    │   │   ├── system-and-notifications.md
    │   │   └── worklets.md
    │   ├── gestures/
    │   │   ├── SKILL.md
    │   │   ├── continuous-gestures.md
    │   │   ├── gesture-composition.md
    │   │   ├── gestures.md
    │   │   ├── swipeable-and-drawer.md
    │   │   ├── tap-handling.md
    │   │   └── testing.md
    │   ├── multithreading/
    │   │   ├── SKILL.md
    │   │   ├── setup-and-advanced.md
    │   │   ├── shared-memory.md
    │   │   └── threading-api.md
    │   ├── on-device-ai/
    │   │   ├── SKILL.md
    │   │   ├── llm.md
    │   │   ├── setup.md
    │   │   ├── speech.md
    │   │   └── vision.md
    │   ├── rich-text/
    │   │   └── SKILL.md
    │   ├── svg/
    │   │   ├── SKILL.md
    │   │   ├── svg.md
    │   │   └── when-to-use.md
    │   └── jsi/
    │       ├── SKILL.md
    │       ├── overview.md
    │       ├── core-types.md
    │       ├── casting-and-serialization.md
    │       ├── threading-safety.md
    │       ├── calling-js-and-async.md
    │       ├── performance.md
    │       ├── setup-and-templates.md
    │       ├── module-approaches.md
    │       ├── cpp-memory-patterns.md
    │       └── debugging-and-pitfalls.md
    ├── README.md
    └── SKILL.md
```

The top-level `SKILL.md` acts as a table of contents. Reference files load only when relevant to the current task, keeping the context window focused.

## Contributing

### Adding a new sub-skill

1. Create a directory under `skills/react-native-best-practices/references/<topic>/`
2. Add a `SKILL.md` with YAML frontmatter (`name`, `description`) and a references table
3. Place detailed documentation in sibling `.md` files
4. Register the sub-skill in `skills/react-native-best-practices/SKILL.md`

### Writing guidelines

- Keep `SKILL.md` files under 500 lines. Split longer content into reference files.
- Write descriptions that include trigger terms (what phrases should activate this skill).
- Use concrete code examples over abstract explanations.
- Target the New Architecture and current library versions.

## License

MIT

## Skills are created by Software Mansion

Since 2012 [Software Mansion](https://swmansion.com) is a software agency with experience in building web and mobile apps. We are Core React Native Contributors and experts in dealing with all kinds of React Native issues. We can help you build your next dream product – [Hire us](https://swmansion.com/contact?utm_source=skills&utm_medium=readme).
