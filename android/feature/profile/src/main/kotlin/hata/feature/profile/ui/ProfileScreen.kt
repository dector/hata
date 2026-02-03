package hata.feature.profile.ui

import androidx.compose.foundation.background
import androidx.compose.foundation.layout.Arrangement
import androidx.compose.foundation.layout.Box
import androidx.compose.foundation.layout.Column
import androidx.compose.foundation.layout.ColumnScope
import androidx.compose.foundation.layout.Row
import androidx.compose.foundation.layout.Spacer
import androidx.compose.foundation.layout.fillMaxSize
import androidx.compose.foundation.layout.fillMaxWidth
import androidx.compose.foundation.layout.height
import androidx.compose.foundation.layout.padding
import androidx.compose.foundation.layout.size
import androidx.compose.foundation.layout.width
import androidx.compose.foundation.rememberScrollState
import androidx.compose.foundation.shape.CircleShape
import androidx.compose.foundation.shape.RoundedCornerShape
import androidx.compose.foundation.verticalScroll
import androidx.compose.material.icons.Icons
import androidx.compose.material.icons.automirrored.filled.ArrowBack
import androidx.compose.material.icons.filled.Person
import androidx.compose.material3.AlertDialog
import androidx.compose.material3.Button
import androidx.compose.material3.ButtonDefaults
import androidx.compose.material3.Card
import androidx.compose.material3.CardDefaults
import androidx.compose.material3.ExperimentalMaterial3Api
import androidx.compose.material3.Icon
import androidx.compose.material3.IconButton
import androidx.compose.material3.Scaffold
import androidx.compose.material3.Text
import androidx.compose.material3.TextButton
import androidx.compose.material3.TopAppBar
import androidx.compose.material3.TopAppBarDefaults
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
import java.time.Instant
import java.time.ZoneId
import java.time.format.DateTimeFormatter


