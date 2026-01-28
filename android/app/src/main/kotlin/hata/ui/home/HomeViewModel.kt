package hata.ui.home

import androidx.lifecycle.ViewModel
import androidx.lifecycle.viewModelScope
import hata.data.api.RemoteConfigurationService
import hata.data.api.RemoteConfigurationServiceImpl
import hata.data.models.DeviceType
import hata.ui.components.DeviceCard
import kotlinx.coroutines.flow.MutableStateFlow
import kotlinx.coroutines.flow.StateFlow
import kotlinx.coroutines.flow.asStateFlow
import kotlinx.coroutines.flow.update
import kotlinx.coroutines.launch


class HomeViewModel(
    private val configurationService: RemoteConfigurationService = RemoteConfigurationServiceImpl(
        baseUrl = "https://gist.githubusercontent.com/dector/abd9d949a769dc7b0392e7f6529ad148/raw/",
    ),
) : ViewModel() {

    private val _uiState = MutableStateFlow<HomeUiState>(HomeUiState.Loading)
    val uiState: StateFlow<HomeUiState> = _uiState.asStateFlow()

    init {
        loadConfiguration()
    }

    fun loadConfiguration() {
        viewModelScope.launch {
            _uiState.value = HomeUiState.Loading

            configurationService.loadConfiguration()
                .onSuccess { configuration ->
                    val home = configuration.homes.firstOrNull()
                    val devices = configuration.homes
                        .flatMap { it.rooms }
                        .flatMap { room ->
                            room.devices.map { device -> room to device }
                        }
                        .map { (room, device) ->
                            DeviceCard.Generic(
                                title = "[${room.name}] ${device.name}",
                                status = "Ready",
                                isOn = false,
                                icon = when (device.type) {
                                    DeviceType.Light -> DeviceCard.Icon.Light
                                    DeviceType.Unknown -> DeviceCard.Icon.Default
                                },
                            )
                        }

                    _uiState.update {
                        HomeUiState.WithData(
                            data = HomeDisplayData(
                                home = HomeDisplay(name = home?.name ?: "Home"),
                                devices = devices,
                            ),
                        )
                    }
                }
                .onFailure { error ->
                    _uiState.update {
                        HomeUiState.Error(
                            error.message ?: "Unknown error occurred",
                        )
                    }
                }
        }
    }
}
