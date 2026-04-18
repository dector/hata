# Agent Guidelines for Hata Android Project

This document provides essential information for AI coding agents working on the Hata Android application.

## Project Overview

**Tech Stack:**
- Language: Kotlin (version 2.3.0)
- Build System: Gradle (9.0.0 AGP)
- UI Framework: Jetpack Compose (BOM 2026.01.00)
- UI Components: Material Design 3
- Min SDK: 24 (Android 7.0)
- Target SDK: 36
- Java Compatibility: Java 11

**Application ID:** `space.dector.hata`

## Build Commands

### Build & Assemble
```bash
./gradlew build                    # Full build with tests
./gradlew assemble                 # Build without running tests
./gradlew assembleDebug            # Build debug APK
./gradlew assembleRelease          # Build release APK
./gradlew clean                    # Clean build artifacts
./gradlew bundle                   # Create app bundle (AAB)
```

### Install, Deploy & Run
```bash
./gradlew installDebug             # Build + install debug APK on connected device/emulator
./gradlew uninstallDebug           # Uninstall debug build
adb devices                        # List connected devices
adb shell am start -n space.dector.hata/hata.app.AppActivity   # Launch app activity
adb shell monkey -p space.dector.hata -c android.intent.category.LAUNCHER 1  # Launch via launcher intent
```

## Test Commands

### Unit Tests (JUnit)
```bash
./gradlew test                     # Run all unit tests
./gradlew testDebugUnitTest        # Run debug unit tests
./gradlew testDebugUnitTest --tests "hata.ExampleUnitTest"           # Run specific test class
./gradlew testDebugUnitTest --tests "hata.ExampleUnitTest.addition_isCorrect"  # Run single test
./gradlew testDebugUnitTest --debug-jvm                              # Debug tests (port 5005)
./gradlew testDebugUnitTest --tests "*Test" --fail-fast              # Stop on first failure
```

### Instrumented Tests (Android)
```bash
./gradlew connectedDebugAndroidTest                                  # Run on connected device
./gradlew connectedAndroidTest                                       # All flavors on all devices
./gradlew connectedCheck                                             # All device checks
```

### Lint & Code Quality
```bash
./gradlew lint                     # Run lint on default variant
./gradlew lintDebug                # Lint debug variant with output
./gradlew lintRelease              # Lint release variant
./gradlew lintFix                  # Auto-fix safe issues
./gradlew updateLintBaseline       # Update lint baseline
./gradlew check                    # All checks (lint + tests)
```

## Project Structure

```
android/
├── app/
│   ├── src/
│   │   ├── main/
│   │   │   ├── kotlin/hata/              # Main source code
│   │   │   │   ├── AppActivity.kt        # Main entry activity
│   │   │   │   └── ui/theme/             # Theme & styling
│   │   │   ├── res/                      # Android resources
│   │   │   └── AndroidManifest.xml
│   │   ├── test/kotlin/hata/             # Unit tests (JUnit)
│   │   └── androidTest/kotlin/hata/      # Instrumented tests
│   ├── build.gradle.kts                  # App module config
│   └── proguard-rules.pro
├── gradle/
│   └── libs.versions.toml                # Version catalog
├── build.gradle.kts                      # Root build config
└── settings.gradle.kts
```

## Code Style Guidelines

### Package Structure
- Root package: `hata`
- UI components: `hata.ui.*`
- Theme: `hata.ui.theme`
- Follow feature-based or layer-based organization as project grows

### Kotlin Style

**Naming Conventions:**
- Classes: `PascalCase` (e.g., `AppActivity`)
- Functions: `camelCase` (e.g., `onCreate()`)
- Composables: `PascalCase` (e.g., `Greeting()`, `HataTheme()`)
- Constants/vals at file level: `PascalCase` for colors (e.g., `Purple80`)
- Private vals: `camelCase` or `PascalCase` depending on context

**Imports:**
- Organize imports: Android SDK → AndroidX → Compose → Third-party → Project
- No wildcard imports (avoid `import foo.*`)
- Group related imports together with blank lines between groups

**Formatting:**
- Indentation: 4 spaces
- Max line length: Use judgment, but wrap long parameter lists
- Trailing commas: Use for multi-line parameter lists
- Modifier order: Follow standard Kotlin conventions

**Composables:**
- Use `@Composable` annotation
- Add `@Preview` or `@PreviewLightDark` for preview functions
- Default parameters should be `Modifier = Modifier`
- Preview functions should use helper wrappers (e.g., `preview {}`)

**Types:**
- Use type inference where clear: `val name = "value"`
- Explicit types for public APIs: `fun foo(): String`
- Use nullable types (`?`) appropriately
- Prefer `val` over `var` when possible

### Error Handling
- Use `try-catch` for expected errors
- Don't catch generic `Exception` without re-throwing or logging
- Use `Result` type for operations that can fail
- Handle Android lifecycle appropriately (e.g., in `onCreate()`)

### Comments
- Use KDoc (`/** */`) for public APIs
- Inline comments (`//`) for complex logic explanation
- Avoid obvious comments that restate the code
- Keep comments up-to-date with code changes

## Dependencies

Managed via Gradle Version Catalog (`gradle/libs.versions.toml`):

**Core Dependencies:**
- `androidx.core:core-ktx`
- `androidx.lifecycle:lifecycle-runtime-ktx`
- `androidx.activity:activity-compose`

**Compose:**
- BOM-managed versions (no explicit version numbers in app build.gradle)
- Material3 for UI components
- UI tooling for previews and debugging

**Testing:**
- JUnit 4 for unit tests
- AndroidX Test for instrumentation
- Espresso for UI testing
- Compose UI test for Compose testing

### Adding Dependencies
1. Add version to `[versions]` section in `libs.versions.toml`
2. Add library to `[libraries]` section
3. Reference in `app/build.gradle.kts` using `libs.`

## Testing Patterns

**Unit Tests (app/src/test/):**
- Use JUnit 4: `@Test` annotation
- Import assertions: `import org.junit.Assert.*`
- Name format: `methodName_expectedBehavior()` (e.g., `addition_isCorrect()`)

**Instrumented Tests (app/src/androidTest/):**
- Use `@RunWith(AndroidJUnit4::class)`
- Access app context: `InstrumentationRegistry.getInstrumentation().targetContext`
- Compose UI tests: Use `androidx.compose.ui.test` APIs

## Build Configuration

**Gradle Files:**
- Use Kotlin DSL (`.gradle.kts`)
- Version catalog for dependency management
- Plugins applied via `alias(libs.plugins.*)` syntax

**Build Variants:**
- Debug: Default development build with tooling
- Release: Minified/optimized (ProGuard configured but disabled by default)

## Common Tasks

**Add a new Composable:**
1. Create in appropriate package under `hata.ui.*`
2. Add `@Composable` annotation
3. Follow Material3 design patterns
4. Add `@Preview` for development visibility

**Add a new test:**
- Unit test: Create in `app/src/test/kotlin/hata/`
- Instrumented: Create in `app/src/androidTest/kotlin/hata/`
- Match package structure of code under test

**Update dependencies:**
1. Edit `gradle/libs.versions.toml`
2. Sync Gradle
3. Test affected features

## Notes for Agents

- Always run `./gradlew build` after significant changes
- Use `--tests` flag for quick iteration on specific tests
- Prefer Compose over XML layouts (project uses Compose)
- Enable EdgeToEdge in activities (already done in `AppActivity`)
- Use Material3 theming system (`MaterialTheme.colorScheme.*`)
- Theme supports dynamic colors (Android 12+) and light/dark modes
