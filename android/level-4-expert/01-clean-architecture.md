# Chapter 1: Clean Architecture — MVVM and MVI

## What Is Clean Architecture?

Clean Architecture separates your app into concentric layers, each with a single responsibility and a dependency rule: **inner layers never depend on outer layers**.

```
    ┌──────────────────────────────────┐
    │          UI Layer                │  ← Activity, Fragment, ViewModel
    │  ┌───────────────────────────┐   │
    │  │      Domain Layer         │   │  ← UseCases, Repository Interfaces
    │  │  ┌────────────────────┐   │   │
    │  │  │    Data Layer      │   │   │  ← Repository Impl, API, Database
    │  │  └────────────────────┘   │   │
    │  └───────────────────────────┘   │
    └──────────────────────────────────┘

    Dependency rule: →  (UI depends on Domain; Data depends on Domain)
    Domain depends on NOTHING else
```

---

## The Three Layers

### 1. Domain Layer (Pure Kotlin)

The heart of the app. No Android imports. No framework dependencies.

```
domain/
├── model/
│   └── Note.kt              ← Business entity
├── repository/
│   └── NoteRepository.kt    ← Interface (contract)
└── usecase/
    ├── GetNotesUseCase.kt
    ├── AddNoteUseCase.kt
    ├── DeleteNoteUseCase.kt
    └── SearchNotesUseCase.kt
```

### 2. Data Layer

Implements the domain interfaces. Knows about Room, Retrofit, etc.

```
data/
├── local/
│   ├── NoteDao.kt
│   ├── NoteEntity.kt
│   └── NoteDatabase.kt
├── remote/
│   ├── NoteApiService.kt
│   └── dto/NoteDto.kt
├── mapper/
│   └── NoteMapper.kt
└── repository/
    └── NoteRepositoryImpl.kt  ← Implements NoteRepository
```

### 3. UI Layer

Knows about ViewModel and AndroidX. Does NOT know about Room or Retrofit.

```
ui/
├── notes/
│   ├── NoteListFragment.kt
│   ├── NoteListViewModel.kt
│   └── NoteAdapter.kt
└── detail/
    ├── NoteDetailFragment.kt
    └── NoteDetailViewModel.kt
```

---

## Use Cases (Interactors)

Use cases encapsulate a single piece of business logic:

```kotlin
// domain/usecase/GetNotesUseCase.kt
class GetNotesUseCase @Inject constructor(
    private val repository: NoteRepository
) {
    operator fun invoke(): Flow<List<Note>> = repository.getAllNotes()
}

// domain/usecase/AddNoteUseCase.kt
class AddNoteUseCase @Inject constructor(
    private val repository: NoteRepository
) {
    suspend operator fun invoke(title: String, content: String): Result<Unit> {
        if (title.isBlank()) return Result.failure(IllegalArgumentException("Title cannot be empty"))
        return runCatching {
            repository.saveNote(Note(title = title.trim(), content = content.trim()))
        }
    }
}

// domain/usecase/DeleteNoteUseCase.kt
class DeleteNoteUseCase @Inject constructor(
    private val repository: NoteRepository
) {
    suspend operator fun invoke(note: Note) = repository.deleteNote(note)
}
```

---

## MVVM — Model-View-ViewModel

The most common Android pattern. ViewModel holds UI state; View observes it.

```
View (Fragment) ──observes──► ViewModel ──calls──► UseCase ──calls──► Repository
                              (holds state)
     │                                                                      │
     └─────────────────────events──────────────────────────────────────────┘
```

