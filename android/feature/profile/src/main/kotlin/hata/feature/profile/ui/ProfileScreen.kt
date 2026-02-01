package hata.feature.profile.ui

import androidx.compose.foundation.background
import androidx.compose.foundation.layout.Arrangement
import androidx.compose.foundation.layout.Box
import androidx.compose.foundation.layout.Column
import androidx.compose.foundation.layout.Row
import androidx.compose.foundation.layout.Spacer
import androidx.compose.foundation.layout.fillMaxSize
import androidx.compose.foundation.layout.fillMaxWidth
import androidx.compose.foundation.layout.height
import androidx.compose.foundation.layout.padding
import androidx.compose.foundation.layout.size
import androidx.compose.foundation.shape.CircleShape
import androidx.compose.foundation.shape.RoundedCornerShape
import androidx.compose.material.icons.Icons
import androidx.compose.material.icons.automirrored.filled.ArrowBack
import androidx.compose.material.icons.filled.Person
import androidx.compose.material3.AlertDialog
import androidx.compose.material3.Button
import androidx.compose.material3.ButtonDefaults
import androidx.compose.material3.Card
import androidx.compose.material3.CardDefaults
import androidx.compose.material3.Icon
import androidx.compose.material3.IconButton
import androidx.compose.material3.Text
import androidx.compose.material3.TextButton
import androidx.compose.runtime.Composable
import androidx.compose.runtime.LaunchedEffect
import androidx.compose.runtime.getValue
import androidx.compose.ui.Alignment
import androidx.compose.ui.Modifier
import androidx.compose.ui.draw.clip
import androidx.compose.ui.graphics.Color
import androidx.compose.ui.text.font.FontWeight
import androidx.compose.ui.tooling.preview.Preview
import androidx.compose.ui.unit.dp
import androidx.compose.ui.unit.sp
import androidx.lifecycle.compose.collectAsStateWithLifecycle
import androidx.lifecycle.viewmodel.compose.viewModel
import hata.ui.theme.HataAccentColors
import hata.ui.theme.HataColors
import hata.ui.utils.preview


@Composable
fun ProfileScreen(
    vm: ProfileViewModel = viewModel(),
    onNavigateBack: () -> Unit = {},
) {
    val state by vm.uiState.collectAsStateWithLifecycle()
    val navigationEvent by vm.navigationEvent.collectAsStateWithLifecycle()

    LaunchedEffect(Unit) {
        vm.onDispatch(ProfileUiAction.Init)
    }

    LaunchedEffect(navigationEvent) {
        when (navigationEvent) {
            is ProfileNavigationEvent.NavigateBack -> {
                vm.onNavigationEventConsumed()
                onNavigateBack()
            }

            null -> {}
        }
    }

    ProfileScreenUI(
        state = state,
        dispatch = vm::onDispatch,
    )
}

@Composable
private fun ProfileScreenUI(
    state: ProfileUiState,
    dispatch: (ProfileUiAction) -> Unit = {},
) {
    Box(
        modifier = Modifier
            .fillMaxSize()
            .background(HataColors.background),
    ) {
        when (state) {
            is ProfileUiState.Init -> {}

            is ProfileUiState.Loaded -> {
                ProfileContent(
                    userName = state.userName,
                    serverUrl = state.serverUrl,
                    accountStatus = state.accountStatus,
                    appVersion = state.appVersion,
                    showLogoutDialog = state.showLogoutDialog,
                    onLogout = { dispatch(ProfileUiAction.Logout) },
                    onShowLogoutDialog = { dispatch(ProfileUiAction.ShowLogoutDialog) },
                    onDismissLogoutDialog = { dispatch(ProfileUiAction.DismissLogoutDialog) },
                )
            }
        }

        // Back button
        IconButton(
            onClick = { dispatch(ProfileUiAction.Back) },
            modifier = Modifier
                .padding(16.dp)
                .align(Alignment.TopStart),
        ) {
            Icon(
                imageVector = Icons.AutoMirrored.Filled.ArrowBack,
                contentDescription = "Back",
                tint = Color.White,
            )
        }
    }
}

@Composable
private fun ProfileContent(
    userName: String,
    serverUrl: String,
    accountStatus: String,
    appVersion: String,
    showLogoutDialog: Boolean,
    onLogout: () -> Unit,
    onShowLogoutDialog: () -> Unit,
    onDismissLogoutDialog: () -> Unit,
) {

    Column(
        modifier = Modifier
            .fillMaxSize()
            .padding(32.dp),
        horizontalAlignment = Alignment.CenterHorizontally,
        verticalArrangement = Arrangement.SpaceBetween,
    ) {
        Spacer(modifier = Modifier.height(48.dp))

        Column(
            horizontalAlignment = Alignment.CenterHorizontally,
        ) {
            UserAvatar()

            Spacer(modifier = Modifier.height(24.dp))

            UserNameDisplay(userName = userName)
        }

        Spacer(modifier = Modifier.height(48.dp))

        ProfileInfoCard(
            serverUrl = serverUrl,
            accountStatus = accountStatus,
            appVersion = appVersion,
        )

        Spacer(modifier = Modifier.weight(1f))

        LogoutButton(onClick = onShowLogoutDialog)
    }

    // Logout Confirmation Dialog
    if (showLogoutDialog) {
        LogoutDialog(
            onConfirm = {
                onDismissLogoutDialog()
                onLogout()
            },
            onDismiss = onDismissLogoutDialog,
        )
    }
}

