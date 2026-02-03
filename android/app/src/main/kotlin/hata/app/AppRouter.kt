package hata.app

import android.util.Log
import androidx.compose.animation.EnterTransition
import androidx.compose.animation.ExitTransition
import androidx.compose.animation.togetherWith
import androidx.compose.runtime.Composable
import androidx.compose.runtime.LaunchedEffect
import androidx.compose.runtime.getValue
import androidx.lifecycle.compose.collectAsStateWithLifecycle
import androidx.navigation3.runtime.NavEntry
import androidx.navigation3.ui.NavDisplay
import hata.feature.notifications.ui.NotificationsScreen
import hata.feature.profile.ui.ProfileScreen
import hata.feature.session.model.Session
import hata.feature.session.repository.SessionRepository
import hata.navigation.Route
import hata.navigation.internal.AppNavigator
import hata.ui.home.HomeScreen
import hata.ui.home.Mockup1UI
import hata.ui.login.LoginScreen
import kotlinx.coroutines.channels.consumeEach


@Composable
fun AppRouter(
    sessionRepository: SessionRepository,
    navigator: AppNavigator,
) {
    LaunchNavigationSideEffects(navigator)
    LaunchSessionSideEffects(sessionRepository, navigator)

    NavDisplay(
        backStack = navigator.backStack,
        onBack = navigator::internalGoBack,
        transitionSpec = { EnterTransition.None togetherWith ExitTransition.None },
        popTransitionSpec = { EnterTransition.None togetherWith ExitTransition.None },
        entryProvider = { key ->
            Log.e("+++", "Router: $key")
            when (key) {
                is Route.Init -> NavEntry(key) {
                    val session by sessionRepository
                        .observeSession()
                        .collectAsStateWithLifecycle(initialValue = null)
                    InitScreen(session, navigator)
                }

                is Route.Login -> NavEntry(key) {
                    LoginScreen()
                }

                is Route.Home -> NavEntry(key) {
                    HomeScreen()
                }

                is Route.Profile -> NavEntry(key) {
                    ProfileScreen()
                }

                is Route.Notifications -> NavEntry(key) {
                    NotificationsScreen()
                }

                is Route.Mockup -> NavEntry(key) {
                    Mockup1UI()
                }

                is Route.Back -> NavEntry(key) {
                    // ignore
                }

                else -> {
                    error("Unknown route: $key")
                }
            }
        },
    )
}

@Composable
private fun LaunchSessionSideEffects(
    sessionRepository: SessionRepository,
    navigator: AppNavigator,
) {
    val session by sessionRepository
        .observeSession()
        .collectAsStateWithLifecycle(initialValue = null)
    LaunchedEffect(session) {
        onSessionUpdated(session, navigator)
    }
}

@Composable
private fun LaunchNavigationSideEffects(navigator: AppNavigator) {
    LaunchedEffect(Unit) {
        navigator.updates.consumeEach { route ->
            navigator.consume(route)
        }
    }
}

@Composable
private fun InitScreen(
    session: Session?,
    navigator: AppNavigator,
) {
    LaunchedEffect(session) {
        val newRoute = when {
            session != null -> Route.Home
            else -> Route.Login
        }
        navigator.replaceAll(newRoute)
    }
}

private fun onSessionUpdated(
    session: Session?,
    navigator: AppNavigator,
) {
    val currentRoute = navigator.currentRoute()
    if (session != null && currentRoute is Route.Login) {
        navigator.replaceAll(Route.Home)
    } else if (session == null && currentRoute !is Route.Init && currentRoute !is Route.Login) {
        // Session was cleared (logout) - navigate to Init to run all checks
        navigator.replaceAll(Route.Init)
    }
}
