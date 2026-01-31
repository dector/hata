package hata.data.api

import hata.data.models.ServerInfo


interface ServerService {
    suspend fun connect(serverUrl: String): Result<ServerInfo>
    suspend fun login(serverUrl: String, username: String, password: String): Result<String>
}
