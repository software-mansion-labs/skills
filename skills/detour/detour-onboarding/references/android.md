# Detour — Android SDK Reference

## Requirements

- `compileSdk` 34+
- Android Gradle Plugin 8.1.4+
- Gradle 8.0+
- Kotlin 1.9.0+
- JVM target 11+

Also ensure your `build.gradle (app)` has a `namespace` declaration (required by AGP 8.x):

```kotlin
android {
    namespace = "com.yourcompany.yourapp"
    compileSdk = 34
}
```

## Installation

`app/build.gradle.kts`:

```kotlin
dependencies {
    implementation("com.swmansion.detour:detour-sdk:1.1.0")
}
```

No custom Maven repository needed — published on Maven Central.

### local.properties

If the file `android/local.properties` does not exist, create it and add:

```
sdk.dir=/Users/YOUR_USERNAME/Library/Android/sdk
```

Gradle cannot find the Android SDK without this file. It is git-ignored by default.

---

## App Links

**Always copy the intent-filter snippet from the Detour Dashboard** (App Configuration tab). The dashboard generates a snippet with your unique `pathPrefix` (e.g. `/VnAqasAabE`). If you write the intent-filter manually without this pathPrefix, App Links will not work.

The generated snippet looks like this:

```xml
<activity android:name=".MainActivity" android:exported="true">

  <!-- App Links — COPY EXACT SNIPPET FROM DETOUR DASHBOARD → App Configuration -->
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

### Verify package name matches exactly

The package name in the Detour Dashboard (App Configuration tab) must match exactly what is in your app's `build.gradle` `applicationId` / `namespace`. A mismatch causes `assetlinks.json` verification to fail silently.

---

## SDK Setup

### 1. Store credentials

`local.properties` (git-ignored):

```
DETOUR_API_KEY=your_publishable_key
DETOUR_APP_ID=your_app_id
```

`app/build.gradle.kts`:

```kotlin
buildConfigField("String", "DETOUR_API_KEY", "\"${project.findProperty("DETOUR_API_KEY") ?: ""}\"")
buildConfigField("String", "DETOUR_APP_ID",  "\"${project.findProperty("DETOUR_APP_ID")  ?: ""}\"")
```

### 2. Handle links in Activity

```kotlin
import com.swmansion.detour.Detour
import com.swmansion.detour.DetourConfig
import com.swmansion.detour.DetourDelegate
import com.swmansion.detour.models.LinkResult
import com.swmansion.detour.models.LinkType

class MainActivity : AppCompatActivity() {

    private val detourDelegate = DetourDelegate(
        lifecycleOwner = this,
        config = DetourConfig(
            apiKey = BuildConfig.DETOUR_API_KEY,
            appId  = BuildConfig.DETOUR_APP_ID
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
            is LinkResult.Success        -> {
                // result.type:     LinkType.DEFERRED | VERIFIED | SCHEME
                // result.route:    "/articles/2?ref=campaign"
                // result.pathname: "/articles/2"
                // result.params:   mapOf("ref" to "campaign")
                navigateTo(result.route, result.params)
            }
            is LinkResult.NoLink         -> { /* normal launch */ }
            is LinkResult.NotFirstLaunch -> { /* deferred already consumed */ }
            is LinkResult.Error          -> { Log.e("Detour", "Error", result.exception) }
        }
    }
}
```

No `Application` class changes needed — `Detour.initialize` lives in `MainActivity.onCreate`.

### Deferred-only manual flow

```kotlin
lifecycleScope.launch {
    val result = Detour.getDeferredLink()
    if (result is LinkResult.Success) navigateTo(result.route, result.params)
}
```

---

## Testing on Android 12+

### Enabling App Links for debug builds

On Android 12+, debug APKs installed via ADB have App Links user selection defaulting to Disabled. Enable it manually after each install:

```bash
adb shell pm set-app-links --package com.yourapp.package 2 all
```

This resets on every reinstall. Production Play Store installs enable it automatically.

If App Links verification is stuck (status `1024 / NO_RESPONSE`), force it:

```bash
adb shell pm set-app-links --package com.yourapp.package 2 all
# Then verify the status:
adb shell pm get-app-links --package com.yourapp.package
```

### onNewIntent when app is already in background

If the app is in the background when you tap an App Link, `onNewIntent` fires but the link context may not be ready yet, causing a `context is not ready` error. During development, **kill the app before testing** rather than using the back button to background it.

### Deferred link testing rules

1. Uninstall the app before clicking the link — clicking while installed resolves as App Link, not deferred
2. Each test needs a fresh link — the same deferred link won't fire twice on the same device

---

## Analytics

Auto-initialized with `Detour.initialize()`.

```kotlin
// Standard events — use DetourEventNames enum
DetourAnalytics.logEvent(DetourEventNames.Purchase, mapOf("revenue" to 29.99, "currency" to "USD"))
DetourAnalytics.logEvent(DetourEventNames.ViewItem, mapOf("item_id" to "article_2"))

// Custom / non-standard activity — use logRetention
DetourAnalytics.logRetention("article_opened")
DetourAnalytics.logRetention("app_open")
```

---

## Keeping this reference current

If the user asks about the latest version, installation steps, or you suspect anything above may be outdated, fetch the live README:

```
https://github.com/software-mansion-labs/android-detour/blob/main/README.md
```

Cross-reference against the sections above and prefer the live README if they conflict.
