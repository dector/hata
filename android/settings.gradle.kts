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

include(":integrations:wiz")

include(":app")
