# React Native — Branch / AppsFlyer → Detour

## Contents
- Universal / App Links — native config (iOS, Android) and Expo config
- SDK Installation
- SDK Initialization & Deep Link Handling — provider setup, Expo Router, `+native-intent`, React Navigation (linking adapter + auth-gated)
- Deferred Deep Links (First Install)
- Analytics — event mapping and code
- Keeping this reference current

## Universal / App Links

### Native config — iOS

**Before (Branch):** `applinks:yourapp.app.link`
**Before (AppsFlyer):** `applinks:yourapp.onelink.me`

**After (Detour)** — Xcode → Associated Domains:
```
applinks:YOUR_ORG.godetour.link
```

Or via Expo config:
```json
{
  "expo": {
    "ios": {
      "associatedDomains": ["applinks:YOUR_ORG.godetour.link"]
    }
  }
}
```

### Native config — Android (AndroidManifest.xml)

**After (Detour):**
```xml
<activity android:name=".MainActivity" android:exported="true">

    <!-- App Links (https) -->
    <intent-filter android:autoVerify="true">
        <action android:name="android.intent.action.VIEW" />
        <category android:name="android.intent.category.DEFAULT" />
        <category android:name="android.intent.category.BROWSABLE" />
        <data android:scheme="https" android:host="YOUR_ORG.godetour.link" />
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

Or via Expo config:
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

---

## SDK Installation

```bash
# Remove old SDK
npm uninstall react-native-branch
# or
npm uninstall react-native-appsflyer

# Add Detour
npm install @swmansion/react-native-detour
```

For bare React Native (no Expo), run `pod install` after adding the package.

---

## SDK Initialization & Deep Link Handling

### Before (Branch)
```tsx
import branch from 'react-native-branch';

useEffect(() => {
    const unsubscribe = branch.subscribe({
        onOpenComplete: ({ error, params }) => {
            if (error || !params?.['+clicked_branch_link']) return;
            const canonicalUrl = params['$canonical_url'];
            const customRoute = params['my_route_key'];
            // navigate...
        }
    });
    return () => unsubscribe();
}, []);
```

### Before (AppsFlyer)
```tsx
import appsFlyer from 'react-native-appsflyer';

useEffect(() => {
    appsFlyer.initSdk({
        devKey: 'YOUR_AF_DEV_KEY',
        isDebug: false,
        appId: 'YOUR_APPLE_APP_ID',
        onDeepLinkListener: true,
    });

    appsFlyer.onDeepLink((res) => {
        const deepLinkValue = res.deeplink?.deep_link_value;
        // navigate to deepLinkValue
    });

    appsFlyer.onInstallConversionData((res) => {
        // deferred deep link — first install only
        const deepLinkValue = res.data?.deep_link_value;
    });
}, []);
```

### After (Detour) — provider setup (same for all navigation libraries)

```tsx
// app/_layout.tsx (or App.tsx root)
import { DetourProvider } from '@swmansion/react-native-detour';

const config = {
    apiKey: 'YOUR_API_KEY',        // Detour Dashboard → API Configuration
    appID: 'YOUR_APP_ID',          // Detour Dashboard → API Configuration
    shouldUseClipboard: true,      // iOS only — improves deferred link matching
    linkProcessingMode: 'all' as const,
};

export default function RootLayout() {
    return (
        <DetourProvider config={config}>
            <YourNavigator />
        </DetourProvider>
    );
}
```

### After (Detour) — link handling: Expo Router

```tsx
import { useDetourContext } from '@swmansion/react-native-detour';
import { useRouter, usePathname } from 'expo-router';

export default function Navigator() {
    const { isLinkProcessed, link, clearLink } = useDetourContext();
    const router = useRouter();
    const pathname = usePathname();

    useEffect(() => {
        if (!isLinkProcessed || !link) return;

        // link.type:     'deferred' | 'verified' | 'scheme'
        // link.route:    "/products/123?color=red"
        // link.pathname: "/products/123"
        // link.params:   { color: "red" }

        if (pathname !== link.pathname) {
            router.replace({ pathname: link.pathname, params: link.params });
        }
        clearLink();
    }, [isLinkProcessed, link, pathname]);

    return <YourScreens />;
}
```

#### `+native-intent` (Expo Router only — recommended for short link resolution)

Expo Router supports a special `app/+native-intent.tsx` file that intercepts incoming URLs before the router handles them. This is the right place to resolve Detour short links — the router receives the already-resolved route rather than the raw short link URL.

Use this approach when:
- Your links are short links (e.g. `acme.godetour.link/abc123`) that need to be resolved to a route before navigation
- You want to keep link resolution out of your component tree

```tsx
// app/+native-intent.tsx
import { createDetourNativeIntentHandler } from '@swmansion/react-native-detour';

