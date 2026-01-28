package hata.data.models

import com.squareup.moshi.JsonClass


@JsonClass(generateAdapter = true)
data class RemoteConfiguration(
    val homes: List<Home>,
)
