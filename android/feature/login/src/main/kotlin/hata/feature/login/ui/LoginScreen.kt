package hata.feature.login.ui

import androidx.compose.foundation.background
import androidx.compose.foundation.layout.Arrangement
import androidx.compose.foundation.layout.Box
import androidx.compose.foundation.layout.Column
import androidx.compose.foundation.layout.Spacer
import androidx.compose.foundation.layout.fillMaxSize
import androidx.compose.foundation.layout.fillMaxWidth
import androidx.compose.foundation.layout.height
import androidx.compose.foundation.layout.padding
import androidx.compose.foundation.shape.RoundedCornerShape
import androidx.compose.material3.Button
import androidx.compose.material3.ButtonDefaults
import androidx.compose.material3.CircularProgressIndicator
import androidx.compose.material3.OutlinedTextField
import androidx.compose.material3.OutlinedTextFieldDefaults
import androidx.compose.material3.Text
import androidx.compose.material3.TextButton
import androidx.compose.runtime.Composable
import androidx.compose.runtime.LaunchedEffect
import androidx.compose.runtime.getValue
import androidx.compose.ui.Alignment
import androidx.compose.ui.Modifier
import androidx.compose.ui.graphics.Color
import androidx.compose.ui.text.font.FontWeight
import androidx.compose.ui.text.input.PasswordVisualTransformation
import androidx.compose.ui.text.input.VisualTransformation
import androidx.compose.ui.text.style.TextAlign
import androidx.compose.ui.tooling.preview.Preview
import androidx.compose.ui.unit.dp
import androidx.compose.ui.unit.sp
import androidx.lifecycle.compose.collectAsStateWithLifecycle
import androidx.lifecycle.viewmodel.compose.viewModel
import hata.ui.theme.HataColors
import hata.ui.utils.preview


@Composable
fun LoginScreen(
    vm: LoginViewModel = viewModel(),
) {
    val state by vm.uiState.collectAsStateWithLifecycle()

    LaunchedEffect(Unit) {
        vm.onDispatch(LoginUiAction.Init)
    }

    LoginScreenUI(
        state = state,
        dispatch = vm::onDispatch,
    )
}

@Composable
private fun LoginScreenUI(
    state: LoginUiState,
    dispatch: (LoginUiAction) -> Unit = {},
) {
    Box(
        modifier = Modifier
            .fillMaxSize()
            .background(HataColors.background),
        contentAlignment = Alignment.Center,
    ) {
        when (state) {
            is LoginUiState.Init -> {}

            is LoginUiState.ServerInput -> {
                ServerInputStep(
                    serverUrl = state.serverUrl,
                    error = state.error,
                    isConnectEnabled = state.isConnectEnabled,
                    onServerUrlChange = { dispatch(LoginUiAction.UpdateServerUrl(it)) },
                    onConnect = { dispatch(LoginUiAction.Connect(state.serverUrl)) },
                )
            }

            is LoginUiState.Connecting -> {
                LoadingIndicator(message = "Connecting to server...")
            }

            is LoginUiState.CredentialsInput -> {
                CredentialsInputStep(
                    serverName = state.serverName,
                    username = state.username,
                    password = state.password,
                    error = state.error,
                    isLoginEnabled = state.isLoginEnabled,
                    onUsernameChange = { dispatch(LoginUiAction.UpdateUsername(it)) },
                    onPasswordChange = { dispatch(LoginUiAction.UpdatePassword(it)) },
                    onLogin = {
                        dispatch(
                            LoginUiAction.Login(
                                serverUrl = state.serverUrl,
                                username = state.username,
                                password = state.password,
                            ),
                        )
                    },
                    onChangeServer = { dispatch(LoginUiAction.ChangeServer) },
                )
            }

            is LoginUiState.LoggingIn -> {
                LoadingIndicator(message = "Logging in...")
            }

            is LoginUiState.Success -> {}
        }
    }
}

@Composable
private fun ServerInputStep(
    serverUrl: String,
    error: String?,
    isConnectEnabled: Boolean,
    onServerUrlChange: (String) -> Unit,
    onConnect: () -> Unit,
) {
    Column(
        modifier = Modifier
            .fillMaxWidth()
            .padding(horizontal = 32.dp),
        horizontalAlignment = Alignment.CenterHorizontally,
    ) {
        LoginTextField(
            value = serverUrl,
            onValueChange = onServerUrlChange,
            label = "Server URL",
            modifier = Modifier.fillMaxWidth(),
        )

        if (error != null) {
            Spacer(modifier = Modifier.height(8.dp))
            ErrorMessage(error)
        }

        Spacer(modifier = Modifier.height(24.dp))

        LoginButton(
            onClick = onConnect,
            text = "Connect",
            enabled = isConnectEnabled,
            modifier = Modifier
                .fillMaxWidth()
                .height(56.dp),
        )
    }
}

