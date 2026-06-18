# Mini Project: Task List App

## Overview

Build a fully functional task list app using everything covered in Level 1.

**Features:**
- View a list of tasks in a RecyclerView
- Add new tasks via a FAB and dialog
- Mark tasks as complete (checkbox)
- Delete a task with undo (Snackbar)
- Material 3 design with proper theming

---

## Project Structure

```
TaskListApp/
├── app/src/main/
│   ├── java/com/yourname/tasklist/
│   │   ├── MainActivity.kt
│   │   ├── model/
│   │   │   └── Task.kt
│   │   └── adapter/
│   │       └── TaskAdapter.kt
│   └── res/
│       ├── layout/
│       │   ├── activity_main.xml
│       │   ├── item_task.xml
│       │   └── dialog_add_task.xml
│       ├── menu/
│       │   └── main_menu.xml
│       └── values/
│           ├── themes.xml
│           ├── colors.xml
│           └── strings.xml
```

---

## Step 1: `build.gradle.kts` Dependencies

```kotlin
dependencies {
    implementation(libs.androidx.core.ktx)
    implementation(libs.androidx.appcompat)
    implementation(libs.material)
    implementation(libs.androidx.constraintlayout)
}
```

Enable View Binding:

```kotlin
android {
    buildFeatures {
        viewBinding = true
    }
}
```

---

## Step 2: Data Model

`Task.kt`:

```kotlin
data class Task(
    val id: Long = System.currentTimeMillis(),
    val title: String,
    val description: String = "",
    var isCompleted: Boolean = false
)
```

---

## Step 3: Item Layout

`res/layout/item_task.xml`:

```xml
<?xml version="1.0" encoding="utf-8"?>
<com.google.android.material.card.MaterialCardView
    xmlns:android="http://schemas.android.com/apk/res/android"
    xmlns:app="http://schemas.android.com/apk/res-auto"
    android:layout_width="match_parent"
    android:layout_height="wrap_content"
    app:cardCornerRadius="8dp"
    app:cardElevation="2dp"
    android:layout_marginHorizontal="16dp"
    android:layout_marginBottom="8dp">

    <androidx.constraintlayout.widget.ConstraintLayout
        android:layout_width="match_parent"
        android:layout_height="wrap_content"
        android:padding="16dp">

        <CheckBox
            android:id="@+id/cbComplete"
            android:layout_width="wrap_content"
            android:layout_height="wrap_content"
            app:layout_constraintStart_toStartOf="parent"
            app:layout_constraintTop_toTopOf="parent"
            app:layout_constraintBottom_toBottomOf="parent" />

        <TextView
            android:id="@+id/tvTitle"
            android:layout_width="0dp"
            android:layout_height="wrap_content"
            android:textSize="16sp"
            android:textStyle="bold"
            app:layout_constraintStart_toEndOf="@id/cbComplete"
            app:layout_constraintEnd_toEndOf="parent"
            app:layout_constraintTop_toTopOf="parent"
            android:layout_marginStart="12dp" />

        <TextView
            android:id="@+id/tvDescription"
            android:layout_width="0dp"
            android:layout_height="wrap_content"
            android:textSize="14sp"
            android:textColor="?attr/colorOnSurfaceVariant"
            android:visibility="gone"
            app:layout_constraintStart_toStartOf="@id/tvTitle"
            app:layout_constraintEnd_toEndOf="parent"
            app:layout_constraintTop_toBottomOf="@id/tvTitle"
            android:layout_marginTop="4dp" />

    </androidx.constraintlayout.widget.ConstraintLayout>

</com.google.android.material.card.MaterialCardView>
```

---

## Step 4: Add Task Dialog Layout

`res/layout/dialog_add_task.xml`:

```xml
<?xml version="1.0" encoding="utf-8"?>
<LinearLayout
    xmlns:android="http://schemas.android.com/apk/res/android"
    android:layout_width="match_parent"
    android:layout_height="wrap_content"
    android:orientation="vertical"
    android:padding="8dp">

    <com.google.android.material.textfield.TextInputLayout
        android:id="@+id/tilTitle"
        android:layout_width="match_parent"
        android:layout_height="wrap_content"
        android:hint="Task title"
        style="@style/Widget.Material3.TextInputLayout.OutlinedBox">

        <com.google.android.material.textfield.TextInputEditText
            android:id="@+id/etTitle"
            android:layout_width="match_parent"
            android:layout_height="wrap_content"
            android:imeOptions="actionNext"
            android:inputType="textCapSentences" />

    </com.google.android.material.textfield.TextInputLayout>

    <com.google.android.material.textfield.TextInputLayout
        android:id="@+id/tilDescription"
        android:layout_width="match_parent"
        android:layout_height="wrap_content"
        android:hint="Description (optional)"
        android:layout_marginTop="8dp"
        style="@style/Widget.Material3.TextInputLayout.OutlinedBox">

        <com.google.android.material.textfield.TextInputEditText
            android:id="@+id/etDescription"
            android:layout_width="match_parent"
            android:layout_height="wrap_content"
            android:inputType="textMultiLine"
            android:minLines="2" />

    </com.google.android.material.textfield.TextInputLayout>

</LinearLayout>
```

