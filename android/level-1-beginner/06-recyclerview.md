# Chapter 6: RecyclerView — Deep Dive

## What Is RecyclerView?

`RecyclerView` is AndroidX's high-performance list component. It **recycles** (reuses) off-screen item views instead of creating new ones, making it efficient even for thousands of items.

**Analogy:** Imagine a news ticker. Instead of printing new paper for every headline, the ticker reuses the same physical strip — just printing new content on it. RecyclerView does the same with View objects.

```kotlin
implementation("androidx.recyclerview:recyclerview:1.3.2")
```

---

## RecyclerView Architecture

RecyclerView has four key parts:

```
┌───────────────────────────────────────────┐
│              RecyclerView                 │
│                                           │
│  ┌─────────────────────────────────────┐  │
│  │         LayoutManager               │  │  ← Decides HOW items are arranged
│  │  (Linear, Grid, StaggeredGrid)      │  │
│  └─────────────────────────────────────┘  │
│                                           │
│  ┌─────────────────────────────────────┐  │
│  │            Adapter                  │  │  ← Bridges data ↔ Views
│  │  (creates + binds ViewHolders)      │  │
│  └─────────────────────────────────────┘  │
│                                           │
│  ┌─────────────────────────────────────┐  │
│  │           ViewHolder                │  │  ← Holds references to one item's views
│  └─────────────────────────────────────┘  │
│                                           │
│  ┌─────────────────────────────────────┐  │
│  │         ItemAnimator                │  │  ← Animates add/remove/change
│  └─────────────────────────────────────┘  │
└───────────────────────────────────────────┘
```

---

## Step-by-Step: Building a RecyclerView

### Step 1: Define the data model

```kotlin
data class Contact(
    val id: Long,
    val name: String,
    val phone: String,
    val avatarUrl: String?
)
```

### Step 2: Create the item layout

`res/layout/item_contact.xml`:

```xml
<?xml version="1.0" encoding="utf-8"?>
<androidx.constraintlayout.widget.ConstraintLayout
    xmlns:android="http://schemas.android.com/apk/res/android"
    xmlns:app="http://schemas.android.com/apk/res-auto"
    android:layout_width="match_parent"
    android:layout_height="wrap_content"
    android:padding="16dp">

    <ImageView
        android:id="@+id/ivAvatar"
        android:layout_width="48dp"
        android:layout_height="48dp"
        android:src="@drawable/ic_person"
        app:layout_constraintStart_toStartOf="parent"
        app:layout_constraintTop_toTopOf="parent"
        app:layout_constraintBottom_toBottomOf="parent" />

    <TextView
        android:id="@+id/tvName"
        android:layout_width="0dp"
        android:layout_height="wrap_content"
        android:textSize="16sp"
        android:textStyle="bold"
        app:layout_constraintStart_toEndOf="@id/ivAvatar"
        app:layout_constraintEnd_toEndOf="parent"
        app:layout_constraintTop_toTopOf="@id/ivAvatar"
        android:layout_marginStart="12dp" />

    <TextView
        android:id="@+id/tvPhone"
        android:layout_width="0dp"
        android:layout_height="wrap_content"
        android:textSize="14sp"
        android:textColor="@color/gray"
        app:layout_constraintStart_toStartOf="@id/tvName"
        app:layout_constraintEnd_toEndOf="parent"
        app:layout_constraintTop_toBottomOf="@id/tvName"
        android:layout_marginTop="2dp" />

</androidx.constraintlayout.widget.ConstraintLayout>
```

### Step 3: Create the ViewHolder and Adapter

```kotlin
import androidx.recyclerview.widget.RecyclerView
import android.view.LayoutInflater
import android.view.ViewGroup

class ContactAdapter(
    private val contacts: List<Contact>,
    private val onContactClick: (Contact) -> Unit
) : RecyclerView.Adapter<ContactAdapter.ContactViewHolder>() {

    inner class ContactViewHolder(
        private val binding: ItemContactBinding
    ) : RecyclerView.ViewHolder(binding.root) {

        fun bind(contact: Contact) {
            binding.tvName.text = contact.name
            binding.tvPhone.text = contact.phone

            binding.root.setOnClickListener {
                onContactClick(contact)
            }
        }
    }

    override fun onCreateViewHolder(parent: ViewGroup, viewType: Int): ContactViewHolder {
        val binding = ItemContactBinding.inflate(
            LayoutInflater.from(parent.context),
            parent,
            false  // IMPORTANT: always false here — RecyclerView handles attachment
        )
        return ContactViewHolder(binding)
    }

    override fun onBindViewHolder(holder: ContactViewHolder, position: Int) {
        holder.bind(contacts[position])
    }

    override fun getItemCount(): Int = contacts.size
}
```

