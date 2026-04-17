---
name: expo-horizon
description: "Software Mansion's guide for migrating Expo SDK apps to Meta Quest using expo-horizon packages. Use when adding Meta Quest or Meta Horizon OS support to an existing Expo or React Native project. Trigger on: Meta Quest, Horizon OS, Quest 2, Quest 3, Quest 3S, VR app, expo-horizon-core, expo-horizon-location, expo-horizon-notifications, build flavors for Quest, panel sizing, VR headtracking, Horizon App ID, quest build variant, isHorizonDevice, isHorizonBuild, migrate expo-location to Quest, migrate expo-notifications to Quest, Meta Horizon Store publishing, or any task involving running an Expo app on Meta Quest hardware."
---

# Expo Horizon: Migrating Expo SDK to Meta Quest

Software Mansion's production guide for adding Meta Quest support to Expo apps using the [expo-horizon](https://github.com/software-mansion-labs/expo-horizon) packages.

**This skill does not bundle a copy of the docs.** For any task below, always webfetch the linked official README or Meta documentation page to get up-to-date installation steps, plugin options, API surface, and feature matrices. This skill only captures the decision tree, critical rules, and non-obvious gotchas that agents routinely miss.

## Decision Tree

```
What do you need to do?
│
├── Starting from scratch or adding Quest support to an existing Expo app?
│   └── Webfetch: expo-horizon-core README
│       ├── Install expo-horizon-core
│       ├── Configure the config plugin (horizonAppId, panel size, supportedDevices)
│       ├── Add quest/mobile build scripts
│       └── Add runtime device detection (isHorizonDevice, isHorizonBuild)
│
├── Need location services on Quest?
│   └── Webfetch: expo-horizon-location README
│       ├── Replace expo-location with expo-horizon-location
│       ├── Review the feature support matrix
│       └── Guard unsupported calls (heading, geocoding, geofencing, background)
│
├── Need push notifications on Quest?
│   └── Webfetch: expo-horizon-notifications README
│       ├── Replace expo-notifications with expo-horizon-notifications
│       ├── Configure horizonAppId in expo-horizon-core
│       ├── Use getDevicePushTokenAsync (Expo Push Service is not supported)
│       └── Skip badge counts (not supported on Quest)
│
└── Need to build, run, or publish for Quest?
    └── Webfetch: expo-horizon-core README (build variants) + Meta docs below
        ├── Build variants: questDebug, questRelease, mobileDebug, mobileRelease
        ├── Meta Quest Developer Hub (device management, sideloading)
        └── Meta Horizon Store manifest requirements
```

## Critical Rules

- **Always install `expo-horizon-core` first.** It is required by all other expo-horizon packages and sets up the `quest`/`mobile` build flavors that other packages depend on.

- **Use `quest` build variants only on Meta Quest devices.** Running `questDebug` or `questRelease` builds on standard Android phones is unsupported and will behave unexpectedly.

- **Set `supportedDevices` in the config plugin.** This is required for Meta Horizon Store submission. Use pipe-separated values: `"quest2|quest3|quest3s"`.

- **Run `npx expo prebuild --clean` after any plugin config change.** The config plugin modifies native project files at prebuild time. Stale native projects will not reflect your changes.

- **Replace imports, not just packages.** When migrating from `expo-location` or `expo-notifications`, update all import statements to use the new package names (`expo-horizon-location`, `expo-horizon-notifications`).

- **Quest has no GPS, magnetic sensors, or Geocoder.** Features like heading, geocoding, reverse geocoding, and geofencing are unavailable on Quest. Guard these calls with `ExpoHorizon.isHorizonDevice` or `ExpoHorizon.isHorizonBuild`.

- **Push notifications require `horizonAppId`.** Without it, `getDevicePushTokenAsync` will not return a valid token on Quest devices. Use `getDevicePushTokenAsync` (not `getExpoPushTokenAsync`) on Quest; send the returned `{ type: 'horizon', data }` token to your backend and deliver via Meta's push service.

- **`isHorizonDevice` vs `isHorizonBuild`.** Use `isHorizonDevice` for runtime hardware checks (physical Quest detection). Use `isHorizonBuild` for build-time feature gating (which native code was compiled in).

- **Expo Go is not supported.** You must use custom development builds via `npx expo prebuild`.

## Official References

Always webfetch the raw markdown (`raw.githubusercontent.com/...`) if the HTML view does not render the source; the raw URL is the source of truth.

| Topic | Official source |
|-------|-----------------|
| Repo overview and package list | [expo-horizon README](https://github.com/software-mansion-labs/expo-horizon/blob/main/README.md) |
| Install, config plugin options, runtime API, native module access | [expo-horizon-core README](https://github.com/software-mansion-labs/expo-horizon/blob/main/expo-horizon-core/README.md) |
| Location migration, limitations, feature support matrix | [expo-horizon-location README](https://github.com/software-mansion-labs/expo-horizon/blob/main/expo-horizon-location/README.md) |
| Push notifications migration, token types, feature support matrix | [expo-horizon-notifications README](https://github.com/software-mansion-labs/expo-horizon/blob/main/expo-horizon-notifications/README.md) |
| Example app wiring for all three packages | [expo-horizon example README](https://github.com/software-mansion-labs/expo-horizon/blob/main/example/README.md) |
| Panel sizing guidelines (dp values, orientation, letterboxing) | [Meta Panel Sizing](https://developers.meta.com/horizon/documentation/android-apps/panel-sizing) |
| Meta Horizon Store manifest checklist for publishing | [Publish Mobile Manifest](https://developers.meta.com/horizon/resources/publish-mobile-manifest/) |
| Device management, casting, sideloading, ADB | [Meta Quest Developer Hub](https://developers.meta.com/horizon/documentation/android-apps/meta-quest-developer-hub) |
| Server-side push delivery via Meta's push service | [Horizon OS push notifications](https://developers.meta.com/horizon/documentation/android-apps/ps-user-notifications/) |
