package hata.data.api

import hata.data.models.ServerInfo
import retrofit2.HttpException
import java.io.IOException


class RealServerProbeApi : ServerProbeApi {

    override suspend fun connect(serverUrl: String): Result<ServerInfo> {
        return try {
            val api = ApiClient
                .createRetrofit(serverUrl)
                .create(HataApi::class.java)
            val response = api.ping()

            Result.success(
                ServerInfo(
                    serverUrl = serverUrl,
                    serverName = response.serverName,
                    version = response.version,
                ),
            )
        } catch (e: HttpException) {
            Result.failure(
                Exception(ApiClient.parseErrorMessage(e) ?: "Failed to connect to server"),
            )
        } catch (e: IOException) {
            Result.failure(
                Exception("Network error: ${e.message ?: "Unable to reach server"}"),
            )
        } catch (e: Exception) {
            Result.failure(
                Exception("Connection failed: ${e.message ?: "Unknown error"}"),
            )
        }
    }
}
