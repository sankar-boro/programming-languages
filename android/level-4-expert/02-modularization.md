# Chapter 2: Modularization

## What Is Modularization?

Modularization splits a monolithic Android app into separate Gradle modules, each with a single responsibility.

**Analogy:** A monolith is one giant kitchen where all chefs share one stove. Modular is separate kitchens — each team owns their space, and the main kitchen just coordinates.

---

## Why Modularize?

| Benefit | What It Means |
|---------|--------------|
| Build speed | Gradle only rebuilds changed modules |
| Parallel builds | Modules without dependencies build simultaneously |
| Team scalability | Teams own separate modules with clear contracts |
| Reusability | Share a `:core:ui` module between multiple features |
| Strict boundaries | Modules can only use what's in their dependency graph |
| Testing isolation | Test each module in isolation |

---

## Module Types

```
:app                    ← App module — wires everything together
:feature:notes          ← Notes feature
:feature:settings       ← Settings feature
:feature:search         ← Search feature
:core:ui                ← Shared UI components, themes
:core:data              ← Repository implementations
:core:domain            ← Use cases, interfaces, domain models
:core:network           ← Retrofit setup
:core:database          ← Room setup
:core:testing           ← Shared test utilities
```

---

## Creating a Module

1. Right-click the project root → **New → Module**
2. Select **Android Library** (for modules that use Android APIs) or **Java or Kotlin Library** (for pure Kotlin)
3. Name it following your convention: `feature-notes` or `core-ui`

---

## Module-Level `build.gradle.kts`

```kotlin
// :core:domain — pure Kotlin module
plugins {
    id("java-library")
    alias(libs.plugins.kotlin.jvm)
}

dependencies {
    implementation(libs.kotlinx.coroutines.core)
}
```

```kotlin
// :feature:notes — Android module
plugins {
    alias(libs.plugins.android.library)
    alias(libs.plugins.kotlin.android)
    alias(libs.plugins.hilt.android)
    id("com.google.devtools.ksp")
}

android {
    namespace = "com.yourname.feature.notes"
    compileSdk = 35
    defaultConfig { minSdk = 24 }
}

dependencies {
    implementation(project(":core:domain"))
    implementation(project(":core:ui"))
    implementation(libs.hilt.android)
    ksp(libs.hilt.compiler)
}
```

---

## `settings.gradle.kts` — Declare All Modules

```kotlin
rootProject.name = "MyApp"

include(":app")

// Core modules
include(":core:domain")
include(":core:data")
include(":core:database")
include(":core:network")
include(":core:ui")
include(":core:testing")

// Feature modules
include(":feature:notes")
include(":feature:settings")
include(":feature:search")
```

---

## Dependency Graph

```
:app
 ├── :feature:notes
 │    ├── :core:domain
 │    └── :core:ui
 ├── :feature:settings
 │    ├── :core:domain
 │    └── :core:ui
 ├── :core:data   (wires :core:domain interfaces to implementations)
 │    ├── :core:domain
 │    ├── :core:database
 │    └── :core:network
 └── :core:ui
```

**Key rule:** Feature modules NEVER depend on other feature modules directly. They communicate through `:core:domain` interfaces or the `:app` module.

---

## Convention Plugins — Avoiding Duplicate Build Config

Instead of copy-pasting the same Gradle configuration across 10 modules, use **convention plugins**:

`build-logic/convention/build.gradle.kts`:

```kotlin
plugins {
    `kotlin-dsl`
}

dependencies {
    compileOnly(libs.android.gradlePlugin)
    compileOnly(libs.kotlin.gradlePlugin)
    compileOnly(libs.ksp.gradlePlugin)
}
```

`build-logic/convention/src/main/kotlin/AndroidLibraryConventionPlugin.kt`:

```kotlin
class AndroidLibraryConventionPlugin : Plugin<Project> {
    override fun apply(target: Project) {
        with(target) {
            pluginManager.apply("com.android.library")
            pluginManager.apply("org.jetbrains.kotlin.android")

            extensions.configure<LibraryExtension> {
                compileSdk = 35
                defaultConfig {
                    minSdk = 24
                    testInstrumentationRunner = "androidx.test.runner.AndroidJUnitRunner"
                }
                compileOptions {
                    sourceCompatibility = JavaVersion.VERSION_11
                    targetCompatibility = JavaVersion.VERSION_11
                }
            }
        }
    }
}
```

Now each library module just uses:

