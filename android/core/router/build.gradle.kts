plugins {
    alias(libs.plugins.hata.android.library)
    alias(libs.plugins.ksp)
    alias(libs.plugins.hilt)
}

android {
    namespace = "hata.core.router"
}

dependencies {
    implementation(libs.kotlinx.coroutines.core)
    
    // Hilt
    implementation(libs.hilt.android)
    ksp(libs.hilt.compiler)
}
