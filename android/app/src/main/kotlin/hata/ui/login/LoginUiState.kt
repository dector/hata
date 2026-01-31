package hata.ui.login


sealed interface LoginUiState {
    data object Init : LoginUiState

    data class ServerInput(
        val serverUrl: String = "",
        val error: String? = null,
    ) : LoginUiState {
        val isConnectEnabled: Boolean
            get() = serverUrl.isNotBlank()
    }

    data object Connecting : LoginUiState

    data class CredentialsInput(
        val serverUrl: String,
        val serverName: String,
        val username: String = "",
        val password: String = "",
        val error: String? = null,
    ) : LoginUiState {
        val isLoginEnabled: Boolean
            get() = username.isNotBlank() && password.isNotBlank()
    }

    data object LoggingIn : LoginUiState

    data object Success : LoginUiState
}
