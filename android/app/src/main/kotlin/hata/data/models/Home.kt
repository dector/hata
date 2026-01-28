package hata.data.models

import com.squareup.moshi.JsonClass


@JsonClass(generateAdapter = true)
data class Home(
    val id: String,
    val name: String,
    val rooms: List<Room>,
)
