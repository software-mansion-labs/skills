---
name: migrate-to-css-animations
description: "Migrates React Native Reanimated shared value animations (useAnimatedStyle, useSharedValue, withTiming, withRepeat, withSequence, withDelay) to CSS transitions and animations, applying only conversions that keep behavior identical and explaining why the rest stay on shared values. Use when asked to migrate, convert, port or audit Reanimated animation code for CSS, or to simplify, modernize or drop worklets from animations. Requires Reanimated 4.x; reads the installed version and emits only what that version supports."
---

# Migrate shared value animations to CSS

Convert only where behavior stays identical; coverage is not the goal. A conversion that typechecks and runs without errors can still animate differently, so neither proves it is correct: the walk below and the checks after converting do.

For the CSS API itself read `../animations/animations.md` (feature availability by version, transitions, animations, keyframes, timing functions, callbacks) and, from 4.5.0, `../animations/css-pseudo-selectors.md`. Each question of the walk names the reference file that holds its rule, mapping and example; load only the files the site needs.

## 0. Version and scope

- Installed Reanimated: `node_modules/react-native-reanimated/package.json` or the lockfile, never the `package.json` range. 3.x has no CSS API: stop and say so, do not upgrade. On 4.x note the exact version; the feature table in `../animations/animations.md` decides what may be emitted, and a site that needs a feature the installed version lacks stays on shared values.
- Scope: the files or directories the user named. With none named, the project's own source (not `node_modules`); say so when reporting the inventory.

The shape of the request decides how much of this skill runs:

- One site named by the user ("migrate the fade in `Header.tsx`"): check the version, walk that site, convert it or explain in one paragraph why it stays. Ask the step 1 questions only when that site raises them.
- A question or an audit ("can this be CSS?", "what could move?"): inventory, walk every site, report; change nothing.
- A migration of a directory or the whole project ("migrate `src/screens` to CSS"): every step below, inventory through report.

## 1. Inventory, then ask once

Record every animation site in scope: `useAnimatedStyle`, `useAnimatedProps`, `useDerivedValue`, a shared value passed directly in `style` or as a prop (`style={{ opacity: sv }}`, `<AnimatedCircle r={r} />`), and every `with*` call. For each: file, component, animated properties, where the values come from. Report the counts. Nothing animated in scope: say so and stop.

Then ask once, in one message, only the questions the inventory raised:

- Scope: confirm the files you will touch.
- Partial components: when a component has both sites that can move and sites that must stay, migrate the ones that can (that component then mixes CSS and shared values) or leave the whole component on shared values? Judged per component, not per file; a file with several independent components can end up with some migrated and some not either way.
- Reduced motion: every `with*` animation follows the system Reduce Motion setting by default (`ReduceMotion.System`: when the setting is enabled the animation completes at once at its resting value, the target for `withTiming`, `withSequence` and a non-reversed `withRepeat`, the start for a reversed `withRepeat` with an even or infinite count); CSS transitions and animations ignore the setting. Ask: keep that behavior (every migrated site gets a `useReducedMotion()` guard, `references/reduced-motion.md`; the hook reads the system setting once at app start and ignores `<ReducedMotionConfig>`) or drop it (users with Reduce Motion enabled will see the migrated animations play)? Do not ask when the source already settles it for every site: `ReduceMotion.Never` on each `with*`, or a `<ReducedMotionConfig>` with a fixed `mode={ReduceMotion.Never}`, means no guard; `ReduceMotion.Always` on each `with*`, or a fixed `mode={ReduceMotion.Always}`, means the animation never played for anyone, so emit the reduced form (duration `1`, `references/reduced-motion.md`) with no guard; a `mode` that changes at runtime keeps every site on shared values; `mode={ReduceMotion.System}` is the same as no config, ask.
- Easing: only when a timing curve has no exact CSS form. Show the sampled `linear(...)` form for the first such site next to the nearest approximation and its max error, and ask which to use (`references/easing.md`).
- Colors: only when a color animates. CSS interpolates colors straight in sRGB; `withTiming` on a color and `interpolateColor` interpolate gamma-corrected (gamma 2.2), so the middle of the run differs while both ends match: black to white passes through `#808080` in CSS and `#bababa` on the shared value, up to 72/255 apart a quarter of the way in; `#eee` to `#ccc` differs by 1/255. Show the first color site's two midpoint colors and ask whether the difference is acceptable. A no keeps every color site on shared values.

