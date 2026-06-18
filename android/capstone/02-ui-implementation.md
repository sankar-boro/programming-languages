# Capstone Part 2: UI Implementation

## Theme and Styling

`core/ui/src/main/res/values/themes.xml`:

```xml
<resources>
    <style name="Theme.FeedFlow" parent="Theme.Material3.DayNight.NoActionBar">
        <item name="colorPrimary">@color/md_theme_primary</item>
        <item name="colorOnPrimary">@color/md_theme_on_primary</item>
        <item name="colorPrimaryContainer">@color/md_theme_primary_container</item>
        <item name="colorSecondary">@color/md_theme_secondary</item>
        <item name="colorSurface">@color/md_theme_surface</item>
        <item name="colorBackground">@color/md_theme_background</item>
        <item name="colorError">@color/md_theme_error</item>
        <item name="android:fontFamily">@font/inter</item>
        <item name="fontFamily">@font/inter</item>
        <item name="textInputStyle">@style/Widget.Material3.TextInputLayout.OutlinedBox</item>
    </style>
</resources>
```

`core/ui/src/main/res/values/colors.xml`:

```xml
<resources>
    <color name="md_theme_primary">#0061A4</color>
    <color name="md_theme_on_primary">#FFFFFF</color>
    <color name="md_theme_primary_container">#D1E4FF</color>
    <color name="md_theme_secondary">#535F70</color>
    <color name="md_theme_surface">#FDFCFF</color>
    <color name="md_theme_background">#FDFCFF</color>
    <color name="md_theme_error">#BA1A1A</color>
</resources>
```

---

## MainActivity Layout

`app/src/main/res/layout/activity_main.xml`:

```xml
<?xml version="1.0" encoding="utf-8"?>
<androidx.coordinatorlayout.widget.CoordinatorLayout
    xmlns:android="http://schemas.android.com/apk/res/android"
    xmlns:app="http://schemas.android.com/apk/res-auto"
    android:layout_width="match_parent"
    android:layout_height="match_parent">

    <androidx.fragment.app.FragmentContainerView
        android:id="@+id/navHostFragment"
        android:name="androidx.navigation.fragment.NavHostFragment"
        android:layout_width="match_parent"
        android:layout_height="match_parent"
        app:defaultNavHost="true"
        app:navGraph="@navigation/main_nav_graph"
        app:layout_behavior="@string/appbar_scrolling_view_behavior" />

    <com.google.android.material.bottomnavigation.BottomNavigationView
        android:id="@+id/bottomNav"
        android:layout_width="match_parent"
        android:layout_height="wrap_content"
        android:layout_gravity="bottom"
        app:menu="@menu/bottom_nav_menu"
        app:labelVisibilityMode="labeled" />

</androidx.coordinatorlayout.widget.CoordinatorLayout>
```

---

## Article List Item Layout

`feature/headlines/src/main/res/layout/item_article.xml`:

```xml
<?xml version="1.0" encoding="utf-8"?>
<com.google.android.material.card.MaterialCardView
    xmlns:android="http://schemas.android.com/apk/res/android"
    xmlns:app="http://schemas.android.com/apk/res-auto"
    android:layout_width="match_parent"
    android:layout_height="wrap_content"
    app:cardCornerRadius="12dp"
    app:cardElevation="0dp"
    app:strokeWidth="1dp"
    app:strokeColor="?attr/colorOutline"
    android:layout_marginHorizontal="16dp"
    android:layout_marginBottom="8dp"
    android:clickable="true"
    android:focusable="true">

    <androidx.constraintlayout.widget.ConstraintLayout
        android:layout_width="match_parent"
        android:layout_height="wrap_content"
        android:padding="12dp">

        <ImageView
            android:id="@+id/ivThumbnail"
            android:layout_width="80dp"
            android:layout_height="80dp"
            android:scaleType="centerCrop"
            android:src="@drawable/placeholder_news"
            app:layout_constraintTop_toTopOf="parent"
            app:layout_constraintStart_toStartOf="parent"
            android:background="@drawable/bg_image_placeholder" />

        <TextView
            android:id="@+id/tvSource"
            android:layout_width="0dp"
            android:layout_height="wrap_content"
            android:textSize="12sp"
            android:textColor="?attr/colorPrimary"
            android:textStyle="bold"
            app:layout_constraintTop_toTopOf="@id/ivThumbnail"
            app:layout_constraintStart_toEndOf="@id/ivThumbnail"
            app:layout_constraintEnd_toEndOf="parent"
            android:layout_marginStart="12dp" />

        <TextView
            android:id="@+id/tvTitle"
            android:layout_width="0dp"
            android:layout_height="wrap_content"
            android:textSize="14sp"
            android:textStyle="bold"
            android:maxLines="3"
            android:ellipsize="end"
            app:layout_constraintTop_toBottomOf="@id/tvSource"
            app:layout_constraintStart_toStartOf="@id/tvSource"
            app:layout_constraintEnd_toEndOf="parent"
            android:layout_marginTop="4dp" />

        <TextView
            android:id="@+id/tvPublishedAt"
            android:layout_width="0dp"
            android:layout_height="wrap_content"
            android:textSize="12sp"
            android:textColor="?attr/colorOnSurfaceVariant"
            app:layout_constraintTop_toBottomOf="@id/tvTitle"
            app:layout_constraintStart_toStartOf="@id/tvSource"
            app:layout_constraintEnd_toEndOf="parent"
            app:layout_constraintBottom_toBottomOf="@id/ivThumbnail"
            android:layout_marginTop="4dp" />

        <ImageButton
            android:id="@+id/btnBookmark"
            android:layout_width="wrap_content"
            android:layout_height="wrap_content"
            android:background="?attr/selectableItemBackgroundBorderless"
            android:padding="8dp"
            android:src="@drawable/ic_bookmark_border"
            app:layout_constraintBottom_toBottomOf="parent"
            app:layout_constraintEnd_toEndOf="parent"
            app:layout_constraintTop_toBottomOf="@id/ivThumbnail" />

    </androidx.constraintlayout.widget.ConstraintLayout>

</com.google.android.material.card.MaterialCardView>
```

