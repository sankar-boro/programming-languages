# Chapter 4: Lifecycle

## What Is a Lifecycle?

An Android Lifecycle is the set of states an Activity or Fragment moves through from creation to destruction. AndroidX provides the `Lifecycle` API to observe and respond to these state changes in a clean, decoupled way.

---

## Activity Lifecycle States

```
        User launches app
               │
               ▼
         ┌─────────┐
         │ CREATED │  ◄── onCreate()
         └────┬────┘
              │
              ▼
         ┌─────────┐
         │ STARTED │  ◄── onStart()
         └────┬────┘
              │
              ▼
         ┌─────────┐
         │ RESUMED │  ◄── onResume()  ← App is visible and interactive
         └────┬────┘
              │
         User presses Home / another app comes to foreground
              │
              ▼
         ┌─────────┐
         │ STARTED │  ◄── onPause() → back to STARTED
         └────┬────┘
              │
         App no longer visible
              │
              ▼
         ┌─────────┐
         │ CREATED │  ◄── onStop() → back to CREATED
         └────┬────┘
              │
         System kills app / user navigates away
              │
              ▼
         ┌───────────┐
         │ DESTROYED │  ◄── onDestroy()
         └───────────┘
```

---

## Lifecycle Callbacks

```kotlin
class MainActivity : AppCompatActivity() {

    override fun onCreate(savedInstanceState: Bundle?) {
        super.onCreate(savedInstanceState)
        // First creation — set up views, bind data
    }

    override fun onStart() {
        super.onStart()
        // Activity becomes visible — start animations, register receivers
    }

    override fun onResume() {
        super.onResume()
        // Activity is interactive — start camera, begin foreground updates
    }

    override fun onPause() {
        super.onPause()
        // Another Activity partially covers this one — pause heavy work
    }

    override fun onStop() {
        super.onStop()
        // Activity is no longer visible — release camera, unregister receivers
    }

    override fun onDestroy() {
        super.onDestroy()
        // Activity is finishing — clean up final resources
    }

    override fun onSaveInstanceState(outState: Bundle) {
        super.onSaveInstanceState(outState)
        // Save temporary UI state (e.g., scroll position, input text)
        outState.putString("draft", binding.etInput.text.toString())
    }

    override fun onRestoreInstanceState(savedInstanceState: Bundle) {
        super.onRestoreInstanceState(savedInstanceState)
        // Restore UI state after recreation
        binding.etInput.setText(savedInstanceState.getString("draft"))
    }
}
```

---

## Fragment Lifecycle (Additional Callbacks)

Fragments have a more complex lifecycle due to their view lifecycle:

```kotlin
class TaskListFragment : Fragment() {

    override fun onAttach(context: Context) { super.onAttach(context) }
    override fun onCreate(savedInstanceState: Bundle?) { super.onCreate(savedInstanceState) }

    override fun onCreateView(
        inflater: LayoutInflater,
        container: ViewGroup?,
        savedInstanceState: Bundle?
    ): View {
        _binding = FragmentTaskListBinding.inflate(inflater, container, false)
        return binding.root
    }

    override fun onViewCreated(view: View, savedInstanceState: Bundle?) {
        super.onViewCreated(view, savedInstanceState)
        // Set up UI, start observing — view is ready
    }

    override fun onStart() { super.onStart() }
    override fun onResume() { super.onResume() }
    override fun onPause() { super.onPause() }
    override fun onStop() { super.onStop() }

    override fun onDestroyView() {
        super.onDestroyView()
        _binding = null  // CRITICAL: clear binding to avoid memory leak
    }

    override fun onDestroy() { super.onDestroy() }
    override fun onDetach() { super.onDetach() }
}
```

**Key rule:** Clear `_binding` in `onDestroyView()` — not `onDestroy()`. The view can be destroyed while the Fragment object lives on (in the back stack).

---

## LifecycleOwner and LifecycleObserver

Instead of putting everything in lifecycle callbacks inside Activity/Fragment, move lifecycle-aware logic to a separate class using `LifecycleObserver`.

```kotlin
import androidx.lifecycle.DefaultLifecycleObserver
import androidx.lifecycle.LifecycleOwner

class LocationTracker(
    private val context: Context
) : DefaultLifecycleObserver {

    private var locationManager: LocationManager? = null

    override fun onStart(owner: LifecycleOwner) {
        // Start tracking when Activity is visible
        locationManager = context.getSystemService(Context.LOCATION_SERVICE) as LocationManager
        // start updates...
    }

    override fun onStop(owner: LifecycleOwner) {
        // Stop tracking when Activity goes to background
        locationManager?.removeUpdates(locationListener)
        locationManager = null
    }
}

// In Activity — clean, no tracking code in Activity callbacks
class MapsActivity : AppCompatActivity() {

    private val locationTracker by lazy { LocationTracker(this) }

    override fun onCreate(savedInstanceState: Bundle?) {
        super.onCreate(savedInstanceState)
        lifecycle.addObserver(locationTracker)
        // Tracker responds to start/stop automatically
    }
}
```

