plugins {
    id("com.android.application")
}

android {
    compileSdk {
        version = release(Sdk.Compile)
    }

    defaultConfig {
        minSdk = Sdk.Min
        targetSdk = Sdk.Target
        testInstrumentationRunner = Tests.Runner
    }
}
