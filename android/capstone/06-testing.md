# Capstone Part 6: Testing

## Test Structure

```
feedflow/
├── feature/headlines/
│   └── src/
│       ├── test/                    ← JVM unit tests
│       │   └── HeadlinesViewModelTest.kt
│       └── androidTest/             ← Instrumented tests
│           └── HeadlinesFragmentTest.kt
├── core/data/
│   └── src/test/
│       └── ArticleRepositoryTest.kt
└── core/testing/
    └── src/main/
        ├── FakeArticleRepository.kt
        ├── MainDispatcherRule.kt
        └── ArticleFactory.kt
```

---

## Test Infrastructure

### `ArticleFactory.kt`

```kotlin
object ArticleFactory {
    fun create(
        id: String = "https://example.com/article/${UUID.randomUUID()}",
        title: String = "Test Article",
        description: String? = "Test description",
        sourceName: String = "Test Source",
        isBookmarked: Boolean = false
    ) = Article(
        id = id,
        title = title,
        description = description,
        url = id,
        imageUrl = null,
        sourceName = sourceName,
        publishedAt = "2024-01-01T00:00:00Z",
        content = null,
        isBookmarked = isBookmarked
    )

    fun createList(count: Int): List<Article> =
        (1..count).map { i -> create(title = "Article $i") }
}
```

### `FakeArticleRepository.kt`

```kotlin
class FakeArticleRepository : ArticleRepository {

    private val articles = MutableStateFlow<List<Article>>(emptyList())
    private val bookmarks = MutableStateFlow<List<Article>>(emptyList())

    var syncCallCount = 0
    var shouldThrowError = false
    var networkDelay = 0L

    fun setArticles(list: List<Article>) { articles.value = list }

    override fun getTopHeadlines(category: NewsCategory?): Flow<PagingData<Article>> {
        return flowOf(PagingData.from(articles.value))
    }

    override fun searchArticles(query: String): Flow<PagingData<Article>> {
        val filtered = articles.value.filter {
            it.title.contains(query, ignoreCase = true)
        }
        return flowOf(PagingData.from(filtered))
    }

    override fun getBookmarkedArticles(): Flow<List<Article>> = bookmarks

    override suspend fun getArticleById(id: String): Article? =
        articles.value.find { it.id == id }

    override suspend fun toggleBookmark(article: Article) {
        if (shouldThrowError) throw IOException("Database error")
        delay(networkDelay)
        val current = bookmarks.value.toMutableList()
        val existing = current.indexOfFirst { it.id == article.id }
        if (existing >= 0) current.removeAt(existing)
        else current.add(article.copy(isBookmarked = true))
        bookmarks.value = current
    }

    override suspend fun syncArticles() {
        if (shouldThrowError) throw IOException("Network error")
        syncCallCount++
    }
}
```

---

## ViewModel Unit Tests

```kotlin
class HeadlinesViewModelTest {

    @get:Rule
    val mainDispatcherRule = MainDispatcherRule()

    private lateinit var viewModel: HeadlinesViewModel
    private val fakeRepository = FakeArticleRepository()

    @Before
    fun setup() {
        fakeRepository.setArticles(ArticleFactory.createList(5))
        viewModel = HeadlinesViewModel(
            getTopHeadlinesUseCase = GetTopHeadlinesUseCase(fakeRepository)
        )
    }

    @Test
    fun `initial state has no category selected`() {
        assertThat(viewModel.state.value.selectedCategory).isNull()
    }

    @Test
    fun `setCategory updates state`() {
        viewModel.setCategory(NewsCategory.TECHNOLOGY)

        assertThat(viewModel.state.value.selectedCategory).isEqualTo(NewsCategory.TECHNOLOGY)
    }

    @Test
    fun `onArticleClicked emits NavigateToDetail event`() = runTest {
        val events = mutableListOf<HeadlinesEvent>()
        val job = launch { viewModel.events.collect { events.add(it) } }

        viewModel.onArticleClicked("article_id_1")

        job.cancel()

        assertThat(events).hasSize(1)
        assertThat(events[0]).isInstanceOf(HeadlinesEvent.NavigateToDetail::class.java)
        assertThat((events[0] as HeadlinesEvent.NavigateToDetail).articleId)
            .isEqualTo("article_id_1")
    }

    @Test
    fun `toggleBookmark with database error emits error event`() = runTest {
        fakeRepository.shouldThrowError = true
        val article = ArticleFactory.create()
        val events = mutableListOf<HeadlinesEvent>()
        val job = launch { viewModel.events.collect { events.add(it) } }

        viewModel.toggleBookmark(article)

        job.cancel()

        val messages = events.filterIsInstance<HeadlinesEvent.ShowMessage>()
        assertThat(messages).isNotEmpty()
    }

    @Test
    fun `setCategory to null shows all articles`() {
        viewModel.setCategory(NewsCategory.TECHNOLOGY)
        viewModel.setCategory(null)

        assertThat(viewModel.state.value.selectedCategory).isNull()
    }
}
```

