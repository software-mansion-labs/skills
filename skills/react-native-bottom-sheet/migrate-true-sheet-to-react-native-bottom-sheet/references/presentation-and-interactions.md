# Presentation and interactions

## Choose the host and preserve navigation

True Sheet uses the system presentation on iOS and Android. Start with
`ModalBottomSheet` for a presented sheet. Its default portal needs
`BottomSheetProvider`. Put shared providers above it or pass screen values
and navigation callbacks from the owner. Use `windowOverlay` when the sheet
must cover a native modal. It preserves React ancestry while using a separate
presentation layer. Use inline `BottomSheet` when the sheet belongs in the
screen's layout. When removing `TrueSheetOverlay`, place toasts, menus, and
dialogs deliberately relative to the provider or window layer. The destination
documents no general-purpose overlay equivalent.
[True Sheet overlays](https://sheet.lodev09.com/guides/overlays) ·
[Destination comparison](https://docs.swmansion.com/react-native-bottom-sheet/docs/how-it-compares/) ·
[Destination modal hosts and context](https://docs.swmansion.com/react-native-bottom-sheet/docs/modal-sheets/#window-overlay)

True Sheet hides and restores a covered sheet that was below full expansion.
Its `dismiss` can close the sheet and everything above it, `dismissStack`
closes only the sheets above it, and `dismissAll` acts in the current
presentation context. Destination modals stack by opening order within a layer,
with window overlays above portal sheets. App state owns any required
hide/restore, cascading closure, or focus callbacks. When source sheets are
routes, retain params, history, back behavior, and deep links in the navigation
layer. `createTrueSheetNavigator`, `useTrueSheetNavigation`, and Expo Router
sheet layouts have no documented direct target replacement. Check v3's
`withLayoutContext` integration and v4's dedicated Expo Router entry point.
[True Sheet stacking](https://sheet.lodev09.com/guides/stacking) ·
[True Sheet navigation](https://sheet.lodev09.com/guides/navigation) ·
[True Sheet v4 migration](https://sheet.lodev09.com/migration) ·
[Destination stacking](https://docs.swmansion.com/react-native-bottom-sheet/docs/modal-sheets/#stacking)

Check when children mount and when state should reset. True Sheet v4 documents
lazy mounting on first presentation by default on native. The destination
keeps its view and children mounted while closed by default. Its `lazy` mounts
content on open and unmounts it after settling closed, which resets internal
child state. Keep screen and wrapper lifecycles explicit and review autofocus.
[True Sheet configuration](https://sheet.lodev09.com/reference/configuration) ·
[Destination lazy content](https://docs.swmansion.com/react-native-bottom-sheet/docs/modal-sheets/#lazy-content)

## Rebuild dismissal and accessibility

Set `scrimColor` for visible dimming and per-detent `scrimOpacity` for the
source's dimming threshold or background interaction. Target scrim color
defaults to transparent, and positive detent opacity can still make it
interactive. Use zero detent opacity where background interaction is needed
and verify it in the chosen host. Source `dimOpacity={0}` blocks background
touches. Source `dimmed={false}` allows them. Source `dimBlur` has no
documented direct scrim equivalent.
[True Sheet configuration](https://sheet.lodev09.com/reference/configuration) ·
[Destination scrim](https://docs.swmansion.com/react-native-bottom-sheet/docs/modal-sheets/#scrim)

Default destination dismissal needs positive current scrim opacity and a
non-programmatic closed detent. For a guarded dismissal, make every zero-height
detent programmatic, such as `{ value: 0, programmatic: true }`, to block
snap-to-closed by drag while
retaining code-driven closure. Use
`onCloseRequest` to decide eligible scrim, Android Back, physical Escape on
Android/web, and accessibility requests. It has no documented request-source
argument. Set the index on approval. Dragging
commits through `onIndexChange`. True Sheet v4 `onDismissAttempt` includes
blocked drag attempts and excludes dim taps. The target has no documented
blocked-drag attempt callback: a blocked swipe does not invoke `onCloseRequest`.
Provide an accessible close control using the confirmation flow. There are no
documented direct equivalents to `dismissThreshold`,
`draggable={false}`, or the source Android `onBackPress` contract.
[True Sheet events](https://sheet.lodev09.com/reference/events) ·
[Destination programmatic detents](https://docs.swmansion.com/react-native-bottom-sheet/docs/detents-and-index/#programmatic-only-detents) ·
[Destination close requests](https://docs.swmansion.com/react-native-bottom-sheet/docs/modal-sheets/#close-requests)

Review how users focus, resize, and close each sheet. Localize the destination
scrim's dismiss label with `scrimAccessibilityLabel`. Screen-reader dismissal
requires positive current scrim opacity and an available close action. Use app
controls for required grabber actions, labels, and detent announcements when
the source supplies them.
[True Sheet accessibility options](https://sheet.lodev09.com/reference/configuration) ·
[Destination accessible dismissal](https://docs.swmansion.com/react-native-bottom-sheet/docs/modal-sheets/#accessible-dismissal)

## Restore scrolling, keyboard behavior, and chrome

Keep regular React Native scrollables and refs used by the app. Replace v3
`scrollable` or v4 `scrollableRef` wiring with the destination's native
coordination. The destination defaults to `expand: 'handoff'` and
`collapse: 'initial'`. Consider `expand: 'none'` for source settings that
prevent expansion from a scrollable. Gestures starting outside scrollables
follow normal sheet behavior. Check platform handoff limits for nonstandard
scrollers.
[True Sheet v3 scrolling](https://sheet.lodev09.com/v3/guides/scrolling) ·
[True Sheet v4 scrolling](https://sheet.lodev09.com/guides/scrolling) ·
[Destination scroll coordination](https://docs.swmansion.com/react-native-bottom-sheet/docs/scroll-coordination/)

The destination delegates keyboard adjustment to the app. Plan how input
focus, scrolling, safe-area insets, and any pinned action behave as the
keyboard changes. Source v4 coordinates these with its scrollable and footer
layout. Choose host resizing, a keyboard-aware scroller, or padding as needed,
accounting for each inset once. If animated app padding drives a
`'content'` detent, use `animateContentHeight={false}` to follow that
measurement without a second animation.
[True Sheet keyboard](https://sheet.lodev09.com/guides/keyboard) ·
[Destination keyboard handling](https://docs.swmansion.com/react-native-bottom-sheet/docs/keyboard-handling/)

Put source content `style` on the corresponding content wrapper. Compose
headers, footers, and grabbers in the destination content and account for their
height and insets. A footer at the bottom of a tallest-detent filling child can
be offscreen at shorter detents; sticky placement needs app-owned layout.
Render color or a custom background in a filling `surface`
and put radius and supported shadows in `surfaceStyle`. Per-detent insets can
recreate a floating card, and destination 1.1+ wide-screen props can
approximate width and alignment. Treat exact native Liquid Glass, blur
transitions, form/page presentation, and system grabber behavior as product
requirements that may need a different presentation or an explicit compromise.
[True Sheet configuration](https://sheet.lodev09.com/reference/configuration) ·
[Destination surface](https://docs.swmansion.com/react-native-bottom-sheet/docs/surface/) ·
[Destination insets and wide screens](https://docs.swmansion.com/react-native-bottom-sheet/docs/detents-and-index/#insets) ·
[Destination 1.1.0 release](https://swmansion.com/changelog/react-native-bottom-sheet-1-1-0/)
