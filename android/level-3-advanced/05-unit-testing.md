# Chapter 5: Unit Testing

## What Is Unit Testing?

A unit test verifies a single, isolated piece of code — a function or class — in isolation from its dependencies. Dependencies are replaced with **fakes** or **mocks**.

**Why test?**
- Catch bugs before they reach production
- Document expected behavior
- Enable safe refactoring
- Prevent regressions

```kotlin
// Test dependencies
testImplementation("junit:junit:4.13.2")
testImplementation("org.mockito.kotlin:mockito-kotlin:5.4.0")
testImplementation("com.google.truth:truth:1.4.4")
testImplementation("org.jetbrains.kotlinx:kotlinx-coroutines-test:1.8.1")
testImplementation("androidx.arch.core:core-testing:2.2.0")
```

---

## Test Structure: Arrange → Act → Assert

```kotlin
@Test
fun `adding a task increases task count`() {
    // Arrange — set up preconditions
    val viewModel = TaskViewModel(FakeTaskRepository())

    // Act — do the thing being tested
    viewModel.addTask("Buy milk")

    // Assert — verify the outcome
    assertThat(viewModel.tasks.value).hasSize(1)
    assertThat(viewModel.tasks.value[0].title).isEqualTo("Buy milk")
}
```

---

## Testing a ViewModel

### The ViewModel Under Test

```kotlin
@HiltViewModel
class NoteViewModel @Inject constructor(
    private val repository: NoteRepository
) : ViewModel() {

    private val _uiState = MutableStateFlow(NoteUiState())
    val uiState: StateFlow<NoteUiState> = _uiState.asStateFlow()

    fun addNote(title: String, content: String) {
        if (title.isBlank()) {
            _uiState.update { it.copy(error = "Title cannot be empty") }
            return
        }
        viewModelScope.launch {
            repository.saveNote(Note(title = title, content = content))
            _uiState.update { it.copy(error = null) }
        }
    }

    fun loadNotes() {
        viewModelScope.launch {
            _uiState.update { it.copy(isLoading = true) }
            try {
                val notes = repository.getNotesList()
                _uiState.update { it.copy(notes = notes, isLoading = false) }
            } catch (e: Exception) {
                _uiState.update { it.copy(error = e.message, isLoading = false) }
            }
        }
    }
}
```

### Fake Repository

```kotlin
class FakeNoteRepository : NoteRepository {

    private val notes = mutableListOf<Note>()
    var shouldThrow = false

    override fun getAllNotes(): Flow<List<Note>> = flow {
        emit(notes.toList())
    }

    override suspend fun getNotesList(): List<Note> {
        if (shouldThrow) throw IOException("Network error")
        return notes.toList()
    }

    override suspend fun saveNote(note: Note): Long {
        notes.add(note.copy(id = notes.size.toLong() + 1))
        return notes.last().id
    }

    override suspend fun deleteNote(note: Note) {
        notes.removeAll { it.id == note.id }
    }

    fun setNotes(noteList: List<Note>) {
        notes.clear()
        notes.addAll(noteList)
    }
}
```

### `MainDispatcherRule`

```kotlin
class MainDispatcherRule(
    val testDispatcher: TestDispatcher = UnconfinedTestDispatcher()
) : TestWatcher() {
    override fun starting(description: Description) {
        Dispatchers.setMain(testDispatcher)
    }
    override fun finished(description: Description) {
        Dispatchers.resetMain()
    }
}
```

### ViewModel Tests

```kotlin
class NoteViewModelTest {

    @get:Rule
    val mainDispatcherRule = MainDispatcherRule()

    private lateinit var viewModel: NoteViewModel
    private lateinit var fakeRepository: FakeNoteRepository

    @Before
    fun setup() {
        fakeRepository = FakeNoteRepository()
        viewModel = NoteViewModel(fakeRepository)
    }

    @Test
    fun `addNote with blank title sets error state`() {
        viewModel.addNote("", "Some content")

        val state = viewModel.uiState.value
        assertThat(state.error).isEqualTo("Title cannot be empty")
    }

    @Test
    fun `addNote with valid title clears error`() = runTest {
        fakeRepository.setNotes(emptyList())

        viewModel.addNote("Test Note", "Content")

        assertThat(viewModel.uiState.value.error).isNull()
    }

    @Test
    fun `loadNotes populates notes in state`() = runTest {
        fakeRepository.setNotes(listOf(
            Note(id = 1, title = "First"),
            Note(id = 2, title = "Second")
        ))

        viewModel.loadNotes()

        val state = viewModel.uiState.value
        assertThat(state.notes).hasSize(2)
        assertThat(state.isLoading).isFalse()
    }

    @Test
    fun `loadNotes with exception sets error state`() = runTest {
        fakeRepository.shouldThrow = true

        viewModel.loadNotes()

        val state = viewModel.uiState.value
        assertThat(state.error).isEqualTo("Network error")
        assertThat(state.isLoading).isFalse()
    }

    @Test
    fun `loadNotes shows loading then hides it`() = runTest {
        val states = mutableListOf<NoteUiState>()
        val job = launch { viewModel.uiState.collect { states.add(it) } }

        viewModel.loadNotes()

        job.cancel()

        // First state: isLoading = true
        assertThat(states[1].isLoading).isTrue()
        // Final state: isLoading = false
        assertThat(states.last().isLoading).isFalse()
    }
}
```

---

## Testing a Repository

