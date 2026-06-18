# Chapter 6: UI Testing with Espresso

## What Is Espresso?

Espresso is AndroidX's UI testing framework. It runs on a real device or emulator, directly interacting with your app's UI — finding views, performing actions, and verifying assertions.

```kotlin
androidTestImplementation("androidx.test.espresso:espresso-core:3.6.1")
androidTestImplementation("androidx.test.espresso:espresso-contrib:3.6.1")
androidTestImplementation("androidx.test:runner:1.6.2")
androidTestImplementation("androidx.test:rules:1.6.1")
androidTestImplementation("androidx.test.ext:junit:1.2.1")
androidTestImplementation("androidx.fragment:fragment-testing:1.8.5")
```

---

## Espresso Basics

Every Espresso test follows three steps:

```kotlin
// 1. Find the view (ViewMatcher)
onView(withId(R.id.btnLogin))

// 2. Perform an action
    .perform(click())

// 3. Check an assertion
onView(withId(R.id.tvWelcome))
    .check(matches(withText("Welcome, Alice!")))
```

---

## Common ViewMatchers

```kotlin
// By ID
onView(withId(R.id.btnSubmit))

// By text
onView(withText("Submit"))
onView(withText(R.string.submit))

// By content description
onView(withContentDescription("Profile picture"))

// By hint
onView(withHint("Enter email"))

// By ancestor/descendant
onView(allOf(withText("Delete"), isDescendantOfA(withId(R.id.cardItem))))

// In RecyclerView at position
onView(withRecyclerView(R.id.rvNotes).atPosition(0))
```

---

## Common ViewActions

```kotlin
// Tap
perform(click())

// Long press
perform(longClick())

// Type text
perform(typeText("Hello World"))

// Clear and retype
perform(clearText(), typeText("New text"))

// Scroll to view (for scrollable containers)
perform(scrollTo())

// Swipe
perform(swipeLeft())
perform(swipeRight())
perform(swipeUp())
perform(swipeDown())

// Close keyboard
perform(closeSoftKeyboard())

// Replace text (faster than typeText)
perform(replaceText("instant text"))

// Press back
pressBack()
```

---

## Common ViewAssertions

```kotlin
// Visibility
check(matches(isDisplayed()))
check(matches(not(isDisplayed())))
check(matches(withEffectiveVisibility(Visibility.GONE)))

// Text content
check(matches(withText("Expected text")))
check(matches(not(withText("Wrong text"))))

// Enable state
check(matches(isEnabled()))
check(matches(not(isEnabled())))

// Checked state
check(matches(isChecked()))
check(matches(isNotChecked()))

// Error text (on TextInputLayout)
check(matches(hasErrorText("Field is required")))

// Count in RecyclerView
onView(withId(R.id.rvNotes))
    .check(RecyclerViewItemCountAssertion(3))
```

---

## A Complete Espresso Test — Notes App

```kotlin
@RunWith(AndroidJUnit4::class)
class NoteListFragmentTest {

    @get:Rule
    val activityRule = ActivityScenarioRule(MainActivity::class.java)

    @Test
    fun notesScreen_displaysEmptyState_whenNoNotes() {
        onView(withId(R.id.tvEmpty))
            .check(matches(isDisplayed()))

        onView(withId(R.id.tvEmpty))
            .check(matches(withText("No notes yet.\nTap + to add one.")))
    }

    @Test
    fun clickFab_opensAddNoteDialog() {
        onView(withId(R.id.fab))
            .perform(click())

        onView(withText("New Note"))
            .check(matches(isDisplayed()))
    }

    @Test
    fun addNote_appearsInList() {
        // Open add dialog
        onView(withId(R.id.fab)).perform(click())

        // Type the note title
        onView(withId(R.id.etTitle))
            .perform(typeText("Buy groceries"), closeSoftKeyboard())

        // Type content
        onView(withId(R.id.etContent))
            .perform(typeText("Milk, eggs, bread"), closeSoftKeyboard())

        // Tap Add
        onView(withText("Add")).perform(click())

        // Verify the note appears in the list
        onView(withId(R.id.rvNotes))
            .check(matches(hasDescendant(withText("Buy groceries"))))
    }

    @Test
    fun addNote_withEmptyTitle_showsError() {
        onView(withId(R.id.fab)).perform(click())

        // Don't type title — just press Add
        onView(withText("Add")).perform(click())

        onView(withId(R.id.tilTitle))
            .check(matches(hasErrorText("Title is required")))
    }
}
```

---

## Testing a Fragment in Isolation

Use `FragmentScenario` to launch a Fragment without a full Activity:

```kotlin
@RunWith(AndroidJUnit4::class)
class NoteEditorFragmentTest {

    @Test
    fun saveButton_disabledWhenTitleEmpty() {
        val args = bundleOf("noteId" to -1L)

        launchFragmentInContainer<NoteEditorFragment>(
            fragmentArgs = args,
            themeResId = R.style.Theme_MyApp
        ) {
            // Fragment is created with Hilt — use HiltFragmentScenario for Hilt
        }

        onView(withId(R.id.btnSave))
            .check(matches(not(isEnabled())))
    }

    @Test
    fun typingTitle_enablesSaveButton() {
        launchFragmentInContainer<NoteEditorFragment>(
            fragmentArgs = bundleOf("noteId" to -1L),
            themeResId = R.style.Theme_MyApp
        )

        onView(withId(R.id.etTitle))
            .perform(typeText("My Note"), closeSoftKeyboard())

        onView(withId(R.id.btnSave))
            .check(matches(isEnabled()))
    }
}
```

