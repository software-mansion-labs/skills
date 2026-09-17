---
name: migrate-to-css-animations
description: "Migrates React Native Reanimated shared value animations (useAnimatedStyle, useSharedValue, withTiming, withRepeat, withSequence, withDelay) to CSS transitions and animations, applying only conversions that keep behavior identical and explaining why the rest stay on shared values. Use when asked to migrate, convert, port or audit Reanimated animation code for CSS, or to simplify, modernize or drop worklets from animations. Requires Reanimated 4.x; reads the installed version and emits only what that version supports."
---

# Migrate shared value animations to CSS

Convert only where behavior stays identical; coverage is not the goal. A conversion that typechecks and runs without errors can still animate differently, so neither proves it is correct: the walk below and the checks after converting do.

For the CSS API itself read `../animations/animations.md` (feature availability by version, transitions, animations, timing functions, callbacks) and, from 4.5.0, `../animations/css-pseudo-selectors.md`. Each question of the walk names the reference file that holds its rule, mapping and example; load only the files the site needs.

## 0. Version and scope

- Installed Reanimated: `node_modules/react-native-reanimated/package.json` or the lockfile, never the `package.json` range. 3.x has no CSS API: stop and say so, do not upgrade. On 4.x note the exact version; the feature table in `../animations/animations.md` decides what may be emitted, and a site that needs a feature the installed version lacks stays on shared values.
- Scope: the files or directories the user named. With none named, the project's own source (not `node_modules`); say so when reporting the inventory.

What the request asks for decides how much of this skill runs:

- One site ("migrate this spinner"): check the version, walk that site, convert it or explain in one paragraph why it stays. Ask the step 1 questions only if the site raises them.
- A question or an audit ("can this be CSS?", "what could move?"): inventory, walk every site, report; change nothing.
- A directory or the whole app: every step below.

## 1. Inventory, then ask once

Record every animation site in scope: `useAnimatedStyle`, `useAnimatedProps`, `useDerivedValue`, a shared value passed directly in `style` or as a prop (`style={{ opacity: sv }}`, `<AnimatedCircle r={r} />`), and every `with*` call. For each: file, component, animated properties, where the values come from. Report the counts. Nothing animated in scope: say so and stop.

Then ask once, in one message:

- Scope: confirm the files you will touch.
- Partial result: when only some sites can move, migrate those and leave the rest, or change nothing unless every site moves?
- Reduced motion: `with*` animations follow the system Reduce Motion setting on their own; CSS ignores it. Keep that behavior (every migrated site gets a `useReducedMotion()` guard, `references/reduced-motion.md`) or drop it (users with Reduce Motion on will see the animations)? Skip the question when the source already settles it for every site: `ReduceMotion.Never` or `ReduceMotion.Always` on each `with*`, or a `<ReducedMotionConfig>` whose `mode` is a fixed `Never` or `Always`. A `mode={ReduceMotion.System}` config is the same as no config (ask); a `mode` that changes at runtime keeps every site on shared values.
- Easing: only when the inventory found a timing curve with no exact CSS form. Show the sampled `linear(...)` form for the first such site next to the nearest approximation and its max error, and ask which to use (`references/easing.md`).

The answers hold for the whole run. Every other judgment call is a Needs approval row in the report.

## 2. Decide per site

| Verdict | Meaning |
|---|---|
| Migrate | behavior identical; apply |
| Needs approval | behavior changes in a way the user may accept, or the walk does not cover the site; show the proposed code, apply only after a yes |
| Keep on shared values | CSS cannot express it; leave it, give the reason in one sentence |

Walk the questions in order for every site, and where a site animates several properties judge each property on its own. Keep on shared values ends the walk for that property; Needs approval is noted and the walk continues; a property that reaches the end with nothing noted is Migrate, and with notes it is Needs approval; a recorded easing substitution or a stated color gap is a row annotation, not a note. A pattern the questions do not cover that still looks convertible: Needs approval, with the proposal and a plain statement that the walk does not cover it and you are not sure the behavior is identical.

A hook gets the verdict its properties share. When only some properties of a hook, or only some hooks of one component, can move: Needs approval, proposing the split (the movable properties become CSS on the element, the rest stay in the hook) and saying that the CSS side and the shared value side then run on separate clocks and may drift by a frame against each other. When no property can move, the hook is Keep on shared values with the reasons grouped.

