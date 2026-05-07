---
name: rnrepo-configuration
description: "RNRepo configuration options: opting out specific libraries with denyList, disabling the plugin temporarily or permanently, Expo Fingerprint ignore, custom React Native directory, refreshing Gradle cache, and GPG artifact verification."
---

# RNRepo Configuration

## Opting out specific libraries

Create `rnrepo.config.json` in the React Native project root:

```json
{
  "denyList": {
    "android": ["library-name-1", "library-name-2"],
    "ios": ["library-name-3", "library-name-4"]
  }
}
```

Libraries in the deny list always build from source. Use this for libraries with **native patches** (Objective-C, Java, or Kotlin changes), build-time feature flags, or persistent C++ compatibility issues. JS-only patches do not require opting out - they don't affect the native artifact.

---

## Disabling RNRepo for a single build

Set `DISABLE_RNREPO` to any value before the build command:

```bash
# Android
DISABLE_RNREPO=1 ./gradlew :app:assembleDebug

# iOS (requires pod reinstall)
DISABLE_RNREPO=1 pod install
```

This skips the plugin entirely. Use for debugging build failures, not permanent exclusion.

---

## Expo Fingerprint / EAS — ignore RNRepo cache

The CocoaPods plugin caches downloaded xcframeworks in `node_modules/<pkg>/.rnrepo-cache`, which changes the fingerprint between local and CI machines. Exclude it:

**`.fingerprintignore`:**
```text
**/.rnrepo-cache
```

Or in `fingerprint.config.js` — see the [Expo Fingerprint docs](https://docs.expo.dev/versions/latest/sdk/fingerprint/#fingerprintignore).

---

## Custom React Native directory (Android)

When `node_modules/` or `rnrepo.config.json` is not at the standard path, tell the plugin where to look.

Via environment variable (takes precedence):
```bash
REACT_NATIVE_ROOT_DIR=/path/to/rn-root ./gradlew :app:assembleDebug
```

Via `gradle.properties`:
```properties
REACT_NATIVE_ROOT_DIR=/path/to/rn-root
```

---

## Refreshing the Gradle cache

Force re-download of all prebuilt artifacts:
```bash
./gradlew :app:assembleDebug --refresh-dependencies
```

Manually clear the RNRepo Gradle cache:
```bash
rm -rf ~/.gradle/caches/modules-2/metadata-2.107/descriptors/org.rnrepo.public
rm -rf ~/.gradle/caches/modules-2/files-2.1/org.rnrepo.public
rm -rf ~/.gradle/caches/modules-2/metadata-2.107/descriptors/org.rnrepo.tools
rm -rf ~/.gradle/caches/modules-2/files-2.1/org.rnrepo.tools
```

---

## GPG artifact verification (enterprise)

Every artifact is signed with Software Mansion's GPG key.

1. Import the key:
   ```bash
   curl https://keys.openpgp.org/vks/v1/by-fingerprint/6CBF6E07EBA0219DF11C9F78C9ED010ADBD95DFE | gpg --import
   ```
2. Configure Gradle dependency verification per the [Android docs](https://developer.android.com/build/dependency-verification).
3. Store the fingerprint (`6CBF6E07EBA0219DF11C9F78C9ED010ADBD95DFE`) in your policy tooling to alert on key rotations.

---

## `FAIL_ON_PROJECT_REPOS` / `PREFER_SETTINGS` Gradle mode

If Gradle rejects project-level repository additions, manually register the Maven repo in your root `build.gradle`:

```groovy
allprojects {
    repositories {
        maven { url "https://packages.rnrepo.org/releases" }
    }
}
```

Make sure this `allprojects` block appears **before** applying the RNRepo plugin.
