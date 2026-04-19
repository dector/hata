package hata.data.api

open class ApiException(
    message: String,
    val code: String,
) : Exception(message)

class DeviceNoAckException(
    message: String,
) : ApiException(
    message = message,
    code = CODE,
) {
    companion object {
        const val CODE = "device-no-ack"
    }
}
