# Setup and compatibility

React Native Bottom Sheet requires React Native 0.82 or later and React 19.1.1 or later. On iOS and Android, use the New Architecture (Fabric) and a native build. For Expo, check the React and React Native versions for your SDK. [Bottom Sheet requirements](https://docs.swmansion.com/react-native-bottom-sheet/docs/getting-started/#requirements) · [Expo SDK version table](https://docs.expo.dev/versions/latest/#each-expo-sdk-version-depends-on-a-react-native-version)

## Install

Install the package in the app workspace using its existing package manager.

```sh
npm install @swmansion/react-native-bottom-sheet
```

- **Expo on iOS or Android.** Expo Go does not include the library’s native code. Use a development build. For a local build, run `npx expo run:ios` or `npx expo run:android`. Rebuild an existing development client after installing the library. [Bottom Sheet requirements](https://docs.swmansion.com/react-native-bottom-sheet/docs/getting-started/#requirements) · [Expo development builds](https://docs.expo.dev/develop/development-builds/introduction/#rebuild-when-native-code-changes)
- **Bare React Native.** Rebuild the native app after installation. On iOS, run `cd ios && pod install` first. [Bottom Sheet installation](https://docs.swmansion.com/react-native-bottom-sheet/docs/getting-started/#installation)
- **Web.** Install `react-dom` and `react-native-web`. In Expo, run `npx expo install react-dom react-native-web`. If you're setting up Expo web support, follow Expo's setup for `@expo/metro-runtime` as well. [Bottom Sheet web setup](https://docs.swmansion.com/react-native-bottom-sheet/docs/getting-started/#web) · [Expo web dependencies](https://docs.expo.dev/workflow/web/#install-web-dependencies)

The library can run without Reanimated or React Native Gesture Handler. [How it compares](https://docs.swmansion.com/react-native-bottom-sheet/docs/how-it-compares/)

## Provider setup

Wrap the app tree with `BottomSheetProvider` for the default `ModalBottomSheet` portal. Inline `BottomSheet` and `ModalBottomSheet` with `windowOverlay` can render without the provider. Place the context providers used by portal content above `BottomSheetProvider`. [Bottom Sheet installation](https://docs.swmansion.com/react-native-bottom-sheet/docs/getting-started/#installation) · [Modal sheet context](https://docs.swmansion.com/react-native-bottom-sheet/docs/modal-sheets/#react-context)
