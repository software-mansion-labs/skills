# Callbacks

Question 7 of the walk. `with*` completion callbacks map onto the `onCSS*` props from Reanimated 4.6.0 (`../animations/animations.md`, Callbacks). Below 4.6.0 an observable callback keeps the site on shared values. A log-only callback is dropped without a row.

The event props, their payloads and timing are in `../animations/animations.md`, Callbacks; they are component props, not style keys. With several animations on one element, branch on the event's `animationName`, comparing it against the `.name` of the rule object you pass in `animationName`: a `css.keyframes()` rule from module scope, or the restart rule held in state (`references/imperative-control.md`). A plain keyframes object gets a generated name you cannot reference. When a replay swaps the rule, the Cancel event of the run it interrupts carries the previous rule's name.

## Mapping

| Source | CSS | Notes |
|---|---|---|
| `withTiming(v, cfg, cb)` on a transitioned property | `onCSSTransitionEnd` for `cb(true)`, `onCSSTransitionCancel` for `cb(false)` | one handler serves several properties and every direction: branch on `propertyName`, and on the current target when only one direction had a callback |
| `withTiming(v, cfg, cb)` on a keyframe animation | `onCSSAnimationEnd` for `cb(true)`, `onCSSAnimationCancel` for `cb(false)` | |
| `withSequence(a, withTiming(v, cfg, cb))` (callback on the last child) | `onCSSAnimationEnd` | a callback on an earlier child: no per-keyframe event; Needs approval, proposing a timer at that step's offset, else Keep on shared values |
| `withRepeat(withTiming(v, cfg, innerCb), n)` | `onCSSAnimationIteration` for every repetition except the last, which is `onCSSAnimationEnd`; `onCSSAnimationCancel` for the `false` call on cancel | `innerCb` ran after each repetition, including the last: handle both events. Under reduced motion a reversed even or infinite repeat never ran `innerCb` at all, while the CSS End still fires once |
| `withRepeat(anim, n, reverse, outerCb)` | `onCSSAnimationEnd` | `outerCb` ran once after the last repetition; with `n <= 0` it ran only with `false` on cancel (`onCSSAnimationCancel`) and with `true` under reduced motion, which the `animationIterationCount: reduced ? 1 : 'infinite'` form reaches as `onCSSAnimationEnd` |
| `withDelay(ms, withTiming(v, cfg, cb))` | as the inner animation | |
| `cb(false)` from `cancelAnimation` or a retarget | `onCSSTransitionCancel` / `onCSSAnimationCancel` | a retarget of a running transition cancels the old one and runs a new one, so Cancel then Run/Start fire |
| `cb(true)` when the target already equals the current value | nothing: CSS fires no event when nothing changes | call the handler directly in that branch: `if (next === current) onDone(); else setValue(next)`, and say so in the row. The shared value fired it at once only for a bare `withTiming`; wrapped in `withDelay`, `withSequence` or `withRepeat` it waited the whole timeline first |
| the same target written again mid-flight | the running transition is left alone and fires one End | the shared value fired the first callback with `false` and the second with `true`; Needs approval when the site counts on the `false` call |
| a callback that writes another shared value (`withTiming(0, cfg, () => { other.value = withTiming(1); })`) | the `onCSS*` handler runs on the JS thread | set state when the chained site migrates too; otherwise assign the shared value from the handler, one frame later than the shared value did |
| a callback that reads the shared value it animated | read the React state that replaced it; the End event carries no value | |

Callbacks fire for pseudo-selector driven transitions too. A transition removed by a zero effective duration (4.3.0+) fires nothing, so a reduced-motion guard must shorten to `1`, never `0` (`references/reduced-motion.md`).

## Examples

Chained animation:

```tsx
// Before
opacity.value = withTiming(0, { duration: 200 }, (finished) => {
  if (finished) scheduleOnRN(setVisible, false);
});
```

```tsx
// After: the property transitions to the new target; the End event chains the state change
<Animated.View
  onCSSTransitionEnd={({ propertyName }) => {
    if (propertyName === 'opacity' && hidden) setVisible(false);
  }}
  style={[styles.box, { opacity: hidden ? 0 : 1, transitionProperty: 'opacity', transitionDuration: 200, transitionTimingFunction: 'ease-in-out' }]}
/>
```

Needs approval row: setting `hidden` while it is already `true` fired `finished: true` on the shared value and fires nothing here; the row proposes `if (hidden) setVisible(false); else setHidden(true)`. The row also records `inOut(quad)` to `'ease-in-out'` (0.012). The End handler fires after the fade-in too, hence the `hidden` check.

Repetition counter:

```tsx
// Before
progress.value = withRepeat(withTiming(1, { duration: 500 }, () => scheduleOnRN(tick)), 3);
```

```tsx
// After: tick ran three times, so both events call it; pulse is a css.keyframes() rule from module scope
<Animated.View
  onCSSAnimationIteration={tick}
  onCSSAnimationEnd={({ animationName }) => { if (animationName === pulse.name) tick(); }}
  style={{ animationName: pulse, animationDuration: 500, animationIterationCount: 3, animationTimingFunction: 'ease-in-out', animationFillMode: 'forwards' }}
/>
```

Under reduced motion the shared value ran one repetition and called `tick` once, and so does the reduced-motion form from `references/reduced-motion.md` (one End, no Iteration). The row records `inOut(quad)` to `'ease-in-out'` (0.012).
