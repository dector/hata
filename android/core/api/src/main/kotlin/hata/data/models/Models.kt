package hata.data.models

import com.squareup.moshi.Json
import com.squareup.moshi.JsonClass


//region Remote configuration

@JsonClass(generateAdapter = true)
data class RemoteConfiguration(
    val homes: List<Home>,
)

@JsonClass(generateAdapter = true)
data class Home(
    val id: String,
    val name: String,
    val rooms: List<Room>,
)

@JsonClass(generateAdapter = true)
data class Room(
    val id: String,
    val name: String,
    val devices: List<Device>,
)

@JsonClass(generateAdapter = true)
data class Device(
    val id: String,
    val name: String,
    val type: DeviceType,
    val integration: Integration? = null,
)

enum class DeviceType {
    @Json(name = "light")
    Light,

    @Json(name = "unknown")
    Unknown,
}

@JsonClass(generateAdapter = true)
data class Integration(
    val id: String,
    val data: IntegrationData,
)

@JsonClass(generateAdapter = true)
data class IntegrationData(
    val type: String,
    val ip: String,
)

//endregion

//region Server info

data class ServerInfo(
    val serverUrl: String,
    val serverName: String,
    val version: String? = null,
)

//endregion
