# Chapter 4: DataStore

## What Is DataStore?

DataStore is AndroidX's replacement for `SharedPreferences`. It stores key-value pairs (Preferences DataStore) or typed objects (Proto DataStore) using Kotlin Flow and coroutines.

**Why replace SharedPreferences?**

| | SharedPreferences | DataStore |
|--|-------------------|-----------|
| Threading | Not safe (can corrupt data) | Coroutine-safe |
| Main thread | Blocks on first read | Never blocks |
| Error handling | Silent failures | Throws exceptions via Flow |
| Type safety | No | Yes (Proto DataStore) |
| Transactions | No | Yes |

```kotlin
implementation("androidx.datastore:datastore-preferences:1.1.1")
// For Proto DataStore:
implementation("androidx.datastore:datastore:1.1.1")
```

---

## Preferences DataStore

### Creating the DataStore

```kotlin
import androidx.datastore.core.DataStore
import androidx.datastore.preferences.core.*
import androidx.datastore.preferences.preferencesDataStore
import android.content.Context

// Extension property — creates one DataStore per Context (usually Application)
val Context.dataStore: DataStore<Preferences> by preferencesDataStore(name = "settings")
```

### Defining Keys

```kotlin
object PreferencesKeys {
    val DARK_MODE = booleanPreferencesKey("dark_mode")
    val USER_NAME = stringPreferencesKey("user_name")
    val FONT_SIZE = intPreferencesKey("font_size")
    val LAST_SYNC = longPreferencesKey("last_sync")
    val NOTIFICATIONS_ENABLED = booleanPreferencesKey("notifications_enabled")
}
```

### Reading Data

```kotlin
class SettingsRepository @Inject constructor(
    private val dataStore: DataStore<Preferences>
) {

    // Flow that emits whenever preferences change
    val isDarkMode: Flow<Boolean> = dataStore.data
        .catch { exception ->
            if (exception is IOException) {
                emit(emptyPreferences())
            } else {
                throw exception
            }
        }
        .map { preferences ->
            preferences[PreferencesKeys.DARK_MODE] ?: false
        }

    val userName: Flow<String> = dataStore.data
        .catch { emit(emptyPreferences()) }
        .map { it[PreferencesKeys.USER_NAME] ?: "" }

    val fontSize: Flow<Int> = dataStore.data
        .catch { emit(emptyPreferences()) }
        .map { it[PreferencesKeys.FONT_SIZE] ?: 16 }

    // Combined preference object
    val userSettings: Flow<UserSettings> = dataStore.data
        .catch { emit(emptyPreferences()) }
        .map { prefs ->
            UserSettings(
                isDarkMode = prefs[PreferencesKeys.DARK_MODE] ?: false,
                userName = prefs[PreferencesKeys.USER_NAME] ?: "",
                fontSize = prefs[PreferencesKeys.FONT_SIZE] ?: 16,
                notificationsEnabled = prefs[PreferencesKeys.NOTIFICATIONS_ENABLED] ?: true
            )
        }
}

data class UserSettings(
    val isDarkMode: Boolean,
    val userName: String,
    val fontSize: Int,
    val notificationsEnabled: Boolean
)
```

### Writing Data

```kotlin
class SettingsRepository @Inject constructor(
    private val dataStore: DataStore<Preferences>
) {

    suspend fun setDarkMode(enabled: Boolean) {
        dataStore.edit { preferences ->
            preferences[PreferencesKeys.DARK_MODE] = enabled
        }
    }

    suspend fun setUserName(name: String) {
        dataStore.edit { preferences ->
            preferences[PreferencesKeys.USER_NAME] = name
        }
    }

    suspend fun updateLastSync() {
        dataStore.edit { preferences ->
            preferences[PreferencesKeys.LAST_SYNC] = System.currentTimeMillis()
        }
    }

    // Atomic update — read-modify-write safely
    suspend fun incrementLoginCount() {
        dataStore.edit { preferences ->
            val current = preferences[intPreferencesKey("login_count")] ?: 0
            preferences[intPreferencesKey("login_count")] = current + 1
        }
    }

    suspend fun clearAllSettings() {
        dataStore.edit { it.clear() }
    }
}
```

---

## Providing DataStore with Hilt

```kotlin
@Module
@InstallIn(SingletonComponent::class)
object DataStoreModule {

    @Provides
    @Singleton
    fun provideDataStore(@ApplicationContext context: Context): DataStore<Preferences> {
        return context.dataStore
    }
}
```

---

## ViewModel with DataStore

```kotlin
@HiltViewModel
class SettingsViewModel @Inject constructor(
    private val settingsRepository: SettingsRepository
) : ViewModel() {

    val userSettings: StateFlow<UserSettings> = settingsRepository.userSettings
        .stateIn(
            scope = viewModelScope,
            started = SharingStarted.WhileSubscribed(5_000),
            initialValue = UserSettings(isDarkMode = false, userName = "", fontSize = 16,
                notificationsEnabled = true)
        )

    fun toggleDarkMode() {
        viewModelScope.launch {
            settingsRepository.setDarkMode(!userSettings.value.isDarkMode)
        }
    }

    fun setUserName(name: String) {
        viewModelScope.launch {
            settingsRepository.setUserName(name)
        }
    }
}
```

---

## Settings Screen Fragment

