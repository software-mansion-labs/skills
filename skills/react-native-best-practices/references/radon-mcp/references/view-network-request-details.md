---
name: view-network-request-details
description: "Best practices for using the view_network_request_details tool in Radon IDE. Returns full details of a specific network request including headers, body, and metadata. Use after view_network_logs to drill into a specific request for debugging. Trigger on: 'request details', 'request headers', 'response body', 'request body', 'inspect request', 'API response', or any follow-up to viewing network logs where you need the full details of a particular request."
---

# view_network_request_details

Returns the full details of a specific network request, including headers, body, and metadata. Use after `view_network_logs` to inspect a request identified by its `requestId`.

## Tool signature

```
view_network_request_details({ requestId: "<string>" })
```

**Input:**

- `requestId` (required): The ID of the network request, obtained from the output of `view_network_logs`.

**Returns:** Pretty-printed JSON containing the full request and response data.

---

## When to use

- After `view_network_logs` reveals a suspicious or failed request.
- To inspect request headers, query parameters, and request body.
- To read the response body and verify the data the app received.
- To debug authentication issues by checking auth-related headers (note: sensitive values are redacted).
- To verify content types, caching headers, or other HTTP metadata.

## When NOT to use

- As a first step — always call `view_network_logs` first to identify which request to inspect.
- To get an overview of all network activity — use `view_network_logs` for that.

---

## Understanding the output

The output is a pretty-printed JSON object representing the full network log entry, with two important transformations applied automatically:

### Sensitive header redaction

Headers matching any of these patterns (case-insensitive) are replaced with `[redacted]`:

- `auth` (matches `Authorization`, `X-Auth-Token`, etc.)
- `cookie` (matches `Cookie`, `Set-Cookie`)
- `token`
- `secret`
- `key` (matches `X-API-Key`, etc.)
- `session`
- `credential`

This applies to both request and response headers.

Example:

```json
{
  "Authorization": "[redacted]",
  "Content-Type": "application/json",
  "X-API-Key": "[redacted]"
}
```

### Large response body truncation

If the response body exceeds 1000 characters when serialized, it is replaced with a placeholder:

```
[RESPONSE CONTENT HIDDEN. MIME TYPE: application/json. ORIGINAL SIZE: 45230 CHARACTERS]
```

This prevents large API responses from overwhelming the AI context window.

---

## Best practices

### Always call view_network_logs first

The `requestId` is only available from the `view_network_logs` output. The typical workflow is:

1. Call `view_network_logs({ pageIndex: "latest" })` to see recent requests.
2. Identify the problematic request by its URL, status code, or timing.
3. Call `view_network_request_details({ requestId: "<id>" })` with the ID from step 2.

### Be aware of redacted headers

When debugging authentication issues, remember that auth headers show as `[redacted]`. You can still verify:

- Whether the header exists at all (presence vs absence).
- The header name and structure.
- Non-sensitive headers that might affect auth (e.g., `Origin`, `Referer`).

If you need the actual auth header value, you'll need to inspect it through the code or the Radon IDE Network panel directly.

### Handle truncated response bodies

If the response body is truncated, you know:

- The MIME type of the response.
- The original character count of the response.

For large responses, focus on the status code and headers. If the response content is critical, suggest the user check the Network panel in Radon IDE for the full payload.

---

## Error handling

If the request ID doesn't exist:

```
Request with id abc123 does not exist.
```

If the device is off:

```
Could not retrieve network logs!
The development device is likely turned off.
Please turn on the Radon IDE emulator before proceeding.
```

**Resolution:** Verify the `requestId` was copied correctly from `view_network_logs`. If the device error appears, ensure Radon IDE is running.
