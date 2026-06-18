# Chapter 7: Material Components

## What Is Material Design?

Material Design is Google's design system — a set of guidelines and components that define how Android apps should look and feel. Material Components for Android is the AndroidX library that implements it.

```kotlin
implementation("com.google.android.material:material:1.12.0")
```

---

## Setting Up the Material Theme

Your app theme must inherit from a Material theme for components to work correctly.

`res/values/themes.xml`:

```xml
<resources>
    <style name="Theme.MyApp" parent="Theme.Material3.DayNight.NoActionBar">
        <item name="colorPrimary">@color/md_theme_primary</item>
        <item name="colorOnPrimary">@color/md_theme_on_primary</item>
        <item name="colorSecondary">@color/md_theme_secondary</item>
        <item name="colorSurface">@color/md_theme_surface</item>
        <item name="colorError">@color/md_theme_error</item>
    </style>
</resources>
```

`AndroidManifest.xml`:

```xml
<application
    android:theme="@style/Theme.MyApp">
```

---

## Material 3 Core Color Tokens

Material 3 uses semantic color roles rather than hardcoded colors:

| Role | Usage |
|------|-------|
| `colorPrimary` | Main brand color — buttons, highlights |
| `colorSecondary` | Accent color |
| `colorSurface` | Card backgrounds, sheets |
| `colorBackground` | Page background |
| `colorError` | Error states |
| `colorOnPrimary` | Text/icons on top of `colorPrimary` |

Use `?attr/colorPrimary` in XML to reference these dynamically.

---

## MaterialButton

```xml
<!-- Filled button (default) -->
<com.google.android.material.button.MaterialButton
    android:id="@+id/btnSubmit"
    android:layout_width="wrap_content"
    android:layout_height="wrap_content"
    android:text="Submit"
    style="@style/Widget.Material3.Button" />

<!-- Outlined button -->
<com.google.android.material.button.MaterialButton
    android:id="@+id/btnCancel"
    android:layout_width="wrap_content"
    android:layout_height="wrap_content"
    android:text="Cancel"
    style="@style/Widget.Material3.Button.OutlinedButton" />

<!-- Text button -->
<com.google.android.material.button.MaterialButton
    style="@style/Widget.Material3.Button.TextButton"
    android:layout_width="wrap_content"
    android:layout_height="wrap_content"
    android:text="Learn More" />

<!-- Icon button -->
<com.google.android.material.button.MaterialButton
    android:layout_width="wrap_content"
    android:layout_height="wrap_content"
    android:text="Add"
    app:icon="@drawable/ic_add" />
```

---

## TextInputLayout + TextInputEditText

The Material way to implement input fields — provides floating labels, error messages, character counters, and prefix/suffix icons.

```xml
<com.google.android.material.textfield.TextInputLayout
    android:id="@+id/tilEmail"
    android:layout_width="0dp"
    android:layout_height="wrap_content"
    android:hint="Email address"
    app:startIconDrawable="@drawable/ic_email"
    app:endIconMode="clear_text"
    style="@style/Widget.Material3.TextInputLayout.OutlinedBox">

    <com.google.android.material.textfield.TextInputEditText
        android:id="@+id/etEmail"
        android:layout_width="match_parent"
        android:layout_height="wrap_content"
        android:inputType="textEmailAddress"
        android:imeOptions="actionNext" />

</com.google.android.material.textfield.TextInputLayout>
```

Show/hide error in code:

```kotlin
// Show error
binding.tilEmail.error = "Please enter a valid email"
binding.tilEmail.isErrorEnabled = true

// Clear error
binding.tilEmail.error = null
binding.tilEmail.isErrorEnabled = false

// Get text
val email = binding.etEmail.text?.toString()?.trim() ?: ""
```

---

## MaterialCardView

