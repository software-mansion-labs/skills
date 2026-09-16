---
name: migrate-to-css-animations
description: "Migrates React Native Reanimated hook animations (useAnimatedStyle, useSharedValue, withTiming, withRepeat, withSequence, withDelay) to CSS transitions and animations, applying only conversions that keep behavior identical and explaining why the rest stay on hooks. Use when asked to migrate, convert, port or audit Reanimated animation code for CSS, or to simplify, modernize or drop worklets from animations. Requires Reanimated 4.x; reads the installed version and emits only what that version supports."
---

# Migrate Reanimated hook animations to CSS

Convert only where behavior stays identical; coverage is not the goal. A conversion that typechecks and runs without errors can still animate differently, so neither proves it is correct: the walk below and the checks after converting do.

"Hook animations" here means what `../animations/animations.md` calls shared value animations. For the CSS API itself read `../animations/animations.md` (feature availability by version, transitions, animations, timing functions, callbacks) and, from 4.5.0, `../animations/css-pseudo-selectors.md`. Each question of the walk names the reference file that holds its rule, mapping and example; load only the files the site needs.

## 0. Version and scope

- Installed Reanimated: `node_modules/react-native-reanimated/package.json` or the lockfile, never the `package.json` range. 3.x has no CSS API: stop and say so, do not upgrade. On 4.x note the exact version; the feature table in `../animations/animations.md` decides what may be emitted, and a site that needs a feature the installed version lacks stays on hooks.
- Scope: the files or directories the user named. With none named, the project's own source (not `node_modules`); say so when reporting the inventory.

| Request | Steps |
|---|---|
| one site | 0, 2, 3, 4, 5 in short form (or explain why it stays); ask the easing question from step 1 only when step 5 needs it |
| "can this be CSS?", an audit | 0, 1, 2, 5, no edits |
| a directory or the app | all |

## 1. Inventory, then ask once

Record every `useAnimatedStyle`, `useAnimatedProps`, `useDerivedValue` and `with*` site in scope: file, component, animated properties, driver. Report the counts. Nothing animated in scope: say so and stop.

Then ask once, in one message: confirm the scope; whether a partial result is wanted when only some sites can move (migrate those and leave the rest, or change nothing unless every site moves); and, only when the inventory found a timing curve with no exact CSS form, whether to sample it into `linear()` or use the nearest approximation (`references/easing.md`). The answers hold for the whole run. Every other judgment call is a Needs approval row in the report, asked in one batch after classification.

## 2. Decide per site

| Verdict | Meaning |
|---|---|
| Migrate | behavior identical; apply |
| Needs approval | behavior changes in a way the user may accept, or the walk does not cover the site; show the proposed code, apply only after a yes |
| Keep on hooks | CSS cannot express it; leave it, give the reason in one sentence |

Walk the questions in order for every `useAnimatedStyle`/`useAnimatedProps` hook. Keep on hooks ends the walk; Needs approval is noted and the walk continues; a hook that reaches the end with nothing noted is Migrate, and with notes it is Needs approval. One hook is all or nothing: a property that must stay keeps the whole hook (two hooks on one component may get different verdicts). A pattern the questions do not cover that still looks convertible: Needs approval, with the proposal and a plain statement that the walk does not cover it and you are not sure the behavior is identical.

```
1. Is the value written per frame or on the UI thread?      references/drivers.md
   |-- YES -> Keep on hooks
   `-- NO  -> continue (a JS-thread write, or a prop/state read in the hook body)

2. Is it withSpring or withDecay anywhere in the composition, or withClamp?
   |-- YES -> Keep on hooks (no CSS spring, decay or clamp)
   `-- NO  -> continue

3. Does CSS animate every property on every platform the project targets, at the
   installed version? (supported-properties docs, feature table in animations.md)
   |-- NO, the property is a keyword -> note Needs approval, continue
   |-- NO, for any other reason -> Keep on hooks
   `-- YES -> continue

4. Is each animated value a * driver + b between two endpoints of the same kind?
   |-- NO  -> Keep on hooks
   `-- YES -> continue

5. Does the timing curve map to CSS exactly?                references/easing.md
   |-- exact row -> continue
   `-- no exact row -> apply the answer from step 1, record the substitution and its
       max error in the site's row (no Needs approval note for this alone), continue

6. Does code cancel, pause, reverse or restart the animation?
   |-- YES, except cancelAnimation in an unmount cleanup -> Keep on hooks
   `-- NO  -> continue

7. Does a completion callback do something observable?
   |-- YES -> Keep on hooks
   `-- NO  -> continue

