# Flutter — AppsFlyer → Detour

> Branch does not have an official Flutter SDK, so this file covers AppsFlyer → Detour only.

## Contents
- Universal / App Links
- SDK Installation
- SDK Initialization & Deep Link Handling
- Deferred Deep Links (First Install)
- Analytics
- Keeping this reference current

## Universal / App Links

### Native config — iOS

**Before (AppsFlyer):** `applinks:yourapp.onelink.me` in Associated Domains

**After (Detour)** — `ios/Runner/Runner.entitlements`:
```xml
<key>com.apple.developer.associated-domains</key>
<array>
    <string>applinks:YOUR_ORG.godetour.link</string>
</array>
```

Custom URI scheme in `ios/Runner/Info.plist` (optional):
```xml
<key>CFBundleURLTypes</key>
<array>
    <dict>
        <key>CFBundleURLSchemes</key>
        <array>
            <string>yourapp</string>
        </array>
    </dict>
</array>
```

### Native config — Android

**After (Detour)** — `android/app/src/main/AndroidManifest.xml`:
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

---

## SDK Installation

```yaml
# pubspec.yaml
dependencies:
  # Remove:
  # appsflyer_sdk: ^x.x.x

  # Add (use the latest version — check pub.dev / the README):
  detour_flutter_plugin: ^LATEST_VERSION
```

```bash
flutter pub get
```

---

## SDK Initialization & Deep Link Handling

### Before (AppsFlyer)
```dart
import 'package:appsflyer_sdk/appsflyer_sdk.dart';

class _MyAppState extends State<MyApp> {
    late AppsflyerSdk _appsflyerSdk;

    @override
    void initState() {
        super.initState();

        final options = AppsFlyerOptions(
            afDevKey: 'YOUR_AF_DEV_KEY',
            appId: '123456789',
            showDebug: false,
        );
        _appsflyerSdk = AppsflyerSdk(options);

        // Deferred deep link — first install
        _appsflyerSdk.onInstallConversionData((res) {
            final deepLinkValue = res['payload']?['deep_link_value'];
        });

        // Direct deep link — app already installed
        _appsflyerSdk.onDeepLinking((res) {
            final deepLinkValue = res.deepLink?['deep_link_value'];
        });

        _appsflyerSdk.initSdk(
            registerConversionDataCallback: true,
            registerOnDeepLinkingCallback: true,
        );
    }
}
```

### After (Detour)
```dart
import 'package:detour_flutter_plugin/detour_flutter_plugin.dart';

class _MyAppState extends State<MyApp> {
    final _detour = DetourService();

    @override
    void initState() {
        super.initState();
        _detour.addListener(_onDetourChanged);
        _startDetour();
    }

    Future<void> _startDetour() async {
        await _detour.start(
            const DetourConfig(
                apiKey: 'YOUR_API_KEY',        // Detour Dashboard → API Configuration
                appID: 'YOUR_APP_ID',          // Detour Dashboard → API Configuration
                shouldUseClipboard: true,      // iOS only — improves deferred link matching
                linkProcessingMode: LinkProcessingMode.all,
            ),
        );
    }

    void _onDetourChanged() {
        final intent = _detour.pendingIntent;
        if (intent == null) return;

        // intent.link.type:     LinkType.deferred | verified | scheme
        // intent.link.route:    "/products/123?color=red"
        // intent.link.pathname: "/products/123"
        // intent.link.params:   {"color": "red"}
        // intent.source:        DetourIntentSource.initial | runtime | manual

        context.go(intent.link.route);
        _detour.consumePendingIntent();
    }

    @override
    void dispose() {
        _detour.removeListener(_onDetourChanged);
        _detour.dispose();
        super.dispose();
    }
}
```

---

## Deferred Deep Links (First Install)

AppsFlyer had separate `onInstallConversionData` (deferred) and `onDeepLinking` (direct) callbacks.

With Detour, both come through `_onDetourChanged`. Check `intent.link.type`:

```dart
void _onDetourChanged() {
    final intent = _detour.pendingIntent;
    if (intent == null) return;

    if (intent.link.type == LinkType.deferred) {
        // First install — user clicked a Detour link before installing the app
    }

    context.go(intent.link.route);
    _detour.consumePendingIntent();
}
```

`intent.source == DetourIntentSource.initial` identifies the link resolved on app start.
`intent.source == DetourIntentSource.runtime` identifies links received while the app was already running.

---

## Analytics

### Event mapping

| AppsFlyer | Detour (`DetourEventName`) |
|-----------|----------------------------|
| `AFEvent.purchase` / `"af_purchase"` | `DetourEventName.purchase` |
| `AFEvent.addToCart` / `"af_add_to_cart"` | `DetourEventName.addToCart` |
| `AFEvent.contentView` / `"af_content_view"` | `DetourEventName.viewItem` |
| `AFEvent.initiatedCheckout` | `DetourEventName.beginCheckout` |
| `AFEvent.addPaymentInfo` | `DetourEventName.addPaymentInfo` |
| `AFEvent.login` / `"af_login"` | `DetourEventName.login` |
| `AFEvent.completeRegistration` | `DetourEventName.signUp` |
| `AFEvent.search` / `"af_search"` | `DetourEventName.search` |
| `AFEvent.share` | `DetourEventName.share` |
| custom event name | `logRetention('name')` (no enum entry — see below) |

### Before (AppsFlyer)
```dart
await _appsflyerSdk.logEvent('af_purchase', {
    'af_revenue': '29.99',
    'af_currency': 'USD',
    'product_id': 'abc123',
});
```

### After (Detour)
```dart
await _detour.logEvent(
    DetourEventName.purchase,
    data: {
        'revenue': 29.99,
        'currency': 'USD',
        'product_id': 'abc123',
    },
);
```

Custom / non-standard events — `logEvent` accepts only `DetourEventName` enum values (a custom
string is a compile error); custom names must use `logRetention`, which takes only a name
(no `data` payload):
```dart
await _detour.logRetention('promo_banner_tapped');
```

Retention / session events:
```dart
await _detour.logRetention('app_open');
```

---

## Keeping this reference current

If the user asks about the latest version, installation steps, or you suspect anything above may be outdated, fetch the live README:

```
https://github.com/software-mansion-labs/detour-flutter-plugin/blob/main/README.md
```

Cross-reference against the sections above and prefer the live README if they conflict.
