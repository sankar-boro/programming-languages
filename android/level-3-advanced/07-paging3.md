# Chapter 7: Paging 3 Library

## What Is Paging 3?

Paging 3 loads large datasets in small chunks (pages), allowing RecyclerView to display thousands of items without loading them all at once.

**Use Paging when:**
- Loading a list from a paginated API (page 1, page 2, etc.)
- Loading large datasets from Room
- Combining network + Room with caching (RemoteMediator)

```kotlin
implementation("androidx.paging:paging-runtime-ktx:3.3.4")
implementation("androidx.paging:paging-compose:3.3.4")     // For Compose
testImplementation("androidx.paging:paging-testing:3.3.4")  // For tests
```

---

## Paging 3 Architecture

```
┌─────────────────────────────────────────────────────────┐
│                         UI Layer                         │
│  PagingDataAdapter  ←  Flow<PagingData<T>>  ←  ViewModel│
└─────────────────────────────────────────────────────────┘
                              ↑
┌─────────────────────────────────────────────────────────┐
│                      Domain Layer                        │
│           Pager(config, pagingSourceFactory)             │
└─────────────────────────────────────────────────────────┘
                              ↑
┌─────────────────────────────────────────────────────────┐
│                       Data Layer                         │
│    PagingSource (network/database/combined)              │
└─────────────────────────────────────────────────────────┘
```

---

## Step 1: Create a PagingSource

`PagingSource` defines how to load one page of data:

```kotlin
import androidx.paging.PagingSource
import androidx.paging.PagingState

class NewsPagingSource(
    private val api: NewsApiService,
    private val query: String
) : PagingSource<Int, Article>() {

    override suspend fun load(params: LoadParams<Int>): LoadResult<Int, Article> {
        val page = params.key ?: 1  // Start from page 1

        return try {
            val response = api.getArticles(
                query = query,
                page = page,
                pageSize = params.loadSize
            )

            LoadResult.Page(
                data = response.articles,
                prevKey = if (page == 1) null else page - 1,
                nextKey = if (response.articles.isEmpty()) null else page + 1
            )
        } catch (e: IOException) {
            LoadResult.Error(e)
        } catch (e: HttpException) {
            LoadResult.Error(e)
        }
    }

    // Called when the paging library needs to re-run the most recent load
    override fun getRefreshKey(state: PagingState<Int, Article>): Int? {
        return state.anchorPosition?.let { anchor ->
            state.closestPageToPosition(anchor)?.prevKey?.plus(1)
                ?: state.closestPageToPosition(anchor)?.nextKey?.minus(1)
        }
    }
}
```

---

## Step 2: Create the Pager in Repository

```kotlin
class NewsRepository @Inject constructor(
    private val api: NewsApiService
) {
    private val pagingConfig = PagingConfig(
        pageSize = 20,              // Items per page
        enablePlaceholders = false,  // Don't show placeholder items
        prefetchDistance = 5,       // Load next page when 5 items from end
        initialLoadSize = 40        // First load: fetch 2 pages worth
    )

    fun getTopHeadlines(query: String): Flow<PagingData<Article>> {
        return Pager(
            config = pagingConfig,
            pagingSourceFactory = { NewsPagingSource(api, query) }
        ).flow
    }
}
```

---

## Step 3: ViewModel

```kotlin
@HiltViewModel
class NewsViewModel @Inject constructor(
    private val repository: NewsRepository
) : ViewModel() {

    private val _searchQuery = MutableStateFlow("android")

    val articles: Flow<PagingData<Article>> = _searchQuery
        .flatMapLatest { query ->
            repository.getTopHeadlines(query)
        }
        .cachedIn(viewModelScope)  // Cache pages in ViewModel scope — survives rotation

    fun setQuery(query: String) {
        _searchQuery.value = query
    }
}
```

`cachedIn(viewModelScope)` is critical — without it, rotating the device causes a full reload.

---

## Step 4: PagingDataAdapter

