---
name: reload-application
description: "Best practices for using the reload_application tool in Radon IDE. Triggers a reload of the app running in the development emulator with three methods: reloadJs (JS bundle reload), restartProcess (native process restart), and rebuild (full native rebuild). Use when debugging state issues, after code changes, when the app crashes, or when native changes require a rebuild. Trigger on: 'reload app', 'restart app', 'rebuild', 'hot reload', 'reset state', 'app crashed', 'app frozen', 'refresh', 'reloadJs', 'restartProcess', or any request to restart or reload the running application."
---

# reload_application

Triggers a reload of the app running in the Radon IDE development emulator. Supports three reload methods with increasing scope and cost.

## Tool signature

```
reload_application({ reloadMethod: "reloadJs" | "restartProcess" | "rebuild" })
```

**Input:**

- `reloadMethod` (required): One of three reload strategies.

**Returns:** `"App reloaded successfully."` or an error message.

---

## Reload methods

### `reloadJs` — JS bundle reload (fastest)

Reloads the JavaScript bundle from Metro without restarting the native app process.

**Use when:**

- You changed JS/TS/JSX/TSX files and want to see the result.
- You want to reset the JS state of the app (clear React state, context, stores).
- Hot Module Replacement (HMR) did not pick up a change.
- The app is in a broken JS state but the native side is fine.

**Does NOT help when:**

- Native modules need to be re-initialized.
- You changed native code (Objective-C, Swift, Java, Kotlin).
- You added or removed a native dependency.

---

### `restartProcess` — native process restart (moderate)

Kills and relaunches the native app process on the device without rebuilding.

**Use when:**

- A native library or component is in a bugged state.
- The app has a stale native module cache.
- You need a full cold start of the app (as if the user force-quit and reopened it).
- `reloadJs` did not fix the issue and you suspect a native-side problem.

**Does NOT help when:**

- You changed native code that requires recompilation.
- You added a new native dependency that needs linking.

---

### `rebuild` — full rebuild (slowest)

Rebuilds both the JS bundle and the native parts of the app, then reinstalls and launches it on the device.

**Use when:**

- You made changes to native code (iOS or Android).
- You added, removed, or updated a native dependency.
- You changed build configuration (Xcode scheme, Gradle build type, etc.).
- You ran `expo prebuild` or changed the Expo config.
- Both `reloadJs` and `restartProcess` failed to resolve the issue.

**Be aware:**

- This is significantly slower than the other methods — it triggers a full native build.
- Use this as a last resort for JS-only changes.

---

## Best practices

### Choose the lightest method that solves the problem

Always start with `reloadJs`. If the issue persists, escalate to `restartProcess`. Only use `rebuild` when native changes are involved or the other methods fail. This saves significant time during development.

```
reloadJs (seconds) → restartProcess (seconds) → rebuild (minutes)
```

### Use after code changes

After making code changes that fix a bug, call `reload_application` with the appropriate method to verify the fix. Then use `view_screenshot` or `view_application_logs` to confirm the app is in a healthy state.

### Use to reset state when debugging

When debugging state-dependent issues, `reloadJs` is the quickest way to get a clean slate. This resets all React state, context providers, and in-memory stores.

### The tool waits for the app to be ready

`reload_application` does not return until the app status transitions back to "running". This means you can safely call `view_screenshot` or other inspection tools immediately after — the app will be ready.

---

## Error handling

If the development device is off:

```
Could not reload the app!
The development device is likely turned off.
Please turn on the Radon IDE emulator before proceeding.
```

If the reload fails:

```
Failed to reload the app. Details: <error details>
```

**Resolution:** For device-off errors, ask the user to start a device in Radon IDE. For reload failures, check `view_application_logs` for build or runtime errors.
