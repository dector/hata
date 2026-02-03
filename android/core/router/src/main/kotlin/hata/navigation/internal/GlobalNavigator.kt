package hata.navigation.internal

import androidx.compose.runtime.mutableStateListOf
import hata.navigation.Navigator
import hata.navigation.Route
import kotlinx.coroutines.channels.Channel
import kotlinx.coroutines.channels.ReceiveChannel
import kotlin.concurrent.atomics.AtomicBoolean
import kotlin.concurrent.atomics.ExperimentalAtomicApi


@OptIn(ExperimentalAtomicApi::class)
class GlobalNavigator : AppNavigator {
    private val navigationInProgress = AtomicBoolean(false)

    private val _updates = Channel<Route>(Channel.BUFFERED)
    override val updates: ReceiveChannel<Route> = _updates

    override val backStack = mutableStateListOf<Route>(Route.Init)

    override fun currentRoute(): Route? {
        return backStack.lastOrNull()
    }

    override fun consume(route: Route) {
        when (route) {
            is Route.Back -> internalGoBack()
            else -> backStack.add(route)
        }
        onNavigationConsumed()
    }

    override fun replaceAll(route: Route) {
        backStack.clear()
        backStack.add(route)
    }

    override fun internalGoBack() {
        if (backStack.size <= 1) return

        backStack.removeLastOrNull()
    }

    private fun onNavigationConsumed() {
        navigationInProgress.store(false)
    }

    // ---

    val exposed = object : Navigator {

        override fun goTo(route: Route) {
            if (!navigationInProgress.compareAndSet(expectedValue = false, newValue = true)) {
                return
            }

            val result = _updates.trySend(route)
            if (result.isFailure) {
                navigationInProgress.store(false)
            }
        }

        override fun goBack() {
            goTo(Route.Back)
        }
    }
}
