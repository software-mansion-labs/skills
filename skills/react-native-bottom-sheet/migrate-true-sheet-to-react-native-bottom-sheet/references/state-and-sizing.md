# State and sizing

## Remap detents and indices

True Sheet v4 documents up to three sorted detents, with `[0.5, 1]` as the
default. Its numbers are fractions of screen height. Destination numbers are
absolute layout units. Convert fractions to percent strings that resolve
against the destination host's available height after top and bottom insets.
Add a zero-height detent when the sheet must close. With
`[0, '50%', '100%']`, source open indices `0` and `1` map to `1` and `2`,
and destination index `0` closes the sheet. Remap every consumer of those
indices, including styles, initial selection, and analytics.
[True Sheet detents](https://sheet.lodev09.com/guides/detents) ·
[Destination detents](https://docs.swmansion.com/react-native-bottom-sheet/docs/detents-and-index/)

Use `'content'` for an `'auto'` sheet whose children can be measured naturally.
True Sheet v4 measures a connected scrollable's content for `'auto'` and
lays out filling content at the visible detent, including while dragging.
The v3 `scrollable` path lacks `'auto'` sizing. The destination lays out
filling children against the tallest detent. If `'content'` is tallest, a
filling list can use the entire available height. Choose an app sizing strategy
when preserving short-list sizing matters; a fixed upper bound alone does not
reproduce the source's content-sized list.
For `'peek'` and `TrueSheetPeek`, preserve the compact layout with an explicit
height or app measurement. A whole-content
detent does not measure to the peek marker. Revisit any `maxContentHeight`
constraint and ensure the resolved detents, including bottom insets, stay
ordered as content and host geometry change.
[True Sheet detents](https://sheet.lodev09.com/guides/detents) ·
[True Sheet v3 scrolling](https://sheet.lodev09.com/v3/guides/scrolling) ·
[True Sheet peek](https://sheet.lodev09.com/guides/peeking) ·
[Destination detents](https://docs.swmansion.com/react-native-bottom-sheet/docs/detents-and-index/)

True Sheet v4 follows the current detent when it survives a list change and
reports the new index. The destination keeps the controlled index and throws if
it is out of range. If the app needs to preserve a detent's identity, derive
its new index from the updated array and pass both in the same render. Define
the fallback when that detent disappears.
[True Sheet changing detents](https://sheet.lodev09.com/guides/detents#changing-detents) ·
[Destination index events](https://docs.swmansion.com/react-native-bottom-sheet/docs/detents-and-index/#index-events)

## Move commands to the state owner

Map `present(sourceIndex)` to selecting an open target index, `resize` to
selecting another open index, and `dismiss` to selecting the closed index.
Preserve the source's requirement that resize acts on an already presented
sheet if callers rely on it. For named native calls or web `useTrueSheet`
callers, route the action to the app's sheet owner through its existing
callbacks or context. Source `initialDetentIndex` defaults to `-1`. Open on
mount by initializing the target index to an open detent. `animateIn={false}`
covers presentation on first mount. Later index-driven opens animate. The
destination docs describe no equivalent per-command animation flag.
[True Sheet methods](https://sheet.lodev09.com/reference/methods) ·
[True Sheet configuration](https://sheet.lodev09.com/reference/configuration) ·
[Destination initial presentation](https://docs.swmansion.com/react-native-bottom-sheet/docs/detents-and-index/#initial-presentation)

Preserve command no-ops too. In v4, `present` on an already presented sheet
does not resize it, and `dismiss` on an already dismissed sheet is a no-op.
Track presentation state before applying those commands so a repeated open
does not change the selected detent. An already completed action cannot wait
for a presumed new `onSettle`; handle that path at the app's command boundary.
[True Sheet v4 command timing](https://sheet.lodev09.com/migration#13-present-and-dismiss-resolve-after-the-animation-web)

This fragment assumes a `BottomSheetProvider` around the destination portal.
`Details` and `sheetBackground` come from the app. Source index `0` and target
index `1` both select the shorter open detent.
[True Sheet usage](https://sheet.lodev09.com/usage) ·
[Destination modal and scrim](https://docs.swmansion.com/react-native-bottom-sheet/docs/modal-sheets/)

```tsx
const sheet = useRef<TrueSheet>(null);
const open = () => sheet.current?.present(0);

<TrueSheet ref={sheet} detents={[0.5, 1]}>
  <Details onClose={() => sheet.current?.dismiss()} />
</TrueSheet>;
```

```tsx
const [index, setIndex] = useState(0);
const open = () => setIndex(1);

<ModalBottomSheet
  detents={[0, '50%', '100%']}
  index={index}
  onIndexChange={setIndex}
  scrimColor="rgba(0, 0, 0, 0.5)"
  surface={
    <View style={[StyleSheet.absoluteFill, { backgroundColor: sheetBackground }]} />
  }
>
  <Details onClose={() => setIndex(0)} />
</ModalBottomSheet>;
```

## Preserve timing and motion

Update controlled state from `onIndexChange` when a user snap commits.
Programmatic changes already come from the state owner. Move settled work
from `onDetentChange` and after-open or after-close handlers to `onSettle`
where the required signal is arrival at the target detent. Gate it by the
intended transition and previously settled state where needed. An index write
does not await motion. Review every
`await present/resize/dismiss` and will/did handler for what it sequences.
True Sheet v4 changed web `present` and `dismiss` to resolve after animation,
and iOS `onWillDismiss` now reports committed dismissal. The destination docs
provide no direct counterparts for native content readiness, focus/blur, or
drag phase callbacks.
[True Sheet events](https://sheet.lodev09.com/reference/events) ·
[True Sheet v4 timing changes](https://sheet.lodev09.com/migration) ·
[Destination index events](https://docs.swmansion.com/react-native-bottom-sheet/docs/detents-and-index/#index-events)

True Sheet tracks screen Y from the top. The destination reports distance
from the host bottom to the sheet top, plus fractional destination index.
Recalculate origins, directions, and detent thresholds. Source `realtime`
and normalized detent values lack a documented direct field mapping.
When retaining Reanimated, use the documented native event binding through
`wrapNativeView` and `useEvent` if the consumer needs UI-thread motion.
[True Sheet events](https://sheet.lodev09.com/reference/events) ·
[True Sheet Reanimated](https://sheet.lodev09.com/guides/reanimated) ·
[Destination position tracking](https://docs.swmansion.com/react-native-bottom-sheet/docs/position-tracking/)
