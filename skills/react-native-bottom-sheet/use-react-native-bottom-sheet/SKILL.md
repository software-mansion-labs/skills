---
name: use-react-native-bottom-sheet
description: Integrates and troubleshoots @swmansion/react-native-bottom-sheet. Use when the request explicitly concerns Software Mansion's React Native Bottom Sheet or when the affected code imports this package.
---

# Use React Native Bottom Sheet

Confirm the package from the affected import and app dependencies when the
request names only a bottom sheet. Keep fixes for other packages in their
existing API. If the package remains unclear, ask which library is intended.
Package replacement requires an explicit request.

## Resolve the version

Find the resolved version in the importing app's lockfile and keep it unless
the task calls for a change. For a new dependency, check the published release
and the app's React Native, React, Expo, and platform requirements before
choosing a version. Use Software Mansion's official docs for the public API
and release notes to check when behavior was introduced or changed. State any
uncertainty when the published guides cannot establish older release behavior.
[Official docs](https://docs.swmansion.com/react-native-bottom-sheet/docs/) ·
[Product changelog](https://swmansion.com/changelog/react-native-bottom-sheet/)

## Load the relevant reference

For a new integration, read Setup and Usage. For other requests, load the
references that address the task.

- [Setup](references/setup.md) for requirements, Expo, bare React Native, and web.
- [Usage](references/usage.md) for inline and modal placement, context, controlled
  detents, surfaces, dismissal, and motion events.
- [Interactions](references/interactions.md) for scroll and gesture handoff,
  keyboard avoidance, and platform input and accessibility.
- [Troubleshooting](references/troubleshooting.md) for install, layout, navigation,
  or interaction symptoms.

Check the behavior that matters against the upstream page linked in the
reference. Report the resolved version when known and distinguish documented
behavior from any result verified in the app.
