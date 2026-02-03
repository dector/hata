plugins {
    alias(libs.plugins.hata.android.library)
    alias(libs.plugins.kotlin.serialization)
    alias(libs.plugins.ksp)
    alias(libs.plugins.hilt)
}

android {
    namespace = "hata.core.router"
}

dependencies {
    implementation(libs.kotlinx.coroutines.core)
    implementation(libs.kotlinx.serialization.json)

    implementation(platform(libs.androidx.compose.bom))
    implementation(libs.androidx.compose.runtime)

    // Hilt
    implementation(libs.hilt.android)
    ksp(libs.hilt.compiler)
}