export const redirectSystemPath = createDetourNativeIntentHandler({
    apiKey: 'YOUR_API_KEY',
    appID: 'YOUR_APP_ID',
});
```

When using `+native-intent`, set `linkProcessingMode: 'deferred-only'` in `DetourProvider` so the SDK doesn't double-process links that were already resolved at the intent layer:

```tsx
const config = {
    apiKey: 'YOUR_API_KEY',
    appID: 'YOUR_APP_ID',
    shouldUseClipboard: true,
    linkProcessingMode: 'deferred-only' as const, // intent layer handles the rest
};
```

### After (Detour) — link handling: React Navigation

Do **not** subscribe to `useDetourContext` and call `navigation.navigate` with a hand-written
pathname-to-screen mapping. Detour plugs directly into React Navigation's
[built-in linking system](https://reactnavigation.org/docs/deep-linking?config=static#integrating-with-other-tools)
via `getInitialURL` + `subscribe`, so React Navigation parses the route and navigates for you —
the same way it would for any other deep link source.

The SDK exposes a linking adapter from the package: `Detour.getInitialURL`,
`Detour.addEventListener`, and the `DETOUR_LINKING_PREFIX` constant. Wire them into
`NavigationContainer`'s `linking` prop. `DetourProvider` **must** be mounted above
`NavigationContainer` for the adapter to receive resolved links.

```tsx
import {
    type Config,
    DETOUR_LINKING_PREFIX,
    Detour,
    DetourProvider,
} from '@swmansion/react-native-detour';
import {
    type LinkingOptions,
    NavigationContainer,
} from '@react-navigation/native';

// Map routes to screens declaratively — React Navigation handles parsing.
// This replaces any manual pathname → screen mapping.
const linkingConfig: LinkingOptions<RootStackParamList>['config'] = {
    screens: {
        Home: '',
        Details: 'details',
        NotFound: '*',
    },
};

function AppNavigator() {
    const linking: LinkingOptions<RootStackParamList> = {
        prefixes: [DETOUR_LINKING_PREFIX],
        config: linkingConfig,
        async getInitialURL() {
            return await Detour.getInitialURL();
        },
        subscribe(listener) {
            const subscription = Detour.addEventListener('url', ({ url }) => {
                listener(url);
            });
            return () => subscription.remove();
        },
    };

    return (
        <NavigationContainer linking={linking}>
            <YourScreens />
        </NavigationContainer>
    );
}

export function App() {
    return (
        <DetourProvider config={detourConfig}>
            <AppNavigator />
        </DetourProvider>
    );
}
```

`DETOUR_LINKING_PREFIX` is `"detour://"` — an internal adapter prefix for Detour-resolved routes.
React Navigation matches it against `linkingConfig.screens` and navigates automatically. There is
**no** manual `pathnameToScreenName` step.

See the [`react-navigation` example](https://github.com/software-mansion-labs/react-native-detour/blob/main/examples/react-navigation/README.md)
for a minimal working setup.

#### Auth-gated deep links (React Navigation)

For apps with sign-in / onboarding gates, let React Navigation hold the pending deep link until the
target screen becomes reachable. Render screens conditionally on auth state and opt in to React
Navigation's pending-link behavior with `UNSTABLE_routeNamesChangeBehavior="lastUnhandled"` on the
navigator. This is the React Navigation equivalent of Expo Router's `Stack.Protected`.

```tsx
<Stack.Navigator UNSTABLE_routeNamesChangeBehavior="lastUnhandled">
    {isSignedIn
        ? isOnboardingCompleted
            ? <>
                  <Stack.Screen name="Tabs" component={TabNavigator} />
                  <Stack.Screen name="Details" component={Details} />
              </>
            : <Stack.Screen name="Onboarding" component={Onboarding} />
        : <Stack.Screen name="SignIn" component={SignIn} />}
    <Stack.Screen name="NotFound" component={NotFound} />