```java
// Java version
public class ContactAdapter extends RecyclerView.Adapter<ContactAdapter.ContactViewHolder> {

    private final List<Contact> contacts;
    private final OnContactClickListener listener;

    public interface OnContactClickListener {
        void onClick(Contact contact);
    }

    public ContactAdapter(List<Contact> contacts, OnContactClickListener listener) {
        this.contacts = contacts;
        this.listener = listener;
    }

    public static class ContactViewHolder extends RecyclerView.ViewHolder {
        private final ItemContactBinding binding;

        public ContactViewHolder(ItemContactBinding binding) {
            super(binding.getRoot());
            this.binding = binding;
        }

        public void bind(Contact contact, OnContactClickListener listener) {
            binding.tvName.setText(contact.getName());
            binding.tvPhone.setText(contact.getPhone());
            binding.getRoot().setOnClickListener(v -> listener.onClick(contact));
        }
    }

    @Override
    public ContactViewHolder onCreateViewHolder(ViewGroup parent, int viewType) {
        ItemContactBinding binding = ItemContactBinding.inflate(
            LayoutInflater.from(parent.getContext()), parent, false
        );
        return new ContactViewHolder(binding);
    }

    @Override
    public void onBindViewHolder(ContactViewHolder holder, int position) {
        holder.bind(contacts.get(position), listener);
    }

    @Override
    public int getItemCount() { return contacts.size(); }
}
```

### Step 4: Set up RecyclerView in the Activity/Fragment

`activity_main.xml`:

```xml
<androidx.recyclerview.widget.RecyclerView
    android:id="@+id/rvContacts"
    android:layout_width="0dp"
    android:layout_height="0dp"
    app:layout_constraintTop_toTopOf="parent"
    app:layout_constraintBottom_toBottomOf="parent"
    app:layout_constraintStart_toStartOf="parent"
    app:layout_constraintEnd_toEndOf="parent" />
```

`MainActivity.kt`:

```kotlin
class MainActivity : AppCompatActivity() {

    private lateinit var binding: ActivityMainBinding

    override fun onCreate(savedInstanceState: Bundle?) {
        super.onCreate(savedInstanceState)
        binding = ActivityMainBinding.inflate(layoutInflater)
        setContentView(binding.root)

        val contacts = listOf(
            Contact(1, "Alice Johnson", "+1 555-0101", null),
            Contact(2, "Bob Smith", "+1 555-0102", null),
            Contact(3, "Carol White", "+1 555-0103", null)
        )

        val adapter = ContactAdapter(contacts) { contact ->
            Toast.makeText(this, "Clicked: ${contact.name}", Toast.LENGTH_SHORT).show()
        }

        binding.rvContacts.apply {
            layoutManager = LinearLayoutManager(this@MainActivity)
            this.adapter = adapter
            // Optional: improve performance when item size doesn't change
            setHasFixedSize(true)
        }
    }
}
```

---

## DiffUtil — Efficient List Updates

Never call `notifyDataSetChanged()` — it redraws every visible item. Use `DiffUtil` to compute the minimal changes.

```kotlin
class ContactDiffCallback(
    private val oldList: List<Contact>,
    private val newList: List<Contact>
) : DiffUtil.Callback() {

    override fun getOldListSize() = oldList.size
    override fun getNewListSize() = newList.size

    override fun areItemsTheSame(oldPos: Int, newPos: Int): Boolean {
        return oldList[oldPos].id == newList[newPos].id
    }

    override fun areContentsTheSame(oldPos: Int, newPos: Int): Boolean {
        return oldList[oldPos] == newList[newPos]
    }
}
```

---

## ListAdapter — DiffUtil Built-In (Preferred)

`ListAdapter` handles DiffUtil automatically using `AsyncListDiffer`:

```kotlin
class ContactAdapter(
    private val onContactClick: (Contact) -> Unit
) : ListAdapter<Contact, ContactAdapter.ContactViewHolder>(DIFF_CALLBACK) {

    companion object {
        private val DIFF_CALLBACK = object : DiffUtil.ItemCallback<Contact>() {
            override fun areItemsTheSame(old: Contact, new: Contact): Boolean {
                return old.id == new.id
            }
            override fun areContentsTheSame(old: Contact, new: Contact): Boolean {
                return old == new
            }
        }
    }

    inner class ContactViewHolder(
        private val binding: ItemContactBinding
    ) : RecyclerView.ViewHolder(binding.root) {

        fun bind(contact: Contact) {
            binding.tvName.text = contact.name
            binding.tvPhone.text = contact.phone
            binding.root.setOnClickListener { onContactClick(contact) }
        }
    }

    override fun onCreateViewHolder(parent: ViewGroup, viewType: Int): ContactViewHolder {
        val binding = ItemContactBinding.inflate(
            LayoutInflater.from(parent.context), parent, false
        )
        return ContactViewHolder(binding)
    }

    override fun onBindViewHolder(holder: ContactViewHolder, position: Int) {
        holder.bind(getItem(position))
    }
}
```

Usage — just submit a new list:

```kotlin
adapter.submitList(newContactsList)
// DiffUtil computes changes on a background thread, applies them on main
```

---

## LayoutManagers

