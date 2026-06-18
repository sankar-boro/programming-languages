# Chapter 2: WorkManager

## What Is WorkManager?

WorkManager is AndroidX's library for **deferrable, guaranteed background work** — tasks that must complete even if the app exits or the device restarts.

**Use WorkManager for:**
- Uploading logs or analytics batches
- Syncing data with a server
- Processing images (compressing, uploading)
- Sending deferred notifications
- Periodic database cleanup

**Do NOT use WorkManager for:**
- Work that must execute immediately (use coroutines directly)
- Exact-time alarms (use `AlarmManager`)
- Long-running foreground services (use `ForegroundService`)

```kotlin
implementation("androidx.work:work-runtime-ktx:2.9.1")
```

---

## Core Concepts

| Concept | What It Is |
|---------|-----------|
| `Worker` / `CoroutineWorker` | The unit of work — what to do |
| `WorkRequest` | Wraps the Worker + constraints + retry policy |
| `WorkManager` | Schedules and manages WorkRequests |
| `Constraints` | Conditions that must be met before work runs |
| `WorkInfo` | Status and output data of a request |

---

## Creating a Worker

```kotlin
import androidx.work.CoroutineWorker
import androidx.work.WorkerParameters
import android.content.Context

class SyncNotesWorker(
    context: Context,
    params: WorkerParameters
) : CoroutineWorker(context, params) {

    override suspend fun doWork(): Result {
        return try {
            // Receive input data
            val userId = inputData.getString(KEY_USER_ID)
                ?: return Result.failure()

            // Do the actual work — runs on a background thread automatically
            val api = RetrofitClient.notesApi
            val remoteNotes = api.fetchNotes(userId)

            // Save to database
            val db = NoteDatabase.getInstance(applicationContext)
            db.noteDao().insertAll(remoteNotes.map { it.toEntity() })

            // Return success with output data
            val outputData = workDataOf(KEY_SYNC_COUNT to remoteNotes.size)
            Result.success(outputData)

        } catch (e: IOException) {
            // Retry on network errors
            if (runAttemptCount < 3) Result.retry()
            else Result.failure()
        } catch (e: Exception) {
            Result.failure()
        }
    }

    companion object {
        const val KEY_USER_ID = "user_id"
        const val KEY_SYNC_COUNT = "sync_count"
    }
}
```

Three possible results:
- `Result.success()` — Work completed successfully
- `Result.failure()` — Work failed, don't retry
- `Result.retry()` — Work failed, try again (backoff policy applies)

---

## Scheduling a One-Time Request

```kotlin
import androidx.work.*

val inputData = workDataOf(
    SyncNotesWorker.KEY_USER_ID to "user_123"
)

val constraints = Constraints.Builder()
    .setRequiredNetworkType(NetworkType.CONNECTED)  // Only run on network
    .setRequiresBatteryNotLow(true)                 // Not when battery is critical
    .setRequiresStorageNotLow(true)                 // Not when storage is low
    .build()

val syncRequest = OneTimeWorkRequestBuilder<SyncNotesWorker>()
    .setInputData(inputData)
    .setConstraints(constraints)
    .setBackoffCriteria(
        BackoffPolicy.EXPONENTIAL,
        WorkRequest.MIN_BACKOFF_MILLIS,
        TimeUnit.MILLISECONDS
    )
    .addTag("sync_notes")  // Tag for querying/cancelling
    .build()

WorkManager.getInstance(context).enqueue(syncRequest)
```

---

## Scheduling Periodic Work

```kotlin
val periodicSync = PeriodicWorkRequestBuilder<SyncNotesWorker>(
    repeatInterval = 1,
    repeatIntervalTimeUnit = TimeUnit.HOURS,
    flexTimeInterval = 15,         // Can run in the last 15 min of the interval
    flexTimeIntervalUnit = TimeUnit.MINUTES
)
    .setConstraints(Constraints.Builder()
        .setRequiredNetworkType(NetworkType.CONNECTED)
        .build())
    .addTag("periodic_sync")
    .build()

WorkManager.getInstance(context).enqueueUniquePeriodicWork(
    "notes_sync",                        // Unique name — prevents duplicates
    ExistingPeriodicWorkPolicy.KEEP,     // KEEP existing, or UPDATE/CANCEL_AND_REENQUEUE
    periodicSync
)
```

---

## Observing Work Status

```kotlin
// By unique work name
WorkManager.getInstance(context)
    .getWorkInfosForUniqueWorkLiveData("notes_sync")
    .observe(viewLifecycleOwner) { workInfoList ->
        val workInfo = workInfoList?.firstOrNull() ?: return@observe
        when (workInfo.state) {
            WorkInfo.State.RUNNING -> showSyncProgress()
            WorkInfo.State.SUCCEEDED -> {
                val count = workInfo.outputData.getInt(
                    SyncNotesWorker.KEY_SYNC_COUNT, 0
                )
                showSyncSuccess("Synced $count notes")
            }
            WorkInfo.State.FAILED -> showSyncError()
            WorkInfo.State.ENQUEUED -> showSyncPending()
            else -> Unit
        }
    }

// By tag
WorkManager.getInstance(context)
    .getWorkInfosByTagLiveData("sync_notes")
    .observe(viewLifecycleOwner) { list -> /* ... */ }
```

---

## Chaining Work

WorkManager supports sequential and parallel chains:

