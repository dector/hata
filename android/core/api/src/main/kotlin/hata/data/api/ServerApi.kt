package hata.data.api

import hata.data.api.models.ApiDevice
import hata.data.api.models.ApiHouseInfo
import hata.data.api.models.ApiShoppingItemInfo
import hata.data.api.models.ApiShoppingListInfo


interface ServerApi {
    val serverUrl: String

    suspend fun login(username: String, password: String): Result<String>
    suspend fun fetchHouse(): Result<Unit>
    suspend fun fetchHouses(): Result<List<ApiHouseInfo>>

    suspend fun fetchDevicesByHouse(
        houseId: String,
    ): Result<List<ApiDevice>>

    suspend fun fetchDevicesByUser(): Result<List<ApiDevice>>

    suspend fun setDeviceState(
        houseId: String,
        deviceId: String,
        newState: String,
    ): Result<String>

    suspend fun fetchShoppingListsByHouse(
        houseId: String,
    ): Result<List<ApiShoppingListInfo>>

    suspend fun fetchShoppingItemsByList(
        houseId: String,
        listId: String,
    ): Result<List<ApiShoppingItemInfo>>

    suspend fun fetchDefaultShoppingListItems(): Result<List<ApiShoppingItemInfo>>
}
