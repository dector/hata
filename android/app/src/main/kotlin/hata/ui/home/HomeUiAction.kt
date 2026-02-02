package hata.ui.home


sealed interface HomeUiAction {

    data object Init : HomeUiAction
    data object NavigateToProfile : HomeUiAction
    data object NavigateToNotifications : HomeUiAction
}
