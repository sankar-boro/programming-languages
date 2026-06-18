# Chapter 5: Navigation Component

## What Is the Navigation Component?

The Navigation Component is an AndroidX library for managing in-app navigation — Fragment transactions, back stack, deep links, and safe argument passing — with a visual graph editor and type-safe APIs.

```kotlin
implementation("androidx.navigation:navigation-fragment-ktx:2.8.4")
implementation("androidx.navigation:navigation-ui-ktx:2.8.4")
```

Also add the Safe Args Gradle plugin for type-safe navigation arguments:

```kotlin
// build.gradle.kts (project level)
plugins {
    alias(libs.plugins.navigation.safeargs.kotlin) apply false
}

// app/build.gradle.kts
plugins {
    alias(libs.plugins.navigation.safeargs.kotlin)
}

// libs.versions.toml
[plugins]
navigation-safeargs-kotlin = { id = "androidx.navigation.safeargs.kotlin", version.ref = "navigationFragment" }
```

---

## Core Concepts

| Concept | What It Is |
|---------|-----------|
| NavGraph | XML file declaring all destinations and actions |
| NavHostFragment | The container where fragments are shown |
| NavController | The object you call to navigate |
| Destination | A Fragment or Activity in the graph |
| Action | A connection between two destinations |
| Argument | Typed data passed to a destination |

---

## Step 1: Create the Navigation Graph

Right-click `res` → New → Android Resource File → Resource type: Navigation → `nav_graph.xml`

```xml
<?xml version="1.0" encoding="utf-8"?>
<navigation
    xmlns:android="http://schemas.android.com/apk/res/android"
    xmlns:app="http://schemas.android.com/apk/res-auto"
    xmlns:tools="http://schemas.android.com/tools"
    android:id="@+id/nav_graph"
    app:startDestination="@id/taskListFragment">

    <fragment
        android:id="@+id/taskListFragment"
        android:name="com.yourname.app.TaskListFragment"
        android:label="Tasks"
        tools:layout="@layout/fragment_task_list">

        <!-- Action: navigate to detail, passing a task ID -->
        <action
            android:id="@+id/action_list_to_detail"
            app:destination="@id/taskDetailFragment"
            app:enterAnim="@anim/slide_in_right"
            app:exitAnim="@anim/slide_out_left"
            app:popEnterAnim="@anim/slide_in_left"
            app:popExitAnim="@anim/slide_out_right" />

    </fragment>

    <fragment
        android:id="@+id/taskDetailFragment"
        android:name="com.yourname.app.TaskDetailFragment"
        android:label="Task Detail"
        tools:layout="@layout/fragment_task_detail">

        <!-- Argument: the task ID passed to this fragment -->
        <argument
            android:name="taskId"
            app:argType="long" />

        <!-- Optional argument with default value -->
        <argument
            android:name="isEditing"
            app:argType="boolean"
            android:defaultValue="false" />

    </fragment>

    <fragment
        android:id="@+id/settingsFragment"
        android:name="com.yourname.app.SettingsFragment"
        android:label="Settings" />

</navigation>
```

---

## Step 2: Add NavHostFragment to Activity Layout

`activity_main.xml`:

```xml
<?xml version="1.0" encoding="utf-8"?>
<androidx.constraintlayout.widget.ConstraintLayout
    xmlns:android="http://schemas.android.com/apk/res/android"
    xmlns:app="http://schemas.android.com/apk/res-auto"
    android:layout_width="match_parent"
    android:layout_height="match_parent">

    <com.google.android.material.appbar.AppBarLayout
        android:id="@+id/appBarLayout"
        android:layout_width="match_parent"
        android:layout_height="wrap_content"
        app:layout_constraintTop_toTopOf="parent">

        <com.google.android.material.appbar.MaterialToolbar
            android:id="@+id/toolbar"
            android:layout_width="match_parent"
            android:layout_height="?attr/actionBarSize" />

    </com.google.android.material.appbar.AppBarLayout>

    <!-- This fragment hosts the navigation graph -->
    <androidx.fragment.app.FragmentContainerView
        android:id="@+id/navHostFragment"
        android:name="androidx.navigation.fragment.NavHostFragment"
        android:layout_width="0dp"
        android:layout_height="0dp"
        app:defaultNavHost="true"
        app:navGraph="@navigation/nav_graph"
        app:layout_constraintTop_toBottomOf="@id/appBarLayout"
        app:layout_constraintBottom_toBottomOf="parent"
        app:layout_constraintStart_toStartOf="parent"
        app:layout_constraintEnd_toEndOf="parent" />

    <com.google.android.material.bottomnavigation.BottomNavigationView
        android:id="@+id/bottomNav"
        android:layout_width="0dp"
        android:layout_height="wrap_content"
        app:menu="@menu/bottom_nav_menu"
        app:layout_constraintBottom_toBottomOf="parent"
        app:layout_constraintStart_toStartOf="parent"
        app:layout_constraintEnd_toEndOf="parent" />

</androidx.constraintlayout.widget.ConstraintLayout>
```

`app:defaultNavHost="true"` makes the system back button work with the navigation back stack.

---

## Step 3: Set Up NavController in MainActivity

```kotlin
class MainActivity : AppCompatActivity() {

    private lateinit var binding: ActivityMainBinding
    private lateinit var navController: NavController

    override fun onCreate(savedInstanceState: Bundle?) {
        super.onCreate(savedInstanceState)
        binding = ActivityMainBinding.inflate(layoutInflater)
        setContentView(binding.root)

        // Get NavController from the NavHostFragment
        val navHostFragment = supportFragmentManager
            .findFragmentById(R.id.navHostFragment) as NavHostFragment
        navController = navHostFragment.navController

        // Connect Toolbar with NavController
        val appBarConfig = AppBarConfiguration(
            // These destinations will show a hamburger icon, not a back arrow
            topLevelDestinationIds = setOf(
                R.id.taskListFragment,
                R.id.settingsFragment
            )
        )
        setupActionBarWithNavController(navController, appBarConfig)

        // Connect Bottom Navigation with NavController
        binding.bottomNav.setupWithNavController(navController)
    }

    override fun onSupportNavigateUp(): Boolean {
        return navController.navigateUp() || super.onSupportNavigateUp()
    }
}
```

