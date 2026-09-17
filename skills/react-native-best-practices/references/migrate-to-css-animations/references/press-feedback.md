# Reversal: press, hover, focus and toggles

Question 8 of the walk. A site driven by `onPressIn`/`onPressOut`, `onHoverIn`/`onHoverOut`, `onFocus`/`onBlur` or a gesture's begin/end pair is interaction state, not application state. Two migration shapes exist, and both differ from the shared value version in one way: a CSS transition reversed mid-flight takes the shortened return leg (`../animations/animations.md`, CSS Transitions, Rules), while the shared value took whatever duration the site gave the return write, usually the full one. Every such site is therefore Needs approval, with the difference stated for the site's duration: the return takes the duration times the eased progress at the release, so a 200ms `'ease-in-out'` fade released 80ms in returns in about 66ms instead of 200ms (80ms with `'linear'`).

## Which shape

```
Is the styled element the one the finger or pointer lands on?
(the AnimatedPressable or AnimatedTextInput itself, or a child that fills it)
|-- YES -> Reanimated 4.5.0+ (4.6.0 for a react-native-svg element, where the pseudo
|          objects go in animatedProps: only that path injects the SVG hit-test responder
|          on native, style works on web only): a pseudo-selector
|          (../animations/css-pseudo-selectors.md)
|          Below 4.5.0: React state from the handlers on the same element
`-- NO  -> the feedback is on a descendant or an ancestor of the pressed element:
           React state from the handlers, or a Pressable render prop styling the child
           (../animations/animations.md, Simple gesture feedback)
```

| Hook driver | Selector | Notes |
|---|---|---|
| `onPressIn` / `onPressOut`, gesture `onBegin` / `onFinalize` with `.runOnJS(true)` | `:active` (`:active-deepest` for a container that must stay still while an inner control is pressed) | |
| `onHoverIn` / `onHoverOut` | `:hover` | |
| `onFocus` / `onBlur` on a `TextInput` | `:focus` | |
| focus of any descendant | `:focus-within` | |
| a pressed state that stays after release (selected, toggled) | none; it is application state | a transition driven by `useState` |
| a gesture `onUpdate` value (drag distance, scale) | none | continuous input: stays on shared values |

What each selector matches per platform, and which other pseudo-classes are web only, is in `../animations/css-pseudo-selectors.md`; a site that needs a web-only one stays on state.

## Rules that carry over

- Write `default` for every property the selector changes when the property also has a value elsewhere: a pseudo object without `default` rests at the property's built-in default, not at the value it replaced (`../animations/css-pseudo-selectors.md`, Rules).
- Keep the transition config on the element (`transitionProperty`, `transitionDuration`, `transitionTimingFunction`); a selector only switches the value.
- Keep the `Pressable` and its handlers: `onPress` still handles the action, only the visual handlers go. Swapping the `Pressable` for an `Animated.View` changes the element tree and fails the post-conversion checks.
- Below 4.5.0 keep `onPressIn`/`onPressOut` and write `useState` instead of the shared value; the transition then follows the state.

## Examples

Press scale, Reanimated 4.5.0+:

```tsx
// Before
const scale = useSharedValue(1);
const style = useAnimatedStyle(() => ({ transform: [{ scale: scale.value }] }));
<AnimatedPressable
  onPressIn={() => { scale.value = withTiming(0.96, { duration: 120 }); }}
  onPressOut={() => { scale.value = withTiming(1, { duration: 120 }); }}
  onPress={onPress}
  style={[styles.button, style]}
/>
```

```tsx
// After: the Pressable keeps onPress; the visual handlers and the shared value go
<AnimatedPressable
  onPress={onPress}
  style={[
    styles.button,
    {
      transform: { default: [{ scale: 1 }], ':active': [{ scale: 0.96 }] },
      transitionProperty: 'transform',
      transitionDuration: 120,
      transitionTimingFunction: 'ease-in-out',
    },
  ]}
/>
```

Needs approval row: release 60ms into the 120ms press-in now returns in about 60ms instead of 120ms; the row also records `inOut(quad)` to `'ease-in-out'` (0.012).

Highlight on a child of the pressed row, any 4.x:

```tsx
// Before: opacity on the icon, driven by the row's press handlers
```

```tsx
// After: state from the same handlers; the icon transitions
const [pressed, setPressed] = useState(false);
<Pressable onPressIn={() => setPressed(true)} onPressOut={() => setPressed(false)} onPress={onPress}>
  <Animated.View style={[styles.icon, { opacity: pressed ? 0.5 : 1, transitionProperty: 'opacity', transitionDuration: 120, transitionTimingFunction: 'ease-in-out' }]} />
</Pressable>
```

Focus ring on a text input, 4.5.0+:

```tsx
<AnimatedTextInput
  style={[
    styles.input,
    {
      borderColor: { default: '#ccc', ':focus': '#1e3a8a' },
      transitionProperty: 'borderColor',
      transitionDuration: 150,
      transitionTimingFunction: 'ease-in-out',
    },
  ]}
/>
```

## Other toggles

Expand/collapse, show/hide, a switch driven by state: the transition, with the shortened return stated in the row (reversed 80ms into a 200ms `'ease-in-out'` transition, the shared value took the 200ms the site gave it and CSS takes about 66ms). A driver behind a timer longer than the duration, or a rare event (rotation, navigation, network): Migrate and state the difference. A site that became an animation at question 4 or 6 does not retarget: a flip mid-flight starts the incoming rule from its first keyframe, where the shared value tweened back from the current value; Needs approval, stating that jump.
