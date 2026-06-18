# Chapter 7: Introduction to Jetpack Compose

## What Is Jetpack Compose?

Jetpack Compose is Android's modern **declarative UI toolkit**. Instead of describing UI in XML and imperatively updating it, you write Kotlin functions that describe what the UI should look like for a given state — and Compose automatically recomposes (re-renders) when state changes.

```kotlin
implementation("androidx.compose.ui:ui:1.7.6")
implementation("androidx.compose.material3:material3:1.3.1")
implementation("androidx.compose.ui:ui-tooling-preview:1.7.6")
implementation("androidx.activity:activity-compose:1.9.3")
implementation("androidx.lifecycle:lifecycle-viewmodel-compose:2.8.7")
```

---

## XML vs Compose — The Core Difference

**XML (Imperative):**
```kotlin
// You tell the system HOW to change
if (isLoading) {
    progressBar.visibility = View.VISIBLE
    recyclerView.visibility = View.GONE
} else {
    progressBar.visibility = View.GONE
    recyclerView.visibility = View.VISIBLE
}
```

**Compose (Declarative):**
```kotlin
// You describe WHAT it should look like
if (isLoading) {
    CircularProgressIndicator()
} else {
    TaskList(tasks)
}
// Compose figures out what changed and updates the UI
```

---

## Your First Composable

```kotlin
import androidx.compose.material3.*
import androidx.compose.runtime.*
import androidx.compose.foundation.layout.*
import androidx.compose.ui.Modifier
import androidx.compose.ui.unit.dp

// A @Composable function is a UI component
@Composable
fun Greeting(name: String) {
    Text(
        text = "Hello, $name!",
        style = MaterialTheme.typography.headlineMedium
    )
}

// Preview in Android Studio — no need to run on a device
@Preview(showBackground = true)
@Composable
fun GreetingPreview() {
    MyAppTheme {
        Greeting("World")
    }
}
```

---

## Setting Up Compose in an Activity

```kotlin
class MainActivity : ComponentActivity() {

    override fun onCreate(savedInstanceState: Bundle?) {
        super.onCreate(savedInstanceState)

        setContent {
            // Your app's theme wraps everything
            MyAppTheme {
                // Surface = background color from theme
                Surface(
                    modifier = Modifier.fillMaxSize(),
                    color = MaterialTheme.colorScheme.background
                ) {
                    TaskListScreen()
                }
            }
        }
    }
}
```

---

## State in Compose

Compose recomposes (re-runs) composables when their state changes. State is declared with `remember` and `mutableStateOf`:

```kotlin
@Composable
fun Counter() {
    // 'remember' keeps the state across recompositions
    var count by remember { mutableStateOf(0) }

    Column(
        horizontalAlignment = Alignment.CenterHorizontally,
        modifier = Modifier.padding(16.dp)
    ) {
        Text(
            text = "Count: $count",
            style = MaterialTheme.typography.headlineLarge
        )

        Spacer(modifier = Modifier.height(8.dp))

        Button(onClick = { count++ }) {
            Text("Increment")
        }
    }
}
```

---

## Connecting Compose to ViewModel

```kotlin
@Composable
fun TaskListScreen(
    viewModel: TaskViewModel = viewModel()
) {
    val uiState by viewModel.uiState.collectAsStateWithLifecycle()

    TaskListContent(
        tasks = uiState.tasks,
        isLoading = uiState.isLoading,
        onAddTask = { title -> viewModel.addTask(title) },
        onToggleComplete = { task -> viewModel.toggleComplete(task) },
        onDeleteTask = { task -> viewModel.deleteTask(task) }
    )
}

// Separate "content" composable — doesn't know about ViewModel, easy to test/preview
@Composable
fun TaskListContent(
    tasks: List<Task>,
    isLoading: Boolean,
    onAddTask: (String) -> Unit,
    onToggleComplete: (Task) -> Unit,
    onDeleteTask: (Task) -> Unit
) {
    if (isLoading) {
        Box(modifier = Modifier.fillMaxSize(), contentAlignment = Alignment.Center) {
            CircularProgressIndicator()
        }
    } else {
        LazyColumn(
            modifier = Modifier.fillMaxSize(),
            contentPadding = PaddingValues(16.dp),
            verticalArrangement = Arrangement.spacedBy(8.dp)
        ) {
            items(tasks, key = { it.id }) { task ->
                TaskItem(
                    task = task,
                    onToggleComplete = { onToggleComplete(task) },
                    onDelete = { onDeleteTask(task) }
                )
            }
        }
    }
}
```

---

## Common Composables

### Layouts

```kotlin
// Vertical stack
Column(
    modifier = Modifier.padding(16.dp),
    verticalArrangement = Arrangement.spacedBy(8.dp),
    horizontalAlignment = Alignment.CenterHorizontally
) {
    Text("Item 1")
    Text("Item 2")
}

// Horizontal stack
Row(
    modifier = Modifier.fillMaxWidth(),
    horizontalArrangement = Arrangement.SpaceBetween,
    verticalAlignment = Alignment.CenterVertically
) {
    Text("Left")
    Text("Right")
}

// Stack (like FrameLayout)
Box(modifier = Modifier.fillMaxSize()) {
    Image(painter = painterResource(R.drawable.bg), contentDescription = null,
        modifier = Modifier.fillMaxSize())
    Text("Overlay", modifier = Modifier.align(Alignment.Center))
}
```

### Lists

