# Chapter 3: Offline-First Apps

## What Is Offline-First?

An offline-first app works fully without a network connection. When connectivity returns, it syncs changes with the server. The local database is always the source of truth for the UI.

**Why it matters:**
- Users in areas with poor connectivity
- Faster UI — no waiting for network
- Better user experience — no empty states or spinners
- Reliability — the app always works

---

## Offline-First Architecture

```
UI Layer ──► ViewModel ──► Repository
                                │
                    ┌───────────┴───────────┐
                    ▼                       ▼
              Local DB (Room)         Remote API
              (source of truth)       (sync source)
```

The Repository pattern is the key:
- **UI always reads from Room**
- **Repository decides when to fetch from network**
- **Network results write to Room**, then Room emits the update to the UI

---

## Single Source of Truth with Room + Retrofit

```kotlin
class NoteRepository @Inject constructor(
    private val dao: NoteDao,
    private val api: NoteApiService,
    private val connectivityObserver: ConnectivityObserver
) {
    // UI observes Room — always up to date from any source
    val notes: Flow<List<Note>> = dao.getAllNotes()
        .map { entities -> entities.map { it.toDomain() } }

    // Background sync — called by WorkManager or on app launch
    suspend fun sync() {
        if (!connectivityObserver.isOnline()) return
        try {
            val remoteNotes = api.fetchNotes()
            dao.upsertAll(remoteNotes.map { it.toEntity() })
        } catch (e: IOException) {
            // Sync failed — local data still shown
        }
    }

    // Write locally first, then sync to server
    suspend fun saveNote(note: Note) {
        val entity = note.toEntity().copy(syncStatus = SyncStatus.PENDING)
        dao.insertNote(entity)
        // Try to sync immediately if online
        if (connectivityObserver.isOnline()) {
            trySyncNote(entity)
        }
        // Otherwise WorkManager will sync later
    }
}
```

---

## Sync Status Tracking

Track which records need to be synced:

```kotlin
enum class SyncStatus {
    SYNCED,     // In sync with server
    PENDING,    // Created/modified locally, not yet sent to server
    DELETED     // Deleted locally, pending server deletion
}

@Entity(tableName = "notes")
data class NoteEntity(
    @PrimaryKey val id: String = UUID.randomUUID().toString(),
    val title: String,
    val content: String,
    val updatedAt: Long = System.currentTimeMillis(),
    val syncStatus: SyncStatus = SyncStatus.PENDING,
    val serverId: String? = null  // null until server assigns an ID
)
```

---

## ConnectivityObserver

```kotlin
import android.net.ConnectivityManager
import android.net.Network
import android.net.NetworkCapabilities
import kotlinx.coroutines.channels.awaitClose
import kotlinx.coroutines.flow.Flow
import kotlinx.coroutines.flow.callbackFlow
import kotlinx.coroutines.flow.distinctUntilChanged

interface ConnectivityObserver {
    val isConnected: Flow<Boolean>
    fun isOnline(): Boolean
}

class NetworkConnectivityObserver @Inject constructor(
    @ApplicationContext private val context: Context
) : ConnectivityObserver {

    private val connectivityManager =
        context.getSystemService(Context.CONNECTIVITY_SERVICE) as ConnectivityManager

    override val isConnected: Flow<Boolean> = callbackFlow {
        val callback = object : ConnectivityManager.NetworkCallback() {
            override fun onAvailable(network: Network) { trySend(true) }
            override fun onLost(network: Network) { trySend(false) }
            override fun onUnavailable() { trySend(false) }
        }

        val request = NetworkRequest.Builder()
            .addCapability(NetworkCapabilities.NET_CAPABILITY_INTERNET)
            .build()

        connectivityManager.registerNetworkCallback(request, callback)

        // Emit current state immediately
        trySend(isOnline())

        awaitClose { connectivityManager.unregisterNetworkCallback(callback) }
    }.distinctUntilChanged()

    override fun isOnline(): Boolean {
        val network = connectivityManager.activeNetwork ?: return false
        val capabilities = connectivityManager.getNetworkCapabilities(network) ?: return false
        return capabilities.hasCapability(NetworkCapabilities.NET_CAPABILITY_INTERNET) &&
               capabilities.hasCapability(NetworkCapabilities.NET_CAPABILITY_VALIDATED)
    }
}
```

---

## Sync Strategy — WorkManager

```kotlin
@HiltWorker
class SyncWorker @AssistedInject constructor(
    @Assisted context: Context,
    @Assisted params: WorkerParameters,
    private val repository: NoteRepository
) : CoroutineWorker(context, params) {

    override suspend fun doWork(): Result {
        return try {
            repository.syncPendingNotes()
            repository.syncFromServer()
            Result.success()
        } catch (e: IOException) {
            if (runAttemptCount < 3) Result.retry() else Result.failure()
        }
    }

    companion object {
        fun schedule(context: Context) {
            val request = PeriodicWorkRequestBuilder<SyncWorker>(15, TimeUnit.MINUTES)
                .setConstraints(
                    Constraints.Builder()
                        .setRequiredNetworkType(NetworkType.CONNECTED)
                        .build()
                )
                .build()

            WorkManager.getInstance(context).enqueueUniquePeriodicWork(
                "notes_sync",
                ExistingPeriodicWorkPolicy.KEEP,
                request
            )
        }

        fun syncNow(context: Context) {
            val request = OneTimeWorkRequestBuilder<SyncWorker>()
                .setConstraints(
                    Constraints.Builder()
                        .setRequiredNetworkType(NetworkType.CONNECTED)
                        .build()
                )
                .build()
            WorkManager.getInstance(context).enqueue(request)
        }
    }
}
```