```
1. Where do the target values come from?                    references/drivers.md
   (with* itself produces the frames in between; CSS replaces that part)
   |-- from continuous input: scroll position, a moving finger (gesture onUpdate/
   |   onChange), a sensor, the keyboard, a frame callback computing new targets
   |   every frame -> Keep on shared values
   |-- from a discrete UI-thread event: a worklet gesture onStart/onEnd/onFinalize
   |   without .runOnJS(true) -> note Needs approval: propose moving the write to the
   |   JS thread (.runOnJS(true), or runOnJS, from 4.1.0 scheduleOnRN, inside the
   |   worklet) and holding the target in state; the animation then starts after a JS
   |   round trip instead of on the next UI frame; continue
   |-- through a runOnUI/scheduleOnUI body -> not a source: classify the target the
   |   body writes (a literal, state or prop from the calling handler or effect -> the
   |   JS thread arm, set state where runOnUI was called; another shared value ->
   |   question 10; UI-thread work of its own, scrollTo or a gesture state -> Keep)
   `-- from the JS thread: state, props, handlers, effects, timers, a gesture end
       with .runOnJS(true), also when the write travels through useDerivedValue or
       useAnimatedReaction before it reaches the style -> continue

2. Is it withSpring or withDecay anywhere in the composition?
   |-- YES -> Keep on shared values (no CSS spring or decay)
   `-- NO  -> continue. A withClamp around timing animations never triggers when
             every value stays inside the bounds: both endpoints inside, and an easing
             within 0..1 or one whose overshoot stays inside too (Easing.back undershoots
             the start by a tenth of the distance, Easing.elastic overshoots the end by
             0.066 of it): drop it; a clamp that can trigger -> Keep on shared values
             (CSS cannot clamp). The same for a clamp spelled as Extrapolation.CLAMP,
             clamp() or Math.min/Math.max around the driver, or interpolateColor

3. Does CSS animate the property on every platform the project targets, at the
   installed version? (supported-properties docs, feature table in animations.md)
   |-- NO, the property is a keyword flipped at a state change -> render it
   |   conditionally, leave it out of transitionProperty; continue
   |-- NO, a keyword flipped at > 0.5 of a 0..1 numeric driver -> note Needs approval
   |   proposing transitionBehavior: 'allow-discrete' and saying when CSS flips it
   |   (the midpoint; display around none flips at the start when showing and at the
   |   end when hiding); continue
   |-- NO, a keyword flipped at any other threshold, or a flip that breaks layout
   |   -> Keep on shared values
   |-- NO, for any other reason -> Keep on shared values
   `-- YES -> continue

4. Is each animated value a straight line between its two endpoints
   (a * driver + b, both endpoints the same kind of value)?
   |-- YES, a number, length or percentage -> continue
   |-- YES, a color -> the swing decides (continue; wide swings note Needs approval:
   |   CSS lerps sRGB, the shared value interpolated gamma-corrected)
   |-- NO, two or more drivers feed one property -> entries of one compound property
   |   (a transform array, shadowOffset) written together with one config are one
   |   state object: continue; the same entries with different configs -> note Needs
   |   approval proposing one config and naming the entry whose timing changes, or
   |   Keep on shared values; a.value * b.value or base.value + offset.value -> Keep
   |   on shared values (no single state value to transition)
   |-- NO, a number and a keyword (300 to 'auto') -> Keep on shared values (nothing
   |   tweens, it jumps to the target)
   |-- NO, a percentage string at both ends ('0%' to '100%') -> continue (both tween)
   |-- NO, a number and a percentage -> note Needs approval: CSS resolves the percentage
   |   (against the parent; translateX/translateY, border radii, gaps and transform
   |   origin against the view itself) and tweens; the shared value did not: from a
   |   number start it rendered NaN every frame and jumped to the target at the end
   |   (the migration fixes that), from a percentage start it kept the % and ended at
   |   the wrong value ('50%' to 300 ended at '300%'); say which; continue
   |-- NO, a multi-stop interpolate -> propose a keyframe animation with a keyframe at
   |   every stop, show it: exact when the driver's easing is linear or has an exact row
   |   in references/easing.md (value-functions.md says how it splits per interval),
   |   otherwise note Needs approval with the error; continue
   `-- NO, any other curve -> note Needs approval: propose a keyframe animation that
       samples the value curve every 5 to 10 percent, show it with the step and the
       error; continue
   For a state-driven site that became an animation here, say that a flip mid-flight
   restarts the incoming rule instead of retargeting.

5. Does the timing curve map to CSS exactly?                references/easing.md
   |-- exact row -> continue
   `-- no exact row -> apply the answer from step 1 (sampled linear() or the nearest
       approximation), record the substitution and its max error in the site's row,
       continue