```kotlin
plugins {
    id("yourapp.android.library")  // applies everything above
    id("yourapp.android.hilt")     // adds Hilt setup
}
```

---

## Module Navigation

With Navigation Component, each feature module can contribute its own navigation graph:

`:feature:notes` → `res/navigation/notes_nav_graph.xml`:

```xml
<navigation
    android:id="@+id/notes_nav_graph"
    app:startDestination="@id/noteListFragment">

    <fragment android:id="@+id/noteListFragment"
        android:name="com.yourname.feature.notes.NoteListFragment">
        <action android:id="@+id/action_to_detail"
            app:destination="@id/noteDetailFragment" />
    </fragment>

    <fragment android:id="@+id/noteDetailFragment"
        android:name="com.yourname.feature.notes.NoteDetailFragment" />
</navigation>
```

`:app` → `res/navigation/main_nav_graph.xml`:

```xml
<navigation
    android:id="@+id/main_nav_graph"
    app:startDestination="@id/notes_nav_graph">

    <!-- Include feature nav graphs -->
    <include app:graph="@navigation/notes_nav_graph" />
    <include app:graph="@navigation/settings_nav_graph" />

</navigation>
```

---

## Hilt in Multi-Module Setup

Each module can have its own Hilt modules. The `:app` module acts as the root:

```kotlin
// :core:database — provides Room
@Module
@InstallIn(SingletonComponent::class)
object DatabaseModule {
    @Provides @Singleton
    fun provideDatabase(@ApplicationContext ctx: Context): AppDatabase = ...
}

// :core:data — provides Repository implementation
@Module
@InstallIn(SingletonComponent::class)
abstract class RepositoryModule {
    @Binds @Singleton
    abstract fun bindNoteRepository(impl: NoteRepositoryImpl): NoteRepository
}

// :app — Application class ties it all together
@HiltAndroidApp
class MyApplication : Application()
```

Hilt automatically discovers all `@Module`s in the dependency graph.

---

## `api` vs `implementation` Dependencies

```kotlin
dependencies {
    // 'api' — exposes the dependency to modules that depend on this module
    api(project(":core:domain"))

    // 'implementation' — NOT exposed; internal to this module only
    implementation(project(":core:database"))
    implementation(libs.retrofit)
}
```

Use `api` only for types that appear in the module's public API. Prefer `implementation` everywhere else.

---

## Build Performance Tips

```kotlin
// gradle.properties
org.gradle.parallel=true            // Build modules in parallel
org.gradle.caching=true             // Cache build outputs
org.gradle.configuration-cache=true // Cache configuration phase
android.enableBuildCache=true
```

---

## Common Mistakes

### Mistake 1: Feature-to-feature dependencies

```kotlin
// WRONG — :feature:search depends on :feature:notes
dependencies {
    implementation(project(":feature:notes"))
}

// CORRECT — both depend on :core:domain interfaces
dependencies {
    implementation(project(":core:domain"))
}
```

### Mistake 2: Putting business logic in `:app`

The `:app` module should ONLY wire things together (Hilt setup, navigation graph, Application class). No business logic.

### Mistake 3: Over-modularizing early

Don't split into 15 modules on day one. Start with `:app`, `:core:domain`, `:core:data`, and one or two feature modules. Modularize as the app grows.

---

## Interview Questions

**Q1: What is the main benefit of modularization for build speed?**

> Gradle only rebuilds modules that changed. In a monolith, changing one file triggers recompilation of the entire app. With modules, only the changed module and its dependents rebuild — CI build times can drop from 10 minutes to 2 minutes.

**Q2: What is the dependency rule in a modular app?**

> Feature modules depend only on `:core:domain` (interfaces) and `:core:ui` (shared UI). They never depend on each other. The `:app` module wires feature modules together but contains no business logic. The `:core:data` module implements `:core:domain` interfaces.

**Q3: What is a convention plugin?**

> A convention plugin is a Gradle plugin that bundles common build configuration (compileSdk, minSdk, Kotlin options, common dependencies). Instead of duplicating 30 lines of Gradle config in every module's `build.gradle.kts`, each module just applies `id("yourapp.android.library")`.

---

## Summary

- Modularization improves build speed, team scalability, and code boundaries
- Feature modules depend on `:core:domain`, not on each other
- Use `api` for public dependencies, `implementation` for internal ones
- Convention plugins eliminate Gradle config duplication across modules
- Don't over-modularize early — start with 3–5 modules and grow from there

**Next:** [Chapter 3 — Offline-First Apps](./03-offline-first-apps.md)
