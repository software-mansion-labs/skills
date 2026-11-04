---
name: on-device-ai
description: >-
  Routes on-device AI tasks in React Native and Expo apps to the react-native-executorch skill.
  Use when building offline chat, vision, speech, embeddings or retrieval features, choosing
  pre-exported models, or integrating custom ExecuTorch models and pipelines.
---

# On-device AI with React Native ExecuTorch

Use the `react-native-executorch` skill for app setup, model selection, inference, resource lifetime,
custom models and troubleshooting. It covers both ready-made task hooks with pre-exported models
and custom PyTorch models with TypeScript pipelines.

If the skill is not available, install it from the app's project directory, selecting the coding agent
in use when prompted:

```bash
npx skills add software-mansion-labs/skills -s react-native-executorch
```

Then read its `SKILL.md` and the references relevant to the task before implementing the integration.

For direct guidance, read
[Getting Started](https://docs.swmansion.com/react-native-executorch/docs/fundamentals/getting-started.md)
and use documentation matching the app's installed version. Fetch the relevant task and API pages as
needed. Follow [ExecuTorch Export Recipes](https://github.com/software-mansion-labs/executorch-export-recipes)
for exporting custom models. Verify the app's runtime and backend compatibility before integration.
