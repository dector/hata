package hata.data.models

import com.squareup.moshi.JsonClass


@JsonClass(generateAdapter = true)
data class Integration(
    val id: String,
    val data: IntegrationData,
)
