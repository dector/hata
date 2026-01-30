plugins {
    alias(libs.plugins.android.library)
}

android {
    namespace = "hata.core.annotations"
    compileSdk {
        version = release(36)
    }

    defaultConfig {
        minSdk = 24
        testInstrumentationRunner = "androidx.test.runner.AndroidJUnitRunner"
    }

    compileOptions {
        sourceCompatibility = JavaVersion.VERSION_11
        targetCompatibility = JavaVersion.VERSION_11
    }
}

dependencies {
    // Only dependency: javax.inject for @Qualifier
    compileOnly("javax.inject:javax.inject:1")

    // Testing
    testImplementation(libs.junit4)
}
