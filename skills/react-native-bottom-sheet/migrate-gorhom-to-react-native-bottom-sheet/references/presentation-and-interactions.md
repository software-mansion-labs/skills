# Presentation and interaction adaptations

These are migration decisions, not prop-for-prop equivalents. Apply only the
sections corresponding to features used in the app.

## Modal lifecycle and presentation

| Gorhom usage | Migration decision |
| --- | --- |
| `BottomSheetModal` / `BottomSheetModalProvider` | Use `ModalBottomSheet` / `BottomSheetProvider` for destination portal presentation. |
| `present(data?)` | Store the payload in app state and select an open index. |
| `dismiss()` / `enableDismissOnClose` | Set a zero-height detent index. Source dismiss unmounts, and `enableDismissOnClose` defaults true. Destination content persists by default; choose `lazy` for internal content reset. |
| `onDismiss` | Source unmount notification differs from destination `onSettle` motion completion. Keep unmount-dependent work in lifecycle cleanup; guard close-completion effects against duplicate calls or reopening. |
| `dismiss(key?)` / `dismissAll()` | Preserve app-owned actions for named sheets, the last presented sheet when no key is supplied, or all sheets. |

[Gorhom modal methods](https://gorhom.dev/react-native-bottom-sheet/modal/methods)
· [Gorhom modal props](https://gorhom.dev/react-native-bottom-sheet/modal/props)
· [Gorhom modal hook](https://gorhom.dev/react-native-bottom-sheet/modal/hooks)
· [Destination modal lifecycle](https://docs.swmansion.com/react-native-bottom-sheet/docs/modal-sheets/#lazy-content)

Keep the destination component mounted through closing. `lazy` removes its
native view, children, and surface after settling closed; it has no effect
without a closed detent. Portal content reads context above
`BottomSheetProvider`; pass screen-local values/actions explicitly or
provide values again inside. `windowOverlay` retains call-site ancestry and
covers native modals without the provider. Choose it for required native-window
coverage, not as a mechanical rename of Gorhom's `containerComponent`.
[Destination context and coverage](https://docs.swmansion.com/react-native-bottom-sheet/docs/modal-sheets/#react-context)
· [Window overlay](https://docs.swmansion.com/react-native-bottom-sheet/docs/modal-sheets/#window-overlay)
· [Completion events](https://docs.swmansion.com/react-native-bottom-sheet/docs/detents-and-index/#index-events)

Destination modals stack by open order. Preserve app-specific minimize, replace,
and restore behavior explicitly. Gorhom's v5-labeled `stackBehavior` page says
v3-only while its v5 release notes discuss that prop; flag this documentation
conflict and establish the source requirement instead of assuming a default.
[Gorhom stacking prop](https://gorhom.dev/react-native-bottom-sheet/modal/props#stackbehavior)
· [Gorhom v5 release](https://gorhom.dev/react-native-bottom-sheet/blog/bottom-sheet-v5#breaking-changes)
· [Destination stacking](https://docs.swmansion.com/react-native-bottom-sheet/docs/modal-sheets/#stacking)

## Backdrop, dismissal, and accessibility

Gorhom's optional `BottomSheetBackdrop` controls thresholds, touch-through,
and press actions. Audit default indices `1`/`0`, custom targets, and `onPress`
side effects: it runs before the action and is ignored for `pressBehavior="none"`.
The destination's transparent default scrim can still handle dismissal;
preserve both appearance and interaction deliberately using
`scrimColor`, per-detent `scrimOpacity`, and app close logic.
[Gorhom backdrop](https://gorhom.dev/react-native-bottom-sheet/components/bottomsheetbackdrop)
· [Destination scrim](https://docs.swmansion.com/react-native-bottom-sheet/docs/modal-sheets/#scrim)

Default modal dismissal needs positive scrim opacity and a non-programmatic
closed detent. `onCloseRequest` shares eligible scrim, screen-reader, Back, and
Escape requests without a documented source argument; a no-op consumes eligible
Back/Escape input. Drag commits bypass it. To guard dirty edits, make every zero detent
programmatic, confirm in app code, then set the closed index. A blocked swipe
does not request confirmation; provide an accessible close control. Keep
positive opacity for scrim/screen-reader requests and review platform limits.
[Destination close and accessibility rules](https://docs.swmansion.com/react-native-bottom-sheet/docs/modal-sheets/#close-requests)

Map a source backdrop's collapse/custom action to this shared callback only if
the same policy should govern other dismissal inputs. Exact touch-through,
arbitrary backdrop UI, and threshold parity are not established by the scrim
guide. Preserve required accessible labels via `scrimAccessibilityLabel`.

## Scrollables and gesture ownership

Replace `BottomSheetView` with a regular `View`, and Gorhom's integrated
`ScrollView`, `FlatList`, and `SectionList` wrappers with their React Native
components. Preserve list data, keys, refs, refresh, and layout requirements;
remove Gorhom-specific `focusHook` plumbing rather than forwarding it. Give a
filling list a bounded upper detent. Inventory `BottomSheetVirtualizedList`,
`BottomSheetFlashList`, and `useBottomSheetScrollableCreator` separately rather
than assuming destination support for every scroller.
[Gorhom scrollables](https://gorhom.dev/react-native-bottom-sheet/scrollables)
· [Gorhom FlatList props](https://gorhom.dev/react-native-bottom-sheet/components/bottomsheetflatlist)
· [Library comparison](https://docs.swmansion.com/react-native-bottom-sheet/docs/how-it-compares/)
· [Destination viewport sizing](https://docs.swmansion.com/react-native-bottom-sheet/docs/detents-and-index/)

Destination `scrollCoordination` defaults to upward `'handoff'` and downward
`'initial'`. Downward scrolling may need a new gesture at the top to collapse;
choose `'handoff'` explicitly if one-motion collapse is required. Android
third-party scrollers need the documented nested-scrolling contract; web
handoff needs an overflowing scroll ancestor. Enable nested scrolling on those
Android scrollers; unsupported children fall back to `'initial'`. Map source
gesture wrappers, `waitFor`, `simultaneousHandlers`, panning toggles, and overdrag requirements
individually; there is no documented rename for that configuration set.
[Gorhom gesture configuration](https://gorhom.dev/react-native-bottom-sheet/props#gesture-configuration)
· [Destination scroll modes and platform support](https://docs.swmansion.com/react-native-bottom-sheet/docs/scroll-coordination/)

## Keyboard and custom chrome

Gorhom's `BottomSheetTextInput` and keyboard props integrate with its own
keyboard machinery. Replacing the input alone loses that policy. The
destination accepts regular app inputs but leaves avoidance to the app:
retain or choose a resize, padding, or keyboard-aware-scroll strategy, including
blur/restore and sticky actions.
For app-animated padding measured by a `'content'` detent, set
`animateContentHeight={false}`. Apply each keyboard inset once.
[Gorhom keyboard guide](https://gorhom.dev/react-native-bottom-sheet/keyboard-handling)
· [Destination keyboard guide](https://docs.swmansion.com/react-native-bottom-sheet/docs/keyboard-handling/)

Provide a filling `surface` element and use `surfaceStyle` for shared clipping,
radii, and shadows. Recompose any handle and footer in app content, preserving
their layout and keyboard behavior instead of passing Gorhom component props.
Map floating-card gaps to per-detent `insets`; distinguish those gaps from
content safe-area padding and review the destination's different `topInset`
default. [Gorhom custom handle](https://gorhom.dev/react-native-bottom-sheet/custom-handle)
· [Gorhom custom footer](https://gorhom.dev/react-native-bottom-sheet/custom-footer)
· [Destination surface](https://docs.swmansion.com/react-native-bottom-sheet/docs/surface/)
· [Destination insets](https://docs.swmansion.com/react-native-bottom-sheet/docs/detents-and-index/#insets)

The guides do not establish identical handle-only dragging, keyboard-following
footers, custom gesture handlers, or modal switch/replace semantics. Keep such
requirements visible in the migration report until an app-owned adaptation or
accepted behavior change resolves them.
