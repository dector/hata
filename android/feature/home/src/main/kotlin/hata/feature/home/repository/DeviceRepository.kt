package hata.feature.home.repository

import hata.data.models.Device


interface DeviceRepository {

    suspend fun listByHouse(houseId: String): List<Device>
    suspend fun listByUser(): List<Device>
}
