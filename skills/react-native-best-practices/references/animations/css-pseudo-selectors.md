# CSS Pseudo-Selectors

Pseudo-selectors drive interaction state (press, hover, focus) from inside the style, without React state. They work on any `Animated` component (and on `react-native-svg` elements from 4.6.0, declared in `animatedProps`). Docs: [pseudo-selectors](https://docs.swmansion.com/react-native-reanimated/docs/css-transitions/pseudo-selectors).

**Requires Reanimated 4.5.0.** Below that the pseudo object is not understood; use `Pressable` + React state instead (`animations.md`, Simple gesture feedback).

Supported everywhere: `:hover`, `:active`, `:active-deepest`, `:focus` and `:focus-within`. The remaining CSS pseudo-classes (`:focus-visible`, `:disabled`, `:checked`, ...) work on web only; on native they are dropped with a dev-only warning.

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

Selector keys go inside each property; a top-level `':active': { ... }` block is not valid. Pseudo-selectors style the pressed element itself. To style its descendants, the `Pressable` itself below 4.5.0, or an ancestor, use the `Pressable`-based approaches described under Simple gesture feedback in `animations.md`.

### Rules

- A selector only switches the value. The switch animates when the property has a transition (`transitionProperty` and a duration, as for any transition); without one it snaps.
- `default` is the resting value. Leaving it out means the selector animates back to the property's own default value, for example transparent for a color or `1` for `opacity`; the pseudo object owns the property, so a value from an earlier entry of the style array is not used.
- While a selector matches, its values win over re-renders and regular transitions of the same property. A CSS animation targeting that property still overrides it. Do not drive the same property with an animated style as well; the two do not compose.
- When several selectors match at once, the one further right in this fixed order wins, whatever the order of keys in the object: `:focus-within < :focus < :hover < :active < :active-deepest`.
- The transition callbacks (`onCSSTransition*`, `animations.md`) fire for selector-driven transitions too.

### Per-selector behavior

- `:active` matches the pressed element **and every ancestor declaring `:active`**, so a card with `:active` also reacts when a button inside it is pressed. `:active-deepest` matches only the innermost element under the finger that declares a press selector, never an ancestor: put it on a container that should react to presses on its own area but stay still while an inner control is pressed.
- `:hover` follows the pointer on mouse, trackpad and stylus. On touch a touch-down turns it on for the touched element and its ancestors declaring `:hover`, and it stays on until a later touch lands outside the element, or a touch that moved past the touch slop is released elsewhere. It has no `-deepest` variant. For press feedback use `:active`.
- `:focus` matches the element that holds focus itself: on iOS only a text input being edited, on Android and web any focusable element (a `View` with `focusable`, keyboard or D-pad focus). Put it on `createAnimatedComponent(TextInput)`; a wrapping `View` never matches, give it `:focus-within`, which matches while the element or any descendant holds focus.
- On iOS `:active` and `:active-deepest` never fire on an element whose resting `opacity` is `0.01` or lower: UIKit skips those views when hit-testing. Write `0.02`, visually identical and touchable.

### 4.5.x differences

Two 4.5.x behaviors were fixed in 4.6.0: `:hover` reacted only to a hovering pointer (a finger never triggered it), and a re-render could overwrite a matched selector's value (no property lock). Do not rely on either on 4.5.x.