The answers hold for the whole run. Everything else that needs a decision becomes a Needs approval row in the report.

## 2. Decide per site

| Verdict | Meaning |
|---|---|
| Migrate | behavior identical; apply |
| Needs approval | behavior changes in a way the user may accept, or the walk does not cover the site; show the proposed code, apply only after a yes. A declined row stays on shared values |
| Keep on shared values | CSS cannot express it; leave it, give the reason in one sentence |

Walk the questions in order for every site, and where a site animates several properties judge each property on its own. Keep on shared values ends the walk for that property, and so does the state-switched keyword arm of question 3, as Migrate; a Needs approval arm adds a note and the walk continues; a property that reaches the end with no note is Migrate, with a note it is Needs approval. The easing and color forms chosen in step 1 are already approved: write them in the site's row, they add no note. A pattern the questions do not cover that still looks convertible: Needs approval, with the proposal and a plain statement that the walk does not cover it and you are not sure the behavior is identical.

A hook's verdict follows its properties: Keep on shared values when every property is Keep (reasons grouped); Migrate when every property reached the end with no note; Needs approval otherwise. When some properties are Keep and others can move, the Needs approval row proposes the split (the movable properties become CSS on the element, the rest stay in the hook) and says that the CSS side and the shared value side then run on separate clocks and may drift by a frame against each other; the same holds for two hooks on one component. Entries of one compound property (the transform array, shadowOffset) are one property with one verdict: when any entry stays on shared values the whole property stays, naming that entry as the reason, because the element carries the entries as one style key (a CSS transform and a hook transform on the same element replace each other whole) and a wrapper view would change the element tree (step 3).

Pick the form before walking, per property; question 4 overrides it when the value is not a straight line between its endpoints. A **transition**: a trigger outside the animation (React state, a prop, a press, a toggle) moves the property from one resting value to another, and that trigger becomes React state; an effect or handler that writes a `with*` whenever the trigger changes is still a transition. An **animation**: after a start signal the timeline runs by itself (a loop, a `withSequence`, a one-shot play on mount with no later trigger, or a curve with stops between its endpoints, question 4). `references/transitions-and-animations.md` has the edge cases; questions 6 and 8 depend on the form.

The walk calls the driver whatever moves the property: before migration the shared value a `with*` is assigned to (`progress`, `pressed`, `offset`) together with the handler or effect that assigns it, after migration the state or prop that replaces them; the driver's easing is that `with*` call's easing. Its own writes are the plain `with*` target writes; a `cancelAnimation`, pause or restart around them, in the same handler or elsewhere, is code outside the driver (question 6).

