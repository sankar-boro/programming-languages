# Chapter 3: Coroutines and AndroidX

## What Are Coroutines?

Coroutines are Kotlin's way of writing asynchronous, non-blocking code in a sequential style. Instead of callbacks or RxJava chains, you write code that looks synchronous but executes asynchronously.

```kotlin
// Callback hell (old way)
api.fetchNotes(userId, object : Callback<List<Note>> {
    override fun onSuccess(notes: List<Note>) {
        db.saveNotes(notes, object : Callback<Unit> {
            override fun onSuccess(result: Unit) {
                runOnUiThread { updateUI(notes) }
            }
        })
    }
    override fun onFailure(error: Throwable) { showError(error) }
})

// Coroutine (modern way)
viewModelScope.launch {
    val notes = api.fetchNotes(userId)    // suspends, doesn't block
    db.saveNotes(notes)                   // suspends, doesn't block
    updateUI(notes)                        // back on main thread
}
```

```kotlin
// Dependencies
implementation("org.jetbrains.kotlinx:kotlinx-coroutines-android:1.8.1")
implementation("org.jetbrains.kotlinx:kotlinx-coroutines-core:1.8.1")
```

---

## Key Concepts

| Concept | What It Is |
|---------|-----------|
| `suspend` | Marks a function that can be paused without blocking the thread |
| `CoroutineScope` | Defines the lifetime of coroutines launched within it |
| `CoroutineContext` | Carries dispatcher and job; controls how/where coroutines run |
| `Dispatcher` | The thread pool where the coroutine runs |
| `Job` | A handle to a running coroutine; can be cancelled |

---

## Dispatchers

```kotlin
// Main — UI thread; use for UI updates, ViewModels
Dispatchers.Main

// IO — optimized for I/O: network, disk, database
Dispatchers.IO

// Default — CPU-bound work: sorting, parsing large data
Dispatchers.Default

// Unconfined — starts on calling thread; avoid in production
Dispatchers.Unconfined
```

---

## `viewModelScope` — The Most Important Scope in Android

```kotlin
class NoteViewModel @Inject constructor(
    private val repository: NoteRepository
) : ViewModel() {

    fun loadNotes() {
        viewModelScope.launch {
            // Launched on Main by default
            val notes = withContext(Dispatchers.IO) {
                // Switched to IO for the database call
                repository.getAllNotesList()
            }
            // Back on Main automatically
            _uiState.update { it.copy(notes = notes) }
        }
    }
}
```

`viewModelScope` is cancelled automatically in `ViewModel.onCleared()`. No leaks.

---

## `lifecycleScope` — UI-Scoped Coroutines

```kotlin
class NoteListFragment : Fragment() {

    override fun onViewCreated(view: View, savedInstanceState: Bundle?) {
        super.onViewCreated(view, savedInstanceState)

        // Cancelled when Fragment view is DESTROYED
        viewLifecycleOwner.lifecycleScope.launch {
            viewLifecycleOwner.repeatOnLifecycle(Lifecycle.State.STARTED) {
                viewModel.uiState.collect { state ->
                    renderState(state)
                }
            }
        }
    }
}
```

---

## Structured Concurrency

Structured concurrency means coroutines form a parent-child hierarchy. If a parent is cancelled, all children are cancelled:

```kotlin
viewModelScope.launch {
    // These two run concurrently, both are children of this scope
    val notesDeferred = async { repository.getNotes() }
    val categoriesDeferred = async { repository.getCategories() }

    // Wait for both
    val notes = notesDeferred.await()
    val categories = categoriesDeferred.await()

    _uiState.update { it.copy(notes = notes, categories = categories) }
}
// If viewModelScope is cancelled, both async blocks are cancelled too
```

---

## `launch` vs `async`

| | `launch` | `async` |
|--|----------|---------|
| Returns | `Job` | `Deferred<T>` |
| Result | Fire-and-forget | Produces a value |
| Use | Side effects | Parallel computation |

