# Chapter 4: Performance Optimization

## The Performance Trinity

1. **Rendering** — Does the UI draw at 60fps (16ms per frame)?
2. **Memory** — Does the app use memory responsibly and release it?
3. **Threading** — Is work happening on the right thread?

---

## 1. Rendering Performance

### Hierarchy Viewer and Layout Inspector

Use **Android Studio → Layout Inspector** to inspect view hierarchies at runtime. Deep nesting = slow rendering.

**Goal: Keep hierarchy depth < 5 levels**

```xml
<!-- SLOW — 4 nested LinearLayouts -->
<LinearLayout>
  <LinearLayout>
    <LinearLayout>
      <LinearLayout>
        <TextView />
      </LinearLayout>
    </LinearLayout>
  </LinearLayout>
</LinearLayout>

<!-- FAST — flat ConstraintLayout -->
<ConstraintLayout>
  <TextView />
</ConstraintLayout>
```

### `setHasFixedSize(true)` in RecyclerView

```kotlin
// If item count doesn't affect RecyclerView's size
binding.rvNotes.setHasFixedSize(true)
```

### Use `DiffUtil` in RecyclerView Adapters

Never call `notifyDataSetChanged()`. Use `ListAdapter` which runs `DiffUtil` on a background thread.

### Avoid Over-drawing

Use **GPU Overdraw** (Developer Options → Debug GPU overdraw) to find views painted multiple times.

```xml
<!-- WRONG — opaque background over opaque background = overdraw -->
<LinearLayout android:background="@color/white">
  <View android:background="@color/white" />
</LinearLayout>

<!-- CORRECT — remove redundant backgrounds -->
```

---

## 2. Memory Leaks and Management

### LeakCanary

Add LeakCanary to detect memory leaks during development:

```kotlin
debugImplementation("com.squareup.leakcanary:leakcanary-android:2.14")
```

LeakCanary automatically detects leaked Activities, Fragments, ViewModels, and more — no code changes needed.

### Common Leak Sources

#### Leak 1: Static reference to Context

```kotlin
// LEAK — static field holds Activity
companion object {
    var instance: Context? = null
}

// SAFE — use ApplicationContext for long-lived references
companion object {
    var appContext: Context? = null
}
// In Application.onCreate:
appContext = applicationContext
```

#### Leak 2: Non-static inner class in static context

```kotlin
// LEAK — inner Runnable holds reference to outer Activity
Handler().postDelayed(object : Runnable {
    override fun run() {
        // 'this@MainActivity' is held after Activity is destroyed
        binding.tvTimer.text = "Done"
    }
}, 5000)

// SAFE — use WeakReference or coroutines
lifecycleScope.launch {
    delay(5000)
    binding.tvTimer.text = "Done"  // lifecycleScope is cancelled on destroy
}
```

#### Leak 3: Fragment binding reference

```kotlin
// LEAK
class MyFragment : Fragment() {
    private var binding: FragmentMyBinding? = null

    override fun onDestroyView() {
        super.onDestroyView()
        // Missing: binding = null  ← LEAK
    }
}

// SAFE
override fun onDestroyView() {
    super.onDestroyView()
    _binding = null  // ← Required
}
```

#### Leak 4: Unregistered callbacks/listeners

```kotlin
// LEAK — registered but never unregistered
connectivityManager.registerNetworkCallback(request, callback)

// SAFE — unregister in lifecycle
override fun onStop() {
    super.onStop()
    connectivityManager.unregisterNetworkCallback(callback)
}

// BETTER — use DefaultLifecycleObserver or callbackFlow that handles cleanup
```

---

## 3. Threading

### Strict Mode (Development Only)

```kotlin
if (BuildConfig.DEBUG) {
    StrictMode.setThreadPolicy(
        StrictMode.ThreadPolicy.Builder()
            .detectDiskReads()
            .detectDiskWrites()
            .detectNetwork()
            .penaltyLog()
            .penaltyDeath()  // Crash the app when violation detected
            .build()
    )
}
```

### Never Block the Main Thread

```kotlin
// WRONG — disk read on main thread (StrictMode will catch this)
val file = File(filesDir, "data.json").readText()

// CORRECT — move to IO dispatcher
viewModelScope.launch(Dispatchers.IO) {
    val file = File(filesDir, "data.json").readText()
    withContext(Dispatchers.Main) { updateUI(file) }
}
```

---

## 4. App Startup Time

### App Startup Library

```kotlin
implementation("androidx.startup:startup-runtime:1.1.1")
```

Initialize components lazily:

```kotlin
class TimberInitializer : Initializer<Unit> {
    override fun create(context: Context) {
        if (BuildConfig.DEBUG) Timber.plant(Timber.DebugTree())
    }
    override fun dependencies(): List<Class<out Initializer<*>>> = emptyList()
}
```

`AndroidManifest.xml`:

```xml
<provider
    android:name="androidx.startup.InitializationProvider"
    android:authorities="${applicationId}.androidx-startup"
    android:exported="false"
    tools:node="merge">
    <meta-data
        android:name="com.yourname.TimberInitializer"
        android:value="androidx.startup" />
</provider>
```

### Baseline Profiles

Baseline Profiles pre-compile critical code paths, improving startup time by up to 40%:

