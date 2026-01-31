package hata.data.models


data class ServerInfo(
    val serverUrl: String,
    val serverName: String,
    val version: String? = null,
)
