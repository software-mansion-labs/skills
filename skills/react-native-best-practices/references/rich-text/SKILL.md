---
name: rich-text
description: "Software Mansion's best practices for rich text in React Native using react-native-enriched-html (formerly react-native-enriched) and react-native-enriched-markdown. Use when building rich text editors, formatted text inputs, HTML or Markdown renderers, or any feature requiring inline styling, mentions, links, structured text editing, or Markdown display. Trigger on: 'rich text editor', 'rich text input', 'text editor', 'react-native-enriched', 'react-native-enriched-html', 'react-native-enriched-markdown', 'EnrichedTextInput', 'EnrichedText', 'EnrichedMarkdownText', 'formatted text input', 'WYSIWYG', 'mentions input', 'text formatting toolbar', 'context menu items', 'text shortcuts', 'markdown renderer', 'markdown display', 'render markdown', 'display markdown natively', 'render HTML natively', 'LaTeX math', 'GFM tables', or any request to build rich text editing or Markdown rendering in React Native."
---

# Rich Text in React Native

Software Mansion's production patterns for rich text editing and rendering in React Native.

There are two libraries that cover rich text use cases. Both provide an **input** (editing) component and a **display** (rendering) component - the difference is the underlying content format:

| Library                          | Format   | Input component             | Display component      |
| -------------------------------- | -------- | --------------------------- | ---------------------- |
| `react-native-enriched-html`     | HTML     | `EnrichedTextInput`         | `EnrichedText`         |
| `react-native-enriched-markdown` | Markdown | `EnrichedMarkdownTextInput` | `EnrichedMarkdownText` |

Both libraries require the React Native New Architecture (Fabric). `react-native-enriched-html` supports iOS, Android and Web; `react-native-enriched-markdown` supports iOS, Android, macOS, and Web.

> **Package rename**: `react-native-enriched` was renamed to `react-native-enriched-html` in v1.0.0. The old npm package is deprecated and receives no updates. Migration has no breaking changes - swap the dependency and update imports. If a project still depends on `react-native-enriched`, recommend migrating.

## Choosing the right library

Choose based on the **content format** your app stores and exchanges:

- **Content is HTML**, (bold, italic, mentions, links, inline images, headings, lists, text alignment, context menu items): use `react-native-enriched-html`
- **Content is Markdown** (chat messages, documentation, AI responses): use `react-native-enriched-markdown`
- **Don't mix the two libraries in one flow**: an input's output should feed its sibling display component directly (HTML input -> HTML display, Markdown input -> Markdown display)

## react-native-enriched-html (HTML editor + display)

`EnrichedTextInput` is a native, uncontrolled rich text input. It directly interacts with platform-specific components for performance, meaning it does not use React state for its value. Its HTML output is rendered 1:1 by the companion `EnrichedText` display component.

```bash
npm install react-native-enriched-html
# Expo: npx expo install react-native-enriched-html && npx expo prebuild (won't work in Expo Go)
```

### Basic usage

```tsx
import { EnrichedTextInput } from "react-native-enriched-html";
import type {
  EnrichedTextInputInstance,
  OnChangeStateEvent,
} from "react-native-enriched-html";
import { useState, useRef } from "react";
import { View, Button, StyleSheet } from "react-native";

export default function RichEditor() {
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
        title={stylesState?.bold.isActive ? "Unbold" : "Bold"}
        color={stylesState?.bold.isActive ? "green" : "gray"}
        onPress={() => ref.current?.toggleBold()}
      />
    </View>
  );
}
```

### Key concepts

**Toggling styles via ref**: All formatting is applied imperatively through the ref. Call `ref.current?.toggleBold()`, `ref.current?.toggleItalic()`, etc.

**Style detection via onChangeState**: The `onChangeState` callback fires whenever the style state changes (e.g., cursor moves into bold text). Each style reports three properties:

- `isActive`: The style is applied at the current selection (highlight the toolbar button)
- `isBlocking`: The style is blocked by another active style (disable the toolbar button)
- `isConflicting`: The style conflicts with another active style (toggling it removes the conflicting style)

**Inline vs paragraph styles**:

