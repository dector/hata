package hata.data.models

import com.squareup.moshi.JsonClass


@JsonClass(generateAdapter = true)
data class Device(
    val id: String,
    val name: String,
    val type: DeviceType
)
