package hata.feature.notifications.repository

import android.content.SharedPreferences
import androidx.core.content.edit
import hata.core.annotations.NotificationsPreferences
import hata.feature.notifications.model.Notification
import kotlinx.coroutines.flow.Flow
import kotlinx.coroutines.flow.MutableStateFlow
import kotlinx.coroutines.flow.asStateFlow
import kotlinx.serialization.json.Json


// Early implementation. To be replaced.
internal class PrefsNotificationsRepository(
    @param:NotificationsPreferences
    private val sharedPreferences: SharedPreferences,
) : NotificationsRepository {

    private val notificationsFlow = MutableStateFlow(loadNotifications())

    override fun observeNotifications(): Flow<List<Notification>> {
        return notificationsFlow.asStateFlow()
    }

    override suspend fun getNotifications(): List<Notification> {
        return notificationsFlow.value
    }

    override suspend fun markAsRead(notificationId: String) {
        val updatedNotifications = notificationsFlow.value.map { notification ->
            if (notification.id == notificationId) {
                notification.copy(isRead = true)
            } else {
                notification
            }
        }
        saveNotifications(updatedNotifications)
    }

    override suspend fun markAllAsRead() {
        val updatedNotifications = notificationsFlow.value.map { notification ->
            notification.copy(isRead = true)
        }
        saveNotifications(updatedNotifications)
    }

    override suspend fun addNotification(notification: Notification) {
        val updatedNotifications = notificationsFlow.value + notification
        saveNotifications(updatedNotifications)
    }

    override suspend fun hasUnreadNotifications(): Boolean {
        return notificationsFlow.value.any { !it.isRead }
    }

    private fun saveNotifications(notifications: List<Notification>) {
        sharedPreferences.edit {
            putString(KEY_NOTIFICATIONS, Json.encodeToString(notifications))
            apply()
        }
        notificationsFlow.value = notifications
    }

    private fun loadNotifications(): List<Notification> {
        val notificationsJson = sharedPreferences.getString(KEY_NOTIFICATIONS, null)
        return if (notificationsJson != null) {
            try {
                Json.decodeFromString(notificationsJson)
            } catch (e: Exception) {
                emptyList()
            }
        } else {
            emptyList()
        }
    }

    private companion object {
        const val KEY_NOTIFICATIONS = "notifications"
    }
}
