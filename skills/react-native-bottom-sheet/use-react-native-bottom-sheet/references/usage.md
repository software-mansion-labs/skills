# Presentation and sizing

The examples use the documented API of `@swmansion/react-native-bottom-sheet`.
[Getting started](https://docs.swmansion.com/react-native-bottom-sheet/docs/getting-started/)

## Choose a host

| Choose | When | Placement and context |
| --- | --- | --- |
| `BottomSheet` | The sheet belongs to a region of the screen. | Mount it in the view it should cover. The host fills that view, stays in the screen's React tree, and accepts host `style` such as `zIndex`. |
| `ModalBottomSheet` | The sheet should cover the app's screen content. | Wrap the app with `BottomSheetProvider`. Its portal renders the sheet inside the provider. |
| `ModalBottomSheet` with `windowOverlay` | The sheet should cover a native modal, such as a native-stack screen presented as `modal`. | The sheet covers the full window and keeps its React ancestry at the call site. The provider is optional. |

[Inline sheets](https://docs.swmansion.com/react-native-bottom-sheet/docs/inline-sheets/) ·
[Window overlay](https://docs.swmansion.com/react-native-bottom-sheet/docs/modal-sheets/#window-overlay)

- Put context providers shared by portal sheets above `BottomSheetProvider`.
  For values provided closer to the sheet owner, read them in the owner and
  pass props or provide the value again inside the sheet. Read navigation in
  the owner and pass callbacks such as `onDone={() => navigation.goBack()}`.
  [React context](https://docs.swmansion.com/react-native-bottom-sheet/docs/modal-sheets/#react-context)
- Within a presentation layer, the most recently opened modal is on top,
  including sheets reopened while closing. Window overlays appear above portal
  sheets. On Android, a window overlay that does not handle Back forwards it
  to the app, where a portal sheet can receive it.
  [Stacking](https://docs.swmansion.com/react-native-bottom-sheet/docs/modal-sheets/#stacking)

## Control the index

- Use a finite, zero-based integer within the detent array. With the default
  `[0, 'content']`, index `0` is closed and index `1` opens to the measured
  content height. Numeric `0` and `'0%'` are closed detent values.
- Set `index` to the desired detent to open, close, or snap from app code.
  Update the controlled `index` from `onIndexChange` when a drag commits or a
  default modal close request starts a snap. App changes to `index` skip that
  callback. Use `onSettle` only for work that depends on snap completion.
- When changing `detents`, the sheet keeps the selected index and moves to its
  new height as needed. If the new array would put the index out of range,
  update both in the same render. Invalid indices and empty detent arrays throw
  synchronously in JavaScript.
  [Index events](https://docs.swmansion.com/react-native-bottom-sheet/docs/detents-and-index/#index-events)
- A sheet that mounts at an open detent animates in by default. Set
  `animateIn={false}` to show it there immediately. Later opens still animate.
  [Initial presentation](https://docs.swmansion.com/react-native-bottom-sheet/docs/detents-and-index/#initial-presentation)

## Choose detents and size content

- Use finite non-negative layout units such as `[0, 320]`, percentage strings
  such as `'40%'`, or measured `'content'`. Percentages accept unsigned values
  from `0%` through `100%` without whitespace. Keep mixed detents in ascending
  order by resolved height plus bottom inset across the content sizes and host
  heights your app supports.
- Use intrinsic content sizing for a compact `'content'` sheet. Children with
  `flex: 1` fill the content area sized to the tallest reachable detent. With
  `'content'` as the tallest detent, a filling child can use the full available
  height. A fixed or percentage upper detent bounds a filling scrollable to
  that detent's resolved height.
  [Detents and index](https://docs.swmansion.com/react-native-bottom-sheet/docs/detents-and-index/)
- Detent objects take `value`, optional `programmatic` and `insets`, and
  modal-only `scrimOpacity`. A detent with `programmatic: true` is reachable by
  setting `index` and is skipped during drag snapping.
  [Programmatic-only detents](https://docs.swmansion.com/react-native-bottom-sheet/docs/detents-and-index/#programmatic-only-detents)
- Percentages resolve against host height after subtracting the top inset and
  that detent's bottom inset. The default `topInset` is the portion of the
  status bar overlapping the host.
  [Insets](https://docs.swmansion.com/react-native-bottom-sheet/docs/detents-and-index/#insets) ·
  [Top inset](https://docs.swmansion.com/react-native-bottom-sheet/docs/detents-and-index/#top-inset)
- Use detent `insets` for gaps around a floating card. For wide layouts,
  `maxWidth` accepts a positive number and caps the width including side
  insets. Omit it or pass `Infinity` for no cap. `placement` defaults to
  `'center'` and also accepts direction-aware `'start'` and `'end'`. The
  library lays out the surface and content at the widest card width and
  scales both to fit narrower detents.
  [Wide screens](https://docs.swmansion.com/react-native-bottom-sheet/docs/detents-and-index/#wide-screens)

## Style the body and scrim

- Put a filling element such as a `View` styled with `StyleSheet.absoluteFill`
  in `surface`. The library sizes it independently of the children so the
  fill covers the sheet as content height changes. Use
  `animateContentHeight={false}` to follow the children's height immediately
  when they animate and the active detent is `'content'`.
- Set corner radii on `surfaceStyle` to clip the fill and content together.
  Set platform-supported card shadows there too. On an inline `BottomSheet`,
  use `style` to adjust the host.
  [Surface and platform shadow support](https://docs.swmansion.com/react-native-bottom-sheet/docs/surface/) ·
  [Inline host style](https://docs.swmansion.com/react-native-bottom-sheet/docs/inline-sheets/)
- The modal's `scrimColor` defaults to transparent. Detent `scrimOpacity`
  ranges from 0 to 1 and interpolates linearly. It defaults to 0 at closed
  detents and 1 at open detents. A transparent color can still provide a
  dismiss target when the current scrim opacity is positive and the sheet
  can handle the request.
  [Modal scrim](https://docs.swmansion.com/react-native-bottom-sheet/docs/modal-sheets/#scrim)

## Close and preserve state

- With positive current scrim opacity and a non-programmatic closed detent,
  the default modal close response snaps there and calls `onIndexChange`.
  Dragging to close calls `onIndexChange` and bypasses `onCloseRequest`.
- Supply `onCloseRequest` when the app must decide whether to close.
  The handler leaves the index under app control. Programmatic closed
  detents are skipped during drag snapping. For a sheet with one zero detent,
  use `{ value: 0, programmatic: true }` to prevent drag dismissal. A blocked
  swipe does not call the handler. On acceptance, set the closed index.
  [Close requests and platform limits](https://docs.swmansion.com/react-native-bottom-sheet/docs/modal-sheets/#close-requests) ·
  [Programmatic detents](https://docs.swmansion.com/react-native-bottom-sheet/docs/detents-and-index/#programmatic-only-detents)
- For a confirmation flow, give the sheet a reachable close control that
  runs the same check.
- Keep the sheet mounted through a closing animation. By default its
  native view, children, and surface stay mounted while closed. With
  `lazy`, both sheet variants mount those parts in the render that selects
  an open detent and unmount them after settling closed, resetting state
  held in the content. Without a closed detent, `lazy` has no effect. Add a
  zero-height detent for an animated close or manage removal in app code.
  [Lazy content](https://docs.swmansion.com/react-native-bottom-sheet/docs/modal-sheets/#lazy-content)

## Track motion

- `onPositionChange` provides a continuous event. Read
  `event.nativeEvent.position` for the distance from the host's bottom to the
  sheet's top in layout units, and `event.nativeEvent.index` for the fractional
  detent index. A lazy sheet emits no position events while unmounted.
  Use `onSettle` to observe when it reaches the closed detent.
- To drive animation on the native UI thread with Reanimated, pass
  `wrapNativeView={Animated.createAnimatedComponent}` and a `useEvent`
  worklet subscribed to `onPositionChange`. The worklet receives the
  unwrapped `event.position` and `event.index`. Pass a stable wrapper
  function so the worklet binds to the native sheet view on the first
  render. On web, sheet motion runs in JavaScript.
  [Position tracking and Reanimated integration](https://docs.swmansion.com/react-native-bottom-sheet/docs/position-tracking/) ·
  [Web support](https://docs.swmansion.com/react-native-bottom-sheet/docs/getting-started/#web)
- Set `settleDuration` to a finite positive value in milliseconds to tune
  the spring. It defaults to 450 ms. A settle starting at zero velocity
  takes the full duration, while a fast fling toward the detent can
  finish sooner. Use `onSettle` for side effects that depend on completion.
  [Settle animation](https://docs.swmansion.com/react-native-bottom-sheet/docs/settle-animation/) ·
  [Index events](https://docs.swmansion.com/react-native-bottom-sheet/docs/detents-and-index/#index-events)

## Controlled modal example

```tsx
import { useState } from 'react';
import { Button, StyleSheet, Text, View } from 'react-native';
import {
  BottomSheetProvider,
  ModalBottomSheet,
} from '@swmansion/react-native-bottom-sheet';

export default function App() {
  const [index, setIndex] = useState(0);

  return (
    <View style={{ flex: 1 }}>
      <BottomSheetProvider>
        <Button title="Show details" onPress={() => setIndex(1)} />
        <ModalBottomSheet
          detents={[0, 'content']}
          index={index}
          onIndexChange={setIndex}
          scrimColor="rgba(0, 0, 0, 0.35)"
          surface={
            <View style={[StyleSheet.absoluteFill, { backgroundColor: 'white' }]} />
          }
          surfaceStyle={{ borderTopLeftRadius: 20, borderTopRightRadius: 20 }}
        >
          <View style={{ padding: 24 }}>
            <Text>Order details</Text>
            <Button title="Close" onPress={() => setIndex(0)} />
          </View>
        </ModalBottomSheet>
      </BottomSheetProvider>
    </View>
  );
}
```

In an app, add bottom safe-area padding to sheet content where needed.
[Modal sheets example](https://docs.swmansion.com/react-native-bottom-sheet/docs/modal-sheets/)