```xml
<com.google.android.material.card.MaterialCardView
    android:id="@+id/cardProfile"
    android:layout_width="match_parent"
    android:layout_height="wrap_content"
    app:cardCornerRadius="12dp"
    app:cardElevation="4dp"
    app:strokeWidth="1dp"
    app:strokeColor="?attr/colorOutline"
    android:layout_margin="8dp">

    <LinearLayout
        android:layout_width="match_parent"
        android:layout_height="wrap_content"
        android:orientation="vertical"
        android:padding="16dp">

        <TextView
            android:layout_width="wrap_content"
            android:layout_height="wrap_content"
            android:text="Card Title"
            android:textSize="18sp"
            android:textStyle="bold" />

        <TextView
            android:layout_width="wrap_content"
            android:layout_height="wrap_content"
            android:text="Card subtitle or description"
            android:textSize="14sp"
            android:layout_marginTop="4dp" />

    </LinearLayout>
</com.google.android.material.card.MaterialCardView>
```

---

## Top App Bar (Toolbar)

```xml
<com.google.android.material.appbar.AppBarLayout
    android:id="@+id/appBarLayout"
    android:layout_width="match_parent"
    android:layout_height="wrap_content">

    <com.google.android.material.appbar.MaterialToolbar
        android:id="@+id/toolbar"
        android:layout_width="match_parent"
        android:layout_height="?attr/actionBarSize"
        app:title="My App"
        app:subtitle="Subtitle"
        app:navigationIcon="@drawable/ic_back"
        app:menu="@menu/main_menu" />

</com.google.android.material.appbar.AppBarLayout>
```

```kotlin
setSupportActionBar(binding.toolbar)
binding.toolbar.setNavigationOnClickListener {
    onBackPressedDispatcher.onBackPressed()
}
```

---

## Bottom Navigation

```xml
<com.google.android.material.bottomnavigation.BottomNavigationView
    android:id="@+id/bottomNav"
    android:layout_width="match_parent"
    android:layout_height="wrap_content"
    app:menu="@menu/bottom_nav_menu"
    app:layout_constraintBottom_toBottomOf="parent"
    app:layout_constraintStart_toStartOf="parent"
    app:layout_constraintEnd_toEndOf="parent" />
```

`res/menu/bottom_nav_menu.xml`:

```xml
<menu xmlns:android="http://schemas.android.com/apk/res/android">
    <item
        android:id="@+id/nav_home"
        android:icon="@drawable/ic_home"
        android:title="Home" />
    <item
        android:id="@+id/nav_search"
        android:icon="@drawable/ic_search"
        android:title="Search" />
    <item
        android:id="@+id/nav_profile"
        android:icon="@drawable/ic_person"
        android:title="Profile" />
</menu>
```

```kotlin
binding.bottomNav.setOnItemSelectedListener { item ->
    when (item.itemId) {
        R.id.nav_home -> { /* show home */ true }
        R.id.nav_search -> { /* show search */ true }
        R.id.nav_profile -> { /* show profile */ true }
        else -> false
    }
}
```

---

## FloatingActionButton (FAB)

```xml
<com.google.android.material.floatingactionbutton.FloatingActionButton
    android:id="@+id/fab"
    android:layout_width="wrap_content"
    android:layout_height="wrap_content"
    android:src="@drawable/ic_add"
    android:contentDescription="Add new item"
    app:layout_constraintBottom_toBottomOf="parent"
    app:layout_constraintEnd_toEndOf="parent"
    android:layout_margin="16dp" />

<!-- Extended FAB with label -->
<com.google.android.material.floatingactionbutton.ExtendedFloatingActionButton
    android:id="@+id/fabExtended"
    android:layout_width="wrap_content"
    android:layout_height="wrap_content"
    android:text="New Note"
    app:icon="@drawable/ic_add"
    app:layout_constraintBottom_toBottomOf="parent"
    app:layout_constraintEnd_toEndOf="parent"
    android:layout_margin="16dp" />
```

---

## Snackbar

