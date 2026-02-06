package hata.data.api

import hata.data.api.models.ApiDevice


interface ServerApi {
    val serverUrl: String

    suspend fun login(username: String, password: String): Result<String>
    suspend fun fetchHouse(): Result<Unit>

    suspend fun fetchDevicesByHouse(
        houseId: String,
    ): Result<List<ApiDevice>>

    suspend fun fetchDevicesByUser(): Result<List<ApiDevice>>
}
