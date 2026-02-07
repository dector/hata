package hata.feature.home.repository.local

import androidx.room.Entity
import androidx.room.PrimaryKey


@Entity(tableName = "devices")
data class DeviceEntity(
    @PrimaryKey
    val id: String,
    val name: String,
    val integrationId: String,
    val integrationDataJson: String?,
    val state: String,
    val houseId: String?,
)
