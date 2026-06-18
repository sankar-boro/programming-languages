# Capstone Part 4: Database Layer

## Room Setup

`core/database/src/main/kotlin/com/yourname/feedflow/core/database/`:

### `ArticleEntity.kt`

```kotlin
@Entity(
    tableName = "articles",
    indices = [
        Index(value = ["published_at"]),
        Index(value = ["category"]),
        Index(value = ["is_bookmarked"])
    ]
)
data class ArticleEntity(
    @PrimaryKey
    val id: String,  // URL
    val title: String,
    val description: String?,
    val url: String,
    @ColumnInfo(name = "image_url")
    val imageUrl: String?,
    @ColumnInfo(name = "source_name")
    val sourceName: String,
    @ColumnInfo(name = "published_at")
    val publishedAt: String,
    val content: String?,
    val category: String?,
    @ColumnInfo(name = "is_bookmarked")
    val isBookmarked: Boolean = false,
    @ColumnInfo(name = "cached_at")
    val cachedAt: Long = System.currentTimeMillis()
)
```

### `RemoteKeyEntity.kt` (for RemoteMediator)

```kotlin
@Entity(tableName = "remote_keys")
data class RemoteKeyEntity(
    @PrimaryKey
    val articleId: String,
    val prevPage: Int?,
    val nextPage: Int?,
    val category: String?
)
```

### `ArticleDao.kt`

```kotlin
@Dao
interface ArticleDao {

    @Query("""
        SELECT * FROM articles 
        WHERE (:category IS NULL OR category = :category)
        ORDER BY published_at DESC
    """)
    fun getArticlesPaged(category: String?): PagingSource<Int, ArticleEntity>

    @Query("SELECT * FROM articles WHERE is_bookmarked = 1 ORDER BY cached_at DESC")
    fun getBookmarkedArticles(): Flow<List<ArticleEntity>>

    @Query("SELECT * FROM articles WHERE id = :id")
    suspend fun getArticleById(id: String): ArticleEntity?

    @Insert(onConflict = OnConflictStrategy.REPLACE)
    suspend fun upsertArticles(articles: List<ArticleEntity>)

    @Query("UPDATE articles SET is_bookmarked = :bookmarked WHERE id = :id")
    suspend fun setBookmarked(id: String, bookmarked: Boolean)

    @Query("DELETE FROM articles WHERE is_bookmarked = 0 AND category = :category")
    suspend fun clearNonBookmarkedForCategory(category: String?)

    @Query("DELETE FROM articles WHERE is_bookmarked = 0")
    suspend fun clearNonBookmarked()
}
```

### `RemoteKeyDao.kt`

```kotlin
@Dao
interface RemoteKeyDao {
    @Insert(onConflict = OnConflictStrategy.REPLACE)
    suspend fun insertAll(keys: List<RemoteKeyEntity>)

    @Query("SELECT * FROM remote_keys WHERE article_id = :articleId")
    suspend fun getRemoteKey(articleId: String): RemoteKeyEntity?

    @Query("DELETE FROM remote_keys WHERE category = :category")
    suspend fun clearByCategory(category: String?)

    @Query("DELETE FROM remote_keys")
    suspend fun clearAll()
}
```

### `FeedFlowDatabase.kt`

```kotlin
@Database(
    entities = [ArticleEntity::class, RemoteKeyEntity::class],
    version = 1,
    exportSchema = true
)
abstract class FeedFlowDatabase : RoomDatabase() {

    abstract fun articleDao(): ArticleDao
    abstract fun remoteKeyDao(): RemoteKeyDao

    companion object {
        @Volatile private var INSTANCE: FeedFlowDatabase? = null

        fun getInstance(context: Context): FeedFlowDatabase =
            INSTANCE ?: synchronized(this) {
                Room.databaseBuilder(
                    context.applicationContext,
                    FeedFlowDatabase::class.java,
                    "feedflow.db"
                ).build().also { INSTANCE = it }
            }
    }
}
```

---

## Remote Mediator (Offline-First Paging)