---

## Step 5: Adapter

`TaskAdapter.kt`:

```kotlin
package com.yourname.tasklist.adapter

import android.graphics.Paint
import android.view.LayoutInflater
import android.view.ViewGroup
import androidx.core.view.isVisible
import androidx.recyclerview.widget.DiffUtil
import androidx.recyclerview.widget.ListAdapter
import androidx.recyclerview.widget.RecyclerView
import com.yourname.tasklist.databinding.ItemTaskBinding
import com.yourname.tasklist.model.Task

class TaskAdapter(
    private val onCheckedChange: (Task, Boolean) -> Unit
) : ListAdapter<Task, TaskAdapter.TaskViewHolder>(TaskDiffCallback()) {

    inner class TaskViewHolder(
        private val binding: ItemTaskBinding
    ) : RecyclerView.ViewHolder(binding.root) {

        fun bind(task: Task) {
            binding.tvTitle.text = task.title

            if (task.description.isNotBlank()) {
                binding.tvDescription.text = task.description
                binding.tvDescription.isVisible = true
            } else {
                binding.tvDescription.isVisible = false
            }

            // Set checkbox without triggering listener
            binding.cbComplete.setOnCheckedChangeListener(null)
            binding.cbComplete.isChecked = task.isCompleted

            // Strikethrough completed tasks
            if (task.isCompleted) {
                binding.tvTitle.paintFlags =
                    binding.tvTitle.paintFlags or Paint.STRIKE_THRU_TEXT_FLAG
                binding.tvTitle.alpha = 0.5f
            } else {
                binding.tvTitle.paintFlags =
                    binding.tvTitle.paintFlags and Paint.STRIKE_THRU_TEXT_FLAG.inv()
                binding.tvTitle.alpha = 1.0f
            }

            binding.cbComplete.setOnCheckedChangeListener { _, isChecked ->
                onCheckedChange(task, isChecked)
            }
        }
    }

    override fun onCreateViewHolder(parent: ViewGroup, viewType: Int): TaskViewHolder {
        val binding = ItemTaskBinding.inflate(
            LayoutInflater.from(parent.context), parent, false
        )
        return TaskViewHolder(binding)
    }

    override fun onBindViewHolder(holder: TaskViewHolder, position: Int) {
        holder.bind(getItem(position))
    }
}

class TaskDiffCallback : DiffUtil.ItemCallback<Task>() {
    override fun areItemsTheSame(oldItem: Task, newItem: Task) = oldItem.id == newItem.id
    override fun areContentsTheSame(oldItem: Task, newItem: Task) = oldItem == newItem
}
```

---

## Step 6: Main Layout

`res/layout/activity_main.xml`:

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

        <com.google.android.material.appbar.MaterialToolbar
            android:id="@+id/toolbar"
            android:layout_width="match_parent"
            android:layout_height="?attr/actionBarSize"
            app:title="My Tasks"
            app:menu="@menu/main_menu" />

    </com.google.android.material.appbar.AppBarLayout>

    <androidx.recyclerview.widget.RecyclerView
        android:id="@+id/rvTasks"
        android:layout_width="match_parent"
        android:layout_height="match_parent"
        android:paddingTop="8dp"
        android:paddingBottom="80dp"
        android:clipToPadding="false"
        app:layout_behavior="@string/appbar_scrolling_view_behavior" />

    <TextView
        android:id="@+id/tvEmpty"
        android:layout_width="wrap_content"
        android:layout_height="wrap_content"
        android:layout_gravity="center"
        android:text="No tasks yet.\nTap + to add one."
        android:textAlignment="center"
        android:textSize="16sp"
        android:textColor="?attr/colorOnSurfaceVariant"
        android:visibility="gone" />

    <com.google.android.material.floatingactionbutton.FloatingActionButton
        android:id="@+id/fab"
        android:layout_width="wrap_content"
        android:layout_height="wrap_content"
        android:layout_gravity="bottom|end"
        android:layout_margin="16dp"
        android:src="@drawable/ic_add"
        android:contentDescription="Add task" />

