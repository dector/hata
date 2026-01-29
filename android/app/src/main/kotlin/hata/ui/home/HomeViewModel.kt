package hata.ui.home

import androidx.lifecycle.ViewModel
import androidx.lifecycle.viewModelScope
import hata.data.api.RemoteConfigurationServiceImpl
import hata.domain.usecases.LoadRemoteConfigurationUseCase
import kotlinx.coroutines.flow.MutableStateFlow
import kotlinx.coroutines.flow.StateFlow
import kotlinx.coroutines.flow.update
import kotlinx.coroutines.launch


class HomeViewModel(
    private val loadRemoteConfigurationUseCase: LoadRemoteConfigurationUseCase = create(),
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

private inline fun <reified T> create(): T {
    return when (T::class) {
        LoadRemoteConfigurationUseCase::class -> LoadRemoteConfigurationUseCase(
            configurationService = RemoteConfigurationServiceImpl(
                baseUrl = "https://gist.githubusercontent.com/dector/abd9d949a769dc7b0392e7f6529ad148/raw/",
            ),
        ) as T

        else -> error("Unknown type: ${T::class}")
    }
}
