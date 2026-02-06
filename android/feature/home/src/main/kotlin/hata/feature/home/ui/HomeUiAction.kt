package hata.feature.home.ui


sealed interface HomeUiAction {

    data object Init : HomeUiAction
    data object UpdateDevicesStatus : HomeUiAction
    data object NavigateToProfile : HomeUiAction
    data object NavigateToNotifications : HomeUiAction
}
