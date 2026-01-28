package hata.data.models

import com.squareup.moshi.JsonClass


@JsonClass(generateAdapter = true)
data class Room(
    val id: String,
    val name: String,
    val devices: List<Device>,
)
