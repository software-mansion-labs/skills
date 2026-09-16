# CSS Pseudo-Selectors

Pseudo-selectors drive interaction state (press, hover, focus) from inside the style, without React state.

**Requires Reanimated 4.5.0.** Below that a pseudo object is not understood and breaks the property it sits on (a color throws at mount). Use `Pressable` + React state instead (`animations.md`, Simple gesture feedback). For the selector table and types, webfetch the [pseudo-selectors docs](https://docs.swmansion.com/react-native-reanimated/docs/css-transitions/pseudo-selectors).

Native selectors are `:hover`, `:active`, `:active-deepest`, `:focus` and `:focus-within`, and they work on any `Animated` component, no `Pressable` needed. Every other selector (`:focus-visible`, `:disabled`, `:checked`, ...) is web-only; on native its value is dropped with a dev-only warning. Use one only when the code targets web exclusively.

```tsx
import { Pressable } from 'react-native-gesture-handler';
import Animated from 'react-native-reanimated';

const AnimatedPressable = Animated.createAnimatedComponent(Pressable);

<AnimatedPressable
  onPress={onPress}
  style={{
    backgroundColor: { default: '#eee', ':active': '#ccc' },
    transform: { default: [{ scale: 1 }], ':active': [{ scale: 0.96 }] },
    transitionProperty: ['backgroundColor', 'transform'],
    transitionDuration: 150,
  }}
/>
```

Selector keys go inside each property; a top-level `':active': { ... }` block is not valid. Timing comes from the top level of the style, never from inside the pseudo object.

### Rules

- Include every pseudo-styled property in `transitionProperty`. One styled by a selector but missing from the list changes instantly.
- Write `default` unless the resting value is the property's own default. The pseudo object owns the property, so an omitted `default` falls back to that default, never to `StyleSheet.create` or an earlier style in the array.
- Leave no other writer on a property you moved into a pseudo object. From 4.6.0 the property is locked while its selector matches (a re-render or an ordinary transition cannot change it); on 4.5.x a re-render overwrites the matched value.
- Later selector wins: `:focus-within < :focus < :hover < :active < :active-deepest`.
- On `react-native-svg` elements (`:hover` and `:active`, from 4.6.0) the pseudo objects and `transition*` settings go in `animatedProps`, not `style`.

### Per-selector behavior

- `:active` matches the pressed element **and every ancestor declaring `:active`**, so a card with `:active` also reacts when a button inside it is pressed. `:active-deepest` matches only the innermost element under the finger that declares a press selector, never an ancestor: put it on a container that should react to presses on its own area but stay still while an inner control is pressed.
- `:hover` on 4.5.x native reacts only to a hovering pointer (mouse, trackpad, stylus); a finger never triggers it. From 4.6.0 a touch-down turns it on for the touched element and its ancestors declaring `:hover`, and it stays on until a later touch lands outside the element, or a touch that moved past the touch slop is released elsewhere. It has no `-deepest` variant. For press feedback use `:active`.
- `:focus` matches the element that holds focus itself: on iOS only a text input being edited, on Android and web any focusable element (a `View` with `focusable`, keyboard or D-pad focus). Put it on `createAnimatedComponent(TextInput)`; a wrapping `View` never matches, give it `:focus-within`, which matches while the element or any descendant holds focus.
- On iOS `:active` and `:active-deepest` never fire on an element whose resting `opacity` is `0.01` or lower: UIKit skips those views when hit-testing. Write `0.02`, visually identical and touchable. SVG hit-tests its own way and is exempt. From 4.6.0 a dev warning flags an `opacity` pseudo object whose `default` is that low.

### Keep the element that receives the touch

There is no `Animated.Pressable`; the namespace exports only `FlatList`, `Image`, `ScrollView`, `Text` and `View`. Use `createAnimatedComponent(Pressable)` as above. Swapping a `Pressable` for a plain `Animated.View` to get `:active` drops the press handlers and its accessibility props.

When the styled element is a descendant the finger may not land on, or the pressed value also depends on other React state, use `Pressable`'s render prop:

```tsx
<Pressable onPress={onPress}>
  {({ pressed }) => (
    <Animated.Text
      style={{
        color: pressed ? '#000' : '#888',
        transitionProperty: 'color',
        transitionDuration: 150,
      }}>
      Press me
    </Animated.Text>
  )}
</Pressable>
```