```kotlin
// MVVM ViewModel
@HiltViewModel
class NoteListViewModel @Inject constructor(
    private val getNotesUseCase: GetNotesUseCase,
    private val addNoteUseCase: AddNoteUseCase,
    private val deleteNoteUseCase: DeleteNoteUseCase
) : ViewModel() {

    data class UiState(
        val notes: List<Note> = emptyList(),
        val isLoading: Boolean = false,
        val error: String? = null
    )

    private val _uiState = MutableStateFlow(UiState())
    val uiState: StateFlow<UiState> = _uiState.asStateFlow()

    init {
        observeNotes()
    }

    private fun observeNotes() {
        viewModelScope.launch {
            getNotesUseCase()
                .catch { e -> _uiState.update { it.copy(error = e.message) } }
                .collect { notes -> _uiState.update { it.copy(notes = notes) } }
        }
    }

    fun addNote(title: String, content: String) {
        viewModelScope.launch {
            addNoteUseCase(title, content).fold(
                onSuccess = { /* Note added — Flow will emit update */ },
                onFailure = { e -> _uiState.update { it.copy(error = e.message) } }
            )
        }
    }

    fun deleteNote(note: Note) {
        viewModelScope.launch {
            deleteNoteUseCase(note)
        }
    }

    fun clearError() {
        _uiState.update { it.copy(error = null) }
    }
}
```

---

## MVI — Model-View-Intent

MVI is stricter than MVVM: the View sends **intents** (user actions) to the ViewModel, which produces a new **state**. State is immutable and the only source of truth.

```
View ──intent──► ViewModel ──reduces──► new State ──► View
                                 (pure function)
```

```kotlin
// MVI with sealed classes

// Intents — user actions
sealed class NoteIntent {
    data class AddNote(val title: String, val content: String) : NoteIntent()
    data class DeleteNote(val note: Note) : NoteIntent()
    data class SearchNotes(val query: String) : NoteIntent()
    object RefreshNotes : NoteIntent()
    object ClearError : NoteIntent()
}

// State — full snapshot of what the screen shows
data class NoteState(
    val notes: List<Note> = emptyList(),
    val filteredNotes: List<Note> = emptyList(),
    val searchQuery: String = "",
    val isLoading: Boolean = false,
    val error: String? = null,
    val isAdding: Boolean = false
) {
    val displayedNotes: List<Note>
        get() = if (searchQuery.isBlank()) notes else filteredNotes
}

// ViewModel — processes intents, emits states
@HiltViewModel
class NoteListMviViewModel @Inject constructor(
    private val getNotesUseCase: GetNotesUseCase,
    private val addNoteUseCase: AddNoteUseCase,
    private val deleteNoteUseCase: DeleteNoteUseCase
) : ViewModel() {

    private val _state = MutableStateFlow(NoteState())
    val state: StateFlow<NoteState> = _state.asStateFlow()

    init {
        handleIntent(NoteIntent.RefreshNotes)
    }

    fun handleIntent(intent: NoteIntent) {
        viewModelScope.launch {
            when (intent) {
                is NoteIntent.RefreshNotes -> loadNotes()
                is NoteIntent.AddNote -> addNote(intent.title, intent.content)
                is NoteIntent.DeleteNote -> deleteNote(intent.note)
                is NoteIntent.SearchNotes -> search(intent.query)
                NoteIntent.ClearError -> _state.update { it.copy(error = null) }
            }
        }
    }

    private suspend fun loadNotes() {
        _state.update { it.copy(isLoading = true) }
        getNotesUseCase()
            .catch { e -> _state.update { it.copy(error = e.message, isLoading = false) } }
            .collect { notes ->
                _state.update { it.copy(notes = notes, isLoading = false) }
            }
    }

    private suspend fun addNote(title: String, content: String) {
        _state.update { it.copy(isAdding = true) }
        addNoteUseCase(title, content).fold(
            onSuccess = { _state.update { it.copy(isAdding = false) } },
            onFailure = { e -> _state.update { it.copy(error = e.message, isAdding = false) } }
        )
    }

    private suspend fun deleteNote(note: Note) {
        deleteNoteUseCase(note)
    }

    private fun search(query: String) {
        val filtered = _state.value.notes.filter {
            it.title.contains(query, ignoreCase = true) ||
            it.content.contains(query, ignoreCase = true)
        }
        _state.update { it.copy(searchQuery = query, filteredNotes = filtered) }
    }
}
```

**Usage in Fragment (MVI):**

