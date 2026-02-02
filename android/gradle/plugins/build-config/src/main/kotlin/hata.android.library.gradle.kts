plugins {
    id("com.android.library")
}

android {
    compileSdk {
        version = release(Sdk.Compile)
    }

    defaultConfig {
        minSdk = Sdk.Min
        testInstrumentationRunner = Tests.Runner
    }

    compileOptions {
        sourceCompatibility = Java.Version
        targetCompatibility = Java.Version
    }
}