```
1. Where do the target values come from?                    references/drivers.md
   (with* itself produces the frames between two targets; CSS replaces that part)
   |-- continuous input: a scroll position, a moving finger (a gesture's onUpdate, or
   |   onChange in Gesture Handler 2, writing the value), a sensor, the keyboard, a
   |   frame callback computing new targets every frame -> Keep on shared values
   |-- a worklet gesture callback that fires once per interaction (onBegin, onStart or
   |   onActivate, onEnd or onDeactivate, onFinalize), on a gesture without runOnJS: true,
   |   writing a target for a property that onUpdate does not also write: the callback
   |   runs on the UI thread, and CSS needs the target in React state -> note Needs
   |   approval: propose scheduleOnRN(setTarget, value) inside that callback only
   |   (imported from react-native-worklets, which react-native-reanimated does not
   |   re-export; on 4.0.x, which has no scheduleOnRN, runOnJS(setTarget)(value), a
   |   deprecated but still exported function), never runOnJS: true on the gesture, which
   |   would move every callback, onUpdate included, to the JS thread; the animation then
   |   starts after a JS round trip and a render instead of on the next UI frame. A Pan
   |   or Pinch activate/deactivate pair is this arm. Exception: a press pair
   |   (onPressIn/onPressOut, or a Tap or LongPress gesture's onBegin/onFinalize; a Tap's
   |   onStart and onEnd both fire at release, so they are not one) writing the pressed
   |   and the rest value of a property on the component that owns the style adds no
   |   note here: question 8 names its CSS form; continue
   |-- through a runOnUI/scheduleOnUI body -> not a source: classify the target the
   |   body writes (a literal, state or prop from the calling handler or effect -> the
   |   JS thread arm, set state where runOnUI was called; another shared value -> a
   |   chain, follow it back to its first write (references/drivers.md, Chains);
   |   UI-thread work of its own, scrollTo or a gesture state -> Keep)
   `-- the JS thread: state, props, handlers, effects, timers, any callback of a gesture
       with runOnJS: true or whose callbacks are not worklets (references/drivers.md
       says when that happens), also when the write travels through useDerivedValue or
       useAnimatedReaction before it reaches the style -> continue

2. Is there a spring, a decay, or a clamp?
   |-- withSpring or withDecay anywhere in the composition -> Keep on shared values
   |   (CSS has no spring and no decay)
   |-- a clamp: withClamp, clamp(), Extrapolation.CLAMP passed to interpolate,
   |   Math.min/Math.max around the driver, interpolateColor (it always clamps its
   |   input)
   |   |-- it can never trigger: both endpoints inside the bounds and every frame too,
   |   |   which holds for an easing that stays within 0..1 and for an overshooting
   |   |   easing whose excursion stays inside (Easing.back() with its default argument
   |   |   undershoots the start by 0.100 of the distance, 4s^3/(27(s+1)^2) for another
   |   |   s; Easing.elastic() with its default overshoots the end by 0.066 of it, more
   |   |   for a larger bounciness; a cubic bezier stays inside the range spanned by 0, 1
   |   |   and its y control points and reaches less than they do, so compute its
   |   |   extreme when a control point lies outside 0..1: the back() row's
   |   |   cubicBezier(1/3, 0, 2/3, -0.567) bottoms out at -0.100) -> drop the clamp;
   |   |   continue
   |   `-- it can trigger -> Keep on shared values (CSS cannot clamp)
   `-- neither -> continue

3. Does CSS animate the property on every platform the project targets, at the
   installed version? (supported-properties docs, feature table in animations.md)
   |-- YES -> continue
   |-- NO, but the property is a keyword (display, position, flexDirection, ...)
   |   switched from React state (display: open ? 'flex' : 'none') -> it never
   |   animated, it switched: set it from state in the static style and do not list it
   |   in transitionProperty; that property is Migrate and its walk ends here
   |-- NO, a keyword switched at a threshold of a numeric driver that animates
   |   (display: progress.value > 0.5 ? 'flex' : 'none') -> CSS switches a keyword
   |   during a transition only with transitionBehavior: 'allow-discrete', and then at
   |   the midpoint of the transition
   |   |-- display -> note Needs approval: with allow-discrete, display leaves 'none' at
   |   |   the start of the transition and enters 'none' at its end, not at the source
   |   |   threshold, so the element stays visible while it fades; say so; continue
   |   |-- another keyword, threshold 0.5 -> continue (the midpoint is the threshold)
   |   `-- another keyword, any other threshold -> note Needs approval: CSS switches at
   |       the midpoint, not at the source threshold; continue
   `-- NO, for any other reason (no interpolator for the property on that platform or
       version) -> Keep on shared values

