package hata.ui.profile

import androidx.lifecycle.ViewModel
import androidx.lifecycle.viewModelScope
import dagger.hilt.android.lifecycle.HiltViewModel
import hata.BuildConfig
import hata.core.annotations.IoDispatcher
import hata.data.repositories.SessionRepository
import kotlinx.coroutines.CoroutineDispatcher
import kotlinx.coroutines.flow.MutableStateFlow
import kotlinx.coroutines.flow.StateFlow
import kotlinx.coroutines.flow.asStateFlow
import kotlinx.coroutines.launch
import javax.inject.Inject


@HiltViewModel
class ProfileViewModel @Inject constructor(
    private val sessionRepository: SessionRepository,
    @param:IoDispatcher
    private val dispatcher: CoroutineDispatcher,
) : ViewModel() {

    private val _uiState = MutableStateFlow<ProfileUiState>(ProfileUiState.Init)
    val uiState: StateFlow<ProfileUiState> = _uiState.asStateFlow()

    private val _navigationEvent = MutableStateFlow<ProfileNavigationEvent?>(null)
    val navigationEvent: StateFlow<ProfileNavigationEvent?> = _navigationEvent.asStateFlow()

    fun onDispatch(action: ProfileUiAction) {
        when (action) {
            is ProfileUiAction.Init -> handleInit()
            is ProfileUiAction.Logout -> handleLogout()
            is ProfileUiAction.Back -> handleBack()
            is ProfileUiAction.ShowLogoutDialog -> handleShowLogoutDialog()
            is ProfileUiAction.DismissLogoutDialog -> handleDismissLogoutDialog()
        }
    }

    fun onNavigationEventConsumed() {
        _navigationEvent.value = null
    }

    private fun handleInit() {
        viewModelScope.launch(dispatcher) {
            val session = sessionRepository.getSession()
            if (session != null) {
                _uiState.value = ProfileUiState.Loaded(
                    userName = session.user.name,
                    serverUrl = session.serverUrl,
                    accountStatus = if (session.isActive) "Active" else "Inactive",
                    appVersion = BuildConfig.VERSION_NAME,
                )
            }
        }
    }

    private fun handleLogout() {
        viewModelScope.launch(dispatcher) {
            sessionRepository.clearSession()
        }
    }

    private fun handleBack() {
        _navigationEvent.value = ProfileNavigationEvent.NavigateBack
    }

    private fun handleShowLogoutDialog() {
        val currentState = _uiState.value
        if (currentState is ProfileUiState.Loaded) {
            _uiState.value = currentState.copy(showLogoutDialog = true)
        }
    }

    private fun handleDismissLogoutDialog() {
        val currentState = _uiState.value
        if (currentState is ProfileUiState.Loaded) {
            _uiState.value = currentState.copy(showLogoutDialog = false)
        }
    }
}

sealed interface ProfileUiState {
    data object Init : ProfileUiState

    data class Loaded(
        val userName: String,
        val serverUrl: String,
        val accountStatus: String,
        val appVersion: String,
        val showLogoutDialog: Boolean = false,
    ) : ProfileUiState
}

sealed interface ProfileUiAction {
    data object Init : ProfileUiAction
    data object Logout : ProfileUiAction
    data object Back : ProfileUiAction
    data object ShowLogoutDialog : ProfileUiAction
    data object DismissLogoutDialog : ProfileUiAction
}

sealed interface ProfileNavigationEvent {
    data object NavigateBack : ProfileNavigationEvent
}
