# Chapter 5: Multi-Module Apps — Patterns and Practices

## Moving from Modularization Concepts to Implementation

Chapter 2 explained WHY to modularize. This chapter shows HOW to structure a real multi-module app with all the patterns you need for production.

---

## Complete Module Structure (Production Template)

```
MyApp/
├── app/                         ← Entry point: wires everything
├── build-logic/                 ← Convention plugins
│   └── convention/
│       └── src/main/kotlin/
│           ├── AndroidLibraryConventionPlugin.kt
│           ├── AndroidLibraryHiltConventionPlugin.kt
│           └── AndroidFeatureConventionPlugin.kt
│
├── core/
│   ├── common/                  ← Kotlin-only utilities (Result, extensions)
│   ├── domain/                  ← Models, interfaces, use cases
│   ├── data/                    ← Repository implementations
│   ├── database/                ← Room database, DAOs, entities
│   ├── network/                 ← Retrofit, OkHttp, DTOs
│   ├── ui/                      ← Shared composables, themes, resources
│   ├── testing/                 ← Shared test fakes, rules, builders
│   └── datastore/               ← DataStore setup and preferences
│
└── feature/
    ├── notes/                   ← Notes feature
    ├── search/                  ← Search feature
    ├── settings/                ← Settings feature
    └── onboarding/              ← Onboarding flow
```

---

## Convention Plugin Implementation

`build-logic/convention/src/main/kotlin/AndroidFeatureConventionPlugin.kt`:

```kotlin
class AndroidFeatureConventionPlugin : Plugin<Project> {
    override fun apply(target: Project) {
        with(target) {
            pluginManager.apply {
                apply("com.android.library")
                apply("org.jetbrains.kotlin.android")
                apply("com.google.dagger.hilt.android")
                apply("com.google.devtools.ksp")
            }

            extensions.configure<LibraryExtension> {
                compileSdk = 35
                defaultConfig {
                    minSdk = 24
                    testInstrumentationRunner = "com.yourname.core.testing.HiltTestRunner"
                }
                buildFeatures {
                    viewBinding = true
                    compose = true
                }
                composeOptions {
                    kotlinCompilerExtensionVersion = "1.5.8"
                }
            }

            dependencies {
                // Every feature gets these automatically
                add("implementation", project(":core:domain"))
                add("implementation", project(":core:ui"))
                add("implementation", "com.google.dagger:hilt-android:2.51.1")
                add("ksp", "com.google.dagger:hilt-android-compiler:2.51.1")
                add("implementation", "androidx.hilt:hilt-navigation-fragment:1.2.0")
                add("implementation", "androidx.lifecycle:lifecycle-viewmodel-ktx:2.8.7")
                add("implementation", "androidx.lifecycle:lifecycle-runtime-ktx:2.8.7")
                add("testImplementation", project(":core:testing"))
                add("androidTestImplementation", project(":core:testing"))
            }
        }
    }
}
```

Usage in each feature module:

```kotlin
// :feature:notes/build.gradle.kts
plugins {
    id("yourapp.android.feature")  // That's it — everything above is included
}

android {
    namespace = "com.yourname.feature.notes"
}

dependencies {
    // Only module-specific additions
    implementation(project(":core:database"))
}
```

---

## `:core:common` — Shared Utilities

```kotlin
// core/common/src/main/kotlin/result/Result.kt
sealed interface Result<out T> {
    data class Success<T>(val data: T) : Result<T>
    data class Error(val exception: Throwable, val message: String? = null) : Result<Nothing>
    data object Loading : Result<Nothing>
}

fun <T> Result<T>.onSuccess(block: (T) -> Unit): Result<T> {
    if (this is Result.Success) block(data)
    return this
}

fun <T> Result<T>.onError(block: (Throwable, String?) -> Unit): Result<T> {
    if (this is Result.Error) block(exception, message)
    return this
}

// core/common/src/main/kotlin/dispatcher/DispatcherProvider.kt
interface DispatcherProvider {
    val main: CoroutineDispatcher
    val io: CoroutineDispatcher
    val default: CoroutineDispatcher
}

class DefaultDispatcherProvider @Inject constructor() : DispatcherProvider {
    override val main: CoroutineDispatcher = Dispatchers.Main
    override val io: CoroutineDispatcher = Dispatchers.IO
    override val default: CoroutineDispatcher = Dispatchers.Default
}

// Test double
class TestDispatcherProvider : DispatcherProvider {
    val testDispatcher = StandardTestDispatcher()
    override val main: CoroutineDispatcher = testDispatcher
    override val io: CoroutineDispatcher = testDispatcher
    override val default: CoroutineDispatcher = testDispatcher
}
```

---

## `:core:testing` — Shared Test Utilities