---

## ArticleAdapter (PagingDataAdapter)

```kotlin
class ArticleAdapter(
    private val onArticleClick: (Article) -> Unit,
    private val onBookmarkClick: (Article) -> Unit
) : PagingDataAdapter<Article, ArticleAdapter.ArticleViewHolder>(DIFF_CALLBACK) {

    companion object {
        val DIFF_CALLBACK = object : DiffUtil.ItemCallback<Article>() {
            override fun areItemsTheSame(old: Article, new: Article) = old.id == new.id
            override fun areContentsTheSame(old: Article, new: Article) = old == new
        }
    }

    inner class ArticleViewHolder(
        private val binding: ItemArticleBinding
    ) : RecyclerView.ViewHolder(binding.root) {

        fun bind(article: Article) {
            binding.tvTitle.text = article.title
            binding.tvSource.text = article.sourceName
            binding.tvPublishedAt.text = article.publishedAt.formatRelative()

            article.imageUrl?.let { url ->
                binding.ivThumbnail.load(url) {
                    placeholder(R.drawable.placeholder_news)
                    error(R.drawable.placeholder_news)
                    crossfade(true)
                }
            }

            binding.btnBookmark.setImageResource(
                if (article.isBookmarked) R.drawable.ic_bookmark_filled
                else R.drawable.ic_bookmark_border
            )

            binding.root.setOnClickListener { onArticleClick(article) }
            binding.btnBookmark.setOnClickListener { onBookmarkClick(article) }
        }
    }

    override fun onCreateViewHolder(parent: ViewGroup, viewType: Int): ArticleViewHolder {
        val binding = ItemArticleBinding.inflate(
            LayoutInflater.from(parent.context), parent, false
        )
        return ArticleViewHolder(binding)
    }

    override fun onBindViewHolder(holder: ArticleViewHolder, position: Int) {
        getItem(position)?.let { holder.bind(it) }
    }
}
```

---

## Headlines Fragment

```kotlin
@AndroidEntryPoint
class HeadlinesFragment : Fragment(R.layout.fragment_headlines) {

    private var _binding: FragmentHeadlinesBinding? = null
    private val binding get() = _binding!!

    private val viewModel: HeadlinesViewModel by viewModels()
    private lateinit var adapter: ArticleAdapter

    override fun onViewCreated(view: View, savedInstanceState: Bundle?) {
        super.onViewCreated(view, savedInstanceState)
        _binding = FragmentHeadlinesBinding.bind(view)

        setupToolbar()
        setupCategoryChips()
        setupRecyclerView()
        observeArticles()
        observeEvents()
    }

    private fun setupToolbar() {
        binding.toolbar.title = "FeedFlow"
    }

    private fun setupCategoryChips() {
        NewsCategory.entries.forEach { category ->
            val chip = Chip(requireContext()).apply {
                text = category.displayName
                isCheckable = true
                style = "@style/Widget.Material3.Chip.Filter"
            }
            binding.chipGroupCategories.addView(chip)
            chip.setOnCheckedChangeListener { _, isChecked ->
                if (isChecked) viewModel.setCategory(category)
            }
        }
    }

    private fun setupRecyclerView() {
        adapter = ArticleAdapter(
            onArticleClick = { article ->
                viewModel.onArticleClicked(article)
            },
            onBookmarkClick = { article ->
                viewModel.toggleBookmark(article)
            }
        )

        binding.rvArticles.adapter = adapter.withLoadStateFooter(
            footer = LoadingStateAdapter { adapter.retry() }
        )
        binding.rvArticles.layoutManager = LinearLayoutManager(requireContext())
        binding.rvArticles.setHasFixedSize(true)

        adapter.addLoadStateListener { loadState ->
            binding.progressBar.isVisible = loadState.refresh is LoadState.Loading
            binding.layoutError.isVisible = loadState.refresh is LoadState.Error
            binding.rvArticles.isVisible = loadState.refresh is LoadState.NotLoading
        }

        binding.btnRetry.setOnClickListener { adapter.retry() }

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

    private fun observeEvents() {
        viewLifecycleOwner.lifecycleScope.launch {
            viewLifecycleOwner.repeatOnLifecycle(Lifecycle.State.STARTED) {
                viewModel.events.collect { event ->
                    when (event) {
                        is HeadlinesEvent.NavigateToDetail -> {
                            val action = HeadlinesFragmentDirections
                                .actionHeadlinesToDetail(articleId = event.articleId)
                            findNavController().navigate(action)
                        }
                        is HeadlinesEvent.ShowMessage ->
                            Snackbar.make(binding.root, event.message, Snackbar.LENGTH_SHORT)
                                .setAnchorView(binding.root)
                                .show()
                    }
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

**Next:** [Part 3 — API Integration](./03-api-integration.md)
