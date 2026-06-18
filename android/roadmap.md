# AndroidX Learning Roadmap

> A structured 120-day plan to go from beginner to production-ready Android developer.

---

## The 30 / 60 / 90 / 120 Day Framework

| Milestone | Goal |
|-----------|------|
| Day 30 | Build and understand a multi-screen Android app with polished UI |
| Day 60 | Build data-driven apps with architecture, persistence, and navigation |
| Day 90 | Write testable, production-grade code with DI and background work |
| Day 120 | Ship a modular, scalable, CI/CD-backed real-world app |

---

## Week-by-Week Plan

### Weeks 1–2: Foundations (Beginner — Part 1)

**Goal:** Set up your environment, understand AndroidX, and build your first real screen.

| Day | Task |
|-----|------|
| 1 | Read [What Is AndroidX](./level-1-beginner/01-what-is-androidx.md). Install Android Studio. |
| 2 | Read [Project Setup](./level-1-beginner/02-project-setup.md). Create your first empty project. |
| 3 | Read [Gradle Configuration](./level-1-beginner/03-gradle-configuration.md). Update deps in your project. |
| 4 | Read [Core & AppCompat](./level-1-beginner/04-core-and-appcompat.md). Run app on device/emulator. |
| 5 | Build: A single screen app using AppCompatActivity + ActionBar |
| 6–7 | Review + experiment. Write the Chapter 1–4 interview answers. |
| 8 | Read [ConstraintLayout](./level-1-beginner/05-constraint-layout.md) |
| 9 | Practice: Recreate a login screen using ConstraintLayout only |
| 10 | Read [RecyclerView](./level-1-beginner/06-recyclerview.md) |
| 11 | Practice: Display a hardcoded list of items using RecyclerView |
| 12 | Read [Material Components](./level-1-beginner/07-material-components.md) |
| 13 | Apply Material theming to your RecyclerView app |
| 14 | Read [Theming & Styling](./level-1-beginner/08-theming-and-styling.md) |

**Checkpoint:** Can you build a screen with a toolbar, a list, and Material components? ✓

---

### Weeks 3–4: Beginner Mini Project

| Day | Task |
|-----|------|
| 15–20 | Build the [Task List Mini Project](./level-1-beginner/mini-project-task-list-app.md) |
| 21 | Code review your own project: Does it follow Material guidelines? |
| 22–23 | Refactor UI using ConstraintLayout where applicable |
| 24–25 | Add RecyclerView with a custom adapter |
| 26–27 | Answer all beginner interview questions |
| 28–30 | Buffer: Revisit any weak areas from Weeks 1–4 |

---

### Weeks 5–6: Architecture Components (Intermediate — Part 1)

| Day | Task |
|-----|------|
| 31 | Read [ViewModel](./level-2-intermediate/01-viewmodel.md) |
| 32 | Add ViewModel to your Task List app — survive rotation |
| 33 | Read [LiveData](./level-2-intermediate/02-livedata.md) |
| 34 | Connect LiveData to your UI |
| 35 | Read [StateFlow](./level-2-intermediate/03-stateflow.md) |
| 36 | Migrate one LiveData usage to StateFlow |
| 37 | Read [Lifecycle](./level-2-intermediate/04-lifecycle.md) |
| 38–39 | Practice: Implement lifecycle-aware logging in your app |
| 40 | Read [Jetpack Compose Intro](./level-2-intermediate/07-jetpack-compose-intro.md) |
| 41–42 | Build a small Compose screen alongside your XML screen |

---

### Weeks 7–8: Navigation + Room (Intermediate — Part 2)

| Day | Task |
|-----|------|
| 43 | Read [Navigation Component](./level-2-intermediate/05-navigation-component.md) |
| 44 | Add multi-screen navigation to your app |
| 45 | Read [Room Database](./level-2-intermediate/06-room-database.md) |
| 46 | Add a Room database, replace hardcoded list |
| 47 | Connect Room → ViewModel → LiveData → UI |
| 48–54 | Build the [Notes App Mini Project](./level-2-intermediate/mini-project-notes-app.md) |
| 55–60 | Buffer, review, answer interview questions |

**Checkpoint:** Can you build a multi-screen CRUD app that persists data across app restarts? ✓

---

### Weeks 9–10: DI + Background Work (Advanced — Part 1)

| Day | Task |
|-----|------|
| 61 | Read [Hilt](./level-3-advanced/01-hilt-dependency-injection.md) |
| 62–63 | Add Hilt to your Notes app |
| 64 | Read [Coroutines & AndroidX](./level-3-advanced/03-coroutines-and-androidx.md) |
| 65 | Migrate Room calls to use coroutines |
| 66 | Read [WorkManager](./level-3-advanced/02-workmanager.md) |
| 67 | Add a background sync task using WorkManager |
| 68 | Read [DataStore](./level-3-advanced/04-datastore.md) |
| 69 | Replace SharedPreferences with DataStore |
| 70 | Refactor, connect all pieces |

---

### Weeks 11–12: Testing (Advanced — Part 2)

| Day | Task |
|-----|------|
| 71 | Read [Unit Testing](./level-3-advanced/05-unit-testing.md) |
| 72–73 | Write unit tests for your ViewModel and Repository |
| 74 | Read [UI Testing with Espresso](./level-3-advanced/06-ui-testing-espresso.md) |
| 75–76 | Write UI tests for your Notes app |
| 77 | Read [Paging 3](./level-3-advanced/07-paging3.md) |
| 78–84 | Build the [News App Mini Project](./level-3-advanced/mini-project-news-app.md) |

**Checkpoint:** Do your apps have ViewModel tests, Repository tests, and at least 2 UI tests? ✓

---

### Weeks 13–14: Expert Architecture (Expert Level)

| Day | Task |
|-----|------|
| 85 | Read [Clean Architecture](./level-4-expert/01-clean-architecture.md) |
| 86 | Read [Modularization](./level-4-expert/02-modularization.md) |
| 87–88 | Refactor News app into feature modules |
| 89 | Read [Offline-First Apps](./level-4-expert/03-offline-first-apps.md) |
| 90 | Read [Performance Optimization](./level-4-expert/04-performance-optimization.md) |
| 91 | Read [Multi-Module Apps](./level-4-expert/05-multi-module-apps.md) |
| 92–98 | Build the [E-Commerce Module Mini Project](./level-4-expert/mini-project-ecommerce-module.md) |

---

### Weeks 15–16: Production + Capstone

| Day | Task |
|-----|------|
| 99 | Read [CI/CD Basics](./level-4-expert/06-cicd-basics.md) |
| 100 | Read [Production Practices](./level-4-expert/07-production-practices.md) |
| 101–120 | Build the full [Capstone Project](./capstone/overview.md) |

**Final Checkpoint:** You have a portfolio app on GitHub with CI/CD, tests, modular architecture, Room, Hilt, Coroutines, and Navigation. ✓

---

## Daily Habit Checklist

- [ ] Read one chapter
- [ ] Type out at least one code example (no copy-paste)
- [ ] Write a short note on what surprised you today
- [ ] Answer interview questions without looking at answers
- [ ] Push your practice code to GitHub

---

## When You Feel Stuck

1. Re-read the chapter — slowly.
2. Search the [AndroidX release notes](https://developer.android.com/jetpack/androidx/versions).
3. Break the problem into the smallest possible reproducible case.
4. Check the [Android Developers blog](https://medium.com/androiddevelopers).
5. Move on and come back — sometimes time is the best debugger.
