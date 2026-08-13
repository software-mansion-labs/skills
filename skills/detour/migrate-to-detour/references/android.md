# Android — Branch / AppsFlyer → Detour

## Contents
- Universal / App Links
- SDK Installation
- SDK Initialization & Deep Link Handling
- Deferred Deep Links (First Install)
- Analytics
- Testing
- Keeping this reference current

## Universal / App Links

### Before (Branch)
```xml
<!-- AndroidManifest.xml -->
<activity android:name=".MainActivity" android:exported="true">
    <intent-filter android:autoVerify="true">
        <action android:name="android.intent.action.VIEW" />
        <category android:name="android.intent.category.DEFAULT" />
        <category android:name="android.intent.category.BROWSABLE" />
        <data android:scheme="https" android:host="yourapp.app.link" />
        <data android:scheme="https" android:host="yourapp-alternate.app.link" />
    </intent-filter>
</activity>
```

### Before (AppsFlyer)
```xml
<activity android:name=".MainActivity" android:exported="true">
    <intent-filter android:autoVerify="true">
        <action android:name="android.intent.action.VIEW" />
        <category android:name="android.intent.category.DEFAULT" />
        <category android:name="android.intent.category.BROWSABLE" />
        <data android:scheme="https" android:host="yourapp.onelink.me" />
    </intent-filter>
</activity>
```

### After (Detour)
Replace the host with your Detour organization subdomain. Detour auto-hosts `/.well-known/assetlinks.json` — no server setup needed.

**Important:** The Detour Dashboard generates a unique `pathPrefix` for your app (e.g. `/ldukjiuliD`). After creating the app in the dashboard, copy the generated intent filter snippet from the **SDK Setup** tab and use it directly — don't write the intent filter by hand.

The generated snippet will look like this:
```xml
<!-- AndroidManifest.xml -->
<activity android:name=".MainActivity" android:exported="true">

    <!-- App Links (https) — copy exact snippet from Detour Dashboard → SDK Setup -->
    <intent-filter android:autoVerify="true">
        <action android:name="android.intent.action.VIEW" />
        <category android:name="android.intent.category.DEFAULT" />
        <category android:name="android.intent.category.BROWSABLE" />
        <data android:scheme="https" android:host="YOUR_ORG.godetour.link" android:pathPrefix="/YOUR_PATH_PREFIX" />
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

### Requirements

Before adding the SDK, make sure your project meets these minimum versions — Detour SDK requires them:

| Tool | Minimum version |
|------|----------------|
| `compileSdk` | 34 |
| Android Gradle Plugin | 8.1.4 |
| Gradle | 8.0 |
| Kotlin | 1.9.0 |

Also ensure your `build.gradle (app)` has a `namespace` declaration — AGP 8.x no longer reads `package` from `AndroidManifest.xml`:
```kotlin
android {
    namespace = "com.yourcompany.yourapp"
    compileSdk = 34
    // ...
}
```

### 1. Add the Software Mansion Maven repository

Detour is published on Software Mansion's Maven repository. Add it to `settings.gradle`:

```kotlin
// settings.gradle (or settings.gradle.kts)
dependencyResolutionManagement {
    repositories {
        google()
        mavenCentral()
        maven { url = uri("https://maven.swmansion.com") } // required for Detour
    }
}
```

### 2. Add the dependency

```kotlin
// build.gradle (app)
dependencies {
    // Remove:
    // implementation("io.branch.sdk.android:library:5.x.x")
    // implementation("com.appsflyer:af-android-sdk:6.x.x")

    // Add (use the latest version — check Maven Central / the README):
    implementation("com.swmansion.detour:detour-sdk:LATEST_VERSION")
}
```

---

## SDK Initialization & Deep Link Handling

### Before (Branch)
```kotlin
// MyApplication.kt
class MyApplication : Application() {
    override fun onCreate() {
        super.onCreate()
        Branch.getAutoInstance(this)
    }
}

// MainActivity.kt
override fun onStart() {
    super.onStart()
    Branch.sessionBuilder(this)
        .withCallback { branchUniversalObject, linkProperties, error ->
            if (error != null) return@withCallback
            val canonicalUrl = branchUniversalObject?.canonicalUrl
            // navigate using canonicalUrl
        }
        .withData(intent?.data)
        .init()
}

override fun onNewIntent(intent: Intent?) {
    super.onNewIntent(intent)
    setIntent(intent)
    Branch.sessionBuilder(this).withCallback { buo, lp, error -> /* ... */ }.reInit()
}
```

### Before (AppsFlyer)
```kotlin
// MyApplication.kt
class MyApplication : Application() {
    override fun onCreate() {
        super.onCreate()
        AppsFlyerLib.getInstance().init("YOUR_AF_DEV_KEY", null, this)
        AppsFlyerLib.getInstance().subscribeForDeepLink { result ->
            if (result.status == DeepLinkResult.Status.FOUND) {
                val deepLinkValue = result.deepLink?.deepLinkValue
                // navigate using deepLinkValue
            }
        }
        AppsFlyerLib.getInstance().start(this)
    }
}
```

### After (Detour)
```kotlin
import com.swmansion.detour.Detour
import com.swmansion.detour.DetourConfig
import com.swmansion.detour.DetourDelegate
import com.swmansion.detour.models.LinkResult
import com.swmansion.detour.models.LinkType

// MainActivity.kt
class MainActivity : AppCompatActivity() {

    private val detourDelegate = DetourDelegate(
        lifecycleOwner = this,
        config = DetourConfig(
            apiKey = "YOUR_API_KEY",   // Detour Dashboard → API Configuration
            appId = "YOUR_APP_ID"      // Detour Dashboard → API Configuration
        ),
        onLinkResult = { result -> handleLink(result) }
    )

