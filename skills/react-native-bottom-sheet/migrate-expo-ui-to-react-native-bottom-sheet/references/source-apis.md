# Identify the Expo UI sheet API

Read the reference for the app's SDK and actual import. Software Mansion's
comparison uses the universal `BottomSheet`.
[Library comparison](https://docs.swmansion.com/react-native-bottom-sheet/docs/how-it-compares/)

## Universal BottomSheet from `@expo/ui`

SDK 56, 57, and 58 document a modal sheet on Android, iOS, and web, included
in Expo Go. It takes `isPresented` and `onDismiss`. User dismissal calls
`onDismiss`, where the app clears its open state. Programmatic close changes
`isPresented`.

Without `snapPoints` it sizes to content. Use the source's actual opening
height and reachable stops. `'half'` and `'full'` are semantic sizes. iOS
and web honor `{ fraction }` and `{ height }`. Android maps those to the
nearest half/full state, and the half state requires sufficiently tall content.

The SDK 56 example places the trigger and sheet under `Host`. The SDK 57
and 58 examples wrap the Expo UI trigger in `Host` beside the sheet.
SDK 58 documents `contentPadding`, color controls, and Android dismissal
flags beyond the public props listed in SDK 56 and 57.
[SDK 56](https://docs.expo.dev/versions/v56.0.0/sdk/ui/universal/bottomsheet/) ·
[SDK 57](https://docs.expo.dev/versions/v57.0.0/sdk/ui/universal/bottomsheet/) ·
[SDK 58](https://docs.expo.dev/versions/v58.0.0/sdk/ui/universal/bottomsheet/) ·
[Snap points](https://docs.expo.dev/versions/v58.0.0/sdk/ui/universal/bottomsheet/#snap-points)

SDK 58 Android defaults Back and outside-tap dismissal to enabled.
`shouldDismissOnBackPress={false}` also keeps the press from React Native's
`BackHandler`. `shouldDismissOnClickOutside` controls scrim dismissal
separately. `showDragIndicator` controls handle visibility.
[Back](https://docs.expo.dev/versions/v58.0.0/sdk/ui/universal/bottomsheet/#shoulddismissonbackpress) ·
[Outside tap](https://docs.expo.dev/versions/v58.0.0/sdk/ui/universal/bottomsheet/#shoulddismissonclickoutside) ·
[Indicator](https://docs.expo.dev/versions/v58.0.0/sdk/ui/universal/bottomsheet/#showdragindicator)

## SwiftUI BottomSheet from `@expo/ui/swift-ui`

This primitive is documented for iOS and tvOS, including in SDK 55 and 56.
Visibility uses `isPresented` and `onIsPresentedChange(boolean)`. The SDK 56
and 58 references document an optional `onDismiss` for full dismissal. SDK 58
also documents an optional `anchor`. Move its opener into the owning screen
when replacing the sheet. SDK 55 and 56 examples place the opener beside the
sheet under the SwiftUI host.

`fitToContents` opts into content sizing. `presentationDetents` on sheet
content uses `'medium'`, `'large'`, `{ fraction }`, or `{ height }` and can
track a selected detent with `selection` and `onSelectionChange`. Preserve
the selected opening detent independently of visibility. Inspect
`interactiveDismissDisabled` for swipe dismissal and
`presentationBackgroundInteraction` for interaction behind the sheet.
[SDK 55](https://docs.expo.dev/versions/v55.0.0/sdk/ui/swift-ui/bottomsheet/) ·
[SDK 56](https://docs.expo.dev/versions/v56.0.0/sdk/ui/swift-ui/bottomsheet/) ·
[SDK 58](https://docs.expo.dev/versions/v58.0.0/sdk/ui/swift-ui/bottomsheet/) ·
[Detent selection](https://docs.expo.dev/versions/v58.0.0/sdk/ui/swift-ui/bottomsheet/#bottom-sheet-with-detent-selection-tracking) ·
[Background interaction](https://docs.expo.dev/versions/v58.0.0/sdk/ui/swift-ui/bottomsheet/#bottom-sheet-with-background-interaction)

## Compose ModalBottomSheet from `@expo/ui/jetpack-compose`

This Android primitive is documented in SDK 55, 56, and 58. It mounts
conditionally. `onDismissRequest` handles user dismissal. The documented
programmatic path awaits `ref.hide()` before unmounting. Identify callers
of `expand()` and `partialExpand()` and any continuations after the promises
resolve. `skipPartiallyExpanded` removes the partial stop. In SDK 58,
`initialFullyExpanded` opens at expanded while leaving the partial stop
reachable. Check which stops and heights the source content actually shows.
`properties` controls Back and outside taps independently, and
`sheetGesturesEnabled` controls swipe gestures.
[SDK 55](https://docs.expo.dev/versions/v55.0.0/sdk/ui/jetpack-compose/bottomsheet/) ·
[SDK 56](https://docs.expo.dev/versions/v56.0.0/sdk/ui/jetpack-compose/bottomsheet/) ·
[SDK 58](https://docs.expo.dev/versions/v58.0.0/sdk/ui/jetpack-compose/bottomsheet/) ·
[Ref methods](https://docs.expo.dev/versions/v58.0.0/sdk/ui/jetpack-compose/bottomsheet/#modalbottomsheetref) ·
[Dismissal controls](https://docs.expo.dev/versions/v58.0.0/sdk/ui/jetpack-compose/bottomsheet/#non-dismissible-sheet)

## Community adapter from `@expo/ui/community/bottom-sheet`

Expo's Gorhom-compatible adapter presents modally, including its default
`BottomSheet` export. If this import is in scope, map its `-1` closed state
and supported ref commands to a destination zero-height detent and controlled
index, and map `onClose` or `onDismiss` to full-close completion. Inspect
the adapter's SDK-specific compatibility notes and platform behavior,
especially its combined `enablePanDownToClose` effects on scrim, Back, and
Escape. Some props are accepted solely for compatibility.
[SDK 58 adapter](https://docs.expo.dev/versions/v58.0.0/sdk/ui/drop-in-replacements/bottomsheet/#platform-behavior) ·
[Supported exports and compatibility](https://docs.expo.dev/versions/v58.0.0/sdk/ui/drop-in-replacements/bottomsheet/#supported-exports) ·
[SDK 57 adapter](https://docs.expo.dev/versions/v57.0.0/sdk/ui/drop-in-replacements/bottomsheet/) ·
[SDK 56 adapter](https://docs.expo.dev/versions/v56.0.0/sdk/ui/drop-in-replacements/bottomsheet/)
