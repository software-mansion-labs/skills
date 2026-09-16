---
name: migrate-to-css-animations
description: "Migrates React Native Reanimated hook animations (useAnimatedStyle, useSharedValue, withTiming, withRepeat, withSequence, withDelay) to CSS transitions and animations, applying only conversions that keep behavior identical and explaining why the rest stay on hooks. Use when asked to migrate, convert, port or audit Reanimated animation code for CSS, or to simplify, modernize or drop worklets from animations. Requires Reanimated 4.x; reads the installed version and emits only what that version supports."
---

# Migrate Reanimated hook animations to CSS

Convert only where behavior stays identical; coverage is not the goal. A wrong conversion is silent: it typechecks, passes tests and looks right in a diff.

"Hook animations" here means what `../animations/animations.md` calls shared value animations. For the CSS API itself read `../animations/animations.md` (feature availability by version, transitions, animations, timing functions, callbacks) and, from 4.5.0, `../animations/css-pseudo-selectors.md`. Load `references/mapping.md` when converting a site and `references/examples.md` to calibrate output.

## 0. Detect version and platforms

1. Installed Reanimated: `node_modules/react-native-reanimated/package.json` or the lockfile, never the `package.json` range. 3.x has no CSS API: stop and say so, do not upgrade. On 4.x note the exact version; the feature table in `../animations/animations.md` decides what may be emitted, and a site needing a feature under its floor is Keep.
2. Target platforms (`package.json`, app config, web bundler). A `.ios`/`.android`/`.web` suffix, a `Platform.OS` branch or a platform constant narrows a site to those platforms; everything else ships everywhere the project does.
3. Working tree clean, or the user acknowledges the changes are theirs.

Scale to the request: one site, run steps 0 and 2 and convert or explain; "can this be CSS?" or an audit, run steps 0, 1, 2 and 4 and edit nothing; a directory or app, all steps.

## 1. Inventory, then ask once

Record every `useAnimatedStyle`, `useAnimatedProps`, `useDerivedValue` and `with*` site: file, component, animated properties, driver, effective platforms. Report the counts. Nothing animated in scope: say so and stop.

Then ask once, in one message, before applying anything:

- Scope: the files you will touch, nothing outside.
- Colors: wide swings get a different midpoint (`references/mapping.md`, Colors). Migrate all colors, only visually close pairs, or keep colors on hooks?
- Partial migration: when some sites clear and others do not, migrate everything that clears, a named subset, or nothing unless all of it clears?

Answers hold for the whole run.

## 2. Classify: three verdicts per site

| Verdict | When | Action |
|---|---|---|
| Migrate | every check passes and no Keep pattern matches | apply |
| Needs approval | passes but changes a behavior the user may accept, or a judgment call | show the proposed code, apply only after the user says yes |
| Keep on hooks | a check fails or a Keep pattern matches | leave it, give the reason in one sentence |

Each check returns one of the three; the strictest wins. Name the check that permits each migration; if you cannot, do not migrate.

### Checks

