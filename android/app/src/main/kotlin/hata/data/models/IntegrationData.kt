package hata.data.models

import com.squareup.moshi.JsonClass


@JsonClass(generateAdapter = true)
data class IntegrationData(
    val type: String,
    val ip: String,
)