@Composable
private fun UserAvatar() {
    Box(
        modifier = Modifier
            .size(120.dp)
            .clip(CircleShape)
            .background(HataAccentColors.avatarBackground),
        contentAlignment = Alignment.Center,
    ) {
        Icon(
            imageVector = Icons.Default.Person,
            contentDescription = "Profile",
            tint = HataAccentColors.avatarIcon,
            modifier = Modifier.size(64.dp),
        )
    }
}

@Composable
private fun UserNameDisplay(userName: String) {
    Text(
        text = userName,
        color = Color.White,
        fontSize = 32.sp,
        fontWeight = FontWeight.Bold,
    )
}

@Composable
private fun ProfileInfoCard(
    serverUrl: String,
    accountStatus: String,
    appVersion: String,
) {
    Card(
        modifier = Modifier.fillMaxWidth(),
        colors = CardDefaults.cardColors(
            containerColor = HataColors.surface,
        ),
        shape = RoundedCornerShape(16.dp),
    ) {
        Column(
            modifier = Modifier.padding(24.dp),
        ) {
            InfoRow(label = "Server", value = serverUrl)
            Spacer(modifier = Modifier.height(16.dp))
            InfoRow(label = "Status", value = accountStatus)
            Spacer(modifier = Modifier.height(16.dp))
            InfoRow(label = "Version", value = appVersion)
        }
    }
}

@Composable
private fun InfoRow(label: String, value: String) {
    Row(
        modifier = Modifier.fillMaxWidth(),
        horizontalArrangement = Arrangement.SpaceBetween,
    ) {
        Text(
            text = label,
            color = HataColors.onSurfaceVariant,
            fontSize = 16.sp,
        )
        Text(
            text = value,
            color = Color.White,
            fontSize = 16.sp,
            fontWeight = FontWeight.Medium,
        )
    }
}

@Composable
private fun LogoutButton(onClick: () -> Unit) {
    Button(
        onClick = onClick,
        modifier = Modifier
            .fillMaxWidth()
            .height(56.dp),
        colors = ButtonDefaults.buttonColors(
            containerColor = HataColors.error,
            contentColor = HataColors.onPrimary,
        ),
        shape = RoundedCornerShape(12.dp),
    ) {
        Text(
            text = "Logout",
            fontSize = 18.sp,
            fontWeight = FontWeight.SemiBold,
        )
    }
}

@Composable
private fun LogoutDialog(
    onConfirm: () -> Unit,
    onDismiss: () -> Unit,
) {
    AlertDialog(
        onDismissRequest = onDismiss,
        title = {
            Text(
                text = "Logout",
                color = Color.White,
            )
        },
        text = {
            Text(
                text = "Are you sure you want to logout?",
                color = HataColors.onSurfaceVariant,
            )
        },
        confirmButton = {
            TextButton(
                onClick = onConfirm,
                colors = ButtonDefaults.textButtonColors(
                    contentColor = HataColors.error,
                ),
            ) {
                Text("Logout")
            }
        },
        dismissButton = {
            TextButton(
                onClick = onDismiss,
                colors = ButtonDefaults.textButtonColors(
                    contentColor = Color.White,
                ),
            ) {
                Text("Cancel")
            }
        },
        containerColor = HataColors.surface,
    )
}

@Preview
@Composable
private fun Preview_ProfileScreen() = preview {
    ProfileScreenUI(
        state = ProfileUiState.Loaded(
            userName = "Dan",
            serverUrl = "https://hata.local",
            accountStatus = "Active",
            appVersion = "1.0.0",
            showLogoutDialog = false,
        ),
    )
}

@Preview
@Composable
private fun Preview_UserAvatar() = preview {
    UserAvatar()
}

@Preview
@Composable
private fun Preview_UserNameDisplay() = preview {
    UserNameDisplay(userName = "Dan")
}

@Preview
@Composable
private fun Preview_ProfileInfoCard() = preview {
    ProfileInfoCard(
        serverUrl = "https://hata.local",
        accountStatus = "Active",
        appVersion = "1.0.0",
    )
}

@Preview
@Composable
private fun Preview_LogoutButton() = preview {
    LogoutButton(onClick = {})
}

@Preview
@Composable
private fun Preview_LogoutDialog() = preview(pad = 4.dp) {
    LogoutDialog(
        onConfirm = {},
        onDismiss = {},
    )
}
