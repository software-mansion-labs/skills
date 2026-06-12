# iOS — Branch / AppsFlyer → Detour

## Universal Links

### Before (Branch)
Associated Domains entitlement:
```
applinks:yourapp.app.link
applinks:yourapp-alternate.app.link
```

### Before (AppsFlyer)
```
applinks:yourapp.onelink.me
```

### After (Detour)
Replace with your Detour organization subdomain. Detour auto-hosts `/.well-known/apple-app-site-association` — no server setup needed.

**Step 1 — Xcode → Signing & Capabilities → Associated Domains:**
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

**Step 2 — App Store Connect → Associated Domains capability:**

Universal Links also require the capability to be enabled in App Store Connect, not just in Xcode. Without this step, links will work in local debug builds but fail on TestFlight and App Store.

1. Go to [App Store Connect](https://appstoreconnect.apple.com) → **Certificates, Identifiers & Profiles**
2. Select your App ID → **Edit**
3. Under **Capabilities**, enable **Associated Domains**
4. Save — this triggers a new provisioning profile. Re-download and install it in Xcode.

> If you use Xcode's automatic signing, this step may happen automatically when you add the entitlement. Verify the capability is listed as **Enabled** in App Store Connect regardless.

**Info.plist — Custom URI Scheme (optional):**
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

---

## SDK Installation

**Swift Package Manager** — Xcode → File → Add Package Dependencies:
```
https://github.com/software-mansion-labs/ios-detour
```

Or in `Package.swift`:
```swift
.package(url: "https://github.com/software-mansion-labs/ios-detour", from: "1.0.2")
```

Remove the Branch or AppsFlyer package.

---

## SDK Initialization & Deep Link Handling

### Before (Branch)
```swift
// AppDelegate.swift
import Branch

func application(_ application: UIApplication,
                 didFinishLaunchingWithOptions launchOptions: [UIApplication.LaunchOptionsKey: Any]?) -> Bool {
    Branch.getInstance().initSession(launchOptions: launchOptions) { params, error in
        guard error == nil, let params = params as? [String: AnyObject] else { return }
        let canonicalUrl = params["$canonical_url"] as? String
        // navigate to canonicalUrl
    }
    return true
}

func application(_ app: UIApplication, open url: URL,
                 options: [UIApplication.OpenURLOptionsKey: Any] = [:]) -> Bool {
    Branch.getInstance().application(app, open: url, options: options)
    return true
}

func application(_ application: UIApplication,
                 continue userActivity: NSUserActivity,
                 restorationHandler: @escaping ([UIUserActivityRestoring]?) -> Void) -> Bool {
    Branch.getInstance().continue(userActivity)
    return true
}
```

### Before (AppsFlyer)
```swift
// AppDelegate.swift
import AppsFlyerLib

class AppDelegate: UIResponder, UIApplicationDelegate, AppsFlyerLibDelegate, DeepLinkDelegate {

    func application(_ application: UIApplication,
                     didFinishLaunchingWithOptions launchOptions: [UIApplication.LaunchOptionsKey: Any]?) -> Bool {
        AppsFlyerLib.shared().appsFlyerDevKey = "YOUR_AF_DEV_KEY"
        AppsFlyerLib.shared().appleAppID = "YOUR_APPLE_APP_ID"
        AppsFlyerLib.shared().delegate = self
        AppsFlyerLib.shared().deepLinkDelegate = self
        AppsFlyerLib.shared().start()
        return true
    }

    // Deferred deep link (first install)
    func onConversionDataSuccess(_ conversionInfo: [AnyHashable: Any]) {
        guard let isFirstLaunch = conversionInfo["is_first_launch"] as? Bool, isFirstLaunch else { return }
        let deepLinkValue = conversionInfo["deep_link_value"] as? String
        // navigate to deepLinkValue
    }

    // Direct deep link — UDL API
    func didResolveDeepLink(_ deepLinkResult: DeepLinkResult) {
        guard deepLinkResult.status == .found else { return }
        let deepLinkValue = deepLinkResult.deepLink?.deepLinkValue
        // navigate to deepLinkValue
    }

    func application(_ app: UIApplication, open url: URL,
                     options: [UIApplication.OpenURLOptionsKey: Any] = [:]) -> Bool {
        AppsFlyerLib.shared().handleOpen(url, options: options)
        return true
    }

    func application(_ application: UIApplication,
                     continue userActivity: NSUserActivity,
                     restorationHandler: @escaping ([UIUserActivityRestoring]?) -> Void) -> Bool {
        AppsFlyerLib.shared().continue(userActivity, restorationHandler: nil)
        return true
    }
}
```

### After (Detour)
```swift
// AppDelegate.swift
import Detour

@main
class AppDelegate: UIResponder, UIApplicationDelegate {

    let detourConfig = DetourConfig(
        apiKey: "YOUR_API_KEY",    // Detour Dashboard → API Configuration
        appID: "YOUR_APP_ID",      // Detour Dashboard → API Configuration
        shouldUseClipboard: true   // improves deferred link matching on iOS
    )

    func application(_ application: UIApplication,
                     didFinishLaunchingWithOptions launchOptions: [UIApplication.LaunchOptionsKey: Any]?) -> Bool {
        Detour.shared.mountAnalytics(config: detourConfig)

        // Handles both deferred and direct links on cold start
        Detour.shared.resolveInitialLink(config: detourConfig, launchOptions: launchOptions) { result in
            self.handleLink(result)
        }
        return true
    }

    // Custom URI scheme (yourapp://...)
    func application(_ app: UIApplication, open url: URL,
                     options: [UIApplication.OpenURLOptionsKey: Any] = [:]) -> Bool {
        Task { @MainActor in
            let result = await Detour.shared.processLink(url, config: self.detourConfig)
            self.handleLink(result)
        }
        return true
    }

    // Universal Link (https://YOUR_ORG.godetour.link/...)
    func application(_ application: UIApplication,
                     continue userActivity: NSUserActivity,
                     restorationHandler: @escaping ([UIUserActivityRestoring]?) -> Void) -> Bool {
        guard userActivity.activityType == NSUserActivityTypeBrowsingWeb,
              let url = userActivity.webpageURL else { return false }
        Task { @MainActor in
            let result = await Detour.shared.processLink(url, config: self.detourConfig)
            self.handleLink(result)
        }
        return true
    }

    private func handleLink(_ result: DetourResult) {
        guard let link = result.link else { return }

        // link.type:
        //   .deferred — user clicked link before installing
        //   .verified — Universal Link (domain verified)
        //   .scheme   — custom URI scheme
        //
        // link.route    — e.g. "/products/123?color=red"
        // link.pathname — e.g. "/products/123"
        // link.params   — e.g. ["color": "red"]

        NotificationCenter.default.post(name: .detourLinkReceived, object: link)
    }
}
```

**SwiftUI / SceneDelegate apps:**
```swift
// SceneDelegate.swift
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

---

## Deferred Deep Links (First Install)

Branch's `initSession` and AppsFlyer's `onConversionDataSuccess` had separate deferred-link paths.

With Detour, `resolveInitialLink` handles both. Check `link.type == .deferred`:

```swift
Detour.shared.resolveInitialLink(config: detourConfig, launchOptions: launchOptions) { result in
    guard let link = result.link else { return }
    if link.type == .deferred {
        // First install — user clicked a Detour link before installing the app
    }
    // navigate using link.route regardless of type
}
```

For testing, reset the first-launch flag:
```swift
Detour.shared.resetSession(allowDeferredRetry: true)
```

---

## Analytics

### Event mapping

| Branch (`BranchEvent`) | AppsFlyer | Detour (`DetourEventName`) |
|------------------------|-----------|----------------------------|
| `.purchase` | `AFEventPurchase` | `.purchase` |
| `.addToCart` | `AFEventAddToCart` | `.addToCart` |
| `.removeFromCart` | — | `.removeFromCart` |
| `.viewItem` | `AFEventContentView` | `.viewItem` |
| `.initiatePurchase` | `AFEventInitiatedCheckout` | `.beginCheckout` |
| `.addPaymentInfo` | `AFEventAddPaymentInfo` | `.addPaymentInfo` |
| `.login` | `AFEventLogin` | `.login` |
| `.completeRegistration` | `AFEventCompleteRegistration` | `.signUp` |
| `.search` | `AFEventSearch` | `.search` |
| `.share` | `AFEventShare` | `.share` |
| `.invite` | — | `.invite` |
| custom string | custom string | custom string |

### Before (Branch)
```swift
let event = BranchEvent.standardEvent(.purchase)
event.revenue = 29.99
event.currency = BNCCurrency.USD
event.customData = ["product_id": "abc123"]
event.logEvent()
```

### Before (AppsFlyer)
```swift
AppsFlyerLib.shared().logEvent(
    AFEventPurchase,
    withValues: [
        AFEventParamRevenue: 29.99,
        AFEventParamCurrency: "USD",
        "product_id": "abc123"
    ]
)
```

### After (Detour)
```swift
DetourAnalytics.logEvent(.purchase, data: [
    "revenue": 29.99,
    "currency": "USD",
    "product_id": "abc123"
])
```

Custom events:
```swift
DetourAnalytics.logEvent("promo_banner_tapped", data: ["placement": "home_top"])
```

Retention / session events:
```swift
DetourAnalytics.logRetention("app_open")
```

---

## Keeping this reference current

If the user asks about the latest version, installation steps, or you suspect anything above may be outdated, fetch the live README:

```
https://github.com/software-mansion-labs/ios-detour/blob/main/README.md
```

Cross-reference against the sections above and prefer the live README if they conflict.
