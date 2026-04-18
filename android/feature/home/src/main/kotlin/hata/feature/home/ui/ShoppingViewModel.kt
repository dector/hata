package hata.feature.home.ui

import androidx.lifecycle.ViewModel
import dagger.hilt.android.lifecycle.HiltViewModel
import hata.navigation.Navigator
import javax.inject.Inject

@HiltViewModel
class ShoppingViewModel @Inject constructor(
    private val navigator: Navigator,
) : ViewModel() {

    fun onDispatch(action: ShoppingUiAction) {
        when (action) {
            ShoppingUiAction.Back -> navigator.goBack()
        }
    }
}

sealed interface ShoppingUiAction {
    data object Back : ShoppingUiAction
}