```kotlin
// LazyColumn = RecyclerView equivalent
LazyColumn(
    contentPadding = PaddingValues(horizontal = 16.dp, vertical = 8.dp),
    verticalArrangement = Arrangement.spacedBy(8.dp)
) {
    item { Header() }

    items(
        items = tasks,
        key = { task -> task.id }  // Stable keys for animations
    ) { task ->
        TaskItem(task = task)
    }

    item { Footer() }
}

// LazyRow = horizontal list
LazyRow(horizontalArrangement = Arrangement.spacedBy(8.dp)) {
    items(categories) { category ->
        CategoryChip(category)
    }
}
```

### Input

```kotlin
@Composable
fun SearchBar(query: String, onQueryChange: (String) -> Unit) {
    OutlinedTextField(
        value = query,
        onValueChange = onQueryChange,
        label = { Text("Search") },
        leadingIcon = { Icon(Icons.Default.Search, contentDescription = null) },
        trailingIcon = {
            if (query.isNotEmpty()) {
                IconButton(onClick = { onQueryChange("") }) {
                    Icon(Icons.Default.Clear, contentDescription = "Clear")
                }
            }
        },
        singleLine = true,
        modifier = Modifier.fillMaxWidth()
    )
}
```

---

## Task Item Composable

```kotlin
@Composable
fun TaskItem(
    task: Task,
    onToggleComplete: () -> Unit,
    onDelete: () -> Unit,
    modifier: Modifier = Modifier
) {
    Card(
        modifier = modifier.fillMaxWidth(),
        elevation = CardDefaults.cardElevation(defaultElevation = 2.dp)
    ) {
        Row(
            modifier = Modifier
                .fillMaxWidth()
                .padding(16.dp),
            verticalAlignment = Alignment.CenterVertically
        ) {
            Checkbox(
                checked = task.isCompleted,
                onCheckedChange = { onToggleComplete() }
            )

            Spacer(modifier = Modifier.width(12.dp))

            Column(modifier = Modifier.weight(1f)) {
                Text(
                    text = task.title,
                    style = MaterialTheme.typography.titleMedium,
                    textDecoration = if (task.isCompleted)
                        TextDecoration.LineThrough else TextDecoration.None
                )
                if (task.description.isNotBlank()) {
                    Text(
                        text = task.description,
                        style = MaterialTheme.typography.bodySmall,
                        color = MaterialTheme.colorScheme.onSurfaceVariant
                    )
                }
            }

            IconButton(onClick = onDelete) {
                Icon(
                    imageVector = Icons.Default.Delete,
                    contentDescription = "Delete task",
                    tint = MaterialTheme.colorScheme.error
                )
            }
        }
    }
}
```

---

## Compose Theme Setup

```kotlin
// ui/theme/Color.kt
val Purple80 = Color(0xFFD0BCFF)
val PurpleGrey80 = Color(0xFFCCC2DC)
val Pink80 = Color(0xFFEFB8C8)

// ui/theme/Theme.kt
@Composable
fun MyAppTheme(
    darkTheme: Boolean = isSystemInDarkTheme(),
    content: @Composable () -> Unit
) {
    val colorScheme = if (darkTheme) darkColorScheme(
        primary = Purple80,
        secondary = PurpleGrey80,
        tertiary = Pink80
    ) else lightColorScheme(
        primary = Purple40,
        secondary = PurpleGrey40,
        tertiary = Pink40
    )

    MaterialTheme(
        colorScheme = colorScheme,
        typography = Typography,
        content = content
    )
}
```

---

## XML vs Compose — When to Use Which

| | XML + Views | Jetpack Compose |
|--|-------------|-----------------|
| Maturity | Very mature, lots of examples | Stable since 2021, growing rapidly |
| Learning curve | Lower if familiar with Android | Higher initially |
| Custom UI | Complex | Simpler with Canvas API |
| Performance | Good | Equal or better |
| Testing | Requires Espresso | Built-in Compose testing APIs |
| Interop | Works alone | Can mix with XML via `ComposeView`/`AndroidView` |
| Industry trend | Legacy | **The future** |

For **new projects**, start with Compose. For **existing XML projects**, adopt Compose incrementally.

---

## Interoperability: Compose in XML, XML in Compose

**Compose inside XML layout:**

```xml
<androidx.compose.ui.platform.ComposeView
    android:id="@+id/composeView"
    android:layout_width="match_parent"
    android:layout_height="wrap_content" />
```

```kotlin
binding.composeView.setContent {
    MyAppTheme { TaskItem(task = task, ...) }
}
```

**XML View inside Compose:**

```kotlin
@Composable
fun LegacyMapView() {
    AndroidView(
        factory = { context ->
            MapView(context).apply {
                // configure the XML view
            }
        },
        update = { mapView ->
            // update when state changes
        }
    )
}
```

---

## Interview Questions

**Q1: What does "declarative UI" mean in the context of Compose?**

> You describe what the UI should look like for a given state, and Compose figures out the minimal changes needed. Contrast with imperative UI (XML + Views), where you explicitly tell the system how to change: show/hide views, update text, etc.

**Q2: What is recomposition?**

> When state that a Composable reads changes, Compose re-runs that Composable function (recomposes) to update the UI. Compose optimizes this — only Composables that depend on the changed state are recomposed.

**Q3: What is `remember` and why is it needed?**

> `remember` stores a value across recompositions. Without it, every time a Composable recomposes, the variable would be re-initialized to its default. `remember { mutableStateOf(0) }` survives recomposition.

---

## Summary

- Compose is declarative — describe WHAT not HOW
- `@Composable` functions are UI components; they're re-run (recomposed) when state changes
- State uses `remember { mutableStateOf(...) }` or ViewModel + `collectAsStateWithLifecycle()`
- `LazyColumn` is the Compose equivalent of RecyclerView
- Compose and XML views can interop via `ComposeView` and `AndroidView`
- Compose is the direction of Android UI — start with it for new projects

**Next:** [Mini Project — Notes App](./mini-project-notes-app.md)
