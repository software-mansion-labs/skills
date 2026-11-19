---
name: migrate-expo-ui-to-react-native-bottom-sheet
description: Migrates Expo UI BottomSheet to @swmansion/react-native-bottom-sheet. Use for an explicit request to make or assess this one-way replacement, including Expo UI's platform-specific sheets when used by the affected flow.
---

# Migrate Expo UI to React Native Bottom Sheet

Apply this workflow to requests to replace or assess replacing Expo UI sheets
with Software Mansion's library. Keep the existing package for bug reports,
integration questions, and version updates. For an assessment, inspect the app
and report the required changes and blockers without editing it.

## Identify the source and destination

Find the affected workspace, resolved `@expo/ui` version, Expo SDK, destination
release, and target platforms. Follow sheet aliases and re-exports to their
original imports. Read [Source APIs](references/source-apis.md) to identify
the universal, SwiftUI, Compose, or community API and consult its SDK-specific
documentation. Use official release docs for both public contracts.

Check the destination release against the app. The current requirements are
React Native 0.82+, React 19.1.1+, and Fabric with a native build on iOS and
Android. Native use needs a development or native build instead of Expo Go;
an existing development client must include the new native package. Check pods
for bare iOS and React Native Web dependencies for web.
If requirements conflict, report the dependency or platform blocker and defer
the incompatible replacement unless the needed upgrade is already authorized.
[Destination setup](https://docs.swmansion.com/react-native-bottom-sheet/docs/getting-started/) ·
[Release notes](https://swmansion.com/changelog/react-native-bottom-sheet/) ·
[Expo SDK versions](https://docs.expo.dev/versions/latest/#each-expo-sdk-version-depends-on-a-react-native-version) ·
[Rebuilding a development client](https://docs.expo.dev/develop/development-builds/introduction/#rebuild-when-native-code-changes)

## Map and migrate the affected flow

The references describe the destination's documented 1.1 API. Preserve older
pins and confirm each used feature in version-matched docs or release notes
before applying these mappings; 1.0 changed overlay naming and dismissal.
[1.0 API changes](https://swmansion.com/changelog/react-native-bottom-sheet-1-0-0/) ·
[1.1 additions](https://swmansion.com/changelog/react-native-bottom-sheet-1-1-0/)

Inventory each sheet's trigger, opening size, owning screen, host and content
wrappers, dismissal paths, completion side effects, and state that survives
closing. Include shared callers and platform branches.

Read [Migration](references/migration.md) for presentation, controlled detents,
event timing, content and accessibility. For each behavior the app relies on,
identify its documented equivalent, the adaptation the app must own, or the
gap. Use the app's existing package manager and preserve shared wrapper
contracts where possible.

Remove obsolete sheet imports, mocks, and wrappers after checking their
callers. Remove `@expo/ui` from the affected app only when no other uses
remain, including shared components, platform branches, and subpath imports
for hosts, modifiers, or other Expo UI components.

Within the requested verification scope, check opening size, scroll and keyboard behavior,
dismissal and completion, state on reopening, and accessibility. Where sheets
stack, check top-sheet ordering and Back behavior. Report the migrated
callsites, remaining gaps, and which results were verified in the app.
Distinguish documentation findings from app results; documentation-only work
does not require an app build or runtime checks.
