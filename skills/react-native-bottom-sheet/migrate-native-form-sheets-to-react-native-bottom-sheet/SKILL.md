---
name: migrate-native-form-sheets-to-react-native-bottom-sheet
description: >-
  Migrates React Native Screens routes configured with presentation: 'formSheet'
  in React Navigation native-stack or Expo Router Stack to
  @swmansion/react-native-bottom-sheet. Use for an explicit replacement or an
  assessment of that migration.
---

# Migrate native form sheets to React Native Bottom Sheet

Apply this workflow to a requested switch from a native form-sheet route to
Software Mansion's sheet. Keep native form-sheet creation, debugging, and
navigator upgrades in the existing API. For an assessment, report the mapping
and blockers without editing the app.

## Establish scope and compatibility

Find the affected route registrations and inherited options, resolved
`react-native-screens` and navigator or Expo Router versions, Expo SDK if
applicable, platforms, and intended destination release. Use the official docs
for those versions. Record platform-specific options and any doc gaps relevant
to the migration.

Check destination prerequisites before installing. Current guides require
React Native 0.82+, React 19.1.1+, and a Fabric native build on iOS and Android.
Use an Expo development build or bare native build, with pods installed for
bare iOS. Expo Go cannot load the destination native module. On web, check
`react-native-web` and `react-dom`. Report incompatibility and defer the swap
unless framework changes are already in scope.
[Destination setup](https://docs.swmansion.com/react-native-bottom-sheet/docs/getting-started/)
· [Destination releases](https://swmansion.com/changelog/react-native-bottom-sheet/)
· [Expo SDK versions](https://docs.expo.dev/versions/latest/#each-expo-sdk-version-depends-on-a-react-native-version)

## Preserve the route contract

The references describe the destination's documented 1.1 API. Preserve older
pins and confirm each used feature in version-matched docs or release notes
before applying these mappings; 1.0 changed overlay naming and dismissal.
[1.0 API changes](https://swmansion.com/changelog/react-native-bottom-sheet-1-0-0/) ·
[1.1 additions](https://swmansion.com/changelog/react-native-bottom-sheet-1-1-0/)

Read [Navigation and dismissal](references/navigation-and-dismissal.md) before
choosing a host or changing screen options. Inventory open and close callers,
route params and links, history and anchors, focus and cleanup, nested routes,
and back and removal guards. Keep the router responsible for those behaviors
when replacing the presentation. Convert the route to local UI only when the
user authorizes the change to navigation behavior.

Read [Detents and content](references/detents-and-content.md) to map heights,
initial selection, event consumers, scroll and keyboard behavior, insets,
dimming, accessibility, and sheet chrome. Flag any required system behavior
that the selected host and release cannot reproduce.

## Apply and review

When implementation is requested, migrate the route, host, and affected
callers as one change. Keep the existing navigator and its dependencies for
the app's navigation. Check the affected entry and exit paths against the
inventory within the requested verification scope. Report the migrated routes,
any platform or behavior differences, and which behaviors were verified in
the app.
