# Detents and content

These mappings replace presentation behavior for React Navigation native-stack
and Expo Router built-in Stack form sheets. Route removal and linking depend
on the app's navigation integration.

## Heights, selection, and events

- `sheetAllowedDetents` takes an ascending array of fractions or
  `'fitToContents'`, and defaults to `[1.0]`. Android uses up to the first
  three entries. Map the effective open heights to percentage strings or
  `'content'`. Destination numeric detents are points/dp, and its default
  `[0, 'content']` changes the source default full-height behavior. Use
  `'100%'` for the intended full height and check the host geometry. A
  destination zero height is closed. If a source list already contains a zero
  fraction, establish what it does before mapping it to a close state.
- Resolve `sheetInitialDetentIndex` (default `0`, or `'last'`) against
  the effective platform list. Set the controlled destination `index` to
  the intended height on open and update it on each `onIndexChange`.
  For example, `[0.25, 0.5, 1]` initially at `1` becomes
  `[0, '25%', '50%', '100%']` initially at `2` with a closed index of `0`.
  If detents change, keep a valid index in the same render. The destination
  throws on an out-of-range index.
- Source `sheetDetentChange` provides `e.data.index` and `e.data.stable`.
  Expo SDK 58 documents stability as always true on iOS and false during
  Android drag or settling; check the matching source version. Use
  `onIndexChange` for a user-initiated snap before it settles.
  A programmatic index update skips that event. Use `onSettle` for snap
  completion from either cause. Translate indices for existing consumers,
  excluding the synthetic closed index from open-detent analytics and
  filtering initial or programmatic events when required. If motion is
  needed, `onPositionChange` gives position and fractional index. Preserve
  navigation transition side effects separately from sheet motion.

