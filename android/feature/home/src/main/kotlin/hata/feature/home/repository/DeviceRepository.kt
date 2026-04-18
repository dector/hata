package hata.feature.home.repository

import hata.data.models.Device
import hata.feature.home.domain.NewState
import hata.feature.home.domain.ShoppingItem


interface DeviceRepository {

    suspend fun listByHouse(houseId: String): List<Device>
    suspend fun listByUser(): List<Device>
    suspend fun syncByUser(): Result<Unit>
    suspend fun getDefaultShoppingList(): Result<List<ShoppingItem>>
    suspend fun toggleDevice(deviceId: String, newState: NewState): Result<Device>
}
