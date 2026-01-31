package hata

import android.util.Log
import androidx.compose.foundation.layout.PaddingValues
import androidx.compose.foundation.layout.padding
import androidx.compose.runtime.Composable
import androidx.compose.runtime.LaunchedEffect
import androidx.compose.runtime.getValue
import androidx.compose.runtime.mutableStateListOf
import androidx.compose.runtime.remember
import androidx.compose.ui.Modifier
import androidx.lifecycle.compose.collectAsStateWithLifecycle
import androidx.navigation3.runtime.NavEntry
import androidx.navigation3.ui.NavDisplay
import hata.data.repositories.SessionRepository
import hata.ui.home.HomeScreen
import hata.ui.home.Mockup1UI
import hata.ui.login.LoginScreen
import kotlinx.serialization.Serializable


@Composable
fun AppRouter(
    innerPadding: PaddingValues,
    sessionRepository: SessionRepository,
) {
    val backStack = remember { mutableStateListOf<Any>(Route.Init) }
    val session by sessionRepository.observeSession().collectAsStateWithLifecycle(initialValue = null)

    // Navigate to Home when session becomes available
    LaunchedEffect(session) {
        val currentRoute = backStack.lastOrNull()
        if (session != null && currentRoute is Route.Login) {
            backStack.clear()
            backStack.add(Route.Home)
        }
    }

    NavDisplay(
        modifier = Modifier.padding(innerPadding),
        backStack = backStack,
        onBack = { backStack.removeLastOrNull() },
        entryProvider = { key ->
            Log.e("+++", "Router: $key")
            when (key) {
                is Route.Init -> NavEntry(key) {
                    LaunchedEffect(session) {
                        backStack.clear()
                        if (session != null) {
                            backStack.add(Route.Home)
                        } else {
                            backStack.add(Route.Login)
                        }
                    }
                }

                is Route.Login -> NavEntry(key) {
                    LoginScreen()
                }

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
    data object Init : Route

    @Serializable
    data object Login : Route

    @Serializable
    data object Home : Route

    @Serializable
    data object Mockup : Route
}