4. Is each animated value a straight line between its two endpoints?
   (a * driver + b, both endpoints the same kind of value)
   |-- YES, a number at both ends, or a percentage string at both ends -> continue
   |-- YES, a color -> continue (the step 1 answer covers the interpolation difference;
   |   write both midpoint colors in the row)
   |-- NO, a keyword that question 3 routed to allow-discrete -> continue
   |-- NO, two drivers feed one property
   |   |-- entries of one compound property (translateX and translateY of a transform
   |   |   array, width and height of shadowOffset) written together, in one handler
   |   |   with one config -> one state object holds both entries; continue
   |   |-- the same entries written with different durations or easings -> note Needs
   |   |   approval proposing one shared config and naming the entry whose timing
   |   |   changes; continue
   |   `-- combined into one value (a.value * b.value, base.value + offset.value) ->
   |       Keep on shared values (no single state value to transition)
   |-- NO, a number at one end and a keyword at the other (300 to 'auto') -> Keep on
   |   shared values, saying that the site never animated: the shared value could not
   |   interpolate them either (NaN frames, then the keyword when the duration ended),
   |   and CSS would switch to the target instead of interpolating
   |-- NO, a number at one end and a percentage at the other (50 to '100%') -> the
   |   shared value never interpolated these (a number start produced NaN frames, a
   |   percentage start kept the % and ended at '300%' for a target of 300), so the
   |   site is broken as written: note Needs approval saying that CSS will interpolate
   |   it, resolving the percentage against the parent (against the view itself for
   |   translate, border radii, gaps and transform origin); continue
   |-- NO, an interpolate with stops between the endpoints
   |   (interpolate(p, [0, 0.5, 1], [0, 200, 100])) -> not a transition (a transition
   |   goes straight between its endpoints) but a keyframe animation with one keyframe
   |   per stop, holding that stop's output value, at the time t where the driver has
   |   covered the stop's share of its run: solve easing(t) = (stop - start) / (target -
   |   start) with the driver's endpoints of that rule (for Easing.linear, t is the
   |   share itself; the return rule of a two-way driver reaches stop 0.25 at share
   |   0.75). Each interval gets, as the animationTimingFunction of the keyframe that
   |   starts it, the piece of the driver's easing between its two keyframe times,
   |   rescaled to 0..1 in time and value, because CSS eases every interval on its own
   |   clock. Exact for linear (every piece is 'linear'), for Easing.inOut(f) with a
   |   stop at share 0.5 (the pieces are f then Easing.out(f), exact when f has an exact
   |   row in references/easing.md), and for an exact bezier row (subdivide the curve at
   |   the keyframe times and map the control points onto 0..1; exact while the mapped
   |   x controls stay within 0..1, which holds for every named row). Any other curve,
   |   or a piece that fails that check, is sampled with linear() at even stops of the
   |   piece, with the stop count the step 1 answer settled. For a two-way driver emit
   |   one rule per direction chosen by state, attached only after the first change,
   |   with animationFillMode: 'forwards' and the static style at the resting value.
   |   Show it; continue (this settles question 5 for the property)
   `-- NO, any other function of the driver (Math.sin(driver), a lookup table) -> note
       Needs approval: propose a keyframe animation that samples the value curve every
       5 to 10 percent, show it with the step and the error; continue
   A property that became an animation here while its driver is state (a toggle) no
   longer retargets: question 8 notes what a flip mid-flight does.

