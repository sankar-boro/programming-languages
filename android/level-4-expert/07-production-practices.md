# Chapter 7: Production Practices

## Error Handling

### Structured Error Hierarchy

```kotlin
// core/common/src/main/kotlin/error/AppError.kt
sealed class AppError(message: String, cause: Throwable? = null) : Exception(message, cause) {
    class NetworkError(message: String, cause: Throwable? = null) : AppError(message, cause)
    class NotFoundError(val id: String) : AppError("Resource not found: $id")
    class AuthError(message: String) : AppError(message)
    class ValidationError(val field: String, message: String) : AppError(message)
    class UnknownError(cause: Throwable) : AppError(cause.message ?: "Unknown error", cause)
}

// Map throwables to AppError in the data layer
fun Throwable.toAppError(): AppError = when (this) {
    is IOException -> AppError.NetworkError("Network unavailable", this)
    is HttpException -> when (code()) {
        401 -> AppError.AuthError("Unauthorized")
        404 -> AppError.NotFoundError("remote_resource")
        else -> AppError.NetworkError("Server error: ${code()}", this)
    }
    else -> AppError.UnknownError(this)
}
```

### Global Error Handler in ViewModel

```kotlin
@HiltViewModel
class BaseViewModel @Inject constructor() : ViewModel() {
    protected val _uiState = MutableStateFlow(UiState())
    val uiState: StateFlow<UiState> = _uiState.asStateFlow()

    protected fun launchWithErrorHandling(
        block: suspend () -> Unit
    ) = viewModelScope.launch {
        try {
            block()
        } catch (e: AppError.NetworkError) {
            _uiState.update { it.copy(error = UiError.Network(e.message ?: "Network error")) }
        } catch (e: AppError.AuthError) {
            _uiState.update { it.copy(error = UiError.Auth) }
        } catch (e: CancellationException) {
            throw e
        } catch (e: Exception) {
            _uiState.update { it.copy(error = UiError.Unknown(e.message)) }
        }
    }
}
```

---

## Logging with Timber

```kotlin
implementation("com.jakewharton.timber:timber:5.0.1")

// Application setup
@HiltAndroidApp
class MyApplication : Application() {
    override fun onCreate() {
        super.onCreate()
        if (BuildConfig.DEBUG) {
            Timber.plant(Timber.DebugTree())
        } else {
            Timber.plant(CrashReportingTree())
        }
    }
}

// CrashReportingTree — sends errors to Firebase Crashlytics in production
class CrashReportingTree : Timber.Tree() {
    override fun log(priority: Int, tag: String?, message: String, t: Throwable?) {
        if (priority == Log.ERROR || priority == Log.WARN) {
            // Send to crash reporting
            FirebaseCrashlytics.getInstance().apply {
                t?.let { recordException(it) }
                log("[$tag] $message")
            }
        }
    }
}

// Usage throughout the app
Timber.d("Loading notes for user: $userId")
Timber.e(exception, "Failed to sync notes")
Timber.w("Slow network detected, falling back to cache")
```

---

## Firebase Crashlytics Integration

```kotlin
// app/build.gradle.kts
plugins {
    id("com.google.gms.google-services")
    id("com.google.firebase.crashlytics")
}

dependencies {
    implementation(platform("com.google.firebase:firebase-bom:33.5.1"))
    implementation("com.google.firebase:firebase-crashlytics-ktx")
    implementation("com.google.firebase:firebase-analytics-ktx")
}

// Record non-fatal errors
fun recordError(error: Throwable, context: Map<String, String> = emptyMap()) {
    FirebaseCrashlytics.getInstance().apply {
        context.forEach { (key, value) -> setCustomKey(key, value) }
        recordException(error)
    }
}

// Set user context
fun setUserContext(userId: String) {
    FirebaseCrashlytics.getInstance().setUserId(userId)
}
```

---

## Secure Coding Practices

### Storing Sensitive Data

```kotlin
// NEVER store sensitive data in SharedPreferences or DataStore plaintext
// Use EncryptedSharedPreferences for tokens

implementation("androidx.security:security-crypto:1.1.0-alpha06")

val masterKey = MasterKey.Builder(context)
    .setKeyScheme(MasterKey.KeyScheme.AES256_GCM)
    .build()

val encryptedPrefs = EncryptedSharedPreferences.create(
    context,
    "secure_prefs",
    masterKey,
    EncryptedSharedPreferences.PrefKeyEncryptionScheme.AES256_SIV,
    EncryptedSharedPreferences.PrefValueEncryptionScheme.AES256_GCM
)

encryptedPrefs.edit { putString("auth_token", token) }
val token = encryptedPrefs.getString("auth_token", null)
```

### Certificate Pinning

```kotlin
val okHttpClient = OkHttpClient.Builder()
    .certificatePinner(
        CertificatePinner.Builder()
            .add("api.yourapp.com", "sha256/AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA=")
            .build()
    )
    .build()
```

### Obfuscation Rules

`proguard-rules.pro`:

```proguard
# Keep data classes (Room entities, API DTOs)
-keep class com.yourname.app.data.local.entity.** { *; }
-keep class com.yourname.app.data.remote.dto.** { *; }

# Keep Retrofit interfaces
-keep interface com.yourname.app.data.remote.** { *; }

# Keep Hilt-generated classes
-keep class dagger.hilt.** { *; }
-keep @dagger.hilt.android.HiltAndroidApp class *

# Keep Parcelable implementations
-keep class * implements android.os.Parcelable { *; }

# Keep Serializable
-keepnames class * implements java.io.Serializable
```

---

## App Architecture Checklist

