package hata.data.api

import hata.data.models.ServerInfo


interface ServerProbeApi {

    suspend fun connect(serverUrl: String): Result<ServerInfo>
}