---

## Repository Integration Tests

```kotlin
@RunWith(AndroidJUnit4::class)
class ArticleRepositoryTest {

    private lateinit var database: FeedFlowDatabase
    private lateinit var repository: ArticleRepositoryImpl

    @Before
    fun setup() {
        database = Room.inMemoryDatabaseBuilder(
            ApplicationProvider.getApplicationContext(),
            FeedFlowDatabase::class.java
        ).allowMainThreadQueries().build()

        repository = ArticleRepositoryImpl(
            api = MockNewsApiService(),
            database = database,
            ioDispatcher = UnconfinedTestDispatcher()
        )
    }

    @After
    fun tearDown() { database.close() }

    @Test
    fun `toggleBookmark saves bookmark in database`() = runTest {
        val article = ArticleFactory.create(id = "https://test.com/1")

        // Insert article first
        database.articleDao().upsertArticles(listOf(article.toEntity()))

        repository.toggleBookmark(article)

        val bookmarks = repository.getBookmarkedArticles().first()
        assertThat(bookmarks).hasSize(1)
        assertThat(bookmarks[0].id).isEqualTo(article.id)
    }

    @Test
    fun `toggleBookmark on bookmarked article removes it`() = runTest {
        val article = ArticleFactory.create(id = "https://test.com/1", isBookmarked = true)
        database.articleDao().upsertArticles(listOf(article.toEntity(isBookmarked = true)))

        repository.toggleBookmark(article)

        val bookmarks = repository.getBookmarkedArticles().first()
        assertThat(bookmarks).isEmpty()
    }
}
```

---

## UI Tests (Espresso)

```kotlin
@HiltAndroidTest
@RunWith(AndroidJUnit4::class)
class HeadlinesFragmentTest {

    @get:Rule(order = 0)
    val hiltRule = HiltAndroidRule(this)

    @get:Rule(order = 1)
    val activityRule = ActivityScenarioRule(MainActivity::class.java)

    @Inject
    lateinit var articleRepository: ArticleRepository

    @Before
    fun setup() {
        hiltRule.inject()
    }

    @Test
    fun headlines_screen_shows_articles() {
        // Articles are loaded via RemoteMediator which uses fake network in tests
        // Verify the RecyclerView is visible and has items
        onView(withId(R.id.rvArticles))
            .check(matches(isDisplayed()))
    }

    @Test
    fun clicking_article_navigates_to_detail() {
        // Wait for list to load, then click first item
        onView(withId(R.id.rvArticles))
            .perform(
                RecyclerViewActions.actionOnItemAtPosition<ArticleAdapter.ArticleViewHolder>(
                    0, click()
                )
            )

        // Verify detail screen is shown
        onView(withId(R.id.tvTitle))
            .check(matches(isDisplayed()))
    }

    @Test
    fun bookmark_button_changes_icon_on_click() {
        // Find bookmark button in first item and click it
        onView(
            RecyclerViewMatcher(R.id.rvArticles)
                .atPositionOnView(0, R.id.btnBookmark)
        ).perform(click())

        // Snackbar confirmation appears
        onView(withText(containsString("bookmark")))
            .check(matches(isDisplayed()))
    }
}

// Test module — replaces real network with mock
@Module
@TestInstallIn(
    components = [SingletonComponent::class],
    replaces = [NetworkModule::class]
)
object TestNetworkModule {
    @Provides @Singleton
    fun provideNewsApiService(): NewsApiService = MockNewsApiService()
}

@Module
@TestInstallIn(
    components = [SingletonComponent::class],
    replaces = [DatabaseModule::class]
)
object TestDatabaseModule {
    @Provides @Singleton
    fun provideDatabase(@ApplicationContext ctx: Context): FeedFlowDatabase =
        Room.inMemoryDatabaseBuilder(ctx, FeedFlowDatabase::class.java)
            .allowMainThreadQueries()
            .build()
}
```

---

## Coverage Goals

| Layer | Coverage Target |
|-------|----------------|
| Domain (use cases) | 100% — pure Kotlin, easy to test |
| Data (repository) | >80% — integration tests with in-memory Room |
| UI (ViewModel) | >80% — unit tests with fake repository |
| UI (Fragment) | Critical flows only — bookmark, navigation, error state |

Run coverage report:

```bash
./gradlew testDebugUnitTestCoverage
```

---

**Next:** [Part 7 — Production Polish](./07-production-polish.md)
