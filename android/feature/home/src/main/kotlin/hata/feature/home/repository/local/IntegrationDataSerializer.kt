package hata.feature.home.repository.local

import com.squareup.moshi.JsonAdapter
import com.squareup.moshi.Moshi
import com.squareup.moshi.Types


internal object IntegrationDataSerializer {

    private val mapAdapter: JsonAdapter<Map<String, Any?>> = Moshi.Builder()
        .build()
        .adapter<Map<String, Any?>>(mapType())

    fun serialize(data: Map<String, Any?>?): String? {
        return data?.let { mapAdapter.toJson(it) }
    }

    fun deserialize(raw: String?): Map<String, Any?>? {
        return raw?.let { value ->
            runCatching { mapAdapter.fromJson(value) }
                .getOrNull()
        }
    }
}

private fun mapType() = Types.newParameterizedType(
    Map::class.java,
    String::class.java,
    Any::class.java,
)
