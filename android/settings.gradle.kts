pluginManagement {
    repositories {
        google {
            content {
                includeGroupByRegex("com\\.android.*")
                includeGroupByRegex("com\\.google.*")
                includeGroupByRegex("androidx.*")
            }
        }
        mavenCentral()
        gradlePluginPortal()
    }
}

plugins {
    id("org.gradle.toolchains.foojay-resolver-convention") version "1.0.0"
}

dependencyResolutionManagement {
    repositoriesMode.set(RepositoriesMode.FAIL_ON_PROJECT_REPOS)
    repositories {
        google()
        mavenCentral()
    }
}

rootProject.name = "Hata"

enableFeaturePreview("TYPESAFE_PROJECT_ACCESSORS")

includeBuild("gradle/plugins/build-config")

include(":core:annotations")
include(":core:api")
include(":core:ui")
include(":core:router")

include(":feature:session")
include(":feature:login")
include(":feature:profile")
include(":feature:notifications")
include(":feature:home")

include(":integrations:wiz")

include(":app")
