# Chapter 5: ConstraintLayout

## What Is ConstraintLayout?

`ConstraintLayout` is a flat, powerful layout system that positions views using **constraints** — rules that define where a view is relative to its parent or other views.

**Analogy:** Think of each view as a rubber ball in a box. Constraints are elastic bands connecting the ball to the walls or other balls. Without constraints, the ball collapses to the top-left corner (0,0). Add constraints and the ball is anchored in position.

```kotlin
implementation("androidx.constraintlayout:constraintlayout:2.2.0")
```

---

## Why ConstraintLayout Over LinearLayout / RelativeLayout?

| Feature | LinearLayout | RelativeLayout | ConstraintLayout |
|---------|-------------|----------------|-----------------|
| Nesting required | High (complex UIs) | Medium | **None — flat hierarchy** |
| Performance | Poor with deep nesting | Medium | **Best — single measure pass** |
| Complex positioning | Difficult | Medium | **Easy** |
| Motion/animation | No | No | **Yes (MotionLayout)** |
| Visual editor support | Basic | Basic | **Excellent** |

Flat hierarchy = fewer view measurements = faster layout rendering.

---

## Constraint Fundamentals

Every view in a `ConstraintLayout` needs **at least one horizontal and one vertical constraint**, or it will snap to (0,0).

### Constraint Sides

Each view has four constrainable sides:
- `layout_constraintTop_to...`
- `layout_constraintBottom_to...`
- `layout_constraintStart_to...`
- `layout_constraintEnd_to...`

### Constraint Targets

Constraints can connect to:
- `parent` — the `ConstraintLayout` itself
- Another view's `@id/viewId`

---

## Basic Example: Centering a View

```xml
<androidx.constraintlayout.widget.ConstraintLayout
    xmlns:android="http://schemas.android.com/apk/res/android"
    xmlns:app="http://schemas.android.com/apk/res-auto"
    android:layout_width="match_parent"
    android:layout_height="match_parent">

    <TextView
        android:id="@+id/tvTitle"
        android:layout_width="wrap_content"
        android:layout_height="wrap_content"
        android:text="Hello World"
        android:textSize="24sp"
        app:layout_constraintTop_toTopOf="parent"
        app:layout_constraintBottom_toBottomOf="parent"
        app:layout_constraintStart_toStartOf="parent"
        app:layout_constraintEnd_toEndOf="parent" />

</androidx.constraintlayout.widget.ConstraintLayout>
```

Constraining all four sides to `parent` centers the view both horizontally and vertically.

---

## Positioning Views Relative to Each Other

```xml
<androidx.constraintlayout.widget.ConstraintLayout
    xmlns:android="http://schemas.android.com/apk/res/android"
    xmlns:app="http://schemas.android.com/apk/res-auto"
    android:layout_width="match_parent"
    android:layout_height="match_parent"
    android:padding="16dp">

    <ImageView
        android:id="@+id/ivAvatar"
        android:layout_width="56dp"
        android:layout_height="56dp"
        android:src="@drawable/ic_person"
        app:layout_constraintTop_toTopOf="parent"
        app:layout_constraintStart_toStartOf="parent" />

    <TextView
        android:id="@+id/tvName"
        android:layout_width="0dp"
        android:layout_height="wrap_content"
        android:text="Alice Johnson"
        android:textSize="18sp"
        android:textStyle="bold"
        app:layout_constraintTop_toTopOf="@id/ivAvatar"
        app:layout_constraintStart_toEndOf="@id/ivAvatar"
        app:layout_constraintEnd_toEndOf="parent"
        android:layout_marginStart="12dp" />

    <TextView
        android:id="@+id/tvSubtitle"
        android:layout_width="0dp"
        android:layout_height="wrap_content"
        android:text="Android Developer"
        android:textSize="14sp"
        app:layout_constraintTop_toBottomOf="@id/tvName"
        app:layout_constraintStart_toStartOf="@id/tvName"
        app:layout_constraintEnd_toEndOf="parent"
        android:layout_marginTop="4dp" />

</androidx.constraintlayout.widget.ConstraintLayout>
```

Key: `android:layout_width="0dp"` means **"fill the space defined by my constraints"** (match constraints).

---

## Bias: Fine-Tuning Position

When a view is constrained on opposite sides, it's centered (50/50 bias). You can shift it:

