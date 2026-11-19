---
name: migrate-true-sheet-to-react-native-bottom-sheet
description: Migrates @lodev09/react-native-true-sheet to @swmansion/react-native-bottom-sheet. Use for an explicit request to make or assess this one-way replacement.
---

# Migrate True Sheet to React Native Bottom Sheet

Apply this workflow only to a requested True Sheet replacement with Software
Mansion's library. An import, source bug report, or True Sheet version upgrade
leaves the existing package in place. Reverse migrations and other source
libraries have their own contracts. For an assessment, report the mappings
and blockers without editing the app.

## Establish the migration boundary

Find the affected app workspace, resolved True Sheet version, intended
destination release, platforms, and shared sheet wrappers. The current True
Sheet docs cover v4. Use the v3 docs for legacy code, especially scrolling,
layout, navigation, and transition timing. Confirm older behavior in
version-matched official docs before applying a mapping.
[True Sheet v4 docs](https://sheet.lodev09.com/install) ·
[True Sheet v3 docs](https://sheet.lodev09.com/v3/intro) ·
[v3 to v4 changes](https://sheet.lodev09.com/migration)

Check the destination release against the app's React Native, React, Expo, and
platform versions. The current guides require React Native 0.82+ and React
19.1.1+. iOS and Android require Fabric and a native build; Expo Go cannot
run the destination. Check bare iOS pods and React Native Web
dependencies as applicable. If requirements conflict, report the blocker and
defer the swap unless the framework upgrade is already part of the request.
[Destination setup](https://docs.swmansion.com/react-native-bottom-sheet/docs/getting-started/) ·
[Destination releases](https://swmansion.com/changelog/react-native-bottom-sheet/)

Inventory package imports and subpaths, providers, named methods, ref callers,
and wrappers. For each affected sheet, record its opening detent and payload,
completion side effects, mount or reset behavior, owning screen, and relied-on
dismissal, stacking, keyboard, and scrolling behavior.

## Map the used behavior

The references describe the destination's documented 1.1 API. Preserve older
pins and confirm each used feature in version-matched docs or release notes
before applying these mappings; 1.0 changed overlay naming and dismissal.
[1.0 API changes](https://swmansion.com/changelog/react-native-bottom-sheet-1-0-0/) ·
[1.1 additions](https://swmansion.com/changelog/react-native-bottom-sheet-1-1-0/)

- Read [State and sizing](references/state-and-sizing.md) for commands, detents,
  transition timing, content measurement, and motion events.
- Read [Presentation and interactions](references/presentation-and-interactions.md)
  for host selection, navigation, overlay ordering, dismissal, scrolling,
  keyboard handling, and appearance.

For each used feature, give the documented mapping, the adaptation the app
would own, or the unresolved gap. Identify which gaps block a direct
replacement. Preserve wrapper contracts where feasible and review callers
together when timing or behavior changes.

## Apply within scope

For implementation, use the app's package manager and destination setup. A
coherent group of sheets and their callers can migrate while others retain
True Sheet. Remove the source package, imports, providers, and package-specific
mocks or patches after checking remaining consumers. True Sheet web can use
Radix peers, and sheet navigation and motion can use additional dependencies.
Keep dependencies and configuration still used elsewhere.
[True Sheet web](https://sheet.lodev09.com/guides/web) ·
[True Sheet navigation](https://sheet.lodev09.com/guides/navigation) ·
[True Sheet Reanimated](https://sheet.lodev09.com/guides/reanimated)

Report the affected sheets, material behavior changes, blockers, and what was
verified in the app. For assessment alone, identify the work needed to preserve
each required behavior and label conclusions drawn from documentation.
