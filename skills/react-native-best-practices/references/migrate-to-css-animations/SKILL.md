---
name: migrate-to-css-animations
description: "Migrates React Native Reanimated hook animations (useAnimatedStyle, useSharedValue, withTiming, withRepeat, withSequence, withDelay) to CSS transitions and animations, applying only conversions that keep behavior identical and explaining why the rest stay on hooks. Use when asked to migrate, convert, port or audit Reanimated animation code for CSS, or to simplify, modernize or drop worklets from animations. Requires Reanimated 4.x; reads the installed version and emits only what that version supports."
---

# Migrate Reanimated hook animations to CSS

Convert only where behavior stays identical; coverage is not the goal. A conversion that typechecks and runs without errors can still animate differently, so neither proves it is correct: the decision tree below and the checks after converting do.

"Hook animations" here means what `../animations/animations.md` calls shared value animations. For the CSS API itself read `../animations/animations.md` (feature availability by version, transitions, animations, timing functions, callbacks) and, from 4.5.0, `../animations/css-pseudo-selectors.md`. Load `references/mapping.md` when converting a site and `references/examples.md` to calibrate output.

## 0. Version and scope

- Installed Reanimated: `node_modules/react-native-reanimated/package.json` or the lockfile, never the `package.json` range. 3.x has no CSS API: stop and say so, do not upgrade. On 4.x note the exact version; the feature table in `../animations/animations.md` decides what may be emitted, and a site that needs a feature the installed version lacks stays on hooks.
- Scope: the files or directories the user named. With none named, the project's own source (not `node_modules`); say so when reporting the inventory.

| Request | Steps |
|---|---|
| one site | 0, 2, 3 (or explain why it stays) |
| "can this be CSS?", an audit | 0, 1, 2, 5, no edits |
| a directory or the app | all |

## 1. Inventory, then ask once

Record every `useAnimatedStyle`, `useAnimatedProps`, `useDerivedValue` and `with*` site in scope: file, component, animated properties, driver. Report the counts. Nothing animated in scope: say so and stop.

Then ask once, in one message: confirm the scope, and whether a partial result is wanted when only some sites can move (migrate those and leave the rest, or change nothing unless every site moves). The answers hold for the whole run. Every other judgment call is a Needs approval row in the report, asked in one batch after classification.

## 2. Decide per site

Three verdicts:

| Verdict | Meaning |
|---|---|
| Migrate | behavior identical; apply |
| Needs approval | behavior changes in a way the user may accept, or the tree does not cover the site; show the proposed code, apply only after a yes |
| Keep on hooks | CSS cannot express it; leave it, give the reason in one sentence |

Walk the questions in order for every `useAnimatedStyle`/`useAnimatedProps` hook. Keep on hooks ends the walk; Needs approval is noted and the walk continues; a hook that reaches the end with nothing noted is Migrate, and with notes it is Needs approval. One hook is all or nothing: a property that must stay keeps the whole hook (two hooks on one component may get different verdicts). A pattern the questions do not cover that still looks convertible: Needs approval, with the proposal and a note that the tree does not cover it.

```
1. Is the value written per frame or on the UI thread?
   (scroll handler, gesture onUpdate/onChange, gesture callbacks without .runOnJS(true)
    in Gesture Handler 2 and 3, useAnimatedReaction, useFrameCallback, sensors,
    useAnimatedKeyboard, a scheduleOnUI/runOnUI body, a JS loop assigning every frame)
   |-- YES -> Keep on hooks
   `-- NO  -> continue: the driver is a JS-thread write (useEffect, onPress/onChange,
             a timer, a network callback, gesture onEnd/onFinalize with .runOnJS(true)).
             Writer not traceable in scope -> note Needs approval, continue.

2. Is it withSpring or withDecay, or withClamp around one of them?
   |-- YES -> Keep on hooks (no CSS spring or decay)
   `-- NO  -> continue. withClamp around withTiming: both endpoints inside the bounds ->
             drop the clamp; an endpoint outside -> Keep on hooks.

3. Does CSS animate every property on every platform the project targets, at the
   installed version? (references/mapping.md, Properties)
   |-- NO  -> Keep on hooks. A discrete (keyword) property -> note Needs approval,
   |          proposing transitionBehavior: 'allow-discrete' and saying when it flips.
   `-- YES -> continue

4. Is each animated value an affine function of the driver?
   (a * driver + b, or interpolate with fixed stops; references/mapping.md, Value functions)
   |-- NO  -> transition: Keep on hooks (a transition only tweens two endpoints);
   |          animation: note Needs approval, sampling the function into keyframes.
   `-- YES -> continue

5. Does the timing curve map to CSS exactly? (references/mapping.md, Easing)
   |-- exact row, or the user accepted linear() sampling -> continue
   `-- approximation only -> note Needs approval with the max error, continue

