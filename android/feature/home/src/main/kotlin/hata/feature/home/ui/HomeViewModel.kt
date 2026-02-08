package hata.feature.home.ui

import androidx.lifecycle.ViewModel
import androidx.lifecycle.viewModelScope
import dagger.hilt.android.lifecycle.HiltViewModel
import hata.feature.home.domain.LoadDevicesUseCase
import hata.feature.home.domain.NewState
import hata.feature.home.domain.ToggleDeviceUseCase
import hata.feature.notifications.repository.NotificationsRepository
import hata.navigation.Navigator
import hata.navigation.Route
import kotlinx.coroutines.flow.MutableStateFlow
import kotlinx.coroutines.flow.StateFlow
import kotlinx.coroutines.flow.asStateFlow
import kotlinx.coroutines.flow.collect
import kotlinx.coroutines.flow.update
import kotlinx.coroutines.launch
import javax.inject.Inject


@HiltViewModel
class HomeViewModel @Inject constructor(
    private val loadDevicesUseCase: LoadDevicesUseCase,
    private val toggleDeviceUseCase: ToggleDeviceUseCase,
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
                loadDevices()
            }

            is HomeUiAction.UpdateDevicesStatus -> loadDevices()
            is HomeUiAction.ToggleDevice -> toggleDevice(
                deviceId = action.deviceId,
                newState = action.newState,
            )

            is HomeUiAction.NavigateToProfile -> navigator.goTo(Route.Profile)
            is HomeUiAction.NavigateToNotifications -> navigator.goTo(Route.Notifications)
        }
    }

    private fun toggleDevice(deviceId: String, newState: NewState) {
        val previousState = _uiState.value as? HomeUiState.WithData ?: return
        val desiredIsOn = newState == NewState.On

        _uiState.update { state ->
            val withData = state as? HomeUiState.WithData ?: return@update state
            withData.copy(
                data = withData.data.copy(
                    devices = withData.data.devices.map { device ->
                        if (device.id == deviceId) {
                            device.copy(
                                isOn = desiredIsOn,
                                status = device.status.replacePowerState(desiredIsOn),
                            )
                        } else {
                            device
                        }
                    },
                ),
            )
        }

        viewModelScope.launch {
            toggleDeviceUseCase.run(deviceId, newState)
                .onFailure {
                    _uiState.value = previousState
                }
        }
    }

    private fun loadDevices() {
        viewModelScope.launch {
            _uiState.value = HomeUiState.Loading

            try {
                val result = loadDevicesUseCase.run()
                val devices = result.devices
                val showHouseId = devices.mapNotNull { it.houseId }.distinct().size > 1
                val connectionStatus = when (result.syncError) {
                    null -> ServerConnectionStatus.Online
                    else -> ServerConnectionStatus.Offline
                }
                val cards = devices.map { device ->
                    device.toDeviceCard(showHouseId = showHouseId)
                }
                _uiState.update {
                    HomeUiState.WithData(
                        data = HomeDisplayData(
                            home = HomeDisplay(name = "Home"),
                            devices = cards,
                            connectionStatus = connectionStatus,
                        ),
                    )
                }
            } catch (error: Exception) {
                _uiState.update {
                    HomeUiState.Error(
                        error.message ?: "Failed to load devices",
                        connectionStatus = ServerConnectionStatus.Offline,
                    )
                }
            }
        }
    }
}

private fun String.replacePowerState(isOn: Boolean): String {
    val parts = split(" | ").toMutableList()
    val stateIndex = parts.indexOfFirst {
        it.equals("on", ignoreCase = true) ||
            it.equals("off", ignoreCase = true) ||
            it.equals("unknown", ignoreCase = true)
    }
    if (stateIndex == -1) return this

    parts[stateIndex] = if (isOn) "On" else "Off"
    return parts.joinToString(" | ")
}
