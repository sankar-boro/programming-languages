# Mini Project: Notes App

## Overview

Build a multi-screen Notes app that persists data using Room, navigates with Navigation Component, and manages state with ViewModel + StateFlow.

**Features:**
- Notes list screen with search
- Add/edit note screen
- Delete with undo (Snackbar)
- Data persists across app restarts (Room)
- Navigation between screens (Navigation Component)
- MVVM with ViewModel + StateFlow

---

## Architecture

```
UI Layer
├── NoteListFragment (displays list + search)
├── NoteEditorFragment (add/edit a note)
└── NoteViewModel (shared ViewModel)

Domain Layer
└── Note.kt (domain model)

Data Layer
├── NoteEntity.kt (Room entity)
├── NoteDao.kt
├── NoteDatabase.kt
├── NoteMapper.kt (Entity ↔ Domain)
└── NoteRepository.kt
```

---

## Dependencies (`app/build.gradle.kts`)

```kotlin
plugins {
    alias(libs.plugins.android.application)
    alias(libs.plugins.kotlin.android)
    alias(libs.plugins.navigation.safeargs.kotlin)
    id("com.google.devtools.ksp")
}

dependencies {
    implementation(libs.androidx.core.ktx)
    implementation(libs.androidx.appcompat)
    implementation(libs.material)
    implementation(libs.androidx.constraintlayout)

    // Navigation
    implementation(libs.androidx.navigation.fragment.ktx)
    implementation(libs.androidx.navigation.ui.ktx)

    // Lifecycle
    implementation(libs.androidx.lifecycle.viewmodel.ktx)
    implementation(libs.androidx.lifecycle.livedata.ktx)
    implementation(libs.androidx.lifecycle.runtime.ktx)

    // Room
    implementation(libs.androidx.room.runtime)
    implementation(libs.androidx.room.ktx)
    ksp(libs.androidx.room.compiler)

    // Coroutines
    implementation(libs.kotlinx.coroutines.android)
}
```

---

## Data Layer

### `NoteEntity.kt`

```kotlin
@Entity(tableName = "notes")
data class NoteEntity(
    @PrimaryKey(autoGenerate = true)
    val id: Long = 0,
    val title: String,
    val content: String,
    val createdAt: Long = System.currentTimeMillis(),
    val updatedAt: Long = System.currentTimeMillis()
)
```

### `NoteDao.kt`

```kotlin
@Dao
interface NoteDao {

    @Query("SELECT * FROM notes ORDER BY updated_at DESC")
    fun getAllNotes(): Flow<List<NoteEntity>>

    @Query("""
        SELECT * FROM notes 
        WHERE title LIKE '%' || :query || '%' 
           OR content LIKE '%' || :query || '%'
        ORDER BY updated_at DESC
    """)
    fun searchNotes(query: String): Flow<List<NoteEntity>>

    @Query("SELECT * FROM notes WHERE id = :id")
    suspend fun getNoteById(id: Long): NoteEntity?

    @Insert(onConflict = OnConflictStrategy.REPLACE)
    suspend fun insertNote(note: NoteEntity): Long

    @Update
    suspend fun updateNote(note: NoteEntity)

    @Delete
    suspend fun deleteNote(note: NoteEntity)
}
```

### `NoteDatabase.kt`

```kotlin
@Database(entities = [NoteEntity::class], version = 1, exportSchema = true)
abstract class NoteDatabase : RoomDatabase() {

    abstract fun noteDao(): NoteDao

    companion object {
        @Volatile private var INSTANCE: NoteDatabase? = null

        fun getInstance(context: Context): NoteDatabase =
            INSTANCE ?: synchronized(this) {
                Room.databaseBuilder(context.applicationContext,
                    NoteDatabase::class.java, "notes.db")
                    .build()
                    .also { INSTANCE = it }
            }
    }
}
```

### `Note.kt` (Domain Model)

```kotlin
data class Note(
    val id: Long = 0,
    val title: String,
    val content: String,
    val updatedAt: Long = System.currentTimeMillis()
)
```

### `NoteMapper.kt`

```kotlin
fun NoteEntity.toDomain() = Note(id, title, content, updatedAt)
fun Note.toEntity() = NoteEntity(id, title, content, updatedAt = updatedAt)
```

### `NoteRepository.kt`

```kotlin
class NoteRepository(private val dao: NoteDao) {

    fun getAllNotes(): Flow<List<Note>> =
        dao.getAllNotes().map { list -> list.map { it.toDomain() } }

    fun searchNotes(query: String): Flow<List<Note>> =
        dao.searchNotes(query).map { list -> list.map { it.toDomain() } }

    suspend fun getNoteById(id: Long): Note? = dao.getNoteById(id)?.toDomain()

    suspend fun saveNote(note: Note): Long = dao.insertNote(note.toEntity())

    suspend fun updateNote(note: Note) = dao.updateNote(note.toEntity())

    suspend fun deleteNote(note: Note) = dao.deleteNote(note.toEntity())
}
```

---

## ViewModel

