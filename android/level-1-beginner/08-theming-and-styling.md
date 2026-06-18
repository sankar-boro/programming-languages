# Chapter 8: Theming and Styling

## Style vs Theme — The Key Distinction

| | Style | Theme |
|--|-------|-------|
| Scope | Applied to a **single view** | Applied to an **Activity or Application** |
| Purpose | Visual appearance of one widget | Global defaults for the entire UI |
| XML attribute | `style="@style/MyStyle"` | `android:theme="@style/Theme.MyApp"` |

**Analogy:** A style is like a clothing item. A theme is like a dress code for the entire venue.

---

## Defining Styles

`res/values/styles.xml`:

```xml
<resources>

    <!-- Base text style -->
    <style name="TextStyle.Title">
        <item name="android:textSize">24sp</item>
        <item name="android:textStyle">bold</item>
        <item name="android:textColor">?attr/colorOnBackground</item>
        <item name="android:letterSpacing">0.01</item>
    </style>

    <!-- Inherited style — dot notation means inheritance -->
    <style name="TextStyle.Title.Large">
        <item name="android:textSize">32sp</item>
    </style>

    <!-- Card style -->
    <style name="CardStyle.Default">
        <item name="cardCornerRadius">12dp</item>
        <item name="cardElevation">2dp</item>
        <item name="android:layout_margin">8dp</item>
    </style>

</resources>
```

Apply to a view:

```xml
<TextView
    android:id="@+id/tvTitle"
    android:layout_width="wrap_content"
    android:layout_height="wrap_content"
    style="@style/TextStyle.Title"
    android:text="Hello" />
```

---

## Defining a Theme

`res/values/themes.xml`:

```xml
<resources>

    <style name="Theme.MyApp" parent="Theme.Material3.DayNight.NoActionBar">

        <!-- Primary brand color -->
        <item name="colorPrimary">@color/purple_500</item>
        <item name="colorPrimaryVariant">@color/purple_700</item>
        <item name="colorOnPrimary">@color/white</item>

        <!-- Secondary brand color -->
        <item name="colorSecondary">@color/teal_200</item>
        <item name="colorOnSecondary">@color/black</item>

        <!-- Status bar color -->
        <item name="android:statusBarColor">?attr/colorPrimaryVariant</item>

        <!-- Default text input style -->
        <item name="textInputStyle">
            @style/Widget.Material3.TextInputLayout.OutlinedBox
        </item>

        <!-- Default button style -->
        <item name="materialButtonStyle">
            @style/Widget.Material3.Button
        </item>

    </style>

</resources>
```

---

## Dark Theme / Night Mode

Define matching night values:

`res/values-night/themes.xml`:

```xml
<resources>
    <style name="Theme.MyApp" parent="Theme.Material3.DayNight.NoActionBar">
        <item name="colorPrimary">@color/purple_200</item>
        <item name="colorPrimaryVariant">@color/purple_700</item>
        <item name="colorOnPrimary">@color/black</item>
        <item name="colorSecondary">@color/teal_200</item>
        <item name="colorOnSecondary">@color/black</item>
        <item name="android:statusBarColor">?attr/colorPrimaryVariant</item>
    </style>
</resources>
```

Toggle at runtime:

```kotlin
// In a settings toggle or splash screen
fun applyTheme(isDark: Boolean) {
    AppCompatDelegate.setDefaultNightMode(
        if (isDark) AppCompatDelegate.MODE_NIGHT_YES
        else AppCompatDelegate.MODE_NIGHT_NO
    )
}
```

Persist the preference and restore on app start:

```kotlin
// Application.onCreate
override fun onCreate() {
    super.onCreate()
    val prefs = getSharedPreferences("settings", MODE_PRIVATE)
    val isDark = prefs.getBoolean("dark_mode", false)
    AppCompatDelegate.setDefaultNightMode(
        if (isDark) AppCompatDelegate.MODE_NIGHT_YES
        else AppCompatDelegate.MODE_NIGHT_FOLLOW_SYSTEM
    )
}
```

---

## Color Resources

`res/values/colors.xml`:

