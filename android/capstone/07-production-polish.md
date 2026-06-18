# Capstone Part 7: Production Polish

## Background Sync with WorkManager

```kotlin
@HiltWorker
class ArticleSyncWorker @AssistedInject constructor(
    @Assisted context: Context,
    @Assisted params: WorkerParameters,
    private val repository: ArticleRepository
) : CoroutineWorker(context, params) {

    override suspend fun doWork(): Result {
        return try {
            setForeground(createForegroundInfo())
            repository.syncArticles()
            Result.success()
        } catch (e: IOException) {
            if (runAttemptCount < 3) Result.retry() else Result.failure()
        } catch (e: Exception) {
            Timber.e(e, "Sync failed")
            Result.failure()
        }
    }

    private fun createForegroundInfo(): ForegroundInfo {
        val notification = NotificationCompat.Builder(applicationContext, SYNC_CHANNEL_ID)
            .setContentTitle("Syncing FeedFlow...")
            .setSmallIcon(R.drawable.ic_sync)
            .setProgress(0, 0, true)
            .setSilent(true)
            .build()
        return ForegroundInfo(SYNC_NOTIFICATION_ID, notification)
    }

    companion object {
        const val SYNC_CHANNEL_ID = "feedflow_sync"
        const val SYNC_NOTIFICATION_ID = 1001

        fun schedulePeriodicSync(context: Context) {
            val request = PeriodicWorkRequestBuilder<ArticleSyncWorker>(
                repeatInterval = 1,
                repeatIntervalTimeUnit = TimeUnit.HOURS,
                flexTimeInterval = 15,
                flexTimeIntervalUnit = TimeUnit.MINUTES
            )
                .setConstraints(
                    Constraints.Builder()
                        .setRequiredNetworkType(NetworkType.CONNECTED)
                        .setRequiresBatteryNotLow(true)
                        .build()
                )
                .build()

            WorkManager.getInstance(context).enqueueUniquePeriodicWork(
                "article_sync",
                ExistingPeriodicWorkPolicy.KEEP,
                request
            )
        }
    }
}
```

---

## CI/CD Pipeline

`.github/workflows/ci.yml`:

```yaml
name: FeedFlow CI

on:
  push:
    branches: [main, develop]
  pull_request:
    branches: [main]

jobs:
  test-and-lint:
    name: Tests + Lint
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4

      - uses: actions/setup-java@v4
        with:
          java-version: '17'
          distribution: 'temurin'

      - uses: actions/cache@v4
        with:
          path: |
            ~/.gradle/caches
            ~/.gradle/wrapper
          key: gradle-${{ hashFiles('**/*.gradle*', '**/gradle-wrapper.properties') }}

      - name: Run unit tests
        run: ./gradlew testDebugUnitTest

      - name: Run lint
        run: ./gradlew lintDebug

      - name: Upload test results
        if: always()
        uses: actions/upload-artifact@v4
        with:
          name: test-reports
          path: '**/build/reports/'

  release:
    name: Release Build
    runs-on: ubuntu-latest
    needs: test-and-lint
    if: startsWith(github.ref, 'refs/tags/v')
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-java@v4
        with: { java-version: '17', distribution: 'temurin' }
      - uses: actions/cache@v4
        with:
          path: ~/.gradle
          key: gradle-${{ hashFiles('**/*.gradle*') }}

      - name: Decode keystore
        run: echo "${{ secrets.KEYSTORE_B64 }}" | base64 -d > app/feedflow.jks

      - name: Build release bundle
        run: ./gradlew bundleRelease
        env:
          KEYSTORE_PATH: feedflow.jks
          KEYSTORE_PASSWORD: ${{ secrets.KEYSTORE_PASSWORD }}
          KEY_ALIAS: ${{ secrets.KEY_ALIAS }}
          KEY_PASSWORD: ${{ secrets.KEY_PASSWORD }}
          NEWS_API_KEY: ${{ secrets.NEWS_API_KEY }}

      - name: Upload to Play Store internal track
        uses: r0adkll/upload-google-play@v1
        with:
          serviceAccountJsonPlainText: ${{ secrets.PLAY_SERVICE_ACCOUNT_JSON }}
          packageName: com.yourname.feedflow
          releaseFiles: app/build/outputs/bundle/release/*.aab
          track: internal
```

---

## ProGuard Rules

