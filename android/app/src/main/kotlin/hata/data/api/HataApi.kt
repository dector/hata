package hata.data.api

import hata.data.api.models.LoginRequest
import hata.data.api.models.LoginResponse
import hata.data.api.models.PingResponse
import retrofit2.http.Body
import retrofit2.http.GET
import retrofit2.http.POST

interface HataApi {
    @GET("/api/latest/ping")
    suspend fun ping(): PingResponse

    @POST("/api/latest/auth/login")
    suspend fun login(@Body request: LoginRequest): LoginResponse
}