    override fun onCreate(savedInstanceState: Bundle?) {
        super.onCreate(savedInstanceState)
        Detour.initialize(this, detourDelegate.config)
        detourDelegate.onCreate(intent)
    }

    override fun onNewIntent(intent: Intent) {
        super.onNewIntent(intent)
        setIntent(intent)
        detourDelegate.onNewIntent(intent)
    }

    private fun handleLink(result: LinkResult) {
        when (result) {
            is LinkResult.Success -> {
                // result.type:
                //   LinkType.DEFERRED  — user clicked link before installing
                //   LinkType.VERIFIED  — App Link (domain verified)
                //   LinkType.SCHEME    — custom URI scheme
                //
                // result.route    — e.g. "/products/123?color=red"
                // result.pathname — e.g. "/products/123"
                // result.params   — e.g. mapOf("color" to "red")

                navigateTo(result.route, result.params)
            }
            is LinkResult.NoLink -> { /* normal launch */ }
            is LinkResult.NotFirstLaunch -> { /* deferred already consumed */ }
            is LinkResult.Error -> { Log.e("Detour", "Error", result.exception) }
        }
    }
}
```

No Application class changes needed — `Detour.initialize` lives in `MainActivity.onCreate`.

---

## Deferred Deep Links (First Install)

With Detour, deferred links arrive through the same `onLinkResult` callback. Identify them with `result.type == LinkType.DEFERRED`.

If migrating gradually and need deferred-only mode:

```kotlin
Detour.initialize(
    this,
    DetourConfig(
        apiKey = "YOUR_API_KEY",
        appId = "YOUR_APP_ID",
        linkProcessingMode = LinkProcessingMode.DEFERRED_ONLY
    )
)

lifecycleScope.launch {
    val result = Detour.getDeferredLink()
    handleLink(result)
}
```

---

## Analytics

### Event mapping

| Branch (`BRANCH_STANDARD_EVENT`) | AppsFlyer (`AFInAppEventType`) | Detour (`DetourEventNames`) |
|----------------------------------|--------------------------------|-----------------------------|
| `PURCHASE` | `PURCHASE` | `Purchase` |
| `ADD_TO_CART` | `ADD_TO_CART` | `AddToCart` |
| `REMOVE_FROM_CART` | — | `RemoveFromCart` |
| `VIEW_ITEM` | `CONTENT_VIEW` | `ViewItem` |
| `INITIATE_PURCHASE` | `INITIATED_CHECKOUT` | `BeginCheckout` |
| `ADD_PAYMENT_INFO` | `ADD_PAYMENT_INFO` | `AddPaymentInfo` |
| `LOGIN` | `LOGIN` | `Login` |
| `COMPLETE_REGISTRATION` | `COMPLETE_REGISTRATION` | `SignUp` |
| `SEARCH` | `SEARCH` | `Search` |
| `SHARE` | `SHARE` | `Share` |
| `INVITE` | — | `Invite` |

### Before (Branch)
```kotlin
BranchEvent(BRANCH_STANDARD_EVENT.PURCHASE)
    .setCurrency(CurrencyType.USD)
    .setRevenue(29.99)
    .addCustomDataProperty("product_id", "abc123")
    .logEvent(context)
```

### Before (AppsFlyer)
```kotlin
AppsFlyerLib.getInstance().logEvent(
    context,
    AFInAppEventType.PURCHASE,
    mapOf(
        AFInAppEventParameterName.REVENUE to 29.99,
        AFInAppEventParameterName.CURRENCY to "USD",
        "product_id" to "abc123"
    )
)
```

### After (Detour)
```kotlin
DetourAnalytics.logEvent(
    DetourEventNames.Purchase,
    mapOf(
        "revenue" to 29.99,
        "currency" to "USD",
        "product_id" to "abc123"
    )
)
```

**Custom events** — `DetourEventNames` is an enum, so arbitrary strings are not accepted by `logEvent`. For events that don't have a `DetourEventNames` equivalent, use `logRetention`:
```kotlin
DetourAnalytics.logRetention("promo_banner_tapped")
```

Session / app open events:
```kotlin
DetourAnalytics.logRetention("app_open")
```

---

## Testing

### Android 12+ — enabling App Links for debug builds

On Android 12+, there are two independent layers: domain verification (via `assetlinks.json`) and **user selection**. When you install a debug APK directly (via ADB, not Play Store), user selection defaults to `Disabled` — links won't open the app even if domain verification passed.

Enable it manually after each install:
```bash
adb shell pm set-app-links-user-selection \
  --user 0 --package YOUR_PACKAGE_NAME true YOUR_ORG.godetour.link
```

This resets on every reinstall during development. On production builds installed from the Play Store, it's enabled automatically.

### Testing deferred deep links

Two rules that are easy to get wrong:

1. **Click the link after uninstalling, not before.** The flow is: uninstall app → click Detour link in browser → install app → open app. If you click the link while the app is installed, it resolves as a normal Universal Link, not a deferred one.

2. **Each deferred link test needs a fresh link.** Detour tracks device fingerprints — the same link won't trigger deferred resolution a second time on the same device. Generate a new link from the dashboard for each test run.

---

## Keeping this reference current

If the user asks about the latest version, installation steps, or you suspect anything above may be outdated, fetch the live README:

```
https://github.com/software-mansion-labs/android-detour/blob/main/README.md
```

Cross-reference against the sections above and prefer the live README if they conflict.