6. Does code cancel, pause, reverse or restart the animation from outside the driver?
   (the driver's own return write, a press out or a toggle flipped back, is question 8;
   decide transition or animation first: references/transitions-and-animations.md)
   |-- pause/resume -> animationPlayState 'paused' / 'running'; note Needs approval:
   |   the shared value re-eased the rest over a full duration, CSS resumes where it
   |   paused; continue
   |-- reverse -> a NEW css.keyframes() rule with animationDirection: 'reverse'
   |   (flipping the direction on the existing rule only mirrors the current position),
   |   or the state flipped back for a transition (question 8); note Needs approval
   |   when the forward run can be interrupted: the reversed run starts from the end
   |   keyframe, not from the interrupted value; continue
   |-- restart -> a new css.keyframes() rule with the same keyframes (a new rule
   |   restarts the animation; the reset write before it is that rule's first
   |   keyframe, not a jump for question 9); continue
   |-- retarget of a transitioned property (cancelAnimation then a with* to another
   |   value) -> the new state; the transition retargets from the current value; continue
   |-- cancel then a with* to the same target on a transition -> note Needs approval:
   |   the shared value restarted at full duration, a transition already heading there
   |   does nothing; continue
   |-- cancel of a running animation (a loop, a sequence) then a with* to a rest value
   |   -> note Needs approval, or Keep on shared values: removing animationName snaps
   |   to the static style, nothing tweens from the current animated value; continue
   |-- cancel and hold the current value -> on an animation: animationPlayState
   |   'paused' holds it, continue; on a transition: nothing holds mid-flight, note
   |   Needs approval stating the snap, or Keep on shared values
   `-- none, or cancelAnimation only in an unmount cleanup -> continue

7. Does a completion callback do something observable?
   |-- YES, 4.6.0+ -> map it to the onCSS* props: onCSSTransitionEnd/Cancel for a
   |   transition, onCSSAnimationEnd/Iteration/Cancel for an animation. With several
   |   animations on one element branch on the event's animationName, comparing
   |   against the .name of the rule object you pass in animationName (a module-scope
   |   css.keyframes() rule, or the restart rule held in state from question 6); a plain
   |   keyframes object has no name you can reference. Note Needs approval when the
   |   moments differ (../animations/animations.md, Callbacks); continue
   |-- YES, below 4.6.0 -> Keep on shared values (no CSS callbacks)
   `-- NO  -> continue

8. Can the driver reverse mid-flight? (press in/out, a toggle the user flips)
   |-- YES, the site is a transition -> note Needs approval: a CSS transition reversed
   |   mid-flight takes a shortened return leg (../animations/animations.md, CSS
   |   Transitions, Rules), the shared value took whatever duration the site gave the
   |   return write; continue
   |-- YES, the site is an animation (a loop, a sequence, or one that became an
   |   animation at question 4 or 6) -> note Needs approval:
   |   CSS animations do not retarget, a flip mid-flight starts the incoming rule from
   |   its first keyframe where the shared value tweened back from the current value;
   |   continue
   `-- NO  -> continue

9. Does other code write the same value without a with*?   references/drivers.md
   (sv.value = x in a handler, a mount jump, a conditional style)
   |-- YES, only the reset that question 6 folded into a restart -> nothing to note; continue
   |-- YES, and nothing writes it with a with* either -> the site never animated:
   |          plain state in the static style, no transitionProperty; continue
   |-- YES -> note Needs approval: every change of a transitioned property animates;
   |          propose rendering that write with the transition turned off for the
   |          property in that render (references/drivers.md, Other writers); continue
   `-- NO  -> continue

