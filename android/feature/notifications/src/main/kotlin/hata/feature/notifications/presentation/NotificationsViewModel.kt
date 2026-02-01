package hata.feature.notifications.presentation

import androidx.lifecycle.ViewModel
import androidx.lifecycle.viewModelScope
import dagger.hilt.android.lifecycle.HiltViewModel
import hata.core.annotations.IoDispatcher
import hata.feature.notifications.data.model.Notification
import hata.feature.notifications.data.repository.NotificationsRepository
import kotlinx.coroutines.CoroutineDispatcher
import kotlinx.coroutines.flow.MutableStateFlow
import kotlinx.coroutines.flow.StateFlow
import kotlinx.coroutines.flow.asStateFlow
import kotlinx.coroutines.launch
import javax.inject.Inject


@HiltViewModel
class NotificationsViewModel @Inject constructor(
    private val notificationsRepository: NotificationsRepository,
    @param:IoDispatcher
    private val dispatcher: CoroutineDispatcher,
) : ViewModel() {

    private val _uiState = MutableStateFlow<NotificationsUiState>(NotificationsUiState.Init)
    val uiState: StateFlow<NotificationsUiState> = _uiState.asStateFlow()

    private val _navigationEvent = MutableStateFlow<NotificationsNavigationEvent?>(null)
    val navigationEvent: StateFlow<NotificationsNavigationEvent?> = _navigationEvent.asStateFlow()

    fun onDispatch(action: NotificationsUiAction) {
        when (action) {
            is NotificationsUiAction.Init -> handleInit()
            is NotificationsUiAction.Back -> handleBack()
            is NotificationsUiAction.MarkAsRead -> handleMarkAsRead(action.notificationId)
            is NotificationsUiAction.MarkAllAsRead -> handleMarkAllAsRead()
        }
    }

    fun onNavigationEventConsumed() {
        _navigationEvent.value = null
    }

    private fun handleInit() {
        viewModelScope.launch(dispatcher) {
            notificationsRepository.observeNotifications().collect { notifications ->
                _uiState.value = NotificationsUiState.Loaded(
                    notifications = notifications.sortedByDescending { it.timestamp },
                )
            }
        }
    }

    private fun handleBack() {
        _navigationEvent.value = NotificationsNavigationEvent.NavigateBack
    }

    private fun handleMarkAsRead(notificationId: String) {
        viewModelScope.launch(dispatcher) {
            notificationsRepository.markAsRead(notificationId)
        }
    }

    private fun handleMarkAllAsRead() {
        viewModelScope.launch(dispatcher) {
            notificationsRepository.markAllAsRead()
        }
    }
}

sealed interface NotificationsUiState {
    data object Init : NotificationsUiState

    data class Loaded(
        val notifications: List<Notification>,
    ) : NotificationsUiState
}

sealed interface NotificationsUiAction {
    data object Init : NotificationsUiAction
    data object Back : NotificationsUiAction
    data class MarkAsRead(val notificationId: String) : NotificationsUiAction
    data object MarkAllAsRead : NotificationsUiAction
}

sealed interface NotificationsNavigationEvent {
    data object NavigateBack : NotificationsNavigationEvent
}