@Composable
private fun CredentialsInputStep(
    serverName: String,
    username: String,
    password: String,
    error: String?,
    isLoginEnabled: Boolean,
    onUsernameChange: (String) -> Unit,
    onPasswordChange: (String) -> Unit,
    onLogin: () -> Unit,
    onChangeServer: () -> Unit,
) {
    Column(
        modifier = Modifier
            .fillMaxWidth()
            .padding(horizontal = 32.dp),
        horizontalAlignment = Alignment.CenterHorizontally,
    ) {
        Header(serverName = serverName)

        Spacer(modifier = Modifier.height(48.dp))

        LoginTextField(
            value = username,
            onValueChange = onUsernameChange,
            label = "Username",
            modifier = Modifier.fillMaxWidth(),
        )

        Spacer(modifier = Modifier.height(16.dp))

        LoginTextField(
            value = password,
            onValueChange = onPasswordChange,
            label = "Password",
            modifier = Modifier.fillMaxWidth(),
            visualTransformation = PasswordVisualTransformation(),
        )

        if (error != null) {
            Spacer(modifier = Modifier.height(8.dp))
            ErrorMessage(error)
        }

        Spacer(modifier = Modifier.height(32.dp))

        LoginButton(
            onClick = onLogin,
            text = "Login",
            enabled = isLoginEnabled,
            modifier = Modifier
                .fillMaxWidth()
                .height(56.dp),
        )

        Spacer(modifier = Modifier.height(16.dp))

        TextButton(onClick = onChangeServer) {
            Text(
                text = "Change Server",
                color = Color.White.copy(alpha = 0.7f),
            )
        }
    }
}

@Composable
private fun Header(serverName: String) {
    Column(
        horizontalAlignment = Alignment.CenterHorizontally,
    ) {
        Text(
            text = "Welcome to",
            color = Color.White.copy(alpha = 0.7f),
            fontSize = 20.sp,
        )

        Spacer(modifier = Modifier.height(8.dp))

        Text(
            text = serverName,
            color = Color.White,
            lineHeight = 48.sp,
            fontSize = 48.sp,
            fontWeight = FontWeight.Bold,
            textAlign = TextAlign.Center,
        )
    }
}

@Composable
private fun LoadingIndicator(message: String) {
    Column(
        horizontalAlignment = Alignment.CenterHorizontally,
        verticalArrangement = Arrangement.Center,
    ) {
        CircularProgressIndicator(
            color = HataColors.primary,
        )
        Spacer(modifier = Modifier.height(16.dp))
        Text(
            text = message,
            color = Color.White.copy(alpha = 0.7f),
            fontSize = 16.sp,
        )
    }
}

@Composable
private fun ErrorMessage(message: String) {
    Text(
        text = message,
        color = HataColors.error,
        fontSize = 14.sp,
    )
}

@Composable
private fun LoginTextField(
    modifier: Modifier = Modifier,
    value: String,
    onValueChange: (String) -> Unit,
    label: String,
    visualTransformation: VisualTransformation = VisualTransformation.None,
) {
    OutlinedTextField(
        value = value,
        onValueChange = onValueChange,
        label = { Text(label) },
        modifier = modifier,
        singleLine = true,
        visualTransformation = visualTransformation,
        colors = OutlinedTextFieldDefaults.colors(
            focusedTextColor = Color.White,
            unfocusedTextColor = Color.White,
            focusedBorderColor = HataColors.primary,
            unfocusedBorderColor = HataColors.onSurfaceVariant,
            focusedLabelColor = HataColors.primary,
            unfocusedLabelColor = HataColors.onSurfaceVariant,
            cursorColor = HataColors.primary,
        ),
    )
}

@Composable
private fun LoginButton(
    modifier: Modifier = Modifier,
    onClick: () -> Unit,
    text: String,
    enabled: Boolean,
) {
    Button(
        onClick = onClick,
        modifier = modifier,
        colors = ButtonDefaults.buttonColors(
            containerColor = HataColors.primary,
            contentColor = HataColors.onPrimary,
        ),
        shape = RoundedCornerShape(12.dp),
        enabled = enabled,
    ) {
        Text(
            text = text,
            fontSize = 18.sp,
            fontWeight = FontWeight.SemiBold,
        )
    }
}

@Preview
@Composable
private fun Preview_LoginScreen_ServerInput() = preview {
    LoginScreenUI(
        state = LoginUiState.ServerInput(),
    )
}

@Preview
@Composable
private fun Preview_LoginScreen_CredentialsInput() = preview {
    LoginScreenUI(
        state = LoginUiState.CredentialsInput(
            serverUrl = "http://localhost",
            serverName = "Hata Dev Server",
        ),
    )
}
