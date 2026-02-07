package hata.feature.home.repository.local

import androidx.room.Database
import androidx.room.RoomDatabase


@Database(
    entities = [DeviceEntity::class],
    version = 1,
    exportSchema = false,
)
abstract class DevicesDatabase : RoomDatabase() {

    abstract fun deviceDao(): DeviceDao
}
