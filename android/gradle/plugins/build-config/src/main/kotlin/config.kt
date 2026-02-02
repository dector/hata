@file:Suppress("ConstPropertyName")

import org.gradle.api.JavaVersion


object Sdk {
    const val Compile = 36
    const val Min = 24
    const val Target = 36
}

object Java {
    val Version = JavaVersion.VERSION_11
}

object Tests {
    const val Runner = "androidx.test.runner.AndroidJUnitRunner"
}
