package hata.feature.home.ui

import hata.ui.components.DeviceCard


sealed interface HomeUiState {
    data object Init : HomeUiState
    data object Loading : HomeUiState
    data class WithData(val data: HomeDisplayData) : HomeUiState
    data class Error(val message: String) : HomeUiState
}

data class HomeDisplayData(
    val home: HomeDisplay,
    val devices: List<DeviceCard.Generic>,
)

data class HomeDisplay(
    val name: String,
)
