# Chapter 4: androidx.core and androidx.appcompat

## Overview

These two libraries are the bedrock of every AndroidX project. They provide:

- **`androidx.core`** — Compatibility wrappers for Android OS APIs
- **`androidx.appcompat`** — Backward-compatible Activity, ActionBar, and UI components

---

## androidx.core

### What It Is

`androidx.core` provides helper classes that wrap Android OS APIs in a backward-compatible way. Instead of checking the OS version yourself, `core` handles it for you.

```kotlin
implementation("androidx.core:core-ktx:1.15.0")
```

The `-ktx` variant adds Kotlin extension functions on top of the core Java API — strongly preferred in Kotlin projects.

---

### `ContextCompat` — Safe Context Operations

```kotlin
import androidx.core.content.ContextCompat

// Safe color retrieval (handles API level differences internally)
val color = ContextCompat.getColor(context, R.color.primary)

// Check a permission
val granted = ContextCompat.checkSelfPermission(
    context,
    Manifest.permission.CAMERA
) == PackageManager.PERMISSION_GRANTED

// Get a drawable safely
val icon = ContextCompat.getDrawable(context, R.drawable.ic_home)
```

```java
// Java
int color = ContextCompat.getColor(context, R.color.primary);
boolean granted = ContextCompat.checkSelfPermission(context,
    Manifest.permission.CAMERA) == PackageManager.PERMISSION_GRANTED;
```

**Why not use `context.getColor()` directly?**
`context.getColor()` requires API 23+. `ContextCompat.getColor()` works on API 21+.

---

### `ViewCompat` — Safe View Operations

```kotlin
import androidx.core.view.ViewCompat
import androidx.core.view.WindowInsetsCompat

// Set elevation safely
ViewCompat.setElevation(myView, 8f)

// Handle window insets (keyboard, status bar, navigation bar)
ViewCompat.setOnApplyWindowInsetsListener(rootView) { view, insets ->
    val systemBars = insets.getInsets(WindowInsetsCompat.Type.systemBars())
    view.setPadding(systemBars.left, systemBars.top, systemBars.right, systemBars.bottom)
    insets
}
```

---

### `ActivityCompat` — Permission Requests

```kotlin
import androidx.core.app.ActivityCompat

// Request runtime permission
ActivityCompat.requestPermissions(
    activity,
    arrayOf(Manifest.permission.CAMERA),
    REQUEST_CODE_CAMERA
)

// Check if rationale should be shown
val shouldShow = ActivityCompat.shouldShowRequestPermissionRationale(
    activity,
    Manifest.permission.CAMERA
)
```

```java
// Java
ActivityCompat.requestPermissions(this,
    new String[]{Manifest.permission.CAMERA},
    REQUEST_CODE_CAMERA);
```

---

### `NotificationCompat` — Notifications That Work Everywhere

```kotlin
import androidx.core.app.NotificationCompat
import androidx.core.app.NotificationManagerCompat

val notification = NotificationCompat.Builder(context, CHANNEL_ID)
    .setSmallIcon(R.drawable.ic_notification)
    .setContentTitle("New Message")
    .setContentText("You have a new message from Alice")
    .setPriority(NotificationCompat.PRIORITY_DEFAULT)
    .setAutoCancel(true)
    .build()

NotificationManagerCompat.from(context).notify(NOTIFICATION_ID, notification)
```

---

### Kotlin Extensions from `core-ktx`

`core-ktx` adds idiomatic Kotlin extensions that make common patterns concise:

```kotlin
// core-ktx extensions

// View visibility
view.isVisible = true
view.isGone = false
view.isInvisible = true

// String to Uri
val uri = "https://example.com".toUri()

// Bundle creation
val bundle = bundleOf("key1" to "value", "key2" to 42)

// Color
val color = Color.parseColor("#FF5722")

// Handler post
Handler(Looper.getMainLooper()).post { updateUI() }
// With ktx:
view.post { updateUI() }
```

---

## androidx.appcompat

### What It Is

`androidx.appcompat` provides backward-compatible versions of Android UI components — particularly `AppCompatActivity`, which brings modern features (Toolbar, night mode, vector drawables) to older Android versions.

```kotlin
implementation("androidx.appcompat:appcompat:1.7.0")
```

---

### `AppCompatActivity`

All activities in modern Android apps extend `AppCompatActivity`:

```kotlin
import androidx.appcompat.app.AppCompatActivity
import android.os.Bundle

class MainActivity : AppCompatActivity() {

    override fun onCreate(savedInstanceState: Bundle?) {
        super.onCreate(savedInstanceState)
        setContentView(R.layout.activity_main)

        // Set up the toolbar as the ActionBar
        val toolbar = binding.toolbar
        setSupportActionBar(toolbar)
        supportActionBar?.title = "My App"
    }

    // Handle overflow menu
    override fun onCreateOptionsMenu(menu: Menu): Boolean {
        menuInflater.inflate(R.menu.main_menu, menu)
        return true
    }

    override fun onOptionsItemSelected(item: MenuItem): Boolean {
        return when (item.itemId) {
            R.id.action_settings -> {
                // navigate to settings
                true
            }
            else -> super.onOptionsItemSelected(item)
        }
    }
}
```

---

### Adding a Toolbar

`activity_main.xml`:

