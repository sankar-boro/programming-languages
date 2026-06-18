# Capstone Part 3: API Integration

## NewsAPI Setup

Register at [newsapi.org](https://newsapi.org) for a free API key. Add it to `local.properties`:

```properties
NEWS_API_KEY=your_key_here
```

`app/build.gradle.kts`:

```kotlin
android {
    val localProps = Properties().apply {
        rootProject.file("local.properties").takeIf { it.exists() }?.let {
            load(it.inputStream())
        }
    }

    defaultConfig {
        buildConfigField(
            "String",
            "NEWS_API_KEY",
            "\"${localProps.getProperty("NEWS_API_KEY") ?: System.getenv("NEWS_API_KEY") ?: ""}\""
        )
    }

    buildFeatures {
        buildConfig = true
    }
}
```

---

## Network Module

`core/network/src/main/kotlin/com/yourname/feedflow/core/network/`:

### `NewsApiService.kt`

```kotlin
interface NewsApiService {

    @GET("v2/top-headlines")
    suspend fun getTopHeadlines(
        @Query("country") country: String = "us",
        @Query("category") category: String? = null,
        @Query("page") page: Int = 1,
        @Query("pageSize") pageSize: Int = 20,
        @Query("apiKey") apiKey: String = BuildConfig.NEWS_API_KEY
    ): NewsResponseDto

    @GET("v2/everything")
    suspend fun searchArticles(
        @Query("q") query: String,
        @Query("language") language: String = "en",
        @Query("sortBy") sortBy: String = "publishedAt",
        @Query("page") page: Int = 1,
        @Query("pageSize") pageSize: Int = 20,
        @Query("apiKey") apiKey: String = BuildConfig.NEWS_API_KEY
    ): NewsResponseDto
}
```

### DTOs

```kotlin
// dto/NewsResponseDto.kt
data class NewsResponseDto(
    val status: String,
    val totalResults: Int,
    val articles: List<ArticleDto>
)

// dto/ArticleDto.kt
data class ArticleDto(
    val source: SourceDto,
    val author: String?,
    val title: String?,
    val description: String?,
    val url: String,
    val urlToImage: String?,
    val publishedAt: String,
    val content: String?
)

data class SourceDto(
    val id: String?,
    val name: String
)

// Mapper
fun ArticleDto.toDomain(): Article = Article(
    id = url,
    title = title ?: "No title",
    description = description,
    url = url,
    imageUrl = urlToImage,
    sourceName = source.name,
    publishedAt = publishedAt,
    content = content
)
```

### `NetworkModule.kt`

```kotlin
@Module
@InstallIn(SingletonComponent::class)
object NetworkModule {

    @Provides
    @Singleton
    fun provideOkHttpClient(): OkHttpClient {
        return OkHttpClient.Builder()
            .addInterceptor(
                HttpLoggingInterceptor { message -> Timber.tag("OkHttp").d(message) }.apply {
                    level = if (BuildConfig.DEBUG) HttpLoggingInterceptor.Level.BODY
                            else HttpLoggingInterceptor.Level.NONE
                }
            )
            .connectTimeout(30, TimeUnit.SECONDS)
            .readTimeout(30, TimeUnit.SECONDS)
            .build()
    }

    @Provides
    @Singleton
    fun provideRetrofit(client: OkHttpClient): Retrofit {
        return Retrofit.Builder()
            .baseUrl("https://newsapi.org/")
            .client(client)
            .addConverterFactory(GsonConverterFactory.create())
            .build()
    }

    @Provides
    @Singleton
    fun provideNewsApiService(retrofit: Retrofit): NewsApiService =
        retrofit.create(NewsApiService::class.java)
}
```

---

## Paging Source

`core/data/src/main/kotlin/com/yourname/feedflow/core/data/paging/HeadlinesPagingSource.kt`:

```kotlin
class HeadlinesPagingSource(
    private val api: NewsApiService,
    private val category: String?
) : PagingSource<Int, Article>() {

    override suspend fun load(params: LoadParams<Int>): LoadResult<Int, Article> {
        val page = params.key ?: 1
        return try {
            val response = api.getTopHeadlines(
                category = category,
                page = page,
                pageSize = params.loadSize
            )
            val articles = response.articles
                .filter { it.title != null && it.title != "[Removed]" }
                .map { it.toDomain() }

            LoadResult.Page(
                data = articles,
                prevKey = if (page == 1) null else page - 1,
                nextKey = if (articles.isEmpty()) null else page + 1
            )
        } catch (e: IOException) {
            LoadResult.Error(e)
        } catch (e: HttpException) {
            LoadResult.Error(e)
        }
    }

    override fun getRefreshKey(state: PagingState<Int, Article>): Int? {
        return state.anchorPosition?.let { anchor ->
            state.closestPageToPosition(anchor)?.prevKey?.plus(1)
                ?: state.closestPageToPosition(anchor)?.nextKey?.minus(1)
        }
    }
}
```

---

## Network Error Handling

```kotlin
// core/common/src/main/kotlin/network/NetworkResult.kt
sealed class NetworkResult<out T> {
    data class Success<T>(val data: T) : NetworkResult<T>()
    data class Error(
        val code: Int? = null,
        val message: String?,
        val exception: Throwable? = null
    ) : NetworkResult<Nothing>()
}

suspend fun <T> safeApiCall(
    dispatcher: CoroutineDispatcher = Dispatchers.IO,
    call: suspend () -> T
): NetworkResult<T> = withContext(dispatcher) {
    try {
        NetworkResult.Success(call())
    } catch (e: HttpException) {
        NetworkResult.Error(
            code = e.code(),
            message = e.response()?.errorBody()?.string() ?: e.message()
        )
    } catch (e: IOException) {
        NetworkResult.Error(message = "Network error: ${e.message}", exception = e)
    } catch (e: Exception) {
        NetworkResult.Error(message = e.message, exception = e)
    }
}

// Usage in Repository
override suspend fun syncArticles() {
    val result = safeApiCall { api.getTopHeadlines(pageSize = 50) }
    when (result) {
        is NetworkResult.Success -> {
            val entities = result.data.articles
                .filter { it.title != "[Removed]" }
                .map { it.toEntity() }
            dao.upsertArticles(entities)
        }
        is NetworkResult.Error -> {
            Timber.w("Sync failed: ${result.message}")
        }
    }
}
```

---

**Next:** [Part 4 — Database Layer](./04-database-layer.md)