- Inline styles (bold, italic, underline, strikethrough, inline code) apply to the exact character range selected. With no selection, they apply to the next characters typed.
- Paragraph styles (headings, codeblock, blockquote, lists) apply to entire paragraphs (text between newlines). If the selection spans multiple paragraphs, all are affected.
- Only one paragraph style can be active per paragraph - all paragraph styles conflict with each other. Codeblock additionally blocks all inline styles.

**HTML output**: Get HTML via `ref.current?.getHTML()` (on-demand, returns a Promise) or the `onChangeHtml` callback (continuous, has performance cost for large documents). Prefer `getHTML()` when you only need HTML at save time. The same applies to `onChangeText` - skip it unless you need the plain text continuously.

**Setting content**: Use `defaultValue` prop for initial HTML content, or `ref.current?.setValue(html)` to update imperatively. Both accept the library's own HTML output or raw text. Arbitrary external HTML (Google Docs, Word, web pages) is normalized into the supported tag subset by default (`useHtmlNormalizer` is enabled out of the box; set it to `false` to opt out).

### Supported styles

The `OnChangeStateEvent` key column shows the exact property name on the event object returned by `onChangeState`. Use these keys when reading style state (e.g. `stylesState.strikeThrough.isActive`). Note that casing varies (e.g. `strikeThrough` with capital T, `inlineCode` with capital C).

| Style          | Toggle method                 | `OnChangeStateEvent` key | Type      |
| -------------- | ----------------------------- | ------------------------ | --------- |
| Bold           | `toggleBold()`                | `bold`                   | Inline    |
| Italic         | `toggleItalic()`              | `italic`                 | Inline    |
| Underline      | `toggleUnderline()`           | `underline`              | Inline    |
| Strikethrough  | `toggleStrikeThrough()`       | `strikeThrough`          | Inline    |
| Inline code    | `toggleInlineCode()`          | `inlineCode`             | Inline    |
| H1             | `toggleH1()`                  | `h1`                     | Paragraph |
| H2             | `toggleH2()`                  | `h2`                     | Paragraph |
| H3             | `toggleH3()`                  | `h3`                     | Paragraph |
| H4             | `toggleH4()`                  | `h4`                     | Paragraph |
| H5             | `toggleH5()`                  | `h5`                     | Paragraph |
| H6             | `toggleH6()`                  | `h6`                     | Paragraph |
| Code block     | `toggleCodeBlock()`           | `codeBlock`              | Paragraph |
| Block quote    | `toggleBlockQuote()`          | `blockQuote`             | Paragraph |
| Ordered list   | `toggleOrderedList()`         | `orderedList`            | Paragraph |
| Unordered list | `toggleUnorderedList()`       | `unorderedList`          | Paragraph |
| Checkbox list  | `toggleCheckboxList(checked)` | `checkboxList`           | Paragraph |

### Text alignment

Alignment is not a toggle - set it with `setTextAlignment` and read the current value from the `alignment` string on `OnChangeStateEvent` (`'left' | 'center' | 'right' | 'justify' | 'auto'`):

```tsx
ref.current?.setTextAlignment("center");
// reset to natural alignment
ref.current?.setTextAlignment("auto");
```

Alignment applies to the paragraph(s) at the current selection; inside a list it applies to all contiguous list items. `'justify'` works on iOS only - on Android it falls back to natural alignment.

### Text shortcuts

The `textShortcuts` prop auto-converts typed patterns into styles, like Markdown editors do. Default: `- ` starts an unordered list and `1. ` starts an ordered list. Pass an empty array to disable.

```tsx
<EnrichedTextInput
  ref={ref}
  textShortcuts={[
    { trigger: "- ", style: "unordered_list" },
    { trigger: "1. ", style: "ordered_list" },
    { trigger: "# ", style: "h1" },
    { trigger: "> ", style: "blockquote" },
    { trigger: "**", style: "bold" },
    { trigger: "~~", style: "strikethrough" },
  ]}
/>
```

Paragraph triggers (`h1`-`h6`, `blockquote`, `codeblock`, lists) fire at the start of a plain paragraph. Inline triggers (`bold`, `italic`, `underline`, `strikethrough`, `inline_code`) fire when the closing delimiter is typed around text (e.g. `**text**`). Blocking/conflicting rules still apply.

