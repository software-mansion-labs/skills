# Imperative control

Question 6 of the walk: code that cancels, pauses, reverses or restarts the animation from an effect or a handler. Each has a CSS form; only holding a cancelled transition mid-flight does not.

| Source | CSS | Verdict |
|---|---|---|
| `cancelAnimation(sv)` in an unmount cleanup | nothing; CSS cleans up on unmount | continue |
| `cancelAnimation(sv)` then `sv.value = withTiming(other)` (a retarget) | the new state; the transition retargets from the current value | continue |
| `cancelAnimation(sv)` then `sv.value = withTiming(sameTarget)` | the shared value restarted at full duration; a transition to a target already in flight changes nothing | Needs approval |
| `cancelAnimation(sv)` that holds the current value of an animation (a stop button, a "freeze here") | `animationPlayState: 'paused'` holds the current frame | continue |
| `cancelAnimation(sv)` that holds the current value of a transition | none: nothing holds a transition mid-way, and removing `animationName` snaps to the static style (`../animations/animations.md`, Mount animations) | Needs approval, stating the snap, or Keep on shared values |
| pausing and resuming (the shared value has no pause; sites emulate it with `cancelAnimation` and a later `with*` from the current value) | `animationPlayState: 'paused'` / `'running'` on an animation | Needs approval: the emulation re-eased the rest of the distance over a full duration on resume, CSS resumes where it paused with the time that was left |
| reversing from code (`cancelAnimation` then `withTiming` back to the start) | an animation: a new `css.keyframes()` rule with `animationDirection: 'reverse'`; flipping the direction on the rule already attached only updates its settings, the running animation mirrors in place and does not restart. The reversed run starts from the end keyframe, not from the interrupted value. A transition: the state flipped back, which takes the shortened return leg (`../animations/animations.md`, CSS Transitions, Rules) | Needs approval when the forward run can be interrupted (animation) and always for the transition case; a reverse that only ever follows a finished run: continue |
| stopping a running loop and tweening to rest (`cancelAnimation(pulse); pulse.value = withTiming(1)` when loading ends) | removing `animationName` snaps to the static style; there is no tween from the current animated value | Needs approval, stating the snap, or Keep on shared values |
| restarting a finished animation (`sv.value = 0; sv.value = withTiming(1)`) | a new keyframes rule restarts the animation (`../animations/animations.md`, Defining keyframes): create the rule where the restart happened and attach the animation only once it exists, see the example; the reset write is the rule's first keyframe, not a jump for question 9; a remount or `key` change also restarts but changes the element tree | continue |
| seeking, per-frame speed changes | none | Keep on shared values |

## Example: replay on demand

```tsx
// Before
const shake = () => { x.value = 0; x.value = withSequence(withTiming(-8, { duration: 50 }), withTiming(8, { duration: 50 }), withTiming(0, { duration: 50 })); };
```

```tsx
// After: a rule created per replay, attached only once shake() has run
const [shakeName, setShakeName] = useState<ReturnType<typeof css.keyframes>>();
const shake = () => setShakeName(css.keyframes({ '33%': { transform: [{ translateX: -8 }] }, '66%': { transform: [{ translateX: 8 }] } }));
// on the element: ...(shakeName && { animationName: shakeName, animationDuration: 150, animationTimingFunction: 'ease-in-out' })
```

A `useMemo` keyed on a counter would attach the rule on the first render too, so the element would shake on mount; the state form starts with no animation. The `x.value = 0` reset is the implicit first keyframe (the static `translateX`), so question 9 does not count it. Each of the three `withTiming` segments eased on its own with the default `Easing.inOut(Easing.quad)`; a keyframe timing function applies per interval too, so one `'ease-in-out'` covers all three (0.012 error, `references/easing.md`).