</androidx.coordinatorlayout.widget.CoordinatorLayout>
```

---

## Step 7: MainActivity

`MainActivity.kt`:

```kotlin
package com.yourname.tasklist

import android.os.Bundle
import android.view.LayoutInflater
import androidx.appcompat.app.AppCompatActivity
import androidx.core.view.isVisible
import com.google.android.material.dialog.MaterialAlertDialogBuilder
import com.yourname.tasklist.adapter.TaskAdapter
import com.yourname.tasklist.databinding.ActivityMainBinding
import com.yourname.tasklist.databinding.DialogAddTaskBinding
import com.yourname.tasklist.model.Task
import com.google.android.material.snackbar.Snackbar

class MainActivity : AppCompatActivity() {

    private lateinit var binding: ActivityMainBinding
    private lateinit var adapter: TaskAdapter
    private val tasks = mutableListOf<Task>()

    override fun onCreate(savedInstanceState: Bundle?) {
        super.onCreate(savedInstanceState)
        binding = ActivityMainBinding.inflate(layoutInflater)
        setContentView(binding.root)

        setSupportActionBar(binding.toolbar)
        setupRecyclerView()

        binding.fab.setOnClickListener { showAddTaskDialog() }

        // Load some sample data
        tasks.addAll(
            listOf(
                Task(title = "Buy groceries", description = "Milk, eggs, bread"),
                Task(title = "Read a book"),
                Task(title = "Exercise", isCompleted = true)
            )
        )
        submitTasks()
    }

    private fun setupRecyclerView() {
        adapter = TaskAdapter { task, isChecked ->
            val index = tasks.indexOfFirst { it.id == task.id }
            if (index != -1) {
                tasks[index] = tasks[index].copy(isCompleted = isChecked)
                submitTasks()
            }
        }

        binding.rvTasks.adapter = adapter
        binding.rvTasks.layoutManager =
            androidx.recyclerview.widget.LinearLayoutManager(this)
    }

    private fun showAddTaskDialog() {
        val dialogBinding = DialogAddTaskBinding.inflate(LayoutInflater.from(this))

        MaterialAlertDialogBuilder(this)
            .setTitle("New Task")
            .setView(dialogBinding.root)
            .setPositiveButton("Add") { _, _ ->
                val title = dialogBinding.etTitle.text?.toString()?.trim() ?: ""
                val description = dialogBinding.etDescription.text?.toString()?.trim() ?: ""

                if (title.isBlank()) {
                    dialogBinding.tilTitle.error = "Title is required"
                    return@setPositiveButton
                }

                tasks.add(Task(title = title, description = description))
                submitTasks()
            }
            .setNegativeButton("Cancel", null)
            .show()
    }

    private fun deleteTask(task: Task) {
        val index = tasks.indexOfFirst { it.id == task.id }
        if (index == -1) return

        tasks.removeAt(index)
        submitTasks()

        Snackbar.make(binding.root, "Task deleted", Snackbar.LENGTH_LONG)
            .setAction("Undo") {
                tasks.add(index, task)
                submitTasks()
            }
            .setAnchorView(binding.fab)
            .show()
    }

    private fun submitTasks() {
        adapter.submitList(tasks.toList())
        binding.tvEmpty.isVisible = tasks.isEmpty()
    }
}
```

---

## Challenges to Extend This Project

Once the base is working, extend it:

1. **Swipe to delete** — Add `ItemTouchHelper` to enable swipe-left-to-delete
2. **Sort options** — Add a menu to sort by completion status or alphabetically
3. **Filter chips** — Add All / Active / Completed chips above the list
4. **Long press to edit** — Show the add dialog pre-filled for editing
5. **Animations** — Add item add/remove animations with `DefaultItemAnimator`

---

## Level 1 Checkpoint

Before moving to Level 2, confirm you can:

- [ ] Set up a new AndroidX project from scratch
- [ ] Configure Gradle with the right `compileSdk`, `minSdk`, `targetSdk`
- [ ] Use View Binding to access views
- [ ] Build a UI with ConstraintLayout without using nested layouts
- [ ] Implement a RecyclerView with `ListAdapter` and `DiffUtil`
- [ ] Apply Material components: cards, buttons, TextInputLayout, Snackbar, FAB
- [ ] Support dark/light mode via `values-night` theme

**Next Level:** [Level 2 — Intermediate: Architecture Components](../level-2-intermediate/01-viewmodel.md)
