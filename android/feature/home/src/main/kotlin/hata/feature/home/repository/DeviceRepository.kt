package hata.feature.home.repository

import hata.data.models.Device
import hata.feature.home.domain.NewState


interface DeviceRepository {

    suspend fun listByHouse(houseId: String): List<Device>
    suspend fun listByUser(): List<Device>
    suspend fun syncByUser(): Result<Unit>
    suspend fun toggleDevice(deviceId: String, newState: NewState): Result<Device>
}