```xml
<resources>
    <!-- Brand colors -->
    <color name="purple_500">#FF6200EE</color>
    <color name="purple_700">#FF3700B3</color>
    <color name="teal_200">#FF03DAC5</color>
    <color name="white">#FFFFFFFF</color>
    <color name="black">#FF000000</color>

    <!-- Semantic colors — use these in layouts -->
    <color name="error">#FFB00020</color>
    <color name="success">#FF4CAF50</color>
    <color name="warning">#FFFF9800</color>
</resources>
```

---

## Text Appearance Styles (Material Type Scale)

Material 3 defines a type scale. Use these instead of ad-hoc font sizes:

```xml
<TextView
    android:layout_width="wrap_content"
    android:layout_height="wrap_content"
    android:text="Display Large"
    style="?attr/textAppearanceDisplayLarge" />

<TextView
    android:layout_width="wrap_content"
    android:layout_height="wrap_content"
    android:text="Title Medium"
    style="?attr/textAppearanceTitleMedium" />

<TextView
    android:layout_width="wrap_content"
    android:layout_height="wrap_content"
    android:text="Body Small"
    style="?attr/textAppearanceBodySmall" />
```

| Style | Typical Use |
|-------|------------|
| `textAppearanceDisplayLarge/Medium/Small` | Hero headers |
| `textAppearanceHeadlineLarge/Medium/Small` | Section titles |
| `textAppearanceTitleLarge/Medium/Small` | Card/list titles |
| `textAppearanceBodyLarge/Medium/Small` | Paragraph text |
| `textAppearanceLabelLarge/Medium/Small` | Button labels, captions |

---

## Custom Fonts

Add a font file to `res/font/`:

```
res/
└── font/
    ├── inter_regular.ttf
    ├── inter_medium.ttf
    └── inter_bold.ttf
```

Define a font family:

`res/font/inter.xml`:

```xml
<?xml version="1.0" encoding="utf-8"?>
<font-family xmlns:android="http://schemas.android.com/apk/res/android"
    xmlns:app="http://schemas.android.com/apk/res-auto">

    <font
        android:font="@font/inter_regular"
        android:fontStyle="normal"
        android:fontWeight="400"
        app:font="@font/inter_regular"
        app:fontStyle="normal"
        app:fontWeight="400" />

    <font
        android:font="@font/inter_medium"
        android:fontStyle="normal"
        android:fontWeight="500"
        app:font="@font/inter_medium"
        app:fontStyle="normal"
        app:fontWeight="500" />

    <font
        android:font="@font/inter_bold"
        android:fontStyle="normal"
        android:fontWeight="700"
        app:font="@font/inter_bold"
        app:fontStyle="normal"
        app:fontWeight="700" />

</font-family>
```

Apply globally in theme:

```xml
<style name="Theme.MyApp" parent="Theme.Material3.DayNight.NoActionBar">
    <item name="android:fontFamily">@font/inter</item>
    <item name="fontFamily">@font/inter</item>
    ...
</style>
```

---

## Theme Overlays

Theme overlays let you apply a theme to just a portion of the view hierarchy, without affecting the rest:

```xml
<!-- Dark toolbar even in light theme -->
<com.google.android.material.appbar.MaterialToolbar
    android:id="@+id/toolbar"
    android:layout_width="match_parent"
    android:layout_height="?attr/actionBarSize"
    android:theme="@style/ThemeOverlay.Material3.Dark.ActionBar" />
```

---

## Dimension Resources

Always use `sp` for text and `dp` for everything else. Define dimensions in one place:

`res/values/dimens.xml`:

```xml
<resources>
    <dimen name="spacing_xs">4dp</dimen>
    <dimen name="spacing_sm">8dp</dimen>
    <dimen name="spacing_md">16dp</dimen>
    <dimen name="spacing_lg">24dp</dimen>
    <dimen name="spacing_xl">32dp</dimen>

    <dimen name="text_size_caption">12sp</dimen>
    <dimen name="text_size_body">16sp</dimen>
    <dimen name="text_size_title">20sp</dimen>

    <dimen name="card_corner_radius">12dp</dimen>
    <dimen name="card_elevation">4dp</dimen>
</resources>
```

---

## Shape System

Define reusable shapes:

`res/values/shapes.xml`:

