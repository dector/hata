package hata.feature.home.ui

import androidx.lifecycle.ViewModel
import androidx.lifecycle.viewModelScope
import dagger.hilt.android.lifecycle.HiltViewModel
import hata.feature.home.domain.LoadDevicesUseCase
import hata.feature.notifications.repository.NotificationsRepository
import hata.navigation.Navigator
import hata.navigation.Route
import kotlinx.coroutines.flow.MutableStateFlow
import kotlinx.coroutines.flow.StateFlow
import kotlinx.coroutines.flow.asStateFlow
import kotlinx.coroutines.flow.update
import kotlinx.coroutines.launch
import javax.inject.Inject


@HiltViewModel
class HomeViewModel @Inject constructor(
    private val loadDevicesUseCase: LoadDevicesUseCase,
    private val notificationsRepository: NotificationsRepository,
    private val navigator: Navigator,
) : ViewModel() {

    private val _uiState = MutableStateFlow<HomeUiState>(HomeUiState.Init)
    val uiState: StateFlow<HomeUiState> = _uiState.asStateFlow()

    private val _hasUnreadNotifications = MutableStateFlow(false)
    val hasUnreadNotifications: StateFlow<Boolean> = _hasUnreadNotifications.asStateFlow()

    private fun observeNotifications() {
        viewModelScope.launch {
            notificationsRepository.observeNotifications().collect { notifications ->
                _hasUnreadNotifications.value = notifications.any { !it.isRead }
            }
        }
    }

    fun onDispatch(action: HomeUiAction) {
        when (action) {
            is HomeUiAction.Init -> {
                if (_uiState.value !is HomeUiState.Init) return

                observeNotifications()
                onDispatch(HomeUiAction.UpdateDevicesStatus)
            }

            is HomeUiAction.UpdateDevicesStatus -> loadDevices()

            is HomeUiAction.NavigateToProfile -> navigator.goTo(Route.Profile)
            is HomeUiAction.NavigateToNotifications -> navigator.goTo(Route.Notifications)
        }
    }

    private fun loadDevices() {
        viewModelScope.launch {
            _uiState.value = HomeUiState.Loading

            try {
                val devices = loadDevicesUseCase.run()
                val showHouseId = devices.mapNotNull { it.houseId }.distinct().size > 1
                val cards = devices.map { device ->
                    device.toDeviceCard(showHouseId = showHouseId)
                }
                _uiState.update {
                    HomeUiState.WithData(
                        data = HomeDisplayData(
                            home = HomeDisplay(name = "Home"),
                            devices = cards,
                        ),
                    )
                }
            } catch (error: Exception) {
                _uiState.update {
                    HomeUiState.Error(
                        error.message ?: "Failed to load devices",
                    )
                }
            }
        }
    }
}
