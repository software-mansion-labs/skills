# Migrate presentation and interaction

## Choose the destination host

Use `ModalBottomSheet` with `BottomSheetProvider` for the usual Expo UI modal
flow. Use `windowOverlay` when the sheet must cover a native modal screen.
It keeps call-site React context and can run without the provider. Inline
`BottomSheet` fills its parent and has no scrim. Choose it for a sheet that
belongs to part of the screen.

A modal portal reads context from above `BottomSheetProvider`. Capture
theme, form, and navigation values or callbacks at the owner and pass them
into the sheet when those providers sit below the portal provider. Expo UI's
SwiftUI and universal sheets document nesting on iOS for stacked sheets.
Destination modal stacking follows open order within a presentation layer,
and window overlays appear above portals. Preserve each sheet's owning
state when changing its placement.
[Provider setup](https://docs.swmansion.com/react-native-bottom-sheet/docs/getting-started/#installation) ·
[Window overlay](https://docs.swmansion.com/react-native-bottom-sheet/docs/modal-sheets/#window-overlay) ·
[React context](https://docs.swmansion.com/react-native-bottom-sheet/docs/modal-sheets/#react-context) ·
[Stacking](https://docs.swmansion.com/react-native-bottom-sheet/docs/modal-sheets/#stacking) ·
[Inline sheets](https://docs.swmansion.com/react-native-bottom-sheet/docs/inline-sheets/) ·
[Expo iOS nesting](https://docs.expo.dev/versions/v58.0.0/sdk/ui/universal/bottomsheet/)

## Preserve sizing and event timing

- Reserve a zero-height detent and choose an explicit opening index.
  Use `[0, 'content']` for intrinsic content. Start with `'50%'` for half or
  medium and `'100%'` for full or large. Convert fractions to percentage
  strings and fixed heights to layout-unit numbers. Destination percentages
  use available host height after insets, so compare actual resting heights
  and keep only the Android states the source made reachable. Map Compose
  expand and partial-expand to their open indices. Keep resolved detents in
  ascending order and change the index in the same render if its old value
  would be out of range.
- Set `index` for app-driven open, close, and snap actions. Update it from
  `onIndexChange` when a user snap commits or default modal dismissal starts.
  An app-driven change skips that callback. Preserve user-only Expo
  dismissal effects on a user-triggered transition to closed. If a custom
  close request handler accepts dismissal and sets the index, handle that
  user path there as well. Match completion work to the relevant target and
  pending action in `onSettle`, including awaited expand/partial-expand work.
  Match the pending close when moving SwiftUI dismissal completion,
  Compose's awaited hide continuation, or the community adapter's closing
  callbacks. Run the continuation only when its close settles. Discard it if
  the app cancels that close by reopening.
- Keep the destination sheet mounted through closing. Children and surface
  stay mounted while closed by default. Use `lazy` when content should
  unmount after settling closed and reset on reopening. Keep state that
  must survive in the owner.
  [Detents and geometry](https://docs.swmansion.com/react-native-bottom-sheet/docs/detents-and-index/) ·
  [Index events](https://docs.swmansion.com/react-native-bottom-sheet/docs/detents-and-index/#index-events) ·
  [Lazy content](https://docs.swmansion.com/react-native-bottom-sheet/docs/modal-sheets/#lazy-content)

## Move content and styling

Unwrap `RNHostView` around React Native content when it only bridged the old
sheet boundary. Replace native SwiftUI or Compose components with React Native
layout and controls where appropriate. If an Expo UI trigger or native content
island remains, retain its matching `Host` and platform scope.

Use a filling `surface` for the sheet background and `surfaceStyle` for
clipping, corner radii, and shadows. Apply content padding, safe-area padding,
text colors, and any visual drag indicator in React Native. Match scrim color
and opacity explicitly where they affect the experience. Review native sheet
materials, system chrome, and adaptive appearance for any required
app-provided equivalents. SwiftUI documents tvOS, and the destination lists
iOS, Android, and web. Treat an affected tvOS path as a compatibility gap.
[Expo Host](https://docs.expo.dev/versions/v58.0.0/sdk/ui/universal/host/) ·
[Expo RNHostView](https://docs.expo.dev/versions/v58.0.0/sdk/ui/universal/rnhostview/) ·
[SwiftUI sheet and appearance](https://docs.expo.dev/versions/v58.0.0/sdk/ui/swift-ui/bottomsheet/) ·
[Destination surface](https://docs.swmansion.com/react-native-bottom-sheet/docs/surface/) ·
[Modal safe-area example](https://docs.swmansion.com/react-native-bottom-sheet/docs/modal-sheets/) ·
[Destination platforms](https://docs.swmansion.com/react-native-bottom-sheet/docs/)

Use regular React Native scrollables. Bound a filling list with a fixed or
percentage upper detent. A filling list can make a `'content'` detent take
the available height. Expansion uses same-drag handoff by default, and
collapse uses the scroll position at gesture start. Choose collapse
`'handoff'` where the previous flow relied on one drag continuing into the
sheet at the list edge. Check nested-scrolling support for Android
third-party scrollers and overflow on web.

The app supplies keyboard avoidance. Use host resizing, a keyboard-aware
scrollable, or bottom padding as the layout needs. With animated padding
on measured content, `animateContentHeight={false}` lets a content-sized
sheet follow the changes immediately.
[Expo scrollable content](https://docs.expo.dev/versions/v58.0.0/sdk/ui/universal/bottomsheet/#scrollable-react-native-content) ·
[Scroll coordination](https://docs.swmansion.com/react-native-bottom-sheet/docs/scroll-coordination/) ·
[Keyboard handling](https://docs.swmansion.com/react-native-bottom-sheet/docs/keyboard-handling/)

## Preserve dismissal and accessibility

Default modal dismissal uses positive current scrim opacity and a
non-programmatic zero-height detent. A transparent scrim color can remain
actionable. `onCloseRequest` lets the app accept or decline eligible scrim,
screen-reader, Android Back, and Android/web physical Escape requests. The
documented handler has no source discriminator. Review independent Expo Back
and outside-tap flags individually. Without a handler or default closure,
Back and Escape can continue to the host. A handler that declines the
request consumes eligible Back and Escape while leaving the sheet open.
Zero scrim opacity removes the scrim dismissal target at that position and
changes screen-reader dismissal availability. It can still leave Back and
Escape actionable with a handler.

Dragging to a closed detent reports through `onIndexChange` directly.
For an unsaved-edit guard, mark every zero-height detent
`programmatic: true` so drag cannot commit to one. A blocked swipe does
not call `onCloseRequest`. Route explicit close controls and eligible
close requests through the app's guard and set the closed index on acceptance.
This preserves the open detents for dragging. If the source disabled all
sheet gestures or used SwiftUI background interaction, these guides establish
no direct equivalent. Keep those requirements as gaps until an app-owned
adaptation or an accepted behavior change resolves them.

Provide a reachable close control. Check focus and restoration, background
interaction, and scrim dismissal with VoiceOver, TalkBack, and web keyboard
input as applicable. Use `scrimAccessibilityLabel` for a localized scrim
label. VoiceOver escape and TalkBack's Dismiss action can remain available at
full height when opacity and actionability conditions hold. iOS physical
Escape is outside the documented close request paths.
[Close requests and routing](https://docs.swmansion.com/react-native-bottom-sheet/docs/modal-sheets/#close-requests) ·
[Accessible dismissal](https://docs.swmansion.com/react-native-bottom-sheet/docs/modal-sheets/#accessible-dismissal) ·
[Programmatic-only detents](https://docs.swmansion.com/react-native-bottom-sheet/docs/detents-and-index/#programmatic-only-detents)

## Universal sheet with a React Native list

This example assumes the source opens at half height and can expand to full
height. The fragments share `listProps` for the app's React Native `FlatList`,
including its data and renderer. Preserve the existing trigger and close
controls in the actual app.

Before, with Expo UI's sheet and React Native list bridge

```tsx
import { BottomSheet, RNHostView } from '@expo/ui';
import { useState } from 'react';
import { FlatList } from 'react-native';

const [isPresented, setIsPresented] = useState(false);
// The opener calls setIsPresented(true).

<BottomSheet
  isPresented={isPresented}
  onDismiss={() => setIsPresented(false)}
  snapPoints={['half', 'full']}
>
  <RNHostView>
    <FlatList {...listProps} nestedScrollEnabled style={{ flex: 1 }} />
  </RNHostView>
</BottomSheet>
```

After, under an existing `BottomSheetProvider`

```tsx
import { ModalBottomSheet } from '@swmansion/react-native-bottom-sheet';
import { useState } from 'react';
import { FlatList, StyleSheet, View } from 'react-native';

const [index, setIndex] = useState(0);
// The opener calls setIndex(1).

<ModalBottomSheet
  detents={[0, '50%', '100%']}
  index={index}
  onIndexChange={setIndex}
  scrollCoordination={{ expand: 'handoff', collapse: 'handoff' }}
  scrimColor="rgba(0, 0, 0, 0.35)"
  surface={<View style={[StyleSheet.absoluteFill, { backgroundColor: 'white' }]} />}
>
  <FlatList {...listProps} style={{ flex: 1 }} />
</ModalBottomSheet>
```

The destination heights and scrim above are app choices. Adjust them to the
source screen and theme. Keep any existing content padding and keyboard
behavior when carrying over the list props.
