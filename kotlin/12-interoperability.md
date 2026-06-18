# Chapter 12 — Kotlin and Java Interoperability

> *"Kotlin is 100% interoperable with Java. You don't need to convert your whole project at once."*
> — JetBrains

---

## 12.1 Calling Java from Kotlin

Kotlin can call Java code transparently. Java classes, methods, and fields are accessible directly:

```kotlin
import java.util.Date
import java.util.ArrayList
import java.text.SimpleDateFormat

fun main() {
    // Java class instantiation — no 'new' keyword in Kotlin
    val date = Date()
    println(date)
    
    // Java collection
    val javaList = ArrayList<String>()
    javaList.add("Hello")
    javaList.add("World")
    println(javaList)  // [Hello, World]
    
    // Java static methods
    val formatted = SimpleDateFormat("yyyy-MM-dd").format(date)
    println(formatted)
    
    // Java System class
    println(System.currentTimeMillis())
    println(System.getenv("HOME"))
    println(System.getProperty("java.version"))
}
```

### Java Getters and Setters as Properties

Kotlin converts Java getter/setter methods into properties automatically:

```java
// Java class
public class JavaPerson {
    private String name;
    private int age;
    
    public JavaPerson(String name, int age) {
        this.name = name;
        this.age = age;
    }
    
    public String getName() { return name; }
    public void setName(String name) { this.name = name; }
    public int getAge() { return age; }
    public void setAge(int age) { this.age = age; }
    
    public boolean isEmpty() { return name.isEmpty(); }
}
```

```kotlin
// In Kotlin — getters/setters become properties
val person = JavaPerson("Alice", 30)

// These look like property accesses:
println(person.name)  // calls getName()
println(person.age)   // calls getAge()

person.name = "Bob"   // calls setName()
person.age = 25       // calls setAge()

// is-prefixed methods also become properties:
println(person.isEmpty)  // calls isEmpty() — no ()
```

### Java Static Members

Java's `static` fields and methods are accessed on the class (using companion-like syntax):

```java
// Java
public class MathUtils {
    public static final double TAX_RATE = 0.20;
    
    public static double calculateTax(double amount) {
        return amount * TAX_RATE;
    }
}
```

```kotlin
// Kotlin
println(MathUtils.TAX_RATE)               // 0.20
println(MathUtils.calculateTax(100.0))    // 20.0
```

---

## 12.2 Calling Kotlin from Java

### Kotlin Top-Level Functions

Kotlin top-level functions are compiled as static methods in a class named after the file:

```kotlin
// File: Utils.kt
package com.example

fun greet(name: String) = "Hello, $name!"
const val MAX_SIZE = 100
```

```java
// From Java
import com.example.UtilsKt;

String greeting = UtilsKt.greet("Alice");  // calls static method
int max = UtilsKt.MAX_SIZE;
```

### Renaming the File Class

Use `@JvmName` to rename the generated class:

```kotlin
@file:JvmName("Utils")  // put at the top of the file
package com.example

fun greet(name: String) = "Hello, $name!"
```

```java
// Now you can use:
String greeting = Utils.greet("Alice");  // cleaner!
```

### @JvmStatic — Companion Object Members

By default, companion object members are NOT static in Java. Use `@JvmStatic` to make them accessible as static:

```kotlin
class Config {
    companion object {
        @JvmStatic
        fun defaultConfig() = Config()
        
        @JvmField
        val DEFAULT_TIMEOUT = 30
    }
}
```

```java
// Java
Config config = Config.defaultConfig();  // static call — no Companion.
int timeout = Config.DEFAULT_TIMEOUT;    // field access
```

Without `@JvmStatic`:
```java
// Java without @JvmStatic
Config config = Config.Companion.defaultConfig();  // more verbose
```

### @JvmOverloads — Default Parameters

Kotlin's default parameters don't translate to Java overloads by default. Use `@JvmOverloads`:

