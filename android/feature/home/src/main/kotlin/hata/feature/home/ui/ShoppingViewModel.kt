package hata.feature.home.ui

import androidx.lifecycle.ViewModel
import androidx.lifecycle.viewModelScope
import dagger.hilt.android.lifecycle.HiltViewModel
import hata.feature.home.domain.LoadDefaultShoppingListUseCase
import hata.feature.home.domain.ShoppingItem
import hata.navigation.Navigator
import kotlinx.coroutines.flow.MutableStateFlow
import kotlinx.coroutines.flow.StateFlow
import kotlinx.coroutines.flow.asStateFlow
import kotlinx.coroutines.launch
import javax.inject.Inject

@HiltViewModel
class ShoppingViewModel @Inject constructor(
    private val loadDefaultShoppingListUseCase: LoadDefaultShoppingListUseCase,
    private val navigator: Navigator,
) : ViewModel() {

    private val _uiState = MutableStateFlow<ShoppingUiState>(ShoppingUiState.Init)
    val uiState: StateFlow<ShoppingUiState> = _uiState.asStateFlow()

    fun onDispatch(action: ShoppingUiAction) {
        when (action) {
            ShoppingUiAction.Init -> {
                if (_uiState.value !is ShoppingUiState.Init) return
                loadShoppingList()
            }

            ShoppingUiAction.Retry -> loadShoppingList()
            ShoppingUiAction.Back -> navigator.goBack()
        }
    }

    private fun loadShoppingList() {
        viewModelScope.launch {
            _uiState.value = ShoppingUiState.Loading

            loadDefaultShoppingListUseCase.run()
                .onSuccess { items ->
                    _uiState.value = ShoppingUiState.Loaded(items = items)
                }
                .onFailure { error ->
                    _uiState.value = ShoppingUiState.Error(
                        message = error.message ?: "Failed to load shopping list",
                    )
                }
        }
    }
}

sealed interface ShoppingUiState {
    data object Init : ShoppingUiState
    data object Loading : ShoppingUiState

    data class Loaded(
        val items: List<ShoppingItem>,
    ) : ShoppingUiState

    data class Error(
        val message: String,
    ) : ShoppingUiState
}

sealed interface ShoppingUiAction {
    data object Init : ShoppingUiAction
    data object Retry : ShoppingUiAction
    data object Back : ShoppingUiAction
}
