# Chapter 6: Room Database

## What Is Room?

Room is AndroidX's SQLite ORM (Object-Relational Mapper). It provides compile-time SQL verification, coroutine/Flow support, migration tools, and type safety — all on top of Android's built-in SQLite.

**Analogy:** SQLite is the raw engine. Room is like a well-organized dealership — you interact through clean contracts (DAOs) and never touch the raw parts directly.

```kotlin
implementation("androidx.room:room-runtime:2.6.1")
implementation("androidx.room:room-ktx:2.6.1")
kapt("androidx.room:room-compiler:2.6.1")  // for Java/Kotlin with kapt
// OR with KSP (preferred for Kotlin):
ksp("androidx.room:room-compiler:2.6.1")
```

For KSP support, add the KSP plugin to `build.gradle.kts`:
```kotlin
plugins {
    id("com.google.devtools.ksp") version "2.1.0-1.0.29"
}
```

---

## Room Architecture

```
┌──────────────────────────────────────┐
│                 App                  │
│                                      │
│  Repository → DAO → Room Database   │
│                  ↕                   │
│              SQLite                  │
└──────────────────────────────────────┘
```

Room has three main components:

| Component | Role |
|-----------|------|
| `@Entity` | A data class mapped to a database table |
| `@Dao` | Interface with SQL queries as annotated functions |
| `@Database` | The Room database entry point; holds DAOs |

---

## Step 1: Define Entities

```kotlin
import androidx.room.Entity
import androidx.room.PrimaryKey
import androidx.room.ColumnInfo
import androidx.room.Embedded
import androidx.room.Relation

@Entity(tableName = "tasks")
data class TaskEntity(
    @PrimaryKey(autoGenerate = true)
    val id: Long = 0,

    @ColumnInfo(name = "title")
    val title: String,

    @ColumnInfo(name = "description")
    val description: String = "",

    @ColumnInfo(name = "is_completed")
    val isCompleted: Boolean = false,

    @ColumnInfo(name = "created_at")
    val createdAt: Long = System.currentTimeMillis(),

    @ColumnInfo(name = "due_date")
    val dueDate: Long? = null
)
```

```java
// Java version
@Entity(tableName = "tasks")
public class TaskEntity {
    @PrimaryKey(autoGenerate = true)
    public long id;

    @ColumnInfo(name = "title")
    public String title;

    @ColumnInfo(name = "is_completed")
    public boolean isCompleted;
}
```

---

## Step 2: Define the DAO

```kotlin
import androidx.room.*
import kotlinx.coroutines.flow.Flow

@Dao
interface TaskDao {

    // Flow automatically re-emits when the table changes
    @Query("SELECT * FROM tasks ORDER BY created_at DESC")
    fun getAllTasks(): Flow<List<TaskEntity>>

    @Query("SELECT * FROM tasks WHERE is_completed = 0 ORDER BY created_at DESC")
    fun getActiveTasks(): Flow<List<TaskEntity>>

    @Query("SELECT * FROM tasks WHERE id = :taskId")
    suspend fun getTaskById(taskId: Long): TaskEntity?

    @Query("SELECT * FROM tasks WHERE title LIKE '%' || :query || '%'")
    fun searchTasks(query: String): Flow<List<TaskEntity>>

    @Insert(onConflict = OnConflictStrategy.REPLACE)
    suspend fun insertTask(task: TaskEntity): Long

    @Insert(onConflict = OnConflictStrategy.REPLACE)
    suspend fun insertAll(tasks: List<TaskEntity>)

    @Update
    suspend fun updateTask(task: TaskEntity)

    @Delete
    suspend fun deleteTask(task: TaskEntity)

    @Query("DELETE FROM tasks WHERE id = :taskId")
    suspend fun deleteTaskById(taskId: Long)

    @Query("DELETE FROM tasks WHERE is_completed = 1")
    suspend fun deleteCompletedTasks()

    @Query("UPDATE tasks SET is_completed = :completed WHERE id = :taskId")
    suspend fun setCompleted(taskId: Long, completed: Boolean)

    @Query("SELECT COUNT(*) FROM tasks")
    suspend fun getTaskCount(): Int
}
```