```kotlin
@JvmOverloads
fun connect(host: String, port: Int = 8080, timeout: Int = 30) {
    println("Connecting to $host:$port with timeout $timeout")
}
```

```java
// Java — now has three overloads generated:
connect("localhost");                // uses defaults for port and timeout
connect("localhost", 443);           // uses default for timeout
connect("localhost", 443, 60);       // all parameters
```

### @JvmField — Properties as Fields

Kotlin properties are compiled with getter/setter methods. `@JvmField` exposes them as direct fields:

```kotlin
class User {
    @JvmField
    var name: String = ""  // accessible as field in Java
    
    var age: Int = 0  // accessed via getAge()/setAge() in Java
}
```

```java
// Java
User user = new User();
user.name = "Alice";  // direct field access (with @JvmField)
user.setAge(30);       // method call (without @JvmField)
```

---

## 12.3 Nullability in Interop

This is the most critical area of Kotlin-Java interoperability.

### Platform Types

When Kotlin calls Java code, the Java compiler doesn't tell Kotlin whether a value is null or not. Kotlin uses **platform types** (written as `T!`) — types that carry no nullability information.

```java
// Java method — returns String, but could it be null?
public class JavaService {
    public String getName() { ... }         // Might return null!
    public List<String> getItems() { ... }  // Might return null!
}
```

```kotlin
// Kotlin calling Java:
val service = JavaService()
val name = service.name  // Type is String! (platform type)
                          // Could be String or String?
                          // Kotlin trusts you to know

// These both compile — Kotlin doesn't force null handling:
val len1: Int = name.length    // risky — might NPE if name is null
val len2: Int? = name?.length  // safe — but verbose

// Best practice: assign to an explicitly typed variable
val name: String = service.name   // Non-null assertion
val name: String? = service.name  // Nullable assertion
```

### Java Nullability Annotations

If Java code uses nullability annotations, Kotlin respects them:

```java
import org.jetbrains.annotations.NotNull;
import org.jetbrains.annotations.Nullable;

public class AnnotatedService {
    @NotNull
    public String getRequiredValue() { return "value"; }
    
    @Nullable
    public String getOptionalValue() { return null; }
}
```

```kotlin
val service = AnnotatedService()

val required: String = service.requiredValue    // String (non-null)
val optional: String? = service.optionalValue   // String? (nullable)

// required.length is safe
// optional?.length is necessary
```

Supported annotation packages:
- `org.jetbrains.annotations` (`@NotNull`, `@Nullable`)
- `javax.annotation` (`@Nonnull`, `@Nullable`)
- `androidx.annotation` (Android)
- `org.checkerframework.checker.nullness`

### Nullability Best Practices in Interop

```kotlin
// Pattern 1: Explicit typing at the boundary
fun processFromJava(javaService: JavaService) {
    val name: String = javaService.name ?: return  // handle null at boundary
    val items: List<String> = javaService.items ?: emptyList()
    // From here, use non-nullable types
}

// Pattern 2: Check not null explicitly
fun getUser(service: JavaService): User {
    val json = checkNotNull(service.userJson) { "User JSON must not be null" }
    return parseUser(json)
}
```

---

## 12.4 Checked Exceptions

Java has **checked exceptions** — exceptions declared with `throws` in method signatures. Kotlin has **no checked exceptions** — all exceptions are unchecked.

```java
// Java
public String readFile(String path) throws IOException {
    // Must be declared; callers must handle or rethrow
}
```

```kotlin
// Kotlin — no checked exceptions needed
fun readFile(path: String): String {
    // No 'throws IOException' — all exceptions are unchecked
    return File(path).readText()
}

// Java calling Kotlin — doesn't know about the exception
// May get surprising runtime exceptions
```

### @Throws — Declaring Exceptions for Java

If your Kotlin code is called from Java, use `@Throws` to declare exceptions:

```kotlin
@Throws(IOException::class)
fun readFile(path: String): String = File(path).readText()
```

