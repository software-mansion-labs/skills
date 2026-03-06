---
name: view-component-tree
description: "Best practices for using the view_component_tree tool in Radon IDE. Displays the React component tree (view hierarchy) of the running app, filtered to user-authored components. Use for understanding mounted component structure, resolving layout issues, finding context providers, or mapping file structure to component hierarchy. Trigger on: 'component tree', 'view hierarchy', 'component structure', 'React tree', 'mounted components', 'layout structure', 'where is this component', 'context provider location', 'component hierarchy', or any request to understand the UI structure of the running app."
---

# view_component_tree

Displays the React component tree (view hierarchy) of the running app, filtered to show only user-authored components with source file locations.

## Tool signature

```
view_component_tree()
```

No input parameters.

**Returns:** An indented XML-like text representation of the mounted component tree.

---

## When to use

- When you need a general overview of the UI structure.
- To understand the relationship between the project file structure and the component hierarchy.
- When resolving layout issues — see which components are mounted and how they nest.
- To locate context providers in the tree.
- To find where a specific component is used and how it fits into the broader hierarchy.
- When debugging issues that depend on component mount state.

## When NOT to use

- To see components that are not currently mounted — the tool only shows what is on screen.
- For performance profiling — use the React Profiler instead.
- To read component props or state values — the tree shows structure, not data.

---

## Understanding the output

The output is an indented XML-like tree. Each element includes:

### Component name and source location

```xml
<MyComponent> (src/screens/HomeScreen.tsx:42)
  <Header> (src/components/Header.tsx:15)
  </Header>
</MyComponent>
```

The path is relative to the workspace root. The line number points to the component definition.

### HOC descriptors

Higher-order component wrappers are shown in brackets:

```xml
<MyButton> [memo, forwardRef] (src/components/MyButton.tsx:8)
```

### Text content for leaf nodes

Non-user components that render text show their content:

```xml
<Label> (src/components/Label.tsx:12)
  Hello World
</Label>
```

### Self-closing tags for childless components

```xml
<Icon /> (src/components/Icon.tsx:5)
```

---

## Best practices

### Only user-authored components are shown

The tree is filtered to exclude components from `node_modules` and other non-workspace files. This means:

- React Native primitives (`View`, `Text`, `ScrollView`) are **not** shown unless they are wrapped in a user component that passes them through.
- Third-party library components are excluded.
- The tree focuses on your application's own component structure.

A component is considered "user-related" if it is defined in the workspace (not in `node_modules`) or if it is owned by a component that is user-defined.

### The tree starts from the app entry point

The tool looks for an `__RNIDE_APP_WRAPPER` component as the entry point. If found, the tree starts there. Otherwise, it falls back to the React DevTools root. This means framework wrapper components above your app entry point are typically excluded.

### Only mounted components are visible

If a component is conditionally rendered and its condition is currently `false`, it will not appear in the tree. Navigate to different screens or states in the app to see different parts of the tree.

### Use alongside view_screenshot for layout debugging

When debugging layout issues, combine `view_component_tree` with `view_screenshot`:

1. Take a screenshot to see the visual problem.
2. View the component tree to understand the structure causing it.
3. Use the source file paths from the tree to navigate directly to the relevant code.

---

## Error handling

If the app is not running or DevTools are not accessible:

```
Could not extract the component tree from the app, the app is not running!
The development device is likely turned off.
Please turn on the Radon IDE emulator before proceeding.
```

If the component tree is corrupted:

```
Component tree is corrupted. Tree root could not be found.
```

**Resolution:** Reload the app with `reload_application` using `reloadJs`, then retry.