6. Does code cancel, pause, reverse or restart the animation?
   (references/mapping.md, Imperative control)
   |-- control CSS cannot express -> Keep on hooks
   |-- restart to the same target, pause, reverse from code -> note Needs approval, continue
   `-- none, an unmount cleanup, or a retarget -> continue

7. Does a completion callback do something observable? (references/mapping.md, Callbacks)
   |-- YES, below 4.6.0 -> Keep on hooks
   |-- YES, 4.6.0+ -> note Needs approval with the mapped callback, continue
   `-- NO -> continue

8. Can the driver reverse mid-flight? (press in/out, a toggle the user flips)
   |-- YES -> note Needs approval (references/mapping.md, Reversal), continue
   `-- NO  -> continue

9. Does anything else set the same property without a with*?
   (a direct sv.value = x, a mount jump, a conditional style)
   |-- YES -> note Needs approval: every change of a transitioned property animates,
   |          so the jump becomes a transition; continue
   `-- NO  -> continue

10. Do several hooks read one looping shared value?
   |-- YES -> note Needs approval (references/mapping.md, Transition or animation)
   `-- NO  -> Migrate, or Needs approval if anything was noted
```

Outside the tree, leave these alone: `entering`/`exiting`/`layout` animations, `Keyframe`, shared element transitions (a separate subsystem; judge a `useAnimatedStyle` on the same element on its own), and a `useAnimatedProps` for a non-style prop on a non-SVG component (`text`, `contentOffset`, `scrollEnabled`), which CSS never covers.

Indirection: a `useDerivedValue`, or a `SharedValue` passed as a prop, is not a driver. Follow it to the shared value that is written and classify that write; the derived computation joins the "affine function" question. A derived value with a worklet reader or writer of its own stays on hooks.

Measured layout (`onLayout`, `measure` in an effect) stored in React state is a plain state driver: a transition to the measured target is fine. `measure` read inside the hook every frame is the per-frame case at the top.

Reduced motion: `useReducedMotion()` returns the system setting read once when the app starts; `with*` follow a `<ReducedMotionConfig>` live. No config: guard the CSS site with `useReducedMotion()` (`references/mapping.md`, Reduced motion). A literal `mode` of `Always` or `Never` on the config: apply the matching row and no guard. A `mode` that changes at runtime: Keep on hooks, CSS has no live reduced-motion source.

## 3. Convert

`references/mapping.md` has the `with*`, easing, callback, reduced-motion, color and SVG tables. Rules that hold for every site:

- Follow the transition and animation Rules in `../animations/animations.md`; above all list `transitionProperty` explicitly (never `'all'` or the `transition` shorthand string) and write the timing function. The `withTiming` default is `Easing.inOut(Easing.quad)` at 300ms, not the CSS default `'ease'`; `references/mapping.md`, Easing, says how to map it.
- The driver becomes React state (`useState`, a prop, a store value); this is the migration, not a cost.
- Migrate inside each `Platform.select` arm and keep the structure; enumerate `transitionProperty` per platform when the property set differs.
- A property the hook returns from a value that is never animated (a constant, or a shared value that is never written after its initial value) is static: move it to the static style. Remove the shared values, hooks and imports the conversion killed and nothing else.
- Never put CSS declarations on a plain component.
- If converting uncovers something the inventory missed (a worklet writer after all, a property the docs table does not list, a check that cannot be satisfied), revert that site and classify it again.

After converting a site, confirm each of these against the original:

- first render identical: the static style carries the value the hook painted first (`../animations/animations.md`, Mount animations);
- end state identical, including the value the site rests at under reduced motion;
- re-trigger identical: writing the same target mid-flight is identical unless the original cancelled first; replaying a finished animation needs a remount or a `key` change, which is a visible structural change, so Needs approval;
- the callbacks the original fired still fire, at the same moments;
- unmount mid-animation throws nothing;
- nothing else changed: same element tree, props and handlers (a `Pressable` swapped for `Animated.View` fails this), and every `Platform.select` arm converted.

Migrate three sites of differing shape first, compare them for consistent treatment, then continue. Keep a status file outside the repo (site, verdict, reason) and trust it plus `git diff` after an interruption.

## 4. Verify

Typecheck and lint the changed files; run their tests (replacing a `useAnimatedStyle` result with plain style props changes snapshots and style assertions over that subtree, read the diff to confirm only the intended properties moved). When a simulator, emulator or device is available, run the app and watch the first frame and one re-trigger of each migrated site, where a missing fill mode or a dead driver shows. Say what you verified and what you did not.

## 5. Report

Lead with the numbers: `Migrated 14 sites across 9 files. Left 23: 6 need approval, 17 stay on hooks.` Then:

1. **Applied**: one row per site (file, what it animates, transition or animation, note). The note lists every behavior delta or says `exact`: easing substituted and its max error, reduced-motion guard dropped, a reversal now shorter, a color midpoint moved, a mount step or replay that no longer fires.
2. **Needs approval**: one row per site with the open question and the proposed code. Ask which to apply, by number.
3. **Kept on hooks**, grouped by reason with counts.
4. Before/after code for two or three applied sites of different shape.

Up to about ten sites per group, print every row. Above that, print the counts and the grouped reasons, then ask whether to list all rows, the first ten, or one file at a time. Leave out recipes, memoization advice and per-site prose.