```
✓ Domain layer has zero Android imports
✓ ViewModels don't hold Context references
✓ Fragment binding cleared in onDestroyView()
✓ All coroutines use structured scopes (viewModelScope, lifecycleScope)
✓ CancellationException always re-thrown
✓ Room queries return Flow (reactive) or suspend (one-shot)
✓ Hilt is the only way dependencies are created
✓ No hardcoded API keys in code (use BuildConfig or secrets manager)
✓ Lint passes with zero errors
✓ Unit tests cover ViewModel and Repository
✓ UI tests cover critical user flows
✓ R8 enabled for release builds
✓ Crashlytics integrated for production error tracking
✓ No debug logs in release builds (Timber handles this)
```

---

## Code Quality Tools Configuration

### `.editorconfig`

```ini
root = true

[*.{kt,kts}]
indent_style = space
indent_size = 4
max_line_length = 120
end_of_line = lf
charset = utf-8
trim_trailing_whitespace = true
insert_final_newline = true
```

### `ktlint` Integration

```kotlin
// build.gradle.kts (project level)
plugins {
    id("org.jlleitschuh.gradle.ktlint") version "12.1.2"
}

ktlint {
    version.set("1.5.0")
    android.set(true)
    outputToConsole.set(true)
    filter {
        exclude("**/generated/**")
    }
}
```

Run: `./gradlew ktlintCheck` to check, `./gradlew ktlintFormat` to fix.

---

## Handling API Keys Safely

**In development — `local.properties` (never commit this file):**

```properties
NEWS_API_KEY=your_actual_key_here
```

**`build.gradle.kts` — read and inject as `BuildConfig` field:**

```kotlin
android {
    val localProperties = Properties().apply {
        load(rootProject.file("local.properties").inputStream())
    }

    defaultConfig {
        buildConfigField("String", "NEWS_API_KEY",
            "\"${localProperties.getProperty("NEWS_API_KEY", "")}\"")
    }

    buildFeatures {
        buildConfig = true
    }
}
```

**In CI — inject via GitHub Secrets:**

```yaml
- name: Build
  run: ./gradlew assembleRelease
  env:
    NEWS_API_KEY: ${{ secrets.NEWS_API_KEY }}
```

**Usage:**

```kotlin
const val API_KEY = BuildConfig.NEWS_API_KEY
```

---

## Scalability Patterns

### Feature Flags

```kotlin
interface FeatureFlag {
    val isEnabled: Boolean
}

enum class Feature : FeatureFlag {
    DARK_MODE_V2 {
        override val isEnabled = BuildConfig.DEBUG || RemoteConfig.getBoolean("dark_mode_v2")
    },
    NEW_SEARCH_UI {
        override val isEnabled = RemoteConfig.getBoolean("new_search_ui")
    }
}

// Usage
if (Feature.NEW_SEARCH_UI.isEnabled) {
    showNewSearchUI()
} else {
    showLegacySearchUI()
}
```

### Pagination Strategy

- Use Paging 3 for all lists over 50 items
- Use `pageSize = 20` as default
- Always cache with `cachedIn(viewModelScope)`

### Database Performance

```kotlin
// Add indices for frequently queried columns
@Entity(
    tableName = "notes",
    indices = [
        Index(value = ["updated_at"]),
        Index(value = ["user_id", "updated_at"]),
        Index(value = ["sync_status"])
    ]
)
data class NoteEntity(...)
```

---

## App Bundle vs APK

| | APK | Android App Bundle (AAB) |
|--|-----|--------------------------|
| Size | Full resources for all devices | Optimized per device |
| Upload | One APK | One AAB → Play Store generates APKs |
| Install size | Larger | ~15-50% smaller |
| Required for | Sideloading | Google Play publishing |

Always build AAB for Play Store releases:

```bash
./gradlew bundleRelease
```

---

## Pre-Launch Checklist

```
Before every release:
□ All unit tests pass
□ No lint errors
□ Proguard rules verified (test release build on a device)
□ Version code incremented
□ Version name updated
□ Release notes written
□ API keys not exposed
□ Debug logs disabled in release build type
□ Crashlytics is connected and reporting
□ Release build tested on minimum SDK version device
□ Deep links tested
□ Permission requests tested
□ Edge cases: empty state, error state, offline state
```

---

## Interview Questions

**Q1: How do you handle API keys securely in Android?**

> Store them in `local.properties` (gitignored) for development, inject them as `BuildConfig` fields via Gradle. In CI, inject via encrypted environment secrets. Never hardcode them in source files.

**Q2: What is the difference between Timber and `Log.d()`?**

> `Log.d()` always logs, including in production — it adds overhead and may expose sensitive info. Timber is a wrapper that lets you plant different trees per build type: a `DebugTree` in debug (logs everything) and a `CrashReportingTree` in release (logs only errors to Crashlytics). No production logs leak.

**Q3: How do you ensure no sensitive data is logged in production?**

> Use Timber and plant `Timber.DebugTree()` only in `BuildConfig.DEBUG` builds. In release, plant only `CrashReportingTree` which filters to ERROR/WARN level and redacts sensitive fields. Never use `Log` directly.

---

## Summary

- Use a structured error hierarchy (`AppError` sealed class) for predictable error handling
- Use Timber for logging — never `Log.d()` directly; configure per build type
- Integrate Firebase Crashlytics for production error tracking
- Store sensitive keys in `local.properties` (dev) and CI secrets (prod)
- Use `EncryptedSharedPreferences` for tokens and credentials
- Always run a pre-launch checklist before every release

**Next:** [Mini Project — E-Commerce Module](./mini-project-ecommerce-module.md)
