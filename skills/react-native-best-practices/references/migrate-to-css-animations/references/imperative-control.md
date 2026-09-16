# Imperative control

Question 6 of the walk: code that cancels, pauses, reverses or restarts the animation from an effect or a handler. Each has a CSS form; only holding a cancelled value mid-flight does not.

| Source | CSS | Verdict |
|---|---|---|
| `cancelAnimation(sv)` in an unmount cleanup | nothing; CSS cleans up on unmount | continue |
| `cancelAnimation(sv)` then `sv.value = withTiming(other)` (a retarget) | the new state; the transition retargets from the current value | continue |
| `cancelAnimation(sv)` then `sv.value = withTiming(sameTarget)` | the shared value restarted at full duration; a transition to a target already in flight changes nothing | Needs approval |
| `cancelAnimation(sv)` that holds the current value (a stop button, a "freeze here") | none: removing `animationName` snaps to the static style (`../animations/animations.md`, Mount animations) and a transition cannot be stopped mid-way | Needs approval, stating the snap, or Keep on shared values |
| pausing and resuming (the shared value has no pause; sites emulate it with `cancelAnimation` and a later `with*` from the current value) | `animationPlayState: 'paused'` / `'running'` on an animation | continue |
| reversing from code (`cancelAnimation` then `withTiming` back to the start) | an animation: `animationDirection: 'reverse'` restarts it from the end; a transition: the state flipped back, which takes the shortened return leg (`../animations/animations.md`, CSS Transitions, Rules) | Needs approval for the transition case |
| stopping a running loop and tweening to rest (`cancelAnimation(pulse); pulse.value = withTiming(1)` when loading ends) | removing `animationName` snaps to the static style; there is no tween from the current animated value | Needs approval, stating the snap, or Keep on shared values |
| restarting a finished animation (`sv.value = 0; sv.value = withTiming(1)`) | a new keyframes rule restarts the animation: `animationName: useMemo(() => css.keyframes(frames), [replayCount])` (`../animations/animations.md`, Defining keyframes); a remount or `key` change also works but changes the element tree | continue |
| seeking, per-frame speed changes | none | Keep on shared values |

## Example: replay on demand

```tsx
// Before
const shake = () => { x.value = 0; x.value = withSequence(withTiming(-8, { duration: 50 }), withTiming(8, { duration: 50 }), withTiming(0, { duration: 50 })); };
```

```tsx
// After: a new rule per replay
const [shakeCount, setShakeCount] = useState(0);
const shakeName = useMemo(() => css.keyframes({ '33%': { transform: [{ translateX: -8 }] }, '66%': { transform: [{ translateX: 8 }] } }), [shakeCount]);
// on the element: animationName: shakeName, animationDuration: 150, animationTimingFunction: 'ease-in-out'
// shake() becomes setShakeCount((c) => c + 1)
```
