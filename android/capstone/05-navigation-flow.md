# Capstone Part 5: Navigation Flow

## Navigation Graph

`app/src/main/res/navigation/main_nav_graph.xml`:

```xml
<?xml version="1.0" encoding="utf-8"?>
<navigation
    xmlns:android="http://schemas.android.com/apk/res/android"
    xmlns:app="http://schemas.android.com/apk/res-auto"
    android:id="@+id/main_nav_graph"
    app:startDestination="@id/headlines_nav_graph">

    <!-- Include feature nav graphs -->
    <include app:graph="@navigation/headlines_nav_graph" />
    <include app:graph="@navigation/search_nav_graph" />
    <include app:graph="@navigation/bookmarks_nav_graph" />
    <include app:graph="@navigation/settings_nav_graph" />

    <!-- Article detail — accessible from all features -->
    <fragment
        android:id="@+id/articleDetailFragment"
        android:name="com.yourname.feedflow.feature.detail.ArticleDetailFragment"
        android:label="Article">
        <argument
            android:name="articleId"
            app:argType="string" />
    </fragment>

    <!-- Global actions (accessible from any destination) -->
    <action
        android:id="@+id/action_global_to_detail"
        app:destination="@id/articleDetailFragment"
        app:enterAnim="@anim/slide_in_right"
        app:exitAnim="@anim/slide_out_left"
        app:popEnterAnim="@anim/slide_in_left"
        app:popExitAnim="@anim/slide_out_right" />

</navigation>
```

`feature/headlines/src/main/res/navigation/headlines_nav_graph.xml`:

```xml
<?xml version="1.0" encoding="utf-8"?>
<navigation
    xmlns:android="http://schemas.android.com/apk/res/android"
    xmlns:app="http://schemas.android.com/apk/res-auto"
    android:id="@+id/headlines_nav_graph"
    app:startDestination="@id/headlinesFragment">

    <fragment
        android:id="@+id/headlinesFragment"
        android:name="com.yourname.feedflow.feature.headlines.HeadlinesFragment"
        android:label="Headlines">
    </fragment>

</navigation>
```

---

## Bottom Navigation Menu

`app/src/main/res/menu/bottom_nav_menu.xml`:

```xml
<?xml version="1.0" encoding="utf-8"?>
<menu xmlns:android="http://schemas.android.com/apk/res/android">
    <item
        android:id="@+id/headlines_nav_graph"
        android:icon="@drawable/ic_home"
        android:title="Headlines" />
    <item
        android:id="@+id/search_nav_graph"
        android:icon="@drawable/ic_search"
        android:title="Search" />
    <item
        android:id="@+id/bookmarks_nav_graph"
        android:icon="@drawable/ic_bookmark"
        android:title="Bookmarks" />
    <item
        android:id="@+id/settings_nav_graph"
        android:icon="@drawable/ic_settings"
        android:title="Settings" />
</menu>
```

---

## MainActivity Navigation Setup

```kotlin
@AndroidEntryPoint
class MainActivity : AppCompatActivity() {

    private lateinit var binding: ActivityMainBinding
    private lateinit var navController: NavController

    override fun onCreate(savedInstanceState: Bundle?) {
        super.onCreate(savedInstanceState)
        binding = ActivityMainBinding.inflate(layoutInflater)
        setContentView(binding.root)

        setupNavigation()
    }

    private fun setupNavigation() {
        val navHostFragment = supportFragmentManager
            .findFragmentById(R.id.navHostFragment) as NavHostFragment
        navController = navHostFragment.navController

        // Wire bottom navigation to nav controller
        binding.bottomNav.setupWithNavController(navController)

        // Hide bottom nav on detail screen
        navController.addOnDestinationChangedListener { _, destination, _ ->
            binding.bottomNav.isVisible = when (destination.id) {
                R.id.articleDetailFragment -> false
                else -> true
            }
        }
    }

    override fun onSupportNavigateUp(): Boolean {
        return navController.navigateUp() || super.onSupportNavigateUp()
    }
}
```

---

## Navigate to Detail from Any Feature

The detail destination is global — use a global action:

```kotlin
// In HeadlinesFragment, SearchFragment, BookmarksFragment — same code works
viewModel.events.collect { event ->
    when (event) {
        is HeadlinesEvent.NavigateToDetail -> {
            // Global action — available from any destination
            val action = MainNavGraphDirections.actionGlobalToDetail(
                articleId = event.articleId
            )
            findNavController().navigate(action)
        }
    }
}
```

---

## Article Detail Fragment

```kotlin
@AndroidEntryPoint
class ArticleDetailFragment : Fragment(R.layout.fragment_article_detail) {

    private var _binding: FragmentArticleDetailBinding? = null
    private val binding get() = _binding!!

    private val args: ArticleDetailFragmentArgs by navArgs()
    private val viewModel: ArticleDetailViewModel by viewModels()

    override fun onViewCreated(view: View, savedInstanceState: Bundle?) {
        super.onViewCreated(view, savedInstanceState)
        _binding = FragmentArticleDetailBinding.bind(view)

        setupToolbar()
        viewModel.loadArticle(args.articleId)
        observeState()
    }

    private fun setupToolbar() {
        binding.toolbar.setNavigationOnClickListener {
            findNavController().popBackStack()
        }
    }

    private fun observeState() {
        viewLifecycleOwner.lifecycleScope.launch {
            viewLifecycleOwner.repeatOnLifecycle(Lifecycle.State.STARTED) {
                viewModel.uiState.collect { state ->
                    state.article?.let { article ->
                        binding.tvTitle.text = article.title
                        binding.tvSource.text = article.sourceName
                        binding.tvPublishedAt.text = article.publishedAt
                        binding.tvContent.text = article.content ?: article.description

                        article.imageUrl?.let { url ->
                            binding.ivHero.load(url) {
                                crossfade(true)
                                placeholder(R.drawable.placeholder_news)
                            }
                        }

                        binding.btnBookmark.setIconResource(
                            if (article.isBookmarked) R.drawable.ic_bookmark_filled
                            else R.drawable.ic_bookmark_border
                        )
                    }

                    binding.progressBar.isVisible = state.isLoading
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

## Deep Link Support

Declare deep links in the nav graph for sharing articles:

```xml
<fragment android:id="@+id/articleDetailFragment" ...>
    <deepLink
        app:uri="feedflow://article/{articleId}"
        app:action="android.intent.action.VIEW" />
</fragment>
```

`AndroidManifest.xml`:

```xml
<activity android:name=".MainActivity">
    <intent-filter>
        <action android:name="android.intent.action.MAIN" />
        <category android:name="android.intent.category.LAUNCHER" />
    </intent-filter>
    <nav-graph android:value="@navigation/main_nav_graph" />
</activity>
```

---

**Next:** [Part 6 — Testing](./06-testing.md)
