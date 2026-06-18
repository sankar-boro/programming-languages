# Chapter 6: CI/CD Basics for Android

## What Is CI/CD?

**CI (Continuous Integration):** Every code change automatically triggers a build + test run.
**CD (Continuous Delivery):** Passing builds are automatically deployed to testers or the store.

**Why it matters:**
- Catch bugs in minutes, not days
- Prevent "works on my machine" failures
- Automated release process — no manual APK uploads
- Consistent code quality across the team

---

## GitHub Actions — Most Common for Android

GitHub Actions runs your CI pipeline on every push or pull request.

### Basic CI Workflow

`.github/workflows/ci.yml`:

```yaml
name: Android CI

on:
  push:
    branches: [ main, develop ]
  pull_request:
    branches: [ main, develop ]

jobs:
  test:
    name: Run Tests
    runs-on: ubuntu-latest

    steps:
      - name: Checkout code
        uses: actions/checkout@v4

      - name: Set up JDK 17
        uses: actions/setup-java@v4
        with:
          java-version: '17'
          distribution: 'temurin'

      - name: Cache Gradle packages
        uses: actions/cache@v4
        with:
          path: |
            ~/.gradle/caches
            ~/.gradle/wrapper
          key: ${{ runner.os }}-gradle-${{ hashFiles('**/*.gradle*', '**/gradle-wrapper.properties') }}
          restore-keys: ${{ runner.os }}-gradle-

      - name: Grant execute permission for gradlew
        run: chmod +x gradlew

      - name: Run unit tests
        run: ./gradlew testDebugUnitTest

      - name: Run lint
        run: ./gradlew lintDebug

      - name: Upload test results
        uses: actions/upload-artifact@v4
        if: always()
        with:
          name: test-results
          path: '**/build/reports/tests/'

      - name: Upload lint results
        uses: actions/upload-artifact@v4
        if: always()
        with:
          name: lint-results
          path: '**/build/reports/lint-results-debug.html'
```

---

### Build and Release Workflow

`.github/workflows/release.yml`:

```yaml
name: Release Build

on:
  push:
    tags:
      - 'v*'  # Triggers on tags like v1.0.0, v2.1.3

jobs:
  release:
    name: Build Release AAB
    runs-on: ubuntu-latest

    steps:
      - uses: actions/checkout@v4

      - name: Set up JDK 17
        uses: actions/setup-java@v4
        with:
          java-version: '17'
          distribution: 'temurin'

      - name: Cache Gradle
        uses: actions/cache@v4
        with:
          path: |
            ~/.gradle/caches
            ~/.gradle/wrapper
          key: ${{ runner.os }}-gradle-${{ hashFiles('**/*.gradle*') }}

      - name: Decode Keystore
        run: |
          echo "${{ secrets.KEYSTORE_BASE64 }}" | base64 --decode > app/keystore.jks

      - name: Build Release AAB
        run: ./gradlew bundleRelease
        env:
          KEYSTORE_PASSWORD: ${{ secrets.KEYSTORE_PASSWORD }}
          KEY_ALIAS: ${{ secrets.KEY_ALIAS }}
          KEY_PASSWORD: ${{ secrets.KEY_PASSWORD }}

      - name: Upload AAB to Play Store (Internal Track)
        uses: r0adkll/upload-google-play@v1
        with:
          serviceAccountJsonPlainText: ${{ secrets.PLAY_SERVICE_ACCOUNT_JSON }}
          packageName: com.yourname.app
          releaseFiles: app/build/outputs/bundle/release/*.aab
          track: internal
          status: completed
```

---

## Signing Configuration

Never commit keystores or passwords. Use Gradle to read from environment variables:

`app/build.gradle.kts`:

```kotlin
android {
    signingConfigs {
        create("release") {
            storeFile = file(System.getenv("KEYSTORE_PATH") ?: "keystore.jks")
            storePassword = System.getenv("KEYSTORE_PASSWORD") ?: ""
            keyAlias = System.getenv("KEY_ALIAS") ?: ""
            keyPassword = System.getenv("KEY_PASSWORD") ?: ""
        }
    }

    buildTypes {
        release {
            signingConfig = signingConfigs.getByName("release")
            isMinifyEnabled = true
            isShrinkResources = true
            proguardFiles(
                getDefaultProguardFile("proguard-android-optimize.txt"),
                "proguard-rules.pro"
            )
        }
    }
}
```

Store secrets in **GitHub → Settings → Secrets and variables → Actions**.

---

## Automated Version Bumping

```kotlin
// app/build.gradle.kts — read version from git tags
fun getVersionCode(): Int {
    return try {
        val process = ProcessBuilder("git", "rev-list", "--count", "HEAD")
            .start()
        process.inputStream.bufferedReader().readLine().trim().toInt()
    } catch (e: Exception) { 1 }
}

fun getVersionName(): String {
    return try {
        val process = ProcessBuilder("git", "describe", "--tags", "--abbrev=0")
            .start()
        process.inputStream.bufferedReader().readLine().trim()
    } catch (e: Exception) { "1.0.0" }
}

android {
    defaultConfig {
        versionCode = getVersionCode()
        versionName = getVersionName()
    }
}
```