```kotlin
data class NoteListUiState(
    val notes: List<Note> = emptyList(),
    val isLoading: Boolean = false,
    val searchQuery: String = ""
)

sealed class NoteEvent {
    data class ShowMessage(val message: String) : NoteEvent()
    data class NavigateToEditor(val noteId: Long?) : NoteEvent()
}

class NoteViewModel(
    private val repository: NoteRepository
) : ViewModel() {

    private val _uiState = MutableStateFlow(NoteListUiState(isLoading = true))
    val uiState: StateFlow<NoteListUiState> = _uiState.asStateFlow()

    private val _events = MutableSharedFlow<NoteEvent>()
    val events = _events.asSharedFlow()

    private val searchQuery = MutableStateFlow("")

    val notes: StateFlow<List<Note>> = searchQuery
        .debounce(300)
        .distinctUntilChanged()
        .flatMapLatest { query ->
            if (query.isBlank()) repository.getAllNotes()
            else repository.searchNotes(query)
        }
        .stateIn(viewModelScope, SharingStarted.WhileSubscribed(5_000), emptyList())

    fun setSearchQuery(query: String) {
        searchQuery.value = query
        _uiState.update { it.copy(searchQuery = query) }
    }

    fun onAddNote() {
        viewModelScope.launch {
            _events.emit(NoteEvent.NavigateToEditor(null))
        }
    }

    fun onNoteClicked(note: Note) {
        viewModelScope.launch {
            _events.emit(NoteEvent.NavigateToEditor(note.id))
        }
    }

    fun deleteNote(note: Note) {
        viewModelScope.launch {
            repository.deleteNote(note)
            _events.emit(NoteEvent.ShowMessage("Note deleted"))
        }
    }

    fun saveNote(title: String, content: String, existingId: Long?) {
        if (title.isBlank()) {
            viewModelScope.launch {
                _events.emit(NoteEvent.ShowMessage("Title cannot be empty"))
            }
            return
        }
        viewModelScope.launch {
            if (existingId != null) {
                repository.updateNote(Note(id = existingId, title = title, content = content))
            } else {
                repository.saveNote(Note(title = title, content = content))
            }
            _events.emit(NoteEvent.ShowMessage("Note saved"))
        }
    }
}
```

---

## Navigation Graph (`nav_notes.xml`)

```xml
<navigation
    xmlns:android="http://schemas.android.com/apk/res/android"
    xmlns:app="http://schemas.android.com/apk/res-auto"
    android:id="@+id/nav_notes"
    app:startDestination="@id/noteListFragment">

    <fragment
        android:id="@+id/noteListFragment"
        android:name="com.yourname.notes.NoteListFragment"
        android:label="Notes">
        <action
            android:id="@+id/action_list_to_editor"
            app:destination="@id/noteEditorFragment" />
    </fragment>

    <fragment
        android:id="@+id/noteEditorFragment"
        android:name="com.yourname.notes.NoteEditorFragment"
        android:label="Edit Note">
        <argument
            android:name="noteId"
            app:argType="long"
            android:defaultValue="-1" />
    </fragment>

</navigation>
```

---

## `NoteListFragment`

```kotlin
class NoteListFragment : Fragment(R.layout.fragment_note_list) {

    private var _binding: FragmentNoteListBinding? = null
    private val binding get() = _binding!!

    private val viewModel: NoteViewModel by activityViewModels {
        NoteViewModelFactory(
            NoteRepository(NoteDatabase.getInstance(requireContext()).noteDao())
        )
    }

    private lateinit var adapter: NoteAdapter

    override fun onViewCreated(view: View, savedInstanceState: Bundle?) {
        super.onViewCreated(view, savedInstanceState)
        _binding = FragmentNoteListBinding.bind(view)

        setupRecyclerView()
        setupSearch()
        observeState()
        observeEvents()

        binding.fab.setOnClickListener { viewModel.onAddNote() }
    }

    private fun setupRecyclerView() {
        adapter = NoteAdapter(
            onClick = { note -> viewModel.onNoteClicked(note) },
            onDelete = { note -> viewModel.deleteNote(note) }
        )
        binding.rvNotes.layoutManager = LinearLayoutManager(requireContext())
        binding.rvNotes.adapter = adapter
    }

    private fun setupSearch() {
        binding.etSearch.doAfterTextChanged { text ->
            viewModel.setSearchQuery(text?.toString() ?: "")
        }
    }

    private fun observeState() {
        viewLifecycleOwner.lifecycleScope.launch {
            viewLifecycleOwner.repeatOnLifecycle(Lifecycle.State.STARTED) {
                viewModel.notes.collect { notes ->
                    adapter.submitList(notes)
                    binding.tvEmpty.isVisible = notes.isEmpty()
                }
            }
        }
    }

    private fun observeEvents() {
        viewLifecycleOwner.lifecycleScope.launch {
            viewLifecycleOwner.repeatOnLifecycle(Lifecycle.State.STARTED) {
                viewModel.events.collect { event ->
                    when (event) {
                        is NoteEvent.ShowMessage ->
                            Snackbar.make(binding.root, event.message, Snackbar.LENGTH_SHORT).show()
                        is NoteEvent.NavigateToEditor -> {
                            val action = NoteListFragmentDirections
                                .actionListToEditor(noteId = event.noteId ?: -1L)
                            findNavController().navigate(action)
                        }
                    }
                }
            }
        }
    }

    override fun onDestroyView() {
        super.onDestroyView()
        _binding = null
    }
}
```

---

## Level 2 Checkpoint

Before moving to Level 3, confirm you can:

- [ ] Create a ViewModel that survives rotation
- [ ] Use LiveData or StateFlow to push data to the UI reactively
- [ ] Set up Room with Entity, Dao, Database, and Repository
- [ ] Navigate between fragments using Navigation Component with Safe Args
- [ ] Observe a Room Flow in a ViewModel and expose it as StateFlow
- [ ] Handle one-time events (Snackbar, navigation) with SharedFlow
- [ ] Clear Fragment binding in `onDestroyView()`

**Next Level:** [Level 3 — Advanced](../level-3-advanced/01-hilt-dependency-injection.md)
