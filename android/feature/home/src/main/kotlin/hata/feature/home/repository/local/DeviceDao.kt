package hata.feature.home.repository.local

import androidx.room.Dao
import androidx.room.Insert
import androidx.room.OnConflictStrategy
import androidx.room.Query
import androidx.room.Transaction


@Dao
interface DeviceDao {

    @Query("SELECT * FROM devices ORDER BY name")
    suspend fun listByUser(): List<DeviceEntity>

    @Query("SELECT * FROM devices WHERE houseId = :houseId ORDER BY name")
    suspend fun listByHouse(houseId: String): List<DeviceEntity>

    @Query("SELECT * FROM devices WHERE id = :deviceId LIMIT 1")
    suspend fun findById(deviceId: String): DeviceEntity?

    @Insert(onConflict = OnConflictStrategy.REPLACE)
    suspend fun upsert(device: DeviceEntity)

    @Insert(onConflict = OnConflictStrategy.REPLACE)
    suspend fun upsertAll(devices: List<DeviceEntity>)

    @Query("DELETE FROM devices")
    suspend fun clearAll()

    @Transaction
    suspend fun replaceAll(devices: List<DeviceEntity>) {
        clearAll()
        if (devices.isNotEmpty()) {
            upsertAll(devices)
        }
    }
}
