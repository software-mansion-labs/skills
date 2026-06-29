# Detour — iOS SDK Reference

## Contents
- Requirements
- Installation
- Universal Links
- Custom URL Scheme (optional)
- SDK Setup
- Testing deferred links on iOS
- Analytics
- Keeping this reference current

## Requirements
- iOS 13+, Swift 5.5+

## Installation

**SPM (Xcode):** `File > Add Package Dependencies` then enter:
```
https://github.com/software-mansion-labs/ios-detour
```

**Package.swift:**
```swift
// Pin the current release — check the GitHub releases page / README for the latest version.
.package(url: "https://github.com/software-mansion-labs/ios-detour", from: "LATEST_VERSION")
```

**CocoaPods:**
```ruby
platform :ios, '13.0'
target 'YourAppTarget' do
  use_frameworks!
  pod 'Detour'  # latest; add a version constraint to pin — see the README
end
```
Then run `pod install`.

---

## Universal Links

### 1. Add Associated Domains

**Xcode → Signing & Capabilities → Associated Domains:**
```
applinks:YOUR_ORG.godetour.link
```

> **App Store Connect:** The Associated Domains entitlement must also be enabled in your App ID on App Store Connect (Certificates, Identifiers & Profiles → your App ID → Capabilities → Associated Domains). Without this, the entitlement won't be included in the provisioning profile and Universal Links won't work — even if the entitlements file looks correct locally.

Or in `<AppTarget>.entitlements`:
```xml
<key>com.apple.developer.associated-domains</key>
<array>
  <string>applinks:YOUR_ORG.godetour.link</string>
</array>
```

### 2. Handle in AppDelegate

```swift
// Universal Link (https://YOUR_ORG.godetour.link/...)
func application(
  _ application: UIApplication,
  continue userActivity: NSUserActivity,
  restorationHandler: @escaping ([UIUserActivityRestoring]?) -> Void
) -> Bool {
  guard userActivity.activityType == NSUserActivityTypeBrowsingWeb,
        let url = userActivity.webpageURL else { return false }
  Task { @MainActor in
    let result = await Detour.shared.processLink(url, config: detourConfig)
    handleLink(result)
  }
  return true
}
```

---

## Custom URL Scheme (optional)

`Info.plist`:
```xml
<key>CFBundleURLTypes</key>
<array>
  <dict>
    <key>CFBundleURLSchemes</key>
    <array><string>yourapp</string></array>
  </dict>
</array>
```

`AppDelegate`:
```swift
func application(_ app: UIApplication, open url: URL,
  options: [UIApplication.OpenURLOptionsKey: Any] = [:]) -> Bool {
  Task { @MainActor in
    let result = await Detour.shared.processLink(url, config: detourConfig)
    handleLink(result)
  }
  return true
}
```

---

## SDK Setup

### 1. Config

```swift
import Detour

let detourConfig = DetourConfig(
  apiKey: "YOUR_API_KEY",    // Detour Dashboard → API Configuration
  appID: "YOUR_APP_ID",      // Detour Dashboard → API Configuration
  shouldUseClipboard: true   // improves deferred link matching on iOS
)
```

### 2. Cold start — resolveInitialLink

Call in `application(_:didFinishLaunchingWithOptions:)`:

```swift
func application(
  _ application: UIApplication,
  didFinishLaunchingWithOptions launchOptions: [UIApplication.LaunchOptionsKey: Any]?
) -> Bool {
  // resolveInitialLink also mounts analytics automatically
  Detour.shared.resolveInitialLink(config: detourConfig, launchOptions: launchOptions) { result in
    self.handleLink(result)
  }
  return true
}

private func handleLink(_ result: DetourResult) {
  guard let link = result.link else { return }

  // link.type:     .deferred | .verified | .scheme
  // link.route:    "/articles/2?ref=campaign"
  // link.pathname: "/articles/2"
  // link.params:   ["ref": "campaign"]

  NotificationCenter.default.post(name: .detourLinkReceived, object: link)
}
```

**SwiftUI / SceneDelegate:**
```swift
func scene(_ scene: UIScene,
           willConnectTo session: UISceneSession,
           options connectionOptions: UIScene.ConnectionOptions) {
  Detour.shared.resolveInitialLink(
    config: detourConfig,
    connectionOptions: connectionOptions
  ) { result in
    self.handleLink(result)
  }
}
```

Note: `resolveInitialLink` mounts analytics automatically — no need to call `mountAnalytics` separately.

---

## Testing deferred links on iOS

1. Uninstall the app before clicking the Detour link
2. Click `https://YOUR_ORG.godetour.link/APP_HASH/your/path` in Safari on a real device
3. Install the app, open it — the SDK resolves the link
4. Each test needs a fresh install — the same deferred link won't fire twice on the same device

For testing, reset the first-launch flag without reinstalling:
```swift
Detour.shared.resetSession(allowDeferredRetry: true)
```

---

## Analytics

```swift
// Standard events — use DetourEventName enum
DetourAnalytics.logEvent(.purchase, data: ["revenue": 29.99, "currency": "USD"])
DetourAnalytics.logEvent(.viewItem, data: ["item_id": "article_2"])

// Custom / non-standard activity — use logRetention
DetourAnalytics.logRetention("article_opened")
DetourAnalytics.logRetention("app_open")
```

---

## Keeping this reference current

If the user asks about the latest version, installation steps, or you suspect anything above may be outdated, fetch the live README:

```
https://github.com/software-mansion-labs/ios-detour/blob/main/README.md
```

Cross-reference against the sections above and prefer the live README if they conflict.