---

## Testing RecyclerView

```kotlin
@Test
fun recyclerView_displaysCorrectItems() {
    // Inject test data
    // (In a real app, use Hilt test modules to inject a FakeRepository)

    // Verify item at position 0
    onView(withId(R.id.rvNotes))
        .perform(
            RecyclerViewActions.scrollToPosition<RecyclerView.ViewHolder>(0)
        )

    onView(
        RecyclerViewMatcher(R.id.rvNotes).atPositionOnView(0, R.id.tvTitle)
    ).check(matches(withText("First Note")))
}

@Test
fun swipeToDelete_removesItem() {
    onView(withId(R.id.rvNotes))
        .perform(
            RecyclerViewActions.actionOnItemAtPosition<RecyclerView.ViewHolder>(
                0, GeneralSwipeAction(
                    Swipe.SLOW,
                    GeneralLocation.BOTTOM_RIGHT,
                    GeneralLocation.BOTTOM_LEFT,
                    Press.FINGER
                )
            )
        )

    // Verify undo Snackbar appears
    onView(withText("Note deleted"))
        .check(matches(isDisplayed()))
}
```

---

## Hilt Integration Tests

```kotlin
@HiltAndroidTest
@RunWith(AndroidJUnit4::class)
class NoteListIntegrationTest {

    @get:Rule(order = 0)
    val hiltRule = HiltAndroidRule(this)

    @get:Rule(order = 1)
    val activityRule = ActivityScenarioRule(MainActivity::class.java)

    @Inject
    lateinit var repository: NoteRepository

    @Before
    fun setup() {
        hiltRule.inject()
    }

    @Test
    fun insertedNote_appearsInList() = runBlocking {
        repository.saveNote(Note(title = "Integration Test Note"))

        onView(withId(R.id.rvNotes))
            .check(matches(hasDescendant(withText("Integration Test Note"))))
    }
}
```

Replace database module for tests (in-memory):

```kotlin
@Module
@TestInstallIn(
    components = [SingletonComponent::class],
    replaces = [DatabaseModule::class]
)
object TestDatabaseModule {
    @Provides
    @Singleton
    fun provideInMemoryDatabase(@ApplicationContext context: Context): NoteDatabase =
        Room.inMemoryDatabaseBuilder(context, NoteDatabase::class.java)
            .allowMainThreadQueries()
            .build()
}
```

---

## Idling Resources

When Espresso needs to wait for async work (coroutines, network), use Idling Resources:

```kotlin
// For coroutines
implementation("androidx.test.espresso:espresso-idling-resource:3.6.1")

class CountingIdlingResourceWrapper(name: String) {
    val idlingResource = CountingIdlingResource(name)

    fun increment() = idlingResource.increment()
    fun decrement() = idlingResource.decrement()
}

// Register in test
@Before
fun setup() {
    IdlingRegistry.getInstance().register(myIdlingResource.idlingResource)
}

@After
fun tearDown() {
    IdlingRegistry.getInstance().unregister(myIdlingResource.idlingResource)
}
```

---

## Common Mistakes

### Mistake 1: Not closing the keyboard before asserting

```kotlin
// FAILS — keyboard covers the view
onView(withId(R.id.etTitle)).perform(typeText("Test"))
onView(withId(R.id.btnSave)).perform(click())  // might not find btnSave

// CORRECT
onView(withId(R.id.etTitle)).perform(typeText("Test"), closeSoftKeyboard())
onView(withId(R.id.btnSave)).perform(click())
```

### Mistake 2: Using `Thread.sleep()` instead of Idling Resources

```kotlin
// WRONG — flaky, breaks on slow devices
Thread.sleep(2000)
onView(withId(R.id.tvResult)).check(matches(withText("Done")))

// CORRECT — use Idling Resource or test with fake data that resolves synchronously
```

### Mistake 3: Testing implementation details (verify exact view IDs the user doesn't see)

Espresso tests should verify user-visible behavior — text, visibility, clickability — not internal IDs.

---

## Interview Questions

**Q1: What is Espresso and where do Espresso tests run?**

> Espresso is AndroidX's UI testing framework. Espresso tests are instrumented tests — they run on a real device or emulator, inside the app's process, allowing direct interaction with the UI.

**Q2: How do you handle asynchronous operations in Espresso tests?**

> Use `IdlingResource` — register a resource that signals to Espresso when the app is idle. Espresso waits for all registered idling resources to be idle before proceeding with the next test action.

**Q3: What is the difference between unit tests and instrumented tests?**

> Unit tests run on the JVM (fast, no device needed) and test logic in isolation. Instrumented tests run on a device/emulator (slow, requires device) and test real Android behavior including UI. Place unit tests in `src/test/`, instrumented tests in `src/androidTest/`.

---

## Summary

- Espresso tests run on device: `onView().perform().check()`
- Use `ActivityScenarioRule` to launch activities in tests
- Use `FragmentScenario` to test Fragments in isolation
- Use `@HiltAndroidTest` with `TestInstallIn` to replace real modules with test doubles
- Use `RecyclerViewActions` for RecyclerView interaction
- Never use `Thread.sleep()` — use Idling Resources

**Next:** [Chapter 7 — Paging 3](./07-paging3.md)