---

## Lint Configuration

```kotlin
// app/build.gradle.kts
android {
    lint {
        abortOnError = true       // Fail CI on lint errors
        warningsAsErrors = false   // Treat warnings as warnings, not errors
        checkReleaseBuilds = true
        disable += setOf(
            "GradleDependency",   // Disable specific checks that are too noisy
            "ObsoleteSdkInt"
        )
        htmlReport = true
        htmlOutput = file("${buildDir}/reports/lint/lint.html")
    }
}
```

`.github/workflows/ci.yml` — add lint check:

```yaml
- name: Run Lint
  run: ./gradlew lintDebug

- name: Check lint result
  run: |
    if grep -q "error" app/build/reports/lint/lint.html; then
      echo "Lint errors found!"
      exit 1
    fi
```

---

## Detekt — Static Analysis for Kotlin

```kotlin
// build.gradle.kts (project level)
plugins {
    id("io.gitlab.arturbosch.detekt") version "1.23.6"
}

detekt {
    buildUponDefaultConfig = true
    config.setFrom(files("$rootDir/config/detekt/detekt.yml"))
}
```

`config/detekt/detekt.yml`:

```yaml
style:
  MaxLineLength:
    maxLineLength: 120
  MagicNumber:
    active: true
    ignoreNumbers: ['-1', '0', '1', '2']
complexity:
  LongMethod:
    threshold: 60
  LargeClass:
    threshold: 600
```

CI step:

```yaml
- name: Run Detekt
  run: ./gradlew detekt
```

---

## Firebase App Distribution (Beta Testing)

```yaml
- name: Distribute to testers
  uses: wzieba/Firebase-Distribution-Github-Action@v1
  with:
    appId: ${{ secrets.FIREBASE_APP_ID }}
    serviceCredentialsFileContent: ${{ secrets.FIREBASE_CREDENTIALS }}
    groups: internal-testers
    file: app/build/outputs/apk/debug/app-debug.apk
    releaseNotes: "Build from ${{ github.sha }}"
```

---

## Sample Complete CI Pipeline

```yaml
name: Full CI Pipeline

on:
  pull_request:
    branches: [main]

jobs:
  quality:
    name: Code Quality
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-java@v4
        with: { java-version: '17', distribution: 'temurin' }
      - uses: actions/cache@v4
        with:
          path: ~/.gradle
          key: gradle-${{ hashFiles('**/*.gradle*') }}
      - run: ./gradlew lintDebug detekt
      - uses: actions/upload-artifact@v4
        if: always()
        with:
          name: quality-reports
          path: '**/build/reports/'

  unit-tests:
    name: Unit Tests
    runs-on: ubuntu-latest
    needs: quality
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-java@v4
        with: { java-version: '17', distribution: 'temurin' }
      - uses: actions/cache@v4
        with:
          path: ~/.gradle
          key: gradle-${{ hashFiles('**/*.gradle*') }}
      - run: ./gradlew testDebugUnitTest
      - uses: actions/upload-artifact@v4
        if: always()
        with:
          name: test-results
          path: '**/build/reports/tests/'

  build:
    name: Build Debug APK
    runs-on: ubuntu-latest
    needs: unit-tests
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-java@v4
        with: { java-version: '17', distribution: 'temurin' }
      - uses: actions/cache@v4
        with:
          path: ~/.gradle
          key: gradle-${{ hashFiles('**/*.gradle*') }}
      - run: ./gradlew assembleDebug
      - uses: actions/upload-artifact@v4
        with:
          name: debug-apk
          path: app/build/outputs/apk/debug/app-debug.apk
```

---

## Interview Questions

**Q1: What is the difference between CI and CD?**

> CI (Continuous Integration) automatically builds and tests code on every change. CD (Continuous Delivery) automatically deploys the built artifact to a distribution channel (Firebase, Play Store internal track) after passing CI. CD extends CI by automating the release step.

**Q2: Why should signing credentials never be committed to a repository?**

> The keystore and passwords, if exposed, allow anyone to sign and distribute an app under your identity on the Play Store. They must be stored as encrypted secrets in the CI system and injected as environment variables at build time.

**Q3: What is Lint and why run it in CI?**

> Android Lint is a static analysis tool that checks for potential bugs, usability issues, and style violations in your code and resources. Running it in CI ensures issues are caught automatically before code merges, preventing common mistakes from reaching production.

---

## Summary

- GitHub Actions is the most common CI for Android — runs on every push/PR
- Always cache Gradle for faster builds
- Store signing credentials as encrypted GitHub Secrets
- Run lint + unit tests on every PR; release builds on tag pushes
- Use Detekt for Kotlin-specific static analysis
- Use Firebase App Distribution for beta builds before Play Store

**Next:** [Chapter 7 — Production Practices](./07-production-practices.md)
