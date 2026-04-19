package hata.feature.home.ui

import androidx.lifecycle.ViewModel
import androidx.lifecycle.viewModelScope
import dagger.hilt.android.lifecycle.HiltViewModel
import hata.feature.home.domain.LoadDefaultShoppingListUseCase
import hata.feature.home.domain.SetShoppingItemPurchasedUseCase
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
    private val setShoppingItemPurchasedUseCase: SetShoppingItemPurchasedUseCase,
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
            is ShoppingUiAction.MarkPurchased -> markItemPurchased(action.itemId)
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

    private fun markItemPurchased(itemId: String) {
        val state = _uiState.value as? ShoppingUiState.Loaded ?: return
        val currentItem = state.items.firstOrNull { it.id == itemId } ?: return
        if (currentItem.isChecked) return

        viewModelScope.launch {
            setShoppingItemPurchasedUseCase.run(itemId = itemId)
                .onSuccess { updatedItem ->
                    val latestState = _uiState.value as? ShoppingUiState.Loaded ?: return@onSuccess
                    _uiState.value = latestState.copy(
                        items = latestState.items.map { item ->
                            if (item.id == itemId) updatedItem else item
                        },
                        errorMessage = null,
                    )
                }
                .onFailure { error ->
                    val latestState = _uiState.value as? ShoppingUiState.Loaded ?: return@onFailure
                    _uiState.value = latestState.copy(
                        errorMessage = error.message ?: "Failed to mark item as purchased",
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
        val errorMessage: String? = null,
    ) : ShoppingUiState

    data class Error(
        val message: String,
    ) : ShoppingUiState
}

sealed interface ShoppingUiAction {
    data object Init : ShoppingUiAction
    data object Retry : ShoppingUiAction
    data object Back : ShoppingUiAction
    data class MarkPurchased(val itemId: String) : ShoppingUiAction
}
