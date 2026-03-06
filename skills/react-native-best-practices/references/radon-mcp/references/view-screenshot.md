---
name: view-screenshot
description: "Best practices for using the view_screenshot tool in Radon IDE. Captures a PNG screenshot of the app running in the development viewport. Use when you need to inspect the current visual state of the app — verifying UI changes, checking layout, confirming styling, or diagnosing visual bugs. Trigger on: 'screenshot', 'what does the app look like', 'show me the screen', 'visual state', 'UI preview', 'check the layout', 'see the device', 'view the app', or any request to visually inspect the running application."
---

# view_screenshot

Captures a screenshot of the app running in the Radon IDE development viewport and returns it as a PNG image.

## Tool signature

```
view_screenshot()
```

No input parameters.

**Returns:** PNG image of the current device screen.

---

## When to use

- After making UI or styling changes, to confirm they rendered correctly.
- When the user asks "what does the app look like" or wants to verify the current visual state.
- When diagnosing layout issues, misaligned components, or visual regressions.
- As a follow-up to `reload_application` to confirm the app recovered after a crash.
- When `view_application_logs` alone is not enough context for debugging a visual problem.

## When NOT to use

- To debug logic errors that have no visual manifestation — use `view_application_logs` instead.
- Repeatedly in a tight loop — one screenshot per logical step is sufficient.

---

## Best practices

### Pair with logs for full context

When debugging an issue, call `view_application_logs` alongside `view_screenshot`. The logs tool already includes a screenshot when the preview is ready, so calling both is redundant only when you already have the logs output. If you only need the visual state and not the logs, `view_screenshot` is the lighter-weight option.

### Use after reload

After calling `reload_application`, use `view_screenshot` to verify the app has returned to a healthy visual state. The reload tool waits for the app status to return to "running", but a screenshot confirms the UI is actually rendered correctly.

### Interpret the image carefully

The screenshot captures the full device viewport with the current rotation setting. Keep in mind:

- The image reflects the device's current orientation (portrait, landscape left, landscape right, or upside-down portrait).
- Dark mode vs light mode depends on the device appearance setting.
- The content size setting affects text and element scaling.

---

## Error handling

If the development viewport device is turned off, the tool returns a text error:

```
Could not capture a screenshot!
The development viewport device is likely turned off.
Please turn on the Radon IDE emulator before proceeding.
```

**Resolution:** Ask the user to open the Radon IDE panel and start a device before retrying.