```kotlin
// Basic
Snackbar.make(binding.root, "Item deleted", Snackbar.LENGTH_LONG).show()

// With action (undo pattern)
Snackbar.make(binding.root, "Item deleted", Snackbar.LENGTH_LONG)
    .setAction("Undo") {
        // restore the item
    }
    .setAnchorView(binding.fab)  // appears above the FAB
    .show()

// Error state
Snackbar.make(binding.root, "Network error", Snackbar.LENGTH_INDEFINITE)
    .setAction("Retry") { retryRequest() }
    .show()
```

---

## ChipGroup and Chip

```xml
<com.google.android.material.chip.ChipGroup
    android:id="@+id/chipGroup"
    android:layout_width="wrap_content"
    android:layout_height="wrap_content"
    app:singleSelection="true">

    <com.google.android.material.chip.Chip
        android:id="@+id/chipAll"
        android:layout_width="wrap_content"
        android:layout_height="wrap_content"
        android:text="All"
        style="@style/Widget.Material3.Chip.Filter" />

    <com.google.android.material.chip.Chip
        android:id="@+id/chipRecent"
        android:layout_width="wrap_content"
        android:layout_height="wrap_content"
        android:text="Recent"
        style="@style/Widget.Material3.Chip.Filter" />

</com.google.android.material.chip.ChipGroup>
```

```kotlin
binding.chipGroup.setOnCheckedStateChangeListener { group, checkedIds ->
    val selectedId = checkedIds.firstOrNull()
    when (selectedId) {
        R.id.chipAll -> loadAll()
        R.id.chipRecent -> loadRecent()
    }
}
```

---

## Dialog (Material AlertDialog)

```kotlin
import com.google.android.material.dialog.MaterialAlertDialogBuilder

MaterialAlertDialogBuilder(context)
    .setTitle("Delete item?")
    .setMessage("This action cannot be undone.")
    .setPositiveButton("Delete") { dialog, _ ->
        deleteItem()
        dialog.dismiss()
    }
    .setNegativeButton("Cancel") { dialog, _ ->
        dialog.cancel()
    }
    .show()
```

---

## Common Mistakes

### Mistake 1: Using `android.app.AlertDialog` instead of Material

```kotlin
// WRONG — doesn't match Material theme
AlertDialog.Builder(context)

// CORRECT
MaterialAlertDialogBuilder(context)
```

### Mistake 2: Not setting `style` on TextInputLayout

Without the correct style, the TextInputLayout may not show an outline or fill box:

```xml
<!-- CORRECT — explicitly set the style -->
<com.google.android.material.textfield.TextInputLayout
    style="@style/Widget.Material3.TextInputLayout.OutlinedBox"
    ...>
```

### Mistake 3: Not anchoring Snackbar to FAB

The Snackbar will overlap the FAB without `setAnchorView()`.

---

## Interview Questions

**Q1: What is the difference between `TextInputLayout` and a plain `EditText`?**

> `TextInputLayout` wraps an `EditText` to provide floating label animation, error messages, character counter, prefix/suffix, and start/end icons — all following Material Design guidelines. A plain `EditText` provides none of these affordances.

**Q2: What Material component should you use for success/error status messages?**

> `Snackbar` for transient messages (auto-dismiss). For persistent messages or form validation errors, use `TextInputLayout.error` or a `MaterialAlertDialog`.

**Q3: Why inherit from `Theme.Material3.DayNight` instead of `Theme.AppCompat`?**

> `Theme.Material3.DayNight` provides Material 3 color tokens and component styles. Using it ensures consistent theming across all Material components. `Theme.AppCompat` is the older Material 1 baseline and lacks MD3 features like dynamic color.

---

## Summary

- Always inherit from `Theme.Material3.DayNight.NoActionBar` in your theme
- Use `TextInputLayout` for all text input — never bare `EditText` in production
- `MaterialCardView` for cards, `MaterialButton` for buttons, `BottomNavigationView` for nav tabs
- Use `Snackbar` for transient feedback, anchored to FAB when present
- Use `MaterialAlertDialogBuilder` for dialogs

**Next:** [Chapter 8 — Theming and Styling](./08-theming-and-styling.md)
