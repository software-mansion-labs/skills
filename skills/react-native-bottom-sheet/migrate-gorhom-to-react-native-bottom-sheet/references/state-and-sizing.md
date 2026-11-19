# State, sizing, and commands

The rows below are migration adaptations derived from the two public APIs,
not an upstream compatibility layer. Resolve the app's Gorhom version first.

## Rebuild the index table

Gorhom v5 defaults to initial index `0`; set `index={-1}` to mount closed.
Dynamic sizing defaults to enabled and can insert a sorted content-height stop,
changing numeric index meanings. Record the resolved stops used by commands,
callbacks, backdrop thresholds, and animations before choosing destination
indices. [Gorhom configuration](https://gorhom.dev/react-native-bottom-sheet/props#configuration)
· [Dynamic sizing](https://gorhom.dev/react-native-bottom-sheet/dynamic-sizing)

The destination uses zero-based array indices, not `-1`. Add a zero-height
detent for animated closing; `[0, 'content']` is the default. Keep detents
ordered by resolved height plus bottom inset at every supported size. Adding
`'content'` is an explicit sizing choice, not a renamed dynamic-sizing switch.
Mixed content and fixed stops can change order as content grows.
[Destination detents](https://docs.swmansion.com/react-native-bottom-sheet/docs/detents-and-index/)

Only when the exact ordered open-stop list is retained and one closed stop is
prepended does an old open index become `oldIndex + 1`. Map all other cases by
intended height. Destination percentages use host height minus top and detent
bottom insets, so matching strings alone does not establish equal geometry.
Revisit `maxDynamicContentSize`, shared-value `snapPoints`, and measured-layout
hooks as explicit sizing requirements; don't forward their
props. V4's `CONTENT_HEIGHT` and `useBottomSheetDynamicSnapPoints` belong to
the old measurement model. [Gorhom size limit](https://gorhom.dev/react-native-bottom-sheet/props#maxdynamiccontentsize)
· [V4 dynamic hook](https://gorhom.dev/react-native-bottom-sheet/v4/hooks#usebottomsheetdynamicsnappoints)

## Replace commands with controlled state

| Gorhom call | Destination adaptation |
| --- | --- |
| `snapToIndex(i)` | Set state to the mapped detent index. |
| `expand()` / `collapse()` | Select the intended largest / smallest open stop. A newly added zero stop makes index `0` closed, not the old collapsed height. |
| `close()` | Select the explicit closed detent. |
| `snapToPosition(value)` | Model the required position as a detent; use `programmatic: true` if drag snapping must skip it. Keep resolved order valid. |
| `forceClose()` | Closing state alone does not establish Gorhom's interruption-prevention guarantee. Flag that requirement for design. |
| `useBottomSheet()` callers | Pass app-owned state/actions through the existing wrapper, props, or app context; use documented position events for motion consumers. |

[Gorhom methods](https://gorhom.dev/react-native-bottom-sheet/methods)
· [Gorhom sheet hook](https://gorhom.dev/react-native-bottom-sheet/hooks#usebottomsheet)
· [Destination control and programmatic detents](https://docs.swmansion.com/react-native-bottom-sheet/docs/detents-and-index/#programmatic-only-detents)

Port source `onChange` consumers according to purpose. Feed `onIndexChange`
back into the controlled state. It reports drag commits and default modal close
snaps, not app assignments to `index`. Put app-command
side effects in the app's action path. Use `onSettle` for completed movement;
it is not a rename of Gorhom's before-animation `onAnimate`. If `detents` shrink,
update an out-of-range index in the same render. `animateIn` controls initial
open mounting, not all later openings.
[Destination index events](https://docs.swmansion.com/react-native-bottom-sheet/docs/detents-and-index/#index-events)
· [Gorhom callbacks](https://gorhom.dev/react-native-bottom-sheet/props#callbacks)

## Fixed-height inline example

These component fragments assume React's hooks and React Native's `Button`,
`StyleSheet`, `Text`, and `View` are imported. Their parent fills the intended
sheet region. The source fragment assumes Gorhom's configured dependencies and
an existing `GestureHandlerRootView` ancestor. It explicitly disables dynamic
sizing.

Before, with the default `BottomSheet` import from `@gorhom/bottom-sheet`:

```tsx
const sheetRef = useRef<BottomSheet>(null);

return (
  <>
    <Button title="Open" onPress={() => sheetRef.current?.snapToIndex(0)} />
    <BottomSheet
      ref={sheetRef}
      index={-1}
      snapPoints={[240, 480]}
      enableDynamicSizing={false}
      enablePanDownToClose
    >
      <Text>Details</Text>
      <Button title="Close" onPress={() => sheetRef.current?.close()} />
    </BottomSheet>
  </>
);
```

After, with the named `BottomSheet` import from
`@swmansion/react-native-bottom-sheet`:

```tsx
const [index, setIndex] = useState(0);

return (
  <>
    <Button title="Open" onPress={() => setIndex(1)} />
    <BottomSheet
      detents={[0, 240, 480]}
      index={index}
      onIndexChange={setIndex}
      surface={<View style={[StyleSheet.absoluteFill, { backgroundColor: 'white' }]} />}
    >
      <Text>Details</Text>
      <Button title="Close" onPress={() => setIndex(0)} />
    </BottomSheet>
  </>
);
```

The old first open stop is now index `1`; expanding selects `2`. This preserves
the example's drag-to-close intent. Where closing must be app-only, make the
zero detent programmatic and provide an accessible app close action.
[Gorhom inline API](https://gorhom.dev/react-native-bottom-sheet/usage)
· [Destination inline API](https://docs.swmansion.com/react-native-bottom-sheet/docs/inline-sheets/)
· [Destination surface](https://docs.swmansion.com/react-native-bottom-sheet/docs/surface/)

## Rebuild motion consumers

Gorhom exposes `animatedIndex`, `animatedPosition`, and animation-config hooks.
The destination's continuous `onPositionChange` provides
`event.nativeEvent.index` and a position measured from the host bottom to the
sheet top; the event index is fractional. Lazy sheets emit no position events
while unmounted. Recalculate interpolation domains from that definition and the
new index table; copying a raw position value does not establish equal coordinates.
For UI-thread consumers, follow the documented stable `wrapNativeView` plus
Reanimated `useEvent` pattern, whose event is unwrapped. Keep Reanimated when
the app still uses it. [Gorhom animated nodes](https://gorhom.dev/react-native-bottom-sheet/props#animated-nodes)
· [Destination position tracking](https://docs.swmansion.com/react-native-bottom-sheet/docs/position-tracking/)

Per-call timing/spring configs and `overrideReduceMotion` have no direct mapping
established by these guides. The destination documents a positive
`settleDuration`, not Gorhom's animation config shape. Preserve reduced-motion
requirements and report any unresolved behavior rather than promising parity.
[Gorhom animation hooks](https://gorhom.dev/react-native-bottom-sheet/hooks)
· [Gorhom reduced motion](https://gorhom.dev/react-native-bottom-sheet/props#overridereducemotion)
· [Destination settle animation](https://docs.swmansion.com/react-native-bottom-sheet/docs/settle-animation/)
