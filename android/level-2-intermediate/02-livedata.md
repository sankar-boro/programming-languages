# Chapter 2: LiveData

## The Problem LiveData Solves

In Chapter 1, we manually called `displayTasks()` every time data changed. This is fragile — you have to remember every place that modifies data and manually refresh the UI after each one. It also has a lifecycle problem: what if the UI updates while the Activity is in the background?

**Analogy:** LiveData is like a news subscription. Instead of you checking the newspaper every hour (polling), the newspaper calls you whenever there is an update — but only when you're home (the observer is active).

---

## What LiveData Is

`LiveData` is a lifecycle-aware observable data holder. It:

1. Notifies observers when data changes
2. Only notifies when the observer's lifecycle is `STARTED` or `RESUMED` (not in background, not destroyed)
3. Automatically removes observers when the lifecycle is destroyed (no memory leaks)

```kotlin
implementation("androidx.lifecycle:lifecycle-livedata-ktx:2.8.7")
```

---

## MutableLiveData vs LiveData

```kotlin
// MutableLiveData: writable — lives in the ViewModel
private val _tasks = MutableLiveData<List<Task>>()

// LiveData: read-only — exposed to the UI
val tasks: LiveData<List<Task>> = _tasks
```

The pattern of exposing `LiveData` (read-only) while keeping `MutableLiveData` private prevents the UI from directly modifying state — the ViewModel controls all mutations.

---

## ViewModel with LiveData

```kotlin
class TaskViewModel : ViewModel() {

    private val _tasks = MutableLiveData<List<Task>>(emptyList())
    val tasks: LiveData<List<Task>> = _tasks

    private val _isLoading = MutableLiveData<Boolean>(false)
    val isLoading: LiveData<Boolean> = _isLoading

    private val _error = MutableLiveData<String?>()
    val error: LiveData<String?> = _error

    private val taskList = mutableListOf<Task>()

    fun addTask(title: String, description: String = "") {
        taskList.add(Task(id = System.currentTimeMillis(), title = title, description = description))
        _tasks.value = taskList.toList()  // setValue — must be on main thread
    }

    fun toggleComplete(taskId: Long) {
        val index = taskList.indexOfFirst { it.id == taskId }
        if (index != -1) {
            taskList[index] = taskList[index].copy(isCompleted = !taskList[index].isCompleted)
            _tasks.value = taskList.toList()
        }
    }

    fun deleteTask(taskId: Long) {
        taskList.removeAll { it.id == taskId }
        _tasks.value = taskList.toList()
    }

    fun loadFromNetwork() {
        _isLoading.value = true
        // From a background thread:
        // _tasks.postValue(fetchedTasks)
        // _isLoading.postValue(false)
    }
}
```

---

## Observing LiveData in Activity/Fragment

```kotlin
class MainActivity : AppCompatActivity() {

    private val viewModel: TaskViewModel by viewModels()
    private lateinit var binding: ActivityMainBinding
    private lateinit var adapter: TaskAdapter

    override fun onCreate(savedInstanceState: Bundle?) {
        super.onCreate(savedInstanceState)
        binding = ActivityMainBinding.inflate(layoutInflater)
        setContentView(binding.root)

        setupRecyclerView()
        observeViewModel()

        binding.fab.setOnClickListener { showAddDialog() }
    }

    private fun observeViewModel() {
        // observe() — automatically lifecycle-aware
        viewModel.tasks.observe(this) { taskList ->
            adapter.submitList(taskList)
            binding.tvEmpty.isVisible = taskList.isEmpty()
        }

        viewModel.isLoading.observe(this) { isLoading ->
            binding.progressBar.isVisible = isLoading
        }

        viewModel.error.observe(this) { error ->
            error?.let {
                Snackbar.make(binding.root, it, Snackbar.LENGTH_LONG).show()
            }
        }
    }
}
```

In Fragment, use `viewLifecycleOwner` instead of `this`:

```kotlin
override fun onViewCreated(view: View, savedInstanceState: Bundle?) {
    super.onViewCreated(view, savedInstanceState)

    // CORRECT — tied to the fragment's view lifecycle
    viewModel.tasks.observe(viewLifecycleOwner) { taskList ->
        adapter.submitList(taskList)
    }
}
```

> **Never use `this` (the Fragment) as the lifecycle owner in `onViewCreated`**. Use `viewLifecycleOwner`. Fragments can exist without a view (e.g., in the back stack), causing duplicate observers.

---

## Java Version

```java
// ViewModel
public class TaskViewModel extends ViewModel {
    private final MutableLiveData<List<Task>> tasks =
        new MutableLiveData<>(Collections.emptyList());

    public LiveData<List<Task>> getTasks() {
        return tasks;
    }

    public void addTask(String title) {
        List<Task> current = new ArrayList<>(
            tasks.getValue() != null ? tasks.getValue() : Collections.emptyList()
        );
        current.add(new Task(System.currentTimeMillis(), title, "", false));
        tasks.setValue(current);
    }
}

// Activity
viewModel.getTasks().observe(this, taskList -> {
    adapter.submitList(taskList);
});
```

---

## `setValue` vs `postValue`

| Method | Thread | Use When |
|--------|--------|----------|
| `_tasks.value = data` | Main thread only | Updating from ViewModel directly |
| `_tasks.postValue(data)` | Any thread | Updating from a coroutine or background thread |

