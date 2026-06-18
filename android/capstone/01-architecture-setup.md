# Capstone Part 1: Architecture Setup

## Step 1: Create the Project

1. Open Android Studio → New Project → Empty Views Activity
2. Name: `FeedFlow`, Package: `com.yourname.feedflow`
3. Language: Kotlin, Min SDK: API 24, Target: API 35

---

## Step 2: `settings.gradle.kts`

```kotlin
pluginManagement {
    includeBuild("build-logic")
    repositories {
        google()
        mavenCentral()
        gradlePluginPortal()
    }
}

dependencyResolutionManagement {
    repositoriesMode.set(RepositoriesMode.FAIL_ON_PROJECT_REPOS)
    repositories {
        google()
        mavenCentral()
    }
}

rootProject.name = "FeedFlow"

include(":app")
include(":core:common")
include(":core:domain")
include(":core:data")
include(":core:database")
include(":core:network")
include(":core:ui")
include(":core:testing")
include(":feature:headlines")
include(":feature:search")
include(":feature:bookmarks")
include(":feature:settings")
```

---

## Step 3: `gradle/libs.versions.toml`

```toml
[versions]
agp = "8.7.3"
kotlin = "2.1.0"
coreKtx = "1.15.0"
appcompat = "1.7.0"
material = "1.12.0"
constraintlayout = "2.2.0"
lifecycle = "2.8.7"
navigation = "2.8.4"
room = "2.6.1"
paging = "3.3.4"
hilt = "2.51.1"
hiltNavigation = "1.2.0"
retrofit = "2.11.0"
okhttp = "4.12.0"
coroutines = "1.8.1"
datastore = "1.1.1"
work = "2.9.1"
timber = "5.0.1"
coil = "2.7.0"
ksp = "2.1.0-1.0.29"
junit = "4.13.2"
junitExt = "1.2.1"
espresso = "3.6.1"
mockito = "5.4.0"
truth = "1.4.4"
coroutinesTest = "1.8.1"
leakcanary = "2.14"
compose = "1.7.6"
composeMaterial3 = "1.3.1"
activityCompose = "1.9.3"
startup = "1.1.1"

[libraries]
# AndroidX Core
androidx-core-ktx = { group = "androidx.core", name = "core-ktx", version.ref = "coreKtx" }
androidx-appcompat = { group = "androidx.appcompat", name = "appcompat", version.ref = "appcompat" }
material = { group = "com.google.android.material", name = "material", version.ref = "material" }
androidx-constraintlayout = { group = "androidx.constraintlayout", name = "constraintlayout", version.ref = "constraintlayout" }
androidx-startup = { group = "androidx.startup", name = "startup-runtime", version.ref = "startup" }

# Lifecycle
androidx-lifecycle-viewmodel-ktx = { group = "androidx.lifecycle", name = "lifecycle-viewmodel-ktx", version.ref = "lifecycle" }
androidx-lifecycle-runtime-ktx = { group = "androidx.lifecycle", name = "lifecycle-runtime-ktx", version.ref = "lifecycle" }
androidx-lifecycle-livedata-ktx = { group = "androidx.lifecycle", name = "lifecycle-livedata-ktx", version.ref = "lifecycle" }

# Navigation
androidx-navigation-fragment-ktx = { group = "androidx.navigation", name = "navigation-fragment-ktx", version.ref = "navigation" }
androidx-navigation-ui-ktx = { group = "androidx.navigation", name = "navigation-ui-ktx", version.ref = "navigation" }

# Room
androidx-room-runtime = { group = "androidx.room", name = "room-runtime", version.ref = "room" }
androidx-room-ktx = { group = "androidx.room", name = "room-ktx", version.ref = "room" }
androidx-room-paging = { group = "androidx.room", name = "room-paging", version.ref = "room" }
androidx-room-compiler = { group = "androidx.room", name = "room-compiler", version.ref = "room" }

# Paging
androidx-paging-runtime = { group = "androidx.paging", name = "paging-runtime-ktx", version.ref = "paging" }

# Hilt
hilt-android = { group = "com.google.dagger", name = "hilt-android", version.ref = "hilt" }
hilt-compiler = { group = "com.google.dagger", name = "hilt-android-compiler", version.ref = "hilt" }
hilt-navigation-fragment = { group = "androidx.hilt", name = "hilt-navigation-fragment", version.ref = "hiltNavigation" }
hilt-work = { group = "androidx.hilt", name = "hilt-work", version.ref = "hiltNavigation" }
hilt-compiler-androidx = { group = "androidx.hilt", name = "hilt-compiler", version.ref = "hiltNavigation" }

# Retrofit
retrofit = { group = "com.squareup.retrofit2", name = "retrofit", version.ref = "retrofit" }
retrofit-gson = { group = "com.squareup.retrofit2", name = "converter-gson", version.ref = "retrofit" }
okhttp-logging = { group = "com.squareup.okhttp3", name = "logging-interceptor", version.ref = "okhttp" }

# Coroutines
kotlinx-coroutines-android = { group = "org.jetbrains.kotlinx", name = "kotlinx-coroutines-android", version.ref = "coroutines" }
kotlinx-coroutines-core = { group = "org.jetbrains.kotlinx", name = "kotlinx-coroutines-core", version.ref = "coroutines" }

# DataStore
androidx-datastore-preferences = { group = "androidx.datastore", name = "datastore-preferences", version.ref = "datastore" }

# WorkManager
androidx-work-runtime-ktx = { group = "androidx.work", name = "work-runtime-ktx", version.ref = "work" }

# Logging
timber = { group = "com.jakewharton.timber", name = "timber", version.ref = "timber" }

# Image Loading
coil = { group = "io.coil-kt", name = "coil", version.ref = "coil" }

# Testing
junit = { group = "junit", name = "junit", version.ref = "junit" }
androidx-junit = { group = "androidx.test.ext", name = "junit", version.ref = "junitExt" }
androidx-espresso-core = { group = "androidx.test.espresso", name = "espresso-core", version.ref = "espresso" }
mockito-kotlin = { group = "org.mockito.kotlin", name = "mockito-kotlin", version.ref = "mockito" }
truth = { group = "com.google.truth", name = "truth", version.ref = "truth" }
kotlinx-coroutines-test = { group = "org.jetbrains.kotlinx", name = "kotlinx-coroutines-test", version.ref = "coroutinesTest" }
androidx-paging-testing = { group = "androidx.paging", name = "paging-testing", version.ref = "paging" }
hilt-android-testing = { group = "com.google.dagger", name = "hilt-android-testing", version.ref = "hilt" }
leakcanary = { group = "com.squareup.leakcanary", name = "leakcanary-android", version.ref = "leakcanary" }

[plugins]
android-application = { id = "com.android.application", version.ref = "agp" }
android-library = { id = "com.android.library", version.ref = "agp" }
kotlin-android = { id = "org.jetbrains.kotlin.android", version.ref = "kotlin" }
kotlin-jvm = { id = "org.jetbrains.kotlin.jvm", version.ref = "kotlin" }
hilt-android = { id = "com.google.dagger.hilt.android", version.ref = "hilt" }
navigation-safeargs = { id = "androidx.navigation.safeargs.kotlin", version.ref = "navigation" }
ksp = { id = "com.google.devtools.ksp", version.ref = "ksp" }
```

