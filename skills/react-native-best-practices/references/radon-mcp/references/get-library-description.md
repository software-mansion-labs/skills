---
name: get-library-description
description: "Best practices for using the get_library_description tool in Radon IDE. Returns a detailed description of an npm library and its use cases. Use when evaluating whether to use a specific library, understanding what a dependency does, or recommending libraries for a particular task. Trigger on: 'what is this library', 'what does X package do', 'should I use', 'library description', 'npm package info', 'describe library', 'library use cases', or any request to understand the purpose and capabilities of a specific npm package."
---

# get_library_description

Returns a detailed description of an npm library and its use cases. Queries the Radon AI backend with the library's npm package name.

## Tool signature

```
get_library_description({ library_npm_name: "<npm package name>" })
```

**Input:**

- `library_npm_name` (required): The npm package name (e.g., `"react-native-reanimated"`, `"expo-camera"`, `"@react-navigation/native"`).

**Returns:** Text containing a detailed description of the library and its use cases.

---

## When to use

- When the user asks "what does library X do?" or "should I use library Y?".
- When evaluating which library to recommend for a specific task.
- To understand a dependency found in the project's `package.json`.
- Before adding a new dependency, to confirm it's the right tool for the job.
- When comparing alternatives for a particular feature.

## When NOT to use

- For API documentation and usage examples — use `query_documentation` instead.
- When you already know the library well and the user didn't ask about it.
- For libraries unrelated to React Native / Expo — the backend is focused on the RN ecosystem.

---

## Best practices

### Use the exact npm package name

The tool expects the npm package name as published on the npm registry:

```
// Correct
get_library_description({ library_npm_name: "react-native-reanimated" })
get_library_description({ library_npm_name: "@react-navigation/native" })
get_library_description({ library_npm_name: "expo-camera" })

// Incorrect
get_library_description({ library_npm_name: "Reanimated" })
get_library_description({ library_npm_name: "react navigation" })
```

### Use for evaluation, not implementation

This tool gives you a high-level understanding of what a library does and when to use it. For actual implementation guidance (API usage, configuration, code examples), follow up with `query_documentation`.

Typical workflow:

1. `get_library_description` — understand what the library does and whether it fits the need.
2. `query_documentation` — get specific API docs and usage patterns.

### Helpful for project dependency audits

When the user asks "what are all these dependencies?" or wants to understand their `package.json`, use this tool to describe each relevant dependency.

---

## Requirements

- A valid Radon IDE license is required. Without one, the tool returns:

```
You have to have a valid Radon IDE license to use the get_library_description tool.
```

- Network connectivity is required — the tool queries the Radon AI backend service.

---

## Error handling

If the license is invalid or missing, the tool returns a license-required message.

If the backend is unreachable (network failure), the tool will fail with a network error.

**Resolution:** Ensure the user has an active Radon IDE license and internet connectivity.
