# Chapter 1: ViewModel

## The Problem ViewModel Solves

Android Activities and Fragments are destroyed and recreated on configuration changes (screen rotation, language change, multi-window resize). Without ViewModel, this causes:

1. Lost UI state — the user typed something, rotated, it's gone
2. Duplicate network requests — loading is re-triggered on every rotation
3. Memory leaks — background work holds a reference to the dead Activity

**Analogy:** Imagine a whiteboard (Activity) in a room. Every time you walk out and back in (configuration change), someone erases the whiteboard. ViewModel is a separate locker — it survives the room being cleared.

---

## What ViewModel Is

`ViewModel` is an AndroidX lifecycle-aware class that holds and manages UI-related data. It survives configuration changes and is destroyed only when the Activity/Fragment is permanently finished.

```kotlin
implementation("androidx.lifecycle:lifecycle-viewmodel-ktx:2.8.7")
```

---

## ViewModel Lifecycle

```
Activity Created ──────────────────────────────────► Activity Destroyed (back press)
       │                                                      │
       ▼                                                      ▼
  ViewModel Created                                    ViewModel.onCleared()
       │
  Rotation ──► Activity Destroyed ──► Activity Recreated
       │             │
       └─────────────┘
         ViewModel SURVIVES
```

The ViewModel lives as long as the Activity's **scope** is alive — not the individual Activity instance.

---

## Creating a ViewModel

```kotlin
import androidx.lifecycle.ViewModel

class TaskViewModel : ViewModel() {

    // State lives here — survives rotation
    private val _tasks = mutableListOf<Task>()
    val tasks: List<Task> get() = _tasks

    private var nextId = 1L

    fun addTask(title: String, description: String = "") {
        _tasks.add(Task(id = nextId++, title = title, description = description))
    }

    fun toggleComplete(taskId: Long) {
        val index = _tasks.indexOfFirst { it.id == taskId }
        if (index != -1) {
            _tasks[index] = _tasks[index].copy(
                isCompleted = !_tasks[index].isCompleted
            )
        }
    }

    fun deleteTask(taskId: Long) {
        _tasks.removeAll { it.id == taskId }
    }

    // Called when ViewModel is destroyed — clean up resources here
    override fun onCleared() {
        super.onCleared()
        // Cancel coroutines, close connections, etc.
    }
}
```

---

## Connecting ViewModel to an Activity

```kotlin
import androidx.activity.viewModels

class MainActivity : AppCompatActivity() {

    // The 'by viewModels()' delegate handles creation and caching
    private val viewModel: TaskViewModel by viewModels()

    private lateinit var binding: ActivityMainBinding

    override fun onCreate(savedInstanceState: Bundle?) {
        super.onCreate(savedInstanceState)
        binding = ActivityMainBinding.inflate(layoutInflater)
        setContentView(binding.root)

        // ViewModel holds the data — just display it
        displayTasks()

        binding.fab.setOnClickListener {
            viewModel.addTask("New Task")
            displayTasks()
        }
    }

    private fun displayTasks() {
        // For now, manually refresh — Chapter 2 (LiveData) fixes this
        adapter.submitList(viewModel.tasks.toList())
    }
}
```

---

## Connecting ViewModel to a Fragment

```kotlin
import androidx.fragment.app.viewModels
import androidx.fragment.app.activityViewModels

class TaskListFragment : Fragment() {

    // Scoped to this Fragment's lifecycle
    private val viewModel: TaskViewModel by viewModels()

    // Scoped to the parent Activity — SHARED between fragments
    private val sharedViewModel: SharedViewModel by activityViewModels()

    override fun onViewCreated(view: View, savedInstanceState: Bundle?) {
        super.onViewCreated(view, savedInstanceState)

        binding.btnAdd.setOnClickListener {
            viewModel.addTask("New task")
        }
    }
}
```

`activityViewModels()` is how two fragments share data — both get the same ViewModel instance.

---

## Java Version

```java
public class TaskViewModel extends ViewModel {

    private final MutableList<Task> tasks = new ArrayList<>();

    public List<Task> getTasks() {
        return Collections.unmodifiableList(tasks);
    }

    public void addTask(String title) {
        tasks.add(new Task(System.currentTimeMillis(), title, "", false));
    }

    @Override
    protected void onCleared() {
        super.onCleared();
    }
}

// In Activity:
public class MainActivity extends AppCompatActivity {

    private TaskViewModel viewModel;

    @Override
    protected void onCreate(Bundle savedInstanceState) {
        super.onCreate(savedInstanceState);

        viewModel = new ViewModelProvider(this).get(TaskViewModel.class);

        // Use viewModel.getTasks()
    }
}
```