---

## Step 4: Application Class

`app/src/main/java/com/yourname/feedflow/FeedFlowApplication.kt`:

```kotlin
@HiltAndroidApp
class FeedFlowApplication : Application(), Configuration.Provider {

    @Inject
    lateinit var workerFactory: HiltWorkerFactory

    override val workManagerConfiguration: Configuration
        get() = Configuration.Builder()
            .setWorkerFactory(workerFactory)
            .build()

    override fun onCreate() {
        super.onCreate()
        setupTimber()
        setupWorkManager()
    }

    private fun setupTimber() {
        if (BuildConfig.DEBUG) {
            Timber.plant(Timber.DebugTree())
        }
        // Add Crashlytics tree in production
    }

    private fun setupWorkManager() {
        ArticleSyncWorker.schedulePeriodicSync(this)
    }
}
```

---

## Step 5: Domain Layer Models

`:core:domain/src/main/kotlin/com/yourname/feedflow/domain/model/`:

```kotlin
// Article.kt
data class Article(
    val id: String,  // url is used as unique ID
    val title: String,
    val description: String?,
    val url: String,
    val imageUrl: String?,
    val sourceName: String,
    val publishedAt: String,
    val content: String?,
    val isBookmarked: Boolean = false
)

// NewsCategory.kt
enum class NewsCategory(val displayName: String, val apiValue: String) {
    GENERAL("General", "general"),
    TECHNOLOGY("Technology", "technology"),
    BUSINESS("Business", "business"),
    SCIENCE("Science", "science"),
    HEALTH("Health", "health"),
    SPORTS("Sports", "sports"),
    ENTERTAINMENT("Entertainment", "entertainment")
}
```

`:core:domain/src/main/kotlin/com/yourname/feedflow/domain/repository/`:

```kotlin
// ArticleRepository.kt
interface ArticleRepository {
    fun getTopHeadlines(category: NewsCategory?): Flow<PagingData<Article>>
    fun searchArticles(query: String): Flow<PagingData<Article>>
    fun getBookmarkedArticles(): Flow<List<Article>>
    suspend fun getArticleById(id: String): Article?
    suspend fun toggleBookmark(article: Article)
    suspend fun syncArticles()
}
```

---

## Step 6: Base ViewModel

`:core:ui/src/main/kotlin/com/yourname/feedflow/core/ui/BaseViewModel.kt`:

```kotlin
abstract class BaseViewModel : ViewModel() {

    protected fun <T> Flow<T>.asStateFlowIn(
        scope: CoroutineScope = viewModelScope,
        initial: T
    ): StateFlow<T> = stateIn(scope, SharingStarted.WhileSubscribed(5_000), initial)

    protected fun launchSafe(
        onError: (Throwable) -> Unit = {},
        block: suspend CoroutineScope.() -> Unit
    ) = viewModelScope.launch {
        try {
            block()
        } catch (e: CancellationException) {
            throw e
        } catch (e: Exception) {
            Timber.e(e)
            onError(e)
        }
    }
}
```

---

**Next:** [Part 2 — UI Implementation](./02-ui-implementation.md)
