package hata

import androidx.compose.foundation.layout.PaddingValues
import androidx.compose.foundation.layout.padding
import androidx.compose.runtime.Composable
import androidx.compose.runtime.mutableStateListOf
import androidx.compose.runtime.remember
import androidx.compose.ui.Modifier
import androidx.navigation3.runtime.NavEntry
import androidx.navigation3.ui.NavDisplay
import hata.ui.home.HomeScreen
import hata.ui.home.Mockup1UI
import kotlinx.serialization.Serializable


@Composable
fun AppRouter(innerPadding: PaddingValues) {
    val backStack = remember { mutableStateListOf<Any>(Route.Home) }

    NavDisplay(
        modifier = Modifier.padding(innerPadding),
        backStack = backStack,
        onBack = { backStack.removeLastOrNull() },
        entryProvider = { key ->
            when (key) {
                is Route.Home -> NavEntry(key) {
                    HomeScreen()
                }

                is Route.Mockup -> NavEntry(key) {
                    Mockup1UI()
                }

                else -> {
                    error("Unknown route: $key")
                }
            }
        },
    )
}

sealed interface Route {
    @Serializable
    data object Home : Route

    @Serializable
    data object Mockup : Route
}
