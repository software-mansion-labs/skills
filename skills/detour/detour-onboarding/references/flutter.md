# Detour — Flutter Plugin Reference

## Contents
- Requirements
- Installation
- Universal / App Links
- SDK Setup
- Analytics
- Keeping this reference current

## Requirements
- Flutter 3.3.0+, Dart ^3.11.1
- Android minSdk 24, iOS 13+

## Installation

`pubspec.yaml`:
```yaml
dependencies:
  # Use the latest version — check pub.dev / the README (linked below).
  detour_flutter_plugin: ^LATEST_VERSION
```

```bash
flutter pub get
cd ios && pod install && cd ..
```

---

## Universal / App Links

### iOS

`ios/Runner/Runner.entitlements`:
```xml
<key>com.apple.developer.associated-domains</key>
<array>
  <string>applinks:YOUR_ORG.godetour.link</string>
</array>
```

`ios/Runner/Info.plist` (if using custom scheme):
```xml
<key>CFBundleURLTypes</key>
<array>
  <dict>
    <key>CFBundleURLSchemes</key>
    <array><string>yourapp</string></array>
  </dict>
</array>
```

If overriding AppDelegate/SceneDelegate callbacks, **always call `super`**.

### Android

`android/app/src/main/AndroidManifest.xml` inside main `<activity>`:
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

## SDK Setup

`DetourService` extends `ChangeNotifier` — use `addListener` to react to changes.

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

`DetourIntentSource.initial` — link resolved on app start (deferred or cold-start Universal Link).
`DetourIntentSource.runtime` — link received while app was already running.

---

## Analytics

```dart
// Standard event (use DetourEventName enum)
await _detour.logEvent(DetourEventName.purchase, data: {
  'revenue': 29.99,
  'currency': 'USD',
  'product_id': 'abc123',
});

// Custom / non-standard event — logEvent accepts only DetourEventName enum values
// (a custom string is a compile error). Use logRetention for custom names.
// Note: logRetention takes only a name — no data payload.
await _detour.logRetention('promo_banner_tapped');

// Retention
await _detour.logRetention('app_open');
```

---

## Keeping this reference current

If the user asks about the latest version, installation steps, or you suspect anything above may be outdated, fetch the live README:

```
https://github.com/software-mansion-labs/detour-flutter-plugin/blob/main/README.md
```

Cross-reference against the sections above and prefer the live README if they conflict.