```kotlin
val compressImages = OneTimeWorkRequestBuilder<CompressImagesWorker>().build()
val uploadImages = OneTimeWorkRequestBuilder<UploadImagesWorker>().build()
val notifyServer = OneTimeWorkRequestBuilder<NotifyServerWorker>().build()

// Sequential chain
WorkManager.getInstance(context)
    .beginWith(compressImages)
    .then(uploadImages)
    .then(notifyServer)
    .enqueue()

// Parallel then sequential
val downloadA = OneTimeWorkRequestBuilder<DownloadAWorker>().build()
val downloadB = OneTimeWorkRequestBuilder<DownloadBWorker>().build()
val merge = OneTimeWorkRequestBuilder<MergeWorker>().build()

WorkManager.getInstance(context)
    .beginWith(listOf(downloadA, downloadB))  // Run in parallel
    .then(merge)                               // Run after both complete
    .enqueue()
```

---

## Injecting Dependencies with Hilt

WorkManager workers need dependencies. Use `HiltWorker`:

```kotlin
@HiltWorker
class SyncNotesWorker @AssistedInject constructor(
    @Assisted context: Context,
    @Assisted params: WorkerParameters,
    private val repository: NoteRepository,
    private val api: NotesApiService
) : CoroutineWorker(context, params) {

    override suspend fun doWork(): Result {
        return try {
            val remoteNotes = api.fetchNotes()
            repository.syncNotes(remoteNotes)
            Result.success()
        } catch (e: Exception) {
            if (runAttemptCount < 3) Result.retry() else Result.failure()
        }
    }
}
```

Configure WorkManager to use Hilt:

```kotlin
@HiltAndroidApp
class MyApplication : Application(), Configuration.Provider {

    @Inject
    lateinit var workerFactory: HiltWorkerFactory

    override val workManagerConfiguration: Configuration
        get() = Configuration.Builder()
            .setWorkerFactory(workerFactory)
            .build()
}
```

Remove `WorkManager` auto-initialization from manifest:

```xml
<provider
    android:name="androidx.startup.InitializationProvider"
    android:authorities="${applicationId}.androidx-startup"
    android:exported="false"
    tools:node="merge">
    <meta-data
        android:name="androidx.work.WorkManagerInitializer"
        android:value="androidx.startup"
        tools:node="remove" />
</provider>
```

Add Hilt WorkManager dependency:
```kotlin
implementation("androidx.hilt:hilt-work:1.2.0")
ksp("androidx.hilt:hilt-compiler:1.2.0")
```

---

## Cancelling Work

```kotlin
val workManager = WorkManager.getInstance(context)

// Cancel by tag
workManager.cancelAllWorkByTag("sync_notes")

// Cancel by unique name
workManager.cancelUniqueWork("notes_sync")

// Cancel specific request by ID
workManager.cancelWorkById(workRequest.id)

// Cancel ALL work (use carefully)
workManager.cancelAllWork()
```

---

## Foreground Work (Long-Running)

For work that takes >10 minutes or needs to show a notification:

```kotlin
class LongSyncWorker(context: Context, params: WorkerParameters)
    : CoroutineWorker(context, params) {

    override suspend fun doWork(): Result {
        setForeground(createForegroundInfo())
        // Long-running work here
        return Result.success()
    }

    private fun createForegroundInfo(): ForegroundInfo {
        val notification = NotificationCompat.Builder(applicationContext, CHANNEL_ID)
            .setContentTitle("Syncing notes...")
            .setSmallIcon(R.drawable.ic_sync)
            .setProgress(0, 0, true)
            .build()

        return ForegroundInfo(NOTIFICATION_ID, notification)
    }
}
```

---

## Common Mistakes

### Mistake 1: Using WorkManager for immediate one-shot work

```kotlin
// WRONG — WorkManager adds overhead; use coroutines for immediate work
WorkManager.getInstance(context).enqueue(
    OneTimeWorkRequestBuilder<QuickSaveWorker>().build()
)

// CORRECT — immediate one-shot work
viewModelScope.launch {
    repository.saveNote(note)
}
```

### Mistake 2: Not handling `Result.retry()` correctly

```kotlin
// WRONG — retries infinitely
override suspend fun doWork(): Result {
    return try { ... Result.success() }
    catch (e: IOException) { Result.retry() }
}

// CORRECT — limit retries
catch (e: IOException) {
    if (runAttemptCount < 3) Result.retry() else Result.failure()
}
```

---

## Interview Questions

**Q1: What is the difference between WorkManager and a coroutine for background work?**

> Coroutines are cancelled if the process dies — they don't survive app exit or device restart. WorkManager guarantees execution even if the app is killed or the device reboots, using the appropriate system mechanism (JobScheduler, AlarmManager, etc.) under the hood.

**Q2: When would you use `enqueueUniquePeriodicWork` with `KEEP` policy?**

> When you want to schedule a periodic task (like hourly sync) and ensure only one instance runs — new calls are ignored if the work is already scheduled. `CANCEL_AND_REENQUEUE` would reset the timer, which may skip an imminent sync.

**Q3: What are `Constraints` in WorkManager?**

> Constraints specify conditions that must be true before a Worker runs — e.g., network available, battery not low, device charging. WorkManager queues the work and waits until all constraints are met.

---

## Summary

- WorkManager is for guaranteed, deferrable background work that survives process death
- Use `CoroutineWorker` for Kotlin-native async work
- Use `Constraints` to specify execution conditions
- Use `enqueueUniquePeriodicWork` for recurring tasks to prevent duplicate scheduling
- Chain work sequentially or in parallel with `beginWith().then()`
- Inject dependencies into Workers with `@HiltWorker`

**Next:** [Chapter 3 — Coroutines and AndroidX](./03-coroutines-and-androidx.md)
