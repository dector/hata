plugins {
    alias(libs.plugins.hata.android.library)
    alias(libs.plugins.compose)
}

android {
    namespace = "hata.core.ui"

    buildFeatures {
        compose = true
    }
}

dependencies {
    // Compose dependencies
    implementation(platform(libs.androidx.compose.bom))
    implementation(libs.androidx.compose.ui)
    implementation(libs.androidx.compose.ui.graphics)
    implementation(libs.androidx.compose.ui.tooling.preview)
    implementation(libs.androidx.compose.material3)

    debugImplementation(libs.androidx.compose.ui.tooling)

    // Testing
    testImplementation(libs.junit4)
}
