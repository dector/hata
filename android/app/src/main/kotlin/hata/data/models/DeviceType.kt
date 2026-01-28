package hata.data.models

import com.squareup.moshi.Json


enum class DeviceType {
    @Json(name = "light")
    Light,

    @Json(name = "unknown")
    Unknown,
}
