package hata.ui.home

import android.util.Log
import androidx.lifecycle.ViewModel
import androidx.lifecycle.viewModelScope
import dagger.hilt.android.lifecycle.HiltViewModel
import hata.domain.usecases.LoadRemoteConfigurationUseCase
import hata.integrations.wiz.Result
import hata.integrations.wiz.WizControl
import hata.integrations.wiz.WizDevice
import kotlinx.coroutines.flow.MutableStateFlow
import kotlinx.coroutines.flow.StateFlow
import kotlinx.coroutines.flow.update
import kotlinx.coroutines.launch
import javax.inject.Inject

private const val TAG = "HomeViewModel"

@HiltViewModel
class HomeViewModel @Inject constructor(
    private val loadRemoteConfigurationUseCase: LoadRemoteConfigurationUseCase,
) : ViewModel() {

    val uiState: StateFlow<HomeUiState>
        field = MutableStateFlow<HomeUiState>(HomeUiState.Init)

    fun onDispatch(action: HomeUiAction) {
        when (action) {
            is HomeUiAction.Init -> {
                if (uiState.value !is HomeUiState.Init) return

                loadConfiguration()
            }

            else -> {}
        }
    }

    fun loadConfiguration() {
        viewModelScope.launch {
            uiState.value = HomeUiState.Loading

            loadRemoteConfigurationUseCase()
                .onSuccess { data ->
                    uiState.update {
                        HomeUiState.WithData(data = data)
                    }

                    // Check Wiz device statuses after loading data
                    checkWizDeviceStatuses()
                }
                .onFailure { error ->
                    uiState.update {
                        HomeUiState.Error(
                            error.message ?: "Unknown error occurred",
                        )
                    }
                }
        }
    }

    private fun checkWizDeviceStatuses() {
        viewModelScope.launch {
            // Get devices from the current state
            val currentState = uiState.value
            if (currentState !is HomeUiState.WithData) {
                Log.w(TAG, "Cannot check Wiz devices: state is not WithData")
                return@launch
            }

            Log.d(TAG, "Checking Wiz device statuses...")

            // Note: We need access to the original RemoteConfiguration to get integration details
            // For now, we'll parse the device titles to extract IPs
            // Format is "[RoomName] DeviceName" with status "[IntegrationName] IP"
            currentState.data.devices.forEachIndexed { index, deviceCard ->
                // Check if this is a Wiz device by looking at the status field
                if (deviceCard.status.startsWith("[Wiz]")) {
                    val ip = deviceCard.status.substringAfter("[Wiz] ").trim()

                    if (ip.isNotEmpty()) {
                        checkWizDeviceStatus(index, deviceCard.title, ip)
                    }
                }
            }
        }
    }

    private suspend fun checkWizDeviceStatus(deviceIndex: Int, deviceName: String, ip: String) {
        Log.d(TAG, "Checking status for device '$deviceName' at IP: $ip")

        val wizDevice = WizDevice(
            name = deviceName,
            ip = ip,
            type = "Light",
            mac = "unknown"
        )

        val control = WizControl(wizDevice)

        when (val result = control.getState()) {
            is Result.Success -> {
                val state = result.data
                Log.i(TAG, """
                    Wiz Device Status - $deviceName ($ip):
                      Power: ${if (state.state) "ON" else "OFF"}
                      Brightness: ${state.dimming}%
                      Color Temp: ${state.temp}K
                      RGB: (${state.r}, ${state.g}, ${state.b})
                      Scene: ${state.sceneId}
                """.trimIndent())

                // Update the device state in uiState
                updateDeviceStatus(deviceIndex, state.state)
            }
            is Result.Error -> {
                Log.e(TAG, "Failed to get status for device '$deviceName' at $ip: ${result.exception.message}")
            }
        }
    }

    private fun updateDeviceStatus(deviceIndex: Int, isOn: Boolean) {
        uiState.update { currentState ->
            if (currentState is HomeUiState.WithData) {
                val updatedDevices = currentState.data.devices.toMutableList()
                val device = updatedDevices[deviceIndex]
                updatedDevices[deviceIndex] = device.copy(isOn = isOn)

                HomeUiState.WithData(
                    data = currentState.data.copy(devices = updatedDevices)
                )
            } else {
                currentState
            }
        }
    }
}