### Context menu items

Extend the native text editing menu (the one with Copy/Paste/Cut) with custom items via the `contextMenuItems` prop (Android and iOS 16+). Each item's `onPress` receives the selected `text`, the `selection` range, and the current `styleState` (an `OnChangeStateEvent`), so you can act on the selection directly:

```tsx
<EnrichedTextInput
  ref={ref}
  contextMenuItems={[
    {
      text: "Paste Link",
      visible: true, // optional, defaults to true
      onPress: ({ text, selection, styleState }) => {
        if (!styleState.link.isBlocking) {
          ref.current?.setLink(selection.start, selection.end, text, url);
        }
      },
    },
  ]}
/>
```

On iOS items appear in array order before the system items. On Android the order is not guaranteed and items may land in a submenu depending on the device manufacturer.

### Links

Links are detected automatically (customizable via `linkRegex` prop) or applied manually:

```tsx
// Set link on selected text
ref.current?.setLink(selection.start, selection.end, selectedText, url);

// Remove link
ref.current?.removeLink(start, end);
```

Use `onChangeSelection` to get selection position and `onLinkDetected` to detect when the cursor is near a link.

### Mentions

Mentions support custom indicators (default: `@`). Set custom indicators via the `mentionIndicators` prop.

```tsx
<EnrichedTextInput
  ref={ref}
  mentionIndicators={["@", "#"]}
  onStartMention={(indicator) => {
    /* show picker */
  }}
  onChangeMention={({ indicator, text }) => {
    /* filter list */
  }}
  onEndMention={(indicator) => {
    /* hide picker */
  }}
/>;

// Complete the mention when user selects from picker
ref.current?.setMention("@", "John Doe", { userId: "123" });

// Or start a mention programmatically (e.g. from a toolbar button)
ref.current?.startMention("@");
```

Custom `attributes` passed to `setMention` are preserved through HTML parsing, so you can round-trip IDs. Use `onMentionDetected` to know when the cursor is near an existing mention.

### Inline images

```tsx
ref.current?.setImage(imageUri, width, height);
```

Images are inserted at the cursor position (or replace selected text) and affect line height. You are responsible for providing correct dimensions.

Handle images pasted by the user with `onPasteImages`, which returns the URI, MIME type, and dimensions of each pasted image/GIF - typically you upload them and re-insert with `setImage`.

### Submit behavior

By default the return key inserts a newline. For chat-like inputs, use `submitBehavior` with `onSubmitEditing`:

```tsx
<EnrichedTextInput
  ref={ref}
  submitBehavior="submit" // or 'blurAndSubmit'; default 'newline'
  onSubmitEditing={async () => {
    const html = await ref.current?.getHTML();
    sendMessage(html);
  }}
/>
```

### Styling the input

Two separate props control appearance:

- `style` - container-level layout and base typography (dimensions, padding, borders, `fontSize`, `color`, `fontFamily`, ...). A subset of `TextStyle`; `textAlign` and `textDecorationLine` are not supported here.
- `htmlStyle` - per-tag styling of the rich text content itself:

```tsx
<EnrichedTextInput
  ref={ref}
  style={{ fontSize: 16, padding: 12 }}
  htmlStyle={{
    h1: { fontSize: 28, bold: true },
    a: { color: "#007AFF", textDecorationLine: "underline" },
    code: { color: "#E91E63", backgroundColor: "#F5F5F5" },
    codeblock: {
      backgroundColor: "#1E1E1E",
      color: "#D4D4D4",
      borderRadius: 8,
    },
    blockquote: { borderColor: "#007AFF", gapWidth: 12 },
    // Per-indicator mention styling (or pass one config for all indicators)
    mention: {
      "@": { color: "blue", backgroundColor: "transparent" },
      "#": { color: "green", backgroundColor: "transparent" },
    },
  }}
/>
```

Supported `htmlStyle` keys: `h1`-`h6`, `blockquote`, `codeblock`, `code`, `a`, `mention`, `ol`, `ul`, `ulCheckbox`.