```java
// Java now sees the declaration:
try {
    String content = YourKt.readFile("data.txt");
} catch (IOException e) {
    // handle
}
```

---

## 12.5 Collections Interop

Kotlin and Java share the same collection types at the JVM level, but Kotlin adds read-only views:

```kotlin
// Kotlin's List<String> is actually java.util.List<String> at runtime
val kotlinList: List<String> = listOf("a", "b", "c")

// This is the same object! Just with a different interface.
val javaList: java.util.List<String> = kotlinList as java.util.List<String>

// Kotlin's MutableList is java.util.ArrayList under the hood
val mutableList: MutableList<String> = mutableListOf("a", "b", "c")
// same as:
val arrayList: java.util.ArrayList<String> = mutableListOf("a", "b", "c") as java.util.ArrayList<String>
```

### Converting Between Collections

```kotlin
import java.util.*

// Java collections to Kotlin
val javaArrayList = ArrayList<String>().apply { add("a"); add("b") }
val kotlinList: List<String> = javaArrayList.toList()      // immutable copy
val kotlinMutable: MutableList<String> = javaArrayList.toMutableList()

// Kotlin to Java (often not needed — same underlying type)
val kotlinList2 = listOf("x", "y")
val javaList2: java.util.List<String> = kotlinList2  // works directly

// Arrays
val kotlinArray = arrayOf("a", "b", "c")
val javaArray: Array<String> = kotlinArray   // same

val intArray = intArrayOf(1, 2, 3)           // int[] in Java
val integerArray = arrayOf(1, 2, 3)          // Integer[] in Java (boxed)
```

---

## 12.6 Kotlin-Specific Features in Java

Some Kotlin features are not directly accessible from Java without adapters:

### Data Classes

```kotlin
data class Point(val x: Double, val y: Double)
```

From Java:
```java
Point point = new Point(1.0, 2.0);
double x = point.getX();           // works
double y = point.getY();           // works
String s = point.toString();       // works
Point copy = point.copy(3.0, point.getY());  // works — component functions
                                              // but copy() has odd signature in Java
// Destructuring (component1, component2) is accessible but awkward in Java
```

### Extension Functions

Extension functions are static functions in Java:

```kotlin
// Kotlin
fun String.shout() = uppercase() + "!!!"
```

```java
// Java — must call as static function
String result = StringKt.shout("hello");  // or whatever the file class is
// You CANNOT call "hello".shout() from Java
```

### Named and Default Arguments

Not available in Java (unless `@JvmOverloads`):

```kotlin
fun connect(host: String = "localhost", port: Int = 8080) { }
```

```java
// Java must provide all arguments (or use @JvmOverloads overloads)
ConnectKt.connect("myserver", 9090);  // all arguments required
// ConnectKt.connect(port = 9090);  // DOESN'T COMPILE IN JAVA
```

### Object Declarations (Singletons)

```kotlin
object Config {
    val timeout = 30
    fun getConnectionString() = "jdbc://localhost"
}
```

```java
// Access via INSTANCE field:
int timeout = Config.INSTANCE.getTimeout();  // note INSTANCE
String conn = Config.INSTANCE.getConnectionString();

// With @JvmStatic:
// int timeout = Config.timeout;  // if annotated
```

---

## 12.7 Working with Java Frameworks

### Using Java-Style Builders

```java
// Java builder pattern (common in Java frameworks)
AlertDialog.Builder builder = new AlertDialog.Builder(context);
builder.setTitle("Alert");
builder.setMessage("Are you sure?");
builder.setPositiveButton("OK", listener);
AlertDialog dialog = builder.create();
```

```kotlin
// Kotlin's apply makes this clean:
val dialog = AlertDialog.Builder(context).apply {
    setTitle("Alert")
    setMessage("Are you sure?")
    setPositiveButton("OK", listener)
}.create()
```

### Lambda for SAM Interfaces

```java
// Java
button.setOnClickListener(new View.OnClickListener() {
    @Override
    public void onClick(View v) {
        System.out.println("Clicked!");
    }
});
```