[React Navigation options](https://reactnavigation.org/docs/native-stack-navigator/#form-sheet-options)
· [Expo form-sheet detents](https://docs.expo.dev/router/advanced/modals/#configuring-sheet-detents)
· [Expo SDK 58 event payload](https://docs.expo.dev/versions/v58.0.0/sdk/router/#nativestacknavigationeventmap)
· [Destination detents and index events](https://docs.swmansion.com/react-native-bottom-sheet/docs/detents-and-index/)
· [Destination position tracking](https://docs.swmansion.com/react-native-bottom-sheet/docs/position-tracking/)

## Layout, scroll, and keyboard

- Scope native flex behavior to the source version. React Navigation's
  form-sheet guide warns that top-level `flex: 1` with `fitToContents` can
  make the sheet disappear on Android. Its iOS guidance says `fitToContents`
  works with `flex: 1` and fixed-fraction detents can leave the content at its
  intrinsic height. Expo documents support for iOS
  `flex: 1` with numeric detents since SDK 55 and calls for explicit content
  sizing with `fitToContents`. For the destination, measure intrinsic
  children for a compact `'content'` detent. Filling children are laid out
  against the tallest reachable detent. When that detent is `'content'`,
  a filling list can grow to the full available height. Use a fixed or
  percentage upper detent to bound the viewport when needed.
  [React Navigation flex](https://reactnavigation.org/docs/native-stack-navigator/#flex-1-limitations)
  · [Expo flex](https://docs.expo.dev/router/advanced/modals/#using-flex-1-with-custom-detents)
  · [Destination layout](https://docs.swmansion.com/react-native-bottom-sheet/docs/detents-and-index/)
- Source iOS `sheetExpandsWhenScrolledToEdge` defaults to true and requires
  the scroll view on the first-child chain. Use standard React Native
  scrollables in the destination. The default `scrollCoordination` expands
  with `'handoff'` and collapses with `'initial'`. To keep scrolling from
  expanding the sheet, set the expand direction to `'none'`; drags starting
  elsewhere retain normal sheet behavior. Check Android nested scrolling
  and third-party scrollers against supported handoff modes before dropping
  existing workarounds.
  [Source scroll](https://reactnavigation.org/docs/native-stack-navigator/#scroll-behavior)
  · [Source Android limitations](https://reactnavigation.org/docs/native-stack-navigator/#android-specific-limitations)
  · [Destination scroll coordination](https://docs.swmansion.com/react-native-bottom-sheet/docs/scroll-coordination/)
- Source Android `sheetShouldOverflowTopInset` defaults to false. When true,
  the source uses the full stack height. Destination percentages use the
  host height less `topInset` and the detent's `insets.bottom`. The default
  `topInset` is status-bar overlap with the host. Use `topInset={0}` if
  the sheet should reach the host top and pad controls for safe areas as
  needed. Use internal bottom padding for safe-area content and a detent
  bottom inset when the card itself needs a gap.
  [Source insets](https://reactnavigation.org/docs/native-stack-navigator/#sheetshouldoverflowtopinset)
  · [Destination insets](https://docs.swmansion.com/react-native-bottom-sheet/docs/detents-and-index/#insets)
- The app supplies destination keyboard avoidance. Choose host resize,
  keyboard-aware scrolling, or bottom padding for the actual presentation,
  including the separate Android dialog when using `windowOverlay`.
  Avoid applying an inset twice. Source Android `sheetResizeAnimationEnabled`
  controls native `fitToContents` resize animation and defaults to true.
  With a destination `'content'` detent and app-animated padding,
  `animateContentHeight={false}` lets the sheet follow the measured height.
  Check focused inputs and actions at partial detents because native keyboard
  resizing, viewport changes, and timing can differ.
  [Source resize](https://reactnavigation.org/docs/native-stack-navigator/#sheetresizeanimationenabled)
  · [Destination keyboard](https://docs.swmansion.com/react-native-bottom-sheet/docs/keyboard-handling/)
  · [Window overlay](https://docs.swmansion.com/react-native-bottom-sheet/docs/modal-sheets/#window-overlay)

## Surface, footer, and scrim

- Move `contentStyle.backgroundColor` into a filling `surface`, and set
  corner radii and platform-supported shadows in `surfaceStyle`. Source
  radius falls back to the system default, Android elevation defaults to
  `24`, and the iOS grabber defaults to hidden. Compose a handle, title,
  or actions when used. Destination detent insets, `maxWidth`, and
  `placement` (wide sizing since 1.1.0) can approximate the original
  footprint. OS materials, native chrome, and automatic platform adaptation
  may require keeping the native presentation for exact parity.
  [Source appearance](https://reactnavigation.org/docs/native-stack-navigator/#customizing-appearance)
  · [Destination surface](https://docs.swmansion.com/react-native-bottom-sheet/docs/surface/)
  · [Wide screens](https://docs.swmansion.com/react-native-bottom-sheet/docs/detents-and-index/#wide-screens)
  · [Destination releases](https://swmansion.com/changelog/react-native-bottom-sheet/)
- Expo documents Android's experimental `unstable_sheetFooter` as visible
  at every detent. Implement any required pinned footer in the app layout.
  Simply appending it to children laid out against the tallest detent can
  leave it outside the visible area at shorter detents. Check the last list
  item, keyboard, and bottom safe area together. Android's source native
  form-sheet headers and nested stack navigators have platform limits.
  [Expo sheet footer](https://docs.expo.dev/router/advanced/modals/#sheet-footer-android)
  · [Expo Android limitations](https://docs.expo.dev/router/advanced/modals/#android-limitations)
  · [Destination layout](https://docs.swmansion.com/react-native-bottom-sheet/docs/detents-and-index/)
- Map `sheetLargestUndimmedDetentIndex` to modal detent `scrimOpacity`
  and an explicit `scrimColor` for visible dimming. The source defaults
  to `'none'` (dim every open detent). Use zero opacity at open detents
  through a numeric threshold, or at all open detents for `'last'`.
  Translate the threshold after adding a closed detent. The destination
  interpolates opacity, and a transparent color with positive opacity can
  still accept scrim taps. A zero-opacity scrim removes scrim and screen-reader
  dismissal through that path. Keep an accessible close control, preserve
  localized dismissal labels with `scrimAccessibilityLabel`, and check
  background interaction separately, especially under native route or window
  overlays. On iOS the source can dim after a keyboard-driven resize without
  a detent change.
  [Source dimming](https://reactnavigation.org/docs/native-stack-navigator/#controlling-dimming)
  · [Destination scrim](https://docs.swmansion.com/react-native-bottom-sheet/docs/modal-sheets/#scrim)
  · [Accessible dismissal](https://docs.swmansion.com/react-native-bottom-sheet/docs/modal-sheets/#accessible-dismissal)