```kotlin
// launch — for work with no return value
viewModelScope.launch {
    repository.deleteNote(note)
}

// async — for parallel work that produces a value
viewModelScope.launch {
    val result = async { repository.getNotes() }
    val count = async { repository.getNoteCount() }

    val notes = result.await()
    val total = count.await()
}
```

---

## Exception Handling

```kotlin
// Option 1: try-catch
viewModelScope.launch {
    try {
        val notes = repository.fetchFromNetwork()
        _uiState.update { it.copy(notes = notes) }
    } catch (e: IOException) {
        _uiState.update { it.copy(error = "Network error: ${e.message}") }
    } catch (e: CancellationException) {
        throw e  // Always re-throw CancellationException
    }
}

// Option 2: CoroutineExceptionHandler
val exceptionHandler = CoroutineExceptionHandler { _, throwable ->
    _uiState.update { it.copy(error = throwable.message) }
}

viewModelScope.launch(exceptionHandler) {
    val notes = repository.fetchFromNetwork()
    _uiState.update { it.copy(notes = notes) }
}
```

> Always re-throw `CancellationException`. Catching and swallowing it breaks structured concurrency.

---

## `withContext` — Switching Dispatchers

```kotlin
suspend fun processLargeDataset(data: List<RawData>): List<ProcessedData> {
    return withContext(Dispatchers.Default) {
        // CPU-intensive work on Default dispatcher
        data.map { it.process() }
    }
    // Automatically returns to the calling dispatcher
}

// In ViewModel (Main dispatcher)
viewModelScope.launch {
    val processed = processLargeDataset(rawData)  // switches to Default, then back
    _uiState.update { it.copy(data = processed) }  // runs on Main
}
```

---

## Flow — Coroutine-Based Streams

`Flow` is a cold, asynchronous data stream — it emits multiple values over time.

```kotlin
// Cold flow — starts when collected
fun getNotesStream(): Flow<List<Note>> = flow {
    while (true) {
        emit(repository.getNotes())
        delay(30_000)  // Refresh every 30 seconds
    }
}

// Transform flows
val filteredNotes: Flow<List<Note>> = repository.getAllNotes()
    .map { notes -> notes.filter { !it.isDeleted } }
    .catch { e -> emit(emptyList()) }  // Error handling in Flow
    .onEach { notes -> log("Loaded ${notes.size} notes") }

// Collect in ViewModel
viewModelScope.launch {
    filteredNotes.collect { notes ->
        _uiState.update { it.copy(notes = notes) }
    }
}
```

---

## Room + Coroutines Integration

Room natively supports coroutines and Flow:

```kotlin
@Dao
interface NoteDao {
    // Suspend function — use for one-shot queries
    @Insert
    suspend fun insertNote(note: NoteEntity): Long

    // Flow — reactive, re-emits on table changes
    @Query("SELECT * FROM notes")
    fun getAllNotes(): Flow<List<NoteEntity>>
}

// In Repository — just map and return
class NoteRepository @Inject constructor(private val dao: NoteDao) {
    val notes: Flow<List<Note>> = dao.getAllNotes()
        .map { entities -> entities.map { it.toDomain() } }

    suspend fun insertNote(note: Note) = dao.insertNote(note.toEntity())
}
```

---

## Retrofit + Coroutines Integration

```kotlin
interface NotesApiService {
    @GET("notes")
    suspend fun fetchNotes(@Query("userId") userId: String): List<NoteDto>

    @POST("notes")
    suspend fun createNote(@Body note: NoteDto): NoteDto
}

// Repository usage
class NoteRepository @Inject constructor(
    private val api: NotesApiService,
    private val dao: NoteDao
) {
    suspend fun syncNotes(userId: String) {
        val remoteNotes = api.fetchNotes(userId)  // suspend — no callback
        dao.insertAll(remoteNotes.map { it.toEntity() })
    }
}
```