```kotlin
// From background thread — use postValue
viewModelScope.launch(Dispatchers.IO) {
    val data = repository.loadData()
    _tasks.postValue(data)  // switches to main thread automatically
}
```

---

## Transformations

### `map` — Transform emitted values

```kotlin
// Transform a list of Tasks into a count
val taskCount: LiveData<Int> = tasks.map { it.size }

// Transform a list to only completed tasks
val completedTasks: LiveData<List<Task>> = tasks.map { list ->
    list.filter { it.isCompleted }
}
```

### `switchMap` — React to another LiveData

```kotlin
private val _searchQuery = MutableLiveData<String>("")

val filteredTasks: LiveData<List<Task>> = _searchQuery.switchMap { query ->
    if (query.isBlank()) tasks
    else tasks.map { list ->
        list.filter { it.title.contains(query, ignoreCase = true) }
    }
}

fun setSearchQuery(query: String) {
    _searchQuery.value = query
}
```

### `MediatorLiveData` — Combine multiple sources

```kotlin
val combinedState: MediatorLiveData<UiState> = MediatorLiveData<UiState>().apply {
    addSource(tasks) { taskList ->
        value = UiState(tasks = taskList, isLoading = isLoading.value ?: false)
    }
    addSource(isLoading) { loading ->
        value = UiState(tasks = tasks.value ?: emptyList(), isLoading = loading)
    }
}
```

---

## One-Time Events with LiveData

LiveData re-delivers the last value to new observers. For one-time events (navigation, Snackbar), this is a problem:

```kotlin
// A wrapper that marks the value as "consumed"
class Event<out T>(private val content: T) {
    private var hasBeenHandled = false

    fun getContentIfNotHandled(): T? {
        return if (hasBeenHandled) null
        else {
            hasBeenHandled = true
            content
        }
    }
}

// In ViewModel
private val _navigateToDetail = MutableLiveData<Event<Long>>()
val navigateToDetail: LiveData<Event<Long>> = _navigateToDetail

fun onTaskClicked(taskId: Long) {
    _navigateToDetail.value = Event(taskId)
}

// In Activity/Fragment
viewModel.navigateToDetail.observe(viewLifecycleOwner) { event ->
    event.getContentIfNotHandled()?.let { taskId ->
        // Navigate to detail — called only once
        findNavController().navigate(
            TaskListFragmentDirections.actionToDetail(taskId)
        )
    }
}
```

> **Note:** In modern code, `StateFlow` (Chapter 3) and `Channel`/`SharedFlow` are better alternatives for one-time events than the `Event` wrapper pattern.

---

## LiveData vs StateFlow — When to Use Which

| | LiveData | StateFlow |
|--|----------|-----------|
| Coroutine-native | No | Yes |
| Lifecycle-aware (built-in) | Yes | Needs `repeatOnLifecycle` |
| Initial value required | No | Yes |
| Kotlin-only | No | Yes |
| Android-free (pure Kotlin) | No | Yes |

Use **LiveData** when you want the simplest solution with built-in lifecycle handling.
Use **StateFlow** when you're already in coroutines and want a pure Kotlin solution (Chapter 3).

---

## Common Mistakes

### Mistake 1: Using `this` (Fragment) instead of `viewLifecycleOwner`

```kotlin
// WRONG — causes duplicate observers if fragment is recreated
viewModel.tasks.observe(this) { ... }

// CORRECT
viewModel.tasks.observe(viewLifecycleOwner) { ... }
```

### Mistake 2: Calling `setValue` from a background thread

```kotlin
// CRASH — setValue must be called on main thread
background { _tasks.value = result }

// CORRECT
background { _tasks.postValue(result) }
```

### Mistake 3: Exposing `MutableLiveData` publicly

```kotlin
// WRONG — UI can mutate state directly
val tasks = MutableLiveData<List<Task>>()

// CORRECT — UI gets a read-only view
private val _tasks = MutableLiveData<List<Task>>()
val tasks: LiveData<List<Task>> = _tasks
```

---

## Interview Questions

**Q1: Why is LiveData lifecycle-aware and why does that matter?**

> LiveData only notifies active observers (those whose lifecycle is in STARTED or RESUMED state). This prevents updating a destroyed UI and automatically clears observers when the lifecycle owner is destroyed — eliminating memory leaks.

**Q2: What is the difference between `observe()` and `observeForever()`?**

> `observe()` takes a `LifecycleOwner` and auto-removes the observer when the lifecycle is destroyed. `observeForever()` never auto-removes — you must manually call `removeObserver()` or you'll leak memory.

**Q3: Why expose `LiveData` instead of `MutableLiveData` from a ViewModel?**

> To enforce unidirectional data flow. The ViewModel is the only source of truth for state changes. Exposing a read-only `LiveData` prevents UI components from bypassing ViewModel business logic and mutating state directly.

**Q4: When would you use `Transformations.switchMap`?**

> When the source LiveData you want to observe depends on another LiveData value. For example, loading a different Room query when the user's search input (itself a LiveData) changes.

---

## Summary

- LiveData is a lifecycle-aware observable — automatically stops delivering when the observer is inactive
- Use `MutableLiveData` privately in ViewModel; expose `LiveData` publicly
- In Fragments, always use `viewLifecycleOwner`, never `this`
- Use `setValue` from the main thread, `postValue` from background threads
- Use `map`, `switchMap`, and `MediatorLiveData` for derived/combined state

**Next:** [Chapter 3 — StateFlow](./03-stateflow.md)
