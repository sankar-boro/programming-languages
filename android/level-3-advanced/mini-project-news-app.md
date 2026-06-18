# Mini Project: News App

## Overview

Build a production-quality News app using everything from Level 3.

**Features:**
- Paginated news list from a REST API
- Room cache for offline reading
- Hilt dependency injection throughout
- Search with debounce
- Bookmarking articles (Room)
- ViewModel + StateFlow + coroutines
- Unit tests for ViewModel and Repository
- UI tests for the list screen

---

## Architecture

```
:app module
├── di/
│   ├── NetworkModule.kt       ← Retrofit, OkHttp
│   ├── DatabaseModule.kt      ← Room
│   └── RepositoryModule.kt    ← Repository bindings
│
├── data/
│   ├── remote/
│   │   ├── NewsApiService.kt
│   │   └── dto/ArticleDto.kt
│   ├── local/
│   │   ├── NewsDatabase.kt
│   │   ├── ArticleDao.kt
│   │   └── entity/ArticleEntity.kt
│   ├── repository/
│   │   └── NewsRepositoryImpl.kt
│   └── paging/
│       ├── NewsPagingSource.kt
│       └── NewsRemoteMediator.kt
│
├── domain/
│   ├── model/Article.kt
│   └── repository/NewsRepository.kt
│
└── ui/
    ├── news/
    │   ├── NewsFragment.kt
    │   ├── NewsViewModel.kt
    │   └── ArticleAdapter.kt
    └── detail/
        ├── ArticleDetailFragment.kt
        └── ArticleDetailViewModel.kt
```

---

## Dependencies (`app/build.gradle.kts`)

```kotlin
plugins {
    alias(libs.plugins.android.application)
    alias(libs.plugins.kotlin.android)
    alias(libs.plugins.hilt.android)
    alias(libs.plugins.navigation.safeargs.kotlin)
    id("com.google.devtools.ksp")
}

dependencies {
    // AndroidX Core
    implementation(libs.androidx.core.ktx)
    implementation(libs.androidx.appcompat)
    implementation(libs.material)
    implementation(libs.androidx.constraintlayout)

    // Navigation
    implementation(libs.androidx.navigation.fragment.ktx)
    implementation(libs.androidx.navigation.ui.ktx)

    // Lifecycle
    implementation(libs.androidx.lifecycle.viewmodel.ktx)
    implementation(libs.androidx.lifecycle.runtime.ktx)

    // Room
    implementation(libs.androidx.room.runtime)
    implementation(libs.androidx.room.ktx)
    implementation(libs.androidx.room.paging)
    ksp(libs.androidx.room.compiler)

    // Paging
    implementation(libs.androidx.paging.runtime.ktx)

    // Hilt
    implementation(libs.hilt.android)
    ksp(libs.hilt.compiler)
    implementation(libs.hilt.navigation.fragment)

    // Retrofit
    implementation(libs.retrofit)
    implementation(libs.retrofit.converter.gson)
    implementation(libs.okhttp.logging.interceptor)

    // Coroutines
    implementation(libs.kotlinx.coroutines.android)

    // Testing
    testImplementation(libs.junit)
    testImplementation(libs.truth)
    testImplementation(libs.mockito.kotlin)
    testImplementation(libs.kotlinx.coroutines.test)
    testImplementation(libs.androidx.paging.testing)
    androidTestImplementation(libs.androidx.junit)
    androidTestImplementation(libs.androidx.espresso.core)
    androidTestImplementation(libs.hilt.android.testing)
    kspAndroidTest(libs.hilt.compiler)
}
```

---

## Domain Layer

### `Article.kt`

```kotlin
data class Article(
    val id: String,  // URL as unique identifier
    val title: String,
    val description: String?,
    val url: String,
    val imageUrl: String?,
    val sourceName: String,
    val publishedAt: String,
    val isBookmarked: Boolean = false
)
```

### `NewsRepository.kt` (interface)

```kotlin
interface NewsRepository {
    fun getTopHeadlines(query: String): Flow<PagingData<Article>>
    fun getBookmarkedArticles(): Flow<List<Article>>
    suspend fun toggleBookmark(article: Article)
    suspend fun getArticleByUrl(url: String): Article?
}
```

---

## Network Layer

### `NewsApiService.kt`

