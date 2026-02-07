package hata.app

import android.util.Log
import androidx.compose.animation.EnterTransition
import androidx.compose.animation.ExitTransition
import androidx.compose.animation.togetherWith
import androidx.compose.runtime.Composable
import androidx.compose.runtime.LaunchedEffect
import androidx.compose.runtime.getValue
import androidx.lifecycle.compose.collectAsStateWithLifecycle
import androidx.navigation3.runtime.entryProvider
import androidx.navigation3.ui.NavDisplay
import hata.feature.home.repository.DeviceCacheCleaner
import hata.feature.login.ui.LoginScreen
import hata.feature.notifications.ui.NotificationsScreen
import hata.feature.profile.ui.ProfileScreen
import hata.feature.session.model.Session
import hata.feature.session.repository.SessionRepository
import hata.navigation.Route
import hata.navigation.internal.AppNavigator
import hata.feature.home.ui.HomeScreen
import hata.feature.home.ui.Mockup1UI
import kotlinx.coroutines.channels.consumeEach


@Composable
fun AppRouter(
    sessionRepository: SessionRepository,
    deviceCacheCleaner: DeviceCacheCleaner,
    navigator: AppNavigator,
) {
    LaunchNavigationSideEffects(navigator)
    LaunchSessionSideEffects(
        sessionRepository = sessionRepository,
        deviceCacheCleaner = deviceCacheCleaner,
        navigator = navigator,
    )

    NavDisplay(
        backStack = navigator.backStack,
        onBack = navigator::internalGoBack,
        transitionSpec = { EnterTransition.None togetherWith ExitTransition.None },
        popTransitionSpec = { EnterTransition.None togetherWith ExitTransition.None },
        entryProvider = entryProvider {
            Log.e("+++", "Router: ${navigator.backStack.lastOrNull()}")

            entry<Route.Init> {
                val session by sessionRepository
                    .observeSession()
                    .collectAsStateWithLifecycle(initialValue = null)
                InitScreen(session, navigator)
            }

            entry<Route.Login> {
                LoginScreen()
            }

            entry<Route.Home> {
                HomeScreen()
            }

            entry<Route.Profile> {
                ProfileScreen()
            }

            entry<Route.Notifications> {
                NotificationsScreen()
            }

            entry<Route.Mockup> {
                Mockup1UI()
            }

            entry<Route.Back> {
                // ignore
            }
        },
    )
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

//region Navigation

@Composable
private fun LaunchNavigationSideEffects(navigator: AppNavigator) {
    LaunchedEffect(Unit) {
        navigator.updates.consumeEach { route ->
            navigator.consume(route)
        }
    }
}

//endregion

//region Session Change

@Composable
private fun LaunchSessionSideEffects(
    sessionRepository: SessionRepository,
    deviceCacheCleaner: DeviceCacheCleaner,
    navigator: AppNavigator,
) {
    val session by sessionRepository
        .observeSession()
        .collectAsStateWithLifecycle(initialValue = null)
    LaunchedEffect(session) {
        if (session == null) {
            deviceCacheCleaner.clearCurrentUserCache()
        }
        onSessionUpdated(session, navigator)
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

//endregion
