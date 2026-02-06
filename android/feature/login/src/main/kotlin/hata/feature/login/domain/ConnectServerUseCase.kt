package hata.feature.login.domain

import hata.data.api.ServerProbeApi
import hata.data.models.ServerInfo
import javax.inject.Inject


class ConnectServerUseCase @Inject constructor(
    private val serverProbeApi: ServerProbeApi,
) {

    suspend fun run(serverUrl: String): Result<ServerInfo> {
        return serverProbeApi.connect(serverUrl)
    }
}