@Composable
fun ProfileScreen(
    vm: ProfileViewModel = viewModel(),
) {
    val state by vm.uiState.collectAsStateWithLifecycle()

    LaunchedEffect(Unit) {
        vm.onDispatch(ProfileUiAction.Init)
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
    Scaffold(
        topBar = {
            TopBar(dispatch)
        },
        containerColor = HataColors.background,
    ) { paddingValues ->
        Box(
            modifier = Modifier
                .fillMaxSize()
                .padding(paddingValues),
        ) {
            when (state) {
                is ProfileUiState.Init -> {}

                is ProfileUiState.Loaded -> {
                    ProfileContent(
                        userName = state.userName,
                        serverUrl = state.serverUrl,
                        accountStatus = state.accountStatus,
                        appVersion = state.appVersion,
                        appPackage = state.appPackage,
                        sessionCreatedAt = state.sessionCreatedAt,
                        showLogoutDialog = state.showLogoutDialog,
                        onLogout = { dispatch(ProfileUiAction.Logout) },
                        onShowLogoutDialog = { dispatch(ProfileUiAction.ShowLogoutDialog) },
                        onDismissLogoutDialog = { dispatch(ProfileUiAction.DismissLogoutDialog) },
                    )
                }
            }
        }
    }
}

@Composable
@OptIn(ExperimentalMaterial3Api::class)
private fun TopBar(dispatch: (ProfileUiAction) -> Unit) {
    TopAppBar(
        title = { Text("Profile") },
        navigationIcon = {
            IconButton(onClick = { dispatch(ProfileUiAction.Back) }) {
                Icon(
                    imageVector = Icons.AutoMirrored.Filled.ArrowBack,
                    contentDescription = "Back",
                )
            }
        },
        colors = TopAppBarDefaults.topAppBarColors(
            containerColor = HataColors.background,
            titleContentColor = Color.White,
            navigationIconContentColor = Color.White,
        ),
    )
}

@Composable
private fun ProfileContent(
    userName: String,
    serverUrl: String,
    accountStatus: String,
    appVersion: String,
    appPackage: String,
    sessionCreatedAt: Long,
    showLogoutDialog: Boolean,
    onLogout: () -> Unit,
    onShowLogoutDialog: () -> Unit,
    onDismissLogoutDialog: () -> Unit,
) {
    Box(
        modifier = Modifier.fillMaxSize(),
    ) {
        Column(
            modifier = Modifier
                .fillMaxSize()
                .verticalScroll(rememberScrollState())
                .padding(start = 32.dp, end = 32.dp, top = 0.dp, bottom = 120.dp),
            horizontalAlignment = Alignment.CenterHorizontally,
        ) {
            Spacer(modifier = Modifier.height(16.dp))

            Column(
                horizontalAlignment = Alignment.CenterHorizontally,
            ) {
                UserAvatar()

                Spacer(modifier = Modifier.height(24.dp))

                UserNameDisplay(userName = userName)
            }

            Spacer(modifier = Modifier.height(48.dp))

            ProfileInfoSections(
                serverUrl = serverUrl,
                accountStatus = accountStatus,
                appVersion = appVersion,
                appPackage = appPackage,
                sessionCreatedAt = sessionCreatedAt,
            )
        }

        LogoutButton(
            onClick = onShowLogoutDialog,
            modifier = Modifier
                .align(Alignment.BottomCenter)
                .padding(start = 32.dp, end = 32.dp, bottom = 24.dp),
        )
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
private fun ProfileInfoSections(
    serverUrl: String,
    accountStatus: String,
    appVersion: String,
    appPackage: String,
    sessionCreatedAt: Long,
) {
    Column(
        modifier = Modifier.fillMaxWidth(),
    ) {
        SectionTitle(title = "User")
        InfoCard {
            InfoRow(label = "Created", value = formatSessionCreatedAt(sessionCreatedAt))
        }

        Spacer(modifier = Modifier.height(20.dp))

        SectionTitle(title = "Server")
        InfoCard {
            InfoRow(label = "URL", value = serverUrl)
            Spacer(modifier = Modifier.height(12.dp))
            StatusRow(label = "Status", status = accountStatus)
        }

        Spacer(modifier = Modifier.height(20.dp))

        SectionTitle(title = "App")
        InfoCard {
            InfoRow(label = "Package", value = appPackage)
            Spacer(modifier = Modifier.height(12.dp))
            InfoRow(label = "Version", value = appVersion)
        }
    }
}

@Composable
private fun InfoCard(content: @Composable ColumnScope.() -> Unit) {
    Card(
        modifier = Modifier.fillMaxWidth(),
        colors = CardDefaults.cardColors(
            containerColor = HataColors.surface,
        ),
        shape = RoundedCornerShape(16.dp),
    ) {
        Column(
            modifier = Modifier.padding(24.dp),
            content = content,
        )
    }
}

@Composable
private fun SectionTitle(title: String) {
    Text(
        text = title,
        color = HataColors.onSurfaceVariant,
        fontSize = 12.sp,
        fontWeight = FontWeight.SemiBold,
    )
    Spacer(modifier = Modifier.height(8.dp))
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
private fun StatusRow(label: String, status: String) {
    Row(
        modifier = Modifier.fillMaxWidth(),
        horizontalArrangement = Arrangement.SpaceBetween,
        verticalAlignment = Alignment.CenterVertically,
    ) {
        Text(
            text = label,
            color = HataColors.onSurfaceVariant,
            fontSize = 16.sp,
        )
        Row(
            verticalAlignment = Alignment.CenterVertically,
        ) {
            StatusDot(color = statusColor(status))
            Spacer(modifier = Modifier.width(8.dp))
            Text(
                text = status,
                color = Color.White,
                fontSize = 16.sp,
                fontWeight = FontWeight.Medium,
            )
        }
    }
}

@Composable
private fun StatusDot(color: Color) {
    Box(
        modifier = Modifier
            .size(8.dp)
            .clip(CircleShape)
            .background(color),
    )
}

private fun formatSessionCreatedAt(timestamp: Long): String {
    val formatter = DateTimeFormatter.ofPattern("MMM d, yyyy • HH:mm")
    val instant = Instant.ofEpochMilli(timestamp)
    return formatter.format(instant.atZone(ZoneId.systemDefault()))
}

private fun statusColor(status: String): Color {
    return when (status.lowercase()) {
        "active" -> HataColors.primary
        "inactive" -> HataColors.error
        "warning", "pending" -> Color(0xFFF2C94C)
        else -> Color(0xFFF2C94C)
    }
}

@Composable
private fun LogoutButton(
    onClick: () -> Unit,
    modifier: Modifier = Modifier,
) {
    Button(
        onClick = onClick,
        modifier = modifier
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
            appPackage = "space.dector.hata",
            sessionCreatedAt = System.currentTimeMillis(),
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
private fun Preview_ProfileInfoSections() = preview {
    ProfileInfoSections(
        serverUrl = "https://hata.local",
        accountStatus = "Active",
        appVersion = "1.0.0",
        appPackage = "space.dector.hata",
        sessionCreatedAt = System.currentTimeMillis(),
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
