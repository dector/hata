plugins {
    alias(libs.plugins.hata.android.library)
    alias(libs.plugins.kotlin.serialization)
}

android {
    namespace = "hata.integrations.wiz"
}

dependencies {
    // Core dependencies
    implementation(libs.kotlinx.coroutines.core)
    implementation(projects.core.annotations)
    implementation(libs.kotlinx.serialization.json)

    // Testing
    testImplementation(libs.junit4)
    testImplementation(libs.mockk)
    testImplementation(libs.kotest.assertions)
    testImplementation(libs.kotlinx.coroutines.test)
}