```kotlin
@OptIn(ExperimentalPagingApi::class)
class ArticleRemoteMediator(
    private val category: String?,
    private val api: NewsApiService,
    private val database: FeedFlowDatabase
) : RemoteMediator<Int, ArticleEntity>() {

    override suspend fun initialize(): InitializeAction {
        val lastSync = database.articleDao().getOldestCachedAt(category)
        val cacheTimeout = TimeUnit.HOURS.toMillis(1)
        return if (lastSync != null && System.currentTimeMillis() - lastSync < cacheTimeout) {
            InitializeAction.SKIP_INITIAL_REFRESH
        } else {
            InitializeAction.LAUNCH_INITIAL_REFRESH
        }
    }

    override suspend fun load(
        loadType: LoadType,
        state: PagingState<Int, ArticleEntity>
    ): MediatorResult {
        return try {
            val page = when (loadType) {
                LoadType.REFRESH -> 1
                LoadType.PREPEND -> return MediatorResult.Success(endOfPaginationReached = true)
                LoadType.APPEND -> {
                    val lastItem = state.lastItemOrNull()
                        ?: return MediatorResult.Success(endOfPaginationReached = true)
                    database.remoteKeyDao().getRemoteKey(lastItem.id)?.nextPage
                        ?: return MediatorResult.Success(endOfPaginationReached = true)
                }
            }

            val response = api.getTopHeadlines(
                category = category,
                page = page,
                pageSize = state.config.pageSize
            )

            val articles = response.articles
                .filter { it.title != null && it.title != "[Removed]" }

            database.withTransaction {
                if (loadType == LoadType.REFRESH) {
                    database.remoteKeyDao().clearByCategory(category)
                    database.articleDao().clearNonBookmarkedForCategory(category)
                }

                val remoteKeys = articles.map { dto ->
                    RemoteKeyEntity(
                        articleId = dto.url,
                        prevPage = if (page == 1) null else page - 1,
                        nextPage = if (articles.isEmpty()) null else page + 1,
                        category = category
                    )
                }
                database.remoteKeyDao().insertAll(remoteKeys)
                database.articleDao().upsertArticles(
                    articles.map { it.toEntity(category) }
                )
            }

            MediatorResult.Success(endOfPaginationReached = articles.isEmpty())
        } catch (e: IOException) {
            MediatorResult.Error(e)
        } catch (e: HttpException) {
            MediatorResult.Error(e)
        }
    }
}
```

---

## Repository Implementation

```kotlin
class ArticleRepositoryImpl @Inject constructor(
    private val api: NewsApiService,
    private val database: FeedFlowDatabase,
    @IoDispatcher private val ioDispatcher: CoroutineDispatcher
) : ArticleRepository {

    private val pagingConfig = PagingConfig(
        pageSize = 20,
        enablePlaceholders = false,
        prefetchDistance = 5,
        initialLoadSize = 40
    )

    @OptIn(ExperimentalPagingApi::class)
    override fun getTopHeadlines(category: NewsCategory?): Flow<PagingData<Article>> {
        val categoryValue = category?.apiValue
        return Pager(
            config = pagingConfig,
            remoteMediator = ArticleRemoteMediator(categoryValue, api, database),
            pagingSourceFactory = { database.articleDao().getArticlesPaged(categoryValue) }
        ).flow.map { pagingData ->
            pagingData.map { entity -> entity.toDomain() }
        }
    }

    override fun searchArticles(query: String): Flow<PagingData<Article>> {
        return Pager(
            config = pagingConfig,
            pagingSourceFactory = { SearchPagingSource(api, query) }
        ).flow.map { pagingData ->
            pagingData.map { it.toDomain() }
        }
    }

    override fun getBookmarkedArticles(): Flow<List<Article>> =
        database.articleDao().getBookmarkedArticles()
            .map { entities -> entities.map { it.toDomain() } }
            .flowOn(ioDispatcher)

    override suspend fun getArticleById(id: String): Article? =
        withContext(ioDispatcher) {
            database.articleDao().getArticleById(id)?.toDomain()
        }

    override suspend fun toggleBookmark(article: Article) = withContext(ioDispatcher) {
        database.articleDao().setBookmarked(article.id, !article.isBookmarked)
    }
}
```

---

## Database Hilt Module

```kotlin
@Module
@InstallIn(SingletonComponent::class)
object DatabaseModule {

    @Provides
    @Singleton
    fun provideDatabase(@ApplicationContext context: Context): FeedFlowDatabase =
        FeedFlowDatabase.getInstance(context)

    @Provides
    fun provideArticleDao(db: FeedFlowDatabase): ArticleDao = db.articleDao()

    @Provides
    fun provideRemoteKeyDao(db: FeedFlowDatabase): RemoteKeyDao = db.remoteKeyDao()
}

@Module
@InstallIn(SingletonComponent::class)
abstract class RepositoryModule {

    @Binds
    @Singleton
    abstract fun bindArticleRepository(impl: ArticleRepositoryImpl): ArticleRepository
}
```

---

**Next:** [Part 5 — Navigation Flow](./05-navigation-flow.md)
