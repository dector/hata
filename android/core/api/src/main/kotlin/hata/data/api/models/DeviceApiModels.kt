package hata.data.api.models

import com.squareup.moshi.JsonClass


@JsonClass(generateAdapter = true)
data class ApiDeviceListResponse(
    val devices: List<ApiDevice> = emptyList(),
)

@JsonClass(generateAdapter = true)
data class ApiDevice(
    val id: String,
    val name: String,
    val integration: ApiDeviceIntegration? = null,
    val state: String,
    val house: ApiHouseRef? = null,
)

@JsonClass(generateAdapter = true)
data class ApiDeviceIntegration(
    val id: String,
    val data: Map<String, Any?>? = null,
)

@JsonClass(generateAdapter = true)
data class ApiHouseRef(
    val id: String,
)