```kotlin
interface NewsApiService {
    @GET("v2/everything")
    suspend fun getArticles(
        @Query("q") query: String,
        @Query("page") page: Int,
        @Query("pageSize") pageSize: Int,
        @Query("apiKey") apiKey: String = BuildConfig.NEWS_API_KEY
    ): NewsResponse
}

data class NewsResponse(
    val status: String,
    val totalResults: Int,
    val articles: List<ArticleDto>
)

data class ArticleDto(
    val title: String?,
    val description: String?,
    val url: String,
    val urlToImage: String?,
    val publishedAt: String,
    val source: SourceDto
)

data class SourceDto(val id: String?, val name: String)
```

### `NetworkModule.kt`

```kotlin
@Module
@InstallIn(SingletonComponent::class)
object NetworkModule {

    @Provides
    @Singleton
    fun provideOkHttpClient(): OkHttpClient = OkHttpClient.Builder()
        .addInterceptor(HttpLoggingInterceptor().apply {
            level = if (BuildConfig.DEBUG) HttpLoggingInterceptor.Level.BODY
                    else HttpLoggingInterceptor.Level.NONE
        })
        .connectTimeout(30, TimeUnit.SECONDS)
        .readTimeout(30, TimeUnit.SECONDS)
        .build()

    @Provides
    @Singleton
    fun provideRetrofit(okHttpClient: OkHttpClient): Retrofit = Retrofit.Builder()
        .baseUrl("https://newsapi.org/")
        .client(okHttpClient)
        .addConverterFactory(GsonConverterFactory.create())
        .build()

    @Provides
    @Singleton
    fun provideNewsApiService(retrofit: Retrofit): NewsApiService =
        retrofit.create(NewsApiService::class.java)
}
```

---

## Data Layer

### `ArticleEntity.kt`

```kotlin
@Entity(tableName = "articles")
data class ArticleEntity(
    @PrimaryKey
    val url: String,
    val title: String,
    val description: String?,
    val imageUrl: String?,
    val sourceName: String,
    val publishedAt: String,
    val isBookmarked: Boolean = false,
    val cachedAt: Long = System.currentTimeMillis()
)
```

### `ArticleDao.kt`

```kotlin
@Dao
interface ArticleDao {

    @Query("SELECT * FROM articles ORDER BY published_at DESC")
    fun getArticlesPaged(): PagingSource<Int, ArticleEntity>

    @Query("SELECT * FROM articles WHERE is_bookmarked = 1 ORDER BY published_at DESC")
    fun getBookmarkedArticles(): Flow<List<ArticleEntity>>

    @Query("SELECT * FROM articles WHERE url = :url")
    suspend fun getArticleByUrl(url: String): ArticleEntity?

    @Insert(onConflict = OnConflictStrategy.REPLACE)
    suspend fun insertAll(articles: List<ArticleEntity>)

    @Query("UPDATE articles SET is_bookmarked = :bookmarked WHERE url = :url")
    suspend fun setBookmarked(url: String, bookmarked: Boolean)

    @Query("DELETE FROM articles WHERE is_bookmarked = 0")
    suspend fun clearNonBookmarked()
}
```

---

## ViewModel

```kotlin
@HiltViewModel
class NewsViewModel @Inject constructor(
    private val repository: NewsRepository
) : ViewModel() {

    private val _searchQuery = MutableStateFlow("android")
    val searchQuery: StateFlow<String> = _searchQuery.asStateFlow()

    val articles: Flow<PagingData<Article>> = _searchQuery
        .debounce(500)
        .distinctUntilChanged()
        .flatMapLatest { query ->
            repository.getTopHeadlines(query)
        }
        .cachedIn(viewModelScope)

    val bookmarkedArticles: StateFlow<List<Article>> = repository
        .getBookmarkedArticles()
        .stateIn(viewModelScope, SharingStarted.WhileSubscribed(5_000), emptyList())

    fun setSearchQuery(query: String) {
        _searchQuery.value = query
    }

    fun toggleBookmark(article: Article) {
        viewModelScope.launch {
            repository.toggleBookmark(article)
        }
    }
}
```

---

## Level 3 Checkpoint

Before moving to Level 4, confirm you can:

- [ ] Set up Hilt in an app with Application, Activities, Fragments, and ViewModels
- [ ] Write `@Module` with `@Provides` for Retrofit, Room, and Repository
- [ ] Use `WorkManager` to schedule a background sync task
- [ ] Use `viewModelScope.launch` with proper `try-catch` for coroutines
- [ ] Replace `SharedPreferences` with `DataStore`
- [ ] Write ViewModel unit tests with `FakeRepository` and `MainDispatcherRule`
- [ ] Write Espresso UI tests with `ActivityScenarioRule`
- [ ] Implement `PagingDataAdapter` with `LoadState` handling

**Next Level:** [Level 4 — Expert: Architecture and Production](../level-4-expert/01-clean-architecture.md)