```kotlin
class NoteRepositoryTest {

    // Real in-memory Room database
    private lateinit var database: NoteDatabase
    private lateinit var dao: NoteDao
    private lateinit var repository: NoteRepository

    @Before
    fun setup() {
        database = Room.inMemoryDatabaseBuilder(
            ApplicationProvider.getApplicationContext(),
            NoteDatabase::class.java
        ).allowMainThreadQueries().build()

        dao = database.noteDao()
        repository = NoteRepository(dao)
    }

    @After
    fun tearDown() {
        database.close()
    }

    @Test
    fun `saveNote persists and is retrievable`() = runTest {
        val note = Note(title = "Test", content = "Content")
        repository.saveNote(note)

        val notes = repository.getAllNotes().first()
        assertThat(notes).hasSize(1)
        assertThat(notes[0].title).isEqualTo("Test")
    }

    @Test
    fun `deleteNote removes note from database`() = runTest {
        val note = Note(title = "To delete")
        val id = repository.saveNote(note)

        repository.deleteNote(note.copy(id = id))

        val notes = repository.getAllNotes().first()
        assertThat(notes).isEmpty()
    }

    @Test
    fun `getAllNotes emits updates when data changes`() = runTest {
        val results = mutableListOf<List<Note>>()
        val job = launch { repository.getAllNotes().collect { results.add(it) } }

        repository.saveNote(Note(title = "First"))
        repository.saveNote(Note(title = "Second"))

        job.cancel()

        assertThat(results.last()).hasSize(2)
    }
}
```

---

## Mocking with Mockito

When creating a real fake is too complex, use Mockito:

```kotlin
import org.mockito.kotlin.*

class NoteViewModelMockTest {

    @get:Rule
    val mainDispatcherRule = MainDispatcherRule()

    private val mockRepository: NoteRepository = mock()
    private lateinit var viewModel: NoteViewModel

    @Before
    fun setup() {
        viewModel = NoteViewModel(mockRepository)
    }

    @Test
    fun `loadNotes calls repository`() = runTest {
        whenever(mockRepository.getNotesList()).thenReturn(emptyList())

        viewModel.loadNotes()

        verify(mockRepository).getNotesList()
    }

    @Test
    fun `deleteNote calls repository with correct note`() = runTest {
        val note = Note(id = 1, title = "Test")

        viewModel.deleteNote(note)

        verify(mockRepository).deleteNote(note)
    }

    @Test
    fun `error from repository propagates to UI state`() = runTest {
        whenever(mockRepository.getNotesList()).thenThrow(IOException("Network error"))

        viewModel.loadNotes()

        assertThat(viewModel.uiState.value.error).contains("Network error")
    }
}
```

---

## Test Architecture Best Practices

### Test Naming Convention

```kotlin
// Format: `functionName with condition returns result`
@Test fun `addNote with empty title returns error`()
@Test fun `loadNotes on success emits notes list`()
@Test fun `deleteNote removes note and shows snackbar`()
```

### Test Pyramid

```
          ┌───────┐
          │  E2E  │  ← Few (slow, flaky, expensive)
         /─────────\
        │ Integration│
       /─────────────\
      │   Unit Tests  │  ← Many (fast, reliable, cheap)
     /─────────────────\
```

- **70% unit tests** — fast, no Android framework, test business logic
- **20% integration tests** — Room in-memory, multiple components
- **10% UI tests** — full Espresso, test the whole screen

### Fake vs Mock

| | Fake | Mock |
|--|------|------|
| What | A real but simplified implementation | A recording object that captures calls |
| Use | When the dependency has complex behavior | When you only need to verify calls |
| Example | `FakeNoteRepository : NoteRepository` | `mock<NoteRepository>()` |

---

## Common Mistakes

### Mistake 1: Not using `MainDispatcherRule`

```kotlin
// FAILS — viewModelScope uses Dispatchers.Main, which isn't set in tests
@Test fun `test vm`() {
    viewModel.loadNotes()
}

// CORRECT — set Main dispatcher for tests
@get:Rule val mainDispatcherRule = MainDispatcherRule()
```

### Mistake 2: Testing implementation details

```kotlin
// WRONG — tests HOW, not WHAT
verify(mockRepository, times(1)).getNotesList()

// BETTER — test observable behavior
assertThat(viewModel.uiState.value.notes).hasSize(3)
```

### Mistake 3: Slow tests from using real Android components

Unit tests should run on the JVM, not a device. Use `ApplicationProvider.getApplicationContext()` only in integration tests, not unit tests.

---

## Interview Questions

**Q1: What is the difference between a fake and a mock?**

> A fake is a simplified but working implementation of an interface (e.g., an in-memory list instead of a database). A mock is a generated spy that records method calls so you can verify interactions. Use fakes for testing behavior; use mocks to verify specific interactions.

**Q2: Why do you need `MainDispatcherRule` in ViewModel tests?**

> `viewModelScope` uses `Dispatchers.Main` internally. In tests, `Dispatchers.Main` doesn't exist — the test rules override it with a `TestDispatcher` that runs coroutines synchronously in tests.

**Q3: What is the test pyramid and why does it matter?**

> The test pyramid says to have many unit tests (fast, cheap), fewer integration tests, and minimal E2E tests (slow, expensive). An inverted pyramid (many E2E, few unit) results in slow CI, flaky tests, and poor feedback loops.

---

## Summary

- Unit tests verify isolated logic; use `FakeRepository` or Mockito for dependencies
- Always use `MainDispatcherRule` to replace `Dispatchers.Main` in ViewModel tests
- Use `runTest` for coroutine-based tests — it controls virtual time
- Use in-memory Room database for repository integration tests
- Follow Arrange → Act → Assert structure

**Next:** [Chapter 6 — UI Testing with Espresso](./06-ui-testing-espresso.md)
