# Chapter 2 — Setting Up and Running Kotlin

> *"The best way to learn a programming language is to write programs in it."*
> — Dennis Ritchie

---

## 2.1 Installing Kotlin

Kotlin can be installed and run in several ways. Choose the method that fits your workflow.

### Method 1: SDKMAN! (Recommended for Linux/macOS)

SDKMAN! is a tool for managing parallel versions of SDKs on Unix-based systems.

```bash
# Install SDKMAN! (if not already installed)
curl -s "https://get.sdkman.io" | bash
source "$HOME/.sdkman/bin/sdkman-init.sh"

# Install the latest Kotlin compiler
sdk install kotlin

# Verify installation
kotlin -version
# Output: Kotlin version 2.0.x-release-xxx (JRE 17.x.x)
```

### Method 2: Homebrew (macOS)

```bash
brew update
brew install kotlin

# Verify
kotlin -version
kotlinc -version
```

### Method 3: Manual Installation (All Platforms)

1. Download the Kotlin compiler from the [Kotlin releases page](https://github.com/JetBrains/kotlin/releases)
2. Download `kotlin-compiler-x.x.x.zip`
3. Extract to a directory, e.g., `/opt/kotlin`
4. Add to PATH:

```bash
# Linux/macOS — add to ~/.bashrc or ~/.zshrc
export PATH=$PATH:/opt/kotlin/bin

# Windows (PowerShell)
$env:PATH += ";C:\kotlin\bin"
```

### Method 4: JetBrains IntelliJ IDEA

IntelliJ IDEA (Community or Ultimate) includes Kotlin support out of the box — no separate installation needed. This is the most feature-rich development environment for Kotlin.

1. Download IntelliJ IDEA Community (free) from jetbrains.com
2. Create a new Kotlin project: `File → New → Project → Kotlin`
3. The Kotlin SDK is bundled automatically

### Verifying Your Installation

```bash
kotlinc -version
# kotlin-compiler-daemon: info: using kotlin home: /usr/local/lib/kotlin
# kotlinc-jvm version 2.0.x (JRE 17.x.x)

kotlin -version
# Kotlin version 2.0.x
```

### Prerequisites: Java

Kotlin on the JVM requires a **Java Development Kit (JDK)**. Kotlin supports JDK 8 and above. For modern development, use JDK 17 or 21 (LTS releases).

```bash
# Check Java version
java -version
# openjdk version "21.0.x"

# JAVA_HOME should be set
echo $JAVA_HOME
# /usr/lib/jvm/java-21-openjdk
```

---

## 2.2 The Kotlin REPL

The REPL (Read-Eval-Print Loop) lets you experiment with Kotlin interactively, without creating files.

### Starting the REPL

```bash
kotlinc-jvm
# Welcome to Kotlin version 2.0.x (JRE 21.0.x)
# Type :help for help, :quit for quit
# >>>
```

### Using the REPL

```kotlin
>>> val message = "Hello, Kotlin!"
>>> println(message)
Hello, Kotlin!

>>> 2 + 2
res1: kotlin.Int = 4

>>> fun greet(name: String) = "Hello, $name!"
>>> greet("World")
res2: kotlin.String = Hello, World!

>>> (1..10).filter { it % 2 == 0 }
res3: kotlin.collections.List<kotlin.Int> = [2, 4, 6, 8, 10]
```

### REPL Commands

| Command | Action |
|---------|--------|
| `:help` | Show help |
| `:quit` | Exit the REPL |
| `:load <file>` | Load a Kotlin script file |
| `:type <expr>` | Show the type of an expression |
| `Tab` | Auto-complete |

### Kotlin REPL in IntelliJ IDEA

IntelliJ IDEA has a built-in Kotlin REPL:
- `Tools → Kotlin → Kotlin REPL`
- Or press `Shift+Ctrl+Alt+K` (Windows/Linux) / `Shift+Cmd+Opt+K` (macOS)

This REPL is more powerful — it has syntax highlighting, auto-completion, and access to your project's classes.

---

## 2.3 Running Kotlin from the Command Line

### Compiling and Running a Kotlin File

Create a file called `hello.kt`:

```kotlin
fun main() {
    println("Hello, World!")
}
```

Compile it:

```bash
kotlinc hello.kt -include-runtime -d hello.jar
```

Run it:

```bash
java -jar hello.jar
# Hello, World!
```

The `-include-runtime` flag bundles the Kotlin runtime into the JAR. Without it, you'd need the Kotlin runtime on the classpath separately.

### Running Kotlin Scripts Directly

Kotlin supports script files with the `.kts` extension. These don't need a `main` function:

```kotlin
// greet.kts
val name = "Kotlin"
println("Hello, $name!")
println("Kotlin version is awesome.")
```

Run directly without compiling:

```bash
kotlinc -script greet.kts
# Hello, Kotlin!
# Kotlin version is awesome.
```

### The kotlinc Command Reference

```bash
# Compile to JAR
kotlinc MyFile.kt -include-runtime -d output.jar

# Compile to class files (without runtime)
kotlinc MyFile.kt -d output/

# Run a script
kotlinc -script MyScript.kts

# Run the Kotlin REPL
kotlinc-jvm

# Show all compiler options
kotlinc -help
```

---

## 2.4 The Kotlin Playground

The **Kotlin Playground** (play.kotlinlang.org) is an online environment where you can write, run, and share Kotlin code without any local installation.

Features:
- Full Kotlin support with syntax highlighting
- Multiple files
- Different Kotlin versions selectable
- Shareable URLs for code snippets
- Embeddable in websites and documentation

This is the fastest way to try Kotlin concepts from this book without setup overhead.

---

## 2.5 Basic Program Structure

### The main Function

Every Kotlin program begins execution at the `main` function:

```kotlin
fun main() {
    println("Hello, World!")
}
```

Unlike Java, Kotlin's `main` function is a **top-level function** — it does not need to be inside a class. No `public static void` required.

### main with Arguments

If you need command-line arguments:

```kotlin
fun main(args: Array<String>) {
    if (args.isEmpty()) {
        println("No arguments provided")
    } else {
        println("Arguments: ${args.joinToString(", ")}")
    }
}
```

### Packages

Kotlin files can belong to packages, just like Java:

```kotlin
package com.example.greetings

fun greet(name: String) = "Hello, $name!"

fun main() {
    println(greet("Kotlin"))
}
```

Unlike Java, the file name does NOT need to match the class name, and the directory structure does NOT need to match the package name (though it's conventional to do so).

### Imports

```kotlin
package com.example

import kotlin.math.sqrt
import kotlin.math.PI
import java.util.Date           // Java classes work directly
import java.util.Collections.*  // star import

fun main() {
    println(sqrt(16.0))  // 4.0
    println(PI)          // 3.141592653589793
}
```

### Top-Level Declarations

Kotlin allows functions, variables, and even classes at the top level — outside of any class:

```kotlin
// Top-level constant
const val MAX_SIZE = 100

// Top-level variable
var globalCount = 0

// Top-level function
fun doubleIt(n: Int) = n * 2

// Top-level class
class Calculator {
    fun add(a: Int, b: Int) = a + b
}

fun main() {
    println(doubleIt(5))          // 10
    println(Calculator().add(3, 4))  // 7
}
```

This is a significant departure from Java, where everything must be inside a class.

### Comments

```kotlin
// Single-line comment

/*
   Multi-line comment
   Can span several lines
*/

/**
 * KDoc comment — for generating documentation
 * @param name The person's name
 * @return A greeting string
 */
fun greet(name: String): String = "Hello, $name!"

fun main() {
    // println("This line is commented out")
    println(greet("Kotlin"))  // Hello, Kotlin!
}
```

### Semicolons Are Optional

Kotlin does not require semicolons at the end of statements. The compiler uses **newlines** to determine statement boundaries.

```kotlin
// These are all valid:
val a = 1
val b = 2; val c = 3  // semicolons can separate multiple statements on one line

// But this is preferred:
val a = 1
val b = 2
val c = 3
```

---

## 2.6 Kotlin Toolchain Overview

Understanding the Kotlin toolchain helps you build, test, and run Kotlin code effectively.

### The Kotlin Compiler (kotlinc)

The Kotlin compiler (`kotlinc`) transforms `.kt` source files into:

- **JVM bytecode** (`.class` files, `.jar` archives) — for running on any JVM
- **JavaScript** — via Kotlin/JS for browser/Node.js targets
- **Native binaries** — via Kotlin/Native for iOS, Linux, macOS, Windows

```
[.kt source files]
      │
      ▼
  [kotlinc]
      │
      ├──► [JVM bytecode / .jar]  ──► runs on JVM
      ├──► [JavaScript]           ──► runs in browser/Node.js
      └──► [Native binary]        ──► runs natively (no JVM)
```

### Build Tools

For real projects, you use build tools rather than calling `kotlinc` directly:

**Gradle (Kotlin DSL)** — most common for JVM and Android projects:
```kotlin
// build.gradle.kts
plugins {
    kotlin("jvm") version "2.0.0"
}

dependencies {
    implementation(kotlin("stdlib"))
}
```

**Maven** — used in enterprise JVM projects:
```xml
<!-- pom.xml -->
<dependency>
    <groupId>org.jetbrains.kotlin</groupId>
    <artifactId>kotlin-stdlib</artifactId>
    <version>2.0.0</version>
</dependency>
```

### The Kotlin Standard Library

The **Kotlin Standard Library** (stdlib) is automatically available in most setups. It provides:

- Extension functions on Java standard library types
- Collections API
- String utilities
- Math utilities
- Coroutines support (separate library)
- I/O helpers

You rarely need to import stdlib items explicitly — they're imported by default or discovered by the IDE.

### The K2 Compiler (Kotlin 2.0+)

Kotlin 2.0 introduced the **K2 compiler**, a complete rewrite of the compiler frontend:

- **2-3x faster compilation** in many real-world projects
- Better IDE integration
- Foundation for future language features
- More consistent error messages

You don't need to do anything special to use K2 — it's the default in Kotlin 2.0+.

### Kotlin Targets

| Target | Description | Typical Use |
|--------|-------------|-------------|
| Kotlin/JVM | Compiles to JVM bytecode | Server, Android, desktop |
| Kotlin/JS | Compiles to JavaScript | Browser, Node.js |
| Kotlin/Native | Compiles to native code | iOS, Linux, macOS, Windows |
| Kotlin Multiplatform | Share code across targets | Cross-platform libraries |

This book focuses on **Kotlin/JVM** — the most common and fully-featured target.

---

## A Complete "Hello World" Walk-Through

Let's walk through every aspect of the simplest Kotlin program:

```kotlin
// File: Hello.kt
// Package declaration (optional at top level)
package com.example

// The main function — program entry point
fun main() {
    // println is a standard library function
    // It prints to stdout with a newline
    println("Hello, World!")
    
    // print doesn't add a newline
    print("Hello ")
    print("again!")  // Output: Hello again! (on same line)
    println()        // Just a newline
    
    // String template
    val name = "Kotlin"
    println("Hello, $name!")  // Hello, Kotlin!
    
    // Expression in template
    val x = 5
    println("$x squared is ${x * x}")  // 5 squared is 25
}
```

To run this:

```bash
kotlinc Hello.kt -include-runtime -d hello.jar
java -jar hello.jar

# Output:
# Hello, World!
# Hello again!
# Hello, Kotlin!
# 5 squared is 25
```

---

## Summary

Kotlin can be installed via SDKMAN!, Homebrew, or manually, and comes bundled with IntelliJ IDEA. The REPL (`kotlinc-jvm`) allows interactive exploration. Kotlin programs use a top-level `main` function as the entry point — no class needed. The `kotlinc` compiler converts `.kt` files to JVM bytecode, JavaScript, or native binaries. The Kotlin Playground at play.kotlinlang.org is the fastest way to experiment with the language.

---

## Key Takeaways

- Kotlin requires a JDK (Java Development Kit) when targeting the JVM
- The Kotlin REPL is excellent for rapid experimentation
- `kotlinc -script file.kts` runs Kotlin scripts without explicit compilation
- Top-level functions, variables, and classes are a key difference from Java
- Semicolons are optional — use newlines to separate statements
- The K2 compiler (default in Kotlin 2.0) brings significant compilation speed improvements
- Build tools like Gradle or Maven are used in real projects

---

## Practice Questions

### Conceptual
1. What is the purpose of the `-include-runtime` flag when compiling Kotlin?
2. What is the difference between a `.kt` file and a `.kts` file?
3. Why doesn't Kotlin require `main` to be inside a class?
4. What is SDKMAN! and what problem does it solve?

### Code Exercises

**Exercise 1:** Write a Kotlin script (`intro.kts`) that:
- Declares your name and age
- Prints a formatted introduction using string templates
- Prints the result of a mathematical expression

**Exercise 2:** Create a `Calculator.kt` file with:
- A top-level constant `PI = 3.14159`
- A top-level function `circleArea(radius: Double): Double`
- A `main` function that prints the area for radius = 5.0

**Exercise 3:** Experiment in the REPL:
- Define a list of numbers
- Filter for even numbers
- Print the result
- Note the type that is inferred

**Exercise 4:** Explore the Kotlin Playground:
- Visit play.kotlinlang.org
- Run the default example
- Modify it to print "Hello from Chapter 2!"
- Share the URL with a friend

---

*Next: [Chapter 3 — Variables, Types, and Operators](03-basic-syntax.md)*