```xml
<resources>
    <style name="ShapeAppearance.MyApp.SmallComponent"
        parent="ShapeAppearance.Material3.SmallComponent">
        <item name="cornerFamily">rounded</item>
        <item name="cornerSize">8dp</item>
    </style>

    <style name="ShapeAppearance.MyApp.LargeComponent"
        parent="ShapeAppearance.Material3.LargeComponent">
        <item name="cornerFamily">rounded</item>
        <item name="cornerSize">16dp</item>
    </style>
</resources>
```

Apply in theme:

```xml
<item name="shapeAppearanceSmallComponent">
    @style/ShapeAppearance.MyApp.SmallComponent
</item>
```

---

## Complete Theme Setup (Production Template)

```xml
<!-- res/values/themes.xml -->
<resources>
    <style name="Theme.MyApp" parent="Theme.Material3.DayNight.NoActionBar">

        <!-- Colors -->
        <item name="colorPrimary">@color/md_theme_light_primary</item>
        <item name="colorOnPrimary">@color/md_theme_light_onPrimary</item>
        <item name="colorPrimaryContainer">@color/md_theme_light_primaryContainer</item>
        <item name="colorOnPrimaryContainer">@color/md_theme_light_onPrimaryContainer</item>
        <item name="colorSecondary">@color/md_theme_light_secondary</item>
        <item name="colorOnSecondary">@color/md_theme_light_onSecondary</item>
        <item name="colorSurface">@color/md_theme_light_surface</item>
        <item name="colorOnSurface">@color/md_theme_light_onSurface</item>
        <item name="colorBackground">@color/md_theme_light_background</item>
        <item name="colorError">@color/md_theme_light_error</item>

        <!-- Typography -->
        <item name="android:fontFamily">@font/inter</item>
        <item name="fontFamily">@font/inter</item>

        <!-- Shapes -->
        <item name="shapeAppearanceSmallComponent">
            @style/ShapeAppearance.MyApp.SmallComponent
        </item>
        <item name="shapeAppearanceMediumComponent">
            @style/ShapeAppearance.MyApp.MediumComponent
        </item>

        <!-- Window -->
        <item name="android:windowLightStatusBar">true</item>
        <item name="android:navigationBarColor">?attr/colorSurface</item>

        <!-- Component defaults -->
        <item name="textInputStyle">
            @style/Widget.Material3.TextInputLayout.OutlinedBox
        </item>

    </style>
</resources>
```

---

## Common Mistakes

### Mistake 1: Hardcoding colors in layouts

```xml
<!-- WRONG — breaks dark mode, hard to rebrand -->
android:textColor="#212121"
android:background="#FFFFFF"

<!-- CORRECT — semantic, theme-aware -->
android:textColor="?attr/colorOnSurface"
android:background="?attr/colorSurface"
```

### Mistake 2: Using `dp` for text sizes

```xml
<!-- WRONG — ignores user's font size accessibility preference -->
android:textSize="16dp"

<!-- CORRECT -->
android:textSize="16sp"
```

### Mistake 3: Defining colors directly in the theme instead of in `colors.xml`

Keep `themes.xml` referencing only `@color/` names. Keep actual values in `colors.xml` for traceability.

---

## Interview Questions

**Q1: What is the difference between a style and a theme?**

> A style applies appearance attributes to a single view. A theme provides default attributes for a whole Activity or Application — it sets the backdrop that all views inherit from.

**Q2: Why use `?attr/colorPrimary` instead of `@color/purple_500` in layouts?**

> `?attr/` references the current theme's value, which adapts to light/dark mode automatically. `@color/` references a hardcoded color, which won't change with the theme.

**Q3: What is `sp` and why use it for text sizes?**

> `sp` (scale-independent pixels) respects the user's system font size setting for accessibility. `dp` (density-independent pixels) does not — it stays fixed regardless of the user's preferences.

---

## Summary

- Styles apply to single views; themes apply to activities or the entire app
- Always use `?attr/colorXxx` attributes instead of hardcoded colors for theme-awareness
- Use `sp` for text, `dp` for everything else
- Use the Material type scale (`?attr/textAppearanceBodyMedium` etc.) for consistent typography
- Define dark theme in `res/values-night/themes.xml`

**Next:** [Mini Project — Task List App](./mini-project-task-list-app.md)
