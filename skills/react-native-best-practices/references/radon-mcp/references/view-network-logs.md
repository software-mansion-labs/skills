---
name: view-network-logs
description: "Best practices for using the view_network_logs tool in Radon IDE. Returns a paginated list of all network requests made by the app, including method, URL, status, duration, and size. Use when debugging networking issues, inspecting API calls, or verifying that the app communicates correctly with backend services. Trigger on: 'network logs', 'network requests', 'API calls', 'HTTP requests', 'network inspector', 'fetch requests', 'network traffic', 'request failed', 'status code', '404', '500', 'CORS', or any request to inspect network activity from the running app."
---

# view_network_logs

Returns a paginated list of all network requests made by the running app, including method, URL, status code, duration, and size.

## Tool signature

```
view_network_logs({ pageIndex: "latest" | "<number>" })
```

**Input:**

- `pageIndex` (required): Either `"latest"` for the most recent page, or a 0-based page index string (e.g., `"0"`, `"1"`, `"2"`).

**Returns:** Formatted text listing network requests for the requested page, with a page indicator header.

---

## When to use

- When debugging API call failures (wrong status codes, timeouts, missing responses).
- To verify that the app is making the correct requests to the right endpoints.
- When the user reports that data is not loading or is stale.
- To check request timing and identify slow API calls.
- As the first step before drilling into a specific request with `view_network_request_details`.

## When NOT to use

- To inspect request/response headers and bodies — use `view_network_request_details` after identifying the request.
- When the issue is clearly a build or JS error — use `view_application_logs` instead.

---

## Understanding the output

### Page header

```
=== NETWORK LOGS (page 3/5) ===
```

Indicates current page and total page count. Each page contains up to 50 entries.

### Entry format

Each network request is formatted as:

```
{id: <requestId>} "METHOD URL" STATUS statusText TYPE SIZE DURATION
```

Example:

```
{id: abc123} "GET https://api.example.com/users" 200 OK json 1.2kB 150ms
```

The `requestId` is used to drill into details with `view_network_request_details`.

---

## Best practices

### Start with `"latest"` for recent issues

When the user reports a current problem, use `pageIndex: "latest"` to see the most recent network activity. This is almost always the right starting point.

```
view_network_logs({ pageIndex: "latest" })
```

### Use numeric page indexes for historical search

If you need to find an older request, start from page `"0"` and work forward. The page header tells you the total page count.

### Follow up with view_network_request_details

Once you identify a suspicious request in the logs (wrong status code, unexpected URL, etc.), note its `requestId` and call `view_network_request_details` to see the full headers, body, and metadata.

### Ensure the network inspector is enabled

If no logs are returned, the tool will report:

```
No network traffic recorded. Make sure the network inspector is enabled before accessing these logs.
```

The network inspector must be active for requests to be captured. It is enabled by default when the Network panel is available in Radon IDE.

---

## Error handling

If the device is off:

```
Could not retrieve network logs!
The development device is likely turned off.
Please turn on the Radon IDE emulator before proceeding.
```

If the network inspector plugin is unavailable:

```
Network inspector plugin is not available.
```

If the page index is invalid:

```
"abc" is not a valid page index value.
```

If the page index is out of range:

```
Page index out of range. Valid range: 0-4 (5 pages total).
```

**Resolution:** For device/plugin errors, ensure Radon IDE is running with a device and the Network panel is available. For pagination errors, use `"latest"` or a valid page number.
