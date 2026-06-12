# React Native — Branch / AppsFlyer → Detour

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

```tsx
import { useDetourContext } from '@swmansion/react-native-detour';
import { useNavigation } from '@react-navigation/native';

export function RootNavigator() {
    const { isLinkProcessed, link, clearLink } = useDetourContext();
    const navigation = useNavigation();

    useEffect(() => {
        if (!isLinkProcessed || !link) return;

        // link.pathname maps to your screen — adjust to your navigator structure
        // e.g. "/products/123" → navigate to "Product" screen with { id: "123" }
        navigation.navigate(
            pathnameToScreenName(link.pathname),
            link.params
        );
        clearLink();
    }, [isLinkProcessed, link]);

    return <YourScreens />;
}
```

React Navigation doesn't use file-based routing, so you need to map `link.pathname` to your screen names manually. Keep that mapping close to your navigator definition. See the [`react-navigation` example](https://github.com/software-mansion-labs/react-native-detour/blob/main/examples/react-navigation/README.md) for a complete reference.

---

## Deferred Deep Links (First Install)

With Detour, deferred links arrive through the same `link` object from `useDetourContext`. Check `link.type === 'deferred'`:

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

`DetourProvider` handles first-launch detection internally — no separate API call needed.

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
| custom string | custom string | custom string |

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
import { DetourAnalytics } from '@swmansion/react-native-detour';

DetourAnalytics.logEvent('purchase', {
    revenue: 29.99,
    currency: 'USD',
    product_id: 'abc123',
});
```

Custom events:
```tsx
DetourAnalytics.logEvent('promo_banner_tapped', { placement: 'home_top' });
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
