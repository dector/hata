plugins {
    alias(libs.plugins.hata.android.library)
}

android {
    namespace = "hata.core.annotations"
}

dependencies {
    compileOnly(libs.javax.inject)

    // Testing
    testImplementation(libs.junit4)
}
