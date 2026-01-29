package hata.integrations.wiz

import kotlinx.serialization.SerialName
import kotlinx.serialization.Serializable


@Serializable
data class WizDevice(
    val name: String,
    val ip: String,
    val type: String,
    val mac: String
)

@Serializable
data class PilotState(
    val state: Boolean,
    val dimming: Int = 0,
    val r: Int = 0,
    val g: Int = 0,
    val b: Int = 0,
    val temp: Int = 0,
    @SerialName("sceneId") val sceneId: Int = 0
)

@Serializable
internal data class RegistrationResponse(
    val method: String,
    val env: String? = null,
    val result: RegistrationResult
)

@Serializable
internal data class RegistrationResult(
    val mac: String,
    val ip: String? = null,
    @SerialName("productName") val productName: String? = null,
    @SerialName("modelName") val modelName: String? = null,
    @SerialName("fwVersion") val firmwareVersion: String? = null,
    @SerialName("moduleName") val moduleName: String? = null
)

@Serializable
internal data class GetPilotResponse(
    val method: String,
    val env: String? = null,
    val result: PilotState
)
