package hata.ui.home


sealed interface HomeUiAction {

    data object Init : HomeUiAction
}