---

## Testing Coroutines

```kotlin
import kotlinx.coroutines.test.runTest
import kotlinx.coroutines.test.StandardTestDispatcher
import kotlinx.coroutines.test.TestCoroutineScheduler

class NoteViewModelTest {

    private val testDispatcher = StandardTestDispatcher()

    @get:Rule
    val mainDispatcherRule = MainDispatcherRule(testDispatcher)

    private lateinit var viewModel: NoteViewModel
    private val fakeRepository = FakeNoteRepository()

    @Before
    fun setup() {
        viewModel = NoteViewModel(fakeRepository)
    }

    @Test
    fun `loading notes updates ui state`() = runTest {
        fakeRepository.setNotes(listOf(Note(title = "Test")))

        val state = viewModel.uiState.value
        assertThat(state.notes).hasSize(1)
        assertThat(state.notes[0].title).isEqualTo("Test")
    }
}

// MainDispatcherRule — sets Dispatchers.Main in tests
class MainDispatcherRule(private val dispatcher: TestCoroutineDispatcher) : TestWatcher() {
    override fun starting(description: Description) {
        Dispatchers.setMain(dispatcher)
    }
    override fun finished(description: Description) {
        Dispatchers.resetMain()
    }
}
```

---

## Common Mistakes

### Mistake 1: Blocking the main thread with `runBlocking`

```kotlin
// WRONG — blocks main thread, causes ANR
override fun onCreate(...) {
    runBlocking { viewModel.loadNotes() }
}

// CORRECT
override fun onCreate(...) {
    lifecycleScope.launch { viewModel.loadNotes() }
}
```

### Mistake 2: Not using the correct dispatcher for work type

```kotlin
// WRONG — network call on Main thread = NetworkOnMainThreadException
viewModelScope.launch {
    val data = api.fetchNotes()
}

// CORRECT — Retrofit suspend functions switch to IO automatically
// but for custom IO, use withContext
viewModelScope.launch {
    val data = withContext(Dispatchers.IO) { api.fetchNotes() }
}
```

### Mistake 3: Catching `CancellationException`

```kotlin
// WRONG — breaks coroutine cancellation
catch (e: Exception) {
    // handles CancellationException silently
}

// CORRECT
catch (e: CancellationException) {
    throw e  // Let cancellation propagate
} catch (e: Exception) {
    handleError(e)
}
```

---

## Interview Questions

**Q1: What is the difference between `launch` and `async` in coroutines?**

> `launch` starts a coroutine for fire-and-forget work — it returns a `Job` with no result. `async` starts a coroutine that produces a result — it returns a `Deferred<T>`, and you call `.await()` to get the value. Use `async` for parallel computation.

**Q2: Why should `CancellationException` always be re-thrown?**

> Coroutine cancellation propagates via `CancellationException`. Catching and swallowing it means the coroutine thinks it's still running — breaks structured concurrency and prevents proper cleanup in parent scopes.

**Q3: What is the difference between `Flow` and `LiveData`?**

> `Flow` is Kotlin-native, coroutine-based, and can emit multiple values. It has powerful operators (`map`, `filter`, `flatMapLatest`) and works outside Android. `LiveData` is lifecycle-aware out of the box but Android-specific and less composable. Prefer `Flow` in the data/domain layers; convert to `StateFlow` for the UI layer.

---

## Summary

- Coroutines let you write async code sequentially without callbacks
- Use `viewModelScope` in ViewModels, `lifecycleScope` in UI components
- `launch` for fire-and-forget; `async`/`await` for parallel result-producing work
- Switch dispatchers with `withContext(Dispatchers.IO)` for I/O work
- Room and Retrofit both have native coroutine/suspend support
- Always re-throw `CancellationException`

**Next:** [Chapter 4 — DataStore](./04-datastore.md)
