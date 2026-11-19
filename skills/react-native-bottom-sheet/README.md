# React Native Bottom Sheet Skills

Integration and migration guidance for [React Native Bottom Sheet](https://docs.swmansion.com/react-native-bottom-sheet/), Software Mansion's bottom-sheet library for iOS, Android, and React Native Web.

## Skills

### [use-react-native-bottom-sheet](./use-react-native-bottom-sheet/)

Install `@swmansion/react-native-bottom-sheet` and build inline or modal sheets. Covers Expo and bare React Native setup, provider placement, controlled detents, surfaces, scroll coordination, keyboard behavior, and dismissal.

### [migrate-gorhom-to-react-native-bottom-sheet](./migrate-gorhom-to-react-native-bottom-sheet/)

Assess or migrate existing `@gorhom/bottom-sheet` usage to `@swmansion/react-native-bottom-sheet`. Covers compatibility, index and command mapping, modal lifecycle, interaction changes, and dependency cleanup. This skill handles only the Gorhom-to-Software-Mansion direction.

### [migrate-true-sheet-to-react-native-bottom-sheet](./migrate-true-sheet-to-react-native-bottom-sheet/)

Assess or migrate `@lodev09/react-native-true-sheet` usage. Covers version-specific detents, asynchronous commands, lifecycle and stacking, navigation ownership, scrolling, keyboard behavior, and native presentation gaps.

### [migrate-expo-ui-to-react-native-bottom-sheet](./migrate-expo-ui-to-react-native-bottom-sheet/)

Assess or migrate Expo UI sheets. Identifies the SDK and universal, SwiftUI, Compose, or community import before mapping state, sizing, dismissal, native content hosts, and appearance.

### [migrate-native-form-sheets-to-react-native-bottom-sheet](./migrate-native-form-sheets-to-react-native-bottom-sheet/)

Assess or migrate native-stack `presentation: 'formSheet'` through React Navigation or Expo Router. Preserves route ownership, links, params, and history while mapping sheet visuals, detents, interactions, and coordinated dismissal.

Each migration skill handles its named source into `@swmansion/react-native-bottom-sheet`. A source-library fix or upgrade does not imply replacement.

[Getting started](https://docs.swmansion.com/react-native-bottom-sheet/docs/getting-started/) · [Product changelog](https://swmansion.com/changelog/react-native-bottom-sheet/)
