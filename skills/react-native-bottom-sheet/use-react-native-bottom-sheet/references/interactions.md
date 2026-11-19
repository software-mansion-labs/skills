# Scroll, keyboard, and modal interactions

Examples are component fragments; their state, data, and styles come from the app.

## Nested scrolling

Use a regular React Native `ScrollView`, `FlatList`, or `SectionList`. Use a fixed or percentage upper detent to give a filling list a bounded viewport. A filling list with `'content'` as the tallest detent can use the full available height. [Detents and index](https://docs.swmansion.com/react-native-bottom-sheet/docs/detents-and-index/) · [React Native ScrollView](https://reactnative.dev/docs/scrollview)

```tsx
<BottomSheet
  detents={[0, '35%', '80%']}
  index={index}
  onIndexChange={setIndex}
  scrollCoordination={{ expand: 'handoff', collapse: 'initial' }}
  surface={<View style={[StyleSheet.absoluteFill, styles.surface]} />}
>
  <FlatList style={{ flex: 1 }} data={items} renderItem={renderItem} />
</BottomSheet>
```

The example shows the default. An upward drag can expand the sheet and continue into the list. A downward drag that starts in a scrolled list stays with the list through its visual start. A new downward drag at that edge can collapse the sheet if it can move lower. Set `collapse: 'handoff'` to let the same drag do both. Use `'none'` for a direction the list should keep. A string applies to both directions. The scrollable also keeps a vertical start in a direction the sheet cannot move. Drags starting outside a scrollable follow normal sheet behavior. [Scroll coordination modes](https://docs.swmansion.com/react-native-bottom-sheet/docs/scroll-coordination/#modes)

- **iOS.** Handoff supports `UIScrollView` subclasses. A reversing drag can transfer in a `handoff` direction after any bounce or refresh pull clears.
- **Android.** The sheet temporarily enables nested scrolling for standard React Native vertical scrollables when needed and restores it after the touch. Third-party scrollers need to implement Android's nested-scrolling-child contract and enable nested scrolling. Scrollers without that support, such as `WebView`, fall back to `initial`.
- **Web.** The browser keeps starts that are at least as horizontal as vertical. Handoff needs a touch ancestor with `overflow-y: auto` or `scroll` that actually overflows.
- **Gesture edges.** Native diagonal starts inside vertical scrollers follow vertical coordination. A horizontal scroller keeps sideways drags. Inverted lists use their visual edge. Native bounce and refresh can take over when the sheet reaches its boundary. [Modes](https://docs.swmansion.com/react-native-bottom-sheet/docs/scroll-coordination/#modes) · [Platform support](https://docs.swmansion.com/react-native-bottom-sheet/docs/scroll-coordination/#platform-support)

## Keyboard

The app supplies keyboard avoidance. Options include letting the host resize, using an app-provided keyboard-aware scroller, or adding bottom padding. For a `'content'` detent with app-animated keyboard padding, apply the padding to the measured content and set `animateContentHeight={false}`. The sheet then follows the changing height without adding its own resize animation. [Keyboard handling](https://docs.swmansion.com/react-native-bottom-sheet/docs/keyboard-handling/) · [Surface](https://docs.swmansion.com/react-native-bottom-sheet/docs/surface/)

This default modal fragment assumes an existing `BottomSheetProvider`. Import
`Animated` from React Native; `keyboardPadding` is an app-supplied
`Animated.Value` animated with `useNativeDriver: false`, since padding changes
layout. [React Native animation limits](https://reactnative.dev/docs/animations#caveats)

```tsx
<ModalBottomSheet
  detents={[0, 'content']}
  index={index}
  onIndexChange={setIndex}
  animateContentHeight={false}
  surface={<View style={[StyleSheet.absoluteFill, styles.surface]} />}
>
  <Animated.View style={{ paddingBottom: keyboardPadding }}>
    <TextInput value={name} onChangeText={setName} accessibilityLabel="Name" />
  </Animated.View>
</ModalBottomSheet>
```

If a form action inside a `ScrollView` takes a second tap while the keyboard is open, check `keyboardShouldPersistTaps`. With `'handled'`, a tap handled by the action skips automatic keyboard dismissal. On iOS, `automaticallyAdjustKeyboardInsets` can update a bounded `ScrollView`'s insets when the keyboard changes size. [Tap handling](https://reactnative.dev/docs/scrollview#keyboardshouldpersisttaps) · [Keyboard insets](https://reactnative.dev/docs/scrollview#automaticallyadjustkeyboardinsets)

**App layout guidance.** Apply each keyboard inset once. On Android, `windowOverlay` presents the sheet in a separate dialog. [Window overlay](https://docs.swmansion.com/react-native-bottom-sheet/docs/modal-sheets/#window-overlay)

## Dismissal and accessibility

Default modal dismissal needs positive current scrim opacity and a non-programmatic closed detent; it reports the selected index through `onIndexChange`. A custom `onCloseRequest` leaves the decision and index update to the app. Dragging to a closed detent reports through `onIndexChange` directly. For inline `BottomSheet`, set the controlled `index` to close from app code. [Close requests](https://docs.swmansion.com/react-native-bottom-sheet/docs/modal-sheets/#close-requests) · [Inline sheets](https://docs.swmansion.com/react-native-bottom-sheet/docs/inline-sheets/)

- **Scrim tap and screen-reader dismissal.** These require positive current scrim opacity and either an `onCloseRequest` handler or a non-programmatic closed detent. A transparent `scrimColor` can still allow dismissal. While the scrim can be tapped to request dismissal, screen readers expose its Dismiss button. Localize it with `scrimAccessibilityLabel`. At full host height the button has no scrim area, but VoiceOver's escape gesture and TalkBack's Dismiss action remain available under the same opacity and actionability conditions. TalkBack requires focus inside the sheet. VoiceOver can consume its escape gesture when the scrim has positive opacity and there is no close action. [Accessible dismissal](https://docs.swmansion.com/react-native-bottom-sheet/docs/modal-sheets/#accessible-dismissal)
- **Back and physical Escape.** A regular or committed predictive Android Back action can request close. Cancelling predictive Back leaves the sheet alone. Android and web accept a completed, unmodified physical Escape press. An `onCloseRequest` handler can receive eligible Back or Escape at zero scrim opacity. A visibly closing sheet consumes further Back or Escape through settle. A top portal can still block lower portals when Back passes through to the app or Escape uses normal key handling. On Android, portal Escape is best effort when focus is outside the portal. [Input behavior](https://docs.swmansion.com/react-native-bottom-sheet/docs/modal-sheets/#input-behavior) · [Presentation ordering](https://docs.swmansion.com/react-native-bottom-sheet/docs/modal-sheets/#presentation-ordering) · [Limitations](https://docs.swmansion.com/react-native-bottom-sheet/docs/modal-sheets/#limitations)

## Safe areas

Use bottom safe-area padding inside the sheet where interactive content reaches the system area. Set a detent's `insets.bottom` when the card itself needs an external gap. If you set `topInset={0}` to draw under the status bar, give interactive content enough top safe-area padding. [Modal sheets example](https://docs.swmansion.com/react-native-bottom-sheet/docs/modal-sheets/) · [Insets](https://docs.swmansion.com/react-native-bottom-sheet/docs/detents-and-index/#insets) · [Top inset](https://docs.swmansion.com/react-native-bottom-sheet/docs/detents-and-index/#top-inset)
