# Chapter 1: Hilt — Dependency Injection

## What Is Dependency Injection?

Dependency Injection (DI) is a design pattern where objects receive their dependencies from outside rather than creating them internally.

**Without DI (tight coupling):**
```kotlin
class NoteRepository {
    // Creates its own dependency — hard to test, hard to swap
    private val dao = NoteDatabase.getInstance(context).noteDao()
}
```

**With DI (loose coupling):**
```kotlin
class NoteRepository(private val dao: NoteDao) {
    // Receives its dependency — easy to test with a fake DAO
}
```

**Analogy:** Think of a restaurant kitchen. Without DI, every chef buys their own ingredients. With DI, there's a supplier (injector) who delivers exactly what each chef needs — the chef focuses on cooking, not sourcing.

---

## Why Hilt?

Hilt is Google's recommended DI library for Android, built on top of Dagger. It:
- Generates all DI boilerplate at compile time
- Integrates with `ViewModel`, `WorkManager`, `Navigation`
- Removes the need for `ViewModelFactory` boilerplate
- Scopes dependencies to Android components (Activity, Fragment, ViewModel)

```kotlin
// libs.versions.toml
[versions]
hilt = "2.51.1"

[libraries]
hilt-android = { group = "com.google.dagger", name = "hilt-android", version.ref = "hilt" }
hilt-compiler = { group = "com.google.dagger", name = "hilt-android-compiler", version.ref = "hilt" }
hilt-navigation-fragment = { group = "androidx.hilt", name = "hilt-navigation-fragment", version = "1.2.0" }

[plugins]
hilt-android = { id = "com.google.dagger.hilt.android", version.ref = "hilt" }
```

```kotlin
// app/build.gradle.kts
plugins {
    id("com.google.dagger.hilt.android")
    id("com.google.devtools.ksp")
}

dependencies {
    implementation(libs.hilt.android)
    ksp(libs.hilt.compiler)
    implementation(libs.hilt.navigation.fragment)
}
```

---

## Step 1: Annotate Application

Every Hilt app must have an `@HiltAndroidApp` annotation on the Application class:

```kotlin
@HiltAndroidApp
class MyApplication : Application()
```

Register in `AndroidManifest.xml`:
```xml
<application android:name=".MyApplication" ...>
```

---

## Step 2: Annotate Android Components

Activities, Fragments, Services that need injection:

```kotlin
@AndroidEntryPoint
class MainActivity : AppCompatActivity() {
    // Can now inject things here
}

@AndroidEntryPoint
class NoteListFragment : Fragment() {
    // Can now inject things here
}
```

---

## Step 3: Define How Dependencies Are Created — `@Module`

Hilt doesn't know how to create your dependencies. You teach it via modules:

```kotlin
@Module
@InstallIn(SingletonComponent::class)  // Lives as long as the app
object DatabaseModule {

    @Provides
    @Singleton  // Only one instance created ever
    fun provideDatabase(@ApplicationContext context: Context): NoteDatabase {
        return Room.databaseBuilder(context, NoteDatabase::class.java, "notes.db")
            .build()
    }

    @Provides
    @Singleton
    fun provideNoteDao(database: NoteDatabase): NoteDao {
        return database.noteDao()
    }
}

@Module
@InstallIn(SingletonComponent::class)
object RepositoryModule {

    @Provides
    @Singleton
    fun provideNoteRepository(dao: NoteDao): NoteRepository {
        return NoteRepository(dao)
    }
}
```

---

## Step 4: Inject into ViewModel

With Hilt, no `ViewModelFactory` needed:

```kotlin
@HiltViewModel
class NoteViewModel @Inject constructor(
    private val repository: NoteRepository
) : ViewModel() {

    val notes = repository.getAllNotes()
        .stateIn(viewModelScope, SharingStarted.WhileSubscribed(5_000), emptyList())

    fun saveNote(title: String, content: String) {
        viewModelScope.launch {
            repository.saveNote(Note(title = title, content = content))
        }
    }
}
```

Use in Fragment — same as before, no factory needed:

```kotlin
@AndroidEntryPoint
class NoteListFragment : Fragment() {

    // Hilt resolves this automatically
    private val viewModel: NoteViewModel by viewModels()
}
```

---

## Step 5: `@Inject` Constructor — Simplest Way

For classes you own, add `@Inject constructor`:

```kotlin
class NoteRepository @Inject constructor(
    private val dao: NoteDao
) {
    fun getAllNotes() = dao.getAllNotes().map { list -> list.map { it.toDomain() } }
    suspend fun saveNote(note: Note) = dao.insertNote(note.toEntity())
    suspend fun deleteNote(note: Note) = dao.deleteNote(note.toEntity())
}
```

Now Hilt knows how to create `NoteRepository` — no module needed for this class.

---

## Hilt Scopes

| Scope Annotation | Component | Lifetime |
|------------------|-----------|----------|
| `@Singleton` | `SingletonComponent` | App lifetime |
| `@ActivityScoped` | `ActivityComponent` | Activity lifetime |
| `@ViewModelScoped` | `ViewModelComponent` | ViewModel lifetime |
| `@FragmentScoped` | `FragmentComponent` | Fragment lifetime |

```kotlin
@Module
@InstallIn(ViewModelComponent::class)
object ViewModelModule {
    @Provides
    @ViewModelScoped
    fun provideAnalyticsTracker(): AnalyticsTracker = AnalyticsTracker()
}
```

---

## Interfaces and Bindings

When a dependency is an interface, use `@Binds`:

