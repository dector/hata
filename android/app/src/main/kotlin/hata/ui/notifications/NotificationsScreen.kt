package hata.ui.notifications

import androidx.compose.foundation.background
import androidx.compose.foundation.clickable
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
import androidx.compose.foundation.lazy.LazyColumn
import androidx.compose.foundation.lazy.items
import androidx.compose.foundation.shape.CircleShape
import androidx.compose.foundation.shape.RoundedCornerShape
import androidx.compose.material.icons.Icons
import androidx.compose.material.icons.automirrored.filled.ArrowBack
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
import hata.data.models.Notification
import hata.ui.theme.HataColors
import hata.ui.utils.preview
import java.text.SimpleDateFormat
import java.util.Date
import java.util.Locale


@Composable
fun NotificationsScreen(
    vm: NotificationsViewModel = viewModel(),
    onNavigateBack: () -> Unit = {},
) {
    val state by vm.uiState.collectAsStateWithLifecycle()
    val navigationEvent by vm.navigationEvent.collectAsStateWithLifecycle()

    LaunchedEffect(Unit) {
        vm.onDispatch(NotificationsUiAction.Init)
    }

    LaunchedEffect(navigationEvent) {
        when (navigationEvent) {
            is NotificationsNavigationEvent.NavigateBack -> {
                vm.onNavigationEventConsumed()
                onNavigateBack()
            }

            null -> {}
        }
    }

    NotificationsScreenUI(
        state = state,
        dispatch = vm::onDispatch,
    )
}