`app/proguard-rules.pro`:

```proguard
# Keep domain models
-keep class com.yourname.feedflow.core.domain.model.** { *; }

# Keep Room entities
-keep class com.yourname.feedflow.core.database.** { *; }

# Keep Retrofit DTOs
-keep class com.yourname.feedflow.core.network.dto.** { *; }
-keep interface com.yourname.feedflow.core.network.** { *; }

# Gson
-keepattributes Signature
-keepattributes *Annotation*
-dontwarn sun.misc.**
-keep class com.google.gson.** { *; }
-keep class * implements com.google.gson.TypeAdapterFactory

# Retrofit
-keepattributes RuntimeVisibleAnnotations, RuntimeVisibleParameterAnnotations
-keepclassmembers,allowshrinking,allowobfuscation interface * {
    @retrofit2.http.* <methods>;
}

# OkHttp
-dontwarn okhttp3.**
-dontwarn okio.**

# Hilt
-keepclassmembers,allowobfuscation class * {
    @javax.inject.* *;
    @dagger.* *;
    <init>();
}
-keep class dagger.hilt.** { *; }
-keep @dagger.hilt.android.HiltAndroidApp class *

# Coroutines
-keepnames class kotlinx.coroutines.internal.MainDispatcherFactory {}
-keepnames class kotlinx.coroutines.CoroutineExceptionHandler {}

# Navigation
-keepnames class androidx.navigation.fragment.NavHostFragment

# Parcelable / Serializable
-keepclassmembers class * implements android.os.Parcelable {
    public static final ** CREATOR;
}
```

---

## Final App Checklist

```
Architecture:
✓ Domain layer — pure Kotlin, zero Android imports
✓ Data layer — Room + Retrofit + RemoteMediator
✓ UI layer — MVVM + StateFlow + Navigation
✓ Modular structure with feature and core modules

Quality:
✓ Zero lint errors
✓ Zero Detekt violations
✓ Unit tests: ViewModel (>80%) + Repository (>80%)
✓ UI tests: critical paths (navigation, bookmark, error state)

Performance:
✓ Paging 3 with cachedIn for all lists
✓ setHasFixedSize(true) on all RecyclerViews
✓ Images loaded with Coil (caching + sampling)
✓ R8 enabled: isMinifyEnabled + isShrinkResources
✓ LeakCanary — no leaks in debug build

Production:
✓ Crashlytics integrated
✓ Timber logging — debug only
✓ API key via BuildConfig + CI secrets
✓ CI pipeline runs on every PR
✓ Release builds signed and uploaded to Play Store internal track
✓ Background sync via WorkManager
✓ Offline-first — RemoteMediator + Room cache

UX:
✓ Loading states (shimmer or progress indicator)
✓ Error states with retry button
✓ Empty states with meaningful messages
✓ Smooth navigation transitions
✓ Dark mode support
✓ Bookmark persistence across sessions
```

---

## Publishing to Play Store

1. **Create a keystore** (one time):

```bash
keytool -genkey -v -keystore feedflow.jks \
  -keyalg RSA -keysize 2048 -validity 10000 \
  -alias feedflow_key
```

2. **Encode for GitHub Secrets:**

```bash
base64 -i feedflow.jks | tr -d '\n' | pbcopy
```

3. **Add to GitHub Secrets:** `KEYSTORE_B64`, `KEYSTORE_PASSWORD`, `KEY_ALIAS`, `KEY_PASSWORD`

4. **Create Play Store listing** at [play.google.com/console](https://play.google.com/console)

5. **Push a tag to trigger release:**

```bash
git tag v1.0.0
git push origin v1.0.0
```

---

## Congratulations!

You have built FeedFlow — a production-quality Android app demonstrating:

| Skill | Implementation |
|-------|---------------|
| UI | Material 3, ConstraintLayout, RecyclerView |
| Architecture | Clean Arch, MVVM, modular |
| State management | StateFlow, SharedFlow |
| Persistence | Room with offline-first strategy |
| Networking | Retrofit + Paging 3 + RemoteMediator |
| DI | Hilt throughout |
| Background work | WorkManager periodic sync |
| Testing | Unit + integration + UI tests |
| Production | CI/CD, R8, Crashlytics, ProGuard |

**You are now production-ready.**

Push to GitHub, put it on your resume, and ship it to the Play Store.
