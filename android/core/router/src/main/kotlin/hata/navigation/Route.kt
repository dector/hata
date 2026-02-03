package hata.navigation

import kotlinx.serialization.Serializable


sealed interface Route {
    @Serializable
    data object Init : Route

    @Serializable
    data object Login : Route

    @Serializable
    data object Home : Route

    @Serializable
    data object Mockup : Route

    @Serializable
    data object Profile : Route

    @Serializable
    data object Notifications : Route

    @Serializable
    data object Back : Route
}