1. **Driver written on the JS thread, changing discretely.** `useEffect`, a JS callback (`onPress`, `onChange`, timer, network), or a gesture `onEnd`/`onFinalize` with `.runOnJS(true)` writing one target: pass. Any `onUpdate`/`onChange` write, a gesture without `.runOnJS(true)`, a scroll handler, `useAnimatedReaction`, `useFrameCallback`, a sensor or keyboard hook, a `runOnUI`/`scheduleOnUI` body, or a JS loop assigning every frame: Keep, the value changes per frame or on the UI thread. A driver whose origin you cannot trace in scope: Needs approval. Turning the shared value into `useState` is the migration, not a cost.
2. **Values are pure functions of the driver.** Affine arithmetic and `interpolate` with fixed stops become keyframes. Trigonometry, `Math.pow`, modulo, parametric geometry: Keep. `interpolateColor` in `'HSV'` or `'LAB'`: Keep. `interpolate` without `Extrapolation.CLAMP` whose driver can leave the input range: Keep (CSS clamps at the outer keyframes). A ternary on the driver whose result is not wrapped in a `with*` is a step, not a tween: render it conditionally and leave it out of `transitionProperty`.
3. **Indirection resolves to one migratable driver.** `useDerivedValue` or a `SharedValue` prop read only by this style: inline it and classify the root writer. Several style consumers of one loop: Needs approval (`references/mapping.md`, Transition or animation). Any worklet reader or writer: Keep.
4. **Every property animates on the effective platforms** at the installed version. Check the [supported properties](https://docs.swmansion.com/react-native-reanimated/docs/guides/supported-properties) table, never memory. React Native does not render it on a platform either: fine, it was inert (iOS `shadow*` on Android, `elevation` on iOS). React Native renders it and CSS cannot animate it: Keep. Discrete (keyword) properties: Needs approval, proposing `transitionBehavior: 'allow-discrete'` and stating when the value flips (`../animations/animations.md`, Discrete properties); Keep when the flip breaks layout.
5. **No imperative control.** `cancelAnimation`, pausing, reversing or restarting from an effect or handler: Keep. Two uses are not control: `cancelAnimation` in an unmount cleanup, and `cancelAnimation` followed by a `with*` to a different value (a retarget). After `cancelAnimation`, a `with*` to the same target restarts the hook at full duration where CSS ignores the write: Needs approval.
6. **Completion callback reproducible.** Log-only: drop it. Observable (chains an animation, sets state): from 4.6.0 Needs approval, converted per the callback row in `references/mapping.md`; below 4.6.0 Keep. A hook callback also fires `finished: true` when the target already equals the value; CSS fires nothing then. Say so in the row.
7. **The driver cannot reverse the transition mid-flight** (the shortened return leg, `../animations/animations.md`, CSS Transitions, Rules). `onPressIn`/`onPressOut` pairs and controls the user toggles directly: Needs approval, offering the version-appropriate shape (from 4.5.0 `:active` on the pressed element, 4.6.0 for SVG elements; below, a `Pressable` render prop styling the child, or `useState` from `onPressIn`/`onPressOut` when the `Pressable` itself carries the style). `:active` shortens the return leg too. A driver behind a timer longer than the duration, or a rare event (rotation, navigation, network): Migrate and state the difference. Animations have no reversal behavior.
8. **No other writer sets a `transitionProperty` target instantly** (a write with no `with*`: a mount jump, a conditional attach). The jump becomes a glide: Needs approval.

### Keep on hooks, whatever the checks say

| Pattern | Why |
|---|---|
| `withSpring`, `withDecay`, `withClamp` | no CSS spring, no velocity-driven or clamped timeline; never hand-sample a spring into `linear()` |
| `Easing.elastic`, `Easing.bounce`, `Easing.exp`, `poly(n)` for n not 1..3, bare `sin`/`circle` | no exact cubic-bezier form, and this skill approximates only the `inOut` curves listed in `references/mapping.md` |
| a measured layout value (`measure`, `useAnimatedRef`, `onLayout` size) feeding an animated value | targets unknown at write time |
| `useAnimatedProps` for a non-style prop on a non-SVG component (`text`, `contentOffset`, `scrollEnabled`) | CSS covers view, text and image styles, plus SVG attributes from 4.4.0 (`references/mapping.md`, SVG) |
| SVG below 4.4.0 (4.1.0+ only with `EXPERIMENTAL_CSS_ANIMATIONS_FOR_SVG_COMPONENTS`; web below 4.5.0); SVG `transform` arrays; SVG `mask`, `filter`, `marker*`, `pointerEvents` (listed as supported in the docs table but throw `No interpolator factory found`); `fill`/`stroke` as `currentColor` or `url(#...)` | `references/mapping.md`, SVG |
| `entering`/`exiting`/`layout` animations, `Keyframe`, shared element transitions | separate subsystem; leave them and judge a `useAnimatedStyle` on the same element alone |
| `<ReducedMotionConfig>` whose `mode` is not a literal, or is toggled at runtime | `useReducedMotion()` reads the system setting once; `with*` reads the config live. A literal `Always`/`Never` applies the reduced-motion rows in `references/mapping.md` |

Keep is not user-overridable: a lossy substitute is a design change, done in its own labeled diff.

## 3. Convert

`references/mapping.md` has the `with*`, easing, reduced-motion, color and SVG tables. Rules that hold for every site:

- One `useAnimatedStyle` or `useAnimatedProps` is all or nothing: a property that fails a check keeps the whole hook. Separate animated styles on one component may get different verdicts; that is a partial migration and follows the user's answer from step 1.
- Follow the transition and animation Rules in `../animations/animations.md`; above all list `transitionProperty` explicitly (never `'all'` or the `transition` shorthand string) and write the timing function (the `withTiming` default is `Easing.inOut(Easing.quad)` at 300ms, not the CSS default `'ease'`).
- Migrate inside each `Platform.select` arm and keep the structure; enumerate `transitionProperty` per platform when the property set differs.
- Move never-animated values from the deleted hook into the static style. Remove the shared values, hooks and imports the conversion killed and nothing else.
- Never put CSS declarations on a plain component.
- Mid-edit surprise (worklet writer after all, unlisted property, unsatisfiable check): revert that site and re-classify it, Keep when a check fails, Needs approval when it is a judgment call.

After converting, confirm all five:

- first render identical (`../animations/animations.md`, Mount animations);
- end state identical;
- re-trigger identical: the same target mid-flight is identical unless preceded by `cancelAnimation` (check 5); replaying the same animation needs the animation detached for one render or a `key` change, a visible structural change, so Needs approval;
- unmount mid-animation throws nothing;
- nothing else changed: same element tree, props and handlers; a `Pressable` swapped for `Animated.View` fails this.

Migrate three sites of differing shape first, compare them for consistent treatment, then continue. Keep a status file outside the repo (site, verdict, reason) and trust it plus `git diff` after an interruption.

## 4. Report

Lead with the numbers: `Migrated 14 sites across 9 files. Left 23: 6 need approval, 17 stay on hooks.` Then:

1. **Applied**: one row per site (file, what it animates, transition or animation, note). The note lists every behavior delta or says `exact`: easing substituted and its max error, reduced-motion guard dropped, a reversal now shorter, a color midpoint moved, a mount step or replay that no longer fires.
2. **Needs approval**: one row per site with the open question and the proposed code. Ask which to apply, by number.
3. **Kept on hooks**, grouped by reason with counts.
4. Before/after code for two or three applied sites of different shape.

Leave out recipes, memoization advice and per-site prose.

## 5. Verify

Typecheck and lint the changed files; run their tests (replacing a `useAnimatedStyle` result with plain style props changes snapshots and style assertions over that subtree, read the diff to confirm only the intended properties moved); run the app and watch the first frame and one re-trigger, where a missing fill mode or a dead driver shows. Say what you verified and what you did not.