---

## Step 3: Define the Database

```kotlin
import androidx.room.Database
import androidx.room.Room
import androidx.room.RoomDatabase
import android.content.Context

@Database(
    entities = [TaskEntity::class],  // Add all @Entity classes here
    version = 1,
    exportSchema = true  // Exports schema to a JSON file for migration tracking
)
abstract class AppDatabase : RoomDatabase() {

    abstract fun taskDao(): TaskDao

    companion object {
        @Volatile
        private var INSTANCE: AppDatabase? = null

        fun getInstance(context: Context): AppDatabase {
            return INSTANCE ?: synchronized(this) {
                Room.databaseBuilder(
                    context.applicationContext,
                    AppDatabase::class.java,
                    "app_database"
                )
                .fallbackToDestructiveMigration()  // Dev only — destroys data on schema change
                .build()
                .also { INSTANCE = it }
            }
        }
    }
}
```

> In production, replace `fallbackToDestructiveMigration()` with proper migrations.

---

## Step 4: Repository Pattern

The repository abstracts data sources from the ViewModel:

```kotlin
class TaskRepository(private val taskDao: TaskDao) {

    val allTasks: Flow<List<TaskEntity>> = taskDao.getAllTasks()
    val activeTasks: Flow<List<TaskEntity>> = taskDao.getActiveTasks()

    suspend fun insertTask(task: TaskEntity): Long = taskDao.insertTask(task)

    suspend fun updateTask(task: TaskEntity) = taskDao.updateTask(task)

    suspend fun deleteTask(task: TaskEntity) = taskDao.deleteTask(task)

    suspend fun toggleComplete(taskId: Long, currentState: Boolean) {
        taskDao.setCompleted(taskId, !currentState)
    }

    fun searchTasks(query: String): Flow<List<TaskEntity>> = taskDao.searchTasks(query)
}
```

---

## Step 5: ViewModel with Room + Coroutines

```kotlin
class TaskViewModel(
    private val repository: TaskRepository
) : ViewModel() {

    // Room Flow → StateFlow, automatically reflects DB changes
    val tasks: StateFlow<List<TaskEntity>> = repository.allTasks
        .stateIn(
            scope = viewModelScope,
            started = SharingStarted.WhileSubscribed(5_000),
            initialValue = emptyList()
        )

    fun addTask(title: String, description: String) {
        viewModelScope.launch {
            repository.insertTask(
                TaskEntity(title = title, description = description)
            )
        }
    }

    fun toggleComplete(task: TaskEntity) {
        viewModelScope.launch {
            repository.toggleComplete(task.id, task.isCompleted)
        }
    }

    fun deleteTask(task: TaskEntity) {
        viewModelScope.launch {
            repository.deleteTask(task)
        }
    }
}
```

---

## Type Converters

Room can only store primitive types, Strings, and byte arrays. For complex types, use `@TypeConverter`:

```kotlin
import androidx.room.TypeConverter
import java.util.Date

class DateConverter {
    @TypeConverter
    fun fromTimestamp(value: Long?): Date? = value?.let { Date(it) }

    @TypeConverter
    fun toTimestamp(date: Date?): Long? = date?.time
}

// Register in the database:
@Database(entities = [...], version = 1)
@TypeConverters(DateConverter::class)
abstract class AppDatabase : RoomDatabase()
```

---

## Database Migrations

When you change the schema (add a column, rename a table), you must provide a migration:

```kotlin
val MIGRATION_1_2 = object : Migration(1, 2) {
    override fun migrate(database: SupportSQLiteDatabase) {
        // Add a new column to the tasks table
        database.execSQL(
            "ALTER TABLE tasks ADD COLUMN priority INTEGER NOT NULL DEFAULT 0"
        )
    }
}

val MIGRATION_2_3 = object : Migration(2, 3) {
    override fun migrate(database: SupportSQLiteDatabase) {
        database.execSQL(
            "CREATE TABLE IF NOT EXISTS tags (" +
            "id INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL, " +
            "name TEXT NOT NULL)"
        )
    }
}

Room.databaseBuilder(context, AppDatabase::class.java, "app_database")
    .addMigrations(MIGRATION_1_2, MIGRATION_2_3)
    .build()
```

