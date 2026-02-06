package hata.feature.home.ui

import hata.feature.home.domain.NewState


sealed interface HomeUiAction {

    data object Init : HomeUiAction
    data object UpdateDevicesStatus : HomeUiAction
    data object NavigateToProfile : HomeUiAction
    data object NavigateToNotifications : HomeUiAction

    data class ToggleDevice(
        val deviceId: String,
        val newState: NewState,
    ) : HomeUiAction
}
