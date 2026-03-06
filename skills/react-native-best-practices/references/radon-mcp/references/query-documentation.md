---
name: query-documentation
description: "Best practices for using the query_documentation tool in Radon IDE. Returns documentation snippets relevant to a provided query from a curated React Native and Expo knowledge base. Use when you need accurate, up-to-date documentation for React Native APIs, Expo modules, or related libraries. Trigger on: 'how to use', 'documentation for', 'API reference', 'React Native docs', 'Expo docs', 'how does X work in React Native', 'what is the API for', or any question about React Native or Expo library usage that benefits from authoritative documentation."
---

# query_documentation

Returns documentation snippets relevant to a provided text query. Queries a curated React Native and Expo knowledge base hosted on the Radon AI backend.

## Tool signature

```
query_documentation({ text: "<query string>" })
```

**Input:**

- `text` (required): A natural language query describing what documentation you need.

**Returns:** Text containing relevant documentation snippets.

---

## When to use

- When you need accurate, up-to-date API documentation for a React Native or Expo feature.
- Before implementing a feature that uses a library you're not fully familiar with.
- When the user asks "how do I do X in React Native" and you need authoritative guidance.
- To verify correct API usage, prop types, or configuration options.
- When debugging issues that might stem from incorrect API usage.

## When NOT to use

- For general programming questions unrelated to React Native / Expo.
- When you already have high confidence in the API from your training data and the API is unlikely to have changed.
- To look up a specific library's npm package info — use `get_library_description` for that.

---

## Best practices

### Write specific, focused queries

Good queries target a specific API, feature, or use case:

```
// Good
query_documentation({ text: "React Navigation stack navigator configuration options" })
query_documentation({ text: "Expo Camera permissions setup on iOS and Android" })
query_documentation({ text: "FlatList performance optimization techniques" })

// Too vague
query_documentation({ text: "React Native" })
query_documentation({ text: "navigation" })
```

### Include the library name when relevant

If the question is about a specific library, include its name:

```
query_documentation({ text: "react-native-reanimated useSharedValue hook" })
query_documentation({ text: "expo-image-picker launchImageLibraryAsync options" })
```

### Use documentation before writing implementation

When implementing a feature that relies on a React Native or Expo API, query the documentation first to ensure you're using the latest API surface and best practices. APIs change between versions, and the documentation query returns up-to-date information.

### Trust documentation over training data for version-specific details

React Native and Expo evolve rapidly. When there's a discrepancy between what you know from training data and what the documentation says, prefer the documentation — it reflects the current state.

---

## Requirements

- A valid Radon IDE license is required. Without one, the tool returns:

```
You have to have a valid Radon IDE license to use the query_documentation tool.
```

- Network connectivity is required — the tool queries the Radon AI backend service.

---

## Error handling

If the license is invalid or missing, the tool returns a license-required message.

If the backend is unreachable (network failure), the tool will fail with a network error.

**Resolution:** Ensure the user has an active Radon IDE license and internet connectivity.
