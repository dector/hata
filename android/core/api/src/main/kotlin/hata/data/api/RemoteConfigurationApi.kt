package hata.data.api

import hata.data.models.RemoteConfiguration
import retrofit2.http.GET


interface RemoteConfigurationApi {

    @GET(".")
    suspend fun loadConfiguration(): Result<RemoteConfiguration>
}
