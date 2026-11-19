# Navigation and dismissal

## Keep the intended navigation behavior

React Navigation native-stack and Expo Router's built-in `Stack` present
`formSheet` as a screen. The navigator owns the route, including its params,
history, focus, nested screens, and removal. The destination sheet is a
controlled component. Keep route names and files, typed params, links, route guards,
and lifecycle consumers when preserving the route. For an approved local
conversion, move opening payloads and dismissal into the owning screen and
update every former route caller, URL, and cleanup path in scope.
[React Navigation form sheets](https://reactnavigation.org/docs/native-stack-navigator/#form-sheets)
· [React Navigation params](https://reactnavigation.org/docs/params/)
· [Expo Router modal routes](https://docs.expo.dev/router/advanced/modals/#modal-screen-using-expo-router)
· [Expo Router URL params](https://docs.expo.dev/router/reference/url-parameters/)

For a route-backed replacement, an app can keep the screen and render the
target sheet in a different presentation. A full-bounds `transparentModal`
route with a transparent scene is one option when the previous route needs
to stay visible and the navigator and platform support it. Verify that the
intended previous screen exists and that its focus and touch behavior are
appropriate. Coordinate the route animation with the sheet animation and
choose one owner for the backdrop. Replace the original `formSheet` option
for that route to avoid nesting two sheets. The app owns this integration.
Check native headers, system transitions and gestures, status bar treatment,
and nested screen behavior separately from the sheet's visuals.
[React Navigation presentation](https://reactnavigation.org/docs/native-stack-navigator/#presentation)
· [Expo Router transparent overlay route](https://docs.expo.dev/router/advanced/web-modals/#use-a-transparent-overlay-route)
· [Destination comparison](https://docs.swmansion.com/react-native-bottom-sheet/docs/how-it-compares/)

## Choose the host and read route context there

- `BottomSheet` stays within its screen layout. Use it when that layout owns
  the surrounding interaction and dismissal.
- `ModalBottomSheet` renders by default through `BottomSheetProvider`. Place
  the provider at the layer that should host the sheet. A portal in a native
  modal screen is confined to that presentation. The portal sees context
  above its provider, which may exclude the opening route's context. Read
  React Navigation's route/navigation values or Expo Router's
  `useLocalSearchParams` and `useRouter` in the owning screen and pass
  payloads and callbacks to the sheet.
- `windowOverlay` bypasses the provider and can cover an existing native modal.
  It keeps the sheet's React context and changes visual stacking and input
  ownership. On Android it uses a separate dialog. Use it when the sheet
  needs that window-level placement. Modal sheet open order controls visual
  stacking independently of router history.

A closed sheet keeps its children mounted by default. With a closed detent,
`lazy` unmounts the children after close settles and resets their local state.
Keep route focus and removal cleanup in the route owner. Check where focus
lands when the sheet closes and whether form state should reset on reopen.
[Inline sheets](https://docs.swmansion.com/react-native-bottom-sheet/docs/inline-sheets/)
· [Stacking](https://docs.swmansion.com/react-native-bottom-sheet/docs/modal-sheets/#stacking)
· [Window overlay](https://docs.swmansion.com/react-native-bottom-sheet/docs/modal-sheets/#window-overlay)
· [React context](https://docs.swmansion.com/react-native-bottom-sheet/docs/modal-sheets/#react-context)
· [Lazy content](https://docs.swmansion.com/react-native-bottom-sheet/docs/modal-sheets/#lazy-content)

## Coordinate sheet close and route removal

A modal `onCloseRequest` leaves the index under app control. Eligible scrim
taps, screen-reader actions, Android Back, and physical Escape on Android and
web can request close. Dragging to zero reaches `onIndexChange` directly.
Default closing requires positive current scrim opacity and a non-programmatic
closed detent. A custom handler must set the index and has no documented
request-source argument.
Programmatic index changes, navigator actions, and browser Back have their
own paths. A committed predictive Back can ask to close without animating
the sheet during gesture progress. Apply the same confirmation policy to the
explicit close control and eligible sheet inputs. When approval must precede
movement, make every zero-height detent programmatic-only. A blocked swipe does not
call `onCloseRequest`, so provide an accessible close control.
[Target close requests](https://docs.swmansion.com/react-native-bottom-sheet/docs/modal-sheets/#close-requests)
· [Input behavior](https://docs.swmansion.com/react-native-bottom-sheet/docs/modal-sheets/#input-behavior)
· [Programmatic detents](https://docs.swmansion.com/react-native-bottom-sheet/docs/detents-and-index/#programmatic-only-detents)

When the sheet should animate out before its route is removed, accept a
single pending close per route presentation. Keep the route mounted, select
the closed index, and commit the intended removal once at the matching
`onSettle(closedIndex)`. If drag-to-close is allowed, update the controlled
index on the drag commit and use the same pending-close path. For an
intercepted navigator removal (for example header Back, a route gesture,
browser Back, or a parent action), retain and resume the original supported
action after acceptance. Route removal can happen independently of the sheet
and may unmount it before any settle event.

Guard the commit with the same route instance and close attempt. Clear a
pending commit if the sheet reopens or the route is removed or unmounted.
Prevent repeated input, repeated confirmations, and re-entry when resuming a
previously accepted removal. A cancelled close keeps or restores the intended
open index and leaves the route open. Handle
an already closed, unmeasured, or non-animating sheet without waiting for a
new settle event. A cold link or refresh can leave no valid back destination.
Use the app's fallback location or keep a reachable route view when it cannot
go back. A closed sheet on that route can otherwise leave a blank screen.
[Target index events](https://docs.swmansion.com/react-native-bottom-sheet/docs/detents-and-index/#index-events)
· [React Navigation linking](https://reactnavigation.org/docs/configuring-links/#rendering-an-initial-route)
· [Expo Router web modal dismissal](https://docs.expo.dev/router/advanced/web-modals/#use-a-transparent-overlay-route)

### React Navigation native-stack

Use the owning navigator's appropriate action for sheet-originated removal.
`navigation.canGoBack()` can include parent history, so check that Back
implements the intended removal. With `usePreventRemove` from
`@react-navigation/native`, the callback receives `{ data }`. Continue an
accepted blocked action with `navigation.dispatch(data.action)`. A linking
configuration's `initialRouteName` can supply a screen in navigation state
without adding a browser-history entry.
[Navigation object](https://reactnavigation.org/docs/navigation-object/#cangoback)
· [React Navigation removal prevention](https://reactnavigation.org/docs/use-prevent-remove/)
· [Linking initial routes](https://reactnavigation.org/docs/configuring-links/#rendering-an-initial-route)

### Expo Router Stack

Keep `Link` and pathname-based open actions when retaining the route. Preserve
the applicable back or dismiss behavior (`router.canGoBack()` / `router.back()`
or `router.canDismiss()` / `router.dismiss()`). Use the installed SDK's
removal API. Since SDK 58, Expo Router documents `usePreventRemove` from
`expo-router` with a callback receiving `{ repeat }`. Set the prevention
condition to false before calling `repeat()` to continue the blocked action.
For a fresh action chosen after sheet-originated close, or a fallback replacing
the blocked action, set prevention false and call the hook's returned
`disablePrevention()` before the new navigation action. Use `repeat()` for the
original action and the disabling path for a different action, without also
resuming the discarded action. Keep this decision in the matching close
commit so the route can still leave after the sheet has closed.
For SDK 56 and later, use Expo Router's supported entry points rather than
app imports from external `@react-navigation/*` packages.
[Expo navigation](https://docs.expo.dev/router/basics/navigation/)
· [Expo Router API](https://docs.expo.dev/versions/v58.0.0/sdk/router/#router)
· [Expo Router removal prevention](https://docs.expo.dev/router/advanced/prevent-screen-removal/)
· [Disable removal prevention](https://docs.expo.dev/router/advanced/prevent-screen-removal/#disable-removal-prevention)
· [SDK 56 imports](https://docs.expo.dev/router/migrate/sdk-55-to-56/#manual-migration)

Use the installed SDK's anchor configuration to preserve the background route
on deep links. Current Expo Router uses `unstable_settings.anchor` in the
owning layout, including the nested stack where applicable. `withAnchor` can
also load it during in-app navigation into another layout. In SDK 58 the
navigator's `initialRouteName` prop was removed and the settings alias is
deprecated in current docs. Keep existing settings that match the app SDK.
[Deep-linked modals](https://docs.expo.dev/router/advanced/modals/#handle-deep-linked-modals)
· [Router settings](https://docs.expo.dev/router/advanced/router-settings/)
· [SDK 58 migration](https://docs.expo.dev/router/migrate/sdk-57-to-58/#move-initialroutename-to-unstable_settings)

### Web

React Navigation documents native-stack `presentation` for iOS and Android.
The current Expo web guide says `modal` and `formSheet` render as regular
stack routes and describes `transparentModal` for custom overlays. Earlier
Expo Router v6 had experimental opt-in web form sheets. Inspect the installed
SDK, web flags, and platform-specific routes before assuming the existing
behavior. The current guide records the removal of `EXPO_UNSTABLE_WEB_MODAL`,
`webModalStyle`, and the modal CSS variables. Preserve browser history and
app dismissal behavior when introducing the destination's web sheet visuals.
[React Navigation presentation](https://reactnavigation.org/docs/native-stack-navigator/#presentation)
· [Current Expo web modals](https://docs.expo.dev/router/advanced/web-modals/)
· [Earlier Expo Router v6 web modals](https://expo.dev/blog/expo-router-v6#modals-on-web)