### EnrichedText (display component)

`EnrichedText` renders the HTML produced by `EnrichedTextInput` with the same styling API, guaranteeing 1:1 visual consistency between editing and display. Pass the HTML string as `children`:

```tsx
import { EnrichedText } from "react-native-enriched-html";

<EnrichedText
  style={{ fontSize: 16 }}
  htmlStyle={{
    a: { color: "#007AFF", pressColor: "#004999" },
    mention: { color: "blue", pressBackgroundColor: "#E0E0FF" },
  }}
  numberOfLines={3}
  ellipsizeMode="tail"
  selectable
  onLinkPress={({ url }) => Linking.openURL(url)}
  onMentionPress={({ text, indicator, attributes }) =>
    openProfile(attributes.userId)
  }
>
  {html}
</EnrichedText>;
```

Its `htmlStyle` extends the input's with press-state colors (`pressColor` on links, `pressColor`/`pressBackgroundColor` on mentions). Like the input, it normalizes external HTML by default (`useHtmlNormalizer`).

### Full API reference

For the complete API, webfetch the relevant upstream doc:

- [README (usage, supported tags, conflict/blocking tables)](https://github.com/software-mansion/react-native-enriched-html/blob/main/README.md)
- [EnrichedTextInput API Reference (all props, ref methods, events, HtmlStyle)](https://github.com/software-mansion/react-native-enriched-html/blob/main/docs/INPUT_API_REFERENCE.md)
- [EnrichedText API Reference](https://github.com/software-mansion/react-native-enriched-html/blob/main/docs/TEXT_API_REFERENCE.md)
- [EnrichedInputStyle (supported `style` prop properties)](https://github.com/software-mansion/react-native-enriched-html/blob/main/docs/ENRICHED_INPUT_STYLE.md)
- [Web support](https://github.com/software-mansion/react-native-enriched-html/blob/main/docs/WEB.md)

---

## react-native-enriched-markdown (Renderer)

`EnrichedMarkdownText` renders Markdown as fully native text (no WebView). It uses [md4c](https://github.com/mity/md4c) for high-performance CommonMark-compliant parsing.

```bash
npm install react-native-enriched-markdown
```

### Basic usage

```tsx
import { EnrichedMarkdownText } from "react-native-enriched-markdown";
import { Linking } from "react-native";

const markdown = `
# Welcome

This is **bold**, *italic*, and [a link](https://reactnative.dev).

- List item one
- List item two
`;

export default function MarkdownDisplay() {
  return (
    <EnrichedMarkdownText
      markdown={markdown}
      onLinkPress={({ url }) => Linking.openURL(url)}
    />
  );
}
```

### Flavors

| Flavor                 | Features                                                    | Layout                                              |
| ---------------------- | ----------------------------------------------------------- | --------------------------------------------------- |
| `commonmark` (default) | All CommonMark elements, inline math (`$...$`)              | Single TextView                                     |
| `github`               | CommonMark + GFM tables, task lists, block math (`$$...$$`) | Segmented layout (separate TextViews + table views) |

```tsx
<EnrichedMarkdownText
  flavor="github"
  markdown={markdown}
  onLinkPress={({ url }) => Linking.openURL(url)}
/>
```

### Tables (GFM)

Tables require `flavor="github"` and support column alignment, rich text in cells, horizontal scrolling, header styling, alternating row colors, and a long-press context menu.

```tsx
<EnrichedMarkdownText
  flavor="github"
  markdown={tableMarkdown}
  markdownStyle={{
    table: {
      fontSize: 14,
      borderColor: "#E5E7EB",
      borderRadius: 8,
      headerBackgroundColor: "#F3F4F6",
      cellPaddingHorizontal: 12,
      cellPaddingVertical: 8,
    },
  }}
/>
```

### Task lists (GFM)

Task lists require `flavor="github"`. Handle checkbox taps with `onTaskListItemPress`:

```tsx
<EnrichedMarkdownText
  flavor="github"
  markdown={`- [x] Done\n- [ ] Todo`}
  onTaskListItemPress={({ index, checked, text }) => {
    console.log(`Task ${index}: ${checked ? "checked" : "unchecked"}`);
  }}
/>
```

### LaTeX math

- **Inline math** (`$...$`): works in both flavors
- **Block math** (`$$...$$`): requires `flavor="github"`, must be on its own line

Use `String.raw` or double backslashes for LaTeX commands in JS strings.

Disable math to reduce bundle size (iosMath ~2.5 MB on iOS):

```tsx
<EnrichedMarkdownText md4cFlags={{ latexMath: false }} markdown="..." />
```

### Customizing styles

Use the `markdownStyle` prop. Memoize it with `useMemo` to avoid re-renders:

```tsx
import type { MarkdownStyle } from "react-native-enriched-markdown";

const markdownStyle: MarkdownStyle = useMemo(
  () => ({
    paragraph: { fontSize: 16, color: "#333", lineHeight: 24 },
    h1: { fontSize: 32, fontWeight: "bold", marginBottom: 16 },
    code: { color: "#E91E63", backgroundColor: "#F5F5F5" },
    codeBlock: {
      backgroundColor: "#1E1E1E",
      color: "#D4D4D4",
      padding: 16,
      borderRadius: 8,
    },
    link: { color: "#007AFF", underline: true },
    blockquote: { borderColor: "#007AFF", backgroundColor: "#F0F8FF" },
  }),
  []
);
```

Inline elements inherit typography from their parent block (fontSize, fontFamily, color), then add their own styling on top.

### Link handling

```tsx
<EnrichedMarkdownText
  markdown={content}
  onLinkPress={({ url }) => Linking.openURL(url)}
  onLinkLongPress={({ url }) => showShareSheet(url)}
  // On iOS, providing onLinkLongPress disables the system link preview
/>
```

### Additional features

- **Text selection and copy**: Enabled by default (`selectable` prop). Smart Copy provides plain text, Markdown, HTML, RTF, and RTFD on iOS.
- **Accessibility**: VoiceOver and TalkBack support with custom rotors (iOS), semantic heading/link/image navigation, and proper list announcements.
- **RTL**: Automatic on Android. On iOS, call `I18nManager.forceRTL(true)` before rendering.
- **Image caching**: Three-tier caching (memory originals, memory processed variants, disk) with request deduplication.
- **Underline mode**: `md4cFlags={{ underline: true }}` makes `_text_` render as underline instead of italic.

For detailed API documentation, webfetch the relevant page from the upstream docs:

- [API Reference (props, events)](https://github.com/software-mansion-labs/react-native-enriched-markdown/blob/main/docs/API_REFERENCE.md)
- [Styles (MarkdownStyle properties)](https://github.com/software-mansion-labs/react-native-enriched-markdown/blob/main/docs/STYLES.md)
- [Elements Structure (Markdown-to-native mapping)](https://github.com/software-mansion-labs/react-native-enriched-markdown/blob/main/docs/ELEMENTS_STRUCTURE.md)
- [Accessibility (VoiceOver, TalkBack)](https://github.com/software-mansion-labs/react-native-enriched-markdown/blob/main/docs/ACCESSIBILITY.md)
- [RTL Support](https://github.com/software-mansion-labs/react-native-enriched-markdown/blob/main/docs/RTL.md)
- [Image Caching](https://github.com/software-mansion-labs/react-native-enriched-markdown/blob/main/docs/IMAGE_CACHING.md)
- [Copy Options (Smart Copy, copy-as-Markdown)](https://github.com/software-mansion-labs/react-native-enriched-markdown/blob/main/docs/COPY_OPTIONS.md)
- [LaTeX Math](https://github.com/software-mansion-labs/react-native-enriched-markdown/blob/main/docs/LATEX_MATH.md)
- [macOS Support](https://github.com/software-mansion-labs/react-native-enriched-markdown/blob/main/docs/MACOS.md)

---

## Known limitations

### react-native-enriched-html

- Only one level of lists (no nested lists)
- Context menu items require Android or iOS 16+
- `justify` text alignment has no effect on Android (falls back to natural alignment)

### react-native-enriched-markdown

- `flavor="github"` segments text into separate TextViews, so text selection cannot span across segments
