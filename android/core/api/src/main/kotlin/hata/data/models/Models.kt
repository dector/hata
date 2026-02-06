package hata.data.models


enum class DeviceState {
    ON,
    OFF,
    UNKNOWN,
}

data class Device(
    val id: String,
    val name: String,
    val integrationId: String,
    val integrationData: Map<String, Any?>?,
    val state: DeviceState,
    val houseId: String? = null,
)

//region Server info

data class ServerInfo(
    val serverUrl: String,
    val serverName: String,
    val version: String? = null,
)

//endregion
