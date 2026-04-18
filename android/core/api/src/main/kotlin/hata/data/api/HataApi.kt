package hata.data.api

import hata.data.api.models.ApiDeviceListResponse
import hata.data.api.models.ApiDeviceSetStateRequest
import hata.data.api.models.ApiDeviceSetStateResponse
import hata.data.api.models.ApiHouseListResponse
import hata.data.api.models.ApiShoppingItemListResponse
import hata.data.api.models.ApiShoppingListListResponse
import hata.data.api.models.LoginRequest
import hata.data.api.models.LoginResponse
import hata.data.api.models.PingResponse
import retrofit2.http.Body
import retrofit2.http.GET
import retrofit2.http.PATCH
import retrofit2.http.POST
import retrofit2.http.Path


interface HataApi {

    @GET("/api/latest/ping")
    suspend fun ping(): PingResponse

    @POST("/api/latest/auth/login")
    suspend fun login(@Body request: LoginRequest): LoginResponse

    @GET("/api/latest/house")
    suspend fun latestHouse(): ApiHouseListResponse

    @GET("/api/latest/house/{houseId}/device")
    suspend fun devicesByHouse(
        @Path("houseId") houseId: String,
    ): ApiDeviceListResponse

    @GET("/api/latest/device")
    suspend fun devicesByUser(): ApiDeviceListResponse

    @PATCH("/api/latest/house/{houseId}/device/{deviceId}/state")
    suspend fun setDeviceState(
        @Path("houseId") houseId: String,
        @Path("deviceId") deviceId: String,
        @Body request: ApiDeviceSetStateRequest,
    ): ApiDeviceSetStateResponse

    @GET("/api/latest/house/{houseId}/shopping-list")
    suspend fun shoppingListsByHouse(
        @Path("houseId") houseId: String,
    ): ApiShoppingListListResponse

    @GET("/api/latest/house/{houseId}/shopping-list/{listId}/item")
    suspend fun shoppingItemsByList(
        @Path("houseId") houseId: String,
        @Path("listId") listId: String,
    ): ApiShoppingItemListResponse
}