@Composable
private fun NotificationsScreenUI(
    state: NotificationsUiState,
    dispatch: (NotificationsUiAction) -> Unit = {},
) {
    Box(
        modifier = Modifier
            .fillMaxSize()
            .background(HataColors.background),
    ) {
        Column(
            modifier = Modifier
                .fillMaxSize()
                .padding(top = 64.dp),
        ) {
            // Header
            NotificationsHeader(
                hasNotifications = state is NotificationsUiState.Loaded && state.notifications.isNotEmpty(),
                onMarkAllAsRead = { dispatch(NotificationsUiAction.MarkAllAsRead) },
            )

            Spacer(modifier = Modifier.height(16.dp))

            // Content
            when (state) {
                is NotificationsUiState.Init -> {}

                is NotificationsUiState.Loaded -> {
                    if (state.notifications.isEmpty()) {
                        EmptyNotifications()
                    } else {
                        NotificationsList(
                            notifications = state.notifications,
                            onNotificationClick = { notification ->
                                if (!notification.isRead) {
                                    dispatch(NotificationsUiAction.MarkAsRead(notification.id))
                                }
                            },
                        )
                    }
                }
            }
        }

        // Back button
        IconButton(
            onClick = { dispatch(NotificationsUiAction.Back) },
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
private fun NotificationsHeader(
    hasNotifications: Boolean,
    onMarkAllAsRead: () -> Unit,
) {
    Row(
        modifier = Modifier
            .fillMaxWidth()
            .padding(horizontal = 16.dp),
        horizontalArrangement = Arrangement.SpaceBetween,
        verticalAlignment = Alignment.CenterVertically,
    ) {
        Text(
            text = "Notifications",
            color = Color.White,
            fontSize = 28.sp,
            fontWeight = FontWeight.Bold,
        )

        if (hasNotifications) {
            TextButton(onClick = onMarkAllAsRead) {
                Text(
                    text = "Mark all as read",
                    color = HataColors.primary,
                    fontSize = 14.sp,
                )
            }
        }
    }
}

@Composable
private fun NotificationsList(
    notifications: List<Notification>,
    onNotificationClick: (Notification) -> Unit,
) {
    LazyColumn(
        modifier = Modifier.fillMaxSize(),
        verticalArrangement = Arrangement.spacedBy(8.dp),
    ) {
        items(notifications) { notification ->
            NotificationItem(
                notification = notification,
                onClick = { onNotificationClick(notification) },
            )
        }
        item {
            Spacer(modifier = Modifier.height(16.dp))
        }
    }
}

@Composable
private fun NotificationItem(
    notification: Notification,
    onClick: () -> Unit,
) {
    Card(
        modifier = Modifier
            .fillMaxWidth()
            .padding(horizontal = 16.dp)
            .clickable(onClick = onClick),
        colors = CardDefaults.cardColors(
            containerColor = if (notification.isRead) {
                HataColors.surface
            } else {
                HataColors.surfaceVariant
            },
        ),
        shape = RoundedCornerShape(12.dp),
    ) {
        Row(
            modifier = Modifier
                .fillMaxWidth()
                .padding(16.dp),
            horizontalArrangement = Arrangement.spacedBy(12.dp),
        ) {
            // Unread indicator
            if (!notification.isRead) {
                Box(
                    modifier = Modifier
                        .size(8.dp)
                        .clip(CircleShape)
                        .background(HataColors.primary)
                        .align(Alignment.Top)
                        .padding(top = 8.dp),
                )
            } else {
                Spacer(modifier = Modifier.size(8.dp))
            }

            Column(
                modifier = Modifier.weight(1f),
            ) {
                Text(
                    text = notification.title,
                    color = Color.White,
                    fontSize = 16.sp,
                    fontWeight = FontWeight.SemiBold,
                )
                Spacer(modifier = Modifier.height(4.dp))
                Text(
                    text = notification.message,
                    color = HataColors.onSurfaceVariant,
                    fontSize = 14.sp,
                )
                Spacer(modifier = Modifier.height(8.dp))
                Text(
                    text = formatTimestamp(notification.timestamp),
                    color = HataColors.onSurfaceVariant,
                    fontSize = 12.sp,
                )
            }
        }
    }
}

@Composable
private fun EmptyNotifications() {
    Box(
        modifier = Modifier.fillMaxSize(),
        contentAlignment = Alignment.Center,
    ) {
        Column(
            horizontalAlignment = Alignment.CenterHorizontally,
        ) {
            Text(
                text = "No notifications",
                color = HataColors.onSurfaceVariant,
                fontSize = 18.sp,
            )
            Spacer(modifier = Modifier.height(8.dp))
            Text(
                text = "You're all caught up!",
                color = HataColors.onSurfaceVariant,
                fontSize = 14.sp,
            )
        }
    }
}

private fun formatTimestamp(timestamp: Long): String {
    val now = System.currentTimeMillis()
    val diff = now - timestamp

    return when {
        diff < 60_000 -> "Just now"
        diff < 3_600_000 -> "${diff / 60_000}m ago"
        diff < 86_400_000 -> "${diff / 3_600_000}h ago"
        diff < 604_800_000 -> "${diff / 86_400_000}d ago"
        else -> {
            val dateFormat = SimpleDateFormat("MMM dd", Locale.getDefault())
            dateFormat.format(Date(timestamp))
        }
    }
}

@Preview
@Composable
private fun Preview_NotificationsScreen() = preview {
    NotificationsScreenUI(
        state = NotificationsUiState.Loaded(
            notifications = listOf(
                Notification(
                    id = "1",
                    title = "New message",
                    message = "You have a new message from John",
                    timestamp = System.currentTimeMillis() - 3600_000,
                    isRead = false,
                ),
                Notification(
                    id = "2",
                    title = "System update",
                    message = "Your system has been updated successfully",
                    timestamp = System.currentTimeMillis() - 86400_000,
                    isRead = true,
                ),
            ),
        ),
    )
}

@Preview
@Composable
private fun Preview_EmptyNotifications() = preview {
    NotificationsScreenUI(
        state = NotificationsUiState.Loaded(
            notifications = emptyList(),
        ),
    )
}

@Preview
@Composable
private fun Preview_NotificationItem() = preview(pad = 4.dp) {
    NotificationItem(
        notification = Notification(
            id = "1",
            title = "New message",
            message = "You have a new message from John",
            timestamp = System.currentTimeMillis() - 3600_000,
            isRead = false,
        ),
        onClick = {},
    )
}
