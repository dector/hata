package hata.data.api

import hata.data.models.RemoteConfiguration
import retrofit2.http.GET


interface RemoteConfigurationService {
    @GET(".")
    suspend fun loadConfiguration(): Result<RemoteConfiguration>
}