```xml
<?xml version="1.0" encoding="utf-8"?>
<androidx.coordinatorlayout.widget.CoordinatorLayout
    xmlns:android="http://schemas.android.com/apk/res/android"
    xmlns:app="http://schemas.android.com/apk/res-auto"
    android:layout_width="match_parent"
    android:layout_height="match_parent">

    <com.google.android.material.appbar.AppBarLayout
        android:layout_width="match_parent"
        android:layout_height="wrap_content">

        <androidx.appcompat.widget.Toolbar
            android:id="@+id/toolbar"
            android:layout_width="match_parent"
            android:layout_height="?attr/actionBarSize"
            android:background="?attr/colorPrimary"
            android:theme="@style/ThemeOverlay.AppCompat.Dark.ActionBar" />

    </com.google.android.material.appbar.AppBarLayout>

    <!-- Main content -->
    <FrameLayout
        android:layout_width="match_parent"
        android:layout_height="match_parent"
        app:layout_behavior="@string/appbar_scrolling_view_behavior" />

</androidx.coordinatorlayout.widget.CoordinatorLayout>
```

---

### `AppCompatDelegate` — Night Mode / Dark Theme

```kotlin
import androidx.appcompat.app.AppCompatDelegate

// Force dark mode
AppCompatDelegate.setDefaultNightMode(AppCompatDelegate.MODE_NIGHT_YES)

// Force light mode
AppCompatDelegate.setDefaultNightMode(AppCompatDelegate.MODE_NIGHT_NO)

// Follow system setting (recommended)
AppCompatDelegate.setDefaultNightMode(AppCompatDelegate.MODE_NIGHT_FOLLOW_SYSTEM)
```

Call this in `Application.onCreate()` or when the user toggles a setting.

---

### `AppCompatTextView`, `AppCompatButton`, etc.

When you use `<TextView>` in XML, AppCompat **automatically inflates** it as `AppCompatTextView` behind the scenes. This is why Material styles and vector drawable support works on older devices.

You rarely need to use these classes explicitly — they are handled transparently by `AppCompatActivity`.

---

### Vector Drawable Support on API < 21

```kotlin
// In your Application class or Activity.onCreate
AppCompatDelegate.setCompatVectorFromResourcesEnabled(true)
```

In XML, use `app:srcCompat` instead of `android:src` for `ImageView`:

```xml
<ImageView
    android:layout_width="wrap_content"
    android:layout_height="wrap_content"
    app:srcCompat="@drawable/ic_my_vector" />
```

---

## Combining `core` and `appcompat` — A Practical Example

```kotlin
class ProfileActivity : AppCompatActivity() {

    private lateinit var binding: ActivityProfileBinding

    override fun onCreate(savedInstanceState: Bundle?) {
        super.onCreate(savedInstanceState)
        binding = ActivityProfileBinding.inflate(layoutInflater)
        setContentView(binding.root)

        setSupportActionBar(binding.toolbar)
        supportActionBar?.setDisplayHomeAsUpEnabled(true)

        // core-ktx: safe color
        val accent = ContextCompat.getColor(this, R.color.accent)
        binding.tvName.setTextColor(accent)

        // core-ktx: view extension
        binding.progressBar.isVisible = false

        // Request camera permission
        if (ContextCompat.checkSelfPermission(this, Manifest.permission.CAMERA)
            != PackageManager.PERMISSION_GRANTED) {
            ActivityCompat.requestPermissions(
                this,
                arrayOf(Manifest.permission.CAMERA),
                100
            )
        }
    }

    override fun onSupportNavigateUp(): Boolean {
        onBackPressedDispatcher.onBackPressed()
        return true
    }

    override fun onRequestPermissionsResult(
        requestCode: Int,
        permissions: Array<out String>,
        grantResults: IntArray
    ) {
        super.onRequestPermissionsResult(requestCode, permissions, grantResults)
        if (requestCode == 100 && grantResults.isNotEmpty() &&
            grantResults[0] == PackageManager.PERMISSION_GRANTED) {
            // Camera permission granted
        }
    }
}
```

---

## Common Mistakes

### Mistake 1: Extending `Activity` instead of `AppCompatActivity`

```kotlin
// WRONG — loses Material theming, ActionBar, night mode
class MainActivity : Activity()

// CORRECT
class MainActivity : AppCompatActivity()
```

### Mistake 2: Using `android:src` for vector drawables on older APIs

```xml
<!-- May crash on API < 21 -->
<ImageView android:src="@drawable/ic_vector" />

<!-- Safe across all APIs -->
<ImageView app:srcCompat="@drawable/ic_vector" />
```

### Mistake 3: Not calling `setSupportActionBar` before accessing `supportActionBar`

```kotlin
// WRONG — supportActionBar is null without setSupportActionBar
supportActionBar?.title = "Hello"  // Silent no-op

// CORRECT
setSupportActionBar(binding.toolbar)
supportActionBar?.title = "Hello"
```

---

## Interview Questions

**Q1: What is the difference between `ContextCompat.getColor()` and `context.getColor()`?**

> `context.getColor()` requires API 23+. `ContextCompat.getColor()` works on API 21+ by checking the OS version internally and using the appropriate implementation.

**Q2: What does `core-ktx` add over `core`?**

> `core-ktx` adds idiomatic Kotlin extension functions on top of the Java API, enabling concise syntax like `view.isVisible = true`, `bundleOf()`, and `"string".toUri()`.

**Q3: Why does `AppCompatActivity` matter for theming?**

> `AppCompatActivity` hooks into the view inflation process and replaces standard views (`TextView`, `Button`) with their `AppCompat*` counterparts, which support Material theming and backward-compatible features on older Android versions.

---

## Summary

- `androidx.core` provides backward-compatible OS API wrappers (`ContextCompat`, `ViewCompat`, `ActivityCompat`)
- `core-ktx` adds Kotlin extensions for concise, idiomatic code
- `androidx.appcompat` provides `AppCompatActivity`, Toolbar support, and night mode via `AppCompatDelegate`
- Always extend `AppCompatActivity`, not `Activity`

**Next:** [Chapter 5 — ConstraintLayout](./05-constraint-layout.md)
