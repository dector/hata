package hata.ui.components

import androidx.compose.foundation.background
import androidx.compose.foundation.clickable
import androidx.compose.foundation.layout.Arrangement
import androidx.compose.foundation.layout.Box
import androidx.compose.foundation.layout.Column
import androidx.compose.foundation.layout.Row
import androidx.compose.foundation.layout.Spacer
import androidx.compose.foundation.layout.fillMaxWidth
import androidx.compose.foundation.layout.offset
import androidx.compose.foundation.layout.padding
import androidx.compose.foundation.layout.size
import androidx.compose.foundation.layout.width
import androidx.compose.foundation.shape.CircleShape
import androidx.compose.material.icons.Icons
import androidx.compose.material.icons.filled.KeyboardArrowDown
import androidx.compose.material.icons.filled.Notifications
import androidx.compose.material.icons.filled.Person
import androidx.compose.material.icons.filled.Sync
import androidx.compose.material3.Icon
import androidx.compose.material3.Text
import androidx.compose.animation.core.LinearEasing
import androidx.compose.animation.core.RepeatMode
import androidx.compose.animation.core.animateFloat
import androidx.compose.animation.core.infiniteRepeatable
import androidx.compose.animation.core.rememberInfiniteTransition
import androidx.compose.animation.core.tween
import androidx.compose.runtime.Composable
import androidx.compose.runtime.getValue
import androidx.compose.ui.Alignment
import androidx.compose.ui.Modifier
import androidx.compose.ui.draw.clip
import androidx.compose.ui.graphics.graphicsLayer
import androidx.compose.ui.graphics.Color
import androidx.compose.ui.text.font.FontWeight
import androidx.compose.ui.tooling.preview.Preview
import androidx.compose.ui.unit.dp
import androidx.compose.ui.unit.sp
import hata.ui.theme.HataAccentColors
import hata.ui.theme.HataColors
import hata.ui.utils.preview


enum class ConnectionStatus(
    val label: String,
) {
    Unknown("Unknown"),
    Online("Online"),
    Offline("Offline"),
}

@Composable
fun TopBar(
    name: String,
    userName: String? = null,
    connectionStatus: ConnectionStatus = ConnectionStatus.Unknown,
    isSyncing: Boolean = false,
    hasNotifications: Boolean = false,
    onAvatarClick: () -> Unit = {},
    onNotificationsClick: () -> Unit = {},
) {
    Row(
        modifier = Modifier
            .fillMaxWidth()
            .padding(16.dp),
        horizontalArrangement = Arrangement.SpaceBetween,
        verticalAlignment = Alignment.CenterVertically,
    ) {
        Row(
            verticalAlignment = Alignment.CenterVertically,
            horizontalArrangement = Arrangement.spacedBy(8.dp),
        ) {
            HomeDropdown(
                name = name,
                connectionStatus = connectionStatus,
            )
            SyncingIndicator(isSyncing = isSyncing)
        }

        Row(
            horizontalArrangement = Arrangement.spacedBy(12.dp),
        ) {
            NotificationsButton(
                hasNotifications = hasNotifications,
                onClick = onNotificationsClick,
            )
            Avatar(
                userName = userName,
                onClick = onAvatarClick,
            )
        }
    }
}

@Composable
private fun SyncingIndicator(isSyncing: Boolean) {
    Box(
        modifier = Modifier.size(48.dp),
        contentAlignment = Alignment.Center,
    ) {
        if (isSyncing) {
            val transition = rememberInfiniteTransition(label = "syncIconRotation")
            val rotation by transition.animateFloat(
                initialValue = 0f,
                targetValue = 360f,
                animationSpec = infiniteRepeatable(
                    animation = tween(durationMillis = 900, easing = LinearEasing),
                    repeatMode = RepeatMode.Restart,
                ),
                label = "syncIconRotationDegrees",
            )

            Icon(
                imageVector = Icons.Default.Sync,
                contentDescription = "Syncing",
                tint = Color.White,
                modifier = Modifier.graphicsLayer(rotationZ = rotation),
            )
        }
    }
}

@Composable
private fun NotificationsButton(
    hasNotifications: Boolean,
    onClick: () -> Unit = {},
) {
    Box(
        modifier = Modifier
            .size(48.dp)
            .clickable(onClick = onClick),
        contentAlignment = Alignment.Center,
    ) {
        Icon(
            imageVector = Icons.Default.Notifications,
            contentDescription = "Notifications",
            tint = Color.White,
        )
        if (hasNotifications) {
            Box(
                modifier = Modifier
                    .size(8.dp)
                    .offset(x = 8.dp, y = (-8).dp)
                    .clip(CircleShape)
                    .background(HataColors.errorVariant),
            )
        }
    }
}

@Composable
private fun HomeDropdown(
    name: String,
    connectionStatus: ConnectionStatus,
) {
    Column(
        verticalArrangement = Arrangement.spacedBy(4.dp),
    ) {
        Row(
            verticalAlignment = Alignment.CenterVertically,
        ) {
            Text(
                text = name,
                color = Color.White,
                fontSize = 24.sp,
                fontWeight = FontWeight.Bold,
            )
            Spacer(modifier = Modifier.width(8.dp))
            Icon(
                imageVector = Icons.Default.KeyboardArrowDown,
                contentDescription = "Dropdown",
                tint = Color.White,
            )
        }
        ConnectionStatusRow(connectionStatus)
    }
}

@Composable
private fun ConnectionStatusRow(connectionStatus: ConnectionStatus) {
    val indicatorColor = when (connectionStatus) {
        ConnectionStatus.Unknown -> Color(0xFFFFA726)
        ConnectionStatus.Online -> HataColors.primary
        ConnectionStatus.Offline -> HataColors.errorVariant
    }

    Row(
        verticalAlignment = Alignment.CenterVertically,
        horizontalArrangement = Arrangement.spacedBy(6.dp),
    ) {
        Box(
            modifier = Modifier
                .size(8.dp)
                .clip(CircleShape)
                .background(indicatorColor),
        )
        Text(
            text = connectionStatus.label,
            color = Color.White,
            fontSize = 12.sp,
            fontWeight = FontWeight.Medium,
        )
    }
}

@Composable
private fun Avatar(
    userName: String? = null,
    onClick: () -> Unit = {},
) {
    val avatarLetter = userName
        ?.trim()
        ?.firstOrNull { !it.isWhitespace() }
        ?.uppercaseChar()
        ?.toString()

    Box(
        modifier = Modifier
            .size(48.dp)
            .clip(CircleShape)
            .background(HataAccentColors.avatarBackground)
            .clickable(onClick = onClick),
        contentAlignment = Alignment.Center,
    ) {
        if (avatarLetter != null) {
            Text(
                text = avatarLetter,
                color = HataAccentColors.avatarIcon,
                fontSize = 20.sp,
                fontWeight = FontWeight.SemiBold,
            )
        } else {
            Icon(
                imageVector = Icons.Default.Person,
                contentDescription = "Profile",
                tint = HataAccentColors.avatarIcon,
            )
        }
    }
}

@Preview(showBackground = true)
@Composable
fun TopBarPreview() = preview {
    TopBar(
        name = "My Home",
        userName = "Dan",
        connectionStatus = ConnectionStatus.Online,
        isSyncing = true,
        hasNotifications = true,
    )
}