```kotlin
class ArticleAdapter : PagingDataAdapter<Article, ArticleAdapter.ArticleViewHolder>(
    ARTICLE_COMPARATOR
) {

    companion object {
        private val ARTICLE_COMPARATOR = object : DiffUtil.ItemCallback<Article>() {
            override fun areItemsTheSame(old: Article, new: Article) = old.url == new.url
            override fun areContentsTheSame(old: Article, new: Article) = old == new
        }
    }

    inner class ArticleViewHolder(
        private val binding: ItemArticleBinding
    ) : RecyclerView.ViewHolder(binding.root) {

        fun bind(article: Article?) {
            article ?: return
            binding.tvTitle.text = article.title
            binding.tvSource.text = article.source.name
            binding.tvPublishedAt.text = article.publishedAt
        }
    }

    override fun onCreateViewHolder(parent: ViewGroup, viewType: Int): ArticleViewHolder {
        val binding = ItemArticleBinding.inflate(
            LayoutInflater.from(parent.context), parent, false
        )
        return ArticleViewHolder(binding)
    }

    override fun onBindViewHolder(holder: ArticleViewHolder, position: Int) {
        holder.bind(getItem(position))  // getItem() — not positions[position]
    }
}
```

---

## Step 5: Collecting PagingData in Fragment

```kotlin
@AndroidEntryPoint
class NewsFragment : Fragment(R.layout.fragment_news) {

    private var _binding: FragmentNewsBinding? = null
    private val binding get() = _binding!!
    private val viewModel: NewsViewModel by viewModels()
    private lateinit var adapter: ArticleAdapter

    override fun onViewCreated(view: View, savedInstanceState: Bundle?) {
        super.onViewCreated(view, savedInstanceState)
        _binding = FragmentNewsBinding.bind(view)

        setupRecyclerView()
        observeArticles()
    }

    private fun setupRecyclerView() {
        adapter = ArticleAdapter()

        // Add loading state footer
        binding.rvArticles.adapter = adapter.withLoadStateFooter(
            footer = LoadingStateAdapter { adapter.retry() }
        )

        binding.rvArticles.layoutManager = LinearLayoutManager(requireContext())

        // Handle loading states
        adapter.addLoadStateListener { loadState ->
            binding.progressBar.isVisible =
                loadState.refresh is LoadState.Loading

            binding.tvError.isVisible =
                loadState.refresh is LoadState.Error

            if (loadState.refresh is LoadState.Error) {
                val error = (loadState.refresh as LoadState.Error).error
                binding.tvError.text = error.localizedMessage
            }

            binding.rvArticles.isVisible =
                loadState.refresh is LoadState.NotLoading
        }

        // Pull to refresh
        binding.swipeRefresh.setOnRefreshListener {
            adapter.refresh()
            binding.swipeRefresh.isRefreshing = false
        }
    }

    private fun observeArticles() {
        viewLifecycleOwner.lifecycleScope.launch {
            viewLifecycleOwner.repeatOnLifecycle(Lifecycle.State.STARTED) {
                viewModel.articles.collectLatest { pagingData ->
                    adapter.submitData(pagingData)
                }
            }
        }
    }

    override fun onDestroyView() {
        super.onDestroyView()
        _binding = null
    }
}
```

---

## Loading State Footer Adapter

```kotlin
class LoadingStateAdapter(
    private val retry: () -> Unit
) : LoadStateAdapter<LoadingStateAdapter.LoadingViewHolder>() {

    inner class LoadingViewHolder(
        private val binding: ItemLoadStateBinding
    ) : RecyclerView.ViewHolder(binding.root) {

        fun bind(loadState: LoadState) {
            binding.progressBar.isVisible = loadState is LoadState.Loading
            binding.btnRetry.isVisible = loadState is LoadState.Error
            binding.tvError.isVisible = loadState is LoadState.Error

            if (loadState is LoadState.Error) {
                binding.tvError.text = loadState.error.localizedMessage
            }

            binding.btnRetry.setOnClickListener { retry() }
        }
    }

    override fun onCreateViewHolder(parent: ViewGroup, loadState: LoadState): LoadingViewHolder {
        val binding = ItemLoadStateBinding.inflate(
            LayoutInflater.from(parent.context), parent, false
        )
        return LoadingViewHolder(binding)
    }

    override fun onBindViewHolder(holder: LoadingViewHolder, loadState: LoadState) {
        holder.bind(loadState)
    }
}
```

---

## RemoteMediator: Network + Room Cache

For offline-first paging (load from network, cache in Room, show from Room):

