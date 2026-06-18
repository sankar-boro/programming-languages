# Chapter 2: Project Setup in Android Studio

## Installing Android Studio

Download the latest **stable** release from [developer.android.com/studio](https://developer.android.com/studio).

After installation:
1. Open Android Studio
2. Go to **SDK Manager** → install at minimum API 21 (Android 5.0)
3. Go to **AVD Manager** → create a Pixel 6 emulator running API 34

---

## Creating a New AndroidX Project

1. **File → New → New Project**
2. Select template: **Empty Views Activity** (XML) or **Empty Activity** (Compose)
3. Configure:

| Field | Recommended Value |
|-------|------------------|
| Name | `MyFirstApp` |
| Package name | `com.yourname.myfirstapp` |
| Save location | Your projects folder |
| Language | **Kotlin** (preferred) or Java |
| Minimum SDK | **API 24** (covers ~95% of active devices) |

4. Click **Finish** — Android Studio generates a project with AndroidX pre-configured.

---

## Project Structure Overview

```
MyFirstApp/
├── app/
│   ├── src/
│   │   ├── main/
│   │   │   ├── java/com/yourname/myfirstapp/
│   │   │   │   └── MainActivity.kt
│   │   │   ├── res/
│   │   │   │   ├── layout/
│   │   │   │   │   └── activity_main.xml
│   │   │   │   ├── values/
│   │   │   │   │   ├── colors.xml
│   │   │   │   │   ├── strings.xml
│   │   │   │   │   └── themes.xml
│   │   │   │   └── drawable/
│   │   │   └── AndroidManifest.xml
│   │   ├── test/                   ← Unit tests
│   │   └── androidTest/            ← Instrumented (UI) tests
│   └── build.gradle.kts            ← App-level build config
├── build.gradle.kts                ← Project-level build config
├── gradle.properties
├── settings.gradle.kts
└── local.properties
```

---

## Understanding Each Key File

### `AndroidManifest.xml`

Declares the app to the OS: its package name, permissions, activities, and entry point.

```xml
<?xml version="1.0" encoding="utf-8"?>
<manifest xmlns:android="http://schemas.android.com/apk/res/android">

    <!-- Declare permissions here -->
    <uses-permission android:name="android.permission.INTERNET" />

    <application
        android:allowBackup="true"
        android:icon="@mipmap/ic_launcher"
        android:label="@string/app_name"
        android:theme="@style/Theme.MyFirstApp">

        <!-- Every screen (Activity) must be declared here -->
        <activity
            android:name=".MainActivity"
            android:exported="true">
            <intent-filter>
                <!-- This marks MainActivity as the app launcher -->
                <action android:name="android.intent.action.MAIN" />
                <category android:name="android.intent.category.LAUNCHER" />
            </intent-filter>
        </activity>

    </application>
</manifest>
```

### `MainActivity.kt`

The first screen the user sees. Extends `AppCompatActivity` (from AndroidX).

```kotlin
package com.yourname.myfirstapp

import androidx.appcompat.app.AppCompatActivity
import android.os.Bundle

class MainActivity : AppCompatActivity() {

    override fun onCreate(savedInstanceState: Bundle?) {
        super.onCreate(savedInstanceState)
        setContentView(R.layout.activity_main)
    }
}
```

```java
// Java version
package com.yourname.myfirstapp;

import androidx.appcompat.app.AppCompatActivity;
import android.os.Bundle;

public class MainActivity extends AppCompatActivity {

    @Override
    protected void onCreate(Bundle savedInstanceState) {
        super.onCreate(savedInstanceState);
        setContentView(R.layout.activity_main);
    }
}
```

### `activity_main.xml`

The layout file for `MainActivity`. Defines what the screen looks like.

```xml
<?xml version="1.0" encoding="utf-8"?>
<androidx.constraintlayout.widget.ConstraintLayout
    xmlns:android="http://schemas.android.com/apk/res/android"
    xmlns:app="http://schemas.android.com/apk/res-auto"
    android:layout_width="match_parent"
    android:layout_height="match_parent">

    <TextView
        android:id="@+id/tvHello"
        android:layout_width="wrap_content"
        android:layout_height="wrap_content"
        android:text="Hello AndroidX!"
        android:textSize="24sp"
        app:layout_constraintBottom_toBottomOf="parent"
        app:layout_constraintEnd_toEndOf="parent"
        app:layout_constraintStart_toStartOf="parent"
        app:layout_constraintTop_toTopOf="parent" />

</androidx.constraintlayout.widget.ConstraintLayout>
```

---

## Enabling View Binding

View Binding is the modern, safe way to reference XML views in code. It replaces `findViewById`.

Enable it in `app/build.gradle.kts`:

```kotlin
android {
    buildFeatures {
        viewBinding = true
    }
}
```

Then use it in your activity:

```kotlin
class MainActivity : AppCompatActivity() {

    private lateinit var binding: ActivityMainBinding

    override fun onCreate(savedInstanceState: Bundle?) {
        super.onCreate(savedInstanceState)
        binding = ActivityMainBinding.inflate(layoutInflater)
        setContentView(binding.root)

        // Access views by their XML id, no casting needed
        binding.tvHello.text = "Welcome to AndroidX!"
    }
}
```

```java
// Java version
public class MainActivity extends AppCompatActivity {

    private ActivityMainBinding binding;

    @Override
    protected void onCreate(Bundle savedInstanceState) {
        super.onCreate(savedInstanceState);
        binding = ActivityMainBinding.inflate(getLayoutInflater());
        setContentView(binding.getRoot());

        binding.tvHello.setText("Welcome to AndroidX!");
    }
}
```

**Why View Binding over `findViewById`?**

| `findViewById` | View Binding |
|----------------|-------------|
| Returns `View` — requires manual cast | Returns the correct type automatically |
| Can throw `NullPointerException` at runtime | Null safety at compile time |
| No IDE autocomplete for view IDs | Full autocomplete |

---

## Common Mistakes

### Mistake 1: Forgetting to declare an Activity in the Manifest

```
Error: android.content.ActivityNotFoundException
```

Every `Activity` must be registered in `AndroidManifest.xml`.

### Mistake 2: Calling `binding.root` before `inflate`

```kotlin
// WRONG
setContentView(binding.root)  // binding is null here

// CORRECT
binding = ActivityMainBinding.inflate(layoutInflater)
setContentView(binding.root)
```

### Mistake 3: Using `R.id` references that don't exist

If you rename a view ID in XML and forget to update Kotlin/Java code, you get a compile error with View Binding — which is actually better than a silent crash at runtime with `findViewById`.

---

## Exercise

1. Create a new project called `AndroidXPlayground`
2. Enable View Binding
3. Add a `Button` and a `TextView` to `activity_main.xml`
4. When the button is clicked, update the TextView text to "Button clicked!"

```kotlin
// Solution
binding.btnClick.setOnClickListener {
    binding.tvStatus.text = "Button clicked!"
}
```

---

## Interview Questions

**Q1: What is the difference between `Activity` and `AppCompatActivity`?**

> `Activity` is the base Android class. `AppCompatActivity` extends it (from AndroidX) to backport modern features like the ActionBar, Material themes, and day/night mode to older Android versions.

**Q2: What is View Binding and why use it over `findViewById`?**

> View Binding generates a binding class for each XML layout. It provides type-safe, null-safe access to views — eliminating cast errors and NullPointerExceptions that `findViewById` can cause.

**Q3: What does `setContentView` do?**

> It inflates the given layout resource and sets it as the root view of the Activity.

---

## Summary

- Android Studio projects have a clear structure: `src/main`, `res`, `Manifest`, `build.gradle`
- `AppCompatActivity` is the correct base class for all activities in AndroidX projects
- View Binding is the recommended way to access layout views — enable it in `build.gradle.kts`
- Every Activity must be declared in `AndroidManifest.xml`

**Next:** [Chapter 3 — Gradle Configuration](./03-gradle-configuration.md)