---

## Optimistic Updates

Show the change immediately, then sync in the background:

```kotlin
class NoteRepository @Inject constructor(
    private val dao: NoteDao,
    private val api: NoteApiService
) {
    suspend fun deleteNote(note: Note) {
        // Mark as deleted locally — UI updates immediately
        dao.markAsDeleted(note.id)

        // Try server delete in background
        try {
            api.deleteNote(note.serverId!!)
            dao.deleteNote(note.id)  // Fully remove from local DB
        } catch (e: IOException) {
            // Sync will retry later via WorkManager
            // For now, note is marked DELETED and hidden from UI
        }
    }
}
```

---

## Conflict Resolution

When local and server data conflict:

```kotlin
suspend fun resolveConflicts(localNote: NoteEntity, remoteNote: NoteDto) {
    val resolved = when {
        // Server is newer — use server version
        remoteNote.updatedAt > localNote.updatedAt -> remoteNote.toEntity()

        // Local is newer and already pending sync — keep local
        localNote.syncStatus == SyncStatus.PENDING &&
            localNote.updatedAt > remoteNote.updatedAt -> localNote

        // Same timestamp or already synced — server wins
        else -> remoteNote.toEntity()
    }
    dao.upsertNote(resolved)
}
```

---

## UI Indicators for Sync Status

```kotlin
// Show sync status in the UI
@Composable
fun NoteItem(note: Note) {
    Row(verticalAlignment = Alignment.CenterVertically) {
        Text(note.title, modifier = Modifier.weight(1f))

        when (note.syncStatus) {
            SyncStatus.PENDING -> Icon(
                imageVector = Icons.Default.CloudOff,
                contentDescription = "Pending sync",
                tint = MaterialTheme.colorScheme.onSurfaceVariant
            )
            SyncStatus.SYNCED -> { /* No indicator needed */ }
            SyncStatus.DELETED -> { /* Handled elsewhere */ }
        }
    }
}
```

---

## Caching Strategy with `staleWhileRevalidate`

```kotlin
class NoteRepository @Inject constructor(
    private val dao: NoteDao,
    private val api: NoteApiService
) {
    fun getNotesWithRefresh(): Flow<Resource<List<Note>>> = flow {
        // Emit cached data immediately
        val cached = dao.getAllNotesList()
        if (cached.isNotEmpty()) {
            emit(Resource.Success(cached.map { it.toDomain() }, isFromCache = true))
        }

        // Try to refresh from network
        try {
            val remote = api.fetchNotes()
            dao.upsertAll(remote.map { it.toEntity() })
            emit(Resource.Success(remote.map { it.toDomain().also { n -> dao.insertNote(n.toEntity()) } },
                isFromCache = false))
        } catch (e: IOException) {
            if (cached.isEmpty()) {
                emit(Resource.Error("No connection and no cached data"))
            }
            // If we have cache, silently swallow the network error
        }
    }
}

sealed class Resource<T> {
    data class Success<T>(val data: T, val isFromCache: Boolean = false) : Resource<T>()
    data class Error<T>(val message: String, val cachedData: T? = null) : Resource<T>()
    class Loading<T> : Resource<T>()
}
```

---

## Interview Questions

**Q1: What does "single source of truth" mean in an offline-first app?**

> The local database (Room) is the single source of truth. The UI only reads from Room — never directly from the network. The network writes to Room, which then emits updates to the UI automatically via Flow. This ensures consistency regardless of connectivity.

**Q2: How do you handle a user editing a note while offline?**

> Store the note locally with a `SyncStatus.PENDING` flag. When connectivity returns (observed via `ConnectivityManager` or WorkManager constraints), a `SyncWorker` uploads pending changes. The UI shows a "pending sync" indicator.

**Q3: What is an optimistic update?**

> An optimistic update applies the change to the local database immediately (assuming success), giving instant UI feedback. Then the network request happens in the background. If the request fails, the state is rolled back. This makes the app feel faster and more responsive.

---

## Summary

- Offline-first: UI reads from Room (single source of truth), network writes to Room
- Track `SyncStatus` on entities: `PENDING`, `SYNCED`, `DELETED`
- Use `ConnectivityObserver` + WorkManager for background sync
- Optimistic updates give instant UI feedback
- Always handle conflict resolution — server wins by default unless local is newer

**Next:** [Chapter 4 — Performance Optimization](./04-performance-optimization.md)