8. Can the driver reverse mid-flight? (press in/out, a toggle the user flips)
   |-- YES -> note Needs approval: a CSS transition reversed mid-flight takes a
   |          shortened return leg (../animations/animations.md, CSS Transitions, Rules),
   |          the hook played the full duration back; continue
   `-- NO  -> continue

9. Does anything else write the shared value without a with*?  references/drivers.md
   |-- YES, and nothing writes it with a with* either -> the site never animated:
   |          plain state in the static style, no transitionProperty; continue
   |-- YES -> note Needs approval: every change of a transitioned property animates,
   |          so the jump becomes a transition; continue
   `-- NO  -> continue

10. Does anything else read the shared value?                references/drivers.md
   |-- YES -> Keep on hooks, or note Needs approval with a state mirror
   `-- NO  -> Migrate, or Needs approval if anything was noted
```

Outside the walk, leave these alone: `entering`/`exiting`/`layout` animations, `Keyframe`, shared element transitions (a separate subsystem; judge a `useAnimatedStyle` on the same element on its own), and a `useAnimatedProps` for a non-style prop on a non-SVG component (`text`, `contentOffset`, `scrollEnabled`), which CSS never covers.

## 3. Convert

`references/transitions-and-animations.md` maps `with*` compositions to transitions or animations; `references/easing.md` the timing curve; `references/reduced-motion.md` the guard every site gets. Rules that hold for every site:

- Follow the transition and animation Rules in `../animations/animations.md`; above all list `transitionProperty` explicitly (never `'all'` or the `transition` shorthand string) and write the timing function (`references/easing.md` says why the default is never left in place).
- The driver becomes React state (`useState`, a prop, a store value); this is the migration, not a cost.
- Migrate inside each `Platform.select` arm and keep the structure; enumerate `transitionProperty` per platform when the property set differs. A duration or easing computed per trigger (`duration: Math.abs(delta) * k`) becomes state set in the same render as the target, with the formula unchanged.
- A property the hook returns from a value that is never animated (a constant, or a shared value that is never written after its initial value) is static: move it to the static style. Remove the shared values, hooks and imports the conversion killed and nothing else.
- Never put CSS declarations on a plain component.
- If converting uncovers something the inventory missed (a worklet writer after all, a property the docs table does not list, a question that cannot be answered), revert that site and classify it again.

After converting a site, confirm each of these against the original:

- first render identical: the static style carries the value the hook painted first (`../animations/animations.md`, Mount animations);
- end state identical, including the value the site rests at under reduced motion (`references/reduced-motion.md`);
- re-trigger identical: writing the same target mid-flight looks identical unless the original cancelled first; replaying a finished animation needs a new keyframes rule (`useMemo(() => css.keyframes(frames), [replayCount])`, `../animations/animations.md`, Defining keyframes) or a remount, so Needs approval;
- the callbacks the original fired still fire, at the same moments;
- unmount mid-animation throws nothing;
- nothing else changed: same element tree, props and handlers (a `Pressable` swapped for `Animated.View` fails this), and every `Platform.select` arm converted.

Migrate three sites of differing shape first, compare them for consistent treatment, then continue. Keep a status file outside the repo (site, verdict, reason) and trust it plus `git diff` after an interruption.

## 4. Verify

Typecheck and lint the changed files; run their tests (replacing a `useAnimatedStyle` result with plain style props changes snapshots and style assertions over that subtree, read the diff to confirm only the intended properties moved). When a simulator, emulator or device is available, run the app and watch the first frame and one re-trigger of each migrated site, where a missing fill mode or a dead driver shows. Say what you verified and what you did not.

## 5. Report

Lead with the numbers: `Migrated 14 sites across 9 files. Left 23: 6 need approval, 17 stay on hooks.` Then:

1. **Applied**: one row per site (file, what it animates, transition or animation, note). The note lists every behavior delta or says `exact`: easing substituted and its max error, reduced-motion guard dropped, a reversal now shorter, a mount step or replay that no longer fires.
2. **Needs approval**: one row per site with the open question and the proposed code. Ask which to apply, by number.
3. **Kept on hooks**, grouped by reason with counts.
4. Before/after code for two or three applied sites of different shape.

Up to about ten sites per group, print every row. Above that, print the counts and the grouped reasons, then ask whether to list all rows, the first ten, or one file at a time. Leave out recipes, memoization advice and per-site prose.
