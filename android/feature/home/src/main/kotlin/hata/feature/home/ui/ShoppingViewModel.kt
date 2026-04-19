package hata.feature.home.ui

import androidx.lifecycle.ViewModel
import androidx.lifecycle.viewModelScope
import dagger.hilt.android.lifecycle.HiltViewModel
import hata.feature.home.domain.LoadDefaultShoppingListUseCase
import hata.feature.home.domain.SetShoppingItemPurchasedUseCase
import hata.feature.home.domain.ShoppingItem
import hata.navigation.Navigator
import kotlinx.coroutines.flow.MutableSharedFlow
import kotlinx.coroutines.flow.MutableStateFlow
import kotlinx.coroutines.flow.SharedFlow
import kotlinx.coroutines.flow.StateFlow
import kotlinx.coroutines.flow.asSharedFlow
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

    private val _events = MutableSharedFlow<ShoppingUiEvent>()
    val events: SharedFlow<ShoppingUiEvent> = _events.asSharedFlow()

    fun onDispatch(action: ShoppingUiAction) {
        when (action) {
            ShoppingUiAction.Init -> {
                if (_uiState.value !is ShoppingUiState.Init) return
                loadShoppingList(showLoadingState = true)
            }

            ShoppingUiAction.Retry -> loadShoppingList(showLoadingState = true)
            ShoppingUiAction.Refresh -> refreshShoppingList()
            ShoppingUiAction.Back -> navigator.goBack()
            is ShoppingUiAction.SetChecked -> setItemChecked(
                itemId = action.itemId,
                checked = action.checked,
            )
        }
    }

    private fun refreshShoppingList() {
        val state = _uiState.value
        when (state) {
            ShoppingUiState.Loading -> return
            is ShoppingUiState.Loaded -> {
                if (state.isRefreshing) return
                loadShoppingList(showLoadingState = false)
            }
            else -> loadShoppingList(showLoadingState = true)
        }
    }

    private fun loadShoppingList(showLoadingState: Boolean) {
        viewModelScope.launch {
            val previousLoadedState = _uiState.value as? ShoppingUiState.Loaded

            if (!showLoadingState && previousLoadedState != null) {
                _uiState.value = previousLoadedState.copy(
                    isRefreshing = true,
                    errorMessage = null,
                )
            } else {
                _uiState.value = ShoppingUiState.Loading
            }

            loadDefaultShoppingListUseCase.run()
                .onSuccess { items ->
                    _uiState.value = ShoppingUiState.Loaded(items = items.sortedForDisplay())
                }
                .onFailure { error ->
                    if (!showLoadingState && previousLoadedState != null) {
                        _uiState.value = previousLoadedState.copy(
                            isRefreshing = false,
                            errorMessage = error.message ?: "Failed to refresh shopping list",
                        )
                    } else {
                        _uiState.value = ShoppingUiState.Error(
                            message = error.message ?: "Failed to load shopping list",
                        )
                    }
                }
        }
    }

    private fun setItemChecked(itemId: String, checked: Boolean) {
        val state = _uiState.value as? ShoppingUiState.Loaded ?: return
        val currentItem = state.items.firstOrNull { it.id == itemId } ?: return
        if (currentItem.isChecked == checked) return

        viewModelScope.launch {
            setShoppingItemPurchasedUseCase.run(
                itemId = itemId,
                checked = checked,
            )
                .onSuccess { updatedItem ->
                    val latestState = _uiState.value as? ShoppingUiState.Loaded ?: return@onSuccess
                    _uiState.value = latestState.copy(
                        items = latestState.items
                            .map { item -> if (item.id == itemId) updatedItem else item }
                            .sortedForDisplay(),
                        errorMessage = null,
                    )
                    _events.emit(
                        ShoppingUiEvent.ItemCheckedChanged(
                            itemId = updatedItem.id,
                            itemName = updatedItem.name,
                            checked = updatedItem.isChecked,
                        ),
                    )
                }
                .onFailure { error ->
                    val latestState = _uiState.value as? ShoppingUiState.Loaded ?: return@onFailure
                    _uiState.value = latestState.copy(
                        errorMessage = error.message ?: "Failed to update item",
                    )
                }
        }
    }
}

private fun List<ShoppingItem>.sortedForDisplay(): List<ShoppingItem> {
    return sortedWith(
        compareBy<ShoppingItem>({ it.isChecked }, { it.name.lowercase() }, { it.id }),
    )
}

sealed interface ShoppingUiState {
    data object Init : ShoppingUiState
    data object Loading : ShoppingUiState

    data class Loaded(
        val items: List<ShoppingItem>,
        val errorMessage: String? = null,
        val isRefreshing: Boolean = false,
    ) : ShoppingUiState

    data class Error(
        val message: String,
    ) : ShoppingUiState
}

sealed interface ShoppingUiAction {
    data object Init : ShoppingUiAction
    data object Retry : ShoppingUiAction
    data object Refresh : ShoppingUiAction
    data object Back : ShoppingUiAction
    data class SetChecked(
        val itemId: String,
        val checked: Boolean,
    ) : ShoppingUiAction
}

sealed interface ShoppingUiEvent {
    data class ItemCheckedChanged(
        val itemId: String,
        val itemName: String,
        val checked: Boolean,
    ) : ShoppingUiEvent
}
