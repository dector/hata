package hata.feature.notifications.repository

import hata.feature.notifications.model.Notification
import kotlinx.coroutines.flow.Flow


interface NotificationsRepository {
    fun observeNotifications(): Flow<List<Notification>>
    suspend fun getNotifications(): List<Notification>
    suspend fun markAsRead(notificationId: String)
    suspend fun markAllAsRead()
    suspend fun addNotification(notification: Notification)
    suspend fun hasUnreadNotifications(): Boolean
}
