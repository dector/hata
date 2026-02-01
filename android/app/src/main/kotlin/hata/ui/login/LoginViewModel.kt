package hata.ui.login

import android.os.Build
import androidx.lifecycle.ViewModel
import androidx.lifecycle.viewModelScope
import dagger.hilt.android.lifecycle.HiltViewModel
import hata.BuildConfig
import hata.core.annotations.IoDispatcher
import hata.data.api.ServerService
import hata.data.models.ServerInfo
import hata.data.models.Session
import hata.data.models.SessionUser
import hata.data.repositories.SessionRepository
import kotlinx.coroutines.CoroutineDispatcher
import kotlinx.coroutines.flow.MutableStateFlow
import kotlinx.coroutines.flow.StateFlow
import kotlinx.coroutines.flow.asStateFlow
import kotlinx.coroutines.launch
import javax.inject.Inject


@HiltViewModel
class LoginViewModel @Inject constructor(
    private val serverService: ServerService,
    private val sessionRepository: SessionRepository,
    @param:IoDispatcher
    private val dispatcher: CoroutineDispatcher,
) : ViewModel() {

    private val _uiState = MutableStateFlow<LoginUiState>(LoginUiState.Init)
    val uiState: StateFlow<LoginUiState> = _uiState.asStateFlow()

    private var currentServerInfo: ServerInfo? = null

    fun onDispatch(action: LoginUiAction) {
        when (action) {
            is LoginUiAction.Init -> handleInit()
            is LoginUiAction.UpdateServerUrl -> handleUpdateServerUrl(action.url)
            is LoginUiAction.Connect -> handleConnect(action.serverUrl)
            is LoginUiAction.UpdateUsername -> handleUpdateUsername(action.username)
            is LoginUiAction.UpdatePassword -> handleUpdatePassword(action.password)
            is LoginUiAction.Login -> handleLogin(
                action.serverUrl,
                action.username,
                action.password,
            )

            is LoginUiAction.ChangeServer -> handleChangeServer()
        }
    }

    private fun handleInit() {
        // Reset to ServerInput state (handles both first launch and post-logout navigation)
        currentServerInfo = null

        // Prefill server URL for debug builds on emulator
        val serverUrl = if (BuildConfig.DEBUG && isEmulator()) {
            "http://10.0.2.2:8080"
        } else {
            ""
        }

        _uiState.value = LoginUiState.ServerInput(serverUrl = serverUrl)
    }

    private fun isEmulator(): Boolean {
        return Build.FINGERPRINT.startsWith("generic")
            || Build.FINGERPRINT.contains("generic")
            || Build.MODEL.contains("Emulator", ignoreCase = true)
            || Build.MODEL.contains("Android SDK", ignoreCase = true)
            || Build.MODEL.contains("sdk_gphone", ignoreCase = true)
    }

    private fun handleUpdateServerUrl(url: String) {
        val currentState = _uiState.value
        if (currentState is LoginUiState.ServerInput) {
            _uiState.value = currentState.copy(serverUrl = url, error = null)
        }
    }

    private fun handleConnect(serverUrl: String) {
        viewModelScope.launch(dispatcher) {
            _uiState.value = LoginUiState.Connecting

            serverService.connect(serverUrl)
                .onSuccess { serverInfo ->
                    currentServerInfo = serverInfo
                    _uiState.value = LoginUiState.CredentialsInput(
                        serverUrl = serverInfo.serverUrl,
                        serverName = serverInfo.serverName,
                    )
                }
                .onFailure { error ->
                    _uiState.value = LoginUiState.ServerInput(
                        serverUrl = serverUrl,
                        error = error.message ?: "Connection failed",
                    )
                }
        }
    }

    private fun handleUpdateUsername(username: String) {
        val currentState = _uiState.value
        if (currentState is LoginUiState.CredentialsInput) {
            _uiState.value = currentState.copy(username = username, error = null)
        }
    }

    private fun handleUpdatePassword(password: String) {
        val currentState = _uiState.value
        if (currentState is LoginUiState.CredentialsInput) {
            _uiState.value = currentState.copy(password = password, error = null)
        }
    }

    private fun handleLogin(serverUrl: String, username: String, password: String) {
        viewModelScope.launch(dispatcher) {
            _uiState.value = LoginUiState.LoggingIn

            serverService.login(serverUrl, username, password)
                .onSuccess { token ->
                    val session = Session(
                        token = token,
                        user = SessionUser(name = username),
                        isActive = true,
                        serverUrl = serverUrl,
                    )
                    sessionRepository.saveSession(session)
                    _uiState.value = LoginUiState.Success
                }
                .onFailure { error ->
                    currentServerInfo?.let { serverInfo ->
                        _uiState.value = LoginUiState.CredentialsInput(
                            serverUrl = serverInfo.serverUrl,
                            serverName = serverInfo.serverName,
                            error = error.message ?: "Login failed",
                        )
                    } ?: run {
                        _uiState.value = LoginUiState.ServerInput(
                            serverUrl = serverUrl,
                            error = error.message ?: "Login failed",
                        )
                    }
                }
        }
    }

    private fun handleChangeServer() {
        val currentState = _uiState.value
        if (currentState is LoginUiState.CredentialsInput) {
            _uiState.value = LoginUiState.ServerInput(
                serverUrl = currentState.serverUrl,
            )
        }
    }
}
