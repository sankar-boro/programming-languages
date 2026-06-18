# Chapter 3: StateFlow — The Modern Approach

## Why StateFlow?

`LiveData` was designed specifically for Android and requires the `androidx.lifecycle` dependency. `StateFlow` is part of Kotlin Coroutines — framework-agnostic, testable, and the preferred choice in modern Android code.

Google recommends `StateFlow` for new UI state management, especially in apps already using coroutines.

---

## StateFlow vs SharedFlow vs Channel

| | StateFlow | SharedFlow | Channel |
|--|-----------|------------|---------|
| Holds current value | Yes | No | No |
| Replays to new subscribers | Last value | Configurable | 0 (one-shot) |
| Use for | UI state | Events (analytics, etc.) | One-time navigation events |
| Hot/Cold | Hot | Hot | Hot |

---

## StateFlow Basics

```kotlin
import kotlinx.coroutines.flow.MutableStateFlow
import kotlinx.coroutines.flow.StateFlow
import kotlinx.coroutines.flow.asStateFlow
import kotlinx.coroutines.flow.update

class TaskViewModel : ViewModel() {

    // Internal mutable state
    private val _uiState = MutableStateFlow(TaskUiState())

    // Exposed read-only state
    val uiState: StateFlow<TaskUiState> = _uiState.asStateFlow()

    fun addTask(title: String) {
        _uiState.update { currentState ->
            currentState.copy(
                tasks = currentState.tasks + Task(title = title)
            )
        }
    }

    fun setLoading(loading: Boolean) {
        _uiState.update { it.copy(isLoading = loading) }
    }
}

// UI state data class — immutable, represents the full state
data class TaskUiState(
    val tasks: List<Task> = emptyList(),
    val isLoading: Boolean = false,
    val error: String? = null
)
```

The `update { }` lambda applies changes atomically — safe for concurrent updates.

---

## Collecting StateFlow in Activity/Fragment

StateFlow requires a coroutine to collect. The lifecycle-safe way is `repeatOnLifecycle`:

```kotlin
import androidx.lifecycle.lifecycleScope
import androidx.lifecycle.repeatOnLifecycle
import kotlinx.coroutines.launch
import kotlinx.coroutines.flow.collectLatest

class MainActivity : AppCompatActivity() {

    private val viewModel: TaskViewModel by viewModels()

    override fun onCreate(savedInstanceState: Bundle?) {
        super.onCreate(savedInstanceState)
        binding = ActivityMainBinding.inflate(layoutInflater)
        setContentView(binding.root)

        observeState()
    }

    private fun observeState() {
        lifecycleScope.launch {
            // repeatOnLifecycle cancels and restarts collection on lifecycle transitions
            repeatOnLifecycle(Lifecycle.State.STARTED) {
                viewModel.uiState.collect { state ->
                    renderState(state)
                }
            }
        }
    }

    private fun renderState(state: TaskUiState) {
        adapter.submitList(state.tasks)
        binding.progressBar.isVisible = state.isLoading
        binding.tvEmpty.isVisible = state.tasks.isEmpty() && !state.isLoading

        state.error?.let { error ->
            Snackbar.make(binding.root, error, Snackbar.LENGTH_LONG).show()
            viewModel.clearError()
        }
    }
}
```

In Fragment:

```kotlin
override fun onViewCreated(view: View, savedInstanceState: Bundle?) {
    super.onViewCreated(view, savedInstanceState)

    viewLifecycleOwner.lifecycleScope.launch {
        viewLifecycleOwner.repeatOnLifecycle(Lifecycle.State.STARTED) {
            viewModel.uiState.collect { state ->
                renderState(state)
            }
        }
    }
}
```

---

## Why `repeatOnLifecycle` and Not `launchWhenStarted`

`launchWhenStarted` only **suspends** collection when the app is in the background — the coroutine is still alive and holding resources. `repeatOnLifecycle` **cancels** the collection entirely when below the given state and **restarts** it when above — properly releasing resources.

```kotlin
// OUTDATED — still leaks in some scenarios
lifecycleScope.launchWhenStarted {
    viewModel.uiState.collect { ... }
}

// CORRECT — cancels and restarts with lifecycle
lifecycleScope.launch {
    repeatOnLifecycle(Lifecycle.State.STARTED) {
        viewModel.uiState.collect { ... }
    }
}
```

---

## `collectLatest` — Cancel Previous on New Emission

When each emission triggers a slow operation, use `collectLatest` to cancel the in-progress handler:

```kotlin
lifecycleScope.launch {
    repeatOnLifecycle(Lifecycle.State.STARTED) {
        viewModel.searchQuery.collectLatest { query ->
            // If a new query arrives before this completes, this block is cancelled
            val results = repository.search(query)
            adapter.submitList(results)
        }
    }
}
```

---

## One-Time Events with SharedFlow

For navigation commands, Snackbars, and other one-time events:

```kotlin
class TaskViewModel : ViewModel() {

    // SharedFlow for events — no replay = each event handled once
    private val _events = MutableSharedFlow<TaskEvent>()
    val events = _events.asSharedFlow()

    fun onTaskLongPressed(taskId: Long) {
        viewModelScope.launch {
            _events.emit(TaskEvent.ShowDeleteConfirmation(taskId))
        }
    }

    fun onAddSuccess(taskTitle: String) {
        viewModelScope.launch {
            _events.emit(TaskEvent.ShowSnackbar("'$taskTitle' added"))
        }
    }
}

sealed class TaskEvent {
    data class ShowDeleteConfirmation(val taskId: Long) : TaskEvent()
    data class ShowSnackbar(val message: String) : TaskEvent()
    data class NavigateToDetail(val taskId: Long) : TaskEvent()
}

// In Activity/Fragment
lifecycleScope.launch {
    repeatOnLifecycle(Lifecycle.State.STARTED) {
        viewModel.events.collect { event ->
            when (event) {
                is TaskEvent.ShowDeleteConfirmation -> showDeleteDialog(event.taskId)
                is TaskEvent.ShowSnackbar -> showSnackbar(event.message)
                is TaskEvent.NavigateToDetail -> navigateToDetail(event.taskId)
            }
        }
    }
}
```

---

## Full ViewModel Example with StateFlow + SharedFlow

```kotlin
class TaskViewModel(
    private val repository: TaskRepository
) : ViewModel() {

    private val _uiState = MutableStateFlow(TaskUiState())
    val uiState: StateFlow<TaskUiState> = _uiState.asStateFlow()

    private val _events = MutableSharedFlow<TaskEvent>()
    val events = _events.asSharedFlow()

    init {
        loadTasks()
    }

    private fun loadTasks() {
        viewModelScope.launch {
            _uiState.update { it.copy(isLoading = true) }
            try {
                val tasks = repository.getTasks()
                _uiState.update { it.copy(tasks = tasks, isLoading = false) }
            } catch (e: Exception) {
                _uiState.update { it.copy(isLoading = false, error = e.message) }
            }
        }
    }

    fun addTask(title: String) {
        if (title.isBlank()) {
            viewModelScope.launch {
                _events.emit(TaskEvent.ShowSnackbar("Task title cannot be empty"))
            }
            return
        }
        viewModelScope.launch {
            val task = Task(title = title)
            repository.insertTask(task)
            _uiState.update { state ->
                state.copy(tasks = state.tasks + task)
            }
            _events.emit(TaskEvent.ShowSnackbar("Task added"))
        }
    }

    fun clearError() {
        _uiState.update { it.copy(error = null) }
    }
}
```

---

## Flow Operators in ViewModel

```kotlin
class TaskViewModel : ViewModel() {

    private val _searchQuery = MutableStateFlow("")

    // Derived state using Flow operators — updates automatically
    val filteredTasks: StateFlow<List<Task>> = _searchQuery
        .debounce(300)  // Wait 300ms after last keystroke
        .distinctUntilChanged()
        .flatMapLatest { query ->
            if (query.isBlank()) repository.getAllTasks()
            else repository.searchTasks(query)
        }
        .stateIn(
            scope = viewModelScope,
            started = SharingStarted.WhileSubscribed(5_000),
            initialValue = emptyList()
        )

    fun setSearchQuery(query: String) {
        _searchQuery.value = query
    }
}
```

`stateIn` converts a `Flow` to a `StateFlow` so the UI can use `collect` on it.

---

## Common Mistakes

### Mistake 1: Using `lifecycleScope.launch { flow.collect {} }` without `repeatOnLifecycle`

```kotlin
// WRONG — collection continues in background, wasting resources
lifecycleScope.launch {
    viewModel.uiState.collect { renderState(it) }
}

// CORRECT
lifecycleScope.launch {
    repeatOnLifecycle(Lifecycle.State.STARTED) {
        viewModel.uiState.collect { renderState(it) }
    }
}
```

### Mistake 2: Collecting in `onStart()` instead of `repeatOnLifecycle`

Manually managing coroutine start/stop in lifecycle callbacks is error-prone. Use `repeatOnLifecycle`.

### Mistake 3: Emitting to SharedFlow outside a coroutine

```kotlin
// WRONG — emit is a suspend function
_events.emit(TaskEvent.ShowSnackbar("hi"))

// CORRECT
viewModelScope.launch {
    _events.emit(TaskEvent.ShowSnackbar("hi"))
}

// OR use tryEmit for non-suspend contexts (only if buffer is configured)
_events.tryEmit(TaskEvent.ShowSnackbar("hi"))
```

---

## Interview Questions

**Q1: What is the difference between StateFlow and LiveData?**

> `StateFlow` is Kotlin-native, coroutine-based, and requires an initial value. `LiveData` is Android-specific with built-in lifecycle awareness. `StateFlow` needs `repeatOnLifecycle` to be lifecycle-safe, while `LiveData.observe()` handles this automatically. Prefer `StateFlow` in modern Kotlin projects.

**Q2: What is the difference between StateFlow and SharedFlow?**

> `StateFlow` always holds and replays its current value to new collectors — use it for UI state. `SharedFlow` has configurable replay and no mandatory current value — use it for events that should only be processed once.

**Q3: Why use `repeatOnLifecycle` instead of `launchWhenStarted`?**

> `launchWhenStarted` suspends the collection but doesn't cancel the coroutine — the coroutine still exists while the app is backgrounded. `repeatOnLifecycle` cancels the coroutine entirely below the specified state and restarts it above, properly releasing resources.

---

## Summary

- `StateFlow` is the modern alternative to `LiveData` for UI state
- Expose `StateFlow` (read-only); keep `MutableStateFlow` private in the ViewModel
- Use `update { }` for atomic state changes
- Always collect with `repeatOnLifecycle(Lifecycle.State.STARTED)` for lifecycle safety
- Use `SharedFlow` for one-time events (navigation, Snackbars)

**Next:** [Chapter 4 — Lifecycle](./04-lifecycle.md)