```kotlin
// core/testing/src/main/kotlin/FakeNoteRepository.kt
class FakeNoteRepository : NoteRepository {
    private val _notes = MutableStateFlow<List<Note>>(emptyList())

    var shouldThrowError = false
    var networkDelay = 0L

    fun setNotes(notes: List<Note>) { _notes.value = notes }

    override fun getAllNotes(): Flow<List<Note>> = _notes

    override suspend fun saveNote(note: Note): Long {
        if (shouldThrowError) throw IOException("Network error")
        delay(networkDelay)
        val updated = _notes.value + note.copy(id = (_notes.value.size + 1).toLong())
        _notes.value = updated
        return updated.last().id
    }

    override suspend fun deleteNote(note: Note) {
        _notes.value = _notes.value.filter { it.id != note.id }
    }
}

// core/testing/src/main/kotlin/MainDispatcherRule.kt
class MainDispatcherRule(
    val testDispatcher: TestDispatcher = UnconfinedTestDispatcher()
) : TestWatcher() {
    override fun starting(description: Description) = Dispatchers.setMain(testDispatcher)
    override fun finished(description: Description) = Dispatchers.resetMain()
}
```

---

## Navigation in Multi-Module Apps

### Approach 1: Deep Links (Simple)

Each feature module exposes destinations via deep links in its nav graph. The `:app` module includes all graphs:

```xml
<!-- :app main_nav.xml -->
<navigation>
    <include app:graph="@navigation/notes_nav" />
    <include app:graph="@navigation/settings_nav" />
</navigation>

<!-- Navigate between features using deep links -->
findNavController().navigate("yourapp://notes/123".toUri())
```

### Approach 2: Navigation Interface (Decoupled)

Define navigation contracts in `:core:domain`:

```kotlin
// :core:domain
interface NoteNavigator {
    fun openNoteDetail(navController: NavController, noteId: Long)
    fun openSettings(navController: NavController)
}

// :app — implements the interface
class AppNoteNavigator @Inject constructor() : NoteNavigator {
    override fun openNoteDetail(navController: NavController, noteId: Long) {
        navController.navigate(
            MainNavGraphDirections.actionToNoteDetail(noteId)
        )
    }
    override fun openSettings(navController: NavController) {
        navController.navigate(R.id.settings_nav_graph)
    }
}

// :feature:notes — uses the interface (no knowledge of other modules)
@HiltViewModel
class NoteListViewModel @Inject constructor(
    private val navigator: NoteNavigator,
    private val getNotesUseCase: GetNotesUseCase
) : ViewModel()
```

---

## Shared ViewModel Across Modules

```kotlin
// :core:domain — the shared state
data class UserSession(val userId: String, val displayName: String)

// :core:ui — the shared ViewModel
@HiltViewModel
class UserSessionViewModel @Inject constructor(
    private val sessionRepository: SessionRepository
) : ViewModel() {
    val session: StateFlow<UserSession?> = sessionRepository.currentSession
        .stateIn(viewModelScope, SharingStarted.WhileSubscribed(5_000), null)
}

// In any feature Fragment — access via activityViewModels
@AndroidEntryPoint
class NoteListFragment : Fragment() {
    private val sessionViewModel: UserSessionViewModel by activityViewModels()

    override fun onViewCreated(view: View, savedInstanceState: Bundle?) {
        viewLifecycleOwner.lifecycleScope.launch {
            viewLifecycleOwner.repeatOnLifecycle(Lifecycle.State.STARTED) {
                sessionViewModel.session.collect { session ->
                    binding.tvUserName.text = session?.displayName ?: "Guest"
                }
            }
        }
    }
}
```

---

## Dynamic Feature Modules

For large apps where features are downloaded on demand:

```kotlin
// build.gradle.kts for dynamic feature
plugins {
    id("com.android.dynamic-feature")
    alias(libs.plugins.kotlin.android)
}

android {
    namespace = "com.yourname.feature.premium"
}

// :app must declare the dynamic feature
android {
    dynamicFeatures += ":feature:premium"
}

// Install on demand
val splitInstallManager = SplitInstallManagerFactory.create(context)
val request = SplitInstallRequest.newBuilder()
    .addModule("premium")
    .build()

splitInstallManager.startInstall(request)
    .addOnSuccessListener { /* module installed */ }
    .addOnFailureListener { /* installation failed */ }
```

---

## Interview Questions

**Q1: How do you share a ViewModel between multiple feature modules?**

> Define the ViewModel in a shared module (e.g., `:core:ui` or `:core:domain`), inject it with `@HiltViewModel`, and access it via `activityViewModels()` in each Fragment. Both Fragments get the same instance since they share the same Activity.

**Q2: How do you navigate between features in a multi-module app without creating dependencies between feature modules?**

> Option 1: Use deep links — each feature exposes a `Uri` pattern, and features navigate by deep link. Option 2: Define a navigator interface in `:core:domain` and implement it in `:app`. Feature modules use the interface without knowing about other modules.

**Q3: What is a Dynamic Feature Module?**

> A module that is not bundled in the initial APK install but can be downloaded on demand. Useful for rarely-used premium features — keeps the base APK small and downloads the feature only when the user requests it.

---

## Summary

- Convention plugins eliminate boilerplate from each module's `build.gradle.kts`
- `:core:common` provides shared utilities; `:core:testing` provides shared test infrastructure
- Features communicate via `:core:domain` interfaces or deep links — never direct module dependencies
- Use `activityViewModels()` for cross-fragment shared ViewModels
- Dynamic Feature Modules enable on-demand feature delivery

**Next:** [Chapter 6 — CI/CD Basics](./06-cicd-basics.md)
