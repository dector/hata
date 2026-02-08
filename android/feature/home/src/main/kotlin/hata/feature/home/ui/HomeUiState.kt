package hata.feature.home.ui

import hata.ui.components.DeviceCard


sealed interface HomeUiState {
    data object Init : HomeUiState
    data object Loading : HomeUiState
    data class WithData(
        val data: HomeDisplayData,
        val isSyncing: Boolean = false,
    ) : HomeUiState
    data class Error(
        val message: String,
        val connectionStatus: ServerConnectionStatus,
    ) : HomeUiState
}

data class HomeDisplayData(
    val home: HomeDisplay,
    val devices: List<DeviceCard.Generic>,
    val connectionStatus: ServerConnectionStatus,
)

data class HomeDisplay(
    val name: String,
)

enum class ServerConnectionStatus {
    Unknown,
    Online,
    Offline,
}