```kotlin
// Add to app/src/main/baseline-prof.txt
Lcom/yourname/app/MainActivity;
Lcom/yourname/feature/notes/NoteListFragment;
// ... critical classes
```

```kotlin
implementation("androidx.profileinstaller:profileinstaller:1.3.1")
```

---

## 5. Network Performance

### OkHttp Caching

```kotlin
val cacheDir = File(context.cacheDir, "http_cache")
val cache = Cache(cacheDir, 10 * 1024 * 1024)  // 10 MB

val okHttpClient = OkHttpClient.Builder()
    .cache(cache)
    .addInterceptor { chain ->
        val request = if (connectivityObserver.isOnline()) {
            chain.request().newBuilder()
                .header("Cache-Control", "public, max-age=60")  // Cache for 60s
                .build()
        } else {
            chain.request().newBuilder()
                .header("Cache-Control", "public, only-if-cached, max-stale=604800")  // Use cache up to 1 week
                .build()
        }
        chain.proceed(request)
    }
    .build()
```

### Image Loading Efficiency

```kotlin
// Use Coil or Glide — handles caching, sampling, and memory automatically
implementation("io.coil-kt:coil:2.7.0")

// In Compose
AsyncImage(
    model = ImageRequest.Builder(context)
        .data(article.imageUrl)
        .crossfade(true)
        .memoryCachePolicy(CachePolicy.ENABLED)
        .diskCachePolicy(CachePolicy.ENABLED)
        .build(),
    contentDescription = null,
    modifier = Modifier.fillMaxWidth().height(200.dp),
    contentScale = ContentScale.Crop
)

// In XML
imageView.load(article.imageUrl) {
    crossfade(true)
    placeholder(R.drawable.placeholder)
    error(R.drawable.error_image)
    transformations(RoundedCornersTransformation(8f))
}
```

---

## 6. APK Size Optimization

### Enable R8 (mandatory for release)

```kotlin
buildTypes {
    release {
        isMinifyEnabled = true
        isShrinkResources = true
        proguardFiles(
            getDefaultProguardFile("proguard-android-optimize.txt"),
            "proguard-rules.pro"
        )
    }
}
```

### Use Android App Bundles (AAB)

AAB lets Google Play deliver only the resources needed for each device — smaller download size.

```
Build → Build Bundle(s)/APK(s) → Build Bundle(s)
```

### Reduce Image Sizes

- Use **WebP** instead of PNG/JPEG
- Use **vector drawables** for icons
- Use adaptive icons

```xml
<!-- res/mipmap/ic_launcher.xml -->
<adaptive-icon>
    <background android:drawable="@color/ic_launcher_background" />
    <foreground android:drawable="@drawable/ic_launcher_foreground" />
</adaptive-icon>
```

---

## 7. Profiling Tools

| Tool | What It Measures |
|------|-----------------|
| Android Profiler → CPU | Method traces, thread activity |
| Android Profiler → Memory | Heap dumps, allocations |
| Android Profiler → Network | Request timing, size |
| Layout Inspector | View hierarchy, measured dimensions |
| GPU Overdraw (Dev options) | Over-drawn pixels |
| App Startup | Tracing app init time |

---

## Common Performance Anti-Patterns

```kotlin
// 1. Creating objects in onDraw or onBindViewHolder
class MyAdapter : RecyclerView.Adapter<...>() {
    override fun onBindViewHolder(holder, position) {
        // WRONG — creates new Paint on every bind
        val paint = Paint()
        holder.binding.view.paint = paint
    }
}

// 2. Blocking IO on main thread
override fun onCreate(...) {
    val prefs = getSharedPreferences("config", MODE_PRIVATE)
    // SharedPreferences first access can block main thread
}

// 3. Large Bitmap without sampling
val bitmap = BitmapFactory.decodeFile(path)  // Loads full resolution
// Use Coil/Glide instead

// 4. Not using setHasFixedSize
binding.rvList.adapter = adapter  // Missing setHasFixedSize(true)
```

---

## Interview Questions

**Q1: What causes "dropped frames" in Android?**

> Each frame must complete in 16ms (60fps). When the main thread takes longer — due to deep view hierarchies, unnecessary overdraw, I/O on the main thread, or expensive operations in `onBindViewHolder` — frames are dropped, causing jank (visual stuttering).

**Q2: What is a memory leak and how does LeakCanary help?**

> A memory leak is when an object that should be garbage collected is still referenced, preventing GC from reclaiming the memory. LeakCanary monitors Activity/Fragment destruction and uses heap analysis to detect if any references to destroyed components are held. It reports the reference chain.

**Q3: What is R8 and why enable it for release builds?**

> R8 is Google's code shrinker and obfuscator. It removes unused code (shrinking), renames classes/methods to short names (obfuscating), and applies optimizations — all reducing APK size. Enabling it for release builds is a standard best practice.

---

## Summary

- Keep view hierarchies flat — use ConstraintLayout, avoid deep nesting
- Always clear Fragment binding in `onDestroyView()` to prevent leaks
- Use LeakCanary in debug builds to detect memory leaks early
- Never do I/O on the main thread — use `Dispatchers.IO`
- Enable R8 (`isMinifyEnabled = true`, `isShrinkResources = true`) for all release builds
- Use Baseline Profiles to improve cold start time

**Next:** [Chapter 5 — Multi-Module Apps](./05-multi-module-apps.md)