---

## `viewModelScope` and `lifecycleScope`

These coroutine scopes are tied to their respective lifecycle:

```kotlin
class TaskViewModel : ViewModel() {
    fun loadData() {
        viewModelScope.launch {
            // Cancelled when ViewModel.onCleared() is called
            val data = repository.fetchData()
            _uiState.update { it.copy(tasks = data) }
        }
    }
}

class MainActivity : AppCompatActivity() {
    override fun onCreate(savedInstanceState: Bundle?) {
        super.onCreate(savedInstanceState)

        lifecycleScope.launch {
            // Cancelled when Activity is DESTROYED
            delay(2000)
            showWelcomeMessage()
        }
    }
}
```

---

## `ProcessLifecycleOwner` — Whole App Lifecycle

```kotlin
class MyApplication : Application() {

    override fun onCreate() {
        super.onCreate()
        ProcessLifecycleOwner.get().lifecycle.addObserver(AppLifecycleObserver())
    }
}

class AppLifecycleObserver : DefaultLifecycleObserver {

    override fun onStart(owner: LifecycleOwner) {
        // App came to foreground
        Log.d("AppLifecycle", "App is in foreground")
    }

    override fun onStop(owner: LifecycleOwner) {
        // App went to background
        Log.d("AppLifecycle", "App is in background")
    }
}
```

Add dependency:
```kotlin
implementation("androidx.lifecycle:lifecycle-process:2.8.7")
```

---

## `OnBackPressedDispatcher` — Modern Back Navigation

```kotlin
class FormActivity : AppCompatActivity() {

    private val backPressedCallback = object : OnBackPressedCallback(enabled = true) {
        override fun handleOnBackPressed() {
            if (hasUnsavedChanges()) {
                showUnsavedChangesDialog()
            } else {
                // Let the system handle it (pop back stack)
                isEnabled = false
                onBackPressedDispatcher.onBackPressed()
            }
        }
    }

    override fun onCreate(savedInstanceState: Bundle?) {
        super.onCreate(savedInstanceState)
        onBackPressedDispatcher.addCallback(this, backPressedCallback)
    }
}
```

Never override `onBackPressed()` directly — it's deprecated. Use `OnBackPressedDispatcher`.

---

## Common Lifecycle Pitfalls

### Pitfall 1: Memory leak from Fragment binding

```kotlin
// WRONG — binding holds a reference to views after onDestroyView
private var binding: FragmentTaskBinding? = null

override fun onDestroyView() {
    super.onDestroyView()
    // Missing: binding = null
}
```

### Pitfall 2: Starting work in `onResume` that should be in `onCreate`

```kotlin
// WRONG — called every time the screen comes back to foreground
override fun onResume() {
    super.onResume()
    loadData()  // Makes a network call every time!
}

// CORRECT — called only on first creation
override fun onCreate(savedInstanceState: Bundle?) {
    super.onCreate(savedInstanceState)
    if (savedInstanceState == null) {
        loadData()
    }
}
```

### Pitfall 3: Doing heavy work in `onCreate`

`onCreate` runs on the main thread. Network calls, database reads, and large computations belong in `viewModelScope.launch` or `lifecycleScope.launch`.

---

## Interview Questions

**Q1: What is the difference between onStop and onDestroy?**

> `onStop` is called when the Activity is no longer visible (e.g., user goes home) but might return. `onDestroy` is called when the Activity is permanently finishing — the user pressed back, or the system killed it. You can release UI resources in `onStop`; use `onDestroy` for final cleanup.

**Q2: When should you use `DefaultLifecycleObserver` instead of overriding lifecycle methods directly?**

> When you want to move lifecycle-aware logic out of Activity/Fragment into a separate, reusable class. This keeps activities focused on UI coordination and makes the lifecycle-aware component independently testable.

**Q3: What is `viewLifecycleOwner` in a Fragment?**

> `viewLifecycleOwner` is the `LifecycleOwner` tied to the Fragment's view (created in `onCreateView`, destroyed in `onDestroyView`). It is distinct from the Fragment's own lifecycle. Using it for LiveData/StateFlow observation prevents memory leaks when the view is destroyed but the Fragment lives on in the back stack.

---

## Summary

- Activity lifecycle: CREATED → STARTED → RESUMED → STARTED → CREATED → DESTROYED
- Fragments have an additional view lifecycle — always clear binding in `onDestroyView()`
- Use `DefaultLifecycleObserver` to extract lifecycle logic from Activity/Fragment
- Use `viewModelScope` and `lifecycleScope` for coroutines — auto-cancelled on lifecycle end
- Use `OnBackPressedDispatcher` for back navigation — never override `onBackPressed()`

**Next:** [Chapter 5 — Navigation Component](./05-navigation-component.md)
