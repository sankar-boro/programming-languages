# Capstone Project: FeedFlow — A Full-Featured News Reader App

## What You're Building

**FeedFlow** is a production-quality news reader app that demonstrates mastery of all AndroidX topics covered in this guide.

---

## App Features

| Feature | Technology Used |
|---------|----------------|
| Browse top headlines | Retrofit + Paging 3 |
| Search articles | Debounced StateFlow + Paging |
| Bookmark articles | Room Database |
| Offline reading | RemoteMediator offline-first |
| Dark/light theme | DataStore + AppCompatDelegate |
| Multi-screen navigation | Navigation Component + Safe Args |
| Background sync | WorkManager |
| Dependency injection | Hilt |
| Unit tested | JUnit + Mockito + Coroutines Test |
| UI tested | Espresso |

---

## Architecture

```
FeedFlow
├── Clean Architecture (domain/data/ui)
├── MVVM with StateFlow
├── Modular structure (:app, :core:*, :feature:*)
└── Offline-first with Room + RemoteMediator
```

---

## Module Map

```
feedflow/
├── app/                       ← Application, wiring
├── core/
│   ├── common/                ← Result, extensions, dispatchers
│   ├── domain/                ← Article, Bookmark models + interfaces
│   ├── data/                  ← Repository implementations
│   ├── database/              ← Room setup
│   ├── network/               ← Retrofit setup
│   ├── ui/                    ← Shared composables, themes
│   └── testing/               ← Fakes, rules
└── feature/
    ├── headlines/             ← Top headlines list
    ├── search/                ← Search feature
    ├── bookmarks/             ← Saved articles
    └── settings/              ← Theme, preferences
```

---

## Capstone Parts

| Part | File | Focus |
|------|------|-------|
| 1 | [Architecture Setup](./01-architecture-setup.md) | Gradle, modules, Hilt, base classes |
| 2 | [UI Implementation](./02-ui-implementation.md) | XML layouts, Compose, RecyclerView, theming |
| 3 | [API Integration](./03-api-integration.md) | Retrofit, DTOs, error handling |
| 4 | [Database Layer](./04-database-layer.md) | Room, DAOs, RemoteMediator |
| 5 | [Navigation Flow](./05-navigation-flow.md) | Navigation Component, deep links |
| 6 | [Testing](./06-testing.md) | Unit + UI tests |
| 7 | [Production Polish](./07-production-polish.md) | CI/CD, ProGuard, Crashlytics |

---

## Prerequisites

Before starting the Capstone:
- Get a free API key from [newsapi.org](https://newsapi.org)
- Create a new Android Studio project named `FeedFlow`
- Initialize a Git repository

---

## Final Deliverables

By the end, you will have:

1. A working Android app on the Play Store internal testing track
2. A GitHub repository with:
   - Clean Architecture
   - Passing CI pipeline
   - Unit tests with >80% ViewModel coverage
   - UI tests for main flows
3. A `README.md` with screenshots and setup instructions

**Start here:** [Part 1 — Architecture Setup](./01-architecture-setup.md)
