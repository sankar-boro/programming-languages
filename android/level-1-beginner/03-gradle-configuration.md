# Chapter 3: Gradle Configuration for AndroidX Projects

## What Is Gradle?

Gradle is the **build system** for Android. It compiles your code, packages your resources, and produces an APK or AAB.

**Analogy:** Gradle is the chef in your kitchen. You write a recipe (`build.gradle`) and Gradle executes it — gathering ingredients (dependencies), mixing them (compiling), and plating the dish (building the APK).

---

## The Two Build Files

Modern Android projects use **Kotlin DSL** (`*.kts`) for build files.

```
project/
├── build.gradle.kts          ← Project-level (top-level)
├── settings.gradle.kts       ← Module declarations, repository sources
└── app/
    └── build.gradle.kts      ← App-level (module-level)
```

---

## `settings.gradle.kts` — Repository Sources

```kotlin
pluginManagement {
    repositories {
        google()
        mavenCentral()
        gradlePluginPortal()
    }
}

dependencyResolutionManagement {
    repositoriesMode.set(RepositoriesMode.FAIL_ON_PROJECT_REPOS)
    repositories {
        google()        // AndroidX, Google libraries
        mavenCentral()  // Kotlin, most open-source libs
    }
}

rootProject.name = "MyFirstApp"
include(":app")
```

---

## `build.gradle.kts` — Project Level

```kotlin
// Top-level build file: applies plugins to all submodules
plugins {
    alias(libs.plugins.android.application) apply false
    alias(libs.plugins.kotlin.android) apply false
    alias(libs.plugins.kotlin.compose) apply false
}
```

> `apply false` means: "make this plugin available, but don't apply it here — apply it in each module."

---

## `app/build.gradle.kts` — App Level (Full Example)

```kotlin
plugins {
    alias(libs.plugins.android.application)
    alias(libs.plugins.kotlin.android)
}

android {
    namespace = "com.yourname.myfirstapp"
    compileSdk = 35  // Always use the latest stable SDK to compile

    defaultConfig {
        applicationId = "com.yourname.myfirstapp"
        minSdk = 24           // Minimum Android version supported
        targetSdk = 35        // The version you've tested against
        versionCode = 1       // Internal version number (integer, increment each release)
        versionName = "1.0"   // User-facing version string

        testInstrumentationRunner = "androidx.test.runner.AndroidJUnitRunner"
    }

    buildTypes {
        release {
            isMinifyEnabled = true          // Enable ProGuard/R8 shrinking
            proguardFiles(
                getDefaultProguardFile("proguard-android-optimize.txt"),
                "proguard-rules.pro"
            )
        }
        debug {
            isDebuggable = true
            applicationIdSuffix = ".debug"  // Allows debug and release installed simultaneously
        }
    }

    compileOptions {
        sourceCompatibility = JavaVersion.VERSION_11
        targetCompatibility = JavaVersion.VERSION_11
    }

    kotlinOptions {
        jvmTarget = "11"
    }

    buildFeatures {
        viewBinding = true
        // buildConfig = true  // Enable if you need BuildConfig fields
    }
}

dependencies {
    // AndroidX Core
    implementation(libs.androidx.core.ktx)
    implementation(libs.androidx.appcompat)

    // Material Design
    implementation(libs.material)

    // ConstraintLayout
    implementation(libs.androidx.constraintlayout)

    // Testing
    testImplementation(libs.junit)
    androidTestImplementation(libs.androidx.junit)
    androidTestImplementation(libs.androidx.espresso.core)
}
```

---

## Version Catalogs (`libs.versions.toml`)

Modern projects use a **version catalog** to centralize dependency versions.

Location: `gradle/libs.versions.toml`

```toml
[versions]
agp = "8.7.3"
kotlin = "2.1.0"
coreKtx = "1.15.0"
appcompat = "1.7.0"
material = "1.12.0"
constraintlayout = "2.2.0"
lifecycle = "2.8.7"
navigationFragment = "2.8.4"
room = "2.6.1"
hilt = "2.51.1"
junit = "4.13.2"
junitVersion = "1.2.1"
espressoCore = "3.6.1"

[libraries]
androidx-core-ktx = { group = "androidx.core", name = "core-ktx", version.ref = "coreKtx" }
androidx-appcompat = { group = "androidx.appcompat", name = "appcompat", version.ref = "appcompat" }
material = { group = "com.google.android.material", name = "material", version.ref = "material" }
androidx-constraintlayout = { group = "androidx.constraintlayout", name = "constraintlayout", version.ref = "constraintlayout" }
androidx-lifecycle-viewmodel-ktx = { group = "androidx.lifecycle", name = "lifecycle-viewmodel-ktx", version.ref = "lifecycle" }
androidx-lifecycle-livedata-ktx = { group = "androidx.lifecycle", name = "lifecycle-livedata-ktx", version.ref = "lifecycle" }
junit = { group = "junit", name = "junit", version.ref = "junit" }
androidx-junit = { group = "androidx.test.ext", name = "junit", version.ref = "junitVersion" }
androidx-espresso-core = { group = "androidx.test.espresso", name = "espresso-core", version.ref = "espressoCore" }

[plugins]
android-application = { id = "com.android.application", version.ref = "agp" }
kotlin-android = { id = "org.jetbrains.kotlin.android", version.ref = "kotlin" }
```