```kotlin
interface NoteRepository {
    fun getAllNotes(): Flow<List<Note>>
    suspend fun saveNote(note: Note): Long
}

class NoteRepositoryImpl @Inject constructor(
    private val dao: NoteDao
) : NoteRepository {
    override fun getAllNotes() = dao.getAllNotes().map { it.map { e -> e.toDomain() } }
    override suspend fun saveNote(note: Note) = dao.insertNote(note.toEntity())
}

@Module
@InstallIn(SingletonComponent::class)
abstract class RepositoryModule {

    // Binds the interface to its implementation
    @Binds
    @Singleton
    abstract fun bindNoteRepository(impl: NoteRepositoryImpl): NoteRepository
}
```

---

## Multiple Implementations — `@Qualifier`

When two dependencies have the same type, use qualifiers to distinguish:

```kotlin
@Qualifier
@Retention(AnnotationRetention.BINARY)
annotation class IoDispatcher

@Qualifier
@Retention(AnnotationRetention.BINARY)
annotation class MainDispatcher

@Module
@InstallIn(SingletonComponent::class)
object CoroutinesModule {

    @Provides
    @Singleton
    @IoDispatcher
    fun provideIoDispatcher(): CoroutineDispatcher = Dispatchers.IO

    @Provides
    @Singleton
    @MainDispatcher
    fun provideMainDispatcher(): CoroutineDispatcher = Dispatchers.Main
}

// Usage
class NoteRepository @Inject constructor(
    private val dao: NoteDao,
    @IoDispatcher private val ioDispatcher: CoroutineDispatcher
) {
    suspend fun loadNotes() = withContext(ioDispatcher) {
        dao.getAllNotes()
    }
}
```

---

## Injecting into Non-Android Classes

For classes that aren't Activities/Fragments/ViewModels (e.g., a custom manager):

```kotlin
class NotificationHelper @Inject constructor(
    @ApplicationContext private val context: Context
) {
    fun showNotification(title: String, body: String) { /* ... */ }
}

// Used in ViewModel
@HiltViewModel
class NoteViewModel @Inject constructor(
    private val repository: NoteRepository,
    private val notificationHelper: NotificationHelper
) : ViewModel()
```

---

## Full Dependency Graph (Notes App)

```
Application
    │
    ├── NoteDatabase (Singleton)
    │       └── NoteDao (Singleton)
    │               └── NoteRepository (Singleton)
    │                       └── NoteViewModel (@HiltViewModel)
    │                               └── NoteListFragment (@AndroidEntryPoint)
    │
    └── CoroutineDispatcher (Singleton, @IoDispatcher)
            └── NoteRepository (Singleton)
```

---

## Testing with Hilt

Replace real modules with test doubles:

```kotlin
@HiltAndroidTest
class NoteRepositoryTest {

    @get:Rule
    val hiltRule = HiltAndroidRule(this)

    @Inject
    lateinit var repository: NoteRepository

    @Before
    fun setup() { hiltRule.inject() }

    @Test
    fun insertAndRetrieve() = runTest {
        repository.saveNote(Note(title = "Test", content = "Body"))
        val notes = repository.getAllNotes().first()
        assertThat(notes).hasSize(1)
        assertThat(notes[0].title).isEqualTo("Test")
    }
}

// Replace the real database with in-memory for tests
@Module
@TestInstallIn(
    components = [SingletonComponent::class],
    replaces = [DatabaseModule::class]
)
object TestDatabaseModule {
    @Provides
    @Singleton
    fun provideDatabase(@ApplicationContext context: Context): NoteDatabase {
        return Room.inMemoryDatabaseBuilder(context, NoteDatabase::class.java).build()
    }
}
```

---

## Common Mistakes

### Mistake 1: Forgetting `@HiltAndroidApp` on Application

```
Error: [Hilt] @HiltAndroidApp annotation is required on Application
```

### Mistake 2: Using `@Inject` on a class you don't own

For classes you can't add `@Inject constructor` to (Retrofit, OkHttp, Room), use `@Provides` in a `@Module`.

### Mistake 3: Wrong scope — creating expensive objects per-Activity

```kotlin
// WRONG — new database per Activity = catastrophic
@Provides
@ActivityScoped
fun provideDatabase(...): NoteDatabase

// CORRECT — one database for the app
@Provides
@Singleton
fun provideDatabase(...): NoteDatabase
```

---

## Interview Questions

**Q1: What is the difference between `@Inject constructor` and `@Provides`?**

> `@Inject constructor` is used when you own the class — Hilt reads it to know how to create the object. `@Provides` is used in `@Module` for classes you don't own (e.g., Retrofit, Room Database) where you need to write creation logic manually.

**Q2: What does `@InstallIn(SingletonComponent::class)` mean?**

> It tells Hilt which component (and thus which scope) the module belongs to. `SingletonComponent` means the provided bindings live for the entire application lifetime — one instance per app.

**Q3: What is the difference between `@Binds` and `@Provides`?**

> `@Binds` is used in abstract modules to tell Hilt which concrete implementation to use for an interface — it's more efficient since no code is generated. `@Provides` is used in object modules with a function body that creates and returns the dependency.

---

## Summary

- Hilt is Android's recommended DI library — removes `ViewModelFactory` boilerplate
- `@HiltAndroidApp` on Application; `@AndroidEntryPoint` on Activity/Fragment
- `@HiltViewModel` + `@Inject constructor` on ViewModel
- `@Module` + `@Provides` for classes you don't own; `@Inject constructor` for classes you own
- Use scopes (`@Singleton`, `@ViewModelScoped`) to control object lifetime
- Replace modules in tests with `@TestInstallIn`

**Next:** [Chapter 2 — WorkManager](./02-workmanager.md)