</Stack.Navigator>
```

A deep link that arrives while the user is signed-out is parsed, found unreachable (its target
screen isn't currently rendered), and remembered. When the rendered screen set changes — after
sign-in, then again after onboarding — React Navigation retries and lands the user on the target.

> On React Navigation 7 the prop is prefixed `UNSTABLE_`. React Navigation 8 drops the prefix and
> makes it a stable `routeNamesChangeBehavior` API ([upgrade guide](https://reactnavigation.org/docs/8.x/upgrading-from-7.x/)).
> Use the name that matches your installed major version; the behavior is stable and sound.

See the [`react-navigation-advanced` example](https://github.com/software-mansion-labs/react-native-detour/blob/main/examples/react-navigation-advanced/README.md)
for a complete auth + onboarding gated setup.

---

## Deferred Deep Links (First Install)

`DetourProvider` handles first-launch detection internally — no separate API call needed. Deferred
links (a user clicked a Detour link *before* the app was installed) are delivered through the same
channel as regular links; you don't need separate handling to route them.

**Expo Router** — deferred links arrive through the `link` object from `useDetourContext`. Check
`link.type === 'deferred'`:

```tsx
useEffect(() => {
    if (!isLinkProcessed || !link) return;

    if (link.type === 'deferred') {
        // First install — user clicked a Detour link before the app was installed
    }

    router.replace({ pathname: link.pathname, params: link.params });
    clearLink();
}, [isLinkProcessed, link]);
```

**React Navigation** — deferred links flow through the **same linking adapter** as regular links
(`Detour.getInitialURL` / `Detour.addEventListener`). There is nothing extra to wire up: the
`linking` config from the section above already handles them. If you need to branch on link type,
the adapter appends a `linkType` query param (`deferred` | `verified` | `scheme`) to the resolved
route, which lands in the destination screen's params.

---

## Analytics

### Event mapping

| Branch | AppsFlyer | Detour |
|--------|-----------|--------|
| `BranchEvent.Purchase` | `"af_purchase"` | `'purchase'` |
| `BranchEvent.AddToCart` | `"af_add_to_cart"` | `'add_to_cart'` |
| `BranchEvent.RemoveFromCart` | — | `'remove_from_cart'` |
| `BranchEvent.ViewItem` | `"af_content_view"` | `'view_item'` |
| `BranchEvent.InitiatePurchase` | `"af_initiated_checkout"` | `'begin_checkout'` |
| `BranchEvent.Login` | `"af_login"` | `'login'` |
| `BranchEvent.CompleteRegistration` | `"af_complete_registration"` | `'sign_up'` |
| `BranchEvent.Search` | `"af_search"` | `'search'` |
| `BranchEvent.Share` | `"af_share"` | `'share'` |
| custom event name | custom event name | `logRetention('name')` (see note below) |

`logEvent` accepts **only** `DetourEventNames` values (the lowercase strings above, e.g. `'purchase'`,
`'view_item'`). Any event that isn't in that enum — including all of your custom Branch/AppsFlyer
events — must go through `logRetention(name)` instead. `logRetention` takes only an event name and
carries **no** properties payload, so data you previously attached to custom events (e.g. a
`placement` value) is not forwarded; keep that in mind when mapping.

### Before (Branch)
```tsx
import branch, { BranchEvent } from 'react-native-branch';

const event = new BranchEvent(BranchEvent.Purchase, buo, {
    transaction_id: 'txn_123',
    currency: 'USD',
    revenue: 29.99,
    customData: { product_id: 'abc123' },
});
event.logEvent();
```

### Before (AppsFlyer)
```tsx
appsFlyer.logEvent('af_purchase', {
    af_revenue: '29.99',
    af_currency: 'USD',
    product_id: 'abc123',
});
```

### After (Detour)
```tsx
import { DetourAnalytics, DetourEventNames } from '@swmansion/react-native-detour';

// Standard events — use the DetourEventNames enum (or its string-literal value).
DetourAnalytics.logEvent(DetourEventNames.Purchase, {
    revenue: 29.99,
    currency: 'USD',
    product_id: 'abc123',
});
```

Custom / non-standard events — anything not in `DetourEventNames` goes through `logRetention`,
which takes only an event name (no properties payload):
```tsx
DetourAnalytics.logRetention('promo_banner_tapped');
```

Retention / session events:
```tsx
DetourAnalytics.logRetention('app_open');
```

---

## Keeping this reference current

If the user asks about the latest version, installation steps, or you suspect anything above may be outdated, fetch the live README:

```
https://github.com/software-mansion-labs/react-native-detour/blob/main/README.md
```

Cross-reference against the sections above and prefer the live README if they conflict.
