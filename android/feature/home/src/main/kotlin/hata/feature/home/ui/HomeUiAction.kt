package hata.feature.home.ui

import hata.feature.home.domain.NewState


sealed interface HomeUiAction {

    data object Init : HomeUiAction
    data class UpdateDevicesStatus(
        val showSyncFeedback: Boolean = false,
    ) : HomeUiAction
    data object NavigateToProfile : HomeUiAction
    data object NavigateToNotifications : HomeUiAction
    data object NavigateToShopping : HomeUiAction

    data class ToggleDevice(
        val deviceId: String,
        val newState: NewState,
    ) : HomeUiAction
}