---

## Relations: One-to-Many

```kotlin
@Entity(tableName = "categories")
data class CategoryEntity(
    @PrimaryKey(autoGenerate = true)
    val id: Long = 0,
    val name: String
)

@Entity(
    tableName = "tasks",
    foreignKeys = [ForeignKey(
        entity = CategoryEntity::class,
        parentColumns = ["id"],
        childColumns = ["category_id"],
        onDelete = ForeignKey.CASCADE
    )]
)
data class TaskEntity(
    @PrimaryKey(autoGenerate = true)
    val id: Long = 0,
    val title: String,
    val categoryId: Long
)

// Relationship data class
data class CategoryWithTasks(
    @Embedded val category: CategoryEntity,
    @Relation(
        parentColumn = "id",
        entityColumn = "category_id"
    )
    val tasks: List<TaskEntity>
)

// DAO query
@Transaction
@Query("SELECT * FROM categories")
fun getCategoriesWithTasks(): Flow<List<CategoryWithTasks>>
```

---

## Common Mistakes

### Mistake 1: Running Room queries on the main thread

```kotlin
// CRASH — Room throws IllegalStateException if called on main thread
val tasks = taskDao.getAllTasks()

// CORRECT — use suspend functions or Flow
viewModelScope.launch {
    val task = taskDao.getTaskById(id)  // suspend
}
// OR
taskDao.getAllTasks().collect { ... }  // Flow
```

### Mistake 2: Not using `@Transaction` for queries returning relations

```kotlin
// WRONG — Room may see inconsistent state across multiple queries
@Query("SELECT * FROM categories")
fun getCategoriesWithTasks(): Flow<List<CategoryWithTasks>>

// CORRECT
@Transaction
@Query("SELECT * FROM categories")
fun getCategoriesWithTasks(): Flow<List<CategoryWithTasks>>
```

### Mistake 3: Exposing the database directly instead of using a Repository

```kotlin
// WRONG — ViewModel directly depends on Room, hard to test
class ViewModel(db: AppDatabase) {
    val tasks = db.taskDao().getAllTasks()
}

// CORRECT — Repository abstracts data source
class ViewModel(repo: TaskRepository) {
    val tasks = repo.allTasks
}
```

---

## Interview Questions

**Q1: What are the three main components of Room?**

> `@Entity` maps a Kotlin data class to a database table. `@Dao` is an interface with annotated methods that compile down to SQL statements. `@Database` is the abstract class that ties them together and provides access to DAO instances.

**Q2: Why does Room's DAO return `Flow` instead of `List`?**

> `Flow<List<T>>` re-emits automatically whenever the underlying table changes, enabling reactive UI updates. Your UI just collects the Flow — no manual refresh needed.

**Q3: What is the Repository pattern and why use it?**

> The Repository is an abstraction layer between ViewModels and data sources (Room, network APIs). It hides where data comes from, makes ViewModels testable (swap real Room with a fake), and allows combining multiple sources cleanly.

**Q4: What happens if you change the database schema without providing a migration?**

> Room throws `IllegalStateException` on the old device at runtime. If `fallbackToDestructiveMigration()` is set, it destroys and rebuilds the database — users lose all data. Always provide migrations in production.

---

## Summary

- Room has three parts: `@Entity` (table), `@Dao` (queries), `@Database` (entry point)
- Return `Flow` from DAO queries for reactive, auto-updating data
- Use the Repository pattern to abstract Room from ViewModel
- Always provide `Migration` objects when changing the database schema
- Use `@TypeConverter` for complex types (dates, enums, lists)

**Next:** [Chapter 7 — Jetpack Compose Introduction](./07-jetpack-compose-intro.md)