```kotlin
@OptIn(ExperimentalPagingApi::class)
class ArticleRemoteMediator(
    private val query: String,
    private val api: NewsApiService,
    private val database: NewsDatabase
) : RemoteMediator<Int, ArticleEntity>() {

    override suspend fun load(
        loadType: LoadType,
        state: PagingState<Int, ArticleEntity>
    ): MediatorResult {

        val page = when (loadType) {
            LoadType.REFRESH -> 1
            LoadType.PREPEND -> return MediatorResult.Success(endOfPaginationReached = true)
            LoadType.APPEND -> {
                val lastItem = state.lastItemOrNull()
                    ?: return MediatorResult.Success(endOfPaginationReached = true)
                // Get the next page from remote keys
                database.remoteKeyDao().getRemoteKey(lastItem.id)?.nextPage
                    ?: return MediatorResult.Success(endOfPaginationReached = true)
            }
        }

        return try {
            val response = api.getArticles(query, page, state.config.pageSize)

            database.withTransaction {
                if (loadType == LoadType.REFRESH) {
                    database.articleDao().clearAll()
                    database.remoteKeyDao().clearAll()
                }
                val keys = response.articles.map { article ->
                    RemoteKey(id = article.url, nextPage = page + 1)
                }
                database.remoteKeyDao().insertAll(keys)
                database.articleDao().insertAll(response.articles.map { it.toEntity() })
            }

            MediatorResult.Success(endOfPaginationReached = response.articles.isEmpty())
        } catch (e: IOException) {
            MediatorResult.Error(e)
        }
    }
}

// In Repository
fun getArticlesPaged(query: String): Flow<PagingData<Article>> {
    return Pager(
        config = PagingConfig(pageSize = 20),
        remoteMediator = ArticleRemoteMediator(query, api, database),
        pagingSourceFactory = { database.articleDao().getArticlesPaged() }
    ).flow
}
```

---

## Paging in Compose

```kotlin
@Composable
fun NewsScreen(viewModel: NewsViewModel = viewModel()) {
    val articles = viewModel.articles.collectAsLazyPagingItems()

    LazyColumn {
        items(
            count = articles.itemCount,
            key = articles.itemKey { it.url }
        ) { index ->
            val article = articles[index]
            if (article != null) {
                ArticleItem(article = article)
            } else {
                ArticlePlaceholder()
            }
        }

        when (val state = articles.loadState.append) {
            is LoadState.Loading -> item { CircularProgressIndicator() }
            is LoadState.Error -> item {
                RetryButton(onClick = { articles.retry() })
            }
            else -> Unit
        }
    }
}
```

---

## Common Mistakes

### Mistake 1: Not calling `cachedIn(viewModelScope)`

Without `cachedIn`, the paged data is not cached — screen rotation causes a full reload from page 1.

### Mistake 2: Using `getItemCount()` or `positions` directly

```kotlin
// WRONG — position-based access is unreliable in PagingDataAdapter
val item = items[position]

// CORRECT
val item = getItem(position)
```

### Mistake 3: Creating a new Pager on every collection

The `Pager` should be created once in the Repository/ViewModel, not on each `collect`.

---

## Interview Questions

**Q1: What is `PagingSource` and what does it do?**

> `PagingSource<Key, Value>` defines how to load a single page of data. The `load()` function receives a `LoadParams` with the page key and size, and returns a `LoadResult.Page` with the data and adjacent page keys.

**Q2: What is `cachedIn` and why is it important?**

> `cachedIn(viewModelScope)` caches the `PagingData` in the ViewModel's scope. Without it, every new collector (e.g., after rotation) triggers a fresh load from page 1. With it, the paged data survives configuration changes.

**Q3: What is `RemoteMediator` and when would you use it?**

> `RemoteMediator` bridges a network API and a local Room database. The UI always reads from Room, while `RemoteMediator` loads fresh data from the network and writes it to Room. Use it for offline-first apps where you want to cache paginated network data.

---

## Summary

- Paging 3 loads large datasets in pages to avoid loading everything at once
- `PagingSource` defines how to load one page; `Pager` orchestrates loading
- `PagingDataAdapter` handles the paginated list; use `getItem()` not direct positions
- Always call `cachedIn(viewModelScope)` in the ViewModel
- Use `RemoteMediator` for offline-first paging with Room cache
- Handle `LoadState` to show loading/error/empty states

**Next:** [Mini Project — News App](./mini-project-news-app.md)