```kotlin
// Kotlin SAM conversion
button.setOnClickListener {
    println("Clicked!")
}
```

---

## Interop Cheat Sheet

### Kotlin → Java

| Kotlin | Java |
|--------|------|
| `fun` (top-level) | `static` in `FileNameKt` class |
| `companion object` | Inner `Companion` object |
| `@JvmStatic` | Static method/field |
| `@JvmField` | Direct field access |
| `@JvmOverloads` | Generates overloaded methods |
| `@Throws` | Checked exception declaration |
| `object Singleton` | `Singleton.INSTANCE` |
| `val prop` | `getProp()` |
| `var prop` | `getProp()` / `setProp()` |
| `@file:JvmName("X")` | Class named `X` instead of `FileNameKt` |

### Java → Kotlin

| Java | Kotlin |
|------|--------|
| `obj.getName()` | `obj.name` |
| `obj.setName(x)` | `obj.name = x` |
| `obj.isEmpty()` | `obj.isEmpty` |
| Static method | `ClassName.method()` |
| `throws IOException` | Unchecked (add `@Throws` for Java callers) |
| `null` return | Platform type `T!` — handle explicitly |

---

## Summary

Kotlin and Java achieve seamless interoperability because Kotlin compiles to standard JVM bytecode. Java code is called from Kotlin transparently — getters/setters become properties, static members are accessible on the class. Kotlin code is accessible from Java via the generated class name (`FileKt`); `@JvmStatic`, `@JvmField`, `@JvmOverloads`, and `@JvmName` annotations control the generated Java API. The critical null safety challenge arises with Java's **platform types** — types that carry no null information. Use explicit type declarations and null checks at Java-Kotlin boundaries. Checked exceptions are Java-only; use `@Throws` to declare Kotlin exceptions for Java callers.

---

## Key Takeaways

- Java getters/setters are automatically accessible as Kotlin properties
- Kotlin top-level functions compile to static methods in a `FileNameKt` class
- `@JvmStatic`, `@JvmField`, `@JvmOverloads` are needed for idiomatic Java usage of Kotlin
- Platform types (`T!`) carry no null safety — assign explicitly to `T` or `T?` at boundaries
- Kotlin and Java share the same collection types at runtime — no conversion needed
- Checked exceptions are Java-only; use `@Throws` if Kotlin code is called from Java
- Extension functions are static Java methods — not callable as instance methods from Java

---

## Practice Questions

### Conceptual
1. What is a platform type and why is it dangerous?
2. When do you need `@JvmStatic` and when don't you?
3. Why does Kotlin have no checked exceptions?
4. What happens if you call a Java method that returns null and assign it to `val x: String` in Kotlin?
5. How does the `apply` scope function help when working with Java builder patterns?

### Code Exercises

**Exercise 1:** Create a Kotlin file `StringUtils.kt` with:
- A top-level `fun reverseWords(text: String): String`
- A top-level constant `MAX_STRING_LENGTH`
Write the Java interop annotations needed so Java can call it as `StringUtils.reverseWords(text)`.

**Exercise 2:** Write a Kotlin class `DataProcessor` with a companion object that:
- Has a `create()` factory method (accessible as static from Java)
- Has a `DEFAULT_BATCH_SIZE` constant (accessible as field from Java)

**Exercise 3:** Given this Java class with platform types, write a safe Kotlin wrapper:
```java
public class LegacyApi {
    public String fetchName(String id) { ... }  // might return null
    public List<String> fetchTags(String id) { ... }  // might return null
}
```

**Exercise 4:** Write a Kotlin function `@Throws(ValidationException::class) fun validateAge(age: Int)` and explain why the `@Throws` annotation is needed.

**Exercise 5:** Demonstrate the difference between accessing a Kotlin `companion object` member from Java with and without `@JvmStatic`.

---

*Next: [Chapter 13 — Kotlin Internals](13-internals.md)*
