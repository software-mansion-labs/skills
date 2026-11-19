---
name: migrate-gorhom-to-react-native-bottom-sheet
description: Migrates @gorhom/bottom-sheet to @swmansion/react-native-bottom-sheet. Use for this replacement or for assessing how existing Gorhom sheets map to Software Mansion's library.
---

# Migrate Gorhom to React Native Bottom Sheet

Use this workflow for the Gorhom-to-Software-Mansion direction. A Gorhom bug
report or import alone does not request replacement. Keep Gorhom version
upgrades, reverse migrations, and other source libraries outside this workflow.
For an assessment, produce the mapping and blockers without editing the app.

## Establish the migration boundary

Find the affected app workspace, resolved Gorhom version, intended destination
release, platforms, and shared sheet wrappers. The source references below
describe Gorhom v5; for v4 or older code, consult that version's official docs
before applying a mapping. Use each library's official documentation for its
public contract and release notes for version-sensitive changes. State gaps
or conflicting documentation instead of inventing API parity.

Check the selected release against its official requirements. Current guides
require React Native 0.82+ and React 19.1.1+; iOS/Android need Fabric and a
native build, so Expo Go is insufficient. Check Expo SDK alignment and web
dependencies where relevant. If prerequisites are incompatible, report the
blocker and defer the package swap. A sheet migration does not authorize a
framework upgrade.
[Destination setup](https://docs.swmansion.com/react-native-bottom-sheet/docs/getting-started/)
· [Expo SDK versions](https://docs.expo.dev/versions/latest/#each-expo-sdk-version-depends-on-a-react-native-version)
· [Destination releases](https://swmansion.com/changelog/react-native-bottom-sheet/)

Inventory imports, re-exports, refs, hooks, providers, and every caller of a
shared wrapper. Record the intended height behind each index, opening payloads,
close/reset side effects, mounted-state expectations, and interaction features
actually used. Include mocks and package-specific patches in the affected scope.

## Map behavior before replacing imports

- Read [State and sizing](references/state-and-sizing.md) for every migration:
  explicit closed detents, dynamic sizing, imperative commands, callbacks,
  motion consumers, and a fixed-height before/after example.
- Read [Presentation and interactions](references/presentation-and-interactions.md)
  for modal lifecycle, context, stacking, scrims, accessibility, scrolling,
  keyboard handling, or custom chrome used by the app.

For each used feature, identify a documented destination API, an app-owned
adaptation, or an unresolved gap. Resolve behavior-changing choices with the
user when existing requirements cannot determine them. Preserve public wrapper
contracts where feasible and adapt their internals and callers together.

## Apply and finish within scope

When implementation is requested, install the selected destination release
using the app's package manager and follow its native setup. Migrate a coherent
sheet or shared wrapper and its consumers before continuing. Both packages can
remain dependencies while other sheets still use Gorhom.

Remove Gorhom only after the intended workspace has no remaining consumers.
Keep Reanimated, Gesture Handler, their configuration, and root wrappers while
other code needs them; the destination's independence from those libraries is
not evidence that the app can remove them.
[Library comparison](https://docs.swmansion.com/react-native-bottom-sheet/docs/how-it-compares/)

Review the affected callers against the behavior inventory. Use existing app
checks only within the requested scope; documentation-only work needs no
scaffolded app or runtime execution. Report migrated and remaining sheets,
unresolved behavior gaps, documentation consulted, and checks actually run.
Separate documented behavior from behavior verified in the app.
