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
import androidx.compose.material3.Icon
import androidx.compose.material3.Text
import androidx.compose.runtime.Composable
import androidx.compose.ui.Alignment
import androidx.compose.ui.Modifier
import androidx.compose.ui.draw.clip
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
    connectionStatus: ConnectionStatus = ConnectionStatus.Unknown,
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
        HomeDropdown(
            name = name,
            connectionStatus = connectionStatus,
        )

        Row(
            horizontalArrangement = Arrangement.spacedBy(12.dp),
        ) {
            NotificationsButton(
                hasNotifications = hasNotifications,
                onClick = onNotificationsClick,
            )
            Avatar(onClick = onAvatarClick)
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
private fun Avatar(onClick: () -> Unit = {}) {
    Box(
        modifier = Modifier
            .size(48.dp)
            .clip(CircleShape)
            .background(HataAccentColors.avatarBackground)
            .clickable(onClick = onClick),
        contentAlignment = Alignment.Center,
    ) {
        Icon(
            imageVector = Icons.Default.Person,
            contentDescription = "Profile",
            tint = HataAccentColors.avatarIcon,
        )
    }
}

@Preview(showBackground = true)
@Composable
fun TopBarPreview() = preview {
    TopBar(
        name = "My Home",
        connectionStatus = ConnectionStatus.Online,
        hasNotifications = true,
    )
}