```xml
<Button
    android:id="@+id/btnLogin"
    android:layout_width="wrap_content"
    android:layout_height="wrap_content"
    android:text="Login"
    app:layout_constraintTop_toTopOf="parent"
    app:layout_constraintBottom_toBottomOf="parent"
    app:layout_constraintStart_toStartOf="parent"
    app:layout_constraintEnd_toEndOf="parent"
    app:layout_constraintVerticalBias="0.3" />
    <!-- 0.0 = top, 0.5 = center, 1.0 = bottom -->
```

---

## Guidelines

Guidelines are invisible reference lines you constrain views to:

```xml
<androidx.constraintlayout.widget.Guideline
    android:id="@+id/guidelineVertical"
    android:layout_width="wrap_content"
    android:layout_height="wrap_content"
    android:orientation="vertical"
    app:layout_constraintGuide_percent="0.5" />
    <!-- Vertical line at 50% of the width -->

<Button
    android:id="@+id/btnLeft"
    android:layout_width="0dp"
    android:layout_height="wrap_content"
    android:text="Left"
    app:layout_constraintStart_toStartOf="parent"
    app:layout_constraintEnd_toStartOf="@id/guidelineVertical"
    app:layout_constraintTop_toTopOf="parent"
    android:layout_margin="8dp" />

<Button
    android:id="@+id/btnRight"
    android:layout_width="0dp"
    android:layout_height="wrap_content"
    android:text="Right"
    app:layout_constraintStart_toEndOf="@id/guidelineVertical"
    app:layout_constraintEnd_toEndOf="parent"
    app:layout_constraintTop_toTopOf="parent"
    android:layout_margin="8dp" />
```

---

## Chains: Distributing Multiple Views

Chains let you distribute multiple views in a line. Create a chain by constraining views to each other in both directions.

```xml
<Button
    android:id="@+id/btn1"
    android:layout_width="wrap_content"
    android:layout_height="wrap_content"
    android:text="One"
    app:layout_constraintStart_toStartOf="parent"
    app:layout_constraintEnd_toStartOf="@id/btn2"
    app:layout_constraintTop_toTopOf="parent"
    app:layout_constraintHorizontalChainStyle="spread" />
    <!-- Chain styles: spread, spread_inside, packed -->

<Button
    android:id="@+id/btn2"
    android:layout_width="wrap_content"
    android:layout_height="wrap_content"
    android:text="Two"
    app:layout_constraintStart_toEndOf="@id/btn1"
    app:layout_constraintEnd_toStartOf="@id/btn3"
    app:layout_constraintTop_toTopOf="@id/btn1" />

<Button
    android:id="@+id/btn3"
    android:layout_width="wrap_content"
    android:layout_height="wrap_content"
    android:text="Three"
    app:layout_constraintStart_toEndOf="@id/btn2"
    app:layout_constraintEnd_toEndOf="parent"
    app:layout_constraintTop_toTopOf="@id/btn1" />
```

Chain styles:
- `spread` — even spacing between views and edges
- `spread_inside` — even spacing between views, no edge padding
- `packed` — views packed together, centered

---

## Barrier: Align to the Widest of Multiple Views

```xml
<!-- Barrier sits at the end of whichever view is wider -->
<androidx.constraintlayout.widget.Barrier
    android:id="@+id/barrierLabels"
    android:layout_width="wrap_content"
    android:layout_height="wrap_content"
    app:barrierDirection="end"
    app:constraint_referenced_ids="tvLabelName,tvLabelEmail,tvLabelPhone" />

<EditText
    android:id="@+id/etName"
    android:layout_width="0dp"
    android:layout_height="wrap_content"
    app:layout_constraintStart_toEndOf="@id/barrierLabels"
    app:layout_constraintEnd_toEndOf="parent"
    android:layout_marginStart="8dp" />
```

---

## Complete Login Screen Example

