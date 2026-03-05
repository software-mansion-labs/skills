---
name: rich-text
description: "Software Mansion's best practices for rich text editing in React Native apps using react-native-enriched. Use when building rich text editors, formatted text inputs, or any feature requiring inline styling, mentions, links, or structured text editing. Trigger on: 'rich text editor', 'rich text input', 'text editor', 'react-native-enriched', 'EnrichedTextInput', 'formatted text input', 'WYSIWYG', 'mentions input', 'text formatting toolbar', or any request to build a rich text editing experience in React Native."
---

# Rich Text Editing

Software Mansion's production rich text editing patterns for React Native using `react-native-enriched`.

## Use react-native-enriched

`react-native-enriched` is a native rich text editor for React Native. It provides `EnrichedTextInput`, a component that supports inline styles (bold, italic, underline, strikethrough, inline code), block-level formatting (headings, block quotes, code blocks, ordered/unordered/checkbox lists), mentions, links, and images.

```bash
npm install react-native-enriched
```

```tsx
import { EnrichedTextInput } from 'react-native-enriched';
import type {
  EnrichedTextInputInstance,
  OnChangeStateEvent,
} from 'react-native-enriched';
import { useState, useRef } from 'react';
import { View, Button, StyleSheet } from 'react-native';

export default function App() {
  const ref = useRef<EnrichedTextInputInstance>(null);

  const [stylesState, setStylesState] = useState<OnChangeStateEvent | null>();

  return (
    <View style={styles.container}>
      <EnrichedTextInput
        ref={ref}
        onChangeState={(e) => setStylesState(e.nativeEvent)}
        style={styles.input}
      />
      <Button
        title={stylesState?.bold.isActive ? 'Unbold' : 'Bold'}
        color={stylesState?.bold.isActive ? 'green' : 'gray'}
        onPress={() => ref.current?.toggleBold()}
      />
    </View>
  );
}

const styles = StyleSheet.create({
  container: {
    flex: 1,
    justifyContent: 'center',
    alignItems: 'center',
  },
  input: {
    width: '100%',
    fontSize: 20,
    padding: 10,
    maxHeight: 200,
    backgroundColor: 'lightgray',
  },
});
```

---

## References

- [react-native-enriched on GitHub](https://github.com/software-mansion/react-native-enriched)
- [API Reference](https://github.com/software-mansion/react-native-enriched/blob/main/docs/API_REFERENCE.md)
