# Detour — React Native SDK Reference

## Contents
- Installation
- Universal / App Links
- SDK Setup
- Expo Router Integration
- Analytics
- Keeping this reference current

## Installation

```bash
npm install @swmansion/react-native-detour
# or: yarn add / pnpm add / bun add @swmansion/react-native-detour
```

### Peer dependencies

**Expo (managed or bare with expo-modules):**

```bash
npx expo install expo-application expo-clipboard expo-constants expo-device expo-localization @react-native-async-storage/async-storage
```

**Bare React Native (no expo-modules-core):**
Install peer deps manually:

```bash
npm install expo-application expo-clipboard expo-constants expo-device expo-localization @react-native-async-storage/async-storage
```

**async-storage version note:** The official peer dep is `>=2.0.0`. Version 3.x technically satisfies this but uses KMP and requires an Android artifact (`org.asyncstorage.shared_storage:storage-android`) that may not be available on Maven Central — if you hit build errors with 3.x, pin to `^2.x`:

```bash
npm install "@react-native-async-storage/async-storage@^2"
```

### babel-preset warning (bare React Native)

Do not add `babel-preset-expo` to a bare React Native project — it causes Metro errors like `Unable to determine event arguments for "onModeChange"`. Bare RN should use `@react-native/babel-preset` only.

---

## Universal / App Links

### Expo (app.json)

**iOS:**

```json
{
  "expo": {
    "ios": {
      "associatedDomains": ["applinks:YOUR_ORG.godetour.link"]
    }
  }
}
```

**Android:**

```json
{
  "expo": {
    "android": {
      "intentFilters": [
        {
          "action": "VIEW",
          "autoVerify": true,
          "data": [{ "scheme": "https", "host": "YOUR_ORG.godetour.link" }],
          "category": ["BROWSABLE", "DEFAULT"]
        }
      ]
    }
  }
}
```

### Bare React Native — iOS

`<AppTarget>.entitlements`:

```xml
<key>com.apple.developer.associated-domains</key>
<array>
  <string>applinks:YOUR_ORG.godetour.link</string>
</array>
```

`AppDelegate.mm`:

```objc
- (BOOL)application:(UIApplication *)application
    continueUserActivity:(NSUserActivity *)userActivity
    restorationHandler:(void (^)(NSArray<id<UIUserActivityRestoring>> *))restorationHandler {
  return [RCTLinkingManager application:application
                   continueUserActivity:userActivity
                     restorationHandler:restorationHandler];
}
```

### Bare React Native — Android

Copy the **exact snippet from the Detour Dashboard** (App Configuration tab) — it includes the unique `pathPrefix` for your app. Never write the intent-filter by hand. It will look like:

```xml
<activity android:name=".MainActivity" android:exported="true">

  <!-- App Links — COPY FROM DETOUR DASHBOARD, do not write manually -->
  <intent-filter android:autoVerify="true">
    <action android:name="android.intent.action.VIEW" />
    <category android:name="android.intent.category.DEFAULT" />
    <category android:name="android.intent.category.BROWSABLE" />
    <data android:scheme="https"
          android:host="YOUR_ORG.godetour.link"
          android:pathPrefix="/YOUR_APP_HASH" />
  </intent-filter>

  <!-- Optional: custom URI scheme -->
  <intent-filter>
    <action android:name="android.intent.action.VIEW" />
    <category android:name="android.intent.category.DEFAULT" />
    <category android:name="android.intent.category.BROWSABLE" />
    <data android:scheme="yourapp" />
  </intent-filter>

</activity>
```

---

## SDK Setup

### 1. Config

```typescript
// detour.config.ts
import { type Config } from "@swmansion/react-native-detour";

export const detourConfig: Config = {
  apiKey: process.env.EXPO_PUBLIC_DETOUR_API_KEY!,
  appID: process.env.EXPO_PUBLIC_DETOUR_APP_ID!,
  shouldUseClipboard: true,
  linkProcessingMode: "all" as const,
};
```

`.env`:

```
EXPO_PUBLIC_DETOUR_API_KEY=your_publishable_key
EXPO_PUBLIC_DETOUR_APP_ID=your_app_id
```

### 2. Wrap root with DetourProvider

```tsx
// app/_layout.tsx (or App.tsx)
import { DetourProvider } from "@swmansion/react-native-detour";
import { detourConfig } from "./detour.config";

export default function RootLayout() {
  return (
    <DetourProvider config={detourConfig}>
      <YourNavigator />
    </DetourProvider>
  );
}
```

One provider per app tree. On Expo Web it becomes a no-op automatically.

### 3. Handle resolved links

```tsx
import { useDetourContext } from "@swmansion/react-native-detour";

export default function Navigator() {
  const { isLinkProcessed, link, clearLink } = useDetourContext();
  const router = useRouter();
  const pathname = usePathname();

  useEffect(() => {
    if (!isLinkProcessed || !link) return;

    // link.type:     'deferred' | 'verified' | 'scheme'
    // link.route:    "/articles/2?ref=campaign"
    // link.pathname: "/articles/2"
    // link.params:   { ref: "campaign" }

    if (pathname !== link.pathname) {
      router.replace({ pathname: link.pathname, params: link.params });
    }
    clearLink();
  }, [isLinkProcessed, link, pathname]);

  return <YourScreens />;
}
```

Wait for `isLinkProcessed` before navigating — it turns `true` once the SDK has finished its initial link resolution (whether or not a link was found).

---

## Expo Router Integration (recommended if someone already has expo router setup)

Create `app/+native-intent.tsx` to resolve links before the first render:

```typescript
import { createDetourNativeIntentHandler } from "@swmansion/react-native-detour/expo-router";

export const redirectSystemPath = createDetourNativeIntentHandler({
  config: {
    apiKey: process.env.EXPO_PUBLIC_DETOUR_API_KEY!,
    appID: process.env.EXPO_PUBLIC_DETOUR_APP_ID!,
  },
});
```

When using this pattern, limit the provider to deferred-only:

```typescript
export const detourConfig: Config = {
  apiKey: process.env.EXPO_PUBLIC_DETOUR_API_KEY!,
  appID: process.env.EXPO_PUBLIC_DETOUR_APP_ID!,
  shouldUseClipboard: true,
  linkProcessingMode: "deferred-only" as const,
};
```

---

## Analytics

`logEvent` only accepts values from the `DetourEventNames` enum for standard events. Do not pass arbitrary strings like `"article_opened"` — they are not valid event names and will not be tracked correctly.

```typescript
import {
  DetourAnalytics,
  DetourEventNames,
} from "@swmansion/react-native-detour";

// Standard event — use enum
DetourAnalytics.logEvent(DetourEventNames.Purchase, {
  revenue: 29.99,
  currency: "USD",
});
DetourAnalytics.logEvent(DetourEventNames.ViewItem, { item_id: "article_2" });

// Custom / non-standard activity — use logRetention with a descriptive string
DetourAnalytics.logRetention("article_opened");
DetourAnalytics.logRetention("app_open");
```

---

## Keeping this reference current

If the user asks about the latest version, installation steps, or you suspect anything above may be outdated, fetch the live README:

```
https://github.com/software-mansion-labs/react-native-detour/blob/main/README.md
```

Cross-reference against the sections above and prefer the live README if they conflict.