```xml
<?xml version="1.0" encoding="utf-8"?>
<androidx.constraintlayout.widget.ConstraintLayout
    xmlns:android="http://schemas.android.com/apk/res/android"
    xmlns:app="http://schemas.android.com/apk/res-auto"
    android:layout_width="match_parent"
    android:layout_height="match_parent"
    android:padding="24dp">

    <ImageView
        android:id="@+id/ivLogo"
        android:layout_width="80dp"
        android:layout_height="80dp"
        android:src="@mipmap/ic_launcher"
        app:layout_constraintTop_toTopOf="parent"
        app:layout_constraintStart_toStartOf="parent"
        app:layout_constraintEnd_toEndOf="parent"
        app:layout_constraintVerticalBias="0.15"
        app:layout_constraintBottom_toBottomOf="parent" />

    <com.google.android.material.textfield.TextInputLayout
        android:id="@+id/tilEmail"
        android:layout_width="0dp"
        android:layout_height="wrap_content"
        android:hint="Email"
        app:layout_constraintTop_toBottomOf="@id/ivLogo"
        app:layout_constraintStart_toStartOf="parent"
        app:layout_constraintEnd_toEndOf="parent"
        android:layout_marginTop="32dp">

        <com.google.android.material.textfield.TextInputEditText
            android:id="@+id/etEmail"
            android:layout_width="match_parent"
            android:layout_height="wrap_content"
            android:inputType="textEmailAddress" />
    </com.google.android.material.textfield.TextInputLayout>

    <com.google.android.material.textfield.TextInputLayout
        android:id="@+id/tilPassword"
        android:layout_width="0dp"
        android:layout_height="wrap_content"
        android:hint="Password"
        app:passwordToggleEnabled="true"
        app:layout_constraintTop_toBottomOf="@id/tilEmail"
        app:layout_constraintStart_toStartOf="parent"
        app:layout_constraintEnd_toEndOf="parent"
        android:layout_marginTop="16dp">

        <com.google.android.material.textfield.TextInputEditText
            android:id="@+id/etPassword"
            android:layout_width="match_parent"
            android:layout_height="wrap_content"
            android:inputType="textPassword" />
    </com.google.android.material.textfield.TextInputLayout>

    <com.google.android.material.button.MaterialButton
        android:id="@+id/btnLogin"
        android:layout_width="0dp"
        android:layout_height="wrap_content"
        android:text="Log In"
        app:layout_constraintTop_toBottomOf="@id/tilPassword"
        app:layout_constraintStart_toStartOf="parent"
        app:layout_constraintEnd_toEndOf="parent"
        android:layout_marginTop="24dp" />

    <TextView
        android:id="@+id/tvForgotPassword"
        android:layout_width="wrap_content"
        android:layout_height="wrap_content"
        android:text="Forgot password?"
        android:textColor="?attr/colorPrimary"
        app:layout_constraintTop_toBottomOf="@id/btnLogin"
        app:layout_constraintEnd_toEndOf="parent"
        android:layout_marginTop="12dp" />

</androidx.constraintlayout.widget.ConstraintLayout>
```

---

## Common Mistakes

### Mistake 1: Missing constraints (view collapses to 0,0)

```xml
<!-- WRONG — no horizontal or vertical constraint -->
<TextView
    android:layout_width="wrap_content"
    android:layout_height="wrap_content"
    android:text="Hello" />

<!-- CORRECT — minimum viable constraints -->
<TextView
    android:layout_width="wrap_content"
    android:layout_height="wrap_content"
    android:text="Hello"
    app:layout_constraintTop_toTopOf="parent"
    app:layout_constraintStart_toStartOf="parent" />
```

### Mistake 2: Nesting ConstraintLayouts

ConstraintLayout is designed for flat hierarchies. Nesting them defeats the purpose.

### Mistake 3: Using `match_parent` for width/height

```xml
<!-- WRONG — use 0dp for "match constraints" -->
android:layout_width="match_parent"

<!-- CORRECT inside ConstraintLayout -->
android:layout_width="0dp"
```

---

## Interview Questions

**Q1: What does `android:layout_width="0dp"` mean inside ConstraintLayout?**

> It means "match constraints" — the view expands to fill the space defined by its start and end constraints. It is equivalent to `match_parent` in simpler layouts but respects constraint boundaries.

**Q2: What is a chain in ConstraintLayout?**

> A chain is a group of views mutually constrained to each other on the same axis. The chain head (leftmost/topmost) controls the chain style: `spread`, `spread_inside`, or `packed`.

**Q3: Why is ConstraintLayout better than nested LinearLayouts for complex UIs?**

> ConstraintLayout creates a flat view hierarchy — all views are direct children of one parent. Nested LinearLayouts cause multiple measurement passes per level, increasing render time. A flat hierarchy means a single measurement pass regardless of complexity.

---

## Exercises

1. Recreate a profile card: avatar on the left, name and subtitle stacked on the right
2. Build a bottom button bar with three equally-spaced buttons using a chain
3. Build a two-column grid form using guidelines at 50%

---

## Summary

- ConstraintLayout positions views with rules (constraints) rather than nesting
- Every view needs at least one horizontal and one vertical constraint
- `0dp` width/height means "fill constrained space"
- Chains distribute multiple views; Guidelines create reference lines; Barriers adapt to the widest view
- Avoid nesting ConstraintLayouts — keep the hierarchy flat

**Next:** [Chapter 6 — RecyclerView](./06-recyclerview.md)
