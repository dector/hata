package hata.ui.login


sealed interface LoginUiAction {

    data object Init : LoginUiAction

    data class UpdateServerUrl(val url: String) : LoginUiAction

    data class Connect(val serverUrl: String) : LoginUiAction

    data class UpdateUsername(val username: String) : LoginUiAction

    data class UpdatePassword(val password: String) : LoginUiAction

    data class Login(
        val serverUrl: String,
        val username: String,
        val password: String,
    ) : LoginUiAction

    data object ChangeServer : LoginUiAction
}