```kotlin
@AndroidEntryPoint
class SettingsFragment : Fragment(R.layout.fragment_settings) {

    private var _binding: FragmentSettingsBinding? = null
    private val binding get() = _binding!!
    private val viewModel: SettingsViewModel by viewModels()

    override fun onViewCreated(view: View, savedInstanceState: Bundle?) {
        super.onViewCreated(view, savedInstanceState)
        _binding = FragmentSettingsBinding.bind(view)

        viewLifecycleOwner.lifecycleScope.launch {
            viewLifecycleOwner.repeatOnLifecycle(Lifecycle.State.STARTED) {
                viewModel.userSettings.collect { settings ->
                    // Update UI without triggering listeners
                    binding.switchDarkMode.isChecked = settings.isDarkMode
                    binding.switchNotifications.isChecked = settings.notificationsEnabled
                    binding.sliderFontSize.value = settings.fontSize.toFloat()
                }
            }
        }

        binding.switchDarkMode.setOnCheckedChangeListener { _, isChecked ->
            viewModel.toggleDarkMode()
            // Apply theme change immediately
            AppCompatDelegate.setDefaultNightMode(
                if (isChecked) AppCompatDelegate.MODE_NIGHT_YES
                else AppCompatDelegate.MODE_NIGHT_NO
            )
        }
    }

    override fun onDestroyView() {
        super.onDestroyView()
        _binding = null
    }
}
```

---

## Proto DataStore (Strongly Typed)

For complex objects with full type safety, use Proto DataStore:

### 1. Define a `.proto` file

`app/src/main/proto/user_preferences.proto`:

```protobuf
syntax = "proto3";

option java_package = "com.yourname.app";
option java_multiple_files = true;

message UserPreferences {
    bool dark_mode = 1;
    string user_name = 2;
    int32 font_size = 3;
    bool notifications_enabled = 4;
    SortOrder sort_order = 5;

    enum SortOrder {
        UNSPECIFIED = 0;
        BY_DATE = 1;
        BY_TITLE = 2;
    }
}
```

### 2. Create a Serializer

```kotlin
object UserPreferencesSerializer : Serializer<UserPreferences> {
    override val defaultValue: UserPreferences =
        UserPreferences.getDefaultInstance()

    override suspend fun readFrom(input: InputStream): UserPreferences {
        try {
            return UserPreferences.parseFrom(input)
        } catch (exception: InvalidProtocolBufferException) {
            throw CorruptionException("Cannot read proto.", exception)
        }
    }

    override suspend fun writeTo(t: UserPreferences, output: OutputStream) {
        t.writeTo(output)
    }
}
```

### 3. Use It

```kotlin
val Context.userPreferencesStore: DataStore<UserPreferences> by dataStore(
    fileName = "user_prefs.pb",
    serializer = UserPreferencesSerializer
)

// Reading
val darkMode: Flow<Boolean> = context.userPreferencesStore.data
    .map { it.darkMode }

// Writing
suspend fun setDarkMode(enabled: Boolean) {
    context.userPreferencesStore.updateData { current ->
        current.toBuilder().setDarkMode(enabled).build()
    }
}
```

---

## Migrating from SharedPreferences

```kotlin
val dataStore = PreferenceDataStoreFactory.create(
    migrations = listOf(
        SharedPreferencesMigration(context, "old_preferences_name")
    ),
    produceFile = { context.preferencesDataStoreFile("settings") }
)
```

---

## Common Mistakes

### Mistake 1: Creating multiple DataStore instances with the same name

```kotlin
// WRONG — multiple instances for same file = data corruption
class ScreenA {
    val ds = context.dataStore  // ok
}
class ScreenB {
    val ds = context.dataStore  // creates separate instance — WRONG
}

// CORRECT — inject a single instance via Hilt
```

### Mistake 2: Reading DataStore synchronously

```kotlin
// WRONG — blocks the thread
val prefs = runBlocking { dataStore.data.first() }

// CORRECT — collect asynchronously in a lifecycle scope
lifecycleScope.launch {
    dataStore.data.collect { prefs -> updateUI(prefs) }
}
```

### Mistake 3: Not handling `IOException` in the `.catch` operator

DataStore can throw `IOException` if the file is corrupted. Always handle it in `.catch { }` and emit `emptyPreferences()` as a fallback.

---

## Interview Questions

**Q1: What are the advantages of DataStore over SharedPreferences?**

> DataStore is coroutine-safe (no blocking I/O on main thread), handles errors via Flow exceptions rather than silent failures, supports transactions via `edit {}`, and is fully typed with Proto DataStore. SharedPreferences has no thread safety guarantees and can block the main thread on first access.

**Q2: What is the difference between Preferences DataStore and Proto DataStore?**

> Preferences DataStore stores key-value pairs of primitives (Boolean, Int, String). Proto DataStore stores a typed Protocol Buffer object — strongly typed, supports versioning, and handles schema evolution via proto fields.

**Q3: Why use `SharingStarted.WhileSubscribed(5_000)` in `stateIn`?**

> It keeps the upstream (DataStore Flow) active for 5 seconds after the last subscriber disappears — e.g., on screen rotation, the Flow doesn't immediately stop and restart, avoiding a flicker of initial values. After 5 seconds without subscribers, it stops to save resources.

---

## Summary

- DataStore replaces `SharedPreferences` with a coroutine-safe, Flow-based API
- Use Preferences DataStore for simple key-value pairs
- Use Proto DataStore for strongly typed, versioned preferences
- Always handle `IOException` in `.catch { emit(emptyPreferences()) }`
- Provide DataStore as a Singleton via Hilt — never create multiple instances for the same file
- Write with `dataStore.edit { }` — it's atomic and coroutine-safe

**Next:** [Chapter 5 — Unit Testing](./05-unit-testing.md)
