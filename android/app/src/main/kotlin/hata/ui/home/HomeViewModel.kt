package hata.ui.home

import androidx.lifecycle.ViewModel
import androidx.lifecycle.viewModelScope
import dagger.hilt.android.lifecycle.HiltViewModel
import hata.domain.usecases.LoadRemoteConfigurationUseCase
import kotlinx.coroutines.flow.MutableStateFlow
import kotlinx.coroutines.flow.StateFlow
import kotlinx.coroutines.flow.update
import kotlinx.coroutines.launch
import javax.inject.Inject


@HiltViewModel
class HomeViewModel @Inject constructor(
    private val loadRemoteConfigurationUseCase: LoadRemoteConfigurationUseCase,
) : ViewModel() {

    val uiState: StateFlow<HomeUiState>
        field = MutableStateFlow<HomeUiState>(HomeUiState.Init)

    fun onDispatch(action: HomeUiAction) {
        when (action) {
            is HomeUiAction.Init -> {
                if (uiState.value !is HomeUiState.Init) return

                loadConfiguration()
            }

            else -> {}
        }
    }

    fun loadConfiguration() {
        viewModelScope.launch {
            uiState.value = HomeUiState.Loading

            loadRemoteConfigurationUseCase()
                .onSuccess { data ->
                    uiState.update {
                        HomeUiState.WithData(data = data)
                    }
                }
                .onFailure { error ->
                    uiState.update {
                        HomeUiState.Error(
                            error.message ?: "Unknown error occurred",
                        )
                    }
                }
        }
    }
}
