plugins {
    alias(libs.plugins.hata.android.library)
    alias(libs.plugins.ksp)
    alias(libs.plugins.hilt)
}

android {
    namespace = "hata.feature.session"
}

dependencies {
    implementation(projects.core.annotations)
    implementation(projects.core.router)

    // AndroidX
    implementation(libs.androidx.core.ktx)

    // Hilt
    implementation(libs.hilt.android)
    ksp(libs.hilt.compiler)

    // Coroutines
    implementation(libs.kotlinx.coroutines.core)
}
