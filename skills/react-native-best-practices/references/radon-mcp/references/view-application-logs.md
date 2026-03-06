---
name: view-application-logs
description: "Best practices for using the view_application_logs tool in Radon IDE. Returns all build, bundling, and runtime logs from the running app. Use when the user has any issue with the app, when builds are failing, when there are runtime errors or crashes, or whenever logs would help diagnose a problem. Trigger on: 'app logs', 'build logs', 'console errors', 'app crash', 'build failure', 'Metro error', 'runtime error', 'debug logs', 'what went wrong', 'native logs', 'JS logs', or any debugging scenario where log output would be useful."
---

# view_application_logs

Returns all build, bundling, and runtime logs from the app running in Radon IDE. Also includes a screenshot of the current app state when the preview is ready.

## Tool signature

```
view_application_logs()
```

No input parameters.

**Returns:** Multi-part response containing up to 5 log sections (text) and an optional PNG screenshot.

---

## When to use

- **Always** when the user reports any issue with the app — this is the single most useful debugging tool.
- When builds are failing or taking unexpectedly long.
- When the app crashes, freezes, or shows unexpected behavior.
- When Metro bundler reports errors or warnings.
- When native-side errors appear (missing native modules, linking issues, etc.).
- When debugging JS exceptions, unhandled promise rejections, or React errors.

## When NOT to use

- For network-specific issues — use `view_network_logs` instead.
- To only see the visual state without log context — use `view_screenshot` instead.

---

## Understanding the output

The tool collects logs from five separate output channels, each clearly labeled:

### `=== BUILD PROCESS LOGS ===`

Platform-specific build output (Xcode for iOS, Gradle for Android). Look here for:

- Compilation errors and warnings
- Linking failures
- Missing native dependencies
- Code signing issues (iOS)
- Build configuration problems

### `=== JS PACKAGER LOGS ===`

Package manager output (npm, yarn, pnpm, bun). Look here for:

- Dependency installation issues
- Package resolution errors

### `=== METRO LOGS ===`

Metro bundler output. Look here for:

- Bundle compilation errors
- Module resolution failures
- Syntax errors in JS/TS files
- Transform errors
- Hot reload status

### `=== NATIVE-SIDE APP LOGS ===`

Native device logs (Android logcat or iOS system logs). Look here for:

- Native crash reports
- Native module initialization errors
- Platform-specific warnings
- Memory warnings

### `=== JS-SIDE APP LOGS ===`

JavaScript application logs (`console.log`, `console.error`, etc.). Look here for:

- Application-level errors and exceptions
- React component errors
- State management issues
- Custom log output from the app

### Screenshot (when preview is ready)

A PNG screenshot is appended when the app preview is active. This gives visual context alongside the log output, making it easy to correlate errors with the visual state.

---

## Best practices

### Read all sections before diagnosing

Issues often span multiple layers. A native crash might be triggered by a JS error, or a build failure might manifest as a Metro error. Scan all non-empty sections before forming a conclusion.

### Check the section order for build-time vs runtime issues

- **Build-time issues:** Focus on BUILD PROCESS LOGS and JS PACKAGER LOGS.
- **Bundling issues:** Focus on METRO LOGS.
- **Runtime issues:** Focus on NATIVE-SIDE APP LOGS and JS-SIDE APP LOGS.

### Empty sections are normal

Only non-empty log channels are included in the output. If a section is missing, it means that channel has no recorded output — this is expected.

### Use this tool first, then narrow down

When the user reports a problem, call `view_application_logs` as your first step. Based on what you find, you can then decide whether to:

- Call `view_component_tree` for layout issues
- Call `view_network_logs` for API-related errors
- Call `reload_application` to reset the app state
- Suggest code changes based on the errors

---

## Error handling

If Radon IDE is not launched:

```
Couldn't retrieve build logs - Radon IDE is not launched. Open Radon IDE first.
```

If no build has been run:

```
Couldn't retrieve build logs - Radon IDE hasn't run any build.
You need to select a project and a device in Radon IDE panel.
```

**Resolution:** Ask the user to open the Radon IDE panel, select a device, and build/run the project.
