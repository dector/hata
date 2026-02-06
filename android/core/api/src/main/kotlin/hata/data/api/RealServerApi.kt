package hata.data.api

import hata.data.api.models.ApiDevice
import hata.data.api.models.LoginRequest
import retrofit2.HttpException
import java.io.IOException


class RealServerApi(
    override val serverUrl: String,
    tokenProvider: AuthTokenProvider,
) : ServerApi {

    private val api = ApiClient
        .createRetrofit(serverUrl, tokenProvider)
        .create(HataApi::class.java)

    override suspend fun login(
        username: String,
        password: String,
    ): Result<String> {
        return try {
            val response = api.login(
                LoginRequest(
                    username = username,
                    password = password,
                ),
            )

            Result.success(response.session.token)
        } catch (e: HttpException) {
            Result.failure(
                Exception(ApiClient.parseErrorMessage(e) ?: "Login failed"),
            )
        } catch (e: IOException) {
            Result.failure(
                Exception("Network error: ${e.message ?: "Unable to reach server"}"),
            )
        } catch (e: Exception) {
            Result.failure(
                Exception("Login failed: ${e.message ?: "Unknown error"}"),
            )
        }
    }

    override suspend fun fetchHouse(): Result<Unit> {
        return try {
            api.latestHouse()

            Result.success(Unit)
        } catch (e: HttpException) {
            Result.failure(
                Exception(ApiClient.parseErrorMessage(e) ?: "Failed to load latest house"),
            )
        } catch (e: IOException) {
            Result.failure(
                Exception("Network error: ${e.message ?: "Unable to reach server"}"),
            )
        } catch (e: Exception) {
            Result.failure(
                Exception("Failed to load latest house: ${e.message ?: "Unknown error"}"),
            )
        }
    }

    override suspend fun fetchDevicesByHouse(
        houseId: String,
    ): Result<List<ApiDevice>> {
        return try {
            val response = api.devicesByHouse(
                houseId = houseId,
            )

            Result.success(response.devices)
        } catch (e: HttpException) {
            Result.failure(
                Exception(ApiClient.parseErrorMessage(e) ?: "Failed to load devices for house"),
            )
        } catch (e: IOException) {
            Result.failure(
                Exception("Network error: ${e.message ?: "Unable to reach server"}"),
            )
        } catch (e: Exception) {
            Result.failure(
                Exception("Failed to load devices: ${e.message ?: "Unknown error"}"),
            )
        }
    }

    override suspend fun fetchDevicesByUser(
    ): Result<List<ApiDevice>> {
        return try {
            val response = api.devicesByUser()

            Result.success(response.devices)
        } catch (e: HttpException) {
            Result.failure(
                Exception(ApiClient.parseErrorMessage(e) ?: "Failed to load devices"),
            )
        } catch (e: IOException) {
            Result.failure(
                Exception("Network error: ${e.message ?: "Unable to reach server"}"),
            )
        } catch (e: Exception) {
            Result.failure(
                Exception("Failed to load devices: ${e.message ?: "Unknown error"}"),
            )
        }
    }
}