---

## ViewModel with Constructor Parameters — `ViewModelFactory`

By default, ViewModel must have a no-arg constructor. For dependencies, use a factory:

```kotlin
class TaskViewModel(
    private val repository: TaskRepository
) : ViewModel() {

    fun loadTasks() = repository.getAllTasks()
}

class TaskViewModelFactory(
    private val repository: TaskRepository
) : ViewModelProvider.Factory {

    override fun <T : ViewModel> create(modelClass: Class<T>): T {
        if (modelClass.isAssignableFrom(TaskViewModel::class.java)) {
            @Suppress("UNCHECKED_CAST")
            return TaskViewModel(repository) as T
        }
        throw IllegalArgumentException("Unknown ViewModel: ${modelClass.name}")
    }
}

// In Activity:
private val viewModel: TaskViewModel by viewModels {
    TaskViewModelFactory(TaskRepository(database.taskDao()))
}
```

> **Note:** With Hilt (Level 3), you no longer write factories manually — Hilt generates them.

---

## What Should and Shouldn't Go in ViewModel

| Should Go In ViewModel | Should NOT Go In ViewModel |
|-----------------------|---------------------------|
| UI state (list data, loading flag, error state) | `Context` references |
| Business logic | `View` references |
| LiveData / StateFlow | `Activity` or `Fragment` references |
| Coroutine jobs | Resources requiring Context (strings, drawables) |
| Repository calls | Android framework classes |

### Why not Context in ViewModel?

```kotlin
// DANGEROUS — Context leaks the Activity after rotation
class BadViewModel(private val context: Context) : ViewModel()

// SAFE — if you must have Context, use ApplicationContext via AndroidViewModel
class SafeViewModel(application: Application) : AndroidViewModel(application) {
    private val appContext = getApplication<Application>()
}
```

---

## `SavedStateHandle` — Surviving Process Death

`by viewModels()` keeps data through rotations but NOT through process death (system kills your app in the background). For process-death resilience, use `SavedStateHandle`:

```kotlin
class TaskViewModel(
    private val savedStateHandle: SavedStateHandle
) : ViewModel() {

    var searchQuery: String
        get() = savedStateHandle.get<String>("search_query") ?: ""
        set(value) { savedStateHandle.set("search_query", value) }
}

// Inject SavedStateHandle automatically:
private val viewModel: TaskViewModel by viewModels()
// SavedStateHandle is injected by the 'by viewModels()' delegate automatically
// if the constructor parameter is named 'savedStateHandle'
```

---

## Common Mistakes

### Mistake 1: Storing Activity Context in ViewModel

```kotlin
// MEMORY LEAK — ViewModel outlives the Activity
class BadViewModel(private val activity: Activity) : ViewModel()
```

### Mistake 2: Creating ViewModel directly with `TaskViewModel()`

```kotlin
// WRONG — creates a new instance every time, not cached
val viewModel = TaskViewModel()

// CORRECT — framework manages lifecycle and caching
val viewModel: TaskViewModel by viewModels()
```

### Mistake 3: Putting UI logic in ViewModel

ViewModel should not know about Views, Fragments, or how data is displayed. It manages data — the UI observes it.

---

## Interview Questions

**Q1: What is the purpose of ViewModel?**

> ViewModel stores and manages UI-related data in a lifecycle-aware way. It survives configuration changes (like screen rotation), preventing data loss and redundant work.

**Q2: What is the difference between `by viewModels()` and `by activityViewModels()`?**

> `by viewModels()` creates a ViewModel scoped to the Fragment. `by activityViewModels()` creates one scoped to the Activity — meaning all Fragments in that Activity share the same instance. Use the latter for cross-fragment communication.

**Q3: Why should you not pass `Context` to a ViewModel?**

> ViewModel survives configuration changes, but Activity-scoped `Context` is destroyed on rotation. Holding a reference to it causes a memory leak. Use `AndroidViewModel` if you need ApplicationContext, or pass context-free data.

**Q4: When is `ViewModel.onCleared()` called?**

> When the Activity finishes permanently (e.g., the user presses back or `finish()` is called). Not on rotation. Use it to cancel coroutines and clean up resources.

---

## Summary

- ViewModel stores UI state and survives configuration changes
- Access via `by viewModels()` (Fragment-scoped) or `by activityViewModels()` (Activity-scoped)
- Never store Context, Views, or Activity references in ViewModel
- Use `ViewModelFactory` for constructor parameters (until Hilt)
- Use `SavedStateHandle` for process-death resilience

**Next:** [Chapter 2 — LiveData](./02-livedata.md)