```kotlin
// Send intents to ViewModel
binding.btnAdd.setOnClickListener {
    viewModel.handleIntent(
        NoteIntent.AddNote(binding.etTitle.text.toString(), binding.etContent.text.toString())
    )
}

binding.etSearch.doAfterTextChanged { text ->
    viewModel.handleIntent(NoteIntent.SearchNotes(text?.toString() ?: ""))
}

// Observe state
viewLifecycleOwner.lifecycleScope.launch {
    viewLifecycleOwner.repeatOnLifecycle(Lifecycle.State.STARTED) {
        viewModel.state.collect { state ->
            adapter.submitList(state.displayedNotes)
            binding.progressBar.isVisible = state.isLoading
            state.error?.let { error ->
                Snackbar.make(binding.root, error, Snackbar.LENGTH_SHORT).show()
                viewModel.handleIntent(NoteIntent.ClearError)
            }
        }
    }
}
```

---

## MVVM vs MVI — When to Use Which

| | MVVM | MVI |
|--|------|-----|
| State management | Multiple LiveData/StateFlow | Single state object |
| UI logic complexity | Simple–Medium | Medium–Complex |
| Testability | Good | Excellent |
| Debugging | Harder (multiple state sources) | Easy (single state snapshot) |
| Boilerplate | Less | More |
| Best for | Most apps | Complex UIs, lots of state transitions |

---

## Data Mapper Pattern

Keep domain models clean by mapping at layer boundaries:

```kotlin
// data/mapper/NoteMapper.kt
fun NoteEntity.toDomain(): Note = Note(
    id = id,
    title = title,
    content = content,
    updatedAt = updatedAt
)

fun Note.toEntity(): NoteEntity = NoteEntity(
    id = id,
    title = title,
    content = content,
    updatedAt = updatedAt
)

fun NoteDto.toDomain(): Note = Note(
    id = id,
    title = title ?: "",
    content = content ?: "",
    updatedAt = System.currentTimeMillis()
)
```

---

## Repository Implementation

```kotlin
class NoteRepositoryImpl @Inject constructor(
    private val dao: NoteDao,
    private val api: NoteApiService,
    @IoDispatcher private val ioDispatcher: CoroutineDispatcher
) : NoteRepository {

    override fun getAllNotes(): Flow<List<Note>> =
        dao.getAllNotes()
            .map { entities -> entities.map { it.toDomain() } }
            .flowOn(ioDispatcher)

    override suspend fun saveNote(note: Note): Long =
        withContext(ioDispatcher) { dao.insertNote(note.toEntity()) }

    override suspend fun deleteNote(note: Note) =
        withContext(ioDispatcher) { dao.deleteNote(note.toEntity()) }

    override suspend fun syncWithRemote() = withContext(ioDispatcher) {
        val remoteNotes = api.fetchNotes()
        dao.insertAll(remoteNotes.map { it.toDomain().toEntity() })
    }
}
```

---

## Interview Questions

**Q1: What is Clean Architecture and what problem does it solve?**

> Clean Architecture separates an app into layers (UI, Domain, Data) with strict dependency rules — outer layers depend on inner ones, never the reverse. It solves the problem of tightly-coupled code that's hard to test, maintain, or replace. Changing the database shouldn't require touching the UI.

**Q2: What is the difference between MVVM and MVI?**

> MVVM has a ViewModel that holds one or more state streams which the View observes. MVI is more rigid: the View sends typed intents to the ViewModel, which produces a single, immutable state. MVI makes state transitions explicit and easier to debug and test.

**Q3: What is a Use Case / Interactor?**

> A Use Case is a class in the domain layer that encapsulates a single piece of business logic. It takes inputs, calls repository methods, applies business rules, and returns a result. It keeps business logic out of ViewModels (which are Android-specific) and makes it testable with pure Kotlin.

---

## Summary

- Clean Architecture: UI → Domain ← Data (dependencies point inward)
- Domain layer: pure Kotlin, zero Android imports, contains models + interfaces + use cases
- Use Cases encapsulate single business operations
- MVVM: View observes ViewModel state; good for most apps
- MVI: View sends intents; ViewModel emits single state; best for complex UI
- Use mappers at layer boundaries — never let entities leak to the UI

**Next:** [Chapter 2 — Modularization](./02-modularization.md)
