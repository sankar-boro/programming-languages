# Mini Project: E-Commerce Feature Module

## Overview

Build a modular e-commerce **product catalog** feature using everything from Level 4.

**Features:**
- Product list with Paging 3
- Search and filter
- Product detail screen
- Add to cart (shared ViewModel)
- Offline-first (Room cache)
- Clean Architecture (domain/data/ui layers)
- Multi-module structure
- Unit + UI tests

---

## Module Structure

```
:feature:shop
  ├── di/
  │   └── ShopModule.kt
  ├── domain/
  │   ├── model/Product.kt
  │   ├── model/CartItem.kt
  │   ├── repository/ProductRepository.kt
  │   └── usecase/
  │       ├── GetProductsUseCase.kt
  │       ├── GetProductDetailUseCase.kt
  │       └── AddToCartUseCase.kt
  ├── data/
  │   ├── remote/ShopApiService.kt
  │   ├── local/ProductDao.kt
  │   ├── paging/ProductPagingSource.kt
  │   └── repository/ProductRepositoryImpl.kt
  └── ui/
      ├── list/
      │   ├── ProductListFragment.kt
      │   └── ProductListViewModel.kt
      ├── detail/
      │   ├── ProductDetailFragment.kt
      │   └── ProductDetailViewModel.kt
      └── cart/
          └── CartViewModel.kt
```

---

## Domain Layer

```kotlin
// domain/model/Product.kt
data class Product(
    val id: Long,
    val title: String,
    val description: String,
    val price: Double,
    val imageUrl: String,
    val category: String,
    val rating: Float,
    val reviewCount: Int,
    val inStock: Boolean
)

// domain/repository/ProductRepository.kt
interface ProductRepository {
    fun getProducts(category: String?): Flow<PagingData<Product>>
    suspend fun getProductById(id: Long): Product?
    fun searchProducts(query: String): Flow<PagingData<Product>>
}

// domain/usecase/GetProductsUseCase.kt
class GetProductsUseCase @Inject constructor(
    private val repository: ProductRepository
) {
    operator fun invoke(category: String? = null): Flow<PagingData<Product>> =
        repository.getProducts(category)
}

// domain/usecase/AddToCartUseCase.kt
class AddToCartUseCase @Inject constructor(
    private val cartRepository: CartRepository
) {
    suspend operator fun invoke(product: Product, quantity: Int = 1): Result<Unit> {
        if (!product.inStock) return Result.failure(Exception("${product.title} is out of stock"))
        return runCatching { cartRepository.addItem(CartItem(product = product, quantity = quantity)) }
    }
}
```

---

## ViewModel

```kotlin
data class ProductListState(
    val selectedCategory: String? = null,
    val searchQuery: String = "",
    val isSearchActive: Boolean = false
)

sealed class ProductListEvent {
    data class ShowMessage(val message: String) : ProductListEvent()
    data class NavigateToDetail(val productId: Long) : ProductListEvent()
}

@HiltViewModel
class ProductListViewModel @Inject constructor(
    private val getProductsUseCase: GetProductsUseCase,
    private val addToCartUseCase: AddToCartUseCase
) : ViewModel() {

    private val _state = MutableStateFlow(ProductListState())
    val state: StateFlow<ProductListState> = _state.asStateFlow()

    private val _events = MutableSharedFlow<ProductListEvent>()
    val events = _events.asSharedFlow()

    val products: Flow<PagingData<Product>> = _state
        .map { it.selectedCategory to it.searchQuery }
        .distinctUntilChanged()
        .flatMapLatest { (category, query) ->
            if (query.isNotBlank()) {
                // This would require repository.searchProducts(query)
                getProductsUseCase(category)
            } else {
                getProductsUseCase(category)
            }
        }
        .cachedIn(viewModelScope)

    fun setCategory(category: String?) {
        _state.update { it.copy(selectedCategory = category) }
    }

    fun onProductClicked(productId: Long) {
        viewModelScope.launch {
            _events.emit(ProductListEvent.NavigateToDetail(productId))
        }
    }

    fun addToCart(product: Product) {
        viewModelScope.launch {
            addToCartUseCase(product).fold(
                onSuccess = {
                    _events.emit(ProductListEvent.ShowMessage("${product.title} added to cart"))
                },
                onFailure = { error ->
                    _events.emit(ProductListEvent.ShowMessage(error.message ?: "Error"))
                }
            )
        }
    }
}
```

---

## Unit Tests

```kotlin
class ProductListViewModelTest {

    @get:Rule
    val mainDispatcherRule = MainDispatcherRule()

    private lateinit var viewModel: ProductListViewModel
    private val fakeRepository = FakeProductRepository()
    private val fakeCartRepository = FakeCartRepository()

    @Before
    fun setup() {
        viewModel = ProductListViewModel(
            GetProductsUseCase(fakeRepository),
            AddToCartUseCase(fakeCartRepository)
        )
    }

    @Test
    fun `addToCart with out-of-stock product emits error event`() = runTest {
        val outOfStock = Product(id = 1, title = "Sold Out", inStock = false,
            price = 10.0, description = "", imageUrl = "", category = "", rating = 0f, reviewCount = 0)

        val events = mutableListOf<ProductListEvent>()
        val job = launch { viewModel.events.collect { events.add(it) } }

        viewModel.addToCart(outOfStock)

        job.cancel()

        assertThat(events).hasSize(1)
        assertThat((events[0] as ProductListEvent.ShowMessage).message)
            .contains("out of stock")
    }

    @Test
    fun `setCategory updates state`() {
        viewModel.setCategory("Electronics")

        assertThat(viewModel.state.value.selectedCategory).isEqualTo("Electronics")
    }
}
```

---

## Level 4 Checkpoint

Before moving to the Capstone, confirm you can:

- [ ] Structure a project with domain/data/ui layers following Clean Architecture
- [ ] Implement MVVM and understand how MVI differs
- [ ] Create a multi-module Gradle project with feature and core modules
- [ ] Write convention plugins to share build config
- [ ] Implement offline-first with Room as single source of truth
- [ ] Profile the app with Android Profiler and fix a leak
- [ ] Set up GitHub Actions CI with lint + unit tests
- [ ] Use Timber for logging with Crashlytics in production
- [ ] Write a production-grade ProGuard rules file

**Final Step:** [Capstone Project](../capstone/overview.md)
