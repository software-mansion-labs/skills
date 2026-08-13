# Detour Skills

Skills for [Detour](https://detour.swmansion.com) — the open-source deferred deep linking SDK by [Software Mansion](https://swmansion.com/).

Add these skills to give your AI coding agent accurate, up-to-date guidance for setting up and migrating to Detour across iOS, Android, React Native, and Flutter.

## Skills

### [detour-onboarding](./detour-onboarding/)

End-to-end onboarding for developers new to Detour. Covers:

| Phase | What it covers |
|-------|----------------|
| **Dashboard setup** | Account creation, organization and app configuration, generating App ID, Publishable API Key, and App Hash |
| **Universal / App Links** | Registering the Detour domain for direct-open links on iOS and Android |
| **SDK integration** | Initializing the SDK, handling the link callback, `linkProcessingMode`, navigating to the right screen |
| **Analytics** | Automatic click/install tracking, `logEvent` with `DetourEventNames`, `logRetention` |
| **Architecture** | Deterministic (Play Store referrer) and probabilistic (fingerprint scoring) matching, known limitations |
| **Testing** | Deferred link test flow, Universal/App Links testing, Android 12+ adb setup |

### [migrate-to-detour](./migrate-to-detour/)

Structural migration guide for teams switching from Branch or AppsFlyer to Detour. Covers:

| Topic | What it covers |
|-------|----------------|
| **Concept mapping** | How Branch / AppsFlyer concepts (deep link handlers, deferred links, event tracking) map to Detour equivalents |
| **Dashboard migration** | Recreating link campaigns, fallback URLs, and app configuration in the Detour dashboard |
| **Universal / App Links** | Swapping the verification domain from Branch/AppsFlyer to Detour |
| **SDK swap** | Replacing initialization, link handling callbacks, and analytics calls with Detour equivalents |
| **Analytics migration** | Mapping Branch/AppsFlyer event names to `DetourEventNames` and `logRetention` |

## Platforms

Both skills provide platform-specific reference files for:

- **React Native** (Expo and bare workflow)
- **iOS** (native Swift)
- **Android** (native Kotlin / Java)
- **Flutter**

## Structure

```
detour/
├── detour-onboarding/
│   ├── SKILL.md
│   └── references/
│       ├── android.md
│       ├── flutter.md
│       ├── ios.md
│       └── react-native.md
└── migrate-to-detour/
    ├── SKILL.md
    └── references/
        ├── android.md
        ├── flutter.md
        ├── ios.md
        └── react-native.md
```

## Links

- [Detour documentation](https://detour.swmansion.com/docs)
- [Detour dashboard](https://app.godetour.dev)
- [Software Mansion](https://swmansion.com)

---

> Branch and AppsFlyer are trademarks of their respective owners; used for identification and compatibility purposes only; no affiliation or endorsement.