10. Does other code read the shared value?                  references/drivers.md
   (another hook, a worklet such as useAnimatedReaction or a gesture, a child that
    receives it as a prop, JS reading sv.value)
   |-- YES -> those readers lose their source when the value goes: note Needs
   |          approval proposing a state mirror for them, or Keep on shared values
   `-- NO  -> Migrate, or Needs approval if anything was noted
```

Outside the walk, leave these alone: `entering`/`exiting`/`layout` animations, `Keyframe`, shared element transitions (a separate subsystem; judge a `useAnimatedStyle` on the same element on its own), a `useAnimatedProps` for a non-style prop on a non-SVG component (`text`, `contentOffset`, `scrollEnabled`), which CSS never covers, and React Native's own `Animated` API (`Animated.Value`, `Animated.timing`, `Animated.event` from `react-native`): count those sites separately in the inventory and name them in the report, so a file that animates only with them is never reported as having nothing to migrate.

## 3. Convert

`references/transitions-and-animations.md` maps `with*` compositions to transitions or animations; `references/easing.md` the timing curve; `references/reduced-motion.md` the guard when the user kept reduced motion. Rules that hold for every site:

- Follow the transition and animation Rules in `../animations/animations.md`; in the code you emit list `transitionProperty` explicitly (never `'all'` or the `transition` shorthand string) and write the timing function (`references/easing.md` says why the default is never left in place).
- Replacing the shared value with React state (`useState`, a prop, a store value) is expected; do not count it as a downside in the report.
- Migrate inside each `Platform.select` arm and keep the structure; enumerate `transitionProperty` per platform when the property set differs. A duration or easing computed per trigger (`duration: Math.abs(delta) * k`) becomes state set in the same render as the target, with the formula unchanged.
- A property the hook returns from a value that is never animated (a constant, or a shared value that is never written after its initial value) is static: move it to the static style. Remove the shared values, hooks and imports the conversion killed and nothing else.
- Never put CSS declarations on a plain component.
- If converting uncovers something the inventory missed (a continuous input after all, a property the docs table does not list, a question that cannot be answered), revert that site and classify it again.

After converting a site, confirm each of these against the original:

- first render identical: the static style carries the value the hook painted first (`../animations/animations.md`, Mount animations);
- end state identical, including where the site rests under reduced motion when the user kept it (`references/reduced-motion.md`);
- re-trigger identical: writing the same target mid-flight looks identical unless the original cancelled first; a replay restarts only through a new keyframes rule (`css.keyframes(frames)` created per replay, `../animations/animations.md`, Defining keyframes), so that form must be present and must not attach the animation on the first render; a remount or `key` change also restarts it but changes the element tree and fails the last check;
- the callbacks the original fired still fire, at the same moments;
- unmount mid-animation throws nothing;
- nothing else changed: same element tree and props, every handler that did more than write the shared value still attached (a `Pressable` swapped for `Animated.View` fails this), and every `Platform.select` arm converted.

Migrate three sites of differing shape first, compare them for consistent treatment, then continue. Keep a status file outside the repo (site, verdict, reason) and trust it plus `git diff` after an interruption.

## 4. Verify

Typecheck and lint the changed files; run their tests (replacing a `useAnimatedStyle` result with plain style props changes snapshots and style assertions over that subtree, read the diff to confirm only the intended properties moved). When a simulator, emulator or device is available, run the app and watch the first frame and one re-trigger of each migrated site, where a missing fill mode or a dead driver shows. Say what you verified and what you did not.

## 5. Report

Lead with the numbers and the grouped reasons, nothing else:

```
Migrated 14 sites across 9 files. 6 need approval. 17 stay on shared values:
  9 track continuous input, 5 use withSpring, 3 have other readers.
```

Then ask what to expand. Offer: the applied sites (one row each: file, what it animates, transition or animation, and every behavior delta or `exact`), the sites that need approval, or the kept sites. Show the applied rows in pages of ten; walk the Needs approval sites one at a time with the open question and the proposed diff, and apply each on a yes. Leave out recipes, memoization advice and per-site prose.
