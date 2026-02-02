package hata.navigation

import kotlinx.coroutines.channels.Channel
import kotlinx.coroutines.channels.ReceiveChannel
import kotlin.concurrent.atomics.AtomicBoolean
import kotlin.concurrent.atomics.ExperimentalAtomicApi


/**
 * Global navigator for app-level navigation.
 * Injected into ViewModels via Hilt to trigger navigation without callbacks or event consumption.
 */
interface AppNavigator {
    /**
     * Channel to receive navigation commands.
     * AppRouter observes this channel to handle navigation.
     */
    val navigationCommands: ReceiveChannel<NavigationCommand>

    /**
     * Navigate to the Profile screen.
     */
    fun navigateToProfile()

    /**
     * Navigate to the Notifications screen.
     */
    fun navigateToNotifications()

    /**
     * Navigate back in the navigation stack.
     */
    fun navigateBack()

    /**
     * Mark the current navigation as handled so new commands can be sent.
     */
    fun markNavigationHandled()
}

/**
 * Navigation commands that can be sent through the AppNavigator.
 */
sealed interface NavigationCommand {
    data object ToProfile : NavigationCommand
    data object ToNotifications : NavigationCommand
    data object Back : NavigationCommand
}

/**
 * Implementation of AppNavigator using Kotlin Channels.
 * Channels are ideal for one-time events like navigation - they automatically consume values.
 */
@OptIn(ExperimentalAtomicApi::class)
class GlobalAppNavigator : AppNavigator {
    private val _navigationCommands = Channel<NavigationCommand>(Channel.BUFFERED)
    override val navigationCommands: ReceiveChannel<NavigationCommand> = _navigationCommands

    private val navigationInProgress = AtomicBoolean(false)

    override fun navigateToProfile() {
        sendCommand(NavigationCommand.ToProfile)
    }

    override fun navigateToNotifications() {
        sendCommand(NavigationCommand.ToNotifications)
    }

    override fun navigateBack() {
        sendCommand(NavigationCommand.Back)
    }

    override fun markNavigationHandled() {
        navigationInProgress.store(false)
    }

    private fun sendCommand(command: NavigationCommand) {
        if (!navigationInProgress.compareAndSet(expectedValue = false, newValue = true)) {
            return
        }

        val result = _navigationCommands.trySend(command)
        if (result.isFailure) {
            navigationInProgress.store(false)
        }
    }
}
