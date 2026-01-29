# Project Index: Hata Android

Generated: 2026-01-30

## 📁 Project Structure

```
android/
├── app/                           # Main Android application module
│   ├── src/
│   │   ├── main/kotlin/hata/
│   │   │   ├── HataApplication.kt      # Application entry point (Hilt)
│   │   │   ├── AppActivity.kt          # Main Activity (Compose + EdgeToEdge)
│   │   │   ├── AppRouter.kt            # Navigation router (Nav3)
│   │   │   ├── data/                   # Data layer
│   │   │   │   ├── api/                # API services & implementations
│   │   │   │   └── models/             # Data models
│   │   │   ├── domain/                 # Domain layer
│   │   │   │   └── usecases/           # Use cases
│   │   │   ├── ui/                     # UI layer (Compose)
│   │   │   │   ├── components/         # Reusable UI components
│   │   │   │   ├── home/               # Home feature
│   │   │   │   ├── theme/              # Material3 theme
│   │   │   │   └── utils/              # UI utilities
│   │   │   └── di/                     # Dependency injection (Hilt)
│   │   ├── res/                        # Android resources
│   │   ├── test/                       # Unit tests (JUnit)
│   │   └── androidTest/                # Instrumented tests
│   └── build.gradle.kts                # App module config
├── integrations/
│   └── wiz/                            # WiZ smart lights integration module
│       ├── src/main/kotlin/hata/integrations/wiz/
│       │   ├── WizModels.kt            # Data models
│       │   ├── WizDiscovery.kt         # Device discovery (UDP)
│       │   ├── WizControl.kt           # Device control commands
│       │   └── Result.kt               # Result wrapper
│       ├── src/test/                   # WiZ integration tests
│       ├── build.gradle.kts
│       └── README.md                   # Integration documentation
├── gradle/
│   └── libs.versions.toml              # Version catalog
├── build.gradle.kts                    # Root build config
├── settings.gradle.kts                 # Module configuration
└── AGENTS.md                           # AI agent guidelines
```

## 🚀 Entry Points

- **Application**: `app/src/main/kotlin/hata/HataApplication.kt` - Hilt-enabled Application class
- **Main Activity**: `app/src/main/kotlin/hata/AppActivity.kt` - Compose UI entry with EdgeToEdge support
- **Navigation**: `app/src/main/kotlin/hata/AppRouter.kt` - Nav3-based routing (Home, Mockup routes)
- **DI Configuration**: `app/src/main/kotlin/hata/di/AppModule.kt` - Hilt modules for app dependencies

## 📦 Core Modules

### Module: Main App (`app`)
- **Path**: `/app/`
- **Purpose**: Main Android application with home automation UI
- **Key Components**:
  - `HomeViewModel`: Manages home state, loads remote config, checks WiZ device statuses
  - `HomeScreen`: Main UI screen with device cards grid
  - `TopBar`: Reusable top app bar component
  - `DeviceCardsGrid`: Grid layout for device cards
  - `GenericDeviceCard`: Device card UI component

### Module: WiZ Integration (`integrations:wiz`)
- **Path**: `/integrations/wiz/`
- **Purpose**: Standalone module for controlling WiZ smart lights via UDP
- **Exports**:
  - `WizDiscovery.scan()`: Discover WiZ devices on network
  - `WizControl`: Control device (on/off, brightness, RGB, temperature)
  - `WizDevice`: Device data model
  - `Result`: Success/Error wrapper
- **Key Features**:
  - UDP broadcast discovery (port 38899)
  - JSON-based device communication
  - Coroutine-based async operations
  - RGB color control (0-255)
  - Color temperature (2200K-6500K)
  - Brightness control (10-100)

### Module: Data Layer
- **Path**: `app/src/main/kotlin/hata/data/`
- **Key Classes**:
  - `RemoteConfigurationService`: Interface for fetching remote configuration
  - `RemoteConfigurationServiceImpl`: Retrofit-based implementation (GitHub Gist)
  - `Device`, `DeviceType`, `Integration`: Core data models
  - `Home`, `Room`, `RemoteConfiguration`: Configuration models

### Module: Domain Layer
- **Path**: `app/src/main/kotlin/hata/domain/usecases/`
- **Use Cases**:
  - `LoadRemoteConfigurationUseCase`: Fetches and processes remote configuration

### Module: UI Layer
- **Path**: `app/src/main/kotlin/hata/ui/`
- **Key Components**:
  - Theme: Material3 with dynamic color support, light/dark modes
  - Components: Reusable Compose components (TopBar, DeviceCard, Grid)
  - Home: Home screen feature with ViewModel
  - Preview utilities: Preview wrappers for Compose

## 🔧 Configuration

- **settings.gradle.kts**: Defines modules (`:app`, `:integrations:wiz`)
- **build.gradle.kts**: Root build configuration with plugin aliases
- **app/build.gradle.kts**: App module config (plugins, dependencies, build variants)
- **gradle/libs.versions.toml**: Version catalog for all dependencies
- **gradle.properties**: Gradle build properties
- **proguard-rules.pro**: ProGuard/R8 configuration

## 📚 Documentation

- **AGENTS.md**: Comprehensive AI agent guidelines covering:
  - Tech stack overview (Kotlin 2.3.0, AGP 9.0.0, Compose BOM 2026.01.00)
  - Build commands (build, assemble, install, test, lint)
  - Test commands (unit, instrumented, code quality)
  - Project structure
  - Code style guidelines
  - Dependency management
  - Testing patterns
- **integrations/wiz/README.md**: WiZ integration documentation
  - Usage examples
  - API reference
  - Architecture details
  - Permissions required
  - Implementation differences from Go version

