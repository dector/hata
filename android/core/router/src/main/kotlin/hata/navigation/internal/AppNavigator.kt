package hata.navigation.internal

import androidx.compose.runtime.snapshots.SnapshotStateList
import hata.navigation.Route
import kotlinx.coroutines.channels.ReceiveChannel


interface AppNavigator {

    val updates: ReceiveChannel<Route>

    val backStack: SnapshotStateList<Route>

    fun currentRoute(): Route?

    fun consume(route: Route)

    fun replaceAll(route: Route)

    fun internalGoBack() {
        if (backStack.size > 1) {
            backStack.removeLastOrNull()
        }
    }
}