```kotlin
// Vertical list (default)
binding.rvContacts.layoutManager = LinearLayoutManager(this)

// Horizontal list
binding.rvContacts.layoutManager = LinearLayoutManager(
    this, LinearLayoutManager.HORIZONTAL, false
)

// 2-column grid
binding.rvContacts.layoutManager = GridLayoutManager(this, 2)

// Staggered grid (Pinterest-style, different height items)
binding.rvContacts.layoutManager = StaggeredGridLayoutManager(
    2, StaggeredGridLayoutManager.VERTICAL
)
```

---

## Multiple View Types

```kotlin
class FeedAdapter : ListAdapter<FeedItem, RecyclerView.ViewHolder>(FeedDiffCallback()) {

    companion object {
        const val TYPE_POST = 0
        const val TYPE_AD = 1
    }

    override fun getItemViewType(position: Int): Int {
        return when (getItem(position)) {
            is FeedItem.Post -> TYPE_POST
            is FeedItem.Ad -> TYPE_AD
        }
    }

    override fun onCreateViewHolder(parent: ViewGroup, viewType: Int): RecyclerView.ViewHolder {
        return when (viewType) {
            TYPE_POST -> PostViewHolder(
                ItemPostBinding.inflate(LayoutInflater.from(parent.context), parent, false)
            )
            TYPE_AD -> AdViewHolder(
                ItemAdBinding.inflate(LayoutInflater.from(parent.context), parent, false)
            )
            else -> throw IllegalArgumentException("Unknown view type: $viewType")
        }
    }

    override fun onBindViewHolder(holder: RecyclerView.ViewHolder, position: Int) {
        when (val item = getItem(position)) {
            is FeedItem.Post -> (holder as PostViewHolder).bind(item)
            is FeedItem.Ad -> (holder as AdViewHolder).bind(item)
        }
    }
}
```

---

## Item Decorations

Add dividers or spacing without modifying item layouts:

```kotlin
// Built-in divider
binding.rvContacts.addItemDecoration(
    DividerItemDecoration(this, DividerItemDecoration.VERTICAL)
)

// Custom spacing
class SpaceItemDecoration(private val spaceInPx: Int) : RecyclerView.ItemDecoration() {
    override fun getItemOffsets(
        outRect: Rect,
        view: View,
        parent: RecyclerView,
        state: RecyclerView.State
    ) {
        outRect.bottom = spaceInPx
        if (parent.getChildAdapterPosition(view) == 0) {
            outRect.top = spaceInPx
        }
    }
}

binding.rvContacts.addItemDecoration(
    SpaceItemDecoration(resources.getDimensionPixelSize(R.dimen.item_spacing))
)
```

---

## Common Mistakes

### Mistake 1: Calling `notifyDataSetChanged()` for every update

```kotlin
// WRONG — redraws everything, loses animations
contacts = newList
notifyDataSetChanged()

// CORRECT — use ListAdapter.submitList()
adapter.submitList(newList)
```

### Mistake 2: Forgetting to set a LayoutManager

```kotlin
// WRONG — RecyclerView shows nothing without a LayoutManager
binding.rvContacts.adapter = adapter

// CORRECT
binding.rvContacts.layoutManager = LinearLayoutManager(this)
binding.rvContacts.adapter = adapter
```

### Mistake 3: Doing heavy work in `onBindViewHolder`

`onBindViewHolder` is called for every visible item during scroll. Keep it fast — bind data only, no I/O, no large computations.

### Mistake 4: `attachToRoot = true` in `onCreateViewHolder`

```kotlin
// WRONG — causes duplicate parent attachment crash
LayoutInflater.from(parent.context).inflate(R.layout.item_contact, parent, true)

// CORRECT — always false in RecyclerView adapters
LayoutInflater.from(parent.context).inflate(R.layout.item_contact, parent, false)
```

---

## Interview Questions

**Q1: What is the purpose of ViewHolder in RecyclerView?**

> ViewHolder caches view references for a single item. Without it, every `onBindViewHolder` call would invoke `findViewById` — which traverses the entire view hierarchy. ViewHolder prevents that, making scrolling smooth.

**Q2: What is the difference between `notifyDataSetChanged()` and `DiffUtil`?**

> `notifyDataSetChanged()` triggers a full rebind and redraw of all visible items, with no animations. `DiffUtil` computes the minimal set of changes (insertions, deletions, moves, changes) and applies them with animations — better UX and better performance.

**Q3: What is `ListAdapter` and how does it differ from a standard `RecyclerView.Adapter`?**

> `ListAdapter` is a subclass of `RecyclerView.Adapter` that uses `AsyncListDiffer` internally. It runs `DiffUtil` on a background thread automatically when you call `submitList()`, preventing jank on the main thread.

**Q4: When would you use multiple view types in a RecyclerView?**

> When a list needs to show heterogeneous content — e.g., a social feed with posts, ads, and date headers mixed together.

---

## Summary

- RecyclerView = LayoutManager + Adapter + ViewHolder
- Always use `ListAdapter` with `DiffUtil.ItemCallback` for data updates
- Never call `notifyDataSetChanged()` — use `submitList()` instead
- Keep `onBindViewHolder` fast — bind data only
- Use ItemDecoration for dividers/spacing — don't add padding/margins to item layouts for spacing

**Next:** [Chapter 7 — Material Components](./07-material-components.md)