## 🧪 Test Coverage

- **Unit tests**: 5 test files
  - `ExampleUnitTest.kt`: Sample unit test
  - `WizModelsTest.kt`, `WizControlTest.kt`, `WizDiscoveryTest.kt`: WiZ integration tests
- **Instrumented tests**: 1 test file
  - `ExampleInstrumentedTest.kt`: Sample Android instrumentation test
- **Test frameworks**:
  - JUnit 4 (unit tests)
  - AndroidX Test (instrumented tests)
  - Espresso (UI testing)
  - Compose UI Test (Compose testing)

## 🔗 Key Dependencies

### Core
- **Kotlin**: 2.3.0 (with ExplicitBackingFields language feature)
- **AndroidX Core KTX**: Core Android utilities
- **AndroidX Lifecycle Runtime KTX**: Lifecycle-aware components
- **Kotlinx Serialization JSON**: JSON serialization

### Dependency Injection
- **Hilt**: v2.x (Dagger-based DI for Android)
- **Hilt Navigation Compose**: Hilt integration for Compose navigation

### Networking
- **Retrofit**: HTTP client
- **Moshi**: JSON parser with Kotlin support
- **OkHttp**: HTTP client with logging interceptor

### UI
- **Jetpack Compose**: BOM 2026.01.00
- **Material3**: Material Design 3 components
- **Material Icons Extended**: Extended icon set
- **Navigation3 Runtime & UI**: Navigation library (androidx.navigation3)
- **Lifecycle ViewModel Navigation3**: ViewModel integration for Nav3

### Build Tools
- **Android Gradle Plugin**: 9.0.0
- **Kotlin Gradle Plugin**: 2.3.0
- **KSP**: Kotlin Symbol Processing (for annotation processors)

## 📝 Quick Start

### Setup
```bash
# Clone the repository
cd /home/dector/projects/hata/android

# Build the project
./gradlew build

# Install debug build on connected device
./gradlew installDebug
```

### Run Tests
```bash
# Run all unit tests
./gradlew test

# Run specific test class
./gradlew testDebugUnitTest --tests "hata.ExampleUnitTest"

# Run instrumented tests (requires device)
./gradlew connectedDebugAndroidTest

# Run lint checks
./gradlew lintDebug
```

### Development
```bash
# Build debug APK
./gradlew assembleDebug

# Build release APK
./gradlew assembleRelease

# Clean build artifacts
./gradlew clean

# Run all checks (lint + tests)
./gradlew check
```

## 🏗️ Architecture

### Tech Stack
- **Language**: Kotlin 2.3.0
- **Build System**: Gradle 9.0.0 (AGP)
- **UI Framework**: Jetpack Compose (BOM 2026.01.00)
- **Design System**: Material Design 3
- **Navigation**: AndroidX Navigation3
- **Dependency Injection**: Hilt (Dagger)
- **Networking**: Retrofit + OkHttp + Moshi
- **Min SDK**: 24 (Android 7.0)
- **Target SDK**: 36
- **Application ID**: `space.dector.hata`

### Architecture Pattern
- **Clean Architecture**: Separation of concerns with data, domain, and UI layers
- **MVVM**: Model-View-ViewModel pattern for UI
- **Repository Pattern**: Data layer abstraction
- **Use Cases**: Single-responsibility domain operations
- **Dependency Injection**: Hilt for IoC

### Key Design Decisions
1. **Compose-First**: Pure Jetpack Compose UI (no XML layouts)
2. **EdgeToEdge**: Modern edge-to-edge display support
3. **Material3**: Latest Material Design guidelines
4. **Modular Architecture**: Separate integration modules (e.g., `:integrations:wiz`)
5. **Coroutines**: Kotlin coroutines for async operations
6. **StateFlow**: State management in ViewModels
7. **Version Catalog**: Centralized dependency management

## 📊 Statistics

- **Total Kotlin files**: 36
- **Test files**: 5
- **Modules**: 2 (app, integrations:wiz)
- **Lines of documentation**: 420+ (AGENTS.md + integrations/wiz/README.md)
- **Build variants**: Debug, Release
- **Supported SDKs**: Android 7.0 (API 24) to Android 14+ (API 36)

## 🔍 Quick Reference

### Find Files by Purpose
- **ViewModels**: `app/src/main/kotlin/hata/ui/*/ViewModel.kt`
- **UI Screens**: `app/src/main/kotlin/hata/ui/*/Screen.kt`
- **Data Models**: `app/src/main/kotlin/hata/data/models/*.kt`
- **API Services**: `app/src/main/kotlin/hata/data/api/*.kt`
- **DI Modules**: `app/src/main/kotlin/hata/di/*.kt`
- **Unit Tests**: `app/src/test/kotlin/hata/**/*Test.kt`
- **Instrumented Tests**: `app/src/androidTest/kotlin/hata/**/*Test.kt`

### Common Commands
```bash
# Full build with tests
./gradlew build

# Quick build without tests
./gradlew assemble

# Run tests with coverage
./gradlew testDebugUnitTest

# Install and run on device
./gradlew installDebug

# Lint with auto-fix
./gradlew lintFix
```

## 🎯 Current State

The project is in active development with:
- Working WiZ smart light integration with status checking
- Home screen with device cards grid
- Remote configuration loading from GitHub Gist
- Hilt dependency injection setup
- Material3 theming with dynamic colors
- Navigation3 routing between screens
- Basic unit and instrumentation test setup

Recent commits show work on:
- Moving server files to server folder
- WiZ device status updating (#vibe)
- Kotlinizing WiZ integration from Cave project
- Refactoring TopBar into components
- Adding Hilt dependency injection