**Why use version catalogs?**
- One place to update a dependency version for the whole project
- Enables sharing versions across multi-module projects
- IDE autocomplete for `libs.*` references

---

## Key `defaultConfig` Fields Explained

| Field | Purpose |
|-------|---------|
| `applicationId` | Unique identifier for your app on Google Play (like a domain in reverse) |
| `minSdk` | Lowest Android version your app supports |
| `targetSdk` | The Android version you've tested against — affects system behavior |
| `compileSdk` | The Android SDK version used to compile — not visible to users |
| `versionCode` | Integer; must increment with each Play Store release |
| `versionName` | String shown to users (e.g., "2.1.0") |

---

## `minSdk` vs `targetSdk` vs `compileSdk`

```
compileSdk ──────────── What APIs you can call in code
                         (always latest stable)
     │
     ▼
targetSdk ───────────── How the OS treats your app
                         (use same as compileSdk unless you have a reason)
     │
     ▼
minSdk ──────────────── Floor — minimum OS version to install your app
                         (API 24 covers ~95% of active devices)
```

---

## Build Types and Product Flavors

### Build Types

```kotlin
buildTypes {
    debug {
        isDebuggable = true
        applicationIdSuffix = ".debug"
        versionNameSuffix = "-DEBUG"
    }
    release {
        isMinifyEnabled = true
        isShrinkResources = true
        proguardFiles(
            getDefaultProguardFile("proguard-android-optimize.txt"),
            "proguard-rules.pro"
        )
        signingConfig = signingConfigs.getByName("release")
    }
}
```

### Product Flavors (for multi-variant apps)

```kotlin
flavorDimensions += "environment"
productFlavors {
    create("dev") {
        dimension = "environment"
        applicationIdSuffix = ".dev"
        buildConfigField("String", "BASE_URL", "\"https://dev.api.example.com\"")
    }
    create("prod") {
        dimension = "environment"
        buildConfigField("String", "BASE_URL", "\"https://api.example.com\"")
    }
}
```

This creates build variants: `devDebug`, `devRelease`, `prodDebug`, `prodRelease`.

---

## Syncing and Building

| Action | Shortcut |
|--------|---------|
| Sync Gradle | Click "Sync Now" in the notification bar, or File → Sync Project with Gradle Files |
| Build APK | Build → Build Bundle(s)/APK(s) → Build APK(s) |
| Run on device | Shift+F10 (or click ▶) |
| Clean build | Build → Clean Project |

---

## Common Gradle Errors and Fixes

### Error: `Duplicate class kotlin.collections.jdk8`

**Cause:** Kotlin stdlib version conflict.
**Fix:** Add to `gradle.properties`:
```properties
kotlin.stdlib.default.dependency=false
```

### Error: `Minimum supported Gradle version is X`

**Fix:** Update the Gradle wrapper version in `gradle/wrapper/gradle-wrapper.properties`:
```properties
distributionUrl=https\://services.gradle.org/distributions/gradle-8.9-bin.zip
```

### Error: `Could not resolve com.example:library:1.0.0`

**Cause:** Missing repository or wrong version.
**Fix:** Ensure `google()` and `mavenCentral()` are in `settings.gradle.kts` repositories.

### Error: `Android Gradle plugin requires Java 17`

**Fix:** In Android Studio → Settings → Build → Gradle → Gradle JDK → select JDK 17 or 21.

---

## Interview Questions

**Q1: What is the difference between `compileSdk`, `targetSdk`, and `minSdk`?**

> `compileSdk` is the version of the Android SDK the code is compiled against — it determines which APIs are available. `targetSdk` tells the OS what behavior your app expects; new OS behavior changes only apply if `targetSdk` matches. `minSdk` is the minimum Android version that can install the app.

**Q2: What is a version catalog in Gradle and why use one?**

> A `libs.versions.toml` file centralizes all dependency versions in one place. It allows consistent versioning across modules, enables IDE autocomplete, and makes version upgrades a single-file change.

**Q3: What does `isMinifyEnabled = true` do?**

> It enables R8 (the code shrinker/obfuscator), which removes unused code and renames classes to reduce APK size. Always enable it for release builds.

---

## Summary

- Android uses Gradle with Kotlin DSL (`*.kts`) for build configuration
- `settings.gradle.kts` declares repositories and modules
- `app/build.gradle.kts` configures `compileSdk`, `minSdk`, `targetSdk`, dependencies, and build types
- Version catalogs (`libs.versions.toml`) centralize dependency versions
- Always use `isMinifyEnabled = true` for release builds

**Next:** [Chapter 4 — androidx.core and androidx.appcompat](./04-core-and-appcompat.md)