---

## Step 4: Navigate Between Fragments

```kotlin
class TaskListFragment : Fragment() {

    override fun onViewCreated(view: View, savedInstanceState: Bundle?) {
        super.onViewCreated(view, savedInstanceState)

        adapter = TaskAdapter { task ->
            // Type-safe navigation using Safe Args generated directions
            val action = TaskListFragmentDirections
                .actionListToDetail(taskId = task.id, isEditing = false)
            findNavController().navigate(action)
        }

        binding.btnSettings.setOnClickListener {
            // Navigate by action ID
            findNavController().navigate(R.id.action_list_to_settings)
        }
    }
}
```

---

## Step 5: Receive Arguments in Destination Fragment

```kotlin
class TaskDetailFragment : Fragment() {

    // Safe Args generates this — type-safe, no manual bundle parsing
    private val args: TaskDetailFragmentArgs by navArgs()

    override fun onViewCreated(view: View, savedInstanceState: Bundle?) {
        super.onViewCreated(view, savedInstanceState)

        val taskId = args.taskId           // Long
        val isEditing = args.isEditing     // Boolean

        viewModel.loadTask(taskId)

        if (isEditing) {
            binding.etTitle.isEnabled = true
        }
    }
}
```

Without Safe Args (manual bundle):

```kotlin
// Sending — fragile, string key typo = runtime crash
val bundle = Bundle().apply { putLong("taskId", task.id) }
findNavController().navigate(R.id.taskDetailFragment, bundle)

// Receiving
val taskId = arguments?.getLong("taskId") ?: -1L
```

Always prefer Safe Args.

---

## Passing Results Back (Up the Back Stack)

```kotlin
// In detail fragment — set result before popping
findNavController().previousBackStackEntry?.savedStateHandle?.set("deleted", true)
findNavController().popBackStack()

// In list fragment — observe the result
override fun onViewCreated(view: View, savedInstanceState: Bundle?) {
    super.onViewCreated(view, savedInstanceState)

    val savedStateHandle = findNavController().currentBackStackEntry?.savedStateHandle
    savedStateHandle?.getLiveData<Boolean>("deleted")?.observe(viewLifecycleOwner) { deleted ->
        if (deleted) {
            Snackbar.make(binding.root, "Task deleted", Snackbar.LENGTH_SHORT).show()
        }
    }
}
```

---

## Deep Links

Declare in the nav graph:

```xml
<fragment android:id="@+id/taskDetailFragment" ...>
    <deepLink
        android:id="@+id/deepLink"
        app:uri="myapp://tasks/{taskId}" />
</fragment>
```

In AndroidManifest — add to the Activity:

```xml
<activity android:name=".MainActivity">
    <nav-graph android:value="@navigation/nav_graph" />
</activity>
```

Handle programmatic deep link:

```kotlin
val deepLinkRequest = NavDeepLinkRequest.Builder
    .fromUri("myapp://tasks/42".toUri())
    .build()
navController.navigate(deepLinkRequest)
```

---

## Common Mistakes

### Mistake 1: Calling `findNavController()` in `onCreate` (not yet attached)

```kotlin
// WRONG — NavController not available until view is created
override fun onCreate(...) {
    findNavController().navigate(...)
}

// CORRECT
override fun onViewCreated(...) {
    findNavController().navigate(...)
}
```

### Mistake 2: Navigating after the Fragment is detached (causes `IllegalStateException`)

```kotlin
// SAFE navigation — check fragment is still attached
if (isAdded) {
    findNavController().navigate(action)
}

// Even safer with lifecycle check
viewLifecycleOwner.lifecycleScope.launch {
    repeatOnLifecycle(Lifecycle.State.STARTED) {
        // Navigation inside here is always safe
    }
}
```

### Mistake 3: Multiple rapid clicks causing double navigation

```kotlin
// Debounce or check current destination before navigating
binding.btnDetail.setOnClickListener {
    if (findNavController().currentDestination?.id == R.id.taskListFragment) {
        findNavController().navigate(action)
    }
}
```

---

## Interview Questions

**Q1: What is Safe Args and why use it?**

> Safe Args is a Gradle plugin that generates type-safe classes for navigation arguments. Instead of manually putting/getting values from `Bundle` with string keys (error-prone), Safe Args generates `Directions` and `Args` classes with typed properties — compile-time safety.

**Q2: What does `app:defaultNavHost="true"` do?**

> It intercepts the system back button press and delegates it to the `NavController`, so pressing back pops the Fragment back stack correctly instead of finishing the Activity.

**Q3: How do you share data between two fragments using Navigation Component?**

> Use `NavController.previousBackStackEntry?.savedStateHandle` to pass results back up the stack. For forward navigation, use Safe Args. For bidirectional sharing, use `activityViewModels()` with a shared ViewModel.

---

## Summary

- Navigation Component manages Fragment transactions, back stack, deep links, and arguments
- Add `NavHostFragment` to your Activity layout, connect `NavController` to the Toolbar
- Use Safe Args plugin for type-safe argument passing between destinations
- Use `savedStateHandle` on `BackStackEntry` to pass results back to the previous screen
- `app:defaultNavHost="true"` wires up the system back button

**Next:** [Chapter 6 — Room Database](./06-room-database.md)