5. Does the timing curve map to CSS exactly?                references/easing.md
   |-- exact row, or a multi-stop property question 4 already split per interval
   |   -> continue
   `-- no exact row -> apply the step 1 answer (sampled linear() or the nearest
       approximation), write the substitution and its max error in the row; continue

6. Does code outside the driver cancel, pause, reverse or restart the animation?
   (the driver's own return write on a transition, a press out or a toggle flipped
   back, is question 8; stopping a loop when state changes is this question)
   |-- pause and resume -> animationPlayState 'paused' / 'running'; note Needs
   |   approval: the shared value emulated pause with cancelAnimation and re-eased the
   |   remaining distance over a full duration on resume, CSS resumes where it paused
   |   with the time that was left; continue
   |-- reverse (cancelAnimation, then a with* back to the start) -> on an animation: a
   |   new keyframes rule object (css.keyframes() called again; changing the direction
   |   on the rule already attached only mirrors the current position:
   |   ../animations/animations.md, Defining keyframes) with animationDirection:
   |   'reverse'; the reversed run starts from the end keyframe, not from the
   |   interrupted value, so note Needs approval when the forward run can be
   |   interrupted; on a transition: the state flipped back, question 8; continue
   |-- restart (sv.value = 0, then the same with* again) -> a new keyframes rule
   |   object with the same keyframes restarts the animation; create it where the
   |   restart happened and attach the animation only once it exists, so nothing plays
   |   on the first render. The reset write is that rule's first keyframe, not a jump
   |   for question 9; continue
   |-- cancel, then a with* to a different value, on a transition
   |   (cancelAnimation(sv); sv.value = withTiming(other)) -> set the new state; a CSS
   |   transition retargets, that is, it continues from wherever the value is toward
   |   the new target, so cancelAnimation has no counterpart; continue
   |-- cancel, then a with* to the same target, on a transition -> note Needs
   |   approval: the shared value restarted at full duration, a transition already
   |   heading there does nothing; continue
   |-- cancel of a running animation (a loop, a sequence), then a with* to a rest value
   |   (cancelAnimation(pulse); pulse.value = withTiming(1) when loading ends) -> note
   |   Needs approval: removing animationName snaps to the static style, nothing
   |   animates from the current animated value; continue
   |-- cancel that holds the current value (a stop button) -> on an animation:
   |   animationPlayState 'paused' holds it, continue; on a transition: nothing holds
   |   mid-flight, note Needs approval stating the snap to the static style; continue
   `-- none, or cancelAnimation only in an unmount cleanup -> continue

7. Does a completion callback do something observable? (a log-only callback is
   dropped without a row)
   |-- YES, below 4.6.0 -> Keep on shared values (no CSS callbacks)
   `-- YES, 4.6.0+ -> map it to the onCSS* props: onCSSTransitionEnd/Cancel for a
       transition, onCSSAnimationEnd/Iteration/Cancel for an animation. With several
       animations on one element branch on the event's animationName, comparing
       against the .name of the rule object you pass in animationName (a module-scope
       css.keyframes() rule, or the restart rule held in state from question 6); a
       plain keyframes object has no name you can reference. Note Needs approval when
       the moments differ (../animations/animations.md, Callbacks). Then check the
       moments CSS never reports and note Needs approval for each, proposing the
       replacement: a callback on an earlier child of a withSequence (CSS fires End
       once, after the last child) -> a timer at that child's offset started from
       onCSSAnimationStart; a callback on a write whose target already equals the
       current value, including the mount run of an effect whose target equals the
       initial value (the shared value fired it at once with finished true, a
       transition to the value already shown fires nothing) -> call the handler
       directly in the branch that writes the unchanged target, or in a mount effect
       when the initial state already selects that value; continue

8. Can the user reverse the driver within the duration? (press and release, open and
   close, a toggle or a theme switch flipped back)
   |-- NO: the driver is one-shot relative to the duration (a navigation push, one
   |   network result, a rotation) -> continue
   |-- YES, a transition -> note Needs approval: a CSS transition reversed mid-flight
   |   takes a shortened return leg (../animations/animations.md, CSS Transitions,
   |   Rules), the shared value took whatever duration the site gave the return write.
   |   A press pair on the component that owns the style: from 4.5.0 the pressed value
   |   goes under ':active' (../animations/css-pseudo-selectors.md), nothing re-renders,
   |   and the note adds that ':active' also releases once the finger moves about 10
   |   points, where a Pressable stayed pressed until the finger left its press rect;
   |   below 4.5.0 React state set from the press handlers (through scheduleOnRN when
   |   they are worklets); continue
   `-- YES, an animation (a loop, a sequence, or a property that became an animation at
       question 4 or 6) -> note Needs approval: CSS animations do not retarget, a flip
       mid-flight starts the incoming rule from its first keyframe where the shared
       value animated back from the current value (a loop re-entered mid-flight ran
       every cycle from the re-entry value: references/transitions-and-animations.md,
       the withRepeat rows); continue

