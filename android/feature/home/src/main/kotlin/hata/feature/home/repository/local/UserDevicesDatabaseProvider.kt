package hata.feature.home.repository.local

import android.content.Context
import androidx.room.Room
import dagger.hilt.android.qualifiers.ApplicationContext
import hata.feature.home.repository.DeviceCacheCleaner
import hata.feature.session.repository.SessionRepository
import kotlinx.coroutines.sync.Mutex
import kotlinx.coroutines.sync.withLock
import javax.inject.Inject
import javax.inject.Singleton


@Singleton
class UserDevicesDatabaseProvider @Inject constructor(
    @param:ApplicationContext
    private val context: Context,
    private val sessionRepository: SessionRepository,
) : DeviceCacheCleaner {

    private val lock = Mutex()

    private var currentUserId: String? = null
    private var currentDatabase: DevicesDatabase? = null

    suspend fun deviceDao(): DeviceDao {
        return lock.withLock {
            val session = sessionRepository.getSession()
                ?: throw IllegalStateException("No active session")
            val userId = session.user.name
            openDatabase(userId).deviceDao()
        }
    }

    override suspend fun clearCurrentUserCache() {
        lock.withLock {
            val userId = currentUserId ?: return
            closeCurrentDatabase()
            context.deleteDatabase(databaseName(userId))
        }
    }

    private fun openDatabase(userId: String): DevicesDatabase {
        val existingDatabase = currentDatabase
        if (existingDatabase != null && currentUserId == userId) {
            return existingDatabase
        }

        closeCurrentDatabase()

        val db = Room.databaseBuilder(
            context,
            DevicesDatabase::class.java,
            databaseName(userId),
        ).build()

        currentUserId = userId
        currentDatabase = db
        return db
    }

    private fun closeCurrentDatabase() {
        currentDatabase?.close()
        currentDatabase = null
        currentUserId = null
    }

    private fun databaseName(userId: String): String {
        val normalizedUserId = userId.map { char ->
            if (char.isLetterOrDigit() || char == '_' || char == '-') {
                char
            } else {
                '_'
            }
        }.joinToString(separator = "")
        return "data_$normalizedUserId"
    }
}
