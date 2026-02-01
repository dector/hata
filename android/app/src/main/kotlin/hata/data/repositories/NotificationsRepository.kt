package hata.data.repositories

import hata.data.models.Notification
import kotlinx.coroutines.flow.Flow


interface NotificationsRepository {
    fun observeNotifications(): Flow<List<Notification>>
    suspend fun getNotifications(): List<Notification>
    suspend fun markAsRead(notificationId: String)
    suspend fun markAllAsRead()
    suspend fun addNotification(notification: Notification)
    suspend fun hasUnreadNotifications(): Boolean
}