9. Does other code write the same value without a with*?   references/drivers.md
   (sv.value = x in a handler, a mount jump, or a branch in the hook body or a
    style-array entry that gives the property a different value from a condition that
    is not the animated driver, opacity: disabled ? 0.5 : sv.value; a branch on the
    driver itself is question 3 for a keyword and question 4 otherwise, not this)
   |-- YES, only the reset that question 6 folded into a restart -> nothing to note;
   |   continue
   |-- YES, and nothing writes it with a with* either -> the site never animated: plain
   |   state in the static style, no transitionProperty; continue
   |-- YES -> note Needs approval: every change of a transitioned property animates;
   |   propose rendering that write with the transition turned off for the property in
   |   that render (references/drivers.md, Other writers); continue
   `-- NO  -> continue

10. Does other code read the shared value?                  references/drivers.md
   (another hook, a worklet such as useAnimatedReaction or a gesture, a child that
    receives it as a prop, JS reading sv.value)
   |-- YES, another hook in scope reads it -> walk that hook as its own site; question
   |   1's Chains rule collapses both to the same state; continue to the verdict
   |-- YES, a reader that uses the value only at rest (JS reading sv.value in a handler
   |   that cannot run while it animates) -> note Needs approval proposing a state
   |   mirror for it; continue to the verdict
   |-- YES, a reader that uses the value while it animates: a worklet (a
   |   useAnimatedReaction that Chains did not remove, a gesture callback, scrollTo), a
   |   child or library component outside the scope rendering it from a prop, or a
   |   handler that can run mid-flight -> Keep on shared values, naming the reader
   `-- NO  -> Migrate, or Needs approval if anything was noted
```

Outside the walk, leave these alone: `entering`/`exiting`/`layout` animations, `Keyframe`, shared element transitions (a separate subsystem; judge a `useAnimatedStyle` on the same element on its own), a `useAnimatedProps` for a non-style prop on a non-SVG component (`text`, `contentOffset`, `scrollEnabled`), which CSS never covers, and React Native's own `Animated` API (`Animated.Value`, `Animated.timing`, `Animated.event` from `react-native`): count those sites separately in the inventory and name them in the report, so a file that animates only with them is never reported as having nothing to migrate.

## 3. Convert

`references/transitions-and-animations.md` maps `with*` compositions to transitions or animations; `references/easing.md` the timing curve; `references/reduced-motion.md` the guard when the user kept reduced motion. Rules that hold for every site:

- Follow the transition and animation Rules in `../animations/animations.md`; in the code you emit list `transitionProperty` explicitly (never `'all'` or the `transition` shorthand string) and write the timing function (`references/easing.md` says why the default is never left in place).
- Replacing the shared value with React state (`useState`, a prop, a store value) is expected for the value that starts a transition or selects an animation, set once per change; do not count it as a downside in the report. A value that changed every frame stayed on shared values at question 1.
- Migrate inside each `Platform.select` arm and keep the structure; enumerate `transitionProperty` per platform when the property set differs.
- A duration or easing computed when the animation is triggered (`withTiming(target, { duration: Math.abs(target - current) * 2 })`) is computed in the same handler and set in state together with the target, then passed as `transitionDuration`; the formula does not change. When it depends only on props or state, compute it in render instead.
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
