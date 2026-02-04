plugins {
    alias(libs.plugins.hata.android.library)
}

android {
    namespace = "hata.data.api"

    buildFeatures {
        buildConfig = true
    }
}

dependencies {
    implementation(libs.retrofit)
    implementation(libs.retrofit.converter.moshi)

    implementation(libs.moshi)
    implementation(libs.moshi.kotlin)

    implementation(libs.okhttp)
    implementation(libs.okhttp.logging.interceptor)

    // Testing
    testImplementation(libs.junit4)
}
